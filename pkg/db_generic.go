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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
)

const (
	// These two buckets contain FULL data.
	DevicesBucket = "bucket:devices.full"
	StatsBucket   = "bucket:stats"

	StatsKey = "key:stats"
)

type Indexer interface {
	Key() string
}

type DB2Stats struct {
	DevicesCount  int `json:"devices_count"`
	EnrichedCount int `json:"enriched_count"`
}

type DB2[T Indexer] struct {
	tailnetScope string
	hdl          *bolt.DB
}

func NewDB2[T Indexer](tailnetScope string) *DB2[T] {
	return &DB2[T]{
		tailnetScope: tailnetScope,
	}
}

func (d *DB2[T]) TailnetScope() string {
	return d.tailnetScope
}

func (d *DB2[T]) File() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get current user; aborting", "error", err)
	}

	return path.Join(u.HomeDir, fmt.Sprintf("%s.db.bolt", d.tailnetScope))
}

func (d *DB2[T]) Open() error {
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

func (d *DB2[T]) Close() error {
	if d.hdl != nil {
		return d.hdl.Close()
	}
	return nil
}

func (d *DB2[T]) Erase() error {
	fileToDelete := d.File()
	if !strings.HasSuffix(fileToDelete, "db.bolt") {
		return nil
	}

	return deleteFileIfExists(fileToDelete)
}

func (d *DB2[T]) Exists(ctx context.Context) (bool, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	return fileExistsAndIsRecent(d.File(), cfg.CacheTimeout)
}

func (d *DB2[T]) IndexOpaqueItems(ctx context.Context, bucketName string, items []T) error {
	if d.hdl == nil {
		return errors.New("trying to index db when handle to db is nil")
	}

	err := d.hdl.Update(func(tx *bolt.Tx) error {
		bckt, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		// TODO: stats

		for _, item := range items {
			if err := d.put(bckt, item.Key(), item); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

type DBQuery struct {
	PrefixFilters *PrefixFilter
	PrimaryKeys   []string
}

func (d *DB2[T]) LookupOpaqueItem(ctx context.Context, bucketName, primaryKey string) (*T, error) {
	var item *T
	err := d.hdl.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		i, err := d.lookupOpaqueItem(b, primaryKey)
		if err != nil {
			return err
		}
		item = i
		return nil
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (d *DB2[T]) lookupOpaqueItem(bucket *bolt.Bucket, primaryKey string) (*T, error) {
	item, err := d.get(bucket, primaryKey)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// SearchOpaqueItems can generically search with 3 different ways.
// 1. Using one or more primary keys, in which case this is a direct lookup (not technically a search)
// 2. Using the * (all/everything) construct, this is just a full table scan really.
// 3. Using a prefix scan, this is a seek to a segment of the index and should be fast assuming good selectivity.
func (d *DB2[T]) SearchOpaqueItems(ctx context.Context, bucketName string, query DBQuery) ([]T, error) {
	//cfg := CtxAsConfig(ctx, CtxKeyConfig)

	if query.PrefixFilters.Count() == 0 {
		panic("query.PrefixFilter must never be empty, in the case of all it must be: *")
	}

	var items []T

	if err := d.hdl.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New("bucket is unknown: " + bucketName)
		}

		// Search by primary keys, this a direct lookup, the fastest.
		if len(query.PrimaryKeys) > 0 {
			for _, pk := range query.PrimaryKeys {
				item, err := d.lookupOpaqueItem(b, pk)
				if err != nil {
					return err
				}
				items = append(items, *item)
			}
		} else if query.PrefixFilters.IsAll() {
			c := b.Cursor()
			// Search by everything, linear (full-table scan)
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var item T
				if err := json.Unmarshal(v, &item); err != nil {
					return err
				}
				items = append(items, item)
			}
		} else {
			// Since Prefix Filters support OR filters: "foo|bar" we do this in a loop with a new cursor for each prefix.
			// Most of the time, just one iteration occurs.
			for i := 0; i < query.PrefixFilters.Count(); i++ {
				c := b.Cursor()
				prefix := []byte(query.PrefixFilters.PrefixAt(i))
				for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
					var item T
					if err := json.Unmarshal(v, &item); err != nil {
						return err
					}
					items = append(items, item)
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

func (d *DB2[T]) put(bucket *bolt.Bucket, key string, data T) error {
	// Encode as JSON: in the future encode as Proto/more compact form.
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// The key is the data point ID converted to a byte slice.
	// The value is the encoded JSON data.
	return bucket.Put([]byte(key), encoded)
}

func (d *DB2[T]) get(bucket *bolt.Bucket, key string) (T, error) {
	var obj T
	v := bucket.Get([]byte(key))
	err := json.Unmarshal(v, &obj)
	return obj, err
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
