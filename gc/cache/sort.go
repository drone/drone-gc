// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package cache

import (
	"docker.io/go-docker/api/types"
)

type byCreated []*types.ImageSummary

func (a byCreated) Len() int           { return len(a) }
func (a byCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreated) Less(i, j int) bool { return a[i].Created < a[j].Created }

type byCreatedAsc []*types.ImageSummary

func (a byCreatedAsc) Len() int           { return len(a) }
func (a byCreatedAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreatedAsc) Less(i, j int) bool { return a[i].Created > a[j].Created }
