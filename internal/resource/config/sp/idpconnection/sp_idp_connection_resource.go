package spidpconnection

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/datastorerepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/inboundprovisioninguserrepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &spIdpConnectionResource{}
	_ resource.ResourceWithConfigure   = &spIdpConnectionResource{}
	_ resource.ResourceWithImportState = &spIdpConnectionResource{}

	metadataReloadSettingsAttrTypes = map[string]attr.Type{
		"metadata_url_ref":            types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"enable_auto_metadata_update": types.BoolType,
	}

	oidcClientCredentialsAttrTypes = map[string]attr.Type{
		"client_id":     types.StringType,
		"client_secret": types.StringType,
	}

	certViewAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                        types.StringType,
			"serial_number":             types.StringType,
			"subject_dn":                types.StringType,
			"subject_alternative_names": types.SetType{ElemType: types.StringType},
			"issuer_dn":                 types.StringType,
			"valid_from":                types.StringType,
			"expires":                   types.StringType,
			"key_algorithm":             types.StringType,
			"key_size":                  types.Int64Type,
			"signature_algorithm":       types.StringType,
			"version":                   types.Int64Type,
			"sha1_fingerprint":          types.StringType,
			"sha256_fingerprint":        types.StringType,
			"status":                    types.StringType,
			"crypto_provider":           types.StringType,
		},
	}

	x509fileAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"file_data":       types.StringType,
			"crypto_provider": types.StringType,
		},
	}

	connectionCertAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cert_view":                   types.SetType{ElemType: certViewAttrTypes},
			"x509_file":                   x509fileAttrTypes,
			"active_verification_cert":    types.BoolType,
			"primary_verification_cert":   types.BoolType,
			"secondary_verification_cert": types.BoolType,
			"encryption_cert":             types.BoolType,
		},
	}

	signingSettingsAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"signing_pair_key_ref":         types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"active_signing_key_pair_refs": types.ListType{ElemType: types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()}},
			"algorithm":                    types.StringType,
			"include_cert_in_signature":    types.BoolType,
			"include_raw_key_in_signature": types.BoolType,
		},
	}

	usernamePasswordCredentialsAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"username": types.StringType,
			"password": types.StringType,
		},
	}

	outboundBackChannelAuthAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                   types.StringType,
			"http_basic_credentials": usernamePasswordCredentialsAttrTypes,
			"digital_signature":      types.BoolType,
			"ssl_auth_key_pair_ref":  types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"validate_partner_cert":  types.BoolType,
		},
	}

	inboundBackChannelAuthAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                    types.StringType,
			"http_basic_credentials":  usernamePasswordCredentialsAttrTypes,
			"digital_signature":       types.BoolType,
			"verification_subject_dn": types.StringType,
			"verification_issuer_dn":  types.StringType,
			"certs":                   types.SetType{ElemType: connectionCertAttrTypes},
			"require_ssl":             types.BoolType,
		},
	}

	credentialsAttrTypes = map[string]attr.Type{
		"verification_subject_dn":           types.StringType,
		"verification_issuer_dn":            types.StringType,
		"certs":                             types.SetType{ElemType: connectionCertAttrTypes},
		"block_encryption_algorithm":        types.StringType,
		"key_transport_algorithm":           types.StringType,
		"signing_settings":                  signingSettingsAttrTypes,
		"decryption_key_pair_ref":           types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"secondary_decryption_key_pair_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"outbound_back_channel_auth":        outboundBackChannelAuthAttrTypes,
		"inbound_back_channel_auth":         inboundBackChannelAuthAttrTypes,
	}

	contactInfoAttrTypes = map[string]attr.Type{
		"company":    types.StringType,
		"email":      types.StringType,
		"phone":      types.StringType,
		"first_name": types.StringType,
		"last_name":  types.StringType,
	}

	entityIdAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity_id":          types.StringType,
			"entity_description": types.StringType,
		},
	}

	additionalAllowedEntitiesConfigurationAttrTypes = map[string]attr.Type{
		"allow_additional_entities":   types.BoolType,
		"allow_all_entities":          types.BoolType,
		"additional_allowed_entities": types.SetType{ElemType: entityIdAttrTypes},
	}

	extendedPropertiesElemAttrTypes = map[string]attr.Type{
		"values": types.ListType{ElemType: types.StringType},
	}

	oidcRequestParameterAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":                          types.StringType,
			"attribute_value":               types.ObjectType{AttrTypes: sourcetypeidkey.AttrType()},
			"value":                         types.StringType,
			"application_endpoint_override": types.BoolType,
		},
	}

	oidcProviderSettingsAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"scopes":                                types.StringType,
			"authorization_endpoint":                types.StringType,
			"pushed_authorization_request_endpoint": types.StringType,
			"login_type":                            types.StringType,
			"authentication_scheme":                 types.StringType,
			"authentication_signing_algorithm":      types.StringType,
			"request_signing_algorithm":             types.StringType,
			"enable_pkce":                           types.BoolType,
			"token_endpoint":                        types.StringType,
			"user_info_endpoint":                    types.StringType,
			"logout_endpoint":                       types.StringType,
			"jwks_url":                              types.StringType,
			"track_user_sessions_for_logout":        types.BoolType,
			"request_parameters":                    types.ListType{ElemType: oidcRequestParameterAttrTypes},
			"redirect_uri":                          types.StringType,
			"back_channel_logout_uri":               types.StringType,
			"front_channel_logout_uri":              types.StringType,
			"post_logout_redirect_uri":              types.StringType,
		},
	}

	protocolMessageCustomizationAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"context_name":       types.StringType,
			"message_expression": types.StringType,
		},
	}

	urlWhitelistEntryAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"valid_domain":             types.StringType,
			"valid_path":               types.StringType,
			"allow_query_and_fragment": types.BoolType,
			"require_https":            types.BoolType,
		},
	}

	artifactResolverLocationsAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"index": types.Int64Type,
			"url":   types.StringType,
		},
	}

	artifactSettingsAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"lifetime":           types.Int64Type,
			"resolver_locations": types.ListType{ElemType: artifactResolverLocationsAttrTypes},
			"source_id":          types.StringType,
		},
	}

	sloServiceEndPointAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"binding":      types.StringType,
			"url":          types.StringType,
			"response_url": types.StringType,
		},
	}

	idpSSOServiceEndpointAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"binding": types.StringType,
			"url":     types.StringType,
		},
	}

	authnContextMappingAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"local":  types.StringType,
			"remote": types.StringType,
		},
	}

	decryptonPolicyAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"assertion_encrypted":           types.BoolType,
			"attributes_encrypted":          types.BoolType,
			"subject_name_id_encrypted":     types.BoolType,
			"slo_encrypt_subject_name_id":   types.BoolType,
			"slo_subject_name_id_encrypted": types.BoolType,
		},
	}

	idpBrowserSSOAttributeAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"masked": types.BoolType,
		},
	}

	coreAndExtendedAttributesAttrTypes = map[string]attr.Type{
		"core_attributes":     types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
		"extended_attributes": types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
	}

	idpBrowserSSOAttributeContractAttrTypes = types.ObjectType{
		AttrTypes: coreAndExtendedAttributesAttrTypes,
	}

	spAdapterAttributeAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}

	spAdapterAttributeContractAttrTypes = types.ObjectType{
		AttrTypes: coreAndExtendedAttributesAttrTypes,
	}

	spAdapterTargetApplicationInfoAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"application_name":     types.StringType,
			"application_icon_url": types.StringType,
		},
	}

	spAdapterAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                      types.StringType,
			"name":                    types.StringType,
			"plugin_descriptor_ref":   types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"parent_ref":              types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"configuration":           types.ObjectType{AttrTypes: pluginconfiguration.AttrType()},
			"attribute_contract":      spAdapterAttributeContractAttrTypes,
			"target_application_info": spAdapterTargetApplicationInfoAttrTypes,
		},
	}

	spAdapterMappingAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sp_adapter_ref":                types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"restrict_virtual_entity_ids":   types.BoolType,
			"restricted_virtual_entity_ids": types.ListType{ElemType: types.StringType},
			"restrict_virtual_entity_id":    types.BoolType,
			"adapter_override_settings":     spAdapterAttrTypes,
		},
	}

	authenticationPolicyContractMappingAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"authentication_policy_contract_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
			"restrict_virtual_server_ids":        types.BoolType,
			"restricted_virtual_server_ids":      types.ListType{ElemType: types.StringType},
			"attribute_sources":                  types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrType(false)}},
			"attribute_contract_fulfillment":     attributecontractfulfillment.MapType(),
			"issuance_criteria":                  types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
		},
	}

	ssoOAuthMappingAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrType(false)}},
			"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
			"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
		},
	}

	jitProvisioningUserAttributesAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attribute_contract": types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
			"do_attribute_query": types.BoolType,
		},
	}

	jitProvisioningAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"user_attributes": jitProvisioningUserAttributesAttrTypes,
			"user_repository": types.ObjectType{AttrTypes: datastorerepository.ElemAttrType()},
			"event_trigger":   types.StringType,
			"error_handling":  types.StringType,
		},
	}

	idpBrowserSsoAttrTypes = map[string]attr.Type{
		"protocol":                                 types.StringType,
		"oidc_provider_settings":                   oidcProviderSettingsAttrTypes,
		"enabled_profiles":                         types.ListType{ElemType: types.StringType},
		"incoming_bindings":                        types.ListType{ElemType: types.StringType},
		"message_customizations":                   types.ListType{ElemType: protocolMessageCustomizationAttrTypes},
		"url_whitelist_entries":                    types.ListType{ElemType: urlWhitelistEntryAttrTypes},
		"artifact":                                 artifactSettingsAttrTypes,
		"slo_service_endpoints":                    types.ListType{ElemType: sloServiceEndPointAttrTypes},
		"always_sign_artifact_response":            types.BoolType,
		"sso_application_endpoint":                 types.StringType,
		"sso_service_endpoints":                    types.ListType{ElemType: idpSSOServiceEndpointAttrTypes},
		"default_target_url":                       types.StringType,
		"authn_context_mappings":                   types.ListType{ElemType: authnContextMappingAttrTypes},
		"assertions_signed":                        types.BoolType,
		"sign_authn_requests":                      types.BoolType,
		"decryption_policy":                        decryptonPolicyAttrTypes,
		"idp_identity_mapping":                     types.StringType,
		"attribute_contract":                       idpBrowserSSOAttributeContractAttrTypes,
		"adapter_mapings":                          types.ListType{ElemType: spAdapterMappingAttrTypes},
		"authentication_policy_contract_mappings":  types.ListType{ElemType: authenticationPolicyContractMappingAttrTypes},
		"sso_oauth_mapping":                        ssoOAuthMappingAttrTypes,
		"oauth_authentication_policy_contract_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"jit_provisioning":                         jitProvisioningAttrTypes,
	}

	attributeQueryNameMappingAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"local_name":  types.StringType,
			"remote_name": types.StringType,
		},
	}

	idpAttributeQueryPolicyAttrTypes = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"require_signed_response":     types.BoolType,
			"require_signed_assertion":    types.BoolType,
			"require_encrypted_assertion": types.BoolType,
			"sign_attribute_query":        types.BoolType,
			"encrypt_name_id":             types.BoolType,
			"mask_attribute_values":       types.BoolType,
		},
	}

	attributeQueryAttrTypes = map[string]attr.Type{
		"url":           types.StringType,
		"name_mappings": types.ListType{ElemType: attributeQueryNameMappingAttrTypes},
		"policy":        idpAttributeQueryPolicyAttrTypes,
	}

	accessTokenManagerMappingAttrTypes = map[string]attr.Type{
		"access_token_manager_ref":       types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrType(false)}},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
	}

	idpOAuthGrantAttributeMappingAttrTypes = map[string]attr.Type{
		"access_token_manager_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: accessTokenManagerMappingAttrTypes}},
		"idp_oauth_attribute_contract":  idpBrowserSSOAttributeContractAttrTypes,
	}

	spTokenGeneratorMappingAttrTypes = map[string]attr.Type{
		"sp_token_generator_ref":         types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"restricted_virtual_entity_ids":  types.ListType{ElemType: types.StringType},
		"default_mapping":                types.BoolType,
		"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrType(false)}},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
	}

	wsTrustAttrTypes = map[string]attr.Type{
		"attribute_contract":       types.ObjectType{AttrTypes: accessTokenManagerMappingAttrTypes},
		"generate_local_token":     types.BoolType,
		"token_generator_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: spTokenGeneratorMappingAttrTypes}},
	}

	schemaAttributeAttrTypes = map[string]attr.Type{
		"name":           types.StringType,
		"multi_valued":   types.BoolType,
		"types":          types.ListType{ElemType: types.StringType},
		"sub_attributes": types.ListType{ElemType: types.StringType},
	}

	customSchemaAttrTypes = map[string]attr.Type{
		"namespace":  types.StringType,
		"attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: schemaAttributeAttrTypes}},
	}

	usersAttrTypes = map[string]attr.Type{
		"write_users": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_fulfillment": attributecontractfulfillment.MapType(),
			},
		},
		"read_users": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_contract": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"core_attributes":     types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
						"extended_attributes": types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
					},
				},
				"attributes":            types.ListType{ElemType: spAdapterAttributeAttrTypes},
				"attribute_fulfillment": attributecontractfulfillment.MapType(),
			},
		},
	}

	groupsAttrTypes = map[string]attr.Type{
		"write_groups": attributecontractfulfillment.ObjType(),
		"read_groups": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_contract": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"core_attributes":     types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
						"extended_attributes": types.ListType{ElemType: idpBrowserSSOAttributeAttrTypes},
					},
				},
				"attributes":            types.ListType{ElemType: spAdapterAttributeAttrTypes},
				"attribute_fulfillment": attributecontractfulfillment.MapType(),
			},
		},
	}

	inboundProvisioningAttrTypes = map[string]attr.Type{
		"group_support":    types.BoolType,
		"user_repository":  types.ObjectType{AttrTypes: inboundprovisioninguserrepository.ElemAttrType()},
		"custom_schema":    types.ObjectType{AttrTypes: customSchemaAttrTypes},
		"users":            types.ObjectType{AttrTypes: usersAttrTypes},
		"groups":           types.ObjectType{AttrTypes: groupsAttrTypes},
		"action_on_delete": types.StringType,
	}

	tokenGeneratorAttrTypes = map[string]attr.Type{
		"sp_token_generator_ref":         types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType(false)}},
		"default_mapping":                types.BoolType,
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
	}
)

