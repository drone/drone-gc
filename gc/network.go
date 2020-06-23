// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import (
	"context"

	"github.com/docker/docker/api/types"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

func (c *collector) collectNetworks(ctx context.Context) error {
	var result error

	logger := log.Ctx(ctx)
	networks, err := c.client.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot list networks")
		return err
	}

	for _, v := range networks {
		if isProtected(v.Labels) {
			logger.Debug().
				Str("name", v.Name).
				Msg("network is protected")
			continue
		}
		if isExpired(v.Labels) == false {
			logger.Debug().
				Str("name", v.Name).
				Msg("network not expired")
			continue
		}

		logger.Debug().
			Str("name", v.Name).
			Msg("remove network")

		err = c.client.NetworkRemove(ctx, v.Name)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("cannot remove network")
			result = multierror.Append(result, err)
			continue
		}

		logger.Info().
			Str("name", v.Name).
			Msg("network removed")
	}
	return result
}
