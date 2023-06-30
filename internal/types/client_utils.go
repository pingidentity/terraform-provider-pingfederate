package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
)

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
