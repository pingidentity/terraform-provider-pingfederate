package planmodifiers

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Object = &x509FileData{}

type x509FileData struct{}

func (v x509FileData) Description(ctx context.Context) string {
	return "Validates that the formatted_file_data and file_data values match. This is to detect a drift with assigned certificate(s)."
}

func (v x509FileData) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v x509FileData) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}

	var planFileDataStringFormatted, formattedFileDataAsStringFormatted, fileDataBase64Decoded string

	// Remove header, footer, and new lines
	stringReplacer := strings.NewReplacer("-----BEGIN CERTIFICATE-----", "", "-----END CERTIFICATE-----", "", "\n", "")

	// Get file_data from plan
	fileData, ok := req.ConfigValue.Attributes()["file_data"]
	if !ok {
		return
	}
	planFileDataString := fileData.(types.String).ValueString()
	planFileDataStringFormatted = stringReplacer.Replace(planFileDataString)
	base64DecodedFileData, err := base64.StdEncoding.DecodeString(planFileDataString)
	if err == nil {
		// The plan value was base64-encoded, use the decoded value for comparison
		fileDataBase64Decoded = string(base64DecodedFileData)
	}
	fileDataBase64Decoded = stringReplacer.Replace(fileDataBase64Decoded)

	// Get formatted_file_data from state
	formattedFileData, ok := req.StateValue.Attributes()["formatted_file_data"].(types.String)
	if !ok {
		return
	}
	formattedFileDataAsStringFormatted = stringReplacer.Replace(formattedFileData.ValueString())

	// Check if formatted_file_data and file_data strings match, or if formatted_file_data matches original string
	// If they do not, formatted_file_data is set to unknown
	if formattedFileDataAsStringFormatted != fileDataBase64Decoded && formattedFileDataAsStringFormatted != planFileDataStringFormatted {
		var respDiags diag.Diagnostics
		reqConfigAttrs := req.ConfigValue.Attributes()
		reqConfigAttrs["formatted_file_data"] = types.StringUnknown()
		resp.PlanValue, respDiags = types.ObjectValue(req.ConfigValue.AttributeTypes(ctx), reqConfigAttrs)
		if respDiags != nil {
			resp.Diagnostics.AddError(
				"Unable to build plan object",
				"x509_file object did not build properly",
			)
		}
	}
}

func ValidateX509FileData() x509FileData {
	return x509FileData{}
}
