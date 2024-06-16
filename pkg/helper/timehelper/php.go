package timehelper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func Strtotime(str string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func DateISODateSet(year, week, day int) (time.Time, error) {
	firstDateOfYear, err := time.Parse("2006-1-2", fmt.Sprintf("%04d-%d-%d", year, 1, 1))
	if err != nil {
		return time.Time{}, err
	}

	offset := time.Duration(1-firstDateOfYear.Weekday()) * 24 * time.Hour
	firstDateOfFirstWeek := firstDateOfYear.Add(offset)

	return firstDateOfFirstWeek.Add(time.Duration(((week-1)*7+day-1)*24) * time.Hour), nil
}

func DateModify(t time.Time, modify string) (time.Time, error) {
	duration, err := DateIntervalCreateFromDateString(modify)
	if err != nil {
		return t, err
	}
	return t.Add(duration), nil
}

// DateIntervalCreateFromDateString returns a time.Duration from the given string
func DateIntervalCreateFromDateString(str string) (time.Duration, error) {
	reg := regexp.MustCompile("((\\+|\\-)?\\s*(\\d*)\\s*(day|month|year|week|hour|minute|second)s?\\s*)+?")
	matches := reg.FindAllStringSubmatch(str, -1)
	if matches != nil {
		var duration int64
		for _, match := range matches {
			var diff, num int64
			if match[3] == "" {
				num = 1
			} else {
				num, _ = strconv.ParseInt(match[3], 10, 64)
			}
			switch match[4] {
			case "day":
				diff = num * 86400
			case "month":
				diff = num * 86400 * 30
			case "year":
				diff = num * 86400 * 365
			case "week":
				diff = num * 86400 * 7
			case "hour":
				diff = num * 3600
			case "minute":
				diff = num * 60
			case "second":
				diff = num
			}
			if match[2] == "-" {
				diff = -diff
			}
			duration += diff
		}
		return time.Duration(duration) * time.Second, nil
	}
	return 0, errors.New("unsupported string format")
}
