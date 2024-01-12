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
	"strings"

	"github.com/charmbracelet/log"

	"github.com/deckarep/tips/pkg/filtercomp"

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
// filterTagsSupported = mapset.NewSet[string](
//
//	"exit",
//	"ipv4",
//	"ipv6",
//	"os",
//	"tag",
//	"user",
//	"version",
//
// )
)

// tag:peanut, walnut =>
// tag:peanut, chestnut | walnut
// tag:peanut, (walnut | chestnut)
// ((tag:peanut, tag:walnut, tag:cachew) | chestnut)
// (((tag:peanut, tag:walnut) | tag:chestnut) | user:dude@dude.com))

//type KindType int
//
//const (
//	KindTag  KindType = iota
//	KindUser KindType = iota
//)
//
//type Exp interface {
//	Eval(device *WrappedDevice) bool
//}
//
//type TerminalFilterExp struct {
//	kind KindType
//	val  string
//}
//
//func (e *TerminalFilterExp) Eval(device *WrappedDevice) bool {
//	switch e.kind {
//	case KindTag:
//		// TODO: Creates a set for every eval...but we'll optimize later by memoizing or something.
//		tagSet := mapset.NewSet[string](device.Tags...)
//		return tagSet.ContainsOne(e.val)
//	case KindUser:
//		return device.User == e.val
//	default:
//		panic("unknown KindType")
//	}
//}
//
//type FilterCompExp struct {
//	KindComp string
//	val      time.Time
//	left     Exp
//	right    Exp
//}
//
//func (e *FilterCompExp) Eval(device *WrappedDevice) bool {
//	// TODO: flesh this out depending on KindComp and value type.
//	return device.LastSeen.Second() > e.val.Second()
//}
//
//type FilterOrExp struct {
//	left  Exp
//	right Exp
//}
//
//func (e *FilterOrExp) Eval(device *WrappedDevice) bool {
//	return e.left.Eval(device) || e.right.Eval(device)
//}
//
//type FilterAndExp struct {
//	left  Exp
//	right Exp
//}
//
//func (e *FilterAndExp) Eval(device *WrappedDevice) bool {
//	return e.left.Eval(device) && e.right.Eval(device)
//}

//func line(depth int, line string) {
//	prefix := strings.Repeat(" ", depth)
//	println(prefix + line)
//}

