package sourcetypeidkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SourceTypeIdKeyAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"type": basetypes.StringType{},
		"id":   basetypes.StringType{},
	}
}

func SourceTypeIdKeyAttrVal(con context.Context, attrVal attr.Value) basetypes.ObjectValue {
	sourceTypeIdKeyAttrValObj, _ := types.ObjectValueFrom(con, SourceTypeIdKeyAttrType(), attrVal)
	return sourceTypeIdKeyAttrValObj
}
