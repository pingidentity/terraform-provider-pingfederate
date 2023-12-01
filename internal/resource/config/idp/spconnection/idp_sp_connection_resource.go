package idpspconnection

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpSpConnectionResource{}
	_ resource.ResourceWithConfigure   = &idpSpConnectionResource{}
	_ resource.ResourceWithImportState = &idpSpConnectionResource{}

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

// IdpSpConnectionResource is a helper function to simplify the provider implementation.
func IdpSpConnectionResource() resource.Resource {
	return &idpSpConnectionResource{}
}

// idpSpConnectionResource is the resource implementation.
type idpSpConnectionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpSpConnectionResourceModel struct {
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

// GetSchema defines the schema for the resource.
func (r *idpSpConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certsSchema := schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"active_verification_cert": schema.BoolAttribute{
					Optional:    true,
					Description: "Indicates whether this is an active signature verification certificate.",
				},
				"cert_view": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"crypto_provider": schema.StringAttribute{
							Optional:    true,
							Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"LOCAL",
									"HSM",
								),
							},
						},
						"expires": schema.StringAttribute{
							Optional:    true,
							Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
						},
						"id": schema.StringAttribute{
							Optional:    true,
							Description: "The persistent, unique ID for the certificate.",
						},
						"issuer_dn": schema.StringAttribute{
							Optional:    true,
							Description: "The issuer's distinguished name.",
						},
						"key_algorithm": schema.StringAttribute{
							Optional:    true,
							Description: "The public key algorithm.",
						},
						"key_size": schema.Int64Attribute{
							Optional:    true,
							Description: "The public key size.",
						},
						"serial_number": schema.StringAttribute{
							Optional:    true,
							Description: "The serial number assigned by the CA.",
						},
						"sha1fingerprint": schema.StringAttribute{
							Optional:    true,
							Description: "SHA-1 fingerprint in Hex encoding.",
						},
						"sha256fingerprint": schema.StringAttribute{
							Optional:    true,
							Description: "SHA-256 fingerprint in Hex encoding.",
						},
						"signature_algorithm": schema.StringAttribute{
							Optional:    true,
							Description: "The signature algorithm.",
						},
						"status": schema.StringAttribute{
							Optional:    true,
							Description: "Status of the item.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"VALID",
									"EXPIRED",
									"NOT_YET_VALID",
									"REVOKED",
								),
							},
						},
						"subject_alternative_names": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "The subject alternative names (SAN).",
						},
						"subject_dn": schema.StringAttribute{
							Optional:    true,
							Description: "The subject's distinguished name.",
						},
						"valid_from": schema.StringAttribute{
							Optional:    true,
							Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
						},
						"version": schema.Int64Attribute{
							Optional:    true,
							Description: "The X.509 version to which the item conforms.",
						},
					},
					Optional:    true,
					Description: "Certificate details.",
				},
				"encryption_cert": schema.BoolAttribute{
					Optional:    true,
					Description: "Indicates whether to use this cert to encrypt outgoing assertions. Only one certificate in the collection can have this flag set.",
				},
				"primary_verification_cert": schema.BoolAttribute{
					Optional:    true,
					Description: "Indicates whether this is the primary signature verification certificate. Only one certificate in the collection can have this flag set.",
				},
				"secondary_verification_cert": schema.BoolAttribute{
					Optional:    true,
					Description: "Indicates whether this is the secondary signature verification certificate. Only one certificate in the collection can have this flag set.",
				},
				"x509file": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"crypto_provider": schema.StringAttribute{
							Optional:    true,
							Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"LOCAL",
									"HSM",
								),
							},
						},
						"file_data": schema.StringAttribute{
							Required:    true,
							Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
						},
						"id": schema.StringAttribute{
							Optional:    true,
							Description: "The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
						},
					},
					Required:    true,
					Description: "Encoded certificate data.",
				},
			},
		},
		Optional:    true,
		Computed:    true,
		Default:     listdefault.StaticValue(certsDefault),
		Description: "The certificates used for signature verification and XML encryption.",
	}

	httpBasicCredentialsSchema := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"encrypted_password": schema.StringAttribute{
				Optional:    true,
				Description: "For GET requests, this field contains the encrypted password, if one exists.  For POST and PUT requests, if you wish to reuse the existing password, this field should be passed back unchanged.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "User password.  To update the password, specify the plaintext value in this field.  This field will not be populated for GET requests.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "The username.",
			},
		},
		Optional:    true,
		Description: "Username and password credentials.",
	}

	adapterOverrideSettingsAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"masked": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
			},
			"pseudonym": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
			},
		},
	}

	spBrowserSSOAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
			},
			"name_format": schema.StringAttribute{
				Optional:    true,
				Description: "The SAML Name Format for the attribute.",
			},
		},
	}

	wsTrustAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
			},
			"namespace": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "The attribute namespace.  This is required when the Default Token Type is SAML2.0 or SAML1.1 or SAML1.1 for Office 365.",
			},
		},
	}

	messageCustomizationsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"context_name": schema.StringAttribute{
				Optional:    true,
				Description: "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be 'assertion', 'authn-response' or 'authn-request'.",
			},
			"message_expression": schema.StringAttribute{
				Optional:    true,
				Description: "The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.",
			},
		},
	}

	channelsAttributeMappingNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"field_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of target field.",
			},
			"saas_field_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_names": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						//Default:     listdefault.StaticValue(emptyStringList),
						Description: "The list of source attribute names used to generate or map to a target field",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
					},
					"character_case": schema.StringAttribute{
						Optional: true,
						Computed: true,
						//Default:     stringdefault.StaticString("NONE"),
						Description: "The character case of the field value.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"LOWER",
								"UPPER",
								"NONE",
							),
						},
					},
					"create_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						//Default:     booldefault.StaticBool(false),
						Description: "Indicates whether this field is a create only field and cannot be updated.",
					},
					"default_value": schema.StringAttribute{
						Optional: true,
						Computed: true,
						//Default:     stringdefault.StaticString(""),
						Description: "The default value for the target field",
					},
					"expression": schema.StringAttribute{
						Optional: true,
						Computed: true,
						//Default:     stringdefault.StaticString(""),
						Description: "An OGNL expression to obtain a value.",
					},
					"masked": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						//Default:     booldefault.StaticBool(false),
						Description: "Indicates whether the attribute should be masked in server logs.",
					},
					"parser": schema.StringAttribute{
						Optional: true,
						Computed: true,
						//Default:     stringdefault.StaticString("NONE"),
						Description: "Indicates how the field shall be parsed.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"EXTRACT_CN_FROM_DN",
								"EXTRACT_USERNAME_FROM_EMAIL",
								"NONE",
							),
						},
					},
					"trim": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						//Default:     booldefault.StaticBool(false),
						Description: "Indicates whether field should be trimmed before provisioning.",
					},
				},
				Required:    true,
				Description: "The settings that represent how attribute values from source data store will be mapped into Fields specified by the service provider.",
			},
		},
	}

	outboundProvisionTargetSettingsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"inherited": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the configuration field.",
			},
			"value": schema.StringAttribute{
				Optional:    true,
				Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
			},
		},
	}

	schema := schema.Schema{
		Description: "Manages an IdP SP Connection",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Specifies whether the connection is active and ready to process incoming requests. The default value is false.",
			},
			"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"additional_allowed_entities": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"entity_description": schema.StringAttribute{
									Optional:    true,
									Description: "Entity description.",
								},
								"entity_id": schema.StringAttribute{
									Optional:    true,
									Description: "Unique entity identifier.",
								},
							},
						},
						Optional:    true,
						Description: "An array of additional allowed entities or issuers to be accepted during entity or issuer validation.",
					},
					"allow_additional_entities": schema.BoolAttribute{
						Optional:    true,
						Description: "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation.",
					},
					"allow_all_entities": schema.BoolAttribute{
						Optional:    true,
						Description: "Set to true to accept any entity or issuer during entity or issuer validation. (Not Recommended)",
					},
				},
				Optional:    true,
				Description: "Additional allowed entities or issuers configuration. Currently only used in OIDC IdP (RP) connection.",
			},
			"application_icon_url": schema.StringAttribute{
				Optional:    true,
				Description: "The application icon url.",
			},
			"application_name": schema.StringAttribute{
				Optional:    true,
				Description: "The application name.",
			},
			"attribute_query": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
					"attribute_sources":              attributesources.ToSchema(1),
					"attributes": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "The list of attributes that may be returned to the SP in the response to an attribute request.",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
							listvalidator.SizeAtLeast(1),
						},
					},
					"issuance_criteria": issuancecriteria.ToSchema(),
					"policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Optional:    true,
								Description: "Encrypt the assertion.",
							},
							"require_encrypted_name_id": schema.BoolAttribute{
								Optional:    true,
								Description: "Require an encrypted name identifier.",
							},
							"require_signed_attribute_query": schema.BoolAttribute{
								Optional:    true,
								Description: "Require signed attribute query.",
							},
							"sign_assertion": schema.BoolAttribute{
								Optional:    true,
								Description: "Sign the assertion.",
							},
							"sign_response": schema.BoolAttribute{
								Optional:    true,
								Description: "Sign the response.",
							},
						},
						Optional:    true,
						Description: "The attribute query profile's security policy.",
					},
				},
				Optional:    true,
				Description: "The attribute query profile supports SPs in requesting user attributes.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The fully-qualified hostname and port on which your partner's federation deployment runs.",
			},
			"connection_target_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Description: "The connection target type. This field is intended for bulk import/export usage. Changing its value may result in unexpected behavior.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"STANDARD",
						"SALESFORCE",
						"SALESFORCE_CP",
						"SALESFORCE_PP",
						"PINGONE_SCIM11",
					),
				},
			},
			"contact_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Optional:    true,
						Description: "Company name.",
					},
					"email": schema.StringAttribute{
						Optional:    true,
						Description: "Contact email address.",
					},
					"first_name": schema.StringAttribute{
						Optional:    true,
						Description: "Contact first name.",
					},
					"last_name": schema.StringAttribute{
						Optional:    true,
						Description: "Contact last name.",
					},
					"phone": schema.StringAttribute{
						Optional:    true,
						Description: "Contact phone number.",
					},
				},
				Optional:    true,
				Description: "Contact information.",
			},
			"creation_date": schema.StringAttribute{
				Optional: false,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The time at which the connection was created. This property is read only and is ignored on PUT and POST requests.",
			},
			"credentials": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"block_encryption_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "The algorithm used to encrypt assertions sent to this partner. AES_128, AES_256, AES_128_GCM, AES_192_GCM, AES_256_GCM and Triple_DES are supported.",
					},
					"certs":                   certsSchema,
					"decryption_key_pair_ref": resourcelink.SingleNestedAttribute(),
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"certs": certsSchema,
							"digital_signature": schema.BoolAttribute{
								Optional:    true,
								Description: "If incoming or outgoing messages must be signed.",
							},
							"http_basic_credentials": httpBasicCredentialsSchema,
							"require_ssl": schema.BoolAttribute{
								Optional:    true,
								Description: "Incoming HTTP transmissions must use a secure channel.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The back channel authentication type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"INBOUND",
										"OUTBOUND",
									),
								},
							},
							"verification_issuer_dn": schema.StringAttribute{
								Optional:    true,
								Description: "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
							},
							"verification_subject_dn": schema.StringAttribute{
								Optional:    true,
								Description: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
							},
						},
						Optional: true,
					},
					"key_transport_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "The algorithm used to transport keys to this partner. RSA_OAEP, RSA_OAEP_256 and RSA_v15 are supported.",
					},
					"outbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"digital_signature": schema.BoolAttribute{
								Optional:    true,
								Description: "If incoming or outgoing messages must be signed.",
							},
							"http_basic_credentials": httpBasicCredentialsSchema,
							"ssl_auth_key_pair_ref":  resourcelink.SingleNestedAttribute(),
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The back channel authentication type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"INBOUND",
										"OUTBOUND",
									),
								},
							},
							"validate_partner_cert": schema.BoolAttribute{
								Optional:    true,
								Description: "Validate the partner server certificate. Default is true.",
							},
						},
						Optional: true,
					},
					"secondary_decryption_key_pair_ref": resourcelink.SingleNestedAttribute(),
					"signing_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Optional:    true,
								Description: "The algorithm used to sign messages sent to this partner. The default is SHA1withDSA for DSA certs, SHA256withRSA for RSA certs, and SHA256withECDSA for EC certs. For RSA certs, SHA1withRSA, SHA384withRSA, and SHA512withRSA are also supported. For EC certs, SHA384withECDSA and SHA512withECDSA are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
							},
							"alternative_signing_key_pair_refs": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: resourcelink.ToSchema(),
								},
								Optional:    true,
								Description: "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
							},
							"include_cert_in_signature": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Determines whether the signing certificate is included in the signature <KeyInfo> element.",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Optional:    true,
								Description: "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element.",
							},
							"signing_key_pair_ref": resourcelink.SingleNestedAttribute(),
						},
						Optional:    true,
						Description: "Settings related to signing messages sent to this partner.",
					},
					"verification_issuer_dn": schema.StringAttribute{
						Optional:    true,
						Description: "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
					},
					"verification_subject_dn": schema.StringAttribute{
						Optional:    true,
						Description: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
					},
				},
				Optional:    true,
				Description: "The certificates and settings for encryption, signing, and signature verification.",
			},
			"default_virtual_entity_id": schema.StringAttribute{
				Optional:    true,
				Description: "The default alternate entity ID that identifies the local server to this partner. It is required when virtualEntityIds is not empty and must be included in that list.",
			},
			"entity_id": schema.StringAttribute{
				Required:    true,
				Description: "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
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
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "The persistent, unique ID for the connection. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
			},
			"license_connection_group": schema.StringAttribute{
				Optional:    true,
				Description: "The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.",
				// License connection group must not be empty if configured
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"logging_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Description: "The level of transaction logging applicable for this connection. Default is STANDARD.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"NONE",
						"STANDARD",
						"ENHANCED",
						"FULL",
					),
				},
			},
			"metadata_reload_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_auto_metadata_update": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
						Description: "Specifies whether the metadata of the connection will be automatically reloaded. The default value is true.",
					},
					"metadata_url_ref": resourcelink.SingleNestedAttribute(),
				},
				Optional:    true,
				Description: "Configuration settings to enable automatic reload of partner's metadata.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The connection name.",
			},
			"outbound_provision": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"channels": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"active": schema.BoolAttribute{
									Required:    true,
									Description: "Indicates whether the channel is the active channel for this connection.",
								},
								"attribute_mapping_all": schema.SetNestedAttribute{
									NestedObject: channelsAttributeMappingNestedObject,
									Optional:     false,
									Computed:     true,
									PlanModifiers: []planmodifier.Set{
										setplanmodifier.UseStateForUnknown(),
									},
									Description: "The mapping of attributes from the local data store into Fields specified by the service provider. This attribute will include any values set by default by PingFederate.",
								},
								"attribute_mapping": schema.SetNestedAttribute{
									NestedObject: channelsAttributeMappingNestedObject,
									Required:     true,
									Description:  "The mapping of attributes from the local data store into Fields specified by the service provider.",
								},
								"channel_source": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"account_management_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"account_status_algorithm": schema.StringAttribute{
													Required:    true,
													Description: "The account status algorithm name. \nACCOUNT_STATUS_ALGORITHM_AD -  Algorithm name for Active Directory, which uses a bitmap for each user entry. \nACCOUNT_STATUS_ALGORITHM_FLAG - Algorithm name for Oracle Directory Server and other LDAP directories that use a separate attribute to store the user's status. When this option is selected, the Flag Comparison Value and Flag Comparison Status fields should be used.",
													Validators: []validator.String{
														stringvalidator.OneOf(
															"ACCOUNT_STATUS_ALGORITHM_AD",
															"ACCOUNT_STATUS_ALGORITHM_FLAG",
														),
													},
												},
												"account_status_attribute_name": schema.StringAttribute{
													Required:    true,
													Description: "The account status attribute name.",
												},
												"default_status": schema.BoolAttribute{
													Optional:    true,
													Description: "The default status of the account.",
												},
												"flag_comparison_status": schema.BoolAttribute{
													Optional:    true,
													Description: "The flag that represents comparison status.",
												},
												"flag_comparison_value": schema.StringAttribute{
													Optional:    true,
													Description: "The flag that represents comparison value.",
												},
											},
											Required:    true,
											Description: "Account management settings.",
										},
										"base_dn": schema.StringAttribute{
											Required:    true,
											Description: "The base DN where the user records are located.",
										},
										"change_detection_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"changed_users_algorithm": schema.StringAttribute{
													Required:    true,
													Description: "The changed user algorithm. \nACTIVE_DIRECTORY_USN - For Active Directory only, this algorithm queries for update sequence numbers on user records that are larger than the last time records were checked. \nTIMESTAMP - Queries for timestamps on user records that are not older than the last time records were checked. This check is more efficient from the point of view of the PingFederate provisioner but can be more time consuming on the LDAP side, particularly with the Oracle Directory Server. \nTIMESTAMP_NO_NEGATION - Queries for timestamps on user records that are newer than the last time records were checked. This algorithm is recommended for the Oracle Directory Server.",
													Validators: []validator.String{
														stringvalidator.OneOf(
															"ACTIVE_DIRECTORY_USN",
															"TIMESTAMP",
															"TIMESTAMP_NO_NEGATION",
														),
													},
												},
												"group_object_class": schema.StringAttribute{
													Required:    true,
													Description: "The group object class.",
												},
												"time_stamp_attribute_name": schema.StringAttribute{
													Required:    true,
													Description: "The timestamp attribute name.",
												},
												"user_object_class": schema.StringAttribute{
													Required:    true,
													Description: "The user object class.",
												},
												"usn_attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The USN attribute name.",
												},
											},
											Required:    true,
											Description: "Setting to detect changes to a user or a group.",
										},
										"data_source": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to an LDAP datastore."),
										"group_membership_detection": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"group_member_attribute_name": schema.StringAttribute{
													Required:    true,
													Description: "The name of the attribute that represents group members in a group, also known as group member attribute.",
												},
												"member_of_group_attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The name of the attribute that indicates the entity is a member of a group, also known as member of attribute.",
												},
											},
											Required:    true,
											Description: "Settings to detect group memberships.",
										},
										"group_source_location": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"filter": schema.StringAttribute{
													Optional:    true,
													Description: "An LDAP filter.",
												},
												"group_dn": schema.StringAttribute{
													Optional:    true,
													Description: "The group DN for users or groups.",
												},
												"nested_search": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
													Description: "Indicates whether the search is nested.",
												},
											},
											Optional:    true,
											Computed:    true,
											Default:     objectdefault.StaticValue(groupSourceLocationDefault),
											Description: "The location settings that includes a DN and a LDAP filter.",
										},
										"guid_attribute_name": schema.StringAttribute{
											Required:    true,
											Description: "the GUID attribute name.",
										},
										"guid_binary": schema.BoolAttribute{
											Required:    true,
											Description: "Indicates whether the GUID is stored in binary format.",
										},
										"user_source_location": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"filter": schema.StringAttribute{
													Optional:    true,
													Description: "An LDAP filter.",
												},
												"group_dn": schema.StringAttribute{
													Optional:    true,
													Description: "The group DN for users or groups.",
												},
												"nested_search": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
													Description: "Indicates whether the search is nested.",
												},
											},
											Required:    true,
											Description: "The location settings that includes a DN and a LDAP filter.",
										},
									},
									Required:    true,
									Description: "The source data source and LDAP settings.",
								},
								"max_threads": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Default:     int64default.StaticInt64(1),
									Description: "The number of processing threads. The default value is 1.",
								},
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of the channel.",
								},
								"timeout": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Default:     int64default.StaticInt64(60),
									Description: "Timeout, in seconds, for individual user and group provisioning operations on the target service provider. The default value is 60.",
								},
							},
						},
						Required:    true,
						Description: "Includes settings of a source data store, managing provisioning threads and mapping of attributes.",
					},
					"custom_schema": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"multi_valued": schema.BoolAttribute{
											Optional:    true,
											Description: "Indicates whether the attribute is multi-valued.",
										},
										"name": schema.StringAttribute{
											Optional:    true,
											Description: "Name of the attribute.",
										},
										"sub_attributes": schema.ListAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Description: "List of sub-attributes for an attribute.",
										},
										"types": schema.ListAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Description: "Represents the name of each attribute type in case of multi-valued attribute.",
										},
									},
								},
								Optional: true,
							},
							"namespace": schema.StringAttribute{
								Optional: true,
							},
						},
						Optional:    true,
						Description: "Custom SCIM Attributes configuration.",
					},
					"target_settings_all": schema.ListNestedAttribute{
						NestedObject: outboundProvisionTargetSettingsNestedObject,
						Optional:     false,
						Computed:     true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						Description: "Configuration fields that includes credentials to target SaaS application. This attribute will include any values set by default by PingFederate.",
					},
					"target_settings": schema.ListNestedAttribute{
						NestedObject: outboundProvisionTargetSettingsNestedObject,
						Required:     true,
						Description:  "Configuration fields that includes credentials to target SaaS application.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The SaaS plugin type.",
					},
				},
				Optional:    true,
				Description: "Outbound Provisioning allows an IdP to create and maintain user accounts at standards-based partner sites using SCIM as well as select-proprietary provisioning partner sites that are protocol-enabled.",
			},
			"sp_browser_sso": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"adapter_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Optional:    true,
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is false.",
								},
								"adapter_override_settings": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"attribute_contract": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"core_attributes": schema.ListNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Required:     true,
													Description:  "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
												},
												"extended_attributes": schema.ListNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Optional:     true,
													Description:  "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
												},
												"inherited": schema.BoolAttribute{
													Optional:    true,
													Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
												},
												"mask_ognl_values": schema.BoolAttribute{
													Optional:    true,
													Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to false.",
												},
												"unique_user_key_attribute": schema.StringAttribute{
													Optional:    true,
													Description: "The attribute to use for uniquely identify a user's authentication sessions.",
												},
											},
											Optional:    true,
											Description: "A set of attributes exposed by an IdP adapter.",
										},
										"attribute_mapping": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
												"attribute_sources":              attributesources.ToSchema(0),
												"inherited": schema.BoolAttribute{
													Optional:    true,
													Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
												},
												"issuance_criteria": issuancecriteria.ToSchema(),
											},
											Optional:    true,
											Description: "An IdP Adapter Contract Mapping.",
										},
										"authn_ctx_class_ref": schema.StringAttribute{
											Optional:    true,
											Description: "The fixed value that indicates how the user was authenticated.",
										},
										"configuration": pluginconfiguration.ToSchema(),
										"id": schema.StringAttribute{
											Required:    true,
											Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
										},
										"name": schema.StringAttribute{
											Required:    true,
											Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
										},
										"parent_ref":            resourcelink.CompleteSingleNestedAttribute(true, false, false, "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)"),
										"plugin_descriptor_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override."),
									},
									Optional: true,
								},
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
								"attribute_sources":              attributesources.ToSchema(0),
								"idp_adapter_ref":                resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated IdP adapter. Note: This is ignored if adapter overrides for this mapping exists. In this case, the override's parent adapter reference is used."),
								"issuance_criteria":              issuancecriteria.ToSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:    true,
									Description: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Required:    true,
						Description: "A list of adapters that map to outgoing assertions.",
					},
					"always_sign_artifact_response": schema.BoolAttribute{
						Optional:    true,
						Description: "Specify to always sign the SAML ArtifactResponse.",
					},
					"artifact": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"lifetime": schema.Int64Attribute{
								Required:    true,
								Description: "The lifetime of the artifact in seconds.",
							},
							"resolver_locations": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"index": schema.Int64Attribute{
											Required:    true,
											Description: "The priority of the endpoint.",
										},
										"url": schema.StringAttribute{
											Required:    true,
											Description: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
										},
									},
								},
								Required:    true,
								Description: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
							},
							"source_id": schema.StringAttribute{
								Optional:    true,
								Description: "Source ID for SAML1.x connections",
							},
						},
						Optional:    true,
						Description: "The settings for an Artifact binding.",
					},
					"assertion_lifetime": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"minutes_after": schema.Int64Attribute{
								Required:    true,
								Description: "Assertion validity in minutes after the assertion issuance.",
							},
							"minutes_before": schema.Int64Attribute{
								Required:    true,
								Description: "Assertion validity in minutes before the assertion issuance.",
							},
						},
						Required:    true,
						Description: "The timeframe of validity before and after the issuance of the assertion.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Optional:     true,
								Description:  "A list of read-only assertion attributes (for example, SAML_SUBJECT) that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Optional:     true,
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Required:    true,
						Description: "A set of user attributes that the IdP sends in the SAML assertion.",
					},
					"authentication_policy_contract_assertion_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Optional:    true,
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is false.",
								},
								"attribute_contract_fulfillment":     attributecontractfulfillment.ToSchema(true, false),
								"attribute_sources":                  attributesources.ToSchema(0),
								"authentication_policy_contract_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated Authentication Policy Contract."),
								"issuance_criteria":                  issuancecriteria.ToSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:    true,
									Description: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Optional:    true,
						Description: "A list of authentication policy contracts that map to outgoing assertions.",
					},
					"default_target_url": schema.StringAttribute{
						Optional:    true,
						Description: "Default Target URL for SAML1.x connections. For SP connections, this default URL represents the destination on the SP where the user will be directed. For IdP connections, entering a URL in the Default Target URL field overrides the SP Default URL SSO setting.",
					},
					"enabled_profiles": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
					},
					"encryption_policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Optional:    true,
								Description: "Whether the outgoing SAML assertion will be encrypted.",
							},
							"encrypt_slo_subject_name_id": schema.BoolAttribute{
								Optional:    true,
								Description: "Encrypt the name-identifier attribute in outbound SLO messages.  This can be set if the name id is encrypted.",
							},
							"encrypted_attributes": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "The list of outgoing SAML assertion attributes that will be encrypted. The 'encryptAssertion' property takes precedence over this.",
							},
							"slo_subject_name_id_encrypted": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow the encryption of the name-identifier attribute for inbound SLO messages. This can be set if SP initiated SLO is enabled.",
							},
						},
						Optional:    true,
						Description: "Defines what to encrypt in the browser-based SSO profile.",
					},
					"incoming_bindings": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections and it is optional for IdP Connections.",
						Validators: []validator.List{
							listvalidator.UniqueValues(),
						},
					},
					"message_customizations": schema.ListNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Optional:     true,
						Description:  "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"protocol": schema.StringAttribute{
						Required:    true,
						Description: "The browser-based SSO protocol to use.",
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
					"require_signed_authn_requests": schema.BoolAttribute{
						Optional:    true,
						Description: "Require AuthN requests to be signed when received via the POST or Redirect bindings.",
					},
					"sign_assertions": schema.BoolAttribute{
						Optional:    true,
						Description: "Always sign the SAML Assertion.",
					},
					"sign_response_as_required": schema.BoolAttribute{
						Optional:    true,
						Description: "Sign SAML Response as required by the associated binding and encryption policy. Applicable to SAML2.0 only and is defaulted to true. It can be set to false only on SAML2.0 connections when signAssertions is set to true.",
					},
					"slo_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Required:    true,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
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
									Optional:    true,
									Description: "The absolute or relative URL to which logout responses are sent. A relative URL can be specified if a base URL for the connection has been defined.",
								},
								"url": schema.StringAttribute{
									Required:    true,
									Description: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
								},
							},
						},
						Optional:    true,
						Description: "A list of possible endpoints to send SLO requests and responses.",
					},
					"sp_saml_identity_mapping": schema.StringAttribute{
						Optional:    true,
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"PSEUDONYM",
								"STANDARD",
								"TRANSIENT",
							),
						},
					},
					"sp_ws_fed_identity_mapping": schema.StringAttribute{
						Optional:    true,
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP for WS-Federation connection types.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"EMAIL_ADDRESS",
								"USER_PRINCIPLE_NAME",
								"COMMON_NAME",
							),
						},
					},
					"sso_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:    true,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.  Supported bindings are Artifact and POST.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"ARTIFACT",
											"POST",
										),
									},
								},
								"index": schema.Int64Attribute{
									Optional:    true,
									Description: "The priority of the endpoint.",
								},
								"is_default": schema.BoolAttribute{
									Optional:    true,
									Description: "Whether or not this endpoint is the default endpoint. Defaults to false.",
								},
								"url": schema.StringAttribute{
									Required:    true,
									Description: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
								},
							},
						},
						Required:    true,
						Description: "A list of possible endpoints to send assertions to.",
					},
					"url_whitelist_entries": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"allow_query_and_fragment": schema.BoolAttribute{
									Optional:    true,
									Description: "Allow Any Query/Fragment",
								},
								"require_https": schema.BoolAttribute{
									Optional:    true,
									Description: "Require HTTPS",
								},
								"valid_domain": schema.StringAttribute{
									Optional:    true,
									Description: "Valid Domain Name (leading wildcard '*.' allowed)",
								},
								"valid_path": schema.StringAttribute{
									Optional:    true,
									Description: "Valid Path (leave blank to allow any path)",
								},
							},
						},
						Optional:    true,
						Description: "For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled.",
					},
					"ws_fed_token_type": schema.StringAttribute{
						Optional:    true,
						Description: "The WS-Federation Token Type to use.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"SAML11",
								"SAML20",
								"JWT",
							),
						},
					},
					"ws_trust_version": schema.StringAttribute{
						Optional:    true,
						Description: "The WS-Trust version for a WS-Federation connection. The default version is WSTRUST12.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"WSTRUST12",
								"WSTRUST13",
							),
						},
					},
				},
				Optional:    true,
				Description: "The SAML settings used to enable secure browser-based SSO to resources at your partner's site.",
			},
			"type": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Default:     stringdefault.StaticString("SP"),
				Description: "The type of this connection.",
			},
			"virtual_entity_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				Description: "List of alternate entity IDs that identifies the local server to this partner.",
			},
			"ws_trust": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"abort_if_not_fulfilled_from_request": schema.BoolAttribute{
						Optional:    true,
						Description: "If the attribute contract cannot be fulfilled using data from the Request, abort the transaction.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: wsTrustAttribute,
								Optional:     true,
								Description:  "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: wsTrustAttribute,
								Optional:     true,
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Required:    true,
						Description: "A set of user attributes that this server will send in the token.",
					},
					"default_token_type": schema.StringAttribute{
						Optional:    true,
						Description: "The default token type when a web service client (WSC) does not specify in the token request which token type the STS should issue. Defaults to SAML 2.0.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"SAML20",
								"SAML11",
								"SAML11_O365",
							),
						},
					},
					"encrypt_saml2_assertion": schema.BoolAttribute{
						Optional:    true,
						Description: "When selected, the STS encrypts the SAML 2.0 assertion. Applicable only to SAML 2.0 security token.  This option does not apply to OAuth assertion profiles.",
					},
					"generate_key": schema.BoolAttribute{
						Optional:    true,
						Description: "When selected, the STS generates a symmetric key to be used in conjunction with the \"Holder of Key\" (HoK) designation for the assertion's Subject Confirmation Method.  This option does not apply to OAuth assertion profiles.",
					},
					"message_customizations": schema.ListNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Optional:     true,
						Description:  "The message customizations for WS-Trust. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"minutes_after": schema.Int64Attribute{
						Optional:    true,
						Description: "The amount of time after the SAML token was issued during which it is to be considered valid. The default value is 30.",
					},
					"minutes_before": schema.Int64Attribute{
						Optional:    true,
						Description: "The amount of time before the SAML token was issued during which it is to be considered valid. The default value is 5.",
					},
					"oauth_assertion_profiles": schema.BoolAttribute{
						Optional:    true,
						Description: "When selected, four additional token-type requests become available.",
					},
					"partner_service_ids": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "The partner service identifiers.",
					},
					"request_contract_ref": resourcelink.CompleteSingleNestedAttribute(true, false, false, "Request Contract to be used to map attribute values into the security token."),
					"token_processor_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
								"attribute_sources":              attributesources.ToSchema(0),
								"idp_token_processor_ref":        resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated token processor."),
								"issuance_criteria":              issuancecriteria.ToSchema(),
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Required:    true,
						Description: "A list of token processors to validate incoming tokens.",
					},
				},
				Optional:    true,
				Description: "Ws-Trust STS provides security-token validation and creation to extend SSO access to identity-enabled Web Services",
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"connection_id",
		true,
		"The persistent, unique ID for the connection. It can be any combination of [a-zA-Z0-9._-].")
	resp.Schema = schema
}

