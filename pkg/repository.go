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
	"tips/pkg/tailscale_cli"

	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

type RemoteDeviceRepo struct {
	client *tailscale.Client
}

func NewRemoteDeviceRepo(client *tailscale.Client) *RemoteDeviceRepo {
	return &RemoteDeviceRepo{
		client: client,
	}
}

func (r *RemoteDeviceRepo) DevicesResource(ctx context.Context) ([]*WrappedDevice, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	log.Debug("doing remote lookup of devices data")

	// 1. Do tailscale api lookup for devices data.
	ctxTimeOut, cancelTimeout := context.WithTimeout(ctx, cfg.TailscaleAPI.Timeout)
	defer cancelTimeout()
	devList, err := r.client.Devices(ctxTimeOut)
	if err != nil {
		log.Fatal("tailscale api failed during devices lookup: ", err)
	}

	// 2. When available, enrich this data with data from the Tailscale cli, if this is run from a node within the
	// tailnet. NOTE: This data may not be available if this tool is not run within a node on the tailnet.
	enrichedDevices, err := tailscale_cli.GetDevicesState()
	if err != nil {
		log.Debug("unable to get enriched data from tailscale cli, but this is optionally returned",
			"error", err)
	}

	// NOTE: It's not ideal, but a "join" occurs here if it turns out we're operating on a node within the tailnet
	// and enriched device results were returned from the Tailscale CLI app. This may not always be the case! The
	// other thing that is not ideal is that we do a loop again to join the results.
	wrappedDevs := make([]*WrappedDevice, 0, len(devList))
	for _, dev := range devList {
		if enrichedInfo, exists := enrichedDevices[dev.NodeKey]; exists {
			wrappedDevs = append(wrappedDevs, &WrappedDevice{Device: dev, EnrichedInfo: &enrichedInfo})
		} else {
			wrappedDevs = append(wrappedDevs, &WrappedDevice{Device: dev})
		}
	}

	return wrappedDevs, nil
}
