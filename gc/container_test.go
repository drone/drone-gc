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

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestCollectContainers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockContainers := []types.Container{
		{
			ID:      "c3d2a6307f4e",
			Names:   []string{"bar"},
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": "915148800"},
			Created: 359596800,
		},
		// skip missing io.drone.expires label
		{
			ID:      "2b8fd9751c4c",
			Names:   []string{"foo", "bar"},
			State:   "exited",
			Labels:  map[string]string{},
			Created: 359596800,
		},
		// skip whitelisted name
		{
			ID:      "2b8fd9751c4c",
			Names:   []string{"foo", "bar"},
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": "915148800"},
			Created: 359596800,
		},
		// skip drone images
		{
			ID:      "4e38e38c8ce0",
			Names:   []string{"bar"},
			Image:   "drone/drone:latest",
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": "915148800"},
			Created: 359596800,
		},
		// skip recently created containers
		{
			ID:      "481995377a04",
			Names:   []string{"bar"},
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": fmt.Sprint(time.Now().Add(time.Hour).Unix())},
			Created: time.Now().Unix(),
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().ContainerList(gomock.Any(), containerListArgs).Return(mockContainers, nil)
	client.EXPECT().ContainerRemove(gomock.Any(), mockContainers[0].ID, containerRemoveOpts).Return(nil)

	c := New(client,
		WithWhitelist([]string{"foo"}),
	).(*collector)
	err := c.collectContainers(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestCollectContainers_MultiError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockContainers := []types.Container{
		{
			ID:      "c3d2a6307f4e",
			Names:   []string{"bar"},
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": "915148800"},
			Created: 359596800,
		},
		{
			ID:      "2b8fd9751c4c",
			Names:   []string{"foo"},
			State:   "exited",
			Labels:  map[string]string{"io.drone.expires": "915148800"},
			Created: 359596800,
		},
	}
	mockErr := errors.New("cannot remove contianer")

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().ContainerList(gomock.Any(), containerListArgs).Return(mockContainers, nil)
	client.EXPECT().ContainerRemove(gomock.Any(), mockContainers[0].ID, containerRemoveOpts).Return(mockErr)
	client.EXPECT().ContainerRemove(gomock.Any(), mockContainers[1].ID, containerRemoveOpts).Return(nil)

	c := New(client).(*collector)
	err := c.collectContainers(context.Background())
	if err == nil {
		t.Errorf("Expected multi-error returned")
	}
}
