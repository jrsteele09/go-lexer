package lexer

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-lexer/lexer/utils"
	"github.com/pkg/errors"
)

// NumberTokenizer processes numeric literals.
func NumberTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	parsedString := initialString

	return func(r rune) ([]*Token, completed, error) {
		if utils.IsDigit(r) {
			parsedString = parsedString + string(r)
		} else if tf.languageConfig.IsCustomTokenizer(parsedString + string(r)) { // Not a digit, perhaps a custom tokenizer, i.e. hex: "0xFF"
			parsedString += string(r)
			tf.SetTokenizer(tf.languageConfig.customTokenizers[parsedString](tf, parsedString))
			return nil, false, nil
		} else {
			tf.SetOverFlow(r)
			number, err := utils.StringToNumber(parsedString)
			if err != nil {
				return nil, false, errors.Wrap(err, "numberTokenizer stringToNumber")
			}
			switch number.(type) {
			case float64:
				return []*Token{NewToken(NumberLiteral, parsedString, number)}, true, nil
			case int64:
				return []*Token{NewToken(IntegerLiteral, parsedString, number)}, true, nil
			}
			return []*Token{NewToken(NumberLiteral, parsedString, number)}, true, nil
		}
		return nil, false, nil
	}
}

// BinaryTokenizer processes a binary number
func BinaryTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	parsedString := ""
	if initialString == "0" || initialString == "1" {
		parsedString = initialString
	}

	return func(r rune) ([]*Token, completed, error) {
		if utils.IsBinaryDigit(r) {
			parsedString = parsedString + string(r)
		} else if tf.languageConfig.IsCustomTokenizer(parsedString + string(r)) { // Not a digit, perhaps a custom tokenizer, i.e. hex: "0xFF"
			parsedString += string(r)
			tf.SetTokenizer(tf.languageConfig.customTokenizers[parsedString](tf, parsedString))
			return nil, false, nil
		} else {
			tf.SetOverFlow(r)
			number, err := utils.BinaryStringToNumber(parsedString)
			if err != nil {
				return nil, true, fmt.Errorf("BinaryTokenizer BinaryStringToNumber [%w]", err)
			}
			return []*Token{NewToken(IntegerLiteral, parsedString, number)}, true, nil
		}
		return nil, false, nil
	}
}

// HexTokenizer processes hex literals.
func HexTokenizer(tf *TokenCreator, initalString string) TokenizerHandler {
	var parsedString string
	if initalString == "" || initalString == "$" {
		parsedString = "0x"
	}

	return func(r rune) ([]*Token, completed, error) {
		if utils.IsHexDigit(r) {
			parsedString = parsedString + string(r)
		} else {
			tf.SetOverFlow(r)

			number, err := utils.HexToNumber(parsedString)
			if err != nil {
				return nil, false, errors.Wrap(err, "hexTokenizer stringToNumber")
			}
			return []*Token{NewToken(HexLiteral, parsedString, number)}, true, nil
		}
		return nil, false, nil
	}
}

// StringTokenizer processes string literals.
func StringTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	startRune := initialString
	var parsedString string
	tf.parsingString = true

	return func(runeChar rune) ([]*Token, completed, error) {
		if string(runeChar) == startRune {
			tf.parsingString = false
			return []*Token{NewToken(StringLiteral, "", parsedString)}, true, nil
		}

		parsedString += string(runeChar)
		return nil, false, nil
	}
}

// IdentifierTokenizer processes identifiers like variable names.
func IdentifierTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	parsedString := initialString

	return func(runeChar rune) ([]*Token, completed, error) {
		if !utils.IsIdentifierChar(runeChar, len(parsedString), tf.languageConfig.extendedIdentifierRunes, tf.languageConfig.identifierTermination) {
			tf.SetOverFlow(runeChar)

			if tf.commentParser.IsStartOfComment(parsedString) {
				return nil, false, nil
			}
			return []*Token{tf.languageConfig.tokenFromIdentifier(parsedString)}, true, nil
		}
		parsedString += string(runeChar)

		if strings.Contains(tf.languageConfig.identifierTermination, string(runeChar)) {
			return []*Token{tf.languageConfig.tokenFromIdentifier(parsedString)}, true, nil
		}
		return nil, false, nil
	}
}

// SymbolTokenizer processes operators
func SymbolTokenizer(tf *TokenCreator, initialString string) TokenizerHandler {
	symbolsString := string(initialString)

	createToken := func(overflowRune rune) ([]*Token, completed, error) {
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
		symbolTokens := make([]*Token, 0)

		// Parses potentially larger Operator tokens such as "<=" and single symbols such as
		// "+", "-" ...
		i := 0
		for i < len(symbolsString) {
			var longestSymbol string
			for x := i + 1; x < len(symbolsString); x++ {
				symbolStr := symbolsString[i : x+1]
				if _, found := tf.languageConfig.operatorTokens[symbolStr]; found {
					longestSymbol = symbolStr
				} else if tf.commentParser.IsStartOfComment(symbolStr) {
					return nil, false, nil
				}
			}
			if longestSymbol != "" {
				tokenID := tf.languageConfig.operatorTokens[longestSymbol]
				symbolTokens = append(symbolTokens, NewToken(tokenID, longestSymbol, longestSymbol))
				i += len(longestSymbol)
			} else {
				if tf.commentParser.IsStartOfComment(symbolsString) {
					return nil, false, nil
				}
				tokenID, found := tf.languageConfig.symbolTokens[rune(symbolsString[i])]
				if !found {
					return nil, false, fmt.Errorf("unknown symbol %s", string(symbolsString[i]))
				}
				symbolTokens = append(symbolTokens, NewToken(tokenID, string(symbolsString[i]), symbolsString[i]))
				i++
			}
		}

		return symbolTokens, resetTokenizer, nil
	}

	return func(r rune) ([]*Token, completed, error) {
		if tf.languageConfig.IsCustomTokenizer(string(r)) { // CustomTokenizers take priority over symbols
			return createToken(r)
		} else if _, found := tf.languageConfig.symbolTokens[r]; found {
			symbolsString += string(r)
		} else {
			return createToken(r)
		}

		return nil, false, nil
	}
}
