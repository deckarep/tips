package pkg

import (
	"testing"

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

func TestParsePrefixFilter(t *testing.T) {
	// The prefix for everything prefix
	pf := ParsePrefixFilter("*")
	assert.NotNil(t, pf)
	assert.True(t, pf.IsAll())

	// The prefix for everything prefix (the same thing)
	pf = ParsePrefixFilter("@")
	assert.NotNil(t, pf)
	assert.True(t, pf.IsAll())

	// The prefix with multiple OR conditions
	pf = ParsePrefixFilter("foo|bar | baz")
	assert.NotNil(t, pf)
	assert.False(t, pf.IsAll())

	assert.Equal(t, pf.Count(), 3)

	expectedPrefixOrder := []string{"bar", "baz", "foo"}
	for i := 0; i < pf.Count(); i++ {
		assert.Equal(t, pf.PrefixAt(i), expectedPrefixOrder[i])
	}
}
