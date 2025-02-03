// Copyright © 2025 Ping Identity Corporation

package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
)

type sessionApplicationPolicyModel struct {
	IdleTimeoutMins types.Int64 `tfsdk:"idle_timeout_mins"`
	MaxTimeoutMins  types.Int64 `tfsdk:"max_timeout_mins"`
}

func readSessionApplicationPolicyResponse(ctx context.Context, r *client.ApplicationSessionPolicy, state *sessionApplicationPolicyModel) {
	state.IdleTimeoutMins = types.Int64PointerValue(r.IdleTimeoutMins)
	state.MaxTimeoutMins = types.Int64PointerValue(r.MaxTimeoutMins)
}
