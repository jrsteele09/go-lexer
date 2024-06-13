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

	if tf.IsSymbol(r) {
		tf.Tokenizer = tf.symbolTokenizer(r) // Replace the defaultTokenizer with the operatorTokenizer
		return nil, nil

	} else if IsDigit(r) {
		tf.Tokenizer = tf.numberTokenizer(r) // Replace the defaultTokenizer with the numberTokenizer
		return nil, nil

	} else if IsStringQuotes(r) {
		tf.Tokenizer = tf.stringTokenizer(r) // Replace the defaultTokenizer with the stringTokenizer
		return nil, nil

	} else if IsIdentifierChar(r, 0, tf.lexer.language.extraIdentifierRunes) {
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
		} else if r == 'x' && parsedString == "0" { // Swap over to a Hex Tokenizer
			tf.Tokenizer = tf.hexTokenizer()
			return nil, nil
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

// numberTokenizer processes numeric literals.
func (tf *TokenFactory) hexTokenizer() func(r rune) (*Token, error) {
	parsedString := "0x"

	return func(r rune) (*Token, error) {
		if IsHexDigit(r) {
			parsedString = parsedString + string(r)
		} else {
			tf.lexer.overflowRune = &r
			tf.Tokenizer = tf.defaultTokenizer

			number, err := tf.stringToHex(parsedString)
			if err != nil {
				return nil, errors.Wrap(err, "hexTokenizer stringToNumber")
			}
			return NewToken(HexLiteral, parsedString, number), nil
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
		if !IsIdentifierChar(runeChar, len(parsedString), tf.lexer.language.extraIdentifierRunes) {
			tf.lexer.overflowRune = &runeChar
			tf.Tokenizer = tf.defaultTokenizer
			return tf.lexer.language.tokenFromIdentifier(parsedString), nil
		}
		parsedString += string(runeChar)
		return nil, nil
	}
}

// symbolTokenizer processes operators
func (tf *TokenFactory) symbolTokenizer(initialRune rune) func(r rune) (*Token, error) {
	parsedString := string(initialRune)

	createSymbolToken := func(r rune) (*Token, error) {
		tf.lexer.overflowRune = &r
		tf.Tokenizer = tf.defaultTokenizer

		// Parse a single symbol character
		if len(parsedString) == 1 {
			if tokenID, found := tf.lexer.language.symbolTokens[rune(parsedString[0])]; found {
				return NewToken(tokenID, parsedString, parsedString), nil
			}
		} else if len(parsedString) > 1 {
			if tokenID, found := tf.lexer.language.operatorTokens[parsedString]; found {
				return NewToken(tokenID, parsedString, parsedString), nil
			}
		}
		return nil, errors.New("unknown symbol: \"" + parsedString + "\"")
	}

	return func(r rune) (*Token, error) {
		if tf.IsSymbol(r) {
			operator := parsedString + string(r)
			if _, found := tf.lexer.language.operatorTokens[operator]; found { // Check if the operator is a valid token
				parsedString = operator
			} else {
				return createSymbolToken(r)
			}
		} else {
			return createSymbolToken(r)
		}

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

func (tf *TokenFactory) stringToHex(hexString string) (interface{}, error) {
	// Remove the "0x" prefix
	if len(hexString) > 2 && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}

	bitSize := 64

	l := len(hexString)
	if l > 0 && l <= 2 {
		bitSize = 8
	} else if l > 2 && l <= 4 {
		bitSize = 16
	} else if l > 4 && l <= 8 {
		bitSize = 32
	} else if l > 8 {
		bitSize = 64
	}

	// Parse the hexadecimal string to int64
	value, err := strconv.ParseInt(hexString, bitSize, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "stringToHex: failed to parse hex string %s", hexString)
	}

	switch bitSize {
	case 8:
		return int8(value), nil
	case 16:
		return int16(value), nil
	case 32:
		return int32(value), nil
	default:
		return value, nil
	}
}

// IsIdentifierChar checks if a rune is valid in an identifier.
func IsIdentifierChar(runeChar rune, pos int, extraRunes string) bool {
	if unicode.IsLetter(runeChar) {
		return true
	}
	if pos > 0 && IsDigit(runeChar) {
		return true
	}

	return strings.Contains(extraRunes, string(runeChar))
}

// IsDigit checks if a rune is a digit.
func IsDigit(runeChar rune) bool {
	const extraDigits = "."
	return unicode.IsDigit(runeChar) || strings.Contains(extraDigits, string(runeChar))
}

// IsHexDigit checks if a rune is a hex digit.
func IsHexDigit(r rune) bool {
	return unicode.IsDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// IsStringQuotes checks if a rune can start or end a string.
func IsStringQuotes(runeChar rune) bool {
	const startOfString = "\"'`"
	return strings.Contains(startOfString, string(runeChar))
}

// IsSymbol checks if a rune is a symbol.
func (tf *TokenFactory) IsSymbol(runeChar rune) bool {
	_, found := tf.lexer.language.symbolTokens[runeChar]
	return found
}
