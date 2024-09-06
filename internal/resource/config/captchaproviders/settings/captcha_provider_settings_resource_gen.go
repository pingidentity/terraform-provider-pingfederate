// Code generated by ping-terraform-plugin-framework-generator

package captchaproviderssettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &captchaProviderSettingsResource{}
	_ resource.ResourceWithConfigure   = &captchaProviderSettingsResource{}
	_ resource.ResourceWithImportState = &captchaProviderSettingsResource{}
)

func CaptchaProviderSettingsResource() resource.Resource {
	return &captchaProviderSettingsResource{}
}

type captchaProviderSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *captchaProviderSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_captcha_provider_settings"
}

func (r *captchaProviderSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type captchaProviderSettingsResourceModel struct {
	DefaultCaptchaProviderRef types.Object `tfsdk:"default_captcha_provider_ref"`
}

func (r *captchaProviderSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage general CAPTCHA providers settings.",
		Attributes: map[string]schema.Attribute{
			"default_captcha_provider_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Required:    true,
				Description: "Reference to the default CAPTCHA provider, if one is defined.",
			},
		},
	}
}

func (model *captchaProviderSettingsResourceModel) buildClientStruct() (*client.CaptchaProvidersSettings, diag.Diagnostics) {
	result := &client.CaptchaProvidersSettings{}
	// default_captcha_provider_ref
	if !model.DefaultCaptchaProviderRef.IsNull() {
		defaultCaptchaProviderRefValue := &client.ResourceLink{}
		defaultCaptchaProviderRefAttrs := model.DefaultCaptchaProviderRef.Attributes()
		defaultCaptchaProviderRefValue.Id = defaultCaptchaProviderRefAttrs["id"].(types.String).ValueString()
		result.DefaultCaptchaProviderRef = defaultCaptchaProviderRefValue
	}

	return result, nil
}

func (state *captchaProviderSettingsResourceModel) readClientResponse(response *client.CaptchaProvidersSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// default_captcha_provider_ref
	defaultCaptchaProviderRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var defaultCaptchaProviderRefValue types.Object
	if response.DefaultCaptchaProviderRef == nil {
		defaultCaptchaProviderRefValue = types.ObjectNull(defaultCaptchaProviderRefAttrTypes)
	} else {
		defaultCaptchaProviderRefValue, diags = types.ObjectValue(defaultCaptchaProviderRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.DefaultCaptchaProviderRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.DefaultCaptchaProviderRef = defaultCaptchaProviderRefValue
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *captchaProviderSettingsResource) emptyModel() captchaProviderSettingsResourceModel {
	var model captchaProviderSettingsResourceModel
	// default_captcha_provider_ref
	defaultCaptchaProviderRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	model.DefaultCaptchaProviderRef = types.ObjectNull(defaultCaptchaProviderRefAttrTypes)
	return model
}

func (r *captchaProviderSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data captchaProviderSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.CaptchaProvidersAPI.UpdateCaptchaProvidersSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.CaptchaProvidersAPI.UpdateCaptchaProvidersSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the captchaProviderSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *captchaProviderSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data captchaProviderSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.CaptchaProvidersAPI.GetCaptchaProvidersSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Captcha Provider Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the captchaProviderSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *captchaProviderSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data captchaProviderSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.CaptchaProvidersAPI.UpdateCaptchaProvidersSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.CaptchaProvidersAPI.UpdateCaptchaProvidersSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the captchaProviderSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *captchaProviderSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	resp.Diagnostics.AddWarning("Configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}

func (r *captchaProviderSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}