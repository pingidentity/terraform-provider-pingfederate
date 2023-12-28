package spauthenticationpolicycontractmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &spAuthenticationPolicyContractMappingResource{}
	_ resource.ResourceWithConfigure   = &spAuthenticationPolicyContractMappingResource{}
	_ resource.ResourceWithImportState = &spAuthenticationPolicyContractMappingResource{}
)

// SpAuthenticationPolicyContractMappingResource is a helper function to simplify the provider implementation.
func SpAuthenticationPolicyContractMappingResource() resource.Resource {
	return &spAuthenticationPolicyContractMappingResource{}
}

// spAuthenticationPolicyContractMappingResource is the resource implementation.
type spAuthenticationPolicyContractMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type spAuthenticationPolicyContractMappingResourceModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	DefaultTargetResource            types.String `tfsdk:"default_target_resource"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the resource.
func (r *spAuthenticationPolicyContractMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the mapping from an Authentication Policy Contract (APC) to a Service Provider (SP).",
		Attributes: map[string]schema.Attribute{
			"attribute_sources":              attributesources.ToSchema(0),
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
			"issuance_criteria":              issuancecriteria.ToSchema(),
			"source_id": schema.StringAttribute{
				Description: "The id of the Authentication Policy Contract.",
				Required:    true,
			},
			"default_target_resource": schema.StringAttribute{
				Description: "Default target URL for this APC-to-adapter mapping configuration.",
				Optional:    true,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the SP Adapter.",
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

func addOptionalSpAuthenticationPolicyContractMappingResourceFields(ctx context.Context, addRequest *client.ApcToSpAdapterMapping, plan spAuthenticationPolicyContractMappingResourceModel) error {
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
func (r *spAuthenticationPolicyContractMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_authentication_policy_contract_mapping"
}

func (r *spAuthenticationPolicyContractMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readSpAuthenticationPolicyContractMappingResourceResponse(ctx context.Context, r *client.ApcToSpAdapterMapping, state *spAuthenticationPolicyContractMappingResourceModel) diag.Diagnostics {
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

func (r *spAuthenticationPolicyContractMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan spAuthenticationPolicyContractMappingResourceModel

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
	createSpAuthenticationPolicyContractMappingResource := client.NewApcToSpAdapterMapping(attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err = addOptionalSpAuthenticationPolicyContractMappingResourceFields(ctx, createSpAuthenticationPolicyContractMappingResource, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the SP Authentication Policy Contract Mapping Resource", err.Error())
		return
	}

	apiCreateSpAuthenticationPolicyContractMappingResource := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.CreateApcToSpAdapterMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSpAuthenticationPolicyContractMappingResource = apiCreateSpAuthenticationPolicyContractMappingResource.Body(*createSpAuthenticationPolicyContractMappingResource)
	spAuthenticationPolicyContractMappingResponse, httpResp, err := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.CreateApcToSpAdapterMappingExecute(apiCreateSpAuthenticationPolicyContractMappingResource)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the SP Authentication Policy Contract Mapping Resource", err, httpResp)
		return
	}

	// Read the response into the state
	var state spAuthenticationPolicyContractMappingResourceModel

	diags = readSpAuthenticationPolicyContractMappingResourceResponse(ctx, spAuthenticationPolicyContractMappingResponse, &state)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *spAuthenticationPolicyContractMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state spAuthenticationPolicyContractMappingResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSpAuthenticationPolicyContractMappingResource, httpResp, err := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.GetApcToSpAdapterMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the SP Authentication Policy Contract Mapping Resource", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the SP Authentication Policy Contract Mapping Resource", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readSpAuthenticationPolicyContractMappingResourceResponse(ctx, apiReadSpAuthenticationPolicyContractMappingResource, &state)
	resp.Diagnostics.Append(diags...)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *spAuthenticationPolicyContractMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan spAuthenticationPolicyContractMappingResourceModel
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
	updateSpAuthenticationPolicyContractMappingResource := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.UpdateApcToSpAdapterMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewApcToSpAdapterMapping(attributeContractFulfillment, plan.SourceId.ValueString(), plan.TargetId.ValueString())
	err = addOptionalSpAuthenticationPolicyContractMappingResourceFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for SP Authentication Policy Contract Mapping Resource", err.Error())
		return
	}
	updateSpAuthenticationPolicyContractMappingResource = updateSpAuthenticationPolicyContractMappingResource.Body(*createUpdateRequest)
	updateSpAuthenticationPolicyContractMappingResourceResponse, httpResp, err := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.UpdateApcToSpAdapterMappingByIdExecute(updateSpAuthenticationPolicyContractMappingResource)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the SP Authentication Policy Contract Mapping Resource", err, httpResp)
		return
	}
	// Read the response
	var state spAuthenticationPolicyContractMappingResourceModel
	diags = readSpAuthenticationPolicyContractMappingResourceResponse(ctx, updateSpAuthenticationPolicyContractMappingResourceResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *spAuthenticationPolicyContractMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state spAuthenticationPolicyContractMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.DeleteApcToSpAdapterMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the SP Authentication Policy Contract Mapping Resource", err, httpResp)
	}

}
func (r *spAuthenticationPolicyContractMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
