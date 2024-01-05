package policyaction

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

var (
	simplePolicyActionAttrTypes = map[string]attr.Type{
		"context": types.StringType,
	}
	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_contract_fulfillment": types.ObjectType{AttrTypes: attributecontractfulfillment.AttrType()},
		"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: attributesources.ElemAttrType(),
			},
		},
		"issuance_criteria": types.ObjectType{AttrTypes: issuancecriteria.AttrType()},
	}
	attributeRulesAttrTypes = map[string]attr.Type{
		"fallback_to_success": types.BoolType,
		"items": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"attribute_name":   types.StringType,
			"attribute_source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrType()},
			"condition":        types.StringType,
			"expected_value":   types.StringType,
			"expression":       types.StringType,
			"result":           types.StringType,
		}}},
	}

	apcMappingPolicyActionAttrTypes = map[string]attr.Type{
		"context":                            types.StringType,
		"attribute_mapping":                  types.ObjectType{AttrTypes: attributeMappingAttrTypes},
		"authentication_policy_contract_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	authnSelectorPolicyActionAttrTypes = map[string]attr.Type{
		"context":                     types.StringType,
		"authentication_selector_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	authnSourcePolicyActionAttrTypes = map[string]attr.Type{
		"context":         types.StringType,
		"attribute_rules": types.ObjectType{AttrTypes: attributeRulesAttrTypes},
		"authentication_source": types.ObjectType{AttrTypes: map[string]attr.Type{
			"source_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
			"type":       types.StringType,
		}},
		"input_user_id_mapping": types.ObjectType{AttrTypes: map[string]attr.Type{
			"source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrType()},
			"value":  types.StringType,
		}},
		"user_id_authenticated": types.BoolType,
	}
	fragmentPolicyActionAttrTypes = map[string]attr.Type{
		"context":          types.StringType,
		"attribute_rules":  types.ObjectType{AttrTypes: attributeRulesAttrTypes},
		"fragment":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"fragment_mapping": types.ObjectType{AttrTypes: attributeMappingAttrTypes},
	}
	localIdentityMappingPolicyActionAttrTypes = map[string]attr.Type{
		"context":                    types.StringType,
		"inbound_mapping":            types.ObjectType{AttrTypes: attributeMappingAttrTypes},
		"local_identity_ref":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"outbound_attribute_mapping": types.ObjectType{AttrTypes: attributeMappingAttrTypes},
	}

	policyActionAttrTypes = map[string]attr.Type{
		"apc_mapping_policy_action":            types.ObjectType{AttrTypes: apcMappingPolicyActionAttrTypes},
		"authn_selector_policy_action":         types.ObjectType{AttrTypes: authnSelectorPolicyActionAttrTypes},
		"authn_source_policy_action":           types.ObjectType{AttrTypes: authnSourcePolicyActionAttrTypes},
		"continue_policy_action":               types.ObjectType{AttrTypes: simplePolicyActionAttrTypes},
		"done_policy_action":                   types.ObjectType{AttrTypes: simplePolicyActionAttrTypes},
		"fragment_policy_action":               types.ObjectType{AttrTypes: fragmentPolicyActionAttrTypes},
		"local_identity_mapping_policy_action": types.ObjectType{AttrTypes: localIdentityMappingPolicyActionAttrTypes},
		"restart_policy_action":                types.ObjectType{AttrTypes: simplePolicyActionAttrTypes},
	}
)

func AttrTypes() map[string]attr.Type {
	return policyActionAttrTypes
}

func State(ctx context.Context, response *client.PolicyActionAggregation) (types.Object, diag.Diagnostics) {
	var diags, respDiags diag.Diagnostics
	if response == nil {
		diags.AddError("provided client struct is nil", "")
		return types.ObjectNull(policyActionAttrTypes), diags
	}

	attrs := map[string]attr.Value{}
	if response.ApcMappingPolicyAction != nil {
		attrs["apc_mapping_policy_action"], respDiags = types.ObjectValueFrom(ctx, apcMappingPolicyActionAttrTypes, response.ApcMappingPolicyAction)
		diags.Append(respDiags...)
	} else if response.AuthnSelectorPolicyAction != nil {
		attrs["authn_selector_policy_action"], respDiags = types.ObjectValueFrom(ctx, authnSelectorPolicyActionAttrTypes, response.AuthnSelectorPolicyAction)
		diags.Append(respDiags...)
	} else if response.AuthnSourcePolicyAction != nil {
		attrs["authn_source_policy_action"], respDiags = types.ObjectValueFrom(ctx, authnSourcePolicyActionAttrTypes, response.AuthnSourcePolicyAction)
		diags.Append(respDiags...)
	} else if response.ContinuePolicyAction != nil {
		attrs["continue_policy_action"], respDiags = types.ObjectValueFrom(ctx, simplePolicyActionAttrTypes, response.ContinuePolicyAction)
		diags.Append(respDiags...)
	} else if response.DonePolicyAction != nil {
		attrs["done_policy_action"], respDiags = types.ObjectValueFrom(ctx, simplePolicyActionAttrTypes, response.DonePolicyAction)
		diags.Append(respDiags...)
	} else if response.FragmentPolicyAction != nil {
		attrs["fragment_policy_action"], respDiags = types.ObjectValueFrom(ctx, fragmentPolicyActionAttrTypes, response.FragmentPolicyAction)
		diags.Append(respDiags...)
	} else if response.LocalIdentityMappingPolicyAction != nil {
		attrs["local_identity_mapping_policy_action"], respDiags = types.ObjectValueFrom(ctx, localIdentityMappingPolicyActionAttrTypes, response.LocalIdentityMappingPolicyAction)
		diags.Append(respDiags...)
	} else if response.RestartPolicyAction != nil {
		attrs["restart_policy_action"], respDiags = types.ObjectValueFrom(ctx, simplePolicyActionAttrTypes, response.RestartPolicyAction)
		diags.Append(respDiags...)
	} else {
		diags.AddError("no valid non-nil policy action type found in struct", "")
	}

	if diags.HasError() {
		return types.ObjectNull(policyActionAttrTypes), diags
	}

	return types.ObjectValue(policyActionAttrTypes, attrs)
}
