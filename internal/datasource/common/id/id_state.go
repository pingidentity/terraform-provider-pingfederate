package id

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GenerateUUIDStringPointer() *string {
	uuid := uuid.NewString()
	return &uuid
}

func GenerateUUIDToState(id *string) basetypes.StringValue {
	if id == nil {
		return types.StringValue(uuid.NewString())
	} else {
		return types.StringPointerValue(id)
	}
}
