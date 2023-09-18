package lexer

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const (
	NullType = iota
	EOFType
	EndOfLineType
	IntegerLiteral
	NumberLiteral
	StringLiteral
	LastStdLiteral
)

type Lexer struct {
	overflowRune      *rune
	currentCommentEnd *string
	parsingString     bool
	firstLineToken    bool
	language          *LanguageConfig
}

func NewLexer(language *LanguageConfig) *Lexer {
	return &Lexer{
		// tokens:   make([]*Token, 0),
		language: language,
	}
}

func (l *Lexer) TokenizeLine(line string, lineNo uint) ([]Token, error) {
	var tokens []Token

	addNewToken := func(column int, token *Token) {
		if token == nil {
			return
		}
		l.firstLineToken = false
		token.SourceLine = lineNo
		token.SourceColumn = uint(column)
		tokens = append(tokens, *token)
	}

	tokenFactory := l.TokenFactory()

	skipNextChar := false
	l.firstLineToken = true

	for i, r := range line {
		if l.terminateDueToComment() {
			break
		}
		if string(r) == "\n" {
			break
		}
		if skipNextChar {
			skipNextChar = false
			continue
		}
		if l.endOfComment(i, line) {
			skipNextChar = true
			continue
		}
		if l.inComment() {
			continue
		}
		if l.startOfComment(i, line) {
			skipNextChar = true
		}

		addNewToken(i, tokenFactory(r))
		if l.hasRuneOverflow() {
			addNewToken(i, tokenFactory(*l.overflowRune))
			l.overflowRune = nil
		}
	}

	addNewToken(len(line), tokenFactory('\n')) // Finish off parsing the end of the line token
	if len(tokens) > 0 {
		addNewToken(len(line), NewToken(EndOfLineType, "\n", nil))
	}

	return tokens, nil
}

func (l *Lexer) terminateDueToComment() bool {
	if l.currentCommentEnd == nil {
		return false
	}
	if *l.currentCommentEnd != "\n" {
		return false
	}
	l.currentCommentEnd = nil
	return true
}

func (l *Lexer) startOfComment(idx int, line string) bool {
	if l.parsingString {
		return false
	}
	if l.currentCommentEnd != nil {
		return false
	}
	if idx >= len(line)-1 {
		return false
	}
	c := line[idx]
	n := line[idx+1]

	currentAndNextRune := string(string(c) + string(n))

	if comment, ok := l.language.comments[currentAndNextRune]; ok {
		l.currentCommentEnd = &comment
		return true
	}

	return false
}

func (l *Lexer) endOfComment(idx int, line string) bool {
	if l.currentCommentEnd == nil {
		return false
	}
	if idx >= len(line)-1 {
		return false
	}
	c := line[idx]
	n := line[idx+1]
	currentAndNextRune := string(string(c) + string(n))

	if currentAndNextRune != *l.currentCommentEnd {
		return false
	}

	l.currentCommentEnd = nil
	return true
}

func (l *Lexer) inComment() bool {
	return l.currentCommentEnd != nil
}

func (l *Lexer) hasRuneOverflow() bool {
	return !l.inComment() && l.overflowRune != nil
}

func (l *Lexer) TokenFactory() func(runeChar rune) *Token {
	var tokenizer func(r rune) *Token

	return func(runeChar rune) *Token {
		var token *Token

		if tokenizer == nil {
			if l.inComment() {
				return nil
			}
			if singleCharToken := l.language.createSingleCharToken(runeChar); singleCharToken != nil {
				return singleCharToken
			}
			if unicode.IsSpace(runeChar) {
				return nil
			}
			var err error
			if tokenizer, err = l.getTokenizer(runeChar); err != nil {
				log.Println(err)
				return nil
			}
		}

		token = tokenizer(runeChar)
		if token != nil {
			tokenizer = nil
		}
		return token
	}
}

func (l *Lexer) getTokenizer(runeChar rune) (func(rune) *Token, error) {
	if IsDigit(runeChar) {
		return l.numberTokenizer(), nil

	} else if IsStringQuotes(runeChar) {
		return l.stringTokenizer(runeChar), nil

	} else if IsIdentifierChar(runeChar, 0) {
		return l.identifierTokenizer(), nil
	}

	return nil, errors.New("unknown character: \"" + string(runeChar) + "\"")
}

func (l *Lexer) numberTokenizer() func(rune) *Token {
	var parsedString string

	return func(runeChar rune) *Token {
		if IsDigit(runeChar) {
			parsedString = parsedString + string(runeChar)
		} else {
			l.overflowRune = &runeChar
			number, _ := l.stringToNumber(parsedString)
			switch number.(type) {
			case float64:
				return NewToken(NumberLiteral, parsedString, number)
			case int64:
				return NewToken(IntegerLiteral, parsedString, number)
			}
			return NewToken(NumberLiteral, parsedString, number)
		}
		return nil
	}
}

func (l *Lexer) stringToNumber(strNum string) (interface{}, error) {
	if strings.Contains(strNum, ".") {
		return strconv.ParseFloat(strNum, 10)
	}
	return strconv.ParseInt(strNum, 10, 64)
}

func (l *Lexer) stringTokenizer(runeChar rune) func(rune) *Token {
	startRune := string(runeChar)
	openingString := true
	var parsedString string
	l.parsingString = true

	return func(runeChar rune) *Token {
		if openingString {
			openingString = false
			return nil
		}
		if string(runeChar) == startRune {
			l.parsingString = false
			return NewToken(StringLiteral, "", parsedString)
		}

		parsedString += string(runeChar)
		return nil
	}
}

func (l *Lexer) identifierTokenizer() func(rune) *Token {
	var parsedString string

	return func(runeChar rune) *Token {
		if l.language.labelTerminator != nil && l.firstLineToken && runeChar == *l.language.labelTerminator {
			return NewToken(l.language.labelToken, parsedString+string(runeChar), "")
		}
		if !IsIdentifierChar(runeChar, len(parsedString)) {
			l.overflowRune = &runeChar
			return l.language.tokenFromIdentifier(parsedString)
		}
		parsedString += string(runeChar)
		return nil
	}
}

func IsDigit(runeChar rune) bool {
	const extraDigits = "."
	return unicode.IsDigit(runeChar) || strings.Contains(extraDigits, string(runeChar))
}

func IsIdentifierChar(runeChar rune, pos int) bool {
	if unicode.IsLetter(runeChar) {
		return true
	}
	if pos > 0 && IsDigit(runeChar) {
		return true
	}

	return strings.Contains("_#%", string(runeChar))
}

func IsStringQuotes(runeChar rune) bool {
	const startOfString = "\"'`"
	return strings.Contains(startOfString, string(runeChar))
}
