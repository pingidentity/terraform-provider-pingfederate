package resourcelink

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
)

var (
	resourceLinkAttrTypes = map[string]attr.Type{
		"id":       types.StringType,
		"location": types.StringType,
	}

	resourceLinkNoLocationAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
)

func ToState(ctx context.Context, r *client.ResourceLink) (types.Object, diag.Diagnostics) {
	if r == nil {
		return types.ObjectNull(resourceLinkAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValueFrom(ctx, resourceLinkAttrTypes, r)
}

func ToStateNoLocation(r *client.ResourceLink) (types.Object, diag.Diagnostics) {
	if r == nil || r.Id == "" {
		return types.ObjectNull(resourceLinkNoLocationAttrTypes), diag.Diagnostics{}
	}

	objectValue := map[string]attr.Value{
		"id": types.StringValue(r.Id),
	}
	return types.ObjectValue(resourceLinkNoLocationAttrTypes, objectValue)
}

func AttrType() map[string]attr.Type {
	return resourceLinkAttrTypes
}

func AttrTypeNoLocation() map[string]attr.Type {
	return resourceLinkNoLocationAttrTypes
}
