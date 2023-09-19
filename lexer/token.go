// Package lexer provides utilities for lexical analysis.
package lexer

import "fmt"

// TokenIdentifier is a type used to distinguish between different kinds of tokens.
// It's an integer code that gives the token its "identity".
type TokenIdentifier int16

// Token represents a single lexical token in the language being parsed.
// Each token has an identifier, a literal string representation, and an optional value.
// The SourceLine and SourceColumn fields represent the token's position in the source text.
type Token struct {
	ID           TokenIdentifier // The identifier for the type of token.
	Literal      string          // The literal string content of the token.
	Value        interface{}     // The value that the token represents, can be nil.
	SourceLine   uint            // The line in the source text where this token occurs.
	SourceColumn uint            // The column in the source text where this token occurs.
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
