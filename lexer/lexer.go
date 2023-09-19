package lexer

import "github.com/pkg/errors"

const (
	NullType = iota
	EOFType
	EndOfLineType
	IntegerLiteral
	NumberLiteral
	StringLiteral
	LastStdLiteral
)

var (
	newLine rune = '\n'
)

// Lexer performs lexical analysis on a stream of input.
type Lexer struct {
	overflowRune      *rune
	currentCommentEnd *string
	parsingString     bool
	firstLineToken    bool
	language          *LanguageConfig
}

// NewLexer initializes a new Lexer with the given language configuration.
func NewLexer(language *LanguageConfig) *Lexer {
	return &Lexer{
		language: language,
	}
}

// TokenizeLine tokenizes a single line of input and returns an array of tokens.
func (l *Lexer) TokenizeLine(line string, lineNo uint) ([]Token, error) {
	var tokens []Token

	addNewToken := func(column int, token *Token) {
		l.firstLineToken = false
		token.SourceLine = lineNo
		token.SourceColumn = uint(column)
		tokens = append(tokens, *token)
	}

	tokenFactory := NewTokenFactory(l)

	skipNextRune := false
	l.firstLineToken = true

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
		}

		runePtr := &r

		for {
			if token, err := tokenFactory.Tokenizer(*runePtr); err != nil {
				return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
			} else if token != nil {
				addNewToken(i, token)
			}
			if l.hasRuneOverflow() {
				runePtr = l.overflowRune
				l.overflowRune = nil
				continue
			}
			break
		}
	}

	// Process the token at the end of the line
	if token, err := tokenFactory.Tokenizer(newLine); err != nil {
		return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
	} else if token != nil {
		addNewToken(len(line), token)
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
