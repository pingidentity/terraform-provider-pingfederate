// Copyright © 2025 Ping Identity Corporation

package idptokenprocessors

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	attributeContractAttrObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"masked": types.BoolType,
		},
	}
	extendedAttributesDefault, _ = types.SetValue(attributeContractAttrObjectType, nil)
)
