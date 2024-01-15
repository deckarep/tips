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
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRenderRemoteSummary(t *testing.T) {
	ctx := context.Background()

	var b bytes.Buffer
	err := RenderRemoteSummary(ctx, &b, 2, 0, time.Millisecond*333)
	assert.NoError(t, err, "RenderRemoteSummary should have returned no error")

	assert.Equal(t, b.String(), "Finished: successes: 2, failures: 0, elapsed (secs): 0.33\n")

	b.Reset()
	err = RenderRemoteSummary(ctx, &b, 0, 3, time.Millisecond*777)
	assert.NoError(t, err, "RenderRemoteSummary should have returned no error")

	assert.Equal(t, b.String(), "Finished: successes: 0, failures: 3, elapsed (secs): 0.78\n")
}

func TestRenderIPs(t *testing.T) {
	var b bytes.Buffer
	ctx := context.Background()
	cfgCtx := NewConfigCtx()
	cfgCtx.IPsDelimiter = "\n"
	ctx = context.WithValue(ctx, CtxKeyConfig, cfgCtx)

	tv := &GeneralTableView{
		ContextView: ContextView{},
		TailnetView: TailnetView{},
		Self:        &SelfView{},
		Rows: [][]string{
			{"foo", "bar", "127.0.0.1"},
			{"coo", "car", "127.0.0.2"},
			{"soo", "sar", "127.0.0.3"},
			{"too", "tar", "127.0.0.4"},
		},
	}

	err := RenderIPs(ctx, tv, &b)
	assert.NoError(t, err, "RenderIPs should have returned no error")
	assert.Equal(t, b.String(), "127.0.0.1\n127.0.0.2\n127.0.0.3\n127.0.0.4\n")
}

func TestRenderJson(t *testing.T) {
	var b bytes.Buffer
	ctx := context.Background()
	cfgCtx := NewConfigCtx()
	ctx = context.WithValue(ctx, CtxKeyConfig, cfgCtx)

	tv := &GeneralTableView{
		ContextView: ContextView{},
		TailnetView: TailnetView{},
		Self:        &SelfView{},
		Headers: []Header{
			HdrNo, HdrMachine, HdrIpv4,
		},
		Rows: [][]string{
			{"foo", "bar", "127.0.0.1"},
			{"coo", "car", "127.0.0.2"},
			{"soo", "sar", "127.0.0.3"},
			{"too", "tar", "127.0.0.4"},
		},
	}

	err := RenderJson(ctx, tv, &b)
	// TODO: At a minimum, assert no error but I don't want to test JSON encoding really, its already well tested.
	// TODO: But this test should be a little more flushed out.
	assert.NoError(t, err, "RenderJson should have returned no error")
	// assert.Equal(t, b.String(), "blah, blah")
}

func TestRenderLogLine(t *testing.T) {
	var b bytes.Buffer
	ctx := context.Background()
	cfgCtx := NewConfigCtx()
	ctx = context.WithValue(ctx, CtxKeyConfig, cfgCtx)

	scroll := []string{
		"restarting server...",
		"file not found: foo.txt",
		"hello world!",
	}

	// Just simulate a few lines being scrolled on by.
	for i := 0; i < 3; i++ {
		RenderLogLine(ctx, &b, i, false, "blade", "dinky", scroll[i])
	}

	assert.Equal(t, b.String(),
		"dinky >1 (0): restarting server...\ndinky >1 (1): file not found: foo.txt\ndinky >1 (2): hello world!\n")
}
