package prefixcomp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParsePrefixFilter(t *testing.T) {
	ast, err := ParsePrimaryFilter("hello | world [5:10]")
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	fmt.Println(ast.String())
}
