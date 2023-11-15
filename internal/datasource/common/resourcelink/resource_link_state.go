package resourcelink

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

var (
	resourceLinkAttrTypes = map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
)

func ToDataSourceState(ctx context.Context, r *client.ResourceLink) (basetypes.ObjectValue, diag.Diagnostics) {
	if r == nil {
		return types.ObjectNull(resourceLinkAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValueFrom(ctx, resourceLinkAttrTypes, r)
}

func AttrType() map[string]attr.Type {
	return resourceLinkAttrTypes
}
