package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	// Note: much more robust testing of the parser occurs in the filtercomp package.
	// It makes this testing somewhat redundant.

	// A filter expression with imbalanced parenthesis should return an error.
	_, err := ParseFilter("(hello, world")
	assert.Error(t, err, "imbalanced parenthesis detected")

	// An empty filter should return no error with a nil ast.
	ast, err := ParseFilter("")
	assert.NoError(t, err)
	assert.Nil(t, ast)

	// An empty filter should return no error with a nil ast.
	_, err = ParseFilter("how are you doing?")
	assert.Error(t, err, "parser did not run to completion, tokens were not fully consumed")
}
