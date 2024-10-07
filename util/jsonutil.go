package util

import "encoding/json"

func JSONEscapeNoQuotes(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func JSONEscape(i string, withSurroundingQuotes bool) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	if !withSurroundingQuotes {
		return s[1 : len(s)-1]
	}
	return s
}
