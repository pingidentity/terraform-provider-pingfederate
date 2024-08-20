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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/datastorerepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/inboundprovisioninguserrepository"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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
		"client_id":     types.StringType,
		"client_secret": types.StringType,
	}

	credentialsCertsCertViewAttrTypes = map[string]attr.Type{
		"crypto_provider":           types.StringType,
		"expires":                   types.StringType,
		"id":                        types.StringType,
		"issuer_dn":                 types.StringType,
		"key_algorithm":             types.StringType,
		"key_size":                  types.Int64Type,
		"serial_number":             types.StringType,
		"sha1_fingerprint":          types.StringType,
		"sha256_fingerprint":        types.StringType,
		"signature_algorithm":       types.StringType,
		"status":                    types.StringType,
		"subject_alternative_names": types.SetType{ElemType: types.StringType},
		"subject_dn":                types.StringType,
		"valid_from":                types.StringType,
		"version":                   types.Int64Type,
	}
	credentialsCertsX509fileAttrTypes = map[string]attr.Type{
		"crypto_provider":     types.StringType,
		"formatted_file_data": types.StringType,
		"file_data":           types.StringType,
		"id":                  types.StringType,
	}

	credentialsInboundBackChannelAuthCertsCertViewAttrTypes = map[string]attr.Type{
		"crypto_provider":           types.StringType,
		"expires":                   types.StringType,
		"id":                        types.StringType,
		"issuer_dn":                 types.StringType,
		"key_algorithm":             types.StringType,
		"key_size":                  types.Int64Type,
		"serial_number":             types.StringType,
		"sha1_fingerprint":          types.StringType,
		"sha256_fingerprint":        types.StringType,
		"signature_algorithm":       types.StringType,
		"status":                    types.StringType,
		"subject_alternative_names": types.SetType{ElemType: types.StringType},
		"subject_dn":                types.StringType,
		"valid_from":                types.StringType,
		"version":                   types.Int64Type,
	}
	credentialsInboundBackChannelAuthCertsX509fileAttrTypes = map[string]attr.Type{
		"crypto_provider": types.StringType,
		"file_data":       types.StringType,
		"id":              types.StringType,
	}

	credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes = map[string]attr.Type{
		"password": types.StringType,
		"username": types.StringType,
	}
	credentialsInboundBackChannelAuthAttrTypes = map[string]attr.Type{
		"certs":                   types.ListType{ElemType: connectioncert.ObjType()},
		"digital_signature":       types.BoolType,
		"http_basic_credentials":  types.ObjectType{AttrTypes: credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes},
		"require_ssl":             types.BoolType,
		"type":                    types.StringType,
		"verification_issuer_dn":  types.StringType,
		"verification_subject_dn": types.StringType,
	}
	credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes = map[string]attr.Type{
		"password": types.StringType,
		"username": types.StringType,
	}

	credentialsOutboundBackChannelAuthAttrTypes = map[string]attr.Type{
		"digital_signature":      types.BoolType,
		"http_basic_credentials": types.ObjectType{AttrTypes: credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes},
		"ssl_auth_key_pair_ref":  types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"type":                   types.StringType,
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
		"authentication_scheme":                 types.StringType,
		"authentication_signing_algorithm":      types.StringType,
		"authorization_endpoint":                types.StringType,
		"back_channel_logout_uri":               types.StringType,
		"enable_pkce":                           types.BoolType,
		"front_channel_logout_uri":              types.StringType,
		"jwks_url":                              types.StringType,
		"login_type":                            types.StringType,
		"logout_endpoint":                       types.StringType,
		"post_logout_redirect_uri":              types.StringType,
		"pushed_authorization_request_endpoint": types.StringType,
		"redirect_uri":                          types.StringType,
		"request_parameters":                    types.SetType{ElemType: idpBrowserSsoOidcProviderSettingsRequestParametersElementType},
		"request_signing_algorithm":             types.StringType,
		"scopes":                                types.StringType,
		"token_endpoint":                        types.StringType,
		"track_user_sessions_for_logout":        types.BoolType,
		"user_info_endpoint":                    types.StringType,
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

	conditionalCriteriaDefault, _ = types.SetValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source": types.ObjectType{
				AttrTypes: sourcetypeidkey.AttrTypes(),
			},
			"attribute_name": types.StringType,
			"condition":      types.StringType,
			"value":          types.StringType,
			"error_result":   types.StringType,
		},
	}, nil)

	idpBrowserSsoUrlWhitelistEntriesElementType = types.ObjectType{AttrTypes: idpBrowserSsoUrlWhitelistEntriesAttrTypes}
	idpBrowserSsoAttrTypes                      = map[string]attr.Type{
		"adapter_mappings":                         types.ListType{ElemType: idpBrowserSsoAdapterMappingsElementType},
		"always_sign_artifact_response":            types.BoolType,
		"artifact":                                 types.ObjectType{AttrTypes: idpBrowserSsoArtifactAttrTypes},
		"assertions_signed":                        types.BoolType,
		"attribute_contract":                       types.ObjectType{AttrTypes: idpBrowserSsoAttributeContractAttrTypes},
		"authentication_policy_contract_mappings":  types.SetType{ElemType: idpBrowserSsoAuthenticationPolicyContractMappingsElementType},
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

	accessTokenManagerMappingAttrTypes = map[string]attr.Type{
		"access_token_manager_ref":       types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"attribute_sources":              types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria":              types.ObjectType{AttrTypes: issuancecriteria.AttrTypes()},
	}

	idpOAuthGrantAttributeMappingAttrTypes = map[string]attr.Type{
		"access_token_manager_mappings": types.SetType{ElemType: types.ObjectType{AttrTypes: accessTokenManagerMappingAttrTypes}},
		"idp_oauth_attribute_contract": types.ObjectType{AttrTypes: map[string]attr.Type{
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
		"token_generator_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: spTokenGeneratorMappingAttrTypes}},
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
	emptyStringSet, _ = types.SetValue(types.StringType, []attr.Value{})
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
		Description: "Resource to create and manage a SP Idp Connection",
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
					"name_mappings": schema.SetNestedAttribute{
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
						Description:         "The algorithm used to encrypt assertions sent to this partner. AES_128, AES_256, AES_128_GCM, AES_192_GCM, AES_256_GCM and Triple_DES are supported.",
						MarkdownDescription: "The algorithm used to encrypt assertions sent to this partner. AES_128, AES_256, AES_128_GCM, AES_192_GCM, AES_256_GCM and Triple_DES are supported.",
						Validators: []validator.String{
							stringvalidator.OneOf("AES_128", "AES_256", "AES_128_GCM", "AES_192_GCM", "AES_256_GCM", "Triple_DES"),
						},
					},
					"key_transport_algorithm": schema.StringAttribute{
						Optional:            true,
						Description:         "The algorithm used to transport keys to this partner. RSA_OAEP, RSA_OAEP_256 and RSA_v15 are supported.",
						MarkdownDescription: "The algorithm used to transport keys to this partner. RSA_OAEP, RSA_OAEP_256 and RSA_v15 are supported.",
						Validators: []validator.String{
							stringvalidator.OneOf("RSA_OAEP", "RSA_OAEP_256", "RSA_v15"),
						},
					},
					"signing_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"signing_key_pair_ref": schema.SingleNestedAttribute{
								Attributes:          resourcelink.ToSchema(),
								Optional:            true,
								Description:         "A reference to a resource.",
								MarkdownDescription: "A reference to a resource.",
							},
							"alternative_signing_key_pair_refs": schema.SetAttribute{
								ElementType:         types.ObjectType{AttrTypes: resourcelink.AttrType()},
								Optional:            true,
								Description:         "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
								MarkdownDescription: "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
							},
							"algorithm": schema.StringAttribute{
								Optional:            true,
								Description:         "The algorithm used to sign messages sent to this partner. The default is SHA1withDSA for DSA certs, SHA256withRSA for RSA certs, and SHA256withECDSA for EC certs. For RSA certs, SHA1withRSA, SHA384withRSA, SHA512withRSA, SHA256withRSAandMGF1, SHA384withRSAandMGF1 and SHA512withRSAandMGF1 are also supported. For EC certs, SHA384withECDSA and SHA512withECDSA are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
								MarkdownDescription: "The algorithm used to sign messages sent to this partner. The default is SHA1withDSA for DSA certs, SHA256withRSA for RSA certs, and SHA256withECDSA for EC certs. For RSA certs, SHA1withRSA, SHA384withRSA, SHA512withRSA, SHA256withRSAandMGF1, SHA384withRSAandMGF1 and SHA512withRSAandMGF1 are also supported. For EC certs, SHA384withECDSA and SHA512withECDSA are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
							},
							"include_cert_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "Determines whether the signing certificate is included in the signature element.",
								MarkdownDescription: "Determines whether the signing certificate is included in the signature element.",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "Determines whether the element with the raw public key is included in the signature element.",
								MarkdownDescription: "Determines whether the element with the raw public key is included in the signature element.",
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
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "The back channel authentication type.",
								MarkdownDescription: "The back channel authentication type.",
								Validators: []validator.String{
									stringvalidator.OneOf("INBOUND", "OUTBOUND"),
								},
							},
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Optional:            true,
										Description:         "The username.",
										MarkdownDescription: "The username.",
									},
									"password": schema.StringAttribute{
										Optional:            true,
										Sensitive:           false,
										Description:         "User password. To update the password, specify the plaintext value in this field.",
										MarkdownDescription: "User password. To update the password, specify the plaintext value in this field.",
									},
								},
								Optional:            true,
								Description:         "Username and password credentials.",
								MarkdownDescription: "Username and password credentials.",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "If incoming or outgoing messages must be signed.",
								MarkdownDescription: "If incoming or outgoing messages must be signed.",
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
								Description:         "Validate the partner server certificate. Default is true.",
								MarkdownDescription: "Validate the partner server certificate. Default is true.",
								Default:             booldefault.StaticBool(true),
							},
						},
						Optional:            true,
						Description:         "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
						MarkdownDescription: "The SOAP authentication methods when sending or receiving a message using SOAP back channel.",
					},
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "The back channel authentication type.",
								MarkdownDescription: "The back channel authentication type.",
								Validators: []validator.String{
									stringvalidator.OneOf("INBOUND", "OUTBOUND"),
								},
							},
							"http_basic_credentials": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"username": schema.StringAttribute{
										Optional:            true,
										Description:         "The username.",
										MarkdownDescription: "The username.",
									},
									"password": schema.StringAttribute{
										Optional:            true,
										Sensitive:           false,
										Description:         "User password. To update the password, specify the plaintext value in this field.",
										MarkdownDescription: "User password. To update the password, specify the plaintext value in this field.",
									},
								},
								Optional:            true,
								Description:         "Username and password credentials.",
								MarkdownDescription: "Username and password credentials.",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								Description:         "If incoming or outgoing messages must be signed.",
								MarkdownDescription: "If incoming or outgoing messages must be signed.",
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
							"certs": connectioncert.ToSchema(
								"The certificates used for signature verification and XML encryption.",
								false,
								false,
							),
							"require_ssl": schema.BoolAttribute{
								Optional:            true,
								Description:         "Incoming HTTP transmissions must use a secure channel.",
								MarkdownDescription: "Incoming HTTP transmissions must use a secure channel.",
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
				Computed:            true,
				Description:         "Identifier that specifies the message displayed on a user-facing error page.",
				MarkdownDescription: "Identifier that specifies the message displayed on a user-facing error page.",
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
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The plugin instance name.",
											MarkdownDescription: "The plugin instance name.",
										},
										"plugin_descriptor_ref": schema.SingleNestedAttribute{
											Attributes:          resourcelink.ToSchema(),
											Required:            true,
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
											Attributes: map[string]schema.Attribute{
												"application_name": schema.StringAttribute{
													Optional:            true,
													Description:         "The application name.",
													MarkdownDescription: "The application name.",
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
								"issuance_criteria": schema.SingleNestedAttribute{
									Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
									// Computed:    true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.SetNestedAttribute{
											Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											Computed:    true,
											Optional:    true,
											Default:     setdefault.StaticValue(conditionalCriteriaDefault),
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": sourcetypeidkey.ToSchema(false),
													"attribute_name": schema.StringAttribute{
														Description: "The name of the attribute to use in this issuance criterion.",
														Required:    true,
													},
													"condition": schema.StringAttribute{
														Description: "The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.",
														Required:    true,
														Validators: []validator.String{
															stringvalidator.OneOf([]string{"EQUALS", "EQUALS_CASE_INSENSITIVE", "EQUALS_DN", "NOT_EQUAL", "NOT_EQUAL_CASE_INSENSITIVE", "NOT_EQUAL_DN", "MULTIVALUE_CONTAINS", "MULTIVALUE_CONTAINS_CASE_INSENSITIVE", "MULTIVALUE_CONTAINS_DN", "MULTIVALUE_DOES_NOT_CONTAIN", "MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE", "MULTIVALUE_DOES_NOT_CONTAIN_DN"}...),
														},
													},
													"value": schema.StringAttribute{
														Required:    true,
														Description: "The expected value of this issuance criterion.",
													},
													"error_result": schema.StringAttribute{
														Optional:    true,
														Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
												},
											},
										},
										"expression_criteria": schema.SetNestedAttribute{
											Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"expression": schema.StringAttribute{
														Required:    true,
														Description: "The OGNL expression to evaluate.",
													},
													"error_result": schema.StringAttribute{
														Optional:    true,
														Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
													},
												},
											},
										},
									},
								},
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:            true,
									Description:         "Restricts this mapping to specific virtual entity IDs.",
									MarkdownDescription: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "The list of virtual server IDs that this mapping is restricted to.",
									MarkdownDescription: "The list of virtual server IDs that this mapping is restricted to.",
								},
								"sp_adapter_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchema(),
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
						Computed:            true,
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						Description:         "Specify to always sign the SAML ArtifactResponse.",
						MarkdownDescription: "Specify to always sign the SAML ArtifactResponse.",
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
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            false,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Optional:            true,
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
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
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
					"authentication_policy_contract_mappings": schema.SetNestedAttribute{
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
								"authentication_policy_contract_ref": schema.SingleNestedAttribute{
									Attributes:          resourcelink.ToSchema(),
									Required:            true,
									Description:         "A reference to a resource.",
									MarkdownDescription: "A reference to a resource.",
								},
								"issuance_criteria": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"conditional_criteria": schema.SetNestedAttribute{
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
										"expression_criteria": schema.SetNestedAttribute{
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
								"restricted_virtual_server_ids": schema.SetAttribute{
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
					"authn_context_mappings": schema.SetNestedAttribute{
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
							"slo_subject_name_id_encrypted": schema.BoolAttribute{
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
					"enabled_profiles": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						MarkdownDescription: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
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
									"attribute_contract": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"masked": schema.BoolAttribute{
													Optional:            true,
													Computed:            true,
													Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
													MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
													Default:             booldefault.StaticBool(false),
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
															},
															"table_name": schema.StringAttribute{
																Required:            true,
																Description:         "The name of the database table.",
																MarkdownDescription: "The name of the database table.",
															},
															"unique_id_column": schema.StringAttribute{
																Required:            true,
																Description:         "The database column that uniquely identifies the provisioned user on the SP side.",
																MarkdownDescription: "The database column that uniquely identifies the provisioned user on the SP side.",
															},
														},
													},
													"stored_procedure": schema.SingleNestedAttribute{
														Description: "The Stored Procedure SQL method. The procedure is always called for all SSO tokens and \"eventTrigger\" will always be 'ALL_SAML_ASSERTIONS'.",
														Optional:    true,
														Attributes: map[string]schema.Attribute{
															"schema": schema.StringAttribute{
																Required:            true,
																Description:         "Lists the table structure that stores information within a database.",
																MarkdownDescription: "Lists the table structure that stores information within a database.",
															},
															"stored_procedure": schema.StringAttribute{
																Required:            true,
																Description:         "The name of the database stored procedure.",
																MarkdownDescription: "The name of the database stored procedure.",
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
											"base_dn": schema.StringAttribute{
												Optional:            true,
												Description:         "The base DN to search from. If not specified, the search will start at the LDAP's root.",
												MarkdownDescription: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
											},
											"unique_user_id_filter": schema.StringAttribute{
												Required:            true,
												Description:         "The expression that results in a unique user identifier, when combined with the Base DN.",
												MarkdownDescription: "The expression that results in a unique user identifier, when combined with the Base DN.",
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
						Attributes:          resourcelink.ToSchema(),
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
					"slo_service_endpoints": schema.SetNestedAttribute{
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
						Optional:            true,
						Description:         "A mapping in a connection that defines how access tokens are created.",
						MarkdownDescription: "A mapping in a connection that defines how access tokens are created.",
					},
					"idp_oauth_attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"masked": schema.BoolAttribute{
											Optional:            false,
											Computed:            true,
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
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
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
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
							"attributes": schema.SetNestedAttribute{
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
											"core_attributes": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"masked": schema.BoolAttribute{
															Optional:            false,
															Computed:            true,
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Optional:            true,
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
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of additional attributes that are added to the SCIM r.",
												MarkdownDescription: "A list of additional attributes that are added to the SCIM r.",
											},
										},
										Required:            true,
										Description:         "A set of user attributes that the IdP sends in the SCIM r.",
										MarkdownDescription: "A set of user attributes that the IdP sends in the SCIM r.",
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
									"attributes": schema.SetNestedAttribute{
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
						Optional:            true,
						Description:         "Group creation and read configuration.",
						MarkdownDescription: "Group creation and read configuration.",
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
									},
									"unique_user_id_filter": schema.StringAttribute{
										Required:            true,
										Description:         "The expression that results in a unique user identifier, when combined with the Base DN.",
										MarkdownDescription: "The expression that results in a unique user identifier, when combined with the Base DN.",
									},
									"unique_group_id_filter": schema.StringAttribute{
										Required:            true,
										Description:         "The expression that results in a unique group identifier, when combined with the Base DN.",
										MarkdownDescription: "The expression that results in a unique group identifier, when combined with the Base DN.",
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
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
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
															Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
															Default:             booldefault.StaticBool(false),
														},
														"name": schema.StringAttribute{
															Required:            true,
															Description:         "The name of this attribute.",
															MarkdownDescription: "The name of this attribute.",
														},
													},
												},
												Optional:            true,
												Description:         "A list of additional attributes that are added to the SCIM r.",
												MarkdownDescription: "A list of additional attributes that are added to the SCIM r.",
											},
										},
										Required:            true,
										Description:         "A set of user attributes that the IdP sends in the SCIM r.",
										MarkdownDescription: "A set of user attributes that the IdP sends in the SCIM r.",
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
									"attributes": schema.SetNestedAttribute{
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
				Computed:            true,
				Description:         "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				MarkdownDescription: "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
						Attributes:          resourcelink.ToSchema(),
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
				Computed:            true,
				Description:         "List of alternate entity IDs that identifies the local server to this partner.",
				MarkdownDescription: "List of alternate entity IDs that identifies the local server to this partner.",
				Default:             setdefault.StaticValue(emptyStringSet),
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
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
										},
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "The name of this attribute.",
											MarkdownDescription: "The name of this attribute.",
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
											Description:         "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											MarkdownDescription: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
											Default:             booldefault.StaticBool(false),
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
									Description:         "Indicates whether the token generator mapping is the default mapping. The default value is false.",
									MarkdownDescription: "Indicates whether the token generator mapping is the default mapping. The default value is false.",
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
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
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

func (r *spIdpConnectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Compare to version 12.0.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1200 := compare >= 0
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
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func addOptionalSpIdpConnectionFields(ctx context.Context, addRequest *client.IdpConnection, plan spIdpConnectionResourceModel) error {
	addRequest.ErrorPageMsgId = plan.ErrorPageMsgId.ValueStringPointer()
	addRequest.Id = plan.ConnectionId.ValueStringPointer()
	addRequest.Type = plan.Type.ValueStringPointer()
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.BaseUrl = plan.BaseUrl.ValueStringPointer()
	addRequest.DefaultVirtualEntityId = plan.DefaultVirtualEntityId.ValueStringPointer()

	if internaltypes.IsDefined(plan.LicenseConnectionGroup) {
		addRequest.LicenseConnectionGroup = plan.LicenseConnectionGroup.ValueStringPointer()
	}

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

		userRepository := plan.IdpBrowserSso.Attributes()["jit_provisioning"].(types.Object).Attributes()["user_repository"]
		if userRepository != nil {
			jdbcDataStoreRepository := plan.IdpBrowserSso.Attributes()["jit_provisioning"].(types.Object).Attributes()["user_repository"].(types.Object).Attributes()["jdbc"]
			ldapDataStoreRepository := plan.IdpBrowserSso.Attributes()["jit_provisioning"].(types.Object).Attributes()["user_repository"].(types.Object).Attributes()["ldap"]
			if jdbcDataStoreRepository != nil {
				addRequest.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository = &client.JdbcDataStoreRepository{}
				err := json.Unmarshal([]byte(internaljson.FromValue(jdbcDataStoreRepository, false)), addRequest.IdpBrowserSso.JitProvisioning.UserRepository.JdbcDataStoreRepository)
				if err != nil {
					return err
				}
			} else if ldapDataStoreRepository != nil {
				addRequest.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository = &client.LdapDataStoreRepository{}
				err := json.Unmarshal([]byte(internaljson.FromValue(ldapDataStoreRepository, false)), addRequest.IdpBrowserSso.JitProvisioning.UserRepository.LdapDataStoreRepository)
				if err != nil {
					return err
				}
			}
		}

		ssoOAuthMapping := plan.IdpBrowserSso.Attributes()["sso_oauth_mapping"]
		if ssoOAuthMapping != nil {
			addRequest.IdpBrowserSso.SsoOAuthMapping = &client.SsoOAuthMapping{}

			attributeSources := ssoOAuthMapping.(types.Object).Attributes()["attribute_sources"]
			if internaltypes.IsDefined(attributeSources) {
				attributeSourceClientStruct, attributeSourceClientStructErr := attributesources.ClientStruct(attributeSources.(types.Set))
				if attributeSourceClientStructErr != nil {
					return attributeSourceClientStructErr
				}

				addRequest.IdpBrowserSso.SsoOAuthMapping.AttributeSources = attributeSourceClientStruct
			}

			attributeContractFulfillment := ssoOAuthMapping.(types.Object).Attributes()["attribute_contract_fulfillment"]
			if internaltypes.IsDefined(attributeContractFulfillment) {
				attributeContractFulfillmentClientStruct, attributeContractFulfillmentClientStructErr := attributecontractfulfillment.ClientStruct(attributeContractFulfillment.(types.Map))
				if attributeContractFulfillmentClientStructErr != nil {
					return attributeContractFulfillmentClientStructErr
				}
				addRequest.IdpBrowserSso.SsoOAuthMapping.AttributeContractFulfillment = attributeContractFulfillmentClientStruct
			}

			issuanceCriteria := ssoOAuthMapping.(types.Object).Attributes()["issuance_criteria"]
			if internaltypes.IsDefined(issuanceCriteria) {
				issuanceCriteriaClientStruct, issuanceCriteriaClientStructErr := issuancecriteria.ClientStruct(issuanceCriteria.(types.Object))
				if issuanceCriteriaClientStructErr != nil {
					return issuanceCriteriaClientStructErr
				}
				addRequest.IdpBrowserSso.SsoOAuthMapping.IssuanceCriteria = issuanceCriteriaClientStruct
			}
		}

		err := json.Unmarshal([]byte(internaljson.FromValue(plan.IdpBrowserSso, true)), addRequest.IdpBrowserSso)
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
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.WsTrust, true)), addRequest.WsTrust)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.InboundProvisioning) {
		inboundProvisioningAttibutes := plan.InboundProvisioning.Attributes()
		addRequest.InboundProvisioning = &client.IdpInboundProvisioning{}

		// group support
		addRequest.InboundProvisioning.GroupSupport = inboundProvisioningAttibutes["group_support"].(types.Bool).ValueBool()

		// user repository
		userRepository := inboundProvisioningAttibutes["user_repository"]
		if userRepository != nil {
			identityStoreDataStoreRepository := userRepository.(types.Object).Attributes()["identity_store"]
			ldapDataStoreRepository := userRepository.(types.Object).Attributes()["ldap"]
			if identityStoreDataStoreRepository != nil {
				addRequest.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository = &client.IdentityStoreInboundProvisioningUserRepository{}
				addRequest.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository.Type = "IDENTITY_STORE"
				err := json.Unmarshal([]byte(internaljson.FromValue(identityStoreDataStoreRepository, true)), addRequest.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository)
				if err != nil {
					return err
				}
			} else if ldapDataStoreRepository != nil {
				addRequest.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository = &client.LdapInboundProvisioningUserRepository{}
				addRequest.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository.Type = "LDAP"
				err := json.Unmarshal([]byte(internaljson.FromValue(ldapDataStoreRepository, true)), addRequest.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository)
				if err != nil {
					return err
				}
			}
		}

		// custom schema
		customSchema := inboundProvisioningAttibutes["custom_schema"]
		addRequest.InboundProvisioning.CustomSchema = client.Schema{}
		err := json.Unmarshal([]byte(internaljson.FromValue(customSchema, true)), &addRequest.InboundProvisioning.CustomSchema)
		if err != nil {
			return err
		}

		// users
		if internaltypes.IsDefined(inboundProvisioningAttibutes["users"]) {
			addRequest.InboundProvisioning.Users = client.Users{}
			err := json.Unmarshal([]byte(internaljson.FromValue(inboundProvisioningAttibutes["users"], true)), &addRequest.InboundProvisioning.Users)
			if err != nil {
				return err
			}
		}

		// groups
		if internaltypes.IsDefined(inboundProvisioningAttibutes["groups"]) {
			addRequest.InboundProvisioning.Groups = &client.Groups{}
			err := json.Unmarshal([]byte(internaljson.FromValue(inboundProvisioningAttibutes["groups"], true)), &addRequest.InboundProvisioning.Groups)
			if err != nil {
				return err
			}
		}

		// action on delete
		if internaltypes.IsDefined(inboundProvisioningAttibutes["action_on_delete"]) {
			addRequest.InboundProvisioning.ActionOnDelete = inboundProvisioningAttibutes["action_on_delete"].(types.String).ValueStringPointer()
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

	state.Active = types.BoolPointerValue(r.Active)
	state.AdditionalAllowedEntitiesConfiguration, objDiags = types.ObjectValueFrom(ctx, additionalAllowedEntitiesConfigurationAttrTypes, r.AdditionalAllowedEntitiesConfiguration)
	diags.Append(objDiags...)
	state.AttributeQuery, objDiags = types.ObjectValueFrom(ctx, attributeQueryAttrTypes, r.AttributeQuery)
	state.BaseUrl = types.StringPointerValue(r.BaseUrl)
	diags.Append(objDiags...)
	state.ConnectionId = types.StringPointerValue(r.Id)
	state.ContactInfo, objDiags = types.ObjectValueFrom(ctx, contactInfoAttrTypes, r.ContactInfo)
	diags.Append(objDiags...)
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	state.EntityId = types.StringValue(r.EntityId)
	state.ErrorPageMsgId = types.StringPointerValue(r.ErrorPageMsgId)
	state.ExtendedProperties, objDiags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: extendedPropertiesElemAttrTypes}, r.ExtendedProperties)
	diags.Append(objDiags...)
	state.Id = types.StringPointerValue(r.Id)
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	state.MetadataReloadSettings, objDiags = types.ObjectValueFrom(ctx, metadataReloadSettingsAttrTypes, r.MetadataReloadSettings)
	diags.Append(objDiags...)
	state.Name = types.StringValue(r.Name)
	state.OidcClientCredentials = types.ObjectNull(oidcClientCredentialsAttrTypes)
	state.Type = types.StringPointerValue(r.Type)
	state.VirtualEntityIds = internaltypes.GetStringSet(r.VirtualEntityIds)

	// LicenseConnectionGroup
	if r.LicenseConnectionGroup != nil {
		state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	}
	// Credentials
	var credentialsValue types.Object
	if r.Credentials != nil {
		var credentialsCertsValues []attr.Value
		for _, cert := range r.Credentials.Certs {
			for _, certInPlan := range plan.Credentials.Attributes()["certs"].(types.List).Elements() {
				x509FilePlanAttrs := certInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
				x509FileIdPlan := x509FilePlanAttrs["id"].(types.String).ValueString()
				if *cert.X509File.Id == x509FileIdPlan {
					planFileData := x509FilePlanAttrs["file_data"].(types.String)
					credentialsCertsObjValue, objDiags := connectioncert.ToState(ctx, planFileData, cert, &diags)
					diags.Append(objDiags...)
					credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
				}
			}
		}
		credentialsCertsValue, objDiags := types.ListValue(connectioncert.ObjType(), credentialsCertsValues)
		diags.Append(objDiags...)
		var credentialsDecryptionKeyPairRefValue types.Object
		if r.Credentials.DecryptionKeyPairRef == nil {
			credentialsDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.DecryptionKeyPairRef)
			diags.Append(objDiags...)
		}
		var credentialsInboundBackChannelAuthValue types.Object
		if r.Credentials.InboundBackChannelAuth == nil {
			credentialsInboundBackChannelAuthValue = types.ObjectNull(credentialsInboundBackChannelAuthAttrTypes)
		} else {
			var credentialsInboundBackChannelAuthCertsValue types.List
			if r.Credentials.InboundBackChannelAuth.Certs != nil && len(r.Credentials.InboundBackChannelAuth.Certs) > 0 {
				var credentialsInboundBackChannelAuthCertsValues []attr.Value
				for _, ibcaCert := range r.Credentials.InboundBackChannelAuth.Certs {
					for _, ibcaCertInPlan := range plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["certs"].(types.List).Elements() {
						ibcax509FilePlanAttrs := ibcaCertInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
						ibcax509FileIdPlan := ibcax509FilePlanAttrs["id"].(types.String).ValueString()
						if *ibcaCert.X509File.Id == ibcax509FileIdPlan {
							planIbcaX509FileFileData := ibcax509FilePlanAttrs["file_data"].(types.String)
							planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(ctx, planIbcaX509FileFileData, ibcaCert, &diags)
							diags.Append(objDiags...)
							credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
						}
					}
				}
				credentialsInboundBackChannelAuthCertsValue, objDiags = types.ListValue(connectioncert.ObjType(), credentialsInboundBackChannelAuthCertsValues)
				diags.Append(objDiags...)
			} else {
				credentialsInboundBackChannelAuthCertsValue = types.ListNull(connectioncert.ObjType())
			}
			var credentialsInboundBackChannelAuthHttpBasicCredentialsValue types.Object
			if r.Credentials.InboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password string = ""
				if plan != nil {
					passwordFromPlan := plan.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					if internaltypes.IsDefined(passwordFromPlan) {
						password = passwordFromPlan.ValueString()
					} else if state != nil {
						password = state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueString()
					}
				} else if state != nil && internaltypes.IsDefined(state.Credentials) {
					password = state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueString()
				}
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password": types.StringValue(password),
					"username": types.StringPointerValue(r.Credentials.InboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				diags.Append(objDiags...)
			}
			credentialsInboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"certs":                   credentialsInboundBackChannelAuthCertsValue,
				"digital_signature":       types.BoolPointerValue(r.Credentials.InboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials":  credentialsInboundBackChannelAuthHttpBasicCredentialsValue,
				"require_ssl":             types.BoolPointerValue(r.Credentials.InboundBackChannelAuth.RequireSsl),
				"type":                    types.StringValue(r.Credentials.InboundBackChannelAuth.Type),
				"verification_issuer_dn":  types.StringPointerValue(r.Credentials.InboundBackChannelAuth.VerificationIssuerDN),
				"verification_subject_dn": types.StringPointerValue(r.Credentials.InboundBackChannelAuth.VerificationSubjectDN),
			})
			diags.Append(objDiags...)
		}
		var credentialsOutboundBackChannelAuthValue types.Object
		if r.Credentials.OutboundBackChannelAuth == nil {
			credentialsOutboundBackChannelAuthValue = types.ObjectNull(credentialsOutboundBackChannelAuthAttrTypes)
		} else {
			var credentialsOutboundBackChannelAuthHttpBasicCredentialsValue types.Object
			if r.Credentials.OutboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password string = ""
				if plan != nil {
					passwordFromPlan := plan.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					if internaltypes.IsDefined(passwordFromPlan) {
						password = passwordFromPlan.ValueString()
					} else if state != nil {
						password = state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueString()
					}
				} else if state != nil && internaltypes.IsDefined(state.Credentials) {
					password = state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String).ValueString()
				}
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password": types.StringPointerValue(&password),
					"username": types.StringPointerValue(r.Credentials.OutboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				diags.Append(objDiags...)
			}
			var credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue types.Object
			if r.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef == nil {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
			} else {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef)
				diags.Append(objDiags...)
			}
			credentialsOutboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"digital_signature":      types.BoolPointerValue(r.Credentials.OutboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials": credentialsOutboundBackChannelAuthHttpBasicCredentialsValue,
				"ssl_auth_key_pair_ref":  credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue,
				"type":                   types.StringValue(r.Credentials.OutboundBackChannelAuth.Type),
				"validate_partner_cert":  types.BoolPointerValue(r.Credentials.OutboundBackChannelAuth.ValidatePartnerCert),
			})
			diags.Append(objDiags...)
		}
		var credentialsSecondaryDecryptionKeyPairRefValue types.Object
		if r.Credentials.SecondaryDecryptionKeyPairRef == nil {
			credentialsSecondaryDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsSecondaryDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(ctx, r.Credentials.SecondaryDecryptionKeyPairRef)
			diags.Append(objDiags...)
		}
		var credentialsSigningSettingsValue types.Object
		if r.Credentials.SigningSettings == nil {
			credentialsSigningSettingsValue = types.ObjectNull(credentialsSigningSettingsAttrTypes)
		} else {
			var credentialsSigningSettingsAlternativeSigningKeyPairRefsValues []attr.Value
			for _, credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue := range r.Credentials.SigningSettings.AlternativeSigningKeyPairRefs {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue, objDiags := resourcelink.ToState(ctx, &credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue)
				diags.Append(objDiags...)
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValues = append(credentialsSigningSettingsAlternativeSigningKeyPairRefsValues, credentialsSigningSettingsAlternativeSigningKeyPairRefsValue)
			}
			var credentialsSigningSettingsAlternativeSigningKeyPairRefsValue types.Set
			if len(credentialsSigningSettingsAlternativeSigningKeyPairRefsValues) > 0 {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue, objDiags = types.SetValue(credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType, credentialsSigningSettingsAlternativeSigningKeyPairRefsValues)
				diags.Append(objDiags...)
			} else {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue = types.SetNull(credentialsSigningSettingsAlternativeSigningKeyPairRefsElementType)
			}

			credentialsSigningSettingsSigningKeyPairRefValue, objDiags := resourcelink.ToState(ctx, &r.Credentials.SigningSettings.SigningKeyPairRef)
			diags.Append(objDiags...)
			credentialsSigningSettingsValue, objDiags = types.ObjectValue(credentialsSigningSettingsAttrTypes, map[string]attr.Value{
				"algorithm":                         types.StringPointerValue(r.Credentials.SigningSettings.Algorithm),
				"alternative_signing_key_pair_refs": credentialsSigningSettingsAlternativeSigningKeyPairRefsValue,
				"include_cert_in_signature":         types.BoolPointerValue(r.Credentials.SigningSettings.IncludeCertInSignature),
				"include_raw_key_in_signature":      types.BoolPointerValue(r.Credentials.SigningSettings.IncludeRawKeyInSignature),
				"signing_key_pair_ref":              credentialsSigningSettingsSigningKeyPairRefValue,
			})
			diags.Append(objDiags...)
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
		diags.Append(objDiags...)
	} else {
		credentialsValue = types.ObjectNull(credentialsAttrTypes)
	}
	state.Credentials = credentialsValue

	// OidcClientCredentials
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

	// IdpBrowserSso
	var idpBrowserSsoValue types.Object
	if r.IdpBrowserSso == nil {
		idpBrowserSsoValue = types.ObjectNull(idpBrowserSsoAttrTypes)
	} else {
		var idpBrowserSsoAdapterMappingsValues []attr.Value
		var idpBrowserSsoAdapterMappingsValue types.List
		if r.IdpBrowserSso.AdapterMappings != nil {
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
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
								"name": types.StringValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Name),
							})
							diags.Append(objDiags...)
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues = append(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue)
						}
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.ListValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues)
						diags.Append(objDiags...)
						var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues []attr.Value
						for _, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue := range idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.ExtendedAttributes {
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
								"name": types.StringValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Name),
							})
							diags.Append(objDiags...)
							idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues = append(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue)
						}
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.ListValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType, idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues)
						diags.Append(objDiags...)
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue, diags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes, map[string]attr.Value{
							"core_attributes":     idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue,
							"extended_attributes": idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue,
						})
						diags.Append(objDiags...)
					}
					adapterMapping := plan.IdpBrowserSso.Attributes()["adapter_mappings"].(types.Set).Elements()[i].(types.Object).Attributes()
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue, diags := pluginconfiguration.ToState(adapterMapping["adapter_override_settings"].(types.Object).Attributes()["configuration"].(types.Object), &idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Configuration)
					diags.Append(objDiags...)
					var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue types.Object
					if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef == nil {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue = types.ObjectNull(resourcelink.AttrType())
					} else {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue, diags = types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
							"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef.Id),
						})
						diags.Append(objDiags...)
					}
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue, diags := types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
						"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.PluginDescriptorRef.Id),
					})
					diags.Append(objDiags...)
					var idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue types.Object
					if idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo == nil {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue = types.ObjectNull(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes)
					} else {
						idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue, diags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoAttrTypes, map[string]attr.Value{
							"application_icon_url": types.StringPointerValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo.ApplicationIconUrl),
							"application_name":     types.StringPointerValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.TargetApplicationInfo.ApplicationName),
						})
						diags.Append(objDiags...)
					}
					idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue, diags = types.ObjectValue(idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes, map[string]attr.Value{
						"attribute_contract":      idpBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue,
						"configuration":           idpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue,
						"id":                      types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Id),
						"name":                    types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Name),
						"parent_ref":              idpBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue,
						"plugin_descriptor_ref":   idpBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue,
						"target_application_info": idpBrowserSsoAdapterMappingsAdapterOverrideSettingsTargetApplicationInfoValue,
					})
					diags.Append(objDiags...)
				}
				idpBrowserSsoAdapterMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(ctx, &idpBrowserSsoAdapterMappingsResponseValue.AttributeContractFulfillment)
				diags.Append(objDiags...)

				idpBrowserSsoAdapterMappingsAttributeSourcesValue, diags := attributesources.ToState(ctx, idpBrowserSsoAdapterMappingsResponseValue.AttributeSources)
				diags.Append(objDiags...)

				var idpBrowserSsoAdapterMappingsIssuanceCriteriaValue types.Object
				if idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria != nil && (len(idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria.ConditionalCriteria) > 0 || len(idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria.ExpressionCriteria) > 0) {
					idpBrowserSsoAdapterMappingsIssuanceCriteriaValue, diags = issuancecriteria.ToState(ctx, idpBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria)
					diags.Append(objDiags...)
				} else {
					idpBrowserSsoAdapterMappingsIssuanceCriteriaValue = types.ObjectNull(issuancecriteria.AttrTypes())
				}

				idpBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(ctx, types.StringType, idpBrowserSsoAdapterMappingsResponseValue.RestrictedVirtualEntityIds)
				diags.Append(objDiags...)
				idpBrowserSsoAdapterMappingsSpAdapterRefValue, diags := types.ObjectValue(resourcelink.AttrType(), map[string]attr.Value{
					"id": types.StringValue(idpBrowserSsoAdapterMappingsResponseValue.SpAdapterRef.Id),
				})
				diags.Append(objDiags...)
				idpBrowserSsoAdapterMappingsValue, diags := types.ObjectValue(idpBrowserSsoAdapterMappingsAttrTypes, map[string]attr.Value{
					"adapter_override_settings":      idpBrowserSsoAdapterMappingsAdapterOverrideSettingsValue,
					"attribute_contract_fulfillment": idpBrowserSsoAdapterMappingsAttributeContractFulfillmentValue,
					"attribute_sources":              idpBrowserSsoAdapterMappingsAttributeSourcesValue,
					"issuance_criteria":              idpBrowserSsoAdapterMappingsIssuanceCriteriaValue,
					"restrict_virtual_entity_ids":    types.BoolPointerValue(idpBrowserSsoAdapterMappingsResponseValue.RestrictVirtualEntityIds),
					"restricted_virtual_entity_ids":  idpBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue,
					"sp_adapter_ref":                 idpBrowserSsoAdapterMappingsSpAdapterRefValue,
				})
				diags.Append(objDiags...)
				idpBrowserSsoAdapterMappingsValues = append(idpBrowserSsoAdapterMappingsValues, idpBrowserSsoAdapterMappingsValue)
			}
			idpBrowserSsoAdapterMappingsValue, diags = types.ListValue(idpBrowserSsoAdapterMappingsElementType, idpBrowserSsoAdapterMappingsValues)
			diags.Append(objDiags...)
		} else {
			idpBrowserSsoAdapterMappingsValue = types.ListNull(idpBrowserSsoAdapterMappingsElementType)
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
		diags.Append(objDiags...)

		// IdpBrowserSSO Attribute Contract
		idpBrowserSsoAttributeContractValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoAttributeContractAttrTypes, r.IdpBrowserSso.AttributeContract)
		diags.Append(objDiags...)

		// IdpBrowserSSO Authentication Policy Contract Mappings
		idpBrowserSsoAuthenticationPolicyContractMappingsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoAuthenticationPolicyContractMappingsElementType, r.IdpBrowserSso.AuthenticationPolicyContractMappings)
		diags.Append(objDiags...)

		// IdpBrowserSSO AuthnContextMappings
		idpBrowserSsoAuthnContextMappingsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoAuthnContextMappingsElementType, r.IdpBrowserSso.AuthnContextMappings)
		diags.Append(objDiags...)

		// IdpBrowserSSO Decryption Policy
		idpBrowserSsoDecryptionPolicyValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoDecryptionPolicyAttrTypes, r.IdpBrowserSso.DecryptionPolicy)
		diags.Append(objDiags...)

		// IdpBrowserSSO Enabled Profiles
		idpBrowserSsoEnabledProfilesValue, objDiags := types.SetValueFrom(ctx, types.StringType, r.IdpBrowserSso.EnabledProfiles)
		diags.Append(objDiags...)

		// IdpBrowserSSO Incoming Bindings
		idpBrowserSsoIncomingBindingsValue, objDiags := types.SetValueFrom(ctx, types.StringType, r.IdpBrowserSso.IncomingBindings)
		diags.Append(objDiags...)

		// IdpBrowserSSO JIT Provisioning
		var idpBrowserSsoJitProvisioningValue types.Object
		if r.IdpBrowserSso.JitProvisioning == nil {
			idpBrowserSsoJitProvisioningValue = types.ObjectNull(idpBrowserSsoJitProvisioningAttrTypes)
		} else {
			idpBrowserSsoJitProvisioningValue, objDiags = types.ObjectValueFrom(ctx, idpBrowserSsoJitProvisioningAttrTypes, r.IdpBrowserSso.JitProvisioning)
			diags.Append(objDiags...)
		}

		// IdpBrowserSSO Message Customizations
		idpBrowserSsoMessageCustomizationsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoMessageCustomizationsElementType, r.IdpBrowserSso.MessageCustomizations)
		diags.Append(objDiags...)

		// IdpBrowserSSO OAuth Authentication Policy Contract Ref
		idpBrowserSsoOauthAuthenticationPolicyContractRefValue, objDiags := resourcelink.ToState(ctx, r.IdpBrowserSso.OauthAuthenticationPolicyContractRef)
		diags.Append(objDiags...)

		// IdpBrowserSSO OIDC Provider Settings
		idpBrowserSsoOidcProviderSettingsValue, objDiags := types.ObjectValueFrom(ctx, idpBrowserSsoOidcProviderSettingsAttrTypes, r.IdpBrowserSso.OidcProviderSettings)
		diags.Append(objDiags...)

		// IdpBrowserSSO SLO Service Endpoints
		idpBrowserSsoSloServiceEndpointsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoSloServiceEndpointsElementType, r.IdpBrowserSso.SloServiceEndpoints)
		diags.Append(objDiags...)

		// IdpBrowserSSO SSO Service Endpoints
		idpBrowserSsoSsoServiceEndpointsValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoSsoServiceEndpointsElementType, r.IdpBrowserSso.SsoServiceEndpoints)
		diags.Append(objDiags...)

		// IdpBrowserSSO URL Whitelist Entries
		idpBrowserSsoUrlWhitelistEntriesValue, objDiags := types.SetValueFrom(ctx, idpBrowserSsoUrlWhitelistEntriesElementType, r.IdpBrowserSso.UrlWhitelistEntries)
		diags.Append(objDiags...)

		// IdpBrowserSSO SSO OAuth Mapping
		var idpBrowserSsoSsoOauthMappingValue types.Object
		if r.IdpBrowserSso.SsoOAuthMapping == nil {
			idpBrowserSsoSsoOauthMappingValue = types.ObjectNull(idpBrowserSsoSsoOauthMappingAttrTypes)
		} else {
			idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentValue, objDiags := attributecontractfulfillment.ToState(ctx, &r.IdpBrowserSso.SsoOAuthMapping.AttributeContractFulfillment)
			diags.Append(objDiags...)

			var idpBrowserSsoSsoOauthMappingAttributeSourcesValue types.Set
			if r.IdpBrowserSso.SsoOAuthMapping.AttributeSources != nil {
				idpBrowserSsoSsoOauthMappingAttributeSourcesValue, objDiags = attributesources.ToState(ctx, r.IdpBrowserSso.SsoOAuthMapping.AttributeSources)
				diags.Append(objDiags...)
			}

			var idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue types.Object
			if r.IdpBrowserSso.SsoOAuthMapping.IssuanceCriteria != nil {
				idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue, objDiags = issuancecriteria.ToState(ctx, r.IdpBrowserSso.SsoOAuthMapping.IssuanceCriteria)
				diags.Append(objDiags...)
			}

			idpBrowserSsoSsoOauthMappingValue, objDiags = types.ObjectValue(idpBrowserSsoSsoOauthMappingAttrTypes, map[string]attr.Value{
				"attribute_contract_fulfillment": idpBrowserSsoSsoOauthMappingAttributeContractFulfillmentValue,
				"attribute_sources":              idpBrowserSsoSsoOauthMappingAttributeSourcesValue,
				"issuance_criteria":              idpBrowserSsoSsoOauthMappingIssuanceCriteriaValue,
			})
			diags.Append(objDiags...)
		}

		idpBrowserSsoValue, diags = types.ObjectValue(idpBrowserSsoAttrTypes, map[string]attr.Value{
			"adapter_mappings":                         idpBrowserSsoAdapterMappingsValue,
			"always_sign_artifact_response":            idpBrowserSsoAlwaysSignArtifactResponse,
			"artifact":                                 idpBrowserSsoArtifactValue,
			"assertions_signed":                        types.BoolPointerValue(r.IdpBrowserSso.AssertionsSigned),
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
			"sign_authn_requests":                      types.BoolPointerValue(r.IdpBrowserSso.SignAuthnRequests),
			"slo_service_endpoints":                    idpBrowserSsoSloServiceEndpointsValue,
			"sso_application_endpoint":                 types.StringPointerValue(r.IdpBrowserSso.SsoApplicationEndpoint),
			"sso_oauth_mapping":                        idpBrowserSsoSsoOauthMappingValue,
			"sso_service_endpoints":                    idpBrowserSsoSsoServiceEndpointsValue,
			"url_whitelist_entries":                    idpBrowserSsoUrlWhitelistEntriesValue,
		})
		diags.Append(objDiags...)
	}
	state.IdpBrowserSso = idpBrowserSsoValue

	// IdpOAuthGrantAttributeMapping
	state.IdpOAuthGrantAttributeMapping, objDiags = types.ObjectValueFrom(ctx, idpOAuthGrantAttributeMappingAttrTypes, r.IdpOAuthGrantAttributeMapping)
	diags.Append(objDiags...)
	var inboundProvisioningValue types.Object

	// InboundProvisioning
	var inboundProvisioningGroupsValue types.Object
	if r.InboundProvisioning == nil {
		inboundProvisioningValue = types.ObjectNull(inboundProvisioningAttrTypes)
	} else {
		var inboundProvisioningCustomSchemaAttributesValues []attr.Value
		for _, inboundProvisioningCustomSchemaAttributesResponseValue := range r.InboundProvisioning.CustomSchema.Attributes {
			inboundProvisioningCustomSchemaAttributesSubAttributesValue, objDiags := types.SetValueFrom(ctx, types.StringType, inboundProvisioningCustomSchemaAttributesResponseValue.SubAttributes)
			diags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesTypesValue, objDiags := types.SetValueFrom(ctx, types.StringType, inboundProvisioningCustomSchemaAttributesResponseValue.Types)
			diags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesValue, objDiags := types.ObjectValue(inboundProvisioningCustomSchemaAttributesAttrTypes, map[string]attr.Value{
				"multi_valued":   types.BoolPointerValue(inboundProvisioningCustomSchemaAttributesResponseValue.MultiValued),
				"name":           types.StringPointerValue(inboundProvisioningCustomSchemaAttributesResponseValue.Name),
				"sub_attributes": inboundProvisioningCustomSchemaAttributesSubAttributesValue,
				"types":          inboundProvisioningCustomSchemaAttributesTypesValue,
			})
			diags.Append(objDiags...)
			inboundProvisioningCustomSchemaAttributesValues = append(inboundProvisioningCustomSchemaAttributesValues, inboundProvisioningCustomSchemaAttributesValue)
		}
		inboundProvisioningCustomSchemaAttributesValue, objDiags := types.SetValue(inboundProvisioningCustomSchemaAttributesElementType, inboundProvisioningCustomSchemaAttributesValues)
		diags.Append(objDiags...)
		inboundProvisioningCustomSchemaValue, objDiags := types.ObjectValue(inboundProvisioningCustomSchemaAttrTypes, map[string]attr.Value{
			"attributes": inboundProvisioningCustomSchemaAttributesValue,
			"namespace":  types.StringPointerValue(r.InboundProvisioning.CustomSchema.Namespace),
		})
		diags.Append(objDiags...)
		if r.InboundProvisioning.Groups != nil {
			var inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeContract.CoreAttributes {
				inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue.Masked),
					"name":   types.StringValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesResponseValue.Name),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesElementType, inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValues)
			diags.Append(objDiags...)
			var inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeContract.ExtendedAttributes {
				inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
					"masked": types.BoolPointerValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue.Masked),
					"name":   types.StringValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesResponseValue.Name),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesElementType, inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValues)
			diags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsAttributeContractValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeContractAttrTypes, map[string]attr.Value{
				"core_attributes":     inboundProvisioningGroupsReadGroupsAttributeContractCoreAttributesValue,
				"extended_attributes": inboundProvisioningGroupsReadGroupsAttributeContractExtendedAttributesValue,
			})
			diags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues := make(map[string]attr.Value)
			for key, inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue := range r.InboundProvisioning.Groups.ReadGroups.AttributeFulfillment {
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
					"id":   types.StringPointerValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Source.Id),
					"type": types.StringValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Source.Type),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentAttrTypes, map[string]attr.Value{
					"source": inboundProvisioningGroupsReadGroupsAttributeFulfillmentSourceValue,
					"value":  types.StringValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentResponseValue.Value),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues[key] = inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue
			}
			inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningGroupsReadGroupsAttributeFulfillmentElementType, inboundProvisioningGroupsReadGroupsAttributeFulfillmentValues)
			diags.Append(objDiags...)
			var inboundProvisioningGroupsReadGroupsAttributesValues []attr.Value
			for _, inboundProvisioningGroupsReadGroupsAttributesResponseValue := range r.InboundProvisioning.Groups.ReadGroups.Attributes {
				inboundProvisioningGroupsReadGroupsAttributesValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttributesAttrTypes, map[string]attr.Value{
					"name": types.StringValue(inboundProvisioningGroupsReadGroupsAttributesResponseValue.Name),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsReadGroupsAttributesValues = append(inboundProvisioningGroupsReadGroupsAttributesValues, inboundProvisioningGroupsReadGroupsAttributesValue)
			}
			inboundProvisioningGroupsReadGroupsAttributesValue, objDiags := types.SetValue(inboundProvisioningGroupsReadGroupsAttributesElementType, inboundProvisioningGroupsReadGroupsAttributesValues)
			diags.Append(objDiags...)
			inboundProvisioningGroupsReadGroupsValue, objDiags := types.ObjectValue(inboundProvisioningGroupsReadGroupsAttrTypes, map[string]attr.Value{
				"attribute_contract":    inboundProvisioningGroupsReadGroupsAttributeContractValue,
				"attribute_fulfillment": inboundProvisioningGroupsReadGroupsAttributeFulfillmentValue,
				"attributes":            inboundProvisioningGroupsReadGroupsAttributesValue,
			})
			diags.Append(objDiags...)
			inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues := make(map[string]attr.Value)
			for key, inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue := range r.InboundProvisioning.Groups.WriteGroups.AttributeFulfillment {
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
					"id":   types.StringPointerValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Source.Id),
					"type": types.StringValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Source.Type),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentAttrTypes, map[string]attr.Value{
					"source": inboundProvisioningGroupsWriteGroupsAttributeFulfillmentSourceValue,
					"value":  types.StringValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentResponseValue.Value),
				})
				diags.Append(objDiags...)
				inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues[key] = inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue
			}
			inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningGroupsWriteGroupsAttributeFulfillmentElementType, inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValues)
			diags.Append(objDiags...)
			inboundProvisioningGroupsWriteGroupsValue, objDiags := types.ObjectValue(inboundProvisioningGroupsWriteGroupsAttrTypes, map[string]attr.Value{
				"attribute_fulfillment": inboundProvisioningGroupsWriteGroupsAttributeFulfillmentValue,
			})
			diags.Append(objDiags...)
			inboundProvisioningGroupsValue, objDiags = types.ObjectValue(inboundProvisioningGroupsAttrTypes, map[string]attr.Value{
				"read_groups":  inboundProvisioningGroupsReadGroupsValue,
				"write_groups": inboundProvisioningGroupsWriteGroupsValue,
			})
			diags.Append(objDiags...)
		} else {
			inboundProvisioningGroupsValue = types.ObjectNull(inboundProvisioningGroupsAttrTypes)
		}

		var identityStoreInboundProvisioningUserRepository, ldapInboundProvisioningUserRepository types.Object
		if r.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository != nil {
			identityStoreProvisionerRef, objDiags := resourcelink.ToState(ctx, &r.InboundProvisioning.UserRepository.IdentityStoreInboundProvisioningUserRepository.IdentityStoreProvisionerRef)
			diags.Append(objDiags...)
			identityStoreInboundProvisioningUserRepository, objDiags = types.ObjectValue(inboundprovisioninguserrepository.IdentityStoreInboundProvisioningUserRepositoryAttrType(), map[string]attr.Value{
				"identity_store_provisioner_ref": identityStoreProvisionerRef,
			})
			diags.Append(objDiags...)

			ldapInboundProvisioningUserRepository = types.ObjectNull(inboundprovisioninguserrepository.LdapInboundProvisioningUserRepositoryAttrType())
		} else if r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository != nil {
			dataStoreRef, objDiags := resourcelink.ToState(ctx, &r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.DataStoreRef)
			diags.Append(objDiags...)
			ldapInboundProvisioningUserRepository, objDiags = types.ObjectValue(inboundprovisioninguserrepository.LdapInboundProvisioningUserRepositoryAttrType(), map[string]attr.Value{
				"base_dn":                types.StringPointerValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.BaseDn),
				"data_store_ref":         dataStoreRef,
				"unique_user_id_filter":  types.StringValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.UniqueUserIdFilter),
				"unique_group_id_filter": types.StringValue(r.InboundProvisioning.UserRepository.LdapInboundProvisioningUserRepository.UniqueGroupIdFilter),
			})
			diags.Append(objDiags...)
			identityStoreInboundProvisioningUserRepository = types.ObjectNull(inboundprovisioninguserrepository.IdentityStoreInboundProvisioningUserRepositoryAttrType())
		}

		inboundProvisioningUserRepositoryAttrValue := map[string]attr.Value{
			"identity_store": identityStoreInboundProvisioningUserRepository,
			"ldap":           ldapInboundProvisioningUserRepository,
		}

		inboundProvisioningUserRepositoryValue, objDiags := types.ObjectValue(inboundprovisioninguserrepository.ElemAttrType(), inboundProvisioningUserRepositoryAttrValue)
		diags.Append(objDiags...)

		var inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeContract.CoreAttributes {
			inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue.Masked),
				"name":   types.StringValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesResponseValue.Name),
			})
			diags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues = append(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributeContractCoreAttributesElementType, inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValues)
		diags.Append(objDiags...)
		var inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeContract.ExtendedAttributes {
			inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue.Masked),
				"name":   types.StringValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesResponseValue.Name),
			})
			diags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues = append(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesElementType, inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValues)
		diags.Append(objDiags...)
		inboundProvisioningUsersReadUsersAttributeContractValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     inboundProvisioningUsersReadUsersAttributeContractCoreAttributesValue,
			"extended_attributes": inboundProvisioningUsersReadUsersAttributeContractExtendedAttributesValue,
		})
		diags.Append(objDiags...)
		inboundProvisioningUsersReadUsersAttributeFulfillmentValues := make(map[string]attr.Value)
		for key, inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue := range r.InboundProvisioning.Users.ReadUsers.AttributeFulfillment {
			inboundProvisioningUsersReadUsersAttributeFulfillmentSourceValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeFulfillmentSourceAttrTypes, map[string]attr.Value{
				"id":   types.StringPointerValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Source.Id),
				"type": types.StringValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Source.Type),
			})
			diags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeFulfillmentValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributeFulfillmentAttrTypes, map[string]attr.Value{
				"source": inboundProvisioningUsersReadUsersAttributeFulfillmentSourceValue,
				"value":  types.StringValue(inboundProvisioningUsersReadUsersAttributeFulfillmentResponseValue.Value),
			})
			diags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributeFulfillmentValues[key] = inboundProvisioningUsersReadUsersAttributeFulfillmentValue
		}
		inboundProvisioningUsersReadUsersAttributeFulfillmentValue, objDiags := types.MapValue(inboundProvisioningUsersReadUsersAttributeFulfillmentElementType, inboundProvisioningUsersReadUsersAttributeFulfillmentValues)
		diags.Append(objDiags...)
		var inboundProvisioningUsersReadUsersAttributesValues []attr.Value
		for _, inboundProvisioningUsersReadUsersAttributesResponseValue := range r.InboundProvisioning.Users.ReadUsers.Attributes {
			inboundProvisioningUsersReadUsersAttributesValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(inboundProvisioningUsersReadUsersAttributesResponseValue.Name),
			})
			diags.Append(objDiags...)
			inboundProvisioningUsersReadUsersAttributesValues = append(inboundProvisioningUsersReadUsersAttributesValues, inboundProvisioningUsersReadUsersAttributesValue)
		}
		inboundProvisioningUsersReadUsersAttributesValue, objDiags := types.SetValue(inboundProvisioningUsersReadUsersAttributesElementType, inboundProvisioningUsersReadUsersAttributesValues)
		diags.Append(objDiags...)
		inboundProvisioningUsersReadUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersReadUsersAttrTypes, map[string]attr.Value{
			"attribute_contract":    inboundProvisioningUsersReadUsersAttributeContractValue,
			"attribute_fulfillment": inboundProvisioningUsersReadUsersAttributeFulfillmentValue,
			"attributes":            inboundProvisioningUsersReadUsersAttributesValue,
		})
		diags.Append(objDiags...)

		inboundProvisioningUsersWriteUsersAttributeFulfillmentValue, objDiags := attributecontractfulfillment.ToState(ctx, &r.InboundProvisioning.Users.WriteUsers.AttributeFulfillment)
		diags.Append(objDiags...)
		inboundProvisioningUsersWriteUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersWriteUsersAttrTypes, map[string]attr.Value{
			"attribute_fulfillment": inboundProvisioningUsersWriteUsersAttributeFulfillmentValue,
		})
		diags.Append(objDiags...)
		inboundProvisioningUsersValue, objDiags := types.ObjectValue(inboundProvisioningUsersAttrTypes, map[string]attr.Value{
			"read_users":  inboundProvisioningUsersReadUsersValue,
			"write_users": inboundProvisioningUsersWriteUsersValue,
		})
		diags.Append(objDiags...)
		inboundProvisioningValue, objDiags = types.ObjectValue(inboundProvisioningAttrTypes, map[string]attr.Value{
			"action_on_delete": types.StringPointerValue(r.InboundProvisioning.ActionOnDelete),
			"custom_schema":    inboundProvisioningCustomSchemaValue,
			"group_support":    types.BoolValue(r.InboundProvisioning.GroupSupport),
			"groups":           inboundProvisioningGroupsValue,
			"user_repository":  inboundProvisioningUserRepositoryValue,
			"users":            inboundProvisioningUsersValue,
		})
		diags.Append(objDiags...)
	}
	state.InboundProvisioning = inboundProvisioningValue

	// WsTrust
	if r.WsTrust != nil {
		var tokenGeneratorMappings []basetypes.ObjectValue
		for _, tokenGeneratorMapping := range r.WsTrust.TokenGeneratorMappings {
			spTokenGeneratorRef, objDiags := resourcelink.ToState(ctx, &tokenGeneratorMapping.SpTokenGeneratorRef)
			diags.Append(objDiags...)

			var attributeSources basetypes.SetValue
			attributeSources, objDiags = attributesources.ToState(ctx, tokenGeneratorMapping.AttributeSources)
			diags.Append(objDiags...)

			attributeContractFulfillment, objDiags := attributecontractfulfillment.ToState(ctx, &tokenGeneratorMapping.AttributeContractFulfillment)
			diags.Append(objDiags...)

			issuanceCriteria, objDiags := issuancecriteria.ToState(ctx, tokenGeneratorMapping.IssuanceCriteria)
			diags.Append(objDiags...)

			var restrictedVirtualEntityIds types.Set
			if len(tokenGeneratorMapping.RestrictedVirtualEntityIds) > 0 {
				restrictedVirtualEntityIds, objDiags = types.SetValueFrom(ctx, types.StringType, tokenGeneratorMapping.RestrictedVirtualEntityIds)
				diags.Append(objDiags...)
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
			diags.Append(objDiags...)
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
		diags.Append(objDiags...)

		var tokenGeneratorMappingsList types.List
		if tokenGeneratorMappings != nil {
			tokenGeneratorMappingsList, objDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: tokenGeneratorAttrTypes}, tokenGeneratorMappings)
			diags.Append(objDiags...)
		}

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
	var plan, state spIdpConnectionResourceModel

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

	apiCreateSpIdpConnection := r.apiClient.SpIdpConnectionsAPI.CreateConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSpIdpConnection = apiCreateSpIdpConnection.Body(*createSpIdpConnection)
	spIdpConnectionResponse, httpResp, err := r.apiClient.SpIdpConnectionsAPI.CreateConnectionExecute(apiCreateSpIdpConnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the SpIdpConnection", err, httpResp)
		return
	}

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
	diags = readSpIdpConnectionResponse(ctx, apiReadSpIdpConnection, &state, &state)
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
