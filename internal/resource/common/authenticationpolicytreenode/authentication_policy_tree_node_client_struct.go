package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
)

func ClientStruct(planNode types.Object) (*client.AuthenticationPolicyTreeNode, error) {
	//TODO nil/undefined checks
	rootNode := client.AuthenticationPolicyTreeNode{}
	rootNodeAttrs := planNode.Attributes()
	action, err := policyaction.ClientStruct(rootNodeAttrs["policy_action"].(types.Object))
	if err != nil {
		return nil, err
	}
	//TODO children
	rootNode.Action = *action
	return &rootNode, nil
}
