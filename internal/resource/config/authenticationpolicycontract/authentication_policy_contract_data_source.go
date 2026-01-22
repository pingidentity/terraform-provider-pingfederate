// Copyright © 2025 Ping Identity Corporation

package authenticationpolicycontract

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationPolicyContractDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationPolicyContractDataSource{}
)

// AuthenticationPolicyContractDataSource is a helper function to simplify the provider implementation.
func AuthenticationPolicyContractDataSource() datasource.DataSource {
	return &authenticationPolicyContractDataSource{}
}

// authenticationPolicyContractDataSource is the datasource implementation.
type authenticationPolicyContractDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationPolicyContractDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an authentication policy contract.",
		Attributes: map[string]schema.Attribute{
			"core_attributes": schema.SetNestedAttribute{
				Description: "A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
							Optional: false,
						},
					},
				},
			},
			"extended_attributes": schema.SetNestedAttribute{
				Description: "A list of additional attributes as needed.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
							Optional: false,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Description: "The Authentication Policy contract name. Name is unique.",
				Computed:    true,
				Optional:    false,
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"contract_id",
		true,
		"The persistent, unique ID for the authentication policy contract.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *authenticationPolicyContractDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policy_contract"
}

func (r *authenticationPolicyContractDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *authenticationPolicyContractDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationPolicyContractModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPolicyContracts, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.GetAuthenticationPolicyContract(config.AuthContext(ctx, r.providerConfig), state.ContractId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an authentication policy contract", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readAuthenticationPolicyContractsResponse(ctx, apiReadAuthenticationPolicyContracts, &state, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
