package policyaction

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/sourcetypeidkey"
)

// Common schema across all policy actions
func commonPolicyActionSchema() map[string]schema.Attribute {
	commonPolicyActionSchema := map[string]schema.Attribute{}
	commonPolicyActionSchema["context"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The result context.",
	}
	return commonPolicyActionSchema
}

// TODO probably make common across other resources
func commonAttributeMappingAttr() schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
			"attribute_sources":              attributesources.ToDataSourceSchema(),
			"issuance_criteria":              issuancecriteria.ToDataSourceSchema(),
		},
		Optional:    false,
		Computed:    true,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
}

func commonAttributeRulesAttr() schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"fallback_to_success": schema.BoolAttribute{
				Optional:    false,
				Computed:    true,
				Description: "When all the rules fail, you may choose to default to the general success action or fail. Default to success.",
			},
			"items": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attribute_name": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The name of the attribute to use in this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
						},
						"attribute_source": sourcetypeidkey.ToDataSourceSchema(),
						"condition": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The condition that will be applied to the attribute's expected value. This field is required if the Attribute Source type is not 'EXPRESSION'.",
						},
						"expected_value": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The expected value of this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
						},
						"expression": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The expression of this attribute rule. This field is required if the Attribute Source type is 'EXPRESSION'.",
						},
						"result": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The result of this attribute rule.",
						},
					},
				},
				Optional:    false,
				Computed:    true,
				Description: "The actual list of attribute rules.",
			},
		},
		Optional:    false,
		Computed:    true,
		Description: "A collection of attribute rules",
	}
}

// Complete schemas for the individual types of policy action

func apcMappingPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_mapping"] = commonAttributeMappingAttr()
	attrs["authentication_policy_contract_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToDataSourceSchema(),
		Optional:    false,
		Computed:    true,
		Description: "A reference to a resource.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    false,
		Computed:    true,
		Description: "An authentication policy contract selection action.",
	}
}

func authnSelectorPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["authentication_selector_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToDataSourceSchema(),
		Optional:    false,
		Computed:    true,
		Description: "A reference to a resource.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    false,
		Computed:    true,
		Description: "An authentication selector selection action.",
	}
}

func authnSourcePolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_rules"] = commonAttributeRulesAttr()
	attrs["authentication_source"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"source_ref": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToDataSourceSchema(),
				Optional:    false,
				Computed:    true,
				Description: "A reference to a resource.",
			},
			"type": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The type of this authentication source.",
			},
		},
		Optional:    false,
		Computed:    true,
		Description: "An authentication source (IdP adapter or IdP connection).",
	}
	attrs["input_user_id_mapping"] = schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"source": sourcetypeidkey.ToDataSourceSchema(),
			"value": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The value for this attribute.",
			},
		},
		Optional:    false,
		Computed:    true,
		Description: "Defines how an attribute in an attribute contract should be populated.",
	}
	attrs["user_id_authenticated"] = schema.BoolAttribute{
		Optional:    false,
		Computed:    true,
		Description: "Indicates whether the user ID obtained by the user ID mapping is authenticated.",
	}
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    false,
		Computed:    true,
		Description: "An authentication source (IdP adapter or IdP connection).",
	}
}

func continuePolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    false,
		Computed:    true,
		Description: "The continue selection action.",
	}
}

func donePolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    false,
		Computed:    true,
		Description: "The done selection action.",
	}
}

func fragmentPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["attribute_rules"] = commonAttributeRulesAttr()
	attrs["fragment"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToDataSourceSchema(),
		Optional:    false,
		Computed:    true,
		Description: "A reference to a resource.",
	}
	attrs["fragment_mapping"] = commonAttributeMappingAttr()
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    false,
		Computed:    true,
		Description: "A authentication policy fragment selection action.",
	}
}

func localIdentityMappingPolicyActionSchema() schema.SingleNestedAttribute {
	attrs := commonPolicyActionSchema()
	attrs["inbound_mapping"] = commonAttributeMappingAttr()
	attrs["local_identity_ref"] = schema.SingleNestedAttribute{
		Attributes:  resourcelink.ToDataSourceSchema(),
		Optional:    false,
		Computed:    true,
		Description: "A reference to a resource.",
	}
	attrs["outbound_attribute_mapping"] = commonAttributeMappingAttr()
	return schema.SingleNestedAttribute{
		Attributes:  attrs,
		Optional:    false,
		Computed:    true,
		Description: "A local identity profile selection action.",
	}
}

func restartPolicyActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes:  commonPolicyActionSchema(),
		Optional:    false,
		Computed:    true,
		Description: "The restart selection action.",
	}
}

// Schema for the polymorphic attribute allowing you to specify a single policy action type
func DataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The result action.",
		Optional:    false,
		Computed:    true,
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
