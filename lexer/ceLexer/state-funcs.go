package ceLexer

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/lexer"
	"github.com/rs/zerolog/log"
	"unicode"
)

const (
	EmptyToken lexer.TokenType = iota
	ErrorToken
	StringToken
	IdentifierToken
	IntegerToken
	DecimalToken

	GenericPunctuationToken
	LParenPunctuationToken
	RParenPunctuationToken
	SquareLParenPunctuationToken
	SquareRParenPunctuationToken
	CurlyLParenPunctuationToken
	CurlyRParenPunctuationToken
	SingleQuotePunctuationToken
	AtPunctuationToken
	LtSymbolToken
	GtSymbolToken
	CommaPunctuationToken
	ColonPunctuationToken
	SemiColoPunctuationToken
	DotPunctuationToken
	PipeSymbolToken
	PlusSymbolToken
	MinusPunctuationToken
	StarPunctuationToken
	PercentPunctuationToken
	XorSymbolToken
	AmpersandPunctuationToken
	DoubleQuotePunctuationToken
	SlashPunctuationToken
	EqualSymbolToken
	EqualToOperatorToken
	ExclamationMarkPunctuationToken
	NotEqualToOperatorToken
	LteOperatorToken
	GteOperatorToken
	XmlTagAutoCLoseElementPunctuationToken
	XmlTagStartCLoseElementPunctuationToken
	LogicalAndOperatorToken
	LogicalOrOperatorToken
)

var punctuations = map[rune]lexer.TokenType{
	'(':  LParenPunctuationToken,
	')':  RParenPunctuationToken,
	'[':  SquareLParenPunctuationToken,
	']':  SquareRParenPunctuationToken,
	'{':  CurlyLParenPunctuationToken,
	'}':  CurlyRParenPunctuationToken,
	'\'': SingleQuotePunctuationToken,
	'@':  AtPunctuationToken,
	'<':  LtSymbolToken,
	'>':  GtSymbolToken,
	',':  CommaPunctuationToken,
	';':  SemiColoPunctuationToken,
	':':  ColonPunctuationToken,
	'.':  DotPunctuationToken,
	'=':  EqualSymbolToken,
	'!':  ExclamationMarkPunctuationToken,
	'|':  PipeSymbolToken,
	'/':  SlashPunctuationToken,
	'"':  DoubleQuotePunctuationToken,
	'+':  PlusSymbolToken,
	'-':  MinusPunctuationToken,
	'*':  StarPunctuationToken,
	'%':  PercentPunctuationToken,
	'^':  XorSymbolToken,
	'&':  AmpersandPunctuationToken,
}

