// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import (
	"path"
	"strconv"
	"strings"
	"time"

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
	switch {
	case strings.HasPrefix(image, "docker.io/drone/"):
		return true
	case strings.HasPrefix(image, "docker.io/v2tec/watchtower:"):
		return true
	default:
		return false
	}
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

func isExpired(labels map[string]string) bool {
	l, ok := labels["io.drone.expires"]
	if !ok {
		return false
	}
	i, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		return true
	}
	t := time.Unix(i, 0)
	return time.Now().After(t)
}

func isProtected(labels map[string]string) bool {
	return labels["io.drone.protected"] == "true"
}
