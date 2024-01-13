package pkg

import (
	"context"
	"testing"

	"github.com/deckarep/tips/pkg/slicecomp"

	"github.com/deckarep/tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"

	"github.com/stretchr/testify/assert"
)

func TestProcessDevicesTable(t *testing.T) {
	ctx := context.Background()
	cfgCtx := NewConfigCtx()

	// Apply sorting
	cfgCtx.SortOrder = ParseSortString("name:asc")
	// Apply slicing
	slice, err := slicecomp.ParseSlice("[0:1]", 0)
	assert.NoError(t, err)

	cfgCtx.Slice = slice

	ctx = context.WithValue(ctx, CtxKeyConfig, cfgCtx)
	ctx = context.WithValue(ctx, CtxKeyUserQuery, "*")

	devList := []*WrappedDevice{
		{
			Device: tailscale.Device{
				Addresses:     []string{"", "127.0.0.1"},
				Name:          "deadbeef",
				ID:            "deadbeef",
				User:          "user@foo.com",
				Tags:          []string{"tag:peanut", "tag:walnut"},
				ClientVersion: "1.22.3",
				Hostname:      "deadbeef",
				NodeKey:       "deadbeef",
				OS:            "linux",
			},
			EnrichedInfo: &tailscale_cli.DeviceInfo{
				NodeKey:           "deadbeef",
				DNSName:           "",
				HasExitNodeOption: false,
				IsSelf:            true,
				Online:            true,
			},
		},
		{
			Device: tailscale.Device{
				Addresses:     []string{"", "127.0.0.2"},
				Name:          "badbeef",
				ID:            "badbeef",
				User:          "user@foo.com",
				Tags:          []string{"tag:peanut"},
				ClientVersion: "1.22.3",
				Hostname:      "badbeef",
				NodeKey:       "badbeef",
				OS:            "windows",
			},
			EnrichedInfo: &tailscale_cli.DeviceInfo{
				NodeKey:           "badbeef",
				DNSName:           "",
				HasExitNodeOption: false,
				IsSelf:            true,
				Online:            true,
			},
		},
	}

	// Single slice - 0:1
	tv, err := ProcessDevicesTable(ctx, devList)
	assert.NoError(t, err)

	assert.Equal(t, len(tv.Rows), 1, "the general table view should have a single row")

	// slice from 0:(everything else)
	slice, err = slicecomp.ParseSlice("[0:]", 0)
	assert.NoError(t, err)

	cfgCtx.Slice = slice
	tv, err = ProcessDevicesTable(ctx, devList)
	assert.NoError(t, err)

	assert.Equal(t, len(tv.Rows), 2, "the general table view should have a single row")

	// slice from :1
	slice, err = slicecomp.ParseSlice("[:1]", 0)
	assert.NoError(t, err)

	cfgCtx.Slice = slice
	tv, err = ProcessDevicesTable(ctx, devList)
	assert.NoError(t, err)

	assert.Equal(t, len(tv.Rows), 1, "the general table view should have a single row")

	// slice from 0:50 - overly large slice.
	slice, err = slicecomp.ParseSlice("[0:50]", 0)
	assert.NoError(t, err)
	cfgCtx.Slice = slice
	tv, err = ProcessDevicesTable(ctx, devList)
	assert.NoError(t, err)

	assert.Equal(t, len(tv.Rows), 2, "the general table view should have a single row")
}
