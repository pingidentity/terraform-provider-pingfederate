// Copyright © 2025 Ping Identity Corporation

package serversettingslogsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsLoggingDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsLoggingDataSource{}
)

// ServerSettingsLoggingDataSource is a helper function to simplify the provider implementation.
func ServerSettingsLoggingDataSource() datasource.DataSource {
	return &serverSettingsLoggingDataSource{}
}

// serverSettingsLoggingDataSource is the datasource implementation.
type serverSettingsLoggingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsLoggingDataSourceModel struct {
	LogCategories types.Set `tfsdk:"log_categories"`
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsLoggingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the log category. This field must match one of the category IDs defined in log4j-categories.xml.",
							Computed:    true,
							Optional:    false,
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category.",
							Computed:    true,
							Optional:    false,
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category.",
							Computed:    true,
							Optional:    false,
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled.",
							Computed:    true,
							Optional:    false,
						},
					},
				},
			},
		},
	}
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *serverSettingsLoggingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_logging"
}

func (r *serverSettingsLoggingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsLoggingDataSourceResponse(ctx context.Context, r *client.LogSettings, state *serverSettingsLoggingDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.LogCategories, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, r.LogCategories)
	return diags
}

func (r *serverSettingsLoggingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsLoggingDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsLogging, httpResp, err := r.apiClient.ServerSettingsAPI.GetLogSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
		return
	}

	diags = readServerSettingsLoggingDataSourceResponse(ctx, apiReadServerSettingsLogging, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}
