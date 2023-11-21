package tokenprocessortotokengeneratormapping

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
	_ resource.Resource                = &tokenProcessorToTokenGeneratorMappingResource{}
	_ resource.ResourceWithConfigure   = &tokenProcessorToTokenGeneratorMappingResource{}
	_ resource.ResourceWithImportState = &tokenProcessorToTokenGeneratorMappingResource{}
)

// TokenProcessorToTokenGeneratorMappingResource is a helper function to simplify the provider implementation.
func TokenProcessorToTokenGeneratorMappingResource() resource.Resource {
	return &tokenProcessorToTokenGeneratorMappingResource{}
}

// tokenProcessorToTokenGeneratorMappingResource is the resource implementation.
type tokenProcessorToTokenGeneratorMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type tokenProcessorToTokenGeneratorMappingResourceModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	Id                               types.String `tfsdk:"id"`
	DefaultTargetResource            types.String `tfsdk:"default_target_resource"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the resource.
func (r *tokenProcessorToTokenGeneratorMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Token Processor To Token Generator Mapping",
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
			"attribute_sources":              attributesources.ToSchema(0),
			"default_target_resource": schema.StringAttribute{
				Description: "Default target URL for this Token Processor to Token Generator mapping configuration.",
				Optional:    true,
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group.",
				Optional:    true,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator.",
				Required:    true,
			},
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Processor.",
				Required:    true,
			},
			"issuance_criteria": issuancecriteria.ToSchema(),
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx context.Context, addRequest *client.TokenToTokenMapping, plan tokenProcessorToTokenGeneratorMappingResourceModel) error {
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

	if internaltypes.IsDefined(plan.DefaultTargetResource) {
		addRequest.DefaultTargetResource = plan.DefaultTargetResource.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.LicenseConnectionGroupAssignment) {
		addRequest.LicenseConnectionGroupAssignment = plan.LicenseConnectionGroupAssignment.ValueStringPointer()
	}

	return nil

}

// Metadata returns the resource type name.
func (r *tokenProcessorToTokenGeneratorMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_processor_to_token_generator_mapping"
}

func (r *tokenProcessorToTokenGeneratorMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readTokenProcessorToTokenGeneratorMappingResponse(ctx context.Context, r *client.TokenToTokenMapping, state *tokenProcessorToTokenGeneratorMappingResourceModel, plan tokenProcessorToTokenGeneratorMappingResourceModel) diag.Diagnostics {
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
	state.DefaultTargetResource = types.StringPointerValue(r.DefaultTargetResource)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
	return diags
}

func (r *tokenProcessorToTokenGeneratorMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tokenProcessorToTokenGeneratorMappingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributeContractFulfillment, err := attributecontractfulfillment.ClientStruct(plan.AttributeContractFulfillment)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build attribute contract fulfillment request object:", err.Error())
		return
	}
	createTokenProcessorToTokenGeneratorMapping := client.NewTokenToTokenMapping(attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err = addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx, createTokenProcessorToTokenGeneratorMapping, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TokenProcessorToTokenGeneratorMapping", err.Error())
		return
	}

	apiCreateTokenProcessorToTokenGeneratorMapping := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.CreateTokenToTokenMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateTokenProcessorToTokenGeneratorMapping = apiCreateTokenProcessorToTokenGeneratorMapping.Body(*createTokenProcessorToTokenGeneratorMapping)
	tokenProcessorToTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.CreateTokenToTokenMappingExecute(apiCreateTokenProcessorToTokenGeneratorMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the TokenProcessorToTokenGeneratorMapping", err, httpResp)
		return
	}

	// Read the response into the state
	var state tokenProcessorToTokenGeneratorMappingResourceModel

	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, tokenProcessorToTokenGeneratorMappingsResponse, &state, plan)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenProcessorToTokenGeneratorMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tokenProcessorToTokenGeneratorMappingResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadTokenProcessorToTokenGeneratorMapping, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.GetTokenToTokenMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the TokenProcessorToTokenGeneratorMapping", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  TokenProcessorToTokenGeneratorMapping", err, httpResp)
		}
	}

	// Read the response into the state
	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, apiReadTokenProcessorToTokenGeneratorMapping, &state, state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tokenProcessorToTokenGeneratorMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan tokenProcessorToTokenGeneratorMappingResourceModel
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
	updateTokenProcessorToTokenGeneratorMapping := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.UpdateTokenToTokenMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewTokenToTokenMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TokenProcessorToTokenGeneratorMapping", err.Error())
		return
	}

	updateTokenProcessorToTokenGeneratorMapping = updateTokenProcessorToTokenGeneratorMapping.Body(*createUpdateRequest)
	updateTokenProcessorToTokenGeneratorMappingResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.UpdateTokenToTokenMappingByIdExecute(updateTokenProcessorToTokenGeneratorMapping)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating TokenProcessorToTokenGeneratorMapping", err, httpResp)
		return
	}

	// Read the response
	var state tokenProcessorToTokenGeneratorMappingResourceModel
	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, updateTokenProcessorToTokenGeneratorMappingResponse, &state, plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *tokenProcessorToTokenGeneratorMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tokenProcessorToTokenGeneratorMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.DeleteTokenToTokenMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Token Processor to Token Generator Mapping", err, httpResp)
		return
	}

}

func (r *tokenProcessorToTokenGeneratorMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
