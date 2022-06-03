package dotnotation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type IndexingTypeEnum string

const (
	None       IndexingTypeEnum = "none"
	Empty                       = "empty"
	Add                         = "add"
	IndexValue                  = "ndx-value"
	First                       = "first"
	Last                        = "last"
)

type Element struct {
	Name          string
	IndexingType  IndexingTypeEnum
	IndexingValue int
}

type DotPath struct {
	BasePath string
	Elems    []Element
}

var IndexedRegExpr = regexp.MustCompile("\\[([0-9]+|\\+|first|last)\\]")

func NewPath(p string) (DotPath, error) {
	xp := DotPath{}

	els := strings.Split(p, ".")
	elsInfo := make([]Element, len(els))
	var sb strings.Builder
	for i, e := range els {

		withIndexAnnotation := false
		a := Element{IndexingType: None}
		if ndx := strings.Index(e, "["); ndx >= 0 {
			a.Name = e[0:ndx]
			withIndexAnnotation = true
		} else {
			a.Name = e
		}

		if withIndexAnnotation {
			var err error
			a.IndexingType, a.IndexingValue, err = getIndexAnnotation(e)
			if err != nil {
				return xp, err
			}
		}

		if i > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(a.Name)

		/*
			if strings.HasSuffix(e, "[]") {
				a.IndexingType = Empty
			} else {

			}

			if strings.HasSuffix(e, "[+]") {
				a.IndexingType = Add
			} else if strings.HasSuffix(e, "[]") {
				a.IndexingType = Empty
			} else {
				matches := IndexedRegExpr.FindAllSubmatch([]byte(e), -1)
				if len(matches) > 0 {
					a.IndexingType = IndexValue

					var err error
					a.IndexingValue, err = strconv.Atoi(string(matches[0][1]))
					if err != nil {
						return xp, err
					}
				}
			}
		*/

		elsInfo[i] = a
	}

	xp.BasePath = sb.String()
	xp.Elems = elsInfo
	return xp, nil
}

func getIndexAnnotation(expr string) (IndexingTypeEnum, int, error) {

	if strings.HasSuffix(expr, "[]") {
		return Empty, 0, nil
	}

	matches := IndexedRegExpr.FindAllSubmatch([]byte(expr), -1)
	if len(matches) > 0 {
		switch string(matches[0][1]) {
		case "+":
			return Add, 0, nil
		case "first":
			return First, 0, nil
		case "last":
			return Last, 0, nil
		}

		var err error
		ndxVal, err := strconv.Atoi(string(matches[0][1]))
		if err != nil {
			return Empty, 0, err
		}

		return IndexValue, ndxVal, nil
	}

	return Empty, 0, fmt.Errorf("unmatched index annotation expression %s", expr)
}
