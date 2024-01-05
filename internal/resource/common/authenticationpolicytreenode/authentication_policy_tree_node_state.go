package authenticationpolicytreenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

var (
	rootNodeAttrTypes = map[string]attr.Type{
		"policy_action": types.ObjectType{AttrTypes: policyaction.AttrTypes()},
		//TODO children
	}
)

func State(ctx context.Context, node *client.AuthenticationPolicyTreeNode) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if node == nil {
		diags.AddError("provided authentication policy tree node is nil", "")
		return types.ObjectNull(rootNodeAttrTypes), diags
	}
	var attrValues = map[string]attr.Value{}

	attrValues["policy_action"], diags = policyaction.State(ctx, &node.Action)
	if diags.HasError() {
		return types.ObjectNull(rootNodeAttrTypes), diags
	}

	//TODO children

	return types.ObjectValue(rootNodeAttrTypes, attrValues)
}
