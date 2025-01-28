package ceLexer_test

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/lexer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/lexer/ceLexer"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"testing"
	"unicode"
)

func TestUnicodeRune(t *testing.T) {

	rs := []rune{
		'(',
		')',
		'[',
		']',
		'{',
		'}',
		'\'',
		'@',
		'<',
		'>',
		',',
		';',
		':',
		'.',
		'=',
		'!',
		'|',
		'/',
		'"',
		'+',
		'-',
		'*',
		'%',
		'^',
		'&',
	}

	for _, r := range rs {
		switch {
		case unicode.IsLetter(r):
			t.Log("Unicode letter", r, string(r))
			// case r == '"':
			// t.Log("double quote", r, string(r))
		case unicode.IsSymbol(r):
			t.Log("symbol", r, string(r))
		case unicode.IsPunct(r):
			t.Log("Unicode punctuation", r, string(r))
		default:
			t.Log("Unicode else", r, string(r))
		}
	}
}

type testExpectedToken struct {
	tokType lexer.TokenType
	val     string
}

type testCase struct {
	text           string
	expectedTokens []testExpectedToken
}

func Test_BuiltinStateFuncs(t *testing.T) {
	const semLogContext = "test::built-in-state-funcs"
	var err error

	cases := []testCase{
		{
			text: `len(123., 37, "hello", @annotation)`,
			expectedTokens: []testExpectedToken{
				{ceLexer.IdentifierToken, "len"},
				{ceLexer.LParenPunctuationToken, "("},
				{ceLexer.DecimalToken, "123."},
				{ceLexer.CommaPunctuationToken, ","},
				{ceLexer.IntegerToken, "37"},
				{ceLexer.CommaPunctuationToken, ","},
				{ceLexer.StringToken, "\"hello\""},
				{ceLexer.CommaPunctuationToken, ","},
				{ceLexer.IdentifierToken, "@annotation"},
				{ceLexer.RParenPunctuationToken, ")"},
			},
		},
		{
			text: `{ src="newTest.kt" }`,
			expectedTokens: []testExpectedToken{
				{ceLexer.CurlyLParenPunctuationToken, "{"},
				{ceLexer.IdentifierToken, "src"},
				{ceLexer.EqualSymbolToken, "="},
				{ceLexer.StringToken, "\"newTest.kt\""},
				{ceLexer.CurlyRParenPunctuationToken, "}"},
			},
		},
	}

	for i, tc := range cases {

		var l *lexer.L
		l, err = ceLexer.NewLexer(tc.text)
		require.NoError(t, err)

		for j, c := range tc.expectedTokens {
			tok, done := l.NextToken()
			if done {
				err = errors.New("expected there to be more tokens, but there weren't")
				log.Fatal().Err(err).Int("_i", i).Int("_j", j).Str("et", ceLexer.TokenTypeString(c.tokType)).Str("tv", tok.Value).Str("tt", ceLexer.TokenTypeString(tok.Type)).Msg(semLogContext)
				return
			}

			if c.tokType != tok.Type {
				err := errors.New("token type != than expected")
				log.Fatal().Err(err).Int("_i", i).Int("_j", j).Str("et", ceLexer.TokenTypeString(c.tokType)).Str("tv", tok.Value).Str("tt", ceLexer.TokenTypeString(tok.Type)).Msg(semLogContext)
				return
			}

			if c.val != tok.Value {
				t.Errorf("Expected %q but got %q", c.val, tok.Value)
				err := errors.New("token value != than expected")
				log.Fatal().Err(err).Int("_i", i).Int("_j", j).Str("ev", c.val).Str("tv", tok.Value).Str("tt", ceLexer.TokenTypeString(tok.Type)).Msg(semLogContext)

				return
			}
		}

		tok, done := l.NextToken()
		if !done {
			t.Error("Expected the lexer to be done, but it wasn't.")
			return
		}

		if tok != nil {
			t.Errorf("Did not expect a token, but got %v", ceLexer.TokenTypeString(tok.Type))
			return
		}
	}
}
