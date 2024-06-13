package stringssplitter

import "strings"

type quoteMap map[string]struct{}

func (q quoteMap) isQuote(s string) bool {
	_, found := q[s]
	return found
}

func (q quoteMap) startQuote(s string) bool {
	return q.isQuote(s[:1])
}

func (q quoteMap) endQuote(s, currentQuote string) bool {
	return currentQuote != "" && s[len(s)-1:] == currentQuote
}

func newQuoteMap(quotes ...string) quoteMap {
	q := make(quoteMap)
	for _, quote := range quotes {
		q[quote] = struct{}{}
	}
	return q
}

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
func StringSplitter(splitToken string, listOfQuotes ...string) func(str string) []string {
	quoteMap := newQuoteMap(listOfQuotes...)

	return func(str string) []string {
		if strings.TrimSpace(str) == "" {
			return []string{}
		}

		listOfStrings := strings.Split(strings.TrimSpace(str), splitToken)

		if len(listOfStrings) < 2 {
			return listOfStrings
		}

		idx := 0
		var currentQuote string
		stringList := make([]string, len(listOfStrings))

		for i := range listOfStrings {
			s := listOfStrings[i]
			if currentQuote == "" && quoteMap.startQuote(strings.TrimLeft(s, " ")) {
				s = strings.TrimLeft(s, " ")
				currentQuote = s[:1]
			}
			if quoteMap.endQuote(strings.TrimRight(s, " "), currentQuote) {
				s = strings.TrimRight(s, " ")
				currentQuote = ""
			}
			stringList[idx] = (stringList[idx] + s)
			if currentQuote == "" {
				idx++
			} else {
				stringList[idx] = (stringList[idx] + splitToken) // Put back the split token
			}
		}

		return stringList[:idx]
	}
}