func (r *idpSpConnectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *idpSpConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil || state == nil {
		return
	}

	if internaltypes.IsDefined(plan.OutboundProvision) && internaltypes.IsDefined(state.OutboundProvision) {
		// If the plan for target_settings has changed, then set target_settings_all to Unknown.
		planAttrs := plan.OutboundProvision.Attributes()
		stateAttrs := state.OutboundProvision.Attributes()
		planTargetSettings := planAttrs["target_settings"].(types.List)
		stateTargetSettings := stateAttrs["target_settings"].(types.List)
		if !planTargetSettings.Equal(stateTargetSettings) {
			planAttrs["target_settings_all"] = types.ListUnknown(targetSettingsElemAttrType)
		}

		// If the plan for channels has changed, then set attribute_mapping_all to Unknown.
		planChannels := planAttrs["channels"].(types.List)
		stateChannels := stateAttrs["channels"].(types.List)
		if !planChannels.Equal(stateChannels) {
			newPlanChannels := []attr.Value{}
			for _, channel := range planChannels.Elements() {
				channelAttrs := channel.(types.Object).Attributes()
				channelAttrs["attribute_mapping_all"] = types.SetUnknown(attributeMappingElemAttrTypes)
				newChannel, respDiags := types.ObjectValue(channelsElemAttrType.AttrTypes, channelAttrs)
				resp.Diagnostics.Append(respDiags...)
				newPlanChannels = append(newPlanChannels, newChannel)
			}
			planAttrs["channels"], respDiags = types.ListValue(channelsElemAttrType, newPlanChannels)
			resp.Diagnostics.Append(respDiags...)
		}

		plan.OutboundProvision, respDiags = types.ObjectValue(outboundProvisionAttrTypes, planAttrs)
		resp.Diagnostics.Append(respDiags...)

		// Update plan
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}

	if internaltypes.IsDefined(state.OutboundProvision) != internaltypes.IsDefined(plan.OutboundProvision) {
		// PF can't add or remove outbound_provision from a sp connection
		resp.RequiresReplace = []path.Path{
			path.Root("outbound_provision"),
		}
	}
}

