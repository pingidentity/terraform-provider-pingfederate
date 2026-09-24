// Copyright © 2026 Ping Identity Corporation

package upgradeladder

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestParseLane(t *testing.T) {
	// Parse the *latest* and *earliest* supported patch versions from the real
	// supported list — no per-lane variables to maintain when versions change.
	latestVersion := latestSupportedVersion()
	latestLane := string(version.MajorMinor(latestVersion))
	earliestVersion := string(earliestSupportedVersion())
	earliestParts := strings.Split(string(earliestSupportedVersion()), ".")
	// A lane below the oldest supported one is unsupported; derive it from the
	// supported list (one minor below the oldest supported lane) rather than a
	// hand-written version.
	belowOldestLane := earliestParts[0] + "." + mustDecrement(t, earliestParts[1]) + ".4"

	tests := []struct {
		name           string
		productVersion string
		want           string
		wantErr        bool
	}{
		{"two digits", latestLane, latestLane, false},
		{"three digits", string(latestVersion), latestLane, false},
		{"earliest supported patch", earliestVersion, string(version.MajorMinor(earliestSupportedVersion())), false},
		{"four digits errors (provider contract)", string(latestVersion) + ".5", "", true},
		{"untrimmed errors (provider contract)", " " + latestLane + " ", "", true},
		{"unsupported lane below oldest", belowOldestLane, "", true},
		{"empty", "", "", true},
		{"one digit", earliestParts[0], "", true},
	}
	for _, tt := range tests {
		t.Run(tt.productVersion, func(t *testing.T) {
			got, err := ParseLane(tt.productVersion)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseLane(%q) error = %v, wantErr %v", tt.productVersion, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("ParseLane(%q) = %q, want %q", tt.productVersion, got, tt.want)
			}
		})
	}

	// Every supported two-digit lane must parse to itself.
	for _, lane := range version.SupportedMajorMinorVersions() {
		t.Run("lane "+string(lane), func(t *testing.T) {
			got, err := ParseLane(string(lane))
			if err != nil {
				t.Fatalf("ParseLane(%q) unexpected error: %v", lane, err)
			}
			if got != string(lane) {
				t.Fatalf("ParseLane(%q) = %q, want %q", lane, got, lane)
			}
		})
	}
}

// latestSupportedVersion returns the newest patch version in the supported
// list.
func latestSupportedVersion() version.SupportedVersion {
	sortedVersions := version.SupportedVersions()
	return sortedVersions[len(sortedVersions)-1]
}

// earliestSupportedVersion returns the oldest supported patch version.
func earliestSupportedVersion() version.SupportedVersion {
	return version.SupportedVersions()[0]
}

// TestLadderTableCoversSupportedLanes ties the ladder table to
// internal/version: every PingFederate major.minor supported by the current
// provider must have at least one ladder rung. When a new version is added
// to internal/version, this test fails until release rows are appended to
// ladderTable — the mitigation for the manual step, not a silent skip.
func TestLadderTableCoversSupportedLanes(t *testing.T) {
	for _, majorMinorVersion := range version.SupportedMajorMinorVersions() {
		lane := string(majorMinorVersion)
		if _, ok := laneFloorVersion(lane); !ok {
			t.Errorf("PingFederate lane '%s' is supported in internal/version but has no ladder table entry: "+
				"verify which released provider versions support it (git show v<tag>:internal/version/version.go) "+
				"and append its rows to ladderTable", lane)
		}
	}
	// Reverse check: every table row's MaxPFMinor must be a lane internal/version supports.
	for _, entry := range ladderTable {
		if !version.IsValid(string(entry.MaxPFMinor)) {
			t.Errorf("ladder rung %s references MaxPFMinor '%s', which is not supported in internal/version — "+
				"the lane was likely dropped as EOL; remove the stale rows", entry.ProviderVersion, entry.MaxPFMinor)
		}
	}
}

