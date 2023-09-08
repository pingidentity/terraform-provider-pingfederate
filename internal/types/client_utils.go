package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
)

var (
	resourceLinkAttrTypes = map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
)

func ToRequestResourceLink(planObj basetypes.ObjectValue) *client.ResourceLink {
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

func ToStateResourceLink(ctx context.Context, r *client.ResourceLink, diags *diag.Diagnostics) basetypes.ObjectValue {
	if r == nil {
		return types.ObjectNull(resourceLinkAttrTypes)
	}
	linkObjectValue, objectValueFromDiags := types.ObjectValueFrom(ctx, resourceLinkAttrTypes, r)
	diags.Append(objectValueFromDiags...)
	return linkObjectValue
}
