package app

import (
	"sort"
	"strings"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

type SortSpec struct {
	Field     string
	Direction SortDirection
}

// Parse the sort string and return a slice of SortSpec
func ParseSortString(sortString string) []SortSpec {
	var specs []SortSpec
	fields := strings.Split(sortString, ",")
	for _, field := range fields {
		parts := strings.Split(strings.ToLower(strings.TrimSpace(field)), ":")
		if len(parts) != 2 {
			continue // or handle the error
		}
		direction := Ascending
		if parts[1] == "dsc" {
			direction = Descending
		}
		specs = append(specs, SortSpec{Field: strings.ToUpper(parts[0]), Direction: direction})
	}
	return specs
}

func dynamicSortDevices(slice []tailscale.Device, specs []SortSpec) {
	sort.Slice(slice, func(i, j int) bool {
		for _, spec := range specs {
			switch spec.Field {
			case "MACHINE":
				if spec.Direction == Ascending {
					if slice[i].Name != slice[j].Name {
						return slice[i].Name < slice[j].Name
					}
				} else {
					if slice[i].Name != slice[j].Name {
						return slice[i].Name > slice[j].Name
					}
				}
			case "USER":
				if spec.Direction == Ascending {
					if slice[i].User != slice[j].User {
						return slice[i].User < slice[j].User
					}
				} else {
					if slice[i].User != slice[j].User {
						return slice[i].User > slice[j].User
					}
				}
				// Add cases for other fields...
			}
		}
		return false
	})
}
