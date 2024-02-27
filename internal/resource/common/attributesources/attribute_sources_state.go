package attributesources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func CommonAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	commonAttrSourceAttrType := map[string]attr.Type{}
	commonAttrSourceAttrType["type"] = types.StringType
	if includeDataStoreRefResourceLinkLocation {
		commonAttrSourceAttrType["data_store_ref"] = types.ObjectType{AttrTypes: resourcelink.AttrType()}
	} else {
		commonAttrSourceAttrType["data_store_ref"] = types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()}
	}
	commonAttrSourceAttrType["id"] = types.StringType
	commonAttrSourceAttrType["description"] = types.StringType
	commonAttrSourceAttrType["attribute_contract_fulfillment"] = attributecontractfulfillment.MapType()
	return commonAttrSourceAttrType
}

func CustomAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	customAttrSourceAttrType := CommonAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	customAttrSourceAttrType["filter_fields"] = types.ListType{ElemType: types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"value": types.StringType,
			"name":  types.StringType,
		},
	}}
	return customAttrSourceAttrType
}

func JdbcAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	jdbcAttributeSourceAttrType := CommonAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	jdbcAttributeSourceAttrType["schema"] = types.StringType
	jdbcAttributeSourceAttrType["table"] = types.StringType
	jdbcAttributeSourceAttrType["column_names"] = types.ListType{ElemType: types.StringType}
	jdbcAttributeSourceAttrType["filter"] = types.StringType
	return jdbcAttributeSourceAttrType
}

func LdapAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	ldapAttrSourceAttrType := CommonAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	ldapAttrSourceAttrType["base_dn"] = types.StringType
	ldapAttrSourceAttrType["search_scope"] = types.StringType
	ldapAttrSourceAttrType["search_filter"] = types.StringType
	ldapAttrSourceAttrType["search_attributes"] = types.ListType{ElemType: types.StringType}
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

func AttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	return map[string]attr.Type{
		"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: ElemAttrType(includeDataStoreRefResourceLinkLocation),
			},
		},
	}
}

func ElemAttrType(includeDataStoreRefResourceLinkLocation bool) map[string]attr.Type {
	return map[string]attr.Type{
		"custom_attribute_source": types.ObjectType{
			AttrTypes: CustomAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation),
		},
		"jdbc_attribute_source": types.ObjectType{
			AttrTypes: JdbcAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation),
		},
		"ldap_attribute_source": types.ObjectType{
			AttrTypes: LdapAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation),
		},
	}
}

