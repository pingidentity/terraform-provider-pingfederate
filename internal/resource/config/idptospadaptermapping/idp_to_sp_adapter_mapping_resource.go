package idptospadaptermapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *idpToSpAdapterMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model idpToSpAdapterMappingResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
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
			resp.Diagnostics.AddError("attribute_contract_fulfillment.subject is required", "")
		}
	}
}
