package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributemapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	attributeAttrTypes = map[string]attr.Type{
		"name":                 types.StringType,
		"include_in_id_token":  types.BoolType,
		"include_in_user_info": types.BoolType,
		"multi_valued":         types.BoolType,
	}
	attributesSetAttrType = types.SetType{
		ElemType: types.ObjectType{AttrTypes: attributeAttrTypes},
	}

	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes":     attributesSetAttrType,
		"extended_attributes": attributesSetAttrType,
	}

	scopeAttributeMappingsElemAttrTypes = map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}

	scopeAttributeMappingsDefault, _  = types.MapValue(types.ObjectType{AttrTypes: scopeAttributeMappingsElemAttrTypes}, nil)
	emptyExtendedAttributesDefault, _ = types.SetValue(types.ObjectType{AttrTypes: attributeAttrTypes}, nil)
)

type oauthOpenIdConnectPolicyModel struct {
	Id                          types.String `tfsdk:"id"`
	PolicyId                    types.String `tfsdk:"policy_id"`
	Name                        types.String `tfsdk:"name"`
	AccessTokenManagerRef       types.Object `tfsdk:"access_token_manager_ref"`
	IdTokenLifetime             types.Int64  `tfsdk:"id_token_lifetime"`
	IncludeSriInIdToken         types.Bool   `tfsdk:"include_sri_in_id_token"`
	IncludeUserInfoInIdToken    types.Bool   `tfsdk:"include_user_info_in_id_token"`
	IncludeSHashInIdToken       types.Bool   `tfsdk:"include_s_hash_in_id_token"`
	ReturnIdTokenOnRefreshGrant types.Bool   `tfsdk:"return_id_token_on_refresh_grant"`
	ReissueIdTokenInHybridFlow  types.Bool   `tfsdk:"reissue_id_token_in_hybrid_flow"`
	AttributeContract           types.Object `tfsdk:"attribute_contract"`
	AttributeMapping            types.Object `tfsdk:"attribute_mapping"`
	ScopeAttributeMappings      types.Map    `tfsdk:"scope_attribute_mappings"`
	IncludeX5tInIdToken         types.Bool   `tfsdk:"include_x5t_in_id_token"`
	IdTokenTypHeaderValue       types.String `tfsdk:"id_token_typ_header_value"`
}

func readOauthOpenIdConnectPolicyResponse(ctx context.Context, response *client.OpenIdConnectPolicy, state *oauthOpenIdConnectPolicyModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(response.Id)
	state.PolicyId = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)

	state.AccessTokenManagerRef, diags = resourcelink.ToState(ctx, &response.AccessTokenManagerRef)
	respDiags.Append(diags...)

	state.IdTokenLifetime = types.Int64PointerValue(response.IdTokenLifetime)
	state.IncludeSriInIdToken = types.BoolPointerValue(response.IncludeSriInIdToken)
	state.IncludeUserInfoInIdToken = types.BoolPointerValue(response.IncludeUserInfoInIdToken)
	state.IncludeSHashInIdToken = types.BoolPointerValue(response.IncludeSHashInIdToken)
	state.ReturnIdTokenOnRefreshGrant = types.BoolPointerValue(response.ReturnIdTokenOnRefreshGrant)
	state.ReissueIdTokenInHybridFlow = types.BoolPointerValue(response.ReissueIdTokenInHybridFlow)
	state.IncludeX5tInIdToken = types.BoolPointerValue(response.IncludeX5tInIdToken)
	if response.IdTokenTypHeaderValue != nil && *response.IdTokenTypHeaderValue == "" {
		// PF can return an empty string for a nil value, so treat that as null here
		state.IdTokenTypHeaderValue = types.StringNull()
	} else {
		state.IdTokenTypHeaderValue = types.StringPointerValue(response.IdTokenTypHeaderValue)
	}

	state.AttributeContract, diags = types.ObjectValueFrom(ctx, attributeContractAttrTypes, response.AttributeContract)
	respDiags.Append(diags...)

	// Attribute mapping
	attributeMappingValues := map[string]attr.Value{}

	// Build attribute_contract_fulfillment value
	attributeMappingValues["attribute_contract_fulfillment"], diags = attributecontractfulfillment.ToState(ctx, &response.AttributeMapping.AttributeContractFulfillment)
	respDiags.Append(diags...)

	// Build issuance_criteria value
	attributeMappingValues["issuance_criteria"], diags = issuancecriteria.ToState(ctx, response.AttributeMapping.IssuanceCriteria)
	respDiags.Append(diags...)

	// Build attribute_sources value
	attributeMappingValues["attribute_sources"], diags = attributesources.ToState(ctx, response.AttributeMapping.AttributeSources)
	respDiags.Append(diags...)

	// Build complete attribute mapping value
	state.AttributeMapping, diags = types.ObjectValue(attributemapping.AttrTypes(), attributeMappingValues)
	respDiags.Append(diags...)

	state.ScopeAttributeMappings, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: scopeAttributeMappingsElemAttrTypes}, response.ScopeAttributeMappings)
	respDiags.Append(diags...)
	return respDiags
}
