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
	//TODO children
	rootNode.Action = *action
	return &rootNode, nil
}
