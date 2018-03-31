// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	"github.com/drone/drone-gc/gc"
	"github.com/drone/drone-gc/gc/cache"
	"github.com/drone/signal"

	"docker.io/go-docker"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	Once       bool          `envconfig:"GC_ONCE"`
	Debug      bool          `envconfig:"GC_DEBUG"`
	Color      bool          `envconfig:"GC_DEBUG_COLOR"`
	Pretty     bool          `envconfig:"GC_DEBUG_PRETTY"`
	Images     []string      `envconfig:"GC_IGNORE_IMAGES"`
	Containers []string      `envconfig:"GC_IGNORE_CONTAINERS"`
	Interval   time.Duration `envconfig:"GC_INTERVAL" default:"5m"`
	Threshold  int64         `envconfig:"GC_THRESHOLD" default:"5000000000"`
}

func main() {
	cfg := new(config)
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("Cannot load configuration variables")
	}

	client, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal().Err(err).
			Msg("Cannot create Docker client")
	}

	initLogger(cfg)
	ctx := log.Logger.WithContext(context.Background())
	ctx = signal.WithContext(ctx)

	collector := gc.New(
		cache.Wrap(ctx, client),
		gc.WithImageWhitelist(gc.ReservedImages),
		gc.WithImageWhitelist(cfg.Images),
		gc.WithThreshold(cfg.Threshold),
		gc.WithWhitelist(gc.ReservedNames),
		gc.WithWhitelist(cfg.Containers),
	)
	if cfg.Once {
		collector.Collect(ctx)
	} else {
		log.Info().
			Strs("ignore-containers", cfg.Containers).
			Strs("ignore-images", cfg.Images).
			Dur("interval", cfg.Interval).
			Msg("starting the garbage collector")

		gc.Schedule(ctx, collector, cfg.Interval)
	}
}

func initLogger(cfg *config) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if cfg.Pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !cfg.Color,
			},
		)
	}
}
