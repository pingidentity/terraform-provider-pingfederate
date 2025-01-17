package spidpconnection

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/datastorerepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/inboundprovisioninguserrepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &spIdpConnectionResource{}
	_ resource.ResourceWithConfigure   = &spIdpConnectionResource{}
	_ resource.ResourceWithImportState = &spIdpConnectionResource{}

	metadataReloadSettingsAttrTypes = map[string]attr.Type{
		"metadata_url_ref":            types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"enable_auto_metadata_update": types.BoolType,
	}

	oidcClientCredentialsAttrTypes = map[string]attr.Type{
		"client_id":        types.StringType,
		"client_secret":    types.StringType,
		"encrypted_secret": types.StringType,
	}

	credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes = map[string]attr.Type{
		"password":           types.StringType,
		"encrypted_password": types.StringType,
		"username":           types.StringType,
	}
	credentialsInboundBackChannelAuthAttrTypes = map[string]attr.Type{
		"certs":                   types.ListType{ElemType: connectioncert.ObjType()},
		"digital_signature":       types.BoolType,
		"http_basic_credentials":  types.ObjectType{AttrTypes: credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes},
		"require_ssl":             types.BoolType,
		"verification_issuer_dn":  types.StringType,
		"verification_subject_dn": types.StringType,
	}
	credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes = map[string]attr.Type{
		"password":           types.StringType,
		"encrypted_password": types.StringType,
		"username":           types.StringType,
	}

	credentialsOutboundBackChannelAuthAttrTypes = map[string]attr.Type{
		"digital_signature":      types.BoolType,
		"http_basic_credentials": types.ObjectType{AttrTypes: credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes},
		"ssl_auth_key_pair_ref":  types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"validate_partner_cert":  types.BoolType,
	}

	credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType = types.ObjectType{AttrTypes: resourcelink.AttrType()}

	credentialsSigningSettingsAttrTypes = map[string]attr.Type{
		"algorithm":                         types.StringType,
		"alternative_signing_key_pair_refs": types.SetType{ElemType: credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType},
		"include_cert_in_signature":         types.BoolType,
		"include_raw_key_in_signature":      types.BoolType,
		"signing_key_pair_ref":              types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	credentialsAttrTypes = map[string]attr.Type{
		"block_encryption_algorithm":        types.StringType,
		"certs":                             types.ListType{ElemType: connectioncert.ObjType()},
		"decryption_key_pair_ref":           types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"inbound_back_channel_auth":         types.ObjectType{AttrTypes: credentialsInboundBackChannelAuthAttrTypes},
		"key_transport_algorithm":           types.StringType,
		"outbound_back_channel_auth":        types.ObjectType{AttrTypes: credentialsOutboundBackChannelAuthAttrTypes},
		"secondary_decryption_key_pair_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"signing_settings":                  types.ObjectType{AttrTypes: credentialsSigningSettingsAttrTypes},
		"verification_issuer_dn":            types.StringType,
		"verification_subject_dn":           types.StringType,
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
		"values": types.SetType{ElemType: types.StringType},
	}

	// idp_browser_sso
	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes = map[string]attr.Type{
		"name": types.StringType,
	}
	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType   = types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes}
	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes = map[string]attr.Type{
		"name": types.StringType,
	}
	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType = types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes}
	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes                     = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType},
	}

	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes = map[string]attr.Type{
		"application_icon_url": types.StringType,
		"application_name":     types.StringType,
	}

	idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes = map[string]attr.Type{
		"attribute_contract":      types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes},
		"configuration":           types.ObjectType{AttrTypes: pluginconfiguration.AttrTypes()},
		"id":                      types.StringType,
		"name":                    types.StringType,
		"parent_ref":              types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"plugin_descriptor_ref":   types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"target_application_info": types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes},
	}

	idpBrowserSsoAdapterMappingsAttrTypes = map[string]attr.Type{
		"adapter_override_settings":      types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"attribute_sources":              types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrTypes()},
		"restrict_virtual_entity_ids":    types.BoolType,
		"restricted_virtual_entity_ids":  types.SetType{ElemType: types.StringType},
		"sp_adapter_ref":                 types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	idpBrowserSsoAdapterMappingsElementType         = types.ObjectType{AttrTypes: idpBrowserSsoAdapterMappingsAttrTypes}
	idpBrowserSsoArtifactResolverLocationsAttrTypes = map[string]attr.Type{
		"index": types.Int64Type,
		"url":   types.StringType,
	}
	idpBrowserSsoArtifactResolverLocationsElementType = types.ObjectType{AttrTypes: idpBrowserSsoArtifactResolverLocationsAttrTypes}
	idpBrowserSsoArtifactAttrTypes                    = map[string]attr.Type{
		"lifetime":           types.Int64Type,
		"resolver_locations": types.SetType{ElemType: idpBrowserSsoArtifactResolverLocationsElementType},
		"source_id":          types.StringType,
	}
	idpBrowserSsoAttributeContractCoreAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	idpBrowserSsoAttributeContractCoreAttributesElementType   = types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractCoreAttributesAttrTypes}
	idpBrowserSsoAttributeContractExtendedAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	idpBrowserSsoAttributeContractExtendedAttributesElementType = types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractExtendedAttributesAttrTypes}
	idpBrowserSsoAttributeContractAttrTypes                     = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: idpBrowserSsoAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: idpBrowserSsoAttributeContractExtendedAttributesElementType},
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentElementType = types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentAttrTypes}

	idpBrowserSsoAuthenticationPolicyContractMappingsAuthenticationPolicyContractRefAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaAttrTypes = map[string]attr.Type{
		"attribute_name": types.StringType,
		"condition":      types.StringType,
		"error_result":   types.StringType,
		"source":         types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaSourceAttrTypes},
		"value":          types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaElementType = types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaAttrTypes}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaExpressionCriteriaAttrTypes    = map[string]attr.Type{
		"error_result": types.StringType,
		"expression":   types.StringType,
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaExpressionCriteriaElementType = types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaExpressionCriteriaAttrTypes}
	idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaAttrTypes                     = map[string]attr.Type{
		"conditional_criteria": types.SetType{ElemType: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaConditionalCriteriaElementType},
		"expression_criteria":  types.SetType{ElemType: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaExpressionCriteriaElementType},
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsAttrTypes = map[string]attr.Type{
		"attribute_contract_fulfillment":     types.MapType{ElemType: idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":                  types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"authentication_policy_contract_ref": types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsAuthenticationPolicyContractRefAttrTypes},
		"issuance_criteria":                  types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaAttrTypes},
		"restrict_virtual_server_ids":        types.BoolType,
		"restricted_virtual_server_ids":      types.SetType{ElemType: types.StringType},
	}
	idpBrowserSsoAuthenticationPolicyContractMappingsElementType = types.ObjectType{AttrTypes: idpBrowserSsoAuthenticationPolicyContractMappingsAttrTypes}
	idpBrowserSsoAuthnContextMappingsAttrTypes                   = map[string]attr.Type{
		"local":  types.StringType,
		"remote": types.StringType,
	}
	idpBrowserSsoAuthnContextMappingsElementType = types.ObjectType{AttrTypes: idpBrowserSsoAuthnContextMappingsAttrTypes}
	idpBrowserSsoDecryptionPolicyAttrTypes       = map[string]attr.Type{
		"assertion_encrypted":           types.BoolType,
		"attributes_encrypted":          types.BoolType,
		"slo_encrypt_subject_name_id":   types.BoolType,
		"slo_subject_name_id_encrypted": types.BoolType,
		"subject_name_id_encrypted":     types.BoolType,
	}
	idpBrowserSsoJitProvisioningUserAttributesAttributeContractAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	idpBrowserSsoJitProvisioningUserAttributesAttributeContractElementType = types.ObjectType{AttrTypes: idpBrowserSsoJitProvisioningUserAttributesAttributeContractAttrTypes}
	idpBrowserSsoJitProvisioningUserAttributesAttrTypes                    = map[string]attr.Type{
		"attribute_contract": types.SetType{ElemType: idpBrowserSsoJitProvisioningUserAttributesAttributeContractElementType},
		"do_attribute_query": types.BoolType,
	}

	idpBrowserSsoJitProvisioningAttrTypes = map[string]attr.Type{
		"error_handling":  types.StringType,
		"event_trigger":   types.StringType,
		"user_attributes": types.ObjectType{AttrTypes: idpBrowserSsoJitProvisioningUserAttributesAttrTypes},
		"user_repository": types.ObjectType{AttrTypes: datastorerepository.ElemAttrType()},
	}
	idpBrowserSsoMessageCustomizationsAttrTypes = map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}
	idpBrowserSsoMessageCustomizationsElementType              = types.ObjectType{AttrTypes: idpBrowserSsoMessageCustomizationsAttrTypes}
	idpBrowserSsoOauthAuthenticationPolicyContractRefAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
	idpBrowserSsoOidcProviderSettingsRequestParametersAttributeValueSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	idpBrowserSsoOidcProviderSettingsRequestParametersAttributeValueAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: idpBrowserSsoOidcProviderSettingsRequestParametersAttributeValueSourceAttrTypes},
		"value":  types.StringType,
	}
	idpBrowserSsoOidcProviderSettingsRequestParametersAttrTypes = map[string]attr.Type{
		"application_endpoint_override": types.BoolType,
		"attribute_value":               types.ObjectType{AttrTypes: idpBrowserSsoOidcProviderSettingsRequestParametersAttributeValueAttrTypes},
		"name":                          types.StringType,
		"value":                         types.StringType,
	}
	idpBrowserSsoOidcProviderSettingsRequestParametersElementType = types.ObjectType{AttrTypes: idpBrowserSsoOidcProviderSettingsRequestParametersAttrTypes}
	idpBrowserSsoOidcProviderSettingsAttrTypes                    = map[string]attr.Type{
		"authentication_scheme":                        types.StringType,
		"authentication_signing_algorithm":             types.StringType,
		"authorization_endpoint":                       types.StringType,
		"back_channel_logout_uri":                      types.StringType,
		"enable_pkce":                                  types.BoolType,
		"front_channel_logout_uri":                     types.StringType,
		"jwks_url":                                     types.StringType,
		"jwt_secured_authorization_response_mode_type": types.StringType,
		"login_type":                                   types.StringType,
		"logout_endpoint":                              types.StringType,
		"post_logout_redirect_uri":                     types.StringType,
		"pushed_authorization_request_endpoint":        types.StringType,
		"redirect_uri":                                 types.StringType,
		"request_parameters":                           types.SetType{ElemType: idpBrowserSsoOidcProviderSettingsRequestParametersElementType},
		"request_signing_algorithm":                    types.StringType,
		"scopes":                                       types.StringType,
		"token_endpoint":                               types.StringType,
		"track_user_sessions_for_logout":               types.BoolType,
		"user_info_endpoint":                           types.StringType,
	}
	idpBrowserSsoSloServiceEndpointsAttrTypes = map[string]attr.Type{
		"binding":      types.StringType,
		"response_url": types.StringType,
		"url":          types.StringType,
	}
	idpBrowserSsoSloServiceEndpointsElementType                             = types.ObjectType{AttrTypes: idpBrowserSsoSloServiceEndpointsAttrTypes}
	idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentElementType = types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentAttrTypes}

	idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaAttrTypes = map[string]attr.Type{
		"attribute_name": types.StringType,
		"condition":      types.StringType,
		"error_result":   types.StringType,
		"source":         types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaSourceAttrTypes},
		"value":          types.StringType,
	}
	idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaElementType = types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaAttrTypes}
	idpBrowserSsoSsoOauthMappingIssuanceCriteriaExpressionCriteriaAttrTypes    = map[string]attr.Type{
		"error_result": types.StringType,
		"expression":   types.StringType,
	}
	idpBrowserSsoSsoOauthMappingIssuanceCriteriaExpressionCriteriaElementType = types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingIssuanceCriteriaExpressionCriteriaAttrTypes}
	idpBrowserSsoSsoOauthMappingIssuanceCriteriaAttrTypes                     = map[string]attr.Type{
		"conditional_criteria": types.SetType{ElemType: idpBrowserSsoSsoOauthMappingIssuanceCriteriaConditionalCriteriaElementType},
		"expression_criteria":  types.SetType{ElemType: idpBrowserSsoSsoOauthMappingIssuanceCriteriaExpressionCriteriaElementType},
	}
	idpBrowserSsoSsoOauthMappingAttrTypes = map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"issuance_criteria":              types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingIssuanceCriteriaAttrTypes},
	}
	idpBrowserSsoSsoServiceEndpointsAttrTypes = map[string]attr.Type{
		"binding": types.StringType,
		"url":     types.StringType,
	}
	idpBrowserSsoSsoServiceEndpointsElementType = types.ObjectType{AttrTypes: idpBrowserSsoSsoServiceEndpointsAttrTypes}
	idpBrowserSsoUrlWhitelistEntriesAttrTypes   = map[string]attr.Type{
		"allow_query_and_fragment": types.BoolType,
		"require_https":            types.BoolType,
		"valid_domain":             types.StringType,
		"valid_path":               types.StringType,
	}

	idpBrowserSsoUrlWhitelistEntriesElementType = types.ObjectType{AttrTypes: idpBrowserSsoUrlWhitelistEntriesAttrTypes}
	idpBrowserSsoAttrTypes                      = map[string]attr.Type{
		"adapter_mappings":                         types.ListType{ElemType: idpBrowserSsoAdapterMappingsElementType},
		"always_sign_artifact_response":            types.BoolType,
		"artifact":                                 types.ObjectType{AttrTypes: idpBrowserSsoArtifactAttrTypes},
		"assertions_signed":                        types.BoolType,
		"attribute_contract":                       types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractAttrTypes},
		"authentication_policy_contract_mappings":  types.ListType{ElemType: idpBrowserSsoAuthenticationPolicyContractMappingsElementType},
		"authn_context_mappings":                   types.SetType{ElemType: idpBrowserSsoAuthnContextMappingsElementType},
		"decryption_policy":                        types.ObjectType{AttrTypes: idpBrowserSsoDecryptionPolicyAttrTypes},
		"default_target_url":                       types.StringType,
		"enabled_profiles":                         types.SetType{ElemType: types.StringType},
		"idp_identity_mapping":                     types.StringType,
		"incoming_bindings":                        types.SetType{ElemType: types.StringType},
		"jit_provisioning":                         types.ObjectType{AttrTypes: idpBrowserSsoJitProvisioningAttrTypes},
		"message_customizations":                   types.SetType{ElemType: idpBrowserSsoMessageCustomizationsElementType},
		"oauth_authentication_policy_contract_ref": types.ObjectType{AttrTypes: idpBrowserSsoOauthAuthenticationPolicyContractRefAttrTypes},
		"oidc_provider_settings":                   types.ObjectType{AttrTypes: idpBrowserSsoOidcProviderSettingsAttrTypes},
		"protocol":                                 types.StringType,
		"sign_authn_requests":                      types.BoolType,
		"slo_service_endpoints":                    types.SetType{ElemType: idpBrowserSsoSloServiceEndpointsElementType},
		"sso_application_endpoint":                 types.StringType,
		"sso_oauth_mapping":                        types.ObjectType{AttrTypes: idpBrowserSsoSsoOauthMappingAttrTypes},
		"sso_service_endpoints":                    types.SetType{ElemType: idpBrowserSsoSsoServiceEndpointsElementType},
		"url_whitelist_entries":                    types.SetType{ElemType: idpBrowserSsoUrlWhitelistEntriesElementType},
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
		"name_mappings": types.SetType{ElemType: attributeQueryNameMappingAttrTypes},
		"policy":        idpAttributeQueryPolicyAttrTypes,
	}

	spTokenGeneratorMappingAttrTypes = map[string]attr.Type{
		"sp_token_generator_ref":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"restricted_virtual_entity_ids":  types.SetType{ElemType: types.StringType},
		"default_mapping":                types.BoolType,
		"attribute_sources":              types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrTypes()},
	}

	wsTrustAttrTypes = map[string]attr.Type{
		"attribute_contract": types.ObjectType{AttrTypes: map[string]attr.Type{
			"core_attributes": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":   types.StringType,
					"masked": types.BoolType,
				},
			}},
			"extended_attributes": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":   types.StringType,
					"masked": types.BoolType,
				},
			}},
		}},
		"generate_local_token":     types.BoolType,
		"token_generator_mappings": types.SetType{ElemType: types.ObjectType{AttrTypes: spTokenGeneratorMappingAttrTypes}},
	}

	// inbound_provisioning
	inboundProvisioningCustomSchemaAttributesAttrTypes = map[string]attr.Type{
		"multi_valued":   types.BoolType,
		"name":           types.StringType,
		"sub_attributes": types.SetType{ElemType: types.StringType},
		"types":          types.SetType{ElemType: types.StringType},
	}
	inboundProvisioningCustomSchemaAttributesElementType = types.ObjectType{AttrTypes: inboundProvisioningCustomSchemaAttributesAttrTypes}
	inboundProvisioningCustomSchemaAttrTypes             = map[string]attr.Type{
		"attributes": types.SetType{ElemType: inboundProvisioningCustomSchemaAttributesElementType},
		"namespace":  types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesElementType   = types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesAttrTypes}
	inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesElementType = types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesAttrTypes}
	inboundProvisioningGroupsReadGroupsAttributeContractAttrTypes                     = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesElementType},
	}
	inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributeFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributeFulfillmentElementType = types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributeFulfillmentAttrTypes}
	inboundProvisioningGroupsReadGroupsAttributesAttrTypes             = map[string]attr.Type{
		"name": types.StringType,
	}
	inboundProvisioningGroupsReadGroupsAttributesElementType = types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributesAttrTypes}
	inboundProvisioningGroupsReadGroupsAttrTypes             = map[string]attr.Type{
		"attribute_contract":    types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttributeContractAttrTypes},
		"attribute_fulfillment": types.MapType{ElemType: inboundProvisioningGroupsReadGroupsAttributeFulfillmentElementType},
		"attributes":            types.SetType{ElemType: inboundProvisioningGroupsReadGroupsAttributesElementType},
	}
	inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	inboundProvisioningGroupsWriteGroupsAttributeFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	inboundProvisioningGroupsWriteGroupsAttributeFulfillmentElementType = types.ObjectType{AttrTypes: inboundProvisioningGroupsWriteGroupsAttributeFulfillmentAttrTypes}
	inboundProvisioningGroupsWriteGroupsAttrTypes                       = map[string]attr.Type{
		"attribute_fulfillment": types.MapType{ElemType: inboundProvisioningGroupsWriteGroupsAttributeFulfillmentElementType},
	}
	inboundProvisioningGroupsAttrTypes = map[string]attr.Type{
		"read_groups":  types.ObjectType{AttrTypes: inboundProvisioningGroupsReadGroupsAttrTypes},
		"write_groups": types.ObjectType{AttrTypes: inboundProvisioningGroupsWriteGroupsAttrTypes},
	}

	inboundProvisioningUsersReadUsersAttributeContractCoreAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	inboundProvisioningUsersReadUsersAttributeContractCoreAttributesElementType   = types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributeContractCoreAttributesAttrTypes}
	inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesAttrTypes = map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesElementType = types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesAttrTypes}
	inboundProvisioningUsersReadUsersAttributeContractAttrTypes                     = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: inboundProvisioningUsersReadUsersAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesElementType},
	}
	inboundProvisioningUsersReadUsersAttributeFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	inboundProvisioningUsersReadUsersAttributeFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributeFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	inboundProvisioningUsersReadUsersAttributeFulfillmentElementType = types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributeFulfillmentAttrTypes}
	inboundProvisioningUsersReadUsersAttributesAttrTypes             = map[string]attr.Type{
		"name": types.StringType,
	}
	inboundProvisioningUsersReadUsersAttributesElementType = types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributesAttrTypes}
	inboundProvisioningUsersReadUsersAttrTypes             = map[string]attr.Type{
		"attribute_contract":    types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttributeContractAttrTypes},
		"attribute_fulfillment": types.MapType{ElemType: inboundProvisioningUsersReadUsersAttributeFulfillmentElementType},
		"attributes":            types.SetType{ElemType: inboundProvisioningUsersReadUsersAttributesElementType},
	}
	inboundProvisioningUsersWriteUsersAttributeFulfillmentSourceAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
	inboundProvisioningUsersWriteUsersAttributeFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: inboundProvisioningUsersWriteUsersAttributeFulfillmentSourceAttrTypes},
		"value":  types.StringType,
	}
	inboundProvisioningUsersWriteUsersAttributeFulfillmentElementType = types.ObjectType{AttrTypes: inboundProvisioningUsersWriteUsersAttributeFulfillmentAttrTypes}
	inboundProvisioningUsersWriteUsersAttrTypes                       = map[string]attr.Type{
		"attribute_fulfillment": types.MapType{ElemType: inboundProvisioningUsersWriteUsersAttributeFulfillmentElementType},
	}
	inboundProvisioningUsersAttrTypes = map[string]attr.Type{
		"read_users":  types.ObjectType{AttrTypes: inboundProvisioningUsersReadUsersAttrTypes},
		"write_users": types.ObjectType{AttrTypes: inboundProvisioningUsersWriteUsersAttrTypes},
	}
	inboundProvisioningAttrTypes = map[string]attr.Type{
		"action_on_delete": types.StringType,
		"custom_schema":    types.ObjectType{AttrTypes: inboundProvisioningCustomSchemaAttrTypes},
		"group_support":    types.BoolType,
		"groups":           types.ObjectType{AttrTypes: inboundProvisioningGroupsAttrTypes},
		"user_repository":  types.ObjectType{AttrTypes: inboundprovisioninguserrepository.ElemAttrType()},
		"users":            types.ObjectType{AttrTypes: inboundProvisioningUsersAttrTypes},
	}

	tokenGeneratorAttrTypes = map[string]attr.Type{
		"sp_token_generator_ref":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"attribute_sources":              types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"default_mapping":                types.BoolType,
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrTypes()},
		"restricted_virtual_entity_ids":  types.SetType{ElemType: types.StringType},
	}

	customId = "connection_id"
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
	VirtualEntityIds                       types.Set    `tfsdk:"virtual_entity_ids"`
	WsTrust                                types.Object `tfsdk:"ws_trust"`
}

