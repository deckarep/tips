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
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type TextASTCheckType int

const (
	// EqualityCheck is a simple direct == or Contains check.
	EqualityCheck TextASTCheckType = 0
	// PrefixCheck is a strings.HasPrefix check, linear scan when a set is checked.
	PrefixCheck TextASTCheckType = 1
	// SuffixCheck is a strings.HasSuffix check, linear scan when a set is checked.
	SuffixCheck TextASTCheckType = 2
	// PrefixCheck | SuffixCheck is an Index check, linear scan when a set is checked.
)

type AST interface {
	Eval(mapset.Set[string]) bool
}

type TextAST struct {
	checkType TextASTCheckType
	val       string
}

func (t *TextAST) Eval(s mapset.Set[string]) bool {
	if t.checkType == EqualityCheck {
		// Normal path check.
		return s.Contains(t.val)
	} else {
		// A prefix check must linearly scan the whole set
		// but these are tiny containers.
		var hasPrefix bool
		s.Each(func(item string) bool {
			if (t.checkType & (PrefixCheck | SuffixCheck)) == (PrefixCheck | SuffixCheck) {
				// When prefix/suffix check is applied it becomes an Index check.
				if strings.Index(item, t.val) > -1 {
					hasPrefix = true
				}
			} else if t.checkType&PrefixCheck == PrefixCheck {
				if strings.HasPrefix(item, t.val) {
					hasPrefix = true
				}
			} else {
				if strings.HasSuffix(item, t.val) {
					hasPrefix = true
				}
			}
			return false
		})
		return hasPrefix
	}
}

type OrAST struct {
	left  AST
	right AST
}

func (o *OrAST) Eval(s mapset.Set[string]) bool {
	return o.left.Eval(s) || o.right.Eval(s)
}

type AndAST struct {
	left  AST
	right AST
}

func (a *AndAST) Eval(s mapset.Set[string]) bool {
	return a.left.Eval(s) && a.right.Eval(s)
}

type ParenAST struct {
	exp AST
}

func (p *ParenAST) Eval(s mapset.Set[string]) bool {
	return p.exp.Eval(s)
}

type NegatedAST struct {
	exp AST
}

func (n *NegatedAST) Eval(s mapset.Set[string]) bool {
	return !n.exp.Eval(s)
}

func DumpAST(node AST, indent int) {
	if node == nil {
		return
	}

	// Prepare the indentation string
	indentStr := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *TextAST:
		// TODO: this isn't printing flags that get set
		fmt.Printf("%s- Text: %s\n", indentStr, n.val)
	case *OrAST:
		fmt.Printf("%s- OR\n", indentStr)
		DumpAST(n.left, indent+1)
		DumpAST(n.right, indent+1)
	case *AndAST:
		fmt.Printf("%s- AND\n", indentStr)
		DumpAST(n.left, indent+1)
		DumpAST(n.right, indent+1)
	case *ParenAST:
		fmt.Printf("%s- Parentheses\n", indentStr)
		DumpAST(n.exp, indent+1)
	default:
		fmt.Printf("%s- Unknown AST Type\n", indentStr)
	}
}
