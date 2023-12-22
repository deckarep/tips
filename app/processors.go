package app

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tips/pkg/tailscale_cli"
	"tips/pkg/ui"

	"github.com/dustin/go-humanize"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

const (
	NoHdr       = "No"
	MachineHdr  = "Machine"
	AddressHdr  = "Address"
	TagsHdr     = "Tags"
	UserHdr     = "User"
	VersionHdr  = "Version"
	LastSeenHdr = "LastSeen"
	ExitNodeHdr = "Exit Node"
)

var (
	DefaultHeader = []string{NoHdr, MachineHdr, AddressHdr, TagsHdr, UserHdr, VersionHdr, ExitNodeHdr, LastSeenHdr}
	nowField      = fmt.Sprintf("%s now", ui.Styles.Green.Render(ui.Dot))
	checkField    = fmt.Sprintf("%s yes", ui.Styles.Green.Render(ui.Checkmark))
)

// ProcessDevicesTable will apply sorting (if required), slicing (if required) and the massage/transformation of data to produce a final
// `*DevicesTable` that has everything required to render.
func ProcessDevicesTable(ctx context.Context, devList []tailscale.Device,
	devEnriched map[string]tailscale_cli.DeviceInfo) (*GeneralTableView, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	// 1. Filter - if user requested any with the --filter flag
	filteredDevList := executeFilters(ctx, devList, devEnriched)

	// 2. Sort - based on user's configured setting or --sort flag
	// If at least one dynamic sort was defined, then apply it.
	if len(cfg.SortOrder) > 0 {
		dynamicSortDevices(filteredDevList, cfg.SortOrder)
	}

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
			CLIElapsed: time.Second * 3, //TODO measure this time.
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
	// TODO: I need to remove columns based on --columns flag.

	var (
		hasExitNodeOption = "no"
		version           = fmt.Sprintf("%s - %s", strings.Split(d.ClientVersion, "-")[0], d.OS)

		timeAgo = humanize.Time(d.LastSeen.Time)
		// Remove all tag: prefixes, and join the tags as a comma delimited string.
		tags = strings.Replace(strings.Join(d.Tags, ", "), "tag:", "", -1)

		// d.Name is the fully qualified DNS name, but we just shorten it and this is the name used
		// that takes precedence when the user overrides the name.
		easyName = strings.Split(d.Name, ".")[0]
	)

	seenAgo := timeAgo
	//if strings.Contains(seenAgo, "seconds") {
	//	// https://github.com/tailscale/tailscale/pull/3534/files
	//	seenAgo = fmt.Sprintf("• %s", seenAgo)
	//}

	num := fmt.Sprintf("%04d", idx)

	// Enriched results are only available when run from a node that is in the tailnet cluster itself.
	// In other words, the external Tailscale API only provides so much data.

	if enrichedDev, ok := enrichedResults[d.NodeKey]; ok {
		if enrichedDev.Online {
			seenAgo = nowField
		}
		if enrichedDev.HasExitNodeOption {
			hasExitNodeOption = checkField
		}
	}

	return []string{num, easyName, d.Addresses[0], tags, d.User, version, hasExitNodeOption, seenAgo}
}