// GetSchema defines the schema for the resource.
func (r *spIdpConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a partner Identity Provider connection",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies whether the connection is active and ready to process incoming requests. The default value is `false`.",
				MarkdownDescription: "Specifies whether the connection is active and ready to process incoming requests. The default value is `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Additional allowed entities or issuers configuration. Currently only used in OIDC IdP (RP) connection.",
				Attributes: map[string]schema.Attribute{
					"allow_additional_entities": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation. The default value is `false`.",
						MarkdownDescription: "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation. The default value is `false`.",
					},
					"allow_all_entities": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Set to true to accept any entity or issuer during entity or issuer validation. (Not Recommended). The default value is `false`.",
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"entity_description": schema.StringAttribute{
									Optional:            true,
									Description:         "Entity description.",
									MarkdownDescription: "Entity description.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
			"attribute_query": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"name_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"local_name": schema.StringAttribute{
									Required:            true,
									Description:         "The local attribute name.",
									MarkdownDescription: "The local attribute name.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"remote_name": schema.StringAttribute{
									Required:            true,
									Description:         "The remote attribute name as defined by the attribute authority.",
									MarkdownDescription: "The remote attribute name as defined by the attribute authority.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
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
								Description:         "Require signed r.",
								MarkdownDescription: "Require signed r.",
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:            true,
				Description:         "The attribute query profile supports local applications in requesting user attributes from an attribute authority.",
				MarkdownDescription: "The attribute query profile supports local applications in requesting user attributes from an attribute authority.",
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"connection_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The persistent, unique ID for the connection. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.",
				Validators: []validator.String{
					configvalidators.PingFederateId(),
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contact_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Optional:            true,
						Description:         "Company name.",
						MarkdownDescription: "Company name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact email address.",
						MarkdownDescription: "Contact email address.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"phone": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact phone number.",
						MarkdownDescription: "Contact phone number.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"first_name": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact first name.",
						MarkdownDescription: "Contact first name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"last_name": schema.StringAttribute{
						Optional:            true,
						Description:         "Contact last name.",
						MarkdownDescription: "Contact last name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
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
						Description:         "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
						MarkdownDescription: "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"verification_subject_dn": schema.StringAttribute{
						Optional:            true,
						Description:         "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the `certs` array.",
						MarkdownDescription: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the `certs` array.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"certs": connectioncert.ToSchema("The certificates used for signature verification and XML encryption.", true),
					"block_encryption_algorithm": schema.StringAttribute{
						Optional:            true,
						Description:         "The algorithm used to encrypt assertions sent to this partner. Options are `AES_128`, `AES_256`, `AES_128_GCM`, `AES_192_GCM`, `AES_256_GCM`, `Triple_DES`.",
						MarkdownDescription: "The algorithm used to encrypt assertions sent to this partner. Options are `AES_128`, `AES_256`, `AES_128_GCM`, `AES_192_GCM`, `AES_256_GCM`, `Triple_DES`.",
						Validators: []validator.String{
							stringvalidator.OneOf("AES_128", "AES_256", "AES_128_GCM", "AES_192_GCM", "AES_256_GCM", "Triple_DES"),
						},
					},
					"key_transport_algorithm": schema.StringAttribute{
						Optional:            true,
						Description:         "The algorithm used to transport keys to this partner. Options are `RSA_OAEP`, `RSA_OAEP_256`, `RSA_v15`.",
						MarkdownDescription: "The algorithm used to transport keys to this partner. Options are `RSA_OAEP`, `RSA_OAEP_256`, `RSA_v15`.",
						Validators: []validator.String{
							stringvalidator.OneOf("RSA_OAEP", "RSA_OAEP_256", "RSA_v15"),
						},
					},
					"signing_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"signing_key_pair_ref": schema.SingleNestedAttribute{
								Attributes:          resourcelink.ToSchema(),
								Optional:            true,
								Description:         "A reference to a signing key pair.",
								MarkdownDescription: "A reference to a signing key pair.",
							},
							"alternative_signing_key_pair_refs": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: resourcelink.ToSchema(),
								},
								Optional:            true,
								Description:         "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
								MarkdownDescription: "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
							},
							"algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "The algorithm used to sign messages sent to this partner. The default is `SHA1withDSA` for DSA certs, `SHA256withRSA` for RSA certs, and `SHA256withECDSA` for EC certs. For RSA certs, `SHA1withRSA`, `SHA384withRSA`, `SHA512withRSA`, `SHA256withRSAandMGF1`, `SHA384withRSAandMGF1` and `SHA512withRSAandMGF1` are also supported. For EC certs, `SHA384withECDSA` and `SHA512withECDSA` are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
								MarkdownDescription: "The algorithm used to sign messages sent to this partner. The default is `SHA1withDSA` for DSA certs, `SHA256withRSA` for RSA certs, and `SHA256withECDSA` for EC certs. For RSA certs, `SHA1withRSA`, `SHA384withRSA`, `SHA512withRSA`, `SHA256withRSAandMGF1`, `SHA384withRSAandMGF1` and `SHA512withRSAandMGF1` are also supported. For EC certs, `SHA384withECDSA` and `SHA512withECDSA` are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"include_cert_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "Determines whether the signing certificate is included in the signature <KeyInfo> element. The default value is `false`.",
								MarkdownDescription: "Determines whether the signing certificate is included in the signature <KeyInfo> element. The default value is `false`.",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element.",
								MarkdownDescription: "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element.",
							},
						},
						Optional:            true,
						Description:         "Settings related to signing messages sent to this partner.",
						MarkdownDescription: "Settings related to signing messages sent to this partner.",
					},
					"decryption_key_pair_ref": schema.SingleNestedAttribute{
						Attributes:          resourcelink.ToSchema(),
						Optional:            true,
						Description:         "A reference to a resource.",
						MarkdownDescription: "A reference to a resource.",
					},
					"secondary_decryption_key_pair_ref": schema.SingleNestedAttribute{
						Attributes:          resourcelink.ToSchema(),
						Optional:            true,
						Description:         "A reference to a resource.",
						MarkdownDescription: "A reference to a resource.",
					},
					"outbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Optional:            true,
										Description:         "The username.",
										MarkdownDescription: "The username.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"password": schema.StringAttribute{
										Optional:            true,
										Sensitive:           true,
										Description:         "User password. Either this attribute or `encrypted_password` must be specified.",
										MarkdownDescription: "User password. Either this attribute or `encrypted_password` must be specified.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"encrypted_password": schema.StringAttribute{
										Description:         "Encrypted user password. Either this attribute or `password` must be specified.",
										MarkdownDescription: "Encrypted user password. Either this attribute or `password` must be specified.",
										Optional:            true,
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("password")),
										},
									},
								},
								Optional:            true,
								Description:         "Username and password credentials.",
								MarkdownDescription: "Username and password credentials.",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "If incoming or outgoing messages must be signed. The default value is `false`.",
								MarkdownDescription: "If incoming or outgoing messages must be signed. The default value is `false`.",
							},
							"ssl_auth_key_pair_ref": schema.SingleNestedAttribute{
								Attributes:          resourcelink.ToSchema(),
								Optional:            true,
								Description:         "A reference to a resource.",
								MarkdownDescription: "A reference to a resource.",
							},
							"validate_partner_cert": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Validate the partner server certificate. Default is `true`.",
								MarkdownDescription: "Validate the partner server certificate. Default is `true`.",
								Default:             booldefault.StaticBool(true),
							},
						},
						Optional:            true,
						Description:         "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
						MarkdownDescription: "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
					},
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Optional:            true,
										Description:         "The username.",
										MarkdownDescription: "The username.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"password": schema.StringAttribute{
										Optional:            true,
										Sensitive:           true,
										Description:         "User password. Either this attribute or `encrypted_password` must be specified.",
										MarkdownDescription: "User password. Either this attribute or `encrypted_password` must be specified.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"encrypted_password": schema.StringAttribute{
										Description:         "Encrypted user password. Either this attribute or `password` must be specified.",
										MarkdownDescription: "Encrypted user password. Either this attribute or `password` must be specified.",
										Optional:            true,
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("password")),
										},
									},
								},
								Optional:            true,
								Description:         "Username and password credentials.",
								MarkdownDescription: "Username and password credentials.",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "If incoming or outgoing messages must be signed. The default value is `false`.",
								MarkdownDescription: "If incoming or outgoing messages must be signed. The default value is `false`.",
							},
							"verification_subject_dn": schema.StringAttribute{
								Optional:            true,
								Description:         "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the `certs` array.",
								MarkdownDescription: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the `certs` array.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"verification_issuer_dn": schema.StringAttribute{
								Optional:            true,
								Description:         "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
								MarkdownDescription: "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"certs": connectioncert.ToSchema("The certificates used for signature verification and XML encryption.", false),
							"require_ssl": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Incoming HTTP transmissions must use a secure channel. The default value is `false`.",
								MarkdownDescription: "Incoming HTTP transmissions must use a secure channel. The default value is `false`.",
							},
						},
						Optional:            true,
						Description:         "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
						MarkdownDescription: "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
					},
				},
				Optional:            true,
				Description:         "The certificates and settings for encryption, signing, and signature verification.",
				MarkdownDescription: "The certificates and settings for encryption, signing, and signature verification.",
			},
			"default_virtual_entity_id": schema.StringAttribute{
				Optional:            true,
				Description:         "The default alternate entity ID that identifies the local server to this partner. It is required when `virtual_entity_ids` is not empty and must be included in that list.",
				MarkdownDescription: "The default alternate entity ID that identifies the local server to this partner. It is required when `virtual_entity_ids` is not empty and must be included in that list.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"entity_id": schema.StringAttribute{
				Required:            true,
				Description:         "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
				MarkdownDescription: "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"error_page_msg_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Identifier that specifies the message displayed on a user-facing error page. Defaults to `errorDetail.spSsoFailure` for browser SSO connections, null otherwise.",
				MarkdownDescription: "Identifier that specifies the message displayed on a user-facing error page. Defaults to `errorDetail.spSsoFailure` for browser SSO connections, null otherwise.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"extended_properties": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "A Set of values",
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
								"adapter_override_settings": schema.SingleNestedAttribute{
									Optional:            true,
									Description:         "An SP adapter instance.",
									MarkdownDescription: "An SP adapter instance.",
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required:            true,
											Description:         "The ID of the plugin instance. The ID cannot be modified once the instance is created.",
											MarkdownDescription: "The ID of the plugin instance. The ID cannot be modified once the instance is created.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"name": schema.StringAttribute{
											Optional:            true,
											Description:         "The plugin instance name.",
											MarkdownDescription: "The plugin instance name.",
										},
										"plugin_descriptor_ref": schema.SingleNestedAttribute{
											Attributes:          resourcelink.ToSchemaNoLengthValidator(),
											Optional:            true,
											Description:         "Reference to the plugin descriptor for this instance.",
											MarkdownDescription: "Reference to the plugin descriptor for this instance.",
										},
										"parent_ref": schema.SingleNestedAttribute{
											Attributes:          resourcelink.ToSchema(),
											Optional:            true,
											Description:         "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances.",
											MarkdownDescription: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances.",
										},
										"configuration": pluginconfiguration.ToSchema(),
										"attribute_contract": schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{
												"core_attributes": schema.SetNestedAttribute{
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"name": schema.StringAttribute{
																Optional:            false,
																Computed:            true,
																Description:         "The name of this attribute.",
																MarkdownDescription: "The name of this attribute.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
														},
													},
													Optional:            false,
													Computed:            true,
													Description:         "A list of read-only assertion attributes that are automatically populated by the SP adapter descriptor.",
													MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by the SP adapter descriptor.",
													PlanModifiers: []planmodifier.Set{
														setplanmodifier.UseStateForUnknown(),
													},
												},
												"extended_attributes": schema.SetNestedAttribute{
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"name": schema.StringAttribute{
																Required:            true,
																Description:         "The name of this attribute.",
																MarkdownDescription: "The name of this attribute.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
														},
													},
													Optional:            true,
													Description:         "A list of additional attributes that can be returned by the SP adapter.",
													MarkdownDescription: "A list of additional attributes that can be returned by the SP adapter.",
												},
											},
										},
										"target_application_info": schema.SingleNestedAttribute{
											Optional: true,
											Computed: true,
											Default: objectdefault.StaticValue(types.ObjectValueMust(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes, map[string]attr.Value{
												"application_name":     types.StringNull(),
												"application_icon_url": types.StringNull(),
											})),
											Attributes: map[string]schema.Attribute{
												"application_name": schema.StringAttribute{
													Optional:            true,
													Description:         "The application name.",
													MarkdownDescription: "The application name.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"application_icon_url": schema.StringAttribute{
													Optional:            true,
													Description:         "The application icon URL.",
													MarkdownDescription: "The application icon URL.",
												},
											},
										},
									},
								},
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
								"attribute_sources":              attributesources.ToSchema(0, false),
								"issuance_criteria":              issuancecriteria.ToSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									Description:         "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
									MarkdownDescription: "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
								},
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
								"sp_adapter_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchemaNoLengthValidator(),
									Optional:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
							},
						},
						Optional:            true,
						Description:         "A list of adapters that map to incoming assertions.",
						MarkdownDescription: "A list of adapters that map to incoming assertions.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
					"always_sign_artifact_response": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						Description:         "Specify to always sign the SAML ArtifactResponse. Default is `false`.",
						MarkdownDescription: "Specify to always sign the SAML ArtifactResponse. Default is `false`.",
					},
					"artifact": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"lifetime": schema.Int64Attribute{
								Optional:            true,
								Description:         "The lifetime of the artifact in seconds.",
								MarkdownDescription: "The lifetime of the artifact in seconds.",
							},
							"resolver_locations": schema.SetNestedAttribute{
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
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
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
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
						Optional:            true,
						Description:         "The settings for an Artifact binding.",
						MarkdownDescription: "The settings for an Artifact binding.",
					},
					"assertions_signed": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         "Specify whether the incoming SAML assertions are signed rather than the entire SAML response being signed. The default value is `false`.",
						MarkdownDescription: "Specify whether the incoming SAML assertions are signed rather than the entire SAML response being signed. The default value is `false`.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            false,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            false,
								Computed:            true,
								Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
							"extended_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            true,
								Computed:            true,
								Description:         "A list of additional attributes that are present in the incoming assertion.",
								MarkdownDescription: "A list of additional attributes that are present in the incoming assertion.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractExtendedAttributesAttrTypes}, nil)),
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "A set of user attributes that the IdP sends in the SAML assertion.",
						MarkdownDescription: "A set of user attributes that the IdP sends in the SAML assertion.",
					},
					"authentication_policy_contract_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
								"attribute_sources":              attributesources.ToSchema(0, false),
								"authentication_policy_contract_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchema(),
									Required:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
								"issuance_criteria": issuancecriteria.ToSchema(),
								"restrict_virtual_server_ids": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									Description:         "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
									MarkdownDescription: "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
								},
								"restricted_virtual_server_ids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
									Description:         "The list of virtual server IDs that this mapping is restricted to. The default value is an empty set.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to. The default value is an empty set.",
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Default:             listdefault.StaticValue(types.ListValueMust(idpBrowserSsoAuthenticationPolicyContractMappingsElementType, nil)),
						Description:         "A list of Authentication Policy Contracts that map to incoming assertions.",
						MarkdownDescription: "A list of Authentication Policy Contracts that map to incoming assertions.",
					},
					"authn_context_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"local": schema.StringAttribute{
									Optional:            true,
									Description:         "The local authentication context value.",
									MarkdownDescription: "The local authentication context value.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"remote": schema.StringAttribute{
									Optional:            true,
									Description:         "The remote authentication context value.",
									MarkdownDescription: "The remote authentication context value.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
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
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Specify whether the incoming SAML assertion is encrypted for an IdP connection. The default value is `false`.",
								MarkdownDescription: "Specify whether the incoming SAML assertion is encrypted for an IdP connection. The default value is `false`.",
							},
							"attributes_encrypted": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Specify whether one or more incoming SAML attributes are encrypted for an IdP connection. The default value is `false`.",
								MarkdownDescription: "Specify whether one or more incoming SAML attributes are encrypted for an IdP connection. The default value is `false`.",
							},
							"slo_encrypt_subject_name_id": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Encrypt the Subject Name ID in SLO messages to the IdP. The default value is `false`.",
								MarkdownDescription: "Encrypt the Subject Name ID in SLO messages to the IdP. The default value is `false`.",
							},
							"slo_subject_name_id_encrypted": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Allow encrypted Subject Name ID in SLO messages from the IdP. The default value is `false`.",
								MarkdownDescription: "Allow encrypted Subject Name ID in SLO messages from the IdP. The default value is `false`.",
							},
							"subject_name_id_encrypted": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Specify whether the incoming Subject Name ID is encrypted for an IdP connection. The default value is `false`.",
								MarkdownDescription: "Specify whether the incoming Subject Name ID is encrypted for an IdP connection. The default value is `false`.",
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Defines what to decrypt in the browser-based SSO profile.",
						MarkdownDescription: "Defines what to decrypt in the browser-based SSO profile.",
					},
					"default_target_url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Description:         "The default target URL for this connection. If defined, this overrides the default URL. The default value is an empty string.",
						MarkdownDescription: "The default target URL for this connection. If defined, this overrides the default URL. The default value is an empty string.",
					},
					"enabled_profiles": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						MarkdownDescription: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
					},
					"idp_identity_mapping": schema.StringAttribute{
						Required:            true,
						Description:         "Defines the process in which users authenticated by the IdP are associated with user accounts local to the SP. Options are `ACCOUNT_LINKING`, `ACCOUNT_MAPPING`, `NONE`.",
						MarkdownDescription: "Defines the process in which users authenticated by the IdP are associated with user accounts local to the SP. Options are `ACCOUNT_LINKING`, `ACCOUNT_MAPPING`, `NONE`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ACCOUNT_MAPPING",
								"ACCOUNT_LINKING",
								"NONE",
							),
						},
					},
					"incoming_bindings": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections and it is optional for IdP Connections.",
						MarkdownDescription: "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections and it is optional for IdP Connections.",
					},
					"jit_provisioning": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"error_handling": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Specify behavior when provisioning request fails. The default is `CONTINUE_SSO`. Options are `ABORT_SSO`, `CONTINUE_SSO`.",
								MarkdownDescription: "Specify behavior when provisioning request fails. The default is `CONTINUE_SSO`. Options are `ABORT_SSO`, `CONTINUE_SSO`.",
								Default:             stringdefault.StaticString("CONTINUE_SSO"),
								Validators: []validator.String{
									stringvalidator.OneOf(
										"CONTINUE_SSO",
										"ABORT_SSO",
									),
								},
							},
							"event_trigger": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Specify when provisioning occurs during assertion processing. The default is `NEW_USER_ONLY`. Options are `ALL_SAML_ASSERTIONS`, `NEW_USER_ONLY`.",
								MarkdownDescription: "Specify when provisioning occurs during assertion processing. The default is `NEW_USER_ONLY`. Options are `ALL_SAML_ASSERTIONS`, `NEW_USER_ONLY`.",
								Default:             stringdefault.StaticString("NEW_USER_ONLY"),
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NEW_USER_ONLY",
										"ALL_SAML_ASSERTIONS",
									),
								},
							},
							"user_attributes": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"masked": schema.BoolAttribute{
													Computed:            true,
													Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
													MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
												},
												"name": schema.StringAttribute{
													Computed:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
											},
										},
										Computed: true,
										PlanModifiers: []planmodifier.Set{
											setplanmodifier.UseStateForUnknown(),
										},
										Description:         "A list of user attributes that the IdP sends in the SAML assertion.",
										MarkdownDescription: "A list of user attributes that the IdP sends in the SAML assertion.",
									},
									"do_attribute_query": schema.BoolAttribute{
										Optional:            true,
										Computed:            true,
										Default:             booldefault.StaticBool(false),
										Description:         "Specify whether to use only attributes from the SAML Assertion or retrieve additional attributes from the IdP. The default is `false`.",
										MarkdownDescription: "Specify whether to use only attributes from the SAML Assertion or retrieve additional attributes from the IdP. The default is `false`.",
									},
								},
								Required: true,
							},
							"user_repository": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"jdbc": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"data_store_ref": schema.SingleNestedAttribute{
												Attributes:          resourcelink.ToSchema(),
												Required:            true,
												Description:         "Reference to the associated data store.",
												MarkdownDescription: "Reference to the associated data store.",
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
																	Validators: []validator.String{
																		stringvalidator.LengthAtLeast(1),
																	},
																},
																"type": schema.StringAttribute{
																	Required:            true,
																	Description:         "The source type of this key. Options are `ACCOUNT_LINK`, `ACTOR_TOKEN`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
																	MarkdownDescription: "The source type of this key. Options are `ACCOUNT_LINK`, `ACTOR_TOKEN`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
															Optional:            true,
															Computed:            true,
															Default:             stringdefault.StaticString(""),
															Description:         "The value for this attribute.",
															MarkdownDescription: "The value for this attribute.",
														},
													},
												},
												Required:            true,
												Description:         "The user repository attribute mapping.",
												MarkdownDescription: "The user repository attribute mapping.",
											},
											"sql_method": schema.SingleNestedAttribute{
												Description:         "The method to map attributes from the assertion directly to database table columns or to stored-procedure parameters.",
												MarkdownDescription: "The method to map attributes from the assertion directly to database table columns or to stored-procedure parameters.",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"table": schema.SingleNestedAttribute{
														Optional:            true,
														Description:         "The Table SQL method.",
														MarkdownDescription: "The Table SQL method.",
														Attributes: map[string]schema.Attribute{
															"schema": schema.StringAttribute{
																Required:            true,
																Description:         "Lists the table structure that stores information within a database.",
																MarkdownDescription: "Lists the table structure that stores information within a database.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
															"table_name": schema.StringAttribute{
																Required:            true,
																Description:         "The name of the database table.",
																MarkdownDescription: "The name of the database table.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
															"unique_id_column": schema.StringAttribute{
																Required:            true,
																Description:         "The database column that uniquely identifies the provisioned user on the SP side.",
																MarkdownDescription: "The database column that uniquely identifies the provisioned user on the SP side.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
														},
													},
													"stored_procedure": schema.SingleNestedAttribute{
														Description: "The Stored Procedure SQL method. The procedure is always called for all SSO tokens and `event_trigger` will always be `ALL_SAML_ASSERTIONS`.",
														Optional:    true,
														Attributes: map[string]schema.Attribute{
															"schema": schema.StringAttribute{
																Required:            true,
																Description:         "Lists the table structure that stores information within a database.",
																MarkdownDescription: "Lists the table structure that stores information within a database.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
															"stored_procedure": schema.StringAttribute{
																Required:            true,
																Description:         "The name of the database stored procedure.",
																MarkdownDescription: "The name of the database stored procedure.",
																Validators: []validator.String{
																	stringvalidator.LengthAtLeast(1),
																},
															},
														},
													},
												},
											},
										},
										Optional:            true,
										Description:         "JDBC data store user repository.",
										MarkdownDescription: "JDBC data store user repository.",
									},
									"ldap": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"data_store_ref": schema.SingleNestedAttribute{
												Attributes:          resourcelink.ToSchema(),
												Description:         "Reference to the associated data store.",
												MarkdownDescription: "Reference to the associated data store.",
												Required:            true,
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
																	Validators: []validator.String{
																		stringvalidator.LengthAtLeast(1),
																	},
																},
																"type": schema.StringAttribute{
																	Required:            true,
																	Description:         "The source type of this key. Options are `ACCOUNT_LINK`, `ACTOR_TOKEN`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
																	MarkdownDescription: "The source type of this key. Options are `ACCOUNT_LINK`, `ACTOR_TOKEN`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
															Optional:            true,
															Computed:            true,
															Default:             stringdefault.StaticString(""),
															Description:         "The value for this attribute.",
															MarkdownDescription: "The value for this attribute.",
														},
													},
												},
												Required:            true,
												Description:         "The user repository attribute mapping.",
												MarkdownDescription: "The user repository attribute mapping.",
											},
											"base_dn": schema.StringAttribute{
												Optional:            true,
												Description:         "The base DN to search from. If not specified, the search will start at the LDAP's root.",
												MarkdownDescription: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"unique_user_id_filter": schema.StringAttribute{
												Required:            true,
												Description:         "The expression that results in a unique user identifier, when combined with the Base DN.",
												MarkdownDescription: "The expression that results in a unique user identifier, when combined with the Base DN.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
										Optional: true,
										Validators: []validator.Object{
											objectvalidator.ExactlyOneOf(
												path.MatchRelative().AtParent().AtName("jdbc"),
												path.MatchRelative().AtParent().AtName("ldap"),
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
					"message_customizations": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"context_name": schema.StringAttribute{
									Optional:            true,
									Description:         "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be `assertion`, `authn-response` or `authn-request`.",
									MarkdownDescription: "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be `assertion`, `authn-response` or `authn-request`.",
									Validators: []validator.String{
										stringvalidator.OneOf("assertion", "authn-request", "authn-response"),
									},
								},
								"message_expression": schema.StringAttribute{
									Optional:            true,
									Description:         "The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.",
									MarkdownDescription: "The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Optional:            true,
						Description:         "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
						MarkdownDescription: "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"oauth_authentication_policy_contract_ref": schema.SingleNestedAttribute{
						Attributes:          resourcelink.ToSchema(),
						Optional:            true,
						Description:         "A reference to a resource.",
						MarkdownDescription: "A reference to a resource.",
					},
					"oidc_provider_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"authentication_scheme": schema.StringAttribute{
								Optional:            true,
								Description:         "The OpenID Connect Authentication Scheme. This is required for Authentication using Code Flow. Options are `BASIC`, `CLIENT_SECRET_JWT`, `POST`, `PRIVATE_KEY_JWT`.",
								MarkdownDescription: "The OpenID Connect Authentication Scheme. This is required for Authentication using Code Flow. Options are `BASIC`, `CLIENT_SECRET_JWT`, `POST`, `PRIVATE_KEY_JWT`.",
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
								Description:         "The authentication signing algorithm for token endpoint PRIVATE_KEY_JWT or CLIENT_SECRET_JWT authentication. Asymmetric algorithms are allowed for PRIVATE_KEY_JWT and symmetric algorithms are allowed for CLIENT_SECRET_JWT. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11. Options are `NONE`, `ES256`, `ES384`, `ES512`, `HS256`, `HS384`, `HS512`, `PS256`, `PS384`, `PS512` `RS256`, `RS384`, `RS512`.",
								MarkdownDescription: "The authentication signing algorithm for token endpoint PRIVATE_KEY_JWT or CLIENT_SECRET_JWT authentication. Asymmetric algorithms are allowed for PRIVATE_KEY_JWT and symmetric algorithms are allowed for CLIENT_SECRET_JWT. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11. Options are `NONE`, `ES256`, `ES384`, `ES512`, `HS256`, `HS384`, `HS512`, `PS256`, `PS384`, `PS512` `RS256`, `RS384`, `RS512`.",
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
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"back_channel_logout_uri": schema.StringAttribute{
								Computed:            true,
								Description:         "The Back-Channel Logout URI. This read-only parameter is available when user sessions are tracked for logout.",
								MarkdownDescription: "The Back-Channel Logout URI. This read-only parameter is available when user sessions are tracked for logout.",
							},
							"enable_pkce": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Enable Proof Key for Code Exchange (PKCE). When enabled, the client sends an SHA-256 code challenge and corresponding code verifier to the OpenID Provider during the authorization code flow. The default value is `false`.",
								MarkdownDescription: "Enable Proof Key for Code Exchange (PKCE). When enabled, the client sends an SHA-256 code challenge and corresponding code verifier to the OpenID Provider during the authorization code flow. The default value is `false`.",
							},
							"front_channel_logout_uri": schema.StringAttribute{
								Optional:            false,
								Computed:            true,
								Description:         "The Front-Channel Logout URI. This is a read-only parameter.",
								MarkdownDescription: "The Front-Channel Logout URI. This is a read-only parameter.",
							},
							"jwks_url": schema.StringAttribute{
								Required:            true,
								Description:         "URL of the OpenID Provider's JSON Web Key Set [JWK] document.",
								MarkdownDescription: "URL of the OpenID Provider's JSON Web Key Set [JWK] document.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"jwt_secured_authorization_response_mode_type": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "The OpenId Connect JWT Secured Authorization Response Mode (JARM). The supported values are: <br>  `DISABLED`: Authorization responses will not be encoded using JARM. This is the default value. <br>  `QUERY_JWT`: query.jwt <br> `FORM_POST_JWT`: form_post.jwt <br><br> Note: `QUERY_JWT` must not be used in conjunction with loginType POST or  POST_AT unless the response JWT is encrypted to prevent token leakage in the URL. Supported in PingFederate `12.1` and later.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									stringvalidator.OneOf("DISABLED", "QUERY_JWT", "FORM_POST_JWT"),
								},
							},
							"login_type": schema.StringAttribute{
								Required:            true,
								Description:         "The OpenID Connect login type. These values maps to: \n CODE: Authentication using Code Flow \n  POST: Authentication using Form Post \n  POST_AT: Authentication using Form Post with Access Token. Options are `CODE`, `POST`, `POST_AT`.",
								MarkdownDescription: "The OpenID Connect login type. These values maps to: \n CODE: Authentication using Code Flow \n  POST: Authentication using Form Post \n  POST_AT: Authentication using Form Post with Access Token. Options are `CODE`, `POST`, `POST_AT`.",
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
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"post_logout_redirect_uri": schema.StringAttribute{
								Computed:            true,
								Description:         "The Post-Logout Redirect URI, where the OpenID Provider may redirect the user when RP-Initiated Logout has completed. This is a read-only parameter.",
								MarkdownDescription: "The Post-Logout Redirect URI, where the OpenID Provider may redirect the user when RP-Initiated Logout has completed. This is a read-only parameter.",
							},
							"pushed_authorization_request_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's OAuth 2.0 Pushed Authorization Request Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's OAuth 2.0 Pushed Authorization Request Endpoint.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"redirect_uri": schema.StringAttribute{
								Computed:            true,
								Description:         "The redirect URI. This is a read-only parameter.",
								MarkdownDescription: "The redirect URI. This is a read-only parameter.",
							},
							"request_parameters": schema.SetNestedAttribute{
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
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
															MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
													Optional:            true,
													Computed:            true,
													Default:             stringdefault.StaticString(""),
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
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"value": schema.StringAttribute{
											Optional:            true,
											Description:         "A request parameter value. A parameter can have either a value or a attribute value but not both. Value set here will be converted to an attribute value of source type TEXT. An empty value will be converted to attribute value of source type NO_MAPPING.",
											MarkdownDescription: "A request parameter value. A parameter can have either a value or a attribute value but not both. Value set here will be converted to an attribute value of source type TEXT. An empty value will be converted to attribute value of source type NO_MAPPING.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            true,
								Computed:            true,
								Default:             setdefault.StaticValue(types.SetValueMust(idpBrowserSsoOidcProviderSettingsRequestParametersElementType, nil)),
								Description:         "A list of request parameters. Request parameters with same name but different attribute values are treated as a multi-valued request parameter.",
								MarkdownDescription: "A list of request parameters. Request parameters with same name but different attribute values are treated as a multi-valued request parameter.",
							},
							"request_signing_algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "The request signing algorithm. Required only if you wish to use signed requests. Only asymmetric algorithms are allowed. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11. Options are `ES256`, `ES384`, `ES512`, `HS256`, `HS384`, `HS512`, `NONE`, `PS256`, `PS384`, `PS512`, `RS256`, `RS384`, `RS512`.",
								MarkdownDescription: "The request signing algorithm. Required only if you wish to use signed requests. Only asymmetric algorithms are allowed. For RSASSA-PSS signing algorithm, PingFederate must be integrated with a hardware security module (HSM) or Java 11. Options are `ES256`, `ES384`, `ES512`, `HS256`, `HS384`, `HS512`, `NONE`, `PS256`, `PS384`, `PS512`, `RS256`, `RS384`, `RS512`.",
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
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"token_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's OAuth 2.0 Token Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's OAuth 2.0 Token Endpoint.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
							"track_user_sessions_for_logout": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Determines whether PingFederate tracks a logout entry when a user signs in, so that the user session can later be terminated via a logout request from the OP. This setting must also be enabled in order for PingFederate to send an RP-initiated logout request to the OP during SLO. Default value is `false`.",
								MarkdownDescription: "Determines whether PingFederate tracks a logout entry when a user signs in, so that the user session can later be terminated via a logout request from the OP. This setting must also be enabled in order for PingFederate to send an RP-initiated logout request to the OP during SLO. Default value is `false`.",
							},
							"user_info_endpoint": schema.StringAttribute{
								Optional:            true,
								Description:         "URL of the OpenID Provider's UserInfo Endpoint.",
								MarkdownDescription: "URL of the OpenID Provider's UserInfo Endpoint.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
						Optional: true,

						Description:         "The OpenID Provider settings.",
						MarkdownDescription: "The OpenID Provider settings.",
					},
					"protocol": schema.StringAttribute{
						Required:            true,
						Description:         "The browser-based SSO protocol to use. Options are `OIDC`, `SAML10`, `SAML11`, `SAML20`, `WSFED`.",
						MarkdownDescription: "The browser-based SSO protocol to use. Options are `OIDC`, `SAML10`, `SAML11`, `SAML20`, `WSFED`.",
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
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         "Determines whether SAML authentication requests should be signed.",
						MarkdownDescription: "Determines whether SAML authentication requests should be signed.",
					},
					"slo_service_endpoints": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:            true,
									Description:         "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.",
									MarkdownDescription: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.",
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"url": schema.StringAttribute{
									Required:            true,
									Description:         "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
									MarkdownDescription: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Optional:            true,
						Description:         "A list of possible endpoints to send SLO requests and responses.",
						MarkdownDescription: "A list of possible endpoints to send SLO requests and responses.",
					},
					"sso_application_endpoint": schema.StringAttribute{
						Optional: false,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"type": schema.StringAttribute{
													Required:            true,
													Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
													MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
											Optional:            true,
											Computed:            true,
											Default:             stringdefault.StaticString(""),
											Description:         "The value for this attribute.",
											MarkdownDescription: "The value for this attribute.",
										},
									},
								},
								Required:            true,
								Description:         "A list of mappings from attribute names to their fulfillment values.",
								MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
							},
							"attribute_sources": attributesources.ToSchema(0, false),
							"issuance_criteria": issuancecriteria.ToSchema(),
						},
						Optional:            true,
						Description:         "IdP Browser SSO OAuth Attribute Mapping",
						MarkdownDescription: "IdP Browser SSO OAuth Attribute Mapping",
					},
					"sso_service_endpoints": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:            true,
									Description:         "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.",
									MarkdownDescription: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.",
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Optional:            true,
						Description:         "The IdP SSO endpoints that define where to send your authentication requests. Only required for SP initiated SSO. This is required for SAML x.x and WS-FED Connections.",
						MarkdownDescription: "The IdP SSO endpoints that define where to send your authentication requests. Only required for SP initiated SSO. This is required for SAML x.x and WS-FED Connections.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"url_whitelist_entries": schema.SetNestedAttribute{
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"valid_path": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString(""),
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
					"access_token_manager_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_token_manager_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchema(),
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
														Validators: []validator.String{
															stringvalidator.LengthAtLeast(1),
														},
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
														MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
												Computed:            true,
												Optional:            true,
												Default:             stringdefault.StaticString(""),
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": attributesources.ToSchema(0, false),
								"issuance_criteria": issuancecriteria.ToSchema(),
							},
						},
						Required:            true,
						Description:         "A mapping in a connection that defines how access tokens are created.",
						MarkdownDescription: "A mapping in a connection that defines how access tokens are created.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"idp_oauth_attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            false,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Optional:            false,
											Computed:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            false,
								Computed:            true,
								Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
							"extended_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            true,
								Computed:            true,
								Default:             setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractExtendedAttributesAttrTypes}, nil)),
								Description:         "A list of additional attributes that are present in the incoming assertion.",
								MarkdownDescription: "A list of additional attributes that are present in the incoming assertion.",
							},
						},
						Required:            true,
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
						Description:         "Specify behavior of how SCIM DELETE requests are handled. Options are `DISABLE_USER`, `PERMANENTLY_DELETE_USER`.",
						MarkdownDescription: "Specify behavior of how SCIM DELETE requests are handled. Options are `DISABLE_USER`, `PERMANENTLY_DELETE_USER`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"DISABLE_USER",
								"PERMANENTLY_DELETE_USER",
							),
						},
					},
					"custom_schema": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attributes": schema.SetNestedAttribute{
								Description: "A custom SCIM attribute.",
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
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"sub_attributes": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											Description:         "List of sub-attributes for an attribute.",
											MarkdownDescription: "List of sub-attributes for an attribute.",
										},
										"types": schema.SetAttribute{
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
								Description: "Custom SCIM namespace.",
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
						Required:            true,
						Description:         "Custom SCIM Attributes configuration.",
						MarkdownDescription: "Custom SCIM Attributes configuration.",
					},
					"group_support": schema.BoolAttribute{
						Required:            true,
						Description:         "Specify support for provisioning of groups. Must be `true` to configure `groups` attribute.",
						MarkdownDescription: "Specify support for provisioning of groups. Must be `true` to configure `groups` attribute.",
					},
					"groups": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"read_groups": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"core_attributes": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            false,
															Computed:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Optional:            true,
															Computed:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
													},
												},
												Optional:            false,
												Computed:            true,
												Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
												MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
											},
											"extended_attributes": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Computed:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
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
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
															MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
													Optional:            true,
													Computed:            true,
													Default:             stringdefault.StaticString(""),
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
									"attributes": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
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
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
															MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
													Optional:            true,
													Computed:            true,
													Default:             stringdefault.StaticString(""),
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
						Optional:            true,
						Description:         "Group creation and read configuration. Requires `group_support` to be `true`.",
						MarkdownDescription: "Group creation and read configuration. Requires `group_support` to be `true`.",
					},
					"user_repository": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"identity_store": schema.SingleNestedAttribute{
								Optional:            true,
								Description:         "Identity Store Provisioner data store user repository.",
								MarkdownDescription: "Identity Store Provisioner data store user repository.",
								Attributes: map[string]schema.Attribute{
									"identity_store_provisioner_ref": schema.SingleNestedAttribute{
										Attributes:          resourcelink.ToSchema(),
										Required:            true,
										Description:         "Identity Store Provisioner data store user repository.",
										MarkdownDescription: "Identity Store Provisioner data store user repository.",
									},
								},
								Validators: []validator.Object{
									objectvalidator.ExactlyOneOf(
										path.MatchRelative().AtParent().AtName("identity_store"),
										path.MatchRelative().AtParent().AtName("ldap"),
									),
								},
							},
							"ldap": schema.SingleNestedAttribute{
								Description:         "LDAP Active Directory data store user repository.",
								MarkdownDescription: "LDAP Active Directory data store user repository.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"data_store_ref": schema.SingleNestedAttribute{
										Attributes:          resourcelink.ToSchema(),
										Required:            true,
										Description:         "Reference to the associated data store.",
										MarkdownDescription: "Reference to the associated data store.",
									},
									"base_dn": schema.StringAttribute{
										Optional:            true,
										Description:         "The base DN to search from. If not specified, the search will start at the LDAP's root.",
										MarkdownDescription: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"unique_user_id_filter": schema.StringAttribute{
										Required:            true,
										Description:         "The expression that results in a unique user identifier, when combined with the Base DN.",
										MarkdownDescription: "The expression that results in a unique user identifier, when combined with the Base DN.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"unique_group_id_filter": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										Default:             stringdefault.StaticString(""),
										Description:         "The expression that results in a unique group identifier, when combined with the Base DN. Only required when configuring the `inbound_provisioning.groups` attribute. Otherwise should not be set.",
										MarkdownDescription: "The expression that results in a unique group identifier, when combined with the Base DN. Only required when configuring the `inbound_provisioning.groups` attribute. Otherwise should not be set.",
									},
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
											"core_attributes": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            false,
															Computed:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
														},
														"name": schema.StringAttribute{
															Optional:            false,
															Computed:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            false,
												Computed:            true,
												Description:         "A list of read-only assertion attributes that are automatically populated by PingFederate.",
												MarkdownDescription: "A list of read-only assertion attributes that are automatically populated by PingFederate.",
											},
											"extended_attributes": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            true,
															Computed:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
													},
												},
												Optional: true,
												Computed: true,
												Default: setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
													"name":   types.StringType,
													"masked": types.BoolType,
												}}, nil)),
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
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
															MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
													Optional:            true,
													Computed:            true,
													Default:             stringdefault.StaticString(""),
													Description:         "The value for this attribute.",
													MarkdownDescription: "The value for this attribute.",
												},
											},
										},
										Required:            true,
										Description:         "A list of user repository mappings from attribute names to their fulfillment values.",
										MarkdownDescription: "A list of user repository mappings from attribute names to their fulfillment values.",
									},
									"attributes": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Required:            true,
													Description:         "The name of this attribute.",
													MarkdownDescription: "The name of this attribute.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
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
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
															},
														},
														"type": schema.StringAttribute{
															Required:            true,
															Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
															MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
													Optional:            true,
													Computed:            true,
													Default:             stringdefault.StaticString(""),
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"license_connection_group": schema.StringAttribute{
				Optional:            true,
				Description:         "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				MarkdownDescription: "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"logging_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The level of transaction logging applicable for this connection. Default is `STANDARD`. Options are `ENHANCED`, `FULL`, `NONE`, `STANDARD`. If the `idp_connection_transaction_logging_override` attribute is set to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource, then this attribute must be set to the same value.",
				MarkdownDescription: "The level of transaction logging applicable for this connection. Default is `STANDARD`. Options are `ENHANCED`, `FULL`, `NONE`, `STANDARD`. If the `idp_connection_transaction_logging_override` attribute is set to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource, then this attribute must be set to the same value.",
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
						Attributes:          resourcelink.ToSchema(),
					},
					"enable_auto_metadata_update": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies whether the metadata of the connection will be automatically reloaded. The default value is `true`.",
						MarkdownDescription: "Specifies whether the metadata of the connection will be automatically reloaded. The default value is `true`.",
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"client_secret": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						Description:         "The OpenID Connect client secret. Only one of `client_secret` or `encrypted_secret` can be set.",
						MarkdownDescription: "The OpenID Connect client secret. Only one of `client_secret` or `encrypted_secret` can be set.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"encrypted_secret": schema.StringAttribute{
						Description: "Encrypted OpenID Connect client secret. Only one of `client_secret` or `encrypted_secret` can be set.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("client_secret")),
						},
					},
				},
				Optional:            true,
				Description:         "The OpenID Connect Client Credentials settings. This is required for an OIDC Connection.",
				MarkdownDescription: "The OpenID Connect Client Credentials settings. This is required for an OIDC Connection.",
			},
			"virtual_entity_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "List of alternate entity IDs that identifies the local server to this partner.",
				MarkdownDescription: "List of alternate entity IDs that identifies the local server to this partner.",
			},
			"ws_trust": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
								Optional:            true,
								Description:         "A list of assertion attributes that are automatically populated by PingFederate.",
								MarkdownDescription: "A list of assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
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
						Description:         "Indicates whether a local token needs to be generated. The default value is `false`.",
						MarkdownDescription: "Indicates whether a local token needs to be generated. The default value is `false`.",
					},
					"token_generator_mappings": schema.SetNestedAttribute{
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
														Validators: []validator.String{
															configvalidators.PingFederateId(),
														},
													},
													"type": schema.StringAttribute{
														Required:            true,
														Description:         "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
														MarkdownDescription: "The source type of this key. Options are `ACTOR_TOKEN`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `ATTRIBUTE_QUERY`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `EXPRESSION`, `FRAGMENT`, `IDENTITY_STORE_GROUP`, `IDENTITY_STORE_USER`, `IDP_CONNECTION`, `INPUTS`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `LOCAL_IDENTITY_PROFILE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `OAUTH_PERSISTENT_GRANT`, `PASSWORD_CREDENTIAL_VALIDATOR`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `REQUEST`, `SCIM_GROUP`, `SCIM_USER`, `SUBJECT_TOKEN`, `TEXT`, `TOKEN`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `TRACKED_HTTP_PARAMS`.",
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
												Computed:            true,
												Optional:            true,
												Description:         "The value for this attribute.",
												MarkdownDescription: "The value for this attribute.",
												Default:             stringdefault.StaticString(""),
											},
										},
									},
									Required:            true,
									Description:         "A list of mappings from attribute names to their fulfillment values.",
									MarkdownDescription: "A list of mappings from attribute names to their fulfillment values.",
								},
								"attribute_sources": attributesources.ToSchema(0, false),
								"default_mapping": schema.BoolAttribute{
									Optional:            true,
									Description:         "Indicates whether the token generator mapping is the default mapping. The default value is `false`.",
									MarkdownDescription: "Indicates whether the token generator mapping is the default mapping. The default value is `false`.",
								},
								"issuance_criteria": issuancecriteria.ToSchema(),
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
								"sp_token_generator_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchema(),
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

	resp.Schema = schema
}

