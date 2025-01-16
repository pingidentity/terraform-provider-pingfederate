package planmodifiers

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.List = &x509FileData{}

type x509FileData struct{}

func (v x509FileData) Description(ctx context.Context) string {
	return "Validates that the formatted_file_data and file_data values match. This is to detect a drift with assigned certificate(s)."
}

func (v x509FileData) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v x509FileData) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}

	if len(req.PlanValue.Elements()) == 0 || len(req.StateValue.Elements()) == 0 {
		return
	}

	var finalPlanElements []attr.Value
	var respDiags diag.Diagnostics
	for i, planElement := range req.PlanValue.Elements() {
		if i >= len(req.StateValue.Elements()) {
			continue
		}
		stateValueAttrs := req.StateValue.Elements()[i].(types.Object).Attributes()
		planValue := planElement.(types.Object)
		planValueAttrs := planValue.Attributes()
		stateX509Value := stateValueAttrs["x509_file"].(types.Object)
		planX509Value := planValueAttrs["x509_file"].(types.Object)
		stateCertViewValue := stateValueAttrs["cert_view"].(types.Object)
		planCertViewValue := planValueAttrs["cert_view"].(types.Object)

		var planFileDataStringFormatted, formattedFileDataAsStringFormatted, fileDataBase64Decoded string

		// Remove header, footer, and new lines
		stringReplacer := strings.NewReplacer("-----BEGIN CERTIFICATE-----", "", "-----END CERTIFICATE-----", "", "\n", "")

		// Get file_data from plan
		fileData, ok := planX509Value.Attributes()["file_data"]
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
		formattedFileData, ok := stateX509Value.Attributes()["formatted_file_data"].(types.String)
		if !ok {
			return
		}
		formattedFileDataAsStringFormatted = stringReplacer.Replace(formattedFileData.ValueString())

		// Check if formatted_file_data and file_data strings match, or if formatted_file_data matches original string
		// If they do not, formatted_file_data is set to unknown
		if formattedFileDataAsStringFormatted != fileDataBase64Decoded && formattedFileDataAsStringFormatted != planFileDataStringFormatted {
			reqPlanAttrs := planX509Value.Attributes()
			reqPlanAttrs["formatted_file_data"] = types.StringUnknown()
			planX509Value, respDiags = types.ObjectValue(planX509Value.AttributeTypes(ctx), reqPlanAttrs)
			if respDiags.HasError() {
				resp.Diagnostics.AddError(
					"Unable to build plan object",
					"x509_file object did not build properly",
				)
			}
			resp.Diagnostics.Append(respDiags...)
			planCertViewValue = types.ObjectUnknown(planCertViewValue.AttributeTypes(ctx))
		} else {
			reqPlanAttrs := planX509Value.Attributes()
			reqPlanAttrs["formatted_file_data"] = types.StringValue(formattedFileData.ValueString())
			planX509Value, respDiags = types.ObjectValue(planX509Value.AttributeTypes(ctx), reqPlanAttrs)
			resp.Diagnostics.Append(respDiags...)
			if respDiags.HasError() {
				resp.Diagnostics.AddError(
					"Unable to build plan object",
					"x509_file object did not build properly",
				)
			}
			planCertViewValue, respDiags = types.ObjectValue(planCertViewValue.AttributeTypes(ctx), stateCertViewValue.Attributes())
			resp.Diagnostics.Append(respDiags...)
			if respDiags.HasError() {
				resp.Diagnostics.AddError(
					"Unable to build plan object",
					"cert_view object did not build properly",
				)
			}
		}

		if resp.Diagnostics.HasError() {
			return
		}

		// Build the final element
		planValueAttrs["x509_file"] = planX509Value
		planValueAttrs["cert_view"] = planCertViewValue
		if planValueAttrs["active_verification_cert"].IsUnknown() {
			planValueAttrs["active_verification_cert"] = types.BoolValue(false)
		}
		if planValueAttrs["primary_verification_cert"].IsUnknown() {
			planValueAttrs["primary_verification_cert"] = types.BoolValue(false)
		}
		if planValueAttrs["secondary_verification_cert"].IsUnknown() {
			planValueAttrs["secondary_verification_cert"] = types.BoolValue(false)
		}
		if planValueAttrs["encryption_cert"].IsUnknown() {
			planValueAttrs["encryption_cert"] = types.BoolValue(false)
		}

		finalPlanValue, respDiags := types.ObjectValue(planValue.AttributeTypes(ctx), planValueAttrs)
		if respDiags.HasError() {
			resp.Diagnostics.AddError(
				"Unable to build plan object",
				"certs element did not build properly",
			)
		}
		resp.Diagnostics.Append(respDiags...)
		finalPlanElements = append(finalPlanElements, finalPlanValue)
	}

	resp.PlanValue, respDiags = types.ListValue(req.PlanValue.ElementType(ctx), finalPlanElements)
	if respDiags.HasError() {
		resp.Diagnostics.AddError(
			"Unable to build plan object",
			"certs object did not build properly",
		)
	}
	resp.Diagnostics.Append(respDiags...)
}

func ValidateX509FileData() x509FileData {
	return x509FileData{}
}
