package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(aDuration string, defaultDuration time.Duration) time.Duration {

	if aDuration != "" {
		d, err := time.ParseDuration(aDuration)
		if err != nil {
			d = defaultDuration
		}
		return d
	}

	return defaultDuration
}

func ParseISO8601Duration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

var timeLayout = "15:04"
var TimeParseError = errors.New(`TimeParseError: should be a string formatted as "15:04"`)

type TimeOfDay struct {
	hour   int
	minute int
}

func NewTimeOfDay(t time.Time) TimeOfDay {
	return TimeOfDay{hour: t.Hour(), minute: t.Minute()}
}

func NewTimeOfDayFromString(s string) (TimeOfDay, error) {

	if len(s) == 0 {
		return TimeOfDay{}, nil
	}

	switch len(s) {
	case 5:
		break
	case 7:
		// len(`"23:59"`) == 7
		s = s[1:6]
		break
	default:
		return TimeOfDay{}, TimeParseError
	}

	ret, err := time.Parse(timeLayout, s)
	if err != nil {
		return TimeOfDay{}, err
	}

	return NewTimeOfDay(ret), nil
}

func (td *TimeOfDay) MarshalJSON() ([]byte, error) {
	t := time.Date(1, 1, 1, td.hour, td.minute, 0, 0, time.Local)
	return []byte(`"` + t.Format(timeLayout) + `"`), nil
}

func (td *TimeOfDay) IsZero() bool {
	return td.hour == 0 && td.minute == 0
}

func (td *TimeOfDay) After(t TimeOfDay) bool {

	if td.hour != t.hour {
		return td.hour > t.hour
	}

	if td.minute != t.minute {
		return td.minute > t.minute
	}

	return true
}

func (td *TimeOfDay) Before(t TimeOfDay) bool {

	if td.hour != t.hour {
		return td.hour < t.hour
	}

	if td.minute != t.minute {
		return td.minute < t.minute
	}

	return true
}

func (td *TimeOfDay) UnmarshalJSON(b []byte) error {

	ret, err := NewTimeOfDayFromString(string(b))
	if err != nil {
		return err
	}

	td.hour = ret.hour
	td.minute = ret.minute
	return nil
}

func TimeOfDayBetween(theTime time.Time, t1 TimeOfDay, t2 TimeOfDay) bool {

	td := NewTimeOfDay(theTime)

	comparisonMode := true
	if !t1.IsZero() && !t2.IsZero() {
		if t1.After(t2) {
			t1, t2 = t2, t1
			comparisonMode = false
		}
	}

	rc := true
	if comparisonMode {
		if rc && !t1.IsZero() {
			rc = rc && td.After(t1)
		}

		if rc && !t2.IsZero() {
			rc = rc && td.Before(t2)
		}
	} else {
		rc = td.Before(t1) || td.After(t2)
	}

	return rc
}

var TimeRangeParseError = errors.New(`TimeParseError: should be a string formatted as "18:05-19:00"`)

type TimeOfDayRange struct {
	compMode string
	start    TimeOfDay
	end      TimeOfDay
}

type TimeOfDayRanges []TimeOfDayRange

var TimeOfDayRangePatternRegexp = regexp.MustCompile("^(([0-9]{2}):([0-9]{2}))?(-)(([0-9]{2}):([0-9]{2}))?$")
var FullDayRange = TimeOfDayRange{end: TimeOfDay{23, 59}, compMode: "in"}
var FullDayRanges = []TimeOfDayRange{FullDayRange}

func NewTimeOfDayRangesFromString(s string) (TimeOfDayRanges, error) {

	trs := make([]TimeOfDayRange, 0)

	if len(s) == 0 {
		return FullDayRanges, nil
	}

	stringRanges := strings.Split(s, ",")
	for _, r := range stringRanges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		tr, err := NewTimeOfDayRangeFromString(r)
		if err == nil {
			trs = append(trs, tr)
		} else {
			return FullDayRanges, err
		}
	}

	if len(trs) == 0 {
		return FullDayRanges, nil
	}

	return trs, nil
}

func NewTimeOfDayRangeFromString(s string) (TimeOfDayRange, error) {

	tr := TimeOfDayRange{}

	if len(s) == 0 {
		return FullDayRange, nil
	}

	matches := TimeOfDayRangePatternRegexp.FindStringSubmatch(s)
	if len(matches) != 8 {
		return FullDayRange, TimeRangeParseError
	}

	if matches[1] != "" {
		h, _ := strconv.Atoi(matches[2])
		m, _ := strconv.Atoi(matches[3])

		if h > 23 || m > 59 {
			return FullDayRange, TimeRangeParseError
		}

		tr.start = TimeOfDay{hour: h, minute: m}
	}

	if matches[5] != "" {
		h, _ := strconv.Atoi(matches[6])
		m, _ := strconv.Atoi(matches[7])

		if h > 23 || m > 59 {
			return FullDayRange, TimeRangeParseError
		}

		tr.end = TimeOfDay{hour: h, minute: m}
	} else {
		tr.end = TimeOfDay{hour: 23, minute: 59}
	}

	if tr.end.Before(tr.start) {
		tr.compMode = "out"
		tr.start, tr.end = tr.end, tr.start
	} else {
		tr.compMode = "in"
	}

	return tr, nil
}

func (trs TimeOfDayRanges) InRange(theTime time.Time) bool {
	for _, tr := range trs {
		if tr.InRange(theTime) {
			return true
		}
	}

	return false
}

func (tr *TimeOfDayRange) InRange(theTime time.Time) bool {

	td := NewTimeOfDay(theTime)

	rc := false
	switch tr.compMode {
	case "in":
		rc = td.After(tr.start) && td.Before(tr.end)
	case "out":
		rc = td.Before(tr.start) || td.After(tr.end)
	default:
		rc = td.After(tr.start) && td.Before(tr.end)
	}

	return rc
}

func DayCompare(t1 time.Time, t2 time.Time) int {
	if t1.Year() == t2.Year() {
		if t1.Month() == t2.Month() {
			return t1.Day() - t2.Day()
		}

		return int(t1.Month() - t2.Month())
	}

	return t1.Year() - t2.Year()
}
