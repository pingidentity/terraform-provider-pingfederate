package oauthtokenexchangegeneratorgroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *oauthTokenExchangeGeneratorGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *oauthTokenExchangeGeneratorGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}

	// Exactly one generator mapping must be marked as default
	numDefaults := 0
	// Each pair of requested_token_type and token_generator.id must be unique
	uniquePairs := map[string]bool{}
	for _, mapping := range plan.GeneratorMappings.Elements() {
		mappingAttrs := mapping.(types.Object).Attributes()
		if mappingAttrs["default_mapping"].(types.Bool).ValueBool() {
			numDefaults++
		}
		tokenType := mappingAttrs["requested_token_type"].(types.String).ValueString()
		tokenGeneratorAttributes := mappingAttrs["token_generator"].(types.Object).Attributes()
		tokenGeneratorId := tokenGeneratorAttributes["id"].(types.String).ValueString()
		pair := fmt.Sprintf("%s, %s", tokenType, tokenGeneratorId)
		if uniquePairs[pair] {
			resp.Diagnostics.AddError("Each generator mapping pair of `requested_token_type` and `token_generator.id` must be unique.", fmt.Sprintf("Duplicate pair: %s", pair))
		}
		uniquePairs[pair] = true
	}
	if numDefaults != 1 {
		resp.Diagnostics.AddError("Exactly one generator mapping must be marked as default by setting `default_mapping` to `true`.", fmt.Sprintf("Found %d default mappings", numDefaults))
	}
}
