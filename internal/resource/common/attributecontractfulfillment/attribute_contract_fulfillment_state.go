package attributecontractfulfillment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttrType() map[string]attr.Type {
	attributeContractFulfillmentAttrType := map[string]attr.Type{}
	attributeContractFulfillmentAttrType["source"] = basetypes.ObjectType{AttrTypes: sourcetypeidkey.AttrType()}
	attributeContractFulfillmentAttrType["value"] = basetypes.StringType{}
	return attributeContractFulfillmentAttrType
}

func ObjType() basetypes.ObjectType {
	return basetypes.ObjectType{
		AttrTypes: AttrType(),
	}
}

func MapType() basetypes.MapType {
	return basetypes.MapType{ElemType: types.ObjectType{AttrTypes: AttrType()}}
}

func ToState(con context.Context, attributeContractFulfillmentFromClient map[string]client.AttributeFulfillmentValue) basetypes.MapValue {
	attributeContractFulfillmentToState, _ := types.MapValueFrom(con, ObjType(), attributeContractFulfillmentFromClient)
	return attributeContractFulfillmentToState
}
