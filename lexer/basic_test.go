package lexer_test

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-lexer/lexer"
	"github.com/stretchr/testify/require"
)

// Constant declarations for different types of tokens
const (
	IfStatementToken     = 0x8B
	LetStatementToken    = 0x88
	ForStatementToken    = 0x81
	ToStatementToken     = 0xA4
	NextStatementToken   = 0x82
	AddSymbolToken       = 0xAA
	DivideSymbolToken    = 0xAD
	MinusSymbolToken     = 0xAB
	AsterixSymbolToken   = 0xAC
	EqualsSymbolToken    = 0xB2
	LeftParenthesis      = 0xF0
	RightParenthesis     = 0xF1
	LabelToken           = 0xF2
	ColonToken           = 0xF3
	StringVariableToken  = 0xF4
	IntegerVariableToken = 0xF5
	NumberVariableToken  = 0xF6
)

// KeywordTokens defines keyword to token mappings
var KeywordTokens = map[string]lexer.TokenIdentifier{
	"if":   IfStatementToken,
	"let":  LetStatementToken,
	"for":  ForStatementToken,
	"to":   ToStatementToken,
	"next": NextStatementToken,
}

// SingleCharTokens defines single character to token mappings
var SingleCharTokens = map[rune]lexer.TokenIdentifier{
	'(': LeftParenthesis,
	')': RightParenthesis,
	'+': AddSymbolToken,
	'/': DivideSymbolToken,
	'-': MinusSymbolToken,
	'*': AsterixSymbolToken,
	'=': EqualsSymbolToken,
	':': ColonToken,
}

// comments defines comment syntax mappings
var comments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
}

// TestBasicExpression tests a basic tokenization of a line
func TestBasicExpression(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine("let a = 10", 0)
	require.NoError(t, err)
	require.Equal(t, 4, len(tokens))
}

// TestIdentifierExpressionWithComment tests tokenization when a comment follows an identifier
func TestIdentifierExpressionWithComment(t *testing.T) {
	sourceCode := "abc	// Define an identifier"
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(tokens))
}

// TestIdentifierExpressionWithComment
func TestCrossLineComments(t *testing.T) {
	sourceCode := []string{
		"abc	/* multiline comments",
		"def */",
		"ghi",
	}
	l := NewBasicLexer()
	allTokens := make([]lexer.Token, 0)
	for _, s := range sourceCode {
		tokens, err := l.TokenizeLine(s, 0)
		require.NoError(t, err)
		allTokens = append(allTokens, tokens...)

	}
	require.Equal(t, 2, len(allTokens))
}

// TestFloatLexer tests tokenization of floating point numbers
func TestFloatLexer(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "100.23"
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(tokens))
}

// NewBasicLexer constructs a new Lexer using predefined language settings
func NewBasicLexer() *lexer.Lexer {
	ll := lexer.NewLexerLanguage(
		lexer.WithSingleRuneMap(SingleCharTokens),
		lexer.WithCommentMap(comments),
		lexer.WithTokenCreators(identifierToken),
		lexer.WithLabelSettings(':', LabelToken),
	)
	return lexer.NewLexer(ll)
}

// identifierToken handles the token creation for identifiers based on custom rules
func identifierToken(identifier string) *lexer.Token {
	if tokenID, foundBasicKeyword := KeywordTokens[strings.ToLower(identifier)]; foundBasicKeyword {
		return lexer.NewToken(tokenID, identifier, nil)
	}
	if validIntegerVariableName(identifier) {
		return lexer.NewToken(IntegerVariableToken, identifier, 0)
	} else if validStringVariableName(identifier) {
		return lexer.NewToken(StringVariableToken, identifier, "")
	}
	return nil
}

// validStringVariableName checks if an identifier is a valid string variable name
func validStringVariableName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !lexer.IsIdentifierChar([]rune(identifier)[0], 0) {
		return false
	}
	if !strings.HasSuffix(identifier, "$") {
		return false
	}
	return true
}

// validIntegerVariableName checks if an identifier is a valid integer variable name
func validIntegerVariableName(identifier string) bool {
	if validStringVariableName(identifier) || validLabelName(identifier) {
		return false
	}
	return true
}

// validLabelName checks if an identifier is a valid label name
func validLabelName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if lexer.IsIdentifierChar([]rune(identifier)[0], 0) {
		return false
	}
	if !strings.HasSuffix(identifier, ":") {
		return false
	}
	return true
}
