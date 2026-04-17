// Copyright © 2026 Ping Identity Corporation

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = &customDataStoreFieldValuePlanModifier{}

type customDataStoreFieldValuePlanModifier struct{}

func (m customDataStoreFieldValuePlanModifier) Description(ctx context.Context) string {
	return "Preserves configured values, keeps prior state only for unknown replans, and otherwise falls back to an empty string for omitted custom data store field values."
}

func (m customDataStoreFieldValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m customDataStoreFieldValuePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsUnknown() {
		if !req.ConfigValue.IsNull() {
			resp.PlanValue = req.ConfigValue
			return
		}

		resp.PlanValue = types.StringValue("")
		return
	}

	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	resp.PlanValue = types.StringValue("")
}

func CustomDataStoreFieldValue() planmodifier.String {
	return customDataStoreFieldValuePlanModifier{}
}
