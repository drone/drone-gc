// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"context"

	"docker.io/go-docker/api/types/filters"
	"github.com/rs/zerolog/log"
)

func (c *collector) collectNetworks(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Debug().
		Msg("prune networks")

	report, err := c.client.NetworksPrune(ctx, networkPruneArgs)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot prune networks")
		return err
	}

	logger.Debug().
		Msg("networks pruned")

	for _, network := range report.NetworksDeleted {
		logger.Info().
			Str("network", network).
			Msg("network deleted")
	}
	return nil
}

var networkPruneArgs = filters.NewArgs(
	filters.KeyValuePair{
		Key:   "until",
		Value: "1h",
	},
)
