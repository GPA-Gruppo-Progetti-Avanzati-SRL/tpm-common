package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"

	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

func TestResolveConfigValue(t *testing.T) {

	os.Setenv("MY_VAR", "a my_var value")
	os.Setenv("THIRD_ONE", "a third value")
	os.Setenv("SECOND_ONE", "")

	sarr := []string{
		"${MY_VAR}",
		"${NOT_DEFINED_VAR}",
		"my var: ${MY_VAR}, second one: ${SECOND_ONE}, third var: ${THIRD_ONE}",
	}

	for _, s := range sarr {
		t.Logf("resolving %s --> %s", s, util.ResolveConfigValue(s))
	}
}

var fileContent = []byte(`
First line of file ${MY_VAR},
  Second line of file ${MY_VAR} - ${THIRD_ONE}
Third line: ${MY_VAR}, ${MY_VAR}, ${SECOND_ONE}
 this is not a {var}
`)

func TestResolveEnvVars(t *testing.T) {
	os.Setenv("MY_VAR", "a my_var value")
	os.Setenv("THIRD_ONE", "a third value")
	os.Setenv("SECOND_ONE", "")

	fn := "./test-file.txt"
	err := ioutil.WriteFile(fn, fileContent, fs.ModePerm)
	require.NoError(t, err)
	defer os.Remove(fn)

	b, err := util.ReadFileAndResolveEnvVars(fn)
	require.NoError(t, err)

	t.Log(string(b))

	fn, b, err = util.ReadConfig("NOT_EXISTENT_VAR", fileContent, true)
	require.NoError(t, err)

	t.Log("filename: ", fn)
	t.Log(string(b))
}

func TestRegexp(t *testing.T) {

	sarr := []string{
		"${MY_VAR}",
		"string without vars",
		"my var ${MY_VAR} Whatever ${SECOND_ONE}",
	}

	for _, s := range sarr {
		t.Logf("expression: %s", s)
		matches := util.ConfigValueRegexp.FindAllSubmatch([]byte(s), -1)
		for i, m := range matches {
			t.Logf("[%d] Len match: %d --> %s --> %s", i, len(m), m[1], m[2])
		}
	}
}
