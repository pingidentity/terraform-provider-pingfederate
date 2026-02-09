// Copyright © 2026 Ping Identity Corporation

package oauthauthenticationpolicycontractmappings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *oauthAuthenticationPolicyContractMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model *oauthAuthenticationPolicyContractMappingResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if model == nil || resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(model.AttributeContractFulfillment) {
		userKeyFound := false
		userNameFound := false
		for key := range model.AttributeContractFulfillment.Elements() {
			if key == "USER_KEY" {
				userKeyFound = true
			}
			if key == "USER_NAME" {
				userNameFound = true
			}
		}

		if !userKeyFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("attribute_contract_fulfillment"),
				providererror.InvalidAttributeConfiguration,
				"attribute_contract_fulfillment.USER_KEY is required")
		}
		if !userNameFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("attribute_contract_fulfillment"),
				providererror.InvalidAttributeConfiguration,
				"attribute_contract_fulfillment.USER_NAME is required")
		}
	}
}
