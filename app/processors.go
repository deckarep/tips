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

const (
	NoHdr       = "No"
	MachineHdr  = "Machine"
	AddressHdr  = "Address"
	TagsHdr     = "Tags"
	UserHdr     = "User"
	VersionHdr  = "Version"
	LastSeenHdr = "LastSeen"
)

var (
	DefaultHeader = []string{NoHdr, MachineHdr, AddressHdr, TagsHdr, UserHdr, VersionHdr, LastSeenHdr}
)

// ProcessDevicesTable will apply sorting (if required), slicing (if required) and the massage/transformation of data to produce a final
// `*DevicesTable` that has everything required to render.
func ProcessDevicesTable(ctx context.Context, devList []tailscale.Device, devEnriched map[string]tailscale_cli.DeviceInfo) (*GeneralTableView, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	// 1. Filter -  if user requested any with the --filter flag

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
	var primFilter = strings.ToLower(cfg.PrimaryFilter)
	for _, dev := range devList {
		// PrimaryFilter when not '*' (everything), by default this is a case-insensitive prefix filter.
		if cfg.PrimaryFilter != PrimaryFilterAll &&
			!strings.HasPrefix(strings.ToLower(dev.Hostname), primFilter) {
			continue
		}

		// Filter by 'os' when provided.
		if f, exists := cfg.Filters["os"]; exists {
			if !f.Contains(strings.ToLower(dev.OS)) {
				continue
			}
		}

		// Filter by 'tag' when provided.
		if f, exists := cfg.Filters["tag"]; exists {
			normalizedTags := normalizeTags(dev.Tags)
			//spew.Dump(normalizedTags, f)
			if (len(dev.Tags) == 0) || !f.Contains(normalizedTags...) {
				continue
			}
		}

		// Filter by 'user' when provided.
		if f, exists := cfg.Filters["user"]; exists {
			if !f.Contains(strings.ToLower(dev.User)) {
				continue
			}
		}

		// Filter by 'version' when provided.
		if f, exists := cfg.Filters["version"]; exists {
			// For now, just filter on the first portion of the version which has the format: 1.xx.1
			semanticVersion := strings.Split(dev.ClientVersion, "-")[0]
			if !f.Contains(strings.ToLower(semanticVersion)) {
				continue
			}
		}

		filteredDevList = append(filteredDevList, dev)
	}

	// 2. Sort - based on user's configured setting or --sort flag
	slices.SortFunc(filteredDevList, func(a, b tailscale.Device) int {
		// TODO: Must be able to do this from configuration logic from context.
		if n := cmp.Compare(a.Name, b.Name); n != 0 {
			return n
		}
		return cmp.Compare(a.Name, b.Name)
	})

	// 3. Slice - if provided via the --slice flag or configured, slice the results according to
	// Go's standard slicing convention.
	var slicedDevList = filteredDevList
	if cfg.Slice.IsDefined() {
		if cfg.Slice.From != nil && cfg.Slice.To != nil {
			slicedDevList = filteredDevList[*cfg.Slice.From:*cfg.Slice.To]
		} else if cfg.Slice.From != nil {
			slicedDevList = filteredDevList[*cfg.Slice.From:]
		} else {
			slicedDevList = filteredDevList[:*cfg.Slice.To]
		}
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
	tbl.Rows = make([][]string, 0, len(slicedDevList))

	for idx, dev := range slicedDevList {
		tbl.Rows = append(tbl.Rows, getRow(idx, dev, devEnriched))
	}

	return tbl, nil
}

func getHeaders(enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
	// TODO: I need to remove columns based on --columns flag.

	// Currently there's no difference on the headers returned when enrichedResults are present.
	if len(enrichedResults) > 0 {
		return DefaultHeader
	}
	return DefaultHeader
}

func getRow(idx int, d tailscale.Device, enrichedResults map[string]tailscale_cli.DeviceInfo) []string {
	// You can also add tables row-by-row

	// TODO: I need to remove columns based on --columns flag.

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
