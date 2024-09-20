// Code generated by ping-terraform-plugin-framework-generator

package oauthidpadaptermappings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &oauthIdpAdapterMappingResource{}
	_ resource.ResourceWithConfigure   = &oauthIdpAdapterMappingResource{}
	_ resource.ResourceWithImportState = &oauthIdpAdapterMappingResource{}

	customId = "mapping_id"
)

func OauthIdpAdapterMappingResource() resource.Resource {
	return &oauthIdpAdapterMappingResource{}
}

type oauthIdpAdapterMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *oauthIdpAdapterMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_idp_adapter_mapping"
}

func (r *oauthIdpAdapterMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type oauthIdpAdapterMappingResourceModel struct {
	AttributeContractFulfillment types.Map    `tfsdk:"attribute_contract_fulfillment"`
	AttributeSources             types.Set    `tfsdk:"attribute_sources"`
	Id                           types.String `tfsdk:"id"`
	IdpAdapterRef                types.Object `tfsdk:"idp_adapter_ref"`
	IssuanceCriteria             types.Object `tfsdk:"issuance_criteria"`
	MappingId                    types.String `tfsdk:"mapping_id"`
}

func (r *oauthIdpAdapterMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage IdP adapter mappings.",
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchemaWithSuffix(true, false, false, " Map values `USER_NAME` and `USER_KEY` are required.  If extended attributes are configured on the persistent grant contract (for example, using the `pingfederate_oauth_auth_server_settings` resource), these must also be configured as map keys."),
			"attribute_sources":              attributesources.ToSchema(0, false),
			"idp_adapter_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "The ID of the resource.",
					},
				},
				Optional:    false,
				Computed:    true,
				Description: "Read only reference to the associated IdP adapter.",
			},
			"issuance_criteria": issuancecriteria.ToSchema(),
			"mapping_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the adapter mapping. This field is immutable and will trigger a replacement plan if changed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *oauthIdpAdapterMappingResourceModel) buildClientStruct() (*client.IdpAdapterMapping, error) {
	result := &client.IdpAdapterMapping{}
	var err error
	// attribute_contract_fulfillment
	result.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(model.AttributeContractFulfillment)
	if err != nil {
		return nil, err
	}

	// attribute_sources
	result.AttributeSources, err = attributesources.ClientStruct(model.AttributeSources)
	if err != nil {
		return nil, err
	}

	// issuance_criteria
	result.IssuanceCriteria, err = issuancecriteria.ClientStruct(model.IssuanceCriteria)
	if err != nil {
		return nil, err
	}

	// mapping_id
	result.Id = model.MappingId.ValueString()
	return result, nil
}

func (state *oauthIdpAdapterMappingResourceModel) readClientResponse(response *client.IdpAdapterMapping) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
	// attribute_contract_fulfillment
	attributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &response.AttributeContractFulfillment)
	respDiags.Append(diags...)

	state.AttributeContractFulfillment = attributeContractFulfillmentValue
	// attribute_sources
	attributeSourcesValue, diags := attributesources.ToState(context.Background(), response.AttributeSources)
	respDiags.Append(diags...)

	state.AttributeSources = attributeSourcesValue
	// idp_adapter_ref
	idpAdapterRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var idpAdapterRefValue types.Object
	if response.IdpAdapterRef == nil {
		idpAdapterRefValue = types.ObjectNull(idpAdapterRefAttrTypes)
	} else {
		idpAdapterRefValue, diags = types.ObjectValue(idpAdapterRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.IdpAdapterRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.IdpAdapterRef = idpAdapterRefValue
	// issuance_criteria
	issuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), response.IssuanceCriteria)
	respDiags.Append(diags...)

	state.IssuanceCriteria = issuanceCriteriaValue
	// mapping_id
	state.MappingId = types.StringValue(response.Id)
	return respDiags
}

func (r *oauthIdpAdapterMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data oauthIdpAdapterMappingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, err := data.buildClientStruct()
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to build client struct for the oauthIdpAdapterMapping: "+err.Error())
		return
	}
	apiCreateRequest := r.apiClient.OauthIdpAdapterMappingsAPI.CreateIdpAdapterMapping(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthIdpAdapterMappingsAPI.CreateIdpAdapterMappingExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the oauthIdpAdapterMapping", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthIdpAdapterMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data oauthIdpAdapterMappingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.OauthIdpAdapterMappingsAPI.GetIdpAdapterMapping(config.AuthContext(ctx, r.providerConfig), data.MappingId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth IdP Adapter Mapping", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the oauthIdpAdapterMapping", err, httpResp, &customId)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthIdpAdapterMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data oauthIdpAdapterMappingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, err := data.buildClientStruct()
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to build client struct for the oauthIdpAdapterMapping: "+err.Error())
		return
	}
	apiUpdateRequest := r.apiClient.OauthIdpAdapterMappingsAPI.UpdateIdpAdapterMapping(config.AuthContext(ctx, r.providerConfig), data.MappingId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthIdpAdapterMappingsAPI.UpdateIdpAdapterMappingExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the oauthIdpAdapterMapping", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthIdpAdapterMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data oauthIdpAdapterMappingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.OauthIdpAdapterMappingsAPI.DeleteIdpAdapterMapping(config.AuthContext(ctx, r.providerConfig), data.MappingId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the oauthIdpAdapterMapping", err, httpResp, &customId)
	}
}

func (r *oauthIdpAdapterMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to mapping_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("mapping_id"), req, resp)
}
