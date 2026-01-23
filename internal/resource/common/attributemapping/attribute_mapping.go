// Copyright © 2025 Ping Identity Corporation

package attributemapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
)

var (
	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_sources": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesources.AttrTypes(),
			},
		},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria": types.ObjectType{
			AttrTypes: issuancecriteria.AttrTypes(),
		},
	}
)

func AttrTypes() map[string]attr.Type {
	return attributeMappingAttrTypes
}

func ToSchema(required bool) schema.SingleNestedAttribute {
	return toSchemaInternal(required, true)
}

func ToSchemaNoValueDefault(required bool) schema.SingleNestedAttribute {
	return toSchemaInternal(required, false)
}

func toSchemaInternal(required, includeValueDefault bool) schema.SingleNestedAttribute {
	var attributeSources schema.Attribute
	if includeValueDefault {
		attributeSources = attributesources.ToSchema(0, false)
	} else {
		attributeSources = attributesources.ToSchemaNoValueDefault(0, false)
	}
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, true),
			"attribute_sources":              attributeSources,
			"issuance_criteria":              issuancecriteria.ToSchema(),
		},
		Required:    required,
		Optional:    !required,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
}

func ToState(con context.Context, attributeMappingFromClient *configurationapi.AttributeMapping) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if attributeMappingFromClient == nil {
		return types.ObjectNull(attributeMappingAttrTypes), diags
	}

	attributeMappingState := map[string]attr.Value{}

	attributeContractFulfillment, objDiags := attributecontractfulfillment.ToState(con, &attributeMappingFromClient.AttributeContractFulfillment)
	diags = append(diags, objDiags...)

	attributeSources, objDiags := attributesources.ToState(con, attributeMappingFromClient.AttributeSources)
	diags = append(diags, objDiags...)

	issuanceCriteria, objDiags := issuancecriteria.ToState(con, attributeMappingFromClient.IssuanceCriteria)
	diags = append(diags, objDiags...)

	attributeMappingState["attribute_contract_fulfillment"] = attributeContractFulfillment
	attributeMappingState["attribute_sources"] = attributeSources
	attributeMappingState["issuance_criteria"] = issuanceCriteria

	attributeMappingObject, objDiags := types.ObjectValue(attributeMappingAttrTypes, attributeMappingState)
	diags = append(diags, objDiags...)

	return attributeMappingObject, diags
}
