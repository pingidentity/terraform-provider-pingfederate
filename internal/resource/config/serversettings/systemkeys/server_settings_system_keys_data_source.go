package serversettingssystemkeys

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsSystemKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsSystemKeysDataSource{}
)

// ServerSettingsSystemKeysDataSource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysDataSource() datasource.DataSource {
	return &serverSettingsSystemKeysDataSource{}
}

// serverSettingsSystemKeysDataSource is the resource implementation.
type serverSettingsSystemKeysDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsSystemKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the server system keys.",
		Attributes: map[string]schema.Attribute{
			"current": schema.SingleNestedAttribute{
				Description: "Current SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Previous SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"pending": schema.SingleNestedAttribute{
				Description: "Pending SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *serverSettingsSystemKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys"
}

func (r *serverSettingsSystemKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *serverSettingsSystemKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsSystemKeysModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadServerSettingsSystemKeys, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readServerSettingsSystemKeysResponse(ctx, apiReadServerSettingsSystemKeys, &state, pointers.String("server_settings_system_keys_id"))
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
