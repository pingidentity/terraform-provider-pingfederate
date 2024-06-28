package sptargeturlmappings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *spTargetUrlMappingsResourceModel) setNullObjectValues() {
	itemsRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	itemsAttrTypes := map[string]attr.Type{
		"ref":  types.ObjectType{AttrTypes: itemsRefAttrTypes},
		"type": types.StringType,
		"url":  types.StringType,
	}
	itemsElementType := types.ObjectType{AttrTypes: itemsAttrTypes}
	m.Items = types.ListNull(itemsElementType)
}
