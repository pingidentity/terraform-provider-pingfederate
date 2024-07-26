package authenticationpoliciessettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPoliciesSettingsResource{}
	_ resource.ResourceWithConfigure   = &authenticationPoliciesSettingsResource{}
	_ resource.ResourceWithImportState = &authenticationPoliciesSettingsResource{}
)

// AuthenticationPoliciesSettingsResource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesSettingsResource() resource.Resource {
	return &authenticationPoliciesSettingsResource{}
}

// authenticationPoliciesSettingsResource is the resource implementation.
type authenticationPoliciesSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *authenticationPoliciesSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Authentication Policies Settings",
		Attributes: map[string]schema.Attribute{
			"enable_idp_authn_selection": schema.BoolAttribute{
				Description: "Enable IdP authentication policies. Default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_sp_authn_selection": schema.BoolAttribute{
				Description: "Enable SP authentication policies. Default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func addOptionalAuthenticationPoliciesSettingsFields(addRequest *client.AuthenticationPoliciesSettings, plan authenticationPoliciesSettingsModel) {
	addRequest.EnableIdpAuthnSelection = plan.EnableIdpAuthnSelection.ValueBoolPointer()
	addRequest.EnableSpAuthnSelection = plan.EnableSpAuthnSelection.ValueBoolPointer()
}

// Metadata returns the resource type name.
func (r *authenticationPoliciesSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_settings"
}

func (r *authenticationPoliciesSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *authenticationPoliciesSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPoliciesSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationPoliciesSettings := client.NewAuthenticationPoliciesSettings()
	addOptionalAuthenticationPoliciesSettingsFields(createAuthenticationPoliciesSettings, plan)

	apiCreateAuthenticationPoliciesSettings := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPoliciesSettings = apiCreateAuthenticationPoliciesSettings.Body(*createAuthenticationPoliciesSettings)
	authenticationPoliciesSettingsResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettingsExecute(apiCreateAuthenticationPoliciesSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication policies settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state authenticationPoliciesSettingsModel

	readAuthenticationPoliciesSettingsResponse(authenticationPoliciesSettingsResponse, &state, nil)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPoliciesSettings, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetAuthenticationPolicySettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the authentication policies settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the authentication policies settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readAuthenticationPoliciesSettingsResponse(apiReadAuthenticationPoliciesSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan authenticationPoliciesSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authenticationPoliciesSettingsModel
	req.State.Get(ctx, &state)
	updateAuthenticationPoliciesSettings := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthenticationPoliciesSettings()
	addOptionalAuthenticationPoliciesSettingsFields(createUpdateRequest, plan)

	updateAuthenticationPoliciesSettings = updateAuthenticationPoliciesSettings.Body(*createUpdateRequest)
	updateAuthenticationPoliciesSettingsResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettingsExecute(updateAuthenticationPoliciesSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication policies settings", err, httpResp)
		return
	}

	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read the response
	readAuthenticationPoliciesSettingsResponse(updateAuthenticationPoliciesSettingsResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *authenticationPoliciesSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *authenticationPoliciesSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
