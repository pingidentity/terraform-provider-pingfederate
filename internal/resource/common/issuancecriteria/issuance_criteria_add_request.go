package issuancecriteria

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ConditionalCriteriaToRequest(issuanceCriteria basetypes.ObjectValue) ([]client.ConditionalIssuanceCriteriaEntry, error) {
	conditionalCriteria := []client.ConditionalIssuanceCriteriaEntry{}
	conditionalCriteriaErr := json.Unmarshal([]byte(internaljson.FromValue(issuanceCriteria.Attributes()["conditional_criteria"].(types.List), true)), &conditionalCriteria)
	if conditionalCriteriaErr != nil {
		return nil, conditionalCriteriaErr
	}
	return conditionalCriteria, nil
}

func ExpressionCriteriaToRequest(issuanceCriteria basetypes.ObjectValue) ([]client.ExpressionIssuanceCriteriaEntry, error) {
	expressionCriteria := []client.ExpressionIssuanceCriteriaEntry{}
	expressionCriteriaErr := json.Unmarshal([]byte(internaljson.FromValue(issuanceCriteria.Attributes()["expression_criteria"].(types.List), true)), &expressionCriteria)
	if expressionCriteriaErr != nil {
		return nil, expressionCriteriaErr
	}
	return expressionCriteria, nil
}

func ToRequest(issuanceCriteria basetypes.ObjectValue) (*client.IssuanceCriteria, error) {
	// conditional criteria
	var conditionalCriteriaErr error
	newIssuanceCriteria := client.NewIssuanceCriteriaWithDefaults()
	if internaltypes.IsDefined(issuanceCriteria.Attributes()["conditional_criteria"]) {
		newIssuanceCriteria.ConditionalCriteria, conditionalCriteriaErr = ConditionalCriteriaToRequest(issuanceCriteria)
		if conditionalCriteriaErr != nil {
			return nil, conditionalCriteriaErr
		}
	}

	// expression criteria
	var expressionCriteriaErr error
	if internaltypes.IsDefined(issuanceCriteria.Attributes()["expression_criteria"]) {
		newIssuanceCriteria.ExpressionCriteria, expressionCriteriaErr = ExpressionCriteriaToRequest(issuanceCriteria)
		if expressionCriteriaErr != nil {
			return nil, expressionCriteriaErr
		}
	}
	return newIssuanceCriteria, nil
}
