package oauthauthserversettingsscopescommonscope

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsScopesCommonScopeResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsScopesCommonScopeResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsScopesCommonScopeResource{}
)

// OauthAuthServerSettingsScopesCommonScopeResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsScopesCommonScopeResource() resource.Resource {
	return &oauthAuthServerSettingsScopesCommonScopeResource{}
}

// oauthAuthServerSettingsScopesCommonScopeResource is the resource implementation.
type oauthAuthServerSettingsScopesCommonScopeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsScopesCommonScopeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description:        "Manages a common scope in the authorization server settings.",
		DeprecationMessage: "This resource is deprecated and will be removed in a future release. Use the `pingfederate_oauth_auth_server_settings` resource instead.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the scope. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the scope that appears when the user is prompted for authorization.",
				Required:    true,
			},
			"dynamic": schema.BoolAttribute{
				Description: "True if the scope is dynamic. (Defaults to false)",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsScopesCommonScopeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *oauthAuthServerSettingsScopesCommonScopeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	if plan.Dynamic.ValueBool() && (plan.Name.ValueString() != "" || !plan.Name.IsNull()) {
		{
			containsAsteriskPrefix := strings.Index(plan.Name.ValueString(), "*")
			if containsAsteriskPrefix == -1 {
				resp.Diagnostics.AddAttributeError(
					path.Root("name"),
					providererror.InvalidAttributeConfiguration,
					"The name must include a \"*\" when set to dynamic")
			}
		}
	}
}

func addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx context.Context, addRequest *client.ScopeEntry, plan oauthAuthServerSettingsScopesCommonScopeModel) error {

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
func (r *oauthAuthServerSettingsScopesCommonScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_common_scope"
}

func (r *oauthAuthServerSettingsScopesCommonScopeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *oauthAuthServerSettingsScopesCommonScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsScopesCommonScopeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAuthServerSettingsScopesCommonScopes := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx, createOauthAuthServerSettingsScopesCommonScopes, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OAuth Auth Server Settings Scopes Common Scope: "+err.Error())
		return
	}

	apiCreateOauthAuthServerSettingsScopesCommonScopes := r.apiClient.OauthAuthServerSettingsAPI.AddCommonScope(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettingsScopesCommonScopes = apiCreateOauthAuthServerSettingsScopesCommonScopes.Body(*createOauthAuthServerSettingsScopesCommonScopes)
	oauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.AddCommonScopeExecute(apiCreateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAuthServerSettingsScopesCommonScopeModel

	readOauthAuthServerSettingsScopesCommonScopeResponse(ctx, oauthAuthServerSettingsScopesCommonScopesResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAuthServerSettingsScopesCommonScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAuthServerSettingsScopesCommonScopeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettingsScopesCommonScopes, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetCommonScope(config.AuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Auth Server Settings Common Scopes", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Auth Server Settings Scopes Common Scope", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesCommonScopeResponse(ctx, apiReadOauthAuthServerSettingsScopesCommonScopes, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthAuthServerSettingsScopesCommonScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAuthServerSettingsScopesCommonScopeModel
	req.State.Get(ctx, &state)
	updateOauthAuthServerSettingsScopesCommonScopes := r.apiClient.OauthAuthServerSettingsAPI.UpdateCommonScope(config.AuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesCommonScopesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OAuth Auth Server Settings Scopes Common Scope: "+err.Error())
		return
	}

	updateOauthAuthServerSettingsScopesCommonScopes = updateOauthAuthServerSettingsScopesCommonScopes.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsScopesCommonScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateCommonScopeExecute(updateOauthAuthServerSettingsScopesCommonScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}

	// Read the response
	readOauthAuthServerSettingsScopesCommonScopeResponse(ctx, updateOauthAuthServerSettingsScopesCommonScopesResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthAuthServerSettingsScopesCommonScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthAuthServerSettingsScopesCommonScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.RemoveCommonScope(config.AuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a OAuth Auth Server Settings Scopes Common Scope", err, httpResp)
	}
}

func (r *oauthAuthServerSettingsScopesCommonScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
