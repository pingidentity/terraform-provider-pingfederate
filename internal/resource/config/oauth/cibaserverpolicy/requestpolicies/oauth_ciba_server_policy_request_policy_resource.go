package oauthcibaserverpolicyrequestpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
)

var (
	alternativeLoginHintTokenIssuersElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"issuer":   types.StringType,
			"jwks":     types.StringType,
			"jwks_url": types.StringType,
		},
	}
	alternativeLoginHintTokenIssuersDefault, _ = types.ListValue(alternativeLoginHintTokenIssuersElemType, nil)

	attributeElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}
	extendedAttributesDefault, _ = types.ListValue(attributeElemType, nil)

	identityHintSubjectCoreAttribute, _ = types.ObjectValue(attributeElemType.AttrTypes, map[string]attr.Value{
		"name": types.StringValue("IDENTITY_HINT_SUBJECT"),
	})
	coreAttributesDefault, _ = types.ListValue(attributeElemType, []attr.Value{identityHintSubjectCoreAttribute})

	identityHintContractDefault, _ = types.ObjectValue(map[string]attr.Type{
		"core_attributes":     types.ListType{ElemType: attributeElemType},
		"extended_attributes": types.ListType{ElemType: attributeElemType},
	}, map[string]attr.Value{
		"core_attributes":     coreAttributesDefault,
		"extended_attributes": extendedAttributesDefault,
	})

	identityHintContractFulfillmentAttrTypes = map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributecontractfulfillment.AttrTypes(),
			},
		},
		"attribute_sources": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesources.AttrTypes(),
			},
		},
		"issuance_criteria": types.ObjectType{
			AttrTypes: issuancecriteria.AttrTypes(),
		},
	}
)

func (r *oauthCibaServerPolicyRequestPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the identity_hint_contract has changed, invalidate the identity_hint_contract_fulfillment value if it is not included in the config,
	// since it is computed and uses UseStateForUnknown.
	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() && !req.Plan.Raw.Equal(req.State.Raw) {
		var plan, config, state oauthCibaServerPolicyRequestPolicyResourceModel
		req.Config.Get(ctx, &config)
		// Ensure the attribute was not explicitly specified in the config by the user
		if config.IdentityHintContractFulfillment.IsNull() {
			req.Plan.Get(ctx, &plan)
			req.State.Get(ctx, &state)
			// If the identity_hint_contract has changed, then invalidate the computed contract_fulfillment
			if !plan.IdentityHintContract.Equal(state.IdentityHintContract) {
				plan.IdentityHintContractFulfillment = types.ObjectUnknown(identityHintContractFulfillmentAttrTypes)
				resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
			}
		}
	}
}
