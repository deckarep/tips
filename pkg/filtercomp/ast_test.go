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
	"testing"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/stretchr/testify/assert"
)

func TestTextAST_Eval(t *testing.T) {

	type testCase struct {
		inputTxt    string
		matchTxt    string
		checkType   TextASTCheckType
		shouldMatch bool
	}

	cases := []testCase{
		{inputTxt: "foo", matchTxt: "foo", checkType: EqualityCheck, shouldMatch: true},
		{inputTxt: "foo", matchTxt: "bar", checkType: EqualityCheck, shouldMatch: false},

		{inputTxt: "foo", matchTxt: "foobar", checkType: PrefixCheck, shouldMatch: true},
		{inputTxt: "foo", matchTxt: "doobar", checkType: PrefixCheck, shouldMatch: false},

		{inputTxt: "foo", matchTxt: "poofoo", checkType: SuffixCheck, shouldMatch: true},
		{inputTxt: "foo", matchTxt: "doodoo", checkType: SuffixCheck, shouldMatch: false},

		{inputTxt: "foo", matchTxt: "goofoodie", checkType: PrefixCheck | SuffixCheck, shouldMatch: true},
		{inputTxt: "foo", matchTxt: "goooodie", checkType: PrefixCheck | SuffixCheck, shouldMatch: false},
	}

	for _, tc := range cases {
		node := &TextAST{
			checkType: tc.checkType,
			val:       tc.inputTxt,
		}

		assert.Equal(t, node.val, tc.inputTxt)

		s := mapset.NewSet[string](tc.matchTxt)
		assert.Equal(t, tc.shouldMatch, node.Eval(s))
	}
}

func TestParenAST_Eval(t *testing.T) {
	node := &ParenAST{
		exp: &TextAST{
			val: "foo",
		},
	}

	s := mapset.NewSet[string]("foo")
	assert.True(t, node.Eval(s))
}

func TestAndAST_Eval(t *testing.T) {
	node := &AndAST{
		left: &TextAST{
			val: "foo",
		},
		right: &TextAST{
			val: "bar",
		},
	}

	s := mapset.NewSet[string]("foo", "bar")
	assert.True(t, node.Eval(s))

	s = mapset.NewSet[string]("bar")
	assert.False(t, node.Eval(s))
}

func TestOrAST_Eval(t *testing.T) {
	node := &OrAST{
		left: &TextAST{
			val: "foo",
		},
		right: &TextAST{
			val: "bar",
		},
	}

	s := mapset.NewSet[string]("foo")
	assert.True(t, node.Eval(s))

	s = mapset.NewSet[string]("bar")
	assert.True(t, node.Eval(s))
}

func TestNegatedAST_Eval(t *testing.T) {
	node := &NegatedAST{
		exp: &TextAST{
			val: "foo",
		},
	}

	s := mapset.NewSet[string]("foo")
	assert.False(t, node.Eval(s))
}

func TestDumpAST(t *testing.T) {
	filter := "(foo | bar), !baz"
	tokens := Tokenize([]byte(filter))
	p := NewParser(tokens)
	ast, err := p.Parse()

	assert.NoError(t, err)
	assert.NotNil(t, ast)

	DumpAST(ast, 4)

	// This should be ok too and just do nothing.
	DumpAST(nil, 4)
}
