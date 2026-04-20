// Copyright © 2026 Ping Identity Corporation

package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/authenticationpolicytreenode"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

var childrenDescription = "The nodes inside the authentication policy tree node of type AuthenticationPolicyTreeNode."

func rootNodeObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: authenticationpolicytreenode.GetRootNodeAttrTypes()}
}

func nodeObjectType(depth int) types.ObjectType {
	attrs := map[string]attr.Type{
		"action": types.ObjectType{AttrTypes: policyaction.AttrTypes()},
	}
	if depth < authenticationpolicytreenode.MaxPolicyNodeRecursiveDepth {
		attrs["children"] = types.ListType{ElemType: nodeObjectType(depth + 1)}
	}
	return types.ObjectType{AttrTypes: attrs}
}

func childrenListType(depth int) types.ListType {
	return types.ListType{ElemType: nodeObjectType(depth)}
}

func DataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		CustomType: rootNodeObjectType(),
		Attributes: map[string]schema.Attribute{
			"action": policyaction.ToSchema(),
			"children": schema.ListNestedAttribute{
				Optional:     false,
				Computed:     true,
				CustomType:   childrenListType(1),
				Description:  childrenDescription,
				NestedObject: buildSchema(1),
			},
		},
		Optional:    false,
		Computed:    true,
		Description: "The beginning action for the authentication fragment policy.",
	}
}

func buildSchema(depth int) schema.NestedAttributeObject {
	attrs := map[string]schema.Attribute{
		"action": policyaction.ToSchema(),
	}
	if depth < authenticationpolicytreenode.MaxPolicyNodeRecursiveDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Optional:     false,
			Computed:     true,
			CustomType:   childrenListType(depth + 1),
			Description:  childrenDescription,
			NestedObject: buildSchema(depth + 1),
		}
	}
	return schema.NestedAttributeObject{
		Attributes: attrs,
		CustomType: nodeObjectType(depth),
	}
}
