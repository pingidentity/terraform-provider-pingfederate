// Copyright © 2025 Ping Identity Corporation

package keypairssigningrotationsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (r *keypairsSigningKeyRotationSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *keypairsSigningKeyRotationSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if plan == nil {
		return
	}

	// The activation buffer must be less than or equal to the creation buffer
	if internaltypes.IsDefined(plan.ActivationBufferDays) && internaltypes.IsDefined(plan.CreationBufferDays) && plan.ActivationBufferDays.ValueInt64() > plan.CreationBufferDays.ValueInt64() {
		resp.Diagnostics.AddError("activation_buffer_days must be less than or equal to creation_buffer_days",
			fmt.Sprintf("activation_buffer_days: %d, creation_buffer_days: %d", plan.ActivationBufferDays.ValueInt64(), plan.CreationBufferDays.ValueInt64()))
	}
}