func punctuationToken(r rune) lexer.TokenType {
	tt, ok := punctuations[r]
	if !ok {
		return GenericPunctuationToken
	}

	return tt
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isPunct(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isIdentifierCharacter(r rune, firstCharacter bool) bool {
	if firstCharacter {
		return (r >= 'a' && r <= 'z') || r == '_' || (r >= 'A' && r <= 'Z')
	}

	return (r >= 'a' && r <= 'z') || r == '_' || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func nextIs(l *lexer.L, r rune) bool {
	if l.Peek() == r {
		l.Next()
		return true
	}

	return false
}

func ZeroState(l *lexer.L) lexer.StateFunc {
	const semLogContext = "ce-lexer::zero-state"
	log.Trace().Msg(semLogContext)

	var stateFn lexer.StateFunc
	r := l.Next()
	if r == lexer.EOFRune {
		return nil
	}

	switch {
	case isSpace(r):
		stateFn = WhitespaceState
	case isIdentifierCharacter(r, true):
		stateFn = IdentifierState
	case isDigit(r):
		stateFn = NumberState
	case isPunct(r):
		stateFn = ZeroState
		tt := punctuationToken(r)
		switch tt {
		case DoubleQuotePunctuationToken:
			stateFn = StringState
		case LtSymbolToken:
			switch l.Next() {
			case '=':
				l.Emit(LteOperatorToken)
			case '>':
				l.Emit(NotEqualToOperatorToken)
			case '/':
				l.Emit(XmlTagStartCLoseElementPunctuationToken)
			default:
				l.Rewind()
				l.Emit(tt)
			}
		case DotPunctuationToken:
			pr := l.Peek()
			if pr >= '0' && pr <= '9' {
				stateFn = DecimalState
			} else {
				l.Emit(tt)
			}
		case GtSymbolToken:
			if nextIs(l, '=') {
				l.Emit(GteOperatorToken)
			} else {
				l.Emit(tt)
			}
		case EqualSymbolToken:
			if nextIs(l, '=') {
				l.Emit(EqualToOperatorToken)
			} else {
				l.Emit(tt)
			}
		case ExclamationMarkPunctuationToken:
			if nextIs(l, '=') {
				l.Emit(NotEqualToOperatorToken)
			} else {
				l.Emit(tt)
			}
		case SlashPunctuationToken:
			if nextIs(l, '>') {
				l.Emit(XmlTagAutoCLoseElementPunctuationToken)
			} else {
				l.Emit(SlashPunctuationToken)
			}
		case AtPunctuationToken:
			pr := l.Peek()
			if isIdentifierCharacter(pr, true) {
				stateFn = AtIdentifierState
			} else {
				l.Emit(tt)
			}
		case AmpersandPunctuationToken:
			if nextIs(l, '&') {
				l.Emit(LogicalAndOperatorToken)
			} else {
				l.Emit(tt)
			}
		case PipeSymbolToken:
			if nextIs(l, '|') {
				l.Emit(LogicalOrOperatorToken)
			} else {
				l.Emit(tt)
			}
		default:
			l.Emit(tt)
		}
	default:
		l.Emit(ErrorToken)
		stateFn = nil
	}

	return stateFn
}

func WhitespaceState(l *lexer.L) lexer.StateFunc {

	const semLogContext = "ce-lexer::whitespace"
	log.Trace().Msg(semLogContext)

	/*
		    r := l.Next()
			if r == lexer.EOFRune {
				return nil
			}

			if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
				l.Error(fmt.Sprintf("unexpected token %q", r))
				return nil
			}
	*/

	l.Take(" \t\n\r")
	l.Ignore()

	return ZeroState
}

func NumberState(l *lexer.L) lexer.StateFunc {

	const semLogContext = "ce-lexer::number-state"
	log.Trace().Msg(semLogContext)

	r := l.Next()
	for r >= '0' && r <= '9' {
		r = l.Next()
	}

	if r == lexer.EOFRune {
		l.Emit(IntegerToken)
		return nil
	}

	if r == '.' {
		return DecimalState(l)
	}

	if isSpace(r) || isPunct(r) {
		l.Rewind()
		l.Emit(IntegerToken)
		return ZeroState
	}

	l.Emit(ErrorToken)
	return nil
}

func DecimalState(l *lexer.L) lexer.StateFunc {

	const semLogContext = "ce-lexer::decimal-state"
	log.Trace().Msg(semLogContext)

	r := l.Next()
	for r >= '0' && r <= '9' {
		r = l.Next()
	}

	if r == lexer.EOFRune {
		l.Emit(DecimalToken)
		return nil
	}

	if r == '.' {
		l.Emit(ErrorToken)
		return nil
	}

	if isSpace(r) || isPunct(r) {
		l.Rewind()
		l.Emit(DecimalToken)
		return ZeroState
	}

	l.Emit(ErrorToken)
	return nil
}

func IdentifierState(l *lexer.L) lexer.StateFunc {

	const semLogContext = "ce-lexer::identifier-state"
	log.Trace().Msg(semLogContext)

	r := l.Next()
	for isIdentifierCharacter(r, false) {
		r = l.Next()
	}

	if r == '.' {
		return DottedIdentifierState
	}

	l.Rewind()
	l.Emit(IdentifierToken)

	if r == lexer.EOFRune {
		return nil
	}

	return ZeroState
}

func DottedIdentifierState(l *lexer.L) lexer.StateFunc {

	const semLogContext = "ce-lexer::dotted-identifier-state"
	log.Trace().Msg(semLogContext)

	r := l.Next()
	if !isIdentifierCharacter(r, true) {
		l.Emit(ErrorToken)
		return nil
	}

	r = l.Next()
	for isIdentifierCharacter(r, false) {
		r = l.Next()
	}

	if r == '.' {
		return DottedIdentifierState
	}

	l.Rewind()
	l.Emit(IdentifierToken)

	if r == lexer.EOFRune {
		return nil
	}

	return ZeroState
}

func AtIdentifierState(l *lexer.L) lexer.StateFunc {
	const semLogContext = "ce-lexer::at-identifier-state"
	log.Trace().Msg(semLogContext)

	var tt lexer.TokenType
	r := l.Next()
	if (r >= 'a' && r <= 'z') || r == '_' || (r >= 'A' && r <= 'Z') {
		r = l.Next()
		for (r >= 'a' && r <= 'z') || r == '_' || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			r = l.Next()
		}
		tt = IdentifierToken
	} else {
		tt = GenericPunctuationToken
	}

	l.Rewind()
	l.Emit(tt)

	if r == lexer.EOFRune {
		return nil
	}

	return ZeroState
}

func StringState(l *lexer.L) lexer.StateFunc {
	const semLogContext = "ce-lexer::string-state"
	log.Trace().Msg(semLogContext)

	r := l.Next()
	for r != lexer.EOFRune && r != '"' {
		if r == '\\' {
			r = l.Next()
			if r == '"' {
				r = l.Next()
			}
		}
		r = l.Next()
	}

	if r == lexer.EOFRune {
		return l.Errorf(ErrorToken, "string not properly terminated")
	}

	l.Emit(StringToken)
	return ZeroState
}
