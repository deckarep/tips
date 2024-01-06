package pkg

import (
	"testing"
)

func TestParseSlice(t *testing.T) {
	// Neither defined
	s := ParseSlice(":", 0)
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
	s = ParseSlice("0:", 0)

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
	s = ParseSlice(":5", 0)

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
	s = ParseSlice("0:5", 0)

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