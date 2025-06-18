// Copyright © 2025 Ping Identity Corporation

package idptospadaptermapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *idpToSpAdapterMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model *idpToSpAdapterMappingResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if model == nil || resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(model.AttributeContractFulfillment) {
		subjectKeyFound := false
		for key := range model.AttributeContractFulfillment.Elements() {
			if key == "subject" {
				subjectKeyFound = true
				break
			}
		}

		if !subjectKeyFound {
			resp.Diagnostics.AddAttributeError(
				path.Root("attribute_contract_fulfillment"),
				providererror.InvalidAttributeConfiguration,
				"attribute_contract_fulfillment.subject is required")
		}
	}
}
