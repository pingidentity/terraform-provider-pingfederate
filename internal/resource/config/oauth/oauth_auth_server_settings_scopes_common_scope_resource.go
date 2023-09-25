package oauth

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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
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
	resp.Schema = schema.Schema{
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
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan oauthAuthServerSettingsScopesCommonScopesResourceModel
	req.Plan.Get(ctx, &plan)
	if plan.Dynamic.ValueBool() && (plan.Name.ValueString() != "" || !plan.Name.IsNull()) {
		{
			containsAsteriskPrefix := strings.Index(plan.Name.ValueString(), "*")
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
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_common_scope"
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Oauth Auth Server Settings Scopes Common Scope", err.Error())
		return
	}
	_, requestErr := createOauthAuthServerSettingsScopesCommonScopes.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Oauth Auth Server Settings Scopes Common Scope: %s", requestErr.Error())
	}

	apiCreateOauthAuthServerSettingsScopesCommonScopes := r.apiClient.OauthAuthServerSettingsApi.AddCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettingsScopesCommonScopes = apiCreateOauthAuthServerSettingsScopesCommonScopes.Body(*createOauthAuthServerSettingsScopesCommonScopes)
	oauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsApi.AddCommonScopeExecute(apiCreateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Oauth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}
	_, responseErr := oauthAuthServerSettingsScopesCommonScopesResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Oauth Auth Server Settings Scopes Common Scope: %s", responseErr.Error())
	}

	// Read the response into the state
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel

	readOauthAuthServerSettingsScopesCommonScopesResponse(ctx, oauthAuthServerSettingsScopesCommonScopesResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettingsScopesCommonScopes, httpResp, err := r.apiClient.OauthAuthServerSettingsApi.GetCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an Oauth AuthServerSettingsScopesCommonScopes", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OauthAuthServerSettingsScopesCommonScopes", err, httpResp)
		}
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
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	updateOauthAuthServerSettingsScopesCommonScopes := r.apiClient.OauthAuthServerSettingsApi.UpdateCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewScopeEntry(plan.Id.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Oauth Auth Server Settings Scopes Common Scope", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Oauth Auth Server Settings Scopes Common Scope: %s", requestErr.Error())
	}
	updateOauthAuthServerSettingsScopesCommonScopes = updateOauthAuthServerSettingsScopesCommonScopes.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsApi.UpdateCommonScopeExecute(updateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Oauth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateOauthAuthServerSettingsScopesCommonScopesResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Oauth Auth Server Settings Scopes Common Scope: %s", responseErr.Error())
	}
	// Read the response
	readOauthAuthServerSettingsScopesCommonScopesResponse(ctx, updateOauthAuthServerSettingsScopesCommonScopesResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthAuthServerSettingsScopesCommonScopesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAuthServerSettingsApi.RemoveCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Oauth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}
}

func (r *oauthAuthServerSettingsScopesCommonScopesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
