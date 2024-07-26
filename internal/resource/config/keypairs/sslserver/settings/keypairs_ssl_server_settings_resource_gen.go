// Code generated by ping-terraform-plugin-framework-generator

package keypairssslserversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &keypairsSslServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &keypairsSslServerSettingsResource{}
	_ resource.ResourceWithImportState = &keypairsSslServerSettingsResource{}
)

func KeypairsSslServerSettingsResource() resource.Resource {
	return &keypairsSslServerSettingsResource{}
}

type keypairsSslServerSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSslServerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_server_settings"
}

func (r *keypairsSslServerSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSslServerSettingsResourceModel struct {
	ActiveAdminConsoleCerts  types.Set    `tfsdk:"active_admin_console_certs"`
	ActiveRuntimeServerCerts types.Set    `tfsdk:"active_runtime_server_certs"`
	AdminConsoleCertRef      types.Object `tfsdk:"admin_console_cert_ref"`
	RuntimeServerCertRef     types.Object `tfsdk:"runtime_server_cert_ref"`
}

func (r *keypairsSslServerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the SSL server certificate settings.",
		Attributes: map[string]schema.Attribute{
			"active_admin_console_certs": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the resource.",
						},
					},
				},
				Required:    true,
				Description: "The active SSL Server Certificate Key pairs for PF Administrator Console. Must not be empty and must contain a reference to the cert configured in `admin_console_cert_ref.id`.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"active_runtime_server_certs": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the resource.",
						},
					},
				},
				Required:    true,
				Description: "The active SSL Server Certificate Key pairs for Runtime Server. Must not be empty and must contain a reference to the cert configured in `runtime_server_cert_ref.id`.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"admin_console_cert_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Required:    true,
				Description: "Reference to the default SSL Server Certificate Key pair active for PF Administrator Console.",
			},
			"runtime_server_cert_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Required:    true,
				Description: "Reference to the default SSL Server Certificate Key pair active for Runtime Server.",
			},
		},
	}
}

func (model *keypairsSslServerSettingsResourceModel) buildClientStruct() *client.SslServerSettings {
	result := &client.SslServerSettings{}
	// active_admin_console_certs
	result.ActiveAdminConsoleCerts = []client.ResourceLink{}
	for _, activeAdminConsoleCertsElement := range model.ActiveAdminConsoleCerts.Elements() {
		activeAdminConsoleCertsValue := client.ResourceLink{}
		activeAdminConsoleCertsAttrs := activeAdminConsoleCertsElement.(types.Object).Attributes()
		activeAdminConsoleCertsValue.Id = activeAdminConsoleCertsAttrs["id"].(types.String).ValueString()
		result.ActiveAdminConsoleCerts = append(result.ActiveAdminConsoleCerts, activeAdminConsoleCertsValue)
	}

	// active_runtime_server_certs
	result.ActiveRuntimeServerCerts = []client.ResourceLink{}
	for _, activeRuntimeServerCertsElement := range model.ActiveRuntimeServerCerts.Elements() {
		activeRuntimeServerCertsValue := client.ResourceLink{}
		activeRuntimeServerCertsAttrs := activeRuntimeServerCertsElement.(types.Object).Attributes()
		activeRuntimeServerCertsValue.Id = activeRuntimeServerCertsAttrs["id"].(types.String).ValueString()
		result.ActiveRuntimeServerCerts = append(result.ActiveRuntimeServerCerts, activeRuntimeServerCertsValue)
	}

	// admin_console_cert_ref
	adminConsoleCertRefValue := client.ResourceLink{}
	adminConsoleCertRefAttrs := model.AdminConsoleCertRef.Attributes()
	adminConsoleCertRefValue.Id = adminConsoleCertRefAttrs["id"].(types.String).ValueString()
	result.AdminConsoleCertRef = adminConsoleCertRefValue

	// runtime_server_cert_ref
	runtimeServerCertRefValue := client.ResourceLink{}
	runtimeServerCertRefAttrs := model.RuntimeServerCertRef.Attributes()
	runtimeServerCertRefValue.Id = runtimeServerCertRefAttrs["id"].(types.String).ValueString()
	result.RuntimeServerCertRef = runtimeServerCertRefValue

	return result
}