func (r *spIdpConnectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.0.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast1200 := compare >= 0
	// Compare to version 12.1.0 of PF
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast1210 := compare >= 0
	var plan *spIdpConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if plan == nil {
		return
	}

	// If any of these fields are set by the user and the PF version is not new enough, throw an error
	if !pfVersionAtLeast1200 {
		if internaltypes.IsDefined(plan.IdpBrowserSso) {
			oidcProviderSettings := plan.IdpBrowserSso.Attributes()["oidc_provider_settings"]
			if internaltypes.IsDefined(oidcProviderSettings) {
				frontChannelLogoutUri := oidcProviderSettings.(types.Object).Attributes()["front_channel_logout_uri"]
				if internaltypes.IsDefined(frontChannelLogoutUri) {
					version.AddUnsupportedAttributeError("idp_browser_sso.oidc_provider_settings.front_channel_logout_uri",
						r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
				}
				logoutEndpoint := oidcProviderSettings.(types.Object).Attributes()["logout_endpoint"]
				if internaltypes.IsDefined(logoutEndpoint) {
					version.AddUnsupportedAttributeError("idp_browser_sso.oidc_provider_settings.logout_endpoint",
						r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
				}
				postLogoutRedirectUri := oidcProviderSettings.(types.Object).Attributes()["post_logout_redirect_uri"]
				if internaltypes.IsDefined(postLogoutRedirectUri) {
					version.AddUnsupportedAttributeError("idp_browser_sso.oidc_provider_settings.post_logout_redirect_uri",
						r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
				}
			}
		}
	}
	if !pfVersionAtLeast1210 {
		if internaltypes.IsDefined(plan.IdpBrowserSso) {
			oidcProviderSettings := plan.IdpBrowserSso.Attributes()["oidc_provider_settings"]
			if internaltypes.IsDefined(oidcProviderSettings) {
				jwtSecuredAuthorizationResponseModeType := oidcProviderSettings.(types.Object).Attributes()["jwt_secured_authorization_response_mode_type"]
				if internaltypes.IsDefined(jwtSecuredAuthorizationResponseModeType) {
					version.AddUnsupportedAttributeError("idp_browser_sso.oidc_provider_settings.jwt_secured_authorization_response_mode_type",
						r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
				}
			}
		}
	}

	// Ensure that group attributes have appropriate values
	if internaltypes.IsDefined(plan.InboundProvisioning) {
		inboundProvisioningAttrs := plan.InboundProvisioning.Attributes()
		groupSupportEnabled := inboundProvisioningAttrs["group_support"].(types.Bool).ValueBool()
		groupsConfigured := internaltypes.IsDefined(inboundProvisioningAttrs["groups"])
		if !groupSupportEnabled && groupsConfigured {
			resp.Diagnostics.AddAttributeError(path.Root("inbound_provisioning"),
				providererror.InvalidAttributeConfiguration, "`inbound_provisioning.group_support` must be set to `true` to configure inbound provisioning groups.")
		}
		if groupSupportEnabled && !groupsConfigured {
			resp.Diagnostics.AddAttributeError(path.Root("inbound_provisioning"),
				providererror.InvalidAttributeConfiguration, "If `inbound_provisioning.group_support` is set to `true`, then `inbound_provisioning.groups` must be configured.")
		}

		userRepositoryLdap := inboundProvisioningAttrs["user_repository"].(types.Object).Attributes()["ldap"]
		if internaltypes.IsDefined(userRepositoryLdap) && len(userRepositoryLdap.(types.Object).Attributes()["unique_group_id_filter"].(types.String).ValueString()) == 0 && groupSupportEnabled {
			resp.Diagnostics.AddAttributeError(path.Root("inbound_provisioning.user_repository.ldap"),
				providererror.InvalidAttributeConfiguration, "`inbound_provisioning.user_repository.ldap.unique_group_id_filter` must be set if `inbound_provisioning.group_support` is set to `true`.")
		}
	}

	// Set default for jwt_secured_authorization_response_mode_type if version is 12.1+
	planModified := false
	var diags diag.Diagnostics
	if internaltypes.IsDefined(plan.IdpBrowserSso) {
		browserSsoAttributes := plan.IdpBrowserSso.Attributes()
		oidcProviderSettings := browserSsoAttributes["oidc_provider_settings"].(types.Object)
		if internaltypes.IsDefined(oidcProviderSettings) {
			oidcProviderSettingsAttributes := oidcProviderSettings.Attributes()
			jwtSecuredAuthorizationResponseModeType := oidcProviderSettingsAttributes["jwt_secured_authorization_response_mode_type"]
			if jwtSecuredAuthorizationResponseModeType.IsUnknown() {
				if pfVersionAtLeast1210 {
					oidcProviderSettingsAttributes["jwt_secured_authorization_response_mode_type"] = types.StringValue("DISABLED")
				} else {
					oidcProviderSettingsAttributes["jwt_secured_authorization_response_mode_type"] = types.StringNull()
				}
				oidcProviderSettings, diags = types.ObjectValue(oidcProviderSettings.AttributeTypes(ctx), oidcProviderSettingsAttributes)
				resp.Diagnostics.Append(diags...)
				browserSsoAttributes["oidc_provider_settings"] = oidcProviderSettings
				plan.IdpBrowserSso, diags = types.ObjectValue(plan.IdpBrowserSso.AttributeTypes(ctx), browserSsoAttributes)
				resp.Diagnostics.Append(diags...)
				planModified = true
			}
		}
	}

	// Set default for additional_allowed_entities_configuration for OIDC connections
	if plan.AdditionalAllowedEntitiesConfiguration.IsUnknown() {
		if internaltypes.IsDefined(plan.IdpBrowserSso) &&
			plan.IdpBrowserSso.Attributes()["protocol"].(types.String).ValueString() == "OIDC" {
			additionalAllowedEntities, diags := types.SetValue(entityIdAttrTypes, nil)
			resp.Diagnostics.Append(diags...)
			plan.AdditionalAllowedEntitiesConfiguration, diags = types.ObjectValue(additionalAllowedEntitiesConfigurationAttrTypes, map[string]attr.Value{
				"additional_allowed_entities": additionalAllowedEntities,
				"allow_additional_entities":   types.BoolValue(false),
				"allow_all_entities":          types.BoolValue(false),
			})
			resp.Diagnostics.Append(diags...)
		} else {
			plan.AdditionalAllowedEntitiesConfiguration = types.ObjectNull(additionalAllowedEntitiesConfigurationAttrTypes)
		}
		planModified = true
	}

	// Set default for decryption_policy for non-OIDC connections
	if internaltypes.IsDefined(plan.IdpBrowserSso) {
		browserSsoAttrs := plan.IdpBrowserSso.Attributes()
		if browserSsoAttrs["decryption_policy"].IsUnknown() {
			if browserSsoAttrs["protocol"].(types.String).ValueString() != "OIDC" {
				browserSsoAttrs["decryption_policy"], diags = types.ObjectValue(idpBrowserSsoDecryptionPolicyAttrTypes, map[string]attr.Value{
					"assertion_encrypted":           types.BoolValue(false),
					"attributes_encrypted":          types.BoolValue(false),
					"slo_encrypt_subject_name_id":   types.BoolValue(false),
					"slo_subject_name_id_encrypted": types.BoolValue(false),
					"subject_name_id_encrypted":     types.BoolValue(false),
				})
				resp.Diagnostics.Append(diags...)
			} else {
				browserSsoAttrs["decryption_policy"] = types.ObjectNull(idpBrowserSsoDecryptionPolicyAttrTypes)
			}
			plan.IdpBrowserSso, diags = types.ObjectValue(plan.IdpBrowserSso.AttributeTypes(ctx), browserSsoAttrs)
			resp.Diagnostics.Append(diags...)
			planModified = true
		}
	}

	// Similar logic for virtual_entity_ids for non-OIDC connections
	if plan.VirtualEntityIds.IsUnknown() {
		if internaltypes.IsDefined(plan.IdpBrowserSso) && plan.IdpBrowserSso.Attributes()["protocol"].(types.String).ValueString() == "OIDC" {
			plan.VirtualEntityIds = types.SetNull(types.StringType)
		} else {
			plan.VirtualEntityIds, diags = types.SetValue(types.StringType, nil)
			resp.Diagnostics.Append(diags...)
		}
		planModified = true
	}

	// Set default for error_page_msg_id
	if plan.ErrorPageMsgId.IsUnknown() {
		if internaltypes.IsDefined(plan.IdpBrowserSso) {
			plan.ErrorPageMsgId = types.StringValue("errorDetail.spSsoFailure")
		} else {
			plan.ErrorPageMsgId = types.StringNull()
		}
		planModified = true
	}

	if planModified {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}

	// If the entity_id has been changed, then mark corresponding attributes as unknown
	var state *spIdpConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}

	if (plan.EntityId.IsUnknown() || !plan.EntityId.Equal(state.EntityId)) && internaltypes.IsDefined(plan.IdpBrowserSso) {
		browserSsoAttributes := plan.IdpBrowserSso.Attributes()
		oidcProviderSettings := browserSsoAttributes["oidc_provider_settings"].(types.Object)
		if internaltypes.IsDefined(oidcProviderSettings) {
			oidcProviderSettingsAttributes := oidcProviderSettings.Attributes()
			oidcProviderSettingsAttributes["back_channel_logout_uri"] = types.StringUnknown()
			oidcProviderSettingsAttributes["front_channel_logout_uri"] = types.StringUnknown()
			oidcProviderSettingsAttributes["post_logout_redirect_uri"] = types.StringUnknown()
			oidcProviderSettingsAttributes["redirect_uri"] = types.StringUnknown()
			var diags diag.Diagnostics
			oidcProviderSettings, diags = types.ObjectValue(oidcProviderSettings.AttributeTypes(ctx), oidcProviderSettingsAttributes)
			resp.Diagnostics.Append(diags...)
			browserSsoAttributes["oidc_provider_settings"] = oidcProviderSettings
		}
		browserSsoAttributes["sso_application_endpoint"] = types.StringUnknown()
		plan.IdpBrowserSso, diags = types.ObjectValue(plan.IdpBrowserSso.AttributeTypes(ctx), browserSsoAttributes)
		resp.Diagnostics.Append(diags...)
		planModified = true
	}

	// If the attribute_contract changes, mark the jit_provisioning.user_attributes.attribute_contract as unknown
	var planAttributeContract, stateAttributeContract attr.Value
	if internaltypes.IsDefined(plan.IdpBrowserSso) {
		planBrowserSsoAttributes := plan.IdpBrowserSso.Attributes()
		planAttributeContract = planBrowserSsoAttributes["attribute_contract"]
		if internaltypes.IsDefined(state.IdpBrowserSso) {
			stateAttributeContract = state.IdpBrowserSso.Attributes()["attribute_contract"]
		}

		if (planAttributeContract != nil && !planAttributeContract.Equal(stateAttributeContract)) || (stateAttributeContract != nil && !stateAttributeContract.Equal(planAttributeContract)) {
			planJitProvisioning := planBrowserSsoAttributes["jit_provisioning"]
			if internaltypes.IsDefined(planJitProvisioning) {
				userAttrs := planJitProvisioning.(types.Object).Attributes()["user_attributes"]
				if internaltypes.IsDefined(userAttrs) {
					userAttrsAttrs := userAttrs.(types.Object).Attributes()
					userAttrsAttrs["attribute_contract"] = types.SetUnknown(idpBrowserSsoJitProvisioningUserAttributesAttributeContractElementType)
					userAttrsUpdated, diags := types.ObjectValue(userAttrs.(types.Object).AttributeTypes(ctx), userAttrsAttrs)
					resp.Diagnostics.Append(diags...)

					jitProvisioningAttrs := planJitProvisioning.(types.Object).Attributes()
					jitProvisioningAttrs["user_attributes"] = userAttrsUpdated
					jitProvisioningUpdated, diags := types.ObjectValue(planJitProvisioning.(types.Object).AttributeTypes(ctx), jitProvisioningAttrs)
					resp.Diagnostics.Append(diags...)

					planBrowserSsoAttributes["jit_provisioning"] = jitProvisioningUpdated
					plan.IdpBrowserSso, diags = types.ObjectValue(plan.IdpBrowserSso.AttributeTypes(ctx), planBrowserSsoAttributes)
					resp.Diagnostics.Append(diags...)

					planModified = true
				}
			}
		}
	}

	// Handle the encrypted OIDC client secret
	if internaltypes.IsDefined(plan.OidcClientCredentials) {
		planOidcClientCredentials := plan.OidcClientCredentials.Attributes()
		stateOidcClientCredentials := state.OidcClientCredentials.Attributes()
		if !internaltypes.IsDefined(planOidcClientCredentials["client_secret"]) && planOidcClientCredentials["encrypted_secret"].IsUnknown() {
			planOidcClientCredentials["encrypted_secret"] = types.StringNull()
			plan.OidcClientCredentials, diags = types.ObjectValue(plan.OidcClientCredentials.AttributeTypes(ctx), planOidcClientCredentials)
			resp.Diagnostics.Append(diags...)
			planModified = true
		} else if !planOidcClientCredentials["client_secret"].Equal(stateOidcClientCredentials["client_secret"]) {
			planOidcClientCredentials["encrypted_secret"] = types.StringUnknown()
			plan.OidcClientCredentials, diags = types.ObjectValue(plan.OidcClientCredentials.AttributeTypes(ctx), planOidcClientCredentials)
			resp.Diagnostics.Append(diags...)
			planModified = true
		}
	}

	// Handle the computed attributes within oidc_provider_settings
	if internaltypes.IsDefined(plan.IdpBrowserSso) && internaltypes.IsDefined(state.IdpBrowserSso) {
		planSsoAttrs := plan.IdpBrowserSso.Attributes()
		stateSsoAttrs := state.IdpBrowserSso.Attributes()
		if internaltypes.IsDefined(planSsoAttrs["oidc_provider_settings"]) && internaltypes.IsDefined(stateSsoAttrs["oidc_provider_settings"]) {
			planOidcProviderSettings := planSsoAttrs["oidc_provider_settings"].(types.Object)
			stateOidcProviderSettings := stateSsoAttrs["oidc_provider_settings"].(types.Object)
			planOidcProviderSettingsAttrs := planOidcProviderSettings.Attributes()
			stateOidcProviderSettingsAttrs := stateOidcProviderSettings.Attributes()

			// Check if all the non-computed attributes are unchanged
			if planOidcProviderSettingsAttrs["authentication_scheme"].Equal(stateOidcProviderSettingsAttrs["authentication_scheme"]) &&
				planOidcProviderSettingsAttrs["authentication_signing_algorithm"].Equal(stateOidcProviderSettingsAttrs["authentication_signing_algorithm"]) &&
				planOidcProviderSettingsAttrs["authorization_endpoint"].Equal(stateOidcProviderSettingsAttrs["authorization_endpoint"]) &&
				planOidcProviderSettingsAttrs["enable_pkce"].Equal(stateOidcProviderSettingsAttrs["enable_pkce"]) &&
				planOidcProviderSettingsAttrs["jwks_url"].Equal(stateOidcProviderSettingsAttrs["jwks_url"]) &&
				planOidcProviderSettingsAttrs["jwt_secured_authorization_response_mode_type"].Equal(stateOidcProviderSettingsAttrs["jwt_secured_authorization_response_mode_type"]) &&
				planOidcProviderSettingsAttrs["login_type"].Equal(stateOidcProviderSettingsAttrs["login_type"]) &&
				planOidcProviderSettingsAttrs["logout_endpoint"].Equal(stateOidcProviderSettingsAttrs["logout_endpoint"]) &&
				planOidcProviderSettingsAttrs["pushed_authorization_request_endpoint"].Equal(stateOidcProviderSettingsAttrs["pushed_authorization_request_endpoint"]) &&
				planOidcProviderSettingsAttrs["request_parameters"].Equal(stateOidcProviderSettingsAttrs["request_parameters"]) &&
				planOidcProviderSettingsAttrs["request_signing_algorithm"].Equal(stateOidcProviderSettingsAttrs["request_signing_algorithm"]) &&
				planOidcProviderSettingsAttrs["scopes"].Equal(stateOidcProviderSettingsAttrs["scopes"]) &&
				planOidcProviderSettingsAttrs["token_endpoint"].Equal(stateOidcProviderSettingsAttrs["token_endpoint"]) &&
				planOidcProviderSettingsAttrs["track_user_sessions_for_logout"].Equal(stateOidcProviderSettingsAttrs["track_user_sessions_for_logout"]) &&
				planOidcProviderSettingsAttrs["user_info_endpoint"].Equal(stateOidcProviderSettingsAttrs["user_info_endpoint"]) {
				// Keep the state values for computed attributes
				planOidcProviderSettingsAttrs["back_channel_logout_uri"] = types.StringPointerValue(stateOidcProviderSettingsAttrs["back_channel_logout_uri"].(types.String).ValueStringPointer())
				planOidcProviderSettingsAttrs["front_channel_logout_uri"] = types.StringPointerValue(stateOidcProviderSettingsAttrs["front_channel_logout_uri"].(types.String).ValueStringPointer())
				planOidcProviderSettingsAttrs["redirect_uri"] = types.StringPointerValue(stateOidcProviderSettingsAttrs["redirect_uri"].(types.String).ValueStringPointer())
				planOidcProviderSettingsAttrs["post_logout_redirect_uri"] = types.StringPointerValue(stateOidcProviderSettingsAttrs["post_logout_redirect_uri"].(types.String).ValueStringPointer())

				planOidcProviderSettings, diags := types.ObjectValue(planOidcProviderSettings.AttributeTypes(ctx), planOidcProviderSettingsAttrs)
				resp.Diagnostics.Append(diags...)
				planSsoAttrs["oidc_provider_settings"] = planOidcProviderSettings
				plan.IdpBrowserSso, diags = types.ObjectValue(plan.IdpBrowserSso.AttributeTypes(ctx), planSsoAttrs)
				resp.Diagnostics.Append(diags...)
				planModified = true
			}
		}
	}

	if planModified {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

func addOptionalSpIdpConnectionFields(ctx context.Context, addRequest *client.IdpConnection, plan spIdpConnectionResourceModel) diag.Diagnostics {
	var err error
	var respDiags diag.Diagnostics
	addRequest.ErrorPageMsgId = plan.ErrorPageMsgId.ValueStringPointer()
	if !plan.ConnectionId.IsUnknown() {
		addRequest.Id = plan.ConnectionId.ValueStringPointer()
	}
	addRequest.Type = utils.Pointer("IDP")
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.BaseUrl = plan.BaseUrl.ValueStringPointer()
	addRequest.DefaultVirtualEntityId = plan.DefaultVirtualEntityId.ValueStringPointer()

	if internaltypes.IsDefined(plan.LicenseConnectionGroup) {
		addRequest.LicenseConnectionGroup = plan.LicenseConnectionGroup.ValueStringPointer()
	}

	addRequest.LoggingMode = plan.LoggingMode.ValueStringPointer()

	if internaltypes.IsDefined(plan.VirtualEntityIds) {
		var virtualIdentitySlice []string
		plan.VirtualEntityIds.ElementsAs(ctx, &virtualIdentitySlice, false)
		addRequest.VirtualEntityIds = virtualIdentitySlice
	}

	// oidc_client_credentials
	if !plan.OidcClientCredentials.IsNull() {
		oidcClientCredentialsValue := &client.OIDCClientCredentials{}
		oidcClientCredentialsAttrs := plan.OidcClientCredentials.Attributes()
		oidcClientCredentialsValue.ClientId = oidcClientCredentialsAttrs["client_id"].(types.String).ValueString()
		oidcClientCredentialsValue.ClientSecret = oidcClientCredentialsAttrs["client_secret"].(types.String).ValueStringPointer()
		addRequest.OidcClientCredentials = oidcClientCredentialsValue
	}

	// metadata_reload_settings
	if !plan.MetadataReloadSettings.IsNull() {
		metadataReloadSettingsValue := &client.ConnectionMetadataUrl{}
		metadataReloadSettingsAttrs := plan.MetadataReloadSettings.Attributes()
		metadataReloadSettingsValue.EnableAutoMetadataUpdate = metadataReloadSettingsAttrs["enable_auto_metadata_update"].(types.Bool).ValueBoolPointer()
		metadataReloadSettingsMetadataUrlRefValue := client.ResourceLink{}
		metadataReloadSettingsMetadataUrlRefAttrs := metadataReloadSettingsAttrs["metadata_url_ref"].(types.Object).Attributes()
		metadataReloadSettingsMetadataUrlRefValue.Id = metadataReloadSettingsMetadataUrlRefAttrs["id"].(types.String).ValueString()
		metadataReloadSettingsValue.MetadataUrlRef = metadataReloadSettingsMetadataUrlRefValue
		addRequest.MetadataReloadSettings = metadataReloadSettingsValue
	}

	if internaltypes.IsDefined(plan.Credentials) {
		addRequest.Credentials = &client.ConnectionCredentials{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Credentials, true)), addRequest.Credentials)
		if err != nil {
			respDiags.AddError("Error building client struct for credentials", err.Error())
		}
		if addRequest.Credentials.InboundBackChannelAuth != nil {
			addRequest.Credentials.InboundBackChannelAuth.Type = "INBOUND"
		}
		if addRequest.Credentials.OutboundBackChannelAuth != nil {
			addRequest.Credentials.OutboundBackChannelAuth.Type = "OUTBOUND"
		}
	}

	// contact_info
	if !plan.ContactInfo.IsNull() {
		contactInfoValue := &client.ContactInfo{}
		contactInfoAttrs := plan.ContactInfo.Attributes()
		contactInfoValue.Company = contactInfoAttrs["company"].(types.String).ValueStringPointer()
		contactInfoValue.Email = contactInfoAttrs["email"].(types.String).ValueStringPointer()
		contactInfoValue.FirstName = contactInfoAttrs["first_name"].(types.String).ValueStringPointer()
		contactInfoValue.LastName = contactInfoAttrs["last_name"].(types.String).ValueStringPointer()
		contactInfoValue.Phone = contactInfoAttrs["phone"].(types.String).ValueStringPointer()
		addRequest.ContactInfo = contactInfoValue
	}

	// additional_allowed_entities_configuration
	if internaltypes.IsDefined(plan.AdditionalAllowedEntitiesConfiguration) {
		additionalAllowedEntitiesConfigurationValue := &client.AdditionalAllowedEntitiesConfiguration{}
		additionalAllowedEntitiesConfigurationAttrs := plan.AdditionalAllowedEntitiesConfiguration.Attributes()
		additionalAllowedEntitiesConfigurationValue.AdditionalAllowedEntities = []client.Entity{}
		for _, additionalAllowedEntitiesElement := range additionalAllowedEntitiesConfigurationAttrs["additional_allowed_entities"].(types.Set).Elements() {
			additionalAllowedEntitiesValue := client.Entity{}
			additionalAllowedEntitiesAttrs := additionalAllowedEntitiesElement.(types.Object).Attributes()
			additionalAllowedEntitiesValue.EntityDescription = additionalAllowedEntitiesAttrs["entity_description"].(types.String).ValueStringPointer()
			additionalAllowedEntitiesValue.EntityId = additionalAllowedEntitiesAttrs["entity_id"].(types.String).ValueStringPointer()
			additionalAllowedEntitiesConfigurationValue.AdditionalAllowedEntities = append(additionalAllowedEntitiesConfigurationValue.AdditionalAllowedEntities, additionalAllowedEntitiesValue)
		}
		additionalAllowedEntitiesConfigurationValue.AllowAdditionalEntities = additionalAllowedEntitiesConfigurationAttrs["allow_additional_entities"].(types.Bool).ValueBoolPointer()
		additionalAllowedEntitiesConfigurationValue.AllowAllEntities = additionalAllowedEntitiesConfigurationAttrs["allow_all_entities"].(types.Bool).ValueBoolPointer()
		addRequest.AdditionalAllowedEntitiesConfiguration = additionalAllowedEntitiesConfigurationValue
	}

	// extended_properties
	if !plan.ExtendedProperties.IsNull() {
		addRequest.ExtendedProperties = &map[string]client.ParameterValues{}
		for key, extendedPropertiesElement := range plan.ExtendedProperties.Elements() {
			extendedPropertiesValue := client.ParameterValues{}
			extendedPropertiesAttrs := extendedPropertiesElement.(types.Object).Attributes()
			if !extendedPropertiesAttrs["values"].IsNull() {
				extendedPropertiesValue.Values = []string{}
				for _, valuesElement := range extendedPropertiesAttrs["values"].(types.Set).Elements() {
					extendedPropertiesValue.Values = append(extendedPropertiesValue.Values, valuesElement.(types.String).ValueString())
				}
			}
			(*addRequest.ExtendedProperties)[key] = extendedPropertiesValue
		}
	}

	// idp_browser_sso
	if !plan.IdpBrowserSso.IsNull() {
		idpBrowserSsoValue := &client.IdpBrowserSso{}
		idpBrowserSsoAttrs := plan.IdpBrowserSso.Attributes()
		idpBrowserSsoValue.AdapterMappings = []client.SpAdapterMapping{}
		for _, adapterMappingsElement := range idpBrowserSsoAttrs["adapter_mappings"].(types.List).Elements() {
			adapterMappingsValue := client.SpAdapterMapping{}
			adapterMappingsAttrs := adapterMappingsElement.(types.Object).Attributes()
			if !adapterMappingsAttrs["adapter_override_settings"].IsNull() {
				adapterMappingsAdapterOverrideSettingsValue := &client.SpAdapter{}
				adapterMappingsAdapterOverrideSettingsAttrs := adapterMappingsAttrs["adapter_override_settings"].(types.Object).Attributes()
				if !adapterMappingsAdapterOverrideSettingsAttrs["attribute_contract"].IsNull() {
					adapterMappingsAdapterOverrideSettingsAttributeContractValue := &client.SpAdapterAttributeContract{}
					adapterMappingsAdapterOverrideSettingsAttributeContractAttrs := adapterMappingsAdapterOverrideSettingsAttrs["attribute_contract"].(types.Object).Attributes()
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes = []client.SpAdapterAttribute{}
					for _, coreAttributesElement := range adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
						coreAttributesValue := client.SpAdapterAttribute{}
						coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
						coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
						adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes = append(adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes, coreAttributesValue)
					}
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes = []client.SpAdapterAttribute{}
					for _, extendedAttributesElement := range adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
						extendedAttributesValue := client.SpAdapterAttribute{}
						extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
						extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
						adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes = append(adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
					}
					adapterMappingsAdapterOverrideSettingsValue.AttributeContract = adapterMappingsAdapterOverrideSettingsAttributeContractValue
				}
				adapterMappingsAdapterOverrideSettingsConfigurationValue, err := pluginconfiguration.ClientStruct(adapterMappingsAdapterOverrideSettingsAttrs["configuration"].(types.Object))
				if err != nil {
					respDiags.AddError("Error building client struct for configuration", err.Error())
				} else {
					adapterMappingsAdapterOverrideSettingsValue.Configuration = *adapterMappingsAdapterOverrideSettingsConfigurationValue
				}
				adapterMappingsAdapterOverrideSettingsValue.Id = adapterMappingsAdapterOverrideSettingsAttrs["id"].(types.String).ValueString()
				adapterMappingsAdapterOverrideSettingsValue.Name = adapterMappingsAdapterOverrideSettingsAttrs["name"].(types.String).ValueString()
				if !adapterMappingsAdapterOverrideSettingsAttrs["parent_ref"].IsNull() {
					adapterMappingsAdapterOverrideSettingsParentRefValue := &client.ResourceLink{}
					adapterMappingsAdapterOverrideSettingsParentRefAttrs := adapterMappingsAdapterOverrideSettingsAttrs["parent_ref"].(types.Object).Attributes()
					adapterMappingsAdapterOverrideSettingsParentRefValue.Id = adapterMappingsAdapterOverrideSettingsParentRefAttrs["id"].(types.String).ValueString()
					adapterMappingsAdapterOverrideSettingsValue.ParentRef = adapterMappingsAdapterOverrideSettingsParentRefValue
				}
				adapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue := client.ResourceLink{}
				adapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrs := adapterMappingsAdapterOverrideSettingsAttrs["plugin_descriptor_ref"].(types.Object).Attributes()
				adapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue.Id = adapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrs["id"].(types.String).ValueString()
				adapterMappingsAdapterOverrideSettingsValue.PluginDescriptorRef = adapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue
				if !adapterMappingsAdapterOverrideSettingsAttrs["target_application_info"].IsNull() {
					adapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue := &client.SpAdapterTargetApplicationInfo{}
					adapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrs := adapterMappingsAdapterOverrideSettingsAttrs["target_application_info"].(types.Object).Attributes()
					adapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue.ApplicationIconUrl = adapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrs["application_icon_url"].(types.String).ValueStringPointer()
					adapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue.ApplicationName = adapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrs["application_name"].(types.String).ValueStringPointer()
					adapterMappingsAdapterOverrideSettingsValue.TargetApplicationInfo = adapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue
				}
				adapterMappingsValue.AdapterOverrideSettings = adapterMappingsAdapterOverrideSettingsValue
			}
			adapterMappingsValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(adapterMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
			}
			adapterMappingsValue.AttributeSources, err = attributesources.ClientStruct(adapterMappingsAttrs["attribute_sources"].(types.Set))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_sources", err.Error())
			}
			adapterMappingsValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(adapterMappingsAttrs["issuance_criteria"].(types.Object))
			if err != nil {
				respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
			}
			adapterMappingsValue.RestrictVirtualEntityIds = adapterMappingsAttrs["restrict_virtual_entity_ids"].(types.Bool).ValueBoolPointer()
			if !adapterMappingsAttrs["restricted_virtual_entity_ids"].IsNull() {
				adapterMappingsValue.RestrictedVirtualEntityIds = []string{}
				for _, restrictedVirtualEntityIdsElement := range adapterMappingsAttrs["restricted_virtual_entity_ids"].(types.Set).Elements() {
					adapterMappingsValue.RestrictedVirtualEntityIds = append(adapterMappingsValue.RestrictedVirtualEntityIds, restrictedVirtualEntityIdsElement.(types.String).ValueString())
				}
			}
			if !adapterMappingsAttrs["sp_adapter_ref"].IsNull() {
				adapterMappingsSpAdapterRefValue := &client.ResourceLink{}
				adapterMappingsSpAdapterRefAttrs := adapterMappingsAttrs["sp_adapter_ref"].(types.Object).Attributes()
				adapterMappingsSpAdapterRefValue.Id = adapterMappingsSpAdapterRefAttrs["id"].(types.String).ValueString()
				adapterMappingsValue.SpAdapterRef = adapterMappingsSpAdapterRefValue
			}
			idpBrowserSsoValue.AdapterMappings = append(idpBrowserSsoValue.AdapterMappings, adapterMappingsValue)
		}
		idpBrowserSsoValue.AlwaysSignArtifactResponse = idpBrowserSsoAttrs["always_sign_artifact_response"].(types.Bool).ValueBoolPointer()
		if !idpBrowserSsoAttrs["artifact"].IsNull() {
			idpBrowserSsoArtifactValue := &client.ArtifactSettings{}
			idpBrowserSsoArtifactAttrs := idpBrowserSsoAttrs["artifact"].(types.Object).Attributes()
			idpBrowserSsoArtifactValue.Lifetime = idpBrowserSsoArtifactAttrs["lifetime"].(types.Int64).ValueInt64Pointer()
			idpBrowserSsoArtifactValue.ResolverLocations = []client.ArtifactResolverLocation{}
			for _, resolverLocationsElement := range idpBrowserSsoArtifactAttrs["resolver_locations"].(types.Set).Elements() {
				resolverLocationsValue := client.ArtifactResolverLocation{}
				resolverLocationsAttrs := resolverLocationsElement.(types.Object).Attributes()
				resolverLocationsValue.Index = resolverLocationsAttrs["index"].(types.Int64).ValueInt64()
				resolverLocationsValue.Url = resolverLocationsAttrs["url"].(types.String).ValueString()
				idpBrowserSsoArtifactValue.ResolverLocations = append(idpBrowserSsoArtifactValue.ResolverLocations, resolverLocationsValue)
			}
			idpBrowserSsoArtifactValue.SourceId = idpBrowserSsoArtifactAttrs["source_id"].(types.String).ValueStringPointer()
			idpBrowserSsoValue.Artifact = idpBrowserSsoArtifactValue
		}
		idpBrowserSsoValue.AssertionsSigned = idpBrowserSsoAttrs["assertions_signed"].(types.Bool).ValueBoolPointer()
		if internaltypes.IsDefined(idpBrowserSsoAttrs["attribute_contract"]) {
			idpBrowserSsoAttributeContractValue := &client.IdpBrowserSsoAttributeContract{}
			idpBrowserSsoAttributeContractAttrs := idpBrowserSsoAttrs["attribute_contract"].(types.Object).Attributes()
			idpBrowserSsoAttributeContractValue.CoreAttributes = []client.IdpBrowserSsoAttribute{}
			for _, coreAttributesElement := range idpBrowserSsoAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
				coreAttributesValue := client.IdpBrowserSsoAttribute{}
				coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
				coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
				idpBrowserSsoAttributeContractValue.CoreAttributes = append(idpBrowserSsoAttributeContractValue.CoreAttributes, coreAttributesValue)
			}
			idpBrowserSsoAttributeContractValue.ExtendedAttributes = []client.IdpBrowserSsoAttribute{}
			for _, extendedAttributesElement := range idpBrowserSsoAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.IdpBrowserSsoAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				idpBrowserSsoAttributeContractValue.ExtendedAttributes = append(idpBrowserSsoAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
			idpBrowserSsoValue.AttributeContract = idpBrowserSsoAttributeContractValue
		}
		idpBrowserSsoValue.AuthenticationPolicyContractMappings = []client.AuthenticationPolicyContractMapping{}
		for _, authenticationPolicyContractMappingsElement := range idpBrowserSsoAttrs["authentication_policy_contract_mappings"].(types.List).Elements() {
			authenticationPolicyContractMappingsValue := client.AuthenticationPolicyContractMapping{}
			authenticationPolicyContractMappingsAttrs := authenticationPolicyContractMappingsElement.(types.Object).Attributes()
			authenticationPolicyContractMappingsValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(authenticationPolicyContractMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
			}
			authenticationPolicyContractMappingsValue.AttributeSources, err = attributesources.ClientStruct(authenticationPolicyContractMappingsAttrs["attribute_sources"].(types.Set))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_sources", err.Error())
			}
			authenticationPolicyContractMappingsAuthenticationPolicyContractRefValue := client.ResourceLink{}
			authenticationPolicyContractMappingsAuthenticationPolicyContractRefAttrs := authenticationPolicyContractMappingsAttrs["authentication_policy_contract_ref"].(types.Object).Attributes()
			authenticationPolicyContractMappingsAuthenticationPolicyContractRefValue.Id = authenticationPolicyContractMappingsAuthenticationPolicyContractRefAttrs["id"].(types.String).ValueString()
			authenticationPolicyContractMappingsValue.AuthenticationPolicyContractRef = authenticationPolicyContractMappingsAuthenticationPolicyContractRefValue
			authenticationPolicyContractMappingsValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(authenticationPolicyContractMappingsAttrs["issuance_criteria"].(types.Object))
			if err != nil {
				respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
			}
			authenticationPolicyContractMappingsValue.RestrictVirtualServerIds = authenticationPolicyContractMappingsAttrs["restrict_virtual_server_ids"].(types.Bool).ValueBoolPointer()
			if !authenticationPolicyContractMappingsAttrs["restricted_virtual_server_ids"].IsNull() {
				authenticationPolicyContractMappingsValue.RestrictedVirtualServerIds = []string{}
				for _, restrictedVirtualServerIdsElement := range authenticationPolicyContractMappingsAttrs["restricted_virtual_server_ids"].(types.Set).Elements() {
					authenticationPolicyContractMappingsValue.RestrictedVirtualServerIds = append(authenticationPolicyContractMappingsValue.RestrictedVirtualServerIds, restrictedVirtualServerIdsElement.(types.String).ValueString())
				}
			}
			idpBrowserSsoValue.AuthenticationPolicyContractMappings = append(idpBrowserSsoValue.AuthenticationPolicyContractMappings, authenticationPolicyContractMappingsValue)
		}
		for _, authnContextMappingsElement := range idpBrowserSsoAttrs["authn_context_mappings"].(types.Set).Elements() {
			authnContextMappingsValue := client.AuthnContextMapping{}
			authnContextMappingsAttrs := authnContextMappingsElement.(types.Object).Attributes()
			authnContextMappingsValue.Local = authnContextMappingsAttrs["local"].(types.String).ValueStringPointer()
			authnContextMappingsValue.Remote = authnContextMappingsAttrs["remote"].(types.String).ValueStringPointer()
			idpBrowserSsoValue.AuthnContextMappings = append(idpBrowserSsoValue.AuthnContextMappings, authnContextMappingsValue)
		}
		if !idpBrowserSsoAttrs["decryption_policy"].IsNull() {
			idpBrowserSsoDecryptionPolicyValue := &client.DecryptionPolicy{}
			idpBrowserSsoDecryptionPolicyAttrs := idpBrowserSsoAttrs["decryption_policy"].(types.Object).Attributes()
			idpBrowserSsoDecryptionPolicyValue.AssertionEncrypted = idpBrowserSsoDecryptionPolicyAttrs["assertion_encrypted"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoDecryptionPolicyValue.AttributesEncrypted = idpBrowserSsoDecryptionPolicyAttrs["attributes_encrypted"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoDecryptionPolicyValue.SloEncryptSubjectNameID = idpBrowserSsoDecryptionPolicyAttrs["slo_encrypt_subject_name_id"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoDecryptionPolicyValue.SloSubjectNameIDEncrypted = idpBrowserSsoDecryptionPolicyAttrs["slo_subject_name_id_encrypted"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoDecryptionPolicyValue.SubjectNameIdEncrypted = idpBrowserSsoDecryptionPolicyAttrs["subject_name_id_encrypted"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoValue.DecryptionPolicy = idpBrowserSsoDecryptionPolicyValue
		}
		idpBrowserSsoValue.DefaultTargetUrl = idpBrowserSsoAttrs["default_target_url"].(types.String).ValueStringPointer()
		if !idpBrowserSsoAttrs["enabled_profiles"].IsNull() {
			idpBrowserSsoValue.EnabledProfiles = []string{}
			for _, enabledProfilesElement := range idpBrowserSsoAttrs["enabled_profiles"].(types.Set).Elements() {
				idpBrowserSsoValue.EnabledProfiles = append(idpBrowserSsoValue.EnabledProfiles, enabledProfilesElement.(types.String).ValueString())
			}
		}
		idpBrowserSsoValue.IdpIdentityMapping = idpBrowserSsoAttrs["idp_identity_mapping"].(types.String).ValueString()
		if !idpBrowserSsoAttrs["incoming_bindings"].IsNull() {
			idpBrowserSsoValue.IncomingBindings = []string{}
			for _, incomingBindingsElement := range idpBrowserSsoAttrs["incoming_bindings"].(types.Set).Elements() {
				idpBrowserSsoValue.IncomingBindings = append(idpBrowserSsoValue.IncomingBindings, incomingBindingsElement.(types.String).ValueString())
			}
		}
		if !idpBrowserSsoAttrs["jit_provisioning"].IsNull() {
			idpBrowserSsoJitProvisioningValue := &client.JitProvisioning{}
			idpBrowserSsoJitProvisioningAttrs := idpBrowserSsoAttrs["jit_provisioning"].(types.Object).Attributes()
			idpBrowserSsoJitProvisioningValue.ErrorHandling = idpBrowserSsoJitProvisioningAttrs["error_handling"].(types.String).ValueStringPointer()
			idpBrowserSsoJitProvisioningValue.EventTrigger = idpBrowserSsoJitProvisioningAttrs["event_trigger"].(types.String).ValueStringPointer()
			idpBrowserSsoJitProvisioningUserAttributesValue := client.JitProvisioningUserAttributes{}
			idpBrowserSsoJitProvisioningUserAttributesAttrs := idpBrowserSsoJitProvisioningAttrs["user_attributes"].(types.Object).Attributes()
			idpBrowserSsoJitProvisioningUserAttributesValue.AttributeContract = []client.IdpBrowserSsoAttribute{}
			for _, attributeContractElement := range idpBrowserSsoJitProvisioningUserAttributesAttrs["attribute_contract"].(types.Set).Elements() {
				attributeContractValue := client.IdpBrowserSsoAttribute{}
				attributeContractAttrs := attributeContractElement.(types.Object).Attributes()
				attributeContractValue.Masked = attributeContractAttrs["masked"].(types.Bool).ValueBoolPointer()
				attributeContractValue.Name = attributeContractAttrs["name"].(types.String).ValueString()
				idpBrowserSsoJitProvisioningUserAttributesValue.AttributeContract = append(idpBrowserSsoJitProvisioningUserAttributesValue.AttributeContract, attributeContractValue)
			}
			idpBrowserSsoJitProvisioningUserAttributesValue.DoAttributeQuery = idpBrowserSsoJitProvisioningUserAttributesAttrs["do_attribute_query"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoJitProvisioningValue.UserAttributes = idpBrowserSsoJitProvisioningUserAttributesValue
			idpBrowserSsoJitProvisioningUserRepositoryValue := client.DataStoreRepositoryAggregation{}
			idpBrowserSsoJitProvisioningUserRepositoryAttrs := idpBrowserSsoJitProvisioningAttrs["user_repository"].(types.Object).Attributes()
			if !idpBrowserSsoJitProvisioningUserRepositoryAttrs["jdbc"].IsNull() {
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue := &client.JdbcDataStoreRepository{}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryAttrs := idpBrowserSsoJitProvisioningUserRepositoryAttrs["jdbc"].(types.Object).Attributes()
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefValue := client.ResourceLink{}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefAttrs := idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryAttrs["data_store_ref"].(types.Object).Attributes()
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefValue.Id = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefAttrs["id"].(types.String).ValueString()
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue.DataStoreRef = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefValue
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue.JitRepositoryAttributeMapping = map[string]client.AttributeFulfillmentValue{}
				for key, jitRepositoryAttributeMappingElement := range idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryAttrs["jit_repository_attribute_mapping"].(types.Map).Elements() {
					jitRepositoryAttributeMappingValue := client.AttributeFulfillmentValue{}
					jitRepositoryAttributeMappingAttrs := jitRepositoryAttributeMappingElement.(types.Object).Attributes()
					jitRepositoryAttributeMappingSourceValue := client.SourceTypeIdKey{}
					jitRepositoryAttributeMappingSourceAttrs := jitRepositoryAttributeMappingAttrs["source"].(types.Object).Attributes()
					jitRepositoryAttributeMappingSourceValue.Id = jitRepositoryAttributeMappingSourceAttrs["id"].(types.String).ValueStringPointer()
					jitRepositoryAttributeMappingSourceValue.Type = jitRepositoryAttributeMappingSourceAttrs["type"].(types.String).ValueString()
					jitRepositoryAttributeMappingValue.Source = jitRepositoryAttributeMappingSourceValue
					jitRepositoryAttributeMappingValue.Value = jitRepositoryAttributeMappingAttrs["value"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue.JitRepositoryAttributeMapping[key] = jitRepositoryAttributeMappingValue
				}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue := client.SqlMethod{}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodAttrs := idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryAttrs["sql_method"].(types.Object).Attributes()
				if !idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodAttrs["stored_procedure"].IsNull() {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue := &client.StoredProcedure{}
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureAttrs := idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodAttrs["stored_procedure"].(types.Object).Attributes()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue.Schema = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureAttrs["schema"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue.StoredProcedure = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureAttrs["stored_procedure"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue.StoredProcedure = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue
				}
				if !idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodAttrs["table"].IsNull() {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue := &client.Table{}
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableAttrs := idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodAttrs["table"].(types.Object).Attributes()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue.Schema = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableAttrs["schema"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue.TableName = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableAttrs["table_name"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue.UniqueIdColumn = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableAttrs["unique_id_column"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue.Table = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue
				}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue.SqlMethod = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue.Type = "JDBC"
				idpBrowserSsoJitProvisioningUserRepositoryValue.JdbcDataStoreRepository = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue
			}
			if !idpBrowserSsoJitProvisioningUserRepositoryAttrs["ldap"].IsNull() {
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue := &client.LdapDataStoreRepository{}
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryAttrs := idpBrowserSsoJitProvisioningUserRepositoryAttrs["ldap"].(types.Object).Attributes()
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.BaseDn = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryAttrs["base_dn"].(types.String).ValueStringPointer()
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefValue := client.ResourceLink{}
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefAttrs := idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryAttrs["data_store_ref"].(types.Object).Attributes()
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefValue.Id = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefAttrs["id"].(types.String).ValueString()
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.DataStoreRef = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefValue
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.JitRepositoryAttributeMapping = map[string]client.AttributeFulfillmentValue{}
				for key, jitRepositoryAttributeMappingElement := range idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryAttrs["jit_repository_attribute_mapping"].(types.Map).Elements() {
					jitRepositoryAttributeMappingValue := client.AttributeFulfillmentValue{}
					jitRepositoryAttributeMappingAttrs := jitRepositoryAttributeMappingElement.(types.Object).Attributes()
					jitRepositoryAttributeMappingSourceValue := client.SourceTypeIdKey{}
					jitRepositoryAttributeMappingSourceAttrs := jitRepositoryAttributeMappingAttrs["source"].(types.Object).Attributes()
					jitRepositoryAttributeMappingSourceValue.Id = jitRepositoryAttributeMappingSourceAttrs["id"].(types.String).ValueStringPointer()
					jitRepositoryAttributeMappingSourceValue.Type = jitRepositoryAttributeMappingSourceAttrs["type"].(types.String).ValueString()
					jitRepositoryAttributeMappingValue.Source = jitRepositoryAttributeMappingSourceValue
					jitRepositoryAttributeMappingValue.Value = jitRepositoryAttributeMappingAttrs["value"].(types.String).ValueString()
					idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.JitRepositoryAttributeMapping[key] = jitRepositoryAttributeMappingValue
				}
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.Type = "LDAP"
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue.UniqueUserIdFilter = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryAttrs["unique_user_id_filter"].(types.String).ValueString()
				idpBrowserSsoJitProvisioningUserRepositoryValue.LdapDataStoreRepository = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue
			}
			idpBrowserSsoJitProvisioningValue.UserRepository = idpBrowserSsoJitProvisioningUserRepositoryValue
			idpBrowserSsoValue.JitProvisioning = idpBrowserSsoJitProvisioningValue
		}
		idpBrowserSsoValue.MessageCustomizations = []client.ProtocolMessageCustomization{}
		for _, messageCustomizationsElement := range idpBrowserSsoAttrs["message_customizations"].(types.Set).Elements() {
			messageCustomizationsValue := client.ProtocolMessageCustomization{}
			messageCustomizationsAttrs := messageCustomizationsElement.(types.Object).Attributes()
			messageCustomizationsValue.ContextName = messageCustomizationsAttrs["context_name"].(types.String).ValueStringPointer()
			messageCustomizationsValue.MessageExpression = messageCustomizationsAttrs["message_expression"].(types.String).ValueStringPointer()
			idpBrowserSsoValue.MessageCustomizations = append(idpBrowserSsoValue.MessageCustomizations, messageCustomizationsValue)
		}
		if !idpBrowserSsoAttrs["oauth_authentication_policy_contract_ref"].IsNull() {
			idpBrowserSsoOauthAuthenticationPolicyContractRefValue := &client.ResourceLink{}
			idpBrowserSsoOauthAuthenticationPolicyContractRefAttrs := idpBrowserSsoAttrs["oauth_authentication_policy_contract_ref"].(types.Object).Attributes()
			idpBrowserSsoOauthAuthenticationPolicyContractRefValue.Id = idpBrowserSsoOauthAuthenticationPolicyContractRefAttrs["id"].(types.String).ValueString()
			idpBrowserSsoValue.OauthAuthenticationPolicyContractRef = idpBrowserSsoOauthAuthenticationPolicyContractRefValue
		}
		if !idpBrowserSsoAttrs["oidc_provider_settings"].IsNull() {
			idpBrowserSsoOidcProviderSettingsValue := &client.OIDCProviderSettings{}
			idpBrowserSsoOidcProviderSettingsAttrs := idpBrowserSsoAttrs["oidc_provider_settings"].(types.Object).Attributes()
			idpBrowserSsoOidcProviderSettingsValue.AuthenticationScheme = idpBrowserSsoOidcProviderSettingsAttrs["authentication_scheme"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.AuthenticationSigningAlgorithm = idpBrowserSsoOidcProviderSettingsAttrs["authentication_signing_algorithm"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.AuthorizationEndpoint = idpBrowserSsoOidcProviderSettingsAttrs["authorization_endpoint"].(types.String).ValueString()
			idpBrowserSsoOidcProviderSettingsValue.BackChannelLogoutUri = idpBrowserSsoOidcProviderSettingsAttrs["back_channel_logout_uri"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.EnablePKCE = idpBrowserSsoOidcProviderSettingsAttrs["enable_pkce"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoOidcProviderSettingsValue.JwksURL = idpBrowserSsoOidcProviderSettingsAttrs["jwks_url"].(types.String).ValueString()
			idpBrowserSsoOidcProviderSettingsValue.JwtSecuredAuthorizationResponseModeType = idpBrowserSsoOidcProviderSettingsAttrs["jwt_secured_authorization_response_mode_type"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.LoginType = idpBrowserSsoOidcProviderSettingsAttrs["login_type"].(types.String).ValueString()
			idpBrowserSsoOidcProviderSettingsValue.LogoutEndpoint = idpBrowserSsoOidcProviderSettingsAttrs["logout_endpoint"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.PostLogoutRedirectUri = idpBrowserSsoOidcProviderSettingsAttrs["post_logout_redirect_uri"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.PushedAuthorizationRequestEndpoint = idpBrowserSsoOidcProviderSettingsAttrs["pushed_authorization_request_endpoint"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.RedirectUri = idpBrowserSsoOidcProviderSettingsAttrs["redirect_uri"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.RequestParameters = []client.OIDCRequestParameter{}
			for _, requestParametersElement := range idpBrowserSsoOidcProviderSettingsAttrs["request_parameters"].(types.Set).Elements() {
				requestParametersValue := client.OIDCRequestParameter{}
				requestParametersAttrs := requestParametersElement.(types.Object).Attributes()
				requestParametersValue.ApplicationEndpointOverride = requestParametersAttrs["application_endpoint_override"].(types.Bool).ValueBool()
				requestParametersAttributeValueValue := client.AttributeFulfillmentValue{}
				requestParametersAttributeValueAttrs := requestParametersAttrs["attribute_value"].(types.Object).Attributes()
				requestParametersAttributeValueSourceValue := client.SourceTypeIdKey{}
				requestParametersAttributeValueSourceAttrs := requestParametersAttributeValueAttrs["source"].(types.Object).Attributes()
				requestParametersAttributeValueSourceValue.Id = requestParametersAttributeValueSourceAttrs["id"].(types.String).ValueStringPointer()
				requestParametersAttributeValueSourceValue.Type = requestParametersAttributeValueSourceAttrs["type"].(types.String).ValueString()
				requestParametersAttributeValueValue.Source = requestParametersAttributeValueSourceValue
				requestParametersAttributeValueValue.Value = requestParametersAttributeValueAttrs["value"].(types.String).ValueString()
				requestParametersValue.AttributeValue = requestParametersAttributeValueValue
				requestParametersValue.Name = requestParametersAttrs["name"].(types.String).ValueString()
				requestParametersValue.Value = requestParametersAttrs["value"].(types.String).ValueStringPointer()
				idpBrowserSsoOidcProviderSettingsValue.RequestParameters = append(idpBrowserSsoOidcProviderSettingsValue.RequestParameters, requestParametersValue)
			}
			idpBrowserSsoOidcProviderSettingsValue.RequestSigningAlgorithm = idpBrowserSsoOidcProviderSettingsAttrs["request_signing_algorithm"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.Scopes = idpBrowserSsoOidcProviderSettingsAttrs["scopes"].(types.String).ValueString()
			idpBrowserSsoOidcProviderSettingsValue.TokenEndpoint = idpBrowserSsoOidcProviderSettingsAttrs["token_endpoint"].(types.String).ValueStringPointer()
			idpBrowserSsoOidcProviderSettingsValue.TrackUserSessionsForLogout = idpBrowserSsoOidcProviderSettingsAttrs["track_user_sessions_for_logout"].(types.Bool).ValueBoolPointer()
			idpBrowserSsoOidcProviderSettingsValue.UserInfoEndpoint = idpBrowserSsoOidcProviderSettingsAttrs["user_info_endpoint"].(types.String).ValueStringPointer()
			idpBrowserSsoValue.OidcProviderSettings = idpBrowserSsoOidcProviderSettingsValue
		}
		idpBrowserSsoValue.Protocol = idpBrowserSsoAttrs["protocol"].(types.String).ValueString()
		idpBrowserSsoValue.SignAuthnRequests = idpBrowserSsoAttrs["sign_authn_requests"].(types.Bool).ValueBoolPointer()
		idpBrowserSsoValue.SloServiceEndpoints = []client.SloServiceEndpoint{}
		for _, sloServiceEndpointsElement := range idpBrowserSsoAttrs["slo_service_endpoints"].(types.Set).Elements() {
			sloServiceEndpointsValue := client.SloServiceEndpoint{}
			sloServiceEndpointsAttrs := sloServiceEndpointsElement.(types.Object).Attributes()
			sloServiceEndpointsValue.Binding = sloServiceEndpointsAttrs["binding"].(types.String).ValueStringPointer()
			sloServiceEndpointsValue.ResponseUrl = sloServiceEndpointsAttrs["response_url"].(types.String).ValueStringPointer()
			sloServiceEndpointsValue.Url = sloServiceEndpointsAttrs["url"].(types.String).ValueString()
			idpBrowserSsoValue.SloServiceEndpoints = append(idpBrowserSsoValue.SloServiceEndpoints, sloServiceEndpointsValue)
		}
		idpBrowserSsoValue.SsoApplicationEndpoint = idpBrowserSsoAttrs["sso_application_endpoint"].(types.String).ValueStringPointer()
		if !idpBrowserSsoAttrs["sso_oauth_mapping"].IsNull() {
			idpBrowserSsoSsoOauthMappingValue := &client.SsoOAuthMapping{}
			idpBrowserSsoSsoOauthMappingAttrs := idpBrowserSsoAttrs["sso_oauth_mapping"].(types.Object).Attributes()
			idpBrowserSsoSsoOauthMappingValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(idpBrowserSsoSsoOauthMappingAttrs["attribute_contract_fulfillment"].(types.Map))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
			}
			idpBrowserSsoSsoOauthMappingValue.AttributeSources, err = attributesources.ClientStruct(idpBrowserSsoSsoOauthMappingAttrs["attribute_sources"].(types.Set))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_sources", err.Error())
			}
			idpBrowserSsoSsoOauthMappingValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(idpBrowserSsoSsoOauthMappingAttrs["issuance_criteria"].(types.Object))
			if err != nil {
				respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
			}
			idpBrowserSsoValue.SsoOAuthMapping = idpBrowserSsoSsoOauthMappingValue
		}
		idpBrowserSsoValue.SsoServiceEndpoints = []client.IdpSsoServiceEndpoint{}
		for _, ssoServiceEndpointsElement := range idpBrowserSsoAttrs["sso_service_endpoints"].(types.Set).Elements() {
			ssoServiceEndpointsValue := client.IdpSsoServiceEndpoint{}
			ssoServiceEndpointsAttrs := ssoServiceEndpointsElement.(types.Object).Attributes()
			ssoServiceEndpointsValue.Binding = ssoServiceEndpointsAttrs["binding"].(types.String).ValueStringPointer()
			ssoServiceEndpointsValue.Url = ssoServiceEndpointsAttrs["url"].(types.String).ValueString()
			idpBrowserSsoValue.SsoServiceEndpoints = append(idpBrowserSsoValue.SsoServiceEndpoints, ssoServiceEndpointsValue)
		}
		for _, urlWhitelistEntriesElement := range idpBrowserSsoAttrs["url_whitelist_entries"].(types.Set).Elements() {
			urlWhitelistEntriesValue := client.UrlWhitelistEntry{}
			urlWhitelistEntriesAttrs := urlWhitelistEntriesElement.(types.Object).Attributes()
			urlWhitelistEntriesValue.AllowQueryAndFragment = urlWhitelistEntriesAttrs["allow_query_and_fragment"].(types.Bool).ValueBoolPointer()
			urlWhitelistEntriesValue.RequireHttps = urlWhitelistEntriesAttrs["require_https"].(types.Bool).ValueBoolPointer()
			urlWhitelistEntriesValue.ValidDomain = urlWhitelistEntriesAttrs["valid_domain"].(types.String).ValueStringPointer()
			urlWhitelistEntriesValue.ValidPath = urlWhitelistEntriesAttrs["valid_path"].(types.String).ValueStringPointer()
			idpBrowserSsoValue.UrlWhitelistEntries = append(idpBrowserSsoValue.UrlWhitelistEntries, urlWhitelistEntriesValue)
		}
		addRequest.IdpBrowserSso = idpBrowserSsoValue
	}

	// attribute_query
	if !plan.AttributeQuery.IsNull() {
		attributeQueryValue := &client.IdpAttributeQuery{}
		attributeQueryAttrs := plan.AttributeQuery.Attributes()
		attributeQueryValue.NameMappings = []client.AttributeQueryNameMapping{}
		for _, nameMappingsElement := range attributeQueryAttrs["name_mappings"].(types.Set).Elements() {
			nameMappingsValue := client.AttributeQueryNameMapping{}
			nameMappingsAttrs := nameMappingsElement.(types.Object).Attributes()
			nameMappingsValue.LocalName = nameMappingsAttrs["local_name"].(types.String).ValueString()
			nameMappingsValue.RemoteName = nameMappingsAttrs["remote_name"].(types.String).ValueString()
			attributeQueryValue.NameMappings = append(attributeQueryValue.NameMappings, nameMappingsValue)
		}
		if !attributeQueryAttrs["policy"].IsNull() {
			attributeQueryPolicyValue := &client.IdpAttributeQueryPolicy{}
			attributeQueryPolicyAttrs := attributeQueryAttrs["policy"].(types.Object).Attributes()
			attributeQueryPolicyValue.EncryptNameId = attributeQueryPolicyAttrs["encrypt_name_id"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.MaskAttributeValues = attributeQueryPolicyAttrs["mask_attribute_values"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.RequireEncryptedAssertion = attributeQueryPolicyAttrs["require_encrypted_assertion"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.RequireSignedAssertion = attributeQueryPolicyAttrs["require_signed_assertion"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.RequireSignedResponse = attributeQueryPolicyAttrs["require_signed_response"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.SignAttributeQuery = attributeQueryPolicyAttrs["sign_attribute_query"].(types.Bool).ValueBoolPointer()
			attributeQueryValue.Policy = attributeQueryPolicyValue
		}
		attributeQueryValue.Url = attributeQueryAttrs["url"].(types.String).ValueString()
		addRequest.AttributeQuery = attributeQueryValue
	}

	// idp_oauth_grant_attribute_mapping
	if !plan.IdpOAuthGrantAttributeMapping.IsNull() {
		idpOauthGrantAttributeMappingValue := &client.IdpOAuthGrantAttributeMapping{}
		idpOauthGrantAttributeMappingAttrs := plan.IdpOAuthGrantAttributeMapping.Attributes()
		idpOauthGrantAttributeMappingValue.AccessTokenManagerMappings = []client.AccessTokenManagerMapping{}
		for _, accessTokenManagerMappingsElement := range idpOauthGrantAttributeMappingAttrs["access_token_manager_mappings"].(types.Set).Elements() {
			accessTokenManagerMappingsValue := client.AccessTokenManagerMapping{}
			accessTokenManagerMappingsAttrs := accessTokenManagerMappingsElement.(types.Object).Attributes()
			if !accessTokenManagerMappingsAttrs["access_token_manager_ref"].IsNull() {
				accessTokenManagerMappingsAccessTokenManagerRefValue := &client.ResourceLink{}
				accessTokenManagerMappingsAccessTokenManagerRefAttrs := accessTokenManagerMappingsAttrs["access_token_manager_ref"].(types.Object).Attributes()
				accessTokenManagerMappingsAccessTokenManagerRefValue.Id = accessTokenManagerMappingsAccessTokenManagerRefAttrs["id"].(types.String).ValueString()
				accessTokenManagerMappingsValue.AccessTokenManagerRef = accessTokenManagerMappingsAccessTokenManagerRefValue
			}
			accessTokenManagerMappingsValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(accessTokenManagerMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
			}
			accessTokenManagerMappingsValue.AttributeSources, err = attributesources.ClientStruct(accessTokenManagerMappingsAttrs["attribute_sources"].(types.Set))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_sources", err.Error())
			}
			accessTokenManagerMappingsValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(accessTokenManagerMappingsAttrs["issuance_criteria"].(types.Object))
			if err != nil {
				respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
			}
			idpOauthGrantAttributeMappingValue.AccessTokenManagerMappings = append(idpOauthGrantAttributeMappingValue.AccessTokenManagerMappings, accessTokenManagerMappingsValue)
		}
		if !idpOauthGrantAttributeMappingAttrs["idp_oauth_attribute_contract"].IsNull() {
			idpOauthGrantAttributeMappingIdpOauthAttributeContractValue := &client.IdpOAuthAttributeContract{}
			idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrs := idpOauthGrantAttributeMappingAttrs["idp_oauth_attribute_contract"].(types.Object).Attributes()
			idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.CoreAttributes = []client.IdpBrowserSsoAttribute{}
			for _, coreAttributesElement := range idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
				coreAttributesValue := client.IdpBrowserSsoAttribute{}
				coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
				coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
				idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.CoreAttributes = append(idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.CoreAttributes, coreAttributesValue)
			}
			idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.ExtendedAttributes = []client.IdpBrowserSsoAttribute{}
			for _, extendedAttributesElement := range idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.IdpBrowserSsoAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.ExtendedAttributes = append(idpOauthGrantAttributeMappingIdpOauthAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
			idpOauthGrantAttributeMappingValue.IdpOAuthAttributeContract = idpOauthGrantAttributeMappingIdpOauthAttributeContractValue
		}
		addRequest.IdpOAuthGrantAttributeMapping = idpOauthGrantAttributeMappingValue
	}

	// ws_trust
	if !plan.WsTrust.IsNull() {
		wsTrustValue := &client.IdpWsTrust{}
		wsTrustAttrs := plan.WsTrust.Attributes()
		wsTrustAttributeContractValue := client.IdpWsTrustAttributeContract{}
		wsTrustAttributeContractAttrs := wsTrustAttrs["attribute_contract"].(types.Object).Attributes()
		wsTrustAttributeContractValue.CoreAttributes = []client.IdpWsTrustAttribute{}
		for _, coreAttributesElement := range wsTrustAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.IdpWsTrustAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			wsTrustAttributeContractValue.CoreAttributes = append(wsTrustAttributeContractValue.CoreAttributes, coreAttributesValue)
		}
		wsTrustAttributeContractValue.ExtendedAttributes = []client.IdpWsTrustAttribute{}
		for _, extendedAttributesElement := range wsTrustAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.IdpWsTrustAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			wsTrustAttributeContractValue.ExtendedAttributes = append(wsTrustAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		wsTrustValue.AttributeContract = wsTrustAttributeContractValue
		wsTrustValue.GenerateLocalToken = wsTrustAttrs["generate_local_token"].(types.Bool).ValueBool()
		wsTrustValue.TokenGeneratorMappings = []client.SpTokenGeneratorMapping{}
		for _, tokenGeneratorMappingsElement := range wsTrustAttrs["token_generator_mappings"].(types.Set).Elements() {
			tokenGeneratorMappingsValue := client.SpTokenGeneratorMapping{}
			tokenGeneratorMappingsAttrs := tokenGeneratorMappingsElement.(types.Object).Attributes()
			tokenGeneratorMappingsValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(tokenGeneratorMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
			}
			tokenGeneratorMappingsValue.AttributeSources, err = attributesources.ClientStruct(tokenGeneratorMappingsAttrs["attribute_sources"].(types.Set))
			if err != nil {
				respDiags.AddError("Error building client struct for attribute_sources", err.Error())
			}
			tokenGeneratorMappingsValue.DefaultMapping = tokenGeneratorMappingsAttrs["default_mapping"].(types.Bool).ValueBoolPointer()
			tokenGeneratorMappingsValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(tokenGeneratorMappingsAttrs["issuance_criteria"].(types.Object))
			if err != nil {
				respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
			}
			if !tokenGeneratorMappingsAttrs["restricted_virtual_entity_ids"].IsNull() {
				tokenGeneratorMappingsValue.RestrictedVirtualEntityIds = []string{}
				for _, restrictedVirtualEntityIdsElement := range tokenGeneratorMappingsAttrs["restricted_virtual_entity_ids"].(types.Set).Elements() {
					tokenGeneratorMappingsValue.RestrictedVirtualEntityIds = append(tokenGeneratorMappingsValue.RestrictedVirtualEntityIds, restrictedVirtualEntityIdsElement.(types.String).ValueString())
				}
			}
			tokenGeneratorMappingsSpTokenGeneratorRefValue := client.ResourceLink{}
			tokenGeneratorMappingsSpTokenGeneratorRefAttrs := tokenGeneratorMappingsAttrs["sp_token_generator_ref"].(types.Object).Attributes()
			tokenGeneratorMappingsSpTokenGeneratorRefValue.Id = tokenGeneratorMappingsSpTokenGeneratorRefAttrs["id"].(types.String).ValueString()
			tokenGeneratorMappingsValue.SpTokenGeneratorRef = tokenGeneratorMappingsSpTokenGeneratorRefValue
			wsTrustValue.TokenGeneratorMappings = append(wsTrustValue.TokenGeneratorMappings, tokenGeneratorMappingsValue)
		}
		addRequest.WsTrust = wsTrustValue
	}

	// inbound_provisioning
	if !plan.InboundProvisioning.IsNull() {
		inboundProvisioningValue := &client.IdpInboundProvisioning{}
		inboundProvisioningAttrs := plan.InboundProvisioning.Attributes()
		inboundProvisioningValue.ActionOnDelete = inboundProvisioningAttrs["action_on_delete"].(types.String).ValueStringPointer()
		inboundProvisioningCustomSchemaValue := client.Schema{}
		inboundProvisioningCustomSchemaAttrs := inboundProvisioningAttrs["custom_schema"].(types.Object).Attributes()
		inboundProvisioningCustomSchemaValue.Attributes = []client.SchemaAttribute{}
		for _, attributesElement := range inboundProvisioningCustomSchemaAttrs["attributes"].(types.Set).Elements() {
			attributesValue := client.SchemaAttribute{}
			attributesAttrs := attributesElement.(types.Object).Attributes()
			attributesValue.MultiValued = attributesAttrs["multi_valued"].(types.Bool).ValueBoolPointer()
			attributesValue.Name = attributesAttrs["name"].(types.String).ValueStringPointer()
			if !attributesAttrs["sub_attributes"].IsNull() {
				attributesValue.SubAttributes = []string{}
				for _, subAttributesElement := range attributesAttrs["sub_attributes"].(types.Set).Elements() {
					attributesValue.SubAttributes = append(attributesValue.SubAttributes, subAttributesElement.(types.String).ValueString())
				}
			}
			if !attributesAttrs["types"].IsNull() {
				attributesValue.Types = []string{}
				for _, typesElement := range attributesAttrs["types"].(types.Set).Elements() {
					attributesValue.Types = append(attributesValue.Types, typesElement.(types.String).ValueString())
				}
			}
			inboundProvisioningCustomSchemaValue.Attributes = append(inboundProvisioningCustomSchemaValue.Attributes, attributesValue)
		}
		if !inboundProvisioningCustomSchemaAttrs["namespace"].IsUnknown() {
			inboundProvisioningCustomSchemaValue.Namespace = inboundProvisioningCustomSchemaAttrs["namespace"].(types.String).ValueStringPointer()
		}
		inboundProvisioningValue.CustomSchema = inboundProvisioningCustomSchemaValue
		inboundProvisioningValue.GroupSupport = inboundProvisioningAttrs["group_support"].(types.Bool).ValueBool()
		if !inboundProvisioningAttrs["groups"].IsNull() {
			inboundProvisioningGroupsValue := &client.Groups{}
			inboundProvisioningGroupsAttrs := inboundProvisioningAttrs["groups"].(types.Object).Attributes()
			inboundProvisioningGroupsReadGroupsValue := client.ReadGroups{}
			inboundProvisioningGroupsReadGroupsAttrs := inboundProvisioningGroupsAttrs["read_groups"].(types.Object).Attributes()
			inboundProvisioningGroupsReadGroupsAttributeContractValue := client.IdpInboundProvisioningAttributeContract{}
			inboundProvisioningGroupsReadGroupsAttributeContractAttrs := inboundProvisioningGroupsReadGroupsAttrs["attribute_contract"].(types.Object).Attributes()
			// PF requires core_attributes to be set, even though the property is read-only.
			// Provide a placeholder value here to prevent the API from returning an error.
			inboundProvisioningGroupsReadGroupsAttributeContractValue.CoreAttributes = []client.IdpInboundProvisioningAttribute{
				{
					Name: "placeholder",
				},
			}
			inboundProvisioningGroupsReadGroupsAttributeContractValue.ExtendedAttributes = []client.IdpInboundProvisioningAttribute{}
			for _, extendedAttributesElement := range inboundProvisioningGroupsReadGroupsAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.IdpInboundProvisioningAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				inboundProvisioningGroupsReadGroupsAttributeContractValue.ExtendedAttributes = append(inboundProvisioningGroupsReadGroupsAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsValue.AttributeContract = inboundProvisioningGroupsReadGroupsAttributeContractValue
			inboundProvisioningGroupsReadGroupsValue.AttributeFulfillment = map[string]client.AttributeFulfillmentValue{}
			for key, attributeFulfillmentElement := range inboundProvisioningGroupsReadGroupsAttrs["attribute_fulfillment"].(types.Map).Elements() {
				attributeFulfillmentValue := client.AttributeFulfillmentValue{}
				attributeFulfillmentAttrs := attributeFulfillmentElement.(types.Object).Attributes()
				attributeFulfillmentSourceValue := client.SourceTypeIdKey{}
				attributeFulfillmentSourceAttrs := attributeFulfillmentAttrs["source"].(types.Object).Attributes()
				attributeFulfillmentSourceValue.Id = attributeFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
				attributeFulfillmentSourceValue.Type = attributeFulfillmentSourceAttrs["type"].(types.String).ValueString()
				attributeFulfillmentValue.Source = attributeFulfillmentSourceValue
				attributeFulfillmentValue.Value = attributeFulfillmentAttrs["value"].(types.String).ValueString()
				inboundProvisioningGroupsReadGroupsValue.AttributeFulfillment[key] = attributeFulfillmentValue
			}
			inboundProvisioningGroupsReadGroupsValue.Attributes = []client.Attribute{}
			for _, attributesElement := range inboundProvisioningGroupsReadGroupsAttrs["attributes"].(types.Set).Elements() {
				attributesValue := client.Attribute{}
				attributesAttrs := attributesElement.(types.Object).Attributes()
				attributesValue.Name = attributesAttrs["name"].(types.String).ValueString()
				inboundProvisioningGroupsReadGroupsValue.Attributes = append(inboundProvisioningGroupsReadGroupsValue.Attributes, attributesValue)
			}
			inboundProvisioningGroupsValue.ReadGroups = inboundProvisioningGroupsReadGroupsValue
			inboundProvisioningGroupsWriteGroupsValue := client.WriteGroups{}
			inboundProvisioningGroupsWriteGroupsAttrs := inboundProvisioningGroupsAttrs["write_groups"].(types.Object).Attributes()
			inboundProvisioningGroupsWriteGroupsValue.AttributeFulfillment = map[string]client.AttributeFulfillmentValue{}
			for key, attributeFulfillmentElement := range inboundProvisioningGroupsWriteGroupsAttrs["attribute_fulfillment"].(types.Map).Elements() {
				attributeFulfillmentValue := client.AttributeFulfillmentValue{}
				attributeFulfillmentAttrs := attributeFulfillmentElement.(types.Object).Attributes()
				attributeFulfillmentSourceValue := client.SourceTypeIdKey{}
				attributeFulfillmentSourceAttrs := attributeFulfillmentAttrs["source"].(types.Object).Attributes()
				attributeFulfillmentSourceValue.Id = attributeFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
				attributeFulfillmentSourceValue.Type = attributeFulfillmentSourceAttrs["type"].(types.String).ValueString()
				attributeFulfillmentValue.Source = attributeFulfillmentSourceValue
				attributeFulfillmentValue.Value = attributeFulfillmentAttrs["value"].(types.String).ValueString()
				inboundProvisioningGroupsWriteGroupsValue.AttributeFulfillment[key] = attributeFulfillmentValue
			}
			inboundProvisioningGroupsValue.WriteGroups = inboundProvisioningGroupsWriteGroupsValue
			inboundProvisioningValue.Groups = inboundProvisioningGroupsValue
		}
		inboundProvisioningUserRepositoryValue := client.InboundProvisioningUserRepositoryAggregation{}
		inboundProvisioningUserRepositoryAttrs := inboundProvisioningAttrs["user_repository"].(types.Object).Attributes()
		if !inboundProvisioningUserRepositoryAttrs["identity_store"].IsNull() {
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryValue := &client.IdentityStoreInboundProvisioningUserRepository{}
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryAttrs := inboundProvisioningUserRepositoryAttrs["identity_store"].(types.Object).Attributes()
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryIdentityStoreProvisionerRefValue := client.ResourceLink{}
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryIdentityStoreProvisionerRefAttrs := inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryAttrs["identity_store_provisioner_ref"].(types.Object).Attributes()
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryIdentityStoreProvisionerRefValue.Id = inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryIdentityStoreProvisionerRefAttrs["id"].(types.String).ValueString()
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryValue.IdentityStoreProvisionerRef = inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryIdentityStoreProvisionerRefValue
			inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryValue.Type = "IDENTITY_STORE"
			inboundProvisioningUserRepositoryValue.IdentityStoreInboundProvisioningUserRepository = inboundProvisioningUserRepositoryIdentityStoreInboundProvisioningUserRepositoryValue
		}
		if !inboundProvisioningUserRepositoryAttrs["ldap"].IsNull() {
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue := &client.LdapInboundProvisioningUserRepository{}
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryAttrs := inboundProvisioningUserRepositoryAttrs["ldap"].(types.Object).Attributes()
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue.BaseDn = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryAttrs["base_dn"].(types.String).ValueStringPointer()
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryDataStoreRefValue := client.ResourceLink{}
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryDataStoreRefAttrs := inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryAttrs["data_store_ref"].(types.Object).Attributes()
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryDataStoreRefValue.Id = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryDataStoreRefAttrs["id"].(types.String).ValueString()
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue.DataStoreRef = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryDataStoreRefValue
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue.Type = "LDAP"
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue.UniqueGroupIdFilter = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryAttrs["unique_group_id_filter"].(types.String).ValueString()
			inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue.UniqueUserIdFilter = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryAttrs["unique_user_id_filter"].(types.String).ValueString()
			inboundProvisioningUserRepositoryValue.LdapInboundProvisioningUserRepository = inboundProvisioningUserRepositoryLdapInboundProvisioningUserRepositoryValue
		}
		inboundProvisioningValue.UserRepository = inboundProvisioningUserRepositoryValue
		inboundProvisioningUsersValue := client.Users{}
		inboundProvisioningUsersAttrs := inboundProvisioningAttrs["users"].(types.Object).Attributes()
		inboundProvisioningUsersReadUsersValue := client.ReadUsers{}
		inboundProvisioningUsersReadUsersAttrs := inboundProvisioningUsersAttrs["read_users"].(types.Object).Attributes()
		inboundProvisioningUsersReadUsersAttributeContractValue := client.IdpInboundProvisioningAttributeContract{}
		inboundProvisioningUsersReadUsersAttributeContractAttrs := inboundProvisioningUsersReadUsersAttrs["attribute_contract"].(types.Object).Attributes()
		// PF requires core_attributes to be set, even though the property is read-only.
		// Provide a placeholder value here to prevent the API from returning an error.
		inboundProvisioningUsersReadUsersAttributeContractValue.CoreAttributes = []client.IdpInboundProvisioningAttribute{
			{
				Name: "placeholder",
			},
		}
		inboundProvisioningUsersReadUsersAttributeContractValue.ExtendedAttributes = []client.IdpInboundProvisioningAttribute{}
		for _, extendedAttributesElement := range inboundProvisioningUsersReadUsersAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.IdpInboundProvisioningAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			inboundProvisioningUsersReadUsersAttributeContractValue.ExtendedAttributes = append(inboundProvisioningUsersReadUsersAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		inboundProvisioningUsersReadUsersValue.AttributeContract = inboundProvisioningUsersReadUsersAttributeContractValue
		inboundProvisioningUsersReadUsersValue.AttributeFulfillment = map[string]client.AttributeFulfillmentValue{}
		for key, attributeFulfillmentElement := range inboundProvisioningUsersReadUsersAttrs["attribute_fulfillment"].(types.Map).Elements() {
			attributeFulfillmentValue := client.AttributeFulfillmentValue{}
			attributeFulfillmentAttrs := attributeFulfillmentElement.(types.Object).Attributes()
			attributeFulfillmentSourceValue := client.SourceTypeIdKey{}
			attributeFulfillmentSourceAttrs := attributeFulfillmentAttrs["source"].(types.Object).Attributes()
			attributeFulfillmentSourceValue.Id = attributeFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
			attributeFulfillmentSourceValue.Type = attributeFulfillmentSourceAttrs["type"].(types.String).ValueString()
			attributeFulfillmentValue.Source = attributeFulfillmentSourceValue
			attributeFulfillmentValue.Value = attributeFulfillmentAttrs["value"].(types.String).ValueString()
			inboundProvisioningUsersReadUsersValue.AttributeFulfillment[key] = attributeFulfillmentValue
		}
		inboundProvisioningUsersReadUsersValue.Attributes = []client.Attribute{}
		for _, attributesElement := range inboundProvisioningUsersReadUsersAttrs["attributes"].(types.Set).Elements() {
			attributesValue := client.Attribute{}
			attributesAttrs := attributesElement.(types.Object).Attributes()
			attributesValue.Name = attributesAttrs["name"].(types.String).ValueString()
			inboundProvisioningUsersReadUsersValue.Attributes = append(inboundProvisioningUsersReadUsersValue.Attributes, attributesValue)
		}
		inboundProvisioningUsersValue.ReadUsers = inboundProvisioningUsersReadUsersValue
		inboundProvisioningUsersWriteUsersValue := client.WriteUsers{}
		inboundProvisioningUsersWriteUsersAttrs := inboundProvisioningUsersAttrs["write_users"].(types.Object).Attributes()
		inboundProvisioningUsersWriteUsersValue.AttributeFulfillment = map[string]client.AttributeFulfillmentValue{}
		for key, attributeFulfillmentElement := range inboundProvisioningUsersWriteUsersAttrs["attribute_fulfillment"].(types.Map).Elements() {
			attributeFulfillmentValue := client.AttributeFulfillmentValue{}
			attributeFulfillmentAttrs := attributeFulfillmentElement.(types.Object).Attributes()
			attributeFulfillmentSourceValue := client.SourceTypeIdKey{}
			attributeFulfillmentSourceAttrs := attributeFulfillmentAttrs["source"].(types.Object).Attributes()
			attributeFulfillmentSourceValue.Id = attributeFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
			attributeFulfillmentSourceValue.Type = attributeFulfillmentSourceAttrs["type"].(types.String).ValueString()
			attributeFulfillmentValue.Source = attributeFulfillmentSourceValue
			attributeFulfillmentValue.Value = attributeFulfillmentAttrs["value"].(types.String).ValueString()
			inboundProvisioningUsersWriteUsersValue.AttributeFulfillment[key] = attributeFulfillmentValue
		}
		inboundProvisioningUsersValue.WriteUsers = inboundProvisioningUsersWriteUsersValue
		inboundProvisioningValue.Users = inboundProvisioningUsersValue
		addRequest.InboundProvisioning = inboundProvisioningValue
	}

	return respDiags
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

func readSpIdpConnectionResponse(ctx context.Context, r *client.IdpConnection, plan, state *spIdpConnectionResourceModel, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics

	state.Active = types.BoolPointerValue(r.Active)
	state.AdditionalAllowedEntitiesConfiguration, diags = types.ObjectValueFrom(ctx, additionalAllowedEntitiesConfigurationAttrTypes, r.AdditionalAllowedEntitiesConfiguration)
	respDiags.Append(diags...)
	state.AttributeQuery, diags = types.ObjectValueFrom(ctx, attributeQueryAttrTypes, r.AttributeQuery)
	respDiags.Append(diags...)
	state.BaseUrl = types.StringPointerValue(r.BaseUrl)
	respDiags.Append(diags...)
	state.ConnectionId = types.StringPointerValue(r.Id)
	state.ContactInfo, diags = types.ObjectValueFrom(ctx, contactInfoAttrTypes, r.ContactInfo)
	respDiags.Append(diags...)
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	state.EntityId = types.StringValue(r.EntityId)
	state.ErrorPageMsgId = types.StringPointerValue(r.ErrorPageMsgId)
	state.ExtendedProperties, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: extendedPropertiesElemAttrTypes}, r.ExtendedProperties)
	respDiags.Append(diags...)
	state.Id = types.StringPointerValue(r.Id)
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	// If the plan logging mode does not match the state logging mode, report that the error might be being controlled
	// by the `server_settings_general` resource
	if plan != nil && plan.LoggingMode.ValueString() != state.LoggingMode.ValueString() {
		respDiags.AddAttributeError(path.Root("logging_mode"), providererror.ConflictingValueReturnedError,
			"PingFederate returned a different value for `logging_mode` for this resource than was planned. "+
				"If `idp_connection_transaction_logging_override` is configured to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource,"+
				" `logging_mode` should be configured to the same value in this resource.")
	}
	state.MetadataReloadSettings, diags = types.ObjectValueFrom(ctx, metadataReloadSettingsAttrTypes, r.MetadataReloadSettings)
	respDiags.Append(diags...)
	state.Name = types.StringValue(r.Name)
	if r.VirtualEntityIds == nil {
		if plan != nil && internaltypes.IsDefined(plan.VirtualEntityIds) && len(plan.VirtualEntityIds.Elements()) == 0 {
			state.VirtualEntityIds, diags = types.SetValue(types.StringType, nil)
			respDiags.Append(diags...)
		} else {
			state.VirtualEntityIds = types.SetNull(types.StringType)
		}
	} else {
		state.VirtualEntityIds = internaltypes.GetStringSet(r.VirtualEntityIds)
	}

	// LicenseConnectionGroup
	if r.LicenseConnectionGroup != nil {
		state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	}
	// Credentials
	var credentialsValue types.Object
	if r.Credentials != nil {
		var credentialsCertsValues []attr.Value
		for _, cert := range r.Credentials.Certs {
			if plan.Credentials.Attributes()["certs"] != nil {
				certMatchFound := false
				for _, certInPlan := range plan.Credentials.Attributes()["certs"].(types.List).Elements() {
					x509FilePlanAttrs := certInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
					x509FileIdPlan := x509FilePlanAttrs["id"].(types.String).ValueString()
					if cert.X509File.Id != nil && *cert.X509File.Id == x509FileIdPlan {
						planFileData := x509FilePlanAttrs["file_data"].(types.String)
						credentialsCertsObjValue, objDiags := connectioncert.ToState(ctx, planFileData, cert, &respDiags, isImportRead)
						respDiags.Append(objDiags...)
						credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
						certMatchFound = true
						break
					}
				}
				if !certMatchFound {
					credentialsCertsObjValue, objDiags := connectioncert.ToState(ctx, types.StringNull(), cert, &respDiags, isImportRead)
					respDiags.Append(objDiags...)
					credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
				}
			} else {
				credentialsCertsObjValue, objDiags := connectioncert.ToState(ctx, types.StringNull(), cert, &respDiags, isImportRead)
				respDiags.Append(objDiags...)
				credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
			}
		}
		credentialsCertsValue, objDiags := types.ListValue(connectioncert.ObjType(), credentialsCertsValues)
		respDiags.Append(objDiags...)
		var credentialsDecryptionKeyPairRefValue types.Object
		if r.Credentials.DecryptionKeyPairRef == nil {
			credentialsDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.DecryptionKeyPairRef)
			respDiags.Append(objDiags...)
		}
		var credentialsInboundBackChannelAuthValue types.Object
		if r.Credentials.InboundBackChannelAuth == nil {
			credentialsInboundBackChannelAuthValue = types.ObjectNull(credentialsInboundBackChannelAuthAttrTypes)
		} else {
			var credentialsInboundBackChannelAuthCertsValue types.List
			if len(r.Credentials.InboundBackChannelAuth.Certs) > 0 {
				var credentialsInboundBackChannelAuthCertsValues []attr.Value
				for _, ibcaCert := range r.Credentials.InboundBackChannelAuth.Certs {
					if plan.Credentials.Attributes()["inbound_back_channel_auth"] != nil {
						ibaCertMatch := false
						for _, ibcaCertInPlan := range plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["certs"].(types.List).Elements() {
							ibcax509FilePlanAttrs := ibcaCertInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
							ibcax509FileIdPlan := ibcax509FilePlanAttrs["id"].(types.String).ValueString()
							if ibcaCert.X509File.Id != nil && *ibcaCert.X509File.Id == ibcax509FileIdPlan {
								planIbcaX509FileFileData := ibcax509FilePlanAttrs["file_data"].(types.String)
								planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(ctx, planIbcaX509FileFileData, ibcaCert, &respDiags, isImportRead)
								respDiags.Append(objDiags...)
								credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
								ibaCertMatch = true
								break
							}
						}
						if !ibaCertMatch {
							planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(ctx, types.StringNull(), ibcaCert, &respDiags, isImportRead)
							respDiags.Append(objDiags...)
							credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
						}
					} else {
						planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(ctx, types.StringNull(), ibcaCert, &respDiags, isImportRead)
						respDiags.Append(objDiags...)
						credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
					}
				}
				credentialsInboundBackChannelAuthCertsValue, objDiags = types.ListValue(connectioncert.ObjType(), credentialsInboundBackChannelAuthCertsValues)
				respDiags.Append(objDiags...)
			} else {
				credentialsInboundBackChannelAuthCertsValue = types.ListNull(connectioncert.ObjType())
			}
			var credentialsInboundBackChannelAuthHttpBasicCredentialsValue types.Object
			if r.Credentials.InboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password *string
				if plan != nil && plan.Credentials.Attributes()["inbound_back_channel_auth"] != nil && plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					passwordFromPlan := plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					if internaltypes.IsDefined(passwordFromPlan) {
						password = passwordFromPlan.ValueStringPointer()
					} else if state != nil && internaltypes.IsDefined(state.Credentials) && state.Credentials.Attributes()["inbound_back_channel_auth"] != nil && state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
						password = state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueStringPointer()
					}
				} else if state != nil && internaltypes.IsDefined(state.Credentials) && state.Credentials.Attributes()["inbound_back_channel_auth"] != nil && state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					password = state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueStringPointer()
				}
				encryptedPassword := types.StringPointerValue(r.Credentials.InboundBackChannelAuth.HttpBasicCredentials.EncryptedPassword)
				if plan != nil && plan.Credentials.Attributes()["inbound_back_channel_auth"] != nil && plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					encryptedPasswordFromPlan := plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["encrypted_password"].(types.String)
					if internaltypes.IsDefined(encryptedPasswordFromPlan) {
						encryptedPassword = types.StringValue(encryptedPasswordFromPlan.ValueString())
					}
				}
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password":           types.StringPointerValue(password),
					"encrypted_password": encryptedPassword,
					"username":           types.StringPointerValue(r.Credentials.InboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				respDiags.Append(objDiags...)
			}
			credentialsInboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"certs":                   credentialsInboundBackChannelAuthCertsValue,
				"digital_signature":       types.BoolPointerValue(r.Credentials.InboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials":  credentialsInboundBackChannelAuthHttpBasicCredentialsValue,
				"require_ssl":             types.BoolPointerValue(r.Credentials.InboundBackChannelAuth.RequireSsl),
				"verification_issuer_dn":  types.StringPointerValue(r.Credentials.InboundBackChannelAuth.VerificationIssuerDN),
				"verification_subject_dn": types.StringPointerValue(r.Credentials.InboundBackChannelAuth.VerificationSubjectDN),
			})
			respDiags.Append(objDiags...)
		}
		var credentialsOutboundBackChannelAuthValue types.Object
		if r.Credentials.OutboundBackChannelAuth == nil {
			credentialsOutboundBackChannelAuthValue = types.ObjectNull(credentialsOutboundBackChannelAuthAttrTypes)
		} else {
			var credentialsOutboundBackChannelAuthHttpBasicCredentialsValue types.Object
			if r.Credentials.OutboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password *string
				if plan != nil && plan.Credentials.Attributes()["outbound_back_channel_auth"] != nil && plan.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					passwordFromPlan := plan.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					if internaltypes.IsDefined(passwordFromPlan) {
						password = passwordFromPlan.ValueStringPointer()
					} else if state != nil && internaltypes.IsDefined(state.Credentials) && state.Credentials.Attributes()["outbound_back_channel_auth"] != nil && state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
						password = state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueStringPointer()
					}
				} else if state != nil && internaltypes.IsDefined(state.Credentials) && state.Credentials.Attributes()["outbound_back_channel_auth"] != nil && state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					password = state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueStringPointer()
				}
				encryptedPassword := types.StringPointerValue(r.Credentials.OutboundBackChannelAuth.HttpBasicCredentials.EncryptedPassword)
				if plan != nil && plan.Credentials.Attributes()["outbound_back_channel_auth"] != nil && plan.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					encryptedPasswordFromPlan := plan.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["encrypted_password"].(types.String)
					if internaltypes.IsDefined(encryptedPasswordFromPlan) {
						encryptedPassword = types.StringValue(encryptedPasswordFromPlan.ValueString())
					}
				}
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password":           types.StringPointerValue(password),
					"encrypted_password": encryptedPassword,
					"username":           types.StringPointerValue(r.Credentials.OutboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				respDiags.Append(objDiags...)
			}
			var credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue types.Object
			if r.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef == nil {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
			} else {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef)
				respDiags.Append(objDiags...)
			}
			credentialsOutboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"digital_signature":      types.BoolPointerValue(r.Credentials.OutboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials": credentialsOutboundBackChannelAuthHttpBasicCredentialsValue,
				"ssl_auth_key_pair_ref":  credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue,
				"validate_partner_cert":  types.BoolPointerValue(r.Credentials.OutboundBackChannelAuth.ValidatePartnerCert),
			})
			respDiags.Append(objDiags...)
		}
		var credentialsSecondaryDecryptionKeyPairRefValue types.Object
		if r.Credentials.SecondaryDecryptionKeyPairRef == nil {
			credentialsSecondaryDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsSecondaryDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.SecondaryDecryptionKeyPairRef)
			respDiags.Append(objDiags...)
		}
		var credentialsSigningSettingsValue types.Object
		if r.Credentials.SigningSettings == nil {
			credentialsSigningSettingsValue = types.ObjectNull(credentialsSigningSettingsAttrTypes)
		} else {
			var credentialsSigningSettingsAlternativeSigningKeyPairRefsValues []attr.Value
			for _, credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue := range r.Credentials.SigningSettings.AlternativeSigningKeyPairRefs {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink := &client.ResourceLink{}
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink.Id = credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue.Id
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue, objDiags := resourcelink.ToState(ctx, credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink)
				respDiags.Append(objDiags...)
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValues = append(credentialsSigningSettingsAlternativeSigningKeyPairRefsValues, credentialsSigningSettingsAlternativeSigningKeyPairRefsValue)
			}
			var credentialsSigningSettingsAlternativeSigningKeyPairRefsValue types.Set
			if len(credentialsSigningSettingsAlternativeSigningKeyPairRefsValues) > 0 {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue, objDiags = types.SetValue(credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType, credentialsSigningSettingsAlternativeSigningKeyPairRefsValues)
				respDiags.Append(objDiags...)
			} else {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue = types.SetNull(credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType)
			}

			credentialsSigningSettingsSigningKeyPairRefValue, objDiags := resourcelink.ToState(ctx, &r.Credentials.SigningSettings.SigningKeyPairRef)
			respDiags.Append(objDiags...)

			includeCertInSignature := r.Credentials.SigningSettings.IncludeCertInSignature
			planSigningSettings := plan.Credentials.Attributes()["signing_settings"]
			if includeCertInSignature == nil && plan != nil && internaltypes.IsDefined(planSigningSettings) {
				planIncludeCertInSignature := planSigningSettings.(types.Object).Attributes()["include_cert_in_signature"]
				if internaltypes.IsDefined(planIncludeCertInSignature) && !planIncludeCertInSignature.(types.Bool).ValueBool() {
					includeCertInSignature = utils.Pointer(false)
				}
			}

			credentialsSigningSettingsValue, objDiags = types.ObjectValue(credentialsSigningSettingsAttrTypes, map[string]attr.Value{
				"algorithm":                         types.StringPointerValue(r.Credentials.SigningSettings.Algorithm),
				"alternative_signing_key_pair_refs": credentialsSigningSettingsAlternativeSigningKeyPairRefsValue,
				"include_cert_in_signature":         types.BoolPointerValue(includeCertInSignature),
				"include_raw_key_in_signature":      types.BoolPointerValue(r.Credentials.SigningSettings.IncludeRawKeyInSignature),
				"signing_key_pair_ref":              credentialsSigningSettingsSigningKeyPairRefValue,
			})
			respDiags.Append(objDiags...)
		}
		credentialsValue, objDiags = types.ObjectValue(credentialsAttrTypes, map[string]attr.Value{
			"block_encryption_algorithm":        types.StringPointerValue(r.Credentials.BlockEncryptionAlgorithm),
			"certs":                             credentialsCertsValue,
			"decryption_key_pair_ref":           credentialsDecryptionKeyPairRefValue,
			"inbound_back_channel_auth":         credentialsInboundBackChannelAuthValue,
			"key_transport_algorithm":           types.StringPointerValue(r.Credentials.KeyTransportAlgorithm),
			"outbound_back_channel_auth":        credentialsOutboundBackChannelAuthValue,
			"secondary_decryption_key_pair_ref": credentialsSecondaryDecryptionKeyPairRefValue,
			"signing_settings":                  credentialsSigningSettingsValue,
			"verification_issuer_dn":            types.StringPointerValue(r.Credentials.VerificationIssuerDN),
			"verification_subject_dn":           types.StringPointerValue(r.Credentials.VerificationSubjectDN),
		})
		respDiags.Append(objDiags...)
	} else {
		credentialsValue = types.ObjectNull(credentialsAttrTypes)
	}
	state.Credentials = credentialsValue

	// OidcClientCredentials
	if r.OidcClientCredentials != nil {
		var oidcClientCredentialsAttrValues map[string]attr.Value
		var clientSecret *string
		planOidcCredentialsAttrs := plan.OidcClientCredentials.Attributes()
		if planOidcCredentialsAttrs["client_secret"] != nil {
			clientSecret = plan.OidcClientCredentials.Attributes()["client_secret"].(types.String).ValueStringPointer()
		}

		encryptedSecretVal := planOidcCredentialsAttrs["encrypted_secret"]
		var encryptedSecretToState *string
		if encryptedSecretVal != nil && internaltypes.IsDefined(encryptedSecretVal) {
			encryptedSecretToState = encryptedSecretVal.(types.String).ValueStringPointer()
		} else {
			encryptedSecretToState = r.OidcClientCredentials.EncryptedSecret
		}

		oidcClientCredentialsAttrValues = map[string]attr.Value{
			"client_id":        types.StringValue(r.OidcClientCredentials.ClientId),
			"client_secret":    types.StringPointerValue(clientSecret),
			"encrypted_secret": types.StringPointerValue(encryptedSecretToState),
		}

		state.OidcClientCredentials, diags = types.ObjectValue(oidcClientCredentialsAttrTypes, oidcClientCredentialsAttrValues)
		respDiags.Append(diags...)
	} else {
		state.OidcClientCredentials = types.ObjectNull(oidcClientCredentialsAttrTypes)
	}

	// IdpBrowserSso
	var idpBrowserSsoValue types.Object
	if r.IdpBrowserSso == nil {
		idpBrowserSsoValue = types.ObjectNull(idpBrowserSsoAttrTypes)
	} else {
		var idpBrowserSsoAdapterMappingsValues []attr.Value
		idpBrowserSsoAdapterMappingsValue := types.ListNull(idpBrowserSsoAdapterMappingsElementType)
		if len(r.IdpBrowserSso.AdapterMappings) > 0 {
			for i, idpBrowserSsoAdapterMappingsResponseValue := range r.IdpBrowserSso.AdapterMappings {
				var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue types.Object
				if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings == nil {
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue = types.ObjectNull(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes)
				} else {
					var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue types.Object
					if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract == nil {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue = types.ObjectNull(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes)
					} else {
						var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues []attr.Value
						for _, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue := range idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.CoreAttributes {
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, objDiags := types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
								"name": types.StringValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Name),
							})
							respDiags.Append(objDiags...)
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues = append(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue)
						}
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, objDiags := types.SetValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues)
						respDiags.Append(objDiags...)
						var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues []attr.Value
						for _, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue := range idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.ExtendedAttributes {
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, objDiags := types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
								"name": types.StringValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Name),
							})
							respDiags.Append(objDiags...)
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues = append(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue)
						}
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, objDiags := types.SetValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues)
						respDiags.Append(objDiags...)
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue, objDiags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes, map[string]attr.Value{
							"core_attributes":     idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue,
							"extended_attributes": idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue,
						})
						respDiags.Append(objDiags...)
					}
					var planAdapterMappingOverrideConfiguration types.Object
					if internaltypes.IsDefined(plan.IdpBrowserSso) && internaltypes.IsDefined(plan.IdpBrowserSso.Attributes()["adapter_mappings"]) {
						adapterMappingsList, ok := plan.IdpBrowserSso.Attributes()["adapter_mappings"]
						if ok && len(adapterMappingsList.(types.List).Elements()) > i {
							overrideSettings, ok := adapterMappingsList.(types.List).Elements()[i].(types.Object).Attributes()["adapter_override_settings"]
							if ok {
								configValue, ok := overrideSettings.(types.Object).Attributes()["configuration"]
								if ok {
									planAdapterMappingOverrideConfiguration = configValue.(types.Object)
								}
							}
						}
					}
					adapterOverrideSettingsConfiguration := idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Configuration
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue, objDiags := pluginconfiguration.ToState(planAdapterMappingOverrideConfiguration, &adapterOverrideSettingsConfiguration, isImportRead)
					respDiags.Append(objDiags...)
					var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue types.Object
					if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef == nil {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue = types.ObjectNull(resourcelink.AttrType())
					} else {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue, objDiags = types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
							"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef.Id),
						})
						respDiags.Append(objDiags...)
					}
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue, objDiags := types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
						"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.PluginDescriptorRef.Id),
					})
					respDiags.Append(objDiags...)
					var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue types.Object
					if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo == nil {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue = types.ObjectNull(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes)
					} else {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue, objDiags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes, map[string]attr.Value{
							"application_icon_url": types.StringPointerValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo.ApplicationIconUrl),
							"application_name":     types.StringPointerValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo.ApplicationName),
						})
						respDiags.Append(objDiags...)
					}
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue, objDiags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes, map[string]attr.Value{
						"attribute_contract":      idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue,
						"configuration":           idpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue,
						"id":                      types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Id),
						"name":                    types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Name),
						"parent_ref":              idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue,
						"plugin_descriptor_ref":   idpBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue,
						"target_application_info": idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue,
					})
					respDiags.Append(objDiags...)
				}
				idpBrowserSsoAdapterMappingsResponseValueAttributeContractFulfillment := idpBrowserSsoAdapterMappingsResponseValue.AttributeContractFulfillment
				idpBrowserSsoAdapterMappingsAttributeContractFulfillmentValue, objDiags := attributecontractfulfillment.ToState(ctx, &idpBrowserSsoAdapterMappingsResponseValueAttributeContractFulfillment)
				respDiags.Append(objDiags...)

				idpBrowserSsoAdapterMappingsAttributeSourcesValue, objDiags := attributesources.ToState(ctx, idpBrowserSsoAdapterMappingsResponseValue.AttributeSources)
				respDiags.Append(objDiags...)

				var idpBrowserSsoAdapterMappingsIssuanceCriteriaValue types.Object
				if idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria != nil {
					idpBrowserSsoAdapterMappingsIssuanceCriteriaValue, objDiags = issuancecriteria.ToState(ctx, idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria)
					respDiags.Append(objDiags...)
				} else {
					idpBrowserSsoAdapterMappingsIssuanceCriteriaValue = types.ObjectNull(issuancecriteria.AttrTypes())
				}

				idpBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue, objDiags := types.SetValueFrom(ctx, types.StringType, idpBrowserSsoAdapterMappingsResponseValue.RestrictedVirtualEntityIds)
				respDiags.Append(objDiags...)
				var idpBrowserSsoAdapterMappingsSpAdapterRefValue types.Object
				if idpBrowserSsoAdapterMappingsResponseValue.SpAdapterRef == nil {
					idpBrowserSsoAdapterMappingsSpAdapterRefValue = types.ObjectNull(resourcelink.AttrType())
				} else {
					idpBrowserSsoAdapterMappingsSpAdapterRefValue, objDiags = types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
						"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.SpAdapterRef.Id),
					})
					respDiags.Append(objDiags...)
				}
				idpBrowserSsoAdapterMappingsValue, objDiags := types.ObjectValue(idpBrowserSsoAdapterMappingsAttrTypes, map[string]attr.Value{
					"adapter_override_settings":      idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue,
					"attribute_contract_fulfillment": idpBrowserSsoAdapterMappingsAttributeContractFulfillmentValue,
					"attribute_sources":              idpBrowserSsoAdapterMappingsAttributeSourcesValue,
					"issuance_criteria":              idpBrowserSsoAdapterMappingsIssuanceCriteriaValue,
					"restrict_virtual_entity_ids":    types.BoolPointerValue(idpBrowserSsoAdapterMappingsResponseValue.RestrictVirtualEntityIds),
					"restricted_virtual_entity_ids":  idpBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue,
					"sp_adapter_ref":                 idpBrowserSsoAdapterMappingsSpAdapterRefValue,
				})
				respDiags.Append(objDiags...)
				idpBrowserSsoAdapterMappingsValues = append(idpBrowserSsoAdapterMappingsValues, idpBrowserSsoAdapterMappingsValue)
			}
			idpBrowserSsoAdapterMappingsValue, diags = types.ListValue(idpBrowserSsoAdapterMappingsElementType, idpBrowserSsoAdapterMappingsValues)
			respDiags.Append(diags...)
		}

		// IdpBrowserSSO Always Sign Artifact Response
		var idpBrowserSsoAlwaysSignArtifactResponse types.Bool
		if r.IdpBrowserSso.AlwaysSignArtifactResponse == nil {
			idpBrowserSsoAlwaysSignArtifactResponse = types.BoolValue(false)
		} else {
			idpBrowserSsoAlwaysSignArtifactResponse = types.BoolPointerValue(r.IdpBrowserSso.AlwaysSignArtifactResponse)
		}

		// IdpBrowserSSO Artifact
		idpBrowserSsoArtifactValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoArtifactAttrTypes, r.IdpBrowserSso.Artifact)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO Attribute Contract
		idpBrowserSsoAttributeContractValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoAttributeContractAttrTypes, r.IdpBrowserSso.AttributeContract)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO Authentication Policy Contract Mappings
		var idpBrowserSsoAuthenticationPolicyContractMappingsValues []attr.Value
		for _, idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue := range r.IdpBrowserSso.AuthenticationPolicyContractMappings {
			idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.AttributeContractFulfillment)
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsAuthenticationPolicyContractRefValue, diags := types.ObjectValue(idpBrowserSsoAuthenticationPolicyContractMappingsAuthenticationPolicyContractRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.AuthenticationPolicyContractRef.Id),
			})
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsRestrictedVirtualServerIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.RestrictedVirtualServerIds)
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsValue, diags := types.ObjectValue(idpBrowserSsoAuthenticationPolicyContractMappingsAttrTypes, map[string]attr.Value{
				"attribute_contract_fulfillment":     idpBrowserSsoAuthenticationPolicyContractMappingsAttributeContractFulfillmentValue,
				"attribute_sources":                  idpBrowserSsoAuthenticationPolicyContractMappingsAttributeSourcesValue,
				"authentication_policy_contract_ref": idpBrowserSsoAuthenticationPolicyContractMappingsAuthenticationPolicyContractRefValue,
				"issuance_criteria":                  idpBrowserSsoAuthenticationPolicyContractMappingsIssuanceCriteriaValue,
				"restrict_virtual_server_ids":        types.BoolPointerValue(idpBrowserSsoAuthenticationPolicyContractMappingsResponseValue.RestrictVirtualServerIds),
				"restricted_virtual_server_ids":      idpBrowserSsoAuthenticationPolicyContractMappingsRestrictedVirtualServerIdsValue,
			})
			respDiags.Append(diags...)
			idpBrowserSsoAuthenticationPolicyContractMappingsValues = append(idpBrowserSsoAuthenticationPolicyContractMappingsValues, idpBrowserSsoAuthenticationPolicyContractMappingsValue)
		}
		idpBrowserSsoAuthenticationPolicyContractMappingsValue, diags := types.ListValue(idpBrowserSsoAuthenticationPolicyContractMappingsElementType, idpBrowserSsoAuthenticationPolicyContractMappingsValues)
		respDiags.Append(diags...)

		// IdpBrowserSSO AuthnContextMappings
		idpBrowserSsoAuthnContextMappingsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoAuthnContextMappingsElementType, r.IdpBrowserSso.AuthnContextMappings)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO Decryption Policy
		idpBrowserSsoDecryptionPolicyValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoDecryptionPolicyAttrTypes, r.IdpBrowserSso.DecryptionPolicy)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO Enabled Profiles
		idpBrowserSsoEnabledProfilesValue, objDiags := types.SetValueFrom(ctx, types.StringType, r.IdpBrowserSso.EnabledProfiles)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO Incoming Bindings
		idpBrowserSsoIncomingBindingsValue, objDiags := types.SetValueFrom(ctx, types.StringType, r.IdpBrowserSso.IncomingBindings)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO JIT Provisioning
		var idpBrowserSsoJitProvisioningValue types.Object
		if r.IdpBrowserSso.JitProvisioning == nil {
			idpBrowserSsoJitProvisioningValue = types.ObjectNull(idpBrowserSsoJitProvisioningAttrTypes)
		} else {
			attributeMappingSourceAttrTypes := map[string]attr.Type{
				"id":   types.StringType,
				"type": types.StringType,
			}
			attributeMappingAttrTypes := map[string]attr.Type{
				"source": types.ObjectType{AttrTypes: attributeMappingSourceAttrTypes},
				"value":  types.StringType,
			}
			sqlMethodStoredProcedureAttrTypes := map[string]attr.Type{
				"schema":           types.StringType,
				"stored_procedure": types.StringType,
			}
			sqlMethodTableAttrTypes := map[string]attr.Type{
				"schema":           types.StringType,
				"table_name":       types.StringType,
				"unique_id_column": types.StringType,
			}
			sqlMethodAttrTypes := map[string]attr.Type{
				"stored_procedure": types.ObjectType{AttrTypes: sqlMethodStoredProcedureAttrTypes},
				"table":            types.ObjectType{AttrTypes: sqlMethodTableAttrTypes},
			}

			var idpBrowserSsoJitProvisioningUserAttributesAttributeContractValues []attr.Value
			for _, idpBrowserSsoJitProvisioningUserAttributesAttributeContractResponseValue := range r.IdpBrowserSso.JitProvisioning.UserAttributes.AttributeContract {
				idpBrowserSsoJitProvisioningUserAttributesAttributeContractValue, objDiags := types.ObjectValue(idpBrowserSsoJitProvisioningUserAttributesAttributeContractAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(idpBrowserSsoJitProvisioningUserAttributesAttributeContractResponseValue.Masked),
					"name":   types.StringValue(idpBrowserSsoJitProvisioningUserAttributesAttributeContractResponseValue.Name),
				})
				respDiags.Append(objDiags...)
				idpBrowserSsoJitProvisioningUserAttributesAttributeContractValues = append(idpBrowserSsoJitProvisioningUserAttributesAttributeContractValues, idpBrowserSsoJitProvisioningUserAttributesAttributeContractValue)
			}
			idpBrowserSsoJitProvisioningUserAttributesAttributeContractValue, objDiags := types.SetValue(idpBrowserSsoJitProvisioningUserAttributesAttributeContractElementType, idpBrowserSsoJitProvisioningUserAttributesAttributeContractValues)
			respDiags.Append(objDiags...)
			idpBrowserSsoJitProvisioningUserAttributesValue, objDiags := types.ObjectValue(idpBrowserSsoJitProvisioningUserAttributesAttrTypes, map[string]attr.Value{
				"attribute_contract": idpBrowserSsoJitProvisioningUserAttributesAttributeContractValue,
				"do_attribute_query": types.BoolPointerValue(r.IdpBrowserSso.JitProvisioning.UserAttributes.DoAttributeQuery),
			})
			respDiags.Append(objDiags...)
			var idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue types.Object
			if r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository == nil {
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue = types.ObjectNull(datastorerepository.JdbcDataStoreRepositoryAttrType())
			} else {
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefValue, objDiags := types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
					"id": types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.DataStoreRef.Id),
				})
				respDiags.Append(objDiags...)
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValues := make(map[string]attr.Value)
				for key, idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingResponseValue := range r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.JitRepositoryAttributeMapping {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingSourceValue, objDiags := types.ObjectValue(attributeMappingSourceAttrTypes, map[string]attr.Value{
						"id":   types.StringPointerValue(idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Source.Id),
						"type": types.StringValue(idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Source.Type),
					})
					respDiags.Append(objDiags...)
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValue, objDiags := types.ObjectValue(attributeMappingAttrTypes, map[string]attr.Value{
						"source": idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingSourceValue,
						"value":  types.StringValue(idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Value),
					})
					respDiags.Append(objDiags...)
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValues[key] = idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValue
				}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValue, objDiags := types.MapValue(types.ObjectType{AttrTypes: attributeMappingAttrTypes}, idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValues)
				respDiags.Append(objDiags...)
				var idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue types.Object
				if r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.StoredProcedure == nil {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue = types.ObjectNull(sqlMethodStoredProcedureAttrTypes)
				} else {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue, objDiags = types.ObjectValue(sqlMethodStoredProcedureAttrTypes, map[string]attr.Value{
						"schema":           types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.StoredProcedure.Schema),
						"stored_procedure": types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.StoredProcedure.StoredProcedure),
					})
					respDiags.Append(objDiags...)
				}
				var idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue types.Object
				if r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.Table == nil {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue = types.ObjectNull(sqlMethodTableAttrTypes)
				} else {
					idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue, objDiags = types.ObjectValue(sqlMethodTableAttrTypes, map[string]attr.Value{
						"schema":           types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.Table.Schema),
						"table_name":       types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.Table.TableName),
						"unique_id_column": types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository.SqlMethod.Table.UniqueIdColumn),
					})
					respDiags.Append(objDiags...)
				}
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue, objDiags := types.ObjectValue(sqlMethodAttrTypes, map[string]attr.Value{
					"stored_procedure": idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodStoredProcedureValue,
					"table":            idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodTableValue,
				})
				respDiags.Append(objDiags...)
				idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue, objDiags = types.ObjectValue(datastorerepository.JdbcDataStoreRepositoryAttrType(), map[string]attr.Value{
					"data_store_ref":                   idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryDataStoreRefValue,
					"jit_repository_attribute_mapping": idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryJitRepositoryAttributeMappingValue,
					"sql_method":                       idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositorySqlMethodValue,
				})
				respDiags.Append(objDiags...)
			}
			var idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue types.Object
			if r.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository == nil {
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue = types.ObjectNull(datastorerepository.LdapDataStoreRepositoryAttrType())
			} else {
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefValue, objDiags := types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
					"id": types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository.DataStoreRef.Id),
				})
				respDiags.Append(objDiags...)
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValues := make(map[string]attr.Value)
				for key, idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingResponseValue := range r.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository.JitRepositoryAttributeMapping {
					idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingSourceValue, objDiags := types.ObjectValue(attributeMappingSourceAttrTypes, map[string]attr.Value{
						"id":   types.StringPointerValue(idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Source.Id),
						"type": types.StringValue(idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Source.Type),
					})
					respDiags.Append(objDiags...)
					idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValue, objDiags := types.ObjectValue(attributeMappingAttrTypes, map[string]attr.Value{
						"source": idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingSourceValue,
						"value":  types.StringValue(idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingResponseValue.Value),
					})
					respDiags.Append(objDiags...)
					idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValues[key] = idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValue
				}
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValue, objDiags := types.MapValue(types.ObjectType{AttrTypes: attributeMappingAttrTypes}, idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValues)
				respDiags.Append(objDiags...)
				idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue, objDiags = types.ObjectValue(datastorerepository.LdapDataStoreRepositoryAttrType(), map[string]attr.Value{
					"base_dn":                          types.StringPointerValue(r.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository.BaseDn),
					"data_store_ref":                   idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryDataStoreRefValue,
					"jit_repository_attribute_mapping": idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryJitRepositoryAttributeMappingValue,
					"unique_user_id_filter":            types.StringValue(r.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository.UniqueUserIdFilter),
				})
				respDiags.Append(objDiags...)
			}
			idpBrowserSsoJitProvisioningUserRepositoryValue, objDiags := types.ObjectValue(datastorerepository.ElemAttrType(), map[string]attr.Value{
				"jdbc": idpBrowserSsoJitProvisioningUserRepositoryJdbcDataStoreRepositoryValue,
				"ldap": idpBrowserSsoJitProvisioningUserRepositoryLdapDataStoreRepositoryValue,
			})
			respDiags.Append(objDiags...)
			idpBrowserSsoJitProvisioningValue, objDiags = types.ObjectValue(idpBrowserSsoJitProvisioningAttrTypes, map[string]attr.Value{
				"error_handling":  types.StringPointerValue(r.IdpBrowserSso.JitProvisioning.ErrorHandling),
				"event_trigger":   types.StringPointerValue(r.IdpBrowserSso.JitProvisioning.EventTrigger),
				"user_attributes": idpBrowserSsoJitProvisioningUserAttributesValue,
				"user_repository": idpBrowserSsoJitProvisioningUserRepositoryValue,
			})
			respDiags.Append(objDiags...)
		}

		// IdpBrowserSSO Message Customizations
		idpBrowserSsoMessageCustomizationsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoMessageCustomizationsElementType, r.IdpBrowserSso.MessageCustomizations)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO OAuth Authentication Policy Contract Ref
		idpBrowserSsoOauthAuthenticationPolicyContractRefValue, objDiags := resourcelink.ToState(ctx, r.IdpBrowserSso.OauthAuthenticationPolicyContractRef)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO OIDC Provider Settings
		idpBrowserSsoOidcProviderSettingsValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoOidcProviderSettingsAttrTypes, r.IdpBrowserSso.OidcProviderSettings)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO SLO Service Endpoints
		idpBrowserSsoSloServiceEndpointsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoSloServiceEndpointsElementType, r.IdpBrowserSso.SloServiceEndpoints)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO SSO Service Endpoints
		idpBrowserSsoSsoServiceEndpointsValue := types.SetNull(idpBrowserSsoSsoServiceEndpointsElementType)
		if len(r.IdpBrowserSso.SsoServiceEndpoints) > 0 {
			idpBrowserSsoSsoServiceEndpointsValue, objDiags = types.SetValueFrom(ctx, idpBrowserSsoSsoServiceEndpointsElementType, r.IdpBrowserSso.SsoServiceEndpoints)
			respDiags.Append(objDiags...)
		}

		// IdpBrowserSSO URL Whitelist Entries
		idpBrowserSsoUrlWhitelistEntriesValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoUrlWhitelistEntriesElementType, r.IdpBrowserSso.UrlWhitelistEntries)
		respDiags.Append(objDiags...)

		// IdpBrowserSSO SSO OAuth Mapping
		var idpBrowserSsoSsoOauthMappingValue types.Object
		if r.IdpBrowserSso.SsoOAuthMapping == nil {
			idpBrowserSsoSsoOauthMappingValue = types.ObjectNull(idpBrowserSsoSsoOauthMappingAttrTypes)
		} else {
			idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentValue, objDiags := attributecontractfulfillment.ToState(ctx, &r.IdpBrowserSso.SsoOAuthMapping.AttributeContractFulfillment)
			respDiags.Append(objDiags...)

			idpBrowserSsoSsoOauthMappingAttributeSourcesValue := types.SetNull(types.ObjectType{AttrTypes: attributesources.AttrTypes()})
			if r.IdpBrowserSso.SsoOAuthMapping.AttributeSources != nil {
				idpBrowserSsoSsoOauthMappingAttributeSourcesValue, objDiags = attributesources.ToState(ctx, r.IdpBrowserSso.SsoOAuthMapping.AttributeSources)
				respDiags.Append(objDiags...)
			}

			idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue := types.ObjectNull(issuancecriteria.AttrTypes())
			if r.IdpBrowserSso.SsoOAuthMapping.IssuanceCriteria != nil {
				idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue, objDiags = issuancecriteria.ToState(ctx, r.IdpBrowserSso.SsoOAuthMapping.IssuanceCriteria)
				respDiags.Append(objDiags...)
			}

			idpBrowserSsoSsoOauthMappingValue, objDiags = types.ObjectValue(idpBrowserSsoSsoOauthMappingAttrTypes, map[string]attr.Value{
				"attribute_contract_fulfillment": idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentValue,
				"attribute_sources":              idpBrowserSsoSsoOauthMappingAttributeSourcesValue,
				"issuance_criteria":              idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue,
			})
			respDiags.Append(objDiags...)
		}

		signAuthnRequest := r.IdpBrowserSso.SignAuthnRequests
		if signAuthnRequest == nil && plan != nil && internaltypes.IsDefined(plan.IdpBrowserSso) {
			planSignAuthnRequest := plan.IdpBrowserSso.Attributes()["sign_authn_requests"].(types.Bool)
			if internaltypes.IsDefined(planSignAuthnRequest) && !planSignAuthnRequest.ValueBool() {
				signAuthnRequest = utils.Pointer(false)
			}
		}

		assertionsSigned := r.IdpBrowserSso.AssertionsSigned
		if assertionsSigned == nil && plan != nil && internaltypes.IsDefined(plan.IdpBrowserSso) {
			planAssertionsSigned := plan.IdpBrowserSso.Attributes()["assertions_signed"].(types.Bool)
			if internaltypes.IsDefined(planAssertionsSigned) && !planAssertionsSigned.ValueBool() {
				assertionsSigned = utils.Pointer(false)
			}
		}

		idpBrowserSsoValue, objDiags = types.ObjectValue(idpBrowserSsoAttrTypes, map[string]attr.Value{
			"adapter_mappings":                         idpBrowserSsoAdapterMappingsValue,
			"always_sign_artifact_response":            idpBrowserSsoAlwaysSignArtifactResponse,
			"artifact":                                 idpBrowserSsoArtifactValue,
			"assertions_signed":                        types.BoolPointerValue(assertionsSigned),
			"attribute_contract":                       idpBrowserSsoAttributeContractValue,
			"authentication_policy_contract_mappings":  idpBrowserSsoAuthenticationPolicyContractMappingsValue,
			"authn_context_mappings":                   idpBrowserSsoAuthnContextMappingsValue,
			"decryption_policy":                        idpBrowserSsoDecryptionPolicyValue,
			"default_target_url":                       types.StringPointerValue(r.IdpBrowserSso.DefaultTargetUrl),
			"enabled_profiles":                         idpBrowserSsoEnabledProfilesValue,
			"idp_identity_mapping":                     types.StringValue(r.IdpBrowserSso.IdpIdentityMapping),
			"incoming_bindings":                        idpBrowserSsoIncomingBindingsValue,
			"jit_provisioning":                         idpBrowserSsoJitProvisioningValue,
			"message_customizations":                   idpBrowserSsoMessageCustomizationsValue,
			"oauth_authentication_policy_contract_ref": idpBrowserSsoOauthAuthenticationPolicyContractRefValue,
			"oidc_provider_settings":                   idpBrowserSsoOidcProviderSettingsValue,
			"protocol":                                 types.StringValue(r.IdpBrowserSso.Protocol),
			"sign_authn_requests":                      types.BoolPointerValue(signAuthnRequest),
			"slo_service_endpoints":                    idpBrowserSsoSloServiceEndpointsValue,
			"sso_application_endpoint":                 types.StringPointerValue(r.IdpBrowserSso.SsoApplicationEndpoint),
			"sso_oauth_mapping":                        idpBrowserSsoSsoOauthMappingValue,
			"sso_service_endpoints":                    idpBrowserSsoSsoServiceEndpointsValue,
			"url_whitelist_entries":                    idpBrowserSsoUrlWhitelistEntriesValue,
		})
		respDiags.Append(objDiags...)
	}
	state.IdpBrowserSso = idpBrowserSsoValue

	// idp_oauth_grant_attribute_mapping
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentAttrTypes}
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesAttrTypes}
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttrTypes := map[string]attr.Type{
		"access_token_manager_ref":       types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefAttrTypes},
		"attribute_contract_fulfillment": types.MapType{ElemType: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesElementType},
		"issuance_criteria":              types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingAccessTokenManagerMappingsIssuanceCriteriaAttrTypes},
	}
	idpOauthGrantAttributeMappingAccessTokenManagerMappingsElementType := types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttrTypes}
	idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesAttrTypes}
	idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesAttrTypes}
	idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesElementType},
	}
	idpOauthGrantAttributeMappingAttrTypes := map[string]attr.Type{
		"access_token_manager_mappings": types.SetType{ElemType: idpOauthGrantAttributeMappingAccessTokenManagerMappingsElementType},
		"idp_oauth_attribute_contract":  types.ObjectType{AttrTypes: idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrTypes},
	}
	var idpOauthGrantAttributeMappingValue types.Object
	if r.IdpOAuthGrantAttributeMapping == nil {
		idpOauthGrantAttributeMappingValue = types.ObjectNull(idpOauthGrantAttributeMappingAttrTypes)
	} else {
		var idpOauthGrantAttributeMappingAccessTokenManagerMappingsValues []attr.Value
		for _, idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue := range r.IdpOAuthGrantAttributeMapping.AccessTokenManagerMappings {
			var idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefValue types.Object
			if idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue.AccessTokenManagerRef == nil {
				idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefValue = types.ObjectNull(idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefAttrTypes)
			} else {
				idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefValue, diags = types.ObjectValue(idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefAttrTypes, map[string]attr.Value{
					"id": types.StringValue(idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue.AccessTokenManagerRef.Id),
				})
				respDiags.Append(diags...)
			}
			idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue.AttributeContractFulfillment)
			respDiags.Append(diags...)
			idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			idpOauthGrantAttributeMappingAccessTokenManagerMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), idpOauthGrantAttributeMappingAccessTokenManagerMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			idpOauthGrantAttributeMappingAccessTokenManagerMappingsValue, diags := types.ObjectValue(idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttrTypes, map[string]attr.Value{
				"access_token_manager_ref":       idpOauthGrantAttributeMappingAccessTokenManagerMappingsAccessTokenManagerRefValue,
				"attribute_contract_fulfillment": idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeContractFulfillmentValue,
				"attribute_sources":              idpOauthGrantAttributeMappingAccessTokenManagerMappingsAttributeSourcesValue,
				"issuance_criteria":              idpOauthGrantAttributeMappingAccessTokenManagerMappingsIssuanceCriteriaValue,
			})
			respDiags.Append(diags...)
			idpOauthGrantAttributeMappingAccessTokenManagerMappingsValues = append(idpOauthGrantAttributeMappingAccessTokenManagerMappingsValues, idpOauthGrantAttributeMappingAccessTokenManagerMappingsValue)
		}
		idpOauthGrantAttributeMappingAccessTokenManagerMappingsValue, diags := types.SetValue(idpOauthGrantAttributeMappingAccessTokenManagerMappingsElementType, idpOauthGrantAttributeMappingAccessTokenManagerMappingsValues)
		respDiags.Append(diags...)
		var idpOauthGrantAttributeMappingIdpOauthAttributeContractValue types.Object
		if r.IdpOAuthGrantAttributeMapping.IdpOAuthAttributeContract == nil {
			idpOauthGrantAttributeMappingIdpOauthAttributeContractValue = types.ObjectNull(idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrTypes)
		} else {
			var idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValues []attr.Value
			for _, idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesResponseValue := range r.IdpOAuthGrantAttributeMapping.IdpOAuthAttributeContract.CoreAttributes {
				idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValue, diags := types.ObjectValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesResponseValue.Masked),
					"name":   types.StringValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesResponseValue.Name),
				})
				respDiags.Append(diags...)
				idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValues = append(idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValues, idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValue)
			}
			idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValue, diags := types.SetValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesElementType, idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValues)
			respDiags.Append(diags...)
			var idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValues []attr.Value
			for _, idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesResponseValue := range r.IdpOAuthGrantAttributeMapping.IdpOAuthAttributeContract.ExtendedAttributes {
				idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValue, diags := types.ObjectValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesResponseValue.Masked),
					"name":   types.StringValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesResponseValue.Name),
				})
				respDiags.Append(diags...)
				idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValues = append(idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValues, idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValue)
			}
			idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValue, diags := types.SetValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesElementType, idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValues)
			respDiags.Append(diags...)
			idpOauthGrantAttributeMappingIdpOauthAttributeContractValue, diags = types.ObjectValue(idpOauthGrantAttributeMappingIdpOauthAttributeContractAttrTypes, map[string]attr.Value{
				"core_attributes":     idpOauthGrantAttributeMappingIdpOauthAttributeContractCoreAttributesValue,
				"extended_attributes": idpOauthGrantAttributeMappingIdpOauthAttributeContractExtendedAttributesValue,
			})
			respDiags.Append(diags...)
		}
		idpOauthGrantAttributeMappingValue, diags = types.ObjectValue(idpOauthGrantAttributeMappingAttrTypes, map[string]attr.Value{
			"access_token_manager_mappings": idpOauthGrantAttributeMappingAccessTokenManagerMappingsValue,
			"idp_oauth_attribute_contract":  idpOauthGrantAttributeMappingIdpOauthAttributeContractValue,
		})
		respDiags.Append(diags...)
	}

	state.IdpOAuthGrantAttributeMapping = idpOauthGrantAttributeMappingValue

	// InboundProvisioning
	var inboundProvisioningValue types.Object
	var inboundProvisioningGroupsValue types.Object
	if r.InboundProvisioning == nil {
		inboundProvisioningValue = types.ObjectNull(inboundProvisioningAttrTypes)
	} else {
		var inboundProvisioningCustomSchemaAttributesValues []attr.Value
		for _, inboundProvisioningCustomSchemaAttributesResponseValue := range r.InboundProvisioning.CustomSchema.Attributes {
			inboundProvisioningCustomSchemaAttributesSubAttributesValue, objDiags := types.SetValueFrom(ctx, types.StringType, inboundProvisioningCustomSchemaAttributesResponseValue.SubAttributes)
			respDiags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesTypesValue, objDiags := types.SetValueFrom(ctx, types.StringType, inboundProvisioningCustomSchemaAttributesResponseValue.Types)
			respDiags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesValue, objDiags := types.ObjectValue(inboundProvisioningCustomSchemaAttributesAttrTypes, map[string]attr.Value{
				"multi_valued":   types.BoolPointerValue(inboundProvisioningCustomSchemaAttributesResponseValue.MultiValued),
				"name":           types.StringPointerValue(inboundProvisioningCustomSchemaAttributesResponseValue.Name),
				"sub_attributes": inboundProvisioningCustomSchemaAttributesSubAttributesValue,
				"types":          inboundProvisioningCustomSchemaAttributesTypesValue,
			})
			respDiags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesValues = append(inboundProvisioningCustomSchemaAttributesValues, inboundProvisioningCustomSchemaAttributesValue)
		}
		inboundProvisioningCustomSchemaAttributesValue, objDiags := types.SetValue(inboundProvisioningCustomSchemaAttributesElementType, inboundProvisioningCustomSchemaAttributesValues)
		respDiags.Append(objDiags...)
		inboundProvisioningCustomSchemaValue, objDiags := types.ObjectValue(inboundProvisioningCustomSchemaAttrTypes, map[string]attr.Value{
			"attributes": inboundProvisioningCustomSchemaAttributesValue,
			"namespace":  types.StringPointerValue(r.InboundProvisioning.CustomSchema.Namespace),
		})
		respDiags.Append(objDiags...)
		if r.InboundProvisioning.Groups != nil {
			var inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeContract.CoreAttributes {
				inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue.Masked),
					"name":   types.StringValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue.Name),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesElementType, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues)
			respDiags.Append(objDiags...)
			var inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeContract.ExtendedAttributes {
				inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue.Masked),
					"name":   types.StringValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue.Name),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesElementType, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues)
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsAttributeContractValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractAttrTypes, map[string]attr.Value{
				"core_attributes":     inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue,
				"extended_attributes": inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue,
			})
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues := make(map[string]attr.Value)
			for key, inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeFulfillment {
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
					"id":   types.StringPointerValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Source.Id),
					"type": types.StringValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Source.Type),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentAttrTypes, map[string]attr.Value{
					"source": inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceValue,
					"value":  types.StringValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Value),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues[key] = inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue
			}
			inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentElementType, inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues)
			respDiags.Append(objDiags...)
			var inboundProvisioningGroupsReadGroupsAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.Attributes {
				inboundProvisioningGroupsReadGroupsAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributesAttrTypes, map[string]attr.Value{
					"name": types.StringValue(inboundProvisioningGroupsReadGroupsAttributesResponseValue.Name),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributesValues, inboundProvisioningGroupsReadGroupsAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributesElementType, inboundProvisioningGroupsReadGroupsAttributesValues)
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttrTypes, map[string]attr.Value{
				"attribute_contract":    inboundProvisioningGroupsReadGroupsAttributeContractValue,
				"attribute_fulfillment": inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue,
				"attributes":            inboundProvisioningGroupsReadGroupsAttributesValue,
			})
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues := make(map[string]attr.Value)
			for key, inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue := range r.InboundProvisioning.Groups.WriteGroups.AttributeFulfillment {
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
					"id":   types.StringPointerValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Source.Id),
					"type": types.StringValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Source.Type),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentAttrTypes, map[string]attr.Value{
					"source": inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceValue,
					"value":  types.StringValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Value),
				})
				respDiags.Append(objDiags...)
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues[key] = inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue
			}
			inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentElementType, inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues)
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsWriteGroupsValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttrTypes, map[string]attr.Value{
				"attribute_fulfillment": inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue,
			})
			respDiags.Append(objDiags...)
			inboundProvisioningGroupsValue, objDiags = types.ObjectValue(inboundProvisioningGroupsAttrTypes, map[string]attr.Value{
				"read_groups":  inboundProvisioningGroupsReadGroupsValue,
				"write_groups": inboundProvisioningGroupsWriteGroupsValue,
			})
			respDiags.Append(objDiags...)
		} else {
			inboundProvisioningGroupsValue = types.ObjectNull(inboundProvisioningGroupsAttrTypes)
		}

		var identityStoreInboundProvisioningUserRepository, ldapInboundProvisioningUserRepository types.Object
		if r.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository != nil {
			identityStoreProvisionerRef, objDiags := resourcelink.ToState(ctx, &r.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository.IdentityStoreProvisionerRef)
			respDiags.Append(objDiags...)
			identityStoreInboundProvisioningUserRepository, objDiags = types.ObjectValue(inboundprovisioninguserrepository.IdentityStoreInboundProvisioningUserRepositoryAttrType(), map[string]attr.Value{
				"identity_store_provisioner_ref": identityStoreProvisionerRef,
			})
			respDiags.Append(objDiags...)

			ldapInboundProvisioningUserRepository = types.ObjectNull(inboundprovisioninguserrepository.LdapInboundProvisioningUserRepositoryAttrType())
		} else if r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository != nil {
			dataStoreRef, objDiags := resourcelink.ToState(ctx, &r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.DataStoreRef)
			respDiags.Append(objDiags...)
			ldapInboundProvisioningUserRepository, objDiags = types.ObjectValue(inboundprovisioninguserrepository.LdapInboundProvisioningUserRepositoryAttrType(), map[string]attr.Value{
				"base_dn":                types.StringPointerValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.BaseDn),
				"data_store_ref":         dataStoreRef,
				"unique_user_id_filter":  types.StringValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.UniqueUserIdFilter),
				"unique_group_id_filter": types.StringValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.UniqueGroupIdFilter),
			})
			respDiags.Append(objDiags...)
			identityStoreInboundProvisioningUserRepository = types.ObjectNull(inboundprovisioninguserrepository.IdentityStoreInboundProvisioningUserRepositoryAttrType())
		}

		inboundProvisioningUserRepositoryAttrValue := map[string]attr.Value{
			"identity_store": identityStoreInboundProvisioningUserRepository,
			"ldap":           ldapInboundProvisioningUserRepository,
		}

		inboundProvisioningUserRepositoryValue, objDiags := types.ObjectValue(inboundprovisioninguserrepository.ElemAttrType(), inboundProvisioningUserRepositoryAttrValue)
		respDiags.Append(objDiags...)

		var inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeContract.CoreAttributes {
			inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue.Masked),
				"name":   types.StringValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue.Name),
			})
			respDiags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues = append(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesElementType, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues)
		respDiags.Append(objDiags...)
		var inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeContract.ExtendedAttributes {
			inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue.Masked),
				"name":   types.StringValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue.Name),
			})
			respDiags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues = append(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesElementType, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues)
		respDiags.Append(objDiags...)
		inboundProvisioningUsersReadUsersAttributeContractValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue,
			"extended_attributes": inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(objDiags...)
		inboundProvisioningUsersReadUsersAttributeFulfillmentValues := make(map[string]attr.Value)
		for key, inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeFulfillment {
			inboundProvisioningUsersReadUsersAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
				"id":   types.StringPointerValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Source.Id),
				"type": types.StringValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Source.Type),
			})
			respDiags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeFulfillmentAttrTypes, map[string]attr.Value{
				"source": inboundProvisioningUsersReadUsersAttributeFulfillmentSourceValue,
				"value":  types.StringValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Value),
			})
			respDiags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeFulfillmentValues[key] = inboundProvisioningUsersReadUsersAttributeFulfillmentValue
		}
		inboundProvisioningUsersReadUsersAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningUsersReadUsersAttributeFulfillmentElementType, inboundProvisioningUsersReadUsersAttributeFulfillmentValues)
		respDiags.Append(objDiags...)
		var inboundProvisioningUsersReadUsersAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.Attributes {
			inboundProvisioningUsersReadUsersAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(inboundProvisioningUsersReadUsersAttributesResponseValue.Name),
			})
			respDiags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributesValues = append(inboundProvisioningUsersReadUsersAttributesValues, inboundProvisioningUsersReadUsersAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributesElementType, inboundProvisioningUsersReadUsersAttributesValues)
		respDiags.Append(objDiags...)
		inboundProvisioningUsersReadUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttrTypes, map[string]attr.Value{
			"attribute_contract":    inboundProvisioningUsersReadUsersAttributeContractValue,
			"attribute_fulfillment": inboundProvisioningUsersReadUsersAttributeFulfillmentValue,
			"attributes":            inboundProvisioningUsersReadUsersAttributesValue,
		})
		respDiags.Append(objDiags...)

		inboundProvisioningUsersWriteUsersAttributeFulfillmentValue, objDiags := attributecontractfulfillment.ToState(ctx, &r.InboundProvisioning.Users.WriteUsers.AttributeFulfillment)
		respDiags.Append(objDiags...)
		inboundProvisioningUsersWriteUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersWriteUsersAttrTypes, map[string]attr.Value{
			"attribute_fulfillment": inboundProvisioningUsersWriteUsersAttributeFulfillmentValue,
		})
		respDiags.Append(objDiags...)
		inboundProvisioningUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersAttrTypes, map[string]attr.Value{
			"read_users":  inboundProvisioningUsersReadUsersValue,
			"write_users": inboundProvisioningUsersWriteUsersValue,
		})
		respDiags.Append(objDiags...)
		inboundProvisioningValue, objDiags = types.ObjectValue(inboundProvisioningAttrTypes, map[string]attr.Value{
			"action_on_delete": types.StringPointerValue(r.InboundProvisioning.ActionOnDelete),
			"custom_schema":    inboundProvisioningCustomSchemaValue,
			"group_support":    types.BoolValue(r.InboundProvisioning.GroupSupport),
			"groups":           inboundProvisioningGroupsValue,
			"user_repository":  inboundProvisioningUserRepositoryValue,
			"users":            inboundProvisioningUsersValue,
		})
		respDiags.Append(objDiags...)
	}
	state.InboundProvisioning = inboundProvisioningValue

	// WsTrust
	if r.WsTrust != nil {
		var tokenGeneratorMappings []basetypes.ObjectValue
		for _, tokenGeneratorMapping := range r.WsTrust.TokenGeneratorMappings {
			tokenGeneratorMappingSpTokenGeneratorRef := tokenGeneratorMapping.SpTokenGeneratorRef
			spTokenGeneratorRef, objDiags := resourcelink.ToState(ctx, &tokenGeneratorMappingSpTokenGeneratorRef)
			respDiags.Append(objDiags...)

			var attributeSources basetypes.SetValue
			attributeSources, objDiags = attributesources.ToState(ctx, tokenGeneratorMapping.AttributeSources)
			respDiags.Append(objDiags...)

			tokenGeneratorMappingAttributeContractFulfillment := tokenGeneratorMapping.AttributeContractFulfillment
			attributeContractFulfillment, objDiags := attributecontractfulfillment.ToState(ctx, &tokenGeneratorMappingAttributeContractFulfillment)
			respDiags.Append(objDiags...)

			issuanceCriteria, objDiags := issuancecriteria.ToState(ctx, tokenGeneratorMapping.IssuanceCriteria)
			respDiags.Append(objDiags...)

			var restrictedVirtualEntityIds types.Set
			if len(tokenGeneratorMapping.RestrictedVirtualEntityIds) > 0 {
				restrictedVirtualEntityIds, objDiags = types.SetValueFrom(ctx, types.StringType, tokenGeneratorMapping.RestrictedVirtualEntityIds)
				respDiags.Append(objDiags...)
			} else {
				restrictedVirtualEntityIds = types.SetNull(types.StringType)
			}

			tokenGeneratorAttrValues := map[string]attr.Value{
				"attribute_contract_fulfillment": attributeContractFulfillment,
				"attribute_sources":              attributeSources,
				"default_mapping":                types.BoolPointerValue(tokenGeneratorMapping.DefaultMapping),
				"issuance_criteria":              issuanceCriteria,
				"sp_token_generator_ref":         spTokenGeneratorRef,
				"restricted_virtual_entity_ids":  restrictedVirtualEntityIds,
			}

			tokenGeneratorMappingState, objDiags := types.ObjectValue(tokenGeneratorAttrTypes, tokenGeneratorAttrValues)
			respDiags.Append(objDiags...)
			tokenGeneratorMappings = append(tokenGeneratorMappings, tokenGeneratorMappingState)
		}

		attributeContract, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"core_attributes": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":   types.StringType,
					"masked": types.BoolType,
				},
			}},
			"extended_attributes": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":   types.StringType,
					"masked": types.BoolType,
				},
			}},
		}, r.WsTrust.AttributeContract)
		respDiags.Append(objDiags...)

		var tokenGeneratorMappingsSet types.Set
		if tokenGeneratorMappings != nil {
			tokenGeneratorMappingsSet, objDiags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: tokenGeneratorAttrTypes}, tokenGeneratorMappings)
			respDiags.Append(objDiags...)
		} else {
			tokenGeneratorMappingsSet = types.SetNull(types.ObjectType{AttrTypes: tokenGeneratorAttrTypes})
		}

		wsTrustAttrValues := map[string]attr.Value{
			"attribute_contract":       attributeContract,
			"generate_local_token":     types.BoolValue(r.WsTrust.GenerateLocalToken),
			"token_generator_mappings": tokenGeneratorMappingsSet,
		}

		state.WsTrust, objDiags = types.ObjectValue(wsTrustAttrTypes, wsTrustAttrValues)
		respDiags.Append(objDiags...)
	} else {
		state.WsTrust = types.ObjectNull(wsTrustAttrTypes)
	}

	return respDiags
}

