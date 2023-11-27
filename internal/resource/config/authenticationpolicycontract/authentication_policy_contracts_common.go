package authenticationpolicycontract

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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
}

func readAuthenticationPolicyContractsResponse(ctx context.Context, r *client.AuthenticationPolicyContract, state *authenticationPolicyContractModel, expectedValues *authenticationPolicyContractModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.ContractId = internaltypes.StringTypeOrNil(r.Id, false)
	state.Name = internaltypes.StringTypeOrNil(r.Name, false)

	state.CoreAttributes, respDiags = types.ListValueFrom(ctx, attributeElemAttrType, r.GetCoreAttributes())
	diags.Append(respDiags...)

	state.ExtendedAttributes, respDiags = types.SetValueFrom(ctx, attributeElemAttrType, r.GetExtendedAttributes())
	diags.Append(respDiags...)

	return diags
}
