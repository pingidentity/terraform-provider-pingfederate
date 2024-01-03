package attributecontractfulfillment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttrType() map[string]attr.Type {
	attributeContractFulfillmentAttrType := map[string]attr.Type{}
	attributeContractFulfillmentAttrType["source"] = types.ObjectType{AttrTypes: sourcetypeidkey.AttrType()}
	attributeContractFulfillmentAttrType["value"] = types.StringType
	return attributeContractFulfillmentAttrType
}

func ObjType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: AttrType(),
	}
}

func MapType() types.MapType {
	return types.MapType{ElemType: types.ObjectType{AttrTypes: AttrType()}}
}

func ToState(con context.Context, attributeContractFulfillmentFromClient map[string]client.AttributeFulfillmentValue) (types.Map, diag.Diagnostics) {
	return types.MapValueFrom(con, ObjType(), attributeContractFulfillmentFromClient)
}
