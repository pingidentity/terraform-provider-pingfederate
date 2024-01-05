package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

func Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"policy_action": policyaction.Schema(),
			//TODO children, recursiveness
		},
		Optional:    true,
		Description: "An authentication policy tree node.",
	}
}
