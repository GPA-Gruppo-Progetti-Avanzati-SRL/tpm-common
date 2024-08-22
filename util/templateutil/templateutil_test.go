package templateutil_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/templateutil"
	varResolver "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/stretchr/testify/require"
	"testing"
)

var errTmpl = `{"_id":"66a0d746df0f009f3c298eab","awards":{"nominations":13,"text":"Won 2 Oscars. Another 7 wins \u0026 13 nominations.","wins":9},"cast":["Judy Garland","Frank Morgan","Ray Bolger","Bert Lahr"],"countries":["USA"],"directors":["Victor Fleming","George Cukor","Mervyn LeRoy","Norman Taurog","King Vidor"],"fullplot":"In this charming film based on the popular L. Frank Baum stories, Dorothy and her dog Toto are caught in a tornado's path and somehow end up in the land of Oz. Here she meets some memorable friends and foes in her journey to meet the Wizard of Oz who everyone says can help her return home and possibly grant her new friends their goals of a brain, heart and courage.","genres":["Adventure","Family","Fantasy"],"imdb":{"id":32138,"rating":8.1,"votes":262132},"languages":["English"],"lastupdated":"2015-09-01 00:29:05.490000000","metacritic":100,"num_mflix_comments":139,"plot":"Dorothy Gale is swept away to a magical land in a tornado and embarks on a quest to see the Wizard who can help her return home.","poster":"https://m.media-amazon.com/images/M/MV5BNjUyMTc4MDExMV5BMl5BanBnXkFtZTgwNDg0NDIwMjE@._V1_SY1000_SX677_AL_.jpg","rated":"PASSED","released":"-957916800000","runtime":102,"title":"The Wizard of Oz","tomatoes":{"consensus":"An absolute masterpiece whose groundbreaking visuals and deft storytelling are still every bit as resonant, The Wizard of Oz is a must-see film for young and old.","critic":{"meter":99,"numReviews":108,"rating":9.4},"dvd":"1999-10-19T00:00:00.000Z","fresh":107,"lastUpdated":"2015-09-12T17:53:21.000Z","production":"Warner Bros. Pictures","rotten":1,"viewer":{"meter":89,"numReviews":872115,"rating":3.7},"website":"http://thewizardofoz.warnerbros.com/"},"type":"movie","writers":["Noel Langley (screenplay)","Florence Ryerson (screenplay)","Edgar Allan Woolf (screenplay)","Noel Langley (adaptation)","L. Frank Baum (from the book by)"],"year":1939}

,"doesntExists": false
,"awards-nominations": true

{{if false }} Should not Come Out {}
{{if true }} ,"note": "awards-nomination is there" {}`

func TestErrTmpl(t *testing.T) {
	ti := []templateutil.Info{
		{
			Name:    "body",
			Content: errTmpl,
		},
	}

	pkgTemplate, err := templateutil.Parse(ti, nil)
	require.NoError(t, err)

	_, err = templateutil.Process(pkgTemplate, nil, false)
	require.NoError(t, err)
}

func TestTemplate(t *testing.T) {

	tmplInfo := []templateutil.Info{
		{Name: "t1", Content: "Subst: {{.prop_1}}, {{.prop2}}"},
		{Name: "t2", Content: "Subst: {{.prop_1}}, {{.prop2}}"},
	}

	data := make(map[string]interface{})
	data["prop_1"] = "value_1"
	data["prop_2"] = "value_2"

	for _, ti := range tmplInfo {
		tmpl, err := templateutil.Parse([]templateutil.Info{
			{Name: ti.Name, Content: ti.Content},
		}, nil)
		require.NoError(t, err)

		b, err := templateutil.Process(tmpl, data, false)
		require.NoError(t, err)

		t.Log(string(b))
	}

}

func TestPreprocessTemplate(t *testing.T) {

	tmplInfo := []templateutil.Info{
		{Name: "t1", Content: "Subst: {{.prop_1}}, <%=prop_2>"},
		{Name: "t2", Content: "Subst: <#=prop_1>, {{.prop2}}"},
	}

	data := make(map[string]interface{})
	data["prop_1"] = "value_1"
	data["prop_2"] = "value_2"

	for _, ti := range tmplInfo {

		tmplContent := templateutil.PreprocessVariableReferences(ti.Content, varResolver.AnyVariableReference)
		t.Log("pre-processed template content", tmplContent)
		tmpl, err := templateutil.Parse([]templateutil.Info{
			{Name: ti.Name, Content: tmplContent},
		}, nil)
		require.NoError(t, err)

		b, err := templateutil.Process(tmpl, data, false)
		require.NoError(t, err)

		t.Log(string(b))
	}

}
