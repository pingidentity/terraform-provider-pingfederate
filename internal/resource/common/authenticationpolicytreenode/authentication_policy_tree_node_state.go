package authenticationpolicytreenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

func GetRootNodeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action":   types.ObjectType{AttrTypes: policyaction.AttrTypes()},
		"children": childrenAttrTypes(1),
	}
}

func childrenAttrTypes(depth int) types.ListType {
	attrs := map[string]attr.Type{
		"action": types.ObjectType{AttrTypes: policyaction.AttrTypes()},
	}
	if depth < MaxPolicyNodeRecursiveDepth {
		attrs["children"] = childrenAttrTypes(depth + 1)
	}
	return types.ListType{
		ElemType: types.ObjectType{AttrTypes: attrs},
	}
}

func ToState(ctx context.Context, node *client.AuthenticationPolicyTreeNode) (types.Object, diag.Diagnostics) {
	return recursiveState(ctx, node, 1, GetRootNodeAttrTypes())
}

func recursiveState(ctx context.Context, node *client.AuthenticationPolicyTreeNode, depth int, attrTypes map[string]attr.Type) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if node == nil {
		diags.AddError(providererror.InternalProviderError, "provided authentication policy tree node is nil")
		return types.ObjectNull(attrTypes), diags
	}
	var attrValues = map[string]attr.Value{}

	attrValues["action"], diags = policyaction.ToState(ctx, &node.Action)
	if diags.HasError() {
		return types.ObjectNull(attrTypes), diags
	}

	if depth <= MaxPolicyNodeRecursiveDepth {
		childrenType := attrTypes["children"].(types.ListType).ElemType.(types.ObjectType)
		children := []attr.Value{}
		for i := range node.Children {
			childObj, diags := recursiveState(ctx, &node.Children[i], depth+1, childrenType.AttrTypes)
			if diags.HasError() {
				return types.ObjectNull(attrTypes), diags
			}
			children = append(children, childObj)
		}
		attrValues["children"], diags = types.ListValue(childrenType, children)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}
	}

	return types.ObjectValue(attrTypes, attrValues)
}
