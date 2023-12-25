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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"

	"github.com/charmbracelet/log"

	bolt "go.etcd.io/bbolt"
)

const (

	// These buckets only point to keys where the actual data exists in those respective buckets.
	// TODO:

	// These two buckets contain FULL data.
	devicesBucket  = "devices.full"
	enrichedBucket = "enriched.full"

	statsBucket = "stats"
	statsKey    = "stats"
)

type DBStats struct {
	DevicesCount  int `json:"devices_count"`
	EnrichedCount int `json:"enriched_count"`
}

type DB struct {
	tailnetScope string
	hdl          *bolt.DB
}

func NewDB(tailnetScope string) DB {
	return DB{
		tailnetScope: tailnetScope,
	}
}

func (d *DB) TailnetScope() string {
	return d.tailnetScope
}

func (d *DB) File() string {
	return fmt.Sprintf("%s.db.bolt", d.tailnetScope)
}

func (d *DB) Open() error {
	db, err := bolt.Open(d.File(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	d.hdl = db
	return nil
}

func (d *DB) Close() error {
	if d.hdl != nil {
		return d.hdl.Close()
	}
	return nil
}

func (d *DB) Erase() error {
	// TODO: destroy the local file.
	return nil
}

func (d *DB) Exists(ctx context.Context) (bool, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	return fileExistsAndIsRecent(d.File(), cfg.CacheTimeout)
}

func (d *DB) IndexDevices(ctx context.Context, devList []tailscale.Device, enrichedDevList map[string]tailscale_cli.DeviceInfo) error {
	// Start a writable transaction.
	err := d.hdl.Update(func(tx *bolt.Tx) error {
		// Create all buckets.
		statsBuck, err := tx.CreateBucketIfNotExists([]byte(statsBucket))
		if err != nil {
			return err
		}

		// If the bucket already exists, it will return a reference to it.
		devicesBucket, err := tx.CreateBucketIfNotExists([]byte(devicesBucket))
		if err != nil {
			return err
		}

		enrichedBucket, err := tx.CreateBucketIfNotExists([]byte(enrichedBucket))
		if err != nil {
			return err
		}

		// Record stats
		// TODO: record version
		var stats DBStats
		stats.DevicesCount = len(devList)
		stats.EnrichedCount = len(enrichedDevList)

		encoded, err := json.Marshal(stats)
		if err != nil {
			return err
		}

		err = statsBuck.Put([]byte(statsKey), encoded)
		if err != nil {
			return err
		}

		// Iterate over all devices.
		for _, dev := range devList {
			// Encode as JSON: in the future encode as Proto/more compact form.
			encoded, err := json.Marshal(dev)
			if err != nil {
				return err
			}

			// The key is the data point ID converted to a byte slice.
			// The value is the encoded JSON data.
			err = devicesBucket.Put([]byte(dev.Name), encoded)
			if err != nil {
				return err
			}
		}

		// Iterate over all enriched info.
		for _, enr := range enrichedDevList {
			// Encode as JSON: in the future encode as Proto/more compact form.
			encoded, err := json.Marshal(enr)
			if err != nil {
				return err
			}

			// The key is the data point ID converted to a byte slice.
			// The value is the encoded JSON data.
			err = enrichedBucket.Put([]byte(enr.NodeKey), encoded)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal("error updating the database index", "error", err)
	}

	return nil
}

func (d *DB) FindDevices(ctx context.Context) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	// 0. First populate all devices.
	// TODO: index metadata like the size of devices then we can instantiate with the correct capacity.
	//devList := make([]tailscale.Device, 0)
	var devList []tailscale.Device
	var enrichedDevs map[string]tailscale_cli.DeviceInfo
	err := d.hdl.View(func(tx *bolt.Tx) error {
		// 0. Check for stats
		sb := tx.Bucket([]byte(statsBucket))
		if sb == nil {
			return errors.New("bucket is unknown: " + statsBucket)
		}
		sts := sb.Get([]byte(statsKey))
		var stats DBStats
		err := json.Unmarshal(sts, &stats)
		if err != nil {
			return err
		}

		devList = make([]tailscale.Device, 0, stats.DevicesCount)
		enrichedDevs = make(map[string]tailscale_cli.DeviceInfo, stats.EnrichedCount)

		// 1. Next populate all devices data
		b := tx.Bucket([]byte(devicesBucket))
		if b == nil {
			return errors.New("bucket is unknown: " + devicesBucket)
		}
		c := b.Cursor()

		// This is a linear scan over all key/values
		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			//fmt.Printf("key=%s\n", k)

			var dev tailscale.Device
			err := json.Unmarshal(v, &dev)
			if err != nil {
				return err
			}

			devList = append(devList, dev)
		}

		// 2. Next populate only enriched data, that is needed.
		// No cursor is needed, we only need to get the enriched data for devices that were returned above!
		b = tx.Bucket([]byte(enrichedBucket))

		for _, dev := range devList {
			k := dev.NodeKey
			v := b.Get([]byte(k))
			var dev tailscale_cli.DeviceInfo
			err := json.Unmarshal(v, &dev)
			if err != nil {
				return err
			}

			enrichedDevs[k] = dev
		}

		// TODO: need to set this up.
		// Prefix scan (use this in the future)
		//prefix := []byte("1234")
		//for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		//	fmt.Printf("key=%s, value=%s\n", k, v)
		//}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return devList, enrichedDevs, nil
}

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
