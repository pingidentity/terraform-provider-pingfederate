package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type sessionApplicationPolicyModel struct {
	Id              types.String `tfsdk:"id"`
	IdleTimeoutMins types.Int64  `tfsdk:"idle_timeout_mins"`
	MaxTimeoutMins  types.Int64  `tfsdk:"max_timeout_mins"`
}

func readSessionApplicationPolicyResponse(ctx context.Context, r *client.ApplicationSessionPolicy, state *sessionApplicationPolicyModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.IdleTimeoutMins = types.Int64PointerValue(r.IdleTimeoutMins)
	state.MaxTimeoutMins = types.Int64PointerValue(r.MaxTimeoutMins)
}
