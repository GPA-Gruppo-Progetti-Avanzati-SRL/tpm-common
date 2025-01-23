package util_test

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestIsNumeric(t *testing.T) {

	sarr := []string{
		"T01009140904",
		"1+1",
		"+23.90",
	}

	for _, s := range sarr {
		t.Logf("Is numeric %s --> %t", s, util.IsNumeric(s))
	}
}

type InputWanted struct {
	input  string
	wanted string
	param  int
}

func TestStringJoin(t *testing.T) {

	s := []string{
		"0123456789", "ABCDEFGHIJ",
	}

	ns := util.StringJoin(s, "-", 15)
	require.Equal(t, 15, len(ns), "failed to join")
	t.Log(ns)

	ns = util.StringJoin(s, "-", -15)
	require.Equal(t, 15, len(ns), "failed to join")
	t.Log(ns)

	ns = util.StringJoin(s, "-", 0)
	require.Equal(t, 21, len(ns), "failed to join")
	t.Log(ns)
}

func TestMaxLengh(t *testing.T) {

	assert := assert.New(t)
	var s []InputWanted

	s = []InputWanted{
		{input: "0123456789", wanted: "01234", param: 5},
		{input: "0123456789", wanted: "0123456789", param: 10},
		{input: "0123456789", wanted: "0123456789", param: -10},
		{input: "0123456789", wanted: "789", param: -3},
		{input: "0123456789", wanted: "0123456789", param: 0},
	}

	for _, iw := range s {
		v, _ := util.ToMaxLength(iw.input, iw.param)
		fmt.Printf("%s (%d) --> %s\n", iw.input, iw.param, v)
		assert.Equal(iw.wanted, v, "to max length: strings should match")
	}
}

type InputWanted4PrefixWithWildCard struct {
	input        string
	prefix       string
	wildCardChar byte
	shouldMatch  bool
}

func TestHasPrefixWithWildCard(t *testing.T) {

	a := assert.New(t)
	var s []InputWanted4PrefixWithWildCard

	s = []InputWanted4PrefixWithWildCard{
		{input: "0123456789", prefix: "01234", wildCardChar: '*', shouldMatch: true},
		{input: "0123D56789", prefix: "01234", wildCardChar: '*', shouldMatch: false},
		{input: "0123456789", prefix: "01*34", wildCardChar: '*', shouldMatch: true},
	}

	for _, iw := range s {
		b := util.HasPrefixWithWildCard(iw.input, iw.prefix, iw.wildCardChar)
		a.Equal(iw.shouldMatch, b)
	}
}

func TestPadLengh(t *testing.T) {

	a := assert.New(t)
	var s []InputWanted

	s = []InputWanted{
		{input: "0123456789", wanted: "0123456789", param: 10},
		{input: "0123456789", wanted: "0123456789", param: 7},
		{input: "0123456789", wanted: "-----0123456789", param: -15},
		{input: "0123456789", wanted: "0123456789-----", param: 15},
		{input: "0123456789", wanted: "0123456789", param: 0},
	}

	for _, iw := range s {
		v, _ := util.Pad2Length(iw.input, iw.param, "-")
		fmt.Printf("%s (%d) --> %s\n", iw.input, iw.param, v)
		a.Equal(iw.wanted, v, "pad to length: strings should match")
	}
}

