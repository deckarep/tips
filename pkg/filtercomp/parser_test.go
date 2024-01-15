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
