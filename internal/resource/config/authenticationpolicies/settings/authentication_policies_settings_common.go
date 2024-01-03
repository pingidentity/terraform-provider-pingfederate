package authenticationpoliciessettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type authenticationPoliciesSettingsModel struct {
	Id                      types.String `tfsdk:"id"`
	EnableIdpAuthnSelection types.Bool   `tfsdk:"enable_idp_authn_selection"`
	EnableSpAuthnSelection  types.Bool   `tfsdk:"enable_sp_authn_selection"`
}

func readAuthenticationPoliciesSettingsResponse(r *client.AuthenticationPoliciesSettings, state *authenticationPoliciesSettingsModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.EnableIdpAuthnSelection = types.BoolPointerValue(r.EnableIdpAuthnSelection)
	state.EnableSpAuthnSelection = types.BoolPointerValue(r.EnableSpAuthnSelection)
}
