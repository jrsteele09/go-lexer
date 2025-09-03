package lexer

import "fmt"

// TokenIdentifier is a type used to distinguish between different kinds of tokens.
// It's an integer code that gives the token its "identity".
type TokenIdentifier int16

const (
	// NullType represents a null or undefined token type.
	NullType TokenIdentifier = iota

	// EOFType represents the end-of-file token type.
	EOFType

	// EndOfLineType represents the end-of-line token type.
	EndOfLineType

	// IntegerLiteral represents an integer literal token type.
	IntegerLiteral

	// NumberLiteral represents a floating-point number literal token type.
	NumberLiteral

	// HexLiteral represents a hexadecimal number literal token type
	HexLiteral

	// StringLiteral represents a string literal token type.
	StringLiteral

	// LastStdLiteral serves as a marker for the last standard literal token type.
	// Any custom token types should be declared after this constant.
	LastStdLiteral
)

// Token represents a single lexical token in the language being parsed.
// Each token has an identifier, a literal string representation, and an optional value.
// The SourceLine and SourceColumn fields represent the token's position in the source text.
type Token struct {
	ID           TokenIdentifier // The identifier for the type of token.
	Literal      string          // The literal string content of the token.
	Value        interface{}     // The value that the token represents, can be nil.
	Filename     string
	SourceLine   uint // The line in the source text where this token occurs.
	SourceColumn uint // The column in the source text where this token occurs.
}

// String returns a string representation of a Token instance.
// The representation includes the token's identifier, its literal string, and its value.
func (t Token) String() string {
	return fmt.Sprintf("%d: %s: %v", int(t.ID), t.Literal, t.Value)
}

// NewToken is a constructor function for creating a new Token.
// It takes a TokenIdentifier to specify the type, a string for the literal representation,
// and an optional value that the token represents.
func NewToken(id TokenIdentifier, literal string, value interface{}) *Token {
	return &Token{
		ID:      id,
		Literal: literal,
		Value:   value,
	}
}
