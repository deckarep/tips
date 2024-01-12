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
