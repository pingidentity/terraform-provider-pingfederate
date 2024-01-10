package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

const MaxPolicyNodeRecursiveDepth = 10

var (
	childrenDescription = "The nodes inside the authentication policy tree node of type AuthenticationPolicyTreeNode."
)

func childrenDefault(depth int) types.List {
	baseAttrTypes := map[string]attr.Type{
		"action": types.ObjectType{AttrTypes: policyaction.AttrTypes()},
	}
	if depth < MaxPolicyNodeRecursiveDepth {
		baseAttrTypes["children"] = childrenAttrTypes(depth + 1)
	}

	resp, _ := types.ListValue(types.ObjectType{AttrTypes: baseAttrTypes}, []attr.Value{})
	return resp
}

func Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"action": policyaction.Schema(),
			"children": schema.ListNestedAttribute{
				Optional:     true,
				Computed:     true,
				Default:      listdefault.StaticValue(childrenDefault(1)),
				Description:  childrenDescription,
				NestedObject: buildSchema(1),
			},
		},
		Required:    true,
		Description: "The beginning action for the authentication fragment policy.",
	}
}

func buildSchema(depth int) schema.NestedAttributeObject {
	attrs := map[string]schema.Attribute{
		"action": policyaction.Schema(),
	}
	if depth < MaxPolicyNodeRecursiveDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Optional:     true,
			Computed:     true,
			Default:      listdefault.StaticValue(childrenDefault(depth + 1)),
			Description:  childrenDescription,
			NestedObject: buildSchema(depth + 1),
		}
	}
	return schema.NestedAttributeObject{
		Attributes: attrs,
	}
}
