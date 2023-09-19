package stringssplitter

import "strings"

// StringSplitter returns a closure that takes a string and splits it based on the given delimiter (splitStr).
// The function also takes an optional list of quotes (listOfQuotes) to treat substrings within quotes as a single unit during the split.
//
// For example, given splitStr as "," and listOfQuotes as ['"']:
// - The string 'apple, "orange, banana", cherry' will be split into ['apple', '"orange, banana"', 'cherry']
//
// Parameters:
// - splitStr: The delimiter used for splitting the string.
// - listOfQuotes: Optional list of quote characters to treat substrings within them as a single unit.
//
// Returns:
// A function that takes a string and returns a slice of strings based on the delimiter and quotes.
func StringSplitter(splitStr string, listOfQuotes ...string) func(str string) []string {
	quoteMap := make(map[string]struct{})

	for _, q := range listOfQuotes {
		quoteMap[q] = struct{}{}
	}

	isQuote := func(s string) bool {
		if _, found := quoteMap[s]; found {
			return true
		}
		return false
	}

	startQuote := func(s, endQuote string) bool { return endQuote == "" && isQuote(s[:1]) }
	endQuote := func(s, endQuote string) bool { return endQuote != "" && s[len(s)-1:] == endQuote }

	return func(str string) []string {
		if strings.TrimSpace(str) == "" {
			return []string{}
		}

		split := strings.Split(strings.TrimSpace(str), splitStr)

		if len(split) < 2 {
			return split
		}

		idx := 0
		var endQuoteString, spacer string

		commands := make([]string, len(split))
		for _, s := range split {
			if startQuote(s, endQuoteString) {
				endQuoteString, spacer = s[:1], " "
			}
			commands[idx] = strings.TrimSpace(commands[idx] + spacer + s)
			if endQuote(s, endQuoteString) {
				endQuoteString, spacer = "", ""
			}
			if endQuoteString == "" {
				idx++
			}
		}

		return commands[:idx]
	}
}
