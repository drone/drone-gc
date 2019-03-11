// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package cache

import (
	"testing"

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
	c.push("golang:1.7", 1192233601) // bump hit count x2
	c.push("golang:1.7", 1192233601) // bump hit count x3

	if got, want := len(c.list), 5; got != want {
		t.Errorf("Want %d items in the cache, got %d", want, got)
	}

	want := []*item{
		&item{
			Last: 1192233603,
			Hits: 1,
			Name: "golang:1.9",
		},
		&item{
			Last: 1192233602,
			Hits: 1,
			Name: "golang:1.8",
		},
		&item{
			Last: 1192233601,
			Hits: 3,
			Name: "golang:1.7",
		},
		&item{
			Last: 1192233600,
			Hits: 1,
			Name: "golang:1",
		},
		&item{
			Last: 420681600,
			Hits: 1,
			Name: "busybox:latest",
		},
		// note that we expect the alpine container is
		// removed because the cache limit is 5 items.
	}
	if !cmp.Equal(want, c.list) {
		t.Errorf("Invalid cache order")
	}
}
