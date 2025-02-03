// Copyright © 2025 Ping Identity Corporation

package policyaction

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
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
		err := json.Unmarshal([]byte(internaljson.FromValue(apcMappingPolicyAction, true)), &result.ApcMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		attributeSources, ok := apcMappingPolicyAction.(types.Object).Attributes()["attribute_mapping"].(types.Object).Attributes()["attribute_sources"]
		if ok {
			attributeSourceStruct, err := attributesources.ClientStruct(attributeSources.(types.Set))
			if err != nil {
				return nil, err
			}
			result.ApcMappingPolicyAction.AttributeMapping.AttributeSources = attributeSourceStruct
		}
		result.ApcMappingPolicyAction.Type = "APC_MAPPING"
		return &result, nil
	}
	authnSelectorPolicyAction, ok := attrs["authn_selector_policy_action"]
	if ok && internaltypes.IsDefined(authnSelectorPolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(authnSelectorPolicyAction, true)), &result.AuthnSelectorPolicyAction)
		if err != nil {
			return nil, err
		}
		result.AuthnSelectorPolicyAction.Type = "AUTHN_SELECTOR"
		return &result, nil
	}
	authnSourcePolicyAction, ok := attrs["authn_source_policy_action"]
	if ok && internaltypes.IsDefined(authnSourcePolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(authnSourcePolicyAction, true)), &result.AuthnSourcePolicyAction)
		if err != nil {
			return nil, err
		}
		result.AuthnSourcePolicyAction.Type = "AUTHN_SOURCE"
		return &result, nil
	}
	continuePolicyAction, ok := attrs["continue_policy_action"]
	if ok && internaltypes.IsDefined(continuePolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(continuePolicyAction, true)), &result.ContinuePolicyAction)
		if err != nil {
			return nil, err
		}
		result.ContinuePolicyAction.Type = "CONTINUE"
		return &result, nil
	}
	donePolicyAction, ok := attrs["done_policy_action"]
	if ok && internaltypes.IsDefined(donePolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(donePolicyAction, true)), &result.DonePolicyAction)
		if err != nil {
			return nil, err
		}
		result.DonePolicyAction.Type = "DONE"
		return &result, nil
	}
	fragmentPolicyAction, ok := attrs["fragment_policy_action"]
	if ok && internaltypes.IsDefined(fragmentPolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(fragmentPolicyAction, true)), &result.FragmentPolicyAction)
		if err != nil {
			return nil, err
		}

		attributeSources, ok := fragmentPolicyAction.(types.Object).Attributes()["fragment_mapping"].(types.Object).Attributes()["attribute_sources"]
		if ok {
			attributeSourceStruct, err := attributesources.ClientStruct(attributeSources.(types.Set))
			if err != nil {
				return nil, err
			}
			result.FragmentPolicyAction.FragmentMapping.AttributeSources = attributeSourceStruct
		}
		result.FragmentPolicyAction.Type = "FRAGMENT"
		return &result, nil
	}
	localIdentityMappingPolicyAction, ok := attrs["local_identity_mapping_policy_action"]
	if ok && internaltypes.IsDefined(localIdentityMappingPolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(localIdentityMappingPolicyAction, true)), &result.LocalIdentityMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		outboundAttributeMappingAttributeSources, ok := localIdentityMappingPolicyAction.(types.Object).Attributes()["outbound_attribute_mapping"].(types.Object).Attributes()["attribute_sources"]
		if ok {
			attributeSourceStruct, err := attributesources.ClientStruct(outboundAttributeMappingAttributeSources.(types.Set))
			if err != nil {
				return nil, err
			}
			result.LocalIdentityMappingPolicyAction.OutboundAttributeMapping.AttributeSources = attributeSourceStruct
		}

		inboundAttributeMappingAttributeSources, ok := localIdentityMappingPolicyAction.(types.Object).Attributes()["inbound_mapping"].(types.Object).Attributes()["attribute_sources"]
		if ok {
			attributeSourceStruct, err := attributesources.ClientStruct(inboundAttributeMappingAttributeSources.(types.Set))
			if err != nil {
				return nil, err
			}
			result.LocalIdentityMappingPolicyAction.InboundMapping.AttributeSources = attributeSourceStruct
		}

		result.LocalIdentityMappingPolicyAction.Type = "LOCAL_IDENTITY_MAPPING"
		return &result, nil
	}
	restartPolicyAction, ok := attrs["restart_policy_action"]
	if ok && internaltypes.IsDefined(restartPolicyAction) {
		err := json.Unmarshal([]byte(internaljson.FromValue(restartPolicyAction, true)), &result.RestartPolicyAction)
		if err != nil {
			return nil, err
		}
		result.RestartPolicyAction.Type = "RESTART"
		return &result, nil
	}

	return nil, errors.New("no valid policy action type found when building client struct. Ensure you have specified an action value in your policy tree node")
}
