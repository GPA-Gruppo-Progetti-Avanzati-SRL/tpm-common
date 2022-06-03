package dotnotation_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/dotnotation"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPath(t *testing.T) {

	sarr := []string{
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[].Id",
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[+].Id",
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[12].Id",
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[first].Id",
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[last].Id",
		"CdtrPmtActvtnReq.PmtInf.Dbtr.Id.PrvtId.Othr[err].Id",
	}

	for i, s := range sarr {
		xp, err := dotnotation.NewPath(s)
		require.NoError(t, err)

		t.Log(i, xp)
	}
}
