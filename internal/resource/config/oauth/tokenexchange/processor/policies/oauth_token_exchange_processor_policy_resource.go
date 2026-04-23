// Copyright © 2026 Ping Identity Corporation

package oauthtokenexchangeprocessorpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
)

func (r *oauthTokenExchangeProcessorPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *oauthTokenExchangeProcessorPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if plan == nil {
		return
	}

	attributeContractAttributesAttrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attributeContractAttributesElementType := types.ObjectType{AttrTypes: attributeContractAttributesAttrTypes}
	attributeContractExtendedAttributesDefault, diags := types.SetValue(attributeContractAttributesElementType, nil)
	resp.Diagnostics.Append(diags...)

	// If there is no state, set the default value for extended_attributes
	if state == nil {
		if plan.AttributeContract.IsUnknown() {
			plan.AttributeContract, diags = types.ObjectValue(map[string]attr.Type{
				"core_attributes":     types.SetType{ElemType: attributeContractAttributesElementType},
				"extended_attributes": types.SetType{ElemType: attributeContractAttributesElementType},
			}, map[string]attr.Value{
				"core_attributes":     types.SetUnknown(attributeContractAttributesElementType),
				"extended_attributes": attributeContractExtendedAttributesDefault,
			})
			resp.Diagnostics.Append(diags...)
			resp.Plan.Set(ctx, plan)
		}
	} else if plan.AttributeContract.IsUnknown() && !state.AttributeContract.IsUnknown() && !state.AttributeContract.IsNull() {
		// if the attribute_contract is not defined, maintain the core_attributes value from state,
		// and set extended_attributes to empty set
		stateCoreAttributes := state.AttributeContract.Attributes()["core_attributes"]
		plan.AttributeContract, diags = types.ObjectValue(map[string]attr.Type{
			"core_attributes":     types.SetType{ElemType: attributeContractAttributesElementType},
			"extended_attributes": types.SetType{ElemType: attributeContractAttributesElementType},
		}, map[string]attr.Value{
			"core_attributes":     stateCoreAttributes,
			"extended_attributes": attributeContractExtendedAttributesDefault,
		})
		resp.Diagnostics.Append(diags...)
		resp.Plan.Set(ctx, plan)
	}

}

func (r *oauthTokenExchangeProcessorPolicyResource) getOauthTokenExchangeProcessorPolicyByID(ctx context.Context, policyID string, diagnostics *diag.Diagnostics, action string) (*client.TokenExchangeProcessorPolicy, bool) {
	response, httpResp, err := r.apiClient.OauthTokenExchangeProcessorAPI.GetOauthTokenExchangeProcessorPolicyById(config.AuthContext(ctx, r.providerConfig), policyID).Execute()
	if err != nil {
		config.ReportHttpError(ctx, diagnostics, action, err, httpResp)
		return nil, false
	}

	return response, true
}
