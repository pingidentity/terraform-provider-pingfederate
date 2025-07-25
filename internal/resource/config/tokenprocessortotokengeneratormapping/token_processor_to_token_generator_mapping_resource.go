// Copyright © 2025 Ping Identity Corporation

package tokenprocessortotokengeneratormapping

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
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

// GetSchema defines the schema for the resource.
func (r *tokenProcessorToTokenGeneratorMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the mapping from token processor to a token generator.",
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
			"attribute_sources":              attributesources.ToSchema(0, false),
			"default_target_resource": schema.StringAttribute{
				Description: "Default target URL for this Token Processor to Token Generator mapping configuration.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Processor. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"issuance_criteria": issuancecriteria.ToSchema(),
			"mapping_id": schema.StringAttribute{
				Description: "The id of the Token Processor to Token Generator Mapping.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx context.Context, addRequest *client.TokenToTokenMapping, plan tokenProcessorToTokenGeneratorMappingModel) error {
	if internaltypes.IsDefined(plan.AttributeSources) {
		addRequest.AttributeSources = []client.AttributeSourceAggregation{}
		addRequest.AttributeSources = attributesources.ClientStruct(plan.AttributeSources)
	}

	if internaltypes.IsDefined(plan.IssuanceCriteria) {
		addRequest.IssuanceCriteria = issuancecriteria.ClientStruct(plan.IssuanceCriteria)
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

func (r *tokenProcessorToTokenGeneratorMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tokenProcessorToTokenGeneratorMappingModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributeContractFulfillment := attributecontractfulfillment.ClientStruct(plan.AttributeContractFulfillment)
	createTokenProcessorToTokenGeneratorMapping := client.NewTokenToTokenMapping(attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx, createTokenProcessorToTokenGeneratorMapping, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Token Processor to Token Generator Mapping: "+err.Error())
		return
	}

	apiCreateTokenProcessorToTokenGeneratorMapping := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.CreateTokenToTokenMapping(config.AuthContext(ctx, r.providerConfig))
	apiCreateTokenProcessorToTokenGeneratorMapping = apiCreateTokenProcessorToTokenGeneratorMapping.Body(*createTokenProcessorToTokenGeneratorMapping)
	tokenProcessorToTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.CreateTokenToTokenMappingExecute(apiCreateTokenProcessorToTokenGeneratorMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Token Processor to Token Generator Mapping", err, httpResp)
		return
	}

	// Read the response into the state
	var state tokenProcessorToTokenGeneratorMappingModel

	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, tokenProcessorToTokenGeneratorMappingsResponse, &state, plan)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenProcessorToTokenGeneratorMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tokenProcessorToTokenGeneratorMappingModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadTokenProcessorToTokenGeneratorMapping, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.GetTokenToTokenMappingById(config.AuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Token Processor To Token Generator Mapping", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Token Processor to Token Generator Mapping", err, httpResp)
		}
		return
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

	var plan tokenProcessorToTokenGeneratorMappingModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributeContractFulfillment := &map[string]client.AttributeFulfillmentValue{}
	attributeContractFulfillmentErr := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContractFulfillment, false)), attributeContractFulfillment)
	if attributeContractFulfillmentErr != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to build attribute contract fulfillment request object: "+attributeContractFulfillmentErr.Error())
		return
	}
	updateTokenProcessorToTokenGeneratorMapping := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.UpdateTokenToTokenMappingById(config.AuthContext(ctx, r.providerConfig), plan.MappingId.ValueString())
	createUpdateRequest := client.NewTokenToTokenMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalTokenProcessorToTokenGeneratorMappingFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Token Processor to Token Generator Mapping: "+err.Error())
		return
	}

	updateTokenProcessorToTokenGeneratorMapping = updateTokenProcessorToTokenGeneratorMapping.Body(*createUpdateRequest)
	updateTokenProcessorToTokenGeneratorMappingResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.UpdateTokenToTokenMappingByIdExecute(updateTokenProcessorToTokenGeneratorMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Token Processor to Token Generator Mapping", err, httpResp)
		return
	}

	// Read the response
	var state tokenProcessorToTokenGeneratorMappingModel
	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, updateTokenProcessorToTokenGeneratorMappingResponse, &state, plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenProcessorToTokenGeneratorMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tokenProcessorToTokenGeneratorMappingModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.DeleteTokenToTokenMappingById(config.AuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Token Processor to Token Generator Mapping", err, httpResp)
	}
}

func (r *tokenProcessorToTokenGeneratorMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("mapping_id"), req, resp)
}
