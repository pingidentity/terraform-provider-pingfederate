// Copyright © 2026 Ping Identity Corporation

package clustersettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func (r *clusterSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 13.0.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1300 := compare >= 0
	var plan *clusterSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	// If any of these fields are set by the user and the PF version is not new enough, throw an error
	if !pfVersionAtLeast1300 {
		if internaltypes.IsDefined(plan.ReplicateLogSettingsOnSave) {
			version.AddUnsupportedAttributeError("replicate_log_settings_on_save",
				r.providerConfig.ProductVersion, version.PingFederate1300, &resp.Diagnostics)
		}
	}

	// Set default if version is new enough
	if plan.ReplicateClientsOnSave.IsUnknown() {
		if pfVersionAtLeast1300 {
			plan.ReplicateClientsOnSave = types.BoolValue(false)
		} else {
			plan.ReplicateClientsOnSave = types.BoolNull()
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}
