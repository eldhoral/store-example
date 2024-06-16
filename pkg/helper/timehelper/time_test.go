package timehelper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddDurationFromTimeString(t *testing.T) {
	t.Run("adding 0 current duration 0 total duration returns 00:00:00", func(t *testing.T) {
		current_dur := "00:00:00"
		total_dur := "00:00:00"
		result := AddDurationFromTimeString(current_dur, total_dur)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})

	t.Run("adding current duration to total duration returns the correct amount", func(t *testing.T) {
		current_dur := "10:12:52"
		total_dur := "15:01:14"
		result := AddDurationFromTimeString(current_dur, total_dur)

		assert.NotNil(t, result)
		assert.Equal(t, result, "25:14:06")
	})
}

func TestSet24HourFormatOnTimeUnit(t *testing.T) {
	t.Run("single digit int returns 0 + digit", func(t *testing.T) {
		var number int = 1
		result := set24HourFormatOnTimeUnit(number)

		assert.NotNil(t, result)
		assert.Equal(t, result, "01")
	})

	t.Run("double digit int returns same number", func(t *testing.T) {
		var number int = 12
		result := set24HourFormatOnTimeUnit(number)

		assert.NotNil(t, result)
		assert.Equal(t, result, "12")
	})
}

func TestTimeDuration(t *testing.T) {
	t.Run("duration of 0s returns 0s", func(t *testing.T) {
		var activityDuration float64 = 0
		result := TimeDuration(activityDuration)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})

	t.Run("duration of 100s returns 1 min 40s", func(t *testing.T) {
		var activityDuration float64 = 100
		result := TimeDuration(activityDuration)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:01:40")
	})

	t.Run("duration of 4200s returns 1 hour 10 mins", func(t *testing.T) {
		var activityDuration float64 = 4200
		result := TimeDuration(activityDuration)

		assert.NotNil(t, result)
		assert.Equal(t, result, "01:10:00")
	})

	t.Run("duration of -100s returns 0s", func(t *testing.T) {
		var activityDuration float64 = -100
		result := TimeDuration(activityDuration)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})
}

func TestSplitOvertimeDuration(t *testing.T) {
	t.Run("startDate after endDate returns correct number of time objects", func(t *testing.T) {
		startDate := time.Date(2021, time.October, 3, 15, 0, 0, 0, time.UTC)
		endDate := time.Date(2021, time.October, 3, 12, 0, 0, 0, time.UTC)
		result := SplitOvertimeDuration(startDate, endDate)

		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
	})

	t.Run("startDate before endDate returns correct number of time objects", func(t *testing.T) {
		startDate := time.Date(2021, time.November, 2, 22, 59, 59, 59, time.UTC)
		endDate := time.Date(2021, time.November, 3, 0, 0, 0, 0, time.UTC)
		result := SplitOvertimeDuration(startDate, endDate)

		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result))
	})
}

func TestGetStartDayAndEndDayOfTheWeekByDate(t *testing.T) {
	t.Run("A sunday date returns correct start and end dates", func(t *testing.T) {
		//Date: 20 Mar 2022 (Sunday)
		date := time.Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC)
		result := GetStartDayAndEndDayOfTheWeekByDate(date)

		startDate := result.StartDate
		endDate := result.EndDate

		//Start date: 14 Mar 2022 (Monday)
		assert.NotNil(t, startDate)
		assert.Equal(t, time.Weekday(1), startDate.Weekday())
		assert.Equal(t, int(2022), startDate.Year())
		assert.Equal(t, time.Month(3), startDate.Month())
		assert.Equal(t, int(14), startDate.Day())

		//End date: 21 Mar 2022 (Monday)
		assert.NotNil(t, result.EndDate)
		assert.Equal(t, time.Weekday(1), endDate.Weekday())
		assert.Equal(t, int(2022), endDate.Year())
		assert.Equal(t, time.Month(3), endDate.Month())
		assert.Equal(t, int(21), endDate.Day())
	})

	t.Run("A weekday date returns correct start and end dates", func(t *testing.T) {
		//Date: 1 Dec 2020 (Tuesday)
		date := time.Date(2020, time.December, 1, 0, 0, 0, 0, time.UTC)
		result := GetStartDayAndEndDayOfTheWeekByDate(date)

		startDate := result.StartDate
		endDate := result.EndDate

		//Start date: 30 Nov 2022 (Monday)
		assert.NotNil(t, startDate)
		assert.Equal(t, time.Weekday(1), startDate.Weekday())
		assert.Equal(t, int(2020), startDate.Year())
		assert.Equal(t, time.Month(11), startDate.Month())
		assert.Equal(t, int(30), startDate.Day())

		//End date: 7 Dec 2020 (Monday)
		assert.NotNil(t, result.EndDate)
		assert.Equal(t, time.Weekday(1), endDate.Weekday())
		assert.Equal(t, int(2020), endDate.Year())
		assert.Equal(t, time.Month(12), endDate.Month())
		assert.Equal(t, int(7), endDate.Day())
	})
}

