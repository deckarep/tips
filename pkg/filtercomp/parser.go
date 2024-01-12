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

import "errors"

/*
// Possible TODO:
// - support negation with a !(not) prefix or use a -(minus) (should be very obvious/easy to grok)
// - support conditionals such as: lastseen > 2hr ago

// Future TODO:
// - optimizations or simplifications.

EBNF Definition Below.
WARNING: Keep this updated in lock-step with the parser!
Guaranteed to be free of left recursion: https://bnfplayground.pauliankline.com/

<expression> ::= <factor> <logexp>*
<logexp> ::= ("|" | ",") <factor>
<factor> ::= "!"? (<name> | "(" <expression> ")")
<name> ::= "*"? [a-z]+ "*"?
*/

type Parser struct {
	idx    int
	tokens []Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// Parse is the entrypoint and kick off method to start the parsing process.
func (p *Parser) Parse() (AST, error) {
	if err := p.preflightCheck(); err != nil {
		return nil, err
	}
	// Entry point into the parse tree.
	ast, err := p.parseExp()
	if err != nil {
		return nil, err
	}

	if !p.isEOF() {
		return nil, errors.New("parser did not run to completion, tokens were not fully consumed")
	}

	return ast, nil
}

// preflightCheck does a few early checks before parsing begins. This token pre-scan is perhaps not the most efficient
// but please be real: we're parsing a 'baby' DSL passed in on the command-line. This won't be a bottleneck for a long
// time.
func (p *Parser) preflightCheck() error {
	return p.detectImbalancedParenthesis()
}

// detectImbalancedParenthesis does a quick linear scan through the tokens. In hindsight this should probably
// be checked in the tokenize phase.
func (p *Parser) detectImbalancedParenthesis() error {
	var openParenCount int
	var closedParenCount int
	for _, t := range p.tokens {
		if t.Kind == TokenKindSymbol {
			if t.Name == "(" {
				openParenCount++
			}
			if t.Name == ")" {
				closedParenCount++
			}
		}
	}
	if openParenCount != closedParenCount {
		return errors.New("imbalanced parenthesis detected")
	}
	return nil
}

func (p *Parser) isEOF() bool {
	return p.idx >= len(p.tokens)
}

func (p *Parser) peekToken() Token {
	if p.isEOF() {
		// Return a default Token or handle the EOF scenario
		return Token{Name: "EOF", Kind: -1}
	}
	return p.tokens[p.idx]
}

func (p *Parser) consumeToken() (Token, error) {
	if p.isEOF() {
		// Handle the EOF scenario, maybe return a default Token or panic
		return Token{}, errors.New("attempt to consume Token at EOF")
	}
	t := p.tokens[p.idx]
	p.idx++
	return t, nil
}

func (p *Parser) parseExp() (AST, error) {
	leftFactor, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	// The loop handles zero or more *
	for !p.isEOF() {
		// If a Logical Exp is found, the leftFactor actually becomes the left child.
		logExp, err := p.parseLogexp(leftFactor)
		if err != nil {
			return nil, err
		}
		if logExp == nil {
			// If logExp is nil, it means we didn't find a valid logexp, so break out of the loop
			break
		}
		leftFactor = logExp
	}
	return leftFactor, nil
}

func (p *Parser) parseLogexp(leftFactor AST) (AST, error) {
	t := p.peekToken()
	switch t.Name {
	case "OR":
		_, err := p.consumeToken() // consume 'OR'
		if err != nil {
			return nil, err
		}
		rightFactor, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return &OrAST{
			left:  leftFactor,
			right: rightFactor,
		}, nil
	case "AND":
		_, err := p.consumeToken() // consume 'AND'
		if err != nil {
			return nil, err
		}
		rightFactor, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return &AndAST{
			left:  leftFactor,
			right: rightFactor,
		}, nil
	default:
		// Return nil to indicate that no valid logexp was found
		return nil, nil
	}
}

func (p *Parser) parseFactor() (AST, error) {
	t := p.peekToken()

	var shouldNegate bool
	if t.Name == "!" {
		_, err := p.consumeToken() // ! - consume the negation which is optional.
		if err != nil {
			return nil, err
		}
		shouldNegate = true
	}

	t = p.peekToken()
	var selectedNode AST
	if t.Name == "(" {
		_, err := p.consumeToken() // (
		if err != nil {
			return nil, err
		}
		exp, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		parenAST := &ParenAST{
			exp: exp,
		}
		_, err = p.consumeToken() // )
		if err != nil {
			return nil, err
		}
		selectedNode = parenAST
	} else {
		ast, err := p.parseName()
		if err != nil {
			return nil, err
		}
		selectedNode = ast
	}

	if shouldNegate {
		return &NegatedAST{exp: selectedNode}, nil
	}

	return selectedNode, nil
}

func (p *Parser) parseName() (AST, error) {
	t := p.peekToken()

	// Default type check.
	checkFlags := EqualityCheck

	// Check if there is an optional suffix check as in: *foo
	// Note: Both can be applied as in: *foo* which becomes an Index check.
	if t.Name == "*" {
		_, err := p.consumeToken() // Consume the *
		if err != nil {
			return nil, err
		}
		checkFlags |= SuffixCheck
	}

	t = p.peekToken()
	if t.Kind != TokenKindName {
		return nil, errors.New("expected to be a Name token")
	}

	nameToken, err := p.consumeToken() // Just consume the Name Token for now.
	if err != nil {
		return nil, err
	}

	// Check if there is an optional prefix check as in: foo*
	// Note: Both can be applied as in: *foo* which becomes an Index check.
	t = p.peekToken()
	if t.Name == "*" {
		_, err = p.consumeToken() // Consume the *
		if err != nil {
			return nil, err
		}
		checkFlags |= PrefixCheck
	}

	return &TextAST{val: nameToken.Name, checkType: checkFlags}, nil
}
