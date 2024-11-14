package util

import (
	"encoding/json"
	"strings"
)

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

const (
	JSONStructureMap    = "map"
	JSONStructureArray  = "array"
	JSONStructureString = "string"
	JSONStructureEmpty  = "empty"
)

func JSONStructure(json string) string {
	json = strings.TrimSpace(json)
	if len(json) == 0 {
		return JSONStructureEmpty
	}

	var res string
	switch json[0] {
	case '{':
		res = JSONStructureMap
	case '[':
		res = JSONStructureArray
	default:
		res = JSONStructureString
	}

	return res
}
