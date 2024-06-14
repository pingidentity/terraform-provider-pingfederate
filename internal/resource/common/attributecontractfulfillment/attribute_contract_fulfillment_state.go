package attributecontractfulfillment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttrTypes() map[string]attr.Type {
	attributeContractFulfillmentAttrType := map[string]attr.Type{}
	attributeContractFulfillmentAttrType["source"] = types.ObjectType{AttrTypes: sourcetypeidkey.AttrTypes()}
	attributeContractFulfillmentAttrType["value"] = types.StringType
	return attributeContractFulfillmentAttrType
}

func ObjType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: AttrTypes(),
	}
}

func MapType() types.MapType {
	return types.MapType{ElemType: types.ObjectType{AttrTypes: AttrTypes()}}
}

func ToState(con context.Context, attributeContractFulfillmentFromClient *map[string]client.AttributeFulfillmentValue) (types.Map, diag.Diagnostics) {
	if attributeContractFulfillmentFromClient == nil {
		return types.MapNull(ObjType()), nil
	}
	return types.MapValueFrom(con, ObjType(), attributeContractFulfillmentFromClient)
}
