// Copyright © 2025 Ping Identity Corporation

package protocolmetadatalifetimesettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &protocolMetadataLifetimeSettingsResource{}
	_ resource.ResourceWithConfigure   = &protocolMetadataLifetimeSettingsResource{}
	_ resource.ResourceWithImportState = &protocolMetadataLifetimeSettingsResource{}
)

// ProtocolMetadataLifetimeSettingsResource is a helper function to simplify the provider implementation.
func ProtocolMetadataLifetimeSettingsResource() resource.Resource {
	return &protocolMetadataLifetimeSettingsResource{}
}

// protocolMetadataLifetimeSettingsResource is the resource implementation.
type protocolMetadataLifetimeSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *protocolMetadataLifetimeSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings for the metadata cache duration and reload delay for protocol metadata.",
		Attributes: map[string]schema.Attribute{
			"cache_duration": schema.Int64Attribute{
				Description: "This field adjusts the validity of your metadata in minutes. The default value is `1440` (1 day).",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(1440),
			},
			"reload_delay": schema.Int64Attribute{
				Description: "This field adjusts the frequency of automatic reloading of SAML metadata in minutes. The default value is `1440` (1 day).",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(1440),
			},
		},
	}
	resp.Schema = schema
}

func addOptionalProtocolMetadataLifetimeSettingsFields(ctx context.Context, addRequest *client.MetadataLifetimeSettings, plan protocolMetadataLifetimeSettingsModel) error {

	if internaltypes.IsDefined(plan.CacheDuration) {
		addRequest.CacheDuration = plan.CacheDuration.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ReloadDelay) {
		addRequest.ReloadDelay = plan.ReloadDelay.ValueInt64Pointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *protocolMetadataLifetimeSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protocol_metadata_lifetime_settings"
}

func (r *protocolMetadataLifetimeSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (m *protocolMetadataLifetimeSettingsModel) buildDefaultClientStruct() *client.MetadataLifetimeSettings {
	return &client.MetadataLifetimeSettings{
		CacheDuration: utils.Pointer(int64(1440)),
		ReloadDelay:   utils.Pointer(int64(1440)),
	}
}

func (r *protocolMetadataLifetimeSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan protocolMetadataLifetimeSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Protocol Metadata Lifetime Settings: "+err.Error())
		return
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Protocol Metadata Lifetime Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsModel
	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *protocolMetadataLifetimeSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state protocolMetadataLifetimeSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadProtocolMetadataLifetimeSettings, httpResp, err := r.apiClient.ProtocolMetadataAPI.GetLifetimeSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Protocol Metadata Lifetime Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Protocol Metadata Lifetime Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readProtocolMetadataLifetimeSettingsResponse(ctx, apiReadProtocolMetadataLifetimeSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *protocolMetadataLifetimeSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan protocolMetadataLifetimeSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Protocol Metadata Lifetime Settings: "+err.Error())
		return
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Protocol Metadata Lifetime Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsModel
	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *protocolMetadataLifetimeSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model protocolMetadataLifetimeSettingsModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.ProtocolMetadataAPI.UpdateLifetimeSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the protocol metadata lifetime settings", err, httpResp)
	}
}

func (r *protocolMetadataLifetimeSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState protocolMetadataLifetimeSettingsModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
