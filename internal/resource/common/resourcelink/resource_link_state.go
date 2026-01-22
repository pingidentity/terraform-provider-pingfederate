// Copyright © 2025 Ping Identity Corporation

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
	return types.ObjectValueFrom(ctx, resourceLinkAttrTypes, r)
}

func AttrType() map[string]attr.Type {
	return resourceLinkAttrTypes
}
