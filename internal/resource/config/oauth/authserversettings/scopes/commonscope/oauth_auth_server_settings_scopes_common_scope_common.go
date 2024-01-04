package oauthauthserversettingsscopescommonscope

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
)

type oauthAuthServerSettingsScopesCommonScopeModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Dynamic     types.Bool   `tfsdk:"dynamic"`
}

// Read a OauthAuthServerSettingsScopesCommonScopeResponse object into the model struct
func readOauthAuthServerSettingsScopesCommonScopeResponse(ctx context.Context, r *client.ScopeEntry, state *oauthAuthServerSettingsScopesCommonScopeModel) {
	state.Id = types.StringValue(r.Name)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(r.Description)
	state.Dynamic = types.BoolPointerValue(r.Dynamic)
}
