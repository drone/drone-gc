// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"errors"
	"fmt"
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
			{Name: "a180b24e38ed", Driver: "local", Labels: map[string]string{"io.drone.expires": "915148800"}},
			{Name: "e3d0f1751532", Driver: "local", Labels: map[string]string{"io.drone.expires": fmt.Sprint(time.Now().Add(time.Hour).Unix())}},
			{Name: "bfbf8512f21e", Driver: "local", Labels: nil},
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
			{Name: "a180b24e38ed", Driver: "local", Labels: map[string]string{"io.drone.expires": "915148800"}},
			{Name: "bfbf8512f21e", Driver: "local", Labels: map[string]string{"io.drone.expires": "915148800"}},
		},
	}
	mockErr := errors.New("cannot remove volume")

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
