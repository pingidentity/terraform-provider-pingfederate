package localidentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &localIdentityIdentityProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &localIdentityIdentityProfileDataSource{}
)

// LocalIdentityIdentityProfileDataSource is a helper function to simplify the provider implementation.
func LocalIdentityIdentityProfileDataSource() datasource.DataSource {
	return &localIdentityIdentityProfileDataSource{
		impl: localIdentityProfileDataSource{},
	}
}

// localIdentityIdentityProfileDataSource is the datasource implementation.
type localIdentityIdentityProfileDataSource struct {
	impl localIdentityProfileDataSource
}

// GetSchema defines the schema for the datasource.
func (r *localIdentityIdentityProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_local_identity_identity_profile` datasource is deprecated. Use the `pingfederate_local_identity_profile` datasource instead."
}

// Metadata returns the datasource type name.
func (r *localIdentityIdentityProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_identity_identity_profile"
}

func (r *localIdentityIdentityProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *localIdentityIdentityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
