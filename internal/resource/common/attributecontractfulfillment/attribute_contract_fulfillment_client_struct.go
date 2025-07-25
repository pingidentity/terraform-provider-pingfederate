// Copyright © 2025 Ping Identity Corporation

package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

func ClientStruct(attributeContractFulfillmentAttr types.Map) map[string]client.AttributeFulfillmentValue {
	attributeContractFulfillment := map[string]client.AttributeFulfillmentValue{}
	for key, fulfillment := range attributeContractFulfillmentAttr.Elements() {
		fulfillmentValue := client.AttributeFulfillmentValue{}
		fulfillmentAttrs := fulfillment.(types.Object).Attributes()
		fulfillmentValue.Value = fulfillmentAttrs["value"].(types.String).ValueString()
		fulfillmentValue.Source = client.SourceTypeIdKey{}
		sourceAttrs := fulfillmentAttrs["source"].(types.Object).Attributes()
		fulfillmentValue.Source.Type = sourceAttrs["type"].(types.String).ValueString()
		fulfillmentValue.Source.Id = sourceAttrs["id"].(types.String).ValueStringPointer()
		attributeContractFulfillment[key] = fulfillmentValue
	}
	return attributeContractFulfillment
}
