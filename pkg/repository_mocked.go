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

	"os"
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/charmbracelet/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"golang.org/x/exp/rand"
)

const (
	testDevicesFile = "testmode/devices.json"
)

type MockedDeviceRepo struct {
}

func NewMockedDeviceRepo() *MockedDeviceRepo {
	return &MockedDeviceRepo{}
}

func (r *MockedDeviceRepo) DevicesResource(ctx context.Context) ([]*WrappedDevice, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()
	defer func() {
		cfg.TailscaleAPI.ElapsedTime = time.Since(startTime)
	}()

	f, err := os.Open(testDevicesFile)
	if err != nil {
		log.Fatal("failed to read file", "error", err)
	}
	defer f.Close()

	var devs []tailscale.Device

	// Using jsoniter for now...stdlib json was a bit slower for large blobs.
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(f).Decode(&devs)
	if err != nil {
		log.Fatal("failed to Unmarshal file", "error", err)
	}

	wrappedDevs := make([]*WrappedDevice, 0, len(devs))

	// Total shameful hack in the interest of testing.
	// Enrich the data for now with this made up data so that must boxes appear online.
	var counter int
	for _, dev := range devs {
		isSelf := counter == 0
		isOnline := true
		offersExitNode := false
		if rand.Float32() < 0.01 {
			isOnline = false
		}
		if rand.Float32() < 0.05 {
			offersExitNode = true
		}

		wrappedDevs = append(wrappedDevs, &WrappedDevice{
			Device: dev,
			EnrichedInfo: &tailscale_cli.DeviceInfo{
				DNSName:           "", // currently unused
				HasExitNodeOption: offersExitNode,
				IsSelf:            isSelf,
				Online:            isOnline,
				Tags:              nil, // currently unused
			},
		})
		counter += 1
	}

	return wrappedDevs, nil
}
