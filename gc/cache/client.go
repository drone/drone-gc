// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"context"
	"sort"

	"github.com/drone/drone-gc/gc/internal"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type client struct {
	docker.APIClient
	cache *cache
}

func (c *client) DiskUsage(ctx context.Context) (types.DiskUsage, error) {
	df, err := c.APIClient.DiskUsage(ctx)
	if err != nil {
		return df, err
	}
	for _, image := range df.Images {
		if len(image.RepoTags) == 0 {
			continue
		}
	tags:
		for _, tag := range image.RepoTags {
			tag = internal.ExpandImage(tag)
			unix, ok := c.cache.find(tag)
			if ok {
				image.Created = unix
				break tags
			}
		}
	}
	sort.Sort(byCreated(df.Images))
	return df, err
}
