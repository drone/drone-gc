// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import (
	"context"

	"docker.io/go-docker/api/types/filters"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

func (c *collector) collectVolumes(ctx context.Context) error {
	var result error

	logger := log.Ctx(ctx)
	volumes, err := c.client.VolumeList(ctx, volumeListArgs)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot list volumes")
		return err
	}

	for _, v := range volumes.Volumes {
		if isProtected(v.Labels) {
			logger.Debug().
				Str("name", v.Name).
				Msg("volume is protected")
			continue
		}
		if isExpired(v.Labels) == false {
			logger.Debug().
				Str("name", v.Name).
				Msg("volume not expired")
			continue
		}

		logger.Debug().
			Str("name", v.Name).
			Msg("remove volume")

		err = c.client.VolumeRemove(ctx, v.Name, false)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("cannot remove volume")
			result = multierror.Append(result, err)
			continue
		}

		logger.Info().
			Str("name", v.Name).
			Msg("volume removed")
	}
	return result
}

var volumeListArgs = filters.NewArgs(
	filters.KeyValuePair{
		Key:   "driver",
		Value: "local",
	},
)
