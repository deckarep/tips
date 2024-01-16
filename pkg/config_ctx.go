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

package pkg

import (
	"github.com/deckarep/tips/pkg/prefixcomp"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/deckarep/tips/pkg/filtercomp"
	"github.com/deckarep/tips/pkg/slicecomp"

	"strings"
	"time"
)

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
	Basic          bool
	CacheTimeout   time.Duration
	Columns        mapset.Set[string]
	ColumnsExclude mapset.Set[string]
	Concurrency    int
	Filters        filtercomp.AST
	IPsOutput      bool
	IPsDelimiter   string
	JsonOutput     bool
	NoCache        bool
	NoColor        bool
	PrefixFilter   *prefixcomp.PrimaryFilterAST
	RemoteCmd      string
	Slice          *slicecomp.Slice
	SortOrder      []SortSpec
	Stderr         bool
	Tailnet        string
	CachedElapsed  time.Duration
	TailscaleAPI   TailscaleAPICfgCtx
	TailscaleCLI   TailscaleCLICfgCtx
	Page           int

	TestMode bool
}

func NewConfigCtx() *ConfigCtx {
	return &ConfigCtx{}
}

func (c *ConfigCtx) IsRemoteCommand() bool {
	return len(c.RemoteCmd) > 0
}

func ParseColumns(s string) (mapset.Set[string], mapset.Set[string]) {
	i := mapset.NewSet[string]()
	e := mapset.NewSet[string]()

	if len(strings.TrimSpace(s)) == 0 {
		return i, e
	}

	parts := strings.Split(s, ",")
	for _, p := range parts {
		colName := strings.ToLower(strings.TrimSpace(p))
		if strings.HasPrefix(colName, "-") {
			e.Add(colName[1:])
		} else {
			i.Add(colName)
		}
	}

	return i, e
}
