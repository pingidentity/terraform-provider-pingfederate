package authenticationpoliciessettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
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
			"id": schema.StringAttribute{
				Description: "The ID of the resource.",
				Computed:    true,
				Optional:    false,
			},
			"enable_idp_authn_selection": schema.BoolAttribute{
				Description: "Enable IdP authentication policies.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_sp_authn_selection": schema.BoolAttribute{
				Description: "Enable SP authentication policies.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	// // Set attributes in string list
	// if setOptionalToComputed {
	// 	config.SetAllAttributesToOptionalAndComputed(&schema, []string{"FIX_ME"})
	// }
	// config.AddCommonSchema(&schema, false)
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalAuthenticationPoliciesSettingsFields(ctx context.Context, addRequest *client.AuthenticationPoliciesSettings, plan authenticationPoliciesSettingsModel) error {

	if internaltypes.IsDefined(plan.EnableIdpAuthnSelection) {
		addRequest.EnableIdpAuthnSelection = plan.EnableIdpAuthnSelection.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.EnableSpAuthnSelection) {
		addRequest.EnableSpAuthnSelection = plan.EnableSpAuthnSelection.ValueBoolPointer()
	}

	return nil

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
	err := addOptionalAuthenticationPoliciesSettingsFields(ctx, createAuthenticationPoliciesSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for authentication policies settings", err.Error())
		return
	}
	requestJson, err := createAuthenticationPoliciesSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateAuthenticationPoliciesSettings := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPoliciesSettings = apiCreateAuthenticationPoliciesSettings.Body(*createAuthenticationPoliciesSettings)
	authenticationPoliciesSettingsResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettingsExecute(apiCreateAuthenticationPoliciesSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication policies settings", err, httpResp)
		return
	}
	responseJson, err := authenticationPoliciesSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state authenticationPoliciesSettingsModel

	readAuthenticationPoliciesSettings(ctx, authenticationPoliciesSettingsResponse, &state, nil)
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
	apiReadAuthenticationPoliciesSettings, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetAuthenticationPolicySettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the authentication policies settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the authentication policies settings", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadAuthenticationPoliciesSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readAuthenticationPoliciesSettings(ctx, apiReadAuthenticationPoliciesSettings, &state, id)

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
	updateAuthenticationPoliciesSettings := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthenticationPoliciesSettings()
	err := addOptionalAuthenticationPoliciesSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for authentication policies settings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateAuthenticationPoliciesSettings = updateAuthenticationPoliciesSettings.Body(*createUpdateRequest)
	updateAuthenticationPoliciesSettingsResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateAuthenticationPolicySettingsExecute(updateAuthenticationPoliciesSettings)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication policies settings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAuthenticationPoliciesSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read the response
	readAuthenticationPoliciesSettings(ctx, updateAuthenticationPoliciesSettingsResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *authenticationPoliciesSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *authenticationPoliciesSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
