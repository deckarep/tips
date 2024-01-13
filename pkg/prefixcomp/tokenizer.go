package prefixcomp

import (
	"strings"
	"unicode"
)

// Token types
const (
	TokenEOF = iota
	TokenWord
	TokenOr
	TokenInteger
	TokenLeftBracket
	TokenRightBracket
	TokenColon
	TokenAll
)

// Token represents a lexical token.
type Token struct {
	Type  int
	Value string
}

// Tokenizer holds the state of the scanner.
type Tokenizer struct {
	input string
	pos   int
}

func Tokenize(input string) ([]Token, error) {
	tokenizer := NewTokenizer(input)

	var tokens []Token
	for token := tokenizer.Next(); token.Type != TokenEOF; token = tokenizer.Next() {
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// NewTokenizer returns a new instance of Tokenizer.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: strings.TrimSpace(input)}
}

// Next returns the next token from the input.
func (t *Tokenizer) Next() Token {
	t.skipWhitespace()

	if t.pos >= len(t.input) {
		return Token{Type: TokenEOF}
	}

	switch t.input[t.pos] {
	case '*', '@':
		t.pos++
		// normalize(@) -> *
		return Token{Type: TokenAll, Value: "*"}
	case '[':
		t.pos++
		return Token{Type: TokenLeftBracket, Value: "["}
	case ']':
		t.pos++
		return Token{Type: TokenRightBracket, Value: "]"}
	case ':':
		t.pos++
		return Token{Type: TokenColon, Value: ":"}
	case '|':
		t.pos++
		return Token{Type: TokenOr, Value: "|"}
	default:
		if unicode.IsDigit(rune(t.input[t.pos])) {
			return t.lexInteger()
		} else if unicode.IsLetter(rune(t.input[t.pos])) {
			return t.lexWord()
		}
	}

	// If no token is recognized, advance and return EOF (error handling can be added here)
	t.pos++
	return Token{Type: TokenEOF}
}

// lexInteger scans an integer token.
func (t *Tokenizer) lexInteger() Token {
	start := t.pos
	for t.pos < len(t.input) && unicode.IsDigit(rune(t.input[t.pos])) {
		t.pos++
	}
	return Token{Type: TokenInteger, Value: t.input[start:t.pos]}
}

// lexWord scans a word token.
func (t *Tokenizer) lexWord() Token {
	start := t.pos
	for t.pos < len(t.input) && unicode.IsLetter(rune(t.input[t.pos])) {
		t.pos++
	}
	return Token{Type: TokenWord, Value: t.input[start:t.pos]}
}

// skipWhitespace advances the position over any whitespace.
func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.input) && unicode.IsSpace(rune(t.input[t.pos])) {
		t.pos++
	}
}
