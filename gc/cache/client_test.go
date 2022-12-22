// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package cache

import (
	"context"
	"testing"

	"github.com/drone/drone-gc/gc/internal"
	"github.com/google/go-cmp/cmp"

	"github.com/docker/docker/api/types"
	"github.com/drone/drone-gc/mocks"
	"github.com/golang/mock/gomock"
)

func TestDiskUsage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		Images: []*types.ImageSummary{
			{
				ID: "a180b24e38ed",
				RepoTags: []string{
					"golang:1.0.0",
					"golang:1.0",
					"golang:1",
					"golang:latest",
				},
			},
			{
				ID: "4e38e38c8ce0",
				RepoTags: []string{
					"alpine:latest",
				},
			},
			{
				ID: "481995377a04",
				RepoTags: []string{
					"busybox:latest",
				},
			},
			// the redis image has not record in the cache
			// and therefore will be sorted using its image
			// creation date (which is 0 in this test)
			{
				ID: "6d8c4adbca87",
				RepoTags: []string{
					"redis:latest",
				},
			},
		},
	}

	api := mocks.NewMockAPIClient(controller)
	api.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)

	c := newCache(100)
	c.push(internal.ExpandImage("golang:1"), 1192233600)      // newest
	c.push(internal.ExpandImage("alpine:latest"), 359596800)  // oldest
	c.push(internal.ExpandImage("busybox:latest"), 420681600) // middle

	s := &client{
		APIClient: api,
		cache:     c,
	}

	got, _ := s.DiskUsage(context.Background())

	want := types.DiskUsage{
		Images: []*types.ImageSummary{
			// note the redis container is first in the list
			// because it was not included in the cache, and
			// did not have a Created date set.
			{
				Created:  0,
				ID:       "6d8c4adbca87",
				RepoTags: []string{"redis:latest"},
			},
			{
				Created:  359596800,
				ID:       "4e38e38c8ce0",
				RepoTags: []string{"alpine:latest"},
			},
			{
				Created:  420681600,
				ID:       "481995377a04",
				RepoTags: []string{"busybox:latest"},
			},
			{
				Created:  1192233600,
				ID:       "a180b24e38ed",
				RepoTags: []string{"golang:1.0.0", "golang:1.0", "golang:1", "golang:latest"},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Errorf("Invalid image order")
	}
}
