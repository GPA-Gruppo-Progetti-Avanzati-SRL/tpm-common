package util

import "regexp"

func ExtractCapturedGroupIfMatch(re *regexp.Regexp, s string) string {

	matches := re.FindAllSubmatch([]byte(s), -1)

	for _, m := range matches {
		return string(m[1])
	}

	return s
}
