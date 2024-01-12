package filtercomp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tokens := Tokenize([]byte("(foo | bar, baz)"))
	assert.Equal(t, len(tokens), 7)

	expectedTokens := []Token{
		{Name: "(", Kind: TokenKindSymbol},
		{Name: "foo", Kind: TokenKindName},
		{Name: "OR", Kind: TokenKindLogical},
		{Name: "bar", Kind: TokenKindName},
		{Name: "AND", Kind: TokenKindLogical},
		{Name: "baz", Kind: TokenKindName},
		{Name: ")", Kind: TokenKindSymbol},
	}

	assert.Equal(t, tokens, expectedTokens)
}
