package idpspconnection

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

type idpSpConnectionModel struct {
	SpBrowserSso                           types.Object `tfsdk:"sp_browser_sso"`
	Type                                   types.String `tfsdk:"type"`
	ConnectionId                           types.String `tfsdk:"connection_id"`
	Id                                     types.String `tfsdk:"id"`
	EntityId                               types.String `tfsdk:"entity_id"`
	Name                                   types.String `tfsdk:"name"`
	CreationDate                           types.String `tfsdk:"creation_date"`
	Active                                 types.Bool   `tfsdk:"active"`
	BaseUrl                                types.String `tfsdk:"base_url"`
	DefaultVirtualEntityId                 types.String `tfsdk:"default_virtual_entity_id"`
	VirtualEntityIds                       types.Set    `tfsdk:"virtual_entity_ids"`
	MetadataReloadSettings                 types.Object `tfsdk:"metadata_reload_settings"`
	Credentials                            types.Object `tfsdk:"credentials"`
	ContactInfo                            types.Object `tfsdk:"contact_info"`
	LicenseConnectionGroup                 types.String `tfsdk:"license_connection_group"`
	LoggingMode                            types.String `tfsdk:"logging_mode"`
	AdditionalAllowedEntitiesConfiguration types.Object `tfsdk:"additional_allowed_entities_configuration"`
	ExtendedProperties                     types.Map    `tfsdk:"extended_properties"`
	AttributeQuery                         types.Object `tfsdk:"attribute_query"`
	WsTrust                                types.Object `tfsdk:"ws_trust"`
	ApplicationName                        types.String `tfsdk:"application_name"`
	ApplicationIconUrl                     types.String `tfsdk:"application_icon_url"`
	OutboundProvision                      types.Object `tfsdk:"outbound_provision"`
	ConnectionTargetType                   types.String `tfsdk:"connection_target_type"`
}

