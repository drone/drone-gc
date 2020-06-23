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

func TestCollectNetworks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockNetworks := []types.NetworkResource{
		{Name: "a180b24e38ed", Driver: "bridge", Labels: map[string]string{"io.drone.expires": "915148800"}},
		{Name: "e3d0f1751532", Driver: "bridge", Labels: map[string]string{"io.drone.expires": fmt.Sprint(time.Now().Add(time.Hour).Unix())}},
		{Name: "bfbf8512f21e", Driver: "bridge", Labels: nil},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().NetworkList(gomock.Any(), gomock.Any()).Return(mockNetworks, nil)
	client.EXPECT().NetworkRemove(gomock.Any(), mockNetworks[0].Name).Return(nil)

	c := New(client).(*collector)
	err := c.collectNetworks(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestCollectNetworks_MultiError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockNetworks := []types.NetworkResource{
		{Name: "a180b24e38ed", Driver: "bridge", Labels: map[string]string{"io.drone.expires": "915148800"}},
		{Name: "bfbf8512f21e", Driver: "bridge", Labels: map[string]string{"io.drone.expires": "915148800"}},
	}
	mockErr := errors.New("cannot remove network")

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().NetworkList(gomock.Any(), gomock.Any()).Return(mockNetworks, nil)
	client.EXPECT().NetworkRemove(gomock.Any(), mockNetworks[0].Name).Return(mockErr)
	client.EXPECT().NetworkRemove(gomock.Any(), mockNetworks[1].Name).Return(nil)

	c := New(client).(*collector)
	err := c.collectNetworks(context.Background())
	if err == nil {
		t.Errorf("Expected multi-error returned")
	}
}
