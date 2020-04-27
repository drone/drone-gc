// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gc

import (
	"reflect"
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	expectedMinImageAge, _ := time.ParseDuration("5h")
	c := New(nil,
		WithImageWhitelist([]string{"foo"}),
		WithThreshold(42),
		WithWhitelist([]string{"bar"}),
		WithMinImageAge(expectedMinImageAge),
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

	if got, want := c.minImageAge, expectedMinImageAge; !reflect.DeepEqual(want, got) {
		t.Errorf("Want minImageAge %v, got %v", want, got)
	}
}
