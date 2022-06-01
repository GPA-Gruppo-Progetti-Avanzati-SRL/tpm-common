package dotnotation

import (
	"regexp"
	"strconv"
	"strings"
)

type Element struct {
	Name          string
	IndexingType  string
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

		a := Element{IndexingType: "none"}
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
			a.IndexingType = "plus"
		} else if strings.HasSuffix(e, "[]") {
			a.IndexingType = "empty"
		} else {
			matches := IndexedRegExpr.FindAllSubmatch([]byte(e), -1)
			if len(matches) > 0 {
				a.IndexingType = "index"

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
