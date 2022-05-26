package jsonmask

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
)

type FieldInfo struct {
	Path     string
	MaskType string
	indexes  []int
}

type fieldRegistry map[string]FieldInfo
type JsonMask struct {
	registry map[string]fieldRegistry
}

func NewJsonMask() *JsonMask {
	return &JsonMask{registry: make(map[string]fieldRegistry)}
}

func (jm *JsonMask) Add(ctxName string, fields []FieldInfo) {
	fr := make(fieldRegistry)
	for _, f := range fields {
		up, ndxs := ParsePath(f.Path)
		f.indexes = ndxs
		fr[up] = f
	}

	if len(fr) > 0 {
		jm.registry[ctxName] = fr
	}
}

func (jm *JsonMask) Mask(ctxName string, jsonData []byte) ([]byte, error) {
	if len(jm.registry) == 0 || len(jsonData) == 0 {
		return jsonData, nil
	}

	_, ok := jm.registry[ctxName]
	if !ok {
		return jsonData, nil
	}

	var target interface{}
	err := json.Unmarshal(jsonData, &target)
	if err != nil {
		return nil, err
	}

	err = jm.walkThrough("request", "", nil, target)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (jm *JsonMask) walkThrough(ctxName, path string, parentTarget, target interface{}) error {

	log.Trace().Str("path", path).Msg("walking")

	fi, maskIt := jm.HasToBeMasked(ctxName, path)
	if maskIt {
		log.Info().Str("path", path).Msg("should be masked")
	}

	var err error
	switch v := target.(type) {
	case []interface{}:

		for i, p := range v {
			err = jm.walkThrough(ctxName, strings.Join([]string{path, fmt.Sprintf("[%d]", i)}, "."), target, p)
			if err != nil {
				return err
			}
		}

	case map[string]interface{}:
		for n, p := range v {
			err = jm.walkThrough(ctxName, strings.Join([]string{path, n}, "."), target, p)
			if err != nil {
				return err
			}
		}

	case string:
		if maskIt {
			err = doMask(parentTarget, path, target, fi)
		}
	case float64:
		if maskIt {
			err = doMask(parentTarget, path, target, fi)
		}

	default:
		log.Trace().Str("path", path).Msgf("%T", v)
	}

	return err
}

func (jm *JsonMask) HasToBeMasked(ctxName, path string) (FieldInfo, bool) {
	fr := jm.registry[ctxName]

	up, ndxs := ParsePath(path)
	if fi, ok := fr[up]; ok {
		if len(fi.indexes) > 0 {
			for i, val := range fi.indexes {
				if val != -1 && ndxs[i] != val {
					return FieldInfo{}, false
				}
			}
		}
		return fi, true
	}

	return FieldInfo{}, false
}

var PathIndexedRegexp = regexp.MustCompile("\\.\\[([0-9]*)\\]")

func ParsePath(p string) (string, []int) {

	if !strings.Contains(p, "[") {
		return p, nil
	}

	matches := PathIndexedRegexp.FindAllSubmatch([]byte(p), -1)
	ndxs := make([]int, len(matches))
	for i, m := range matches {
		p = strings.ReplaceAll(p, string(m[0]), "[]")
		m1 := string(m[1])
		if m1 != "" {
			ndxs[i], _ = strconv.Atoi(string(m[1]))
		} else {
			ndxs[i] = -1
		}
	}

	return p, ndxs
}
