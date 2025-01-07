package authenticationapisettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

type authenticationApiSettingsModel struct {
	ApiEnabled                       types.Bool   `tfsdk:"api_enabled"`
	EnableApiDescriptions            types.Bool   `tfsdk:"enable_api_descriptions"`
	RestrictAccessToRedirectlessMode types.Bool   `tfsdk:"restrict_access_to_redirectless_mode"`
	IncludeRequestContext            types.Bool   `tfsdk:"include_request_context"`
	DefaultApplicationRef            types.Object `tfsdk:"default_application_ref"`
}

// Read a AuthenticationApiSettingsResponse object into the model struct
func readAuthenticationApiSettingsResponse(ctx context.Context, r *client.AuthnApiSettings, state *authenticationApiSettingsModel) diag.Diagnostics {
	state.ApiEnabled = types.BoolPointerValue(r.ApiEnabled)
	state.EnableApiDescriptions = types.BoolPointerValue(r.EnableApiDescriptions)
	state.RestrictAccessToRedirectlessMode = types.BoolPointerValue(r.RestrictAccessToRedirectlessMode)
	state.IncludeRequestContext = types.BoolPointerValue(r.IncludeRequestContext)
	resourceLinkObjectValue, valueFromDiags := resourcelink.ToState(ctx, r.DefaultApplicationRef)
	state.DefaultApplicationRef = resourceLinkObjectValue

	return valueFromDiags
}
