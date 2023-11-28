package tokenprocessortotokengeneratormapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tokenProcessorToTokenGeneratorMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &tokenProcessorToTokenGeneratorMappingDataSource{}
)

// TokenProcessorToTokenGeneratorMappingDataSource is a helper function to simplify the provider implementation.
func TokenProcessorToTokenGeneratorMappingDataSource() datasource.DataSource {
	return &tokenProcessorToTokenGeneratorMappingDataSource{}
}

// tokenProcessorToTokenGeneratorMappingDataSource is the datasource implementation.
type tokenProcessorToTokenGeneratorMappingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *tokenProcessorToTokenGeneratorMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the mapping from token processor to a token generator.",
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
			"attribute_sources":              attributesources.ToDataSourceSchema(),
			"default_target_resource": schema.StringAttribute{
				Description: "Default target URL for this Token Processor to Token Generator mapping configuration.",
				Computed:    true,
				Optional:    false,
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group.",
				Computed:    true,
				Optional:    false,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator.",
				Computed:    true,
				Optional:    false,
			},
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Processor.",
				Computed:    true,
				Optional:    false,
			},
			"issuance_criteria": issuancecriteria.ToDataSourceSchema(),
		},
	}
	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema, "mapping_id", true, "ID of Token Processor to Token Generator Mapping.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *tokenProcessorToTokenGeneratorMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_processor_to_token_generator_mapping"
}

func (r *tokenProcessorToTokenGeneratorMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *tokenProcessorToTokenGeneratorMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tokenProcessorToTokenGeneratorMappingModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadTokenProcessorToTokenGeneratorMapping, httpResp, err := r.apiClient.TokenProcessorToTokenGeneratorMappingsAPI.GetTokenToTokenMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Token Processor To Token Generator Mapping", err, httpResp)
	}

	// Read the response into the state
	diags = readTokenProcessorToTokenGeneratorMappingResponse(ctx, apiReadTokenProcessorToTokenGeneratorMapping, &state, state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
