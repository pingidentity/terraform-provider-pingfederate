package attributecontractfulfillment

import (
	"fmt"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

func InitialAttributeContractFulfillment() client.AttributeFulfillmentValue {
	initialAttributecontractfulfillment := *client.NewAttributeFulfillmentValue(
		*client.NewSourceTypeIdKey("TEXT"),
		"value",
	)
	return initialAttributecontractfulfillment
}

func UpdatedAttributeContractFulfillment() client.AttributeFulfillmentValue {
	updatedAttributecontractfulfillment := *client.NewAttributeFulfillmentValue(
		*client.NewSourceTypeIdKey("CONTEXT"),
		"ClientIp",
	)
	return updatedAttributecontractfulfillment
}

func Hcl(aCf *client.AttributeFulfillmentValue) string {
	var builder strings.Builder
	if aCf == nil {
		return ""
	}
	if aCf != nil {
		tf := `
			source = {
				type = "%s"
			},
			value = "%s"
		`
		builder.WriteString(fmt.Sprintf(tf,
			aCf.Source.Type,
			aCf.Value))
	}
	return builder.String()
}
