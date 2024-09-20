package policyaction

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributemapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var (
	simplePolicyActionAttrTypes = map[string]attr.Type{
		"context": types.StringType,
	}
	attributeRulesAttrTypes = map[string]attr.Type{
		"fallback_to_success": types.BoolType,
		"items": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"attribute_name":   types.StringType,
			"attribute_source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrTypes()},
			"condition":        types.StringType,
			"expected_value":   types.StringType,
			"expression":       types.StringType,
			"result":           types.StringType,
		}}},
	}

	apcMappingPolicyActionAttrTypes = map[string]attr.Type{
		"context":                            types.StringType,
		"attribute_mapping":                  types.ObjectType{AttrTypes: attributemapping.AttrTypes()},
		"authentication_policy_contract_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	authnSelectorPolicyActionAttrTypes = map[string]attr.Type{
		"context":                     types.StringType,
		"authentication_selector_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
	authenticationSourceAttrTypes = map[string]attr.Type{
		"source_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"type":       types.StringType,
	}
	inputUserIdMappingAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{AttrTypes: sourcetypeidkey.AttrTypes()},
		"value":  types.StringType,
	}
	authnSourcePolicyActionAttrTypes = map[string]attr.Type{
		"context":               types.StringType,
		"attribute_rules":       types.ObjectType{AttrTypes: attributeRulesAttrTypes},
		"authentication_source": types.ObjectType{AttrTypes: authenticationSourceAttrTypes},
		"input_user_id_mapping": types.ObjectType{AttrTypes: inputUserIdMappingAttrTypes},
		"user_id_authenticated": types.BoolType,
	}
	fragmentPolicyActionAttrTypes = map[string]attr.Type{
		"context":          types.StringType,
		"attribute_rules":  types.ObjectType{AttrTypes: attributeRulesAttrTypes},
		"fragment":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"fragment_mapping": types.ObjectType{AttrTypes: attributemapping.AttrTypes()},
	}
	localIdentityMappingPolicyActionAttrTypes = map[string]attr.Type{
		"context":                    types.StringType,
		"inbound_mapping":            types.ObjectType{AttrTypes: attributemapping.AttrTypes()},
		"local_identity_ref":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"outbound_attribute_mapping": types.ObjectType{AttrTypes: attributemapping.AttrTypes()},
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

func ToState(ctx context.Context, response *client.PolicyActionAggregation) (types.Object, diag.Diagnostics) {
	var diags, respDiags diag.Diagnostics
	if response == nil {
		diags.AddError(providererror.InternalProviderError, "provided client struct is nil")
		return types.ObjectNull(policyActionAttrTypes), diags
	}

	// Default all to null
	attrs := map[string]attr.Value{
		"apc_mapping_policy_action":            types.ObjectNull(apcMappingPolicyActionAttrTypes),
		"authn_selector_policy_action":         types.ObjectNull(authnSelectorPolicyActionAttrTypes),
		"authn_source_policy_action":           types.ObjectNull(authnSourcePolicyActionAttrTypes),
		"continue_policy_action":               types.ObjectNull(simplePolicyActionAttrTypes),
		"done_policy_action":                   types.ObjectNull(simplePolicyActionAttrTypes),
		"fragment_policy_action":               types.ObjectNull(fragmentPolicyActionAttrTypes),
		"local_identity_mapping_policy_action": types.ObjectNull(localIdentityMappingPolicyActionAttrTypes),
		"restart_policy_action":                types.ObjectNull(simplePolicyActionAttrTypes),
	}

	if response.ApcMappingPolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.ApcMappingPolicyAction.Context),
		}
		actionAttrs["authentication_policy_contract_ref"], respDiags = resourcelink.ToState(ctx, &response.ApcMappingPolicyAction.AuthenticationPolicyContractRef)
		diags.Append(respDiags...)
		actionAttrs["attribute_mapping"], respDiags = attributemapping.ToState(ctx, &response.ApcMappingPolicyAction.AttributeMapping)
		diags.Append(respDiags...)
		attrs["apc_mapping_policy_action"], respDiags = types.ObjectValue(apcMappingPolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.AuthnSelectorPolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.AuthnSelectorPolicyAction.Context),
		}
		actionAttrs["authentication_selector_ref"], respDiags = resourcelink.ToState(ctx, &response.AuthnSelectorPolicyAction.AuthenticationSelectorRef)
		diags.Append(respDiags...)
		attrs["authn_selector_policy_action"], respDiags = types.ObjectValue(authnSelectorPolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.AuthnSourcePolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.AuthnSourcePolicyAction.Context),
		}
		actionAttrs["attribute_rules"], respDiags = types.ObjectValueFrom(ctx, attributeRulesAttrTypes, response.AuthnSourcePolicyAction.AttributeRules)
		diags.Append(respDiags...)
		actionAttrs["authentication_source"], respDiags = types.ObjectValueFrom(ctx, authenticationSourceAttrTypes, response.AuthnSourcePolicyAction.AuthenticationSource)
		diags.Append(respDiags...)
		actionAttrs["input_user_id_mapping"], respDiags = types.ObjectValueFrom(ctx, inputUserIdMappingAttrTypes, response.AuthnSourcePolicyAction.InputUserIdMapping)
		diags.Append(respDiags...)
		actionAttrs["user_id_authenticated"] = types.BoolPointerValue(response.AuthnSourcePolicyAction.UserIdAuthenticated)
		attrs["authn_source_policy_action"], respDiags = types.ObjectValue(authnSourcePolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.ContinuePolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.ContinuePolicyAction.Context),
		}
		attrs["continue_policy_action"], respDiags = types.ObjectValue(simplePolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.DonePolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.DonePolicyAction.Context),
		}
		attrs["done_policy_action"], respDiags = types.ObjectValue(simplePolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.FragmentPolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.FragmentPolicyAction.Context),
		}
		actionAttrs["attribute_rules"], respDiags = types.ObjectValueFrom(ctx, attributeRulesAttrTypes, response.FragmentPolicyAction.AttributeRules)
		diags.Append(respDiags...)
		actionAttrs["fragment"], respDiags = resourcelink.ToState(ctx, &response.FragmentPolicyAction.Fragment)
		diags.Append(respDiags...)
		actionAttrs["fragment_mapping"], respDiags = attributemapping.ToState(ctx, response.FragmentPolicyAction.FragmentMapping)
		diags.Append(respDiags...)
		attrs["fragment_policy_action"], respDiags = types.ObjectValue(fragmentPolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.LocalIdentityMappingPolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.LocalIdentityMappingPolicyAction.Context),
		}
		actionAttrs["local_identity_ref"], respDiags = resourcelink.ToState(ctx, &response.LocalIdentityMappingPolicyAction.LocalIdentityRef)
		diags.Append(respDiags...)
		actionAttrs["inbound_mapping"], respDiags = attributemapping.ToState(ctx, response.LocalIdentityMappingPolicyAction.InboundMapping)
		diags.Append(respDiags...)
		actionAttrs["outbound_attribute_mapping"], respDiags = attributemapping.ToState(ctx, &response.LocalIdentityMappingPolicyAction.OutboundAttributeMapping)
		diags.Append(respDiags...)
		attrs["local_identity_mapping_policy_action"], respDiags = types.ObjectValue(localIdentityMappingPolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else if response.RestartPolicyAction != nil {
		actionAttrs := map[string]attr.Value{
			"context": types.StringPointerValue(response.RestartPolicyAction.Context),
		}
		attrs["restart_policy_action"], respDiags = types.ObjectValue(simplePolicyActionAttrTypes, actionAttrs)
		diags.Append(respDiags...)

	} else {
		diags.AddError(providererror.InternalProviderError, "no valid non-nil policy action type found in struct")
	}

	if diags.HasError() {
		return types.ObjectNull(policyActionAttrTypes), diags
	}

	return types.ObjectValue(policyActionAttrTypes, attrs)
}
