package serversettingslogsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsLogSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsLogSettingsDataSource{}
)

// ServerSettingsLogSettingsDataSource is a helper function to simplify the provider implementation.
func ServerSettingsLogSettingsDataSource() datasource.DataSource {
	return &serverSettingsLogSettingsDataSource{
		impl: serverSettingsLoggingDataSource{},
	}
}

// serverSettingsLogSettingsDataSource is the datasource implementation.
type serverSettingsLogSettingsDataSource struct {
	impl serverSettingsLoggingDataSource
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsLogSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_server_settings_log_settings` datasource is deprecated. Use the `pingfederate_server_settings_logging` datasource instead."
}

// Metadata returns the datasource type name.
func (r *serverSettingsLogSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_log_settings"
}

func (r *serverSettingsLogSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *serverSettingsLogSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}
