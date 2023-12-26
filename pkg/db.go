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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"

	"github.com/charmbracelet/log"

	bolt "go.etcd.io/bbolt"
)

const (
	// These two buckets contain FULL data.
	devicesBucket  = "bucket:devices.full"
	enrichedBucket = "bucket:enriched.full"
	statsBucket    = "bucket:stats"

	statsKey = "key:stats"
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
	// Already opened, it's a no-op.
	if d.hdl != nil {
		return nil
	}

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
	fileToDelete := d.File()
	if !strings.HasSuffix(fileToDelete, "db.bolt") {
		return nil
	}

	return deleteFileIfExists(fileToDelete)
}

func (d *DB) Exists(ctx context.Context) (bool, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	return fileExistsAndIsRecent(d.File(), cfg.CacheTimeout)
}

func put(bucket *bolt.Bucket, key string, data any) error {
	// Encode as JSON: in the future encode as Proto/more compact form.
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// The key is the data point ID converted to a byte slice.
	// The value is the encoded JSON data.
	return bucket.Put([]byte(key), encoded)
}

func get[T any](bucket *bolt.Bucket, key string) (T, error) {
	var obj T
	v := bucket.Get([]byte(key))
	err := json.Unmarshal(v, &obj)
	return obj, err
}

func (d *DB) IndexDevices(ctx context.Context, devList []tailscale.Device, enrichedDevList map[string]tailscale_cli.DeviceInfo) error {
	if d.hdl == nil {
		return errors.New("trying to index db when handle to db is nil")
	}

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
		// TODO: perhaps record version so we can delete version files that don't match. Overkill?
		stats := DBStats{
			DevicesCount:  len(devList),
			EnrichedCount: len(enrichedDevList),
		}

		err = put(statsBuck, statsKey, stats)
		if err != nil {
			return err
		}

		// Iterate over all devices.
		for _, dev := range devList {
			err = put(devicesBucket, dev.Name, dev)
			if err != nil {
				return err
			}
		}

		// Iterate over all enriched info.
		for _, enr := range enrichedDevList {
			err = put(enrichedBucket, enr.NodeKey, enr)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal("error creating the database index", "error", err)
	}

	return nil
}

func (d *DB) FindDevices(ctx context.Context) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	// 0. First populate all devices.
	var devList []tailscale.Device
	var enrichedDevs map[string]tailscale_cli.DeviceInfo
	err := d.hdl.View(func(tx *bolt.Tx) error {
		// 0. Check for stats
		sb := tx.Bucket([]byte(statsBucket))
		if sb == nil {
			return errors.New("bucket is unknown: " + statsBucket)
		}
		stats, err := get[DBStats](sb, statsKey)
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

		// Full scan when everything is requested.
		if cfg.PrefixFilter == "*" {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var dev tailscale.Device
				err := json.Unmarshal(v, &dev)
				if err != nil {
					return err
				}

				devList = append(devList, dev)
			}
		} else {
			// Prefix scan (much faster) when a prefix is present.
			prefix := []byte(cfg.PrefixFilter)
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				var dev tailscale.Device
				err := json.Unmarshal(v, &dev)
				if err != nil {
					return err
				}

				devList = append(devList, dev)
			}
		}

		// 2. Next populate only enriched data, that is needed.
		// No cursor is needed, we only need to get the enriched data for devices that were returned above!
		b = tx.Bucket([]byte(enrichedBucket))
		for _, dev := range devList {
			dev, err := get[tailscale_cli.DeviceInfo](b, dev.NodeKey)
			if err != nil {
				return err
			}

			enrichedDevs[dev.NodeKey] = dev
		}

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

// deleteFileIfExists deletes the file if it exists
func deleteFileIfExists(filename string) error {
	// Check if the file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// File does not exist, nothing to do
		return nil
	} else if err != nil {
		// Some other error occurred when trying to get the file info
		return err
	}

	// Delete the file
	return os.Remove(filename)
}
