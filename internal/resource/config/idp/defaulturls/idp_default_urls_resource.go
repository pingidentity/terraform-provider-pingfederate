package idpdefaulturls

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
	_ resource.Resource                = &idpDefaultUrlsResource{}
	_ resource.ResourceWithConfigure   = &idpDefaultUrlsResource{}
	_ resource.ResourceWithImportState = &idpDefaultUrlsResource{}
)

// IdpDefaultUrlsResource is a helper function to simplify the provider implementation.
func IdpDefaultUrlsResource() resource.Resource {
	return &idpDefaultUrlsResource{}
}

// idpDefaultUrlsResource is the resource implementation.
type idpDefaultUrlsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *idpDefaultUrlsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the IdP default URL settings",
		Attributes: map[string]schema.Attribute{
			"confirm_idp_slo": schema.BoolAttribute{
				Description: "Prompt user to confirm Single Logout (SLO).",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idp_error_msg": schema.StringAttribute{
				Description: "Provide the error text displayed in a user's browser when an SSO operation fails.",
				Required:    true,
			},
			"idp_slo_success_url": schema.StringAttribute{
				Description: "Provide the default URL you would like to send the user to when Single Logout has succeeded.",
				Optional:    true,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalIdpDefaultUrlsFields(ctx context.Context, addRequest *client.IdpDefaultUrl, plan idpDefaultUrlsModel) error {
	if internaltypes.IsDefined(plan.ConfirmIdpSlo) {
		addRequest.ConfirmIdpSlo = plan.ConfirmIdpSlo.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IdpSloSuccessUrl) {
		addRequest.IdpSloSuccessUrl = plan.IdpSloSuccessUrl.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IdpErrorMsg) {
		addRequest.IdpErrorMsg = plan.IdpErrorMsg.ValueString()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *idpDefaultUrlsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_default_urls"
}

func (r *idpDefaultUrlsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *idpDefaultUrlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpDefaultUrlsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdpDefaultUrls := client.NewIdpDefaultUrl(plan.IdpErrorMsg.ValueString())
	err := addOptionalIdpDefaultUrlsFields(ctx, createIdpDefaultUrls, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdP default URL settings", err.Error())
		return
	}

	apiCreateIdpDefaultUrls := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateIdpDefaultUrls = apiCreateIdpDefaultUrls.Body(*createIdpDefaultUrls)
	idpDefaultUrlsResponse, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettingsExecute(apiCreateIdpDefaultUrls)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdP default URL settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state idpDefaultUrlsModel
	readIdpDefaultUrlsResponse(ctx, idpDefaultUrlsResponse, &state, nil)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpDefaultUrlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpDefaultUrlsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpDefaultUrls, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.GetDefaultUrl(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the IdP default URL settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the IdP default URL settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readIdpDefaultUrlsResponse(ctx, apiReadIdpDefaultUrls, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpDefaultUrlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan idpDefaultUrlsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIdpDefaultUrls := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewIdpDefaultUrl(plan.IdpErrorMsg.ValueString())
	err := addOptionalIdpDefaultUrlsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdP default URL settings", err.Error())
		return
	}

	updateIdpDefaultUrls = updateIdpDefaultUrls.Body(*createUpdateRequest)
	updateIdpDefaultUrlsResponse, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettingsExecute(updateIdpDefaultUrls)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdP default URL settings", err, httpResp)
		return
	}

	// Get the current state to see how any attributes are changing
	var state idpDefaultUrlsModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readIdpDefaultUrlsResponse(ctx, updateIdpDefaultUrlsResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *idpDefaultUrlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *idpDefaultUrlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
