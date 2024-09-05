package lexer

type completed bool // Used to signal that the current tokenizer is completed
type TokenizerHandler func(r rune) ([]*Token, completed, error)
type TokenizerFunc func(tf *TokenCreator, initialString string) TokenizerHandler

// LanguageConfig is the struct containing the configurations for the lexer.
type LanguageConfig struct {
	keywordTokens           map[string]TokenIdentifier
	operatorTokens          map[string]TokenIdentifier       // Map of operator tokens.
	customTokenizers        map[string]TokenizerFunc         // customer tokenizers allow language custom tokenizers to be added to the lexer
	symbolTokens            map[rune]TokenIdentifier         // Delimeter tokens
	comments                map[string]string                // Map of comment delimiters.
	extendedIdentifierRunes string                           // Map of extra chars (runes) that can be part of an identifier name
	identifierTermination   string                           // Map of chars that will terminate an identifier, for example ":" could be used to identify a label
	tokenCreators           []func(identifier string) *Token // Custom token creators.
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
	if tokenID, ok := ll.keywordTokens[identifier]; ok {
		return NewToken(tokenID, identifier, nil)
	}

	for _, c := range ll.tokenCreators {
		if t := c(identifier); t != nil {
			return t
		}
	}
	return nil
}

func (ll *LanguageConfig) IsCustomTokenizer(parsedString string) bool {
	_, found := ll.customTokenizers[parsedString]
	return found
}

func (ll *LanguageConfig) Tokenizer(parsedString string) TokenizerFunc {
	tokenizer := ll.customTokenizers[parsedString]
	return tokenizer
}
