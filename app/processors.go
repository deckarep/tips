package app

import (
	"cmp"
	"context"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"slices"
	"tips/pkg/tailscale_cli"
)

// ProcessDevicesTable will apply sorting (if required), slicing (if required) and the massage/transformation of data to produce a final
// `*DevicesTable` that has everything required to render.
func ProcessDevicesTable(ctx context.Context, devList []tailscale.Device, devEnriched map[string]tailscale_cli.DeviceInfo) (*GeneralTableView, error) {
	// 1. Sort - here we'll sort based on user's configured setting.
	slices.SortFunc(devList, func(a, b tailscale.Device) int {
		// TODO: Must be able to do this from configuration logic from context.
		if n := cmp.Compare(a.Name, b.Name); n != 0 {
			return n
		}
		return cmp.Compare(a.Name, b.Name)
	})

	// 2. Slice - then slice what gets returned
	// TODO: must be able to slice from user's configured setting.
	devList = devList[0:3]

	// 3. Massage - final transformations here.
	return &GeneralTableView{
		TailnetView: TailnetView{
			Tailnet: "deckarep@gmail.com",
		},
		Headers: []string{"No", "Bar", "Baz"},
		Rows: [][]string{
			{"0", "foo1", "foo2"},
		},
	}, nil
}
