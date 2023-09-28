package protocolmetadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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

type protocolMetadataLifetimeSettingsResourceModel struct {
	Id            types.String `tfsdk:"id"`
	CacheDuration types.Int64  `tfsdk:"cache_duration"`
	ReloadDelay   types.Int64  `tfsdk:"reload_delay"`
}

// GetSchema defines the schema for the resource.
func (r *protocolMetadataLifetimeSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a ProtocolMetadataLifetimeSettings.",
		Attributes: map[string]schema.Attribute{
			"cache_duration": schema.Int64Attribute{
				Description: "This field adjusts the validity of your metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown()},
			},
			"reload_delay": schema.Int64Attribute{
				Description: "This field adjusts the frequency of automatic reloading of SAML metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema)
	resp.Schema = schema
}

func addOptionalProtocolMetadataLifetimeSettingsFields(ctx context.Context, addRequest *client.MetadataLifetimeSettings, plan protocolMetadataLifetimeSettingsResourceModel) error {

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

func readProtocolMetadataLifetimeSettingsResponse(ctx context.Context, r *client.MetadataLifetimeSettings, state *protocolMetadataLifetimeSettingsResourceModel, expectedValues *protocolMetadataLifetimeSettingsResourceModel) {
	//TODO placeholder?
	state.Id = types.StringValue("id")
	state.CacheDuration = types.Int64Value(*r.CacheDuration)
	state.ReloadDelay = types.Int64Value(*r.ReloadDelay)
}

func (r *protocolMetadataLifetimeSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan protocolMetadataLifetimeSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Protocol Metadata Lifetime Settings", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Protocol Metadata Lifetime Settings: %s", requestErr.Error())
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Protocol Metadata Lifetime Settings", err, httpResp)
		return
	}
	_, responseErr := protocolMetadataLifetimeSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Protocol Metadata Lifetime Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsResourceModel

	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *protocolMetadataLifetimeSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state protocolMetadataLifetimeSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadProtocolMetadataLifetimeSettings, httpResp, err := r.apiClient.ProtocolMetadataApi.GetLifetimeSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Protocol Metadata Lifetime Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Protocol Metadata Lifetime Settings", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadProtocolMetadataLifetimeSettings.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Protocol Metadata Lifetime Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	readProtocolMetadataLifetimeSettingsResponse(ctx, apiReadProtocolMetadataLifetimeSettings, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *protocolMetadataLifetimeSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan protocolMetadataLifetimeSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Protocol Metadata Lifetime Settings", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Protocol Metadata Lifetime Settings: %s", requestErr.Error())
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Protocol Metadata Lifetime Settings", err, httpResp)
		return
	}
	_, responseErr := protocolMetadataLifetimeSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Protocol Metadata Lifetime Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsResourceModel

	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *protocolMetadataLifetimeSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *protocolMetadataLifetimeSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
