package oauthAuthServerSettingsScopesCommonScopes

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsScopesCommonScopesResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsScopesCommonScopesResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsScopesCommonScopesResource{}
)

// OauthAuthServerSettingsScopesCommonScopesResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsScopesCommonScopesResource() resource.Resource {
	return &oauthAuthServerSettingsScopesCommonScopesResource{}
}

// oauthAuthServerSettingsScopesCommonScopesResource is the resource implementation.
type oauthAuthServerSettingsScopesCommonScopesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAuthServerSettingsScopesCommonScopesResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Dynamic     types.Bool   `tfsdk:"dynamic"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthAuthServerSettingsScopesCommonScopesResourceSchema(ctx, req, resp, false)
}

func oauthAuthServerSettingsScopesCommonScopesResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a OauthAuthServerSettingsScopesCommonScopes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed attribute tied to the name property of this resource.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "The name of the scope.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the scope that appears when the user is prompted for authorization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic": schema.BoolAttribute{
				Description: "True if the scope is dynamic. (Defaults to false)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "description"})
	}
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model oauthAuthServerSettingsScopesCommonScopesResourceModel
	req.Plan.Get(ctx, &model)
	if model.Dynamic.ValueBool() && (model.Name.ValueString() != "" || !model.Name.IsNull()) {
		{
			containsAsteriskPrefix := strings.Index(model.Name.ValueString(), "*")
			if containsAsteriskPrefix != 0 {
				resp.Diagnostics.AddError("Dynamic property is set to true with Name property incorrectly specified!", "The Name property must be prefixed with an \"*\". For example, \"*example\"")
			}
		}
	}
}

func addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx context.Context, addRequest *client.ScopeEntry, plan oauthAuthServerSettingsScopesCommonScopesResourceModel) error {

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueString()
	}
	if internaltypes.IsDefined(plan.Description) {
		addRequest.Description = plan.Description.ValueString()
	}
	if internaltypes.IsDefined(plan.Dynamic) {
		addRequest.Dynamic = plan.Dynamic.ValueBoolPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_authserversettings_scopes_commonscopes"
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAuthServerSettingsScopesCommonScopesResponse(ctx context.Context, r *client.ScopeEntry, state *oauthAuthServerSettingsScopesCommonScopesResourceModel, expectedValues *oauthAuthServerSettingsScopesCommonScopesResourceModel) {
	state.Id = types.StringValue(r.Name)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(r.Description)
	state.Dynamic = types.BoolPointerValue(r.Dynamic)
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsScopesCommonScopesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAuthServerSettingsScopesCommonScopes := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx, createOauthAuthServerSettingsScopesCommonScopes, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAuthServerSettingsScopesCommonScopes", err.Error())
		return
	}
	requestJson, err := createOauthAuthServerSettingsScopesCommonScopes.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthAuthServerSettingsScopesCommonScopes := r.apiClient.OauthAuthServerSettingsApi.AddCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettingsScopesCommonScopes = apiCreateOauthAuthServerSettingsScopesCommonScopes.Body(*createOauthAuthServerSettingsScopesCommonScopes)
	oauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsApi.AddCommonScopeExecute(apiCreateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthAuthServerSettingsScopesCommonScopes", err, httpResp)
		return
	}
	responseJson, err := oauthAuthServerSettingsScopesCommonScopesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel

	readOauthAuthServerSettingsScopesCommonScopesResponse(ctx, oauthAuthServerSettingsScopesCommonScopesResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthAuthServerSettingsScopesCommonScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOauthAuthServerSettingsScopesCommonScopes(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettingsScopesCommonScopes, httpResp, err := apiClient.OauthAuthServerSettingsApi.GetCommonScope(config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a OauthAuthServerSettingsScopesCommonScopes", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadOauthAuthServerSettingsScopesCommonScopes.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesCommonScopesResponse(ctx, apiReadOauthAuthServerSettingsScopesCommonScopes, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthAuthServerSettingsScopesCommonScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOauthAuthServerSettingsScopesCommonScopes(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan oauthAuthServerSettingsScopesCommonScopesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel
	req.State.Get(ctx, &state)
	updateOauthAuthServerSettingsScopesCommonScopes := apiClient.OauthAuthServerSettingsApi.UpdateCommonScope(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewScopeEntry(plan.Id.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAuthServerSettingsScopesCommonScopes", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthAuthServerSettingsScopesCommonScopes = updateOauthAuthServerSettingsScopesCommonScopes.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := apiClient.OauthAuthServerSettingsApi.UpdateCommonScopeExecute(updateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthAuthServerSettingsScopesCommonScopes", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthAuthServerSettingsScopesCommonScopesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthAuthServerSettingsScopesCommonScopesResponse(ctx, updateOauthAuthServerSettingsScopesCommonScopesResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteOauthAuthServerSettingsScopesCommonScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteOauthAuthServerSettingsScopesCommonScopes(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.OauthAuthServerSettingsApi.RemoveCommonScope(config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a OauthAuthServerSettingsScopesCommonScopes", err, httpResp)
		return
	}

}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
