package lexer

type LanguageOptions func(ll *LanguageConfig)

func WithSingleRuneMap(rm map[rune]TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.singleRuneTokens = rm
	}
}

func WithCommentMap(cm map[string]string) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.comments = cm
	}
}

func WithTokenCreators(tc ...func(identifier string) *Token) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.tokenCreators = make([]func(identifier string) *Token, 0)
		ll.tokenCreators = append(ll.tokenCreators, tc...)
	}
}

func WithLabelSettings(terminator rune, labelTokenID TokenIdentifier) LanguageOptions {
	return func(ll *LanguageConfig) {
		ll.labelTerminator = &terminator
		ll.labelToken = labelTokenID
	}
}

// LanguageConfig
type LanguageConfig struct {
	singleRuneTokens map[rune]TokenIdentifier
	comments         map[string]string
	tokenCreators    []func(identifier string) *Token
	labelTerminator  *rune
	labelToken       TokenIdentifier
}

func NewLexerLanguage(opts ...LanguageOptions) *LanguageConfig {
	ll := &LanguageConfig{}
	for _, opt := range opts {
		opt(ll)
	}
	return ll
}

func (ll *LanguageConfig) createSingleCharToken(r rune) *Token {
	if tokenId, ok := ll.singleRuneTokens[r]; ok {
		return NewToken(tokenId, string(r), nil)
	}
	return nil
}

func (ll *LanguageConfig) tokenFromIdentifier(identifier string) *Token {
	for _, c := range ll.tokenCreators {
		if t := c(identifier); t != nil {
			return t
		}
	}
	return nil
}
