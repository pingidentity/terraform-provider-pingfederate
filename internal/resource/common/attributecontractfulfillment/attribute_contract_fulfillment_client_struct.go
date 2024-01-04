package attributecontractfulfillment

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
)

func ClientStruct(attributeContractFulfillmentAttr types.Map) (map[string]client.AttributeFulfillmentValue, error) {
	attributeContractFulfillment := map[string]client.AttributeFulfillmentValue{}
	attributeContractFulfillmentErr := json.Unmarshal([]byte(internaljson.FromValue(attributeContractFulfillmentAttr, false)), &attributeContractFulfillment)
	return attributeContractFulfillment, attributeContractFulfillmentErr
}
