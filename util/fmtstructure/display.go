package fmtstructure

import (
	"fmt"
	"reflect"
	"strings"
)

type TagInfo struct {
	Key       string
	Value     string
	OmitEmpty bool
}

func (ti TagInfo) IsZero() bool {
	return ti.Value == ""
}

func FormatStructure(name string, x interface{}, wrtr Writer) {
	fmt.Printf("FormatStructure %s (%T):\n", name, x)
	format(name, reflect.ValueOf(x), TagInfo{}, wrtr)
}

func format(path string, v reflect.Value, tgi TagInfo, wrtr Writer) {

	quoted := false
	switch v.Kind() {
	case reflect.Invalid:
		wrtr.WritePathValue(path, v.Kind(), "invalid", "", false, tgi)
		// fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		//isEmpty := false
		//if v.IsZero() {
		//	isEmpty = true
		//}
		if v.Type().String() == "primitive.ObjectID" {
			wrtr.WritePathValue(path, reflect.String, "primitive.ObjectID", fmt.Sprint(v), false, tgi)
		} else {
			wrtr.WritePathValue(path, v.Kind(), "Array", "", v.IsZero(), tgi)
			for i := 0; i < v.Len(); i++ {
				format(fmt.Sprintf("%s[%d]", path, i), v.Index(i), tgi, wrtr)
			}
		}
	case reflect.Struct:
		//isEmpty := false
		//if v.IsZero() {
		//	isEmpty = true
		//}
		wrtr.WritePathValue(path, v.Kind(), "Object", "", v.IsZero(), tgi)
		for i := 0; i < v.NumField(); i++ {
			fn := v.Type().Field(i).Name
			ti := parseTag(v.Type().Field(i).Tag, "bson", "json")
			if ti.IsZero() {
				ti.Value = fn
			}
			fieldPath := fmt.Sprintf("%s.%s", path, ti.Value)
			format(fieldPath, v.Field(i), ti, wrtr)
		}
	case reflect.Map:
		//isEmpty := false
		//if v.IsZero() {
		//	isEmpty = true
		//}
		wrtr.WritePathValue(path, v.Kind(), "Map", "", v.IsZero(), tgi)
		for _, key := range v.MapKeys() {
			format(fmt.Sprintf("%s[%s]", path, formatAtom(key, quoted)), v.MapIndex(key), tgi, wrtr)
		}
	case reflect.Ptr:
		if v.IsNil() {
			wrtr.WritePathValue(path, v.Kind(), "ptr", "nil", true, tgi)
			// fmt.Printf("%s = nil\n", path)
		} else {
			format(fmt.Sprintf("(*%s)", path), v.Elem(), tgi, wrtr)
		}
	case reflect.Interface:
		if v.IsNil() {
			// fmt.Printf("%s = nil\n", path)
			wrtr.WritePathValue(path, v.Kind(), "interface{}", "nil", true, tgi)
		} else {
			//fmt.Printf("### %s.type = %s\n", path, v.Elem().Type())
			wrtr.WritePathValue(path, v.Kind(), fmt.Sprintf("%s", v.Elem().Type()), "", false, tgi)
			format(path+".value", v.Elem(), tgi, wrtr)
		}
	default: // basic types, channels, funcs
		//isEmpty := false
		//if v.IsZero() {
		//	isEmpty = true
		//}
		wrtr.WritePathValue(path, v.Kind(), fmt.Sprintf("%s", v.Type()), formatAtom(v, quoted), v.IsZero(), tgi)
		// fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func parseTag(tag reflect.StructTag, keys ...string) TagInfo {

	tg := TagInfo{}

	var t string
	var ok bool
	for _, k := range keys {
		t, ok = tag.Lookup(k)
		if ok {
			tg.Key = k
			break
		}
	}

	if tg.Key == "" {
		return tg
	}

	tagParts := strings.Split(t, ",")
	tg.Value = tagParts[0]

	for i := 1; i < len(tagParts); i++ {
		if tagParts[i] == "omitempty" {
			tg.OmitEmpty = true
			return tg
		}
	}

	return tg
}
