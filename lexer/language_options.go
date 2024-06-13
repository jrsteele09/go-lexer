package lexer

// LanguageOptions is a function type used for configuring the LanguageConfig.
type LanguageOptions func(ll *LanguageConfig)

// WithOperators is a LanguageOptions function for setting the map of operator tokens
func WithOperators(ot map[string]TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.operatorTokens = ot
	}
}

// WithSymbols is a LanguageOptions function for setting the map of delimeter tokens
func WithSymbols(dt map[rune]TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.symbolTokens = dt
	}
}

// WithCommentMap is a LanguageOptions function for setting the map of comment delimiters.
func WithCommentMap(cm map[string]string) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.comments = cm
	}
}

// WithTokenCreators is a LanguageOptions function for setting custom token creators.
func WithTokenCreators(tc ...func(identifier string) *Token) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.tokenCreators = make([]func(identifier string) *Token, 0)
		ll.tokenCreators = append(ll.tokenCreators, tc...)
	}
}

// WithExtendendedIdentifierRunes is a string of individual runes that can also be used to name identifiers
func WithExtendendedIdentifierRunes(extraIdentifierRunes string, identifierTermination string) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.extendedIdentifierRunes = extraIdentifierRunes
		ll.identifierTermination = identifierTermination
	}
}
