// Copyright © 2026 Ping Identity Corporation

package upgradeladder

import (
	"fmt"
	"strings"
)

// EnvResourceFilter optionally restricts a ladder run to resources whose type
// contains one of the comma-separated substrings, e.g.
// "oauth_server_settings,incoming_proxy_settings".
const EnvResourceFilter = "PINGFEDERATE_UPGRADE_RESOURCES"

// Spec declares one resource's participation in the provider version ladder.
// Each resource's *_gen_test.go under internal/acctest/config declares its
// Spec in a TestUpgradeLadder_<ResourceName> function, reusing the package's
// existing generated HCL funcs — there is no central registry and no frozen
// HCL copy. Ladder tests run only when the build sets the 'upgradeladder'
// tag (see Enabled); otherwise RunUpgradeLadder skips.
type Spec struct {
	// ResourceType is the Terraform resource type name,
	// e.g. "pingfederate_oauth_server_settings".
	ResourceType string

	// HCL returns the configuration applied at every rung. Usually the
	// package's generated MinimalHCL (plus dependency HCL when the resource
	// requires prerequisite resources); the harness strips top-level
	// data-source blocks automatically. The configuration must apply cleanly
	// with the oldest rung of every lane the resource supports (see
	// AvailableSince).
	HCL func() string

	// Allowlist is reserved for a future phase that tolerates attribute-level
	// plan deltas (e.g. server-injected defaults); entries require a JIRA
	// reference in Reason. Rungs currently assert strictly empty plans.
	Allowlist []AllowlistEntry

	// AvailableSince is the provider version the resource first shipped in,
	// e.g. "1.6.0". Rungs older than this are dropped from the resource's
	// ladder. Empty means the resource exists in every supported rung.
	AvailableSince string

	// ClusterModeProbe optionally probes the live server and returns a
	// non-nil error when the resource cannot apply on it (e.g. cluster
	// settings on a standalone server return 403 not_in_clustered_mode, or a
	// required env var is missing). A non-nil probe error skips the ladder
	// for the resource. Nil means no probe.
	//
	// The probe runs before ValidateSpec evaluates HCL, so it may also
	// perform setup the HCL functions depend on — e.g. assigning the package
	// variable from the env it just checked. Probes that skip via a missing
	// env var must be idempotent: both ladders call the same probe.
	ClusterModeProbe func() error
}

// AllowlistEntry is reserved for phase 2 (attribute-level plan-drift
// tolerances); see Spec.Allowlist.
type AllowlistEntry struct {
	// AttributePath is the dotted attribute path the delta is tolerated on.
	AttributePath string
	// Reason must cite the JIRA case or explain why the delta is expected.
	Reason string
}

// ValidateSpec reports whether a Spec is usable by the ladder harness.
func ValidateSpec(spec Spec) error {
	if spec.ResourceType == "" {
		return fmt.Errorf("Spec.ResourceType is required")
	}
	if !strings.HasPrefix(spec.ResourceType, "pingfederate_") {
		return fmt.Errorf("Spec.ResourceType %q must start with 'pingfederate_'", spec.ResourceType)
	}
	if spec.HCL == nil {
		return fmt.Errorf("Spec.HCL is required for %s", spec.ResourceType)
	}
	if !strings.Contains(spec.HCL(), `"`+spec.ResourceType+`" "example"`) {
		return fmt.Errorf("HCL for %s must use the resource label 'example'", spec.ResourceType)
	}
	if spec.AvailableSince != "" && !isKnownRung(spec.AvailableSince) {
		return fmt.Errorf("Spec.AvailableSince %q for %s is not a ladder rung (valid: %s); an unknown value silently disables clamping and the oldest rungs then fail at runtime",
			spec.AvailableSince, spec.ResourceType, ladderTableVersions())
	}
	for _, entry := range spec.Allowlist {
		if strings.TrimSpace(entry.Reason) == "" {
			return fmt.Errorf("allowlist entry for %s attribute %q is missing a Reason (cite the JIRA case)", spec.ResourceType, entry.AttributePath)
		}
	}
	return nil
}

// ResourceFilterMatches reports whether resourceType matches the
// comma-separated substring filter. A filter with no non-empty parts (e.g. ""
// or ",,") selects all.
func ResourceFilterMatches(resourceType string, envValue string) bool {
	hasSubstring := false
	for _, substr := range strings.Split(envValue, ",") {
		trimmed := strings.TrimSpace(substr)
		if trimmed == "" {
			continue
		}
		hasSubstring = true
		if strings.Contains(resourceType, trimmed) {
			return true
		}
	}
	return !hasSubstring
}
