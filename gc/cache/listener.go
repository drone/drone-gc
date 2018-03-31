// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"context"
	"time"

	"github.com/drone/drone-gc/gc/internal"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
)

type listener struct {
	client docker.APIClient
	cache  *cache
}

func (l *listener) listen(ctx context.Context) error {
	// this is an infinite loop that only exites when
	// the context is cancelled (e.g. graceful shutdown).
	// we want to continuously re-connect to the docker
	// event stream if disconnected.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l.do(ctx)
		}
	}
}

func (l *listener) do(ctx context.Context) error {
	eventc, errc := l.client.Events(ctx, eventOpts)
	for {
		select {
		case err := <-errc:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case event := <-eventc:
			if event.Action == "create" && event.Type == "container" {
				name := internal.ExpandImage(event.From)
				l.cache.push(name, time.Now().Unix())
			}
		}
	}
}

var eventOpts = types.EventsOptions{
	Filters: filters.NewArgs(
		filters.KeyValuePair{
			Key:   "type",
			Value: "container",
		},
		filters.KeyValuePair{
			Key:   "event",
			Value: "create",
		},
	),
}
