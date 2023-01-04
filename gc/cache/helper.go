// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package cache

import (
	"context"

	docker "github.com/docker/docker/client"
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
	go l.listen(ctx) // nolint: errcheck
	return &client{
		APIClient: api,
		cache:     c,
	}
}
