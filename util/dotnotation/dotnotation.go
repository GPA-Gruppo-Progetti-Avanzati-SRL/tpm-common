package dotnotation

import (
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

var IndexedRegExpr = regexp.MustCompile("\\[([0-9]+)\\]")

func NewPath(p string) (DotPath, error) {
	xp := DotPath{}

	els := strings.Split(p, ".")
	elsInfo := make([]Element, len(els))
	var sb strings.Builder
	for i, e := range els {

		a := Element{IndexingType: None}
		if ndx := strings.Index(e, "["); ndx >= 0 {
			a.Name = e[0:ndx]
		} else {
			a.Name = e
		}

		if i > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(a.Name)

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

		elsInfo[i] = a
	}

	xp.BasePath = sb.String()
	xp.Elems = elsInfo
	return xp, nil
}
