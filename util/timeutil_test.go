package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"regexp"
	"testing"
	"time"
)

func TestTimeOfDayRange(t *testing.T) {

	timeRanges := []string{
		"00:00-23:59",
		"-18:12",
		"17:12-",
		"17:12-17:13",
		"17:1-17:13",
		"11:00-12:00",
		"11:00-10:00",
		"11:20-10:00",
		"11:20-10:00,13:00-14:00",
	}

	for _, s := range timeRanges {
		tr, err := util.NewTimeOfDayRangesFromString(s)
		if err != nil {
			t.Error(s, err)
		}
		t.Log(s, tr)
	}

}

const rx = "^(([0-9]{2}):([0-9]{2}))?(,)(([0-9]{2}):([0-9]{2}))?$"

func TestTimeOfDayRangeRegexp(t *testing.T) {

	r, err := regexp.Compile(rx)
	if err != nil {
		t.Fatal(err)
	}

	var s string

	s = "00:00,23:59"
	matchRegexp(t, r, s)

	s = ",23:59"
	matchRegexp(t, r, s)

	s = "00:00,"
	matchRegexp(t, r, s)

	s = "00:00,23:5"
	matchRegexp(t, r, s)
}

func matchRegexp(t *testing.T, r *regexp.Regexp, s string) {
	t.Log("matching regexp on: ", s)
	matches := r.FindStringSubmatch(s)
	if len(matches) == 0 {
		t.Log("not matching regexp on: ", s)
	}

	for i, m := range matches {
		t.Logf("[%d] - %s", i, m)
	}
}

func TestParseTime(t *testing.T) {

	s1 := "2022-01-19"
	tm, err := time.Parse("2006-01-02", s1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tm)

	s2 := "2021-10-25T14:52:25.155Z"
	tm, err = time.Parse("2006-01-02T15:04:05.999Z", s2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tm)
}
