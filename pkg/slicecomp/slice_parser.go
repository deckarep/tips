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

package slicecomp

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Slice EBNF
==========
<slice> ::= "[" <digit>? ":" <digit>? "]"
<digit> ::= [0-9]+
*/

type Slice struct {
	From *int
	To   *int
}

func (s *Slice) IsDefined() bool {
	return s != nil && (s.From != nil || s.To != nil)
}

func ParseSlice(input string, page int) (*Slice, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil, nil
	}

	p := NewParser(input)
	return p.ParseSlice()
}

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input}
}

func (p *Parser) ParseSlice() (*Slice, error) {
	if !p.match('[') {
		return nil, fmt.Errorf("expected '['")
	}

	start, err := p.parseDigit()
	if err != nil {
		start = -1 // Indicates no start index
	}

	if !p.match(':') {
		return nil, fmt.Errorf("expected ':'")
	}

	end, err := p.parseDigit()
	if err != nil {
		end = -1 // Indicates no end index
	}

	if !p.match(']') {
		return nil, fmt.Errorf("expected ']'")
	}

	var orNil = func(val int) *int {
		var result int
		if val > -1 {
			result = val
			return &result
		}
		return nil
	}

	from := orNil(start)
	to := orNil(end)

	return &Slice{
		From: from,
		To:   to,
	}, nil
}

func (p *Parser) parseDigit() (int, error) {
	start := p.pos
	for p.pos < len(p.input) && isDigit(p.input[p.pos]) {
		p.pos++
	}
	if start == p.pos {
		return 0, fmt.Errorf("expected digit")
	}
	return strconv.Atoi(p.input[start:p.pos])
}

func (p *Parser) match(expected rune) bool {
	if p.pos < len(p.input) && rune(p.input[p.pos]) == expected {
		p.pos++
		return true
	}
	return false
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
