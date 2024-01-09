package authenticationpolicytreenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

var rootNodeAttrTypes map[string]attr.Type
var rootNodeAttrTypesBuilt = false

func getRootNodeAttrTypes() map[string]attr.Type {
	if rootNodeAttrTypesBuilt {
		return rootNodeAttrTypes
	}
	rootNodeAttrTypes = map[string]attr.Type{
		"action":   types.ObjectType{AttrTypes: policyaction.AttrTypes()},
		"children": buildRootNodeAttrTypesChildren(1),
	}
	rootNodeAttrTypesBuilt = true
	return rootNodeAttrTypes
}

func buildRootNodeAttrTypesChildren(depth int) types.ListType {
	attrs := map[string]attr.Type{
		"action": types.ObjectType{AttrTypes: policyaction.AttrTypes()},
	}
	if depth < maxRecursiveDepth {
		attrs["children"] = buildRootNodeAttrTypesChildren(depth + 1)
	}
	return types.ListType{
		ElemType: types.ObjectType{AttrTypes: attrs},
	}
}

func State(ctx context.Context, node *client.AuthenticationPolicyTreeNode) (types.Object, diag.Diagnostics) {
	return recursiveState(ctx, node, 1, getRootNodeAttrTypes())
}

func recursiveState(ctx context.Context, node *client.AuthenticationPolicyTreeNode, depth int, attrTypes map[string]attr.Type) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if node == nil {
		diags.AddError("provided authentication policy tree node is nil", "")
		return types.ObjectNull(attrTypes), diags
	}
	var attrValues = map[string]attr.Value{}

	attrValues["action"], diags = policyaction.State(ctx, &node.Action)
	if diags.HasError() {
		return types.ObjectNull(attrTypes), diags
	}

	if depth <= maxRecursiveDepth {
		childrenType := attrTypes["children"].(types.ListType).ElemType.(types.ObjectType)
		children := []attr.Value{}
		for _, child := range node.Children {
			childObj, diags := recursiveState(ctx, &child, depth+1, childrenType.AttrTypes)
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
