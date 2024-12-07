package misc

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

var b = []byte(`{"indirizzo": "V_A C/I.SA;BA,R*R-T,? J Y§7)& $"}`)

type IntfA interface {
	MyContract() string
}

type TypeA struct {
}

func (ta TypeA) MyContract() string {
	return "hello"
}

func TestPippo(t *testing.T) {

	var v interface{}
	err := json.Unmarshal(b, &v)
	require.NoError(t, err)

	t.Log(v)
}
