package ceLexer

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/lexer"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EmptyToken-0]
	_ = x[ErrorToken-1]
	_ = x[StringToken-2]
	_ = x[IdentifierToken-3]
	_ = x[IntegerToken-4]
	_ = x[DecimalToken-5]
	_ = x[GenericPunctuationToken-6]
	_ = x[LParenPunctuationToken-7]
	_ = x[RParenPunctuationToken-8]
	_ = x[SquareLParenPunctuationToken-9]
	_ = x[SquareRParenPunctuationToken-10]
	_ = x[CurlyLParenPunctuationToken-11]
	_ = x[CurlyRParenPunctuationToken-12]
	_ = x[SingleQuotePunctuationToken-13]
	_ = x[AtPunctuationToken-14]
	_ = x[LtSymbolToken-15]
	_ = x[GtSymbolToken-16]
	_ = x[CommaPunctuationToken-17]
	_ = x[ColonPunctuationToken-18]
	_ = x[SemiColoPunctuationToken-19]
	_ = x[DotPunctuationToken-20]
	_ = x[PipeSymbolToken-21]
	_ = x[PlusSymbolToken-22]
	_ = x[MinusPunctuationToken-23]
	_ = x[StarPunctuationToken-24]
	_ = x[PercentPunctuationToken-25]
	_ = x[XorSymbolToken-26]
	_ = x[AmpersandPunctuationToken-27]
	_ = x[DoubleQuotePunctuationToken-28]
	_ = x[SlashPunctuationToken-29]
	_ = x[EqualSymbolToken-30]
	_ = x[EqualToOperatorToken-31]
	_ = x[ExclamationMarkPunctuationToken-32]
	_ = x[NotEqualToOperatorToken-33]
	_ = x[LteOperatorToken-34]
	_ = x[GteOperatorToken-35]
	_ = x[XmlTagAutoCLoseElementPunctuationToken-36]
	_ = x[XmlTagStartCLoseElementPunctuationToken-37]
}

const _TokenType_name = "EmptyTokenErrorTokenStringTokenIdentifierTokenIntegerTokenDecimalTokenGenericPunctuationTokenLParenPunctuationTokenRParenPunctuationTokenSquareLParenPunctuationTokenSquareRParenPunctuationTokenCurlyLParenPunctuationTokenCurlyRParenPunctuationTokenSingleQuotePunctuationTokenAtPunctuationTokenLtPunctuationTokenGtPunctuationTokenCommaPunctuationTokenColonPunctuationTokenSemiColoPunctuationTokenDotPunctuationTokenPipePunctuationTokenPlusPunctuationTokenMinusPunctuationTokenStarPunctuationTokenPercentPunctuationTokenXorPunctuationTokenAmpersandPunctuationTokenDoubleQuotePunctuationTokenSlashPunctuationTokenEqualPunctuationTokenEqualToPunctuationTokenExclamationMarkPunctuationTokenNotEqualToPunctuationTokenLtePunctuationTokenGtePunctuationTokenXmlTagAutoCLoseElementPunctuationTokenXmlTagStartCLoseElementPunctuationToken"

var _TokenType_index = [...]uint16{0, 10, 20, 31, 46, 58, 70, 93, 115, 137, 165, 193, 220, 247, 274, 292, 310, 328, 349, 370, 394, 413, 433, 453, 474, 494, 517, 536, 561, 588, 609, 630, 653, 684, 710, 729, 748, 786, 825}

func TokenTypeString(i lexer.TokenType) string {
	if i < 0 || i >= lexer.TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
