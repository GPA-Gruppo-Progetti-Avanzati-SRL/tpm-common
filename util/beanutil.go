package util

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	SemLogPropertyName        = "property-path"
	SemLogPropertySegmentName = "segment-name"
)

func GetProperty(sourceMap map[string]interface{}, propName string) interface{} {

	propSegments := strings.Split(propName, ".")

	m := sourceMap
	for i := 0; i < len(propSegments)-1; i++ {

		mustArray := false
		if strings.HasSuffix(propSegments[i], "[]") {
			mustArray = true
			propSegments[i] = strings.TrimSuffix(propSegments[i], "[]")
		}

		if v, ok := m[propSegments[i]]; ok {
			m = cast2Map(propName, propSegments[i], v, mustArray)
			if m == nil {
				return nil
			}
		} else {
			return nil
		}
	}

	v, ok := m[propSegments[len(propSegments)-1]]
	if ok {
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
		mustArray := false
		if strings.HasSuffix(propSegments[i], "[]") {
			mustArray = true
			propSegments[i] = strings.TrimSuffix(propSegments[i], "[]")
		}

		if v, ok := m[propSegments[i]]; ok {
			m = cast2Map(propertyPath, propSegments[i], v, mustArray)
			if m == nil {
				return fmt.Errorf("error in creating patch metadata for %s", propertyPath)
			}
		} else {
			if mustArray {
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

func cast2Map(propPath string, propName string, v interface{}, mustArray bool) map[string]interface{} {
	var m map[string]interface{}
	switch tv := v.(type) {
	case map[string]interface{}:
		if mustArray {
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
		case float64:
			return int(tv)
		}
	}

	return 0
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
