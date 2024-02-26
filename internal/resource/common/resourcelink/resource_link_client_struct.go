package resourcelink

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(planObj basetypes.ObjectValue, includeLocation bool) (*client.ResourceLink, error) {
	if !internaltypes.IsDefined(planObj) {
		return nil, nil
	}

	objValues := planObj.Attributes()
	objId, ok := objValues["id"]
	if !ok {
		return nil, errors.New("object value missing \"id\" attribute when creating resource link client struct")
	}
	newLink := client.NewResourceLinkWithDefaults()
	idStrValue := objId.(basetypes.StringValue)
	newLink.SetId(idStrValue.ValueString())

	if includeLocation {
		objLoc, ok := objValues["location"]
		if !ok {
			return nil, errors.New("object value missing \"location\" attribute when creating resource link client struct")
		}
		locStrValue := objLoc.(basetypes.StringValue)
		newLink.SetLocation(locStrValue.ValueString())
	}

	return newLink, nil
}
