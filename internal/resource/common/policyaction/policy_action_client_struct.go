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
		return nil, errors.New("provided object is Null or Unknown")
	}

	result := client.PolicyActionAggregation{}
	attrs := object.Attributes()
	apcMappingPolicyAction, ok := attrs["apc_mapping_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(apcMappingPolicyAction, false)), &result.ApcMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	authnSelectorPolicyAction, ok := attrs["authn_selector_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(authnSelectorPolicyAction, false)), &result.AuthnSelectorPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	authnSourcePolicyAction, ok := attrs["authn_source_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(authnSourcePolicyAction, false)), &result.AuthnSourcePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	continuePolicyAction, ok := attrs["continue_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(continuePolicyAction, false)), &result.ContinuePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	donePolicyAction, ok := attrs["done_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(donePolicyAction, false)), &result.DonePolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	fragmentPolicyAction, ok := attrs["fragment_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(fragmentPolicyAction, false)), &result.FragmentPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	localIdentityMappingPolicyAction, ok := attrs["local_identity_mapping_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(localIdentityMappingPolicyAction, false)), &result.LocalIdentityMappingPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	restartPolicyAction, ok := attrs["restart_policy_action"]
	if ok {
		err := json.Unmarshal([]byte(internaljson.FromValue(restartPolicyAction, false)), &result.RestartPolicyAction)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}

	return nil, errors.New("no valid policy action type found when building client struct")
}
