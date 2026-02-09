// Copyright © 2026 Ping Identity Corporation

package id

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

type idModelStruct struct {
	Id types.String `tfsdk:"id"`
}

func GetID(con context.Context, state tfsdk.State) (*string, diag.Diagnostics) {
	idModel := idModelStruct{}
	diags := state.GetAttribute(con, path.Root("id"), &idModel.Id)
	if !internaltypes.IsDefined(idModel.Id) {
		return nil, diags
	} else {
		return idModel.Id.ValueStringPointer(), diags
	}
}

func GenerateUUIDToState(id *string) basetypes.StringValue {
	if id == nil {
		return types.StringValue(uuid.NewString())
	} else {
		return types.StringPointerValue(id)
	}
}
