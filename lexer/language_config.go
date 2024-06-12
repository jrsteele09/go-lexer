package lexer

// LanguageOptions is a function type used for configuring the LanguageConfig.
type LanguageOptions func(ll *LanguageConfig)

// WithSingleRuneMap is a LanguageOptions function for setting the map of single-rune tokens.
func WithSingleRuneMap(rm map[rune]TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.singleRuneTokens = rm
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

// LanguageConfig is the struct containing the configurations for the lexer.
type LanguageConfig struct {
	singleRuneTokens     map[rune]TokenIdentifier         // Map of single-rune tokens.
	comments             map[string]string                // Map of comment delimiters.
	extraIdentifierRunes string                           // Map of extra chars (runes) that can be part of an identifier name
	tokenCreators        []func(identifier string) *Token // Custom token creators.
	labelTerminator      *rune                            // Label terminator rune.
	labelToken           TokenIdentifier                  // Label token identifier.
}

// NewLexerLanguage creates a new LanguageConfig using the provided options.
func NewLexerLanguage(opts ...LanguageOptions) *LanguageConfig {
	ll := &LanguageConfig{}
	for _, opt := range opts {
		opt(ll)
	}
	return ll
}

// tokenFromIdentifier searches for a token matching the given identifier
// using custom token creators and returns the token if found.
func (ll *LanguageConfig) tokenFromIdentifier(identifier string) *Token {
	for _, c := range ll.tokenCreators {
		if t := c(identifier); t != nil {
			return t
		}
	}
	return nil
}
