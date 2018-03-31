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
			Created: 359596800,
		},
		// skip whitelisted name
		{
			ID:      "2b8fd9751c4c",
			Names:   []string{"foo", "bar"},
			State:   "exited",
			Created: 359596800,
		},
		// skip drone images
		{
			ID:      "4e38e38c8ce0",
			Names:   []string{"bar"},
			Image:   "drone/drone:latest",
			State:   "exited",
			Created: 359596800,
		},
		// skip non-exited containers
		{
			ID:      "a180b24e38ed",
			Names:   []string{"bar"},
			State:   "created",
			Created: 359596800,
		},
		// skip recently created containers
		{
			ID:      "481995377a04",
			Names:   []string{"bar"},
			State:   "exited",
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
			Created: 359596800,
		},
		{
			ID:      "2b8fd9751c4c",
			Names:   []string{"foo"},
			State:   "exited",
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
