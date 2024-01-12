package serversettingslogsettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsLogSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsLogSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsLogSettingsResource{}

	logCategoriesAttrTypes = map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"enabled":     types.BoolType,
	}

	logCategoriesDefault, _ = types.SetValue(types.ObjectType{AttrTypes: logCategoriesAttrTypes}, []attr.Value{
		createDefaultLogCategoryObject("xmlsig"),
		createDefaultLogCategoryObject("core"),
		createDefaultLogCategoryObject("requestparams"),
		createDefaultLogCategoryObject("requestheaders"),
		createDefaultLogCategoryObject("trustedcas"),
		createDefaultLogCategoryObject("restdatastore"),
		createDefaultLogCategoryObject("policytree"),
	})
)

func createDefaultLogCategoryObject(id string) types.Object {
	defaultObj, _ := types.ObjectValue(logCategoriesAttrTypes, map[string]attr.Value{
		"id":          types.StringValue(id),
		"enabled":     types.BoolValue(false),
		"name":        types.StringUnknown(),
		"description": types.StringUnknown(),
	})
	return defaultObj
}

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
	Id               types.String `tfsdk:"id"`
	LogCategories    types.Set    `tfsdk:"log_categories"`
	LogCategoriesAll types.Set    `tfsdk:"log_categories_all"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsLogSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. On a PUT request, if a category is not included in the list, it will be disabled.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(logCategoriesDefault),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the log category. This field must match one of the category IDs defined in log4j-categories.xml.",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Optional:    false,
							Computed:    true,
							// Adding these plan modifiers also seems to cause issues with Terraform's set planning logic. Possibly related to the issue linked below
							/*PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},*/
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Optional:    false,
							Computed:    true,
							// Adding these plan modifiers also seems to cause issues with Terraform's set planning logic. Possibly related to the issue linked below
							/*PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},*/
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled. The default is false.",
							Optional:    true,
							Computed:    true,
							// This default causes issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
							// Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
			"log_categories_all": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. On a PUT request, if a category is not included in the list, it will be disabled. This attribute will include any categories not specified in the normal log_categories attribute.",
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the log category. This field must match one of the category IDs defined in log4j-categories.xml.",
							Optional:    false,
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Optional:    false,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only and is ignored for PUT requests.",
							Optional:    false,
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled. The default is false.",
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
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

func readServerSettingsLogSettingsResourceResponse(ctx context.Context, r *client.LogSettings, plan *serverSettingsLogSettingsResourceModel, state *serverSettingsLogSettingsResourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}

	// Build a list of log categories specified in the plan
	plannedIds := map[string]bool{}
	if plan != nil {
		for _, plannedCategory := range plan.LogCategories.Elements() {
			plannedCategoryId := plannedCategory.(types.Object).Attributes()["id"]
			if internaltypes.IsDefined(plannedCategoryId) {
				plannedIds[plannedCategoryId.(types.String).ValueString()] = true
			}
		}
	}

	// Build list of planned and unplanned log categories
	plannedCategories := []client.LogCategorySettings{}
	unplannedCategories := []client.LogCategorySettings{}
	for _, resultCategory := range r.LogCategories {
		_, isInPlan := plannedIds[resultCategory.Id]
		if isInPlan {
			plannedCategories = append(plannedCategories, resultCategory)
		} else {
			unplannedCategories = append(unplannedCategories, resultCategory)
		}
	}

	// Build results
	state.LogCategories, respDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, plannedCategories)
	diags.Append(respDiags...)
	state.LogCategoriesAll, respDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, unplannedCategories)
	diags.Append(respDiags...)
	return diags
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

	apiCreateServerSettingsLogSettings := r.apiClient.ServerSettingsAPI.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsLogSettings = apiCreateServerSettingsLogSettings.Body(*createServerSettingsLogSettings)
	serverSettingsLogSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateLogSettingsExecute(apiCreateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings Log Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsLogSettingsResourceModel
	diags = readServerSettingsLogSettingsResourceResponse(ctx, serverSettingsLogSettingsResponse, &plan, &state, nil)
	resp.Diagnostics.Append(diags...)
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
	apiReadServerSettingsLogSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
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
	diags = readServerSettingsLogSettingsResourceResponse(ctx, apiReadServerSettingsLogSettings, &state, &state, id)
	resp.Diagnostics.Append(diags...)

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

	updateServerSettingsLogSettings := r.apiClient.ServerSettingsAPI.UpdateLogSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewLogSettings()
	err := addOptionalServerSettingsLogSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings Log Settings", err.Error())
		return
	}

	updateServerSettingsLogSettings = updateServerSettingsLogSettings.Body(*createUpdateRequest)
	updateServerSettingsLogSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateLogSettingsExecute(updateServerSettingsLogSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings Log Settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsLogSettingsResourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsLogSettingsResourceResponse(ctx, updateServerSettingsLogSettingsResponse, &plan, &state, id)
	resp.Diagnostics.Append(diags...)

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
