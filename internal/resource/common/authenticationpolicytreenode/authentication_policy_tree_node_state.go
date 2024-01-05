package authenticationpolicytreenode

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
)

var (
	rootNodeAttrTypes = map[string]attr.Type{}
)

func State(node *client.AuthenticationPolicyTreeNode) (types.Object, diag.Diagnostics) {
	var attrValues = map[string]attr.Value{}

	return types.ObjectValue(rootNodeAttrTypes, attrValues)
}
