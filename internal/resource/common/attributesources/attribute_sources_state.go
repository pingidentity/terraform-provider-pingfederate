// Copyright © 2026 Ping Identity Corporation

package attributesources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func commonAttributeSourceAttrType(includeIdAttr bool) map[string]attr.Type {
	commonAttrSourceAttrType := map[string]attr.Type{}
	commonAttrSourceAttrType["type"] = types.StringType
	commonAttrSourceAttrType["data_store_ref"] = types.ObjectType{AttrTypes: resourcelink.AttrType()}
	if includeIdAttr {
		commonAttrSourceAttrType["id"] = types.StringType
	}
	commonAttrSourceAttrType["description"] = types.StringType
	commonAttrSourceAttrType["attribute_contract_fulfillment"] = attributecontractfulfillment.MapType()
	return commonAttrSourceAttrType
}

func customAttributeSourceAttrType(includeIdAttr bool) map[string]attr.Type {
	customAttrSourceAttrType := commonAttributeSourceAttrType(includeIdAttr)
	customAttrSourceAttrType["filter_fields"] = types.SetType{ElemType: types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"value": types.StringType,
			"name":  types.StringType,
		},
	}}
	return customAttrSourceAttrType
}

func jdbcAttributeSourceAttrType(includeIdAttr bool) map[string]attr.Type {
	jdbcAttributeSourceAttrType := commonAttributeSourceAttrType(includeIdAttr)
	jdbcAttributeSourceAttrType["schema"] = types.StringType
	jdbcAttributeSourceAttrType["table"] = types.StringType
	jdbcAttributeSourceAttrType["column_names"] = types.SetType{ElemType: types.StringType}
	jdbcAttributeSourceAttrType["filter"] = types.StringType
	return jdbcAttributeSourceAttrType
}

func ldapAttributeSourceAttrType(includeIdAttr bool) map[string]attr.Type {
	ldapAttrSourceAttrType := commonAttributeSourceAttrType(includeIdAttr)
	ldapAttrSourceAttrType["base_dn"] = types.StringType
	ldapAttrSourceAttrType["search_scope"] = types.StringType
	ldapAttrSourceAttrType["search_filter"] = types.StringType
	ldapAttrSourceAttrType["search_attributes"] = types.SetType{ElemType: types.StringType}
	ldapAttrSourceAttrType["binary_attribute_settings"] = types.MapType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"binary_encoding": types.StringType,
			},
		},
	}
	ldapAttrSourceAttrType["member_of_nested_group"] = types.BoolType
	return ldapAttrSourceAttrType
}

func AttrTypes() map[string]attr.Type {
	return attrTypesInternal(true)
}

func AttrTypesNoId() map[string]attr.Type {
	return attrTypesInternal(false)
}

func attrTypesInternal(includeIdAttr bool) map[string]attr.Type {
	return map[string]attr.Type{
		"custom_attribute_source": types.ObjectType{
			AttrTypes: customAttributeSourceAttrType(includeIdAttr),
		},
		"jdbc_attribute_source": types.ObjectType{
			AttrTypes: jdbcAttributeSourceAttrType(includeIdAttr),
		},
		"ldap_attribute_source": types.ObjectType{
			AttrTypes: ldapAttributeSourceAttrType(includeIdAttr),
		},
	}
}

func ToState(con context.Context, attributeSourcesFromClient []client.AttributeSourceAggregation) (basetypes.SetValue, diag.Diagnostics) {
	return toStateInternal(con, attributeSourcesFromClient, true)
}

func ToStateNoId(con context.Context, attributeSourcesFromClient []client.AttributeSourceAggregation) (basetypes.SetValue, diag.Diagnostics) {
	return toStateInternal(con, attributeSourcesFromClient, false)
}

