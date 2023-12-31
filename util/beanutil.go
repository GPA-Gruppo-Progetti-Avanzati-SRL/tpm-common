package util

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const (
	SemLogPropertyName        = "property-path"
	SemLogPropertySegmentName = "segment-name"

	ArrayItemNoneSpecifier  = -2
	ArrayItemEmptySpecifier = -1
	ArrayItemIndexSpecifier = 0
	ArrayItemStarSpecifier  = -3
)

func beanUtilPropertyNameKind(n string) (string, int) {
	if strings.HasSuffix(n, "[]") {
		return strings.TrimSuffix(n, "[]"), ArrayItemEmptySpecifier
	}

	if strings.HasSuffix(n, "[*]") {
		return strings.TrimSuffix(n, "[*]"), ArrayItemStarSpecifier
	}

	if strings.HasSuffix(n, "]") {
		sarr := strings.Split(strings.TrimSuffix(n, "]"), "[")
		switch len(sarr) {
		case 1:
			return sarr[0], ArrayItemEmptySpecifier
		case 2:
			i, err := strconv.Atoi(sarr[1])
			if err != nil {
				return sarr[0], ArrayItemEmptySpecifier
			}

			if i < 0 {
				return sarr[0], ArrayItemEmptySpecifier
			}

			return sarr[0], i
		default:
			return sarr[0], ArrayItemEmptySpecifier
		}
	}

	return n, ArrayItemNoneSpecifier
}

func GetProperty(sourceMap map[string]interface{}, propName string) interface{} {

	propSegments := strings.Split(propName, ".")

	m := sourceMap
	for i := 0; i < len(propSegments)-1; i++ {

		pn, pnKind := beanUtilPropertyNameKind(propSegments[i])

		/*
			if strings.HasSuffix(propSegments[i], "[]") {
				mustArray = true
				propSegments[i] = strings.TrimSuffix(propSegments[i], "[]")
			}
		*/

		if v, ok := m[pn]; ok {
			m = cast2Map(propName, pn, v, pnKind)
			if m == nil {
				return nil
			}
		} else {
			return nil
		}
	}

	pn, pnKind := beanUtilPropertyNameKind(propSegments[len(propSegments)-1])
	v, ok := m[pn]
	if ok {
		switch tv := v.(type) {
		case []string:
			switch pnKind {
			case ArrayItemNoneSpecifier:
				v = tv[0]
			case ArrayItemEmptySpecifier:
				v = tv[0]
			case ArrayItemStarSpecifier:
				v = strings.Join(tv, ",")
			default:
				if len(tv) > pnKind {
					v = tv[pnKind]
				} else {
					log.Error().Int("len-array", len(tv)).Int("index", pnKind).Msg("array index out of bound")
					v = nil
				}

			}
		case []interface{}:
			switch pnKind {
			case ArrayItemNoneSpecifier:
				v = nil
			case ArrayItemEmptySpecifier:
				v = tv[0]
			case ArrayItemStarSpecifier:
				v = tv[0]
			default:
				if len(tv) > pnKind {
					v = tv[pnKind]
				} else {
					log.Error().Int("len-array", len(tv)).Int("index", pnKind).Msg("array index out of bound")
					v = nil
				}
			}
		default:
			if pnKind != ArrayItemNoneSpecifier {
				v = nil
			}
		}

		return v
		/*
			switch typedValue := v.(type) {
			case int:
				return fmt.Sprintf(ndx.Format, typedValue)
			case string:
				return typedValue
			default:
				log.Error().Interface("value", v).Str("type", fmt.Sprintf("%T", v)).Msg("unknown type")
			}

		*/
	}

	return nil
}

func SetProperty(targetMap map[string]interface{}, propertyPath string, value interface{}) error {

	propSegments := strings.Split(propertyPath, ".")

	m := targetMap
	for i := 0; i < len(propSegments)-1; i++ {
		mustArray := ArrayItemNoneSpecifier
		if strings.HasSuffix(propSegments[i], "[]") {
			mustArray = ArrayItemEmptySpecifier
			propSegments[i] = strings.TrimSuffix(propSegments[i], "[]")
		}

		if v, ok := m[propSegments[i]]; ok {
			m = cast2Map(propertyPath, propSegments[i], v, mustArray)
			if m == nil {
				return fmt.Errorf("error in creating patch metadata for %s", propertyPath)
			}
		} else {
			if mustArray == ArrayItemEmptySpecifier {
				arr := make([]interface{}, 0)
				arrItem := make(map[string]interface{})
				arr = append(arr, arrItem)
				m[propSegments[i]] = arr
				m = arrItem
			} else {
				subMap := make(map[string]interface{})
				m[propSegments[i]] = subMap
				m = subMap
			}
		}
	}

	m[propSegments[len(propSegments)-1]] = value
	return nil
}

func cast2Map(propPath string, propName string, v interface{}, propertyKind int) map[string]interface{} {
	var m map[string]interface{}
	switch tv := v.(type) {
	case map[string]interface{}:
		if propertyKind != ArrayItemNoneSpecifier {
			log.Error().Str(SemLogPropertySegmentName, propName).Str("name", propPath).Msg("expected array and found map")
		} else {
			m = tv
		}
	case []interface{}:
		if len(tv) != 0 {
			if v, ok := tv[0].(map[string]interface{}); ok {
				m = v
			} else {
				log.Error().Str(SemLogPropertySegmentName, propName).Str(SemLogPropertyName, propPath).Msg("error document finding indexed property")
			}
		}
	default:
		log.Error().Str(SemLogPropertySegmentName, propName).Str(SemLogPropertyName, propPath).Msg("expected property should be of type map o array")
	}

	return m
}

func GetIntProperty(sourceMap map[string]interface{}, p string) int {
	i := GetProperty(sourceMap, p)
	if i != nil {
		switch tv := i.(type) {
		case int:
			return tv
		case int32:
			return int(tv)
		case int64:
			return int(tv)
		case float64:
			return int(tv)
		}
	}

	return 0
}

func GetBoolProperty(sourceMap map[string]interface{}, p string, required bool) (bool, error) {
	i := GetProperty(sourceMap, p)
	if i != nil {
		switch tv := i.(type) {
		case bool:
			return tv, nil
		case string:
			switch strings.ToLower(tv) {
			case "true":
				return true, nil
			case "false":
				return false, nil
			case "":
				break
			}
		default:
			return false, fmt.Errorf("the property %s is not bool or string but %T", p, tv)
		}
	}

	if required {
		return false, fmt.Errorf("the property %s is not present", p)
	}

	return false, nil
}

func GetStringProperty(sourceMap map[string]interface{}, propName string) (string, error) {

	v := GetProperty(sourceMap, propName)
	if v == nil {
		return "", fmt.Errorf("%s cannot be found in document", propName)
	}

	propValue, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s  cannot be converted to string: %#v", propName, v)
	}

	return propValue, nil
}

func GetTimeProperty(sourceMap map[string]interface{}, propName string, propLayout string) (time.Time, error) {

	v := GetProperty(sourceMap, propName)
	if v == nil {
		return time.Time{}, fmt.Errorf("%s cannot be found in document", propName)
	}

	propValueStr, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("%s  cannot be converted to string: %#v", propName, v)
	}

	propValue, err := time.Parse(propLayout, propValueStr)
	if err != nil {
		return time.Time{}, err
	}

	return propValue, nil
}
