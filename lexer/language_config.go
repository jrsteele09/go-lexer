package lexer

// LanguageConfig is the struct containing the configurations for the lexer.
type LanguageConfig struct {
	operatorTokens       map[string]TokenIdentifier       // Map of operator tokens.
	symbolTokens         map[rune]TokenIdentifier         // Delimeter tokens
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
