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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func TestParseSortString(t *testing.T) {
	parsedResult := ParseSortString("name:dsc, foo:asc,bar:asc")

	if parsedResult == nil {
		t.Error("expected non-nil parsedResult")
	}

	if len(parsedResult) != 3 {
		t.Errorf("expected 2 results, got: %d", len(parsedResult))
	}

	expectedResults := []struct {
		field string
		dir   SortDirection
	}{
		{field: "NAME", dir: Descending},
		{field: "FOO", dir: Ascending},
		{field: "BAR", dir: Ascending},
	}

	for i, er := range expectedResults {
		if parsedResult[i].Field != er.field || parsedResult[i].Direction != er.dir {
			t.Errorf("expected at idx: %d field: %s dir: %d, got field: %s dir: %d", i,
				er.field, er.dir, parsedResult[i].Field, parsedResult[i].Direction)
		}
	}
}

func TestDynamicSortDevices(t *testing.T) {
	inputDevs := []*WrappedDevice{
		{
			Device: tailscale.Device{
				Name: "b-foo",
				User: "chestnut@foo.com",
			},
			EnrichedInfo: nil,
		},
		{
			Device: tailscale.Device{
				Name: "a-foo",
				User: "peanut@foo.com",
			},
			EnrichedInfo: nil,
		},
		{
			Device: tailscale.Device{
				Name: "c-foo",
				User: "walnut@foo.com",
			},
			EnrichedInfo: nil,
		},
	}

	specs := []SortSpec{{
		Field:     "NAME",
		Direction: Ascending,
	}}

	dynamicSortDevices(inputDevs, specs)

	assert.Equal(t, inputDevs[0].Name, "a-foo")
	assert.Equal(t, inputDevs[1].Name, "b-foo")
	assert.Equal(t, inputDevs[2].Name, "c-foo")

	specs = []SortSpec{{
		Field:     "USER",
		Direction: Descending,
	}}

	dynamicSortDevices(inputDevs, specs)

	assert.Equal(t, inputDevs[0].User, "walnut@foo.com")
	assert.Equal(t, inputDevs[1].User, "peanut@foo.com")
	assert.Equal(t, inputDevs[2].User, "chestnut@foo.com")
}
