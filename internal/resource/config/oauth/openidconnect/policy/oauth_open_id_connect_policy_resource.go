package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithConfigure   = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithImportState = &oauthOpenIdConnectPolicyResource{}
)

// OauthOpenIdConnectPolicyResource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPolicyResource() resource.Resource {
	return &oauthOpenIdConnectPolicyResource{}
}

// oauthOpenIdConnectPolicyResource is the resource implementation.
type oauthOpenIdConnectPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthOpenIdConnectPolicyResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	AccessTokenManagerRef       types.Object `tfsdk:"access_token_manager_ref"`
	IdTokenLifetime             types.Int64  `tfsdk:"id_token_lifetime"`
	IncludeSriInIdToken         types.Bool   `tfsdk:"include_sri_in_id_token"`
	IncludeUserInfoInIdToken    types.Bool   `tfsdk:"include_user_info_in_id_token"`
	IncludeSHashInIdToken       types.Bool   `tfsdk:"include_s_hash_in_id_token"`
	ReturnIdTokenOnRefreshGrant types.Bool   `tfsdk:"return_id_token_on_refresh_grant"`
	ReissueIdTokenInHybridFlow  types.Bool   `tfsdk:"reissue_id_token_in_hybrid_flow"`
	AttributeContract           types.Object `tfsdk:"attribute_contract"`
	AttributeMapping            types.Object `tfsdk:"attribute_mapping"`
	ScopeAttributeMappings      types.Object `tfsdk:"scope_attribute_mappings"`
}

// GetSchema defines the schema for the resource.
func (r *oauthOpenIdConnectPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OpenID Connect Policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The policy ID used internally.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name used for display in UI screens.",
				Required:    true,
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The access token manager associated with this Open ID Connect policy.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"id_token_lifetime": schema.Int64Attribute{
				Description: "The ID Token Lifetime, in minutes. The default value is 5.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"include_sri_in_id_token": schema.BoolAttribute{
				Description: "Determines whether a Session Reference Identifier is included in the ID token.",
				Optional:    true,
			},
			"include_user_info_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the User Info is always included in the ID token",
				Optional:    true,
			},
			"include_s_hash_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the State Hash should be included in the ID token.",
				Optional:    true,
			},
			"return_id_token_on_refresh_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when refresh grant is requested or not.",
				Optional:    true,
			},
			"reissue_id_token_in_hybrid_flow": schema.BoolAttribute{
				Description: "Determines whether a new ID Token should be returned during token request of the hybrid flow.",
				Optional:    true,
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be returned to OAuth clients in response to requests received at the PingFederate UserInfo endpoint.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of read-only attributes (for example, sub) that are automatically populated by PingFederate.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Optional:    true,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Optional:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
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
									Required:    true,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Optional:    true,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Optional:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
			//TODO common?
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"attribute_sources":              attributesources.ToSchema(),
					"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
					"issuance_criteria":              issuancecriteria.ToSchema(),
				},
			},
			"scope_attribute_mappings": schema.MapNestedAttribute{
				Description: "The attribute scope mappings from scopes to attribute names.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.ListAttribute{
							Description: "A List of values.",
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
	resp.Schema = schema
}

func addOptionalOauthOpenIdConnectPolicyFields(ctx context.Context, plan oauthOpenIdConnectPolicyResourceModel) error {
	return nil
}

// Metadata returns the resource type name.
func (r *oauthOpenIdConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policy"
}

func (r *oauthOpenIdConnectPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readOauthOpenIdConnectPolicyResponse(ctx context.Context, state *oauthOpenIdConnectPolicyResourceModel) {
}

func (r *oauthOpenIdConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *oauthOpenIdConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *oauthOpenIdConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *oauthOpenIdConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthOpenIdConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
