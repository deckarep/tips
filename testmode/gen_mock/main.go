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

package main

// This cli just generates fake data...it's meant to create a static test file of devices. You may have to periodically
// re-run this to add support for new fields that we want to process.

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"golang.org/x/exp/rand"
	"net"
	"strings"
	"time"
)

const (
	NUM_DEVICES = 3000
)

var (
	nameFmt                   = "tsc-serv-%04d"
	startingIPVAddress net.IP = net.IPv4(100, 100, 0, 0)
	tagPool                   = []string{"almond", "walnut", "peanut", "pistachio", "pecan", "cachew", "hazelnut"}
	users                     = []string{"admin@foo.net", "admin@foo.org", "jane@foo.net", "roberto@foo.com", "terry@foo.org", "ellen@foo.io"}
	os                        = []string{"iOS", "macOS", "linux", "pc"}

	expiresDate = time.Now().AddDate(0, 6, 0).UTC()
	createdDate = time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	versions    = []string{
		"1.54.1-t0a01efc8f-g3d0598425",
		"1.52.1-t0swffsasd-sdfwwdfsss",
		"1.43.0-t3sswgfdss-ggswwfdfdf",
	}
)

func main() {
	devices := make([]*tailscale.Device, 0, NUM_DEVICES)
	for i := 0; i < NUM_DEVICES; i++ {
		devices = append(devices, createFakeDevice(i))
	}

	b, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal all devices into json with err: ", err)
	}

	fmt.Println(string(b))
}

func createFakeDevice(idx int) *tailscale.Device {
	dev := &tailscale.Device{
		Addresses:                 getIPV4Set(),
		Name:                      fmt.Sprintf(nameFmt, idx),
		ID:                        getID(16),
		Authorized:                randBool(),
		User:                      oneOf(users),
		Tags:                      getTags(),
		KeyExpiryDisabled:         randBool(),
		BlocksIncomingConnections: randBool(),
		ClientVersion:             oneOf(versions),
		Created: tailscale.Time{
			createdDate,
		},
		Expires: tailscale.Time{
			expiresDate,
		},
		Hostname:        fmt.Sprintf(nameFmt, idx),
		IsExternal:      randBool(),
		LastSeen:        tailscale.Time{},
		MachineKey:      getHex("mkey:", 64),
		NodeKey:         getHex("nodekey:", 64),
		OS:              oneOf(os),
		UpdateAvailable: randBool(),
	}
	return dev
}

func getHex(prefix string, length int) string {
	var builder strings.Builder
	builder.WriteString(prefix)

	hexDigits := "0123456789abcdef"

	for i := len(prefix); i < length; i++ {
		randomIndex := rand.Intn(len(hexDigits))
		builder.WriteByte(hexDigits[randomIndex])
	}

	return builder.String()
}

func getID(n int) string {
	var builder strings.Builder

	for i := 0; i < n; i++ {
		digit := rand.Intn(10) // Generates a random integer between 0 and 9
		fmt.Fprintf(&builder, "%d", digit)
	}

	return builder.String()
}

func oneOf[E any](data []E) E {
	n := rand.Intn(len(data))
	return data[n]
}

func getTags() []string {
	var tags []string
	for _, t := range tagPool {
		if rand.Float32() < 0.2 {
			tags = append(tags, t)
		}
	}
	return tags
}

func randBool() bool {
	return rand.Float32() > 0.5
}

func getIPV4Set() []string {
	n := rand.Intn(2) + 1
	var results []string
	for i := 0; i < n; i++ {
		results = append(results, getNewIPV4().String())
	}
	return results
}

func getNewIPV4() net.IP {
	spankingNew := startingIPVAddress
	startingIPVAddress = incrementIP(spankingNew)
	return spankingNew
}

func incrementIP(ip net.IP) net.IP {
	// Create a copy of the IP to avoid modifying the original
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)

	// Increment the IP
	for i := len(newIP) - 1; i >= 0; i-- {
		newIP[i]++
		if newIP[i] != 0 {
			break
		}
	}

	return newIP
}
