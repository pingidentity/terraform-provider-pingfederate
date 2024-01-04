package sessionauthenticationsessionpoliciesglobal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type sessionAuthenticationSessionPoliciesGlobalModel struct {
	Id                         types.String `tfsdk:"id"`
	EnableSessions             types.Bool   `tfsdk:"enable_sessions"`
	PersistentSessions         types.Bool   `tfsdk:"persistent_sessions"`
	HashUniqueUserKeyAttribute types.Bool   `tfsdk:"hash_unique_user_key_attribute"`
	IdleTimeoutMins            types.Int64  `tfsdk:"idle_timeout_mins"`
	IdleTimeoutDisplayUnit     types.String `tfsdk:"idle_timeout_display_unit"`
	MaxTimeoutMins             types.Int64  `tfsdk:"max_timeout_mins"`
	MaxTimeoutDisplayUnit      types.String `tfsdk:"max_timeout_display_unit"`
}

func readSessionAuthenticationSessionPoliciesGlobalResponse(ctx context.Context, r *client.GlobalAuthenticationSessionPolicy, state *sessionAuthenticationSessionPoliciesGlobalModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.EnableSessions = types.BoolValue(r.EnableSessions)
	state.PersistentSessions = types.BoolPointerValue(r.PersistentSessions)
	state.HashUniqueUserKeyAttribute = types.BoolPointerValue(r.HashUniqueUserKeyAttribute)
	state.IdleTimeoutMins = types.Int64PointerValue(r.IdleTimeoutMins)
	state.IdleTimeoutDisplayUnit = types.StringPointerValue(r.IdleTimeoutDisplayUnit)
	state.MaxTimeoutMins = types.Int64PointerValue(r.MaxTimeoutMins)
	state.MaxTimeoutDisplayUnit = types.StringPointerValue(r.MaxTimeoutDisplayUnit)
}
