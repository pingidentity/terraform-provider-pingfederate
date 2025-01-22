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
	var respDiags, diags diag.Diagnostics
	var issuanceCriteriaValue types.Object
	if issuanceCriteriaFromClient == nil {
		issuanceCriteriaValue = types.ObjectNull(AttrTypes())
	} else {
		var conditionalCriteriaFinalValue types.Set
		if issuanceCriteriaFromClient.ConditionalCriteria == nil {
			conditionalCriteriaFinalValue = types.SetNull(ConditionalCriteriaElemType())
		} else {
			var conditionalCriteriaValues []attr.Value
			for _, conditionalCriteriaResponseValue := range issuanceCriteriaFromClient.ConditionalCriteria {
				conditionalCriteriaSourceValue, diags := types.ObjectValue(sourcetypeidkey.AttrTypes(), map[string]attr.Value{
					"id":   types.StringPointerValue(conditionalCriteriaResponseValue.Source.Id),
					"type": types.StringValue(conditionalCriteriaResponseValue.Source.Type),
				})
				respDiags.Append(diags...)
				// PF can return error_result as an empty string rather than null
				errorResult := conditionalCriteriaResponseValue.ErrorResult
				if errorResult != nil && *errorResult == "" {
					errorResult = nil
				}
				conditionalCriteriaValue, diags := types.ObjectValue(ConditionalCriteriaElemType().AttrTypes, map[string]attr.Value{
					"attribute_name": types.StringValue(conditionalCriteriaResponseValue.AttributeName),
					"condition":      types.StringValue(conditionalCriteriaResponseValue.Condition),
					"error_result":   types.StringPointerValue(errorResult),
					"source":         conditionalCriteriaSourceValue,
					"value":          types.StringValue(conditionalCriteriaResponseValue.Value),
				})
				respDiags.Append(diags...)
				conditionalCriteriaValues = append(conditionalCriteriaValues, conditionalCriteriaValue)
			}
			conditionalCriteriaFinalValue, diags = types.SetValue(ConditionalCriteriaElemType(), conditionalCriteriaValues)
			respDiags.Append(diags...)
		}
		var expressionCriteriaFinalValue types.Set
		if issuanceCriteriaFromClient.ExpressionCriteria == nil {
			expressionCriteriaFinalValue = types.SetNull(ExpressionCriteriaElemType())
		} else {
			var expressionCriteriaValues []attr.Value
			for _, expressionCriteriaResponseValue := range issuanceCriteriaFromClient.ExpressionCriteria {
				// PF can return error_result as an empty string rather than null
				errorResult := expressionCriteriaResponseValue.ErrorResult
				if errorResult != nil && *errorResult == "" {
					errorResult = nil
				}
				expressionCriteriaValue, diags := types.ObjectValue(ExpressionCriteriaElemType().AttrTypes, map[string]attr.Value{
					"error_result": types.StringPointerValue(errorResult),
					"expression":   types.StringValue(expressionCriteriaResponseValue.Expression),
				})
				respDiags.Append(diags...)
				expressionCriteriaValues = append(expressionCriteriaValues, expressionCriteriaValue)
			}
			expressionCriteriaFinalValue, diags = types.SetValue(ExpressionCriteriaElemType(), expressionCriteriaValues)
			respDiags.Append(diags...)
		}
		issuanceCriteriaValue, diags = types.ObjectValue(AttrTypes(), map[string]attr.Value{
			"conditional_criteria": conditionalCriteriaFinalValue,
			"expression_criteria":  expressionCriteriaFinalValue,
		})
		respDiags.Append(diags...)
	}
	return issuanceCriteriaValue, respDiags
}