func TestStrings(t *testing.T) {

	a := assert.New(t)

	var s []InputWanted
	var modS string

	// Decamelize
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner_html"},
		{input: "action_name", wanted: "action_name"},
		{input: "css-class-name", wanted: "css-class-name"},
		{input: "my favorite items", wanted: "my favorite items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto bancoposta retail"},
		{input: "camt_029_001_09", wanted: "camt_029_001_09"},
	}

	for _, iw := range s {
		modS = util.Decamelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "decamelize: strings should match")
	}

	// Dasherize
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner-html"},
		{input: "action_name", wanted: "action-name"},
		{input: "css-class-name", wanted: "css-class-name"},
		{input: "my favorite items", wanted: "my-favorite-items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto-bancoposta-retail"},
		{input: "camt_029_001_09", wanted: "camt-029-001-09"},
	}

	for _, iw := range s {
		modS = util.Dasherize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "dasherize: strings should match")
	}

	// Camelize
	s = []InputWanted{
		{input: "innerHTML", wanted: "innerHTML"},
		{input: "action_name", wanted: "actionName"},
		{input: "css-class-name", wanted: "cssClassName"},
		{input: "my favorite items", wanted: "myFavoriteItems"},
		{input: "My Favorite Items", wanted: "myFavoriteItems"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "cONTOBANCOPOSTARETAIL"},
		{input: "conto bancoposta retail", wanted: "contoBancopostaRetail"},
		{input: "camt_029_001_09", wanted: "camt02900109"},
		{input: "camt.029.001.09", wanted: "camt02900109"},
	}

	for _, iw := range s {
		modS = util.Camelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "camelize: strings should match")
	}

	// Classify
	s = []InputWanted{
		{input: "innerHTML", wanted: "InnerHTML"},
		{input: "action_name", wanted: "ActionName"},
		{input: "css-class-name", wanted: "CssClassName"},
		{input: "my favorite items", wanted: "MyFavoriteItems"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "CONTOBANCOPOSTARETAIL"},
		{input: "camt_029_001_09", wanted: "Camt02900109"},
		{input: "camt.029.001.09", wanted: "Camt.029.001.09"},
	}

	for _, iw := range s {
		modS = util.Classify(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "classify: strings should match")
	}

	// Underscore
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner_html"},
		{input: "action_name", wanted: "action_name"},
		{input: "css-class-name", wanted: "css_class_name"},
		{input: "my favorite items", wanted: "my_favorite_items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto_bancoposta_retail"},
		{input: "camt_029_001_09", wanted: "camt_029_001_09"},
	}

	for _, iw := range s {
		modS = util.Underscore(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "underscore: strings should match")
	}

	s = []InputWanted{
		{input: "000034036586", wanted: "34036586"},
		{input: "000", wanted: ""},
	}

	for _, iw := range s {
		modS = util.TrimPrefixCharacters(iw.input, "0")
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "clear prefix")
	}

	s = []InputWanted{
		{input: "03403658600", wanted: "034036586"},
		{input: "000", wanted: ""},
	}

	for _, iw := range s {
		modS = util.TrimSuffixCharacters(iw.input, "0")
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "clear suffix")
	}

}

