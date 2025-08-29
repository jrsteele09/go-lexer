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
	currentTokenizer TokenizerHandler
	commentParser    *comments.CommentParser
	languageConfig   *LanguageConfig
	parsingString    bool
}

// NewTokenCreator initializes and returns a new TokenFactory for a given lexer.
func NewTokenCreator(commentParser *comments.CommentParser, lc *LanguageConfig) *TokenCreator {
	tf := &TokenCreator{commentParser: commentParser, languageConfig: lc}
	tf.SetTokenizer(tf.tokenizerSelector())
	return tf
}

// Tokenize calls the current tokenizer, defaulting to the tokenizer identifier function.
// Once a token has been created, it restores to identifying the type of the next token.
func (tf *TokenCreator) Tokenize(r rune) ([]*Token, error) {
	tokens, completed, err := tf.currentTokenizer(r)
	if err != nil {
		return nil, err
	}
	if completed {
		tf.SetTokenizer(tf.tokenizerSelector())
	}
	return tokens, err
}

// tokenizerSelector is the default tokenization function.
// It identifies tokens based on individual runes.
func (tf *TokenCreator) tokenizerSelector() TokenizerHandler {
	return func(r rune) ([]*Token, completed, error) {
		// Ignore the rune if in a comment or a space
		if tf.commentParser.InComment() || unicode.IsSpace(r) {
			return nil, false, nil
		}

		if tf.languageConfig.IsCustomTokenizer(string(r)) {
			tf.SetTokenizer(tf.languageConfig.Tokenizer(string(r))(tf, string(r)))
			return nil, false, nil
		} else if _, found := tf.languageConfig.symbolTokens[r]; found {
			tf.SetTokenizer(SymbolTokenizer(tf, string(r))) // Replace the defaultTokenizer with the symbolTokenizer
			return nil, false, nil
		} else if unicode.IsDigit(r) {
			tf.SetTokenizer(NumberTokenizer(tf, string(r))) // Replace the defaultTokenizer with the numberTokenizer
			return nil, false, nil

		} else if utils.IsStringQuotes(r) {
			tf.SetTokenizer(StringTokenizer(tf, string(r))) // Replace the defaultTokenizer with the stringTokenizer
			return nil, false, nil

		} else if utils.IsIdentifierChar(r, 0, tf.languageConfig.extendedIdentifierRunes, tf.languageConfig.identifierTermination) {
			tf.SetTokenizer(IdentifierTokenizer(tf, string(r))) // Replace the defaultTokenizer with the identifierTokenizer
			return nil, false, nil
		} else if tf.commentParser.IsStartOfComment(string(r)) { // Check for being in a comment - could be assembly ";"
			return nil, false, nil
		}

		return nil, false, errors.New("unknown character: \"" + string(r) + "\"")
	}
}

func (tf *TokenCreator) SetTokenizer(tokenizer TokenizerHandler) {
	tf.currentTokenizer = tokenizer
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
