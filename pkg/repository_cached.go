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

type CachedRepository struct {
	remoteRepo  *RemoteDeviceRepo
	indexedRepo *DB
}

func NewCachedRepo(remoteRepo *RemoteDeviceRepo, dbRepo *DB) *CachedRepository {
	return &CachedRepository{
		remoteRepo:  remoteRepo,
		indexedRepo: dbRepo,
	}
}

func (c *CachedRepository) DevicesResource(ctx context.Context) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	// 0. Check cache config - return cached results if cache timeout not yet expired.
	existsAndRecent, err := c.indexedRepo.Exists(ctx)
	if err != nil {
		log.Warn("problem checking for bolt db file", "error", err)
	}

	// 1. If the db cache file exists, and we're not asked to expunge the cache then return from the cache.
	if existsAndRecent && !cfg.NoCache {

		if err = c.indexedRepo.Open(); err != nil {
			return nil, nil, err
		}

		// Care is taken to measure just cache retrieval time.
		cachedStartTime := time.Now()
		if devList, enrichedDevs, err := c.indexedRepo.FindDevices(ctx); err == nil {
			log.Info("local db file (db.bolt) was found and recent enough so using this as a cache")
			c.indexedRepo.Close()
			cfg.CachedElapsed = time.Since(cachedStartTime)
			return devList, enrichedDevs, nil
		}
	}

	log.Debug("rebuilding local db cache file", "file", c.indexedRepo.File())
	if err = c.indexedRepo.Erase(); err != nil {
		return nil, nil, err
	}

	if err = c.indexedRepo.Open(); err != nil {
		return nil, nil, err
	}
	defer c.indexedRepo.Close()

	// 2. Do remote lookup if we got here.
	repoStartTime := time.Now()
	devList, enrichedDevices, err := c.remoteRepo.DevicesResource(ctx)
	if err != nil {
		return nil, nil, err
	}
	cfg.TailscaleAPI.ElapsedTime = time.Since(repoStartTime)

	// 3. Index the remotely found data.
	err = c.indexedRepo.IndexDevices(ctx, devList, enrichedDevices)
	if err != nil {
		log.Debug("unable to index the devices", "error", err)
	}

	// 4. Return the data from the db because the db can utilize the index on prefix filters.
	// In the future it may also do other heavyweight filters that we don't have to do in "user space"
	devList, enrichedDevices, err = c.indexedRepo.FindDevices(ctx)
	if err != nil {
		return nil, nil, err
	}

	return devList, enrichedDevices, nil
}
