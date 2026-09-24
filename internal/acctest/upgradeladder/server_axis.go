// Copyright © 2026 Ping Identity Corporation

package upgradeladder

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// EnvServerLanesOverride points RunServerUpgradeLadder at the PingFederate
// lanes to climb, as a comma-separated list in ascending order, e.g.
// "12.2,12.3,13.0,13.1". Empty selects every supported lane. One live server
// per lane must be reachable at the matching entry of PINGFEDERATE_UPGRADE_LANE_HOSTS.
const EnvServerLanesOverride = "PINGFEDERATE_UPGRADE_SERVER_LANES"

// EnvServerLaneHosts is the comma-separated list of https_host values for
// every lane beyond the first; the first lane uses PINGFEDERATE_PROVIDER_HTTPS_HOST.
const EnvServerLaneHosts = "PINGFEDERATE_UPGRADE_LANE_HOSTS"

// EnvServerLadderEnabled opts in to the server-upgrade ladder: the provider
// ladder's make target and CI run the same test functions, and without this
// gate they would need lane containers they don't have. Set to any non-empty
// value (make testserverupgradeacc does).
const EnvServerLadderEnabled = "PINGFEDERATE_UPGRADE_SERVER_LADDER"

// httpsHostEnvVar is the first lane's server host variable.
const httpsHostEnvVar = "PINGFEDERATE_PROVIDER_HTTPS_HOST"

// serverLanes resolves the lane list: PINGFEDERATE_UPGRADE_SERVER_LANES if
// set, else every supported lane in ascending order (from internal/version).
// Each value is validated through version.Parse exactly like the provider's
// own configure-time check, and the list must ascend — the ladder steps
// forward.
func serverLanes() ([]string, error) {
	override := strings.TrimSpace(os.Getenv(EnvServerLanesOverride))
	if override == "" {
		majorMinors := version.SupportedMajorMinorVersions()
		lanes := make([]string, 0, len(majorMinors))
		for _, lane := range majorMinors {
			lanes = append(lanes, string(lane))
		}
		return lanes, nil
	}
	lanes := make([]string, 0, len(override))
	for _, part := range strings.Split(override, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			return nil, fmt.Errorf("empty lane in %s override '%s'", EnvServerLanesOverride, override)
		}
		parsed, diags := version.Parse(trimmed)
		if diags.HasError() {
			return nil, fmt.Errorf("lane '%s' in %s is not a supported PingFederate version (supported: %s)",
				trimmed, EnvServerLanesOverride, strings.Join(supportedLanes(), ", "))
		}
		lanes = append(lanes, string(version.MajorMinor(parsed)))
	}
	for i := 1; i < len(lanes); i++ {
		if compareRungs(lanes[i], lanes[i-1]) <= 0 {
			return nil, fmt.Errorf("%s must be ascending, got '%s' at/after '%s'", EnvServerLanesOverride, lanes[i], lanes[i-1])
		}
	}
	return lanes, nil
}

// laneHosts resolves the per-lane https_host values: the first lane reuses
// PINGFEDERATE_PROVIDER_HTTPS_HOST; PINGFEDERATE_UPGRADE_LANE_HOSTS must list
// one host for every lane beyond the first.
func laneHosts(lanes []string) ([]string, error) {
	baseHost := strings.TrimSpace(os.Getenv(httpsHostEnvVar))
	if baseHost == "" {
		return nil, fmt.Errorf("%s must be set (the first lane's server host)", httpsHostEnvVar)
	}
	hosts := []string{baseHost}
	extra := strings.TrimSpace(os.Getenv(EnvServerLaneHosts))
	if extra == "" {
		if len(lanes) > 1 {
			return nil, fmt.Errorf("%s must list one host per lane beyond the first (%d lanes requested: %s)",
				EnvServerLaneHosts, len(lanes), strings.Join(lanes, ","))
		}
		return hosts, nil
	}
	for _, part := range strings.Split(extra, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			return nil, fmt.Errorf("empty host in %s override '%s'", EnvServerLaneHosts, extra)
		}
		hosts = append(hosts, trimmed)
	}
	if len(hosts) != len(lanes) {
		return nil, fmt.Errorf("%s has %d hosts for %d lanes (%s covers lanes beyond the first): %s",
			EnvServerLaneHosts, len(hosts), len(lanes), httpsHostEnvVar, strings.Join(hosts, ","))
	}
	return hosts, nil
}

