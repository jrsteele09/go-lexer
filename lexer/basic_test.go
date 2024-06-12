package lexer_test

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-lexer/lexer"
	"github.com/stretchr/testify/require"
)

// Constant declarations for different types of tokens
const (
	IfStatementToken     lexer.TokenIdentifier = 0x8B
	LetStatementToken    lexer.TokenIdentifier = 0x88
	ForStatementToken    lexer.TokenIdentifier = 0x81
	ToStatementToken     lexer.TokenIdentifier = 0xA4
	NextStatementToken   lexer.TokenIdentifier = 0x82
	AddSymbolToken       lexer.TokenIdentifier = 0xAA
	DivideSymbolToken    lexer.TokenIdentifier = 0xAD
	MinusSymbolToken     lexer.TokenIdentifier = 0xAB
	AsterixSymbolToken   lexer.TokenIdentifier = 0xAC
	EqualsSymbolToken    lexer.TokenIdentifier = 0xB2
	LeftParenthesis      lexer.TokenIdentifier = 0xF0
	RightParenthesis     lexer.TokenIdentifier = 0xF1
	LabelToken           lexer.TokenIdentifier = 0xF2
	ColonToken           lexer.TokenIdentifier = 0xF3
	StringVariableToken  lexer.TokenIdentifier = 0xF4
	IntegerVariableToken lexer.TokenIdentifier = 0xF5
	NumberVariableToken  lexer.TokenIdentifier = 0xF6
	DollarToken          lexer.TokenIdentifier = 0xF7
	CommaToken           lexer.TokenIdentifier = 0xF8
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
	'$': DollarToken,
	',': CommaToken,
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
	require.Equal(t, LetStatementToken, tokens[0].ID)
	require.Equal(t, IntegerVariableToken, tokens[1].ID)
	require.Equal(t, EqualsSymbolToken, tokens[2].ID)
	require.Equal(t, lexer.IntegerLiteral, tokens[3].ID)
}

// TestIdentifierExpressionWithComment tests tokenization when a comment follows an identifier
func TestIdentifierExpressionWithComment(t *testing.T) {
	sourceCode := "abc	// Define an identifier"
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(tokens))
	require.Equal(t, IntegerVariableToken, tokens[0].ID)
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
	require.Equal(t, lexer.NumberLiteral, tokens[0].ID)
	require.Equal(t, 100.23, tokens[0].Value.(float64))
}

func TestAssemblerTokens(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "LDA ($FF),X"
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 7, len(tokens))
}

func TestHexTokens(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "0x1234"
	tokens, err := l.TokenizeLine(sourceCode, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(tokens))
	require.Equal(t, tokens[0].Value.(int16), int16(0x1234))
}

// NewBasicLexer constructs a new Lexer using predefined language settings
func NewBasicLexer() *lexer.Lexer {
	ll := lexer.NewLexerLanguage(
		lexer.WithSingleRuneMap(SingleCharTokens),
		lexer.WithCommentMap(comments),
		lexer.WithTokenCreators(identifierToken),
		lexer.WithLabelSettings(':', LabelToken),
		lexer.WithExtraIdentifierRunes("_#%"),
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
	if !lexer.IsIdentifierChar([]rune(identifier)[0], 0, "_") {
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
	if !lexer.IsIdentifierChar([]rune(identifier)[0], 0, "") {
		return false
	}
	if !strings.HasSuffix(identifier, ":") {
		return false
	}
	return true
}
