package sessionsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
)

type sessionSettingsModel struct {
	TrackAdapterSessionsForLogout types.Bool  `tfsdk:"track_adapter_sessions_for_logout"`
	RevokeUserSessionOnLogout     types.Bool  `tfsdk:"revoke_user_session_on_logout"`
	SessionRevocationLifetime     types.Int64 `tfsdk:"session_revocation_lifetime"`
}

func readSessionSettingsResponse(ctx context.Context, r *client.SessionSettings, state *sessionSettingsModel) {
	state.TrackAdapterSessionsForLogout = types.BoolPointerValue(r.TrackAdapterSessionsForLogout)
	state.RevokeUserSessionOnLogout = types.BoolPointerValue(r.RevokeUserSessionOnLogout)
	state.SessionRevocationLifetime = types.Int64PointerValue(r.SessionRevocationLifetime)
}
