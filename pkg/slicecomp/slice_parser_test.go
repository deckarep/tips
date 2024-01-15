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

package slicecomp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSlice(t *testing.T) {
	// Empty string is just an empty slice.
	slice, err := ParseSlice(" ", 0)
	assert.NoError(t, err)
	assert.Nil(t, slice)
	assert.False(t, slice.IsDefined())

	// Invalid syntax
	slice, err = ParseSlice("[0:$]", 0)
	assert.Error(t, err)
	assert.Nil(t, slice)

	// Complete and correct
	slice, err = ParseSlice("[0:5]", 0)
	assert.NoError(t, err)
	assert.NotNil(t, slice)
	assert.Equal(t, *slice.From, 0)
	assert.Equal(t, *slice.To, 5)
	assert.True(t, slice.IsDefined())

	// Lower-bound only
	slice, err = ParseSlice("[0:]", 0)
	assert.NoError(t, err)
	assert.NotNil(t, slice)
	assert.Nil(t, slice.To)
	assert.Equal(t, *slice.From, 0)
	assert.True(t, slice.IsDefined())

	// Upper bound only
	slice, err = ParseSlice("[:5]", 0)
	assert.NoError(t, err)
	assert.NotNil(t, slice)
	assert.Nil(t, slice.From)
	assert.Equal(t, *slice.To, 5)
	assert.True(t, slice.IsDefined())
}
