// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"time"

	"docker.io/go-docker"
)

// FilterFunc filters the Docker resource based
// on its labels. If the function returns false,
// the resource is ignored.
type FilterFunc func(map[string]string) bool

// default timeout for the collection cycle.
var timeout = time.Hour

// Collector defines a Docker container garbage collector.
type Collector interface {
	Collect(context.Context) error
}

type collector struct {
	client docker.APIClient

	whitelist []string // reserved containers
	reserved  []string // reserved images
	threshold int64    // target threshold in bytes
	filter    FilterFunc
}

// New returns a garbage collector.
func New(client docker.APIClient, opt ...Option) Collector {
	c := new(collector)
	c.client = client
	for _, o := range opt {
		o(c)
	}
	return c
}

func (c *collector) Collect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c.collectContainers(ctx)
	c.collectDanglingImages(ctx)
	c.collectImages(ctx)
	c.collectNetworks(ctx)
	c.collectVolumes(ctx)
	return nil
}

// Schedule schedules the garbage collector to execute at the
// specified interval duration.
func Schedule(ctx context.Context, collector Collector, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			collector.Collect(ctx)
		}
	}
}
