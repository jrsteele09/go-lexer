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

// WithLabelSettings is a LanguageOptions function for setting label terminators and label tokens.
func WithLabelSettings(terminator rune, labelTokenID TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.labelTerminator = &terminator
		ll.labelToken = labelTokenID
	}
}

// WithExtraIdentifierRunes is a string of individual runes that can also be used to name identifiers
func WithExtraIdentifierRunes(extraIdentifierRunes string) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.extraIdentifierRunes = extraIdentifierRunes
	}
}
