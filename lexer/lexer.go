package lexer

import (
	"bufio"
	"io"

	"github.com/jrsteele09/go-lexer/lexer/comments"
	"github.com/pkg/errors"
)

const (
	newLine = '\n'
)

// Lexer performs lexical analysis on a stream of input.
type Lexer struct {
	language      *LanguageConfig
	commentParser *comments.CommentParser
}

// NewLexer initializes a new Lexer with the given language configuration.
func NewLexer(language *LanguageConfig) *Lexer {
	return &Lexer{
		language:      language,
		commentParser: comments.NewCommentParser(language.comments),
	}
}

// Tokenize reads from an io.Reader line by line, tokenizes each line using TokenizeLine,
// and returns all the generated tokens.
func (l *Lexer) Tokenize(r io.Reader) ([]*Token, error) {
	var allTokens []*Token
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

	allTokens = append(allTokens, NewToken(EOFType, "", nil))
	return allTokens, nil
}

// TokenizeLine tokenizes a single line of input and returns an array of tokens.
func (l *Lexer) TokenizeLine(line string, lineNo uint) ([]*Token, error) {
	var lineTokens []*Token

	addNewTokens := func(column int, tokens []*Token) {
		if tokens == nil {
			return
		}
		for _, token := range tokens {
			token.SourceLine = lineNo
			token.SourceColumn = uint(column)
			lineTokens = append(lineTokens, token)
		}
	}

	addEndOfLine := func() {
		if len(lineTokens) != 0 {
			addNewTokens(len(line), []*Token{NewToken(EndOfLineType, string(newLine), nil)})
		}
	}

	tokenFactory := NewTokenCreator(l.commentParser, l.language)

	for i, r := range line {
		if l.commentParser.InComment() {
			if l.commentParser.IsNewLineComment() {
				break
			}
			l.commentParser.ParseEndOfComment(r)
			continue
		}

		for {
			if tokens, err := tokenFactory.Tokenize(r); err != nil {
				return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
			} else if len(tokens) > 0 {
				addNewTokens(i, tokens)
			}
			if !tokenFactory.HasRuneOverflow() {
				break
			}
			r = tokenFactory.OverflowRune()
		}
	}

	if l.commentParser.InComment() {
		if l.commentParser.IsNewLineComment() {
			l.commentParser.Reset()
		}
		addEndOfLine()
		return lineTokens, nil
	}

	// Need to complete the tokenization process for the last rune,
	// It could be that a tokenizer was in progress when a newline was reached
	token, err := tokenFactory.Tokenize(newLine)
	if err != nil {
		return nil, errors.Wrap(err, "Lexer.TokenizeLine.tokenFactory.Tokenizer")
	}

	addNewTokens(len(line), token)
	addEndOfLine()
	return lineTokens, nil
}
