// Copyright © 2026 Ping Identity Corporation

package upgradeladder

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// stripTopLevelDataSourceBlocks removes the resource-under-test's own trailing
// `data "pingfederate_..." "example"` block from the generated HCL: data
// sources re-read on every plan, which the rungs' empty-plan checks treat as
// drift. Only the resource under test's self data block is stripped; other
// data sources remain and fail the run loudly.
func stripTopLevelDataSourceBlocks(t *testing.T, resourceType string, hclSource string) string {
	t.Helper()
	file, parseDiags := hclwrite.ParseConfig([]byte(hclSource), "ladder_hcl.tf", hcl.Pos{Byte: 0, Line: 1, Column: 1})
	if parseDiags.HasErrors() {
		// Unparseable HCL is left untouched: the terraform rung will fail
		// with its own (clearer) error.
		t.Logf("upgradeladder: could not parse HCL for data-block stripping: %v", parseDiags)
		return hclSource
	}
	stripped := false
	for _, block := range file.Body().Blocks() {
		if block.Type() != "data" || len(block.Labels()) != 2 {
			continue
		}
		if block.Labels()[0] == resourceType && block.Labels()[1] == "example" {
			file.Body().RemoveBlock(block)
			stripped = true
		}
	}
	if stripped {
		return string(file.Bytes())
	}
	return hclSource
}

// RunUpgradeLadder runs the provider version ladder for one resource:
// apply the Spec's HCL with each registry rung from the lane floor, then with
// the local build, using the identical configuration at every rung. Every
// rung asserts an empty post-apply plan (the built-in helper/resource guard);
// a hop that produces plan drift — the "Provider produced inconsistent result
// after apply" class — fails the step.
//
// Rungs older than the Spec's LadderFloor are dropped, so resources that
// can't run on the lane's oldest rungs start their ladder further up.
// Resources are filtered by PINGFEDERATE_UPGRADE_RESOURCES (comma-separated
// substrings; empty selects all).
func RunUpgradeLadder(t *testing.T, spec Spec) {
	t.Helper()
	// Wrapped in a subtest: t.Skip/t.Fatal call runtime.Goexit, which would
	// otherwise unwind this whole goroutine and skip the caller's next
	// statement — every gen_test.go calls RunUpgradeLadder then
	// RunServerUpgradeLadder sequentially, and the two axes must run
	// independently of each other's outcome.
	t.Run("provider", func(t *testing.T) {
		if !Enabled {
			t.Skipf("skipping upgrade ladder for %s: build without the 'upgradeladder' tag (make testupgradeacc)", spec.ResourceType)
			return
		}

		run, err := prepareLadderRun(t, spec)
		if err != nil {
			if _, ok := err.(skipSignal); ok {
				t.Skipf("%v", err)
				return
			}
			t.Fatalf("%v", err)
		}
		run(t)
	})
}

// skipSignal marks conditions that skip (not fail) a ladder run: filter
// exclusion, degenerate rung lists, and server probes.
type skipSignal struct{ cause error }

func (s skipSignal) Error() string { return s.cause.Error() }

// preflightSpec runs the checks both ladder axes share, in the required
// order: the probe before ValidateSpec, since some probes assign package
// variables the HCL functions ValidateSpec evaluates depend on (e.g. the
// certificate resources' file data). Returns a skipSignal for a probe
// rejection, or ValidateSpec's error otherwise.
func preflightSpec(spec Spec) error {
	if spec.ClusterModeProbe != nil {
		if err := spec.ClusterModeProbe(); err != nil {
			return skipSignal{fmt.Errorf("server probe says the resource cannot apply here: %v", err)}
		}
	}
	return ValidateSpec(spec)
}

// prepareLadderRun resolves everything needed before any Terraform call:
// filter, TF_ACC, probe/validation, lane, rung list, and HCL preparation.
// Errors distinguish skip (skipSignal) from hard failure.
func prepareLadderRun(t *testing.T, spec Spec) (func(*testing.T), error) {
	if !ResourceFilterMatches(spec.ResourceType, os.Getenv(EnvResourceFilter)) {
		return nil, skipSignal{fmt.Errorf("skipping upgrade ladder for %s: excluded by %s filter", spec.ResourceType, EnvResourceFilter)}
	}
	// Checked before the probe (a live HTTP call for some resources) so a
	// plain `go test -tags upgradeladder` run without TF_ACC skips instead
	// of making network calls it has no need to make.
	if os.Getenv("TF_ACC") == "" {
		return nil, skipSignal{fmt.Errorf("skipping upgrade ladder for %s: Acceptance tests skipped unless env 'TF_ACC' set", spec.ResourceType)}
	}

	if err := preflightSpec(spec); err != nil {
		if _, ok := err.(skipSignal); ok {
			return nil, skipSignal{fmt.Errorf("skipping upgrade ladder for %s: %v", spec.ResourceType, err)}
		}
		return nil, fmt.Errorf("invalid upgrade-ladder Spec: %w", err)
	}

	lane, err := ParseLane(os.Getenv("PINGFEDERATE_PROVIDER_PRODUCT_VERSION"))
	if err != nil {
		// TF_ACC is confirmed set above, so a bad or missing product version
		// here is a setup error, not a skip.
		return nil, fmt.Errorf("failed to determine PingFederate lane for the upgrade ladder: %w", err)
	}

	rungs, err := BuildLadderFromEnv(lane)
	if err != nil {
		return nil, fmt.Errorf("failed to build upgrade ladder: %w", err)
	}

	rungs = clampRungsForResource(rungs, spec.LadderFloor)
	if len(rungs) < 2 {
		return nil, skipSignal{fmt.Errorf("skipping upgrade ladder for %s: fewer than 2 rungs after clamping (rungs: %v)", spec.ResourceType, rungs)}
	}

	hcl := stripTopLevelDataSourceBlocks(t, spec.ResourceType, spec.HCL())

	t.Logf("Running %d-rung upgrade ladder for %s: %s", len(rungs), spec.ResourceType, strings.Join(rungs, " -> "))

	steps := make([]resource.TestStep, 0, len(rungs))
	for _, rung := range rungs {
		step := resource.TestStep{
			// The identical config at every rung: any plan delta between rungs
			// is version-induced drift, not config change.
			Config: hcl,
		}
		if rung == LocalRung {
			step.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
				"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
			}
		} else {
			step.ExternalProviders = map[string]resource.ExternalProvider{
				"pingfederate": {
					Source:            "registry.terraform.io/pingidentity/pingfederate",
					VersionConstraint: "=" + rung,
				},
			}
		}
		steps = append(steps, step)
	}

	return func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck: func() { acctest.ConfigurationPreCheck(t) },
			Steps:    steps,
		})
	}, nil
}

// clampRungsForResource drops rungs older than the Spec's LadderFloor. An
// empty floor means every supported rung, the default.
func clampRungsForResource(rungs []string, ladderFloor string) []string {
	if ladderFloor == "" {
		return rungs
	}
	floorIdx := -1
	for i, rung := range rungs {
		if rung == ladderFloor {
			floorIdx = i
			break
		}
	}
	if floorIdx == -1 {
		// Unknown or newer-than-ladder floor: keep the ladder intact so the
		// failure is loud rather than silently skipping rungs.
		return rungs
	}
	return rungs[floorIdx:]
}
