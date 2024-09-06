package sessionauthenticationsessionpoliciesglobal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionAuthenticationSessionPoliciesGlobalDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionAuthenticationSessionPoliciesGlobalDataSource{}
)

// SessionAuthenticationSessionPoliciesGlobalDataSource is a helper function to simplify the provider implementation.
func SessionAuthenticationSessionPoliciesGlobalDataSource() datasource.DataSource {
	return &sessionAuthenticationSessionPoliciesGlobalDataSource{
		impl: sessionAuthenticationPoliciesGlobalDataSource{},
	}
}

// sessionAuthenticationSessionPoliciesGlobalDataSource is the datasource implementation.
type sessionAuthenticationSessionPoliciesGlobalDataSource struct {
	impl sessionAuthenticationPoliciesGlobalDataSource
}

// GetSchema defines the schema for the datasource.
func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_session_authentication_session_policies_global` datasource is deprecated. Use the `pingfederate_session_authentication_policies_global` datasource instead."
}

// Metadata returns the datasource type name.
func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_authentication_session_policies_global"
}

func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
