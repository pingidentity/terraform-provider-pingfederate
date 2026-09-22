// Copyright © 2026 Ping Identity Corporation

package version

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestMajorMinor(t *testing.T) {
	// For every supported version, MajorMinor must equal the first two
	// segments — derived from the constant itself, no copies of the values.
	for _, supportedVersion := range getSortedVersions() {
		parts := strings.Split(string(supportedVersion), ".")
		want := SupportedVersion(parts[0] + "." + parts[1])
		if got := MajorMinor(supportedVersion); got != want {
			t.Errorf("MajorMinor(%q) = %q, want %q", supportedVersion, got, want)
		}
	}

	// Structural cases, derived from a supported version rather than
	// hand-written values: a two-segment version passes through unchanged, a
	// single-segment (major-only) version passes through unchanged, and an
	// empty version returns empty.
	supportedVersions := getSortedVersions()
	latestVersion := string(supportedVersions[len(supportedVersions)-1])
	latestVersionParts := strings.Split(latestVersion, ".")
	majorOnly := SupportedVersion(latestVersionParts[0])
	twoSegments := SupportedVersion(latestVersionParts[0] + "." + latestVersionParts[1])

	tests := []struct {
		input SupportedVersion
		want  SupportedVersion
	}{
		{twoSegments, twoSegments},
		{majorOnly, majorOnly},
		{"", ""},
	}
	for _, tt := range tests {
		if got := MajorMinor(tt.input); got != tt.want {
			t.Errorf("MajorMinor(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSupportedMajorMinorVersions(t *testing.T) {
	majorMinors := SupportedMajorMinorVersions()

	// Every supported version must appear in the lane list truncated to its
	// major.minor, so the lane list follows the constants automatically.
	for _, supportedVersion := range getSortedVersions() {
		want := MajorMinor(supportedVersion)
		found := false
		for _, majorMinorVersion := range majorMinors {
			if majorMinorVersion == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedMajorMinorVersions() is missing lane %s (required by supported version %s)", want, supportedVersion)
		}
	}

	// Ascending order via the public Compare helper (Compare requires the
	// full supported version strings, so snap each major.minor to ".0").
	for i := 1; i < len(majorMinors); i++ {
		firstVersion := majorMinors[i-1] + ".0"
		secondVersion := majorMinors[i] + ".0"
		compare, err := Compare(firstVersion, secondVersion)
		if err != nil {
			t.Fatalf("Compare(%s, %s) unexpected error: %v", firstVersion, secondVersion, err)
		}
		if compare >= 0 {
			t.Fatalf("SupportedMajorMinorVersions() is not in ascending order: %s at index %d precedes %s", majorMinors[i-1], i-1, majorMinors[i])
		}
	}
}

func TestIsValid(t *testing.T) {
	// Supported versions: every constant must validate.
	for _, supportedVersion := range getSortedVersions() {
		if !IsValid(string(supportedVersion)) {
			t.Errorf("IsValid(%q) = false for a supported version constant", supportedVersion)
		}
	}

	// Unsupported values, derived from the supported list rather than
	// hand-written versions: an unknown patch of a supported lane, a lane
	// below the oldest supported one (EOL), a lane above the newest one
	// (future), and structural non-versions.
	supportedVersions := getSortedVersions()
	oldestVersion := string(supportedVersions[0])
	oldestParts := strings.Split(oldestVersion, ".")
	newestVersion := string(supportedVersions[len(supportedVersions)-1])
	newestParts := strings.Split(newestVersion, ".")
	belowOldest := SupportedVersion(fmt.Sprintf("%s.%d.0", oldestParts[0], mustInt(t, oldestParts[1])-1))
	aboveNewest := SupportedVersion(newestParts[0] + "." + mustIntString(t, mustInt(t, newestParts[1])+1) + ".0")

	tests := []struct {
		versionString string
		want          bool
	}{
		{oldestParts[0] + "." + oldestParts[1] + ".9", false}, // unknown patch of a supported lane is invalid as-is
		{string(belowOldest), false},                          // EOL version below the oldest supported lane
		{string(aboveNewest), false},                          // future version above the newest supported lane
		{oldestParts[0] + "." + oldestParts[1], false},        // two digits are not a valid version
		{"", false},
		{"nonsense", false},
	}
	for _, tt := range tests {
		if got := IsValid(tt.versionString); got != tt.want {
			t.Errorf("IsValid(%q) = %v, want %v", tt.versionString, got, tt.want)
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1  SupportedVersion
		v2  SupportedVersion
		cmp int
	}{
		{PingFederate1220, PingFederate1220, 0},
		{PingFederate1220, PingFederate1221, -1},
		{PingFederate1228, PingFederate1230, -1},
		{PingFederate1236, PingFederate1300, -1},
		{PingFederate1311, PingFederate1310, 1},
		{PingFederate1310, PingFederate1220, 1},
	}
	for _, tt := range tests {
		compare, err := Compare(tt.v1, tt.v2)
		if err != nil {
			t.Fatalf("Compare(%s, %s) unexpected error: %v", tt.v1, tt.v2, err)
		}
		// Normalize: any negative/positive value satisfies the contract.
		got := compare
		if got > 0 {
			got = 1
		} else if got < 0 {
			got = -1
		}
		if got != tt.cmp {
			t.Errorf("Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, compare, tt.cmp)
		}
	}

	// Invalid versions error rather than compare. The invalid values derive
	// from the supported list (one lane below the oldest, one above the
	// newest) rather than hand-written versions.
	supportedVersions := getSortedVersions()
	oldestParts := strings.Split(string(supportedVersions[0]), ".")
	newestParts := strings.Split(string(supportedVersions[len(supportedVersions)-1]), ".")
	belowOldest := SupportedVersion(oldestParts[0] + "." + mustIntString(t, mustInt(t, oldestParts[1])-1) + ".0")
	aboveNewest := SupportedVersion(newestParts[0] + "." + mustIntString(t, mustInt(t, newestParts[1])+1) + ".0")
	if _, err := Compare(belowOldest, PingFederate1220); err == nil {
		t.Error("Compare with an unsupported first version should error")
	}
	if _, err := Compare(PingFederate1220, aboveNewest); err == nil {
		t.Error("Compare with an unsupported second version should error")
	}
}

func TestParse(t *testing.T) {
	// "x.y" inputs snap to the latest supported patch of that minor — derive
	// those expectations from the version list itself.
	latestPatchOf := func(lane SupportedVersion) SupportedVersion {
		var latest SupportedVersion
		for _, supportedVersion := range getSortedVersions() {
			if MajorMinor(supportedVersion) == lane {
				latest = supportedVersion
			}
		}
		return latest
	}

	// Lanes derive from the supported list rather than hand-written values.
	supportedVersions := getSortedVersions()
	oldestLaneParts := strings.Split(string(supportedVersions[0]), ".")
	newestVersion := string(supportedVersions[len(supportedVersions)-1])
	newestLaneParts := strings.Split(newestVersion, ".")
	newestLane := SupportedVersion(newestLaneParts[0] + "." + newestLaneParts[1])
	oldestLane := SupportedVersion(oldestLaneParts[0] + "." + oldestLaneParts[1])
	belowOldestLane := SupportedVersion(oldestLaneParts[0] + "." + mustIntString(t, mustInt(t, oldestLaneParts[1])-1))

	tests := []struct {
		name          string
		versionString string
		wantVersion   SupportedVersion
		wantError     bool
		wantWarning   bool
	}{
		{"oldest lane two digits snap to its latest patch", string(oldestLane), latestPatchOf(oldestLane), false, false},
		{"newest lane two digits snap to its latest patch", string(newestLane), latestPatchOf(newestLane), false, false},
		{"three digits exact", newestVersion, SupportedVersion(newestVersion), false, false},
		{"earliest version exact", string(supportedVersions[0]), supportedVersions[0], false, false},
		{"unknown patch of oldest lane warns and snaps", oldestLaneParts[0] + "." + oldestLaneParts[1] + ".99", latestPatchOf(oldestLane), false, true},
		{"unsupported lane below oldest errors", string(belowOldestLane), "", true, false},
		{"unsupported lane below oldest three digits errors", string(belowOldestLane) + ".4", "", true, false},
		{"empty errors", "", "", true, false},
		{"one segment errors", newestLaneParts[0], "", true, false},
		{"four segments error", newestVersion + ".0", "", true, false},
		{"non-numeric errors", "abc", "", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.versionString, func(t *testing.T) {
			gotVersion, diags := Parse(tt.versionString)
			if (diags.HasError()) != tt.wantError {
				t.Fatalf("Parse(%q) hasError = %v, want %v (diags: %+v)", tt.versionString, diags.HasError(), tt.wantError, diags.Errors())
			}
			if tt.wantError {
				return
			}
			if gotVersion != tt.wantVersion {
				t.Fatalf("Parse(%q) = %q, want %q", tt.versionString, gotVersion, tt.wantVersion)
			}
			if (diags.WarningsCount() > 0) != tt.wantWarning {
				t.Fatalf("Parse(%q) warning = %v, want %v", tt.versionString, diags.WarningsCount() > 0, tt.wantWarning)
			}
		})
	}
}

func TestAddUnsupportedAttributeError(t *testing.T) {
	diags := diag.Diagnostics{}
	AddUnsupportedAttributeError("client_cert_header_encoding_format", PingFederate1220, PingFederate1230, &diags)
	if !diags.HasError() {
		t.Fatal("AddUnsupportedAttributeError should add an error diagnostic")
	}
	if len(diags.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
	}
	// The detail embeds the versions; assert via the constants so the message
	// follows internal/version values.
	if !strings.Contains(diags.Errors()[0].Detail(),
		"PingFederate version "+string(PingFederate1230)+" or later is required for attribute client_cert_header_encoding_format") {
		t.Errorf("error detail should name the attribute and required version %s, got: %s", PingFederate1230, diags.Errors()[0].Detail())
	}
}

func TestAddUnsupportedResourceError(t *testing.T) {
	diags := diag.Diagnostics{}
	AddUnsupportedResourceError("pingfederate_foo", PingFederate1220, PingFederate1230, &diags)
	if !diags.HasError() {
		t.Fatal("AddUnsupportedResourceError should add an error diagnostic")
	}
	if len(diags.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
	}
	if !strings.Contains(diags.Errors()[0].Detail(),
		"PingFederate version "+string(PingFederate1230)+" or later is required for resource pingfederate_foo") {
		t.Errorf("error detail should name the resource and required version %s, got: %s", PingFederate1230, diags.Errors()[0].Detail())
	}
}

// mustInt parses a non-negative decimal string, failing the test otherwise.
func mustInt(t *testing.T, digits string) int {
	t.Helper()
	value := 0
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			t.Fatalf("expected digits, got %q", digits)
		}
		value = value*10 + int(digit-'0')
	}
	return value
}

// mustIntString renders an int back to decimal digits.
func mustIntString(t *testing.T, value int) string {
	t.Helper()
	if value < 0 {
		t.Fatalf("cannot render negative lane number %d", value)
	}
	return fmt.Sprintf("%d", value)
}
