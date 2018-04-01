// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"time"

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

	now := time.Now()
	for _, v := range volumes.Volumes {
		t, err := time.Parse("2006-01-02T15:04:05Z", v.CreatedAt)
		if err != nil {
			logger.Error().
				Err(err).
				Str("name", v.Name).
				Msg("invalid date time format")
			result = multierror.Append(result, err)
			continue
		}

		if t.Add(time.Hour).After(now) {
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
		Key:   "dangling",
		Value: "true",
	},
	filters.KeyValuePair{
		Key:   "driver",
		Value: "local",
	},
)
