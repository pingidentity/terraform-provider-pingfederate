// Copyright © 2026 Ping Identity Corporation

package pemcertificates

import (
	"encoding/base64"
	"strings"
)

// Compare PEM encoded certificates, handling any formatting done by the PF API
func FileDataEquivalent(planned, apiResponse string) bool {
	var plannedFormatted, apiResponseFormatted, plannedBase64Decoded string

	// Remove header, footer, and new lines
	stringReplacer := strings.NewReplacer("-----BEGIN CERTIFICATE-----", "", "-----END CERTIFICATE-----", "", "\n", "")

	plannedFormatted = stringReplacer.Replace(planned)
	base64DecodedPlannedBytes, err := base64.StdEncoding.DecodeString(planned)
	if err == nil {
		// The plan value was base64-encoded, use the decoded value for comparison
		plannedBase64Decoded = string(base64DecodedPlannedBytes)
	}
	plannedBase64Decoded = stringReplacer.Replace(plannedBase64Decoded)

	apiResponseFormatted = stringReplacer.Replace(apiResponse)

	// If the formatted versions match (base64 decoded or not), these represent the same certificate
	return plannedFormatted == apiResponseFormatted || plannedBase64Decoded == apiResponseFormatted
}
