package oauthauthserversettingsscopesexclusivescope

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
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsScopesExclusiveScopeResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsScopesExclusiveScopeResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsScopesExclusiveScopeResource{}
)

// OauthAuthServerSettingsScopesExclusiveScopeResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsScopesExclusiveScopeResource() resource.Resource {
	return &oauthAuthServerSettingsScopesExclusiveScopeResource{}
}

// oauthAuthServerSettingsScopesExclusiveScopeResource is the resource implementation.
type oauthAuthServerSettingsScopesExclusiveScopeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an exclusive scope in the authorization server settings.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the scope.",
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

	id.ToSchema(&schema)
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) ValidateConfig(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model oauthAuthServerSettingsScopesExclusiveScopeModel
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

func addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx context.Context, addRequest *client.ScopeEntry, plan oauthAuthServerSettingsScopesExclusiveScopeModel) error {

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
func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_exclusive_scope"
}

func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsScopesExclusiveScopeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAuthServerSettingsScopesExclusiveScopes := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx, createOauthAuthServerSettingsScopesExclusiveScopes, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Auth Server Settings Scopes Exclusive Scope", err.Error())
		return
	}

	apiCreateOauthAuthServerSettingsScopesExclusiveScopes := r.apiClient.OauthAuthServerSettingsAPI.AddExclusiveScope(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettingsScopesExclusiveScopes = apiCreateOauthAuthServerSettingsScopesExclusiveScopes.Body(*createOauthAuthServerSettingsScopesExclusiveScopes)
	oauthAuthServerSettingsScopesExclusiveScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.AddExclusiveScopeExecute(apiCreateOauthAuthServerSettingsScopesExclusiveScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAuthServerSettingsScopesExclusiveScopeModel

	readOauthAuthServerSettingsScopesExclusiveScopeResponse(ctx, oauthAuthServerSettingsScopesExclusiveScopesResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAuthServerSettingsScopesExclusiveScopeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettingsScopesExclusiveScopes, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetExclusiveScope(config.AuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesExclusiveScopeResponse(ctx, apiReadOauthAuthServerSettingsScopesExclusiveScopes, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthAuthServerSettingsScopesExclusiveScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAuthServerSettingsScopesExclusiveScopeModel
	req.State.Get(ctx, &state)
	updateOauthAuthServerSettingsScopesExclusiveScopes := r.apiClient.OauthAuthServerSettingsAPI.UpdateExclusiveScope(config.AuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Auth Server Settings Scopes Exclusive Scope", err.Error())
		return
	}

	updateOauthAuthServerSettingsScopesExclusiveScopes = updateOauthAuthServerSettingsScopesExclusiveScopes.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsScopesExclusiveScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateExclusiveScopeExecute(updateOauthAuthServerSettingsScopesExclusiveScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
		return
	}

	// Read the response
	readOauthAuthServerSettingsScopesExclusiveScopeResponse(ctx, updateOauthAuthServerSettingsScopesExclusiveScopesResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthAuthServerSettingsScopesExclusiveScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.RemoveExclusiveScope(config.AuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
	}
}

func (r *oauthAuthServerSettingsScopesExclusiveScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
