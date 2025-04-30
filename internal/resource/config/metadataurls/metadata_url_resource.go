// Copyright © 2025 Ping Identity Corporation

package metadataurls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/pemcertificates"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *metadataUrlResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	// Handle drift detection for the x509_file.file_data value changing outside of terraform
	var plan, state metadataUrlResourceModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	if internaltypes.IsDefined(plan.X509File) && internaltypes.IsDefined(state.X509File) {
		planX509Attrs := plan.X509File.Attributes()
		planFileData := planX509Attrs["file_data"].(types.String).ValueString()
		stateFormattedFileData := state.X509File.Attributes()["formatted_file_data"].(types.String).ValueString()
		if !pemcertificates.FileDataEquivalent(planFileData, stateFormattedFileData) {
			planX509Attrs["formatted_file_data"] = types.StringUnknown()
			var diags diag.Diagnostics
			plan.X509File, diags = types.ObjectValue(plan.X509File.AttributeTypes(ctx), planX509Attrs)
			resp.Diagnostics.Append(diags...)
			plan.CertView = types.ObjectUnknown(plan.CertView.AttributeTypes(ctx))
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		}
	}
}

func (state *metadataUrlResourceModel) readClientResponseX509File(response *client.MetadataUrl) diag.Diagnostics {
	var diags diag.Diagnostics
	x509FileAttrTypes := map[string]attr.Type{
		"crypto_provider":     types.StringType,
		"file_data":           types.StringType,
		"formatted_file_data": types.StringType,
		"id":                  types.StringType,
	}
	var x509FileValue types.Object
	if response.X509File == nil {
		x509FileValue = types.ObjectNull(x509FileAttrTypes)
	} else {
		// Get the current file_data value
		fileDataAttr := types.StringNull()
		if internaltypes.IsDefined(state.X509File) {
			fileDataAttr = state.X509File.Attributes()["file_data"].(types.String)
		}
		// Get the current id value from the cert view - pf will store the value there and won't
		// return the id value in the x509 attribute.
		// Note that this method assumes the response CertView has already been set in state.
		idAttr := types.StringNull()
		if internaltypes.IsDefined(state.CertView) {
			idAttr = state.CertView.Attributes()["id"].(types.String)
		}
		x509FileValue, diags = types.ObjectValue(x509FileAttrTypes, map[string]attr.Value{
			"crypto_provider":     types.StringPointerValue(response.X509File.CryptoProvider),
			"file_data":           fileDataAttr,
			"formatted_file_data": types.StringValue(response.X509File.FileData),
			"id":                  idAttr,
		})
	}

	state.X509File = x509FileValue
	return diags
}
