package serversettingslogsettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
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
		Description: "LogSettings settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. On a PUT request, if a category is not included in the list, it will be disabled.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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

func addOptionalServerSettingsLogSettingsFields(ctx context.Context, addRequest *client.LogSettings, plan serverSettingsLogSettingsDataSourceModel) error {
	if internaltypes.IsDefined(plan.LogCategories) {
		addRequest.LogCategories = []client.LogCategorySettings{}
		for _, logCategoriesSetting := range plan.LogCategories.Elements() {
			unmarshalled := client.LogCategorySettings{}
			err := json.Unmarshal([]byte(internaljson.FromValue(logCategoriesSetting, false)), &unmarshalled)
			if err != nil {
				return err
			}
			addRequest.LogCategories = append(addRequest.LogCategories, unmarshalled)
		}
	}
	return nil

}

// Metadata returns the datasource type name.
func (r *serverSettingsLogSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_log_settings"
}

func (r *serverSettingsLogSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.DataSourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsLogSettingsResponse(ctx context.Context, r *client.LogSettings, state *serverSettingsLogSettingsDataSourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
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
	LogCategorySlice, respDiags := types.SetValue(LogCategorySliceType, LogCategorySliceAttrVal)
	diags.Append(respDiags...)
	state.LogCategories = LogCategorySlice
	return diags
}

func (r *serverSettingsLogSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsLogSettingsDataSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsLogSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
			resp.State.RemoveDataSource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsLogSettingsResponse(ctx, apiReadServerSettingsLogSettings, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsLogSettingsDataSource) ImportState(ctx context.Context, req datasource.ImportStateRequest, resp *datasource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	datasource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
