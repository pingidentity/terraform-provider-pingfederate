package authenticationpoliciespolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/authenticationpolicytreenode"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

type authenticationPoliciesPolicyModel struct {
	Id                              types.String `tfsdk:"id"`
	PolicyId                        types.String `tfsdk:"policy_id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	AuthenticationApiApplicationRef types.Object `tfsdk:"authentication_api_application_ref"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
	RootNode                        types.Object `tfsdk:"root_node"`
	HandleFailuresLocally           types.Bool   `tfsdk:"handle_failures_locally"`
}

func readAuthenticationPolicyResponse(ctx context.Context, r *client.AuthenticationPolicyTree, state *authenticationPoliciesPolicyModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.Id = types.StringPointerValue(r.Id)
	state.PolicyId = types.StringPointerValue(r.Id)
	state.Name = types.StringPointerValue(r.Name)
	state.Description = types.StringPointerValue(r.Description)
	state.Enabled = types.BoolPointerValue(r.Enabled)
	state.HandleFailuresLocally = types.BoolPointerValue(r.HandleFailuresLocally)

	state.AuthenticationApiApplicationRef, respDiags = resourcelink.ToState(ctx, r.AuthenticationApiApplicationRef)
	diags.Append(respDiags...)

	state.RootNode, respDiags = authenticationpolicytreenode.ToState(ctx, r.RootNode)
	diags.Append(respDiags...)

	return diags
}
