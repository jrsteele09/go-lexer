package lexer

type completed bool // Used to signal that the current tokenizer is completed
type TokenizerHandler func(r rune) ([]Token, completed, error)
type TokenizerFunc func(tf *TokenCreator, initialString string) TokenizerHandler

// LanguageConfig is the struct containing the configurations for the lexer.
type LanguageConfig struct {
	Keywords                map[string]TokenIdentifier
	Operators               map[string]TokenIdentifier      // Multi-character operator tokens e.g. "<=", "=="
	PrefixTokenizers        map[string]TokenizerFunc        // Language-specific tokenizers keyed by their trigger string
	Symbols                 map[rune]TokenIdentifier        // Single-rune symbol tokens
	Comments                map[string]string               // Comment delimiters: open -> close, e.g. "//" -> "\n"
	ExtendedIdentifierRunes string                          // Extra runes that are valid inside an identifier name
	IdentifierTermination   string                          // Runes that end an identifier and are included in it, e.g. ":" for labels
	TokenCreators           []func(identifier string) Token // Custom token creators called when no keyword matches
}

// NewLexerLanguage creates a new LanguageConfig from the provided configuration.
func NewLexerLanguage(config LanguageConfig) *LanguageConfig {
	return &config
}

// tokenFromIdentifier searches for a token matching the given identifier
// using custom token creators and returns the token if found.
func (ll *LanguageConfig) tokenFromIdentifier(identifier string) Token {
	if tokenID, ok := ll.Keywords[identifier]; ok {
		return NewToken(tokenID, identifier, nil)
	}

	for _, c := range ll.TokenCreators {
		if t := c(identifier); t.ID != NullType {
			return t
		}
	}
	return Token{}
}

func (ll *LanguageConfig) IsCustomTokenizer(parsedString string) bool {
	_, found := ll.PrefixTokenizers[parsedString]
	return found
}

func (ll *LanguageConfig) Tokenizer(parsedString string) TokenizerFunc {
	tokenizer := ll.PrefixTokenizers[parsedString]
	return tokenizer
}
