package serversettingslogsettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsLoggingResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsLoggingResource{}
	_ resource.ResourceWithImportState = &serverSettingsLoggingResource{}

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

// ServerSettingsLoggingResource is a helper function to simplify the provider implementation.
func ServerSettingsLoggingResource() resource.Resource {
	return &serverSettingsLoggingResource{}
}

// serverSettingsLoggingResource is the resource implementation.
type serverSettingsLoggingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsLoggingResourceModel struct {
	LogCategories    types.Set `tfsdk:"log_categories"`
	LogCategoriesAll types.Set `tfsdk:"log_categories_all"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsLoggingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings related to server logging.",
		Attributes: map[string]schema.Attribute{
			"log_categories": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled.",
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only.",
							Optional:    false,
							Computed:    true,
							// Adding these plan modifiers also seems to cause issues with Terraform's set planning logic. Possibly related to the issue linked below
							/*PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},*/
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category. This field is read-only.",
							Optional:    false,
							Computed:    true,
							// Adding these plan modifiers also seems to cause issues with Terraform's set planning logic. Possibly related to the issue linked below
							/*PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},*/
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled.",
							Required:    true,
							// This default causes issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
							// Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
			"log_categories_all": schema.SetNestedAttribute{
				Description: "The log categories defined for the system and whether they are enabled. This attribute is read-only and will include any categories returned by PingFederate that were not specified in the normal log_categories attribute.",
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the log category. This field must match one of the category IDs defined in log4j-categories.xml.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The description of the log category.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the log category.",
							Optional:    false,
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Determines whether or not the log category is enabled.",
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}
	resp.Schema = schema
}

func addOptionalServerSettingsLoggingFields(ctx context.Context, addRequest *client.LogSettings, plan serverSettingsLoggingResourceModel) error {
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
func (r *serverSettingsLoggingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_logging"
}

func (r *serverSettingsLoggingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsLoggingResourceResponse(ctx context.Context, r *client.LogSettings, plan *serverSettingsLoggingResourceModel, state *serverSettingsLoggingResourceModel, isImport bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

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
		}
		unplannedCategories = append(unplannedCategories, resultCategory)
	}

	// Build results
	// On import, just read directly into log_categories
	if isImport {
		state.LogCategories, respDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, unplannedCategories)
		diags.Append(respDiags...)
	} else {
		state.LogCategories, respDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, plannedCategories)
		diags.Append(respDiags...)
	}
	state.LogCategoriesAll, respDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, unplannedCategories)
	diags.Append(respDiags...)
	return diags
}

func (r *serverSettingsLoggingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsLoggingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServerSettingsLogging := client.NewLogSettings()
	err := addOptionalServerSettingsLoggingFields(ctx, createServerSettingsLogging, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Server Settings Log Settings: "+err.Error())
		return
	}

	apiCreateServerSettingsLogging := r.apiClient.ServerSettingsAPI.UpdateLogSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsLogging = apiCreateServerSettingsLogging.Body(*createServerSettingsLogging)
	serverSettingsLoggingResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateLogSettingsExecute(apiCreateServerSettingsLogging)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings Log Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsLoggingResourceModel
	diags = readServerSettingsLoggingResourceResponse(ctx, serverSettingsLoggingResponse, &plan, &state, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsLoggingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state serverSettingsLoggingResourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadServerSettingsLogging, httpResp, err := r.apiClient.ServerSettingsAPI.GetLogSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Server Settings Log Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings Log Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readServerSettingsLoggingResourceResponse(ctx, apiReadServerSettingsLogging, &state, &state, isImportRead)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsLoggingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsLoggingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateServerSettingsLogging := r.apiClient.ServerSettingsAPI.UpdateLogSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewLogSettings()
	err := addOptionalServerSettingsLoggingFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Server Settings Log Settings: "+err.Error())
		return
	}

	updateServerSettingsLogging = updateServerSettingsLogging.Body(*createUpdateRequest)
	updateServerSettingsLoggingResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateLogSettingsExecute(updateServerSettingsLogging)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings Log Settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsLoggingResourceModel
	diags = readServerSettingsLoggingResourceResponse(ctx, updateServerSettingsLoggingResponse, &plan, &state, false)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsLoggingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	serverLogSettingsClientData := client.NewLogSettings()
	serverLogSettingsApiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateLogSettings(config.AuthContext(ctx, r.providerConfig))
	serverLogSettingsApiUpdateRequest = serverLogSettingsApiUpdateRequest.Body(*serverLogSettingsClientData)
	_, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateLogSettingsExecute(serverLogSettingsApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the Server Log Settings", err, httpResp)
	}
}

func (r *serverSettingsLoggingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState serverSettingsLoggingResourceModel
	emptyState.LogCategories = types.SetNull(types.ObjectType{AttrTypes: logCategoriesAttrTypes})
	emptyState.LogCategoriesAll = types.SetNull(types.ObjectType{AttrTypes: logCategoriesAttrTypes})
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
