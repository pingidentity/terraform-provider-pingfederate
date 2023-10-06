package resourcelink

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

func ToRequest(planObj basetypes.ObjectValue) *client.ResourceLink {
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
