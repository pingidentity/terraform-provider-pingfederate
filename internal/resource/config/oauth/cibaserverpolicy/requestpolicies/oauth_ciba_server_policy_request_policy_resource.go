package oauthcibaserverpolicyrequestpolicies

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	alternativeLoginHintTokenIssuersElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"issuer":   types.StringType,
			"jwks":     types.StringType,
			"jwks_url": types.StringType,
		},
	}
	alternativeLoginHintTokenIssuersDefault, _ = types.SetValue(alternativeLoginHintTokenIssuersElemType, nil)

	attributeElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}
	extendedAttributesDefault, _ = types.SetValue(attributeElemType, nil)

	identityHintSubjectCoreAttribute, _ = types.ObjectValue(attributeElemType.AttrTypes, map[string]attr.Value{
		"name": types.StringValue("IDENTITY_HINT_SUBJECT"),
	})
	coreAttributesDefault, _ = types.SetValue(attributeElemType, []attr.Value{identityHintSubjectCoreAttribute})

	identityHintContractDefault, _ = types.ObjectValue(map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: attributeElemType},
		"extended_attributes": types.SetType{ElemType: attributeElemType},
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

	issuanceCriteriaDefault = types.ObjectValueMust(issuancecriteria.AttrTypes(), map[string]attr.Value{
		"conditional_criteria": types.SetValueMust(issuancecriteria.ConditionalCriteriaElemType(), nil),
		"expression_criteria":  types.SetNull(issuancecriteria.ExpressionCriteriaElemType()),
	})
	attributeSourcesDefault = types.SetValueMust(types.ObjectType{AttrTypes: attributesources.AttrTypes()}, nil)
)

func (r *oauthCibaServerPolicyRequestPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Calculate a default for identity_hint_contract_fulfillment if it is not included in the config
	var plan *oauthCibaServerPolicyRequestPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}

	if plan.IdentityHintContractFulfillment.IsUnknown() {
		subjectSource, diags := types.ObjectValue(sourcetypeidkey.AttrTypes(), map[string]attr.Value{
			"id":   types.StringNull(),
			"type": types.StringValue("REQUEST"),
		})
		resp.Diagnostics.Append(diags...)
		idHintSubjectValue, diags := types.ObjectValue(attributecontractfulfillment.AttrTypes(), map[string]attr.Value{
			"source": subjectSource,
			"value":  types.StringValue("IDENTITY_HINT_SUBJECT"),
		})
		resp.Diagnostics.Append(diags...)
		fulfillmentValues := map[string]attr.Value{
			"IDENTITY_HINT_SUBJECT": idHintSubjectValue,
		}

		if internaltypes.IsDefined(plan.IdentityHintContract) {
			extendedAttrs := plan.IdentityHintContract.Attributes()["extended_attributes"].(types.Set)
			for _, extendedAttr := range extendedAttrs.Elements() {
				attrName := extendedAttr.(types.Object).Attributes()["name"].(types.String).ValueString()
				fulfillmentValue, diags := types.ObjectValue(attributecontractfulfillment.AttrTypes(), map[string]attr.Value{
					"source": subjectSource,
					"value":  types.StringValue(attrName),
				})
				resp.Diagnostics.Append(diags...)
				fulfillmentValues[attrName] = fulfillmentValue
			}
		}

		attributeContractFulfillmentDefault, diags := types.MapValue(
			types.ObjectType{AttrTypes: attributecontractfulfillment.AttrTypes()},
			fulfillmentValues,
		)
		resp.Diagnostics.Append(diags...)

		plan.IdentityHintContractFulfillment, diags = types.ObjectValue(identityHintContractFulfillmentAttrTypes, map[string]attr.Value{
			"attribute_contract_fulfillment": attributeContractFulfillmentDefault,
			"attribute_sources":              attributeSourcesDefault,
			"issuance_criteria":              issuanceCriteriaDefault,
		})
		resp.Diagnostics.Append(diags...)

		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

const maxRetries = 5

// Retry creates for this resource because sometimes PF silently fails to create them
func (r *oauthCibaServerPolicyRequestPolicyResource) exponentialBackOffRetryCreate(ctx context.Context, apiCreateRequest client.ApiCreateCibaServerPolicyRequest, policyId string) (*client.RequestPolicy, *http.Response, error) {
	var responseData *client.RequestPolicy
	var httpResp *http.Response
	var err error
	backOffTime := time.Second

	for i := 0; i < maxRetries; i++ {
		responseData, httpResp, err = r.apiClient.OauthCibaServerPolicyAPI.CreateCibaServerPolicyExecute(apiCreateRequest)
		if err != nil {
			// If PF returned an error, don't retry
			return responseData, httpResp, err
		}

		// If PF returned success, ensure the resource was actually created
		_, readHttpResp, readErr := r.apiClient.OauthCibaServerPolicyAPI.GetCibaServerPolicyById(config.AuthContext(ctx, r.providerConfig), policyId).Execute()
		if readErr == nil || readHttpResp == nil || readHttpResp.StatusCode != 404 {
			// Either the create succeeded or the read returned a non-404 status code
			return responseData, httpResp, err
		}

		backOffTime = backOffTime * 2
	}

	tflog.Info(context.Background(), fmt.Sprintf("Request failed after %d attempts", maxRetries))

	return responseData, httpResp, err
}