// SpIdpConnectionResource is a helper function to simplify the provider implementation.
func SpIdpConnectionResource() resource.Resource {
	return &spIdpConnectionResource{}
}

// spIdpConnectionResource is the resource implementation.
type spIdpConnectionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type spIdpConnectionResourceModel struct {
	Active                                 types.Bool   `tfsdk:"active"`
	AdditionalAllowedEntitiesConfiguration types.Object `tfsdk:"additional_allowed_entities_configuration"`
	AttributeQuery                         types.Object `tfsdk:"attribute_query"`
	BaseUrl                                types.String `tfsdk:"base_url"`
	ContactInfo                            types.Object `tfsdk:"contact_info"`
	ConnectionId                           types.String `tfsdk:"connection_id"`
	Credentials                            types.Object `tfsdk:"credentials"`
	DefaultVirtualEntityId                 types.String `tfsdk:"default_virtual_entity_id"`
	EntityId                               types.String `tfsdk:"entity_id"`
	ErrorPageMsgId                         types.String `tfsdk:"error_page_msg_id"`
	ExtendedProperties                     types.Map    `tfsdk:"extended_properties"`
	Id                                     types.String `tfsdk:"id"`
	IdpBrowserSso                          types.Object `tfsdk:"idp_browser_sso"`
	InboundProvisioning                    types.Object `tfsdk:"inbound_provisioning"`
	IdpOAuthGrantAttributeMapping          types.Object `tfsdk:"idp_oauth_grant_attribute_mapping"`
	LicenseConnectionGroup                 types.String `tfsdk:"license_connection_group"`
	LoggingMode                            types.String `tfsdk:"logging_mode"`
	MetadataReloadSettings                 types.Object `tfsdk:"metadata_reload_settings"`
	Name                                   types.String `tfsdk:"name"`
	OidcClientCredentials                  types.Object `tfsdk:"oidc_client_credentials"`
	Type                                   types.String `tfsdk:"type"`
	VirtualEntityIds                       types.Set    `tfsdk:"virtual_entity_ids"`
	WsTrust                                types.Object `tfsdk:"ws_trust"`
}

