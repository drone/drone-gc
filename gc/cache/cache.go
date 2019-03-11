// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package cache

import (
	"sort"
	"sync"
)

// DefaultCacheSize is the default size of the LFRU cache.
const DefaultCacheSize = 1000

type cache struct {
	mu sync.Mutex

	limit int
	list  []*item
	index map[string]*item
}

type item struct {
	Name string
	Hits int
	Last int64
}

func newCache(limit int) *cache {
	return &cache{
		limit: limit,
		list:  []*item{},
		index: make(map[string]*item),
	}
}

func (c *cache) push(name string, value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.index[name]
	if ok {
		i.Last = value
		i.Hits++
	} else {
		i = &item{
			Name: name,
			Hits: 1,
			Last: value,
		}
		c.list = append(c.list, i)
		c.index[name] = i
	}
	sort.Sort(byLastUsed(c.list))
	if len(c.list) > c.limit {
		c.list = c.list[:c.limit]
	}
}

func (c *cache) find(name string) (v int64, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.index[name]; ok {
		return item.Last, true
	}
	return
}

type byLastUsed []*item

func (a byLastUsed) Len() int           { return len(a) }
func (a byLastUsed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLastUsed) Less(i, j int) bool { return a[i].Last > a[j].Last }
