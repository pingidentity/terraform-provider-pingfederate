package issuancecriteria

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func ToSchema() schema.SingleNestedAttribute {
	conditionalCriteriaDefault, _ := types.SetValue(ConditionalCriteriaElemType(), nil)
	issuanceCriteriaDefault, _ := types.ObjectValue(AttrTypes(), map[string]attr.Value{
		"conditional_criteria": conditionalCriteriaDefault,
		"expression_criteria":  types.SetNull(ExpressionCriteriaElemType()),
	})
	return schema.SingleNestedAttribute{
		Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
		Computed:    true,
		Optional:    true,
		Default:     objectdefault.StaticValue(issuanceCriteriaDefault),
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"error_result": schema.StringAttribute{
							Optional:    true,
							Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"error_result": schema.StringAttribute{
							Optional:    true,
							Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}
