package config

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &tokenProcessorToTokenGeneratorMappingsResource{}
	_ resource.ResourceWithConfigure   = &tokenProcessorToTokenGeneratorMappingsResource{}
	_ resource.ResourceWithImportState = &tokenProcessorToTokenGeneratorMappingsResource{}
)

// TokenProcessorToTokenGeneratorMappingsResource is a helper function to simplify the provider implementation.
func TokenProcessorToTokenGeneratorMappingsResource() resource.Resource {
	return &tokenProcessorToTokenGeneratorMappingsResource{}
}

// tokenProcessorToTokenGeneratorMappingsResource is the resource implementation.
type tokenProcessorToTokenGeneratorMappingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type tokenProcessorToTokenGeneratorMappingsResourceModel struct {
	// AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	Id                               types.String `tfsdk:"id"`
	DefaultTargetResource            types.String `tfsdk:"default_target_resource"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the resource.
func (r *tokenProcessorToTokenGeneratorMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Token Processor To Token Generator Mappings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the Token Processor to Token Generator mapping. This field is read-only and is ignored when passed in with the payload.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attribute_contract_fulfillment": attributecontractfulfillment.AttributeContractFulfillmentSchema(),
			// "attribute_sources":              common.AttributeSourcesSchema(),
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
			"issuance_criteria": issuancecriteria.IssuanceCriteriaSchema(),
		},
	}
}

func addOptionalTokenProcessorToTokenGeneratorMappingsFields(ctx context.Context, addRequest *client.TokenToTokenMapping, plan tokenProcessorToTokenGeneratorMappingsResourceModel) error {

	// if internaltypes.IsDefined(plan.AttributeSources) {
	// 	addRequest.AttributeSources = []client.AttributeSource{}
	// 	for _, attrSource := range plan.AttributeSources.Elements() {
	// 		// if attrSource
	// 		// addRequest.AttributeSources = append(addRequest.AttributeSources, attrSource)
	// 		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeSources, false)), attrSource)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	if internaltypes.IsDefined(plan.IssuanceCriteria) {
		addRequest.IssuanceCriteria = client.NewIssuanceCriteria()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.IssuanceCriteria, true)), addRequest.IssuanceCriteria)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
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
func (r *tokenProcessorToTokenGeneratorMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_processor_to_token_generator_mappings"
}

func (r *tokenProcessorToTokenGeneratorMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readTokenProcessorToTokenGeneratorMappingsResponse(ctx context.Context, r *client.TokenToTokenMapping, state *tokenProcessorToTokenGeneratorMappingsResourceModel) {
	// state.AttributeSources = internaltypes.GetCorrectMethodFromInternalTypesForThis(r.AttributeSources)
	state.AttributeContractFulfillment = attributecontractfulfillment.AttributeContractFulfillmentToState(ctx, r.AttributeContractFulfillment)
	state.IssuanceCriteria = issuancecriteria.IssuanceCriteriaToState(ctx, r.IssuanceCriteria)
	state.SourceId = types.StringValue(r.SourceId)
	state.TargetId = types.StringValue(r.TargetId)
	state.Id = types.StringPointerValue(r.Id)
	state.DefaultTargetResource = types.StringPointerValue(r.DefaultTargetResource)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
}

func (r *tokenProcessorToTokenGeneratorMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tokenProcessorToTokenGeneratorMappingsResourceModel

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
	createTokenProcessorToTokenGeneratorMappings := client.NewTokenToTokenMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalTokenProcessorToTokenGeneratorMappingsFields(ctx, createTokenProcessorToTokenGeneratorMappings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TokenProcessorToTokenGeneratorMappings", err.Error())
		return
	}
	requestJson, err := createTokenProcessorToTokenGeneratorMappings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateTokenProcessorToTokenGeneratorMappings := r.apiClient.TokenProcessorToTokenGeneratorMappingsApi.CreateTokenToTokenMapping(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateTokenProcessorToTokenGeneratorMappings = apiCreateTokenProcessorToTokenGeneratorMappings.Body(*createTokenProcessorToTokenGeneratorMappings)
	tokenProcessorToTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsApi.CreateTokenToTokenMappingExecute(apiCreateTokenProcessorToTokenGeneratorMappings)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the TokenProcessorToTokenGeneratorMappings", err, httpResp)
		return
	}
	responseJson, err := tokenProcessorToTokenGeneratorMappingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state tokenProcessorToTokenGeneratorMappingsResourceModel

	readTokenProcessorToTokenGeneratorMappingsResponse(ctx, tokenProcessorToTokenGeneratorMappingsResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenProcessorToTokenGeneratorMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tokenProcessorToTokenGeneratorMappingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadTokenProcessorToTokenGeneratorMappings, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsApi.GetTokenToTokenMappingById(ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the TokenProcessorToTokenGeneratorMappings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  TokenProcessorToTokenGeneratorMappings", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadTokenProcessorToTokenGeneratorMappings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTokenProcessorToTokenGeneratorMappingsResponse(ctx, apiReadTokenProcessorToTokenGeneratorMappings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tokenProcessorToTokenGeneratorMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan tokenProcessorToTokenGeneratorMappingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state tokenProcessorToTokenGeneratorMappingsResourceModel
	req.State.Get(ctx, &state)

	attributeContractFulfillment := &map[string]client.AttributeFulfillmentValue{}
	attributeContractFulfillmentErr := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContractFulfillment, false)), attributeContractFulfillment)
	if attributeContractFulfillmentErr != nil {
		resp.Diagnostics.AddError("Failed to build attribute contract fulfillment request object:", attributeContractFulfillmentErr.Error())
		return
	}
	updateTokenProcessorToTokenGeneratorMappings := r.apiClient.TokenProcessorToTokenGeneratorMappingsApi.UpdateTokenToTokenMappingById(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewTokenToTokenMapping(*attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err := addOptionalTokenProcessorToTokenGeneratorMappingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TokenProcessorToTokenGeneratorMappings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateTokenProcessorToTokenGeneratorMappings = updateTokenProcessorToTokenGeneratorMappings.Body(*createUpdateRequest)
	updateTokenProcessorToTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsApi.UpdateTokenToTokenMappingByIdExecute(updateTokenProcessorToTokenGeneratorMappings)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating TokenProcessorToTokenGeneratorMappings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateTokenProcessorToTokenGeneratorMappingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readTokenProcessorToTokenGeneratorMappingsResponse(ctx, updateTokenProcessorToTokenGeneratorMappingsResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *tokenProcessorToTokenGeneratorMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *tokenProcessorToTokenGeneratorMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
