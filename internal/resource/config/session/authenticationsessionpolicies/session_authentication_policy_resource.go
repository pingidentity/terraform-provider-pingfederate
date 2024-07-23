package sessionauthenticationsessionpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *sessionAuthenticationPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config sessionAuthenticationPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(config.EnableSessions) && !config.EnableSessions.ValueBool() && internaltypes.IsDefined(config.Persistent) && config.Persistent.ValueBool() {
		resp.Diagnostics.AddError("`persistent` cannot be true when `enable_sessions` is false", "")
	}

	if internaltypes.IsDefined(config.IdleTimeoutMins) != internaltypes.IsDefined(config.MaxTimeoutMins) {
		resp.Diagnostics.AddError("`idle_timeout_mins` and `max_timeout_mins` must either both be defined or both be undefined", "")
	}
}
