package reader

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"strings"
)

type Discriminator interface {
	DiscriminateLine(lineno int, l string, recs []fixedlengthfile.FixedLengthRecordDefinition) (string, error)
}

type DiscriminatorFunc func(lineno int, l string, recs []fixedlengthfile.FixedLengthRecordDefinition) (string, error)

func (f DiscriminatorFunc) DiscriminateLine(lineno int, l string, recs []fixedlengthfile.FixedLengthRecordDefinition) (string, error) {
	return f(lineno, l, recs)
}

func PrefixDiscriminator(lineno int, l string, recs []fixedlengthfile.FixedLengthRecordDefinition) (string, error) {

	var recClue string
	if len(l) < 10 {
		recClue = l
	} else {
		recClue = l[:10]
	}

	for _, r := range recs {
		if r.PrefixDiscriminator == "" {
			return ErrRecordId, fmt.Errorf("record %s doesn't specify a pprefix to discriminate from", r.Id)
		}

		if len(l) >= len(r.PrefixDiscriminator) {
			if strings.Contains(r.PrefixDiscriminator, "*") {
				if util.HasPrefixWithWildCard(l, r.PrefixDiscriminator, '*') {
					return r.Id, nil
				}
			} else if strings.HasPrefix(l, r.PrefixDiscriminator) {
				return r.Id, nil
			}
		}
	}

	return ErrRecordId, fmt.Errorf("cannot discriminate record by prefix ('%s...') for line at %d", recClue, lineno)
}