// GetSchema defines the schema for the resource.
func (r *spIdpConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a SP Idp Connection",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies whether the connection is active and ready to process incoming requests. The default value is false.",
				MarkdownDescription: "Specifies whether the connection is active and ready to process incoming requests. The default value is false.",
				Default:             booldefault.StaticBool(false),
			},
			"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Additional allowed entities or issuers configuration. Currently only used in OIDC IdP (RP) connection.",
				Attributes: map[string]schema.Attribute{
					"allow_additional_entities": schema.BoolAttribute{
						Optional:            true,
						Description:         "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation.",
						MarkdownDescription: "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation.",
					},
					"allow_all_entities": schema.BoolAttribute{
						Optional:    true,
						Description: "Set to true to accept any entity or issuer during entity or issuer validation. (Not Recommended)",
					},
					"additional_allowed_entities": schema.SetNestedAttribute{
						Optional:            true,
						Description:         "An array of additional allowed entities or issuers to be accepted during entity or issuer validation.",
						MarkdownDescription: "An array of additional allowed entities or issuers to be accepted during entity or issuer validation.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"entity_id": schema.StringAttribute{
									Optional:            true,
									Description:         "Unique entity identifier.",
									MarkdownDescription: "Unique entity identifier.",
								},
								"entity_description": schema.StringAttribute{
									Optional:            true,
									Description:         "Entity description.",
									MarkdownDescription: "Entity description.",
								},
							},
						},
					},
				},
			},
			"attribute_query": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"name_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"local_name": schema.StringAttribute{
									Required:            true,
									Description:         "The local attribute name.",
									MarkdownDescription: "The local attribute name.",
								},
								"remote_name": schema.StringAttribute{
									Required:            true,
									Description:         "The remote attribute name as defined by the attribute authority.",
									MarkdownDescription: "The remote attribute name as defined by the attribute authority.",
								},
							},
						},
						Optional:            true,
						Description:         "The attribute name mappings between the SP and the IdP.",
						MarkdownDescription: "The attribute name mappings between the SP and the IdP.",
					},
					"policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_name_id": schema.BoolAttribute{
								Optional:            true,
								Description:         "Encrypt the name identifier.",
								MarkdownDescription: "Encrypt the name identifier.",
							},
							"mask_attribute_values": schema.BoolAttribute{
								Optional:            true,
								Description:         "Mask attributes in log files.",
								MarkdownDescription: "Mask attributes in log files.",
							},
							"require_encrypted_assertion": schema.BoolAttribute{
								Optional:            true,
								Description:         "Require encrypted assertion.",
								MarkdownDescription: "Require encrypted assertion.",
							},
							"require_signed_assertion": schema.BoolAttribute{
								Optional:            true,
								Description:         "Require signed assertion.",
								MarkdownDescription: "Require signed assertion.",
							},
							"require_signed_response": schema.BoolAttribute{
								Optional:            true,
								Description:         "Require signed response.",
								MarkdownDescription: "Require signed response.",
							},
							"sign_attribute_query": schema.BoolAttribute{
								Optional:            true,
								Description:         "Sign the attribute query.",
								MarkdownDescription: "Sign the attribute query.",
							},
						},
						Optional:            true,
						Description:         "The attribute query profile's security policy.",
						MarkdownDescription: "The attribute query profile's security policy.",
					},
					"url": schema.StringAttribute{
						Required:            true,
						Description:         "The URL at your IdP partner's site where attribute queries are to be sent.",
						MarkdownDescription: "The URL at your IdP partner's site where attribute queries are to be sent.",
					},
				},
				Optional:            true,
				Description:         "The attribute query profile supports local applications in requesting user attributes from an attribute authority.",
				MarkdownDescription: "The attribute query profile supports local applications in requesting user attributes from an attribute authority.",
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
						path.Empty().Expression().AtName("attribute_query"),
						path.Empty().Expression().AtName("idp_browser_sso"),
						path.Empty().Expression().AtName("idp_oauth_grant_attribute_mapping"),
						path.Empty().Expression().AtName("inbound_provisioning"),
						path.Empty().Expression().AtName("ws_trust"),
					),
				},
			},
			"base_url": schema.StringAttribute{
				Optional:            true,
				Description:         "The fully-qualified hostname and port on which your partner's federation deployment runs.",
				MarkdownDescription: "The fully-qualified hostname and port on which your partner's federation deployment runs.",
			},
			"contact_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Optional:            true,
						Description:         "Company name.",
						MarkdownDescription: "Company name.",
					},
					"email": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact email address.",
						MarkdownDescription: "Contact email address.",
					},
					"phone": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact phone number.",
						MarkdownDescription: "Contact phone number.",
					},
					"first_name": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact first name.",
						MarkdownDescription: "Contact first name.",
					},
					"last_name": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact last name.",
						MarkdownDescription: "Contact last name.",
					},
				},
				Optional:            true,
				Description:         "Contact information.",
				MarkdownDescription: "Contact information.",
			},
			"credentials": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"verification_issuer_dn": schema.StringAttribute{
						Optional:            true,
						Description:         "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
						MarkdownDescription: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
					},
					"verification_subject_dn": schema.StringAttribute{
						Optional:            true,
						Description:         "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
						MarkdownDescription: "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
					},
					"certs": connectioncert.ToSchema(
						"The certificates used for signature verification and XML encryption.",
						false,
						false,
					),
					"block_encryption_algorithm": schema.StringAttribute{
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"key_transport_algorithm": schema.StringAttribute{
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"signing_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"signing_pair_key_ref": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:            true,
										Description:         "",
										MarkdownDescription: "",
									},
									"type": schema.StringAttribute{
										Optional:            true,
										Description:         "",
										MarkdownDescription: "",
									},
								},
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"active_signing_key_pair_refs": schema.SetAttribute{
								ElementType:         types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"include_cert_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
						},
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"decryption_key_pair_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
						},
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"secondary_decryption_key_pair_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:            true,
								Description:         "",
								MarkdownDescription: "",
							},
						},
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"outbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "",
								MarkdownDescription: "",
								Validators: []validator.String{
									stringvalidator.OneOf("INBOUND", "OUTBOUND"),
								},
							},
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Optional:            true,
										Description:         "",
										MarkdownDescription: "",
									},
									"password": schema.StringAttribute{
										Optional:            true,
										Sensitive:           true,
										Description:         "",
										MarkdownDescription: "",
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
									},
								},
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"ssl_auth_key_pair_ref": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:            true,
										Description:         "",
										MarkdownDescription: "",
									},
								},
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"validate_partner_cert": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "",
								MarkdownDescription: "",
								Default:             booldefault.StaticBool(true),
							},
						},
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "",
								MarkdownDescription: "",
								Validators: []validator.String{
									stringvalidator.OneOf("INBOUND", "OUTBOUND"),
								},
							},
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Required:            true,
										Description:         "",
										MarkdownDescription: "",
									},
									"password": schema.StringAttribute{
										Required:            true,
										Sensitive:           true,
										Description:         "",
										MarkdownDescription: "",
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
									},
								},
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"verification_subject_dn": schema.StringAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"verification_issuer_dn": schema.StringAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"certs": schema.SetAttribute{
								ElementType:         connectionCertAttrTypes,
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
							"require_ssl": schema.BoolAttribute{
								Optional:            true,
								Description:         "",
								MarkdownDescription: "",
							},
						},
						Optional:            true,
						Description:         "",
						MarkdownDescription: "",
					},
				},
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
			},
			"default_virtual_entity_id": schema.StringAttribute{
				Optional:            true,
				Description:         "The default alternate entity ID that identifies the local server to this partner. It is required when virtualEntityIds is not empty and must be included in that list.",
				MarkdownDescription: "The default alternate entity ID that identifies the local server to this partner. It is required when virtualEntityIds is not empty and must be included in that list.",
			},
			"entity_id": schema.StringAttribute{
				Required:            true,
				Description:         "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
				MarkdownDescription: "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
			},
			"error_page_msg_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Identifier that specifies the message displayed on a user-facing error page.",
				MarkdownDescription: "Identifier that specifies the message displayed on a user-facing error page.",
			},
			"extended_properties": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "A List of values",
						},
					},
				},
				Optional:    true,
				Description: "Extended Properties allows to store additional information for IdP/SP Connections. The names of these extended properties should be defined in /extendedProperties.",
			},
			"idp_browser_sso": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"adapter_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"adapter_override_settings": schema.StringAttribute{
									Optional:            true,
									Description:         "An SP adapter instance.",
									MarkdownDescription: "An SP adapter instance.",
								},
								"attribute_contract_fulfillment": schema.MapNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:            true,
														Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key.",
														MarkdownDescription: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:            true,
												Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:            true,
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "A list of configured data stores to look up attributes from.",
									MarkdownDescription: "A list of configured data stores to look up attributes from.",
								},
								"issuance_criteria": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"attribute_name": schema.StringAttribute{
														Required:            true,
														Description:         "The name of the attribute to use in this issuance criterion.",
														MarkdownDescription: "The name of the attribute to use in this issuance criterion.",
													},
													"condition": schema.StringAttribute{
														Required:            true,
														Description:         "The condition that will be applied to the source attribute's value and the expected value.",
														MarkdownDescription: "The condition that will be applied to the source attribute's value and the expected value.",
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
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"source": schema.SingleNestedAttribute{
														Attributes: map[string]schema.Attribute{
															"id": schema.StringAttribute{
																Optional:            true,
																Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															},
															"type": schema.StringAttribute{
																Required:            true,
																Description:         "The source type of this key.",
																MarkdownDescription: "The source type of this key.",
																Validators: []validator.String{
																	stringvalidator.OneOf(
																		"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																		"ACCOUNT_LINK",
																		"ADAPTER",
																		"ASSERTION",
																		"CONTEXT",
																		"CUSTOM_DATA_STORE",
																		"EXPRESSION",
																		"JDBC_DATA_STORE",
																		"LDAP_DATA_STORE",
																		"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																		"MAPPED_ATTRIBUTES",
																		"NO_MAPPING",
																		"TEXT",
																		"TOKEN",
																		"REQUEST",
																		"OAUTH_PERSISTENT_GRANT",
																		"SUBJECT_TOKEN",
																		"ACTOR_TOKEN",
																		"PASSWORD_CREDENTIAL_VALIDATOR",
																		"IDP_CONNECTION",
																		"AUTHENTICATION_POLICY_CONTRACT",
																		"CLAIMS",
																		"LOCAL_IDENTITY_PROFILE",
																		"EXTENDED_CLIENT_METADATA",
																		"EXTENDED_PROPERTIES",
																		"TRACKED_HTTP_PARAMS",
																		"FRAGMENT",
																		"INPUTS",
																		"ATTRIBUTE_QUERY",
																		"IDENTITY_STORE_USER",
																		"IDENTITY_STORE_GROUP",
																		"SCIM_USER",
																		"SCIM_GROUP",
																	),
																},
															},
														},
														Required:            true,
														Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "The expected value of this issuance criterion.",
														MarkdownDescription: "The expected value of this issuance criterion.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											MarkdownDescription: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
										},
										"expression_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"expression": schema.StringAttribute{
														Required:            true,
														Description:         "The OGNL expression to evaluate.",
														MarkdownDescription: "The OGNL expression to evaluate.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											MarkdownDescription: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
										},
									},
									Optional:            true,
									Description:         "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									MarkdownDescription: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
								},
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional: true,

									Description:         "Restricts this mapping to specific virtual entity IDs.",
									MarkdownDescription: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,

									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
								"sp_adapter_ref": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:            true,
											Description:         "The ID of the resource.",
											MarkdownDescription: "The ID of the resource.",
										},
										"location": schema.StringAttribute{
											Optional:            true,
											Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										},
									},
									Required:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of adapters that map to incoming assertions.",
						MarkdownDescription: "A list of adapters that map to incoming assertions.",
					},
					"always_sign_artifact_response": schema.BoolAttribute{
						Optional:            true,
						Description:         "Specify to always sign the SAML ArtifactResponse.",
						MarkdownDescription: "Specify to always sign the SAML ArtifactResponse.",
					},
					"artifact": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"lifetime": schema.Int64Attribute{
								Required:            true,
								Description:         "The lifetime of the artifact in seconds.",
								MarkdownDescription: "The lifetime of the artifact in seconds.",
							},
							"resolver_locations": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"index": schema.Int64Attribute{
											Required:            true,
											Description:         "The priority of the endpoint.",
											MarkdownDescription: "The priority of the endpoint.",
										},
										"url": schema.StringAttribute{
											Required:            true,
											Description:         "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
											MarkdownDescription: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
										},
									},
								},
								Required:            true,
								Description:         "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
								MarkdownDescription: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
							},
							"source_id": schema.StringAttribute{
								Optional:            true,
								Description:         "Source ID for SAML1.x connections",
								MarkdownDescription: "Source ID for SAML1.x connections",
							},
						},
						Optional:            true,
						Description:         "The settings for an Artifact binding.",
						MarkdownDescription: "The settings for an Artifact binding.",
					},
					"assertions_signed": schema.BoolAttribute{
						Optional:            true,
						Description:         "Specify whether the incoming SAML assertions are signed rather than the entire SAML response being signed.",
						MarkdownDescription: "Specify whether the incoming SAML assertions are signed rather than the entire SAML response being signed.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of additional attributes that are present in the incoming assertion.",
								MarkdownDescription: "A list of additional attributes that are present in the incoming assertion.",
							},
						},
						Optional:            true,
						Description:         "A set of user attributes that the IdP sends in the SAML assertion.",
						MarkdownDescription: "A set of user attributes that the IdP sends in the SAML assertion.",
					},
					"authentication_policy_contract_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": schema.MapNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:            true,
														Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key.",
														MarkdownDescription: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:            true,
												Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:            true,
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "A list of configured data stores to look up attributes from.",
									MarkdownDescription: "A list of configured data stores to look up attributes from.",
								},
								"authentication_policy_contract_ref": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:            true,
											Description:         "The ID of the resource.",
											MarkdownDescription: "The ID of the resource.",
										},
										"location": schema.StringAttribute{
											Optional:            true,
											Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										},
									},
									Required:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
								"issuance_criteria": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"attribute_name": schema.StringAttribute{
														Required:            true,
														Description:         "The name of the attribute to use in this issuance criterion.",
														MarkdownDescription: "The name of the attribute to use in this issuance criterion.",
													},
													"condition": schema.StringAttribute{
														Required:            true,
														Description:         "The condition that will be applied to the source attribute's value and the expected value.",
														MarkdownDescription: "The condition that will be applied to the source attribute's value and the expected value.",
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
													"error_result": schema.StringAttribute{
														Optional: true,

														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"source": schema.SingleNestedAttribute{
														Attributes: map[string]schema.Attribute{
															"id": schema.StringAttribute{
																Optional:            true,
																Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															},
															"type": schema.StringAttribute{
																Required:            true,
																Description:         "The source type of this key.",
																MarkdownDescription: "The source type of this key.",
																Validators: []validator.String{
																	stringvalidator.OneOf(
																		"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																		"ACCOUNT_LINK",
																		"ADAPTER",
																		"ASSERTION",
																		"CONTEXT",
																		"CUSTOM_DATA_STORE",
																		"EXPRESSION",
																		"JDBC_DATA_STORE",
																		"LDAP_DATA_STORE",
																		"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																		"MAPPED_ATTRIBUTES",
																		"NO_MAPPING",
																		"TEXT",
																		"TOKEN",
																		"REQUEST",
																		"OAUTH_PERSISTENT_GRANT",
																		"SUBJECT_TOKEN",
																		"ACTOR_TOKEN",
																		"PASSWORD_CREDENTIAL_VALIDATOR",
																		"IDP_CONNECTION",
																		"AUTHENTICATION_POLICY_CONTRACT",
																		"CLAIMS",
																		"LOCAL_IDENTITY_PROFILE",
																		"EXTENDED_CLIENT_METADATA",
																		"EXTENDED_PROPERTIES",
																		"TRACKED_HTTP_PARAMS",
																		"FRAGMENT",
																		"INPUTS",
																		"ATTRIBUTE_QUERY",
																		"IDENTITY_STORE_USER",
																		"IDENTITY_STORE_GROUP",
																		"SCIM_USER",
																		"SCIM_GROUP",
																	),
																},
															},
														},
														Required:            true,
														Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "The expected value of this issuance criterion.",
														MarkdownDescription: "The expected value of this issuance criterion.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											MarkdownDescription: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
										},
										"expression_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"expression": schema.StringAttribute{
														Required:            true,
														Description:         "The OGNL expression to evaluate.",
														MarkdownDescription: "The OGNL expression to evaluate.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											MarkdownDescription: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
										},
									},
									Optional:            true,
									Description:         "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									MarkdownDescription: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
								},
								"restrict_virtual_server_ids": schema.BoolAttribute{
									Optional:            true,
									Description:         "Restricts this mapping to specific virtual entity IDs.",
									MarkdownDescription: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_server_ids": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of Authentication Policy Contracts that map to incoming assertions.",
						MarkdownDescription: "A list of Authentication Policy Contracts that map to incoming assertions.",
					},
					"authn_context_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"local": schema.StringAttribute{
									Optional:            true,
									Description:         "The local authentication context value.",
									MarkdownDescription: "The local authentication context value.",
								},
								"remote": schema.StringAttribute{
									Optional:            true,
									Description:         "The remote authentication context value.",
									MarkdownDescription: "The remote authentication context value.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of authentication context mappings between local and remote values. Applicable for SAML 2.0 and OIDC protocol connections.",
						MarkdownDescription: "A list of authentication context mappings between local and remote values. Applicable for SAML 2.0 and OIDC protocol connections.",
					},
					"decryption_policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"assertion_encrypted": schema.BoolAttribute{
								Optional:            true,
								Description:         "Specify whether the incoming SAML assertion is encrypted for an IdP connection.",
								MarkdownDescription: "Specify whether the incoming SAML assertion is encrypted for an IdP connection.",
							},
							"attributes_encrypted": schema.BoolAttribute{
								Optional:            true,
								Description:         "Specify whether one or more incoming SAML attributes are encrypted for an IdP connection.",
								MarkdownDescription: "Specify whether one or more incoming SAML attributes are encrypted for an IdP connection.",
							},
							"slo_encrypt_subject_name_id": schema.BoolAttribute{
								Optional:            true,
								Description:         "Encrypt the Subject Name ID in SLO messages to the IdP.",
								MarkdownDescription: "Encrypt the Subject Name ID in SLO messages to the IdP.",
							},
							"slo_subject_name_idencrypted": schema.BoolAttribute{
								Optional:            true,
								Description:         "Allow encrypted Subject Name ID in SLO messages from the IdP.",
								MarkdownDescription: "Allow encrypted Subject Name ID in SLO messages from the IdP.",
							},
							"subject_name_id_encrypted": schema.BoolAttribute{
								Optional:            true,
								Description:         "Specify whether the incoming Subject Name ID is encrypted for an IdP connection.",
								MarkdownDescription: "Specify whether the incoming Subject Name ID is encrypted for an IdP connection.",
							},
						},
						Optional:            true,
						Description:         "Defines what to decrypt in the browser-based SSO profile.",
						MarkdownDescription: "Defines what to decrypt in the browser-based SSO profile.",
					},
					"default_target_url": schema.StringAttribute{
						Optional:            true,
						Description:         "The default target URL for this connection. If defined, this overrides the default URL.",
						MarkdownDescription: "The default target URL for this connection. If defined, this overrides the default URL.",
					},
					"enabled_profiles": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						MarkdownDescription: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
					},
					"idp_identity_mapping": schema.StringAttribute{
						Required:            true,
						Description:         "Defines the process in which users authenticated by the IdP are associated with user accounts local to the SP.",
						MarkdownDescription: "Defines the process in which users authenticated by the IdP are associated with user accounts local to the SP.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ACCOUNT_MAPPING",
								"ACCOUNT_LINKING",
								"NONE",
							),
						},
					},
					"incoming_bindings": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections and it is optional for IdP Connections.",
						MarkdownDescription: "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections and it is optional for IdP Connections.",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
					},
					"jit_provisioning": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"error_handling": schema.StringAttribute{
								Optional:            true,
								Description:         "Specify behavior when provisioning request fails. The default is 'CONTINUE_SSO'.",
								MarkdownDescription: "Specify behavior when provisioning request fails. The default is 'CONTINUE_SSO'.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"CONTINUE_SSO",
										"ABORT_SSO",
									),
								},
							},
							"event_trigger": schema.StringAttribute{
								Optional:            true,
								Description:         "Specify when provisioning occurs during assertion processing. The default is 'NEW_USER_ONLY'.",
								MarkdownDescription: "Specify when provisioning occurs during assertion processing. The default is 'NEW_USER_ONLY'.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NEW_USER_ONLY",
										"ALL_SAML_ASSERTIONS",
									),
								},
							},
							"user_attributes": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"masked": schema.BoolAttribute{
													Optional:            true,
													Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
													MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
												},
												"name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
												},
											},
										},
										Optional:            true,
										Description:         "A list of user attributes that the IdP sends in the SAML assertion.",
										MarkdownDescription: "A list of user attributes that the IdP sends in the SAML assertion.",
									},
									"do_attribute_query": schema.BoolAttribute{
										Optional:            true,
										Description:         "Specify whether to use only attributes from the SAML Assertion or retrieve additional attributes from the IdP. The default is false.",
										MarkdownDescription: "Specify whether to use only attributes from the SAML Assertion or retrieve additional attributes from the IdP. The default is false.",
									},
								},
								Required: true,
							},
							"user_repository": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"data_store_ref": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:            true,
												Description:         "The ID of the resource.",
												MarkdownDescription: "The ID of the resource.",
											},
											"location": schema.StringAttribute{
												Optional:            true,
												Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
												MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											},
										},
										Required:            true,
										Description:         "A reference to a resource.",
										MarkdownDescription: "A reference to a resource.",
									},
									"jit_repository_attribute_mapping": schema.MapNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "The user repository attribute mapping.",
										MarkdownDescription: "The user repository attribute mapping.",
									},
									"type": schema.StringAttribute{
										Required:            true,
										Description:         "The data store repository type.",
										MarkdownDescription: "The data store repository type.",
										Validators: []validator.String{
											stringvalidator.OneOf(
												"LDAP",
												"JDBC",
											),
										},
									},
								},
								Required:            true,
								Description:         "Jit Provisioning user repository data store.",
								MarkdownDescription: "Jit Provisioning user repository data store.",
							},
						},
						Optional:            true,
						Description:         "The settings used to specify how and when to provision user accounts.",
						MarkdownDescription: "The settings used to specify how and when to provision user accounts.",
					},
					"message_customizations": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"context_name": schema.StringAttribute{
									Optional:            true,
									Description:         "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be 'assertion', 'authn-response' or 'authn-request'.",
									MarkdownDescription: "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be 'assertion', 'authn-response' or 'authn-request'.",
								},
								"message_expression": schema.StringAttribute{
									Optional:            true,
									Description:         "The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.",
									MarkdownDescription: "The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.",
								},
							},
						},
						Optional:            true,
						Description:         "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
						MarkdownDescription: "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"oauth_authentication_policy_contract_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:            true,
								Description:         "The ID of the resource.",
								MarkdownDescription: "The ID of the resource.",
							},
							"location": schema.StringAttribute{
								Optional: true,

								Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
							},
						},
						Optional:            true,
						Description:         "A reference to a resource.",
						MarkdownDescription: "A reference to a resource.",
					},
					"oidc_provider_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"authentication_scheme": schema.StringAttribute{
								Optional:            true,
								Description:         "The OpenID Connect Authentication Scheme. This is required for Authentication using Code Flow. ",
								MarkdownDescription: "The OpenID Connect Authentication Scheme. This is required for Authentication using Code Flow. ",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"BASIC",
										"POST",
										"PRIVATE_KEY_JWT",
										"CLIENT_SECRET_JWT",
									),
								},
							},
							"authentication_signing_algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "The authentication signing algorithm for token endpoint PRIVATE_KEY_JWT or CLIENT_SECRET_JWT authentication. Asymmetric algorithms are allowed for PRIVATE_KEY_JWT and symmetric algorithms are allowed for CLIENT_SECRET_JWT. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11.",
								MarkdownDescription: "The authentication signing algorithm for token endpoint PRIVATE_KEY_JWT or CLIENT_SECRET_JWT authentication. Asymmetric algorithms are allowed for PRIVATE_KEY_JWT and symmetric algorithms are allowed for CLIENT_SECRET_JWT. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NONE",
										"HS256",
										"HS384",
										"HS512",
										"RS256",
										"RS384",
										"RS512",
										"ES256",
										"ES384",
										"ES512",
										"PS256",
										"PS384",
										"PS512",
									),
								},
							},
							"authorization_endpoint": schema.StringAttribute{
								Required:            true,
								Description:         "URL of the OpenID Provider's OAuth 2.0 Authorization Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's OAuth 2.0 Authorization Endpoint.",
							},
							"back_channel_logout_uri": schema.StringAttribute{
								Optional:            true,
								Description:         "The Back-Channel Logout URI. This read-only parameter is available when user sessions are tracked for logout.",
								MarkdownDescription: "The Back-Channel Logout URI. This read-only parameter is available when user sessions are tracked for logout.",
							},
							"enable_pkce": schema.BoolAttribute{
								Optional:            true,
								Description:         "Enable Proof Key for Code Exchange (PKCE). When enabled, the client sends an SHA-256 code challenge and corresponding code verifier to the OpenID Provider during the authorization code flow.",
								MarkdownDescription: "Enable Proof Key for Code Exchange (PKCE). When enabled, the client sends an SHA-256 code challenge and corresponding code verifier to the OpenID Provider during the authorization code flow.",
							},
							"front_channel_logout_uri": schema.StringAttribute{
								Optional:            true,
								Description:         "The Front-Channel Logout URI. This is a read-only parameter.",
								MarkdownDescription: "The Front-Channel Logout URI. This is a read-only parameter.",
							},
							"jwks_url": schema.StringAttribute{
								Required:            true,
								Description:         "URL of the OpenID Provider's JSON Web Key Set [JWK] document.",
								MarkdownDescription: "URL of the OpenID Provider's JSON Web Key Set [JWK] document.",
							},
							"login_type": schema.StringAttribute{
								Required:            true,
								Description:         "The OpenID Connect login type. These values maps to: <br>  CODE: Authentication using Code Flow <br> POST: Authentication using Form Post <br> POST_AT: Authentication using Form Post with Access Token",
								MarkdownDescription: "The OpenID Connect login type. These values maps to: <br>  CODE: Authentication using Code Flow <br> POST: Authentication using Form Post <br> POST_AT: Authentication using Form Post with Access Token",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"CODE",
										"POST",
										"POST_AT",
									),
								},
							},
							"logout_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's RP-Initiated Logout Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's RP-Initiated Logout Endpoint.",
							},
							"post_logout_redirect_uri": schema.StringAttribute{
								Optional:            true,
								Description:         "The Post-Logout Redirect URI, where the OpenID Provider may redirect the user when RP-Initiated Logout has completed. This is a read-only parameter.",
								MarkdownDescription: "The Post-Logout Redirect URI, where the OpenID Provider may redirect the user when RP-Initiated Logout has completed. This is a read-only parameter.",
							},
							"pushed_authorization_request_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's OAuth 2.0 Pushed Authorization Request Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's OAuth 2.0 Pushed Authorization Request Endpoint.",
							},
							"redirect_uri": schema.StringAttribute{
								Optional:            true,
								Description:         "The redirect URI. This is a read-only parameter.",
								MarkdownDescription: "The redirect URI. This is a read-only parameter.",
							},
							"request_parameters": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"application_endpoint_override": schema.BoolAttribute{
											Required:            true,
											Description:         "Indicates whether the parameter value can be overridden by an Application Endpoint parameter",
											MarkdownDescription: "Indicates whether the parameter value can be overridden by an Application Endpoint parameter",
										},
										"attribute_value": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
											Required:            true,
											Description:         "Defines how an attribute in an attribute contract should be populated.",
											MarkdownDescription: "Defines how an attribute in an attribute contract should be populated.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "Request parameter name.",
											MarkdownDescription: "Request parameter name.",
										},
										"value": schema.StringAttribute{
											Optional:            true,
											Description:         "A request parameter value. A parameter can have either a value or a attribute value but not both. Value set here will be converted to an attribute value of source type TEXT. An empty value will be converted to attribute value of source type NO_MAPPING.",
											MarkdownDescription: "A request parameter value. A parameter can have either a value or a attribute value but not both. Value set here will be converted to an attribute value of source type TEXT. An empty value will be converted to attribute value of source type NO_MAPPING.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of request parameters. Request parameters with same name but different attribute values are treated as a multi-valued request parameter.",
								MarkdownDescription: "A list of request parameters. Request parameters with same name but different attribute values are treated as a multi-valued request parameter.",
							},
							"request_signing_algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "The request signing algorithm. Required only if you wish to use signed requests. Only asymmetric algorithms are allowed. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11.",
								MarkdownDescription: "The request signing algorithm. Required only if you wish to use signed requests. Only asymmetric algorithms are allowed. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NONE",
										"HS256",
										"HS384",
										"HS512",
										"RS256",
										"RS384",
										"RS512",
										"ES256",
										"ES384",
										"ES512",
										"PS256",
										"PS384",
										"PS512",
									),
								},
							},
							"scopes": schema.StringAttribute{
								Required:            true,
								Description:         "Space separated scope values that the OpenID Provider supports.",
								MarkdownDescription: "Space separated scope values that the OpenID Provider supports.",
							},
							"token_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's OAuth 2.0 Token Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's OAuth 2.0 Token Endpoint.",
							},
							"track_user_sessions_for_logout": schema.BoolAttribute{
								Optional:            true,
								Description:         "Determines whether PingFederate tracks a logout entry when a user signs in, so that the user session can later be terminated via a logout request from the OP. This setting must also be enabled in order for PingFederate to send an RP-initiated logout request to the OP during SLO.",
								MarkdownDescription: "Determines whether PingFederate tracks a logout entry when a user signs in, so that the user session can later be terminated via a logout request from the OP. This setting must also be enabled in order for PingFederate to send an RP-initiated logout request to the OP during SLO.",
							},
							"user_info_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's UserInfo Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's UserInfo Endpoint.",
							},
						},
						Optional: true,

						Description:         "The OpenID Provider settings.",
						MarkdownDescription: "The OpenID Provider settings.",
					},
					"protocol": schema.StringAttribute{
						Required:            true,
						Description:         "The browser-based SSO protocol to use.",
						MarkdownDescription: "The browser-based SSO protocol to use.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"SAML20",
								"WSFED",
								"SAML11",
								"SAML10",
								"OIDC",
							),
						},
					},
					"sign_authn_requests": schema.BoolAttribute{
						Optional:            true,
						Description:         "Determines whether SAML authentication requests should be signed.",
						MarkdownDescription: "Determines whether SAML authentication requests should be signed.",
					},
					"slo_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:            true,
									Description:         "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
									MarkdownDescription: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"ARTIFACT",
											"POST",
											"REDIRECT",
											"SOAP",
										),
									},
								},
								"response_url": schema.StringAttribute{
									Optional:            true,
									Description:         "The absolute or relative URL to which logout responses are sent. A relative URL can be specified if a base URL for the connection has been defined.",
									MarkdownDescription: "The absolute or relative URL to which logout responses are sent. A relative URL can be specified if a base URL for the connection has been defined.",
								},
								"url": schema.StringAttribute{
									Required:            true,
									Description:         "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
									MarkdownDescription: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of possible endpoints to send SLO requests and responses.",
						MarkdownDescription: "A list of possible endpoints to send SLO requests and responses.",
					},
					"sso_application_endpoint": schema.StringAttribute{
						Optional:            true,
						Description:         "Application endpoint that can be used to invoke single sign-on (SSO) for the connection. This is a read-only parameter.",
						MarkdownDescription: "Application endpoint that can be used to invoke single sign-on (SSO) for the connection. This is a read-only parameter.",
					},
					"sso_oauth_mapping": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_contract_fulfillment": schema.MapNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"source": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Optional:            true,
													Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
												},
												"type": schema.StringAttribute{
													Required:            true,
													Description:         "The source type of this key.",
													MarkdownDescription: "The source type of this key.",
													Validators: []validator.String{
														stringvalidator.OneOf(
															"TOKEN_EXCHANGE_PROCESSOR_POLICY",
															"ACCOUNT_LINK",
															"ADAPTER",
															"ASSERTION",
															"CONTEXT",
															"CUSTOM_DATA_STORE",
															"EXPRESSION",
															"JDBC_DATA_STORE",
															"LDAP_DATA_STORE",
															"PING_ONE_LDAP_GATEWAY_DATA_STORE",
															"MAPPED_ATTRIBUTES",
															"NO_MAPPING",
															"TEXT",
															"TOKEN",
															"REQUEST",
															"OAUTH_PERSISTENT_GRANT",
															"SUBJECT_TOKEN",
															"ACTOR_TOKEN",
															"PASSWORD_CREDENTIAL_VALIDATOR",
															"IDP_CONNECTION",
															"AUTHENTICATION_POLICY_CONTRACT",
															"CLAIMS",
															"LOCAL_IDENTITY_PROFILE",
															"EXTENDED_CLIENT_METADATA",
															"EXTENDED_PROPERTIES",
															"TRACKED_HTTP_PARAMS",
															"FRAGMENT",
															"INPUTS",
															"ATTRIBUTE_QUERY",
															"IDENTITY_STORE_USER",
															"IDENTITY_STORE_GROUP",
															"SCIM_USER",
															"SCIM_GROUP",
														),
													},
												},
											},
											Required:            true,
											Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
										},
										"value": schema.StringAttribute{
											Required:            true,
											Description:         "The value for this attribute.",
											MarkdownDescription: "The value for this attribute.",
										},
									},
								},
								Required:            true,
								Description:         "A list of mappings from attribute names to their fulfillment values.",
								MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
							},
							"attribute_sources": schema.ListAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								Description:         "A list of configured data stores to look up attributes from.",
								MarkdownDescription: "A list of configured data stores to look up attributes from.",
							},
							"issuance_criteria": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"conditional_criteria": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"attribute_name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of the attribute to use in this issuance criterion.",
													MarkdownDescription: "The name of the attribute to use in this issuance criterion.",
												},
												"condition": schema.StringAttribute{
													Required:            true,
													Description:         "The condition that will be applied to the source attribute's value and the expected value.",
													MarkdownDescription: "The condition that will be applied to the source attribute's value and the expected value.",
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
												"error_result": schema.StringAttribute{
													Optional:            true,
													Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
												},
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The expected value of this issuance criterion.",
													MarkdownDescription: "The expected value of this issuance criterion.",
												},
											},
										},
										Optional:            true,
										Description:         "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
										MarkdownDescription: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
									},
									"expression_criteria": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"error_result": schema.StringAttribute{
													Optional:            true,
													Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
												},
												"expression": schema.StringAttribute{
													Required:            true,
													Description:         "The OGNL expression to evaluate.",
													MarkdownDescription: "The OGNL expression to evaluate.",
												},
											},
										},
										Optional:            true,
										Description:         "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
										MarkdownDescription: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
									},
								},
								Optional:            true,
								Description:         "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
								MarkdownDescription: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
							},
						},
						Optional:            true,
						Description:         "IdP Browser SSO OAuth Attribute Mapping",
						MarkdownDescription: "IdP Browser SSO OAuth Attribute Mapping",
					},
					"sso_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Required:            true,
									Description:         "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
									MarkdownDescription: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"ARTIFACT",
											"POST",
											"REDIRECT",
											"SOAP",
										),
									},
								},
								"url": schema.StringAttribute{
									Required:            true,
									Description:         "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
									MarkdownDescription: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
								},
							},
						},
						Optional:            true,
						Description:         "The IdP SSO endpoints that define where to send your authentication requests. Only required for SP initiated SSO. This is required for SAML x.x and WS-FED Connections.",
						MarkdownDescription: "The IdP SSO endpoints that define where to send your authentication requests. Only required for SP initiated SSO. This is required for SAML x.x and WS-FED Connections.",
					},
					"url_whitelist_entries": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"allow_query_and_fragment": schema.BoolAttribute{
									Optional:            true,
									Description:         "Allow Any Query/Fragment",
									MarkdownDescription: "Allow Any Query/Fragment",
								},
								"require_https": schema.BoolAttribute{
									Optional:            true,
									Description:         "Require HTTPS",
									MarkdownDescription: "Require HTTPS",
								},
								"valid_domain": schema.StringAttribute{
									Optional:            true,
									Description:         "Valid Domain Name (leading wildcard '*.' allowed)",
									MarkdownDescription: "Valid Domain Name (leading wildcard '*.' allowed)",
								},
								"valid_path": schema.StringAttribute{
									Optional:            true,
									Description:         "Valid Path (leave blank to allow any path)",
									MarkdownDescription: "Valid Path (leave blank to allow any path)",
								},
							},
						},
						Optional:            true,
						Description:         "For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled.",
						MarkdownDescription: "For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled.",
					},
				},
				Optional:            true,
				Description:         "The settings used to enable secure browser-based SSO to resources at your site.",
				MarkdownDescription: "The settings used to enable secure browser-based SSO to resources at your site.",
			},
			"idp_oauth_grant_attribute_mapping": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"access_token_manager_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_token_manager_ref": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:            true,
											Description:         "The ID of the resource.",
											MarkdownDescription: "The ID of the resource.",
										},
										"location": schema.StringAttribute{
											Optional:            true,
											Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										},
									},
									Optional:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
								"attribute_contract_fulfillment": schema.MapNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:            true,
														Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key.",
														MarkdownDescription: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:            true,
												Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:            true,
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "A list of configured data stores to look up attributes from.",
									MarkdownDescription: "A list of configured data stores to look up attributes from.",
								},
								"issuance_criteria": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"attribute_name": schema.StringAttribute{
														Required:            true,
														Description:         "The name of the attribute to use in this issuance criterion.",
														MarkdownDescription: "The name of the attribute to use in this issuance criterion.",
													},
													"condition": schema.StringAttribute{
														Required:            true,
														Description:         "The condition that will be applied to the source attribute's value and the expected value.",
														MarkdownDescription: "The condition that will be applied to the source attribute's value and the expected value.",
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
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"source": schema.SingleNestedAttribute{
														Attributes: map[string]schema.Attribute{
															"id": schema.StringAttribute{
																Optional:            true,
																Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															},
															"type": schema.StringAttribute{
																Required:            true,
																Description:         "The source type of this key.",
																MarkdownDescription: "The source type of this key.",
																Validators: []validator.String{
																	stringvalidator.OneOf(
																		"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																		"ACCOUNT_LINK",
																		"ADAPTER",
																		"ASSERTION",
																		"CONTEXT",
																		"CUSTOM_DATA_STORE",
																		"EXPRESSION",
																		"JDBC_DATA_STORE",
																		"LDAP_DATA_STORE",
																		"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																		"MAPPED_ATTRIBUTES",
																		"NO_MAPPING",
																		"TEXT",
																		"TOKEN",
																		"REQUEST",
																		"OAUTH_PERSISTENT_GRANT",
																		"SUBJECT_TOKEN",
																		"ACTOR_TOKEN",
																		"PASSWORD_CREDENTIAL_VALIDATOR",
																		"IDP_CONNECTION",
																		"AUTHENTICATION_POLICY_CONTRACT",
																		"CLAIMS",
																		"LOCAL_IDENTITY_PROFILE",
																		"EXTENDED_CLIENT_METADATA",
																		"EXTENDED_PROPERTIES",
																		"TRACKED_HTTP_PARAMS",
																		"FRAGMENT",
																		"INPUTS",
																		"ATTRIBUTE_QUERY",
																		"IDENTITY_STORE_USER",
																		"IDENTITY_STORE_GROUP",
																		"SCIM_USER",
																		"SCIM_GROUP",
																	),
																},
															},
														},
														Required:            true,
														Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "The expected value of this issuance criterion.",
														MarkdownDescription: "The expected value of this issuance criterion.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											MarkdownDescription: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
										},
										"expression_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"expression": schema.StringAttribute{
														Required:            true,
														Description:         "The OGNL expression to evaluate.",
														MarkdownDescription: "The OGNL expression to evaluate.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											MarkdownDescription: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
										},
									},
									Optional:            true,
									Description:         "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									MarkdownDescription: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
								},
							},
						},
						Optional:            true,
						Description:         "A mapping in a connection that defines how access tokens are created.",
						MarkdownDescription: "A mapping in a connection that defines how access tokens are created.",
					},
					"idp_oauth_attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional: true,

											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of additional attributes that are present in the incoming assertion.",
								MarkdownDescription: "A list of additional attributes that are present in the incoming assertion.",
							},
						},
						Optional:            true,
						Description:         "A set of user attributes that the IdP sends in the OAuth Assertion Grant.",
						MarkdownDescription: "A set of user attributes that the IdP sends in the OAuth Assertion Grant.",
					},
				},
				Optional:            true,
				Description:         "The OAuth Assertion Grant settings used to map from your IdP.",
				MarkdownDescription: "The OAuth Assertion Grant settings used to map from your IdP.",
			},
			"inbound_provisioning": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"action_on_delete": schema.StringAttribute{
						Optional:            true,
						Description:         "Specify behavior of how SCIM DELETE requests are handled.",
						MarkdownDescription: "Specify behavior of how SCIM DELETE requests are handled.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"DISABLE_USER",
								"PERMANENTLY_DELETE_USER",
							),
						},
					},
					"custom_schema": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"multi_valued": schema.BoolAttribute{
											Optional:            true,
											Description:         "Indicates whether the attribute is multi-valued.",
											MarkdownDescription: "Indicates whether the attribute is multi-valued.",
										},
										"name": schema.StringAttribute{
											Optional:            true,
											Description:         "Name of the attribute.",
											MarkdownDescription: "Name of the attribute.",
										},
										"sub_attributes": schema.ListAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											Description:         "List of sub-attributes for an attribute.",
											MarkdownDescription: "List of sub-attributes for an attribute.",
										},
										"types": schema.ListAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											Description:         "Represents the name of each attribute type in case of multi-valued attribute.",
											MarkdownDescription: "Represents the name of each attribute type in case of multi-valued attribute.",
										},
									},
								},
								Optional: true,
								Computed: true,
							},
							"namespace": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
						},
						Required:            true,
						Description:         "Custom SCIM Attributes configuration.",
						MarkdownDescription: "Custom SCIM Attributes configuration.",
					},
					"group_support": schema.BoolAttribute{
						Required:            true,
						Description:         "Specify support for provisioning of groups.",
						MarkdownDescription: "Specify support for provisioning of groups.",
					},
					"groups": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"read_groups": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"core_attributes": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
												MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
											},
											"extended_attributes": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of additional attributes that are added to the SCIM response.",
												MarkdownDescription: "A list of additional attributes that are added to the SCIM response.",
											},
										},
										Required:            true,
										Description:         "A set of user attributes that the IdP sends in the SCIM response.",
										MarkdownDescription: "A set of user attributes that the IdP sends in the SCIM response.",
									},
									"attribute_fulfillment": schema.MapNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
									"attributes": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of LDAP data store attributes to populate a response to a user-provisioning request.",
										MarkdownDescription: "A list of LDAP data store attributes to populate a response to a user-provisioning request.",
									},
								},
								Required:            true,
								Description:         "Group info lookup and respond to incoming SCIM requests configuration.",
								MarkdownDescription: "Group info lookup and respond to incoming SCIM requests configuration.",
							},
							"write_groups": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_fulfillment": schema.MapNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
								},
								Required:            true,
								Description:         "Group creation configuration.",
								MarkdownDescription: "Group creation configuration.",
							},
						},
						Required:            true,
						Description:         "Group creation and read configuration.",
						MarkdownDescription: "Group creation and read configuration.",
					},
					"user_repository": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "The user repository type.",
								MarkdownDescription: "The user repository type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"LDAP",
										"IDENTITY_STORE",
									),
								},
							},
						},
						Required:            true,
						Description:         "SCIM Inbound Provisioning user repository.",
						MarkdownDescription: "SCIM Inbound Provisioning user repository.",
					},
					"users": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"read_users": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"core_attributes": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
												MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
											},
											"extended_attributes": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of additional attributes that are added to the SCIM response.",
												MarkdownDescription: "A list of additional attributes that are added to the SCIM response.",
											},
										},
										Required:            true,
										Description:         "A set of user attributes that the IdP sends in the SCIM response.",
										MarkdownDescription: "A set of user attributes that the IdP sends in the SCIM response.",
									},
									"attribute_fulfillment": schema.MapNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
									"attributes": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of LDAP data store attributes to populate a response to a user-provisioning request.",
										MarkdownDescription: "A list of LDAP data store attributes to populate a response to a user-provisioning request.",
									},
								},
								Required:            true,
								Description:         "User info lookup and respond to incoming SCIM requests configuration.",
								MarkdownDescription: "User info lookup and respond to incoming SCIM requests configuration.",
							},
							"write_users": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_fulfillment": schema.MapNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:            true,
															Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key.",
															MarkdownDescription: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Required:            true,
													Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"value": schema.StringAttribute{
													Required:            true,
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
								},
								Required:            true,
								Description:         "User creation configuration.",
								MarkdownDescription: "User creation configuration.",
							},
						},
						Required:            true,
						Description:         "User creation and read configuration.",
						MarkdownDescription: "User creation and read configuration.",
					},
				},
				Optional:            true,
				Description:         "SCIM Inbound Provisioning specifies how and when to provision user accounts and groups.",
				MarkdownDescription: "SCIM Inbound Provisioning specifies how and when to provision user accounts and groups.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The connection name.",
				MarkdownDescription: "The connection name.",
			},
			"license_connection_group": schema.StringAttribute{
				Optional:            true,
				Description:         "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				MarkdownDescription: "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
			},
			"logging_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The level of transaction logging applicable for this connection. Default is STANDARD.",
				MarkdownDescription: "The level of transaction logging applicable for this connection. Default is STANDARD.",
				Default:             stringdefault.StaticString("STANDARD"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "STANDARD", "ENHANCED", "FULL"),
				},
			},
			"metadata_reload_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Configuration settings to enable automatic reload of partner's metadata.",
				MarkdownDescription: "Configuration settings to enable automatic reload of partner's metadata.",
				Attributes: map[string]schema.Attribute{
					"metadata_url_ref": schema.SingleNestedAttribute{
						Required:            true,
						Description:         "A reference to a resource.",
						MarkdownDescription: "A reference to a resource.",
						Attributes:          resourcelink.ToSchemaNoLocation(),
					},
					"enable_auto_metadata_update": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies whether the metadata of the connection will be automatically reloaded. The default value is true.",
						MarkdownDescription: "Specifies whether the metadata of the connection will be automatically reloaded. The default value is true.",
						Default:             booldefault.StaticBool(true),
					},
				},
			},
			"oidc_client_credentials": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Required:            true,
						Description:         "The OpenID Connect client identitification.",
						MarkdownDescription: "The OpenID Connect client identitification.",
					},
					"client_secret": schema.StringAttribute{
						Required:            true,
						Description:         "The OpenID Connect client secret. To update the client secret, specify the plaintext value in this field.  This field will not be populated for GET requests.",
						MarkdownDescription: "The OpenID Connect client secret. To update the client secret, specify the plaintext value in this field.  This field will not be populated for GET requests.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Optional:            true,
				Description:         "The OpenID Connect Client Credentials settings. This is required for an OIDC Connection.",
				MarkdownDescription: "The OpenID Connect Client Credentials settings. This is required for an OIDC Connection.",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("IDP"),
				Validators: []validator.String{
					stringvalidator.OneOf("IDP", "SP"),
				},
				Description:         "The type of this connection. Default is 'IDP'.",
				MarkdownDescription: "The type of this connection. Default is 'IDP'.",
			},
			"virtual_entity_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "List of alternate entity IDs that identifies the local server to this partner.",
				MarkdownDescription: "List of alternate entity IDs that identifies the local server to this partner.",
			},
			"ws_trust": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional: true,

											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
										},
									},
								},
								Optional:            true,
								Description:         "A list of additional attributes that are receive in the incoming assertion.",
								MarkdownDescription: "A list of additional attributes that are receive in the incoming assertion.",
							},
						},
						Required:            true,
						Description:         "A set of user attributes that this server will receive in the token.",
						MarkdownDescription: "A set of user attributes that this server will receive in the token.",
					},
					"generate_local_token": schema.BoolAttribute{
						Required:            true,
						Description:         "Indicates whether a local token needs to be generated. The default value is false.",
						MarkdownDescription: "Indicates whether a local token needs to be generated. The default value is false.",
					},
					"token_generator_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": schema.MapNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:            true,
														Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key.",
														MarkdownDescription: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:            true,
												Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:            true,
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "A list of configured data stores to look up attributes from.",
									MarkdownDescription: "A list of configured data stores to look up attributes from.",
								},
								"default_mapping": schema.BoolAttribute{
									Optional:            true,
									Description:         "Indicates whether the token generator mapping is the default mapping. The default value is false.",
									MarkdownDescription: "Indicates whether the token generator mapping is the default mapping. The default value is false.",
								},
								"issuance_criteria": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"attribute_name": schema.StringAttribute{
														Required:            true,
														Description:         "The name of the attribute to use in this issuance criterion.",
														MarkdownDescription: "The name of the attribute to use in this issuance criterion.",
													},
													"condition": schema.StringAttribute{
														Required:            true,
														Description:         "The condition that will be applied to the source attribute's value and the expected value.",
														MarkdownDescription: "The condition that will be applied to the source attribute's value and the expected value.",
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
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"source": schema.SingleNestedAttribute{
														Attributes: map[string]schema.Attribute{
															"id": schema.StringAttribute{
																Optional:            true,
																Description:         "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																MarkdownDescription: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
															},
															"type": schema.StringAttribute{
																Required:            true,
																Description:         "The source type of this key.",
																MarkdownDescription: "The source type of this key.",
																Validators: []validator.String{
																	stringvalidator.OneOf(
																		"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																		"ACCOUNT_LINK",
																		"ADAPTER",
																		"ASSERTION",
																		"CONTEXT",
																		"CUSTOM_DATA_STORE",
																		"EXPRESSION",
																		"JDBC_DATA_STORE",
																		"LDAP_DATA_STORE",
																		"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																		"MAPPED_ATTRIBUTES",
																		"NO_MAPPING",
																		"TEXT",
																		"TOKEN",
																		"REQUEST",
																		"OAUTH_PERSISTENT_GRANT",
																		"SUBJECT_TOKEN",
																		"ACTOR_TOKEN",
																		"PASSWORD_CREDENTIAL_VALIDATOR",
																		"IDP_CONNECTION",
																		"AUTHENTICATION_POLICY_CONTRACT",
																		"CLAIMS",
																		"LOCAL_IDENTITY_PROFILE",
																		"EXTENDED_CLIENT_METADATA",
																		"EXTENDED_PROPERTIES",
																		"TRACKED_HTTP_PARAMS",
																		"FRAGMENT",
																		"INPUTS",
																		"ATTRIBUTE_QUERY",
																		"IDENTITY_STORE_USER",
																		"IDENTITY_STORE_GROUP",
																		"SCIM_USER",
																		"SCIM_GROUP",
																	),
																},
															},
														},
														Required:            true,
														Description:         "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														MarkdownDescription: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "The expected value of this issuance criterion.",
														MarkdownDescription: "The expected value of this issuance criterion.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											MarkdownDescription: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
										},
										"expression_criteria": schema.ListNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"error_result": schema.StringAttribute{
														Optional:            true,
														Description:         "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														MarkdownDescription: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
													"expression": schema.StringAttribute{
														Required:            true,
														Description:         "The OGNL expression to evaluate.",
														MarkdownDescription: "The OGNL expression to evaluate.",
													},
												},
											},
											Optional:            true,
											Description:         "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											MarkdownDescription: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
										},
									},
									Optional:            true,
									Description:         "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									MarkdownDescription: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
								"sp_token_generator_ref": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:            true,
											Description:         "The ID of the resource.",
											MarkdownDescription: "The ID of the resource.",
										},
										"location": schema.StringAttribute{
											Optional:            true,
											Description:         "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											MarkdownDescription: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										},
									},
									Required:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of token generators to generate local tokens. Required if a local token needs to be generated.",
						MarkdownDescription: "A list of token generators to generate local tokens. Required if a local token needs to be generated.",
					},
				},
				Optional:            true,
				Description:         "Ws-Trust STS provides validation of incoming tokens which enable SSO access to Web Services. It also allows generation of local tokens for Web Services.",
				MarkdownDescription: "Ws-Trust STS provides validation of incoming tokens which enable SSO access to Web Services. It also allows generation of local tokens for Web Services.",
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "connection_id", false, false, "The persistent, unique ID for the connection. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.")
	resp.Schema = schema
}

