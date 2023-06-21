package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
)

func ToStateAuthenticationPolicyContract(r *client.AuthenticationPolicyContract) (basetypes.SetValue, basetypes.SetValue) {
	var attrType = map[string]attr.Type{"name": types.StringType}
	clientCoreAttributes := r.GetCoreAttributes()
	var caSlice = []attr.Value{}
	cAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientCoreAttributes); i++ {
		cAname := clientCoreAttributes[i].GetName()
		cAnameVal := map[string]attr.Value{"name": types.StringValue(cAname)}
		newCaObj, _ := types.ObjectValue(attrType, cAnameVal)
		caSlice = append(caSlice, newCaObj)
	}
	caSliceOfObj, _ := types.SetValue(cAobjSlice, caSlice)

	clientExtAttributes := r.GetExtendedAttributes()
	var eaSlice = []attr.Value{}
	eAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientExtAttributes); i++ {
		eAname := clientExtAttributes[i].GetName()
		eAnameVal := map[string]attr.Value{"name": types.StringValue(eAname)}
		newEaObj, _ := types.ObjectValue(attrType, eAnameVal)
		eaSlice = append(eaSlice, newEaObj)
	}
	eaSliceOfObj, _ := types.SetValue(eAobjSlice, eaSlice)

	return caSliceOfObj, eaSliceOfObj
}

func ToStateResourceLink(r *client.ResourceLink) (basetypes.StringValue, basetypes.StringValue) {
	clientResourceLinkId := r.GetId()
	clientResourceLocation := r.GetLocation()
	linkIdStringValue := basetypes.StringValue(types.StringValue(clientResourceLinkId))
	linkLocationStringValue := basetypes.StringValue(types.StringValue(clientResourceLocation))
	return linkIdStringValue, linkLocationStringValue
}

// func ToStateApiResult(r *client.ApiResult)
