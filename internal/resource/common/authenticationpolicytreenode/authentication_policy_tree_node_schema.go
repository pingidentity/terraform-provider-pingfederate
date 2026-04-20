// Copyright © 2026 Ping Identity Corporation

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

func ToSchema(description string) schema.SingleNestedAttribute {
	return toSchemaInternal(description, true)
}

func ToSchemaNoValueDefault(description string) schema.SingleNestedAttribute {
	return toSchemaInternal(description, false)
}

func toSchemaInternal(description string, includeValueDefault bool) schema.SingleNestedAttribute {
	var actionSchema schema.Attribute
	if includeValueDefault {
		actionSchema = policyaction.ToSchema()
	} else {
		actionSchema = policyaction.ToSchemaNoValueDefault()
	}
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"action": actionSchema,
			"children": schema.ListNestedAttribute{
				Optional:     true,
				Computed:     true,
				Default:      listdefault.StaticValue(childrenDefault(1)),
				Description:  childrenDescription,
				NestedObject: buildSchema(1, includeValueDefault),
			},
		},
		Required:    true,
		Description: description,
	}
}

func buildSchema(depth int, includeValueDefault bool) schema.NestedAttributeObject {
	attrs := map[string]schema.Attribute{}
	if includeValueDefault {
		attrs["action"] = policyaction.ToSchema()
	} else {
		attrs["action"] = policyaction.ToSchemaNoValueDefault()
	}
	if depth < MaxPolicyNodeRecursiveDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Optional:     true,
			Computed:     true,
			Default:      listdefault.StaticValue(childrenDefault(depth + 1)),
			Description:  childrenDescription,
			NestedObject: buildSchema(depth+1, includeValueDefault),
		}
	}
	return schema.NestedAttributeObject{
		Attributes: attrs,
	}
}
