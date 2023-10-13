package id

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GenerateUUIDToState(id basetypes.StringValue) basetypes.StringValue {
	if id.IsNull() || id.IsUnknown() {
		return types.StringValue(uuid.NewString())
	} else {
		return types.StringValue(id.ValueString())
	}
}
