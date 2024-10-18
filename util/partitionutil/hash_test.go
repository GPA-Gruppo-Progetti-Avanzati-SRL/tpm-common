package partitionutil_test

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/partitionutil"
	"testing"
)

func TestHashPartition(t *testing.T) {

	var sarr []string

	for i := 0; i < 30; i++ {
		sarr = append(sarr, fmt.Sprintf(`{ "NDG" = "01293837%03d" }`, i))
	}

	for _, s := range sarr {
		t.Log("partition is ", partitionutil.HashPartition([]byte(s), 5))
	}

}