func addOptionalSpIdpConnectionFields(ctx context.Context, addRequest *client.IdpConnection, plan spIdpConnectionResourceModel) error {
	addRequest.ErrorPageMsgId = plan.ErrorPageMsgId.ValueStringPointer()
	addRequest.Id = plan.ConnectionId.ValueStringPointer()
	addRequest.Type = plan.Type.ValueStringPointer()
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.BaseUrl = plan.BaseUrl.ValueStringPointer()
	addRequest.DefaultVirtualEntityId = plan.DefaultVirtualEntityId.ValueStringPointer()
	addRequest.LicenseConnectionGroup = plan.LicenseConnectionGroup.ValueStringPointer()
	addRequest.LoggingMode = plan.LoggingMode.ValueStringPointer()

	var virtualIdentitySlice []string
	plan.VirtualEntityIds.ElementsAs(ctx, &virtualIdentitySlice, false)
	addRequest.VirtualEntityIds = virtualIdentitySlice

	if internaltypes.IsDefined(plan.OidcClientCredentials) {
		addRequest.OidcClientCredentials = &client.OIDCClientCredentials{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.OidcClientCredentials, true)), addRequest.OidcClientCredentials)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.MetadataReloadSettings) {
		addRequest.MetadataReloadSettings = &client.ConnectionMetadataUrl{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.MetadataReloadSettings, false)), addRequest.MetadataReloadSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Credentials) {
		addRequest.Credentials = &client.ConnectionCredentials{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Credentials, true)), addRequest.Credentials)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ContactInfo) {
		addRequest.ContactInfo = &client.ContactInfo{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ContactInfo, false)), addRequest.ContactInfo)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AdditionalAllowedEntitiesConfiguration) {
		addRequest.AdditionalAllowedEntitiesConfiguration = &client.AdditionalAllowedEntitiesConfiguration{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AdditionalAllowedEntitiesConfiguration, false)), addRequest.AdditionalAllowedEntitiesConfiguration)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ExtendedProperties) {
		addRequest.ExtendedProperties = &map[string]client.ParameterValues{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ExtendedProperties, true)), &addRequest.ExtendedProperties)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.IdpBrowserSso) {
		addRequest.IdpBrowserSso = &client.IdpBrowserSso{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.IdpBrowserSso, false)), addRequest.IdpBrowserSso)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeQuery) {
		addRequest.AttributeQuery = &client.IdpAttributeQuery{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeQuery, false)), addRequest.AttributeQuery)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.IdpOAuthGrantAttributeMapping) {
		addRequest.IdpOAuthGrantAttributeMapping = &client.IdpOAuthGrantAttributeMapping{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.IdpOAuthGrantAttributeMapping, true)), addRequest.IdpOAuthGrantAttributeMapping)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.WsTrust) {
		addRequest.WsTrust = &client.IdpWsTrust{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.WsTrust, false)), addRequest.WsTrust)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.InboundProvisioning) {
		addRequest.InboundProvisioning = &client.IdpInboundProvisioning{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.InboundProvisioning, false)), addRequest.InboundProvisioning)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *spIdpConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_idp_connection"
}