func (state *keypairsSslServerSettingsResourceModel) readClientResponse(response *client.SslServerSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// active_admin_console_certs
	activeAdminConsoleCertsAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	activeAdminConsoleCertsElementType := types.ObjectType{AttrTypes: activeAdminConsoleCertsAttrTypes}
	var activeAdminConsoleCertsValues []attr.Value
	for _, activeAdminConsoleCertsResponseValue := range response.ActiveAdminConsoleCerts {
		activeAdminConsoleCertsValue, diags := types.ObjectValue(activeAdminConsoleCertsAttrTypes, map[string]attr.Value{
			"id": types.StringValue(activeAdminConsoleCertsResponseValue.Id),
		})
		respDiags.Append(diags...)
		activeAdminConsoleCertsValues = append(activeAdminConsoleCertsValues, activeAdminConsoleCertsValue)
	}
	activeAdminConsoleCertsValue, diags := types.SetValue(activeAdminConsoleCertsElementType, activeAdminConsoleCertsValues)
	respDiags.Append(diags...)

	state.ActiveAdminConsoleCerts = activeAdminConsoleCertsValue
	// active_runtime_server_certs
	activeRuntimeServerCertsAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	activeRuntimeServerCertsElementType := types.ObjectType{AttrTypes: activeRuntimeServerCertsAttrTypes}
	var activeRuntimeServerCertsValues []attr.Value
	for _, activeRuntimeServerCertsResponseValue := range response.ActiveRuntimeServerCerts {
		activeRuntimeServerCertsValue, diags := types.ObjectValue(activeRuntimeServerCertsAttrTypes, map[string]attr.Value{
			"id": types.StringValue(activeRuntimeServerCertsResponseValue.Id),
		})
		respDiags.Append(diags...)
		activeRuntimeServerCertsValues = append(activeRuntimeServerCertsValues, activeRuntimeServerCertsValue)
	}
	activeRuntimeServerCertsValue, diags := types.SetValue(activeRuntimeServerCertsElementType, activeRuntimeServerCertsValues)
	respDiags.Append(diags...)

	state.ActiveRuntimeServerCerts = activeRuntimeServerCertsValue
	// admin_console_cert_ref
	adminConsoleCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	adminConsoleCertRefValue, diags := types.ObjectValue(adminConsoleCertRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.AdminConsoleCertRef.Id),
	})
	respDiags.Append(diags...)

	state.AdminConsoleCertRef = adminConsoleCertRefValue
	// runtime_server_cert_ref
	runtimeServerCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	runtimeServerCertRefValue, diags := types.ObjectValue(runtimeServerCertRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.RuntimeServerCertRef.Id),
	})
	respDiags.Append(diags...)

	state.RuntimeServerCertRef = runtimeServerCertRefValue
	return respDiags
}

func (r *keypairsSslServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSslServerSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.KeyPairsSslServerAPI.UpdateSslServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KeyPairsSslServerAPI.UpdateSslServerSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the keypairsSslServerSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keypairsSslServerSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KeyPairsSslServerAPI.GetSslServerSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the keypairsSslServerSettings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the keypairsSslServerSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslServerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data keypairsSslServerSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.KeyPairsSslServerAPI.UpdateSslServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KeyPairsSslServerAPI.UpdateSslServerSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the keypairsSslServerSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslServerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	resp.Diagnostics.AddWarning("`keypairs_ssl_server_settings` configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}

func (r *keypairsSslServerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState keypairsSslServerSettingsResourceModel
	emptyState.setNullObjectValues()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
