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
go get -u github.com/yourusername/go-lexer
```

## Usage

### Language Configuration

Before tokenizing, you'll need to create a language configuration. You can use various options to customize the lexer's behavior.

```go
import "github.com/jsteele09/go-lexer"

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

var KeywordTokens = map[string]lexer.TokenIdentifier{
	"if":   IfStatementToken,
	"let":  LetStatementToken,
	"for":  ForStatementToken,
	"to":   ToStatementToken,
	"next": NextStatementToken,
}

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

var comments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
}

// Create new language configuration
ll := lexer.NewLexerLanguage(
    lexer.WithSingleRuneMap(SingleCharTokens),
    lexer.WithCommentMap(comments),
    lexer.WithTokenCreators(identifierToken),
    lexer.WithLabelSettings(':', LabelToken),
)

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
	if lexer.IsIdentifierChar([]rune(identifier)[0], 0) {
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

// Tokenize a line
tokens, err := l.TokenizeLine("let x = 42", 1)
if err != nil {
    // Handle error
}
```

## Contributing

We appreciate any contributions to improve `go-lexer`. Please feel free to file issues or submit pull requests.

## License

MIT License. See [LICENSE](LICENSE.md) for more information.
