// Copyright © 2026 Ping Identity Corporation

package upgradeladder_test

import (
	"testing"

	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/upgradeladder"
)

// This file runs WITHOUT the 'upgradeladder' build tag, so Enabled is false
// and both ladder runners skip before reaching any logic — the runners' skip
// behavior under this tag is verified by the tagged targets
// (make testupgradeacc / testserverupgradeacc). What IS verifiable here is
// the exported pure logic those runners compose, e.g. the resource filter.
func TestResourceFilterMatches(t *testing.T) {
	cases := []struct {
		resourceType string
		envValue     string
		want         bool
	}{
		{"pingfederate_session_settings", "", true},
		{"pingfederate_session_settings", "something_else", false},
		{"pingfederate_session_settings", "something_else,session_settings", true},
		{"pingfederate_session_settings", " session_settings ,", true},
	}
	for _, tc := range cases {
		if got := upgradeladder.ResourceFilterMatches(tc.resourceType, tc.envValue); got != tc.want {
			t.Errorf("ResourceFilterMatches(%q, %q) = %v, want %v", tc.resourceType, tc.envValue, got, tc.want)
		}
	}
}
