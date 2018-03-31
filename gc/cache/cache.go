// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"sort"
	"sync"

	"docker.io/go-docker/api/types"
)

// DefaultCacheSize is the default size of the LFRU cache.
const DefaultCacheSize = 1000

type cache struct {
	mu sync.Mutex

	limit int
	list  []*types.ImageSummary
	index map[string]*types.ImageSummary
}

func newCache(limit int) *cache {
	return &cache{
		limit: limit,
		list:  []*types.ImageSummary{},
		index: make(map[string]*types.ImageSummary),
	}
}

func (c *cache) push(name string, value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	image, ok := c.index[name]
	if ok {
		image.Created = value
	} else {
		image = &types.ImageSummary{
			RepoTags: []string{name},
			Created:  value,
		}
		c.list = append(c.list, image)
		c.index[name] = image
	}
	sort.Sort(byCreatedAsc(c.list))
	if len(c.list) > c.limit {
		c.list = c.list[:c.limit]
	}
}

func (c *cache) find(name string) (v int64, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.index[name]; ok {
		return item.Created, true
	}
	return
}
