// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package internal

import (
	"github.com/docker/distribution/reference"
)

// ExpandImage returns the fully qualified image name.
func ExpandImage(name string) string {
	ref, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return name
	}
	return reference.TagNameOnly(ref).String()
}
