package lexer

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
)

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

var (
	newLine = '\n'
)

// Lexer performs lexical analysis on a stream of input.
type Lexer struct {
	overflowRune      *rune
	currentCommentEnd *string
	parsingString     bool
	language          *LanguageConfig
}

// NewLexer initializes a new Lexer with the given language configuration.
func NewLexer(language *LanguageConfig) *Lexer {
	return &Lexer{
		language: language,
	}
}

// Tokenize reads from an io.Reader line by line, tokenizes each line using TokenizeLine,
// and returns all the generated tokens.
func (l *Lexer) Tokenize(r io.Reader) ([]Token, error) {
	var allTokens []Token
	scanner := bufio.NewScanner(r)
	lineNo := uint(1)

	for scanner.Scan() {
		line := scanner.Text()
		tokens, err := l.TokenizeLine(line, lineNo)
		if err != nil {
			return nil, err
		}
		allTokens = append(allTokens, tokens...)
		lineNo++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	allTokens = append(allTokens, *NewToken(EOFType, "", nil))
	return allTokens, nil
}

// TokenizeLine tokenizes a single line of input and returns an array of tokens.
func (l *Lexer) TokenizeLine(line string, lineNo uint) ([]Token, error) {
	var tokens []Token

	addNewToken := func(column int, token *Token) {
		token.SourceLine = lineNo
		token.SourceColumn = uint(column)
		tokens = append(tokens, *token)
	}

	tokenFactory := NewTokenFactory(l)

	skipNextRune := false

	for i, r := range line {
		if l.commentOnEndOfLine() {
			break
		}
		if skipNextRune {
			skipNextRune = false
			continue
		}
		if l.endOfComment(i, line) {
			skipNextRune = true
			continue
		}
		if l.inComment() {
			continue
		}
		if l.startOfComment(i, line) {
			skipNextRune = true
			continue
		}

		for {
			if token, err := tokenFactory.Tokenizer(r); err != nil {
				return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
			} else if token != nil {
				addNewToken(i, token)
			}
			if l.hasRuneOverflow() {
				r = *l.overflowRune
				l.overflowRune = nil
				continue
			}
			break
		}
	}

	// Need to complete the tokenization process for the last rune, its likely that a tokenizer was in progress when the end of the line was reached.
	if token, err := tokenFactory.Tokenizer(newLine); err != nil { // Call tokenization with the NewLine rune to force the tokenizer to complete the tokenization process.
		return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
	} else if token != nil {
		addNewToken(len(line), token)
	}

	if len(tokens) != 0 {
		addNewToken(len(line), NewToken(EndOfLineType, string(newLine), nil))
	}
	return tokens, nil
}

// commentOnEndOfLine checks if the comment ends at the end of a line.
func (l *Lexer) commentOnEndOfLine() bool {
	if l.currentCommentEnd == nil {
		return false
	}
	if *l.currentCommentEnd != string(newLine) {
		return false
	}
	l.currentCommentEnd = nil
	return true
}

// startOfComment identifies the start of a comment in a given line.
func (l *Lexer) startOfComment(idx int, line string) bool {
	if l.parsingString {
		return false
	}
	if l.currentCommentEnd != nil {
		return false
	}
	if idx >= len(line)-1 {
		return false
	}
	c := line[idx]
	n := line[idx+1]

	currentAndNextRune := string(string(c) + string(n))

	if comment, ok := l.language.comments[currentAndNextRune]; ok {
		l.currentCommentEnd = &comment
		return true
	}

	return false
}

// endOfComment identifies the end of a comment in a given line.
func (l *Lexer) endOfComment(idx int, line string) bool {
	if l.currentCommentEnd == nil {
		return false
	}
	if idx >= len(line)-1 {
		return false
	}
	c := line[idx]
	n := line[idx+1]
	currentAndNextRune := string(string(c) + string(n))

	if currentAndNextRune != *l.currentCommentEnd {
		return false
	}

	l.currentCommentEnd = nil
	return true
}

// inComment checks if the Lexer is currently inside a comment.
func (l *Lexer) inComment() bool {
	return l.currentCommentEnd != nil
}

// hasRuneOverflow checks if there's a pending rune to process.
func (l *Lexer) hasRuneOverflow() bool {
	return !l.inComment() && l.overflowRune != nil
}
