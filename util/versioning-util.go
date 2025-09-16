package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type VersionNumber struct {
	Major            int
	Minor            int
	Patch            int
	PreReleasePrefix string
	PreReleaseValue  int
}

func (v VersionNumber) String() string {
	return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Patch, v.PreReleasePrefix, v.PreReleaseValue)
}

func (v VersionNumber) LessThan(other VersionNumber) bool {
	if v.Major < other.Major {
		return true
	}

	if v.Minor < other.Minor {
		return true
	}

	if v.Patch < other.Patch {
		return true
	}

	// is not a pre-release
	if v.PreReleasePrefix == "" && v.PreReleaseValue == 0 {
		return false
	}

	// it is a pre-release but the other not.
	if other.PreReleasePrefix == "" && other.PreReleaseValue == 0 {
		return true
	}

	// both are pre-release
	if v.PreReleasePrefix < other.PreReleasePrefix {
		return true
	}

	if v.PreReleaseValue < other.PreReleaseValue {
		return true
	}

	return false
}

var versionNumberRegexps = regexp.MustCompile(`v?(\d+)\.(\d+)\.(\d+)(?:(-[a-zA-Z.]*)(\d*))?`)

func NewVersionNumberFromString(s string) (VersionNumber, error) {
	const semLogContext = "version-number::new"
	m := versionNumberRegexps.FindAllStringSubmatch(s, -1)

	if len(m) != 1 {
		return VersionNumber{}, errors.New(s + " is not a supported semantic version")
	}

	match := m[0]
	if len(match) < 4 {
		return VersionNumber{}, errors.New(s + " is not a supported semantic version")
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return VersionNumber{}, err
	}

	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return VersionNumber{}, err
	}

	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return VersionNumber{}, err
	}

	var preReleasePrefix string
	if len(match) > 4 {
		preReleasePrefix = match[4]
	}

	var preReleaseNumber int
	if len(match) > 5 && match[5] != "" {
		preReleaseNumber, err = strconv.Atoi(match[5])
		if err != nil {
			return VersionNumber{}, err
		}
	}

	v := VersionNumber{Major: major, Minor: minor, Patch: patch, PreReleasePrefix: preReleasePrefix, PreReleaseValue: preReleaseNumber}
	return v, nil
}
