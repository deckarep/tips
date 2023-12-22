package app

import (
	"context"
	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"strings"
	"tips/pkg/tailscale_cli"
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

func executeFilters(ctx context.Context, devList []tailscale.Device,
	devEnriched map[string]tailscale_cli.DeviceInfo) []tailscale.Device {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	var normalizeTags = func(vals []string) []string {
		var items []string
		for _, s := range vals {
			items = append(items, strings.ToLower(strings.Replace(s, "tag:", "", -1)))
		}
		return items
	}

	var (
		filteredDevList []tailscale.Device
	)

	// Note: for better performance, filtering should be done by the most selective fields first.
	for _, dev := range devList {
		// d.Name is the fully qualified DNS name, but we just shorten it and this is the name used
		// that takes precedence when the user overrides the name.
		easyName := strings.Split(dev.Name, ".")[0]

		// PrimaryFilter is a regex and applied when non-nil.
		if cfg.PrimaryFilter != nil && !cfg.PrimaryFilter.MatchString(easyName) {
			continue
		}

		// Filter by 'tag' when provided - currently only supports full matching.
		if f, exists := cfg.Filters["tag"]; exists {
			normalizedTags := normalizeTags(dev.Tags)

			// If the user does a filter like: 'tag:nil' they want to filter out those rows WITH tags.
			wantsEmpty := f.Contains("nil")

			if !wantsEmpty && (len(dev.Tags) == 0) || !f.Contains(normalizedTags...) {
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
			if enrichedDev, ok := devEnriched[dev.NodeKey]; ok {
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
