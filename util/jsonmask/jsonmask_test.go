package jsonmask_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/jsonmask"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var jsonData = []byte(`
{
  "token-input": "202215131228-RPOL-RPST-DISP-DR-PP-girofondi-5bffff8b-9b74-47cc-95bb-75ba265c618",
  "timestamp": "2022-05-26 17:59:28.972",
  "sotto-sistema": "BPAP-RICARICHE-PPY",
  "identificativo-terminale": "127.0.0.1",
  "societa-consumer": "POSTEPAY",
  "identificativo-canale": "PPAY",
  "operazioni": [
    {
      "progressivo-ope": "1",
      "tipo-operazione": "RPSA",
      "importo-operazione": 5.10,
      "segno": "A",
      "codice-rapporto": "5333171000679775",
      "iban-rapporto": "IT34H3608105138283271283277",
      "codice-controparte": "5333171000653705",
      "iban-controparte": "IT60B3608105138229731829737",
      "motivazione-operazione": "causale ricarica SOGLIA",
      "identificativo-prodotto": "RICARICA",
      "soggetti": [
        {
          "progressivo-sog": "1",
          "ruolo": "ORDINANTE",
          "natura-giuridica": "NPF",
          "codice-fiscale": "PPPGTN80A01H501Q",
          "denominazione": "APPP AGOSTINO"
        },
        {
          "progressivo-sog": "2",
          "ruolo": "BENEFICIARIO",
          "natura-giuridica": "NPF",
          "codice-fiscale": "PPPGTN80A01H501Q",
          "denominazione": "APPP AGOSTINO"
        }
      ],
      "codici-fiscali": [
         "PPPGTN80A01H501Q",
         "PPPGTN80A01H501Q"
      ]
    }
  ]
}
`)

func TestMask(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	fields := []jsonmask.FieldInfo{
		{
			Path: ".operazioni.[].codice-rapporto",
		},
		{
			Path: ".operazioni.[].codice-controparte",
		},
		{
			Path: ".operazioni.[].soggetti.[].codice-fiscale",
		},
		{
			Path: ".operazioni.[].codici-fiscali.[]",
		},
		{
			Path: "*soggetti.[].denominazione",
		},
	}

	jm, err := jsonmask.NewJsonMask()
	require.NoError(t, err)
	jm.Add("request", fields)

	masked, err := jm.Mask("request", jsonData)
	require.NoError(t, err)

	t.Log(string(masked))
}

func TestPathParser(t *testing.T) {

	sarr := []string{
		".lev.[72].subLev.[99].thirdLev.[]",
	}

	for i, s := range sarr {
		ns, indxs := jsonmask.ParsePath(s)
		t.Log(i, ns, indxs)
	}
}

func TestShouldBeMasked(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/*
	 * Issue the map gets a lookupkey that is an indexed path. So we cannot have two overlapping path with different conditions.
	 * In the example below it works as if only the second is considered. For a key I should consider an array of fieldInfos....
	 */
	fields := []jsonmask.FieldInfo{
		{
			Path: ".lev.[72].subLev.[99].thirdLev.[]",
		},
		{
			Path: ".lev.[].subLev.[].thirdLev.[12]",
		},
		{
			Path: "*thirdLev.[]",
		},
	}

	jm, err := jsonmask.NewJsonMask()
	require.NoError(t, err)

	jm.Add("request", fields)

	sarr := []string{
		".lev.[72].subLev.[99].thirdLev.[0]",
		".lev.[71].subLev.[99].thirdLev.[1]",
		".lev.[4].subLev.[99].thirdLev.[12]",
	}

	for i, s := range sarr {
		_, b := jm.HasToBeMasked("request", s)
		t.Log(i, s, b)
	}

}

var cfgData = []byte(`
request:
  name: request
  fields:
     - path: ".operazioni.[].codice-rapporto"
     - path: ".operazioni.[].codice-controparte"
     - path: ".operazioni.[].soggetti.[].codice-fiscale"
     - path: ".operazioni.[].codici-fiscali.[]"
     - path: "*soggetti.[].denominazione"
`)

func TestRead(t *testing.T) {
	jm, err := jsonmask.NewJsonMask(jsonmask.FromData(cfgData))
	require.NoError(t, err)

	masked, err := jm.Mask("request", jsonData)
	require.NoError(t, err)

	t.Log(string(masked))
}

var errData = []byte(`
[0] Get "https://tpm-router-card-inquiry-api-common-card.app.coll2.ocprm.testposte:443/listaCarte/api
/v1/carta/012965854": dial tcp: lookup tpm-router-card-inquiry-api-common-card.app.coll2.ocprm.testposte
: no such host`)

func TestError(t *testing.T) {
	jm, err := jsonmask.NewJsonMask(jsonmask.FromData(cfgData))
	require.NoError(t, err)

	masked, err := jm.Mask("request", errData)
	require.NoError(t, err)

	t.Log(string(masked))
}

var errData1 = []byte(`
{
"aliasCarta" : 00123456
}
`)

func TestError1(t *testing.T) {
	jm, err := jsonmask.NewJsonMask(jsonmask.FromData(cfgData))
	require.NoError(t, err)

	masked, err := jm.Mask("request", errData1)
	require.NoError(t, err)

	t.Log(string(masked))
}
