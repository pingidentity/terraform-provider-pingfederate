package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthOpenIdConnectPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthOpenIdConnectPolicyDataSource{}
)

// OauthOpenIdConnectPolicyDataSource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPolicyDataSource() datasource.DataSource {
	return &oauthOpenIdConnectPolicyDataSource{
		impl: openidConnectPolicyDataSource{},
	}
}

// oauthOpenIdConnectPolicyDataSource is the datasource implementation.
type oauthOpenIdConnectPolicyDataSource struct {
	impl openidConnectPolicyDataSource
}

// GetSchema defines the schema for the datasource.
func (r *oauthOpenIdConnectPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_oauth_open_id_connect_policy` datasource is deprecated. Use the `pingfederate_openid_connect_policy` datasource instead."
}

// Metadata returns the datasource type name.
func (r *oauthOpenIdConnectPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policy"
}

func (r *oauthOpenIdConnectPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *oauthOpenIdConnectPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
