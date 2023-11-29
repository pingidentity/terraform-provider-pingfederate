package oauthtokenexchangetokengeneratormapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
)

type oauthTokenExchangeTokenGeneratorMappingModel struct {
	AttributeSources                 types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	MappingId                        types.String `tfsdk:"mapping_id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

func readOauthTokenExchangeTokenGeneratorMappingResponse(ctx context.Context, r *client.ProcessorPolicyToGeneratorMapping, state *oauthTokenExchangeTokenGeneratorMappingModel, plan oauthTokenExchangeTokenGeneratorMappingModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.AttributeSources, respDiags = attributesources.ToState(ctx, r.AttributeSources)
	diags.Append(respDiags...)
	state.AttributeContractFulfillment, respDiags = attributecontractfulfillment.ToState(ctx, r.AttributeContractFulfillment)
	diags.Append(respDiags...)
	state.IssuanceCriteria, respDiags = issuancecriteria.ToState(ctx, r.IssuanceCriteria)
	diags.Append(respDiags...)
	state.SourceId = types.StringValue(r.SourceId)
	state.TargetId = types.StringValue(r.TargetId)
	state.Id = types.StringPointerValue(r.Id)
	state.MappingId = types.StringPointerValue(r.Id)
	state.LicenseConnectionGroupAssignment = types.StringPointerValue(r.LicenseConnectionGroupAssignment)
	return diags
}
