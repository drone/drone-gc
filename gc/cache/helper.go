// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"context"

	"docker.io/go-docker"
)

// Wrap returns a wrapped copy of the Docker client that
// collects details about image use and sorts the disk usage
// report based on the image last used date, ascending.
func Wrap(ctx context.Context, api docker.APIClient) docker.APIClient {
	c := newCache(DefaultCacheSize)
	l := &listener{
		client: api,
		cache:  c,
	}
	go l.listen(ctx)
	return &client{
		APIClient: api,
		cache:     c,
	}
}
