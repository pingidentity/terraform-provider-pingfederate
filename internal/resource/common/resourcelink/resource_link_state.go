// Copyright © 2026 Ping Identity Corporation

package resourcelink

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
)

var (
	resourceLinkAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
)

func ToState(ctx context.Context, r *client.ResourceLink) (types.Object, diag.Diagnostics) {
	if r == nil {
		return types.ObjectNull(resourceLinkAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValue(resourceLinkAttrTypes, map[string]attr.Value{
		"id": types.StringValue(r.Id),
	})
}

// ToStateMust explicitly constructs the resource link object, panicking on unexpected diags.
// For use where the client struct is known to be valid.
func ToStateMust(r *client.ResourceLink) types.Object {
	if r == nil {
		return types.ObjectNull(resourceLinkAttrTypes)
	}
	return types.ObjectValueMust(resourceLinkAttrTypes, map[string]attr.Value{
		"id": types.StringValue(r.Id),
	})
}

func AttrType() map[string]attr.Type {
	return resourceLinkAttrTypes
}
