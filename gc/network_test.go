// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package gc

import (
	"context"
	"testing"

	"github.com/drone/drone-gc/mocks"

	"docker.io/go-docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestCollectNetworks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockReport := types.NetworksPruneReport{
		NetworksDeleted: []string{"foo", "bar"},
	}
	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().NetworksPrune(gomock.Any(), networkPruneArgs).Return(mockReport, nil)

	c := New(client).(*collector)
	err := c.collectNetworks(context.Background())
	if err != nil {
		t.Error(err)
	}
}
