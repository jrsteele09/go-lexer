package lexer_test

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-lexer/lexer"
	"github.com/jrsteele09/go-lexer/lexer/utils"
	"github.com/stretchr/testify/require"
)

// Constant declarations for different types of tokens
const (
	IfStatementToken        lexer.TokenIdentifier = 0x8B
	LetStatementToken       lexer.TokenIdentifier = 0x88
	ForStatementToken       lexer.TokenIdentifier = 0x81
	ToStatementToken        lexer.TokenIdentifier = 0xA4
	NextStatementToken      lexer.TokenIdentifier = 0x82
	AddSymbolToken          lexer.TokenIdentifier = 0xAA
	DivideSymbolToken       lexer.TokenIdentifier = 0xAD
	MinusSymbolToken        lexer.TokenIdentifier = 0xAB
	MultiplySymbolToken     lexer.TokenIdentifier = 0xAC
	EqualsSymbolToken       lexer.TokenIdentifier = 0xB2
	LeftParenthesis         lexer.TokenIdentifier = 0xF0
	RightParenthesis        lexer.TokenIdentifier = 0xF1
	LabelToken              lexer.TokenIdentifier = 0xF2
	ColonToken              lexer.TokenIdentifier = 0xF3
	StringVariableToken     lexer.TokenIdentifier = 0xF4
	IntegerVariableToken    lexer.TokenIdentifier = 0xF5
	NumberVariableToken     lexer.TokenIdentifier = 0xF6
	DollarToken             lexer.TokenIdentifier = 0xF7
	CommaToken              lexer.TokenIdentifier = 0xF8
	NotEqualToken           lexer.TokenIdentifier = 0xF9
	LeftCurlyBracket        lexer.TokenIdentifier = 0xFA
	RightCurlyBracket       lexer.TokenIdentifier = 0xFB
	LeftSquareBracket       lexer.TokenIdentifier = 0xFC
	RightSquareBracket      lexer.TokenIdentifier = 0xFD
	SemicolonToken          lexer.TokenIdentifier = 0xFE
	PeriodToken             lexer.TokenIdentifier = 0xFF
	DoubleQuoteToken        lexer.TokenIdentifier = 0x100
	SingleQuoteToken        lexer.TokenIdentifier = 0x101
	BacktickToken           lexer.TokenIdentifier = 0x102
	PipeToken               lexer.TokenIdentifier = 0x103
	AmpersandToken          lexer.TokenIdentifier = 0x104
	AsteriskToken           lexer.TokenIdentifier = 0x105
	ExclamationToken        lexer.TokenIdentifier = 0x106
	QuestionToken           lexer.TokenIdentifier = 0x107
	AtToken                 lexer.TokenIdentifier = 0x108
	HashToken               lexer.TokenIdentifier = 0x109
	LessThanToken           lexer.TokenIdentifier = 0x10A
	GreaterThanToken        lexer.TokenIdentifier = 0x10B
	BackslashToken          lexer.TokenIdentifier = 0x10C
	LessThanOrEqualToken    lexer.TokenIdentifier = 0x10D
	GreaterThanOrEqualToken lexer.TokenIdentifier = 0x10E
	IncrementToken          lexer.TokenIdentifier = 0x10F
	EqualityToken           lexer.TokenIdentifier = 0x110
	PercentageToken         lexer.TokenIdentifier = 0x120
)

// KeywordTokens defines keyword to token mappings
var KeywordTokens = map[string]lexer.TokenIdentifier{
	"if":   IfStatementToken,
	"let":  LetStatementToken,
	"for":  ForStatementToken,
	"to":   ToStatementToken,
	"next": NextStatementToken,
}

// OperatorTokens defines the Operator token mappings that can consist of multiple symbol tokens
var OperatorTokens = map[string]lexer.TokenIdentifier{
	"<=": LessThanOrEqualToken,
	">=": GreaterThanOrEqualToken,
	"<>": NotEqualToken,
	"++": IncrementToken,
	"==": EqualityToken,
}

