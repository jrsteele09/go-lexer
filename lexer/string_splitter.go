package lexer

import "strings"

func StringSplitter(splitStr string, listOfQuotes ...string) func(str string) []string {
	quoteMap := make(map[string]struct{})

	for _, q := range listOfQuotes {
		quoteMap[q] = struct{}{}
	}

	isQuote := func(s string) bool {
		if _, found := quoteMap[s]; found {
			return true
		} else {
			return false
		}
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
		var endString, spacer string

		commands := make([]string, len(split))
		for _, s := range split {
			if startQuote(s, endString) {
				endString, spacer = s[:1], " "
			}
			commands[idx] = strings.TrimSpace(commands[idx] + spacer + s)
			if endQuote(s, endString) {
				endString, spacer = "", ""
			}
			if endString == "" {
				idx++
			}
		}

		return commands[:idx]
	}
}
