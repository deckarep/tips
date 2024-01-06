package pkg

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestParseFilter(t *testing.T) {
	parsedResult := ParseFilter("tag:peanut, walnut, pecan")

	if parsedResult == nil {
		t.Error("expected populated filters but got nil")
	}

	if len(parsedResult) != 1 {
		t.Errorf("expected 1 item but got: %d", len(parsedResult))
	}

	const expectedKey = "tag"
	expectedItems := mapset.NewSet[string]("peanut", "walnut", "pecan")

	s, _ := parsedResult[expectedKey]
	if s == nil {
		t.Error("expected non-nil result from ParseFilter")
	}

	if s.Difference(expectedItems).Cardinality() != 0 {
		t.Errorf("expected result to contain key: %s with set items: %s but got: %s", expectedKey, expectedItems.String(), s.String())
	}
}