func (r *spIdpConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSpIdpConnectionResponse(ctx context.Context, r *client.IdpConnection, plan, state *spIdpConnectionResourceModel) diag.Diagnostics {
	var diags, objDiags diag.Diagnostics

	if r.OidcClientCredentials != nil {
		var oidcClientCredentialsAttrValues map[string]attr.Value
		var clientSecret string = ""
		if len(plan.OidcClientCredentials.Attributes()) > 0 && plan.OidcClientCredentials.Attributes()["client_secret"] != nil {
			clientSecret = plan.OidcClientCredentials.Attributes()["client_secret"].(types.String).ValueString()
		}
		oidcClientCredentialsAttrValues = map[string]attr.Value{
			"client_id":     types.StringValue(r.OidcClientCredentials.ClientId),
			"client_secret": types.StringValue(clientSecret),
		}

		state.OidcClientCredentials, objDiags = types.ObjectValue(oidcClientCredentialsAttrTypes, oidcClientCredentialsAttrValues)
		diags.Append(objDiags...)
	} else {
		state.OidcClientCredentials = types.ObjectNull(oidcClientCredentialsAttrTypes)
	}

	state.AdditionalAllowedEntitiesConfiguration, objDiags = types.ObjectValueFrom(ctx, additionalAllowedEntitiesConfigurationAttrTypes, r.AdditionalAllowedEntitiesConfiguration)
	diags.Append(objDiags...)
	state.AttributeQuery, objDiags = types.ObjectValueFrom(ctx, attributeQueryAttrTypes, r.AttributeQuery)
	diags.Append(objDiags...)
	state.ContactInfo, objDiags = types.ObjectValueFrom(ctx, contactInfoAttrTypes, r.ContactInfo)
	diags.Append(objDiags...)
	state.Credentials, objDiags = types.ObjectValueFrom(ctx, credentialsAttrTypes, r.Credentials)
	diags.Append(objDiags...)
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	state.EntityId = types.StringValue(r.EntityId)
	state.ErrorPageMsgId = types.StringPointerValue(r.ErrorPageMsgId)
	state.ExtendedProperties, objDiags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: extendedPropertiesElemAttrTypes}, r.ExtendedProperties)
	diags.Append(objDiags...)
	state.Id = types.StringPointerValue(r.Id)
	state.IdpBrowserSso, objDiags = types.ObjectValueFrom(ctx, idpBrowserSsoAttrTypes, r.IdpBrowserSso)
	diags.Append(objDiags...)
	state.IdpOAuthGrantAttributeMapping, objDiags = types.ObjectValueFrom(ctx, idpOAuthGrantAttributeMappingAttrTypes, r.IdpOAuthGrantAttributeMapping)
	diags.Append(objDiags...)
	state.InboundProvisioning, objDiags = types.ObjectValueFrom(ctx, inboundProvisioningAttrTypes, r.InboundProvisioning)
	diags.Append(objDiags...)
	state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	state.MetadataReloadSettings, objDiags = types.ObjectValueFrom(ctx, metadataReloadSettingsAttrTypes, r.MetadataReloadSettings)
	diags.Append(objDiags...)
	state.Name = types.StringValue(r.Name)
	state.OidcClientCredentials = types.ObjectNull(oidcClientCredentialsAttrTypes)
	state.Type = types.StringPointerValue(r.Type)
	state.VirtualEntityIds = internaltypes.GetStringSet(r.VirtualEntityIds)

	if r.WsTrust != nil {

		var tokenGeneratorMappings []basetypes.ObjectValue
		for _, tokenGeneratorMapping := range r.WsTrust.TokenGeneratorMappings {
			spTokenGeneratorRef, objDiags := resourcelink.ToStateNoLocation(&tokenGeneratorMapping.SpTokenGeneratorRef)
			diags.Append(objDiags...)

			var attributeSources basetypes.ListValue
			attributeSources, objDiags = attributesources.ToState(ctx, tokenGeneratorMapping.AttributeSources, false)
			diags.Append(objDiags...)

			attributeContractFulfillment, objDiags := attributecontractfulfillment.ToState(ctx, tokenGeneratorMapping.AttributeContractFulfillment)
			diags.Append(objDiags...)

			issuanceCriteria, objDiags := issuancecriteria.ToState(ctx, tokenGeneratorMapping.IssuanceCriteria)
			diags.Append(objDiags...)

			tokenGeneratorAttrValues := map[string]attr.Value{
				"sp_token_generator_ref":         spTokenGeneratorRef,
				"attribute_sources":              attributeSources,
				"default_mapping":                types.BoolPointerValue(tokenGeneratorMapping.DefaultMapping),
				"attribute_contract_fulfillment": attributeContractFulfillment,
				"issuance_criteria":              issuanceCriteria,
			}

			tokenGeneratorMappingState, objDiags := types.ObjectValue(tokenGeneratorAttrTypes, tokenGeneratorAttrValues)
			diags.Append(objDiags...)
			tokenGeneratorMappings = append(tokenGeneratorMappings, tokenGeneratorMappingState)
		}

		attributeContract, objDiags := types.ObjectValueFrom(ctx, coreAndExtendedAttributesAttrTypes, r.WsTrust.AttributeContract)
		diags.Append(objDiags...)

		tokenGeneratorMappingsList, objDiags := types.ListValueFrom(ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: tokenGeneratorAttrTypes}}, tokenGeneratorMappings)
		diags.Append(objDiags...)

		wsTrustAttrValues := map[string]attr.Value{
			"attribute_contract":       attributeContract,
			"generate_local_token":     types.BoolValue(r.WsTrust.GenerateLocalToken),
			"token_generator_mappings": tokenGeneratorMappingsList,
		}

		state.WsTrust, objDiags = types.ObjectValue(wsTrustAttrTypes, wsTrustAttrValues)
		diags.Append(objDiags...)
	} else {
		state.WsTrust = types.ObjectNull(wsTrustAttrTypes)
	}

	return diags
}

