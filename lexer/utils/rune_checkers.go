package utils

import (
	"strings"
	"unicode"
)

// IsDigit checks if a rune is a digit.
func IsDigit(runeChar rune) bool {
	const extraDigits = "."
	return unicode.IsDigit(runeChar) || strings.Contains(extraDigits, string(runeChar))
}

// IsBinaryDigit checks if a rune is a 0 or 1.
func IsBinaryDigit(runeChar rune) bool {
	return runeChar == '0' || runeChar == '1'
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

// IsIdentifierChar checks if a rune is valid in an identifier.
func IsIdentifierChar(runeChar rune, pos int, extraRunes string, terminatorRunes string) bool {
	if unicode.IsLetter(runeChar) {
		return true
	}
	if pos > 0 && IsDigit(runeChar) {
		return true
	}
	if pos > 0 && strings.Contains(terminatorRunes, string(runeChar)) {
		return true
	}

	return strings.Contains(extraRunes, string(runeChar))
}
