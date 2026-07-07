// Copyright © 2026 Ping Identity Corporation

package version

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

type SupportedVersion string

// Supported PingFederate versions
const (
	PingFederate1220 SupportedVersion = "12.2.0"
	PingFederate1221 SupportedVersion = "12.2.1"
	PingFederate1222 SupportedVersion = "12.2.2"
	PingFederate1223 SupportedVersion = "12.2.3"
	PingFederate1224 SupportedVersion = "12.2.4"
	PingFederate1225 SupportedVersion = "12.2.5"
	PingFederate1226 SupportedVersion = "12.2.6"
	PingFederate1227 SupportedVersion = "12.2.7"
	PingFederate1228 SupportedVersion = "12.2.8"
	PingFederate1230 SupportedVersion = "12.3.0"
	PingFederate1231 SupportedVersion = "12.3.1"
	PingFederate1232 SupportedVersion = "12.3.2"
	PingFederate1233 SupportedVersion = "12.3.3"
	PingFederate1234 SupportedVersion = "12.3.4"
	PingFederate1235 SupportedVersion = "12.3.5"
	PingFederate1236 SupportedVersion = "12.3.6"
	PingFederate1300 SupportedVersion = "13.0.0"
	PingFederate1301 SupportedVersion = "13.0.1"
	PingFederate1302 SupportedVersion = "13.0.2"
	PingFederate1303 SupportedVersion = "13.0.3"
	PingFederate1310 SupportedVersion = "13.1.0"
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
		PingFederate1220,
		PingFederate1221,
		PingFederate1222,
		PingFederate1223,
		PingFederate1224,
		PingFederate1225,
		PingFederate1226,
		PingFederate1227,
		PingFederate1228,
		PingFederate1230,
		PingFederate1231,
		PingFederate1232,
		PingFederate1233,
		PingFederate1234,
		PingFederate1235,
		PingFederate1236,
		PingFederate1300,
		PingFederate1301,
		PingFederate1302,
		PingFederate1303,
		PingFederate1310,
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

func getLatestPatchForMajorMinorVersion(majorMinorVersionString string) (string, diag.Diagnostics) {
	var respDiags diag.Diagnostics
	sortedVersions := getSortedVersions()
	versionIndex := -1
	switch majorMinorVersionString {
	case "12.2.0":
		// Use the first version prior to 12.3.0
		versionIndex = getSortedVersionIndex(PingFederate1230) - 1
	case "12.3.0":
		// Use the first version prior to 13.0.0
		versionIndex = getSortedVersionIndex(PingFederate1300) - 1
	case "13.0.0":
		// Use the first version prior to 13.1.0
		versionIndex = getSortedVersionIndex(PingFederate1310) - 1
	case "13.1.0":
		// This is the latest major-minor version, so just use the latest patch version available
		versionIndex = len(sortedVersions) - 1
	}
	if versionIndex < 0 || versionIndex >= len(sortedVersions) {
		// This should never happen
		respDiags.AddError(providererror.InternalProviderError, "Unexpected failure determining major-minor PingFederate version")
		return majorMinorVersionString, respDiags
	}
	return string(sortedVersions[versionIndex]), respDiags
}

func Parse(versionString string) (SupportedVersion, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(versionString) == 0 {
		diags.AddAttributeError(
			path.Root("product_version"),
			providererror.InvalidProviderConfiguration, "failed to parse PingFederate version: empty version string")
		return "", diags
	}

	versionDigits := strings.Split(versionString, ".")
	// Expect a version like "x.x" or "x.x.x"
	// If only two digits are supplied, the last one will be assumed to be "0"
	if len(versionDigits) != 2 && len(versionDigits) != 3 {
		diags.AddAttributeError(
			path.Root("product_version"),
			providererror.InvalidProviderConfiguration, "failed to parse PingFederate version '"+versionString+"'. Expected either two digits (e.g. '11.3') or three digits (e.g. '11.3.4')")
		return "", diags
	}
	if len(versionDigits) == 2 {
		if !IsValid(versionString + ".0") {
			// This major minor version isn't supported - fail now
			diags.AddAttributeError(
				path.Root("product_version"),
				providererror.InvalidProviderConfiguration,
				"PingFederate version '"+versionString+"' is not supported in this version of the PingFederate terraform provider.\n"+getSortedVersionsMessage())
			return "", diags
		}
		// Get the latest patch for the major minor version provided
		var respDiags diag.Diagnostics
		versionString, respDiags = getLatestPatchForMajorMinorVersion(versionString + ".0")
		diags.Append(respDiags...)
	}
	if !IsValid(versionString) {
		// Check if the major-minor version is valid
		majorMinorVersionString := versionDigits[0] + "." + versionDigits[1] + ".0"
		if !IsValid(majorMinorVersionString) {
			diags.AddAttributeError(
				path.Root("product_version"),
				providererror.InvalidProviderConfiguration,
				"PingFederate version '"+versionString+"' is not supported in this version of the PingFederate terraform provider.\n"+getSortedVersionsMessage())
			return "", diags
		}
		// The major-minor version is valid, only the patch is invalid. Warn but do not fail, assume the lastest patch version
		var respDiags diag.Diagnostics
		originalVersionString := versionString
		versionString, respDiags = getLatestPatchForMajorMinorVersion(majorMinorVersionString)
		diags.Append(respDiags...)
		diags.AddAttributeWarning(
			path.Root("product_version"),
			"Unrecognized PingFederate patch version in 'product_version' field or 'PINGFEDERATE_PROVIDER_PRODUCT_VERSION' environment variable",
			"PingFederate patch version '"+originalVersionString+"' is not recognized by this version of the PingFederate terraform provider. Assuming the latest patch version supported by the provider: '"+versionString+"'")
	}
	return SupportedVersion(versionString), diags
}

func AddUnsupportedAttributeError(attr string, actualVersion, requiredVersion SupportedVersion, diags *diag.Diagnostics) {
	if diags == nil {
		return
	}

	diags.AddAttributeError(
		path.Root(attr),
		providererror.InvalidProductVersionAttribute,
		fmt.Sprintf("PingFederate version %s or later is required for attribute %s. "+
			"PingFederate version %s was provided via the 'product_version' field in your provider configuration or the 'PINGFEDERATE_PROVIDER_PRODUCT_VERSION' environment variable.", string(requiredVersion), attr, string(actualVersion)))
}

func AddUnsupportedResourceError(resource string, actualVersion, requiredVersion SupportedVersion, diags *diag.Diagnostics) {
	if diags == nil {
		return
	}

	diags.AddError(
		providererror.InvalidProductVersionResource,
		fmt.Sprintf("PingFederate version %s or later is required for resource %s. "+
			"PingFederate version %s was provided via the 'product_version' field in your provider configuration or the 'PINGFEDERATE_PROVIDER_PRODUCT_VERSION' environment variable.", string(requiredVersion), resource, string(actualVersion)))
}
