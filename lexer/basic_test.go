package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// Roughly based on the C64 Basic Language tokens
	IfStatementToken   = 0x8B
	LetStatementToken  = 0x88
	ForStatementToken  = 0x81
	ToStatementToken   = 0xA4
	NextStatementToken = 0x82
	AddSymbolToken     = 0xAA
	DivideSymbolToken  = 0xAD
	MinusSymbolToken   = 0xAB
	AsterixSymbolToken = 0xAC
	EqualsSymbolToken  = 0xB2
	// Additional tokens that aren't part of the C64 Basic tokens
	LeftParenthesis      = 0xF0
	RightParenthesis     = 0xF1
	LabelToken           = 0xF2
	ColonToken           = 0xF3
	StringVariableToken  = 0xF4
	IntegerVariableToken = 0xF5
	NumberVariableToken  = 0xF6
)

var KeywordTokens = map[string]TokenIdentifier{
	"if":   IfStatementToken,
	"let":  LetStatementToken,
	"for":  ForStatementToken,
	"to":   ToStatementToken,
	"next": NextStatementToken,
}

var VariableTypesToTokens = map[string]TokenIdentifier{
	"StringVariableToken":  StringVariableToken,
	"IntegerVariableToken": IntegerVariableToken,
	"IntegerLiteral":       IntegerLiteral,
	"NumberLiteral":        NumberLiteral,
	"StringLiteral":        StringLiteral,
}

var VariableTypes = map[TokenIdentifier]struct{}{
	StringVariableToken:  {},
	IntegerVariableToken: {},
}

var SingleCharTokens = map[rune]TokenIdentifier{
	'(': LeftParenthesis,
	')': RightParenthesis,
	'+': AddSymbolToken,
	'/': DivideSymbolToken,
	'-': MinusSymbolToken,
	'*': AsterixSymbolToken,
	'=': EqualsSymbolToken,
	':': ColonToken,
}

var comments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
}

func TestBasic(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine("let a = 10", 0)
	require.NoError(t, err)
	require.Equal(t, 5, len(tokens))
}

func TestIdentifierExpressionWithComment(t *testing.T) {
	sourceCode := "abc	// Define an identifier"
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(tokens))
}

func TestCrossLineComments(t *testing.T) {
	sourceCode := []string{
		"abc	/* multiline comments",
		"def */",
		"ghi",
	}
	l := NewBasicLexer()
	allTokens := make([]Token, 0)
	for _, s := range sourceCode {
		tokens, err := l.TokenizeLine(s, 0)
		require.NoError(t, err)
		allTokens = append(allTokens, tokens...)

	}
	require.Equal(t, 4, len(allTokens))
}

func TestFloatLexer(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "100.23"
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(tokens))
}

func NewBasicLexer() *Lexer {
	ll := NewLexerLanguage(
		WithSingleRuneMap(SingleCharTokens),
		WithCommentMap(comments),
		WithTokenCreators(identifierToken),
		WithLabelSettings(':', LabelToken),
	)
	return NewLexer(ll)
}

func identifierToken(identifier string) *Token {
	if tokenID, foundBasicKeyword := KeywordTokens[strings.ToLower(identifier)]; foundBasicKeyword {
		return NewToken(tokenID, identifier, nil)
	}
	if validIntegerVariableName(identifier) {
		return NewToken(IntegerVariableToken, identifier, 0)
	} else if validStringVariableName(identifier) {
		return NewToken(StringVariableToken, identifier, "")
	}
	return nil
}

func validStringVariableName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !IsIdentifierChar([]rune(identifier)[0], 0) {
		return false
	}
	if !strings.HasSuffix(identifier, "$") {
		return false
	}
	return true
}

func validIntegerVariableName(identifier string) bool {
	if validStringVariableName(identifier) || validLabelName(identifier) {
		return false
	}
	return true
}

func validLabelName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if IsIdentifierChar([]rune(identifier)[0], 0) {
		return false
	}
	if !strings.HasSuffix(identifier, ":") {
		return false
	}
	return true
}
