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
	"sort"
	"strings"
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

func dynamicSortDevices(slice []*WrappedDevice, specs []SortSpec) {
	sort.SliceStable(slice, func(i, j int) bool {
		for _, spec := range specs {
			switch spec.Field {
			case "MACHINE", "NAME":
				if spec.Direction == Ascending {
					if slice[i].Name != slice[j].Name {
						return slice[i].Name < slice[j].Name
					}
				} else {
					if slice[i].Name != slice[j].Name {
						return slice[i].Name > slice[j].Name
					}
				}
			case "USER", "EMAIL":
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
