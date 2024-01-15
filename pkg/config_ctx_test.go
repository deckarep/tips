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

	i, e = ParseColumns("foo,bar, baz")
	assert.NotNil(t, i)
	assert.NotNil(t, e)
	assert.Equal(t, e.Cardinality(), 0)

	m := mapset.NewSet[string]("foo", "bar", "baz")
	assert.True(t, m.Equal(i))
}
