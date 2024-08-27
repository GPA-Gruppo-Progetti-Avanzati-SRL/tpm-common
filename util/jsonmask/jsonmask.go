package jsonmask

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
)

type fieldRegistry map[string]FieldInfo
type JsonMask struct {
	fullPathsRegistry    map[string]fieldRegistry
	partialPathsRegistry map[string][]FieldInfo
}

func NewJsonMask(opts ...Option) (*JsonMask, error) {
	jm := &JsonMask{fullPathsRegistry: make(map[string]fieldRegistry), partialPathsRegistry: make(map[string][]FieldInfo)}

	bld := builder{}
	for _, o := range opts {
		o(&bld)
	}

	var cfgMap Config
	var err error
	if bld.fn != "" {
		cfgMap, err = readConfigFromFile(bld.fn)
	} else {
		cfgMap, err = readConfig(bld.yamlData)
	}

	if err != nil {
		return nil, err
	}

	for dn, cfg := range cfgMap {
		jm.Add(dn, cfg.Fields)
	}

	return jm, nil
}

func (jm *JsonMask) Add(ctxName string, fields []FieldInfo) {
	fr := make(fieldRegistry)
	var pfr []FieldInfo
	for _, f := range fields {
		f.uxPath, f.indexes = ParsePath(f.Path)
		switch f.Path[0] {
		case '.':
			fr[f.uxPath] = f
		case '*':
			f.Path = f.Path[1:]
			f.uxPath = f.uxPath[1:]
			pfr = append(pfr, f)
		}

	}

	if len(fr) > 0 {
		jm.fullPathsRegistry[ctxName] = fr
	}

	if len(pfr) > 0 {
		jm.partialPathsRegistry[ctxName] = pfr
	}
}

func (jm *JsonMask) Mask(ctxName string, jsonData []byte) ([]byte, error) {

	if /* (len(jm.fullPathsRegistry) == 0 && len(jm.partialPathsRegistry) == 0) || */ len(jsonData) == 0 {
		return jsonData, nil
	}

	// Check to see if the context has some content somewhere
	// A nil map behaves like an empty map when reading
	_, ok := jm.fullPathsRegistry[ctxName]
	if !ok {
		_, ok = jm.partialPathsRegistry[ctxName]
		if !ok {
			return jsonData, nil
		}
	}

	var target interface{}
	err := json.Unmarshal(jsonData, &target)
	if err != nil {
		return nil, err
	}

	err = jm.walkThrough(ctxName, "", nil, target)
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
	up, ndxs := ParsePath(path)

	//if len(jm.fullPathsRegistry) > 0 {
	// A nil map behaves like an empty map when reading
	fr, ok := jm.fullPathsRegistry[ctxName]
	if ok {
		if fi, ok := fr[up]; ok {
			if cmpFieldIndexes(fi.indexes, ndxs) {
				return fi, true
			}

			return FieldInfo{}, false
			/*
				if len(fi.indexes) > 0 {
					for i, val := range fi.indexes {
						if val != -1 && ndxs[i] != val {
							return FieldInfo{}, false
						}
					}
				}
				return fi, true
			*/
		}
	}
	// }

	// if len(jm.partialPathsRegistry) > 0 {
	// A nil map behaves like an empty map when reading
	pfr, ok := jm.partialPathsRegistry[ctxName]
	if ok {
		for _, fi := range pfr {
			if strings.HasSuffix(up, fi.uxPath) {
				if cmpFieldIndexes(fi.indexes, ndxs) {
					return fi, true
				}

				return FieldInfo{}, false
			}
		}
	}
	// }
	return FieldInfo{}, false
}

// cmpFieldIndexes compare two arrays of integers. They are expected of generally being of the same size. If not the first is expected to have fewer elements.
// comparison happens on the rightmost part of the array.
func cmpFieldIndexes(ndxs1, ndxs2 []int) bool {
	if len(ndxs1) > 0 && len(ndxs2) > 0 {
		ndxs2Offset := len(ndxs2) - len(ndxs1)
		if ndxs2Offset < 0 {
			panic(fmt.Errorf("index arrays have wrong elements to compare %v, %v", ndxs1, ndxs2))
		}
		for i := 0; i < len(ndxs1); i++ {
			val := ndxs1[i]
			val2 := ndxs2[i+ndxs2Offset]
			if val != -1 && val2 != val {
				return false
			}
		}
	}

	return true
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
