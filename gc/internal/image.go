// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
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
