// Code generated by ping-terraform-plugin-framework-generator

package serversettingsoutboundprovisioning

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

var (
	_ resource.Resource                = &serverSettingsOutboundProvisioningResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsOutboundProvisioningResource{}
	_ resource.ResourceWithImportState = &serverSettingsOutboundProvisioningResource{}
)

func ServerSettingsOutboundProvisioningResource() resource.Resource {
	return &serverSettingsOutboundProvisioningResource{}
}

type serverSettingsOutboundProvisioningResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *serverSettingsOutboundProvisioningResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_outbound_provisioning"
}

func (r *serverSettingsOutboundProvisioningResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type serverSettingsOutboundProvisioningResourceModel struct {
	DataStoreRef             types.Object `tfsdk:"data_store_ref"`
	SynchronizationFrequency types.Int64  `tfsdk:"synchronization_frequency"`
}

func (r *serverSettingsOutboundProvisioningResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the outbound provisioning settings.",
		Attributes: map[string]schema.Attribute{
			"data_store_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:    true,
				Description: "Reference to the associated data store.",
			},
			"synchronization_frequency": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60),
				Description: "The synchronization frequency in seconds. The default value is `60`.",
			},
		},
	}
}

func (model *serverSettingsOutboundProvisioningResourceModel) buildClientStruct() (*client.OutboundProvisionDatabase, diag.Diagnostics) {
	result := &client.OutboundProvisionDatabase{}
	// data_store_ref
	dataStoreRefValue := client.ResourceLink{}
	if !model.DataStoreRef.IsNull() {
		dataStoreRefAttrs := model.DataStoreRef.Attributes()
		dataStoreRefValue.Id = dataStoreRefAttrs["id"].(types.String).ValueString()
	}
	result.DataStoreRef = dataStoreRefValue

	// synchronization_frequency
	result.SynchronizationFrequency = model.SynchronizationFrequency.ValueInt64Pointer()
	return result, nil
}

// Build a default client struct to reset the resource to its default state
// If necessary, update this function to set any other values that should be present in the default state of the resource
func (model *serverSettingsOutboundProvisioningResource) buildDefaultClientStruct() *client.OutboundProvisionDatabase {
	result := &client.OutboundProvisionDatabase{
		SynchronizationFrequency: utils.Pointer(int64(60)),
	}
	return result
}

func (state *serverSettingsOutboundProvisioningResourceModel) readClientResponse(response *client.OutboundProvisionDatabase) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// data_store_ref
	dataStoreRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var dataStoreRefValue types.Object
	if response.DataStoreRef.Id == "" {
		dataStoreRefValue = types.ObjectNull(dataStoreRefAttrTypes)
	} else {
		dataStoreRefValue, diags = types.ObjectValue(dataStoreRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.DataStoreRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.DataStoreRef = dataStoreRefValue
	// synchronization_frequency
	state.SynchronizationFrequency = types.Int64PointerValue(response.SynchronizationFrequency)
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *serverSettingsOutboundProvisioningResource) emptyModel() serverSettingsOutboundProvisioningResourceModel {
	var model serverSettingsOutboundProvisioningResourceModel
	// data_store_ref
	dataStoreRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	model.DataStoreRef = types.ObjectNull(dataStoreRefAttrTypes)
	return model
}

func (r *serverSettingsOutboundProvisioningResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverSettingsOutboundProvisioningResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the serverSettingsOutboundProvisioning", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsOutboundProvisioningResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverSettingsOutboundProvisioningResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.GetOutBoundProvisioningSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "serverSettingsOutboundProvisioning", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the serverSettingsOutboundProvisioning", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsOutboundProvisioningResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data serverSettingsOutboundProvisioningResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the serverSettingsOutboundProvisioning", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsOutboundProvisioningResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service.
	// Instead this delete method will attempt to set the resource to its default state on the service. If this isn't possible,
	// this method can be replaced with a no-op with a diagnostic warning message about being unable to set to the default state.
	// Update API call logic to reset to default
	defaultClientData := r.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*defaultClientData)
	_, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateOutBoundProvisioningSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the serverSettingsOutboundProvisioning", err, httpResp)
	}
}

func (r *serverSettingsOutboundProvisioningResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
