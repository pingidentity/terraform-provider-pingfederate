package defaulturls

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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

var (
	_ resource.Resource                = &defaultUrlsResource{}
	_ resource.ResourceWithConfigure   = &defaultUrlsResource{}
	_ resource.ResourceWithImportState = &defaultUrlsResource{}
)

func DefaultUrlsResource() resource.Resource {
	return &defaultUrlsResource{}
}

type defaultUrlsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *defaultUrlsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_urls"
}

func (r *defaultUrlsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type defaultUrlsResourceModel struct {
	ConfirmSpSlo     types.Bool   `tfsdk:"confirm_sp_slo"`
	SpSloSuccessUrl  types.String `tfsdk:"sp_slo_success_url"`
	SpSsoSuccessUrl  types.String `tfsdk:"sp_sso_success_url"`
	ConfirmIdpSlo    types.Bool   `tfsdk:"confirm_idp_slo"`
	IdpErrorMsg      types.String `tfsdk:"idp_error_msg"`
	IdpSloSuccessUrl types.String `tfsdk:"idp_slo_success_url"`
}

func (r *defaultUrlsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage IdP and SP default URL settings.",
		Attributes: map[string]schema.Attribute{
			// Sp default URLs attributes
			"confirm_sp_slo": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "SP setting to prompt user to confirm Single Logout (SLO). The default is `false`.",
			},
			"sp_slo_success_url": schema.StringAttribute{
				Optional:    true,
				Description: "SP setting for the default URL you would like to send the user to when Single Logout (SLO) has succeeded.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidUrl(),
				},
			},
			"sp_sso_success_url": schema.StringAttribute{
				Optional:    true,
				Description: "SP setting for the default URL you would like to send the user to when Single Sign On (SSO) has succeeded.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidUrl(),
				},
			},
			// IdP default URLs attributes
			"confirm_idp_slo": schema.BoolAttribute{
				Description: "IdP setting to prompt user to confirm Single Logout (SLO). The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idp_error_msg": schema.StringAttribute{
				Description: "IdP setting for the error text displayed in a user's browser when an SSO operation fails.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"idp_slo_success_url": schema.StringAttribute{
				Description: "Idp setting for the default URL you would like to send the user to when Single Logout has succeeded.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidUrl(),
				},
			},
		},
	}
}

// Validation must occur in modifyplan since it depends on attribute defaults
func (r *defaultUrlsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *defaultUrlsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	if internaltypes.IsDefined(plan.ConfirmIdpSlo) && internaltypes.IsDefined(plan.ConfirmSpSlo) && plan.ConfirmIdpSlo.ValueBool() != plan.ConfirmSpSlo.ValueBool() {
		resp.Diagnostics.AddError(
			providererror.InvalidAttributeConfiguration,
			"`confirm_idp_slo` and `confirm_sp_slo` must be set to the same value")
	}
}

func (model *defaultUrlsResourceModel) buildSpClientStruct() *client.SpDefaultUrls {
	result := &client.SpDefaultUrls{}
	// confirm_sp_slo
	result.ConfirmSlo = model.ConfirmSpSlo.ValueBoolPointer()
	// sp_slo_success_url
	result.SloSuccessUrl = model.SpSloSuccessUrl.ValueStringPointer()
	// sp_sso_success_url
	result.SsoSuccessUrl = model.SpSsoSuccessUrl.ValueStringPointer()
	return result
}

func (model *defaultUrlsResourceModel) buildIdpClientStruct() *client.IdpDefaultUrl {
	result := &client.IdpDefaultUrl{}
	// confirm_idp_slo
	result.ConfirmIdpSlo = model.ConfirmIdpSlo.ValueBoolPointer()
	// idp_slo_success_url
	result.IdpSloSuccessUrl = model.IdpSloSuccessUrl.ValueStringPointer()
	// idp_error_msg
	result.IdpErrorMsg = model.IdpErrorMsg.ValueString()
	return result
}

func (state *defaultUrlsResourceModel) readSpClientResponse(response *client.SpDefaultUrls) diag.Diagnostics {
	// confirm_slo
	state.ConfirmSpSlo = types.BoolPointerValue(response.ConfirmSlo)
	// slo_success_url
	state.SpSloSuccessUrl = types.StringPointerValue(response.SloSuccessUrl)
	// sso_success_url
	state.SpSsoSuccessUrl = types.StringPointerValue(response.SsoSuccessUrl)
	return nil
}

