// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"time"

	"docker.io/go-docker/api/types"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

func (c *collector) collectContainers(ctx context.Context) error {
	var result error

	logger := log.Ctx(ctx)
	containers, err := c.client.ContainerList(ctx, containerListArgs)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot list containers")
		return err
	}

	now := time.Now()
	for _, cc := range containers {
		if skipImage(cc.Image) {
			continue
		}

		if skipState(cc.State) {
			continue
		}

		if matchPatterns(cc.Names, c.whitelist) {
			continue
		}

		t := time.Unix(cc.Created, 0)
		if t.Add(time.Hour).After(now) {
			continue
		}

		logger.Info().
			Strs("name", cc.Names).
			Msg("remove container")

		err := c.client.ContainerRemove(ctx, cc.ID, containerRemoveOpts)
		if err != nil {
			logger.Error().
				Err(err).
				Strs("name", cc.Names).
				Msg("cannot remove container")

			result = multierror.Append(result, err)
			continue
		}

		logger.Info().
			Strs("name", cc.Names).
			Msg("successfully removed container")
	}
	return result
}

var containerListArgs = types.ContainerListOptions{
	All: true,
}

var containerRemoveOpts = types.ContainerRemoveOptions{
	RemoveVolumes: true,
	RemoveLinks:   false,
	Force:         true,
}