func (r *spIdpConnectionResource) warnFor500Err(httpResp *http.Response, diags *diag.Diagnostics) {
	if httpResp != nil && httpResp.StatusCode == 500 {
		diags.AddError(providererror.PingFederateAPIError, "PingFederate API returned a 500 error. Due to a known issue in PingFederate, you may have to set the `virtual_entity_ids` attribute to `[]` to work around this issue")
	}
}

func (r *spIdpConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state spIdpConnectionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSpIdpConnection := client.NewIdpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(addOptionalSpIdpConnectionFields(ctx, createSpIdpConnection, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCreateSpIdpConnection := r.apiClient.SpIdpConnectionsAPI.CreateConnection(config.AuthContext(ctx, r.providerConfig))
	apiCreateSpIdpConnection = apiCreateSpIdpConnection.Body(*createSpIdpConnection)
	spIdpConnectionResponse, httpResp, err := r.apiClient.SpIdpConnectionsAPI.CreateConnectionExecute(apiCreateSpIdpConnection)
	if err != nil {
		r.warnFor500Err(httpResp, &resp.Diagnostics)
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the SpIdpConnection", err, httpResp, &customId)
		return
	}

	diags = readSpIdpConnectionResponse(ctx, spIdpConnectionResponse, &plan, &state, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *spIdpConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state spIdpConnectionResourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSpIdpConnection, httpResp, err := r.apiClient.SpIdpConnectionsAPI.GetConnection(config.AuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Sp Idp Connection", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting a Sp Idp Connection", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	diags = readSpIdpConnectionResponse(ctx, apiReadSpIdpConnection, &state, &state, isImportRead)
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

	updateSpIdpConnection := r.apiClient.SpIdpConnectionsAPI.UpdateConnection(config.AuthContext(ctx, r.providerConfig), plan.ConnectionId.ValueString())
	createUpdateRequest := client.NewIdpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(addOptionalSpIdpConnectionFields(ctx, createUpdateRequest, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSpIdpConnection = updateSpIdpConnection.Body(*createUpdateRequest)
	updateSpIdpConnectionResponse, httpResp, err := r.apiClient.SpIdpConnectionsAPI.UpdateConnectionExecute(updateSpIdpConnection)
	if err != nil {
		r.warnFor500Err(httpResp, &resp.Diagnostics)
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating Sp Idp Connection", err, httpResp, &customId)
		return
	}

	// Read the response
	var state spIdpConnectionResourceModel
	diags = readSpIdpConnectionResponse(ctx, updateSpIdpConnectionResponse, &plan, &state, false)
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
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting a Sp Idp Connection", err, httpResp, &customId)
	}
}

func (r *spIdpConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to connection_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("connection_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
