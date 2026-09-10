// Copyright © 2026 Ping Identity Corporation

package oauthtokenexchangeprocessorpolicies

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
)

const maxRetries = 5

// Retry creates for this resource because sometimes PF silently fails to create them
func (r *oauthTokenExchangeProcessorPolicyResource) exponentialBackOffRetryCreate(ctx context.Context, apiCreateRequest client.ApiCreateOauthTokenExchangeProcessorPolicyRequest, policyId string) (*client.TokenExchangeProcessorPolicy, *http.Response, error) {
	var responseData *client.TokenExchangeProcessorPolicy
	var httpResp *http.Response
	var err error
	backOffTime := time.Second

	for i := 0; i < maxRetries; i++ {
		responseData, httpResp, err = r.apiClient.OauthTokenExchangeProcessorAPI.CreateOauthTokenExchangeProcessorPolicyExecute(apiCreateRequest)
		if err != nil {
			// If PF returned an error, don't retry
			return responseData, httpResp, err
		}

		// If PF returned success, ensure the resource was actually created
		_, readHttpResp, readErr := r.apiClient.OauthTokenExchangeProcessorAPI.GetOauthTokenExchangeProcessorPolicyById(config.AuthContext(ctx, r.providerConfig), policyId).Execute()
		if readErr == nil || readHttpResp == nil || readHttpResp.StatusCode != 404 {
			// Either the create succeeded or the read returned a non-404 status code
			return responseData, httpResp, err
		}

		backOffTime = backOffTime * 2
	}

	tflog.Info(context.Background(), fmt.Sprintf("Request failed after %d attempts", maxRetries))

	return responseData, httpResp, err
}

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
