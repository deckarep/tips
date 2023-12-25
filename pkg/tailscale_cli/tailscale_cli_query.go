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

package tailscale_cli

import (
	"os/exec"
	"strings"
	"tips/pkg/utils"

	"github.com/tidwall/gjson"
)

var (
	// binarySearchPathCandidates should be extended with all known paths that the Tailscale cli can exist.
	// Additionally, the order has an influence on which path will be selected first.
	binarySearchPathCandidates = map[string][]string{
		"linux": {
			"/usr/bin/tailscale",
		},
		"darwin": {
			// When install via Mac App Store.
			"/Applications/Tailscale.app/Contents/MacOS/Tailscale",
		},
	}
)

type DeviceInfo struct {
	NodeKey           string   `json:"node_key"`
	DNSName           string   `json:"dns_name"`
	HasExitNodeOption bool     `json:"has_exit_node_option"`
	IsSelf            bool     `json:"is_self"`
	Online            bool     `json:"online"`
	Tags              []string `json:"tags"`
}

func GetVersion() (string, error) {
	confirmedPath, err := utils.SelectBinaryPath(binarySearchPathCandidates)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(confirmedPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetDevicesState() (map[string]DeviceInfo, error) {
	confirmedPath, err := utils.SelectBinaryPath(binarySearchPathCandidates)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(confirmedPath, "status", "--json")

	// Running the command and capturing its output
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	results := make(map[string]DeviceInfo)

	jo := string(output)

	// TODO: Lazy, dynamic json without checking type assertions. Clean this up.
	var toTagSet = func(results []gjson.Result) []string {
		var tags []string
		for _, t := range results {
			tags = append(tags, strings.Replace(t.String(), "tag:", "", -1))
		}
		return tags
	}

	// Grab the Self info.
	selfNodeKey := gjson.Get(jo, "Self.PublicKey").String()
	selfOnline := gjson.Get(jo, "Self.Online").Bool()
	selfDNSName := gjson.Get(jo, "Self.DNSName").String()
	selfTagSet := toTagSet(gjson.GetMany(jo, "Self.Tags"))
	selfExitNodeOption := gjson.Get(jo, "Self.ExitNodeOption").Bool()

	results[selfNodeKey] = DeviceInfo{
		NodeKey:           selfNodeKey,
		DNSName:           selfDNSName,
		Online:            selfOnline,
		IsSelf:            true,
		Tags:              selfTagSet,
		HasExitNodeOption: selfExitNodeOption,
	}

	// Grab the peers
	peers := gjson.Get(jo, "Peer")
	peers.ForEach(func(key, value gjson.Result) bool {
		exitNodeOption := value.Get("ExitNodeOption").Bool()
		peerNodeKey := value.Get("PublicKey").String()
		peerOnline := value.Get("Online").Bool()
		dnsName := value.Get("DNSName").String()
		tagSet := toTagSet(value.Get("Tags").Array())

		results[peerNodeKey] = DeviceInfo{
			NodeKey:           peerNodeKey,
			DNSName:           dnsName,
			Online:            peerOnline,
			HasExitNodeOption: exitNodeOption,
			IsSelf:            false,
			Tags:              tagSet,
		}
		return true
	})

	return results, nil
}