func TestBuildLadder(t *testing.T) {
	// Every supported lane derives its expected ladder directly from the
	// ladder table via BuildLadder — the test pins the structural contract
	// (order, floor start, local termination) rather than a per-lane copy.
	for _, majorMinorVersion := range version.SupportedMajorMinorVersions() {
		lane := string(majorMinorVersion)
		t.Run(lane, func(t *testing.T) {
			got, err := BuildLadder(lane)
			if err != nil {
				t.Fatalf("BuildLadder(%q) unexpected error: %v", lane, err)
			}
			if len(got) < 2 {
				t.Fatalf("BuildLadder(%q) = %v, want at least 2 rungs", lane, got)
			}
			if got[len(got)-1] != LocalRung {
				t.Fatalf("BuildLadder(%q) must end with LocalRung, got %v", lane, got)
			}

			// The floor must be the first table entry for the lane, and rungs
			// must run oldest -> newest through the last table row.
			floor, _ := laneFloorVersion(lane)
			if got[0] != floor {
				t.Fatalf("BuildLadder(%q) must start at the lane floor %q, got %q", lane, floor, got[0])
			}
			tableRungs := make([]string, 0, len(ladderTable))
			seenFloor := false
			for _, entry := range ladderTable {
				if !seenFloor {
					if string(version.MajorMinor(entry.MaxPFMinor)) != lane {
						continue
					}
					seenFloor = true
				}
				tableRungs = append(tableRungs, entry.ProviderVersion)
			}
			want := append(tableRungs, LocalRung)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("BuildLadder(%q) = %v, want %v", lane, got, want)
			}

			// Every registry rung must be a strict increase over the previous
			// (the ladder climbs one release at a time).
			for i := 1; i < len(got)-1; i++ {
				if compareRungs(got[i], got[i-1]) <= 0 {
					t.Fatalf("BuildLadder(%q) rungs not strictly increasing at %q: %v", lane, got[i], got)
				}
			}
		})
	}

	if _, err := BuildLadder("11.3"); err == nil {
		t.Fatal("BuildLadder for an EOL lane absent from internal/version should error")
	}
}

func TestBuildLadderFromEnv(t *testing.T) {
	// Lanes derive from the real supported list; expected ladders from the
	// ladder table, as in TestBuildLadder.
	newestLane := string(version.MajorMinor(latestSupportedVersion()))
	oldestLane := string(version.MajorMinor(earliestSupportedVersion()))
	fullNewest, err := BuildLadder(newestLane)
	if err != nil {
		t.Fatalf("BuildLadder(%q) unexpected error: %v", newestLane, err)
	}
	fullOldest, err := BuildLadder(oldestLane)
	if err != nil {
		t.Fatalf("BuildLadder(%q) unexpected error: %v", oldestLane, err)
	}
	last2 := fullNewest[len(fullNewest)-2:]

	tests := []struct {
		name    string
		env     string
		lane    string
		want    []string
		wantErr bool
	}{
		{"default full", "", newestLane, fullNewest, false},
		{"explicit full", "full", oldestLane, fullOldest, false},
		{"last2 smoke", "last2", newestLane, last2, false},
		{"explicit list", "1.3.0,1.9.0", newestLane, []string{"1.3.0", "1.9.0", "local"}, false},
		{"explicit with local", "1.9.0,local", newestLane, []string{"1.9.0", "local"}, false},
		{"unknown rung rejected", "1.2.0,local", newestLane, nil, true},
		{"empty part rejected", "1.9.0,,local", newestLane, nil, true},
		{"unknown lane", "", "11.3", nil, true},
		{"misordered rejected", "1.9.0,1.3.0", newestLane, nil, true},
		{"duplicate rejected", "1.9.0,1.9.0", newestLane, nil, true},
		{"local not last rejected", "1.3.0,local,1.9.0", newestLane, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvLadderOverride, tt.env)
			got, err := BuildLadderFromEnv(tt.lane)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildLadderFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("BuildLadderFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClampRungsForResource(t *testing.T) {
	rungs := []string{"1.3.0", "1.6.2", "1.9.0", "local"}

	if got := clampRungsForResource(rungs, ""); !reflect.DeepEqual(got, rungs) {
		t.Fatalf("empty LadderFloor should not clamp, got %v", got)
	}
	if got := clampRungsForResource(rungs, "1.6.2"); !reflect.DeepEqual(got, []string{"1.6.2", "1.9.0", "local"}) {
		t.Fatalf("LadderFloor 1.6.2 should drop earlier rungs, got %v", got)
	}
	// Unknown floor: keep the ladder (loud failure) rather than skip.
	if got := clampRungsForResource(rungs, "2.0.0"); !reflect.DeepEqual(got, rungs) {
		t.Fatalf("unknown LadderFloor should keep rungs, got %v", got)
	}
}

func TestValidateSpec(t *testing.T) {
	validHCL := func() string {
		return `
resource "pingfederate_example_resource" "example" {
}
`
	}

	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: validHCL}); err != nil {
		t.Fatalf("valid Spec rejected: %v", err)
	}
	if err := ValidateSpec(Spec{ResourceType: "", HCL: validHCL}); err == nil {
		t.Error("empty ResourceType should be rejected")
	}
	if err := ValidateSpec(Spec{ResourceType: "not_prefixed", HCL: validHCL}); err == nil {
		t.Error("missing pingfederate_ prefix should be rejected")
	}
	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: nil}); err == nil {
		t.Error("nil HCL should be rejected")
	}
	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: func() string {
		return `
resource "pingfederate_example_resource" "other" {
}
`
	}}); err == nil {
		t.Error("non-example resource label should be rejected")
	}
	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: validHCL, Allowlist: []AllowlistEntry{{AttributePath: "foo"}}}); err == nil {
		t.Error("allowlist entry without Reason should be rejected")
	}
	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: validHCL, LadderFloor: "1.4.5"}); err != nil {
		t.Fatalf("known ladder rung as LadderFloor should pass: %v", err)
	}
	if err := ValidateSpec(Spec{ResourceType: "pingfederate_example_resource", HCL: validHCL, LadderFloor: "1.45"}); err == nil {
		t.Error("typo'd LadderFloor should be rejected (it would silently disable clamping)")
	}
}

