package oauthcibaserverpolicysettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthCibaServerPolicySettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthCibaServerPolicySettingsResource{}
	_ resource.ResourceWithImportState = &oauthCibaServerPolicySettingsResource{}
)

// OauthCibaServerPolicySettingsResource is a helper function to simplify the provider implementation.
func OauthCibaServerPolicySettingsResource() resource.Resource {
	return &oauthCibaServerPolicySettingsResource{}
}

// oauthCibaServerPolicySettingsResource is the resource implementation.
type oauthCibaServerPolicySettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthCibaServerPolicySettingsResourceModel struct {
	DefaultRequestPolicyRef types.Object `tfsdk:"default_request_policy_ref"`
}

// GetSchema defines the schema for the resource.
func (r *oauthCibaServerPolicySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages OAuth CIBA Server Policy Settings",
		Attributes: map[string]schema.Attribute{
			"default_request_policy_ref": resourcelink.CompleteSingleNestedAttribute(
				false,
				false,
				true,
				"Reference to the default request policy, if one is defined.",
			),
		},
	}
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *oauthCibaServerPolicySettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_ciba_server_policy_settings"
}

func (r *oauthCibaServerPolicySettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthCibaServerPolicySettingsResponse(ctx context.Context, r *client.CibaServerPolicySettings, state *oauthCibaServerPolicySettingsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.DefaultRequestPolicyRef, diags = resourcelink.ToState(ctx, r.DefaultRequestPolicyRef)

	// make sure all object type building appends diags
	return diags
}

func (r *oauthCibaServerPolicySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var plan oauthCibaServerPolicySettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthCibaServerPolicySettings := client.NewCibaServerPolicySettings()
	createOauthCibaServerPolicySettings.DefaultRequestPolicyRef, err = resourcelink.ClientStruct(plan.DefaultRequestPolicyRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add default_request_policy_ref to add request for OAuth CIBA Server Policy Settings: "+err.Error())
		return
	}

	apiCreateOauthCibaServerPolicySettings := r.apiClient.OauthCibaServerPolicyAPI.UpdateCibaServerPolicySettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthCibaServerPolicySettings = apiCreateOauthCibaServerPolicySettings.Body(*createOauthCibaServerPolicySettings)
	oauthCibaServerPolicySettingsResponse, httpResp, err := r.apiClient.OauthCibaServerPolicyAPI.UpdateCibaServerPolicySettingsExecute(apiCreateOauthCibaServerPolicySettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth CIBA Server Policy Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthCibaServerPolicySettingsResourceModel

	diags = readOauthCibaServerPolicySettingsResponse(ctx, oauthCibaServerPolicySettingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthCibaServerPolicySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthCibaServerPolicySettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthCibaServerPolicySettings, httpResp, err := r.apiClient.OauthCibaServerPolicyAPI.GetCibaServerPolicySettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth CIBA Server Policy Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth CIBA Server Policy Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOauthCibaServerPolicySettingsResponse(ctx, apiReadOauthCibaServerPolicySettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthCibaServerPolicySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var plan oauthCibaServerPolicySettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOauthCibaServerPolicySettings := r.apiClient.OauthCibaServerPolicyAPI.UpdateCibaServerPolicySettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewCibaServerPolicySettings()
	createUpdateRequest.DefaultRequestPolicyRef, err = resourcelink.ClientStruct(plan.DefaultRequestPolicyRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add default_request_policy_ref to add request for OAuth CIBA Server Policy Settings: "+err.Error())
		return
	}

	updateOauthCibaServerPolicySettings = updateOauthCibaServerPolicySettings.Body(*createUpdateRequest)
	updateOauthCibaServerPolicySettingsResponse, httpResp, err := r.apiClient.OauthCibaServerPolicyAPI.UpdateCibaServerPolicySettingsExecute(updateOauthCibaServerPolicySettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth CIBA Server Policy Settings", err, httpResp)
		return
	}

	// Read the response
	var state oauthCibaServerPolicySettingsResourceModel
	diags = readOauthCibaServerPolicySettingsResponse(ctx, updateOauthCibaServerPolicySettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthCibaServerPolicySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	providererror.WarnConfigurationCannotBeReset("pingfederate_oauth_ciba_server_policy_settings", &resp.Diagnostics)
}

func (r *oauthCibaServerPolicySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState oauthCibaServerPolicySettingsResourceModel
	emptyState.DefaultRequestPolicyRef = types.ObjectNull(resourcelink.AttrType())
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
