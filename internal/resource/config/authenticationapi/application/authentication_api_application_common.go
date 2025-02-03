// Copyright © 2025 Ping Identity Corporation

package authenticationapiapplication

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

type authenticationApiApplicationModel struct {
	Id                           types.String `tfsdk:"id"`
	ApplicationId                types.String `tfsdk:"application_id"`
	Name                         types.String `tfsdk:"name"`
	Url                          types.String `tfsdk:"url"`
	Description                  types.String `tfsdk:"description"`
	AdditionalAllowedOrigins     types.Set    `tfsdk:"additional_allowed_origins"`
	ClientForRedirectlessModeRef types.Object `tfsdk:"client_for_redirectless_mode_ref"`
}

func readAuthenticationApiApplicationResponse(ctx context.Context, r *client.AuthnApiApplication, state *authenticationApiApplicationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Id = types.StringValue(r.Id)
	state.ApplicationId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.Url = types.StringValue(r.Url)
	state.Description = types.StringPointerValue(r.Description)
	state.AdditionalAllowedOrigins = internaltypes.GetStringSet(r.AdditionalAllowedOrigins)
	state.ClientForRedirectlessModeRef, diags = resourcelink.ToState(ctx, r.ClientForRedirectlessModeRef)

	return diags
}
