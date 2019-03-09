// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package gc

import (
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {
	c := New(nil,
		WithImageWhitelist([]string{"foo"}),
		WithThreshold(42),
		WithWhitelist([]string{"bar"}),
	).(*collector)

	if got, want := c.threshold, int64(42); got != want {
		t.Errorf("Want cache threshold %d, got %d", want, got)
	}
	if got, want := c.whitelist, []string{"bar"}; !reflect.DeepEqual(want, got) {
		t.Errorf("Want container whitelist %v, got %v", want, got)
	}
	if got, want := c.reserved, []string{"foo"}; !reflect.DeepEqual(want, got) {
		t.Errorf("Want image whitelist %v, got %v", want, got)
	}
}
