// Copyright © 2025 Ping Identity Corporation

package idpspconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/pemcertificates"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpSpConnectionResource{}
	_ resource.ResourceWithConfigure   = &idpSpConnectionResource{}
	_ resource.ResourceWithImportState = &idpSpConnectionResource{}

	customId = "connection_id"
)

var (
	resourceLinkObjectType = types.ObjectType{AttrTypes: resourcelink.AttrType()}

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

	additionalAllowedEntitiesElemType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"entity_id":          types.StringType,
		"entity_description": types.StringType,
	}}

	spBrowserSsoAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":        types.StringType,
			"name_format": types.StringType,
		},
	}
	attributeContractFulfillmentElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrTypes()},
		"value":  types.StringType,
	}}
	attributeContractFulfillmentAttrType = types.MapType{
		ElemType: attributeContractFulfillmentElemAttrType,
	}
	issuanceCriteriaAttrType = types.ObjectType{
		AttrTypes: issuancecriteria.AttrTypes(),
	}
	idpAdapterAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"pseudonym": types.BoolType,
			"masked":    types.BoolType,
		},
	}
	authenticationPolicyContractAssertionMappingsElemType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"authentication_policy_contract_ref": resourceLinkObjectType,
		"restrict_virtual_entity_ids":        types.BoolType,
		"restricted_virtual_entity_ids":      types.SetType{ElemType: types.StringType},
		"abort_sso_transaction_as_fail_safe": types.BoolType,
		"attribute_sources":                  types.SetType{ElemType: types.ObjectType{AttrTypes: attributesources.AttrTypes()}},
		"attribute_contract_fulfillment":     attributeContractFulfillmentAttrType,
		"issuance_criteria":                  issuanceCriteriaAttrType,
	}}
	sloServiceEndpointsElemType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"binding":      types.StringType,
		"url":          types.StringType,
		"response_url": types.StringType,
	}}

	policyAttrTypes = map[string]attr.Type{
		"sign_response":                  types.BoolType,
		"sign_assertion":                 types.BoolType,
		"encrypt_assertion":              types.BoolType,
		"require_signed_attribute_query": types.BoolType,
		"require_encrypted_name_id":      types.BoolType,
	}

	spWsTrustAttributeAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"namespace": types.StringType,
		},
	}
	messageCustomizationsElemType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}}

	channelSourceLocationAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"group_dn":      types.StringType,
			"filter":        types.StringType,
			"nested_search": types.BoolType,
		},
	}
	targetSettingsElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}}

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
		"target_settings":     types.SetType{ElemType: targetSettingsElemAttrType},
		"target_settings_all": types.SetType{ElemType: targetSettingsElemAttrType},
		"custom_schema":       types.ObjectType{AttrTypes: customSchemaAttrTypes},
		"channels":            types.ListType{ElemType: channelsElemAttrType},
	}

	emptyStringSet, _ = types.SetValue(types.StringType, nil)

	groupSourceLocationDefault, _ = types.ObjectValue(channelSourceLocationAttrType.AttrTypes, map[string]attr.Value{
		"filter":        types.StringNull(),
		"group_dn":      types.StringNull(),
		"nested_search": types.BoolValue(false),
	})

	browserSsoAttributeEmptyDefault, _ = types.SetValue(spBrowserSsoAttributeAttrType, nil)
	adapterAttributeEmptyDefault, _    = types.SetValue(idpAdapterAttributeAttrType, nil)
	wsTrustAttributeEmptyDefault, _    = types.SetValue(spWsTrustAttributeAttrType, nil)
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

