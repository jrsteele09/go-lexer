package lexer

import (
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenNumber
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLeftParen
	TokenRightParen
)

type GPTToken struct {
	Type  TokenType
	Value string
}

type GPTLexer struct {
	input        string
	pos          int
	currentToken GPTToken
}

func (l *GPTLexer) next() {
	if l.pos >= len(l.input) {
		l.currentToken = GPTToken{Type: TokenEOF}
		return
	}

	ch := rune(l.input[l.pos])
	l.pos++

	switch {
	case unicode.IsDigit(ch):
		start := l.pos - 1
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.pos++
		}
		l.currentToken = GPTToken{Type: TokenNumber, Value: l.input[start:l.pos]}
	case ch == '+':
		l.currentToken = GPTToken{Type: TokenPlus}
	case ch == '-':
		l.currentToken = GPTToken{Type: TokenMinus}
	case ch == '*':
		l.currentToken = GPTToken{Type: TokenMultiply}
	case ch == '/':
		l.currentToken = GPTToken{Type: TokenDivide}
	case ch == '(':
		l.currentToken = GPTToken{Type: TokenLeftParen}
	case ch == ')':
		l.currentToken = GPTToken{Type: TokenRightParen}
	default:
		l.currentToken = GPTToken{Type: TokenError}
	}
}

type Parser struct {
	lexer        *GPTLexer
	currentToken GPTToken
}

func (p *Parser) consume(t TokenType) {
	if p.currentToken.Type == t {
		p.lexer.next()
		p.currentToken = p.lexer.currentToken
	} else {
		panic("Unexpected token")
	}
}

func (p *Parser) parseExpression(precedence int) int {
	left := p.parsePrimary()

	for precedence < p.getOperatorPrecedence() {
		op := p.currentToken.Type
		p.consume(op)

		right := p.parseExpression(p.getOperatorPrecedence())

		switch op {
		case TokenPlus:
			left += right
		case TokenMinus:
			left -= right
		case TokenMultiply:
			left *= right
		case TokenDivide:
			left /= right
		}
	}

	return left
}

func (p *Parser) parsePrimary() int {
	switch p.currentToken.Type {
	case TokenNumber:
		value, _ := strconv.Atoi(p.currentToken.Value)
		p.consume(TokenNumber)
		return value
	case TokenLeftParen:
		p.consume(TokenLeftParen)
		value := p.parseExpression(0)
		p.consume(TokenRightParen)
		return value
	default:
		panic("Invalid primary expression")
	}
}

func (p *Parser) getOperatorPrecedence() int {
	switch p.currentToken.Type {
	case TokenPlus, TokenMinus:
		return 1
	case TokenMultiply, TokenDivide:
		return 2
	default:
		return 0
	}
}

func Parse(input string) int {
	lexer := &GPTLexer{input: strings.TrimSpace(input)}
	lexer.next()
	parser := &Parser{lexer: lexer, currentToken: lexer.currentToken}
	return parser.parseExpression(0)
}
