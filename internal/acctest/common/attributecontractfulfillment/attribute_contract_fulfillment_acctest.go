package attributecontractfulfillment

import (
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

func Hcl(aCf *client.AttributeFulfillmentValue) string {
	var builder strings.Builder
	if aCf == nil {
		return ""
	}
	if aCf != nil {
		builder.WriteString("      source = {\n")
		builder.WriteString("        type = \"")
		builder.WriteString(aCf.Source.Type)
		builder.WriteString("\"\n")
		builder.WriteString("      },\n")
		builder.WriteString("      value = \"")
		builder.WriteString(aCf.Value)
		builder.WriteString("\"")
	}
	return builder.String()
}
