package oauthidpadaptermappings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *oauthIdpAdapterMappingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *oauthIdpAdapterMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	// The idp adapter object will be configured based on the provided mapping_id
	mappingId := plan.MappingId.ValueString()
	idpAdapterRefValue, diags := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType,
	}, map[string]attr.Value{
		"id": types.StringValue(mappingId),
	})
	resp.Diagnostics.Append(diags...)

	plan.IdpAdapterRef = idpAdapterRefValue
	resp.Plan.Set(ctx, plan)
}
