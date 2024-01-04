package issuancecriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func ConditionalCriteriaElemType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source": types.ObjectType{
				AttrTypes: sourcetypeidkey.AttrType(),
			},
			"attribute_name": types.StringType,
			"condition":      types.StringType,
			"value":          types.StringType,
			"error_result":   types.StringType,
		},
	}
}

func ExpressionCriteriaElemType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"expression":   types.StringType,
			"error_result": types.StringType,
		},
	}
}

func AttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"conditional_criteria": types.ListType{
			ElemType: ConditionalCriteriaElemType(),
		},
		"expression_criteria": types.ListType{
			ElemType: ExpressionCriteriaElemType(),
		},
	}
}

func ToState(con context.Context, issuanceCriteriaFromClient *client.IssuanceCriteria) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(con, AttrType(), issuanceCriteriaFromClient)
}
