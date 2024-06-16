package timehelper

import (
    "fmt"
    "math"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"

    "store-api/pkg/data/constant"
)

const _24HourFormat = "15:04:05"

// For switch DB +7 and UTC
var noConvertToUTC = os.Getenv("DB_TZ") != "UTC"

type DateRange struct {
    StartDate time.Time
    EndDate   time.Time
}

//AddDurationFromTimeString Append current duration (cd)(hh:mm:ss) with total duration (td)(hh:mm:ss)
func AddDurationFromTimeString(cd string, td string) string {
    splitCD := strings.Split(cd, ":")
    splitTD := strings.Split(td, ":")

    currentDuration, _ := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", splitCD[0], splitCD[1], splitCD[2]))
    totalDuration, _ := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", splitTD[0], splitTD[1], splitTD[2]))

    totalDuration += currentDuration
    totalDurationStr := totalDuration.String()
    hour := "00"
    min := "00"
    sec := "00"

    if strings.Contains(totalDurationStr, "h") {
        splitHour := strings.Split(totalDurationStr, "h")
        hour = splitHour[0]
        totalDurationStr = splitHour[1]
    }
    if strings.Contains(totalDurationStr, "m") {
        splitMin := strings.Split(totalDurationStr, "m")
        min = splitMin[0]
        totalDurationStr = splitMin[1]
    }

    if strings.Contains(totalDurationStr, "s") {
        splitSec := strings.Split(totalDurationStr, "s")
        sec = splitSec[0]
    }

    durationHour, _ := strconv.Atoi(hour)
    durationMinute, _ := strconv.Atoi(min)
    durationSecond, _ := strconv.Atoi(sec)

    return fmt.Sprintf("%s:%s:%s",
        set24HourFormatOnTimeUnit(durationHour),
        set24HourFormatOnTimeUnit(durationMinute),
        set24HourFormatOnTimeUnit(durationSecond))
}

//set24HourFormatOnTimeUnit append "0" on time unit to matches the 24-hour format (e.g. 3:12:4 => 03:12:04)
func set24HourFormatOnTimeUnit(parameter int) string {
    if parameter < 10 {
        return "0" + strconv.Itoa(parameter)
    }
    return strconv.Itoa(parameter)
}

// Get duration, will return (hh:mm:ss)
func TimeDuration(activityDuration float64) string {
    duration := "00:00:00"

    dur := int(activityDuration)

    if dur > 0 {
        modMin := (dur % 3600)
        modHour := dur % 86400
        hour := int(math.Floor(float64(modHour) / 3600))
        min := int(math.Floor(float64(modMin) / 60))
        sec := dur % 60

        duration = fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
    }

    return duration
}

func BuildDateRangeOnCurrentMonth() (time.Time, time.Time) {
    currentTime := time.Now()
    currentYear, currentMonth, _ := currentTime.Date()
    currentLocation := currentTime.Location()
    startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    endOfMonth := startOfMonth.AddDate(0, 1, -1)
    return startOfMonth, endOfMonth
}

/*SplitOvertimeDuration
Iterate over range of days and map into start_date and end_date
Example :
startDate = '2021-10-01 15:00:00'
endDate  = '2021-10-03' 13:00:00'

result : [
	{"start": "2021-10-01 15:00:00", "end" : "2021-10-01 23:59:59"},
	{"start": "2021-10-02 00:00:00", "end" : "2021-10-02 23:59:59"},
	{"start": "2021-10-03 00:00:00", "end" : "2021-10-03 15:00:00"},
]
ref: splitOvertimeDuration function on talenta-core
*/
func SplitOvertimeDuration(startTime time.Time, endTime time.Time) []map[string]time.Time {
    accuDate := make([]map[string]time.Time, 0)
    limitDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 59, endTime.Location())
    for d := startTime; d.Before(limitDate) || d.Equal(limitDate); d = d.AddDate(0, 0, 1) {
        start := d
        end := time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 59, d.Location())

        if start.Equal(startTime) {
            start = time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
        } else {
            start = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
        }
        if d.Day() == endTime.Day() {
            end = endTime
        }
        accuDate = append(accuDate, map[string]time.Time{"start": start, "end": end})
    }

    return accuDate
}

//	GetStartDayAndEndDayOfTheWeekByDate returns start date as previous Monday and end date as next Monday.
//	Assume a week starts on Monday.
func GetStartDayAndEndDayOfTheWeekByDate(date time.Time) DateRange {
    starMonday := 1
    endMonday := 8

    if date.Weekday() == time.Sunday {
        return DateRange{
            StartDate: date.AddDate(0, 0, -6),
            EndDate:   date.AddDate(0, 0, 1),
        }
    }
    dateNumberOnTheWeek := int(date.Weekday())
    return DateRange{
        StartDate: date.AddDate(0, 0, starMonday-dateNumberOnTheWeek),
        EndDate:   date.AddDate(0, 0, endMonday-dateNumberOnTheWeek),
    }
}