func TestStrings2(t *testing.T) {

	a := assert.New(t)

	s := []InputWanted{
		{input: "NUMERO GIORNI VALIDITA PRENOTAZIONE", wanted: "numeroGiorniValiditaPrenotazione"},
		{input: "FLAG RISCOSSIONE SPESE EMISSIONE", wanted: "flagRiscossioneSpeseEmissione"},
		{input: "NUMERO GIORNI STORNO EMISSIONE DIPENDENZA)", wanted: "numeroGiorniStornoEmissioneDipendenza)"},
		{input: "NUMERO GIORNI STORNO EMISSIONE UFF CENTRALE", wanted: "numeroGiorniStornoEmissioneUffCentrale"},
		{input: "FLAG RISCOSSIONE SPESE STAMPA", wanted: "flagRiscossioneSpeseStampa"},
		{input: "NUMERO GIORNI PER PRELEVAMENTO INTERESSI", wanted: "numeroGiorniPerPrelevamentoInteressi"},
		{input: "NUMERO GIORNI STORNO MOV CONTABILI DIPENDENZA", wanted: "numeroGiorniStornoMovContabiliDipendenza"},
		{input: "NUMERO GIORNI STORNO MOV CONTABILI UFF CENTRALE", wanted: "numeroGiorniStornoMovContabiliUffCentrale"},
		{input: "NUMERO GIORNI VALIDITA RICHIESTA AUTORIZZAZIONE", wanted: "numeroGiorniValiditaRichiestaAutorizzazione"},
		{input: "NUMERO GIORNI VALIDITA AUTORIZZAZIONE", wanted: "numeroGiorniValiditaAutorizzazione"},
		{input: "NUMERO ANNI PER PRESCRIZIONE CAPITALE", wanted: "numeroAnniPerPrescrizioneCapitale"},
		{input: "NUMERO ANNI PER PRESCRIZIONE INTERESSI", wanted: "numeroAnniPerPrescrizioneInteressi"},
		{input: "FLAG SEQUENZA CERTIFICATI", wanted: "flagSequenzaCertificati"},
		{input: "FLAG CIRCOLARITA", wanted: "flagCircolarita"},
		{input: "NUMERO GIORNI ARCHIVIAZIONE", wanted: "numeroGiorniArchiviazione"},
		{input: "NUM GG ARCHIVIAZIONE MOV CONT", wanted: "numGgArchiviazioneMovCont"},
		{input: "TIPO COLLEGAMENTO CON CONTI CORRENTI", wanted: "tipoCollegamentoConContiCorrenti"},
		{input: "TIPO COLLEGAMENTO CON ANAGRAFE GENERALE", wanted: "tipoCollegamentoConAnagrafeGenerale"},
		{input: "TIPO COLLEGAMENTO CON CONTABILITA GENERALE", wanted: "tipoCollegamentoConContabilitaGenerale"},
		{input: "TIPO COLLEGAMENTO CON TITOLI", wanted: "tipoCollegamentoConTitoli"},
		{input: "TIPO COLLEGAMENTO CON ESTERO", wanted: "tipoCollegamentoConEstero"},
		{input: "FLAG COLLEGAMENTO TITOLI-ESTERO", wanted: "flagCollegamentoTitoliEstero"},
		{input: "FLAG GESTIONE MODALITA DI REGOLAMENTO", wanted: "flagGestioneModalitaDiRegolamento"},
		{input: "FLAG GESTIONE ANTIRICICLAGGIO", wanted: "flagGestioneAntiriciclaggio"},
		{input: "SOGLIA 1 PER GESTIONE ANTIRICICLAGGIO", wanted: "soglia1PerGestioneAntiriciclaggio"},
		{input: "SOGLIA 2 PER GESTIONE ANTIRICICLAGGIO", wanted: "soglia2PerGestioneAntiriciclaggio"},
		{input: "IMPORTO DI SOGLIA NDG", wanted: "importoDiSogliaNdg"},
		{input: "TABELLA SEQUENZA STAMPE", wanted: "tabellaSequenzaStampe"},
		{input: "FLAG MEMORANDUM", wanted: "flagMemorandum"},
		{input: "NUMERO MASSIMO CERTIFICATI", wanted: "numeroMassimoCertificati"},
		{input: "FLAG DENSITA DI STAMPA", wanted: "flagDensitaDiStampa"},
		{input: "FLAG SABATO FERIALE", wanted: "flagSabatoFeriale"},
		{input: "NUMERO GIORNI PREAVVISO PER TITOLI", wanted: "numeroGiorniPreavvisoPerTitoli"},
		{input: "NUMERO GIORNI PREAVVISO PER ESTERO", wanted: "numeroGiorniPreavvisoPerEstero"},
		{input: "FLAG FINE VINCOLO FESTIVO", wanted: "flagFineVincoloFestivo"},
		{input: "FLAG MODALITA CALCOLO GIORNO SCADENZA", wanted: "flagModalitaCalcoloGiornoScadenza"},
		{input: "FLAG FORZATURA MODALITA REGOLAMENTO DA TITOLI", wanted: "flagForzaturaModalitaRegolamentoDaTitoli"},
		{input: "FLAG EMISSIONE DIRETTA", wanted: "flagEmissioneDiretta"},
		{input: "FLAG EMISSIONE DIFFERITA", wanted: "flagEmissioneDifferita"},
		{input: "FLAG STAMPA", wanted: "flagStampa"},
		{input: "CODICE DIPENDENZA GENERICA", wanted: "codiceDipendenzaGenerica"},
		{input: "CODICE DIVISA LIRE", wanted: "codiceDivisaLire"},
		{input: "CODICE DIVISA EURO", wanted: "codiceDivisaEuro"},
		{input: "CODICE DIVISA RIFERIMENTO", wanted: "codiceDivisaRiferimento"},
		{input: "DENOMINAZIONE SOCIALE", wanted: "denominazioneSociale"},
		{input: "GIORNI RETRODATAZIONE MCT DIP", wanted: "giorniRetrodatazioneMctDip"},
		{input: "CENSIMENTO NDG", wanted: "censimentoNdg"},
		{input: "TIPO CALCOLO PERIODO", wanted: "tipoCalcoloPeriodo"},
	}

	for _, iw := range s {
		modS := util.Camelize(strings.ToLower(iw.input))
		//fmt.Printf("%s --> %s\n", iw.input, modS)
		fmt.Printf("%s\n", modS)
		a.Equal(iw.wanted, modS, "camelize: strings should match")
	}
}
