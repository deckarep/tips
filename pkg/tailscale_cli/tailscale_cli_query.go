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
	"github.com/tidwall/gjson"
	"os/exec"
)

const (
	macOSAppStorePath = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"
	// TODO: other paths depending on where user installs/OS.
)

type DeviceInfo struct {
	DNSName string
	IsSelf  bool
	Online  bool
}

func GetVersion() (string, error) {
	confirmedPath, err := exec.LookPath(macOSAppStorePath)
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

func GetDevicesStatuses() (map[string]DeviceInfo, error) {
	// TODO: try to find the path or use the user's config setting.
	confirmedPath, err := exec.LookPath(macOSAppStorePath)
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
	// Grab the Self info.
	selfNodeKey := gjson.Get(jo, "Self.PublicKey").String()
	selfOnline := gjson.Get(jo, "Self.Online").Bool()
	selfDNSName := gjson.Get(jo, "Self.DNSName").String()
	results[selfNodeKey] = DeviceInfo{DNSName: selfDNSName, Online: selfOnline, IsSelf: true}

	// Grab the peers
	peers := gjson.Get(jo, "Peer")
	peers.ForEach(func(key, value gjson.Result) bool {
		peerNodeKey := value.Get("PublicKey").String()
		peerOnline := value.Get("Online").Bool()
		dnsName := value.Get("DNSName").String()
		results[peerNodeKey] = DeviceInfo{DNSName: dnsName, Online: peerOnline, IsSelf: false}
		return true
	})

	return results, nil
}
