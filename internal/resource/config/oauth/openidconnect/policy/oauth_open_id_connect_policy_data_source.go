package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthOpenIdConnectPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthOpenIdConnectPolicyDataSource{}
)

// OauthOpenIdConnectPolicyDataSource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPolicyDataSource() datasource.DataSource {
	return &oauthOpenIdConnectPolicyDataSource{}
}

// oauthOpenIdConnectPolicyDataSource is the datasource implementation.
type oauthOpenIdConnectPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *oauthOpenIdConnectPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an OpenID Connect Policy.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name used for display in UI screens.",
				Computed:    true,
				Optional:    false,
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The access token manager associated with this Open ID Connect policy.",
				Computed:    true,
				Optional:    false,
				Attributes:  resourcelink.ToDataSourceSchema(),
			},
			"id_token_lifetime": schema.Int64Attribute{
				Description: "The ID Token Lifetime, in minutes. The default value is 5.",
				Computed:    true,
				Optional:    false,
			},
			"include_sri_in_id_token": schema.BoolAttribute{
				Description: "Determines whether a Session Reference Identifier is included in the ID token.",
				Optional:    true,
				Computed:    true,
			},
			"include_user_info_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the User Info is always included in the ID token",
				Computed:    true,
				Optional:    false,
			},
			"include_s_hash_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the State Hash should be included in the ID token.",
				Computed:    true,
				Optional:    false,
			},
			"return_id_token_on_refresh_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when refresh grant is requested or not.",
				Computed:    true,
				Optional:    false,
			},
			"reissue_id_token_in_hybrid_flow": schema.BoolAttribute{
				Description: "Determines whether a new ID Token should be returned during token request of the hybrid flow.",
				Computed:    true,
				Optional:    false,
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be returned to OAuth clients in response to requests received at the PingFederate UserInfo endpoint.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of read-only attributes (for example, sub) that are automatically populated by PingFederate.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Computed:    true,
									Optional:    false,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Computed:    true,
									Optional:    false,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional attributes.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Computed:    true,
									Optional:    false,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Computed:    true,
									Optional:    false,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"attribute_sources":              attributesources.ToDataSourceSchema(),
					"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
					"issuance_criteria":              issuancecriteria.ToDataSourceSchema(),
				},
			},
			"scope_attribute_mappings": schema.MapNestedAttribute{
				Description: "The attribute scope mappings from scopes to attribute names.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.ListAttribute{
							Description: "A List of values.",
							Computed:    true,
							Optional:    false,
							ElementType: types.StringType,
						},
					},
				},
			},
			"include_x5t_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the X.509 thumbprint header should be included in the ID Token. Supported in PF version 11.3 or later.",
				Optional:    false,
				Computed:    true,
			},
			"id_token_typ_header_value": schema.StringAttribute{
				Description: "ID Token Type (typ) Header Value. Supported in PF version 11.3 or later.",
				Optional:    false,
				Computed:    true,
			},
			"last_modified": schema.StringAttribute{
				Description: "The time at which the policy was last changed. This property is read only and is ignored on PUT and POST requests. Supported in PF version 12.0 or later.",
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"policy_id", true, "The policy ID used internally.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *oauthOpenIdConnectPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policy"
}

func (r *oauthOpenIdConnectPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read the datasource information
func (r *oauthOpenIdConnectPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthOpenIdConnectPolicyModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOIDCPolicy, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.GetOIDCPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OIDC Policy", err, httpResp)
		return
	}

	// Read the response into the state
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, apiReadOIDCPolicy, &state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