var (
	resourceLinkObjectType = types.ObjectType{AttrTypes: resourcelink.AttrType()}

	metadataReloadSettingsAttrTypes = map[string]attr.Type{
		"enable_auto_metadata_update": types.BoolType,
		"metadata_url_ref":            resourceLinkObjectType,
	}

	certsListType = types.ListType{
		ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"cert_view": types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":                        types.StringType,
				"serial_number":             types.StringType,
				"subject_dn":                types.StringType,
				"subject_alternative_names": types.ListType{ElemType: types.StringType},
				"issuer_dn":                 types.StringType,
				"valid_from":                types.StringType,
				"expires":                   types.StringType,
				"key_algorithm":             types.StringType,
				"key_size":                  types.Int64Type,
				"signature_algorithm":       types.StringType,
				"version":                   types.Int64Type,
				"sha1fingerprint":           types.StringType,
				"sha256fingerprint":         types.StringType,
				"status":                    types.StringType,
				"crypto_provider":           types.StringType,
			}},
			"x509file": types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":              types.StringType,
				"file_data":       types.StringType,
				"crypto_provider": types.StringType,
			}},
			"active_verification_cert":    types.BoolType,
			"primary_verification_cert":   types.BoolType,
			"secondary_verification_cert": types.BoolType,
			"encryption_cert":             types.BoolType,
		}},
	}
	signingSettingsAttrTypes = map[string]attr.Type{
		"signing_key_pair_ref":              resourceLinkObjectType,
		"alternative_signing_key_pair_refs": types.ListType{ElemType: resourceLinkObjectType},
		"algorithm":                         types.StringType,
		"include_cert_in_signature":         types.BoolType,
		"include_raw_key_in_signature":      types.BoolType,
	}
	credentialsAttrTypes = map[string]attr.Type{
		"block_encryption_algorithm": types.StringType,
		"certs":                      certsListType,
		"decryption_key_pair_ref":    resourceLinkObjectType,
		"inbound_back_channel_auth": types.ObjectType{AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"http_basic_credentials": types.ObjectType{AttrTypes: map[string]attr.Type{
				"username":           types.StringType,
				"password":           types.StringType,
				"encrypted_password": types.StringType,
			}},
			"digital_signature":       types.BoolType,
			"verification_subject_dn": types.StringType,
			"verification_issuer_dn":  types.StringType,
			"certs":                   certsListType,
			"require_ssl":             types.BoolType,
		}},
		"key_transport_algorithm": types.StringType,
		"outbound_back_channel_auth": types.ObjectType{AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"http_basic_credentials": types.ObjectType{AttrTypes: map[string]attr.Type{
				"username":           types.StringType,
				"password":           types.StringType,
				"encrypted_password": types.StringType,
			}},
			"digital_signature":     types.BoolType,
			"ssl_auth_key_pair_ref": resourceLinkObjectType,
			"validate_partner_cert": types.BoolType,
		}},
		"secondary_decryption_key_pair_ref": resourceLinkObjectType,
		"signing_settings":                  types.ObjectType{AttrTypes: signingSettingsAttrTypes},
		"verification_issuer_dn":            types.StringType,
		"verification_subject_dn":           types.StringType,
	}

	contactInfoAttrTypes = map[string]attr.Type{
		"company":    types.StringType,
		"email":      types.StringType,
		"first_name": types.StringType,
		"last_name":  types.StringType,
		"phone":      types.StringType,
	}

	additionalAllowedEntitiesConfigurationAttrTypes = map[string]attr.Type{
		"allow_additional_entities": types.BoolType,
		"allow_all_entities":        types.BoolType,
		"additional_allowed_entities": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"entity_id":          types.StringType,
			"entity_description": types.StringType,
		}}},
	}

	extendedPropertiesElemAttrTypes = map[string]attr.Type{
		"values": types.ListType{ElemType: types.StringType},
	}

	spBrowserSsoAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":        types.StringType,
			"name_format": types.StringType,
		},
	}
	attributeContractFulfillmentElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrType()},
		"value":  types.StringType,
	}}
	attributeContractFulfillmentAttrType = types.MapType{
		ElemType: attributeContractFulfillmentElemAttrType,
	}
	issuanceCriteriaAttrType = types.ObjectType{
		AttrTypes: issuancecriteria.AttrType(),
	}
	idpAdapterAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"pseudonym": types.BoolType,
			"masked":    types.BoolType,
		},
	}
	spBrowserSSOAttrTypes = map[string]attr.Type{
		"protocol":          types.StringType,
		"ws_fed_token_type": types.StringType,
		"ws_trust_version":  types.StringType,
		"enabled_profiles":  types.ListType{ElemType: types.StringType},
		"incoming_bindings": types.ListType{ElemType: types.StringType},
		"message_customizations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"context_name":       types.StringType,
			"message_expression": types.StringType,
		}}},
		"url_whitelist_entries": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"valid_domain":             types.StringType,
			"valid_path":               types.StringType,
			"allow_query_and_fragment": types.BoolType,
			"require_https":            types.BoolType,
		}}},
		"artifact": types.ObjectType{AttrTypes: map[string]attr.Type{
			"lifetime": types.Int64Type,
			"resolver_locations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"index": types.Int64Type,
				"url":   types.StringType,
			}}},
			"source_id": types.StringType,
		}},
		"slo_service_endpoints": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"binding":      types.StringType,
			"url":          types.StringType,
			"response_url": types.StringType,
		}}},
		"default_target_url":            types.StringType,
		"always_sign_artifact_response": types.BoolType,
		"sso_service_endpoints": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"binding":    types.StringType,
			"url":        types.StringType,
			"is_default": types.BoolType,
			"index":      types.Int64Type,
		}}},
		"sp_saml_identity_mapping":      types.StringType,
		"sp_ws_fed_identity_mapping":    types.StringType,
		"sign_response_as_required":     types.BoolType,
		"sign_assertions":               types.BoolType,
		"require_signed_authn_requests": types.BoolType,
		"encryption_policy": types.ObjectType{AttrTypes: map[string]attr.Type{
			"encrypt_assertion":             types.BoolType,
			"encrypted_attributes":          types.ListType{ElemType: types.StringType},
			"encrypt_slo_subject_name_id":   types.BoolType,
			"slo_subject_name_id_encrypted": types.BoolType,
		}},
		"attribute_contract": types.ObjectType{AttrTypes: map[string]attr.Type{
			"core_attributes":     types.ListType{ElemType: spBrowserSsoAttributeAttrType},
			"extended_attributes": types.ListType{ElemType: spBrowserSsoAttributeAttrType},
		}},
		"adapter_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"idp_adapter_ref":               resourceLinkObjectType,
			"restrict_virtual_entity_ids":   types.BoolType,
			"restricted_virtual_entity_ids": types.ListType{ElemType: types.StringType},
			"adapter_override_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":                    types.StringType,
				"name":                  types.StringType,
				"plugin_descriptor_ref": resourceLinkObjectType,
				"parent_ref":            resourceLinkObjectType,
				"configuration":         types.ObjectType{AttrTypes: pluginconfiguration.AttrType()},
				"authn_ctx_class_ref":   types.StringType,
				"attribute_mapping": types.ObjectType{AttrTypes: map[string]attr.Type{
					"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
					"attribute_contract_fulfillment": attributeContractFulfillmentAttrType,
					"issuance_criteria":              issuanceCriteriaAttrType,
					"inherited":                      types.BoolType,
				}},
				"attribute_contract": types.ObjectType{AttrTypes: map[string]attr.Type{
					"core_attributes":           types.ListType{ElemType: idpAdapterAttributeAttrType},
					"extended_attributes":       types.ListType{ElemType: idpAdapterAttributeAttrType},
					"unique_user_key_attribute": types.StringType,
					"mask_ognl_values":          types.BoolType,
					"inherited":                 types.BoolType,
				}},
			}},
			"abort_sso_transaction_as_fail_safe": types.BoolType,
			"attribute_sources":                  types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
			"attribute_contract_fulfillment":     attributeContractFulfillmentAttrType,
			"issuance_criteria":                  issuanceCriteriaAttrType,
		}}},
		"authentication_policy_contract_assertion_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"authentication_policy_contract_ref": resourceLinkObjectType,
			"restrict_virtual_entity_ids":        types.BoolType,
			"restricted_virtual_entity_ids":      types.ListType{ElemType: types.StringType},
			"abort_sso_transaction_as_fail_safe": types.BoolType,
			"attribute_sources":                  types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
			"attribute_contract_fulfillment":     attributeContractFulfillmentAttrType,
			"issuance_criteria":                  issuanceCriteriaAttrType,
		}}},
		"assertion_lifetime": types.ObjectType{AttrTypes: map[string]attr.Type{
			"minutes_before": types.Int64Type,
			"minutes_after":  types.Int64Type,
		}},
	}

	policyAttrTypes = map[string]attr.Type{
		"sign_response":                  types.BoolType,
		"sign_assertion":                 types.BoolType,
		"encrypt_assertion":              types.BoolType,
		"require_signed_attribute_query": types.BoolType,
		"require_encrypted_name_id":      types.BoolType,
	}
	attributeQueryAttrTypes = map[string]attr.Type{
		"attributes":                     types.ListType{ElemType: types.StringType},
		"attribute_contract_fulfillment": attributeContractFulfillmentAttrType,
		"issuance_criteria":              issuanceCriteriaAttrType,
		"policy":                         types.ObjectType{AttrTypes: policyAttrTypes},
		"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
	}

	spWsTrustAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"namespace": types.StringType,
		},
	}
	wsTrustAttrTypes = map[string]attr.Type{
		"partner_service_ids":      types.ListType{ElemType: types.StringType},
		"oauth_assertion_profiles": types.BoolType,
		"default_token_type":       types.StringType,
		"generate_key":             types.BoolType,
		"encrypt_saml2_assertion":  types.BoolType,
		"minutes_before":           types.Int64Type,
		"minutes_after":            types.Int64Type,
		"attribute_contract": types.ObjectType{AttrTypes: map[string]attr.Type{
			"core_attributes":     types.ListType{ElemType: spWsTrustAttributeAttrType},
			"extended_attributes": types.ListType{ElemType: spWsTrustAttributeAttrType},
		}},
		"token_processor_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"idp_token_processor_ref":        resourceLinkObjectType,
			"restricted_virtual_entity_ids":  types.ListType{ElemType: types.StringType},
			"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
			"attribute_contract_fulfillment": attributeContractFulfillmentAttrType,
			"issuance_criteria":              issuanceCriteriaAttrType,
		}}},
		"abort_if_not_fulfilled_from_request": types.BoolType,
		"request_contract_ref":                resourceLinkObjectType,
		"message_customizations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"context_name":       types.StringType,
			"message_expression": types.StringType,
		}}},
	}

	channelSourceLocationAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"group_dn":      types.StringType,
			"filter":        types.StringType,
			"nested_search": types.BoolType,
		},
	}
	customSchemaAttrTypes = map[string]attr.Type{
		"namespace": types.StringType,
		"attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":           types.StringType,
			"multi_valued":   types.BoolType,
			"types":          types.ListType{ElemType: types.StringType},
			"sub_attributes": types.ListType{ElemType: types.StringType},
		}}},
	}
	targetSettingsElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":      types.StringType,
		"value":     types.StringType,
		"inherited": types.BoolType,
	}}
	saasFieldInfoAttrTypes = map[string]attr.Type{
		"attribute_names": types.ListType{ElemType: types.StringType},
		"default_value":   types.StringType,
		"expression":      types.StringType,
		"create_only":     types.BoolType,
		"trim":            types.BoolType,
		"character_case":  types.StringType,
		"parser":          types.StringType,
		"masked":          types.BoolType,
	}
	attributeMappingElemAttrTypes = types.ObjectType{AttrTypes: map[string]attr.Type{
		"field_name":      types.StringType,
		"saas_field_info": types.ObjectType{AttrTypes: saasFieldInfoAttrTypes},
	}}
	channelSourceAttrTypes = map[string]attr.Type{
		"data_source":         resourceLinkObjectType,
		"guid_attribute_name": types.StringType,
		"guid_binary":         types.BoolType,
		"change_detection_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
			"user_object_class":         types.StringType,
			"group_object_class":        types.StringType,
			"changed_users_algorithm":   types.StringType,
			"usn_attribute_name":        types.StringType,
			"time_stamp_attribute_name": types.StringType,
		}},
		"group_membership_detection": types.ObjectType{AttrTypes: map[string]attr.Type{
			"member_of_group_attribute_name": types.StringType,
			"group_member_attribute_name":    types.StringType,
		}},
		"account_management_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
			"account_status_attribute_name": types.StringType,
			"account_status_algorithm":      types.StringType,
			"flag_comparison_value":         types.StringType,
			"flag_comparison_status":        types.BoolType,
			"default_status":                types.BoolType,
		}},
		"base_dn":               types.StringType,
		"user_source_location":  channelSourceLocationAttrType,
		"group_source_location": channelSourceLocationAttrType,
	}
	channelsElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"active":                types.BoolType,
		"channel_source":        types.ObjectType{AttrTypes: channelSourceAttrTypes},
		"attribute_mapping":     types.SetType{ElemType: attributeMappingElemAttrTypes},
		"attribute_mapping_all": types.SetType{ElemType: attributeMappingElemAttrTypes},
		"name":                  types.StringType,
		"max_threads":           types.Int64Type,
		"timeout":               types.Int64Type,
	}}
	outboundProvisionAttrTypes = map[string]attr.Type{
		"type":                types.StringType,
		"target_settings":     types.ListType{ElemType: targetSettingsElemAttrType},
		"target_settings_all": types.ListType{ElemType: targetSettingsElemAttrType},
		"custom_schema":       types.ObjectType{AttrTypes: customSchemaAttrTypes},
		"channels":            types.ListType{ElemType: channelsElemAttrType},
	}

	emptyStringSet, _ = types.SetValue(types.StringType, nil)

	groupSourceLocationDefault, _ = types.ObjectValue(channelSourceLocationAttrType.AttrTypes, map[string]attr.Value{
		"filter":        types.StringNull(),
		"group_dn":      types.StringNull(),
		"nested_search": types.BoolValue(false),
	})

	certsDefault, _ = types.ListValue(certsListType.ElemType, nil)
)
