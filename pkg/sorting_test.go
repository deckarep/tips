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
			},
			EnrichedInfo: nil,
		},
		{
			Device: tailscale.Device{
				Name: "a-foo",
			},
			EnrichedInfo: nil,
		},
		{
			Device: tailscale.Device{
				Name: "c-foo",
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
		Field:     "NAME",
		Direction: Descending,
	}}

	dynamicSortDevices(inputDevs, specs)

	assert.Equal(t, inputDevs[0].Name, "c-foo")
	assert.Equal(t, inputDevs[1].Name, "b-foo")
	assert.Equal(t, inputDevs[2].Name, "a-foo")
}
