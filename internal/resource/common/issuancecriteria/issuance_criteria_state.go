package issuancecriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func ConditionalCriteriaElemType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source": types.ObjectType{
				AttrTypes: sourcetypeidkey.AttrTypes(),
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

func AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"conditional_criteria": types.SetType{
			ElemType: ConditionalCriteriaElemType(),
		},
		"expression_criteria": types.SetType{
			ElemType: ExpressionCriteriaElemType(),
		},
	}
}

func ToState(con context.Context, issuanceCriteriaFromClient *client.IssuanceCriteria) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(con, AttrTypes(), issuanceCriteriaFromClient)
}