func ParseFilter(filter string) (filtercomp.AST, error) {
	tokens := filtercomp.Tokenize([]byte(filter))
	if len(tokens) == 0 {
		return nil, nil
	}

	filterParser := filtercomp.NewParser(tokens)
	ast, err := filterParser.Parse()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

/*
func ParseFilter2(depth int, filter string) (Exp, int) {
	fmt.Println("ParseFilter2 =>", filter)
	var currentNode Exp

	for i := 0; i < len(filter); i++ {
		b := filter[i]
		switch b {
		case '(':
			line(depth, "open paren")
			depth++
			f, consumed := ParseFilter2(depth, filter[i+1:])
			return f, i + consumed
		case ')':
			depth--
			line(depth, "close paren")
			return currentNode, i + 1 //adding 1 because ) is already consumed.
		case '|':
			line(depth, "OR found")
			//left := currentNode
			var right Exp
			var consumed int
			right, consumed = ParseFilter2(depth, filter[i+1:])
			i += (consumed)

			left := currentNode
			currentNode = &FilterOrExp{
				left:  left,
				right: right,
			}

		case ',':
			line(depth, "AND found")
			//left := currentNode
			var right Exp
			var consumed int
			right, consumed = ParseFilter2(depth, filter[i+1:])
			i += (consumed)

			left := currentNode
			currentNode = &FilterAndExp{
				left:  left,
				right: right,
			}

		case ' ':
			// Skip: do nothing.
			// line(depth, "ignore: space")
		default:
			var chars []byte
			for j := i; j < len(filter); j++ {
				var c = filter[j]
				if c == ',' || c == '|' || c == '(' || c == ')' {
					i += 1
					break
				} else if c == ' ' {
					//skip
					continue
				}
				chars = append(chars, c)
				i = j - 1
			}

			line(depth, "word found: "+string(chars))
			termNode := &TerminalFilterExp{
				kind: KindTag,
				val:  string(chars),
			}

			if currentNode == nil {
				//spew.Dump("termNode =>", termNode)
				currentNode = termNode
			} else {
				if cn, ok := currentNode.(*FilterOrExp); ok {
					cn.right = termNode
				}
				if cn, ok := currentNode.(*FilterAndExp); ok {
					cn.right = termNode
				}
				//spew.Dump("currentNode =>", currentNode)
				// An overwrite is occurring here.
				//panic("should not happen")
			}
		}
	}

	fmt.Println("FINAL FUCKING RETURN!!!!")
	return currentNode, len(filter)
}
*/

/*
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
*/

func executeFilters(ctx context.Context, devList []*WrappedDevice) []*WrappedDevice {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	var normalizeTags = func(vals []string) mapset.Set[string] {
		var items []string
		for _, s := range vals {
			items = append(items, strings.Replace(strings.ToLower(s), "tag:", "", -1))
		}
		return mapset.NewSet[string](items...)
	}

	var (
		filteredDevList []*WrappedDevice
	)

	// Note: for better performance, filtering should be done by the most selective fields first.
	for _, dev := range devList {
		if cfg.Filters != nil {

			// Handle tags.
			// TODO: figure out all items to add.
			// TODO: clean this up and standardize the set logic.
			tagSet := normalizeTags(dev.Tags)
			userSet := mapset.NewSet[string](strings.ToLower(dev.User))
			everything := tagSet.Union(userSet)
			everything.Add(strings.ToLower(dev.OS))

			for _, a := range dev.Addresses {
				everything.Add(a)
			}

			semanticVersion := strings.ToLower(strings.Split(dev.ClientVersion, "-")[0])
			everything.Add(semanticVersion)

			fmt.Println(everything)

			// Single shot filter allows complex filter expressions.
			if !cfg.Filters.Eval(everything) {
				log.Info("device was filtered out!")
				continue
			}
			//if !cfg.Filters.Eval(tag) {
			//	//log.Info("device was filter on tags")
			//	//continue
			//}
			//
			//// Handle users
			//users := mapset.NewSet[string](dev.User)
			//if !cfg.Filters.Eval(users) {
			//	log.Info("device was filter on user")
			//	continue
			//}
		}
		// Filter by 'tag' when provided - currently only supports full matching.
		//if f, exists := cfg.Filters["tag"]; exists {
		//	normalizedTags := normalizeTags(dev.Tags)
		//	wantsNoTags := f.Contains("nil") // User wants to filter out rows with tags.
		//
		//	// Determine if the device should be skipped based on tag presence and user's filter.
		//	hasTags := len(dev.Tags) > 0
		//	matchesTags := f.ContainsAny(normalizedTags...)
		//
		//	// Skip device if it doesn't match the filter criteria.
		//	if (wantsNoTags && hasTags) || (!wantsNoTags && !matchesTags) {
		//		continue
		//	}
		//}

		// Filter by 'user' when provided - currently only supports full matching.
		//if f, exists := cfg.Filters["user"]; exists {
		//	if !f.Contains(strings.ToLower(dev.User)) {
		//		continue
		//	}
		//}

		// Filter by 'version' when provided - currently only supports full matching.
		//if f, exists := cfg.Filters["version"]; exists {
		//	// For now, just filter on the first portion of the version which has the format: 1.xx.1
		//	semanticVersion := strings.Split(dev.ClientVersion, "-")[0]
		//	if !f.Contains(strings.ToLower(semanticVersion)) {
		//		continue
		//	}
		//}

		// Filters by ipv4 - currently only supports full matching.
		//if f, exists := cfg.Filters["ipv4"]; exists {
		//	if (len(dev.Addresses) == 0) || !f.ContainsAny(dev.Addresses...) {
		//		continue
		//	}
		//}

		// Filters by ipv6 - currently only supports full matching.
		//if f, exists := cfg.Filters["ipv6"]; exists {
		//	if (len(dev.Addresses) == 0) || !f.ContainsAny(dev.Addresses...) {
		//		continue
		//	}
		//}

		// TODO: additional filters - lastSeen

		// Filter by 'os' when provided - currently only supports full matching.
		//if f, exists := cfg.Filters["os"]; exists {
		//	if !f.Contains(strings.ToLower(dev.OS)) {
		//		continue
		//	}
		//}

		// Filters by exit node - filter query looks like: --filter 'exit:yes|no'
		//if f, exists := cfg.Filters["exit"]; exists {
		//	// NOTE: this information only exists from enriched results.
		//	if enrichedDev := dev.EnrichedInfo; enrichedDev != nil {
		//		if (f.Contains("yes") && !enrichedDev.HasExitNodeOption) ||
		//			(f.Contains("no") && enrichedDev.HasExitNodeOption) {
		//			continue
		//		}
		//	}
		//}

		filteredDevList = append(filteredDevList, dev)
	}

	return filteredDevList
}
