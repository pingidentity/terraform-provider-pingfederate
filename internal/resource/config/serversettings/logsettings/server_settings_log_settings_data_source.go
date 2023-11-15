package serversettingslogsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsLogSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsLogSettingsDataSource{}
)

// ServerSettingsLogSettingsDataSource is a helper function to simplify the provider implementation.
func NewServerSettingsLogSettingsDataSource() datasource.DataSource {
	return &serverSettingsLogSettingsDataSource{}
}

// serverSettingsLogSettingsDataSource is the datasource implementation.
type serverSettingsLogSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsLogSettingsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	LogCategories types.Set    `tfsdk:"log_categories"`
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsLogSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. On a PUT request, if a category is not included in the list, it will be disabled.",
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
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Computed:    true,
							Optional:    false,
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Computed:    true,
							Optional:    false,
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled. The default is false.",
							Computed:    true,
							Optional:    false,
						},
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *serverSettingsLogSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_log_settings"
}

func (r *serverSettingsLogSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsLogSettingsDataSource(ctx context.Context, r *client.LogSettings, state *serverSettingsLogSettingsDataSourceModel) {
	//	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue("server_settings_log_settings_id")
	logCategorySettings := r.GetLogCategories()
	var LogCategorySliceAttrVal = []attr.Value{}
	LogCategorySliceType := types.ObjectType{AttrTypes: logCategoriesAttrTypes}
	for i := 0; i < len(logCategorySettings); i++ {
		logCategoriesAttrValues := map[string]attr.Value{
			"id":          types.StringValue(logCategorySettings[i].Id),
			"name":        types.StringPointerValue(logCategorySettings[i].Name),
			"description": types.StringPointerValue(logCategorySettings[i].Description),
			"enabled":     types.BoolPointerValue(logCategorySettings[i].Enabled),
		}
		LogCategoryObj, _ := types.ObjectValue(logCategoriesAttrTypes, logCategoriesAttrValues)
		LogCategorySliceAttrVal = append(LogCategorySliceAttrVal, LogCategoryObj)
	}
	LogCategorySlice, _ := types.SetValue(LogCategorySliceType, LogCategorySliceAttrVal)
	//	diags.Append(respDiags...)
	state.LogCategories = LogCategorySlice
	// return diags
}

func (r *serverSettingsLogSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsLogSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsLogSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
		return
	}

	readServerSettingsLogSettingsDataSource(ctx, apiReadServerSettingsLogSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
