package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionApplicationSessionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionApplicationSessionPolicyDataSource{}
)

// SessionApplicationSessionPolicyDataSource is a helper function to simplify the provider implementation.
func SessionApplicationSessionPolicyDataSource() datasource.DataSource {
	return &sessionApplicationSessionPolicyDataSource{
		impl: sessionApplicationPolicyDataSource{},
	}
}

// sessionApplicationSessionPolicyDataSource is the datasource implementation.
type sessionApplicationSessionPolicyDataSource struct {
	impl sessionApplicationPolicyDataSource
}

// GetSchema defines the schema for the datasource.
func (r *sessionApplicationSessionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `session_application_session_policy` datasource is deprecated. Use the `session_application_policy` datasource instead."
}

// Metadata returns the datasource type name.
func (r *sessionApplicationSessionPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_application_session_policy"
}

func (r *sessionApplicationSessionPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *sessionApplicationSessionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
