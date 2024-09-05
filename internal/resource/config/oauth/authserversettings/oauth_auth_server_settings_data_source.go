package oauthauthserversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthAuthServerSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthAuthServerSettingsDataSource{}
)

// OauthAuthServerSettingsDataSource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsDataSource() datasource.DataSource {
	return &oauthAuthServerSettingsDataSource{
		impl: oauthServerSettingsDataSource{},
	}
}

// oauthAuthServerSettingsDataSource is the datasource implementation.
type oauthAuthServerSettingsDataSource struct {
	impl oauthServerSettingsDataSource
}

// GetSchema defines the schema for the datasource.
func (r *oauthAuthServerSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_oauth_auth_server_settings` datasource is deprecated. Use the `pingfederate_oauth_server_settings` datasource instead."
}

// Metadata returns the datasource type name.
func (r *oauthAuthServerSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings"
}

func (r *oauthAuthServerSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *oauthAuthServerSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
