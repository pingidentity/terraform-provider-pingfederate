package identitystoreprovisioners

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (state *identityStoreProvisionerResourceModel) readClientResponseAttributeContracts(response *client.IdentityStoreProvisioner, isImportRead bool) diag.Diagnostics {
	var respDiags diag.Diagnostics
	// attribute_contract
	attributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: attributeContractCoreAttributesAttrTypes}
	attributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: attributeContractExtendedAttributesAttrTypes}
	attributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.ListType{ElemType: attributeContractCoreAttributesElementType},
		"core_attributes_all": types.ListType{ElemType: attributeContractCoreAttributesElementType},
		"extended_attributes": types.ListType{ElemType: attributeContractExtendedAttributesElementType},
	}
	var attributeContractValue types.Object
	if response.AttributeContract == nil {
		attributeContractValue = types.ObjectNull(attributeContractAttrTypes)
	} else {
		var attributeContractCoreAttributesAllValues []attr.Value
		var attributeContractCoreAttributesValues []attr.Value
		coreAttributeNamesInPlan := []string{}
		// Only include core_attributes set in the plan in the state core_attributes value, unless this is an import read.
		// On import reads, just read everything into the core_attributes value.
		if internaltypes.IsDefined(state.AttributeContract) {
			for _, coreAttr := range state.AttributeContract.Attributes()["core_attributes"].(types.List).Elements() {
				coreAttributeNamesInPlan = append(coreAttributeNamesInPlan, coreAttr.(types.Object).Attributes()["name"].(types.String).ValueString())
			}
		}
		for _, attributeContractCoreAttributesResponseValue := range response.AttributeContract.CoreAttributes {
			attributeContractCoreAttributesValue, diags := types.ObjectValue(attributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(attributeContractCoreAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractCoreAttributesAllValues = append(attributeContractCoreAttributesAllValues, attributeContractCoreAttributesValue)
			if isImportRead || slices.Contains(coreAttributeNamesInPlan, attributeContractCoreAttributesResponseValue.Name) {
				attributeContractCoreAttributesValues = append(attributeContractCoreAttributesValues, attributeContractCoreAttributesValue)
			}
		}
		attributeContractCoreAttributesAllValue, diags := types.ListValue(attributeContractCoreAttributesElementType, attributeContractCoreAttributesAllValues)
		respDiags.Append(diags...)
		attributeContractCoreAttributesValue, diags := types.ListValue(attributeContractCoreAttributesElementType, attributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var attributeContractExtendedAttributesValues []attr.Value
		for _, attributeContractExtendedAttributesResponseValue := range response.AttributeContract.ExtendedAttributes {
			attributeContractExtendedAttributesValue, diags := types.ObjectValue(attributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(attributeContractExtendedAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractExtendedAttributesValues = append(attributeContractExtendedAttributesValues, attributeContractExtendedAttributesValue)
		}
		attributeContractExtendedAttributesValue, diags := types.ListValue(attributeContractExtendedAttributesElementType, attributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		attributeContractValue, diags = types.ObjectValue(attributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributesValue,
			"core_attributes_all": attributeContractCoreAttributesAllValue,
			"extended_attributes": attributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
	}

	state.AttributeContract = attributeContractValue
	// group_attribute_contract
	var groupAttributeContractValue types.Object
	if response.GroupAttributeContract == nil {
		groupAttributeContractValue = types.ObjectNull(attributeContractAttrTypes)
	} else {
		var groupAttributeContractCoreAttributesAllValues []attr.Value
		var groupAttributeContractCoreAttributesValues []attr.Value
		coreAttributeNamesInPlan := []string{}
		// Only include core_attributes set in the plan in the state core_attributes value, unless this is an import read.
		// On import reads, just read everything into the core_attributes value.
		if internaltypes.IsDefined(state.GroupAttributeContract) {
			for _, coreAttr := range state.GroupAttributeContract.Attributes()["core_attributes"].(types.List).Elements() {
				coreAttributeNamesInPlan = append(coreAttributeNamesInPlan, coreAttr.(types.Object).Attributes()["name"].(types.String).ValueString())
			}
		}
		for _, groupAttributeContractCoreAttributesResponseValue := range response.GroupAttributeContract.CoreAttributes {
			groupAttributeContractCoreAttributesValue, diags := types.ObjectValue(attributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(groupAttributeContractCoreAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			groupAttributeContractCoreAttributesAllValues = append(groupAttributeContractCoreAttributesAllValues, groupAttributeContractCoreAttributesValue)
			if isImportRead || slices.Contains(coreAttributeNamesInPlan, groupAttributeContractCoreAttributesResponseValue.Name) {
				groupAttributeContractCoreAttributesValues = append(groupAttributeContractCoreAttributesValues, groupAttributeContractCoreAttributesValue)
			}
		}
		groupAttributeContractCoreAttributesAllValue, diags := types.ListValue(attributeContractCoreAttributesElementType, groupAttributeContractCoreAttributesAllValues)
		respDiags.Append(diags...)
		groupAttributeContractCoreAttributesValue, diags := types.ListValue(attributeContractCoreAttributesElementType, groupAttributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var groupAttributeContractExtendedAttributesValues []attr.Value
		for _, groupAttributeContractExtendedAttributesResponseValue := range response.GroupAttributeContract.ExtendedAttributes {
			groupAttributeContractExtendedAttributesValue, diags := types.ObjectValue(attributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(groupAttributeContractExtendedAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			groupAttributeContractExtendedAttributesValues = append(groupAttributeContractExtendedAttributesValues, groupAttributeContractExtendedAttributesValue)
		}
		groupAttributeContractExtendedAttributesValue, diags := types.ListValue(attributeContractExtendedAttributesElementType, groupAttributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		groupAttributeContractValue, diags = types.ObjectValue(attributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     groupAttributeContractCoreAttributesValue,
			"core_attributes_all": groupAttributeContractCoreAttributesAllValue,
			"extended_attributes": groupAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
	}

	state.GroupAttributeContract = groupAttributeContractValue
	return respDiags
}
