# go-lexer

[![Go Report Card](https://goreportcard.com/badge/github.com/jrsteele09/go-lexer)](https://goreportcard.com/report/github.com/jrsteele09/go-lexer)
[![GoDoc](https://pkg.go.dev/badge/github.com/jrsteele09/go-lexer)](https://pkg.go.dev/github.com/jrsteele09/go-lexer)

`go-lexer` is a simple lexical analysis tool written in Go. It tokenizes lines of text based on the language configuration provided, including comments, string literals, number literals.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Language Configuration](#language-configuration)
  - [Tokenizing a Line](#tokenizing-a-line)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the package, run:

```bash
go get -u github.com/jrsteele09/go-lexer
```

## Usage

### Language Configuration

Before tokenizing, you'll need to create a language configuration. You can use various options to customize the lexer's behavior.

```go
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
	PercentageToken         lexer.TokenIdentifier = 0x111
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

// SymbolTokens defines single delimiter runes to token mappings
var SymbolTokens = map[rune]lexer.TokenIdentifier{
	'+':  AddSymbolToken,
	'/':  DivideSymbolToken,
	'-':  MinusSymbolToken,
	'*':  MultiplySymbolToken,
	'=':  EqualsSymbolToken,
	'<':  LessThanToken,
	'>':  GreaterThanToken,
	'(':  LeftParenthesis,
	')':  RightParenthesis,
	'{':  LeftCurlyBracket,
	'}':  RightCurlyBracket,
	'[':  LeftSquareBracket,
	']':  RightSquareBracket,
	',':  CommaToken,
	';':  SemicolonToken,
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
var comments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
	";":   "\n",
}

// prefixTokenizers - on detection of the starting string, jump to a specific tokenizer.
// Prefix tokenizers override symbol tokens.
var prefixTokenizers = map[string]lexer.TokenizerFunc{
	"$":  lexer.HexTokenizer,
	"0x": lexer.HexTokenizer,
	"%0": lexer.BinaryTokenizer,
	"%1": lexer.BinaryTokenizer,
}

// NewBasicLexer constructs a new Lexer using predefined language settings
func NewBasicLexer() *lexer.Lexer {
	ll := lexer.NewLexerLanguage(lexer.LanguageConfig{
		Keywords:                KeywordTokens,
		Operators:               OperatorTokens,
		Symbols:                 SymbolTokens,
		Comments:                comments,
		PrefixTokenizers:        prefixTokenizers,
		ExtendedIdentifierRunes: "_",
		IdentifierTermination:   ":$",
		TokenCreators:           []func(string) lexer.Token{labelTokenCreator, integerVariableTokenCreator, stringVariableTokenCreator},
	})
	return lexer.NewLexer(ll)
}

func integerVariableTokenCreator(identifier string) lexer.Token {
	if validStringVariableName(identifier) || validLabelName(identifier) {
		return lexer.Token{}
	}
	return lexer.NewToken(IntegerVariableToken, identifier, nil)
}

func stringVariableTokenCreator(identifier string) lexer.Token {
	if !validStringVariableName(identifier) || validLabelName(identifier) {
		return lexer.Token{}
	}
	return lexer.NewToken(StringVariableToken, identifier, nil)
}

func labelTokenCreator(identifier string) lexer.Token {
	if !validLabelName(identifier) {
		return lexer.Token{}
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


```

### Tokenizing a Line

After creating a language configuration, you can create a new Lexer instance and tokenize lines.

```go
// Create new Lexer instance
l := lexer.NewLexer(languageConfig)

// Tokenize a single line (filename is used in error messages)
tokens, err := l.TokenizeLine("let x = 42", "myfile.bas", 1)
if err != nil {
    // Handle error
}

// Tokenize a full io.Reader (e.g. an open file)
tokens, err = l.Tokenize(reader, "myfile.bas")
if err != nil {
    // Handle error
}
```

## Contributing

We appreciate any contributions to improve `go-lexer`. Please feel free to file issues or submit pull requests.

## License

MIT License. See [LICENSE](LICENSE.md) for more information.
