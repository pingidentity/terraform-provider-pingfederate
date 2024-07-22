// Code generated by ping-terraform-plugin-framework-generator

package spdefaulturls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &spDefaultUrlsResource{}
	_ resource.ResourceWithConfigure   = &spDefaultUrlsResource{}
	_ resource.ResourceWithImportState = &spDefaultUrlsResource{}
)

func SpDefaultUrlsResource() resource.Resource {
	return &spDefaultUrlsResource{}
}

type spDefaultUrlsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *spDefaultUrlsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_default_urls"
}

func (r *spDefaultUrlsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type spDefaultUrlsResourceModel struct {
	ConfirmSlo    types.Bool   `tfsdk:"confirm_slo"`
	SloSuccessUrl types.String `tfsdk:"slo_success_url"`
	SsoSuccessUrl types.String `tfsdk:"sso_success_url"`
}

func (r *spDefaultUrlsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage SP Default URLs. These are values that affect the user's experience when executing SP-initiated SSO operations.",
		Attributes: map[string]schema.Attribute{
			"confirm_slo": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Determines whether the user is prompted to confirm Single Logout (SLO). The default is `false`.",
			},
			"slo_success_url": schema.StringAttribute{
				Optional:    true,
				Description: "Provide the default URL you would like to send the user to when Single Logout (SLO) has succeeded.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidUrl(),
				},
			},
			"sso_success_url": schema.StringAttribute{
				Optional:    true,
				Description: "Provide the default URL you would like to send the user to when Single Sign On (SSO) has succeeded.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidUrl(),
				},
			},
		},
	}
}

func (model *spDefaultUrlsResourceModel) buildClientStruct() *client.SpDefaultUrls {
	result := &client.SpDefaultUrls{}
	// confirm_slo
	result.ConfirmSlo = model.ConfirmSlo.ValueBoolPointer()
	// slo_success_url
	result.SloSuccessUrl = model.SloSuccessUrl.ValueStringPointer()
	// sso_success_url
	result.SsoSuccessUrl = model.SsoSuccessUrl.ValueStringPointer()
	return result
}

func (state *spDefaultUrlsResourceModel) readClientResponse(response *client.SpDefaultUrls) diag.Diagnostics {
	// confirm_slo
	state.ConfirmSlo = types.BoolPointerValue(response.ConfirmSlo)
	// slo_success_url
	state.SloSuccessUrl = types.StringPointerValue(response.SloSuccessUrl)
	// sso_success_url
	state.SsoSuccessUrl = types.StringPointerValue(response.SsoSuccessUrl)
	return nil
}

func (r *spDefaultUrlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data spDefaultUrlsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the spDefaultUrls", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spDefaultUrlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data spDefaultUrlsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.GetSpDefaultUrls(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the spDefaultUrls", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the spDefaultUrls", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spDefaultUrlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data spDefaultUrlsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the spDefaultUrls", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spDefaultUrlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
}

func (r *spDefaultUrlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState spDefaultUrlsResourceModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
