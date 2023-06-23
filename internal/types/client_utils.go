package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func ToRequestResourceLink(con context.Context, planObj basetypes.ObjectValue) *client.ResourceLink {
	objValues := planObj.Attributes()
	objId := objValues["id"]
	objLoc := objValues["location"]
	idStrValue := objId.(basetypes.StringValue)
	locStrValue := objLoc.(basetypes.StringValue)
	newLink := client.NewResourceLinkWithDefaults()
	newLink.SetId(idStrValue.ValueString())
	newLink.SetLocation(locStrValue.ValueString())

	return newLink
}

func ToStateResourceLink(r *client.ResourceLink, diags diag.Diagnostics) basetypes.ObjectValue {
	attrTypes := map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}

	getId := r.GetId()
	getLocation := r.GetLocation()
	attrValues := map[string]attr.Value{
		"id":       StringTypeOrNil(&getId, false),
		"location": StringTypeOrNil(&getLocation, false),
	}

	linkObjectValue := MaptoObjValue(attrTypes, attrValues, diags)
	return linkObjectValue
}
