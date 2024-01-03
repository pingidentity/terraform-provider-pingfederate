package issuancecriteria

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ConditionalCriteriaClientStruct(issuanceCriteria types.Object) ([]client.ConditionalIssuanceCriteriaEntry, error) {
	conditionalCriteria := []client.ConditionalIssuanceCriteriaEntry{}
	conditionalCriteriaErr := json.Unmarshal([]byte(internaljson.FromValue(issuanceCriteria.Attributes()["conditional_criteria"].(types.List), true)), &conditionalCriteria)
	if conditionalCriteriaErr != nil {
		return nil, conditionalCriteriaErr
	}
	return conditionalCriteria, nil
}

func ExpressionCriteriaClientStruct(issuanceCriteria types.Object) ([]client.ExpressionIssuanceCriteriaEntry, error) {
	expressionCriteria := []client.ExpressionIssuanceCriteriaEntry{}
	expressionCriteriaErr := json.Unmarshal([]byte(internaljson.FromValue(issuanceCriteria.Attributes()["expression_criteria"].(types.List), true)), &expressionCriteria)
	if expressionCriteriaErr != nil {
		return nil, expressionCriteriaErr
	}
	return expressionCriteria, nil
}

func ClientStruct(issuanceCriteria types.Object) (*client.IssuanceCriteria, error) {
	// conditional criteria
	var conditionalCriteriaErr error
	newIssuanceCriteria := client.NewIssuanceCriteriaWithDefaults()
	if internaltypes.IsDefined(issuanceCriteria.Attributes()["conditional_criteria"]) {
		newIssuanceCriteria.ConditionalCriteria, conditionalCriteriaErr = ConditionalCriteriaClientStruct(issuanceCriteria)
		if conditionalCriteriaErr != nil {
			return nil, conditionalCriteriaErr
		}
	}

	// expression criteria
	var expressionCriteriaErr error
	if internaltypes.IsDefined(issuanceCriteria.Attributes()["expression_criteria"]) {
		newIssuanceCriteria.ExpressionCriteria, expressionCriteriaErr = ExpressionCriteriaClientStruct(issuanceCriteria)
		if expressionCriteriaErr != nil {
			return nil, expressionCriteriaErr
		}
	}
	return newIssuanceCriteria, nil
}