// laneProviderBlock renders the provider block a server-axis step uses to
// reach one lane's server. The provider is the local build in every step;
// only its target host and product version change.
func laneProviderBlock(host string, lane string) string {
	return fmt.Sprintf(`
provider "pingfederate" {
  https_host      = %q
  product_version = %q
}
`, host, lane)
}

// RunServerUpgradeLadder runs the PingFederate server-upgrade ladder for one
// resource: apply the Spec's HCL against the first lane's server, then step
// the server to each subsequent lane while the Terraform state carries
// forward, re-pointing the provider at each lane's host via the step's
// provider block. The provider stays the local build throughout — this
// ladder exercises the server-upgrade axis, where version-gated provider
// behavior and server-injected defaults change under existing state.
//
// Every step asserts a strictly empty post-apply plan (plugin-testing's
// built-in guard; ExpectNonEmptyPlan cannot express drift tolerance — it is a
// hard XOR). That strictness is the point: a non-empty plan right after an
// apply on a new server is precisely the "Provider produced inconsistent
// result after apply" signature this harness hunts, so drift fails loudly
// rather than being tolerated.
func RunServerUpgradeLadder(t *testing.T, spec Spec) {
	t.Helper()

	if !Enabled {
		t.Skipf("skipping server-upgrade ladder for %s: build without the 'upgradeladder' tag (make testserverupgradeacc)", spec.ResourceType)
		return
	}
	if os.Getenv(EnvServerLadderEnabled) == "" {
		t.Skipf("skipping server-upgrade ladder for %s: %s not set (make testserverupgradeacc) — the shared test function runs both ladders, and only the server-ladder target spins lane containers", spec.ResourceType, EnvServerLadderEnabled)
		return
	}
	if !ResourceFilterMatches(spec.ResourceType, os.Getenv(EnvResourceFilter)) {
		t.Skipf("skipping server-upgrade ladder for %s: excluded by %s filter", spec.ResourceType, EnvResourceFilter)
		return
	}
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	// The probe runs before Spec validation: ValidateSpec evaluates
	// spec.HCL(), and some probes assign the package variables those HCL
	// functions read (e.g. the certificate resources' file data).
	if spec.ClusterModeProbe != nil {
		if err := spec.ClusterModeProbe(); err != nil {
			t.Skipf("skipping server-upgrade ladder for %s: server probe says the resource cannot apply here: %v", spec.ResourceType, err)
			return
		}
	}

	if err := ValidateSpec(spec); err != nil {
		t.Fatalf("invalid server-upgrade ladder Spec: %v", err)
	}

	lanes, err := serverLanes()
	if err != nil {
		t.Fatalf("failed to determine PingFederate lanes: %v", err)
	}
	if len(lanes) < 2 {
		t.Skipf("skipping server-upgrade ladder for %s: needs at least 2 lanes, got: %v", spec.ResourceType, lanes)
		return
	}
	hosts, err := laneHosts(lanes)
	if err != nil {
		t.Fatalf("failed to resolve lane hosts: %v", err)
	}

	strippedHcl := stripTopLevelDataSourceBlocks(t, spec.ResourceType, spec.HCL())

	t.Logf("Running %d-lane server-upgrade ladder for %s: %s", len(lanes), spec.ResourceType, strings.Join(lanes, " -> "))

	steps := make([]resource.TestStep, 0, len(lanes))
	for i, lane := range lanes {
		step := resource.TestStep{
			// The identical resource config at every lane; only the provider
			// block changes, so any plan delta between lanes is
			// server-upgrade-induced, not config change.
			Config: laneProviderBlock(hosts[i], lane) + strippedHcl,
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
			},
		}
		steps = append(steps, step)
	}

	resource.Test(t, resource.TestCase{
		Steps: steps,
	})
}
