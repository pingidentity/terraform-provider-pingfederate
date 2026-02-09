// Copyright © 2026 Ping Identity Corporation

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/pemcertificates"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ planmodifier.List = &handleFormattedCertsListFileData{}

type handleFormattedCertsListFileData struct{}

func (v handleFormattedCertsListFileData) Description(ctx context.Context) string {
	return "Checks that the computed x509_file.formatted_file_data and user-supplied x509_file.file_data values represent the same value in a list of certs. This is to detect drift with server-formatted certificates."
}

func (v handleFormattedCertsListFileData) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v handleFormattedCertsListFileData) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
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
		if i >= len(req.StateValue.Elements()) || planElement.IsUnknown() {
			continue
		}
		stateValueAttrs := req.StateValue.Elements()[i].(types.Object).Attributes()
		planValue := planElement.(types.Object)
		planValueAttrs := planValue.Attributes()
		stateX509Value := stateValueAttrs["x509_file"].(types.Object)
		planX509Value := planValueAttrs["x509_file"].(types.Object)
		stateCertViewValue := stateValueAttrs["cert_view"].(types.Object)
		planCertViewValue := planValueAttrs["cert_view"].(types.Object)

		// Get file_data from plan
		fileData, ok := planX509Value.Attributes()["file_data"]
		if !ok {
			return
		}
		// Get formatted_file_data from state
		formattedFileData, ok := stateX509Value.Attributes()["formatted_file_data"].(types.String)
		if !ok {
			return
		}

		// Check if formatted_file_data and file_data strings match, or if formatted_file_data matches original string
		// If they do not, formatted_file_data is set to unknown
		if !pemcertificates.FileDataEquivalent(fileData.(types.String).ValueString(), formattedFileData.ValueString()) {
			reqPlanAttrs := planX509Value.Attributes()
			reqPlanAttrs["formatted_file_data"] = types.StringUnknown()
			planX509Value, respDiags = types.ObjectValue(planX509Value.AttributeTypes(ctx), reqPlanAttrs)
			if respDiags.HasError() {
				resp.Diagnostics.AddError(
					providererror.InternalProviderError,
					"Failed to construct x509_file object value with unknown formatted_file_data in certs list plan modifier",
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
					providererror.InternalProviderError,
					"Failed to construct x509_file object value with existing formatted_file_data in certs list plan modifier",
				)
			}

			certViewAttrs := stateCertViewValue.Attributes()
			// Handle if the id was changed between the plan and state
			planCertId := planX509Value.Attributes()["id"]
			if !planCertId.IsUnknown() && !planCertId.IsNull() {
				certViewAttrs["id"] = planCertId
			}

			planCertViewValue, respDiags = types.ObjectValue(planCertViewValue.AttributeTypes(ctx), certViewAttrs)
			resp.Diagnostics.Append(respDiags...)
			if respDiags.HasError() {
				resp.Diagnostics.AddError(
					providererror.InternalProviderError,
					"Failed to construct cert_view object value in certs list plan modifier",
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
				providererror.InternalProviderError,
				"Failed to construct cert object value in certs list plan modifier",
			)
		}
		resp.Diagnostics.Append(respDiags...)
		finalPlanElements = append(finalPlanElements, finalPlanValue)
	}

	resp.PlanValue, respDiags = types.ListValue(req.PlanValue.ElementType(ctx), finalPlanElements)
	if respDiags.HasError() {
		resp.Diagnostics.AddError(
			providererror.InternalProviderError,
			"Failed to construct certs list value in certs list plan modifier",
		)
	}
	resp.Diagnostics.Append(respDiags...)
}

func HandleFormattedCertsListFileData() planmodifier.List {
	return handleFormattedCertsListFileData{}
}
