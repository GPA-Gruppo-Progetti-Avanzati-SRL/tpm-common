package fmtstructure

import (
	"fmt"
	"reflect"
)

func PrintfWritePathValue(aPath string, aKind reflect.Kind, aType string, aValue string, isEmpty bool, aTagInfo TagInfo) {

	if !(aTagInfo.OmitEmpty && isEmpty) {
		fmt.Printf("[%t] %s of type %s = %s\n", aTagInfo.OmitEmpty, aPath, aType, aValue)
	}

}
