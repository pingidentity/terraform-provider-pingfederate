// Copyright © 2025 Ping Identity Corporation

package clustersettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func (r *clusterSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.0.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1200 := compare >= 0
	// Compare to version 13.0.0 of PF
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1300 := compare >= 0
	// This endpoint was added in PingFederate 12.0
	if !pfVersionAtLeast1200 {
		version.AddUnsupportedResourceError("pingfederate_cluster_settings", r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
	}
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
}
