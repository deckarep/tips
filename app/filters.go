package app

import (
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
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
		"tag",
		"os",
		"user",
		"version",
		"ipv4",
		"ipv6")
)

func ApplyFilter(filter string) map[string]mapset.Set[string] {
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

	spew.Dump(m)
	return m
}
