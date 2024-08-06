package issuancecriteria

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/sourcetypeidkey"
)

func ToDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
		Optional:    false,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"conditional_criteria": schema.SetNestedAttribute{
				Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": sourcetypeidkey.ToDataSourceSchema(),
						"attribute_name": schema.StringAttribute{
							Description: "The name of the attribute to use in this issuance criterion.",
							Optional:    false,
							Computed:    true,
						},
						"condition": schema.StringAttribute{
							Description: "The condition that will be applied to the source attribute's value and the expected value.",
							Optional:    false,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The expected value of this issuance criterion.",
						},
						"error_result": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
						},
					},
				},
			},
			"expression_criteria": schema.SetNestedAttribute{
				Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria.",
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"expression": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The OGNL expression to evaluate.",
						},
						"error_result": schema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
						},
					},
				},
			},
		},
	}
}
