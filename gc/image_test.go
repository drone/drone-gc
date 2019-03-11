// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
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

func TestCollectImages(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		LayersSize: 850,
		Images: []*types.ImageSummary{
			{
				ID:         "a180b24e38ed",
				Created:    359596800,
				SharedSize: 50,
				Size:       300,
			},
			{
				ID:         "4e38e38c8ce0",
				Created:    359596800,
				SharedSize: 50,
				Size:       300,
			},
			// this image should not be removed since removal
			// of the above two images will put us below the
			// target threshold.
			{
				ID:         "481995377a04",
				Created:    359596800,
				SharedSize: 50,
				Size:       250,
			},
		},
	}
	mockImages := []types.ImageInspect{
		{ID: "a180b24e38ed"},
		{ID: "4e38e38c8ce0"},
		{ID: "481995377a04"},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)
	client.EXPECT().ImageInspectWithRaw(gomock.Any(), mockImages[0].ID).Return(mockImages[0], nil, nil)
	client.EXPECT().ImageInspectWithRaw(gomock.Any(), mockImages[1].ID).Return(mockImages[1], nil, nil)
	// we DO NOT inspect image 481995377a04

	client.EXPECT().ImageRemove(gomock.Any(), mockImages[0].ID, imageRemoveOpts).Return(nil, nil)
	client.EXPECT().ImageRemove(gomock.Any(), mockImages[1].ID, imageRemoveOpts).Return(nil, nil)
	// we DO NOT remove image 481995377a04

	c := New(client, WithThreshold(500)).(*collector)
	err := c.collectImages(context.Background())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that when an error is encountered we move to
// the next image in the list. Errors are aggregated and returned
// at the end of the loop.
func TestCollectImages_MutliError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		LayersSize: 1,
		Images: []*types.ImageSummary{
			{
				ID:      "a180b24e38ed",
				Created: 359596800,
			},
			{
				ID:      "4e38e38c8ce0",
				Created: 359596800,
			},
		},
	}
	mockImages := []types.ImageInspect{
		{ID: "a180b24e38ed"},
		{ID: "4e38e38c8ce0"},
	}
	mockError := errors.New("cannot remove container")

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)
	client.EXPECT().ImageInspectWithRaw(gomock.Any(), mockImages[0].ID).Return(mockImages[0], nil, nil)
	client.EXPECT().ImageInspectWithRaw(gomock.Any(), mockImages[1].ID).Return(mockImages[1], nil, nil)

	client.EXPECT().ImageRemove(gomock.Any(), mockImages[0].ID, imageRemoveOpts).Return(nil, mockError)
	client.EXPECT().ImageRemove(gomock.Any(), mockImages[1].ID, imageRemoveOpts).Return(nil, nil)

	c := New(client).(*collector)
	err := c.collectImages(context.Background())
	if err == nil {
		t.Errorf("Expect multi-error returned")
	}
}

// this test verifies that we do not purge the image cache
// if the cache is already below the target threshold.
func TestCollectImages_BelowThreshold(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		LayersSize: 1, // 1 byte
	}
	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)

	c := New(client, WithThreshold(2)).(*collector)
	err := c.collectImages(context.Background())
	if err != nil {
		t.Error(err)
	}
}

// this test verifies that we do not purge images that are
// in-use by the system or are newly created.
func TestCollectImages_Skip(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		LayersSize: 1,
		Containers: []*types.Container{
			{ImageID: "a180b24e38ed"},
		},
		Images: []*types.ImageSummary{
			// this image is in-use
			{
				ID:      "a180b24e38ed",
				Created: 359596800,
			},
			// this image is newly created
			{
				ID:      "481995377a04",
				Created: time.Now().Unix(),
			},
		},
	}
	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)

	c := New(client).(*collector)
	err := c.collectImages(context.Background())
	if err != nil {
		t.Error(err)
	}
}

// this test verifies that we do not purge images that
// are whitelisted by the user.
func TestCollectImages_SkipWhitelist(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockdf := types.DiskUsage{
		LayersSize: 1,
		Images: []*types.ImageSummary{
			{
				ID:      "a180b24e38ed",
				Created: 359596800,
			},
		},
	}

	mockImageInspect := types.ImageInspect{
		ID: "a180b24e38ed",
		RepoTags: []string{
			"drone/drone:1.0.0",
			"drone/drone:1.0",
			"drone/drone:1",
			"drone/drone:latest",
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().DiskUsage(gomock.Any()).Return(mockdf, nil)
	client.EXPECT().ImageInspectWithRaw(gomock.Any(), mockImageInspect.ID).Return(mockImageInspect, nil, nil)

	c := New(client,
		WithImageWhitelist(
			[]string{"drone/drone:*"},
		),
	).(*collector)
	err := c.collectImages(context.Background())
	if err != nil {
		t.Error(err)
	}
}
