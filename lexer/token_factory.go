package lexer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// TokenFactory manages the creation of tokens for a given lexer.
type TokenFactory struct {
	Tokenizer func(r rune) (*Token, error)
	lexer     *Lexer
}

// NewTokenFactory initializes and returns a new TokenFactory for a given lexer.
func NewTokenFactory(lexer *Lexer) *TokenFactory {
	tf := &TokenFactory{lexer: lexer}
	tf.Tokenizer = tf.defaultTokenizer
	return tf
}

// defaultTokenizer is the default tokenization function. It identifies tokens based on individual runes.
func (tf *TokenFactory) defaultTokenizer(r rune) (*Token, error) {
	// Ignore the rune if in a comment or a space
	if tf.lexer.inComment() || unicode.IsSpace(r) {
		return nil, nil
	}

	// Check if we can create a single rune token
	if tokenID, ok := tf.lexer.language.singleRuneTokens[r]; ok {
		return NewToken(tokenID, string(r), nil), nil
	}

	// Check if we can use a numberTokenizer
	if IsDigit(r) {
		tf.Tokenizer = tf.numberTokenizer(r) // Replace the defaultTokenizer with the numberTokenizer
		return nil, nil

	} else if tf.IsStringQuotes(r) {
		tf.Tokenizer = tf.stringTokenizer(r) // Replace the defaultTokenizer with the stringTokenizer
		return nil, nil

	} else if IsIdentifierChar(r, 0) {
		tf.Tokenizer = tf.identifierTokenizer(r) // Replace the defaultTokenizer with the identifierTokenizer
		return nil, nil
	}

	return nil, errors.New("unknown character: \"" + string(r) + "\"")
}

// numberTokenizer processes numeric literals.
func (tf *TokenFactory) numberTokenizer(initialRune rune) func(r rune) (*Token, error) {
	parsedString := string(initialRune)

	return func(r rune) (*Token, error) {
		if IsDigit(r) {
			parsedString = parsedString + string(r)
		} else {
			tf.lexer.overflowRune = &r
			tf.Tokenizer = tf.defaultTokenizer
			number, err := tf.stringToNumber(parsedString)
			if err != nil {
				return nil, errors.Wrap(err, "numberTokenizer stringToNumber")
			}
			switch number.(type) {
			case float64:
				return NewToken(NumberLiteral, parsedString, number), nil
			case int64:
				return NewToken(IntegerLiteral, parsedString, number), nil
			}
			return NewToken(NumberLiteral, parsedString, number), nil
		}
		return nil, nil
	}
}

// stringTokenizer processes string literals.
func (tf *TokenFactory) stringTokenizer(initialRune rune) func(rune) (*Token, error) {
	startRune := string(initialRune)
	openingString := true
	var parsedString string
	tf.lexer.parsingString = true

	return func(runeChar rune) (*Token, error) {
		if openingString {
			openingString = false
			return nil, nil
		}
		if string(runeChar) == startRune {
			tf.lexer.parsingString = false
			tf.Tokenizer = tf.defaultTokenizer
			return NewToken(StringLiteral, "", parsedString), nil
		}

		parsedString += string(runeChar)
		return nil, nil
	}
}

// identifierTokenizer processes identifiers like variable names.
func (tf *TokenFactory) identifierTokenizer(initialRune rune) func(rune) (*Token, error) {
	parsedString := string(initialRune)

	return func(runeChar rune) (*Token, error) {
		if tf.lexer.language.labelTerminator != nil && tf.lexer.firstLineToken && runeChar == *tf.lexer.language.labelTerminator {
			return NewToken(tf.lexer.language.labelToken, parsedString+string(runeChar), ""), nil
		}
		if !IsIdentifierChar(runeChar, len(parsedString)) {
			tf.lexer.overflowRune = &runeChar
			tf.Tokenizer = tf.defaultTokenizer
			return tf.lexer.language.tokenFromIdentifier(parsedString), nil
		}
		parsedString += string(runeChar)
		return nil, nil
	}
}

// stringToNumber converts a string to a numerical value.
func (tf *TokenFactory) stringToNumber(strNum string) (interface{}, error) {
	if strings.Contains(strNum, ".") {
		return strconv.ParseFloat(strNum, 64)
	}
	return strconv.ParseInt(strNum, 10, 64)
}

// IsIdentifierChar checks if a rune is valid in an identifier.
func IsIdentifierChar(runeChar rune, pos int) bool {
	if unicode.IsLetter(runeChar) {
		return true
	}
	if pos > 0 && IsDigit(runeChar) {
		return true
	}

	return strings.Contains("_#%", string(runeChar))
}

// IsDigit checks if a rune is a digit.
func IsDigit(runeChar rune) bool {
	const extraDigits = "."
	return unicode.IsDigit(runeChar) || strings.Contains(extraDigits, string(runeChar))
}

// IsStringQuotes checks if a rune can start or end a string.
func (tf *TokenFactory) IsStringQuotes(runeChar rune) bool {
	const startOfString = "\"'`"
	return strings.Contains(startOfString, string(runeChar))
}
