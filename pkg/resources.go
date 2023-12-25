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
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func fileExistsAndIsRecent(filePath string, duration time.Duration) (bool, error) {
	// Check if the file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// The file does not exist
		return false, nil
	} else if err != nil {
		// There was some other error getting the file info
		return false, err
	}

	// Check the time since the file was created
	creationTime := info.ModTime()
	if time.Since(creationTime) <= duration {
		// The file is recent enough
		return true, nil
	}

	// The file exists but is not recent
	return false, nil
}

func DevicesResource(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()
	defer func() {
		cfg.TailscaleAPI.ElapsedTime = time.Since(startTime)
	}()

	// 0. Check this index first.
	indexedDB := NewDB("deckarep@gmail.com")

	existsAndRecent, err := indexedDB.Exists(ctx)
	if err != nil {
		log.Warn("problem checking for bolt db file", "error", err)
	}

	err = indexedDB.Open()
	if err != nil {
		return nil, nil, err
	}
	defer indexedDB.Close()

	if true {
		if devList, enrichedDevs, err := indexedDB.FindDevices(ctx); existsAndRecent && err == nil {
			log.Warn("found a bolt file and its recent enough so that will be used!")
			return devList, enrichedDevs, nil
		}
	}
	log.Warn("bolt file has expired or must be regenerated")

	// TODO: 0. Check cache config - return cached results if cache timeout not yet expired.
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

	// 3. Index this data.
	err = indexedDB.IndexDevices(ctx, devList, enrichedDevices)
	if err != nil {
		log.Debug("unable to index the devices", "error", err)
	}

	return devList, enrichedDevices, nil
}