func TestResourceFilterMatches(t *testing.T) {
	cases := []struct {
		resourceType string
		envValue     string
		want         bool
	}{
		{"pingfederate_oauth_server_settings", "", true},
		{"pingfederate_oauth_server_settings", "oauth_server_settings", true},
		{"pingfederate_oauth_server_settings", "oauth_server_settings,session_settings", true},
		{"pingfederate_session_settings", "oauth_server_settings,session_settings", true},
		{"pingfederate_session_settings", "oauth_server_settings", false},
		{"pingfederate_session_settings", " does_not_exist ", false},
		{"pingfederate_session_settings", ",,", true}, // empty parts ignored -> select all
	}
	for _, tc := range cases {
		if got := ResourceFilterMatches(tc.resourceType, tc.envValue); got != tc.want {
			t.Errorf("ResourceFilterMatches(%q, %q) = %v, want %v", tc.resourceType, tc.envValue, got, tc.want)
		}
	}
}

func TestStripTopLevelDataSourceBlocks(t *testing.T) {
	resourceType := "pingfederate_session_settings"
	withData := `
resource "pingfederate_session_settings" "example" {
}
data "pingfederate_session_settings" "example" {
  depends_on = [pingfederate_session_settings.example]
}
`
	stripped := stripTopLevelDataSourceBlocks(t, resourceType, withData)
	if strings.Contains(stripped, `data "pingfederate_session_settings"`) {
		t.Errorf("self data block should be stripped, got:\n%s", stripped)
	}
	if !strings.Contains(stripped, `resource "pingfederate_session_settings" "example"`) {
		t.Errorf("resource block must remain, got:\n%s", stripped)
	}

	// Other pingfederate data sources are NOT stripped (fail loudly at the rung).
	otherData := `
resource "pingfederate_session_settings" "example" {
}
data "pingfederate_oauth_server_settings" "other" {
}
`
	if !strings.Contains(stripTopLevelDataSourceBlocks(t, resourceType, otherData), `data "pingfederate_oauth_server_settings"`) {
		t.Error("non-self data blocks must not be stripped")
	}

	// Unparseable HCL passes through untouched.
	garbage := "not hcl at all {{{"
	if stripTopLevelDataSourceBlocks(t, resourceType, garbage) != garbage {
		t.Error("unparseable HCL should pass through unchanged")
	}
}
