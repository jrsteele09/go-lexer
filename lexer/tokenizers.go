package lexer

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-lexer/lexer/utils"
	"github.com/pkg/errors"
)

// NumberTokenizer processes numeric literals.
func NumberTokenizer(tf *TokenCreator, initialRune rune) func(r rune) ([]*Token, error) {
	parsedString := string(initialRune)

	return func(r rune) ([]*Token, error) {
		if utils.IsDigit(r) {
			parsedString = parsedString + string(r)
		} else if r == 'x' && parsedString == "0" { // Swap over to a Hex Tokenizer
			tf.SetTokenizer(HexTokenizer(tf, r))
			return nil, nil
		} else {
			tf.SetOverFlow(r)
			number, err := utils.StringToNumber(parsedString)
			if err != nil {
				return nil, errors.Wrap(err, "numberTokenizer stringToNumber")
			}
			switch number.(type) {
			case float64:
				return []*Token{NewToken(NumberLiteral, parsedString, number)}, nil
			case int64:
				return []*Token{NewToken(IntegerLiteral, parsedString, number)}, nil
			}
			return []*Token{NewToken(NumberLiteral, parsedString, number)}, nil
		}
		return nil, nil
	}
}

// numberTokenizer processes numeric literals.
func HexTokenizer(tf *TokenCreator, _ rune) func(r rune) ([]*Token, error) {
	parsedString := "0x"

	return func(r rune) ([]*Token, error) {
		if utils.IsHexDigit(r) {
			parsedString = parsedString + string(r)
		} else {
			tf.SetOverFlow(r)

			number, err := utils.HexToNumber(parsedString)
			if err != nil {
				return nil, errors.Wrap(err, "hexTokenizer stringToNumber")
			}
			return []*Token{NewToken(HexLiteral, parsedString, number)}, nil
		}
		return nil, nil
	}
}

// StringTokenizer processes string literals.
func StringTokenizer(tf *TokenCreator, initialRune rune) func(rune) ([]*Token, error) {
	startRune := string(initialRune)
	openingString := true
	var parsedString string
	tf.parsingString = true

	return func(runeChar rune) ([]*Token, error) {
		if openingString {
			openingString = false
			return nil, nil
		}
		if string(runeChar) == startRune {
			tf.parsingString = false
			return []*Token{NewToken(StringLiteral, "", parsedString)}, nil
		}

		parsedString += string(runeChar)
		return nil, nil
	}
}

// IdentifierTokenizer processes identifiers like variable names.
func IdentifierTokenizer(tf *TokenCreator, initialRune rune) func(rune) ([]*Token, error) {
	parsedString := string(initialRune)

	return func(runeChar rune) ([]*Token, error) {
		if !utils.IsIdentifierChar(runeChar, len(parsedString), tf.languageConfig.extendedIdentifierRunes, tf.languageConfig.identifierTermination) {
			tf.SetOverFlow(runeChar)

			if tf.commentParser.IsStartOfComment(parsedString) {
				return nil, nil
			}
			return []*Token{tf.languageConfig.tokenFromIdentifier(parsedString)}, nil
		}
		parsedString += string(runeChar)

		if strings.Contains(tf.languageConfig.identifierTermination, string(runeChar)) {
			return []*Token{tf.languageConfig.tokenFromIdentifier(parsedString)}, nil
		}
		return nil, nil
	}
}

// SymbolTokenizer processes operators
func SymbolTokenizer(tf *TokenCreator, initialRune rune) func(r rune) ([]*Token, error) {
	symbolsString := string(initialRune)

	createToken := func(overflowRune rune) ([]*Token, error) {
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
					return nil, nil
				}
			}
			if longestSymbol != "" {
				tokenID := tf.languageConfig.operatorTokens[longestSymbol]
				symbolTokens = append(symbolTokens, NewToken(tokenID, longestSymbol, longestSymbol))
				i += len(longestSymbol)
			} else {
				tokenID, found := tf.languageConfig.symbolTokens[rune(symbolsString[i])]
				if !found {
					return nil, fmt.Errorf("unknown symbol %s", string(symbolsString[i]))
				}
				symbolTokens = append(symbolTokens, NewToken(tokenID, string(symbolsString[i]), symbolsString[i]))
				i++
			}
		}

		return symbolTokens, nil
	}

	return func(r rune) ([]*Token, error) {
		if tf.languageConfig.IsCustomTokenizer(r) { // CustomTokenizers take priority over symbols
			return createToken(r)
		} else if _, found := tf.languageConfig.symbolTokens[r]; found {
			symbolsString += string(r)
		} else {
			return createToken(r)
		}

		return nil, nil
	}
}
