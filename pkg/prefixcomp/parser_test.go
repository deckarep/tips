package prefixcomp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParsePrefixFilterAll(t *testing.T) {
	ast, err := ParsePrimaryFilter("* [5:10]")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// All is true
	assert.True(t, ast.All)

	// Slice is created and correct.
	assert.True(t, ast.Slice.IsDefined())
	assert.Equal(t, *ast.Slice.From, 5)
	assert.Equal(t, *ast.Slice.To, 10)

	// Implicit all
	ast, err = ParsePrimaryFilter("[5:10]")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// All is true
	assert.True(t, ast.All)

	// Slice is created and correct.
	assert.True(t, ast.Slice.IsDefined())
	assert.Equal(t, *ast.Slice.From, 5)
	assert.Equal(t, *ast.Slice.To, 10)

	// Implicit all NO slice.
	ast, err = ParsePrimaryFilter("")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// All is true
	assert.True(t, ast.All)

	// Slice is created and correct.
	assert.False(t, ast.Slice.IsDefined())
}

func TestParser_ParsePrefixFilterWords(t *testing.T) {
	ast, err := ParsePrimaryFilter("foo | bar | baz [2:19]")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// Words are correct.
	assert.Equal(t, []string{"foo", "bar", "baz"}, ast.Words)

	// All is true
	assert.False(t, ast.All)

	// Slice is created and correct.
	assert.True(t, ast.Slice.IsDefined())
	assert.Equal(t, *ast.Slice.From, 2)
	assert.Equal(t, *ast.Slice.To, 19)

	// Words with no slice.
	ast, err = ParsePrimaryFilter("foo | bar | baz")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// Words are correct.
	assert.Equal(t, []string{"foo", "bar", "baz"}, ast.Words)

	// All is true
	assert.False(t, ast.All)

	// Slice is created and correct.
	assert.False(t, ast.Slice.IsDefined())
}
