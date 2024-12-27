package oauthtokenexchangegeneratorsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthTokenExchangeGeneratorSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthTokenExchangeGeneratorSettingsResource{}
	_ resource.ResourceWithImportState = &oauthTokenExchangeGeneratorSettingsResource{}
)

// OauthTokenExchangeGeneratorSettingsResource is a helper function to simplify the provider implementation.
func OauthTokenExchangeGeneratorSettingsResource() resource.Resource {
	return &oauthTokenExchangeGeneratorSettingsResource{}
}

// oauthTokenExchangeGeneratorSettingsResource is the resource implementation.
type oauthTokenExchangeGeneratorSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthTokenExchangeGeneratorSettingsResourceModel struct {
	DefaultGeneratorGroupRef types.Object `tfsdk:"default_generator_group_ref"`
}

// GetSchema defines the schema for the resource.
func (r *oauthTokenExchangeGeneratorSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Oauth Token Exchange Generator Settings",
		Attributes: map[string]schema.Attribute{
			"default_generator_group_ref": resourcelink.CompleteSingleNestedAttribute(
				false,
				false,
				true,
				"Reference to the default Token Exchange Generator group, if one is defined.",
			),
		},
	}

	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *oauthTokenExchangeGeneratorSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token_exchange_generator_settings"
}

func (r *oauthTokenExchangeGeneratorSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthTokenExchangeGeneratorSettingsResponse(ctx context.Context, r *client.TokenExchangeGeneratorSettings, state *oauthTokenExchangeGeneratorSettingsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.DefaultGeneratorGroupRef, diags = resourcelink.ToState(ctx, r.DefaultGeneratorGroupRef)

	// make sure all object type building appends diags
	return diags
}

func (r *oauthTokenExchangeGeneratorSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var plan oauthTokenExchangeGeneratorSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthTokenExchangeGeneratorSettings := client.NewTokenExchangeGeneratorSettings()
	createOauthTokenExchangeGeneratorSettings.DefaultGeneratorGroupRef, err = resourcelink.ClientStruct(plan.DefaultGeneratorGroupRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add default_generator_group_ref to add request for OAuth Token Exchange Generator Settings: "+err.Error())
		return
	}

	apiCreateOauthTokenExchangeGeneratorSettings := r.apiClient.OauthTokenExchangeGeneratorAPI.UpdateOauthTokenExchangeSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthTokenExchangeGeneratorSettings = apiCreateOauthTokenExchangeGeneratorSettings.Body(*createOauthTokenExchangeGeneratorSettings)
	oauthTokenExchangeGeneratorSettingsResponse, httpResp, err := r.apiClient.OauthTokenExchangeGeneratorAPI.UpdateOauthTokenExchangeSettingsExecute(apiCreateOauthTokenExchangeGeneratorSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Token Exchange Generator Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthTokenExchangeGeneratorSettingsResourceModel

	diags = readOauthTokenExchangeGeneratorSettingsResponse(ctx, oauthTokenExchangeGeneratorSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthTokenExchangeGeneratorSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthTokenExchangeGeneratorSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthTokenExchangeGeneratorSettings, httpResp, err := r.apiClient.OauthTokenExchangeGeneratorAPI.GetOauthTokenExchangeSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Token Exchange Generator Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Token Exchange Generator Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOauthTokenExchangeGeneratorSettingsResponse(ctx, apiReadOauthTokenExchangeGeneratorSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthTokenExchangeGeneratorSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var plan oauthTokenExchangeGeneratorSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createUpdateRequest := client.NewTokenExchangeGeneratorSettings()
	createUpdateRequest.DefaultGeneratorGroupRef, err = resourcelink.ClientStruct(plan.DefaultGeneratorGroupRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add default_generator_group_ref to add request for OAuth Token Exchange Generator Settings: "+err.Error())
		return
	}

	updateOauthTokenExchangeGeneratorSettings := r.apiClient.OauthTokenExchangeGeneratorAPI.UpdateOauthTokenExchangeSettings(config.AuthContext(ctx, r.providerConfig))
	updateOauthTokenExchangeGeneratorSettings = updateOauthTokenExchangeGeneratorSettings.Body(*createUpdateRequest)
	updateOauthTokenExchangeGeneratorSettingsResponse, httpResp, err := r.apiClient.OauthTokenExchangeGeneratorAPI.UpdateOauthTokenExchangeSettingsExecute(updateOauthTokenExchangeGeneratorSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Token Exchange Generator Settings", err, httpResp)
		return
	}

	// Read the response
	var state oauthTokenExchangeGeneratorSettingsResourceModel
	diags = readOauthTokenExchangeGeneratorSettingsResponse(ctx, updateOauthTokenExchangeGeneratorSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthTokenExchangeGeneratorSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	providererror.WarnConfigurationCannotBeReset("pingfederate_oauth_token_exchange_generator_settings", &resp.Diagnostics)
}

func (r *oauthTokenExchangeGeneratorSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState oauthTokenExchangeGeneratorSettingsResourceModel
	emptyState.DefaultGeneratorGroupRef = types.ObjectNull(resourcelink.AttrType())
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
