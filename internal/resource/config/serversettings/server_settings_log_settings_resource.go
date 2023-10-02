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
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
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
	schema := schema.Schema{
		Description: "LogSettings Settings related to server logging.",
		Attributes: map[string]schema.Attribute{
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

	config.AddCommonSchema(&schema)
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
	//TODO placeholder?
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings Log Settings", err.Error())
		return
	}
	_, requestErr := createServerSettingsLogSettings.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Server Settings Log Settings: %s", requestErr.Error())
	}

	apiCreateServerSettingsLogSettings := r.apiClient.ServerSettingsApi.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsLogSettings = apiCreateServerSettingsLogSettings.Body(*createServerSettingsLogSettings)
	serverSettingsLogSettingsResponse, httpResp, err := r.apiClient.ServerSettingsApi.UpdateLogSettingsExecute(apiCreateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings Log Settings", err, httpResp)
		return
	}
	_, responseErr := serverSettingsLogSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings Log Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	var state serverSettingsLogSettingsResourceModel

	readServerSettingsLogSettingsResponse(ctx, serverSettingsLogSettingsResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsLogSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsLogSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsLogSettings, httpResp, err := r.apiClient.ServerSettingsApi.GetLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadServerSettingsLogSettings.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings Log Settings: %s", responseErr.Error())
	}
	// Read the response into the state
	readServerSettingsLogSettingsResponse(ctx, apiReadServerSettingsLogSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsLogSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	updateServerSettingsLogSettings := r.apiClient.ServerSettingsApi.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewLogSettings()
	err := addOptionalServerSettingsLogSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings Log Settings", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Server Settings Log Settings: %s", requestErr.Error())
	}
	updateServerSettingsLogSettings = updateServerSettingsLogSettings.Body(*createUpdateRequest)
	updateServerSettingsLogSettingsResponse, httpResp, err := r.apiClient.ServerSettingsApi.UpdateLogSettingsExecute(updateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings Log Settings", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateServerSettingsLogSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings Log Settings: %s", responseErr.Error())
	}
	// Read the response
	readServerSettingsLogSettingsResponse(ctx, updateServerSettingsLogSettingsResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsLogSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsLogSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
