// Copyright © 2025 Ping Identity Corporation

package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tableAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables": types.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}
)

func AttrType() map[string]attr.Type {
	return configurationAttrTypes
}
