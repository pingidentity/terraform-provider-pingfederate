package oauthtokenexchangetokengeneratormapping

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
	_ datasource.DataSource              = &oauthTokenExchangeTokenGeneratorMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthTokenExchangeTokenGeneratorMappingDataSource{}
)

// OauthTokenExchangeTokenGeneratorMappingDataSource is a helper function to simplify the provider implementation.
func OauthTokenExchangeTokenGeneratorMappingDataSource() datasource.DataSource {
	return &oauthTokenExchangeTokenGeneratorMappingDataSource{}
}

// oauthTokenExchangeTokenGeneratorMappingDataSource is the datasource implementation.
type oauthTokenExchangeTokenGeneratorMappingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthTokenExchangeTokenGeneratorMappingDataSourceModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	MappingId                        types.String `tfsdk:"mapping_id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the datasource.
func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the mapping from a token exchange processor policy to a token generator.",
		Attributes: map[string]schema.Attribute{
			"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
			"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
			"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Exchange Processor policy.",
				Computed:    true,
				Optional:    false,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator",
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
	id.ToDataSourceSchemaCustomId(
		&schema,
		"mapping_id",
		true,
		"The id of the Token Exchange Processor policy to Token Generator mapping.",
	)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token_exchange_token_generator_mapping"
}

func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readOauthTokenExchangeTokenGeneratorMappingDataSourceResponse(ctx context.Context, r *client.ProcessorPolicyToGeneratorMapping, state *oauthTokenExchangeTokenGeneratorMappingDataSourceModel) diag.Diagnostics {
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
	state.MappingId = types.StringPointerValue(r.Id)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
	return diags
}

func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthTokenExchangeTokenGeneratorMappingDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthTokenExchangeTokenGeneratorMapping, httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.GetTokenGeneratorMappingById(config.DetermineAuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Token Exchange Token Generator Mapping", err, httpResp)
	}

	// Read the response into the state
	diags = readOauthTokenExchangeTokenGeneratorMappingDataSourceResponse(ctx, apiReadOauthTokenExchangeTokenGeneratorMapping, &state)
	resp.Diagnostics.Append(diags...)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
