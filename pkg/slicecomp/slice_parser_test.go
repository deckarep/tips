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
