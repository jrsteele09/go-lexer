package lexer

import "fmt"

// TokenIdentifier codes the token
type TokenIdentifier int16

// Token models a language token
type Token struct {
	Id           TokenIdentifier
	Literal      string
	Value        interface{}
	SourceLine   uint
	SourceColumn uint
}

// String returns a string representation of a token
func (t Token) String() string {
	return fmt.Sprintf("%d: %s: %v", int(t.Id), t.Literal, t.Value)
}

// NewToken creates a new instance of a token
func NewToken(id TokenIdentifier, literal string, value interface{}) *Token {
	return &Token{
		Id:      id,
		Literal: literal,
		Value:   value,
	}
}
