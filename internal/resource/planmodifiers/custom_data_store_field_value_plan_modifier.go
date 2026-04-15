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
	return "Preserves configured custom data store field values across replans and falls back to the prior state or an empty string for omitted values."
}

func (m customDataStoreFieldValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m customDataStoreFieldValuePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
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
