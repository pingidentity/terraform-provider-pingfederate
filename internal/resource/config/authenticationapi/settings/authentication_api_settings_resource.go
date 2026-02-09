// Copyright © 2026 Ping Identity Corporation

package authenticationapisettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationApiSettingsResource{}
	_ resource.ResourceWithConfigure   = &authenticationApiSettingsResource{}
	_ resource.ResourceWithImportState = &authenticationApiSettingsResource{}
)

// AuthenticationApiSettingsResource is a helper function to simplify the provider implementation.
func AuthenticationApiSettingsResource() resource.Resource {
	return &authenticationApiSettingsResource{}
}

// authenticationApiSettingsResource is the resource implementation.
type authenticationApiSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *authenticationApiSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the Authentication API application settings.",
		Attributes: map[string]schema.Attribute{
			"api_enabled": schema.BoolAttribute{
				Description: "Enable Authentication API. The default is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_api_descriptions": schema.BoolAttribute{
				Description: "Enable API descriptions. The default is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"default_application_ref": schema.SingleNestedAttribute{
				Description: "Application for non authentication policy use cases",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"restrict_access_to_redirectless_mode": schema.BoolAttribute{
				Description: "Enable restrict access to redirectless mode. The default is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_request_context": schema.BoolAttribute{
				Description: "Includes request context in API responses. The default is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	resp.Schema = schema
}

func addAuthenticationApiSettingsFields(ctx context.Context, addRequest *client.AuthnApiSettings, plan authenticationApiSettingsModel) error {
	if internaltypes.IsDefined(plan.ApiEnabled) {
		addRequest.ApiEnabled = plan.ApiEnabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EnableApiDescriptions) {
		addRequest.EnableApiDescriptions = plan.EnableApiDescriptions.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RestrictAccessToRedirectlessMode) {
		addRequest.RestrictAccessToRedirectlessMode = plan.RestrictAccessToRedirectlessMode.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestContext) {
		addRequest.IncludeRequestContext = plan.IncludeRequestContext.ValueBoolPointer()
	}
	if internaltypes.IsNonEmptyObj(plan.DefaultApplicationRef) {
		addRequestNewLinkObj, err := resourcelink.ClientStruct(plan.DefaultApplicationRef)
		if err != nil {
			return err
		}
		addRequest.DefaultApplicationRef = addRequestNewLinkObj
	}
	return nil

}

// Metadata returns the resource type name.
func (r *authenticationApiSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_api_settings"
}

func (r *authenticationApiSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (m *authenticationApiSettingsModel) buildDefaultClientStruct() *client.AuthnApiSettings {
	return &client.AuthnApiSettings{
		ApiEnabled:                       utils.Pointer(false),
		EnableApiDescriptions:            utils.Pointer(false),
		RestrictAccessToRedirectlessMode: utils.Pointer(false),
		IncludeRequestContext:            utils.Pointer(false),
	}
}

func (r *authenticationApiSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan authenticationApiSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get the current state to see how any attributes are changing
	updateAuthenticationApiSettings := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthnApiSettings()
	err := addAuthenticationApiSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the authentication API settings: "+err.Error())
		return
	}

	updateAuthenticationApiSettings = updateAuthenticationApiSettings.Body(*createUpdateRequest)
	updateAuthenticationApiSettingsResponse, httpResp, err := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettingsExecute(updateAuthenticationApiSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication API settings", err, httpResp)
		return
	}

	// Read the response
	var state authenticationApiSettingsModel
	diags = readAuthenticationApiSettingsResponse(ctx, updateAuthenticationApiSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationApiSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationApiSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationApiSettings, httpResp, err := r.apiClient.AuthenticationApiAPI.GetAuthenticationApiSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Authentication API Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the authentication API settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readAuthenticationApiSettingsResponse(ctx, apiReadAuthenticationApiSettings, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationApiSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan authenticationApiSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get the current state to see how any attributes are changing
	updateAuthenticationApiSettings := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthnApiSettings()
	err := addAuthenticationApiSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the authentication API settings: "+err.Error())
		return
	}

	updateAuthenticationApiSettings = updateAuthenticationApiSettings.Body(*createUpdateRequest)
	updateAuthenticationApiSettingsResponse, httpResp, err := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettingsExecute(updateAuthenticationApiSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the authentication API settings", err, httpResp)
		return
	}

	// Read the response
	var state authenticationApiSettingsModel
	diags = readAuthenticationApiSettingsResponse(ctx, updateAuthenticationApiSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
func (r *authenticationApiSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model authenticationApiSettingsModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.AuthenticationApiAPI.UpdateAuthenticationApiSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the authentication API settings", err, httpResp)
	}
}

func (r *authenticationApiSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState authenticationApiSettingsModel
	emptyState.DefaultApplicationRef = types.ObjectNull(resourcelink.AttrType())
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
