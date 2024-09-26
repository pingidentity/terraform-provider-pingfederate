package oauthopenidconnectsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &openidConnectSettingsResource{}
	_ resource.ResourceWithConfigure   = &openidConnectSettingsResource{}
	_ resource.ResourceWithImportState = &openidConnectSettingsResource{}

	openidConnectSettingsAttrTypes = map[string]attr.Type{
		"track_user_sessions_for_logout": types.BoolType,
		"revoke_user_session_on_logout":  types.BoolType,
		"session_revocation_lifetime":    types.Int64Type,
	}
)

// OpenidConnectSettingsResource is a helper function to simplify the provider implementation.
func OpenidConnectSettingsResource() resource.Resource {
	return &openidConnectSettingsResource{}
}

// openidConnectSettingsResource is the resource implementation.
type openidConnectSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type openidConnectSettingsResourceModel struct {
	DefaultPolicyRef types.Object `tfsdk:"default_policy_ref"`
}

// GetSchema defines the schema for the resource.
func (r *openidConnectSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages OpenID Connect configuration settings",
		Attributes: map[string]schema.Attribute{
			"default_policy_ref": schema.SingleNestedAttribute{
				Description: "Reference to the default policy.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			//TODO verify that missing session_settings from this reosurce doesn't conflict with the session settings resource
		},
	}
	resp.Schema = schema
}

func addOptionalOpenidConnectSettingsFields(ctx context.Context, addRequest *client.OpenIdConnectSettings, plan openidConnectSettingsResourceModel) error {
	var err error

	if internaltypes.IsDefined(plan.DefaultPolicyRef) {
		addRequest.DefaultPolicyRef, err = resourcelink.ClientStruct(plan.DefaultPolicyRef)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *openidConnectSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openid_connect_settings"
}

func (r *openidConnectSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOpenidConnectSettingsResponse(ctx context.Context, r *client.OpenIdConnectSettings, state *openidConnectSettingsResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.DefaultPolicyRef, respDiags = resourcelink.ToState(ctx, r.DefaultPolicyRef)
	diags = append(diags, respDiags...)

	// make sure all object type building appends diags
	return diags
}

func (r *openidConnectSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan openidConnectSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpenidConnectSettings := client.NewOpenIdConnectSettings()
	err := addOptionalOpenidConnectSettingsFields(ctx, createOpenidConnectSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OpenID Connect settings: "+err.Error())
		return
	}
	apiCreateOpenidConnectSettings := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateOpenidConnectSettings = apiCreateOpenidConnectSettings.Body(*createOpenidConnectSettings)
	openidConnectSettingsResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettingsExecute(apiCreateOpenidConnectSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OpenID Connect settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state openidConnectSettingsResourceModel

	diags = readOpenidConnectSettingsResponse(ctx, openidConnectSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *openidConnectSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state openidConnectSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOpenidConnectSettings, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.GetOIDCSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OpenID Connect Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OpenID Connect settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOpenidConnectSettingsResponse(ctx, apiReadOpenidConnectSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *openidConnectSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan openidConnectSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpenidConnectSettings := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewOpenIdConnectSettings()
	err := addOptionalOpenidConnectSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OpenID Connect settings: "+err.Error())
		return
	}

	updateOpenidConnectSettings = updateOpenidConnectSettings.Body(*createUpdateRequest)
	updateOpenidConnectSettingsResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettingsExecute(updateOpenidConnectSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OpenID Connect settings", err, httpResp)
		return
	}

	// Read the response
	var state openidConnectSettingsResourceModel
	diags = readOpenidConnectSettingsResponse(ctx, updateOpenidConnectSettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *openidConnectSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	providererror.WarnConfigurationCannotBeReset("pingfederate_openid_connect_settings", &resp.Diagnostics)
}

func (r *openidConnectSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	//TODO this needs fixed, need to build an empty struct here
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
