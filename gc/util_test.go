// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package gc

import (
	"testing"
)

func TestSkipState(t *testing.T) {
	var tests = []struct {
		state string
		want  bool
	}{
		{"exited", false},
		{"created", true},
		{"running", true},
	}
	for _, test := range tests {
		if got, want := skipState(test.state), test.want; got != want {
			t.Errorf("Want skipState %v, got %v", want, got)
		}
	}
}

func TestSkipImage(t *testing.T) {
	var tests = []struct {
		image string
		want  bool
	}{
		{"drone/drone", true},
		{"docker.io/drone/drone", true},
		{"docker.io/drone/drone:latest", true},
		// agent
		{"drone/agent", true},
		{"docker.io/drone/agent", true},
		{"docker.io/drone/agent:latest", true},
		// autoscaler
		{"drone/autoscaler", true},
		{"docker.io/drone/autoscaler", true},
		{"docker.io/drone/autoscaler:latest", true},
		// gc
		{"drone/gc", true},
		{"docker.io/drone/gc", true},
		{"docker.io/drone/gc:latest", true},
		// misc
		{"golang", false},
		{"alpine", false},
		{"docker.io/library/busybox", false},
	}
	for _, test := range tests {
		if got, want := skipImage(test.image), test.want; got != want {
			t.Errorf("Want skipImage %v, got %v", want, got)
		}
	}
}

func TestMatchPatterns(t *testing.T) {
	var tests = []struct {
		name string
		path string
		want bool
	}{
		{"drone", "drone", true},
		{"drone-server", "drone-*", true},
		{"mini-kube", "*-kube", true},
		{"redis", "redis", true},
		{"redis", "drone-*", false},
		// fully qualified names
		{"drone/drone", "docker.io/drone/drone:latest", true},
		{"redis", "docker.io/library/redis:*", true},
		{"redis", "docker.io/library/redis:1", false}, // tag mismatch
	}
	for _, test := range tests {
		matched := matchPatterns([]string{test.name}, []string{test.path})
		if got, want := matched, test.want; got != want {
			t.Errorf("Want matchPatterns %v, got %v", want, got)
		}
	}
}
