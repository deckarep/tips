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
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func DevicesResource(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()
	defer func() {
		cfg.TailscaleAPI.ElapsedTime = time.Since(startTime)
	}()

	// 0. Check cache - return cached results if cache timeout not yet expired.
	if cfg.NoCache {
		log.Warn("--nocache is not yet support, but this should force a refresh of data in cache")
	}

	// 1. Do tailscale api lookup for devices data.
	ctxTimeOut, cancelTimeout := context.WithTimeout(ctx, cfg.TailscaleAPI.Timeout)
	defer cancelTimeout()
	devList, err := client.Devices(ctxTimeOut)
	if err != nil {
		log.Fatal("tailscale api failed during devices lookup: ", err)
	}

	// 2. When available, enrich this data with data from the Tailscale cli, if this is run from a node within the tailnet.
	enrichedDevices, err := tailscale_cli.GetDevicesState()
	if err != nil {
		log.Debug("unable to get enriched data from tailscale cli", "error", err)
	}

	return devList, enrichedDevices, nil
}
