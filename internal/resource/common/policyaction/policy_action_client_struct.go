// Copyright © 2025 Ping Identity Corporation

package policyaction

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(object types.Object) (*client.PolicyActionAggregation, error) {
	if !internaltypes.IsDefined(object) {
		return nil, errors.New("provided policy action object is Null or Unknown")
	}

	result := client.PolicyActionAggregation{}
	attrs := object.Attributes()
	apcMappingPolicyAction, ok := attrs["apc_mapping_policy_action"]
	if ok && internaltypes.IsDefined(apcMappingPolicyAction) {
		result.ApcMappingPolicyAction = &client.ApcMappingPolicyAction{}
		apcMappingPolicyActionAttrs := apcMappingPolicyAction.(types.Object).Attributes()
		attributeMappingValue := client.AttributeMapping{}
		attributeMappingAttrs := apcMappingPolicyActionAttrs["attribute_mapping"].(types.Object).Attributes()
		attributeMappingValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(attributeMappingAttrs["attribute_contract_fulfillment"].(types.Map))
		attributeMappingValue.AttributeSources = attributesources.ClientStruct(attributeMappingAttrs["attribute_sources"].(types.List))
		attributeMappingValue.IssuanceCriteria = issuancecriteria.ClientStruct(attributeMappingAttrs["issuance_criteria"].(types.Object))
		result.ApcMappingPolicyAction.AttributeMapping = attributeMappingValue
		authenticationPolicyContractRefValue := client.ResourceLink{}
		authenticationPolicyContractRefAttrs := apcMappingPolicyActionAttrs["authentication_policy_contract_ref"].(types.Object).Attributes()
		authenticationPolicyContractRefValue.Id = authenticationPolicyContractRefAttrs["id"].(types.String).ValueString()
		result.ApcMappingPolicyAction.AuthenticationPolicyContractRef = authenticationPolicyContractRefValue
		result.ApcMappingPolicyAction.Context = apcMappingPolicyActionAttrs["context"].(types.String).ValueStringPointer()
		result.ApcMappingPolicyAction.Type = "APC_MAPPING"
		return &result, nil
	}
	authnSelectorPolicyAction, ok := attrs["authn_selector_policy_action"]
	if ok && internaltypes.IsDefined(authnSelectorPolicyAction) {
		result.AuthnSelectorPolicyAction = &client.AuthnSelectorPolicyAction{}
		authnSelectorPolicyActionAttrs := authnSelectorPolicyAction.(types.Object).Attributes()
		authenticationSelectorRefValue := client.ResourceLink{}
		authenticationSelectorRefAttrs := authnSelectorPolicyActionAttrs["authentication_selector_ref"].(types.Object).Attributes()
		authenticationSelectorRefValue.Id = authenticationSelectorRefAttrs["id"].(types.String).ValueString()
		result.AuthnSelectorPolicyAction.AuthenticationSelectorRef = authenticationSelectorRefValue
		result.AuthnSelectorPolicyAction.Context = authnSelectorPolicyActionAttrs["context"].(types.String).ValueStringPointer()
		result.AuthnSelectorPolicyAction.Type = "AUTHN_SELECTOR"
		return &result, nil
	}
	authnSourcePolicyAction, ok := attrs["authn_source_policy_action"]
	if ok && internaltypes.IsDefined(authnSourcePolicyAction) {
		result.AuthnSourcePolicyAction = &client.AuthnSourcePolicyAction{}
		authnSourcePolicyActionAttrs := authnSourcePolicyAction.(types.Object).Attributes()
		if !authnSourcePolicyActionAttrs["attribute_rules"].IsNull() && !authnSourcePolicyActionAttrs["attribute_rules"].IsUnknown() {
			attributeRulesValue := &client.AttributeRules{}
			attributeRulesAttrs := authnSourcePolicyActionAttrs["attribute_rules"].(types.Object).Attributes()
			attributeRulesValue.FallbackToSuccess = attributeRulesAttrs["fallback_to_success"].(types.Bool).ValueBoolPointer()
			if !attributeRulesAttrs["items"].IsNull() && !attributeRulesAttrs["items"].IsUnknown() {
				attributeRulesValue.Items = []client.AttributeRule{}
				for _, itemsElement := range attributeRulesAttrs["items"].(types.List).Elements() {
					itemsValue := client.AttributeRule{}
					itemsAttrs := itemsElement.(types.Object).Attributes()
					itemsValue.AttributeName = itemsAttrs["attribute_name"].(types.String).ValueStringPointer()
					if !itemsAttrs["attribute_source"].IsNull() && !itemsAttrs["attribute_source"].IsUnknown() {
						itemsAttributeSourceValue := &client.SourceTypeIdKey{}
						itemsAttributeSourceAttrs := itemsAttrs["attribute_source"].(types.Object).Attributes()
						itemsAttributeSourceValue.Id = itemsAttributeSourceAttrs["id"].(types.String).ValueStringPointer()
						itemsAttributeSourceValue.Type = itemsAttributeSourceAttrs["type"].(types.String).ValueString()
						itemsValue.AttributeSource = itemsAttributeSourceValue
					}
					itemsValue.Condition = itemsAttrs["condition"].(types.String).ValueStringPointer()
					itemsValue.ExpectedValue = itemsAttrs["expected_value"].(types.String).ValueStringPointer()
					itemsValue.Expression = itemsAttrs["expression"].(types.String).ValueStringPointer()
					itemsValue.Result = itemsAttrs["result"].(types.String).ValueString()
					attributeRulesValue.Items = append(attributeRulesValue.Items, itemsValue)
				}
			}
			result.AuthnSourcePolicyAction.AttributeRules = attributeRulesValue
		}
		authenticationSourceValue := client.AuthenticationSource{}
		authenticationSourceAttrs := authnSourcePolicyActionAttrs["authentication_source"].(types.Object).Attributes()
		sourceRefValue := client.ResourceLink{}
		sourceSourceRefAttrs := authenticationSourceAttrs["source_ref"].(types.Object).Attributes()
		sourceRefValue.Id = sourceSourceRefAttrs["id"].(types.String).ValueString()
		authenticationSourceValue.SourceRef = sourceRefValue
		authenticationSourceValue.Type = authenticationSourceAttrs["type"].(types.String).ValueString()
		result.AuthnSourcePolicyAction.AuthenticationSource = authenticationSourceValue
		result.AuthnSourcePolicyAction.Context = authnSourcePolicyActionAttrs["context"].(types.String).ValueStringPointer()
		if !authnSourcePolicyActionAttrs["input_user_id_mapping"].IsNull() && !authnSourcePolicyActionAttrs["input_user_id_mapping"].IsUnknown() {
			inputUserIdMappingValue := &client.AttributeFulfillmentValue{}
			inputUserIdMappingAttrs := authnSourcePolicyActionAttrs["input_user_id_mapping"].(types.Object).Attributes()
			inputUserIdMappingSourceValue := client.SourceTypeIdKey{}
			inputUserIdMappingSourceAttrs := inputUserIdMappingAttrs["source"].(types.Object).Attributes()
			inputUserIdMappingSourceValue.Id = inputUserIdMappingSourceAttrs["id"].(types.String).ValueStringPointer()
			inputUserIdMappingSourceValue.Type = inputUserIdMappingSourceAttrs["type"].(types.String).ValueString()
			inputUserIdMappingValue.Source = inputUserIdMappingSourceValue
			inputUserIdMappingValue.Value = inputUserIdMappingAttrs["value"].(types.String).ValueString()
			result.AuthnSourcePolicyAction.InputUserIdMapping = inputUserIdMappingValue
		}
		result.AuthnSourcePolicyAction.Type = "AUTHN_SOURCE"
		result.AuthnSourcePolicyAction.UserIdAuthenticated = authnSourcePolicyActionAttrs["user_id_authenticated"].(types.Bool).ValueBoolPointer()
		return &result, nil
	}
	continuePolicyAction, ok := attrs["continue_policy_action"]
	if ok && internaltypes.IsDefined(continuePolicyAction) {
		result.ContinuePolicyAction = &client.ContinuePolicyAction{}
		continuePolicyActionAttrs := continuePolicyAction.(types.Object).Attributes()
		result.ContinuePolicyAction.Context = continuePolicyActionAttrs["context"].(types.String).ValueStringPointer()
		result.ContinuePolicyAction.Type = "CONTINUE"
		return &result, nil
	}
	donePolicyAction, ok := attrs["done_policy_action"]
	if ok && internaltypes.IsDefined(donePolicyAction) {
		result.DonePolicyAction = &client.DonePolicyAction{}
		donePolicyActionAttrs := donePolicyAction.(types.Object).Attributes()
		result.DonePolicyAction.Context = donePolicyActionAttrs["context"].(types.String).ValueStringPointer()
		result.DonePolicyAction.Type = "DONE"
		return &result, nil
	}
	fragmentPolicyAction, ok := attrs["fragment_policy_action"]
	if ok && internaltypes.IsDefined(fragmentPolicyAction) {
		result.FragmentPolicyAction = &client.FragmentPolicyAction{}
		fragmentPolicyActionAttributes := fragmentPolicyAction.(types.Object).Attributes()
		if !fragmentPolicyActionAttributes["attribute_rules"].IsNull() && !fragmentPolicyActionAttributes["attribute_rules"].IsUnknown() {
			attributeRulesValue := &client.AttributeRules{}
			attributeRulesAttrs := fragmentPolicyActionAttributes["attribute_rules"].(types.Object).Attributes()
			attributeRulesValue.FallbackToSuccess = attributeRulesAttrs["fallback_to_success"].(types.Bool).ValueBoolPointer()
			if !attributeRulesAttrs["items"].IsNull() && !attributeRulesAttrs["items"].IsUnknown() {
				attributeRulesValue.Items = []client.AttributeRule{}
				for _, itemsElement := range attributeRulesAttrs["items"].(types.List).Elements() {
					itemsValue := client.AttributeRule{}
					itemsAttrs := itemsElement.(types.Object).Attributes()
					itemsValue.AttributeName = itemsAttrs["attribute_name"].(types.String).ValueStringPointer()
					if !itemsAttrs["attribute_source"].IsNull() && !itemsAttrs["attribute_source"].IsUnknown() {
						itemsAttributeSourceValue := &client.SourceTypeIdKey{}
						itemsAttributeSourceAttrs := itemsAttrs["attribute_source"].(types.Object).Attributes()
						itemsAttributeSourceValue.Id = itemsAttributeSourceAttrs["id"].(types.String).ValueStringPointer()
						itemsAttributeSourceValue.Type = itemsAttributeSourceAttrs["type"].(types.String).ValueString()
						itemsValue.AttributeSource = itemsAttributeSourceValue
					}
					itemsValue.Condition = itemsAttrs["condition"].(types.String).ValueStringPointer()
					itemsValue.ExpectedValue = itemsAttrs["expected_value"].(types.String).ValueStringPointer()
					itemsValue.Expression = itemsAttrs["expression"].(types.String).ValueStringPointer()
					itemsValue.Result = itemsAttrs["result"].(types.String).ValueString()
					attributeRulesValue.Items = append(attributeRulesValue.Items, itemsValue)
				}
			}
			result.FragmentPolicyAction.AttributeRules = attributeRulesValue
		}
		result.FragmentPolicyAction.Context = fragmentPolicyActionAttributes["context"].(types.String).ValueStringPointer()
		fragmentValue := client.ResourceLink{}
		fragmentAttrs := fragmentPolicyActionAttributes["fragment"].(types.Object).Attributes()
		fragmentValue.Id = fragmentAttrs["id"].(types.String).ValueString()
		result.FragmentPolicyAction.Fragment = fragmentValue
		if !fragmentPolicyActionAttributes["fragment_mapping"].IsNull() && !fragmentPolicyActionAttributes["fragment_mapping"].IsUnknown() {
			fragmentMappingValue := &client.AttributeMapping{}
			fragmentMappingAttrs := fragmentPolicyActionAttributes["fragment_mapping"].(types.Object).Attributes()
			fragmentMappingValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(fragmentMappingAttrs["attribute_contract_fulfillment"].(types.Map))
			fragmentMappingValue.AttributeSources = attributesources.ClientStruct(fragmentMappingAttrs["attribute_sources"].(types.List))
			fragmentMappingValue.IssuanceCriteria = issuancecriteria.ClientStruct(fragmentMappingAttrs["issuance_criteria"].(types.Object))
			result.FragmentPolicyAction.FragmentMapping = fragmentMappingValue
		}
		result.FragmentPolicyAction.Type = "FRAGMENT"
		return &result, nil
	}
	localIdentityMappingPolicyAction, ok := attrs["local_identity_mapping_policy_action"]
	if ok && internaltypes.IsDefined(localIdentityMappingPolicyAction) {
		result.LocalIdentityMappingPolicyAction = &client.LocalIdentityMappingPolicyAction{}
		localIdentityMappingPolicyActionAttrs := localIdentityMappingPolicyAction.(types.Object).Attributes()
		result.LocalIdentityMappingPolicyAction.Context = localIdentityMappingPolicyActionAttrs["context"].(types.String).ValueStringPointer()
		if !localIdentityMappingPolicyActionAttrs["inbound_mapping"].IsNull() && !localIdentityMappingPolicyActionAttrs["inbound_mapping"].IsUnknown() {
			inboundMappingValue := &client.AttributeMapping{}
			inboundMappingAttrs := localIdentityMappingPolicyActionAttrs["inbound_mapping"].(types.Object).Attributes()
			inboundMappingValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(inboundMappingAttrs["attribute_contract_fulfillment"].(types.Map))
			inboundMappingValue.AttributeSources = attributesources.ClientStruct(inboundMappingAttrs["attribute_sources"].(types.List))
			inboundMappingValue.IssuanceCriteria = issuancecriteria.ClientStruct(inboundMappingAttrs["issuance_criteria"].(types.Object))
			result.LocalIdentityMappingPolicyAction.InboundMapping = inboundMappingValue
		}
		localIdentityRefValue := client.ResourceLink{}
		localIdentityRefAttrs := localIdentityMappingPolicyActionAttrs["local_identity_ref"].(types.Object).Attributes()
		localIdentityRefValue.Id = localIdentityRefAttrs["id"].(types.String).ValueString()
		result.LocalIdentityMappingPolicyAction.LocalIdentityRef = localIdentityRefValue
		outboundAttributeMappingValue := client.AttributeMapping{}
		outboundAttributeMappingAttrs := localIdentityMappingPolicyActionAttrs["outbound_attribute_mapping"].(types.Object).Attributes()
		outboundAttributeMappingValue.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(outboundAttributeMappingAttrs["attribute_contract_fulfillment"].(types.Map))
		outboundAttributeMappingValue.AttributeSources = attributesources.ClientStruct(outboundAttributeMappingAttrs["attribute_sources"].(types.List))
		outboundAttributeMappingValue.IssuanceCriteria = issuancecriteria.ClientStruct(outboundAttributeMappingAttrs["issuance_criteria"].(types.Object))
		result.LocalIdentityMappingPolicyAction.OutboundAttributeMapping = outboundAttributeMappingValue
		result.LocalIdentityMappingPolicyAction.Type = "LOCAL_IDENTITY_MAPPING"
		return &result, nil
	}
	restartPolicyAction, ok := attrs["restart_policy_action"]
	if ok && internaltypes.IsDefined(restartPolicyAction) {
		result.RestartPolicyAction = &client.RestartPolicyAction{}
		restartPolicyActionAttrs := restartPolicyAction.(types.Object).Attributes()
		result.RestartPolicyAction.Context = restartPolicyActionAttrs["context"].(types.String).ValueStringPointer()
		result.RestartPolicyAction.Type = "RESTART"
		return &result, nil
	}

	return nil, errors.New("no valid policy action type found when building client struct. Ensure you have specified an action value in your policy tree node")
}
