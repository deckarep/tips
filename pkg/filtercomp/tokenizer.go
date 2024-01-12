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
	"bytes"
	"unicode"
)

type TokenKind int

const (
	TokenKindName    TokenKind = iota
	TokenKindSymbol  TokenKind = iota
	TokenKindLogical TokenKind = iota
)

type Token struct {
	Name string
	Kind TokenKind
}

func Tokenize(data []byte) []Token {
	var tokens []Token
	var current bytes.Buffer

	for i := 0; i < len(data); i++ {
		b := data[i]
		switch b {
		case '(', ')', ',', '|', '*':
			// Flush any accumulated text as a Name Token
			if current.Len() > 0 {
				tokens = append(tokens, Token{Name: current.String(), Kind: TokenKindName})
				current.Reset()
			}
			// Append the symbol or logical operator
			var kind TokenKind
			var name = string(b)
			if b == ',' || b == '|' {
				if b == ',' {
					name = "AND"
				} else {
					name = "OR"
				}
				kind = TokenKindLogical
			} else {
				kind = TokenKindSymbol
			}
			tokens = append(tokens, Token{Name: name, Kind: kind})
		default:
			if unicode.IsSpace(rune(b)) {
				// If whitespace, flush the current buffer as a Name Token
				if current.Len() > 0 {
					tokens = append(tokens, Token{Name: current.String(), Kind: TokenKindName})
					current.Reset()
				}
			} else {
				// Otherwise, accumulate the characters
				current.WriteByte(b)
			}
		}
	}

	// Flush any remaining text as a Name Token
	if current.Len() > 0 {
		tokens = append(tokens, Token{Name: current.String(), Kind: TokenKindName})
	}

	return tokens
}
