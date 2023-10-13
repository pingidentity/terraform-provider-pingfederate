package issuancecriteria

import (
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
)

func Hcl(conditionalIssuanceCriteriaEntry *client.ConditionalIssuanceCriteriaEntry) string {
	var builder strings.Builder
	if conditionalIssuanceCriteriaEntry == nil {
		return ""
	}
	if conditionalIssuanceCriteriaEntry != nil {
		builder.WriteString("  issuance_criteria = {\n    conditional_criteria = [\n      {\n")
		builder.WriteString("        error_result = \"")
		builder.WriteString(*conditionalIssuanceCriteriaEntry.ErrorResult)
		builder.WriteString("\"\n        source = {\n          type = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Source.Type)
		builder.WriteString("\"\n        }\n        attribute_name = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.AttributeName)
		builder.WriteString("\"\n        condition      = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Condition)
		builder.WriteString("\"\n        value          = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Value)
		builder.WriteString("\"\n      }\n    ]\n  }\n")
	}
	return builder.String()
}

func ConditionalCriteria() *client.ConditionalIssuanceCriteriaEntry {
	conditionalIssuanceCriteriaEntry := client.NewConditionalIssuanceCriteriaEntry(
		*client.NewSourceTypeIdKey("CONTEXT"), "ClientIp", "EQUALS", "value")
	conditionalIssuanceCriteriaEntry.ErrorResult = pointers.String("error")
	return conditionalIssuanceCriteriaEntry
}
