package idpspconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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

func readIdpSpconnectionResponse(ctx context.Context, r *client.SpConnection, state *idpSpConnectionModel, plan *idpSpConnectionModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.ConnectionId = types.StringPointerValue(r.Id)
	state.Id = types.StringPointerValue(r.Id)
	state.Type = types.StringPointerValue(r.Type)
	state.EntityId = types.StringValue(r.EntityId)
	state.Name = types.StringValue(r.Name)
	state.Active = types.BoolPointerValue(r.Active)
	state.BaseUrl = types.StringPointerValue(r.BaseUrl)
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	state.ApplicationName = types.StringPointerValue(r.ApplicationName)
	state.ApplicationIconUrl = types.StringPointerValue(r.ApplicationIconUrl)
	state.ConnectionTargetType = types.StringPointerValue(r.ConnectionTargetType)

	if r.CreationDate != nil {
		state.CreationDate = types.StringValue(r.CreationDate.Format(time.RFC3339))
	} else {
		state.CreationDate = types.StringNull()
	}

	state.VirtualEntityIds, respDiags = types.SetValueFrom(ctx, types.StringType, r.VirtualEntityIds)
	diags.Append(respDiags...)

	state.MetadataReloadSettings, respDiags = types.ObjectValueFrom(ctx, metadataReloadSettingsAttrTypes, r.MetadataReloadSettings)
	diags.Append(respDiags...)

	state.Credentials, respDiags = types.ObjectValueFrom(ctx, credentialsAttrTypes, r.Credentials)
	diags.Append(respDiags...)
	if r.Credentials != nil && r.Credentials.SigningSettings != nil && r.Credentials.SigningSettings.IncludeCertInSignature == nil {
		// PF returns false for include_cert_in_signature as nil. If nil is returned, just set it to false
		credentialsAttrs := state.Credentials.Attributes()
		signingSettingsAttrs := credentialsAttrs["signing_settings"].(types.Object).Attributes()
		signingSettingsAttrs["include_cert_in_signature"] = types.BoolValue(false)
		newSigningSettings, respDiags := types.ObjectValue(signingSettingsAttrTypes, signingSettingsAttrs)
		diags.Append(respDiags...)
		credentialsAttrs["signing_settings"] = newSigningSettings
		state.Credentials, respDiags = types.ObjectValue(credentialsAttrTypes, credentialsAttrs)
		diags.Append(respDiags...)
	}

	state.ContactInfo, respDiags = types.ObjectValueFrom(ctx, contactInfoAttrTypes, r.ContactInfo)
	diags.Append(respDiags...)

	state.AdditionalAllowedEntitiesConfiguration, respDiags = types.ObjectValueFrom(ctx, additionalAllowedEntitiesConfigurationAttrTypes, r.AdditionalAllowedEntitiesConfiguration)
	diags.Append(respDiags...)

	state.ExtendedProperties, respDiags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: extendedPropertiesElemAttrTypes}, r.ExtendedProperties)
	diags.Append(respDiags...)

	state.SpBrowserSso, respDiags = types.ObjectValueFrom(ctx, spBrowserSSOAttrTypes, r.SpBrowserSso)
	diags.Append(respDiags...)

	if r.AttributeQuery != nil {
		attributeQueryValues := map[string]attr.Value{}
		attributeQueryValues["attributes"], respDiags = types.ListValueFrom(ctx, types.StringType, r.AttributeQuery.Attributes)
		diags.Append(respDiags...)

		attributeQueryValues["attribute_contract_fulfillment"], respDiags = types.MapValueFrom(ctx, attributeContractFulfillmentElemAttrType, r.AttributeQuery.AttributeContractFulfillment)
		diags.Append(respDiags...)

		attributeQueryValues["issuance_criteria"], respDiags = issuancecriteria.ToState(ctx, r.AttributeQuery.IssuanceCriteria)
		diags.Append(respDiags...)

		attributeQueryValues["policy"], respDiags = types.ObjectValueFrom(ctx, policyAttrTypes, r.AttributeQuery.Policy)
		diags.Append(respDiags...)

		attributeQueryValues["attribute_sources"], respDiags = attributesources.ToState(ctx, r.AttributeQuery.AttributeSources)
		diags.Append(respDiags...)

		state.AttributeQuery, respDiags = types.ObjectValueFrom(ctx, attributeQueryAttrTypes, r.AttributeQuery)
		diags.Append(respDiags...)
	} else {
		state.AttributeQuery = types.ObjectNull(attributeQueryAttrTypes)
	}

	state.WsTrust, respDiags = types.ObjectValueFrom(ctx, wsTrustAttrTypes, r.WsTrust)
	diags.Append(respDiags...)

	if r.OutboundProvision != nil {
		outboundProvisionAttrs := map[string]attr.Value{
			"type": types.StringValue(r.OutboundProvision.Type),
		}

		// PF can return extra target_settings that were not included in the request
		plannedTargetSettingsNames := []string{}
		plannedTargetSettingsValues := map[string]string{}
		if internaltypes.IsDefined(plan.OutboundProvision) {
			targetSettings := plan.OutboundProvision.Attributes()["target_settings"].(types.List)
			for _, plannedTargetSettings := range targetSettings.Elements() {
				nameStrVal := plannedTargetSettings.(types.Object).Attributes()["name"].(types.String)
				if internaltypes.IsDefined(nameStrVal) {
					plannedTargetSettingsNames = append(plannedTargetSettingsNames, nameStrVal.ValueString())

					valueStrVal := plannedTargetSettings.(types.Object).Attributes()["value"].(types.String)
					plannedTargetSettingsValues[nameStrVal.ValueString()] = valueStrVal.ValueString()
				}
			}
		}

		targetSettingsSlice := []attr.Value{}
		targetSettingsAllSlice := []attr.Value{}
		for _, targetSettings := range r.OutboundProvision.TargetSettings {
			value := types.StringPointerValue(targetSettings.Value)

			// Check if this object was in the plan
			inPlan := false
			for _, name := range plannedTargetSettingsNames {
				if name == targetSettings.Name {
					inPlan = true

					// If PF returns nil for the value, then it must be encrypted. Just use the value from the plan in that case
					if targetSettings.Value == nil {
						value = types.StringValue(plannedTargetSettingsValues[targetSettings.Name])
					}
					break
				}
			}
			targetSettingsObj, respDiags := types.ObjectValue(targetSettingsElemAttrType.AttrTypes, map[string]attr.Value{
				"name":      types.StringValue(targetSettings.Name),
				"value":     value,
				"inherited": types.BoolPointerValue(targetSettings.Inherited),
			})
			diags.Append(respDiags...)
			if inPlan {
				targetSettingsSlice = append(targetSettingsSlice, targetSettingsObj)
			}
			targetSettingsAllSlice = append(targetSettingsAllSlice, targetSettingsObj)
		}
		outboundProvisionAttrs["target_settings"], respDiags = types.ListValue(targetSettingsElemAttrType, targetSettingsSlice)
		diags.Append(respDiags...)
		outboundProvisionAttrs["target_settings_all"], respDiags = types.ListValue(targetSettingsElemAttrType, targetSettingsAllSlice)
		diags.Append(respDiags...)

		outboundProvisionAttrs["custom_schema"], respDiags = types.ObjectValueFrom(ctx, customSchemaAttrTypes, r.OutboundProvision.CustomSchema)
		diags.Append(respDiags...)

		channels := []types.Object{}
		plannedChannels := []attr.Value{}
		plannedChannelsAttr := plan.OutboundProvision.Attributes()["channels"]
		if plannedChannelsAttr != nil {
			plannedChannels = plannedChannelsAttr.(types.List).Elements()
		}
		numPlannedChannels := len(plannedChannels)
		for i, channel := range r.OutboundProvision.Channels {
			channelAttrs := map[string]attr.Value{
				"active":      types.BoolValue(channel.Active),
				"name":        types.StringValue(channel.Name),
				"max_threads": types.Int64Value(channel.MaxThreads),
				"timeout":     types.Int64Value(channel.Timeout),
			}

			channelAttrs["channel_source"], respDiags = types.ObjectValueFrom(ctx, channelSourceAttrTypes, channel.ChannelSource)
			diags.Append(respDiags...)

			// PF can return extra attribute_mapping elements that were not included in the request
			attributeMappingNamesInPlan := []string{}
			if i < numPlannedChannels {
				plannedChannel := plannedChannels[i].(types.Object)
				plannedMapping := plannedChannel.Attributes()["attribute_mapping"].(types.Set)
				if internaltypes.IsDefined(plannedMapping) {
					for _, mapping := range plannedMapping.Elements() {
						mappingObj := mapping.(types.Object)
						if internaltypes.IsDefined(mappingObj) {
							attributeMappingNamesInPlan = append(attributeMappingNamesInPlan, mappingObj.Attributes()["field_name"].(types.String).ValueString())
						}
					}
				}
			}

			attributeMappingSlice := []attr.Value{}
			attributeMappingAllSlice := []attr.Value{}
			for _, attributeMapping := range channel.AttributeMapping {
				attributeMappingAttrValues := map[string]attr.Value{
					"field_name": types.StringValue(attributeMapping.FieldName),
				}

				attributeMappingAttrValues["saas_field_info"], respDiags = types.ObjectValueFrom(ctx, saasFieldInfoAttrTypes, attributeMapping.SaasFieldInfo)
				diags.Append(respDiags...)

				attributeMappingObj, respDiags := types.ObjectValue(attributeMappingElemAttrTypes.AttrTypes, attributeMappingAttrValues)
				diags.Append(respDiags...)

				// Check if this object was in the plan
				inPlan := false
				for _, attributeMappingNameInPlan := range attributeMappingNamesInPlan {
					if attributeMappingNameInPlan == attributeMapping.FieldName {
						inPlan = true
						break
					}
				}
				if inPlan {
					attributeMappingSlice = append(attributeMappingSlice, attributeMappingObj)
				}
				attributeMappingAllSlice = append(attributeMappingAllSlice, attributeMappingObj)
			}
			channelAttrs["attribute_mapping"], respDiags = types.SetValue(attributeMappingElemAttrTypes, attributeMappingSlice)
			diags.Append(respDiags...)
			channelAttrs["attribute_mapping_all"], respDiags = types.SetValue(attributeMappingElemAttrTypes, attributeMappingAllSlice)
			diags.Append(respDiags...)

			channelObj, respDiags := types.ObjectValue(channelsElemAttrType.AttrTypes, channelAttrs)
			diags.Append(respDiags...)
			channels = append(channels, channelObj)
		}
		outboundProvisionAttrs["channels"], respDiags = types.ListValueFrom(ctx, channelsElemAttrType, channels)
		diags.Append(respDiags...)

		state.OutboundProvision, respDiags = types.ObjectValue(outboundProvisionAttrTypes, outboundProvisionAttrs)
		diags.Append(respDiags...)
	} else {
		state.OutboundProvision = types.ObjectNull(outboundProvisionAttrTypes)
	}

	return diags
}