func TestConvertToValidDate(t *testing.T) {
	t.Run("Date with unpadded month and day returns fully padded date", func(t *testing.T) {
		date := "2006-1-1"
		result := ConvertToValidDate(date)

		assert.NotNil(t, result)
		assert.Equal(t, "2006-01-01", result)
	})

	t.Run("Date with unpadded month returns fully padded date", func(t *testing.T) {
		date := "2006-1-11"
		result := ConvertToValidDate(date)

		assert.NotNil(t, result)
		assert.Equal(t, "2006-01-11", result)
	})

	t.Run("Date with unpadded day returns fully padded date", func(t *testing.T) {
		date := "2006-11-4"
		result := ConvertToValidDate(date)

		assert.NotNil(t, result)
		assert.Equal(t, "2006-11-04", result)
	})
}

//does not format seconds.
func TestConvertToValidTime(t *testing.T) {
	t.Run("Time with unpadded hour and minute returns fully padded time", func(t *testing.T) {
		time := "1:2"
		result := ConvertToValidTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, "01:02", result)
	})

	t.Run("Time with unpadded hour returns fully padded time", func(t *testing.T) {
		time := "1:52"
		result := ConvertToValidTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, "01:52", result)
	})

	t.Run("Time with unpadded minute returns fully padded time", func(t *testing.T) {
		time := "11:9"
		result := ConvertToValidTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, "11:09", result)
	})

}

func TestParseCustomDuration(t *testing.T) {
	t.Run("00:00:00 returns 0 total seconds", func(t *testing.T) {
		time := "00:00:00"
		result, err := ParseCustomDuration(time)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, int(0))
	})

	t.Run("Valid time returns correct total seconds", func(t *testing.T) {
		time := "10:50:15"
		result, err := ParseCustomDuration(time)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, int(39015))
	})

	t.Run("Invalid time (missing seconds) returns 0", func(t *testing.T) {
		time := "01:50"
		result, err := ParseCustomDuration(time)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, int(0))
	})
}

func TestAddDurations(t *testing.T) {
	t.Run("Empty durations returns 0", func(t *testing.T) {
		durations := []string{}
		result := AddDurations(durations)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})

	t.Run("Multiple durations returns correct total seconds", func(t *testing.T) {
		durations := []string{"11:10:00", "21:20:00", "00:00:01", "04:50:15"}
		result := AddDurations(durations)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(134416))
	})
}

func TestGetWeekByDate(t *testing.T) {
	t.Run("Correctly formatted (yyyy-mm-dd) date returns correct WeekByDate object", func(t *testing.T) {
		date := "2021-10-07"
		result, err := GetWeekByDate(date)

		assert.Nil(t, err)
		assert.NotNil(t, result)

		weekByDate := *result

		//Test start date
		assert.Equal(t, 4, weekByDate.StartDate.Day())
		assert.Equal(t, time.October, weekByDate.StartDate.Month())
		assert.Equal(t, 2021, weekByDate.StartDate.Year())

		//End date should be 7 days after start date
		assert.Equal(t, 11, weekByDate.EndDate.Day())
		assert.Equal(t, time.October, weekByDate.EndDate.Month())
		assert.Equal(t, 2021, weekByDate.EndDate.Year())

		//Should have 7 dates between start and end date
		assert.Equal(t, 7, len(weekByDate.Week))

		//Check that each date in week array is consecutive, starting from start date
		n := 4
		for _, week := range weekByDate.Week {
			if assert.Equal(t, n, week.Day()) {
				n += 1
			}
		}
	})

	t.Run("Incorrectly formatted date returns nil", func(t *testing.T) {
		date := "10-01-2020"
		result, err := GetWeekByDate(date)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestConvertToDefaultMysqlTime(t *testing.T) {
	t.Run("Both hours and minutes absent returns 00:00:00", func(t *testing.T) {
		time := "HH:MM"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})

	t.Run("Hours present and minutes absent returns hours:00:00", func(t *testing.T) {
		time := "09:MM"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "09:00:00")
	})

	t.Run("Hours absent and minutes present returns 00:00:00", func(t *testing.T) {
		time := "HH:16"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})

	t.Run("Both hours and minutes present returns hours:minutes:00", func(t *testing.T) {
		time := "23:59"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "23:59:00")
	})

	t.Run("24 hours returns 00:00:00", func(t *testing.T) {
		time := "24:00"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})

	t.Run("Invalid time returns 00:00:00", func(t *testing.T) {
		time := "24"
		result := ConvertToDefaultMysqlTime(time)

		assert.NotNil(t, result)
		assert.Equal(t, result, "00:00:00")
	})
}
