// Copyright © 2025 Ping Identity Corporation

package sessionauthenticationsessionpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *sessionAuthenticationPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config *sessionAuthenticationPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if config == nil || resp.Diagnostics.HasError() {
		return
	}

	if !config.EnableSessions.IsUnknown() && !config.EnableSessions.ValueBool() && config.Persistent.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("persistent"),
			providererror.InvalidAttributeConfiguration,
			"`persistent` cannot be true when `enable_sessions` is false")
	}

	if (config.IdleTimeoutMins.IsNull() && internaltypes.IsDefined(config.MaxTimeoutMins)) ||
		(config.MaxTimeoutMins.IsNull() && internaltypes.IsDefined(config.IdleTimeoutMins)) {
		resp.Diagnostics.AddError(
			providererror.InvalidAttributeConfiguration,
			"`idle_timeout_mins` and `max_timeout_mins` must either both be defined or both be undefined")
	}
}