func TodayStr() string {
    return time.Now().Format(constant.DefaultDateLayout)
}

func LocationOffset(t time.Time) int {
    _, offset := t.Zone()
    return offset / 3600
}

// NowUTC used to get now in UTC
//func NowInUTC() time.Time {
//	return time.Now().UTC()
//}

// NowInUTCAsDateTimeStr used to format application's time (+7) to UTC
func NowInUTCAsDateTimeStr() string {
    if noConvertToUTC {
        return time.Now().Format(constant.DefaultDatetimeLayout)
    }
    return time.Now().UTC().Format(constant.DefaultDatetimeLayout)
}

// IsUnknownTime alias for easy to remember: time = 0001-01-01 00:00:00 +0000 UTC
func IsUnknownTime(t time.Time) bool {
    return t.IsZero()
}

func GetDate(days int) time.Time {
    now := time.Now()
    today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    return today.AddDate(0, 0, days)
}

func TimeToDate(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// NowInUTCAsDateStr used to format application's time (+7) to UTC
func NowInUTCAsDateStr() string {
    if noConvertToUTC {
        return time.Now().Format(constant.DefaultDateLayout)
    }
    return time.Now().UTC().Format(constant.DefaultDateLayout)
}

// NowInUTCAsTimeStr used to format application's time (+7) to UTC
func NowInUTCAsTimeStr() string {
    if noConvertToUTC {
        return time.Now().Format(constant.DefaultTimeLayout)
    }
    return time.Now().UTC().Format(constant.DefaultTimeLayout)
}

// FormatDateToUTC used to format application's time (+7) to UTC
func FormatDateToUTC(utc time.Time) string {
    if noConvertToUTC {
        return utc.Format(constant.DefaultDateLayout)
    }
    return utc.UTC().Format(constant.DefaultDateLayout)
}

//// FormatTimeToUTC used to format application's time (+7) to UTC
//func FormatTimeToUTC(utc time.Time) string {
//	return utc.UTC().Format(constant.DefaultTimeLayout)
//}

// ToUTC used to format application's time (+7) to UTC
func ToUTC(t1 time.Time) time.Time {
    return t1.UTC()
}

// ToLocal get from Dockerfile, +7GMT
func ToLocal(t1 time.Time) time.Time {
    return t1.Local()
}

func FormatToDefaultDatetimeLayout(t time.Time) string {
    return t.Format(constant.DefaultDatetimeLayout)
}

func FormatUTCToDefaultDatetimeLayout(t time.Time) string {
    return t.UTC().Format(constant.DefaultDatetimeLayout)
}

func FormatDate(t time.Time) string {
    return t.Local().Format(constant.DefaultDateLayout)
}

func FormatDatePtr(t *time.Time) string {
    if t != nil {
        return t.Local().Format(constant.DefaultDateLayout)
    }
    return ""
}

func FormatTime(t time.Time) string {
    return t.Local().Format(constant.DefaultTimeLayout)
}

func FormatTimePtr(t *time.Time) string {
    if t != nil {
        return t.Local().Format(constant.DefaultTimeLayout)
    }
    return ""
}

func FormatDynamicTimePtr(t *time.Time, layout string) string {
    if t != nil {
        if layout == "" {
            layout = constant.DefaultTimeLayout
        }
        return t.Local().Format(layout)
    }
    return ""
}

func FormatDateTime(t time.Time) string {
    return t.Local().Format(constant.DefaultDatetimeLayout)
}

func FormatDateTimeLayout(t time.Time, layout string) string {
    return t.Local().Format(layout)
}

func FormatDateTimePtrLayout(t *time.Time, layout string) string {
    if t != nil {
        return (*t).Local().Format(layout)
    }
    return ""
}

func FormatDateTimePtr(t *time.Time) string {
    if t != nil {
        return (*t).Local().Format(constant.DefaultDatetimeLayout)
    }
    return ""
}

func FormatDateTimePtrLayoutDefault(t *time.Time, layout string, defaultValue string) string {
    if t != nil {
        return (*t).Local().Format(layout)
    }
    return defaultValue
}

func Equal(t1 time.Time, t2 time.Time) bool {
    return t1.Equal(t2)
}

func GreaterThanOrEqual(t1 time.Time, t2 time.Time) bool {
    return t1.After(t2) || t1.Equal(t2)
}

func LessThanOrEqual(t1 time.Time, t2 time.Time) bool {
    return t1.Before(t2) || t1.Equal(t2)
}

func AddDays(t time.Time, day int64) time.Time {
    return t.Add(time.Hour * time.Duration(day*24))
}

func AddHours(t time.Time, hr int64) time.Time {
    return t.Add(time.Hour * time.Duration(hr))
}

func AddMinutes(t time.Time, min int64) time.Time {
    return t.Add(time.Minute * time.Duration(min))
}

func AddSeconds(t time.Time, sec int64) time.Time {
    return t.Add(time.Second * time.Duration(sec))
}

func ConvertToValidDate(dateStr string) string {
    // input must be yyyy-mm-dd with possible missing zero pad, i.e. 2006-1-2

    dates := strings.Split(dateStr, "-")
    // add zero pad for month
    if len(dates[1]) < 2 {
        dates[1] = "0" + dates[1]
    }
    // add zero pad for day
    if len(dates[2]) < 2 {
        dates[2] = "0" + dates[2]
    }

    return strings.Join(dates, "-")
}

func ConvertToValidTime(timeStr string) string {
    times := strings.Split(timeStr, ":")

    if len(times[0]) < 2 {
        times[0] = "0" + times[0]
    }

    if len(times[1]) < 2 {
        times[1] = "0" + times[1]
    }

    return strings.Join(times, ":")
}

/**
MySQL TIMEDIFF format is hh:mm:ss, so activity duration must be parsed
to second units, so that can be added or subtracted
*/
func ParseCustomDuration(st string) (int, error) {
    var h, m, s int
    n, err := fmt.Sscanf(st, "%d:%d:%d", &h, &m, &s)
    if err != nil || n != 3 {
        return 0, err
    }
    return h*3600 + m*60 + s, nil
}

func AddDurations(durations []string) int64 {
    totalSeconds := int64(0)
    for _, duration := range durations {
        newSeconds, _ := ParseCustomDuration(duration)
        totalSeconds += int64(newSeconds)
    }

    return totalSeconds
}

type WeekByDate struct {
    StartDate time.Time
    EndDate   time.Time
    Week      []time.Time
}

func GetWeekByDate(strDate string) (*WeekByDate, error) {
    date, err := Strtotime(strDate)

    if err != nil {
        return nil, err
    }

    ret := WeekByDate{}
    tn := time.Unix(date, 0)
    year, weekNumber := tn.ISOWeek()
    startDateObj, _ := DateISODateSet(year, weekNumber, 7)
    endDateObj, _ := DateModify(startDateObj, "+7 days")
    ret.StartDate = AddDays(startDateObj, 1)
    ret.EndDate = AddDays(endDateObj, 1)

    for startDateObj.Before(endDateObj) {
        date := AddDays(startDateObj, 1)
        startDateObj = date
        ret.Week = append(ret.Week, date)
    }

    return &ret, err
}

/**
MySQL treat invalid time to default time,
either by replacing the minute or the hour
*/
func ConvertToDefaultMysqlTime(time string) string {
    /**
    golang playground: https://go.dev/play/p/zED8mvwz7e5
    Convert invalid time to valid mysql time
    @param string time
    @return string

    Example:
    	ConvertToValidMysqlTime("09:MM") -> "09:00"
    	ConvertToValidMysqlTime("HH:20") -> "00:00"
    	ConvertToValidMysqlTime("23:59") -> "23:59"
    	ConvertToValidMysqlTime("24:00") -> "00:00"
    	ConvertToValidMysqlTime("23:60") -> "23:00"
    	ConvertToValidMysqlTime("24") -> "00:00"

    */
    hourRe := regexp.MustCompile("[0-1][0-9]|2[0-3]")
    minuteRe := regexp.MustCompile(`[0-5][0-9]`)

    splitted := strings.Split(time, ":")
    if len(splitted) >= 2 {
        hour := splitted[0]
        minute := splitted[1]
        if !hourRe.MatchString(hour) {
            return "00:00:00"
        }
        if !minuteRe.MatchString(minute) {
            return hour + ":00:00"
        }
    } else {
        return "00:00:00"
    }

    return time + ":00"
}

// GetStartAndEndOfMonth Return the start of the month (YYYY-MM-01 00:00:00)
// and end of the month (YYYY-MM-31 23:59:59)
func GetStartAndEndOfMonth(month time.Month, year int, location *time.Location) (time.Time, time.Time) {
    startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
    endOfMonth := AddHours(startOfMonth.AddDate(0, 1, -1), 24).Add(time.Nanosecond * -1)
    return startOfMonth, endOfMonth
}

func FormatToLocalTime(t time.Time) time.Time {
    nowStr := t.Format(constant.DefaultDateLayout) + " 00:00:00"
    result, _ := time.ParseInLocation(constant.DefaultDatetimeLayout, nowStr, time.Local)

    return result
}
