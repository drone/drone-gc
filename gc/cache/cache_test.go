// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"testing"

	"docker.io/go-docker/api/types"
	"github.com/google/go-cmp/cmp"
)

func TestCache(t *testing.T) {
	c := newCache(5)
	c.push("alpine:latest", 359596800)
	c.push("busybox:latest", 420681600)
	c.push("golang:1", 1192233600)
	c.push("golang:1.9", 1192233603)
	c.push("golang:1.8", 1192233602)
	c.push("golang:1.7", 1192233601)

	if got, want := len(c.list), 5; got != want {
		t.Errorf("Want %d items in the cache, got %d", want, got)
	}

	want := []*types.ImageSummary{
		&types.ImageSummary{
			Created:  1192233603,
			RepoTags: []string{"golang:1.9"},
		},
		&types.ImageSummary{
			Created:  1192233602,
			RepoTags: []string{"golang:1.8"},
		},
		&types.ImageSummary{
			Created:  1192233601,
			RepoTags: []string{"golang:1.7"},
		},
		&types.ImageSummary{
			Created:  1192233600,
			RepoTags: []string{"golang:1"},
		},
		&types.ImageSummary{
			Created:  420681600,
			RepoTags: []string{"busybox:latest"},
		},
		// note that we expect the alpine container is
		// removed because the cache limit is 5 items.
	}
	if !cmp.Equal(want, c.list) {
		t.Errorf("Invalid cache order")
	}
}
