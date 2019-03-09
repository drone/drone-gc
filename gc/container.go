// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package gc

import (
	"context"

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

	for _, cc := range containers {
		if skipImage(cc.Image) {
			continue
		}

		if matchPatterns(cc.Names, c.whitelist) {
			continue
		}

		if isProtected(cc.Labels) {
			logger.Debug().
				Strs("name", cc.Names).
				Msg("container is protected")
			continue
		}

		if isExpired(cc.Labels) == false {
			logger.Debug().
				Strs("name", cc.Names).
				Msg("container not expired")
			continue
		}

		if cc.State != "exited" {
			logger.Debug().
				Strs("name", cc.Names).
				Msg("kill long-running container")
			c.client.ContainerKill(ctx, cc.ID, "9")
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
