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

package pkg

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/deckarep/tips/pkg/slicecomp"

	"github.com/stretchr/testify/assert"
)

func TestParseSlice(t *testing.T) {
	// Neither defined
	s, err := slicecomp.ParseSlice("[:]", 0)
	assert.NoError(t, err)
	if s == nil {
		t.Error("expected slice to be non-nil")
	}

	if s.From != nil && s.To != nil {
		t.Error("expected both from/to to be nil")
	}

	if s.IsDefined() {
		t.Error("expected slice to not be defined")
	}

	// Only lower bound
	s, err = slicecomp.ParseSlice("[0:]", 0)
	assert.NoError(t, err)

	if s.From != nil && *s.From != 0 {
		t.Errorf("expected from to be: %d, got: %d", 0, *s.From)
	}

	if s.To != nil {
		t.Error("expected to as nil but was not")
	}

	if !s.IsDefined() {
		t.Error("expected slice to be defined")
	}

	// Only upper bound
	s, err = slicecomp.ParseSlice("[:5]", 0)
	assert.NoError(t, err)

	if s.From != nil {
		t.Error("expected from as nil but was not")
	}

	if s.To != nil && *s.To != 5 {
		t.Errorf("expected To to be: %d, got: %d", 5, *s.From)
	}

	if !s.IsDefined() {
		t.Error("expected slice to be defined")
	}

	// Both lower and upper bound.
	s, err = slicecomp.ParseSlice("[0:5]", 0)
	assert.NoError(t, err)

	if s.From != nil && *s.From != 0 {
		t.Errorf("expected from to be: %d, got: %d", 0, *s.From)
	}

	if s.To != nil && *s.To != 5 {
		t.Errorf("expected from to be: %d, got: %d", 5, *s.To)
	}

	if !s.IsDefined() {
		t.Error("expected slice to be defined")
	}
}

func TestParseColumns(t *testing.T) {
	// Empty string results in nil.
	i, e := ParseColumns("")
	assert.NotNil(t, i)
	assert.NotNil(t, e)

	// Include set.
	i, e = ParseColumns("foo,bar, baz")
	assert.NotNil(t, i)
	assert.NotNil(t, e)
	assert.Equal(t, e.Cardinality(), 0)

	m := mapset.NewSet[string]("foo", "bar", "baz")
	assert.True(t, m.Equal(i))

	// Exclude set.
	i, e = ParseColumns("-foo,-bar, -baz")
	assert.NotNil(t, i)
	assert.NotNil(t, e)
	assert.Equal(t, i.Cardinality(), 0)

	m = mapset.NewSet[string]("foo", "bar", "baz")
	assert.True(t, m.Equal(e))

	// Both include and exclude set.
	i, e = ParseColumns("-foo, bar, -baz, bon, -fee")
	assert.NotNil(t, i)
	assert.NotNil(t, e)

	m = mapset.NewSet[string]("bar", "bon")
	assert.True(t, m.Equal(i))

	m = mapset.NewSet[string]("foo", "fee", "baz")
	assert.True(t, m.Equal(e))
}
