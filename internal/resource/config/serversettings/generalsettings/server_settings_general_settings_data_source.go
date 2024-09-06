package serversettingsgeneralsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsGeneralSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsGeneralSettingsDataSource{}
)

// ServerSettingsGeneralSettingsDataSource is a helper function to simplify the provider implementation.
func ServerSettingsGeneralSettingsDataSource() datasource.DataSource {
	return &serverSettingsGeneralSettingsDataSource{
		impl: serverSettingsGeneralDataSource{},
	}
}

// serverSettingsGeneralSettingsDataSource is the datasource implementation.
type serverSettingsGeneralSettingsDataSource struct {
	impl serverSettingsGeneralDataSource
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsGeneralSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_server_settings_general_settings` datasource is deprecated. Use the `pingfederate_server_settings_general` datasource instead."
}

// Metadata returns the datasource type name.
func (r *serverSettingsGeneralSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general_settings"
}

func (r *serverSettingsGeneralSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *serverSettingsGeneralSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
