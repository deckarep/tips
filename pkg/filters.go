/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

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

	"strings"

	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"
)

// Formats we will take (whitespace doesn't matter:
// Single filters
// 'tag:peanuts'
// 'tag:peanuts,nuts'
// Multi filters
// 'tag:peanuts, version:1.54.1'
// Future: not based filter?
// 'tag:peanuts, !walnuts (will exclude walnuts)

var (
	filterTagsSupported = mapset.NewSet[string](
		"exit",
		"ipv4",
		"ipv6",
		"os",
		"tag",
		"user",
		"version",
	)
)

func ParseFilter(filter string) map[string]mapset.Set[string] {
	// 0. If no tags, just return.
	filter = strings.TrimSpace(filter)
	m := make(map[string]mapset.Set[string])

	if len(filter) == 0 {
		return m
	}

	// 0. Validate proper use of filters.
	filter = strings.ToLower(filter)
	if !strings.Contains(filter, ":") {
		log.Fatal("a --filter must specify at least one tag:", "tags_supported", filterTagsSupported.String())
	}

	// 1. Parse filters into sets.
	filterParts := strings.Split(filter, ",")
	var activeFilter string
	for _, part := range filterParts {
		part = strings.TrimSpace(part)

		// Skip any unusable whitespace.
		if len(part) == 0 {
			continue
		}

		// If this is a tag, create a Set for it.
		colorIdx := strings.Index(part, ":")
		if colorIdx > -1 {
			tagType := strings.TrimSpace(part[0:colorIdx])
			if filterTagsSupported.Contains(tagType) {
				part = strings.TrimSpace(part)
				activeFilter = tagType
				m[activeFilter] = mapset.NewSet[string]()
			} else {
				log.Fatal("--filter with tag is unsupported", "tag", tagType)
			}
		}

		if colorIdx > -1 {
			// Ensure if a filter prefix is present, remove it ie 'os:'
			part = part[colorIdx+1:]
		}
		part = strings.TrimSpace(part)
		if _, exists := m[activeFilter]; !exists {
			log.Fatalf("logical error because no val exists for key: %q", activeFilter)
		}
		m[activeFilter].Add(part)
	}

	return m
}

func executeFilters(ctx context.Context, devList []*WrappedDevice) []*WrappedDevice {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	var normalizeTags = func(vals []string) []string {
		var items []string
		for _, s := range vals {
			items = append(items, strings.ToLower(strings.Replace(s, "tag:", "", -1)))
		}
		return items
	}

	var (
		filteredDevList []*WrappedDevice
	)

	// Note: for better performance, filtering should be done by the most selective fields first.
	for _, dev := range devList {

		// Filter by 'tag' when provided - currently only supports full matching.
		if f, exists := cfg.Filters["tag"]; exists {
			normalizedTags := normalizeTags(dev.Tags)
			wantsNoTags := f.Contains("nil") // User wants to filter out rows with tags.

			// Determine if the device should be skipped based on tag presence and user's filter.
			hasTags := len(dev.Tags) > 0
			matchesTags := f.ContainsAny(normalizedTags...)

			// Skip device if it doesn't match the filter criteria.
			if (wantsNoTags && hasTags) || (!wantsNoTags && !matchesTags) {
				continue
			}
		}

		// Filter by 'user' when provided - currently only supports full matching.
		if f, exists := cfg.Filters["user"]; exists {
			if !f.Contains(strings.ToLower(dev.User)) {
				continue
			}
		}

		// Filter by 'version' when provided - currently only supports full matching.
		if f, exists := cfg.Filters["version"]; exists {
			// For now, just filter on the first portion of the version which has the format: 1.xx.1
			semanticVersion := strings.Split(dev.ClientVersion, "-")[0]
			if !f.Contains(strings.ToLower(semanticVersion)) {
				continue
			}
		}

		// Filters by ipv4 - currently only supports full matching.
		if f, exists := cfg.Filters["ipv4"]; exists {
			if (len(dev.Addresses) == 0) || !f.ContainsAny(dev.Addresses...) {
				continue
			}
		}

		// Filters by ipv6 - currently only supports full matching.
		if f, exists := cfg.Filters["ipv6"]; exists {
			if (len(dev.Addresses) == 0) || !f.ContainsAny(dev.Addresses...) {
				continue
			}
		}

		// TODO: additional filters - lastSeen

		// Filter by 'os' when provided - currently only supports full matching.
		if f, exists := cfg.Filters["os"]; exists {
			if !f.Contains(strings.ToLower(dev.OS)) {
				continue
			}
		}

		// Filters by exit node - filter query looks like: --filter 'exit:yes|no'
		if f, exists := cfg.Filters["exit"]; exists {
			// NOTE: this information only exists from enriched results.
			if enrichedDev := dev.EnrichedInfo; enrichedDev != nil {
				if (f.Contains("yes") && !enrichedDev.HasExitNodeOption) ||
					(f.Contains("no") && enrichedDev.HasExitNodeOption) {
					continue
				}
			}
		}

		filteredDevList = append(filteredDevList, dev)
	}

	return filteredDevList
}
