package serversettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsLogSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsLogSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsLogSettingsResource{}
)

// ServerSettingsLogSettingsResource is a helper function to simplify the provider implementation.
func ServerSettingsLogSettingsResource() resource.Resource {
	return &serverSettingsLogSettingsResource{}
}

// serverSettingsLogSettingsResource is the resource implementation.
type serverSettingsLogSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsLogSettingsResourceModel struct {
	Id            types.String `tfsdk:"id"`
	LogCategories types.Set    `tfsdk:"log_categories"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsLogSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	serverSettingsLogSettingsResourceSchema(ctx, req, resp, false)
}

func serverSettingsLogSettingsResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "LogSettings Settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder for Terraform",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. On a PUT request, if a category is not included in the list, it will be disabled.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the log category. This field must match one of the category IDs defined in log4j-categories.xml.",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled. The default is false..",
							Required:    true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
	resp.Schema = schema
}

func addOptionalServerSettingsLogSettingsFields(ctx context.Context, addRequest *client.LogSettings, plan serverSettingsLogSettingsResourceModel) error {
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

// Metadata returns the resource type name.
func (r *serverSettingsLogSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_log_settings"
}

func (r *serverSettingsLogSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsLogSettingsResponse(ctx context.Context, r *client.LogSettings, state *serverSettingsLogSettingsResourceModel) {
	state.Id = types.StringValue("id")
	logCategoriesAttrTypes := map[string]attr.Type{
		"id":          basetypes.StringType{},
		"name":        basetypes.StringType{},
		"description": basetypes.StringType{},
		"enabled":     basetypes.BoolType{},
	}
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
	state.LogCategories = LogCategorySlice
}

func (r *serverSettingsLogSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsLogSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServerSettingsLogSettings := client.NewLogSettings()
	err := addOptionalServerSettingsLogSettingsFields(ctx, createServerSettingsLogSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ServerSettingsLogSettings", err.Error())
		return
	}
	requestJson, err := createServerSettingsLogSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateServerSettingsLogSettings := r.apiClient.ServerSettingsApi.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsLogSettings = apiCreateServerSettingsLogSettings.Body(*createServerSettingsLogSettings)
	serverSettingsLogSettingsResponse, httpResp, err := r.apiClient.ServerSettingsApi.UpdateLogSettingsExecute(apiCreateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ServerSettingsLogSettings", err, httpResp)
		return
	}
	responseJson, err := serverSettingsLogSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state serverSettingsLogSettingsResourceModel

	readServerSettingsLogSettingsResponse(ctx, serverSettingsLogSettingsResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serverSettingsLogSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readServerSettingsLogSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readServerSettingsLogSettings(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state serverSettingsLogSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsLogSettings, httpResp, err := apiClient.ServerSettingsApi.GetLogSettings(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a ServerSettingsLogSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadServerSettingsLogSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readServerSettingsLogSettingsResponse(ctx, apiReadServerSettingsLogSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsLogSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateServerSettingsLogSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateServerSettingsLogSettings(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan serverSettingsLogSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state serverSettingsLogSettingsResourceModel
	req.State.Get(ctx, &state)
	updateServerSettingsLogSettings := apiClient.ServerSettingsApi.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewLogSettings()
	err := addOptionalServerSettingsLogSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ServerSettingsLogSettings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateServerSettingsLogSettings = updateServerSettingsLogSettings.Body(*createUpdateRequest)
	updateServerSettingsLogSettingsResponse, httpResp, err := apiClient.ServerSettingsApi.UpdateLogSettingsExecute(updateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating ServerSettingsLogSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateServerSettingsLogSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readServerSettingsLogSettingsResponse(ctx, updateServerSettingsLogSettingsResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsLogSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsLogSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importServerSettingsLogSettingsLocation(ctx, req, resp)
}
func importServerSettingsLogSettingsLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