func (state *defaultUrlsResourceModel) readIdpClientResponse(response *client.IdpDefaultUrl) diag.Diagnostics {
	// confirm_slo
	state.ConfirmIdpSlo = types.BoolPointerValue(response.ConfirmIdpSlo)
	// slo_success_url
	state.IdpSloSuccessUrl = types.StringPointerValue(response.IdpSloSuccessUrl)
	// sso_success_url
	state.IdpErrorMsg = types.StringValue(response.IdpErrorMsg)
	return nil
}

func (r *defaultUrlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data defaultUrlsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// SP
	// Update API call logic, since this is a singleton resource
	spClientData := data.buildSpClientStruct()
	spApiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	spApiUpdateRequest = spApiUpdateRequest.Body(*spClientData)
	spResponseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(spApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the SP default URLs", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readSpClientResponse(spResponseData)...)

	// IdP
	// Update API call logic, since this is a singleton resource
	idpClientData := data.buildIdpClientStruct()
	idpApiUpdateRequest := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettings(config.AuthContext(ctx, r.providerConfig))
	idpApiUpdateRequest = idpApiUpdateRequest.Body(*idpClientData)
	idpResponseData, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettingsExecute(idpApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdP default URLs", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readIdpClientResponse(idpResponseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *defaultUrlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data defaultUrlsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// SP
	// Read API call logic
	spResponseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.GetSpDefaultUrls(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "SP Default URLs", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the SP default URLs", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readSpClientResponse(spResponseData)...)

	// IdP
	// Read API call logic
	idpResponseData, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.GetDefaultUrl(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "IdP Default URLs", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the IdP default URLs", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readIdpClientResponse(idpResponseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *defaultUrlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data defaultUrlsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// SP
	// Update API call logic, since this is a singleton resource
	spClientData := data.buildSpClientStruct()
	spApiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	spApiUpdateRequest = spApiUpdateRequest.Body(*spClientData)
	spResponseData, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(spApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the SP default URLs", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readSpClientResponse(spResponseData)...)

	// IdP
	// Update API call logic, since this is a singleton resource
	idpClientData := data.buildIdpClientStruct()
	idpApiUpdateRequest := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettings(config.AuthContext(ctx, r.providerConfig))
	idpApiUpdateRequest = idpApiUpdateRequest.Body(*idpClientData)
	idpResponseData, httpResp, err := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettingsExecute(idpApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the IdP default URLs", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readIdpClientResponse(idpResponseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (model *defaultUrlsResource) buildDefaultSpClientStruct() *client.SpDefaultUrls {
	result := &client.SpDefaultUrls{}
	// confirm_slo
	result.ConfirmSlo = utils.Pointer(false)
	// slo_success_url
	result.SloSuccessUrl = utils.Pointer("")
	// sso_success_url
	result.SsoSuccessUrl = utils.Pointer("")
	return result
}

func (model *defaultUrlsResource) buildDefaultIdpClientStruct() *client.IdpDefaultUrl {
	result := &client.IdpDefaultUrl{}
	// confirm_slo
	result.ConfirmIdpSlo = utils.Pointer(false)
	// slo_success_url
	result.IdpSloSuccessUrl = utils.Pointer("")
	// sso_success_url
	result.IdpErrorMsg = "errorDetail.idpSsoFailure"
	return result
}

func (r *defaultUrlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	// SP
	spClientData := r.buildDefaultSpClientStruct()
	spApiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	spApiUpdateRequest = spApiUpdateRequest.Body(*spClientData)
	_, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(spApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the SP default URLs", err, httpResp)
	}

	// IdP
	idpClientData := r.buildDefaultIdpClientStruct()
	idpApiUpdateRequest := r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettings(config.AuthContext(ctx, r.providerConfig))
	idpApiUpdateRequest = idpApiUpdateRequest.Body(*idpClientData)
	_, httpResp, err = r.apiClient.IdpDefaultUrlsAPI.UpdateDefaultUrlSettingsExecute(idpApiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the IDP default URLs", err, httpResp)
	}
}

func (r *defaultUrlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState defaultUrlsResourceModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
