package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	parsedResult, err := ParseFilter("peanut, walnut, pecan")
	assert.NoError(t, err)
	assert.NotNil(t, parsedResult)

	if parsedResult == nil {
		t.Error("expected populated filters but got nil")
	}

	// TODO: more robust testing here.
}
