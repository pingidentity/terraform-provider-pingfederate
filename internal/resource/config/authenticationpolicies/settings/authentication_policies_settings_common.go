// Copyright © 2025 Ping Identity Corporation

package authenticationpoliciessettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

type authenticationPoliciesSettingsModel struct {
	EnableIdpAuthnSelection types.Bool `tfsdk:"enable_idp_authn_selection"`
	EnableSpAuthnSelection  types.Bool `tfsdk:"enable_sp_authn_selection"`
}

func readAuthenticationPoliciesSettingsResponse(r *client.AuthenticationPoliciesSettings, state *authenticationPoliciesSettingsModel) {
	state.EnableIdpAuthnSelection = types.BoolPointerValue(r.EnableIdpAuthnSelection)
	state.EnableSpAuthnSelection = types.BoolPointerValue(r.EnableSpAuthnSelection)
}
