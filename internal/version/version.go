package version

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type SupportedVersion string

// Supported PingFederate versions
const (
	PingFederate1120 SupportedVersion = "11.2.0"
	PingFederate1121 SupportedVersion = "11.2.1"
	PingFederate1122 SupportedVersion = "11.2.2"
	PingFederate1123 SupportedVersion = "11.2.3"
	PingFederate1124 SupportedVersion = "11.2.4"
	PingFederate1125 SupportedVersion = "11.2.5"
	PingFederate1126 SupportedVersion = "11.2.6"
	PingFederate1127 SupportedVersion = "11.2.7"
	PingFederate1128 SupportedVersion = "11.2.8"
	PingFederate1130 SupportedVersion = "11.3.0"
	PingFederate1131 SupportedVersion = "11.3.1"
	PingFederate1132 SupportedVersion = "11.3.2"
	PingFederate1133 SupportedVersion = "11.3.3"
	PingFederate1134 SupportedVersion = "11.3.4"
	PingFederate1135 SupportedVersion = "11.3.5"
	PingFederate1200 SupportedVersion = "12.0.0"
	PingFederate1201 SupportedVersion = "12.0.1"
)

func IsValid(versionString string) bool {
	return getSortedVersionIndex(SupportedVersion(versionString)) != -1
}

func getSortedVersionIndex(versionString SupportedVersion) int {
	for i, version := range getSortedVersions() {
		if version == versionString {
			return i
		}
	}
	return -1
}

func getSortedVersions() []SupportedVersion {
	return []SupportedVersion{
		PingFederate1120,
		PingFederate1121,
		PingFederate1122,
		PingFederate1123,
		PingFederate1124,
		PingFederate1125,
		PingFederate1126,
		PingFederate1127,
		PingFederate1128,
		PingFederate1130,
		PingFederate1131,
		PingFederate1132,
		PingFederate1133,
		PingFederate1134,
		PingFederate1135,
		PingFederate1200,
		PingFederate1201,
	}
}

func getSortedVersionsMessage() string {
	message := "Supported versions are: "
	for i, version := range getSortedVersions() {
		message += string(version)
		if i < len(getSortedVersions())-1 {
			message += ", "
		}
	}
	return message
}

// Compare two PingFederate versions. Returns a negative number if the first argument is less than the second,
// zero if they are equal, and a positive number if the first argument is greater than the second
func Compare(version1, version2 SupportedVersion) (int, error) {
	version1Index := getSortedVersionIndex(version1)
	if version1Index == -1 {
		return 0, errors.New("Invalid version: " + string(version1))
	}
	version2Index := getSortedVersionIndex(version2)
	if version2Index == -1 {
		return 0, errors.New("Invalid version: " + string(version2))
	}

	return version1Index - version2Index, nil
}

func Parse(versionString string) (SupportedVersion, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(versionString) == 0 {
		diags.AddError("failed to parse PingFederate version", "empty version string")
		return "", diags
	}

	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x"
	// If only two digits are supplied, the last one will be assumed to be "0"
	if len(versionDigits) != 2 && len(versionDigits) != 3 {
		diags.AddError("failed to parse PingFederate version '"+versionString+"'", "Expected either two digits (e.g. '11.3') or three digits (e.g. '11.3.4')")
		return "", diags
	}
	if len(versionDigits) == 2 {
		versionString += ".0"
	}
	if !IsValid(versionString) {
		// Check if the major-minor version is valid
		majorMinorVersionString := versionDigits[0] + "." + versionDigits[1] + ".0"
		if !IsValid(majorMinorVersionString) {
			diags.AddError("unsupported PingFederate version '"+versionString+"'", getSortedVersionsMessage())
			return "", diags
		}
		// The major-minor version is valid, only the patch is invalid. Warn but do not fail, assume the lastest patch version
		sortedVersions := getSortedVersions()
		versionIndex := -1
		switch majorMinorVersionString {
		case "11.2.0":
			// Use the first version prior to 11.3.0
			versionIndex = getSortedVersionIndex(PingFederate1130) - 1
		case "11.3.0":
			// Use the first version prior to 12.0.0
			versionIndex = getSortedVersionIndex(PingFederate1200) - 1
		case "12.0.0":
			// This is the latest major-minor version, so just use the latest patch version available
			versionIndex = len(sortedVersions) - 1
		}
		if versionIndex < 0 || versionIndex >= len(sortedVersions) {
			// This should never happen
			diags.AddError("Unexpected failure determining major-minor PingFederate version", "")
			return "", diags
		}
		assumedVersion := string(sortedVersions[versionIndex])
		diags.AddWarning("Unrecognized PingFederate version '"+versionString+"'", "Assuming the latest patch version available: '"+assumedVersion+"'")
		versionString = assumedVersion
	}
	return SupportedVersion(versionString), diags
}

func AddUnsupportedAttributeError(attr string, actualVersion, requiredVersion SupportedVersion, diags *diag.Diagnostics) {
	if diags == nil {
		return
	}

	diags.AddError(fmt.Sprintf("Attribute '%s' not supported by PingFederate version %s", attr, string(actualVersion)),
		fmt.Sprintf("PingFederate version %s or later is required for this attribute", string(requiredVersion)))
}
