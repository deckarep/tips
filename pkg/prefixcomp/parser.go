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

package prefixcomp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/deckarep/tips/pkg/slicecomp"
)

/*
https://bnfplayground.pauliankline.com/
WARN: This should always match the parser.

<primaryfilter> ::= (<all> | <filter>) <ws> <slice>?
<all> ::= "*" | "@" | E
<filter> ::= <word> (<ws> <or> <ws> <word> <ws>)*
<slice> ::= "[" <integer>? ":" <integer>? "]"
<ws> ::= " "+ | E
<or> ::= "|"
<word> ::= ([a-z] | [A-Z])+
<integer> ::= [0-9]+
*/

func ParsePrimaryFilter(input string) (*PrimaryFilterAST, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

// PrimaryFilterAST represents the primary filter syntax.
type PrimaryFilterAST struct {
	// When All is true, it represents the '*' and it means ignore the Words entirely.
	All   bool
	Words []string
	Slice *slicecomp.Slice
}

func (p *PrimaryFilterAST) Query() string {
	var filterOn = "*"
	if !p.All {
		filterOn = strings.Join(p.Words, " | ")
	}

	var buf bytes.Buffer
	if p.Slice != nil {
		buf.WriteString("[")
		if p.Slice.From != nil {
			buf.WriteString(fmt.Sprintf("%d", *p.Slice.From))
		}
		buf.WriteString(":")
		if p.Slice.To != nil {
			buf.WriteString(fmt.Sprintf("%d", *p.Slice.To))
		}
		buf.WriteString("]")
	} else {
		buf.WriteString("")
	}

	return fmt.Sprintf("%s%s", filterOn, buf.String())
}

func (p *PrimaryFilterAST) String() string {
	var filterOn = "*"
	if !p.All {
		filterOn = fmt.Sprintf("%v", p.Words)
	}

	return fmt.Sprintf("PrimaryFilter(Words: %v, Slice: %s)", filterOn, sliceToString(p.Slice))
}

func (p *PrimaryFilterAST) IsAll() bool {
	return p.All
}

func (p *PrimaryFilterAST) Count() int {
	return len(p.Words)
}

func (p *PrimaryFilterAST) PrefixAt(idx int) string {
	return p.Words[idx]
}

func sliceToString(s *slicecomp.Slice) string {
	if s == nil {
		return "<nil-slice>"
	}
	const (
		nilVal = "<nil>"
	)
	var fromVal = nilVal
	if s.From != nil {
		fromVal = fmt.Sprintf("%d", *s.From)
	}
	var toVal = nilVal
	if s.To != nil {
		toVal = fmt.Sprintf("%d", *s.To)
	}

	return fmt.Sprintf("(from: %s, to: %s)", fromVal, toVal)
}

// Parser represents a parser.
type Parser struct {
	tokens  []Token
	current int
}

// NewParser creates a new parser.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse parses the tokens into an AST.
func (p *Parser) Parse() (*PrimaryFilterAST, error) {
	return p.parsePrimaryFilter()
}

// parsePrimaryFilter parses a primary filter.
func (p *Parser) parsePrimaryFilter() (*PrimaryFilterAST, error) {
	var words []string
	var slice *slicecomp.Slice
	var err error
	var useAll bool

	if p.isAtEnd() {
		// In this case, there was an empty token stream. Assume useAll=true!
		useAll = true
	} else if p.match(TokenAll) {
		useAll = true
		if p.isAtEnd() {
			// Do nothing: the other branches are not taken on purpose.
		} else if p.match(TokenLeftBracket) {
			slice, err = p.parseSlice()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unexpected token: %v", p.peek())
		}
	} else if p.match(TokenLeftBracket) {
		// Only slice was provided, useAll=true.
		useAll = true
		slice, err = p.parseSlice()
		if err != nil {
			return nil, err
		}
	} else {
		for {
			if p.isAtEnd() {
				break
			} else if p.match(TokenWord) {
				words = append(words, p.previous().Value)
			} else if p.match(TokenOr) {
				if !p.match(TokenWord) {
					return nil, fmt.Errorf("expected word after '|' symbol")
				}
				words = append(words, p.previous().Value)
			} else if p.match(TokenLeftBracket) {
				slice, err = p.parseSlice()
				if err != nil {
					return nil, err
				}
				break
			} else {
				return nil, fmt.Errorf("unexpected token: %v", p.peek())
			}
		}
	}

	return &PrimaryFilterAST{All: useAll, Words: words, Slice: slice}, nil
}

// parseSlice parses a slice.
func (p *Parser) parseSlice() (*slicecomp.Slice, error) {
	var err error

	// Initialize the start/end to -1.
	start := -1
	end := -1

	if p.match(TokenInteger) {
		start, err = strconv.Atoi(p.previous().Value)
		if err != nil {
			return nil, fmt.Errorf("invalid start index: %v", p.previous().Value)
		}
	}

	if !p.match(TokenColon) {
		return nil, fmt.Errorf("expected ':' in slice")
	}

	if p.match(TokenInteger) {
		end, err = strconv.Atoi(p.previous().Value)
		if err != nil {
			return nil, fmt.Errorf("invalid end index: %v", p.previous().Value)
		}
	}

	if !p.match(TokenRightBracket) {
		return nil, fmt.Errorf("expected ']' in slice")
	}

	var orNil = func(in int) *int {
		if in > -1 {
			return &in
		}
		return nil
	}

	return &slicecomp.Slice{From: orNil(start), To: orNil(end)}, nil
}

// match checks if the current token matches the given type.
func (p *Parser) match(t int) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	return false
}

// check checks if the current token is of the given type.
func (p *Parser) check(t int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// advance advances to the next token.
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if the parser has reached the end of the tokens.
func (p *Parser) isAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return false
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// previous returns the previous token.
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