func addOptionalIdpSpconnectionFields(ctx context.Context, addRequest *client.SpConnection, plan idpSpConnectionResourceModel) error {
	addRequest.Id = plan.ConnectionId.ValueStringPointer()
	addRequest.Type = plan.Type.ValueStringPointer()
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.BaseUrl = plan.BaseUrl.ValueStringPointer()
	addRequest.DefaultVirtualEntityId = plan.DefaultVirtualEntityId.ValueStringPointer()
	addRequest.LicenseConnectionGroup = plan.LicenseConnectionGroup.ValueStringPointer()
	addRequest.LoggingMode = plan.LoggingMode.ValueStringPointer()
	addRequest.ApplicationName = plan.ApplicationName.ValueStringPointer()
	addRequest.ApplicationIconUrl = plan.ApplicationIconUrl.ValueStringPointer()
	addRequest.ConnectionTargetType = plan.ConnectionTargetType.ValueStringPointer()

	if internaltypes.IsDefined(plan.VirtualEntityIds) {
		addRequest.VirtualEntityIds = []string{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.VirtualEntityIds, true)), &addRequest.VirtualEntityIds)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.MetadataReloadSettings) {
		addRequest.MetadataReloadSettings = &client.ConnectionMetadataUrl{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.MetadataReloadSettings, true)), &addRequest.MetadataReloadSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Credentials) {
		addRequest.Credentials = &client.ConnectionCredentials{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Credentials, true)), &addRequest.Credentials)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ContactInfo) {
		addRequest.ContactInfo = &client.ContactInfo{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ContactInfo, true)), &addRequest.ContactInfo)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AdditionalAllowedEntitiesConfiguration) {
		addRequest.AdditionalAllowedEntitiesConfiguration = &client.AdditionalAllowedEntitiesConfiguration{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AdditionalAllowedEntitiesConfiguration, true)), &addRequest.AdditionalAllowedEntitiesConfiguration)
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

	if internaltypes.IsDefined(plan.SpBrowserSso) {
		addRequest.SpBrowserSso = &client.SpBrowserSso{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SpBrowserSso, true)), &addRequest.SpBrowserSso)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeQuery) {
		addRequest.AttributeQuery = &client.SpAttributeQuery{}

		addRequest.AttributeQuery.Attributes = []string{}
		for _, attribute := range plan.AttributeQuery.Attributes()["attributes"].(types.List).Elements() {
			addRequest.AttributeQuery.Attributes = append(addRequest.AttributeQuery.Attributes, attribute.(types.String).ValueString())
		}

		addRequest.AttributeQuery.AttributeContractFulfillment = map[string]client.AttributeFulfillmentValue{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeQuery.Attributes()["attribute_contract_fulfillment"], true)), &addRequest.AttributeQuery.AttributeContractFulfillment)
		if err != nil {
			return err
		}

		addRequest.AttributeQuery.IssuanceCriteria, err = issuancecriteria.ClientStruct(plan.AttributeQuery.Attributes()["issuance_criteria"].(types.Object))
		if err != nil {
			return err
		}

		addRequest.AttributeQuery.Policy = &client.SpAttributeQueryPolicy{}
		err = json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeQuery.Attributes()["policy"], true)), &addRequest.AttributeQuery.Policy)
		if err != nil {
			return err
		}

		addRequest.AttributeQuery.AttributeSources, err = attributesources.ClientStruct(plan.AttributeQuery.Attributes()["attribute_sources"].(types.List))
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.WsTrust) {
		addRequest.WsTrust = &client.SpWsTrust{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.WsTrust, true)), &addRequest.WsTrust)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.OutboundProvision) {
		addRequest.OutboundProvision = &client.OutboundProvision{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.OutboundProvision, true)), &addRequest.OutboundProvision)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *idpSpConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_sp_connection"
}

