/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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

	assert.Equal(t, ast.String(), "PrimaryFilter(Words: *, Slice: (from: 5, to: 10))")

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

	assert.Equal(t, ast.String(), "PrimaryFilter(Words: *, Slice: (from: 5, to: 10))")

	// Implicit all NO slice.
	ast, err = ParsePrimaryFilter("")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// All is true
	assert.True(t, ast.All)

	// Slice is created and correct.
	assert.False(t, ast.Slice.IsDefined())
	assert.Equal(t, ast.String(), "PrimaryFilter(Words: *, Slice: <nil-slice>)")
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

	assert.Equal(t, ast.String(), "PrimaryFilter(Words: [foo bar baz], Slice: (from: 2, to: 19))")

	// Words with no slice.
	ast, err = ParsePrimaryFilter("foo | bar")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// Words are correct.
	assert.Equal(t, []string{"foo", "bar"}, ast.Words)

	// All is true
	assert.False(t, ast.All)

	// Slice is created and correct.
	assert.False(t, ast.Slice.IsDefined())

	assert.Equal(t, ast.String(), "PrimaryFilter(Words: [foo bar], Slice: <nil-slice>)")
}
