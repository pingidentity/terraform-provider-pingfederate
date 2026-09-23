// Copyright © 2026 Ping Identity Corporation

package upgradeladder

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestServerLanesDefault(t *testing.T) {
	t.Setenv(EnvServerLanesOverride, "")
	got, err := serverLanes()
	if err != nil {
		t.Fatalf("serverLanes() unexpected error: %v", err)
	}
	want := supportedLaneStrings()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("serverLanes() = %v, want %v", got, want)
	}
	// Lanes must be ascending so the ladder steps forward (two-segment lane
	// keys like "12.3" are not version.Parse-able; pad to three segments).
	for i := 1; i < len(got); i++ {
		prev, cur := got[i-1]+".0", got[i]+".0"
		if _, err := version.Compare(version.SupportedVersion(cur), version.SupportedVersion(prev)); err != nil {
			t.Fatalf("lane pair invalid: %v", err)
		}
	}
}

// supportedLaneStrings renders the supported major.minor list as strings.
func supportedLaneStrings() []string {
	majorMinors := version.SupportedMajorMinorVersions()
	lanes := make([]string, 0, len(majorMinors))
	for _, lane := range majorMinors {
		lanes = append(lanes, string(lane))
	}
	return lanes
}

func TestServerLanesOverride(t *testing.T) {
	// Every lane value derives from internal/version: the oldest and newest
	// supported lanes, plus a version one major below the oldest (unsupported).
	supported := version.SupportedMajorMinorVersions()
	oldest := string(supported[0])
	oldestParts := strings.Split(oldest, ".")
	belowOldest := fmt.Sprintf("%s.%s", mustDecrement(oldestParts[0]), oldestParts[1])
	newest := string(supported[len(supported)-1])
	oldestPatch := string(version.SupportedVersions()[0])

	tests := []struct {
		name    string
		env     string
		want    []string
		wantErr string
	}{
		{"full versions truncate to lanes", oldestPatch + "," + newest + ".1", []string{oldest, newest}, ""},
		{"lanes passthrough", oldest + "," + newest, []string{oldest, newest}, ""},
		{"empty part rejected", oldest + ",," + newest, nil, "empty lane"},
		{"unsupported version rejected", belowOldest + "," + newest, nil, "not a supported PingFederate version"},
		{"garbage rejected", "abc", nil, "not a supported PingFederate version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvServerLanesOverride, tt.env)
			got, err := serverLanes()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("serverLanes() error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("serverLanes() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("serverLanes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustDecrement decrements a decimal major string, so the "unsupported"
// fixture stays relative to the supported list as it evolves.
func mustDecrement(major string) string {
	n := 0
	for _, digit := range major {
		if digit < '0' || digit > '9' {
			panic(fmt.Sprintf("expected digits, got %q", major))
		}
		n = n*10 + int(digit-'0')
	}
	return fmt.Sprintf("%d", n-1)
}

func TestLaneHosts(t *testing.T) {
	t.Run("single lane uses base host", func(t *testing.T) {
		t.Setenv(httpsHostEnvVar, "https://localhost:9999")
		t.Setenv(EnvServerLaneHosts, "")
		got, err := laneHosts(derivedLanes(1))
		if err != nil {
			t.Fatalf("laneHosts() unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []string{"https://localhost:9999"}) {
			t.Fatalf("laneHosts() = %v", got)
		}
	})

	t.Run("one host per extra lane required", func(t *testing.T) {
		t.Setenv(httpsHostEnvVar, "https://localhost:9999")
		t.Setenv(EnvServerLaneHosts, "https://localhost:10099,https://localhost:10199")
		got, err := laneHosts(derivedLanes(3))
		if err != nil {
			t.Fatalf("laneHosts() unexpected error: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("laneHosts() = %v, want 3 hosts", got)
		}
	})

	t.Run("missing hosts rejected", func(t *testing.T) {
		t.Setenv(httpsHostEnvVar, "https://localhost:9999")
		t.Setenv(EnvServerLaneHosts, "")
		if _, err := laneHosts(derivedLanes(2)); err == nil || !strings.Contains(err.Error(), "one host per lane") {
			t.Fatalf("laneHosts() error = %v, want host-count mismatch", err)
		}
	})

	t.Run("empty part in hosts rejected", func(t *testing.T) {
		t.Setenv(httpsHostEnvVar, "https://localhost:9999")
		// The empty host part is rejected before the count check, so any lane
		// count exercises it; derivedLanes(2) with a 2-entry host list keeps
		// the fixture minimal.
		t.Setenv(EnvServerLaneHosts, "https://localhost:10099,")
		if _, err := laneHosts(derivedLanes(3)); err == nil {
			t.Fatal("laneHosts() should reject empty host parts")
		}
	})

	t.Run("missing base host rejected", func(t *testing.T) {
		t.Setenv(httpsHostEnvVar, "")
		t.Setenv(EnvServerLaneHosts, "")
		if _, err := laneHosts(derivedLanes(2)); err == nil || !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("laneHosts() error = %v, want base-host requirement", err)
		}
	})
}

// derivedLanes returns the first n supported lanes, so tests carry no
// hard-coded PingFederate versions and follow the supported list automatically.
func derivedLanes(n int) []string {
	supported := supportedLaneStrings()
	if len(supported) < n {
		panic(fmt.Sprintf("test needs %d supported lanes, internal/version has %d", n, len(supported)))
	}
	return supported[:n]
}

func TestLaneProviderBlock(t *testing.T) {
	host, lane := "https://localhost:10099", string(version.SupportedMajorMinorVersions()[0])
	block := laneProviderBlock(host, lane)
	if !strings.Contains(block, `provider "pingfederate"`) {
		t.Fatalf("laneProviderBlock() missing provider block:\n%s", block)
	}
	if !strings.Contains(block, fmt.Sprintf("%q", host)) || !strings.Contains(block, fmt.Sprintf("%q", lane)) {
		t.Fatalf("laneProviderBlock() must quote host and version:\n%s", block)
	}
}
