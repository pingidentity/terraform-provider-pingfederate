package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
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

func ToStateResourceLink(con context.Context, resLinkVals client.ResourceLink) basetypes.ObjectValue {
	attrTypes := map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
	resourceLink, _ := types.ObjectValueFrom(con, attrTypes, resLinkVals)
	return resourceLink
}

func ResourceLinkStateAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
}