func ToState(con context.Context, attributeSourcesFromClient []client.AttributeSourceAggregation, includeDataStoreRefResourceLinkLocation bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var customAttrSourceAttrTypes = CustomAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	var jdbcAttrSourceAttrTypes = JdbcAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	var ldapAttrSourceAttrTypes = LdapAttributeSourceAttrType(includeDataStoreRefResourceLinkLocation)
	var valueFromDiags diag.Diagnostics

	// Build attribute_sources value
	attrSourceElements := []attr.Value{}
	// This is assuming there won't be any default attribute sources returned by PF and that they will be returned in the same order
	for _, attrSource := range attributeSourcesFromClient {
		attrSourceValues := map[string]attr.Value{}
		if attrSource.CustomAttributeSource != nil {
			customAttrSourceValues := map[string]attr.Value{}
			customAttrSourceValues["filter_fields"], valueFromDiags = types.ListValueFrom(con, customAttrSourceAttrTypes["filter_fields"].(types.ListType).ElemType, attrSource.CustomAttributeSource.FilterFields)
			diags.Append(valueFromDiags...)

			customAttrSourceValues["type"] = types.StringValue("CUSTOM")
			if includeDataStoreRefResourceLinkLocation {
				customAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.CustomAttributeSource.DataStoreRef)
			} else {
				customDataStoreRefAttrValue := map[string]attr.Value{
					"id": types.StringValue(attrSource.CustomAttributeSource.DataStoreRef.Id),
				}
				customAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValue(resourcelink.AttrTypeNoLocation(), customDataStoreRefAttrValue)
			}
			diags.Append(valueFromDiags...)
			customAttrSourceValues["id"] = types.StringPointerValue(attrSource.CustomAttributeSource.Id)
			customAttrSourceValues["description"] = types.StringPointerValue(attrSource.CustomAttributeSource.Description)
			customAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrType()}, attrSource.CustomAttributeSource.AttributeContractFulfillment)
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
			jdbcAttrSourceValues["column_names"], valueFromDiags = types.ListValueFrom(con, types.StringType, attrSource.JdbcAttributeSource.ColumnNames)
			diags.Append(valueFromDiags...)
			jdbcAttrSourceValues["filter"] = types.StringValue(attrSource.JdbcAttributeSource.Filter)
			jdbcAttrSourceValues["type"] = types.StringValue("JDBC")
			if includeDataStoreRefResourceLinkLocation {
				jdbcAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.JdbcAttributeSource.DataStoreRef)
			} else {
				jdbcDataStoreRefAttrValue := map[string]attr.Value{
					"id": types.StringValue(attrSource.JdbcAttributeSource.DataStoreRef.Id),
				}
				jdbcAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValue(resourcelink.AttrTypeNoLocation(), jdbcDataStoreRefAttrValue)
			}
			diags.Append(valueFromDiags...)
			jdbcAttrSourceValues["id"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Id)
			jdbcAttrSourceValues["description"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Description)
			jdbcAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrType()}, attrSource.JdbcAttributeSource.AttributeContractFulfillment)
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
			ldapAttrSourceValues["search_attributes"], valueFromDiags = types.ListValueFrom(con, types.StringType, attrSource.LdapAttributeSource.SearchAttributes)
			diags.Append(valueFromDiags...)
			if attrSource.LdapAttributeSource.BinaryAttributeSettings == nil {
				ldapAttrSourceValues["binary_attribute_settings"] = types.MapNull(ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType)
			} else {
				ldapAttrSourceValues["binary_attribute_settings"], valueFromDiags = types.MapValueFrom(con, ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType, attrSource.LdapAttributeSource.BinaryAttributeSettings)
				diags.Append(valueFromDiags...)
			}
			ldapAttrSourceValues["member_of_nested_group"] = types.BoolPointerValue(attrSource.LdapAttributeSource.MemberOfNestedGroup)
			ldapAttrSourceValues["type"] = types.StringValue(attrSource.LdapAttributeSource.Type)
			if includeDataStoreRefResourceLinkLocation {
				ldapAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(con, resourcelink.AttrType(), attrSource.LdapAttributeSource.DataStoreRef)
			} else {
				ldapDataStoreRefAttrValue := map[string]attr.Value{
					"id": types.StringValue(attrSource.LdapAttributeSource.DataStoreRef.Id),
				}
				ldapAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValue(resourcelink.AttrTypeNoLocation(), ldapDataStoreRefAttrValue)
			}
			diags.Append(valueFromDiags...)
			ldapAttrSourceValues["id"] = types.StringPointerValue(attrSource.LdapAttributeSource.Id)
			ldapAttrSourceValues["description"] = types.StringPointerValue(attrSource.LdapAttributeSource.Description)
			ldapAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(con, types.ObjectType{AttrTypes: attributecontractfulfillment.AttrType()}, attrSource.LdapAttributeSource.AttributeContractFulfillment)
			diags.Append(valueFromDiags...)
			attrSourceValues["ldap_attribute_source"], valueFromDiags = types.ObjectValue(ldapAttrSourceAttrTypes, ldapAttrSourceValues)
			diags.Append(valueFromDiags...)
		} else {
			attrSourceValues["ldap_attribute_source"] = types.ObjectNull(ldapAttrSourceAttrTypes)
		}
		attrSourceElement, valueFromDiags := types.ObjectValue(ElemAttrType(includeDataStoreRefResourceLinkLocation), attrSourceValues)
		diags.Append(valueFromDiags...)
		attrSourceElements = append(attrSourceElements, attrSourceElement)
	}
	attrToState, valueFromDiags := types.ListValue(types.ObjectType{AttrTypes: ElemAttrType(includeDataStoreRefResourceLinkLocation)}, attrSourceElements)
	diags.Append(valueFromDiags...)
	return attrToState, diags
}
