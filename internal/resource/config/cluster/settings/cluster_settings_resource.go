package clustersettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func (r *clusterSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// This endpoint was added in PingFederate 12.0
	// Compare to version 12.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	if !pfVersionAtLeast120 {
		version.AddUnsupportedResourceError("pingfederate_cluster_settings", r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
	}
}
