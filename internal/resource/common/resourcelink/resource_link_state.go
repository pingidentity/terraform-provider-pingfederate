package resourcelink

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
)

func ResourceLinkStateAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
}

func ToStateResourceLink(con context.Context, resLinkVals client.ResourceLink) basetypes.ObjectValue {
	resourceLink, _ := types.ObjectValueFrom(con, ResourceLinkStateAttrType(), resLinkVals)
	return resourceLink
}
