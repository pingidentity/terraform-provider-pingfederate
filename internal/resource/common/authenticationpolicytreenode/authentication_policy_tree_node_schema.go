package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

const maxRecursiveDepth = 10

var childrenDescription = "The nodes inside the authentication policy tree node of type AuthenticationPolicyTreeNode."

func Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"policy_action": policyaction.Schema(),
			"children": schema.ListNestedAttribute{
				Optional:     true,
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
		"policy_action": policyaction.Schema(),
	}
	if depth < maxRecursiveDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Optional:     true,
			Description:  childrenDescription,
			NestedObject: buildSchema(depth + 1),
		}
	}
	return schema.NestedAttributeObject{
		Attributes: attrs,
	}
}
