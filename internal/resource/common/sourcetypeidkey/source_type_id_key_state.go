package sourcetypeidkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
	}
}

func AttrVal(con context.Context, attrVal attr.Value) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(con, AttrTypes(), attrVal)
}