// SymbolTokens defines single delimeter runes to token mappings
var SymbolTokens = map[rune]lexer.TokenIdentifier{
	'+': AddSymbolToken,
	'/': DivideSymbolToken,
	'-': MinusSymbolToken,
	'*': MultiplySymbolToken,
	'=': EqualsSymbolToken,
	'<': LessThanToken,
	'>': GreaterThanToken,
	'(': LeftParenthesis,
	')': RightParenthesis,
	'{': LeftCurlyBracket,
	'}': RightCurlyBracket,
	'[': LeftSquareBracket,
	']': RightSquareBracket,
	',': CommaToken,
	// ';':  SemicolonToken,
	':':  ColonToken,
	'.':  PeriodToken,
	'|':  PipeToken,
	'&':  AmpersandToken,
	'!':  ExclamationToken,
	'?':  QuestionToken,
	'@':  AtToken,
	'#':  HashToken,
	'$':  DollarToken,
	'\\': BackslashToken,
	'%':  PercentageToken,
}

// comments defines comment syntax mappings
// NOTE: Symbols used to construct comments should be in the SymbolTokens
var comments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
	";":   "\n",
}

// Custom tokenizers - On detection of the starting character, jump to a specific tokenizer.
// Custom tokenizers will override symbol tokens.
var customTokenizers = map[string]lexer.TokenizerFunc{
	"$":  lexer.HexTokenizer,
	"0x": lexer.HexTokenizer,
	"%0": lexer.BinaryTokenizer,
	"%1": lexer.BinaryTokenizer,
}

func TestSingleComment(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(";let a = 10", "testfile", 0)
	require.NoError(t, err)
	require.Len(t, tokens, 0)
}

func TestAssemblerTypeMultiLineComments(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(
		`;* = $2000
		;LDA #1`, "testfile", 0)
	require.NoError(t, err)
	require.Len(t, tokens, 0)
}

// TestBasicExpression tests a basic tokenization of a line
func TestBasicExpression(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine("let a = 10", "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 5, len(tokens))
	require.Equal(t, LetStatementToken, tokens[0].ID)
	require.Equal(t, IntegerVariableToken, tokens[1].ID)
	require.Equal(t, EqualsSymbolToken, tokens[2].ID)
	require.Equal(t, lexer.IntegerLiteral, tokens[3].ID)
	require.Equal(t, lexer.EndOfLineType, tokens[4].ID)
}

// TestStrings tests a basic tokenization of a line
func TestStrings(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine("\" hello \",`hello`,'hello',\"'hello'\"", "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 8, len(tokens))
	require.Equal(t, lexer.StringLiteral, tokens[0].ID)
	require.Equal(t, CommaToken, tokens[1].ID)
	require.Equal(t, lexer.StringLiteral, tokens[2].ID)
	require.Equal(t, CommaToken, tokens[3].ID)
	require.Equal(t, lexer.StringLiteral, tokens[4].ID)
	require.Equal(t, "hello", tokens[4].Value)
	require.Equal(t, CommaToken, tokens[5].ID)
	require.Equal(t, lexer.StringLiteral, tokens[6].ID)
	require.Equal(t, "'hello'", tokens[6].Value)
	require.Equal(t, lexer.EndOfLineType, tokens[7].ID)
}

func TestBinary(t *testing.T) {
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine("%01,%10", "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 4, len(tokens))
	require.Equal(t, lexer.IntegerLiteral, tokens[0].ID)
	require.Equal(t, uint8(1), tokens[0].Value)
	require.Equal(t, CommaToken, tokens[1].ID)
	require.Equal(t, lexer.IntegerLiteral, tokens[2].ID)
	require.Equal(t, uint8(2), tokens[2].Value)
	require.Equal(t, lexer.EndOfLineType, tokens[3].ID)
}

// TestIdentifierExpressionWithComment tests tokenization when a comment follows an identifier
func TestIdentifierExpressionWithComment(t *testing.T) {
	sourceCode := "abc	// Define an identifier"
	l := NewBasicLexer()
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(tokens))
	require.Equal(t, IntegerVariableToken, tokens[0].ID)
	require.Equal(t, lexer.EndOfLineType, tokens[1].ID)

}

// TestIdentifierExpressionWithComment
func TestCrossLineComments(t *testing.T) {
	sourceCode := `abc	/* multiline comments
	def */
	ghi`

	l := NewBasicLexer()
	r := strings.NewReader(sourceCode)
	allTokens, err := l.Tokenize(r, "testfile")
	require.NoError(t, err)
	require.Equal(t, 5, len(allTokens)) // Includes two endOfLine tokens + endOfFile token
}

// TestFloatLexer tests tokenization of floating point numbers
func TestFloatLexer(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "100.23"
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(tokens))
	require.Equal(t, lexer.NumberLiteral, tokens[0].ID)
	require.Equal(t, 100.23, tokens[0].Value.(float64))
	require.Equal(t, lexer.EndOfLineType, tokens[1].ID)
}

func TestAssemblerTokens(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "LDA ($FF),X"
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 7, len(tokens))
}

