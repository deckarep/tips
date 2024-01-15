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
	"fmt"
	"time"

	"github.com/deckarep/tips/pkg/ui"

	"github.com/charmbracelet/log"
)

var (
	nowField   = fmt.Sprintf("%s now", ui.Styles.Green.Render(ui.Symbols.Dot))
	checkField = fmt.Sprintf("%s yes", ui.Styles.Green.Render(ui.Symbols.Checkmark))
)

// ProcessDevicesTable will apply sorting (if required), slicing (if required) and the massage/transformation of data to produce a final
// `*DevicesTable` that has everything required to render.
func ProcessDevicesTable(ctx context.Context, devList []*WrappedDevice) (*GeneralTableView, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	// Simple hack to determine if EnrichedInfo even exists.
	hasEnrichedInfo := len(devList) > 0 && devList[0].EnrichedInfo != nil

	// 1. Filter - if user requested any with the --filter flag
	var filteredDevList []*WrappedDevice
	if cfg.Filters != nil {
		filteredDevList = executeFilters(ctx, devList)
	} else {
		filteredDevList = devList
	}

	// 2. Sort - based on user's configured setting or --sort flag
	// If at least one dynamic sort was defined, then apply it.
	if len(cfg.SortOrder) > 0 {
		dynamicSortDevices(filteredDevList, cfg.SortOrder)
	}

	// 3. Slice - if provided via the --slice flag or configured, slice the results according to
	// Go's standard slicing convention.
	var slicedDevList = filteredDevList
	if cfg.Slice.IsDefined() {
		sliceWarnMsgFmt := "upper bound on slice: %d is larger than results len: %d"
		if cfg.Slice.From != nil && cfg.Slice.To != nil {
			// pin the upperbound to something reasonable.
			if *cfg.Slice.To > len(filteredDevList) {
				log.Warnf(sliceWarnMsgFmt, *cfg.Slice.To, len(filteredDevList))
				*cfg.Slice.To = len(filteredDevList)
			}
			slicedDevList = filteredDevList[*cfg.Slice.From:*cfg.Slice.To]
		} else if cfg.Slice.From != nil {
			slicedDevList = filteredDevList[*cfg.Slice.From:]
		} else {
			// pin the upperbound to something reasonable.
			if *cfg.Slice.To > len(filteredDevList) {
				log.Warnf(sliceWarnMsgFmt, *cfg.Slice.To, len(filteredDevList))
				*cfg.Slice.To = len(filteredDevList)
			}
			slicedDevList = filteredDevList[:*cfg.Slice.To]
		}
	}

	hdrs := getHeaders(ctx, hasEnrichedInfo)

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
		Headers: hdrs,
	}

	// Pre-alloc size.
	tbl.Rows = make([][]string, 0, len(slicedDevList))

	for idx, dev := range slicedDevList {
		if dev.EnrichedInfo != nil && dev.EnrichedInfo.IsSelf {
			tbl.Self = &SelfView{
				Index:   idx,
				DNSName: dev.Name,
			}
		}
		tbl.Rows = append(tbl.Rows, getRow(ctx, idx, hdrs, dev))
	}

	return tbl, nil
}

func getHeaders(ctx context.Context, hasEnrichedInfo bool) []Header {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	uniqHeaders := make(map[HeaderMatchName]Header)
	var headers []Header

	// addHdr ensures duplicate headers are not added needlessly.
	var addHdr = func(hdr Header) {
		// Only add what is uniq.
		if _, exists := uniqHeaders[hdr.MatchName]; !exists {
			headers = append(headers, hdr)

			// Track headers to ensure there are no duplicate columns.
			uniqHeaders[hdr.MatchName] = hdr
		}
	}

	for _, hdr := range DefaultColumnSet {
		if (hdr.ReqEnriched && !hasEnrichedInfo) || cfg.ColumnsExclude != nil && cfg.ColumnsExclude.Contains(string(hdr.MatchName)) {
			// 1. Exclude this header if it requires enriched data and we don't have it.
			// 2. Or when user requested to not include it.
			continue
		}

		addHdr(hdr)
	}

	// Any additionally requested headers append at the tail (for now).
	if cfg.Columns != nil {
		cfg.Columns.Each(func(matchName string) bool {
			if hdr, exists := AllHeadersMap[HeaderMatchName(matchName)]; exists {
				addHdr(hdr)
			}
			return false
		})
	}

	return headers
}

func getRow(ctx context.Context, idx int, headers []Header, d *WrappedDevice) []string {
	var results []string

	for _, hdr := range headers {
		results = append(results, d.EvalColumnField(ctx, idx, hdr.MatchName))
	}

	return results
}
