package version

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Supported PingFederate versions
const (
	PingFederate1125 = "11.2.5"
	PingFederate1130 = "11.3.0"
	PingFederate1131 = "11.3.1"
	PingFederate1132 = "11.3.2"
	PingFederate1133 = "11.3.3"
	PingFederate1134 = "11.3.4"
)

func IsValid(versionString string) bool {
	return getSortedVersionIndex(versionString) != -1
}

func getSortedVersionIndex(versionString string) int {
	for i, version := range getSortedVersions() {
		if version == versionString {
			return i
		}
	}
	return -1
}

func getSortedVersions() []string {
	return []string{
		PingFederate1125,
		PingFederate1130,
		PingFederate1131,
		PingFederate1132,
		PingFederate1133,
		PingFederate1134,
	}
}

// Compare two PingFederate versions. Returns a negative number if the first argument is less than the second,
// zero if they are equal, and a positive number if the first argument is greater than the second
func Compare(version1, version2 string) (int, error) {
	version1Index := getSortedVersionIndex(version1)
	if version1Index == -1 {
		return 0, errors.New("Invalid version: " + version1)
	}
	version2Index := getSortedVersionIndex(version2)
	if version2Index == -1 {
		return 0, errors.New("Invalid version: " + version2)
	}

	return version1Index - version2Index, nil
}

func Parse(versionString string) (string, error) {
	if len(versionString) == 0 {
		return versionString, errors.New("failed to parse PingFederate version: empty version string")
	}

	var err error
	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x"
	// If only two digits are supplied, the last one will be assumed to be "0"
	if len(versionDigits) != 2 && len(versionDigits) != 3 {
		return versionString, errors.New("failed to parse PingFederate version '" + versionString + "', Expected either two digits (e.g. '11.3') or three digits (e.g. '11.3.4')")
	}
	if len(versionDigits) == 2 {
		versionString += ".0"
	}
	if !IsValid(versionString) {
		err = errors.New("unsupported PingFederate version: " + versionString)
	}
	return versionString, err
}

func CheckResourceSupported(diagnostics *diag.Diagnostics, minimumVersion, actualVersion, resourceName string) {
	// Check that the version is at least the minimum version
	compare, err := Compare(actualVersion, minimumVersion)
	if err != nil {
		diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	if compare < 0 {
		diagnostics.AddError(resourceName+" is only supported for PingFederate versions "+minimumVersion+" and later", "Found PF version "+actualVersion)
		return
	}
}
