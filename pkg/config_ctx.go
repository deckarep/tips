/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

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

package pkg

import (
	"regexp"
	"strconv"

	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"

	"strings"
	"time"
)

const (
	// PrimaryFilterAll basically means glob: *, but since this expands in the terminal we use @
	PrimaryFilterAll = "@"
)

type SliceCfg struct {
	From *int
	To   *int
}

func (s *SliceCfg) IsDefined() bool {
	return s != nil && (s.From != nil || s.To != nil)
}

func ParseSlice(s string) *SliceCfg {
	if strings.TrimSpace(s) == "" {
		return nil
	}

	sli := &SliceCfg{
		From: nil,
		To:   nil,
	}

	const (
		invalidSliceMsg = "the --slice flag or slice configuration is invalid; format is like a Go slice without the brackets"
	)

	s = strings.TrimSpace(s)

	colonIdx := strings.Index(s, ":")
	if colonIdx == -1 {
		log.Fatal(invalidSliceMsg)
	}

	fromNumStr := s[0:colonIdx]
	toNumStr := s[colonIdx+1:]

	if n, err := strconv.Atoi(strings.TrimSpace(fromNumStr)); err == nil {
		sli.From = &n
	}

	if n, err := strconv.Atoi(strings.TrimSpace(toNumStr)); err == nil {
		sli.To = &n
	}

	return sli
}

type TailscaleAPICfgCtx struct {
	Timeout time.Duration

	// ApiKey for regular authentication
	ApiKey string

	// OAuthClientID for OAuth based login.
	OAuthClientID string
	// OAuthClientSecret for Oauth based login.
	OAuthClientSecret string

	// ElapsedTime records the time this API call took. It's meant to be mutated during the API call and populated then.
	ElapsedTime time.Duration
}

type TailscaleCLICfgCtx struct {
}

type ConfigCtx struct {
	Columns       mapset.Set[string]
	Concurrency   int
	Filters       map[string]mapset.Set[string]
	IPsOutput     bool
	JsonOutput    bool
	NoCache       bool
	NoColor       bool
	PrimaryFilter *regexp.Regexp
	RemoteCmd     string
	Slice         *SliceCfg
	SortOrder     []SortSpec
	Tailnet       string
	TailscaleAPI  TailscaleAPICfgCtx
	TailscaleCLI  TailscaleCLICfgCtx

	TestMode bool
}

func NewConfigCtx() *ConfigCtx {
	return &ConfigCtx{}
}

func (c *ConfigCtx) IsRemoteCommand() bool {
	return len(c.RemoteCmd) > 0
}

func ParseColumns(s string) mapset.Set[string] {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}

	m := mapset.NewSet[string]()

	parts := strings.Split(s, ",")
	for _, p := range parts {
		m.Add(strings.ToLower(strings.TrimSpace(p)))
	}

	return m
}