func TestHexTokens(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "0x1234"
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(tokens))
	require.Equal(t, tokens[0].Value.(uint16), uint16(0x1234))
	require.Equal(t, lexer.EndOfLineType, tokens[1].ID)
}

func TestOperators(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "=<>+-*/+++"
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 9, len(tokens))
	require.Equal(t, EqualsSymbolToken, tokens[0].ID)
	require.Equal(t, NotEqualToken, tokens[1].ID)
	require.Equal(t, AddSymbolToken, tokens[2].ID)
	require.Equal(t, MinusSymbolToken, tokens[3].ID)
	require.Equal(t, MultiplySymbolToken, tokens[4].ID)
	require.Equal(t, DivideSymbolToken, tokens[5].ID)
	require.Equal(t, IncrementToken, tokens[6].ID)
	require.Equal(t, AddSymbolToken, tokens[7].ID)
	require.Equal(t, lexer.EndOfLineType, tokens[8].ID)
}

func TestLabelParsing(t *testing.T) {
	l := NewBasicLexer()
	sourceCode := "test: if a == 10"
	tokens, err := l.TokenizeLine(sourceCode, "testfile", 0)
	require.NoError(t, err)
	require.Equal(t, 6, len(tokens))
	require.Equal(t, LabelToken, tokens[0].ID)
	require.Equal(t, IfStatementToken, tokens[1].ID)
	require.Equal(t, IntegerVariableToken, tokens[2].ID)
	require.Equal(t, EqualityToken, tokens[3].ID)
	require.Equal(t, lexer.IntegerLiteral, tokens[4].ID)
	require.Equal(t, lexer.EndOfLineType, tokens[5].ID)
}

// NewBasicLexer constructs a new Lexer using predefined language settings
func NewBasicLexer() *lexer.Lexer {
	ll := lexer.NewLexerLanguage(
		lexer.WithKeywords(KeywordTokens),
		lexer.WithCustomTokenizers(customTokenizers),
		lexer.WithOperators(OperatorTokens),
		lexer.WithSymbols(SymbolTokens),
		lexer.WithCommentMap(comments),
		lexer.WithSpecializationCreators(labelTokenCreator, integerVariableTokenCreator, stringVariableTokenCreator),
		lexer.WithExtendendedIdentifierRunes("_", ":$"), // Allow underscores in identifiers, but when parsing an identifier, stop at a colon (Enables things like Labels)
	)
	return lexer.NewLexer(ll)
}

func integerVariableTokenCreator(identifier string) *lexer.Token {
	if validStringVariableName(identifier) || validLabelName(identifier) {
		return nil
	}
	return lexer.NewToken(IntegerVariableToken, identifier, nil)
}

func stringVariableTokenCreator(identifier string) *lexer.Token {
	if !validStringVariableName(identifier) || validLabelName(identifier) {
		return nil
	}
	return lexer.NewToken(StringVariableToken, identifier, nil)
}

func labelTokenCreator(identifier string) *lexer.Token {
	if !validLabelName(identifier) {
		return nil
	}
	return lexer.NewToken(LabelToken, identifier, 0)
}

// validStringVariableName checks if an identifier is a valid string variable name
func validStringVariableName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !utils.IsIdentifierChar([]rune(identifier)[0], 0, "_", "") {
		return false
	}
	if !strings.HasSuffix(identifier, "$") {
		return false
	}
	return true
}

// validLabelName checks if an identifier is a valid label name
func validLabelName(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !strings.HasSuffix(identifier, ":") {
		return false
	}
	return true
}
