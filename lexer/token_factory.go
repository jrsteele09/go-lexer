package lexer

import (
	"unicode"

	"github.com/jrsteele09/go-lexer/lexer/comments"
	"github.com/jrsteele09/go-lexer/lexer/utils"
	"github.com/pkg/errors"
)

// TokenCreator manages the creation of tokens for a given lexer.
type TokenCreator struct {
	overflowRune     *rune
	currentTokenizer func(r rune) ([]*Token, error)
	commentParser    *comments.CommentParser
	languageConfig   *LanguageConfig
	parsingString    bool
}

// NewTokenCreator initializes and returns a new TokenFactory for a given lexer.
func NewTokenCreator(commentParser *comments.CommentParser, lc *LanguageConfig) *TokenCreator {
	tf := &TokenCreator{commentParser: commentParser, languageConfig: lc}
	tf.SetTokenizer(tf.tokenCreatorIdentifier())
	return tf
}

// Tokenize calls the current tokenizer, defaulting to the tokenizer identifier function.
// Once a token has been created, it restores to identifying the type of the next token.
func (tf *TokenCreator) Tokenize(r rune) ([]*Token, error) {
	tokens, err := tf.currentTokenizer(r)
	if err != nil {
		return nil, err
	}
	if tokens == nil {
		return nil, nil
	}
	tf.SetTokenizer(tf.tokenCreatorIdentifier())
	return tokens, err
}

// tokenCreatorIdentifier is the default tokenization function.
// It identifies tokens based on individual runes.
func (tf *TokenCreator) tokenCreatorIdentifier() func(r rune) ([]*Token, error) {
	return func(r rune) ([]*Token, error) {
		// Ignore the rune if in a comment or a space
		if tf.commentParser.InComment() || unicode.IsSpace(r) {
			return nil, nil
		}

		if tf.languageConfig.IsCustomTokenizer(r) {
			tf.SetTokenizer(tf.languageConfig.Tokenizer(r)(tf, r))
			return nil, nil
		} else if _, found := tf.languageConfig.symbolTokens[r]; found {
			tf.SetTokenizer(SymbolTokenizer(tf, r)) // Replace the defaultTokenizer with the operatorTokenizer
			return nil, nil
		} else if utils.IsDigit(r) {
			tf.SetTokenizer(NumberTokenizer(tf, r)) // Replace the defaultTokenizer with the numberTokenizer
			return nil, nil

		} else if utils.IsStringQuotes(r) {
			tf.SetTokenizer(StringTokenizer(tf, r)) // Replace the defaultTokenizer with the stringTokenizer
			return nil, nil

		} else if utils.IsIdentifierChar(r, 0, tf.languageConfig.extendedIdentifierRunes, tf.languageConfig.identifierTermination) {
			tf.SetTokenizer(IdentifierTokenizer(tf, r)) // Replace the defaultTokenizer with the identifierTokenizer
			return nil, nil
		}

		return nil, errors.New("unknown character: \"" + string(r) + "\"")
	}
}

func (tf *TokenCreator) SetTokenizer(tokinzer func(r rune) ([]*Token, error)) {
	tf.currentTokenizer = tokinzer
}

func (tf *TokenCreator) SetOverFlow(r rune) {
	tf.overflowRune = &r
}

// HasRuneOverflow checks if there's a pending rune to process.
func (tf *TokenCreator) HasRuneOverflow() bool {
	return !tf.commentParser.InComment() && tf.overflowRune != nil
}

// OverflowRune returns an overflow rune if available, or '\n' if not.
// It should be preceded by a call to HadRuneOverflow
func (tf *TokenCreator) OverflowRune() rune {
	if tf.overflowRune == nil {
		return '\n'
	}
	r := *tf.overflowRune
	tf.overflowRune = nil
	return r
}
