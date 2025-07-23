// Copyright © 2025 Ping Identity Corporation

package issuancecriteria

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

func ClientStruct(issuanceCriteria types.Object) *client.IssuanceCriteria {
	result := &client.IssuanceCriteria{}
	issuanceCriteriaAttrs := issuanceCriteria.Attributes()
	if !issuanceCriteriaAttrs["conditional_criteria"].IsNull() && !issuanceCriteriaAttrs["conditional_criteria"].IsUnknown() {
		result.ConditionalCriteria = []client.ConditionalIssuanceCriteriaEntry{}
		for _, conditionalCriteriaElement := range issuanceCriteriaAttrs["conditional_criteria"].(types.Set).Elements() {
			conditionalCriteriaValue := client.ConditionalIssuanceCriteriaEntry{}
			conditionalCriteriaAttrs := conditionalCriteriaElement.(types.Object).Attributes()
			conditionalCriteriaValue.AttributeName = conditionalCriteriaAttrs["attribute_name"].(types.String).ValueString()
			conditionalCriteriaValue.Condition = conditionalCriteriaAttrs["condition"].(types.String).ValueString()
			conditionalCriteriaValue.ErrorResult = conditionalCriteriaAttrs["error_result"].(types.String).ValueStringPointer()
			conditionalCriteriaSourceValue := client.SourceTypeIdKey{}
			conditionalCriteriaSourceAttrs := conditionalCriteriaAttrs["source"].(types.Object).Attributes()
			conditionalCriteriaSourceValue.Id = conditionalCriteriaSourceAttrs["id"].(types.String).ValueStringPointer()
			conditionalCriteriaSourceValue.Type = conditionalCriteriaSourceAttrs["type"].(types.String).ValueString()
			conditionalCriteriaValue.Source = conditionalCriteriaSourceValue
			conditionalCriteriaValue.Value = conditionalCriteriaAttrs["value"].(types.String).ValueString()
			result.ConditionalCriteria = append(result.ConditionalCriteria, conditionalCriteriaValue)
		}
	}
	if !issuanceCriteriaAttrs["expression_criteria"].IsNull() && !issuanceCriteriaAttrs["expression_criteria"].IsUnknown() {
		result.ExpressionCriteria = []client.ExpressionIssuanceCriteriaEntry{}
		for _, expressionCriteriaElement := range issuanceCriteriaAttrs["expression_criteria"].(types.Set).Elements() {
			expressionCriteriaValue := client.ExpressionIssuanceCriteriaEntry{}
			expressionCriteriaAttrs := expressionCriteriaElement.(types.Object).Attributes()
			expressionCriteriaValue.ErrorResult = expressionCriteriaAttrs["error_result"].(types.String).ValueStringPointer()
			expressionCriteriaValue.Expression = expressionCriteriaAttrs["expression"].(types.String).ValueString()
			result.ExpressionCriteria = append(result.ExpressionCriteria, expressionCriteriaValue)
		}
	}
	return result
}