func toStateInternal(con context.Context, attributeSourcesFromClient []client.AttributeSourceAggregation, includeIdAttr bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var customAttrSourceAttrTypes = customAttributeSourceAttrType(includeIdAttr)
	var jdbcAttrSourceAttrTypes = jdbcAttributeSourceAttrType(includeIdAttr)
	var ldapAttrSourceAttrTypes = ldapAttributeSourceAttrType(includeIdAttr)
	var valueFromDiags diag.Diagnostics

	// Build attribute_sources value
	attrSourceElements := []attr.Value{}
	// This is assuming there won't be any default attribute sources returned by PF and that they will be returned in the same order
	for _, attrSource := range attributeSourcesFromClient {
		attrSourceValues := map[string]attr.Value{}
		if attrSource.CustomAttributeSource != nil {
			customAttrSourceValues := map[string]attr.Value{}
			customAttrSourceValues["filter_fields"], valueFromDiags = types.SetValueFrom(con, customAttrSourceAttrTypes["filter_fields"].(types.SetType).ElemType, attrSource.CustomAttributeSource.FilterFields)
			diags.Append(valueFromDiags...)

			customAttrSourceValues["type"] = types.StringValue("CUSTOM")
			customAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.CustomAttributeSource.DataStoreRef)
			diags.Append(valueFromDiags...)
			if includeIdAttr {
				customAttrSourceValues["id"] = types.StringPointerValue(attrSource.CustomAttributeSource.Id)
			}
			customAttrSourceValues["description"] = types.StringPointerValue(attrSource.CustomAttributeSource.Description)
			customAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrTypes()}, attrSource.CustomAttributeSource.AttributeContractFulfillment)
			diags.Append(valueFromDiags...)
			attrSourceValues["custom_attribute_source"], valueFromDiags = types.ObjectValue(customAttrSourceAttrTypes, customAttrSourceValues)
			diags.Append(valueFromDiags...)
		} else {
			attrSourceValues["custom_attribute_source"] = types.ObjectNull(customAttrSourceAttrTypes)
		}
		if attrSource.JdbcAttributeSource != nil {
			jdbcAttrSourceValues := map[string]attr.Value{}
			jdbcAttrSourceValues["schema"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Schema)
			jdbcAttrSourceValues["table"] = types.StringValue(attrSource.JdbcAttributeSource.Table)
			jdbcAttrSourceValues["column_names"], valueFromDiags = types.SetValueFrom(con, types.StringType, attrSource.JdbcAttributeSource.ColumnNames)
			diags.Append(valueFromDiags...)
			jdbcAttrSourceValues["filter"] = types.StringValue(attrSource.JdbcAttributeSource.Filter)
			jdbcAttrSourceValues["type"] = types.StringValue("JDBC")
			jdbcAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.JdbcAttributeSource.DataStoreRef)
			diags.Append(valueFromDiags...)
			if includeIdAttr {
				jdbcAttrSourceValues["id"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Id)
			}
			jdbcAttrSourceValues["description"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Description)
			jdbcAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrTypes()}, attrSource.JdbcAttributeSource.AttributeContractFulfillment)
			diags.Append(valueFromDiags...)
			attrSourceValues["jdbc_attribute_source"], valueFromDiags = types.ObjectValue(jdbcAttrSourceAttrTypes, jdbcAttrSourceValues)
			diags.Append(valueFromDiags...)
		} else {
			attrSourceValues["jdbc_attribute_source"] = types.ObjectNull(jdbcAttrSourceAttrTypes)
		}
		if attrSource.LdapAttributeSource != nil {
			ldapAttrSourceValues := map[string]attr.Value{}
			ldapAttrSourceValues["base_dn"] = types.StringPointerValue(attrSource.LdapAttributeSource.BaseDn)
			ldapAttrSourceValues["search_scope"] = types.StringValue(attrSource.LdapAttributeSource.SearchScope)
			ldapAttrSourceValues["search_filter"] = types.StringValue(attrSource.LdapAttributeSource.SearchFilter)
			ldapAttrSourceValues["search_attributes"], valueFromDiags = types.SetValueFrom(con, types.StringType, attrSource.LdapAttributeSource.SearchAttributes)
			diags.Append(valueFromDiags...)
			if attrSource.LdapAttributeSource.BinaryAttributeSettings == nil || len(*attrSource.LdapAttributeSource.BinaryAttributeSettings) == 0 {
				ldapAttrSourceValues["binary_attribute_settings"] = types.MapNull(ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType)
			} else {
				ldapAttrSourceValues["binary_attribute_settings"], valueFromDiags = types.MapValueFrom(con, ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType, attrSource.LdapAttributeSource.BinaryAttributeSettings)
				diags.Append(valueFromDiags...)
			}
			ldapAttrSourceValues["member_of_nested_group"] = types.BoolPointerValue(attrSource.LdapAttributeSource.MemberOfNestedGroup)
			ldapAttrSourceValues["type"] = types.StringValue(attrSource.LdapAttributeSource.Type)
			ldapAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.LdapAttributeSource.DataStoreRef)
			diags.Append(valueFromDiags...)
			if includeIdAttr {
				ldapAttrSourceValues["id"] = types.StringPointerValue(attrSource.LdapAttributeSource.Id)
			}
			ldapAttrSourceValues["description"] = types.StringPointerValue(attrSource.LdapAttributeSource.Description)
			ldapAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrTypes()}, attrSource.LdapAttributeSource.AttributeContractFulfillment)
			diags.Append(valueFromDiags...)
			attrSourceValues["ldap_attribute_source"], valueFromDiags = types.ObjectValue(ldapAttrSourceAttrTypes, ldapAttrSourceValues)
			diags.Append(valueFromDiags...)
		} else {
			attrSourceValues["ldap_attribute_source"] = types.ObjectNull(ldapAttrSourceAttrTypes)
		}
		attrSourceElement, valueFromDiags := types.ObjectValue(attrTypesInternal(includeIdAttr), attrSourceValues)
		diags.Append(valueFromDiags...)
		attrSourceElements = append(attrSourceElements, attrSourceElement)
	}
	attrToState, valueFromDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypesInternal(includeIdAttr)}, attrSourceElements)
	diags.Append(valueFromDiags...)
	return attrToState, diags
}
