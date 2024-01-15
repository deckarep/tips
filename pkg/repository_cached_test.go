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
	"context"
	"testing"
	"time"

	"github.com/tailscale/tailscale-client-go/tailscale"

	"github.com/deckarep/tips/pkg/prefixcomp"

	"github.com/stretchr/testify/assert"
)

type fakeDeviceRepo struct {
	funcToCall func(ctx context.Context) ([]*WrappedDevice, error)
}

func (f *fakeDeviceRepo) DevicesResource(ctx context.Context) ([]*WrappedDevice, error) {
	return f.funcToCall(ctx)
}

func TestCachedRepository_DevicesResource(t *testing.T) {
	var remoteCalls int
	var devicesCall = func(ctx context.Context) ([]*WrappedDevice, error) {

		remoteCalls += 1

		return []*WrappedDevice{
			{Device: tailscale.Device{Name: "foo"}},
			{Device: tailscale.Device{Name: "bar"}},
			{Device: tailscale.Device{Name: "baz"}},
		}, nil
	}

	ctx := context.Background()
	cfg := NewConfigCtx()

	// The test should clean up this file.
	const (
		testTailnet = "test@test.com"
	)
	cfg.Tailnet = testTailnet
	// Should be enough time for our test to work correctly.
	cfg.CacheTimeout = time.Minute * 15
	prefAST, err := prefixcomp.ParsePrimaryFilter("*")
	assert.NoError(t, err)
	cfg.PrefixFilter = prefAST

	// Wrap the config in a bow (in the context)
	ctx = context.WithValue(ctx, CtxKeyConfig, cfg)

	// Cleanup: We're creating a quick db object to simply erase the test file upon the completion of this test run.
	// Instantiation doesn't open any files or anything.
	defer func() {
		fakeDB := NewDB2[*WrappedDevice](testTailnet)
		err := fakeDB.Erase()
		assert.NoError(t, err)
	}()

	cachedRepo := NewCachedRepo(&fakeDeviceRepo{funcToCall: devicesCall})
	devs, err := cachedRepo.DevicesResource(ctx)
	assert.NoError(t, err)

	// Check we have 3 items.
	assert.Equal(t, 3, len(devs))
	// Check that the remote call was made on the first run!
	// This is important, it means cache was not utilized.
	assert.Equal(t, remoteCalls, 1)

	// Next invocation, invokes the cache, and expects no remote calls to occur.
	prefAST, err = prefixcomp.ParsePrimaryFilter("foo|bar")
	cfg.PrefixFilter = prefAST
	assert.NoError(t, err)
	devs, err = cachedRepo.DevicesResource(ctx)
	assert.NoError(t, err)
	// Another 1 implies the cache was utilized.
	assert.Equal(t, 1, remoteCalls)
	// We did an OR search only 2 items should return.
	assert.Equal(t, 2, len(devs))

	// Last invocation, invokes the cache, and expects no remote calls to occur.
	prefAST, err = prefixcomp.ParsePrimaryFilter("baz")
	cfg.PrefixFilter = prefAST
	assert.NoError(t, err)
	devs, err = cachedRepo.DevicesResource(ctx)
	assert.NoError(t, err)
	// Another 1 implies the cache was utilized.
	assert.Equal(t, 1, remoteCalls)
	//We did a single search, only 1 item should return.
	assert.Equal(t, 1, len(devs))
}
