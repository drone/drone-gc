// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"path"
	"strings"

	"github.com/drone/drone-gc/gc/internal"
)

func skipState(state string) bool {
	switch state {
	case "exited":
		return false
	default:
		return true
	}
}

func skipImage(image string) bool {
	image = internal.ExpandImage(image)
	return strings.HasPrefix(image, "docker.io/drone/")
}

func matchPatterns(names []string, patterns []string) bool {
	for _, name := range names {
		full := internal.ExpandImage(name)
		for _, pattern := range patterns {
			matched, _ := path.Match(pattern, name)
			if matched {
				return true
			}
			matched, _ = path.Match(pattern, full)
			if matched {
				return true
			}
		}
	}
	return false
}
