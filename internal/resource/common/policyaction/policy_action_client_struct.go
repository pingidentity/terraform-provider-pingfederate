package policyaction

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
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
		policyActionAttrs := apcMappingPolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("APC_MAPPING")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.ApcMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	authnSelectorPolicyAction, ok := attrs["authn_selector_policy_action"]
	if ok && internaltypes.IsDefined(authnSelectorPolicyAction) {
		policyActionAttrs := authnSelectorPolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("AUTHN_SELECTOR")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.AuthnSelectorPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	authnSourcePolicyAction, ok := attrs["authn_source_policy_action"]
	if ok && internaltypes.IsDefined(authnSourcePolicyAction) {
		policyActionAttrs := authnSourcePolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("AUTHN_SOURCE")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.AuthnSourcePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	continuePolicyAction, ok := attrs["continue_policy_action"]
	if ok && internaltypes.IsDefined(continuePolicyAction) {
		policyActionAttrs := continuePolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("CONTINUE")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.ContinuePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	donePolicyAction, ok := attrs["done_policy_action"]
	if ok && internaltypes.IsDefined(donePolicyAction) {
		policyActionAttrs := donePolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("DONE")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.DonePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	fragmentPolicyAction, ok := attrs["fragment_policy_action"]
	if ok && internaltypes.IsDefined(fragmentPolicyAction) {
		policyActionAttrs := fragmentPolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("FRAGMENT")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.FragmentPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	localIdentityMappingPolicyAction, ok := attrs["local_identity_mapping_policy_action"]
	if ok && internaltypes.IsDefined(localIdentityMappingPolicyAction) {
		policyActionAttrs := localIdentityMappingPolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("LOCAL_IDENTITY_MAPPING")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.LocalIdentityMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	restartPolicyAction, ok := attrs["restart_policy_action"]
	if ok && internaltypes.IsDefined(restartPolicyAction) {
		policyActionAttrs := restartPolicyAction.(types.Object).Attributes()
		// Set type attribute required by the policy action struct
		policyActionAttrs["type"] = types.StringValue("RESTART")
		err := json.Unmarshal([]byte(internaljson.FromAttributesMap(policyActionAttrs, true)), &result.RestartPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}

	return nil, errors.New("no valid policy action type found when building client struct. Ensure you have specified an action value in your policy tree node")
}
