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
	"time"

	"github.com/charmbracelet/log"
)

type InnerRepo interface {
	DevicesResource(ctx context.Context) ([]*WrappedDevice, error)
}

type CachedRepository struct {
	innerRepo InnerRepo
}

func NewCachedRepo(innerRepo InnerRepo) *CachedRepository {
	return &CachedRepository{
		innerRepo: innerRepo,
	}
}

func (c *CachedRepository) DevicesResource(ctx context.Context) ([]*WrappedDevice, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	// 0. Check cache config - return cached results if cache timeout not yet expired.

	// Note: The DB is instantiated on demand for flexibility of using the DB generically on different types.
	deviceIndexedRepo := NewDB2[*WrappedDevice](cfg.Tailnet)
	existsAndRecent, err := deviceIndexedRepo.Exists(ctx)
	if err != nil {
		log.Warn("problem checking for bolt db file", "error", err)
	}

	// 1. If the db cache file exists, and we're not asked to expunge the cache then return from the cache.
	if existsAndRecent && !cfg.NoCache {

		if err = deviceIndexedRepo.Open(); err != nil {
			return nil, err
		}

		// Care is taken to measure just cache retrieval time.
		cachedStartTime := time.Now()
		devList, err := deviceIndexedRepo.SearchOpaqueItems(ctx, DevicesBucket, DBQuery2{PrefixFilter: cfg.PrefixFilter})
		if err != nil {
			return nil, err
		}

		log.Debug("local db file (db.bolt) was found and recent enough so using this as a cache")
		deviceIndexedRepo.Close()
		cfg.CachedElapsed = time.Since(cachedStartTime)

		return devList, nil
	}

	log.Debug("destroying and rebuilding local db cache file", "file", deviceIndexedRepo.File())
	if err = deviceIndexedRepo.Erase(); err != nil {
		return nil, err
	}

	if err = deviceIndexedRepo.Open(); err != nil {
		return nil, err
	}
	defer deviceIndexedRepo.Close()

	// 2. Do remote lookup if we got here.
	repoStartTime := time.Now()
	devList, err := c.innerRepo.DevicesResource(ctx)
	if err != nil {
		return nil, err
	}
	cfg.TailscaleAPI.ElapsedTime = time.Since(repoStartTime)

	// 3. Index the remotely found data.
	err = deviceIndexedRepo.IndexOpaqueItems(ctx, DevicesBucket, devList)
	if err != nil {
		log.Debug("unable to index the devices", "error", err)
	}

	// 4. Return the data from the db because the db can utilize the index on prefix filters.
	// In the future it may also do other heavyweight filters that we don't have to do in "user space"
	devList, err = deviceIndexedRepo.SearchOpaqueItems(ctx, DevicesBucket, DBQuery2{PrefixFilter: cfg.PrefixFilter})
	if err != nil {
		return nil, err
	}

	return devList, nil
}
