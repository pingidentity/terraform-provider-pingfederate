package authenticationpolicytreenode

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(planNode types.Object) (*client.AuthenticationPolicyTreeNode, error) {
	if !internaltypes.IsDefined(planNode) {
		return nil, errors.New("plan authentication policy tree node is not defined")
	}

	rootNode := client.AuthenticationPolicyTreeNode{}
	rootNodeAttrs := planNode.Attributes()
	policyAction, ok := rootNodeAttrs["policy_action"]
	if !ok {
		return nil, errors.New("policy_action attribute not defined in plan authentication policy tree node")
	}
	action, err := policyaction.ClientStruct(policyAction.(types.Object))
	if err != nil {
		return nil, err
	}
	children, ok := rootNodeAttrs["children"]
	if ok { // If there is a children attribute, read the children recursively
		rootNode.Children, err = getChildren(children.(types.List))
		if err != nil {
			return nil, err
		}
	}
	rootNode.Action = *action
	return &rootNode, nil
}

func getChildren(planChildren types.List) ([]client.AuthenticationPolicyTreeNode, error) {
	children := []client.AuthenticationPolicyTreeNode{}

	for _, child := range planChildren.Elements() {
		childObj, ok := child.(types.Object)
		if !ok {
			return []client.AuthenticationPolicyTreeNode{}, errors.New("child policy tree node has invalid type - unable to cast to ObjectType")
		}
		childStruct, err := ClientStruct(childObj)
		if err != nil {
			return []client.AuthenticationPolicyTreeNode{}, err
		}
		children = append(children, *childStruct)
	}

	return children, nil
}
