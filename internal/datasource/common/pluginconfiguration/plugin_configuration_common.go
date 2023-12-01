package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":            basetypes.StringType{},
		"value":           basetypes.StringType{},
		"encrypted_value": basetypes.StringType{},
		"inherited":       basetypes.BoolType{},
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": basetypes.BoolType{},
	}

	tableAttrTypes = map[string]attr.Type{
		"name":      basetypes.StringType{},
		"rows":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: rowAttrTypes}},
		"inherited": basetypes.BoolType{},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}
)

func AttrType() map[string]attr.Type {
	return configurationAttrTypes
}
