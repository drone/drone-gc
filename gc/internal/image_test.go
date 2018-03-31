// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package internal

import "testing"

func Test_expandImage(t *testing.T) {
	testdata := []struct {
		from string
		want string
	}{
		{
			from: "golang",
			want: "docker.io/library/golang:latest",
		},
		{
			from: "golang:latest",
			want: "docker.io/library/golang:latest",
		},
		{
			from: "golang:1.0.0",
			want: "docker.io/library/golang:1.0.0",
		},
		{
			from: "library/golang",
			want: "docker.io/library/golang:latest",
		},
		{
			from: "library/golang:latest",
			want: "docker.io/library/golang:latest",
		},
		{
			from: "library/golang:1.0.0",
			want: "docker.io/library/golang:1.0.0",
		},
		{
			from: "index.docker.io/library/golang:1.0.0",
			want: "docker.io/library/golang:1.0.0",
		},
		{
			from: "gcr.io/golang",
			want: "gcr.io/golang:latest",
		},
		{
			from: "gcr.io/golang:1.0.0",
			want: "gcr.io/golang:1.0.0",
		},
		// error cases, return input unmodified
		{
			from: "foo/bar?baz:boo",
			want: "foo/bar?baz:boo",
		},
	}
	for _, test := range testdata {
		got, want := ExpandImage(test.from), test.want
		if got != want {
			t.Errorf("Want image %q expanded to %q, got %q", test.from, want, got)
		}
	}
}
