// Copyright © 2025 Ping Identity Corporation

package policyaction

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributemapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

// Common schema across all policy actions
func commonPolicyActionSchema() map[string]schema.Attribute {
	commonPolicyActionSchema := map[string]schema.Attribute{}
	commonPolicyActionSchema["context"] = schema.StringAttribute{
		Optional:    true,
		Description: "The result context.",
	}
	return commonPolicyActionSchema
}

func commonAttributeRulesAttr() schema.Attribute {
	return schema.SingleNestedAttribute{
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
						"attribute_source": sourcetypeidkey.ToSchemaWithDescription(true, "The source of the attribute, if this attribute is not provided then it is defaulted to be the previous authentication source."),
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
}

// Complete schemas for the individual types of policy action

func apcMappingPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_mapping"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, true),
			"attribute_sources":              attributesources.ToSchema(0, false),
			"issuance_criteria":              issuancecriteria.ToSchema(),
		},
		Required:    true,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
	attrs["authentication_policy_contract_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToSchema(),
		Required:    true,
		Description: "A reference to a resource.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    true,
		Description: "An authentication policy contract selection action.",
	}
}

func authnSelectorPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["authentication_selector_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToSchema(),
		Required:    true,
		Description: "A reference to a resource.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    true,
		Description: "An authentication selector selection action.",
	}
}

func authnSourcePolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_rules"] = commonAttributeRulesAttr()
	attrs["authentication_source"] = schema.SingleNestedAttribute{
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
		Optional:    true,
		Description: "An authentication source (IdP adapter or IdP connection).",
	}
	attrs["input_user_id_mapping"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"source": sourcetypeidkey.ToSchemaWithDescription(false, "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>"),
			"value": schema.StringAttribute{
				Required:    true,
				Description: "The value for this attribute.",
			},
		},
		Optional:    true,
		Description: "Defines how an attribute in an attribute contract should be populated.",
	}
	attrs["user_id_authenticated"] = schema.BoolAttribute{
		Optional:    true,
		Description: "Indicates whether the user ID obtained by the user ID mapping is authenticated.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    true,
		Description: "An authentication source (IdP adapter or IdP connection).",
	}
}

func continuePolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    true,
		Description: "The continue selection action.",
	}
}

func donePolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    true,
		Description: "The done selection action.",
	}
}

func fragmentPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_rules"] = commonAttributeRulesAttr()
	attrs["fragment"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToSchema(),
		Required:    true,
		Description: "A reference to a resource.",
	}
	attrs["fragment_mapping"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": schema.MapNestedAttribute{
				Description: "Defines how an attribute in an attribute contract should be populated.",
				Required:    false,
				Optional:    true,
				Computed:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": sourcetypeidkey.ToSchema(false),
						"value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
							Description: "The value for this attribute.",
						},
					},
				},
			},
			"attribute_sources": attributesources.ToSchema(0, false),
			"issuance_criteria": issuancecriteria.ToSchema(),
		},
		Required:    false,
		Optional:    true,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    true,
		Description: "A authentication policy fragment selection action.",
	}
}

func localIdentityMappingPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["inbound_mapping"] = attributemapping.ToSchema(false)
	attrs["local_identity_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToSchema(),
		Required:    true,
		Description: "A reference to a resource.",
	}
	attrs["outbound_attribute_mapping"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, true),
			"attribute_sources":              attributesources.ToSchema(0, false),
			"issuance_criteria":              issuancecriteria.ToSchema(),
		},
		Required:    true,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    true,
		Description: "A local identity profile selection action.",
	}
}

func restartPolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    true,
		Description: "The restart selection action.",
	}
}

// Schema for the polymorphic attribute allowing you to specify a single policy action type
func ToSchema() schema.SingleNestedAttribute {
	// In the future it may be worth adding validators to ensure only one of the policy action types is set, but
	// currently it causes a big performance hit
	return schema.SingleNestedAttribute{
		Description: "The result action.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"apc_mapping_policy_action":            apcMappingPolicyActionSchema(),
			"authn_selector_policy_action":         authnSelectorPolicyActionSchema(),
			"authn_source_policy_action":           authnSourcePolicyActionSchema(),
			"continue_policy_action":               continuePolicyActionSchema(),
			"done_policy_action":                   donePolicyActionSchema(),
			"fragment_policy_action":               fragmentPolicyActionSchema(),
			"local_identity_mapping_policy_action": localIdentityMappingPolicyActionSchema(),
			"restart_policy_action":                restartPolicyActionSchema(),
		},
	}
}
