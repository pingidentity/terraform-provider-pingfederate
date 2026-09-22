// Copyright © 2026 Ping Identity Corporation

// Package upgradeladder implements cross-version provider upgrade tests: for
// each registered resource, apply a known-good configuration with the oldest
// provider version that supports the current PingFederate server lane, then
// hop the provider one minor version at a time to the local build, asserting
// apply success and an empty post-apply plan after every hop.
package upgradeladder

import (
	"fmt"
	"os"
	"strings"

	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// LocalRung marks the final, in-process rung (the local build of this provider
// repository).
const LocalRung = "local"

// EnvLadderOverride is the environment variable that overrides the ladder
// for a test run: "full" (default), "last2" (final registry rung + local),
// or an explicit comma-separated list of rungs, e.g. "1.9.0,1.10.0,local".
const EnvLadderOverride = "PINGFEDERATE_UPGRADE_LADDER"

// ladderEntry is one row of the version ladder: a released provider version
// and the highest PingFederate major.minor version that provider release can
// configure (its MaxPFMinor, stated as an internal/version constant).
//
// Rows cover every provider minor from v1.3.0 (the first release supporting a
// PingFederate version the current provider still supports) through the latest
// release. Rungs older than v1.3.0 (v1.0.0-v1.2.0) cannot configure any
// supported PingFederate version and are intentionally absent.
//
// Maintenance when a PingFederate version gains support in
// internal/version/version.go: verify the newly supported lane against the
// published provider releases (git show v<tag>:internal/version/version.go)
// and append the new release rows here. When internal/version drops an EOL
// lane, its rows fail their MaxPFMinor compile-time reference until removed,
// keeping the table honest. Proven via
// `git show v<tag>:internal/version/version.go` for every tag.
type ladderEntry struct {
	ProviderVersion string                   // registry version, latest patch of the minor, e.g. "1.4.5"
	MaxPFMinor      version.SupportedVersion // highest PingFederate major.minor configurable, e.g. version.PingFederate1220
}

// MaxPFMinor values are tied to internal/version constants so an EOL lane
// drop in internal/version breaks compilation here until the stale rows are
// removed. When a new lane is added to internal/version,
// TestLadderTableCoversSupportedLanes fails until its release rows are
// appended here — the manual step, made loud.
var ladderTable = []ladderEntry{
	{"1.3.0", version.PingFederate1220}, {"1.4.5", version.PingFederate1220}, {"1.5.0", version.PingFederate1220},
	{"1.6.2", version.PingFederate1230},
	{"1.7.1", version.PingFederate1300}, {"1.8.1", version.PingFederate1300},
	{"1.9.0", version.PingFederate1310}, {"1.10.0", version.PingFederate1310},
}

// laneFloorVersion returns the first rung that supports a PingFederate
// major.minor lane (the earliest table entry whose MaxPFMinor matches).
func laneFloorVersion(lane string) (string, bool) {
	for _, entry := range ladderTable {
		if string(version.MajorMinor(entry.MaxPFMinor)) == lane {
			return entry.ProviderVersion, true
		}
	}
	return "", false
}

// supportedLanes derives the PingFederate major.minor lanes from the
// internal/version supported list, so new versions there automatically extend
// the ladder's lane handling.
func supportedLanes() []string {
	majorMinors := version.SupportedMajorMinorVersions()
	lanes := make([]string, 0, len(majorMinors))
	for _, majorMinorVersion := range majorMinors {
		lanes = append(lanes, string(majorMinorVersion))
	}
	return lanes
}

// ParseLane maps a PINGFEDERATE_PROVIDER_PRODUCT_VERSION value ("12.2",
// "12.2.8") to a supported lane key ("12.2", "12.3", "13.0", "13.1") via the
// shared version.Parse, so the ladder validates product versions exactly like
// the provider itself (ConfigurationPreCheck) does, against the same
// supported-version list in internal/version/version.go.
func ParseLane(productVersion string) (string, error) {
	parsedVersion, diags := version.Parse(productVersion)
	if diags.HasError() {
		var details []string
		for _, errDiag := range diags.Errors() {
			details = append(details, errDiag.Detail())
		}
		return "", fmt.Errorf("cannot determine PingFederate lane from product version '%s': %s",
			productVersion, strings.Join(details, "; "))
	}
	return string(version.MajorMinor(parsedVersion)), nil
}

// BuildLadder returns the provider versions for a lane, oldest first, always
// ending with LocalRung: every rung that supports the lane (its MaxPFMinor is
// at or above the lane), starting at the lane floor. For example, the 12.2
// lane yields
// ["1.3.0","1.4.5","1.5.0","1.6.2","1.7.1","1.8.1","1.9.0","1.10.0",LocalRung].
func BuildLadder(lane string) ([]string, error) {
	lanes := supportedLanes()
	supported := false
	for _, supportedLane := range lanes {
		if lane == supportedLane {
			supported = true
			break
		}
	}
	if !supported {
		return nil, fmt.Errorf("no ladder for PingFederate lane '%s'; supported lanes are: %s", lane, strings.Join(lanes, ", "))
	}

	rungs := make([]string, 0, len(ladderTable)+1)
	seenFloor := false
	for _, entry := range ladderTable {
		if !seenFloor {
			if string(version.MajorMinor(entry.MaxPFMinor)) != lane {
				continue
			}
			seenFloor = true
		}
		rungs = append(rungs, entry.ProviderVersion)
	}
	if !seenFloor {
		return nil, fmt.Errorf("no ladder table entry supports PingFederate lane '%s'", lane)
	}
	rungs = append(rungs, LocalRung)
	return rungs, nil
}

// BuildLadderFromEnv resolves the ladder for a test run: the
// PINGFEDERATE_UPGRADE_LADDER environment variable if set ("full", "last2",
// or explicit rungs), else the full lane ladder.
func BuildLadderFromEnv(lane string) ([]string, error) {
	override := strings.TrimSpace(os.Getenv(EnvLadderOverride))
	switch override {
	case "", "full":
		return BuildLadder(lane)
	case "last2":
		ladder, err := BuildLadder(lane)
		if err != nil {
			return nil, err
		}
		// Keep the final registry rung plus LocalRung.
		return ladder[len(ladder)-2:], nil
	default:
		rungs := make([]string, 0, len(override))
		for _, part := range strings.Split(override, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				return nil, fmt.Errorf("empty rung in %s override '%s'", EnvLadderOverride, override)
			}
			if trimmed == LocalRung {
				rungs = append(rungs, LocalRung)
				continue
			}
			if !isKnownRung(trimmed) {
				return nil, fmt.Errorf("rung '%s' in %s override is not a ladder version (table versions: %s)",
					trimmed, EnvLadderOverride, ladderTableVersions())
			}
			rungs = append(rungs, trimmed)
		}
		if len(rungs) == 0 || rungs[len(rungs)-1] != LocalRung {
			// Always end on the local build.
			rungs = append(rungs, LocalRung)
		}
		return rungs, nil
	}
}

func isKnownRung(version string) bool {
	for _, entry := range ladderTable {
		if entry.ProviderVersion == version {
			return true
		}
	}
	return false
}

func ladderTableVersions() string {
	versions := make([]string, 0, len(ladderTable)+1)
	for _, entry := range ladderTable {
		versions = append(versions, entry.ProviderVersion)
	}
	versions = append(versions, LocalRung)
	return strings.Join(versions, ", ")
}