func (r *spIdpConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan spIdpConnectionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSpIdpConnection := client.NewIdpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	err := addOptionalSpIdpConnectionFields(ctx, createSpIdpConnection, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for SpIdpConnection", err.Error())
		return
	}

	requestJson, err := createSpIdpConnection.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateSpIdpConnection := r.apiClient.SpIdpConnectionsAPI.CreateConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSpIdpConnection = apiCreateSpIdpConnection.Body(*createSpIdpConnection)
	spIdpConnectionResponse, httpResp, err := r.apiClient.SpIdpConnectionsAPI.CreateConnectionExecute(apiCreateSpIdpConnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the SpIdpConnection", err, httpResp)
		return
	}

	// Read the response into the state
	var state spIdpConnectionResourceModel

	diags = readSpIdpConnectionResponse(ctx, spIdpConnectionResponse, &plan, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *spIdpConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state spIdpConnectionResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSpIdpConnection, httpResp, err := r.apiClient.SpIdpConnectionsAPI.GetConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Sp Idp Connection", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Sp Idp Connection", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readSpIdpConnectionResponse(ctx, apiReadSpIdpConnection, nil, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *spIdpConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan spIdpConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSpIdpConnection := r.apiClient.SpIdpConnectionsAPI.UpdateConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ConnectionId.ValueString())
	createUpdateRequest := client.NewIdpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	err := addOptionalSpIdpConnectionFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sp Idp Connection", err.Error())
		return
	}

	updateSpIdpConnection = updateSpIdpConnection.Body(*createUpdateRequest)
	updateSpIdpConnectionResponse, httpResp, err := r.apiClient.SpIdpConnectionsAPI.UpdateConnectionExecute(updateSpIdpConnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Sp Idp Connection", err, httpResp)
		return
	}

	// Read the response
	var state spIdpConnectionResourceModel
	diags = readSpIdpConnectionResponse(ctx, updateSpIdpConnectionResponse, &plan, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *spIdpConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state spIdpConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.SpIdpConnectionsAPI.DeleteConnection(config.AuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Sp Idp Connection", err, httpResp)
	}
}

func (r *spIdpConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to connection_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("connection_id"), req, resp)
}
