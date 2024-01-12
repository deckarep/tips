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
	"strings"

	"github.com/deckarep/tips/pkg/filtercomp"

	mapset "github.com/deckarep/golang-set/v2"
)

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

	for _, dev := range devList {

		// FAT TODO: clean this up and standardize the set logic.
		// TODO: figure out all items to add (everything to filter on).

		// Tags
		everything := normalizeTags(dev.Tags)

		// User
		everything.Add(strings.ToLower(dev.User))

		// OS
		everything.Add(strings.ToLower(dev.OS))

		// Version
		semanticVersion := strings.ToLower(strings.Split(dev.ClientVersion, "-")[0])
		everything.Add(semanticVersion)

		// ipv4/ipv6
		for _, a := range dev.Addresses {
			everything.Add(a)
		}

		// Exit node status
		// I'm somewhat happy with this approach. However, when the data is not enriched this will incorrectly
		// flag everything as a non-exit node.
		if dev.EnrichedInfo != nil && dev.EnrichedInfo.HasExitNodeOption {
			everything.Add("+exit")
		} else {
			everything.Add("-exit")
		}

		// Apply the single-shot filter: allows complex filter expressions.
		if !cfg.Filters.Eval(everything) {
			continue
		}

		// TODO: a way to filter for NO TAGS => !tag
		// TODO: additional filters like lastSeen and conditional filtering: >, >=, <, <=
		// TODO: meta-filters on things like exit node status => +exit

		filteredDevList = append(filteredDevList, dev)
	}

	return filteredDevList
}
