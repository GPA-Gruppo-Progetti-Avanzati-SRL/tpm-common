package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var idPattern = regexp.MustCompile(`^[A-Za-z_\-]*([0-9]*)$`)

func idAsInt64(id string) (int64, error) {

	matches := idPattern.FindStringSubmatch(id)
	if len(matches) != 2 {
		return 0, fmt.Errorf("cannot parse %s as int64", id)
	}

	n, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func ComputeMida(platform string, id string) (string, error) {

	idAsInt, err := idAsInt64(id)
	if err != nil {
		return "", err
	}

	mida := fmt.Sprintf("%s%s%07x", platform, time.Now().Format("06"), idAsInt%268435456)
	return mida, nil
}
