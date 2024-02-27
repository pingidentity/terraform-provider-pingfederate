package spauthenticationpolicycontractmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	datasourceattributecontractfulfillment "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	datasourceattributesources "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	datasourceissuancecriteria "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &spAuthenticationPolicyContractMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &spAuthenticationPolicyContractMappingDataSource{}
)

// SpAuthenticationPolicyContractMappingDataSource is a helper function to simplify the provider implementation.
func SpAuthenticationPolicyContractMappingDataSource() datasource.DataSource {
	return &spAuthenticationPolicyContractMappingDataSource{}
}

// spAuthenticationPolicyContractMappingDataSource is the datasource implementation.
type spAuthenticationPolicyContractMappingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type spAuthenticationPolicyContractMappingDataSourceModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	MappingId                        types.String `tfsdk:"mapping_id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	DefaultTargetResource            types.String `tfsdk:"default_target_resource"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the datasource.
func (r *spAuthenticationPolicyContractMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the mapping from an Authentication Policy Contract (APC) to a Service Provider (SP).",
		Attributes: map[string]schema.Attribute{
			"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
			"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
			"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
			"source_id": schema.StringAttribute{
				Description: "The id of the Authentication Policy Contract.",
				Computed:    true,
				Optional:    false,
			},
			"default_target_resource": schema.StringAttribute{
				Description: "Default target URL for this APC-to-adapter mapping configuration.",
				Computed:    true,
				Optional:    false,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the SP Adapter.",
				Computed:    true,
				Optional:    false,
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group",
				Computed:    true,
				Optional:    false,
			},
		},
	}
	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"mapping_id",
		true,
		"The id of the APC-to-SP Adapter mapping.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *spAuthenticationPolicyContractMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_authentication_policy_contract_mapping"
}

func (r *spAuthenticationPolicyContractMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readSpAuthenticationPolicyContractMappingDataSourceResponse(ctx context.Context, r *client.ApcToSpAdapterMapping, state *spAuthenticationPolicyContractMappingDataSourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.AttributeSources, respDiags = attributesources.ToState(ctx, r.AttributeSources, true)
	diags.Append(respDiags...)
	state.AttributeContractFulfillment, respDiags = attributecontractfulfillment.ToState(ctx, r.AttributeContractFulfillment)
	diags.Append(respDiags...)
	state.IssuanceCriteria, respDiags = issuancecriteria.ToState(ctx, r.IssuanceCriteria)
	diags.Append(respDiags...)
	state.SourceId = types.StringValue(r.SourceId)
	state.TargetId = types.StringValue(r.TargetId)
	state.MappingId = types.StringPointerValue(r.Id)
	state.Id = types.StringPointerValue(r.Id)
	state.DefaultTargetResource = types.StringPointerValue(r.DefaultTargetResource)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
	return diags
}

func (r *spAuthenticationPolicyContractMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spAuthenticationPolicyContractMappingDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSpAuthenticationPolicyContractMappingResource, httpResp, err := r.apiClient.SpAuthenticationPolicyContractMappingsAPI.GetApcToSpAdapterMappingById(config.AuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the SP Authentication Policy Contract Mapping Resource", err, httpResp)
	}

	// Read the response into the state
	diags = readSpAuthenticationPolicyContractMappingDataSourceResponse(ctx, apiReadSpAuthenticationPolicyContractMappingResource, &state)
	resp.Diagnostics.Append(diags...)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
