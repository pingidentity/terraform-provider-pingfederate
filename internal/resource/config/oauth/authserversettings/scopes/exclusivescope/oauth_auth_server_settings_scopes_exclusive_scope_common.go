package oauthauthserversettingsscopesexclusivescope

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

type oauthAuthServerSettingsScopesExclusiveScopeModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Dynamic     types.Bool   `tfsdk:"dynamic"`
}

// Read a OauthAuthServerSettingsScopesExclusiveScopeResponse object into the model struct
func readOauthAuthServerSettingsScopesExclusiveScopeResponse(ctx context.Context, r *client.ScopeEntry, state *oauthAuthServerSettingsScopesExclusiveScopeModel) {
	state.Id = types.StringValue(r.Name)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(r.Description)
	state.Dynamic = types.BoolPointerValue(r.Dynamic)
}
