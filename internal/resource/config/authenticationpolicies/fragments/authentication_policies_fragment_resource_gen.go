package authenticationpoliciesfragments

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPoliciesFragmentResource{}
	_ resource.ResourceWithConfigure   = &authenticationPoliciesFragmentResource{}
	_ resource.ResourceWithImportState = &authenticationPoliciesFragmentResource{}
)

// AuthenticationPoliciesFragmentResource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesFragmentResource() resource.Resource {
	return &authenticationPoliciesFragmentResource{}
}

// authenticationPoliciesFragmentResource is the resource implementation.
type authenticationPoliciesFragmentResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationPoliciesFragmentModel struct {
	Description types.String `tfsdk:"description"`
	FragmentId  types.String `tfsdk:"fragment_id"`
	Id          types.String `tfsdk:"id"`
	Inputs      types.Object `tfsdk:"inputs"`
	Name        types.String `tfsdk:"name"`
	Outputs     types.Object `tfsdk:"outputs"`
	RootNode    types.Object `tfsdk:"root_node"`
}

func (r *authenticationPoliciesFragmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	contextAttr := schema.StringAttribute{
		Optional:    true,
		Description: "The result context.",
	}
	attributeRulesAttr := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"fallback_to_success": schema.BoolAttribute{
				Optional:    true,
				Description: "When all the rules fail, you may choose to default to the general success action or fail. Default to success.",
			},
			"items": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attribute_name": schema.StringAttribute{
							Optional:    true,
							Description: "The name of the attribute to use in this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
						},
						"attribute_source": sourcetypeidkey.ToSchemaWithDescription(false, "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>"),
						"condition": schema.StringAttribute{
							Optional:    true,
							Description: "The condition that will be applied to the attribute's expected value. This field is required if the Attribute Source type is not 'EXPRESSION'.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"EQUALS",
									"EQUALS_CASE_INSENSITIVE",
									"EQUALS_DN",
									"NOT_EQUAL",
									"NOT_EQUAL_CASE_INSENSITIVE",
									"NOT_EQUAL_DN",
									"MULTIVALUE_CONTAINS",
									"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
									"MULTIVALUE_CONTAINS_DN",
									"MULTIVALUE_DOES_NOT_CONTAIN",
									"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
									"MULTIVALUE_DOES_NOT_CONTAIN_DN",
								),
							},
						},
						"expected_value": schema.StringAttribute{
							Optional:    true,
							Description: "The expected value of this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
						},
						"expression": schema.StringAttribute{
							Optional:    true,
							Description: "The expression of this attribute rule. This field is required if the Attribute Source type is 'EXPRESSION'.",
						},
						"result": schema.StringAttribute{
							Required:    true,
							Description: "The result of this attribute rule.",
						},
					},
				},
				Optional:    true,
				Description: "The actual list of attribute rules.",
			},
		},
		Optional:    true,
		Description: "A collection of attribute rules",
	}
	attributeMappingAttr := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
			"attribute_sources":              attributesources.ToSchema(0),
			"issuance_criteria":              issuancecriteria.ToSchema(),
		},
		Required:    true,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
	schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the authentication policy fragment.",
			},
			"inputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToSchema(),
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The authentication policy fragment name. Name is unique.",
			},
			"outputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToSchema(),
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"root_node": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"apc_mapping_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_mapping": attributeMappingAttr,
							"authentication_policy_contract_ref": schema.SingleNestedAttribute{
								Attributes:  resourcelink.ToSchema(),
								Required:    true,
								Description: "A reference to a resource.",
							},
							"context": contextAttr,
						},
						Optional:    true,
						Description: "An authentication policy contract selection action.",
					},
					"authn_selector_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"authentication_selector_ref": schema.SingleNestedAttribute{
								Attributes:  resourcelink.ToSchema(),
								Required:    true,
								Description: "A reference to a resource.",
							},
							"context": contextAttr,
						},
						Optional:    true,
						Description: "An authentication selector selection action.",
					},
					"authn_source_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_rules": attributeRulesAttr,
							"authentication_source": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"source_ref": schema.SingleNestedAttribute{
										Attributes:  resourcelink.ToSchema(),
										Required:    true,
										Description: "A reference to a resource.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "The type of this authentication source.",
										Validators: []validator.String{
											stringvalidator.OneOf(
												"IDP_ADAPTER",
												"IDP_CONNECTION",
											),
										},
									},
								},
								Required:    true,
								Description: "An authentication source (IdP adapter or IdP connection).",
							},
							"input_user_id_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"source": sourcetypeidkey.ToSchemaWithDescription(false, "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>"),
									"value": schema.StringAttribute{
										Required:    true,
										Description: "The value for this attribute.",
									},
								},
								Optional:    true,
								Description: "Defines how an attribute in an attribute contract should be populated.",
							},
							"user_id_authenticated": schema.BoolAttribute{
								Optional:    true,
								Description: "Indicates whether the user ID obtained by the user ID mapping is authenticated.",
							},
						},
						Optional:    true,
						Description: "An authentication source selection action.",
					},
					"continue_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": contextAttr,
						},
						Optional:    true,
						Description: "The continue selection action.",
					},
					"done_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": contextAttr,
						},
						Optional:    true,
						Description: "The done selection action.",
					},
					"fragment_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_rules": attributeRulesAttr,
							"context":         contextAttr,
							"fragment": schema.SingleNestedAttribute{
								Attributes:  resourcelink.ToSchema(),
								Required:    true,
								Description: "A reference to a resource.",
							},
							"fragment_mapping": attributeMappingAttr,
						},
						Optional:    true,
						Description: "A authentication policy fragment selection action.",
					},
					"local_identity_mapping_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context":         contextAttr,
							"inbound_mapping": attributeMappingAttr,
							"local_identity_ref": schema.SingleNestedAttribute{
								Attributes:  resourcelink.ToSchema(),
								Required:    true,
								Description: "A reference to a resource.",
							},
							"outbound_attribute_mapping": attributeMappingAttr,
						},
						Optional:    true,
						Description: "A local identity profile selection action.",
					},
					"restart_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": contextAttr,
						},
						Optional:    true,
						Description: "The restart selection action.",
					},
				},
				Optional:    true,
				Description: "An authentication policy tree node.",
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "fragment_id", false, false, "The authentication policy fragment ID. ID is unique.")
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *authenticationPoliciesFragmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_fragment"
}

func (r *authenticationPoliciesFragmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPoliciesFragmentResponse(ctx context.Context, r *client.Client, state *authenticationPoliciesFragmentModel) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func addOptionalAuthenticationPoliciesFragmentFields(ctx context.Context, addRequest *client.Client, plan authenticationPoliciesFragmentModel) error {

	return nil

}

func (r *authenticationPoliciesFragmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPoliciesFragmentModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the response into the state
	var state authenticationPoliciesFragmentModel
	//diags = readAuthenticationPoliciesFragmentResponse(ctx, oauthClientResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesFragmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesFragmentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesFragmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan authenticationPoliciesFragmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	//diags = resp.State.Set(ctx, &state)
	//resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesFragmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationPoliciesFragmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPoliciesAPI.DeleteFragment(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.FragmentId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Policy Fragment", err, httpResp)
	}
}

func (r *authenticationPoliciesFragmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to fragment_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("fragment_id"), req, resp)
}
