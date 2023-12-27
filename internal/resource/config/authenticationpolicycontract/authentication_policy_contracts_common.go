package authenticationpolicycontract

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
)

var attributeElemAttrType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	},
}

type authenticationPolicyContractModel struct {
	Id                 types.String `tfsdk:"id"`
	ContractId         types.String `tfsdk:"contract_id"`
	Name               types.String `tfsdk:"name"`
	CoreAttributes     types.List   `tfsdk:"core_attributes"`
	ExtendedAttributes types.Set    `tfsdk:"extended_attributes"`
	LastModified       types.String `tfsdk:"last_modified"`
}

func readAuthenticationPolicyContractsResponse(ctx context.Context, r *client.AuthenticationPolicyContract, state *authenticationPolicyContractModel, expectedValues *authenticationPolicyContractModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringPointerValue(r.Id)
	state.ContractId = types.StringPointerValue(r.Id)
	state.Name = types.StringPointerValue(r.Name)

	state.CoreAttributes, respDiags = types.ListValueFrom(ctx, attributeElemAttrType, r.GetCoreAttributes())
	diags.Append(respDiags...)

	state.ExtendedAttributes, respDiags = types.SetValueFrom(ctx, attributeElemAttrType, r.GetExtendedAttributes())
	diags.Append(respDiags...)

	if r.LastModified == nil {
		state.LastModified = types.StringNull()
	} else {
		state.LastModified = types.StringValue(r.LastModified.Format(time.RFC3339))
	}

	return diags
}
