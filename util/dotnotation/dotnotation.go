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
	Path  string
	Elems []Element
}

var IndexedRegExpr = regexp.MustCompile("\\[([0-9]+)\\]")

func NewPath(p string) (DotPath, error) {
	xp := DotPath{Path: p}

	els := strings.Split(p, ".")
	elsInfo := make([]Element, len(els))
	for i, e := range els {

		a := Element{Name: e, IndexingType: "none"}
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

	xp.Elems = elsInfo
	return xp, nil
}

func (xp DotPath) BasePath() string {
	return strings.ReplaceAll(xp.Path, "[]", "")
}
