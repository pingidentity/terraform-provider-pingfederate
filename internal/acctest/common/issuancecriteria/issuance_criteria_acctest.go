// Copyright © 2025 Ping Identity Corporation

package issuancecriteria

import (
	"fmt"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
)

func Hcl(conditionalIssuanceCriteriaEntry *client.ConditionalIssuanceCriteriaEntry) string {
	var builder strings.Builder
	if conditionalIssuanceCriteriaEntry == nil {
		return ""
	} else {
		tf := `
		issuance_criteria = {
			conditional_criteria = [
				{
					error_result = "%s"
					source = {
						type = "%s"
					}
					attribute_name = "%s"
					condition      = "%s"
					value          = "%s"
				}
			]
		}
		`
		builder.WriteString(fmt.Sprintf(tf,
			*conditionalIssuanceCriteriaEntry.ErrorResult,
			conditionalIssuanceCriteriaEntry.Source.Type,
			conditionalIssuanceCriteriaEntry.AttributeName,
			conditionalIssuanceCriteriaEntry.Condition,
			conditionalIssuanceCriteriaEntry.Value))
	}
	return builder.String()
}

func ConditionalCriteria() *client.ConditionalIssuanceCriteriaEntry {
	conditionalIssuanceCriteriaEntry := client.NewConditionalIssuanceCriteriaEntry(
		*client.NewSourceTypeIdKey("CONTEXT"), "ClientIp", "EQUALS", "value")
	conditionalIssuanceCriteriaEntry.ErrorResult = pointers.String("error")
	return conditionalIssuanceCriteriaEntry
}
