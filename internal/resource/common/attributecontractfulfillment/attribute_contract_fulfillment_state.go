package attributecontractfulfillment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttributeContractFulfillmentAttrType() basetypes.ObjectType {
	return basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source": basetypes.ObjectType{
				AttrTypes: sourcetypeidkey.SourceTypeIdKeyAttrType(),
			},
			"value": basetypes.StringType{},
		},
	}
}

func AttributeContractFulfillmentToState(con context.Context, attributeContractFulfillmentFromClient map[string]client.AttributeFulfillmentValue) basetypes.MapValue {
	attributeContractFulfillmentToState, _ := types.MapValueFrom(con, AttributeContractFulfillmentAttrType(), attributeContractFulfillmentFromClient)
	return attributeContractFulfillmentToState
}
