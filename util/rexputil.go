package util

import (
	"github.com/lucasjones/reggen"
	"regexp"
)

func ExtractCapturedGroupIfMatch(re *regexp.Regexp, s string) string {

	matches := re.FindAllSubmatch([]byte(s), -1)

	for _, m := range matches {
		return string(m[1])
	}

	return s
}

func GenerateString(pattern string, n int) (string, error) {

	if n <= 0 {
		n = 10
	}

	str, err := reggen.Generate(pattern, n)
	return str, err
}
