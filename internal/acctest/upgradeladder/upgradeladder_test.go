// Copyright © 2026 Ping Identity Corporation

package upgradeladder_test

import (
	"testing"

	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/upgradeladder"
)

// TestResourceFilterSkips verifies the PINGFEDERATE_UPGRADE_RESOURCES filter
// excludes non-matching resources via skip (not failure). Without TF_ACC the
// filter is the only skip path exercised before the TF_ACC gate.
func TestResourceFilterSkips(t *testing.T) {
	t.Setenv(upgradeladder.EnvResourceFilter, "something_else")
	upgradeladder.RunUpgradeLadder(t, upgradeladder.Spec{
		ResourceType: "pingfederate_session_settings",
		HCL: func() string {
			return `
resource "pingfederate_session_settings" "example" {
}
`
		},
	})
}
