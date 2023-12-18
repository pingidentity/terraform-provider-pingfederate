package resourcelink

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(planObj basetypes.ObjectValue) (*client.ResourceLink, error) {
	if !internaltypes.IsDefined(planObj) {
		return nil, errors.New("null or Unknown object value passed in when creating resource link client struct")
	}

	objValues := planObj.Attributes()
	objId, ok := objValues["id"]
	if !ok {
		return nil, errors.New("object value missing \"id\" attribute when creating resource link client struct")
	}
	objLoc, ok := objValues["location"]
	if !ok {
		return nil, errors.New("object value missing \"location\" attribute when creating resource link client struct")
	}
	idStrValue := objId.(basetypes.StringValue)
	locStrValue := objLoc.(basetypes.StringValue)
	newLink := client.NewResourceLinkWithDefaults()
	newLink.SetId(idStrValue.ValueString())
	newLink.SetLocation(locStrValue.ValueString())

	return newLink, nil
}
