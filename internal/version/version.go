package version

import (
	"errors"
	"strings"
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
	PingFederate1130 SupportedVersion = "11.3.0"
	PingFederate1131 SupportedVersion = "11.3.1"
	PingFederate1132 SupportedVersion = "11.3.2"
	PingFederate1133 SupportedVersion = "11.3.3"
	PingFederate1134 SupportedVersion = "11.3.4"
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
		PingFederate1130,
		PingFederate1131,
		PingFederate1132,
		PingFederate1133,
		PingFederate1134,
	}
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

func Parse(versionString string) (SupportedVersion, error) {
	if len(versionString) == 0 {
		return "", errors.New("failed to parse PingFederate version: empty version string")
	}

	var err error
	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x"
	// If only two digits are supplied, the last one will be assumed to be "0"
	if len(versionDigits) != 2 && len(versionDigits) != 3 {
		return "", errors.New("failed to parse PingFederate version '" + versionString + "', Expected either two digits (e.g. '11.3') or three digits (e.g. '11.3.4')")
	}
	if len(versionDigits) == 2 {
		versionString += ".0"
	}
	if !IsValid(versionString) {
		err = errors.New("unsupported PingFederate version: " + versionString)
	}
	return SupportedVersion(versionString), err
}
