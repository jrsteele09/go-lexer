package lexer

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-lexer/lexer/utils"
	"github.com/pkg/errors"
)

// NumberTokenizer processes numeric literals.
func NumberTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	parsedNumber := initialString

	return func(r rune) ([]Token, completed, error) {
		if utils.IsDigit(r, len(parsedNumber)) {
			parsedNumber = parsedNumber + string(r)
		} else if tf.languageConfig.IsCustomTokenizer(parsedNumber + string(r)) { // Not a digit, perhaps a custom tokenizer, i.e. hex: "0xFF"
			parsedNumber += string(r)
			tf.SetTokenizer(tf.languageConfig.PrefixTokenizers[parsedNumber](tf, parsedNumber))
			return nil, false, nil
		} else {
			tf.SetOverFlow(r)
			number, err := utils.StringToNumber(parsedNumber)
			if err != nil {
				return nil, false, errors.Wrap(err, "numberTokenizer stringToNumber")
			}
			switch number.(type) {
			case float64:
				return []Token{NewToken(NumberLiteral, parsedNumber, number)}, true, nil
			case int64:
				return []Token{NewToken(IntegerLiteral, parsedNumber, number)}, true, nil
			}
			return []Token{NewToken(NumberLiteral, parsedNumber, number)}, true, nil
		}
		return nil, false, nil
	}
}

// BinaryTokenizer processes a binary number
func BinaryTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	var builder strings.Builder

	if initialString == "0" || initialString == "1" {
		builder.WriteString(initialString)
	}

	return func(r rune) ([]Token, completed, error) {
		if utils.IsBinaryDigit(r) {
			builder.WriteRune(r)
			return nil, false, nil
		}

		current := builder.String()

		// Not a digit, perhaps a custom tokenizer, e.g. "0x"
		if tf.languageConfig.IsCustomTokenizer(current + string(r)) {
			builder.WriteRune(r)
			tokenizerPrefix := builder.String()
			tf.SetTokenizer(tf.languageConfig.PrefixTokenizers[tokenizerPrefix](tf, tokenizerPrefix))
			return nil, false, nil
		}

		tf.SetOverFlow(r)

		number, err := utils.BinaryStringToNumber(current)
		if err != nil {
			return nil, true, fmt.Errorf("BinaryTokenizer BinaryStringToNumber [%w]", err)
		}

		return []Token{
			NewToken(IntegerLiteral, current, number),
		}, true, nil
	}
}

// HexTokenizer processes hex literals.
func HexTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	var builder strings.Builder

	if initialString == "" || initialString == "$" {
		builder.WriteString("0x")
	}

	return func(r rune) ([]Token, completed, error) {
		if utils.IsHexDigit(r) {
			builder.WriteRune(r)
			return nil, false, nil
		}

		tf.SetOverFlow(r)

		parsedString := builder.String()

		number, err := utils.HexToNumber(parsedString)
		if err != nil {
			return nil, false, errors.Wrap(err, "HexTokenizer HexToNumber")
		}

		return []Token{
			NewToken(HexLiteral, parsedString, number),
		}, true, nil
	}
}

// StringTokenizer processes string literals.
func StringTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	startRune := initialString

	var builder strings.Builder

	return func(r rune) ([]Token, completed, error) {
		if string(r) == startRune {
			return []Token{
				NewToken(StringLiteral, "", builder.String()),
			}, true, nil
		}

		builder.WriteRune(r)
		return nil, false, nil
	}
}

// IdentifierTokenizer processes identifiers like variable names.
func IdentifierTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	var builder strings.Builder
	builder.WriteString(initialString)

	return func(r rune) ([]Token, completed, error) {
		if !utils.IsIdentifierChar(
			r,
			builder.Len(),
			tf.languageConfig.ExtendedIdentifierRunes,
			tf.languageConfig.IdentifierTermination,
		) {
			tf.SetOverFlow(r)

			identifier := builder.String()

			if tf.commentParser.IsStartOfComment(identifier) {
				return nil, false, nil
			}

			t := tf.languageConfig.tokenFromIdentifier(identifier)
			if t.ID == NullType {
				return nil, true, fmt.Errorf("unknown identifier %s", identifier)
			}

			return []Token{t}, true, nil
		}

		builder.WriteRune(r)

		rStr := string(r)
		if strings.Contains(tf.languageConfig.IdentifierTermination, rStr) {
			identifier := builder.String()

			t := tf.languageConfig.tokenFromIdentifier(identifier)
			if t.ID == NullType {
				return nil, true, fmt.Errorf("unknown identifier %s", identifier)
			}

			return []Token{t}, true, nil
		}

		return nil, false, nil
	}
}

// SymbolTokenizer processes operators
func SymbolTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	symbolsString := string(initialString)

	createToken := func(overflowRune rune) ([]Token, completed, error) {
		// First check that the parsed symbolString + overflowRune doesn't match a custom tokenizer
		resetTokenizer := completed(true)
		if len(symbolsString) > 0 {
			runes := []rune(symbolsString)
			lastRune := runes[len(runes)-1]
			tokenizerStr := string(lastRune) + string(overflowRune) //
			if tf.languageConfig.IsCustomTokenizer(tokenizerStr) {  // CustomTokenizers take priority over symbols
				tf.SetTokenizer(tf.languageConfig.Tokenizer(tokenizerStr)(tf, string(lastRune)))
				symbolsString = string(runes[:len(runes)-1])
				resetTokenizer = false // Don't want to reset as we've just swapped the tokenizer
			} else if tf.languageConfig.IsCustomTokenizer(symbolsString) {
				tf.SetTokenizer(tf.languageConfig.Tokenizer(symbolsString)(tf, symbolsString))
				tf.SetOverFlow(overflowRune)
				return nil, false, nil
			}
		}

		tf.SetOverFlow(overflowRune)
		symbolTokens := make([]Token, 0)

		// Parses potentially larger Operator tokens such as "<=" and single symbols such as
		// "+", "-" ...
		i := 0
		for i < len(symbolsString) {
			var longestSymbol string
			for x := i + 1; x < len(symbolsString); x++ {
				symbolStr := symbolsString[i : x+1]
				if _, found := tf.languageConfig.Operators[symbolStr]; found {
					longestSymbol = symbolStr
				} else if tf.commentParser.IsStartOfComment(symbolStr) {
					return nil, false, nil
				}
			}
			if longestSymbol != "" {
				tokenID := tf.languageConfig.Operators[longestSymbol]
				symbolTokens = append(symbolTokens, NewToken(tokenID, longestSymbol, longestSymbol))
				i += len(longestSymbol)
			} else {
				if tf.commentParser.IsStartOfComment(symbolsString) {
					return nil, false, nil
				}
				tokenID, found := tf.languageConfig.Symbols[rune(symbolsString[i])]
				if !found {
					return nil, false, fmt.Errorf("unknown symbol %s", string(symbolsString[i]))
				}
				symbolTokens = append(symbolTokens, NewToken(tokenID, string(symbolsString[i]), symbolsString[i]))
				i++
			}
		}

		return symbolTokens, resetTokenizer, nil
	}

	return func(r rune) ([]Token, completed, error) {
		if tf.languageConfig.IsCustomTokenizer(string(r)) { // CustomTokenizers take priority over symbols
			return createToken(r)
		} else if tf.commentParser.IsStartOfComment(symbolsString) { // Check for being in a comment - could be assembly ";"
			return nil, true, nil
		} else if _, found := tf.languageConfig.Symbols[r]; found {
			symbolsString += string(r)
		} else {
			return createToken(r)
		}

		return nil, false, nil
	}
}
