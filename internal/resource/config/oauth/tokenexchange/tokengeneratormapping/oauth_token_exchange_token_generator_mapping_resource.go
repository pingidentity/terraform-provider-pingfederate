package oauthtokenexchangetokengeneratormapping

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthTokenExchangeTokenGeneratorMappingResource{}
	_ resource.ResourceWithConfigure   = &oauthTokenExchangeTokenGeneratorMappingResource{}
	_ resource.ResourceWithImportState = &oauthTokenExchangeTokenGeneratorMappingResource{}
)

// OauthTokenExchangeTokenGeneratorMappingResource is a helper function to simplify the provider implementation.
func OauthTokenExchangeTokenGeneratorMappingResource() resource.Resource {
	return &oauthTokenExchangeTokenGeneratorMappingResource{}
}

// oauthTokenExchangeTokenGeneratorMappingResource is the resource implementation.
type oauthTokenExchangeTokenGeneratorMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthTokenExchangeTokenGeneratorMappingResourceModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the resource.
func (r *oauthTokenExchangeTokenGeneratorMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the mapping from a token exchange processor policy to a token generator.",
		Attributes: map[string]schema.Attribute{
			"attribute_sources":              attributesources.ToSchema(0),
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
			"issuance_criteria":              issuancecriteria.ToSchema(),
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Exchange Processor policy.",
				Required:    true,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator",
				Required:    true,
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group",
				Optional:    true,
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalOauthTokenExchangeTokenGeneratorMappingFields(ctx context.Context, addRequest *client.ProcessorPolicyToGeneratorMapping, plan oauthTokenExchangeTokenGeneratorMappingResourceModel) error {
	if internaltypes.IsDefined(plan.AttributeSources) {
		addRequest.AttributeSources = []client.AttributeSourceAggregation{}
		var attributeSourcesErr error
		addRequest.AttributeSources, attributeSourcesErr = attributesources.ClientStruct(plan.AttributeSources)
		if attributeSourcesErr != nil {
			return attributeSourcesErr
		}
	}

	if internaltypes.IsDefined(plan.IssuanceCriteria) {
		addRequest.IssuanceCriteria = client.NewIssuanceCriteria()
		var issuanceCriteriaErr error
		addRequest.IssuanceCriteria, issuanceCriteriaErr = issuancecriteria.ClientStruct(plan.IssuanceCriteria)
		if issuanceCriteriaErr != nil {
			return issuanceCriteriaErr
		}
	}
	if internaltypes.IsDefined(plan.LicenseConnectionGroupAssignment) {
		addRequest.LicenseConnectionGroupAssignment = plan.LicenseConnectionGroupAssignment.ValueStringPointer()
	}

	return nil
}

// Metadata returns the resource type name.
func (r *oauthTokenExchangeTokenGeneratorMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token_exchange_token_generator_mapping"
}

func (r *oauthTokenExchangeTokenGeneratorMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readOauthTokenExchangeTokenGeneratorMappingResourceResponse(ctx context.Context, r *client.ProcessorPolicyToGeneratorMapping, state *oauthTokenExchangeTokenGeneratorMappingResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.AttributeSources, respDiags = attributesources.ToState(ctx, r.AttributeSources)
	diags.Append(respDiags...)
	state.AttributeContractFulfillment, respDiags = attributecontractfulfillment.ToState(ctx, r.AttributeContractFulfillment)
	diags.Append(respDiags...)
	state.IssuanceCriteria, respDiags = issuancecriteria.ToState(ctx, r.IssuanceCriteria)
	diags.Append(respDiags...)
	state.SourceId = types.StringValue(r.SourceId)
	state.TargetId = types.StringValue(r.TargetId)
	state.Id = types.StringPointerValue(r.Id)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
	return diags
}

func (r *oauthTokenExchangeTokenGeneratorMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthTokenExchangeTokenGeneratorMappingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	attributeContractFulfillment := &map[string]client.AttributeFulfillmentValue{}
	attributeContractFulfillmentErr := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContractFulfillment, false)), attributeContractFulfillment)
	if attributeContractFulfillmentErr != nil {
		resp.Diagnostics.AddError("Failed to build attribute contract fulfillment request object:", attributeContractFulfillmentErr.Error())
		return
	}
	createOauthTokenExchangeTokenGeneratorMapping := client.NewProcessorPolicyToGeneratorMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalOauthTokenExchangeTokenGeneratorMappingFields(ctx, createOauthTokenExchangeTokenGeneratorMapping, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an OAuth Token Exchange Token Generator Mapping", err.Error())
		return
	}

	apiCreateOauthTokenExchangeTokenGeneratorMapping := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.CreateTokenGeneratorMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthTokenExchangeTokenGeneratorMapping = apiCreateOauthTokenExchangeTokenGeneratorMapping.Body(*createOauthTokenExchangeTokenGeneratorMapping)
	oauthTokenExchangeTokenGeneratorMappingResponse, httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.CreateTokenGeneratorMappingExecute(apiCreateOauthTokenExchangeTokenGeneratorMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating an OAuth Token Exchange Token Generator Mapping", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthTokenExchangeTokenGeneratorMappingResourceModel

	diags = readOauthTokenExchangeTokenGeneratorMappingResourceResponse(ctx, oauthTokenExchangeTokenGeneratorMappingResponse, &state)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthTokenExchangeTokenGeneratorMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthTokenExchangeTokenGeneratorMappingResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthTokenExchangeTokenGeneratorMapping, httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.GetTokenGeneratorMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Token Exchange Token Generator Mapping", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Token Exchange Token Generator Mapping", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readOauthTokenExchangeTokenGeneratorMappingResourceResponse(ctx, apiReadOauthTokenExchangeTokenGeneratorMapping, &state)
	resp.Diagnostics.Append(diags...)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthTokenExchangeTokenGeneratorMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthTokenExchangeTokenGeneratorMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	attributeContractFulfillment := &map[string]client.AttributeFulfillmentValue{}
	attributeContractFulfillmentErr := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContractFulfillment, false)), attributeContractFulfillment)
	if attributeContractFulfillmentErr != nil {
		resp.Diagnostics.AddError("Failed to build attribute contract fulfillment request object:", attributeContractFulfillmentErr.Error())
		return
	}
	updateOauthTokenExchangeTokenGeneratorMapping := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.UpdateTokenGeneratorMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewProcessorPolicyToGeneratorMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalOauthTokenExchangeTokenGeneratorMappingFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an OAuth Token Exchange Token Generator Mapping", err.Error())
		return
	}

	updateOauthTokenExchangeTokenGeneratorMapping = updateOauthTokenExchangeTokenGeneratorMapping.Body(*createUpdateRequest)
	updateOauthTokenExchangeTokenGeneratorMappingResponse, httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.UpdateTokenGeneratorMappingByIdExecute(updateOauthTokenExchangeTokenGeneratorMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an OAuth Token Exchange Token Generator Mapping", err, httpResp)
		return
	}

	// Read the response
	var state oauthTokenExchangeTokenGeneratorMappingResourceModel
	diags = readOauthTokenExchangeTokenGeneratorMappingResourceResponse(ctx, updateOauthTokenExchangeTokenGeneratorMappingResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthTokenExchangeTokenGeneratorMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthTokenExchangeTokenGeneratorMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.DeleteTokenGeneratorMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Token Exchange Token Generator Mapping", err, httpResp)
	}
}

func (r *oauthTokenExchangeTokenGeneratorMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
