// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import "time"

// Option configures a garbage collector option.
type Option func(*collector)

// WithImageWhitelist returns an option to set an image
// whitelist. This will prevent the garbage collector from
// removing named containers.
func WithImageWhitelist(images []string) Option {
	return func(c *collector) {
		c.reserved = append(c.reserved, images...)
	}
}

// WithWhitelist returns an option to set a whitelist of
// container names. This will prevent the garbage collector
// from removing matching containers.
func WithWhitelist(names []string) Option {
	return func(c *collector) {
		c.whitelist = append(c.whitelist, names...)
	}
}

// WithThreshold returns an option to set a threshold
// for the image cache. The cache will clear images until
// the layer size is below the target threshold.
func WithThreshold(threshold int64) Option {
	return func(c *collector) {
		c.threshold = threshold
	}
}

// WithMinImageAge returns an option to set the minimum
// age a image should be to become a candidate for removal.
// Images younger than this value won't be removed
func WithMinImageAge(minImageAge time.Duration) Option {
	return func(c *collector) {
		c.minImageAge = minImageAge
	}
}

// ReservedImages provides a list of reserved images names
// that should not be removed.
var ReservedImages = []string{
	"drone/drone:*",
	"drone/agent:*",
	"drone/gc:*",
	"drone/autoscaler:*",
}

// ReservedNames provides a list of reserved container names
// that should not be removed.
var ReservedNames = []string{
	"drone",
	"drone-server",
	"agent",
	"drone-agent",
	"gc",
	"drone-gc",
	"autoscaler",
	"autoscale",
	"watchtower",
	"cadvisor",
}
