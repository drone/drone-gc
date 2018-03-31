// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drone/drone-gc/mocks"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/volume"
	"github.com/golang/mock/gomock"
)

func TestCollectVolumes(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockVolumes := volume.VolumesListOKBody{
		Volumes: []*types.Volume{
			{Name: "a180b24e38ed", Driver: "local", CreatedAt: "2018-01-01T00:00:00Z"},
			{Name: "bfbf8512f21e", Driver: "local", CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z")},
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().VolumeList(gomock.Any(), volumeListArgs).Return(mockVolumes, nil)
	client.EXPECT().VolumeRemove(gomock.Any(), mockVolumes.Volumes[0].Name, false).Return(nil)

	c := New(client).(*collector)
	err := c.collectVolumes(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestCollectVolumes_MultiError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockVolumes := volume.VolumesListOKBody{
		Volumes: []*types.Volume{
			{Name: "a180b24e38ed", Driver: "local", CreatedAt: "2018-01-01T00:00:00Z"},
			{Name: "bfbf8512f21e", Driver: "local", CreatedAt: "2018-01-01T00:00:00Z"},
		},
	}
	mockErr := errors.New("cannot remove container")

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().VolumeList(gomock.Any(), volumeListArgs).Return(mockVolumes, nil)
	client.EXPECT().VolumeRemove(gomock.Any(), mockVolumes.Volumes[0].Name, false).Return(mockErr)
	client.EXPECT().VolumeRemove(gomock.Any(), mockVolumes.Volumes[1].Name, false).Return(nil)

	c := New(client).(*collector)
	err := c.collectVolumes(context.Background())
	if err == nil {
		t.Errorf("Expected multi-error returned")
	}
}