// GetSchema defines the schema for the resource.
func (r *idpSpConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	adapterOverrideSettingsAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"masked": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"pseudonym": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.",
			},
		},
	}

	spBrowserSSOAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name_format": schema.StringAttribute{
				Optional:    true,
				Description: "The SAML Name Format for the attribute.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}

	wsTrustAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this attribute.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
				Description: "The context in which the customization will be applied. Depending on the connection type and protocol, this can either be `assertion`, `authn-response` or `authn-request`.",
				Validators: []validator.String{
					stringvalidator.OneOf("assertion", "authn-response", "authn-request"),
				},
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			// Defaults can't be set here due to issues related to https://github.com/hashicorp/terraform-plugin-framework/issues/783
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"expression": schema.StringAttribute{
						Optional: true,
						Computed: true,
						//Default:     stringdefault.StaticString(""),
						Description: "An OGNL expression to obtain a value.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
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
						Description: "Indicates how the field shall be parsed. Options are `NONE`, `EXTRACT_CN_FROM_DN`, `EXTRACT_USERNAME_FROM_EMAIL`.",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the configuration field.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"value": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The value for the configuration field. Either this attribute or `encrypted_value` must be specified.",
			},
			"encrypted_value": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The encrypted value for the configuration field. Either this attribute or `value` must be specified.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("value")),
				},
			},
		},
	}

	schema := schema.Schema{
		Version:     1,
		Description: "Manages an IdP SP Connection",
		Attributes: map[string]schema.Attribute{
			"connection_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the connection. It can be any combination of `[a-zA-Z0-9._-]`. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateId(),
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Specifies whether the connection is active and ready to process incoming requests. The default value is `false`.",
			},
			"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"additional_allowed_entities": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"entity_description": schema.StringAttribute{
									Optional:    true,
									Description: "Entity description.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"entity_id": schema.StringAttribute{
									Optional:    true,
									Description: "Unique entity identifier.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(types.SetValueMust(additionalAllowedEntitiesElemType, nil)),
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"application_name": schema.StringAttribute{
				Optional:    true,
				Description: "The application name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"attribute_query": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
					"attribute_sources":              attributesources.ToSchema(1, false),
					"attributes": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "The list of attributes that may be returned to the SP in the response to an attribute request.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"issuance_criteria": issuancecriteria.ToSchema(),
					"policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Encrypt the assertion. The default value is `false`.",
							},
							"require_encrypted_name_id": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Require an encrypted name identifier. The default value is `false`.",
							},
							"require_signed_attribute_query": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Require signed attribute query. The default value is `false`.",
							},
							"sign_assertion": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Sign the assertion. The default value is `false`.",
							},
							"sign_response": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Sign the response. The default value is `false`.",
							},
						},
						Optional: true,
						Computed: true,
						Default: objectdefault.StaticValue(types.ObjectValueMust(policyAttrTypes, map[string]attr.Value{
							"encrypt_assertion":              types.BoolValue(false),
							"require_encrypted_name_id":      types.BoolValue(false),
							"require_signed_attribute_query": types.BoolValue(false),
							"sign_assertion":                 types.BoolValue(false),
							"sign_response":                  types.BoolValue(false),
						})),
						Description: "The attribute query profile's security policy.",
					},
				},
				Optional:    true,
				Description: "The attribute query profile supports SPs in requesting user attributes.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The fully-qualified hostname and port on which your partner's federation deployment runs.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"connection_target_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Description: "The connection target type. This field is intended for bulk import/export usage. Changing its value may result in unexpected behavior. The default value is `STANDARD`. Options are `STANDARD`, `SALESFORCE`, `SALESFORCE_CP`, `SALESFORCE_PP`, `PINGONE_SCIM11`.",
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email": schema.StringAttribute{
						Optional:    true,
						Description: "Contact email address.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"first_name": schema.StringAttribute{
						Optional:    true,
						Description: "Contact first name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"last_name": schema.StringAttribute{
						Optional:    true,
						Description: "Contact last name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"phone": schema.StringAttribute{
						Optional:    true,
						Description: "Contact phone number.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
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
				Description: "The time at which the connection was created. This property is read only.",
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
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Determines whether the signing certificate is included in the signature <KeyInfo> element. The default value is `false`.",
								MarkdownDescription: "Determines whether the signing certificate is included in the signature <KeyInfo> element. The default value is `false`.",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								Description:         "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element. The default value is `false`.",
								MarkdownDescription: "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element. The default value is `false`.",
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
										Optional:            true,
										Computed:            true,
										Description:         "Encrypted user password. Either this attribute or `password` must be specified.",
										MarkdownDescription: "Encrypted user password. Either this attribute or `password` must be specified.",
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
										Optional:            true,
										Computed:            true,
										Description:         "Encrypted user password. Either this attribute or `password` must be specified.",
										MarkdownDescription: "Encrypted user password. Either this attribute or `password` must be specified.",
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
								Description:         "If incoming or outgoing messages must be signed.",
								MarkdownDescription: "If incoming or outgoing messages must be signed.",
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
				Optional:    true,
				Description: "The default alternate entity ID that identifies the local server to this partner. It is required when `virtual_entity_ids` is not empty and must be included in that list.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"entity_id": schema.StringAttribute{
				Required:    true,
				Description: "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
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
							Computed:    true,
							Default:     setdefault.StaticValue(emptyStringSet),
							Description: "A List of values",
						},
					},
				},
				Optional:    true,
				Description: "Extended Properties allows to store additional information for IdP/SP Connections. The names of these extended properties should be defined in the `pingfederate_extended_properties` resource.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The persistent, unique ID for the connection. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateId(),
				},
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
				Description: "The level of transaction logging applicable for this connection. Default is `STANDARD`. Options are `NONE`, `STANDARD`, `ENHANCED`, `FULL`. If the `sp_connection_transaction_logging_override` attribute is set to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource, then this attribute must be set to the same value.",
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
						Description: "Specifies whether the metadata of the connection will be automatically reloaded. The default value is `true`.",
					},
					"metadata_url_ref": resourcelink.SingleNestedAttribute(),
				},
				Optional:    true,
				Description: "Configuration settings to enable automatic reload of partner's metadata.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The connection name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
													Description: "The account status algorithm name. Options are `ACCOUNT_STATUS_ALGORITHM_AD`, `ACCOUNT_STATUS_ALGORITHM_FLAG`. `ACCOUNT_STATUS_ALGORITHM_AD` -  Algorithm name for Active Directory, which uses a bitmap for each user entry. `ACCOUNT_STATUS_ALGORITHM_FLAG` - Algorithm name for Oracle Directory Server and other LDAP directories that use a separate attribute to store the user's status. When this option is selected, the Flag Comparison Value and Flag Comparison Status fields should be used.",
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"default_status": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(true),
													Description: "The default status of the account. The default value is `true`.",
												},
												"flag_comparison_status": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(true),
													Description: "The flag that represents comparison status. The default value is `true`.",
												},
												"flag_comparison_value": schema.StringAttribute{
													Optional:    true,
													Description: "The flag that represents comparison value.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
											},
											Required:    true,
											Description: "Account management settings.",
										},
										"base_dn": schema.StringAttribute{
											Required:    true,
											Description: "The base DN where the user records are located.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"change_detection_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"changed_users_algorithm": schema.StringAttribute{
													Required:    true,
													Description: "The changed user algorithm. Options are `ACTIVE_DIRECTORY_USN`, `TIMESTAMP`, `TIMESTAMP_NO_NEGATION`. `ACTIVE_DIRECTORY_USN` - For Active Directory only, this algorithm queries for update sequence numbers on user records that are larger than the last time records were checked. `TIMESTAMP` - Queries for timestamps on user records that are not older than the last time records were checked. This check is more efficient from the point of view of the PingFederate provisioner but can be more time consuming on the LDAP side, particularly with the Oracle Directory Server. `TIMESTAMP_NO_NEGATION` - Queries for timestamps on user records that are newer than the last time records were checked. This algorithm is recommended for the Oracle Directory Server.",
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"time_stamp_attribute_name": schema.StringAttribute{
													Required:    true,
													Description: "The timestamp attribute name.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"user_object_class": schema.StringAttribute{
													Required:    true,
													Description: "The user object class.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"usn_attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The USN attribute name.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
											},
											Required:    true,
											Description: "Setting to detect changes to a user or a group.",
										},
										"data_source": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to an LDAP datastore."),
										"group_membership_detection": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"group_member_attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The name of the attribute that represents group members in a group, also known as group member attribute.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"member_of_group_attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The name of the attribute that indicates the entity is a member of a group, also known as member of attribute.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"group_dn": schema.StringAttribute{
													Optional:    true,
													Description: "The group DN for users or groups.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"nested_search": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
													Description: "Indicates whether the search is nested. The default value is `false`.",
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
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"group_dn": schema.StringAttribute{
													Optional:    true,
													Description: "The group DN for users or groups.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"nested_search": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
													Description: "Indicates whether the search is nested. Default is `false`.",
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
									Description: "The number of processing threads. The default value is `1`.",
								},
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of the channel.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"timeout": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Default:     int64default.StaticInt64(60),
									Description: "Timeout, in seconds, for individual user and group provisioning operations on the target service provider. The default value is `60`.",
								},
							},
						},
						Required:    true,
						Description: "Includes settings of a source data store, managing provisioning threads and mapping of attributes.",
					},
					"custom_schema": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attributes": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"multi_valued": schema.BoolAttribute{
											Optional:    true,
											Description: "Indicates whether the attribute is multi-valued.",
										},
										"name": schema.StringAttribute{
											Optional:    true,
											Description: "Name of the attribute.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"sub_attributes": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Computed:    true,
											Default:     setdefault.StaticValue(emptyStringSet),
											Description: "List of sub-attributes for an attribute.",
										},
										"types": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Computed:    true,
											Default:     setdefault.StaticValue(emptyStringSet),
											Description: "Represents the name of each attribute type in case of multi-valued attribute.",
										},
									},
								},
								Optional: true,
								Computed: true,
								Default:  setdefault.StaticValue(types.SetValueMust(attributesElemType, nil)),
							},
							"namespace": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
						Optional:    true,
						Description: "Custom SCIM Attributes configuration.",
					},
					"target_settings_all": schema.SetNestedAttribute{
						NestedObject: outboundProvisionTargetSettingsNestedObject,
						Optional:     false,
						Computed:     true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						Description: "Configuration fields that includes credentials to target SaaS application. This attribute will include any values set by default by PingFederate.",
					},
					"target_settings": schema.SetNestedAttribute{
						NestedObject: outboundProvisionTargetSettingsNestedObject,
						Required:     true,
						Description:  "Configuration fields that includes credentials to target SaaS application.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The SaaS plugin type.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:    true,
				Description: "Outbound Provisioning allows an IdP to create and maintain user accounts at standards-based partner sites using SCIM as well as select-proprietary provisioning partner sites that are protocol-enabled.",
			},
			"sp_browser_sso": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"adapter_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is `false`.",
								},
								"adapter_override_settings": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"attribute_contract": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"core_attributes": schema.SetNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Required:     true,
													Description:  "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
												},
												"extended_attributes": schema.SetNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Optional:     true,
													Computed:     true,
													Default:      setdefault.StaticValue(adapterAttributeEmptyDefault),
													Description:  "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
												},
												"mask_ognl_values": schema.BoolAttribute{
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
													Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to `false`.",
												},
												"unique_user_key_attribute": schema.StringAttribute{
													Optional:    true,
													Description: "The attribute to use for uniquely identify a user's authentication sessions.",
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
											},
											Optional:    true,
											Description: "A set of attributes exposed by an IdP adapter.",
										},
										"attribute_mapping": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
												"attribute_sources":              attributesources.ToSchema(0, false),
												"issuance_criteria":              issuancecriteria.ToSchema(),
											},
											Optional:    true,
											Description: "An IdP Adapter Contract Mapping.",
										},
										"authn_ctx_class_ref": schema.StringAttribute{
											Optional:    true,
											Description: "The fixed value that indicates how the user was authenticated.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"configuration": pluginconfiguration.ToSchema(),
										"id": schema.StringAttribute{
											Required:    true,
											Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
										},
										"name": schema.StringAttribute{
											Required:    true,
											Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
										"parent_ref":            resourcelink.CompleteSingleNestedAttribute(true, false, false, "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)"),
										"plugin_descriptor_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override."),
									},
									Optional: true,
								},
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
								"attribute_sources":              attributesources.ToSchema(0, false),
								"idp_adapter_ref":                resourcelink.CompleteSingleNestedAttribute(true, false, false, "Reference to the associated IdP adapter. Note: This is ignored if adapter overrides for this mapping exists. In this case, the override's parent adapter reference is used."),
								"issuance_criteria":              issuancecriteria.ToSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
								},
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
									Default:     setdefault.StaticValue(emptyStringSet),
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Required:    true,
						Description: "A list of adapters that map to outgoing assertions.",
					},
					"always_sign_artifact_response": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Specify to always sign the SAML ArtifactResponse. The default value is `false`.",
					},
					"artifact": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"lifetime": schema.Int64Attribute{
								Optional:    true,
								Description: "The lifetime of the artifact in seconds.",
							},
							"resolver_locations": schema.SetNestedAttribute{
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
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
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
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Optional:     true,
								Description:  "A list of read-only assertion attributes (for example, SAML_SUBJECT) that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.SetNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Optional:     true,
								Computed:     true,
								Default:      setdefault.StaticValue(browserSsoAttributeEmptyDefault),
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Required:    true,
						Description: "A set of user attributes that the IdP sends in the SAML assertion.",
					},
					"authentication_policy_contract_assertion_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is `false`.",
								},
								"attribute_contract_fulfillment":     attributecontractfulfillment.ToSchema(true, false, false),
								"attribute_sources":                  attributesources.ToSchema(0, false),
								"authentication_policy_contract_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated Authentication Policy Contract."),
								"issuance_criteria":                  issuancecriteria.ToSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "Restricts this mapping to specific virtual entity IDs. The default value is `false`.",
								},
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
									Default:     setdefault.StaticValue(emptyStringSet),
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(types.SetValueMust(authenticationPolicyContractAssertionMappingsElemType, nil)),
						Description: "A list of authentication policy contracts that map to outgoing assertions.",
					},
					"default_target_url": schema.StringAttribute{
						Optional:    true,
						Description: "Default Target URL for SAML1.x connections. This default URL represents the destination on the SP where the user will be directed.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"enabled_profiles": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections. ",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"encryption_policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Whether the outgoing SAML assertion will be encrypted. The default value is `false`.",
							},
							"encrypt_slo_subject_name_id": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Encrypt the name-identifier attribute in outbound SLO messages. This can be set if the name id is encrypted. The default value is `false`.",
							},
							"encrypted_attributes": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
								Default:     setdefault.StaticValue(emptyStringSet),
								Description: "The list of outgoing SAML assertion attributes that will be encrypted. The `encrypt_assertion` property takes precedence over this.",
							},
							"slo_subject_name_id_encrypted": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
								Description: "Allow the encryption of the name-identifier attribute for inbound SLO messages. This can be set if SP initiated SLO is enabled. The default value is `false`.",
							},
						},
						Optional:    true,
						Description: "Defines what to encrypt in the browser-based SSO profile.",
					},
					"incoming_bindings": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"message_customizations": schema.SetNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Optional:     true,
						Computed:     true,
						Default:      setdefault.StaticValue(types.SetValueMust(messageCustomizationsElemType, nil)),
						Description:  "The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"protocol": schema.StringAttribute{
						Required:    true,
						Description: "The browser-based SSO protocol to use. Options are `SAML20`, `WSFED`, `SAML11`, `SAML10`, `OIDC`.",
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
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Require AuthN requests to be signed when received via the POST or Redirect bindings. The default value is `false`.",
					},
					"sign_assertions": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Always sign the SAML Assertion. The default value is `false`.",
					},
					"sign_response_as_required": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Sign SAML Response as required by the associated binding and encryption policy. Applicable to SAML2.0 only and is defaulted to `true`. It can be set to `false` only on SAML2.0 connections when `sign_assertions` is set to `true`.",
					},
					"slo_service_endpoints": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:    true,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.",
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
						Computed:    true,
						Default:     setdefault.StaticValue(types.SetValueMust(sloServiceEndpointsElemType, nil)),
						Description: "A list of possible endpoints to send SLO requests and responses.",
					},
					"sp_saml_identity_mapping": schema.StringAttribute{
						Optional:    true,
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP. Options are `PSEUDONYM`, `STANDARD`, `TRANSIENT`.",
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
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP for WS-Federation connection types. Options are `EMAIL_ADDRESS`, `USER_PRINCIPLE_NAME`, `COMMON_NAME`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"EMAIL_ADDRESS",
								"USER_PRINCIPLE_NAME",
								"COMMON_NAME",
							),
						},
					},
					"sso_service_endpoints": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Optional:    true,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`.",
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
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "Whether or not this endpoint is the default endpoint. Defaults to `false`.",
								},
								"url": schema.StringAttribute{
									Required:    true,
									Description: "The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Required:    true,
						Description: "A list of possible endpoints to send assertions to.",
					},
					"url_whitelist_entries": schema.SetNestedAttribute{
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"valid_path": schema.StringAttribute{
									Optional:    true,
									Description: "Valid Path (leave undefined to allow any path)",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						Optional:    true,
						Description: "For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled.",
					},
					"ws_fed_token_type": schema.StringAttribute{
						Optional:    true,
						Description: "The WS-Federation Token Type to use. Options are `SAML11`, `SAML20`, `JWT`.",
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
						Computed:    true,
						Description: "The WS-Trust version for a WS-Federation connection. The default version is `WSTRUST12`. Options are `WSTRUST12`, `WSTRUST13`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"WSTRUST12",
								"WSTRUST13",
							),
						},
					},
					"sso_application_endpoint": schema.StringAttribute{
						Optional:    false,
						Computed:    true,
						Description: "Application endpoint that can be used to invoke single sign-on (SSO) for the connection. This is a read-only parameter.",
					},
				},
				Optional:    true,
				Description: "The SAML settings used to enable secure browser-based SSO to resources at your partner's site.",
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
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: wsTrustAttribute,
								Optional:     true,
								Description:  "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.SetNestedAttribute{
								NestedObject: wsTrustAttribute,
								Optional:     true,
								Computed:     true,
								Default:      setdefault.StaticValue(wsTrustAttributeEmptyDefault),
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Required:    true,
						Description: "A set of user attributes that this server will send in the token.",
					},
					"default_token_type": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("SAML20"),
						Description: "The default token type when a web service client (WSC) does not specify in the token request which token type the STS should issue. Options are `SAML20`, `SAML11`, `SAML11_O365`. Defaults to `SAML20`.",
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
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "When selected, the STS encrypts the SAML 2.0 assertion. Applicable only to SAML 2.0 security token.  This option does not apply to OAuth assertion profiles. The default value is `false`.",
					},
					"generate_key": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "When selected, the STS generates a symmetric key to be used in conjunction with the \"Holder of Key\" (HoK) designation for the assertion's Subject Confirmation Method.  This option does not apply to OAuth assertion profiles. The default value is `false`.",
					},
					"message_customizations": schema.SetNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Optional:     true,
						Computed:     true,
						Default:      setdefault.StaticValue(types.SetValueMust(messageCustomizationsElemType, nil)),
						Description:  "The message customizations for WS-Trust. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"minutes_after": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(30),
						Description: "The amount of time after the SAML token was issued during which it is to be considered valid. The default value is `30`.",
					},
					"minutes_before": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: "The amount of time before the SAML token was issued during which it is to be considered valid. The default value is `5`.",
					},
					"oauth_assertion_profiles": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "When selected, four additional token-type requests become available. The default value is `false`.",
					},
					"partner_service_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "The partner service identifiers.",
					},
					"request_contract_ref": resourcelink.CompleteSingleNestedAttribute(true, false, false, "Request Contract to be used to map attribute values into the security token."),
					"token_processor_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
								"attribute_sources":              attributesources.ToSchema(0, false),
								"idp_token_processor_ref":        resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated token processor."),
								"issuance_criteria":              issuancecriteria.ToSchema(),
								"restricted_virtual_entity_ids": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
									Default:     setdefault.StaticValue(emptyStringSet),
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
	resp.Schema = schema
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

func (r *idpSpConnectionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config *idpSpConnectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if config == nil {
		return
	}

	virtualIds := config.VirtualEntityIds.Elements()
	if internaltypes.IsDefined(config.DefaultVirtualEntityId) && !config.VirtualEntityIds.IsUnknown() {
		defaultId := config.DefaultVirtualEntityId.ValueString()
		found := false
		anyUnknown := false
		for _, id := range virtualIds {
			if id.IsUnknown() {
				anyUnknown = true
				break
			}
			if defaultId == id.(types.String).ValueString() {
				found = true
				break
			}
		}
		if !found && !anyUnknown {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_virtual_entity_id"),
				providererror.InvalidAttributeConfiguration,
				"The value provided for 'default_virtual_entity_id' must be included in the 'virtual_entity_ids' list. "+
					fmt.Sprintf("The value '%s' is not included in the 'virtual_entity_ids' list.", defaultId))
		}
	} else if len(virtualIds) > 0 && config.DefaultVirtualEntityId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("default_virtual_entity_id"),
			providererror.InvalidAttributeConfiguration,
			"The 'default_virtual_entity_id' attribute must be set when 'virtual_entity_ids' is non-empty.")
	}

	if internaltypes.IsDefined(config.SpBrowserSso) {
		encryptionPolicy := config.SpBrowserSso.Attributes()["encryption_policy"].(types.Object)
		if internaltypes.IsDefined(encryptionPolicy) {
			encryptAssertion := encryptionPolicy.Attributes()["encrypt_assertion"].(types.Bool)
			encryptionAttributes := encryptionPolicy.Attributes()["encrypted_attributes"].(types.Set)
			if encryptAssertion.ValueBool() && len(encryptionAttributes.Elements()) > 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("sp_browser_sso").AtMapKey("encryption_policy").AtMapKey("encrypted_attributes"),
					providererror.InvalidAttributeConfiguration,
					"The 'encrypted_attributes' attribute cannot be configured when 'encrypt_assertion' is set to true.")
			}
		}
	}
}

func (r *idpSpConnectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *idpSpConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil {
		return
	}

	planSpBrowserSsoAttributes := plan.SpBrowserSso.Attributes()
	if internaltypes.IsDefined(plan.SpBrowserSso) {
		// Get the protocol
		spBrowserSsoPlanModified := false
		protocol := planSpBrowserSsoAttributes["protocol"].(types.String).ValueString()
		switch protocol {
		case "SAML20":
			if planSpBrowserSsoAttributes["sign_response_as_required"].IsUnknown() {
				planSpBrowserSsoAttributes["sign_response_as_required"] = types.BoolValue(true)
				spBrowserSsoPlanModified = true
			}
		case "WSFED":
			if planSpBrowserSsoAttributes["ws_trust_version"].IsUnknown() {
				planSpBrowserSsoAttributes["ws_trust_version"] = types.StringValue("WSTRUST12")
				spBrowserSsoPlanModified = true
			}
		}
		if planSpBrowserSsoAttributes["sign_response_as_required"].IsUnknown() {
			planSpBrowserSsoAttributes["sign_response_as_required"] = types.BoolNull()
			spBrowserSsoPlanModified = true
		}
		if planSpBrowserSsoAttributes["ws_trust_version"].IsUnknown() {
			planSpBrowserSsoAttributes["ws_trust_version"] = types.StringNull()
			spBrowserSsoPlanModified = true
		}

		if spBrowserSsoPlanModified {
			// Update the plan for sp_browser_sso
			plan.SpBrowserSso, respDiags = types.ObjectValue(plan.SpBrowserSso.AttributeTypes(ctx), planSpBrowserSsoAttributes)
			resp.Diagnostics.Append(respDiags...)
			// Update plan
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		}
	}

	if state == nil {
		return
	}

	if internaltypes.IsDefined(plan.OutboundProvision) && internaltypes.IsDefined(state.OutboundProvision) {
		// If the plan for target_settings has changed, then set target_settings_all to Unknown.
		planAttrs := plan.OutboundProvision.Attributes()
		stateAttrs := state.OutboundProvision.Attributes()
		planTargetSettings := planAttrs["target_settings"].(types.Set)
		stateTargetSettings := stateAttrs["target_settings"].(types.Set)
		if !planTargetSettings.Equal(stateTargetSettings) {
			planAttrs["target_settings_all"] = types.SetUnknown(targetSettingsElemAttrType)
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

	if (state.OutboundProvision.IsNull() && internaltypes.IsDefined(plan.OutboundProvision)) ||
		(internaltypes.IsDefined(state.OutboundProvision) && plan.OutboundProvision.IsNull()) {
		// PF can't add or remove outbound_provision from a sp connection
		resp.RequiresReplace = []path.Path{
			path.Root("outbound_provision"),
		}
	}

	if internaltypes.IsDefined(plan.SpBrowserSso) && internaltypes.IsDefined(state.SpBrowserSso) && !plan.SpBrowserSso.Equal(state.SpBrowserSso) {
		stateSpBrowserSsoAttributes := state.SpBrowserSso.Attributes()
		spBrowserSsoPlanModified := false
		if internaltypes.IsDefined(planSpBrowserSsoAttributes["adapter_mappings"]) && internaltypes.IsDefined(stateSpBrowserSsoAttributes["adapter_mappings"]) {
			planAdapterMappings := planSpBrowserSsoAttributes["adapter_mappings"].(types.Set)
			stateAdapterMappings := stateSpBrowserSsoAttributes["adapter_mappings"].(types.Set)
			if !planAdapterMappings.Equal(stateAdapterMappings) {
				planAdapterMappingsElems := planAdapterMappings.Elements()
				stateAdapterMappingsElems := stateAdapterMappings.Elements()
				finalPlanAdapterMappingsElems := []attr.Value{}
				for i, planElem := range planAdapterMappingsElems {
					planElemObj := planElem.(types.Object)
					if i < len(stateAdapterMappingsElems) {
						stateElem := stateAdapterMappingsElems[i]
						// If this plan element doesn't equal the state element, invalidate the fields_all and tables_all fields in the configuration
						if !planElem.Equal(stateElem) {
							planElemAttrs := planElemObj.Attributes()
							if internaltypes.IsDefined(planElemAttrs["adapter_override_settings"]) {
								planOverrideSettingsObj := planElemAttrs["adapter_override_settings"].(types.Object)
								planOverrideSettingsAttrs := planOverrideSettingsObj.Attributes()
								if internaltypes.IsDefined(planOverrideSettingsAttrs["configuration"]) {
									planOverrideSettingsAttrs["configuration"], respDiags = pluginconfiguration.MarkComputedAttrsUnknown(
										planOverrideSettingsAttrs["configuration"].(types.Object))
									resp.Diagnostics.Append(respDiags...)
									planElemAttrs["adapter_override_settings"], respDiags = types.ObjectValue(planOverrideSettingsObj.AttributeTypes(ctx), planOverrideSettingsAttrs)
									resp.Diagnostics.Append(respDiags...)
									planElemObj, respDiags = types.ObjectValue(planElemObj.AttributeTypes(ctx), planElemAttrs)
									resp.Diagnostics.Append(respDiags...)
									spBrowserSsoPlanModified = true
									finalPlanAdapterMappingsElems = append(finalPlanAdapterMappingsElems, planElemObj)
								}
							}
						}
					}
				}
				if spBrowserSsoPlanModified {
					planSpBrowserSsoAttributes["adapter_mappings"], respDiags = types.SetValue(planAdapterMappings.ElementType(ctx), finalPlanAdapterMappingsElems)
					resp.Diagnostics.Append(respDiags...)
				}
			}
		}
		if spBrowserSsoPlanModified {
			// Update the plan for sp_browser_sso
			plan.SpBrowserSso, respDiags = types.ObjectValue(plan.SpBrowserSso.AttributeTypes(ctx), planSpBrowserSsoAttributes)
			resp.Diagnostics.Append(respDiags...)
			// Update plan
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		}
	}
}

func addOptionalIdpSpconnectionFields(_ context.Context, addRequest *client.SpConnection, plan idpSpConnectionModel) diag.Diagnostics {
	var respDiags diag.Diagnostics
	addRequest.Id = plan.ConnectionId.ValueStringPointer()
	addRequest.Type = utils.Pointer("SP")
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.BaseUrl = plan.BaseUrl.ValueStringPointer()
	addRequest.DefaultVirtualEntityId = plan.DefaultVirtualEntityId.ValueStringPointer()
	addRequest.LicenseConnectionGroup = plan.LicenseConnectionGroup.ValueStringPointer()
	addRequest.LoggingMode = plan.LoggingMode.ValueStringPointer()
	addRequest.ApplicationName = plan.ApplicationName.ValueStringPointer()
	addRequest.ApplicationIconUrl = plan.ApplicationIconUrl.ValueStringPointer()
	addRequest.ConnectionTargetType = plan.ConnectionTargetType.ValueStringPointer()

	// virtual_entity_ids
	if !plan.VirtualEntityIds.IsNull() {
		addRequest.VirtualEntityIds = []string{}
		for _, virtualEntityIdsElement := range plan.VirtualEntityIds.Elements() {
			addRequest.VirtualEntityIds = append(addRequest.VirtualEntityIds, virtualEntityIdsElement.(types.String).ValueString())
		}
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
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Credentials, true)), &addRequest.Credentials)
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
	if !plan.AdditionalAllowedEntitiesConfiguration.IsNull() {
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

	// sp_browser_sso
	if !plan.SpBrowserSso.IsNull() {
		spBrowserSsoValue := &client.SpBrowserSso{}
		spBrowserSsoAttrs := plan.SpBrowserSso.Attributes()
		spBrowserSsoValue.AdapterMappings = []client.IdpAdapterAssertionMapping{}
		for _, adapterMappingsElement := range spBrowserSsoAttrs["adapter_mappings"].(types.Set).Elements() {
			adapterMappingsValue := client.IdpAdapterAssertionMapping{}
			adapterMappingsAttrs := adapterMappingsElement.(types.Object).Attributes()
			adapterMappingsValue.AbortSsoTransactionAsFailSafe = adapterMappingsAttrs["abort_sso_transaction_as_fail_safe"].(types.Bool).ValueBoolPointer()
			if !adapterMappingsAttrs["adapter_override_settings"].IsNull() {
				adapterMappingsAdapterOverrideSettingsValue := &client.IdpAdapter{}
				adapterMappingsAdapterOverrideSettingsAttrs := adapterMappingsAttrs["adapter_override_settings"].(types.Object).Attributes()
				if !adapterMappingsAdapterOverrideSettingsAttrs["attribute_contract"].IsNull() {
					adapterMappingsAdapterOverrideSettingsAttributeContractValue := &client.IdpAdapterAttributeContract{}
					adapterMappingsAdapterOverrideSettingsAttributeContractAttrs := adapterMappingsAdapterOverrideSettingsAttrs["attribute_contract"].(types.Object).Attributes()
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes = []client.IdpAdapterAttribute{}
					for _, coreAttributesElement := range adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
						coreAttributesValue := client.IdpAdapterAttribute{}
						coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
						coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
						coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
						coreAttributesValue.Pseudonym = coreAttributesAttrs["pseudonym"].(types.Bool).ValueBoolPointer()
						adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes = append(adapterMappingsAdapterOverrideSettingsAttributeContractValue.CoreAttributes, coreAttributesValue)
					}
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes = []client.IdpAdapterAttribute{}
					for _, extendedAttributesElement := range adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
						extendedAttributesValue := client.IdpAdapterAttribute{}
						extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
						extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
						extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
						extendedAttributesValue.Pseudonym = extendedAttributesAttrs["pseudonym"].(types.Bool).ValueBoolPointer()
						adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes = append(adapterMappingsAdapterOverrideSettingsAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
					}
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.MaskOgnlValues = adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["mask_ognl_values"].(types.Bool).ValueBoolPointer()
					adapterMappingsAdapterOverrideSettingsAttributeContractValue.UniqueUserKeyAttribute = adapterMappingsAdapterOverrideSettingsAttributeContractAttrs["unique_user_key_attribute"].(types.String).ValueStringPointer()
					adapterMappingsAdapterOverrideSettingsValue.AttributeContract = adapterMappingsAdapterOverrideSettingsAttributeContractValue
				}
				if !adapterMappingsAdapterOverrideSettingsAttrs["attribute_mapping"].IsNull() {
					adapterMappingsAdapterOverrideSettingsAttributeMappingValue := &client.IdpAdapterContractMapping{}
					adapterMappingsAdapterOverrideSettingsAttributeMappingAttrs := adapterMappingsAdapterOverrideSettingsAttrs["attribute_mapping"].(types.Object).Attributes()
					adapterMappingsAdapterOverrideSettingsAttributeMappingValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(adapterMappingsAdapterOverrideSettingsAttributeMappingAttrs["attribute_contract_fulfillment"].(types.Map))
					adapterMappingsAdapterOverrideSettingsAttributeMappingValue.AttributeSources = attributesources.ClientStruct(adapterMappingsAdapterOverrideSettingsAttributeMappingAttrs["attribute_sources"].(types.Set))
					adapterMappingsAdapterOverrideSettingsAttributeMappingValue.IssuanceCriteria = issuancecriteria.ClientStruct(adapterMappingsAdapterOverrideSettingsAttributeMappingAttrs["issuance_criteria"].(types.Object))
					adapterMappingsAdapterOverrideSettingsValue.AttributeMapping = adapterMappingsAdapterOverrideSettingsAttributeMappingValue
				}
				adapterMappingsAdapterOverrideSettingsValue.AuthnCtxClassRef = adapterMappingsAdapterOverrideSettingsAttrs["authn_ctx_class_ref"].(types.String).ValueStringPointer()
				adapterMappingsAdapterOverrideSettingsConfigurationValue := pluginconfiguration.ClientStruct(adapterMappingsAdapterOverrideSettingsAttrs["configuration"].(types.Object))
				adapterMappingsAdapterOverrideSettingsValue.Configuration = *adapterMappingsAdapterOverrideSettingsConfigurationValue
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
				adapterMappingsValue.AdapterOverrideSettings = adapterMappingsAdapterOverrideSettingsValue
			}
			adapterMappingsValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(adapterMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			adapterMappingsValue.AttributeSources = attributesources.ClientStruct(adapterMappingsAttrs["attribute_sources"].(types.Set))
			if !adapterMappingsAttrs["idp_adapter_ref"].IsNull() {
				adapterMappingsIdpAdapterRefValue := &client.ResourceLink{}
				adapterMappingsIdpAdapterRefAttrs := adapterMappingsAttrs["idp_adapter_ref"].(types.Object).Attributes()
				adapterMappingsIdpAdapterRefValue.Id = adapterMappingsIdpAdapterRefAttrs["id"].(types.String).ValueString()
				adapterMappingsValue.IdpAdapterRef = adapterMappingsIdpAdapterRefValue
			}
			adapterMappingsValue.IssuanceCriteria = issuancecriteria.ClientStruct(adapterMappingsAttrs["issuance_criteria"].(types.Object))
			adapterMappingsValue.RestrictVirtualEntityIds = adapterMappingsAttrs["restrict_virtual_entity_ids"].(types.Bool).ValueBoolPointer()
			if !adapterMappingsAttrs["restricted_virtual_entity_ids"].IsNull() {
				adapterMappingsValue.RestrictedVirtualEntityIds = []string{}
				for _, restrictedVirtualEntityIdsElement := range adapterMappingsAttrs["restricted_virtual_entity_ids"].(types.Set).Elements() {
					adapterMappingsValue.RestrictedVirtualEntityIds = append(adapterMappingsValue.RestrictedVirtualEntityIds, restrictedVirtualEntityIdsElement.(types.String).ValueString())
				}
			}
			spBrowserSsoValue.AdapterMappings = append(spBrowserSsoValue.AdapterMappings, adapterMappingsValue)
		}
		spBrowserSsoValue.AlwaysSignArtifactResponse = spBrowserSsoAttrs["always_sign_artifact_response"].(types.Bool).ValueBoolPointer()
		if !spBrowserSsoAttrs["artifact"].IsNull() {
			spBrowserSsoArtifactValue := &client.ArtifactSettings{}
			spBrowserSsoArtifactAttrs := spBrowserSsoAttrs["artifact"].(types.Object).Attributes()
			spBrowserSsoArtifactValue.Lifetime = spBrowserSsoArtifactAttrs["lifetime"].(types.Int64).ValueInt64Pointer()
			spBrowserSsoArtifactValue.ResolverLocations = []client.ArtifactResolverLocation{}
			for _, resolverLocationsElement := range spBrowserSsoArtifactAttrs["resolver_locations"].(types.Set).Elements() {
				resolverLocationsValue := client.ArtifactResolverLocation{}
				resolverLocationsAttrs := resolverLocationsElement.(types.Object).Attributes()
				resolverLocationsValue.Index = resolverLocationsAttrs["index"].(types.Int64).ValueInt64()
				resolverLocationsValue.Url = resolverLocationsAttrs["url"].(types.String).ValueString()
				spBrowserSsoArtifactValue.ResolverLocations = append(spBrowserSsoArtifactValue.ResolverLocations, resolverLocationsValue)
			}
			spBrowserSsoArtifactValue.SourceId = spBrowserSsoArtifactAttrs["source_id"].(types.String).ValueStringPointer()
			spBrowserSsoValue.Artifact = spBrowserSsoArtifactValue
		}
		spBrowserSsoAssertionLifetimeValue := client.AssertionLifetime{}
		spBrowserSsoAssertionLifetimeAttrs := spBrowserSsoAttrs["assertion_lifetime"].(types.Object).Attributes()
		spBrowserSsoAssertionLifetimeValue.MinutesAfter = spBrowserSsoAssertionLifetimeAttrs["minutes_after"].(types.Int64).ValueInt64()
		spBrowserSsoAssertionLifetimeValue.MinutesBefore = spBrowserSsoAssertionLifetimeAttrs["minutes_before"].(types.Int64).ValueInt64()
		spBrowserSsoValue.AssertionLifetime = spBrowserSsoAssertionLifetimeValue
		spBrowserSsoAttributeContractValue := client.SpBrowserSsoAttributeContract{}
		spBrowserSsoAttributeContractAttrs := spBrowserSsoAttrs["attribute_contract"].(types.Object).Attributes()
		spBrowserSsoAttributeContractValue.CoreAttributes = []client.SpBrowserSsoAttribute{}
		for _, coreAttributesElement := range spBrowserSsoAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.SpBrowserSsoAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			coreAttributesValue.NameFormat = coreAttributesAttrs["name_format"].(types.String).ValueStringPointer()
			spBrowserSsoAttributeContractValue.CoreAttributes = append(spBrowserSsoAttributeContractValue.CoreAttributes, coreAttributesValue)
		}
		spBrowserSsoAttributeContractValue.ExtendedAttributes = []client.SpBrowserSsoAttribute{}
		for _, extendedAttributesElement := range spBrowserSsoAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.SpBrowserSsoAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			extendedAttributesValue.NameFormat = extendedAttributesAttrs["name_format"].(types.String).ValueStringPointer()
			spBrowserSsoAttributeContractValue.ExtendedAttributes = append(spBrowserSsoAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		spBrowserSsoValue.AttributeContract = spBrowserSsoAttributeContractValue
		spBrowserSsoValue.AuthenticationPolicyContractAssertionMappings = []client.AuthenticationPolicyContractAssertionMapping{}
		for _, authenticationPolicyContractAssertionMappingsElement := range spBrowserSsoAttrs["authentication_policy_contract_assertion_mappings"].(types.Set).Elements() {
			authenticationPolicyContractAssertionMappingsValue := client.AuthenticationPolicyContractAssertionMapping{}
			authenticationPolicyContractAssertionMappingsAttrs := authenticationPolicyContractAssertionMappingsElement.(types.Object).Attributes()
			authenticationPolicyContractAssertionMappingsValue.AbortSsoTransactionAsFailSafe = authenticationPolicyContractAssertionMappingsAttrs["abort_sso_transaction_as_fail_safe"].(types.Bool).ValueBoolPointer()
			authenticationPolicyContractAssertionMappingsValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(authenticationPolicyContractAssertionMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			authenticationPolicyContractAssertionMappingsValue.AttributeSources = attributesources.ClientStruct(authenticationPolicyContractAssertionMappingsAttrs["attribute_sources"].(types.Set))
			authenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue := client.ResourceLink{}
			authenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrs := authenticationPolicyContractAssertionMappingsAttrs["authentication_policy_contract_ref"].(types.Object).Attributes()
			authenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue.Id = authenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrs["id"].(types.String).ValueString()
			authenticationPolicyContractAssertionMappingsValue.AuthenticationPolicyContractRef = authenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue
			authenticationPolicyContractAssertionMappingsValue.IssuanceCriteria = issuancecriteria.ClientStruct(authenticationPolicyContractAssertionMappingsAttrs["issuance_criteria"].(types.Object))
			authenticationPolicyContractAssertionMappingsValue.RestrictVirtualEntityIds = authenticationPolicyContractAssertionMappingsAttrs["restrict_virtual_entity_ids"].(types.Bool).ValueBoolPointer()
			if !authenticationPolicyContractAssertionMappingsAttrs["restricted_virtual_entity_ids"].IsNull() {
				authenticationPolicyContractAssertionMappingsValue.RestrictedVirtualEntityIds = []string{}
				for _, restrictedVirtualEntityIdsElement := range authenticationPolicyContractAssertionMappingsAttrs["restricted_virtual_entity_ids"].(types.Set).Elements() {
					authenticationPolicyContractAssertionMappingsValue.RestrictedVirtualEntityIds = append(authenticationPolicyContractAssertionMappingsValue.RestrictedVirtualEntityIds, restrictedVirtualEntityIdsElement.(types.String).ValueString())
				}
			}
			spBrowserSsoValue.AuthenticationPolicyContractAssertionMappings = append(spBrowserSsoValue.AuthenticationPolicyContractAssertionMappings, authenticationPolicyContractAssertionMappingsValue)
		}
		spBrowserSsoValue.DefaultTargetUrl = spBrowserSsoAttrs["default_target_url"].(types.String).ValueStringPointer()
		if !spBrowserSsoAttrs["enabled_profiles"].IsNull() {
			spBrowserSsoValue.EnabledProfiles = []string{}
			for _, enabledProfilesElement := range spBrowserSsoAttrs["enabled_profiles"].(types.Set).Elements() {
				spBrowserSsoValue.EnabledProfiles = append(spBrowserSsoValue.EnabledProfiles, enabledProfilesElement.(types.String).ValueString())
			}
		}
		if !spBrowserSsoAttrs["encryption_policy"].IsNull() {
			spBrowserSsoEncryptionPolicyValue := &client.EncryptionPolicy{}
			spBrowserSsoEncryptionPolicyAttrs := spBrowserSsoAttrs["encryption_policy"].(types.Object).Attributes()
			spBrowserSsoEncryptionPolicyValue.EncryptAssertion = spBrowserSsoEncryptionPolicyAttrs["encrypt_assertion"].(types.Bool).ValueBoolPointer()
			spBrowserSsoEncryptionPolicyValue.EncryptSloSubjectNameId = spBrowserSsoEncryptionPolicyAttrs["encrypt_slo_subject_name_id"].(types.Bool).ValueBoolPointer()
			if !spBrowserSsoEncryptionPolicyAttrs["encrypted_attributes"].IsNull() {
				spBrowserSsoEncryptionPolicyValue.EncryptedAttributes = []string{}
				for _, encryptedAttributesElement := range spBrowserSsoEncryptionPolicyAttrs["encrypted_attributes"].(types.Set).Elements() {
					spBrowserSsoEncryptionPolicyValue.EncryptedAttributes = append(spBrowserSsoEncryptionPolicyValue.EncryptedAttributes, encryptedAttributesElement.(types.String).ValueString())
				}
			}
			spBrowserSsoEncryptionPolicyValue.SloSubjectNameIDEncrypted = spBrowserSsoEncryptionPolicyAttrs["slo_subject_name_id_encrypted"].(types.Bool).ValueBoolPointer()
			spBrowserSsoValue.EncryptionPolicy = spBrowserSsoEncryptionPolicyValue
		}
		if !spBrowserSsoAttrs["incoming_bindings"].IsNull() {
			spBrowserSsoValue.IncomingBindings = []string{}
			for _, incomingBindingsElement := range spBrowserSsoAttrs["incoming_bindings"].(types.Set).Elements() {
				spBrowserSsoValue.IncomingBindings = append(spBrowserSsoValue.IncomingBindings, incomingBindingsElement.(types.String).ValueString())
			}
		}
		spBrowserSsoValue.MessageCustomizations = []client.ProtocolMessageCustomization{}
		for _, messageCustomizationsElement := range spBrowserSsoAttrs["message_customizations"].(types.Set).Elements() {
			messageCustomizationsValue := client.ProtocolMessageCustomization{}
			messageCustomizationsAttrs := messageCustomizationsElement.(types.Object).Attributes()
			messageCustomizationsValue.ContextName = messageCustomizationsAttrs["context_name"].(types.String).ValueStringPointer()
			messageCustomizationsValue.MessageExpression = messageCustomizationsAttrs["message_expression"].(types.String).ValueStringPointer()
			spBrowserSsoValue.MessageCustomizations = append(spBrowserSsoValue.MessageCustomizations, messageCustomizationsValue)
		}
		spBrowserSsoValue.Protocol = spBrowserSsoAttrs["protocol"].(types.String).ValueString()
		spBrowserSsoValue.RequireSignedAuthnRequests = spBrowserSsoAttrs["require_signed_authn_requests"].(types.Bool).ValueBoolPointer()
		spBrowserSsoValue.SignAssertions = spBrowserSsoAttrs["sign_assertions"].(types.Bool).ValueBoolPointer()
		spBrowserSsoValue.SignResponseAsRequired = spBrowserSsoAttrs["sign_response_as_required"].(types.Bool).ValueBoolPointer()
		spBrowserSsoValue.SloServiceEndpoints = []client.SloServiceEndpoint{}
		for _, sloServiceEndpointsElement := range spBrowserSsoAttrs["slo_service_endpoints"].(types.Set).Elements() {
			sloServiceEndpointsValue := client.SloServiceEndpoint{}
			sloServiceEndpointsAttrs := sloServiceEndpointsElement.(types.Object).Attributes()
			sloServiceEndpointsValue.Binding = sloServiceEndpointsAttrs["binding"].(types.String).ValueStringPointer()
			sloServiceEndpointsValue.ResponseUrl = sloServiceEndpointsAttrs["response_url"].(types.String).ValueStringPointer()
			sloServiceEndpointsValue.Url = sloServiceEndpointsAttrs["url"].(types.String).ValueString()
			spBrowserSsoValue.SloServiceEndpoints = append(spBrowserSsoValue.SloServiceEndpoints, sloServiceEndpointsValue)
		}
		spBrowserSsoValue.SpSamlIdentityMapping = spBrowserSsoAttrs["sp_saml_identity_mapping"].(types.String).ValueStringPointer()
		spBrowserSsoValue.SpWsFedIdentityMapping = spBrowserSsoAttrs["sp_ws_fed_identity_mapping"].(types.String).ValueStringPointer()
		spBrowserSsoValue.SsoApplicationEndpoint = spBrowserSsoAttrs["sso_application_endpoint"].(types.String).ValueStringPointer()
		spBrowserSsoValue.SsoServiceEndpoints = []client.SpSsoServiceEndpoint{}
		for _, ssoServiceEndpointsElement := range spBrowserSsoAttrs["sso_service_endpoints"].(types.Set).Elements() {
			ssoServiceEndpointsValue := client.SpSsoServiceEndpoint{}
			ssoServiceEndpointsAttrs := ssoServiceEndpointsElement.(types.Object).Attributes()
			ssoServiceEndpointsValue.Binding = ssoServiceEndpointsAttrs["binding"].(types.String).ValueStringPointer()
			ssoServiceEndpointsValue.Index = ssoServiceEndpointsAttrs["index"].(types.Int64).ValueInt64Pointer()
			ssoServiceEndpointsValue.IsDefault = ssoServiceEndpointsAttrs["is_default"].(types.Bool).ValueBoolPointer()
			ssoServiceEndpointsValue.Url = ssoServiceEndpointsAttrs["url"].(types.String).ValueString()
			spBrowserSsoValue.SsoServiceEndpoints = append(spBrowserSsoValue.SsoServiceEndpoints, ssoServiceEndpointsValue)
		}
		for _, urlWhitelistEntriesElement := range spBrowserSsoAttrs["url_whitelist_entries"].(types.Set).Elements() {
			urlWhitelistEntriesValue := client.UrlWhitelistEntry{}
			urlWhitelistEntriesAttrs := urlWhitelistEntriesElement.(types.Object).Attributes()
			urlWhitelistEntriesValue.AllowQueryAndFragment = urlWhitelistEntriesAttrs["allow_query_and_fragment"].(types.Bool).ValueBoolPointer()
			urlWhitelistEntriesValue.RequireHttps = urlWhitelistEntriesAttrs["require_https"].(types.Bool).ValueBoolPointer()
			urlWhitelistEntriesValue.ValidDomain = urlWhitelistEntriesAttrs["valid_domain"].(types.String).ValueStringPointer()
			urlWhitelistEntriesValue.ValidPath = urlWhitelistEntriesAttrs["valid_path"].(types.String).ValueStringPointer()
			spBrowserSsoValue.UrlWhitelistEntries = append(spBrowserSsoValue.UrlWhitelistEntries, urlWhitelistEntriesValue)
		}
		spBrowserSsoValue.WsFedTokenType = spBrowserSsoAttrs["ws_fed_token_type"].(types.String).ValueStringPointer()
		spBrowserSsoValue.WsTrustVersion = spBrowserSsoAttrs["ws_trust_version"].(types.String).ValueStringPointer()
		addRequest.SpBrowserSso = spBrowserSsoValue
	}

	if internaltypes.IsDefined(plan.AttributeQuery) {
		addRequest.AttributeQuery = &client.SpAttributeQuery{}
		attributeQueryAttrs := plan.AttributeQuery.Attributes()

		addRequest.AttributeQuery.Attributes = []string{}
		for _, attribute := range attributeQueryAttrs["attributes"].(types.Set).Elements() {
			addRequest.AttributeQuery.Attributes = append(addRequest.AttributeQuery.Attributes, attribute.(types.String).ValueString())
		}

		addRequest.AttributeQuery.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(attributeQueryAttrs["attribute_contract_fulfillment"].(types.Map))

		addRequest.AttributeQuery.IssuanceCriteria = issuancecriteria.ClientStruct(attributeQueryAttrs["issuance_criteria"].(types.Object))

		if !attributeQueryAttrs["policy"].IsNull() {
			attributeQueryPolicyValue := &client.SpAttributeQueryPolicy{}
			attributeQueryPolicyAttrs := attributeQueryAttrs["policy"].(types.Object).Attributes()
			attributeQueryPolicyValue.EncryptAssertion = attributeQueryPolicyAttrs["encrypt_assertion"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.RequireEncryptedNameId = attributeQueryPolicyAttrs["require_encrypted_name_id"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.RequireSignedAttributeQuery = attributeQueryPolicyAttrs["require_signed_attribute_query"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.SignAssertion = attributeQueryPolicyAttrs["sign_assertion"].(types.Bool).ValueBoolPointer()
			attributeQueryPolicyValue.SignResponse = attributeQueryPolicyAttrs["sign_response"].(types.Bool).ValueBoolPointer()
			addRequest.AttributeQuery.Policy = attributeQueryPolicyValue
		}

		addRequest.AttributeQuery.AttributeSources = attributesources.ClientStruct(attributeQueryAttrs["attribute_sources"].(types.Set))
	}

	if !plan.WsTrust.IsNull() {
		wsTrustValue := &client.SpWsTrust{}
		wsTrustAttrs := plan.WsTrust.Attributes()
		wsTrustValue.AbortIfNotFulfilledFromRequest = wsTrustAttrs["abort_if_not_fulfilled_from_request"].(types.Bool).ValueBoolPointer()
		wsTrustAttributeContractValue := client.SpWsTrustAttributeContract{}
		wsTrustAttributeContractAttrs := wsTrustAttrs["attribute_contract"].(types.Object).Attributes()
		wsTrustAttributeContractValue.CoreAttributes = []client.SpWsTrustAttribute{}
		for _, coreAttributesElement := range wsTrustAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.SpWsTrustAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			coreAttributesValue.Namespace = coreAttributesAttrs["namespace"].(types.String).ValueString()
			wsTrustAttributeContractValue.CoreAttributes = append(wsTrustAttributeContractValue.CoreAttributes, coreAttributesValue)
		}
		wsTrustAttributeContractValue.ExtendedAttributes = []client.SpWsTrustAttribute{}
		for _, extendedAttributesElement := range wsTrustAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.SpWsTrustAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			extendedAttributesValue.Namespace = extendedAttributesAttrs["namespace"].(types.String).ValueString()
			wsTrustAttributeContractValue.ExtendedAttributes = append(wsTrustAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		wsTrustValue.AttributeContract = wsTrustAttributeContractValue
		wsTrustValue.DefaultTokenType = wsTrustAttrs["default_token_type"].(types.String).ValueStringPointer()
		wsTrustValue.EncryptSaml2Assertion = wsTrustAttrs["encrypt_saml2_assertion"].(types.Bool).ValueBoolPointer()
		wsTrustValue.GenerateKey = wsTrustAttrs["generate_key"].(types.Bool).ValueBoolPointer()
		wsTrustValue.MessageCustomizations = []client.ProtocolMessageCustomization{}
		for _, messageCustomizationsElement := range wsTrustAttrs["message_customizations"].(types.Set).Elements() {
			messageCustomizationsValue := client.ProtocolMessageCustomization{}
			messageCustomizationsAttrs := messageCustomizationsElement.(types.Object).Attributes()
			messageCustomizationsValue.ContextName = messageCustomizationsAttrs["context_name"].(types.String).ValueStringPointer()
			messageCustomizationsValue.MessageExpression = messageCustomizationsAttrs["message_expression"].(types.String).ValueStringPointer()
			wsTrustValue.MessageCustomizations = append(wsTrustValue.MessageCustomizations, messageCustomizationsValue)
		}
		wsTrustValue.MinutesAfter = wsTrustAttrs["minutes_after"].(types.Int64).ValueInt64Pointer()
		wsTrustValue.MinutesBefore = wsTrustAttrs["minutes_before"].(types.Int64).ValueInt64Pointer()
		wsTrustValue.OAuthAssertionProfiles = wsTrustAttrs["oauth_assertion_profiles"].(types.Bool).ValueBoolPointer()
		wsTrustValue.PartnerServiceIds = []string{}
		for _, partnerServiceIdsElement := range wsTrustAttrs["partner_service_ids"].(types.Set).Elements() {
			wsTrustValue.PartnerServiceIds = append(wsTrustValue.PartnerServiceIds, partnerServiceIdsElement.(types.String).ValueString())
		}
		if !wsTrustAttrs["request_contract_ref"].IsNull() {
			wsTrustRequestContractRefValue := &client.ResourceLink{}
			wsTrustRequestContractRefAttrs := wsTrustAttrs["request_contract_ref"].(types.Object).Attributes()
			wsTrustRequestContractRefValue.Id = wsTrustRequestContractRefAttrs["id"].(types.String).ValueString()
			wsTrustValue.RequestContractRef = wsTrustRequestContractRefValue
		}
		wsTrustValue.TokenProcessorMappings = []client.IdpTokenProcessorMapping{}
		for _, tokenProcessorMappingsElement := range wsTrustAttrs["token_processor_mappings"].(types.Set).Elements() {
			tokenProcessorMappingsValue := client.IdpTokenProcessorMapping{}
			tokenProcessorMappingsAttrs := tokenProcessorMappingsElement.(types.Object).Attributes()
			tokenProcessorMappingsValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(tokenProcessorMappingsAttrs["attribute_contract_fulfillment"].(types.Map))
			tokenProcessorMappingsValue.AttributeSources = attributesources.ClientStruct(tokenProcessorMappingsAttrs["attribute_sources"].(types.Set))
			tokenProcessorMappingsIdpTokenProcessorRefValue := client.ResourceLink{}
			tokenProcessorMappingsIdpTokenProcessorRefAttrs := tokenProcessorMappingsAttrs["idp_token_processor_ref"].(types.Object).Attributes()
			tokenProcessorMappingsIdpTokenProcessorRefValue.Id = tokenProcessorMappingsIdpTokenProcessorRefAttrs["id"].(types.String).ValueString()
			tokenProcessorMappingsValue.IdpTokenProcessorRef = tokenProcessorMappingsIdpTokenProcessorRefValue
			tokenProcessorMappingsValue.IssuanceCriteria = issuancecriteria.ClientStruct(tokenProcessorMappingsAttrs["issuance_criteria"].(types.Object))
			if !tokenProcessorMappingsAttrs["restricted_virtual_entity_ids"].IsNull() {
				tokenProcessorMappingsValue.RestrictedVirtualEntityIds = []string{}
				for _, restrictedVirtualEntityIdsElement := range tokenProcessorMappingsAttrs["restricted_virtual_entity_ids"].(types.Set).Elements() {
					tokenProcessorMappingsValue.RestrictedVirtualEntityIds = append(tokenProcessorMappingsValue.RestrictedVirtualEntityIds, restrictedVirtualEntityIdsElement.(types.String).ValueString())
				}
			}
			wsTrustValue.TokenProcessorMappings = append(wsTrustValue.TokenProcessorMappings, tokenProcessorMappingsValue)
		}
		addRequest.WsTrust = wsTrustValue
	}

	// outbound_provision
	if !plan.OutboundProvision.IsNull() {
		outboundProvisionValue := &client.OutboundProvision{}
		outboundProvisionAttrs := plan.OutboundProvision.Attributes()
		outboundProvisionValue.Channels = []client.Channel{}
		for _, channelsElement := range outboundProvisionAttrs["channels"].(types.List).Elements() {
			channelsValue := client.Channel{}
			channelsAttrs := channelsElement.(types.Object).Attributes()
			channelsValue.Active = channelsAttrs["active"].(types.Bool).ValueBool()
			channelsValue.AttributeMapping = []client.SaasAttributeMapping{}
			for _, attributeMappingElement := range channelsAttrs["attribute_mapping"].(types.Set).Elements() {
				attributeMappingValue := client.SaasAttributeMapping{}
				attributeMappingAttrs := attributeMappingElement.(types.Object).Attributes()
				attributeMappingValue.FieldName = attributeMappingAttrs["field_name"].(types.String).ValueString()
				attributeMappingSaasFieldInfoValue := client.SaasFieldConfiguration{}
				attributeMappingSaasFieldInfoAttrs := attributeMappingAttrs["saas_field_info"].(types.Object).Attributes()
				if !attributeMappingSaasFieldInfoAttrs["attribute_names"].IsNull() {
					attributeMappingSaasFieldInfoValue.AttributeNames = []string{}
					for _, attributeNamesElement := range attributeMappingSaasFieldInfoAttrs["attribute_names"].(types.List).Elements() {
						attributeMappingSaasFieldInfoValue.AttributeNames = append(attributeMappingSaasFieldInfoValue.AttributeNames, attributeNamesElement.(types.String).ValueString())
					}
				}
				if !attributeMappingSaasFieldInfoAttrs["character_case"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.CharacterCase = attributeMappingSaasFieldInfoAttrs["character_case"].(types.String).ValueStringPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["create_only"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.CreateOnly = attributeMappingSaasFieldInfoAttrs["create_only"].(types.Bool).ValueBoolPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["default_value"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.DefaultValue = attributeMappingSaasFieldInfoAttrs["default_value"].(types.String).ValueStringPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["expression"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.Expression = attributeMappingSaasFieldInfoAttrs["expression"].(types.String).ValueStringPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["masked"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.Masked = attributeMappingSaasFieldInfoAttrs["masked"].(types.Bool).ValueBoolPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["parser"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.Parser = attributeMappingSaasFieldInfoAttrs["parser"].(types.String).ValueStringPointer()
				}
				if !attributeMappingSaasFieldInfoAttrs["trim"].IsUnknown() {
					attributeMappingSaasFieldInfoValue.Trim = attributeMappingSaasFieldInfoAttrs["trim"].(types.Bool).ValueBoolPointer()
				}
				attributeMappingValue.SaasFieldInfo = attributeMappingSaasFieldInfoValue
				channelsValue.AttributeMapping = append(channelsValue.AttributeMapping, attributeMappingValue)
			}
			channelsChannelSourceValue := client.ChannelSource{}
			channelsChannelSourceAttrs := channelsAttrs["channel_source"].(types.Object).Attributes()
			channelsChannelSourceAccountManagementSettingsValue := client.AccountManagementSettings{}
			channelsChannelSourceAccountManagementSettingsAttrs := channelsChannelSourceAttrs["account_management_settings"].(types.Object).Attributes()
			channelsChannelSourceAccountManagementSettingsValue.AccountStatusAlgorithm = channelsChannelSourceAccountManagementSettingsAttrs["account_status_algorithm"].(types.String).ValueString()
			channelsChannelSourceAccountManagementSettingsValue.AccountStatusAttributeName = channelsChannelSourceAccountManagementSettingsAttrs["account_status_attribute_name"].(types.String).ValueString()
			channelsChannelSourceAccountManagementSettingsValue.DefaultStatus = channelsChannelSourceAccountManagementSettingsAttrs["default_status"].(types.Bool).ValueBoolPointer()
			channelsChannelSourceAccountManagementSettingsValue.FlagComparisonStatus = channelsChannelSourceAccountManagementSettingsAttrs["flag_comparison_status"].(types.Bool).ValueBoolPointer()
			channelsChannelSourceAccountManagementSettingsValue.FlagComparisonValue = channelsChannelSourceAccountManagementSettingsAttrs["flag_comparison_value"].(types.String).ValueStringPointer()
			channelsChannelSourceValue.AccountManagementSettings = channelsChannelSourceAccountManagementSettingsValue
			channelsChannelSourceValue.BaseDn = channelsChannelSourceAttrs["base_dn"].(types.String).ValueString()
			channelsChannelSourceChangeDetectionSettingsValue := client.ChangeDetectionSettings{}
			channelsChannelSourceChangeDetectionSettingsAttrs := channelsChannelSourceAttrs["change_detection_settings"].(types.Object).Attributes()
			channelsChannelSourceChangeDetectionSettingsValue.ChangedUsersAlgorithm = channelsChannelSourceChangeDetectionSettingsAttrs["changed_users_algorithm"].(types.String).ValueString()
			channelsChannelSourceChangeDetectionSettingsValue.GroupObjectClass = channelsChannelSourceChangeDetectionSettingsAttrs["group_object_class"].(types.String).ValueString()
			channelsChannelSourceChangeDetectionSettingsValue.TimeStampAttributeName = channelsChannelSourceChangeDetectionSettingsAttrs["time_stamp_attribute_name"].(types.String).ValueString()
			channelsChannelSourceChangeDetectionSettingsValue.UserObjectClass = channelsChannelSourceChangeDetectionSettingsAttrs["user_object_class"].(types.String).ValueString()
			channelsChannelSourceChangeDetectionSettingsValue.UsnAttributeName = channelsChannelSourceChangeDetectionSettingsAttrs["usn_attribute_name"].(types.String).ValueStringPointer()
			channelsChannelSourceValue.ChangeDetectionSettings = channelsChannelSourceChangeDetectionSettingsValue
			channelsChannelSourceDataSourceValue := client.ResourceLink{}
			channelsChannelSourceDataSourceAttrs := channelsChannelSourceAttrs["data_source"].(types.Object).Attributes()
			channelsChannelSourceDataSourceValue.Id = channelsChannelSourceDataSourceAttrs["id"].(types.String).ValueString()
			channelsChannelSourceValue.DataSource = channelsChannelSourceDataSourceValue
			channelsChannelSourceGroupMembershipDetectionValue := client.GroupMembershipDetection{}
			channelsChannelSourceGroupMembershipDetectionAttrs := channelsChannelSourceAttrs["group_membership_detection"].(types.Object).Attributes()
			channelsChannelSourceGroupMembershipDetectionValue.GroupMemberAttributeName = channelsChannelSourceGroupMembershipDetectionAttrs["group_member_attribute_name"].(types.String).ValueStringPointer()
			channelsChannelSourceGroupMembershipDetectionValue.MemberOfGroupAttributeName = channelsChannelSourceGroupMembershipDetectionAttrs["member_of_group_attribute_name"].(types.String).ValueStringPointer()
			channelsChannelSourceValue.GroupMembershipDetection = channelsChannelSourceGroupMembershipDetectionValue
			if !channelsChannelSourceAttrs["group_source_location"].IsNull() {
				channelsChannelSourceGroupSourceLocationValue := &client.ChannelSourceLocation{}
				channelsChannelSourceGroupSourceLocationAttrs := channelsChannelSourceAttrs["group_source_location"].(types.Object).Attributes()
				channelsChannelSourceGroupSourceLocationValue.Filter = channelsChannelSourceGroupSourceLocationAttrs["filter"].(types.String).ValueStringPointer()
				channelsChannelSourceGroupSourceLocationValue.GroupDN = channelsChannelSourceGroupSourceLocationAttrs["group_dn"].(types.String).ValueStringPointer()
				channelsChannelSourceGroupSourceLocationValue.NestedSearch = channelsChannelSourceGroupSourceLocationAttrs["nested_search"].(types.Bool).ValueBoolPointer()
				channelsChannelSourceValue.GroupSourceLocation = channelsChannelSourceGroupSourceLocationValue
			}
			channelsChannelSourceValue.GuidAttributeName = channelsChannelSourceAttrs["guid_attribute_name"].(types.String).ValueString()
			channelsChannelSourceValue.GuidBinary = channelsChannelSourceAttrs["guid_binary"].(types.Bool).ValueBool()
			channelsChannelSourceUserSourceLocationValue := client.ChannelSourceLocation{}
			channelsChannelSourceUserSourceLocationAttrs := channelsChannelSourceAttrs["user_source_location"].(types.Object).Attributes()
			channelsChannelSourceUserSourceLocationValue.Filter = channelsChannelSourceUserSourceLocationAttrs["filter"].(types.String).ValueStringPointer()
			channelsChannelSourceUserSourceLocationValue.GroupDN = channelsChannelSourceUserSourceLocationAttrs["group_dn"].(types.String).ValueStringPointer()
			channelsChannelSourceUserSourceLocationValue.NestedSearch = channelsChannelSourceUserSourceLocationAttrs["nested_search"].(types.Bool).ValueBoolPointer()
			channelsChannelSourceValue.UserSourceLocation = channelsChannelSourceUserSourceLocationValue
			channelsValue.ChannelSource = channelsChannelSourceValue
			channelsValue.MaxThreads = channelsAttrs["max_threads"].(types.Int64).ValueInt64()
			channelsValue.Name = channelsAttrs["name"].(types.String).ValueString()
			channelsValue.Timeout = channelsAttrs["timeout"].(types.Int64).ValueInt64()
			outboundProvisionValue.Channels = append(outboundProvisionValue.Channels, channelsValue)
		}
		if !outboundProvisionAttrs["custom_schema"].IsNull() {
			outboundProvisionCustomSchemaValue := &client.Schema{}
			outboundProvisionCustomSchemaAttrs := outboundProvisionAttrs["custom_schema"].(types.Object).Attributes()
			outboundProvisionCustomSchemaValue.Attributes = []client.SchemaAttribute{}
			for _, attributesElement := range outboundProvisionCustomSchemaAttrs["attributes"].(types.Set).Elements() {
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
				outboundProvisionCustomSchemaValue.Attributes = append(outboundProvisionCustomSchemaValue.Attributes, attributesValue)
			}
			outboundProvisionCustomSchemaValue.Namespace = outboundProvisionCustomSchemaAttrs["namespace"].(types.String).ValueStringPointer()
			outboundProvisionValue.CustomSchema = outboundProvisionCustomSchemaValue
		}
		outboundProvisionValue.TargetSettings = []client.ConfigField{}
		for _, targetSettingsElement := range outboundProvisionAttrs["target_settings"].(types.Set).Elements() {
			targetSettingsValue := client.ConfigField{}
			targetSettingsAttrs := targetSettingsElement.(types.Object).Attributes()
			targetSettingsValue.Name = targetSettingsAttrs["name"].(types.String).ValueString()
			targetSettingsValue.Value = targetSettingsAttrs["value"].(types.String).ValueStringPointer()
			outboundProvisionValue.TargetSettings = append(outboundProvisionValue.TargetSettings, targetSettingsValue)
		}
		outboundProvisionValue.Type = outboundProvisionAttrs["type"].(types.String).ValueString()
		addRequest.OutboundProvision = outboundProvisionValue
	}

	return respDiags
}

func (state *idpSpConnectionModel) getSpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfiguration(adapterMappingIndex int) types.Object {
	if !internaltypes.IsDefined(state.SpBrowserSso) {
		return types.ObjectNull(pluginconfiguration.AttrTypes())
	}
	adapterMappings := state.SpBrowserSso.Attributes()["adapter_mappings"].(types.Set).Elements()
	if adapterMappingIndex >= len(adapterMappings) {
		return types.ObjectNull(pluginconfiguration.AttrTypes())
	}

	adapterMapping := adapterMappings[adapterMappingIndex].(types.Object)
	adapterOverrideSettings := adapterMapping.Attributes()["adapter_override_settings"].(types.Object)
	if !internaltypes.IsDefined(adapterOverrideSettings) {
		return types.ObjectNull(pluginconfiguration.AttrTypes())
	}

	return adapterOverrideSettings.Attributes()["configuration"].(types.Object)
}

// Returns the attribute_mapping and attribute_mapping_all attributes, where the _all attribute
// includes any values added in the response by PingFed
func (state *idpSpConnectionModel) buildAttributeMappingAttrs(channelName string, responseAttributeMappings []client.SaasAttributeMapping, isImportRead bool) (types.Set, types.Set, diag.Diagnostics) {
	// Get a list of field_name values that were expected based on the state
	var expectedFieldNames []string
	if internaltypes.IsDefined(state.OutboundProvision) && internaltypes.IsDefined(state.OutboundProvision.Attributes()["channels"]) {
		channelsElements := state.OutboundProvision.Attributes()["channels"].(types.List).Elements()
		for _, channel := range channelsElements {
			stateChannelName := channel.(types.Object).Attributes()["name"].(types.String).ValueString()
			if stateChannelName == channelName {
				// Get the attribute_mapping values specified in the state
				for _, attrMapping := range channel.(types.Object).Attributes()["attribute_mapping"].(types.Set).Elements() {
					expectedFieldNames = append(expectedFieldNames, attrMapping.(types.Object).Attributes()["field_name"].(types.String).ValueString())
				}
				break
			}
		}
	}

	var respDiags diag.Diagnostics
	var allAttributeMappings, expectedAttributeMappings []attr.Value
	for _, outboundProvisionChannelsAttributeMappingResponseValue := range responseAttributeMappings {
		outboundProvisionChannelsAttributeMappingSaasFieldInfoAttributeNamesValue, diags := types.ListValueFrom(context.Background(), types.StringType, outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.AttributeNames)
		respDiags.Append(diags...)
		outboundProvisionChannelsAttributeMappingSaasFieldInfoValue, diags := types.ObjectValue(saasFieldInfoAttrTypes, map[string]attr.Value{
			"attribute_names": outboundProvisionChannelsAttributeMappingSaasFieldInfoAttributeNamesValue,
			"character_case":  types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.CharacterCase),
			"create_only":     types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.CreateOnly),
			"default_value":   types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.DefaultValue),
			"expression":      types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Expression),
			"masked":          types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Masked),
			"parser":          types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Parser),
			"trim":            types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Trim),
		})
		respDiags.Append(diags...)
		outboundProvisionChannelsAttributeMappingValue, diags := types.ObjectValue(attributeMappingElemAttrTypes.AttrTypes, map[string]attr.Value{
			"field_name":      types.StringValue(outboundProvisionChannelsAttributeMappingResponseValue.FieldName),
			"saas_field_info": outboundProvisionChannelsAttributeMappingSaasFieldInfoValue,
		})
		respDiags.Append(diags...)
		allAttributeMappings = append(allAttributeMappings, outboundProvisionChannelsAttributeMappingValue)
		if slices.Contains(expectedFieldNames, outboundProvisionChannelsAttributeMappingResponseValue.FieldName) {
			expectedAttributeMappings = append(expectedAttributeMappings, outboundProvisionChannelsAttributeMappingValue)
		}
	}
	attributeMappingAll, diags := types.SetValue(attributeMappingElemAttrTypes, allAttributeMappings)
	respDiags.Append(diags...)
	var attributeMapping types.Set
	// On import, just read the attribute mappings into attribute_mapping directly
	if isImportRead {
		attributeMapping, diags = types.SetValue(attributeMappingElemAttrTypes, allAttributeMappings)
		respDiags.Append(diags...)
	} else {
		attributeMapping, diags = types.SetValue(attributeMappingElemAttrTypes, expectedAttributeMappings)
		respDiags.Append(diags...)
	}
	return attributeMapping, attributeMappingAll, respDiags
}

// Returns the target_settings and target_settings_all attributes, where the _all attribute
// includes any values added in the response by PingFed
func (state *idpSpConnectionModel) buildTargetSettingsAttrs(responseTargetSettings []client.ConfigField, isImportRead bool) (types.Set, types.Set, diag.Diagnostics) {
	// Get a list of target_setting names that were expected based on the state
	expectedTargetSettingsValues := map[string]string{}
	expectedTargetSettingsEncryptedValues := map[string]types.String{}
	if internaltypes.IsDefined(state.OutboundProvision) {
		for _, targetSetting := range state.OutboundProvision.Attributes()["target_settings"].(types.Set).Elements() {
			targetSettingsAttrs := targetSetting.(types.Object).Attributes()
			expectedTargetSettingsValues[targetSettingsAttrs["name"].(types.String).ValueString()] = targetSettingsAttrs["value"].(types.String).ValueString()
			expectedTargetSettingsEncryptedValues[targetSettingsAttrs["name"].(types.String).ValueString()] = targetSettingsAttrs["encrypted_value"].(types.String)
		}
	}

	var respDiags diag.Diagnostics
	var allTargetSettings, expectedTargetSettings []attr.Value
	for _, outboundProvisionTargetSettingsResponseValue := range responseTargetSettings {
		expectedValue, settingInPlan := expectedTargetSettingsValues[outboundProvisionTargetSettingsResponseValue.Name]
		responseValue := types.StringPointerValue(outboundProvisionTargetSettingsResponseValue.Value)
		if settingInPlan && outboundProvisionTargetSettingsResponseValue.Value == nil {
			responseValue = types.StringValue(expectedValue)
		}
		encryptedValue := types.StringPointerValue(outboundProvisionTargetSettingsResponseValue.EncryptedValue)
		if !expectedTargetSettingsEncryptedValues[outboundProvisionTargetSettingsResponseValue.Name].IsUnknown() {
			encryptedValue = types.StringValue(expectedTargetSettingsEncryptedValues[outboundProvisionTargetSettingsResponseValue.Name].ValueString())
		}
		outboundProvisionTargetSettingsValue, diags := types.ObjectValue(targetSettingsElemAttrType.AttrTypes, map[string]attr.Value{
			"name":            types.StringValue(outboundProvisionTargetSettingsResponseValue.Name),
			"value":           responseValue,
			"encrypted_value": encryptedValue,
		})
		respDiags.Append(diags...)
		allTargetSettings = append(allTargetSettings, outboundProvisionTargetSettingsValue)
		if settingInPlan {
			expectedTargetSettings = append(expectedTargetSettings, outboundProvisionTargetSettingsValue)
		}
	}
	targetSettingsAll, diags := types.SetValue(targetSettingsElemAttrType, allTargetSettings)
	respDiags.Append(diags...)
	var targetSettings types.Set
	if isImportRead {
		// On import, just read the target settings into target_settings directly
		targetSettings, diags = types.SetValue(targetSettingsElemAttrType, allTargetSettings)
		respDiags.Append(diags...)
	} else {
		targetSettings, diags = types.SetValue(targetSettingsElemAttrType, expectedTargetSettings)
		respDiags.Append(diags...)
	}
	return targetSettings, targetSettingsAll, respDiags
}

func (state *idpSpConnectionModel) readClientResponse(response *client.SpConnection, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// active
	state.Active = types.BoolPointerValue(response.Active)
	// additional_allowed_entities_configuration
	additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes := map[string]attr.Type{
		"entity_description": types.StringType,
		"entity_id":          types.StringType,
	}
	additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType := types.ObjectType{AttrTypes: additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes}
	additionalAllowedEntitiesConfigurationAttrTypes := map[string]attr.Type{
		"additional_allowed_entities": types.SetType{ElemType: additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType},
		"allow_additional_entities":   types.BoolType,
		"allow_all_entities":          types.BoolType,
	}
	var additionalAllowedEntitiesConfigurationValue types.Object
	if response.AdditionalAllowedEntitiesConfiguration == nil {
		additionalAllowedEntitiesConfigurationValue = types.ObjectNull(additionalAllowedEntitiesConfigurationAttrTypes)
	} else {
		var additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues []attr.Value
		for _, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue := range response.AdditionalAllowedEntitiesConfiguration.AdditionalAllowedEntities {
			additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue, diags := types.ObjectValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes, map[string]attr.Value{
				"entity_description": types.StringPointerValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue.EntityDescription),
				"entity_id":          types.StringPointerValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue.EntityId),
			})
			respDiags.Append(diags...)
			additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues = append(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue)
		}
		additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue, diags := types.SetValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues)
		respDiags.Append(diags...)
		additionalAllowedEntitiesConfigurationValue, diags = types.ObjectValue(additionalAllowedEntitiesConfigurationAttrTypes, map[string]attr.Value{
			"additional_allowed_entities": additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue,
			"allow_additional_entities":   types.BoolPointerValue(response.AdditionalAllowedEntitiesConfiguration.AllowAdditionalEntities),
			"allow_all_entities":          types.BoolPointerValue(response.AdditionalAllowedEntitiesConfiguration.AllowAllEntities),
		})
		respDiags.Append(diags...)
	}

	state.AdditionalAllowedEntitiesConfiguration = additionalAllowedEntitiesConfigurationValue
	// application_icon_url
	state.ApplicationIconUrl = types.StringPointerValue(response.ApplicationIconUrl)
	// application_name
	state.ApplicationName = types.StringPointerValue(response.ApplicationName)
	// attribute_query
	attributeQueryAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	attributeQueryAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: attributeQueryAttributeContractFulfillmentAttrTypes}
	attributeQueryAttributeSourcesAttrTypes := attributesources.AttrTypes()
	attributeQueryAttributeSourcesElementType := types.ObjectType{AttrTypes: attributeQueryAttributeSourcesAttrTypes}
	attributeQueryIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	attributeQueryPolicyAttrTypes := map[string]attr.Type{
		"encrypt_assertion":              types.BoolType,
		"require_encrypted_name_id":      types.BoolType,
		"require_signed_attribute_query": types.BoolType,
		"sign_assertion":                 types.BoolType,
		"sign_response":                  types.BoolType,
	}
	attributeQueryAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: attributeQueryAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: attributeQueryAttributeSourcesElementType},
		"attributes":                     types.SetType{ElemType: types.StringType},
		"issuance_criteria":              types.ObjectType{AttrTypes: attributeQueryIssuanceCriteriaAttrTypes},
		"policy":                         types.ObjectType{AttrTypes: attributeQueryPolicyAttrTypes},
	}
	var attributeQueryValue types.Object
	if response.AttributeQuery == nil {
		attributeQueryValue = types.ObjectNull(attributeQueryAttrTypes)
	} else {
		attributeQueryAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &response.AttributeQuery.AttributeContractFulfillment)
		respDiags.Append(diags...)
		attributeQueryAttributeSourcesValue, diags := attributesources.ToState(context.Background(), response.AttributeQuery.AttributeSources)
		respDiags.Append(diags...)
		attributeQueryAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, response.AttributeQuery.Attributes)
		respDiags.Append(diags...)
		attributeQueryIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), response.AttributeQuery.IssuanceCriteria)
		respDiags.Append(diags...)
		var attributeQueryPolicyValue types.Object
		if response.AttributeQuery.Policy == nil {
			attributeQueryPolicyValue = types.ObjectNull(attributeQueryPolicyAttrTypes)
		} else {
			attributeQueryPolicyValue, diags = types.ObjectValue(attributeQueryPolicyAttrTypes, map[string]attr.Value{
				"encrypt_assertion":              types.BoolPointerValue(response.AttributeQuery.Policy.EncryptAssertion),
				"require_encrypted_name_id":      types.BoolPointerValue(response.AttributeQuery.Policy.RequireEncryptedNameId),
				"require_signed_attribute_query": types.BoolPointerValue(response.AttributeQuery.Policy.RequireSignedAttributeQuery),
				"sign_assertion":                 types.BoolPointerValue(response.AttributeQuery.Policy.SignAssertion),
				"sign_response":                  types.BoolPointerValue(response.AttributeQuery.Policy.SignResponse),
			})
			respDiags.Append(diags...)
		}
		attributeQueryValue, diags = types.ObjectValue(attributeQueryAttrTypes, map[string]attr.Value{
			"attribute_contract_fulfillment": attributeQueryAttributeContractFulfillmentValue,
			"attribute_sources":              attributeQueryAttributeSourcesValue,
			"attributes":                     attributeQueryAttributesValue,
			"issuance_criteria":              attributeQueryIssuanceCriteriaValue,
			"policy":                         attributeQueryPolicyValue,
		})
		respDiags.Append(diags...)
	}

	state.AttributeQuery = attributeQueryValue
	// base_url
	state.BaseUrl = types.StringPointerValue(response.BaseUrl)
	// connection_id
	state.ConnectionId = types.StringPointerValue(response.Id)
	state.Id = types.StringPointerValue(response.Id)
	// connection_target_type
	state.ConnectionTargetType = types.StringPointerValue(response.ConnectionTargetType)
	// contact_info
	contactInfoAttrTypes := map[string]attr.Type{
		"company":    types.StringType,
		"email":      types.StringType,
		"first_name": types.StringType,
		"last_name":  types.StringType,
		"phone":      types.StringType,
	}
	var contactInfoValue types.Object
	if response.ContactInfo == nil {
		contactInfoValue = types.ObjectNull(contactInfoAttrTypes)
	} else {
		contactInfoValue, diags = types.ObjectValue(contactInfoAttrTypes, map[string]attr.Value{
			"company":    types.StringPointerValue(response.ContactInfo.Company),
			"email":      types.StringPointerValue(response.ContactInfo.Email),
			"first_name": types.StringPointerValue(response.ContactInfo.FirstName),
			"last_name":  types.StringPointerValue(response.ContactInfo.LastName),
			"phone":      types.StringPointerValue(response.ContactInfo.Phone),
		})
		respDiags.Append(diags...)
	}

	state.ContactInfo = contactInfoValue
	// creation_date
	if response.CreationDate != nil {
		state.CreationDate = types.StringValue(response.CreationDate.Format(time.RFC3339))
	} else {
		state.CreationDate = types.StringNull()
	}
	// credentials
	var credentialsValue types.Object
	if response.Credentials != nil {
		var credentialsCertsValues []attr.Value
		for _, cert := range response.Credentials.Certs {
			if state.Credentials.Attributes()["certs"] != nil {
				certMatchFound := false
				for _, certInPlan := range state.Credentials.Attributes()["certs"].(types.List).Elements() {
					x509FilePlanAttrs := certInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
					x509FileIdPlan := x509FilePlanAttrs["id"].(types.String).ValueString()
					// If the plan did not define an id for the cert, we'll have to compare the certs via the file_data.
					// Otherwise just compare the ids of the certs
					if (x509FileIdPlan == "" &&
						pemcertificates.FileDataEquivalent(x509FilePlanAttrs["file_data"].(types.String).ValueString(), cert.X509File.FileData)) ||
						(cert.X509File.Id != nil && *cert.X509File.Id == x509FileIdPlan) {
						planFileData := x509FilePlanAttrs["file_data"].(types.String)
						credentialsCertsObjValue, objDiags := connectioncert.ToState(context.Background(), planFileData, cert, &diags, isImportRead)
						respDiags.Append(objDiags...)
						credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
						certMatchFound = true
						break
					}
				}
				if !certMatchFound {
					credentialsCertsObjValue, objDiags := connectioncert.ToState(context.Background(), types.StringNull(), cert, &diags, isImportRead)
					respDiags.Append(objDiags...)
					credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
				}
			} else {
				credentialsCertsObjValue, objDiags := connectioncert.ToState(context.Background(), types.StringNull(), cert, &diags, isImportRead)
				respDiags.Append(objDiags...)
				credentialsCertsValues = append(credentialsCertsValues, credentialsCertsObjValue)
			}
		}
		credentialsCertsValue, objDiags := types.ListValue(connectioncert.ObjType(), credentialsCertsValues)
		respDiags.Append(objDiags...)
		var credentialsDecryptionKeyPairRefValue types.Object
		if response.Credentials.DecryptionKeyPairRef == nil {
			credentialsDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(context.Background(), response.Credentials.DecryptionKeyPairRef)
			respDiags.Append(objDiags...)
		}
		var credentialsInboundBackChannelAuthValue types.Object
		if response.Credentials.InboundBackChannelAuth == nil {
			credentialsInboundBackChannelAuthValue = types.ObjectNull(credentialsInboundBackChannelAuthAttrTypes)
		} else {
			var credentialsInboundBackChannelAuthCertsValue types.List
			if len(response.Credentials.InboundBackChannelAuth.Certs) > 0 {
				var credentialsInboundBackChannelAuthCertsValues []attr.Value
				for _, ibcaCert := range response.Credentials.InboundBackChannelAuth.Certs {
					if state.Credentials.Attributes()["inbound_back_channel_auth"] != nil {
						ibaCertMatch := false
						for _, ibcaCertInPlan := range state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["certs"].(types.List).Elements() {
							ibcax509FilePlanAttrs := ibcaCertInPlan.(types.Object).Attributes()["x509_file"].(types.Object).Attributes()
							ibcax509FileIdPlan := ibcax509FilePlanAttrs["id"].(types.String).ValueString()
							// If the plan did not define an id for the cert, we'll have to compare the certs via the file_data.
							// Otherwise just compare the ids of the certs
							if (ibcax509FileIdPlan == "" &&
								pemcertificates.FileDataEquivalent(ibcax509FilePlanAttrs["file_data"].(types.String).ValueString(), ibcaCert.X509File.FileData)) ||
								(ibcaCert.X509File.Id != nil && *ibcaCert.X509File.Id == ibcax509FileIdPlan) {
								planIbcaX509FileFileData := ibcax509FilePlanAttrs["file_data"].(types.String)
								planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(context.Background(), planIbcaX509FileFileData, ibcaCert, &diags, isImportRead)
								respDiags.Append(objDiags...)
								credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
								ibaCertMatch = true
								break
							}
						}
						if !ibaCertMatch {
							planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(context.Background(), types.StringNull(), ibcaCert, &diags, isImportRead)
							respDiags.Append(objDiags...)
							credentialsInboundBackChannelAuthCertsValues = append(credentialsInboundBackChannelAuthCertsValues, planIbcaX509FileFileDataCertsObjValue)
						}
					} else {
						planIbcaX509FileFileDataCertsObjValue, objDiags := connectioncert.ToState(context.Background(), types.StringNull(), ibcaCert, &diags, isImportRead)
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
			if response.Credentials.InboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password *string
				if state != nil && state.Credentials.Attributes()["inbound_back_channel_auth"] != nil && state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					passwordFromPlan := state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					password = passwordFromPlan.ValueStringPointer()
				}
				encryptedPassword := types.StringPointerValue(response.Credentials.InboundBackChannelAuth.HttpBasicCredentials.EncryptedPassword)
				if state != nil && state.Credentials.Attributes()["inbound_back_channel_auth"] != nil && state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					encryptedPasswordFromPlan := state.Credentials.Attributes()["inbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["encrypted_password"].(types.String)
					if internaltypes.IsDefined(encryptedPasswordFromPlan) {
						encryptedPassword = types.StringValue(encryptedPasswordFromPlan.ValueString())
					}
				}
				credentialsInboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password":           types.StringPointerValue(password),
					"encrypted_password": encryptedPassword,
					"username":           types.StringPointerValue(response.Credentials.InboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				respDiags.Append(objDiags...)
			}
			credentialsInboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsInboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"certs":                   credentialsInboundBackChannelAuthCertsValue,
				"digital_signature":       types.BoolPointerValue(response.Credentials.InboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials":  credentialsInboundBackChannelAuthHttpBasicCredentialsValue,
				"require_ssl":             types.BoolPointerValue(response.Credentials.InboundBackChannelAuth.RequireSsl),
				"verification_issuer_dn":  types.StringPointerValue(response.Credentials.InboundBackChannelAuth.VerificationIssuerDN),
				"verification_subject_dn": types.StringPointerValue(response.Credentials.InboundBackChannelAuth.VerificationSubjectDN),
			})
			respDiags.Append(objDiags...)
		}
		var credentialsOutboundBackChannelAuthValue types.Object
		if response.Credentials.OutboundBackChannelAuth == nil {
			credentialsOutboundBackChannelAuthValue = types.ObjectNull(credentialsOutboundBackChannelAuthAttrTypes)
		} else {
			var credentialsOutboundBackChannelAuthHttpBasicCredentialsValue types.Object
			if response.Credentials.OutboundBackChannelAuth.HttpBasicCredentials == nil {
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue = types.ObjectNull(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes)
			} else {
				var password *string
				if state != nil && state.Credentials.Attributes()["outbound_back_channel_auth"] != nil && state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					passwordFromPlan := state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["password"].(types.String)
					password = passwordFromPlan.ValueStringPointer()
				}
				encryptedPassword := types.StringPointerValue(response.Credentials.OutboundBackChannelAuth.HttpBasicCredentials.EncryptedPassword)
				if state != nil && state.Credentials.Attributes()["outbound_back_channel_auth"] != nil && state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"] != nil {
					encryptedPasswordFromPlan := state.Credentials.Attributes()["outbound_back_channel_auth"].(types.Object).Attributes()["http_basic_credentials"].(types.Object).Attributes()["encrypted_password"].(types.String)
					if internaltypes.IsDefined(encryptedPasswordFromPlan) {
						encryptedPassword = types.StringValue(encryptedPasswordFromPlan.ValueString())
					}
				}
				credentialsOutboundBackChannelAuthHttpBasicCredentialsValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthHttpBasicCredentialsAttrTypes, map[string]attr.Value{
					"password":           types.StringPointerValue(password),
					"encrypted_password": encryptedPassword,
					"username":           types.StringPointerValue(response.Credentials.OutboundBackChannelAuth.HttpBasicCredentials.Username),
				})
				respDiags.Append(objDiags...)
			}
			var credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue types.Object
			if response.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef == nil {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
			} else {
				credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue, objDiags = resourcelink.ToState(context.Background(), response.Credentials.OutboundBackChannelAuth.SslAuthKeyPairRef)
				respDiags.Append(objDiags...)
			}
			credentialsOutboundBackChannelAuthValue, objDiags = types.ObjectValue(credentialsOutboundBackChannelAuthAttrTypes, map[string]attr.Value{
				"digital_signature":      types.BoolPointerValue(response.Credentials.OutboundBackChannelAuth.DigitalSignature),
				"http_basic_credentials": credentialsOutboundBackChannelAuthHttpBasicCredentialsValue,
				"ssl_auth_key_pair_ref":  credentialsOutboundBackChannelAuthSslAuthKeyPairRefValue,
				"validate_partner_cert":  types.BoolPointerValue(response.Credentials.OutboundBackChannelAuth.ValidatePartnerCert),
			})
			respDiags.Append(objDiags...)
		}
		var credentialsSecondaryDecryptionKeyPairRefValue types.Object
		if response.Credentials.SecondaryDecryptionKeyPairRef == nil {
			credentialsSecondaryDecryptionKeyPairRefValue = types.ObjectNull(resourcelink.AttrType())
		} else {
			credentialsSecondaryDecryptionKeyPairRefValue, objDiags = resourcelink.ToState(context.Background(), response.Credentials.SecondaryDecryptionKeyPairRef)
			respDiags.Append(objDiags...)
		}
		var credentialsSigningSettingsValue types.Object
		if response.Credentials.SigningSettings == nil {
			credentialsSigningSettingsValue = types.ObjectNull(credentialsSigningSettingsAttrTypes)
		} else {
			var credentialsSigningSettingsAlternativeSigningKeyPairRefsValues []attr.Value
			for _, credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue := range response.Credentials.SigningSettings.AlternativeSigningKeyPairRefs {
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink := &client.ResourceLink{}
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink.Id = credentialsSigningSettingsAlternativeSigningKeyPairRefsResponseValue.Id
				credentialsSigningSettingsAlternativeSigningKeyPairRefsValue, objDiags := resourcelink.ToState(context.Background(), credentialsSigningSettingsAlternativeSigningKeyPairRefsValueResourceLink)
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

			credentialsSigningSettingsSigningKeyPairRefValue, objDiags := resourcelink.ToState(context.Background(), &response.Credentials.SigningSettings.SigningKeyPairRef)
			respDiags.Append(objDiags...)
			// PF will return false include_cert_in_signature as nil
			includeCertInSignature := false
			if response.Credentials.SigningSettings.IncludeCertInSignature != nil {
				includeCertInSignature = *response.Credentials.SigningSettings.IncludeCertInSignature
			}
			credentialsSigningSettingsValue, objDiags = types.ObjectValue(credentialsSigningSettingsAttrTypes, map[string]attr.Value{
				"algorithm":                         types.StringPointerValue(response.Credentials.SigningSettings.Algorithm),
				"alternative_signing_key_pair_refs": credentialsSigningSettingsAlternativeSigningKeyPairRefsValue,
				"include_cert_in_signature":         types.BoolValue(includeCertInSignature),
				"include_raw_key_in_signature":      types.BoolPointerValue(response.Credentials.SigningSettings.IncludeRawKeyInSignature),
				"signing_key_pair_ref":              credentialsSigningSettingsSigningKeyPairRefValue,
			})
			respDiags.Append(objDiags...)
		}
		credentialsValue, objDiags = types.ObjectValue(credentialsAttrTypes, map[string]attr.Value{
			"block_encryption_algorithm":        types.StringPointerValue(response.Credentials.BlockEncryptionAlgorithm),
			"certs":                             credentialsCertsValue,
			"decryption_key_pair_ref":           credentialsDecryptionKeyPairRefValue,
			"inbound_back_channel_auth":         credentialsInboundBackChannelAuthValue,
			"key_transport_algorithm":           types.StringPointerValue(response.Credentials.KeyTransportAlgorithm),
			"outbound_back_channel_auth":        credentialsOutboundBackChannelAuthValue,
			"secondary_decryption_key_pair_ref": credentialsSecondaryDecryptionKeyPairRefValue,
			"signing_settings":                  credentialsSigningSettingsValue,
			"verification_issuer_dn":            types.StringPointerValue(response.Credentials.VerificationIssuerDN),
			"verification_subject_dn":           types.StringPointerValue(response.Credentials.VerificationSubjectDN),
		})
		respDiags.Append(objDiags...)
	} else {
		credentialsValue = types.ObjectNull(credentialsAttrTypes)
	}
	state.Credentials = credentialsValue
	// default_virtual_entity_id
	state.DefaultVirtualEntityId = types.StringPointerValue(response.DefaultVirtualEntityId)
	// entity_id
	state.EntityId = types.StringValue(response.EntityId)
	// extended_properties
	extendedPropertiesAttrTypes := map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}
	extendedPropertiesElementType := types.ObjectType{AttrTypes: extendedPropertiesAttrTypes}
	var extendedPropertiesValue types.Map
	if response.ExtendedProperties == nil {
		extendedPropertiesValue = types.MapNull(extendedPropertiesElementType)
	} else {
		extendedPropertiesValues := make(map[string]attr.Value)
		for key, extendedPropertiesResponseValue := range *response.ExtendedProperties {
			extendedPropertiesValuesValue, diags := types.SetValueFrom(context.Background(), types.StringType, extendedPropertiesResponseValue.Values)
			respDiags.Append(diags...)
			extendedPropertiesValue, diags := types.ObjectValue(extendedPropertiesAttrTypes, map[string]attr.Value{
				"values": extendedPropertiesValuesValue,
			})
			respDiags.Append(diags...)
			extendedPropertiesValues[key] = extendedPropertiesValue
		}
		extendedPropertiesValue, diags = types.MapValue(extendedPropertiesElementType, extendedPropertiesValues)
		respDiags.Append(diags...)
	}

	state.ExtendedProperties = extendedPropertiesValue
	// license_connection_group
	state.LicenseConnectionGroup = types.StringPointerValue(response.LicenseConnectionGroup)
	// logging_mode
	// If the plan logging mode does not match the state logging mode, report that the error might be being controlled
	// by the `server_settings_general` resource
	if response.LoggingMode != nil && internaltypes.IsDefined(state.LoggingMode) && state.LoggingMode.ValueString() != *response.LoggingMode {
		respDiags.AddAttributeError(path.Root("logging_mode"), providererror.ConflictingValueReturnedError,
			"PingFederate returned a different value for `logging_mode` for this resource than was planned. "+
				"If `sp_connection_transaction_logging_override` is configured to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource,"+
				" `logging_mode` should be configured to the same value in this resource.")
	}
	state.LoggingMode = types.StringPointerValue(response.LoggingMode)
	// metadata_reload_settings
	metadataReloadSettingsMetadataUrlRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	metadataReloadSettingsAttrTypes := map[string]attr.Type{
		"enable_auto_metadata_update": types.BoolType,
		"metadata_url_ref":            types.ObjectType{AttrTypes: metadataReloadSettingsMetadataUrlRefAttrTypes},
	}
	var metadataReloadSettingsValue types.Object
	if response.MetadataReloadSettings == nil {
		metadataReloadSettingsValue = types.ObjectNull(metadataReloadSettingsAttrTypes)
	} else {
		metadataReloadSettingsMetadataUrlRefValue, diags := types.ObjectValue(metadataReloadSettingsMetadataUrlRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.MetadataReloadSettings.MetadataUrlRef.Id),
		})
		respDiags.Append(diags...)
		metadataReloadSettingsValue, diags = types.ObjectValue(metadataReloadSettingsAttrTypes, map[string]attr.Value{
			"enable_auto_metadata_update": types.BoolPointerValue(response.MetadataReloadSettings.EnableAutoMetadataUpdate),
			"metadata_url_ref":            metadataReloadSettingsMetadataUrlRefValue,
		})
		respDiags.Append(diags...)
	}

	state.MetadataReloadSettings = metadataReloadSettingsValue
	// name
	state.Name = types.StringValue(response.Name)
	// outbound_provision
	outboundProvisionChannelsAttributeMappingSaasFieldInfoAttrTypes := map[string]attr.Type{
		"attribute_names": types.ListType{ElemType: types.StringType},
		"character_case":  types.StringType,
		"create_only":     types.BoolType,
		"default_value":   types.StringType,
		"expression":      types.StringType,
		"masked":          types.BoolType,
		"parser":          types.StringType,
		"trim":            types.BoolType,
	}
	outboundProvisionChannelsAttributeMappingAttrTypes := map[string]attr.Type{
		"field_name":      types.StringType,
		"saas_field_info": types.ObjectType{AttrTypes: outboundProvisionChannelsAttributeMappingSaasFieldInfoAttrTypes},
	}
	outboundProvisionChannelsAttributeMappingElementType := types.ObjectType{AttrTypes: outboundProvisionChannelsAttributeMappingAttrTypes}
	outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes := map[string]attr.Type{
		"account_status_algorithm":      types.StringType,
		"account_status_attribute_name": types.StringType,
		"default_status":                types.BoolType,
		"flag_comparison_status":        types.BoolType,
		"flag_comparison_value":         types.StringType,
	}
	outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes := map[string]attr.Type{
		"changed_users_algorithm":   types.StringType,
		"group_object_class":        types.StringType,
		"time_stamp_attribute_name": types.StringType,
		"user_object_class":         types.StringType,
		"usn_attribute_name":        types.StringType,
	}
	outboundProvisionChannelsChannelSourceDataSourceAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes := map[string]attr.Type{
		"group_member_attribute_name":    types.StringType,
		"member_of_group_attribute_name": types.StringType,
	}
	outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes := map[string]attr.Type{
		"filter":        types.StringType,
		"group_dn":      types.StringType,
		"nested_search": types.BoolType,
	}
	outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes := map[string]attr.Type{
		"filter":        types.StringType,
		"group_dn":      types.StringType,
		"nested_search": types.BoolType,
	}
	outboundProvisionChannelsChannelSourceAttrTypes := map[string]attr.Type{
		"account_management_settings": types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes},
		"base_dn":                     types.StringType,
		"change_detection_settings":   types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes},
		"data_source":                 types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceDataSourceAttrTypes},
		"group_membership_detection":  types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes},
		"group_source_location":       types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes},
		"guid_attribute_name":         types.StringType,
		"guid_binary":                 types.BoolType,
		"user_source_location":        types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes},
	}
	outboundProvisionChannelsAttrTypes := map[string]attr.Type{
		"active":                types.BoolType,
		"attribute_mapping":     types.SetType{ElemType: outboundProvisionChannelsAttributeMappingElementType},
		"attribute_mapping_all": types.SetType{ElemType: outboundProvisionChannelsAttributeMappingElementType},
		"channel_source":        types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceAttrTypes},
		"max_threads":           types.Int64Type,
		"name":                  types.StringType,
		"timeout":               types.Int64Type,
	}
	outboundProvisionChannelsElementType := types.ObjectType{AttrTypes: outboundProvisionChannelsAttrTypes}
	outboundProvisionCustomSchemaAttributesAttrTypes := map[string]attr.Type{
		"multi_valued":   types.BoolType,
		"name":           types.StringType,
		"sub_attributes": types.SetType{ElemType: types.StringType},
		"types":          types.SetType{ElemType: types.StringType},
	}
	outboundProvisionCustomSchemaAttributesElementType := types.ObjectType{AttrTypes: outboundProvisionCustomSchemaAttributesAttrTypes}
	outboundProvisionCustomSchemaAttrTypes := map[string]attr.Type{
		"attributes": types.SetType{ElemType: outboundProvisionCustomSchemaAttributesElementType},
		"namespace":  types.StringType,
	}
	outboundProvisionTargetSettingsAttrTypes := map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}
	outboundProvisionTargetSettingsElementType := types.ObjectType{AttrTypes: outboundProvisionTargetSettingsAttrTypes}
	outboundProvisionAttrTypes := map[string]attr.Type{
		"channels":            types.ListType{ElemType: outboundProvisionChannelsElementType},
		"custom_schema":       types.ObjectType{AttrTypes: outboundProvisionCustomSchemaAttrTypes},
		"target_settings":     types.SetType{ElemType: outboundProvisionTargetSettingsElementType},
		"target_settings_all": types.SetType{ElemType: outboundProvisionTargetSettingsElementType},
		"type":                types.StringType,
	}
	var outboundProvisionValue types.Object
	if response.OutboundProvision == nil {
		outboundProvisionValue = types.ObjectNull(outboundProvisionAttrTypes)
	} else {
		var outboundProvisionChannelsValues []attr.Value
		for _, outboundProvisionChannelsResponseValue := range response.OutboundProvision.Channels {
			attributeMapping, attributeMappingAll, diags := state.buildAttributeMappingAttrs(
				outboundProvisionChannelsResponseValue.Name, outboundProvisionChannelsResponseValue.AttributeMapping, isImportRead)
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceAccountManagementSettingsValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes, map[string]attr.Value{
				"account_status_algorithm":      types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.AccountStatusAlgorithm),
				"account_status_attribute_name": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.AccountStatusAttributeName),
				"default_status":                types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.DefaultStatus),
				"flag_comparison_status":        types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.FlagComparisonStatus),
				"flag_comparison_value":         types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.FlagComparisonValue),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceChangeDetectionSettingsValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes, map[string]attr.Value{
				"changed_users_algorithm":   types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.ChangedUsersAlgorithm),
				"group_object_class":        types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.GroupObjectClass),
				"time_stamp_attribute_name": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.TimeStampAttributeName),
				"user_object_class":         types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.UserObjectClass),
				"usn_attribute_name":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.UsnAttributeName),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceDataSourceValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceDataSourceAttrTypes, map[string]attr.Value{
				"id": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.DataSource.Id),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceGroupMembershipDetectionValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes, map[string]attr.Value{
				"group_member_attribute_name":    types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupMembershipDetection.GroupMemberAttributeName),
				"member_of_group_attribute_name": types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupMembershipDetection.MemberOfGroupAttributeName),
			})
			respDiags.Append(diags...)
			var outboundProvisionChannelsChannelSourceGroupSourceLocationValue types.Object
			if outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation == nil {
				outboundProvisionChannelsChannelSourceGroupSourceLocationValue = types.ObjectNull(outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes)
			} else {
				outboundProvisionChannelsChannelSourceGroupSourceLocationValue, diags = types.ObjectValue(outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes, map[string]attr.Value{
					"filter":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.Filter),
					"group_dn":      types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.GroupDN),
					"nested_search": types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.NestedSearch),
				})
				respDiags.Append(diags...)
			}
			outboundProvisionChannelsChannelSourceUserSourceLocationValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes, map[string]attr.Value{
				"filter":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.Filter),
				"group_dn":      types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.GroupDN),
				"nested_search": types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.NestedSearch),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceAttrTypes, map[string]attr.Value{
				"account_management_settings": outboundProvisionChannelsChannelSourceAccountManagementSettingsValue,
				"base_dn":                     types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.BaseDn),
				"change_detection_settings":   outboundProvisionChannelsChannelSourceChangeDetectionSettingsValue,
				"data_source":                 outboundProvisionChannelsChannelSourceDataSourceValue,
				"group_membership_detection":  outboundProvisionChannelsChannelSourceGroupMembershipDetectionValue,
				"group_source_location":       outboundProvisionChannelsChannelSourceGroupSourceLocationValue,
				"guid_attribute_name":         types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.GuidAttributeName),
				"guid_binary":                 types.BoolValue(outboundProvisionChannelsResponseValue.ChannelSource.GuidBinary),
				"user_source_location":        outboundProvisionChannelsChannelSourceUserSourceLocationValue,
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsValue, diags := types.ObjectValue(outboundProvisionChannelsAttrTypes, map[string]attr.Value{
				"active":                types.BoolValue(outboundProvisionChannelsResponseValue.Active),
				"attribute_mapping":     attributeMapping,
				"attribute_mapping_all": attributeMappingAll,
				"channel_source":        outboundProvisionChannelsChannelSourceValue,
				"max_threads":           types.Int64Value(outboundProvisionChannelsResponseValue.MaxThreads),
				"name":                  types.StringValue(outboundProvisionChannelsResponseValue.Name),
				"timeout":               types.Int64Value(outboundProvisionChannelsResponseValue.Timeout),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsValues = append(outboundProvisionChannelsValues, outboundProvisionChannelsValue)
		}
		outboundProvisionChannelsValue, diags := types.ListValue(outboundProvisionChannelsElementType, outboundProvisionChannelsValues)
		respDiags.Append(diags...)
		var outboundProvisionCustomSchemaValue types.Object
		if response.OutboundProvision.CustomSchema == nil {
			outboundProvisionCustomSchemaValue = types.ObjectNull(outboundProvisionCustomSchemaAttrTypes)
		} else {
			var outboundProvisionCustomSchemaAttributesValues []attr.Value
			for _, outboundProvisionCustomSchemaAttributesResponseValue := range response.OutboundProvision.CustomSchema.Attributes {
				outboundProvisionCustomSchemaAttributesSubAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, outboundProvisionCustomSchemaAttributesResponseValue.SubAttributes)
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesTypesValue, diags := types.SetValueFrom(context.Background(), types.StringType, outboundProvisionCustomSchemaAttributesResponseValue.Types)
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesValue, diags := types.ObjectValue(outboundProvisionCustomSchemaAttributesAttrTypes, map[string]attr.Value{
					"multi_valued":   types.BoolPointerValue(outboundProvisionCustomSchemaAttributesResponseValue.MultiValued),
					"name":           types.StringPointerValue(outboundProvisionCustomSchemaAttributesResponseValue.Name),
					"sub_attributes": outboundProvisionCustomSchemaAttributesSubAttributesValue,
					"types":          outboundProvisionCustomSchemaAttributesTypesValue,
				})
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesValues = append(outboundProvisionCustomSchemaAttributesValues, outboundProvisionCustomSchemaAttributesValue)
			}
			outboundProvisionCustomSchemaAttributesValue, diags := types.SetValue(outboundProvisionCustomSchemaAttributesElementType, outboundProvisionCustomSchemaAttributesValues)
			respDiags.Append(diags...)
			outboundProvisionCustomSchemaValue, diags = types.ObjectValue(outboundProvisionCustomSchemaAttrTypes, map[string]attr.Value{
				"attributes": outboundProvisionCustomSchemaAttributesValue,
				"namespace":  types.StringPointerValue(response.OutboundProvision.CustomSchema.Namespace),
			})
			respDiags.Append(diags...)
		}
		targetSettings, targetSettingsAll, diags := state.buildTargetSettingsAttrs(response.OutboundProvision.TargetSettings, isImportRead)
		respDiags.Append(diags...)
		outboundProvisionValue, diags = types.ObjectValue(outboundProvisionAttrTypes, map[string]attr.Value{
			"channels":            outboundProvisionChannelsValue,
			"custom_schema":       outboundProvisionCustomSchemaValue,
			"target_settings":     targetSettings,
			"target_settings_all": targetSettingsAll,
			"type":                types.StringValue(response.OutboundProvision.Type),
		})
		respDiags.Append(diags...)
	}

	state.OutboundProvision = outboundProvisionValue
	// sp_browser_sso
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"masked":    types.BoolType,
		"name":      types.StringType,
		"pseudonym": types.BoolType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"masked":    types.BoolType,
		"name":      types.StringType,
		"pseudonym": types.BoolType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":           types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType},
		"extended_attributes":       types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType},
		"mask_ognl_values":          types.BoolType,
		"unique_user_key_attribute": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesElementType},
		"issuance_criteria":              types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaAttrTypes},
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationAttrTypes := pluginconfiguration.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes := map[string]attr.Type{
		"attribute_contract":    types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes},
		"attribute_mapping":     types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes},
		"authn_ctx_class_ref":   types.StringType,
		"configuration":         types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationAttrTypes},
		"id":                    types.StringType,
		"name":                  types.StringType,
		"parent_ref":            types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes},
		"plugin_descriptor_ref": types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes},
	}
	spBrowserSsoAdapterMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAdapterMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAdapterMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAdapterMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttributeSourcesAttrTypes}
	spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAdapterMappingsAttrTypes := map[string]attr.Type{
		"abort_sso_transaction_as_fail_safe": types.BoolType,
		"adapter_override_settings":          types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes},
		"attribute_contract_fulfillment":     types.MapType{ElemType: spBrowserSsoAdapterMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":                  types.SetType{ElemType: spBrowserSsoAdapterMappingsAttributeSourcesElementType},
		"idp_adapter_ref":                    types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes},
		"issuance_criteria":                  types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsIssuanceCriteriaAttrTypes},
		"restrict_virtual_entity_ids":        types.BoolType,
		"restricted_virtual_entity_ids":      types.SetType{ElemType: types.StringType},
	}
	spBrowserSsoAdapterMappingsElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttrTypes}
	spBrowserSsoArtifactResolverLocationsAttrTypes := map[string]attr.Type{
		"index": types.Int64Type,
		"url":   types.StringType,
	}
	spBrowserSsoArtifactResolverLocationsElementType := types.ObjectType{AttrTypes: spBrowserSsoArtifactResolverLocationsAttrTypes}
	spBrowserSsoArtifactAttrTypes := map[string]attr.Type{
		"lifetime":           types.Int64Type,
		"resolver_locations": types.SetType{ElemType: spBrowserSsoArtifactResolverLocationsElementType},
		"source_id":          types.StringType,
	}
	spBrowserSsoAssertionLifetimeAttrTypes := map[string]attr.Type{
		"minutes_after":  types.Int64Type,
		"minutes_before": types.Int64Type,
	}
	spBrowserSsoAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"name_format": types.StringType,
	}
	spBrowserSsoAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAttributeContractCoreAttributesAttrTypes}
	spBrowserSsoAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"name_format": types.StringType,
	}
	spBrowserSsoAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAttributeContractExtendedAttributesAttrTypes}
	spBrowserSsoAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: spBrowserSsoAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: spBrowserSsoAttributeContractExtendedAttributesElementType},
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesAttrTypes}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes := map[string]attr.Type{
		"abort_sso_transaction_as_fail_safe": types.BoolType,
		"attribute_contract_fulfillment":     types.MapType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":                  types.SetType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesElementType},
		"authentication_policy_contract_ref": types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes},
		"issuance_criteria":                  types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaAttrTypes},
		"restrict_virtual_entity_ids":        types.BoolType,
		"restricted_virtual_entity_ids":      types.SetType{ElemType: types.StringType},
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes}
	spBrowserSsoEncryptionPolicyAttrTypes := map[string]attr.Type{
		"encrypt_assertion":             types.BoolType,
		"encrypt_slo_subject_name_id":   types.BoolType,
		"encrypted_attributes":          types.SetType{ElemType: types.StringType},
		"slo_subject_name_id_encrypted": types.BoolType,
	}
	spBrowserSsoMessageCustomizationsAttrTypes := map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}
	spBrowserSsoMessageCustomizationsElementType := types.ObjectType{AttrTypes: spBrowserSsoMessageCustomizationsAttrTypes}
	spBrowserSsoSloServiceEndpointsAttrTypes := map[string]attr.Type{
		"binding":      types.StringType,
		"response_url": types.StringType,
		"url":          types.StringType,
	}
	spBrowserSsoSloServiceEndpointsElementType := types.ObjectType{AttrTypes: spBrowserSsoSloServiceEndpointsAttrTypes}
	spBrowserSsoSsoServiceEndpointsAttrTypes := map[string]attr.Type{
		"binding":    types.StringType,
		"index":      types.Int64Type,
		"is_default": types.BoolType,
		"url":        types.StringType,
	}
	spBrowserSsoSsoServiceEndpointsElementType := types.ObjectType{AttrTypes: spBrowserSsoSsoServiceEndpointsAttrTypes}
	spBrowserSsoUrlWhitelistEntriesAttrTypes := map[string]attr.Type{
		"allow_query_and_fragment": types.BoolType,
		"require_https":            types.BoolType,
		"valid_domain":             types.StringType,
		"valid_path":               types.StringType,
	}
	spBrowserSsoUrlWhitelistEntriesElementType := types.ObjectType{AttrTypes: spBrowserSsoUrlWhitelistEntriesAttrTypes}
	spBrowserSsoAttrTypes := map[string]attr.Type{
		"adapter_mappings":              types.SetType{ElemType: spBrowserSsoAdapterMappingsElementType},
		"always_sign_artifact_response": types.BoolType,
		"artifact":                      types.ObjectType{AttrTypes: spBrowserSsoArtifactAttrTypes},
		"assertion_lifetime":            types.ObjectType{AttrTypes: spBrowserSsoAssertionLifetimeAttrTypes},
		"attribute_contract":            types.ObjectType{AttrTypes: spBrowserSsoAttributeContractAttrTypes},
		"authentication_policy_contract_assertion_mappings": types.SetType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType},
		"default_target_url":            types.StringType,
		"enabled_profiles":              types.SetType{ElemType: types.StringType},
		"encryption_policy":             types.ObjectType{AttrTypes: spBrowserSsoEncryptionPolicyAttrTypes},
		"incoming_bindings":             types.SetType{ElemType: types.StringType},
		"message_customizations":        types.SetType{ElemType: spBrowserSsoMessageCustomizationsElementType},
		"protocol":                      types.StringType,
		"require_signed_authn_requests": types.BoolType,
		"sign_assertions":               types.BoolType,
		"sign_response_as_required":     types.BoolType,
		"slo_service_endpoints":         types.SetType{ElemType: spBrowserSsoSloServiceEndpointsElementType},
		"sp_saml_identity_mapping":      types.StringType,
		"sp_ws_fed_identity_mapping":    types.StringType,
		"sso_application_endpoint":      types.StringType,
		"sso_service_endpoints":         types.SetType{ElemType: spBrowserSsoSsoServiceEndpointsElementType},
		"url_whitelist_entries":         types.SetType{ElemType: spBrowserSsoUrlWhitelistEntriesElementType},
		"ws_fed_token_type":             types.StringType,
		"ws_trust_version":              types.StringType,
	}
	var spBrowserSsoValue types.Object
	if response.SpBrowserSso == nil {
		spBrowserSsoValue = types.ObjectNull(spBrowserSsoAttrTypes)
	} else {
		var spBrowserSsoAdapterMappingsValues []attr.Value
		for adapterMappingIndex, spBrowserSsoAdapterMappingsResponseValue := range response.SpBrowserSso.AdapterMappings {
			var spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue types.Object
			if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings == nil {
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes)
			} else {
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes)
				} else {
					var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues []attr.Value
					for _, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue := range spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.CoreAttributes {
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
							"masked":    types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Masked),
							"name":      types.StringValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Name),
							"pseudonym": types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Pseudonym),
						})
						respDiags.Append(diags...)
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues = append(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue)
					}
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.SetValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues)
					respDiags.Append(diags...)
					var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues []attr.Value
					for _, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue := range spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.ExtendedAttributes {
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
							"masked":    types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Masked),
							"name":      types.StringValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Name),
							"pseudonym": types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Pseudonym),
						})
						respDiags.Append(diags...)
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues = append(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue)
					}
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.SetValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes, map[string]attr.Value{
						"core_attributes":           spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue,
						"extended_attributes":       spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue,
						"mask_ognl_values":          types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.MaskOgnlValues),
						"unique_user_key_attribute": types.StringPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.UniqueUserKeyAttribute),
					})
					respDiags.Append(diags...)
				}
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes)
				} else {
					contractFulfillment := spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.AttributeContractFulfillment
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.AttributeSources)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.IssuanceCriteria)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes, map[string]attr.Value{
						"attribute_contract_fulfillment": spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentValue,
						"attribute_sources":              spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesValue,
						"issuance_criteria":              spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaValue,
					})
					respDiags.Append(diags...)
				}
				overrideSettingsConfiguration := spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Configuration
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue, diags := pluginconfiguration.ToState(state.getSpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfiguration(adapterMappingIndex), &overrideSettingsConfiguration, isImportRead)
				respDiags.Append(diags...)
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes)
				} else {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes, map[string]attr.Value{
						"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef.Id),
					})
					respDiags.Append(diags...)
				}
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes, map[string]attr.Value{
					"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.PluginDescriptorRef.Id),
				})
				respDiags.Append(diags...)
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes, map[string]attr.Value{
					"attribute_contract":    spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue,
					"attribute_mapping":     spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue,
					"authn_ctx_class_ref":   types.StringPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AuthnCtxClassRef),
					"configuration":         spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue,
					"id":                    types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Id),
					"name":                  types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Name),
					"parent_ref":            spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue,
					"plugin_descriptor_ref": spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue,
				})
				respDiags.Append(diags...)
			}
			contractFulfillment := spBrowserSsoAdapterMappingsResponseValue.AttributeContractFulfillment
			spBrowserSsoAdapterMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			var spBrowserSsoAdapterMappingsIdpAdapterRefValue types.Object
			if spBrowserSsoAdapterMappingsResponseValue.IdpAdapterRef == nil {
				spBrowserSsoAdapterMappingsIdpAdapterRefValue = types.ObjectNull(spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes)
			} else {
				spBrowserSsoAdapterMappingsIdpAdapterRefValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes, map[string]attr.Value{
					"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.IdpAdapterRef.Id),
				})
				respDiags.Append(diags...)
			}
			spBrowserSsoAdapterMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, spBrowserSsoAdapterMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAttrTypes, map[string]attr.Value{
				"abort_sso_transaction_as_fail_safe": types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.AbortSsoTransactionAsFailSafe),
				"adapter_override_settings":          spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue,
				"attribute_contract_fulfillment":     spBrowserSsoAdapterMappingsAttributeContractFulfillmentValue,
				"attribute_sources":                  spBrowserSsoAdapterMappingsAttributeSourcesValue,
				"idp_adapter_ref":                    spBrowserSsoAdapterMappingsIdpAdapterRefValue,
				"issuance_criteria":                  spBrowserSsoAdapterMappingsIssuanceCriteriaValue,
				"restrict_virtual_entity_ids":        types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.RestrictVirtualEntityIds),
				"restricted_virtual_entity_ids":      spBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsValues = append(spBrowserSsoAdapterMappingsValues, spBrowserSsoAdapterMappingsValue)
		}
		spBrowserSsoAdapterMappingsValue, diags := types.SetValue(spBrowserSsoAdapterMappingsElementType, spBrowserSsoAdapterMappingsValues)
		respDiags.Append(diags...)
		var spBrowserSsoArtifactValue types.Object
		if response.SpBrowserSso.Artifact == nil {
			spBrowserSsoArtifactValue = types.ObjectNull(spBrowserSsoArtifactAttrTypes)
		} else {
			var spBrowserSsoArtifactResolverLocationsValues []attr.Value
			for _, spBrowserSsoArtifactResolverLocationsResponseValue := range response.SpBrowserSso.Artifact.ResolverLocations {
				spBrowserSsoArtifactResolverLocationsValue, diags := types.ObjectValue(spBrowserSsoArtifactResolverLocationsAttrTypes, map[string]attr.Value{
					"index": types.Int64Value(spBrowserSsoArtifactResolverLocationsResponseValue.Index),
					"url":   types.StringValue(spBrowserSsoArtifactResolverLocationsResponseValue.Url),
				})
				respDiags.Append(diags...)
				spBrowserSsoArtifactResolverLocationsValues = append(spBrowserSsoArtifactResolverLocationsValues, spBrowserSsoArtifactResolverLocationsValue)
			}
			spBrowserSsoArtifactResolverLocationsValue, diags := types.SetValue(spBrowserSsoArtifactResolverLocationsElementType, spBrowserSsoArtifactResolverLocationsValues)
			respDiags.Append(diags...)
			spBrowserSsoArtifactValue, diags = types.ObjectValue(spBrowserSsoArtifactAttrTypes, map[string]attr.Value{
				"lifetime":           types.Int64PointerValue(response.SpBrowserSso.Artifact.Lifetime),
				"resolver_locations": spBrowserSsoArtifactResolverLocationsValue,
				"source_id":          types.StringPointerValue(response.SpBrowserSso.Artifact.SourceId),
			})
			respDiags.Append(diags...)
		}
		spBrowserSsoAssertionLifetimeValue, diags := types.ObjectValue(spBrowserSsoAssertionLifetimeAttrTypes, map[string]attr.Value{
			"minutes_after":  types.Int64Value(response.SpBrowserSso.AssertionLifetime.MinutesAfter),
			"minutes_before": types.Int64Value(response.SpBrowserSso.AssertionLifetime.MinutesBefore),
		})
		respDiags.Append(diags...)
		var spBrowserSsoAttributeContractCoreAttributesValues []attr.Value
		for _, spBrowserSsoAttributeContractCoreAttributesResponseValue := range response.SpBrowserSso.AttributeContract.CoreAttributes {
			spBrowserSsoAttributeContractCoreAttributesValue, diags := types.ObjectValue(spBrowserSsoAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name":        types.StringValue(spBrowserSsoAttributeContractCoreAttributesResponseValue.Name),
				"name_format": types.StringPointerValue(spBrowserSsoAttributeContractCoreAttributesResponseValue.NameFormat),
			})
			respDiags.Append(diags...)
			spBrowserSsoAttributeContractCoreAttributesValues = append(spBrowserSsoAttributeContractCoreAttributesValues, spBrowserSsoAttributeContractCoreAttributesValue)
		}
		spBrowserSsoAttributeContractCoreAttributesValue, diags := types.SetValue(spBrowserSsoAttributeContractCoreAttributesElementType, spBrowserSsoAttributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var spBrowserSsoAttributeContractExtendedAttributesValues []attr.Value
		for _, spBrowserSsoAttributeContractExtendedAttributesResponseValue := range response.SpBrowserSso.AttributeContract.ExtendedAttributes {
			spBrowserSsoAttributeContractExtendedAttributesValue, diags := types.ObjectValue(spBrowserSsoAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name":        types.StringValue(spBrowserSsoAttributeContractExtendedAttributesResponseValue.Name),
				"name_format": types.StringPointerValue(spBrowserSsoAttributeContractExtendedAttributesResponseValue.NameFormat),
			})
			respDiags.Append(diags...)
			spBrowserSsoAttributeContractExtendedAttributesValues = append(spBrowserSsoAttributeContractExtendedAttributesValues, spBrowserSsoAttributeContractExtendedAttributesValue)
		}
		spBrowserSsoAttributeContractExtendedAttributesValue, diags := types.SetValue(spBrowserSsoAttributeContractExtendedAttributesElementType, spBrowserSsoAttributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		spBrowserSsoAttributeContractValue, diags := types.ObjectValue(spBrowserSsoAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     spBrowserSsoAttributeContractCoreAttributesValue,
			"extended_attributes": spBrowserSsoAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
		var spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues []attr.Value
		for _, spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue := range response.SpBrowserSso.AuthenticationPolicyContractAssertionMappings {
			contractFulfillment := spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AttributeContractFulfillment
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue, diags := types.ObjectValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AuthenticationPolicyContractRef.Id),
			})
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue, diags := types.ObjectValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes, map[string]attr.Value{
				"abort_sso_transaction_as_fail_safe": types.BoolPointerValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AbortSsoTransactionAsFailSafe),
				"attribute_contract_fulfillment":     spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentValue,
				"attribute_sources":                  spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesValue,
				"authentication_policy_contract_ref": spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue,
				"issuance_criteria":                  spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaValue,
				"restrict_virtual_entity_ids":        types.BoolPointerValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.RestrictVirtualEntityIds),
				"restricted_virtual_entity_ids":      spBrowserSsoAuthenticationPolicyContractAssertionMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues = append(spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues, spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue)
		}
		spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue, diags := types.SetValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType, spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues)
		respDiags.Append(diags...)
		spBrowserSsoEnabledProfilesValue, diags := types.SetValueFrom(context.Background(), types.StringType, response.SpBrowserSso.EnabledProfiles)
		respDiags.Append(diags...)
		var spBrowserSsoEncryptionPolicyValue types.Object
		if response.SpBrowserSso.EncryptionPolicy == nil {
			spBrowserSsoEncryptionPolicyValue = types.ObjectNull(spBrowserSsoEncryptionPolicyAttrTypes)
		} else {
			spBrowserSsoEncryptionPolicyEncryptedAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, response.SpBrowserSso.EncryptionPolicy.EncryptedAttributes)
			respDiags.Append(diags...)
			spBrowserSsoEncryptionPolicyValue, diags = types.ObjectValue(spBrowserSsoEncryptionPolicyAttrTypes, map[string]attr.Value{
				"encrypt_assertion":             types.BoolPointerValue(response.SpBrowserSso.EncryptionPolicy.EncryptAssertion),
				"encrypt_slo_subject_name_id":   types.BoolPointerValue(response.SpBrowserSso.EncryptionPolicy.EncryptSloSubjectNameId),
				"encrypted_attributes":          spBrowserSsoEncryptionPolicyEncryptedAttributesValue,
				"slo_subject_name_id_encrypted": types.BoolPointerValue(response.SpBrowserSso.EncryptionPolicy.SloSubjectNameIDEncrypted),
			})
			respDiags.Append(diags...)
		}
		spBrowserSsoIncomingBindingsValue, diags := types.SetValueFrom(context.Background(), types.StringType, response.SpBrowserSso.IncomingBindings)
		respDiags.Append(diags...)
		var spBrowserSsoMessageCustomizationsValues []attr.Value
		for _, spBrowserSsoMessageCustomizationsResponseValue := range response.SpBrowserSso.MessageCustomizations {
			spBrowserSsoMessageCustomizationsValue, diags := types.ObjectValue(spBrowserSsoMessageCustomizationsAttrTypes, map[string]attr.Value{
				"context_name":       types.StringPointerValue(spBrowserSsoMessageCustomizationsResponseValue.ContextName),
				"message_expression": types.StringPointerValue(spBrowserSsoMessageCustomizationsResponseValue.MessageExpression),
			})
			respDiags.Append(diags...)
			spBrowserSsoMessageCustomizationsValues = append(spBrowserSsoMessageCustomizationsValues, spBrowserSsoMessageCustomizationsValue)
		}
		spBrowserSsoMessageCustomizationsValue, diags := types.SetValue(spBrowserSsoMessageCustomizationsElementType, spBrowserSsoMessageCustomizationsValues)
		respDiags.Append(diags...)
		var spBrowserSsoSloServiceEndpointsValues []attr.Value
		for _, spBrowserSsoSloServiceEndpointsResponseValue := range response.SpBrowserSso.SloServiceEndpoints {
			spBrowserSsoSloServiceEndpointsValue, diags := types.ObjectValue(spBrowserSsoSloServiceEndpointsAttrTypes, map[string]attr.Value{
				"binding":      types.StringPointerValue(spBrowserSsoSloServiceEndpointsResponseValue.Binding),
				"response_url": types.StringPointerValue(spBrowserSsoSloServiceEndpointsResponseValue.ResponseUrl),
				"url":          types.StringValue(spBrowserSsoSloServiceEndpointsResponseValue.Url),
			})
			respDiags.Append(diags...)
			spBrowserSsoSloServiceEndpointsValues = append(spBrowserSsoSloServiceEndpointsValues, spBrowserSsoSloServiceEndpointsValue)
		}
		spBrowserSsoSloServiceEndpointsValue, diags := types.SetValue(spBrowserSsoSloServiceEndpointsElementType, spBrowserSsoSloServiceEndpointsValues)
		respDiags.Append(diags...)
		var spBrowserSsoSsoServiceEndpointsValues []attr.Value
		for _, spBrowserSsoSsoServiceEndpointsResponseValue := range response.SpBrowserSso.SsoServiceEndpoints {
			// PF will return nil for false for the is_default boolean
			isDefault := types.BoolValue(false)
			if spBrowserSsoSsoServiceEndpointsResponseValue.IsDefault != nil {
				isDefault = types.BoolPointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.IsDefault)
			}
			spBrowserSsoSsoServiceEndpointsValue, diags := types.ObjectValue(spBrowserSsoSsoServiceEndpointsAttrTypes, map[string]attr.Value{
				"binding":    types.StringPointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.Binding),
				"index":      types.Int64PointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.Index),
				"is_default": isDefault,
				"url":        types.StringValue(spBrowserSsoSsoServiceEndpointsResponseValue.Url),
			})
			respDiags.Append(diags...)
			spBrowserSsoSsoServiceEndpointsValues = append(spBrowserSsoSsoServiceEndpointsValues, spBrowserSsoSsoServiceEndpointsValue)
		}
		spBrowserSsoSsoServiceEndpointsValue, diags := types.SetValue(spBrowserSsoSsoServiceEndpointsElementType, spBrowserSsoSsoServiceEndpointsValues)
		respDiags.Append(diags...)
		var spBrowserSsoUrlWhitelistEntriesValue types.Set
		if response.SpBrowserSso.UrlWhitelistEntries == nil {
			spBrowserSsoUrlWhitelistEntriesValue = types.SetNull(spBrowserSsoUrlWhitelistEntriesElementType)
		} else {
			var spBrowserSsoUrlWhitelistEntriesValues []attr.Value
			for _, spBrowserSsoUrlWhitelistEntriesResponseValue := range response.SpBrowserSso.UrlWhitelistEntries {
				spBrowserSsoUrlWhitelistEntriesValue, diags := types.ObjectValue(spBrowserSsoUrlWhitelistEntriesAttrTypes, map[string]attr.Value{
					"allow_query_and_fragment": types.BoolPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.AllowQueryAndFragment),
					"require_https":            types.BoolPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.RequireHttps),
					"valid_domain":             types.StringPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.ValidDomain),
					"valid_path":               types.StringPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.ValidPath),
				})
				respDiags.Append(diags...)
				spBrowserSsoUrlWhitelistEntriesValues = append(spBrowserSsoUrlWhitelistEntriesValues, spBrowserSsoUrlWhitelistEntriesValue)
			}
			spBrowserSsoUrlWhitelistEntriesValue, diags = types.SetValue(spBrowserSsoUrlWhitelistEntriesElementType, spBrowserSsoUrlWhitelistEntriesValues)
			respDiags.Append(diags...)
		}

		// always_sign_artifact_response, sign_assertions, require_signed_authn_requests can be returned as nil when set to false
		var alwaysSignArtifactResponse, signAssertions, requireSignedAuthnRequests bool
		if response.SpBrowserSso.AlwaysSignArtifactResponse != nil {
			alwaysSignArtifactResponse = *response.SpBrowserSso.AlwaysSignArtifactResponse
		}
		if response.SpBrowserSso.SignAssertions != nil {
			signAssertions = *response.SpBrowserSso.SignAssertions
		}
		if response.SpBrowserSso.RequireSignedAuthnRequests != nil {
			requireSignedAuthnRequests = *response.SpBrowserSso.RequireSignedAuthnRequests
		}

		spBrowserSsoValue, diags = types.ObjectValue(spBrowserSsoAttrTypes, map[string]attr.Value{
			"adapter_mappings":              spBrowserSsoAdapterMappingsValue,
			"always_sign_artifact_response": types.BoolValue(alwaysSignArtifactResponse),
			"artifact":                      spBrowserSsoArtifactValue,
			"assertion_lifetime":            spBrowserSsoAssertionLifetimeValue,
			"attribute_contract":            spBrowserSsoAttributeContractValue,
			"authentication_policy_contract_assertion_mappings": spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue,
			"default_target_url":            types.StringPointerValue(response.SpBrowserSso.DefaultTargetUrl),
			"enabled_profiles":              spBrowserSsoEnabledProfilesValue,
			"encryption_policy":             spBrowserSsoEncryptionPolicyValue,
			"incoming_bindings":             spBrowserSsoIncomingBindingsValue,
			"message_customizations":        spBrowserSsoMessageCustomizationsValue,
			"protocol":                      types.StringValue(response.SpBrowserSso.Protocol),
			"require_signed_authn_requests": types.BoolValue(requireSignedAuthnRequests),
			"sign_assertions":               types.BoolValue(signAssertions),
			"sign_response_as_required":     types.BoolPointerValue(response.SpBrowserSso.SignResponseAsRequired),
			"slo_service_endpoints":         spBrowserSsoSloServiceEndpointsValue,
			"sp_saml_identity_mapping":      types.StringPointerValue(response.SpBrowserSso.SpSamlIdentityMapping),
			"sp_ws_fed_identity_mapping":    types.StringPointerValue(response.SpBrowserSso.SpWsFedIdentityMapping),
			"sso_application_endpoint":      types.StringPointerValue(response.SpBrowserSso.SsoApplicationEndpoint),
			"sso_service_endpoints":         spBrowserSsoSsoServiceEndpointsValue,
			"url_whitelist_entries":         spBrowserSsoUrlWhitelistEntriesValue,
			"ws_fed_token_type":             types.StringPointerValue(response.SpBrowserSso.WsFedTokenType),
			"ws_trust_version":              types.StringPointerValue(response.SpBrowserSso.WsTrustVersion),
		})
		respDiags.Append(diags...)
	}

	state.SpBrowserSso = spBrowserSsoValue
	// virtual_entity_ids
	state.VirtualEntityIds, diags = types.SetValueFrom(context.Background(), types.StringType, response.VirtualEntityIds)
	respDiags.Append(diags...)
	// ws_trust
	wsTrustAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name":      types.StringType,
		"namespace": types.StringType,
	}
	wsTrustAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: wsTrustAttributeContractCoreAttributesAttrTypes}
	wsTrustAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name":      types.StringType,
		"namespace": types.StringType,
	}
	wsTrustAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: wsTrustAttributeContractExtendedAttributesAttrTypes}
	wsTrustAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: wsTrustAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: wsTrustAttributeContractExtendedAttributesElementType},
	}
	wsTrustMessageCustomizationsAttrTypes := map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}
	wsTrustMessageCustomizationsElementType := types.ObjectType{AttrTypes: wsTrustMessageCustomizationsAttrTypes}
	wsTrustRequestContractRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	wsTrustTokenProcessorMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	wsTrustTokenProcessorMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttributeContractFulfillmentAttrTypes}
	wsTrustTokenProcessorMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	wsTrustTokenProcessorMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttributeSourcesAttrTypes}
	wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	wsTrustTokenProcessorMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	wsTrustTokenProcessorMappingsAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: wsTrustTokenProcessorMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: wsTrustTokenProcessorMappingsAttributeSourcesElementType},
		"idp_token_processor_ref":        types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes},
		"issuance_criteria":              types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsIssuanceCriteriaAttrTypes},
		"restricted_virtual_entity_ids":  types.SetType{ElemType: types.StringType},
	}
	wsTrustTokenProcessorMappingsElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttrTypes}
	wsTrustAttrTypes := map[string]attr.Type{
		"abort_if_not_fulfilled_from_request": types.BoolType,
		"attribute_contract":                  types.ObjectType{AttrTypes: wsTrustAttributeContractAttrTypes},
		"default_token_type":                  types.StringType,
		"encrypt_saml2_assertion":             types.BoolType,
		"generate_key":                        types.BoolType,
		"message_customizations":              types.SetType{ElemType: wsTrustMessageCustomizationsElementType},
		"minutes_after":                       types.Int64Type,
		"minutes_before":                      types.Int64Type,
		"oauth_assertion_profiles":            types.BoolType,
		"partner_service_ids":                 types.SetType{ElemType: types.StringType},
		"request_contract_ref":                types.ObjectType{AttrTypes: wsTrustRequestContractRefAttrTypes},
		"token_processor_mappings":            types.SetType{ElemType: wsTrustTokenProcessorMappingsElementType},
	}
	var wsTrustValue types.Object
	if response.WsTrust == nil {
		wsTrustValue = types.ObjectNull(wsTrustAttrTypes)
	} else {
		var wsTrustAttributeContractCoreAttributesValues []attr.Value
		for _, wsTrustAttributeContractCoreAttributesResponseValue := range response.WsTrust.AttributeContract.CoreAttributes {
			wsTrustAttributeContractCoreAttributesValue, diags := types.ObjectValue(wsTrustAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name":      types.StringValue(wsTrustAttributeContractCoreAttributesResponseValue.Name),
				"namespace": types.StringValue(wsTrustAttributeContractCoreAttributesResponseValue.Namespace),
			})
			respDiags.Append(diags...)
			wsTrustAttributeContractCoreAttributesValues = append(wsTrustAttributeContractCoreAttributesValues, wsTrustAttributeContractCoreAttributesValue)
		}
		wsTrustAttributeContractCoreAttributesValue, diags := types.SetValue(wsTrustAttributeContractCoreAttributesElementType, wsTrustAttributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var wsTrustAttributeContractExtendedAttributesValues []attr.Value
		for _, wsTrustAttributeContractExtendedAttributesResponseValue := range response.WsTrust.AttributeContract.ExtendedAttributes {
			wsTrustAttributeContractExtendedAttributesValue, diags := types.ObjectValue(wsTrustAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name":      types.StringValue(wsTrustAttributeContractExtendedAttributesResponseValue.Name),
				"namespace": types.StringValue(wsTrustAttributeContractExtendedAttributesResponseValue.Namespace),
			})
			respDiags.Append(diags...)
			wsTrustAttributeContractExtendedAttributesValues = append(wsTrustAttributeContractExtendedAttributesValues, wsTrustAttributeContractExtendedAttributesValue)
		}
		wsTrustAttributeContractExtendedAttributesValue, diags := types.SetValue(wsTrustAttributeContractExtendedAttributesElementType, wsTrustAttributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		wsTrustAttributeContractValue, diags := types.ObjectValue(wsTrustAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     wsTrustAttributeContractCoreAttributesValue,
			"extended_attributes": wsTrustAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
		var wsTrustMessageCustomizationsValues []attr.Value
		for _, wsTrustMessageCustomizationsResponseValue := range response.WsTrust.MessageCustomizations {
			wsTrustMessageCustomizationsValue, diags := types.ObjectValue(wsTrustMessageCustomizationsAttrTypes, map[string]attr.Value{
				"context_name":       types.StringPointerValue(wsTrustMessageCustomizationsResponseValue.ContextName),
				"message_expression": types.StringPointerValue(wsTrustMessageCustomizationsResponseValue.MessageExpression),
			})
			respDiags.Append(diags...)
			wsTrustMessageCustomizationsValues = append(wsTrustMessageCustomizationsValues, wsTrustMessageCustomizationsValue)
		}
		wsTrustMessageCustomizationsValue, diags := types.SetValue(wsTrustMessageCustomizationsElementType, wsTrustMessageCustomizationsValues)
		respDiags.Append(diags...)
		wsTrustPartnerServiceIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, response.WsTrust.PartnerServiceIds)
		respDiags.Append(diags...)
		var wsTrustRequestContractRefValue types.Object
		if response.WsTrust.RequestContractRef == nil {
			wsTrustRequestContractRefValue = types.ObjectNull(wsTrustRequestContractRefAttrTypes)
		} else {
			wsTrustRequestContractRefValue, diags = types.ObjectValue(wsTrustRequestContractRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(response.WsTrust.RequestContractRef.Id),
			})
			respDiags.Append(diags...)
		}
		var wsTrustTokenProcessorMappingsValues []attr.Value
		for _, wsTrustTokenProcessorMappingsResponseValue := range response.WsTrust.TokenProcessorMappings {
			contractFulfillment := wsTrustTokenProcessorMappingsResponseValue.AttributeContractFulfillment
			wsTrustTokenProcessorMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), wsTrustTokenProcessorMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsIdpTokenProcessorRefValue, diags := types.ObjectValue(wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(wsTrustTokenProcessorMappingsResponseValue.IdpTokenProcessorRef.Id),
			})
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), wsTrustTokenProcessorMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, wsTrustTokenProcessorMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsValue, diags := types.ObjectValue(wsTrustTokenProcessorMappingsAttrTypes, map[string]attr.Value{
				"attribute_contract_fulfillment": wsTrustTokenProcessorMappingsAttributeContractFulfillmentValue,
				"attribute_sources":              wsTrustTokenProcessorMappingsAttributeSourcesValue,
				"idp_token_processor_ref":        wsTrustTokenProcessorMappingsIdpTokenProcessorRefValue,
				"issuance_criteria":              wsTrustTokenProcessorMappingsIssuanceCriteriaValue,
				"restricted_virtual_entity_ids":  wsTrustTokenProcessorMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsValues = append(wsTrustTokenProcessorMappingsValues, wsTrustTokenProcessorMappingsValue)
		}
		wsTrustTokenProcessorMappingsValue, diags := types.SetValue(wsTrustTokenProcessorMappingsElementType, wsTrustTokenProcessorMappingsValues)
		respDiags.Append(diags...)
		// Ensure that nil values are handled as false for encrypt_saml2_assertion and generate_key
		var encryptSaml2Assertion, generateKey bool
		if response.WsTrust.EncryptSaml2Assertion != nil {
			encryptSaml2Assertion = *response.WsTrust.EncryptSaml2Assertion
		}
		if response.WsTrust.GenerateKey != nil {
			generateKey = *response.WsTrust.GenerateKey
		}
		wsTrustValue, diags = types.ObjectValue(wsTrustAttrTypes, map[string]attr.Value{
			"abort_if_not_fulfilled_from_request": types.BoolPointerValue(response.WsTrust.AbortIfNotFulfilledFromRequest),
			"attribute_contract":                  wsTrustAttributeContractValue,
			"default_token_type":                  types.StringPointerValue(response.WsTrust.DefaultTokenType),
			"encrypt_saml2_assertion":             types.BoolValue(encryptSaml2Assertion),
			"generate_key":                        types.BoolValue(generateKey),
			"message_customizations":              wsTrustMessageCustomizationsValue,
			"minutes_after":                       types.Int64PointerValue(response.WsTrust.MinutesAfter),
			"minutes_before":                      types.Int64PointerValue(response.WsTrust.MinutesBefore),
			"oauth_assertion_profiles":            types.BoolPointerValue(response.WsTrust.OAuthAssertionProfiles),
			"partner_service_ids":                 wsTrustPartnerServiceIdsValue,
			"request_contract_ref":                wsTrustRequestContractRefValue,
			"token_processor_mappings":            wsTrustTokenProcessorMappingsValue,
		})
		respDiags.Append(diags...)
	}

	state.WsTrust = wsTrustValue
	return respDiags
}

func (r *idpSpConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpSpConnectionModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdpSpconnection := client.NewSpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(addOptionalIdpSpconnectionFields(ctx, createIdpSpconnection, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCreateIdpSpconnection := r.apiClient.IdpSpConnectionsAPI.CreateSpConnection(config.AuthContext(ctx, r.providerConfig))
	apiCreateIdpSpconnection = apiCreateIdpSpconnection.Body(*createIdpSpconnection)
	idpSpconnectionResponse, httpResp, err := r.apiClient.IdpSpConnectionsAPI.CreateSpConnectionExecute(apiCreateIdpSpconnection)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the IdP SP Connection", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	diags = plan.readClientResponse(idpSpconnectionResponse, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *idpSpConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state idpSpConnectionModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadIdpSpconnection, httpResp, err := r.apiClient.IdpSpConnectionsAPI.GetSpConnection(config.AuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "IdP SP Connection", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting the IdP SP Connection", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	diags = state.readClientResponse(apiReadIdpSpconnection, isImportRead)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpSpConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan idpSpConnectionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIdpSpconnection := r.apiClient.IdpSpConnectionsAPI.UpdateSpConnection(config.AuthContext(ctx, r.providerConfig), plan.ConnectionId.ValueString())
	createUpdateRequest := client.NewSpConnection(plan.EntityId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(addOptionalIdpSpconnectionFields(ctx, createUpdateRequest, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIdpSpconnection = updateIdpSpconnection.Body(*createUpdateRequest)
	updateIdpSpconnectionResponse, httpResp, err := r.apiClient.IdpSpConnectionsAPI.UpdateSpConnectionExecute(updateIdpSpconnection)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the IdP SP Connection", err, httpResp, &customId)
		return
	}

	// Read the response
	diags = plan.readClientResponse(updateIdpSpconnectionResponse, false)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *idpSpConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idpSpConnectionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.IdpSpConnectionsAPI.DeleteSpConnection(config.AuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the IdP SP Connection", err, httpResp, &customId)
	}
}

func (r *idpSpConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to connection_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("connection_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
