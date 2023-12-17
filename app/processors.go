package app

import (
	"cmp"
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"slices"
	"strconv"
	"strings"
	"time"
	"tips/pkg/tailscale_cli"
	"tips/pkg/ui"
)

// ProcessDevicesTable will apply sorting (if required), slicing (if required) and the massage/transformation of data to produce a final
// `*DevicesTable` that has everything required to render.
func ProcessDevicesTable(ctx context.Context, devList []tailscale.Device, devEnriched map[string]tailscale_cli.DeviceInfo) (*GeneralTableView, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

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
	slicedDevList := devList[0:3]

	// 3. Filter - such as by tag.
	// apply filtering logic from user's configured setting.
	//var strSlicerLower = func(vals []string) []string {
	//	var items []string
	//	for _, s := range vals {
	//		items = append(items, strings.ToLower(s))
	//	}
	//	return items
	//}
	var normalizeTags = func(vals []string) []string {
		var items []string
		for _, s := range vals {
			items = append(items, strings.ToLower(strings.Replace(s, "tag:", "", -1)))
		}
		return items
	}
	var filteredDevList []tailscale.Device
	for _, dev := range slicedDevList {
		// Filter by 'os' if provided.
		if f, exists := cfg.Filters["os"]; exists {
			if !f.Contains(strings.ToLower(dev.OS)) {
				continue
			}
		}

		// Filter by 'tag' if provided.
		if f, exists := cfg.Filters["tag"]; exists {
			normalizedTags := normalizeTags(dev.Tags)
			//spew.Dump(normalizedTags, f)
			if (len(dev.Tags) == 0) || !f.Contains(normalizedTags...) {
				continue
			}
		}

		// Filter by 'user' if provided.
		if f, exists := cfg.Filters["user"]; exists {
			if !f.Contains(strings.ToLower(dev.User)) {
				continue
			}
		}

		filteredDevList = append(filteredDevList, dev)
	}

	// 3. Massage/Transform - final transformations here.
	tbl := &GeneralTableView{
		ContextView: ContextView{
			Query:      CtxAsString(ctx, CtxKeyUserQuery),
			APIElapsed: cfg.TailscaleAPI.ElapsedTime,
			CLIElapsed: time.Duration(time.Second * 3), //TODO measure this time.
		},
		TailnetView: TailnetView{
			Tailnet: cfg.Tailnet,
			// This should represent the size of the entire result-set.
			TotalMachines: len(devList),
		},
		SelfView: SelfView{
			Index:   0,
			DNSName: "foo.bar.3234.dns.name.",
		},
		Headers: getHeaders(devEnriched),
	}

	// Pre-alloc size.
	tbl.Rows = make([][]string, 0, len(filteredDevList))

	for idx, dev := range filteredDevList {
		tbl.Rows = append(tbl.Rows, getRow(idx, dev, devEnriched))
	}

	return tbl, nil
}

func getHeaders(enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
	if len(enrichedResults) > 0 {
		return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
	}
	return []string{"No", "Machine", "Address", "Tags", "User", "Version", "LastSeen"}
}

func getRow(idx int, d tailscale.Device, enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
	// You can also add tables row-by-row

	var (
		version = fmt.Sprintf("%s - %s", strings.Split(d.ClientVersion, "-")[0], d.OS)

		timeAgo = humanize.Time(d.LastSeen.Time)
		// Remove all tag: prefixes, and join the tags as a comma delimited string.
		tags = strings.Replace(strings.Join(d.Tags, ","), "tag:", "", -1)
	)

	seenAgo := timeAgo
	//if strings.Contains(seenAgo, "seconds") {
	//	// https://github.com/tailscale/tailscale/pull/3534/files
	//	seenAgo = fmt.Sprintf("• %s", seenAgo)
	//}

	if len(enrichedResults) > 0 {
		if enrichedDev, ok := enrichedResults[d.NodeKey]; ok && enrichedDev.Online {
			seenAgo = fmt.Sprintf("%s now", ui.Styles.Green.Render(ui.Dot))
		}
		return []string{strconv.Itoa(idx), d.Hostname, d.Addresses[0], tags, d.User, version, seenAgo}
	}
	return []string{strconv.Itoa(idx), d.Hostname, d.Addresses[0], tags, d.User, version, seenAgo}
}
