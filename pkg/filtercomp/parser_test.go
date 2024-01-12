package filtercomp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tokens := Tokenize([]byte("(foo | bar, baz)"))
	p := NewParser(tokens)
	ast, err := p.Parse()

	assert.NoError(t, err)
	assert.NotNil(t, ast)

	// Yes, there's a better way to do this...but for now...
	assert.IsType(t, (*ParenAST)(nil), ast)
	if innerAST, ok := ast.(*ParenAST); ok {
		if andAST, ok := innerAST.exp.(*AndAST); ok {
			if leftOrAST, ok := andAST.left.(*OrAST); ok {
				if rightTextAST, ok := andAST.right.(*TextAST); ok {
					if leftInnerTextAST, ok := leftOrAST.left.(*TextAST); ok {
						if rightInnerTextAST, ok := leftOrAST.right.(*TextAST); ok {
							assert.Equal(t, rightTextAST.val, "baz")
							assert.Equal(t, leftInnerTextAST.val, "foo")
							assert.Equal(t, rightInnerTextAST.val, "bar")
						}
					}
				}
			}
		}
	}
}

func TestTextNode(t *testing.T) {
	type pair struct {
		name      string
		checkType TextASTCheckType
	}

	textNodes := []pair{
		{
			name:      "foo",
			checkType: EqualityCheck,
		},
		{
			name:      "bar*",
			checkType: PrefixCheck,
		},
		{
			name:      "*baz",
			checkType: SuffixCheck,
		},
		{
			name:      "*bom*",
			checkType: PrefixCheck | SuffixCheck,
		},
	}

	for _, textNode := range textNodes {
		tokens := Tokenize([]byte(textNode.name))

		p := NewParser(tokens)
		ast, err := p.Parse()

		assert.NoError(t, err)
		assert.NotNil(t, ast)

		assert.IsType(t, (*TextAST)(nil), ast)
		if tn, ok := ast.(*TextAST); ok {
			assert.Equal(t, tn.val, strings.Replace(textNode.name, "*", "", -1))
			assert.Equal(t, tn.checkType, textNode.checkType)
		}
	}
}
