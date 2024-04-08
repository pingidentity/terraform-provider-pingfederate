package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/authenticationpolicytreenode"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

var childrenDescription = "The nodes inside the authentication policy tree node of type AuthenticationPolicyTreeNode."

func DataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"action": policyaction.ToSchema(),
			"children": schema.ListNestedAttribute{
				Optional:     false,
				Computed:     true,
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
			Description:  childrenDescription,
			NestedObject: buildSchema(depth + 1),
		}
	}
	return schema.NestedAttributeObject{
		Attributes: attrs,
	}
}
