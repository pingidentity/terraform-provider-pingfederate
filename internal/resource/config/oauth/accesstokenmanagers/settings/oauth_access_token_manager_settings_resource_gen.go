// Code generated by ping-terraform-plugin-framework-generator

package oauthaccesstokenmanagerssettings

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
	_ resource.Resource                = &oauthAccessTokenManagerSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenManagerSettingsResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenManagerSettingsResource{}
)

func OauthAccessTokenManagerSettingsResource() resource.Resource {
	return &oauthAccessTokenManagerSettingsResource{}
}

type oauthAccessTokenManagerSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *oauthAccessTokenManagerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_manager_settings"
}

func (r *oauthAccessTokenManagerSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type oauthAccessTokenManagerSettingsResourceModel struct {
	DefaultAccessTokenManagerRef types.Object `tfsdk:"default_access_token_manager_ref"`
}

func (r *oauthAccessTokenManagerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the general access token management settings.",
		Attributes: map[string]schema.Attribute{
			"default_access_token_manager_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Required:    true,
				Description: "Reference to the default access token manager, if one is defined.",
			},
		},
	}
}

func (model *oauthAccessTokenManagerSettingsResourceModel) buildClientStruct() *client.AccessTokenManagementSettings {
	result := &client.AccessTokenManagementSettings{}
	// default_access_token_manager_ref
	if !model.DefaultAccessTokenManagerRef.IsNull() {
		defaultAccessTokenManagerRefValue := &client.ResourceLink{}
		defaultAccessTokenManagerRefAttrs := model.DefaultAccessTokenManagerRef.Attributes()
		defaultAccessTokenManagerRefValue.Id = defaultAccessTokenManagerRefAttrs["id"].(types.String).ValueString()
		result.DefaultAccessTokenManagerRef = defaultAccessTokenManagerRefValue
	}

	return result
}

func (state *oauthAccessTokenManagerSettingsResourceModel) readClientResponse(response *client.AccessTokenManagementSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// default_access_token_manager_ref
	defaultAccessTokenManagerRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var defaultAccessTokenManagerRefValue types.Object
	if response.DefaultAccessTokenManagerRef == nil {
		defaultAccessTokenManagerRefValue = types.ObjectNull(defaultAccessTokenManagerRefAttrTypes)
	} else {
		defaultAccessTokenManagerRefValue, diags = types.ObjectValue(defaultAccessTokenManagerRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.DefaultAccessTokenManagerRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.DefaultAccessTokenManagerRef = defaultAccessTokenManagerRefValue
	return respDiags
}

func (r *oauthAccessTokenManagerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data oauthAccessTokenManagerSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.OauthAccessTokenManagersAPI.UpdateOauthAccessTokenManagersSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.UpdateOauthAccessTokenManagersSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the oauthAccessTokenManagerSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAccessTokenManagerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data oauthAccessTokenManagerSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.GetOauthAccessTokenManagersSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the oauthAccessTokenManagerSettings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the oauthAccessTokenManagerSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAccessTokenManagerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data oauthAccessTokenManagerSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.OauthAccessTokenManagersAPI.UpdateOauthAccessTokenManagersSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.UpdateOauthAccessTokenManagersSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the oauthAccessTokenManagerSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAccessTokenManagerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	resp.Diagnostics.AddWarning("Configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}

func (r *oauthAccessTokenManagerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState oauthAccessTokenManagerSettingsResourceModel
	emptyState.setNullObjectValues()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
