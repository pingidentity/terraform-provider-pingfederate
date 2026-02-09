// Copyright © 2026 Ping Identity Corporation

package configstore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *configStoreResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *configStoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if plan == nil || state == nil {
		return
	}

	// If the type of the config store has changed, require replacement
	if internaltypes.IsDefined(plan.ListValue) && state.ListValue.IsNull() {
		resp.RequiresReplace = path.Paths{path.Root("list_value")}
	} else if internaltypes.IsDefined(plan.MapValue) && state.MapValue.IsNull() {
		resp.RequiresReplace = path.Paths{path.Root("map_value")}
	} else if internaltypes.IsDefined(plan.StringValue) && state.StringValue.IsNull() {
		resp.RequiresReplace = path.Paths{path.Root("string_value")}
	}
}

func (model *configStoreResourceModel) setType(configStore *client.ConfigStoreSetting) {
	if configStore.StringValue != nil {
		configStore.Type = "STRING"
	} else if configStore.MapValue != nil {
		configStore.Type = "MAP"
	} else {
		configStore.Type = "LIST"
	}
}