func (r *idpSpConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readIdpSpconnectionResponse(ctx context.Context, r *client.SpConnection, state *idpSpConnectionResourceModel, plan *idpSpConnectionResourceModel) diag.Diagnostics {
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

func (r *idpSpConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpSpConnectionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdpSpconnection := client.NewSpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	err := addOptionalIdpSpconnectionFields(ctx, createIdpSpconnection, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdP SP Connection", err.Error())
		return
	}

	apiCreateIdpSpconnection := r.apiClient.IdpSpConnectionsAPI.CreateSpConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpSpconnection = apiCreateIdpSpconnection.Body(*createIdpSpconnection)
	idpSpconnectionResponse, httpResp, err := r.apiClient.IdpSpConnectionsAPI.CreateSpConnectionExecute(apiCreateIdpSpconnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdP SP Connection", err, httpResp)
		return
	}

	// Read the response into the state
	var state idpSpConnectionResourceModel

	diags = readIdpSpconnectionResponse(ctx, idpSpconnectionResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpSpConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpSpConnectionResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpSpconnection, httpResp, err := r.apiClient.IdpSpConnectionsAPI.GetSpConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the IdP SP Connection", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  IdP SP Connection", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readIdpSpconnectionResponse(ctx, apiReadIdpSpconnection, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpSpConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan idpSpConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIdpSpconnection := r.apiClient.IdpSpConnectionsAPI.UpdateSpConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ConnectionId.ValueString())
	createUpdateRequest := client.NewSpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	err := addOptionalIdpSpconnectionFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdP SP Connection", err.Error())
		return
	}

	updateIdpSpconnection = updateIdpSpconnection.Body(*createUpdateRequest)
	updateIdpSpconnectionResponse, httpResp, err := r.apiClient.IdpSpConnectionsAPI.UpdateSpConnectionExecute(updateIdpSpconnection)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdP SP Connection", err, httpResp)
		return
	}

	// Read the response
	var state idpSpConnectionResourceModel
	diags = readIdpSpconnectionResponse(ctx, updateIdpSpconnectionResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpSpConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idpSpConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.IdpSpConnectionsAPI.DeleteSpConnection(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the IdP SP Connection", err, httpResp)
	}
}

func (r *idpSpConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to connection_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("connection_id"), req, resp)
}
