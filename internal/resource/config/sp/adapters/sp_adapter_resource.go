package spadapters

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
)

var (
	attributeContractAttrObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}
	extendedAttributesDefault, _ = types.SetValue(attributeContractAttrObjectType, nil)
	subjectCoreAttribute, _      = types.ObjectValue(attributeContractAttrObjectType.AttrTypes, map[string]attr.Value{
		"name": types.StringValue("subject"),
	})
	coreAttributesDefault, _ = types.SetValue(attributeContractAttrObjectType, []attr.Value{
		subjectCoreAttribute,
	})
	attributeContractDefault, _ = types.ObjectValue(map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: attributeContractAttrObjectType},
		"extended_attributes": types.SetType{ElemType: attributeContractAttrObjectType},
	}, map[string]attr.Value{
		"core_attributes":     coreAttributesDefault,
		"extended_attributes": extendedAttributesDefault,
	})

	targetApplicationInfoDefault, _ = types.ObjectValue(map[string]attr.Type{
		"application_icon_url": types.StringType,
		"application_name":     types.StringType,
	}, map[string]attr.Value{
		"application_icon_url": types.StringNull(),
		"application_name":     types.StringNull(),
	})
)

func (r *spAdapterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *spAdapterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil || state == nil {
		return
	}

	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)

	resp.Plan.Set(ctx, plan)
}
