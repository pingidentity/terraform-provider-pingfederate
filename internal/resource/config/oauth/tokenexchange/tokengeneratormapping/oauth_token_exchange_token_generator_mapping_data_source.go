package oauthtokenexchangetokengeneratormapping

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

// GetSchema defines the schema for the datasource.
func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the mapping from a token exchange processor policy to a token generator.",
		Attributes: map[string]schema.Attribute{
			"attribute_sources":              attributesources.ToDataSourceSchema(),
			"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
			"issuance_criteria":              issuancecriteria.ToDataSourceSchema(),
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
	id.ToDataSourceSchemaCustomId(
		&schema,
		"id",
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

func (r *oauthTokenExchangeTokenGeneratorMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthTokenExchangeTokenGeneratorMappingModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthTokenExchangeTokenGeneratorMapping, httpResp, err := r.apiClient.OauthTokenExchangeTokenGeneratorMappingsAPI.GetTokenGeneratorMappingById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Token Exchange Token Generator Mapping", err, httpResp)
	}

	// Read the response into the state
	diags = readOauthTokenExchangeTokenGeneratorMappingResponse(ctx, apiReadOauthTokenExchangeTokenGeneratorMapping, &state, state)
	resp.Diagnostics.Append(diags...)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
