// Copyright © 2025 Ping Identity Corporation

package attributesources

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(attributeSourcesAttr basetypes.ListValue) []client.AttributeSourceAggregation {
	attributeSourceAggregation := []client.AttributeSourceAggregation{}
	for _, source := range attributeSourcesAttr.Elements() {
		//Determine which attribute source type this is
		sourceAttrs := source.(types.Object).Attributes()
		attributeSourceInner := client.AttributeSourceAggregation{}
		if internaltypes.IsDefined(sourceAttrs["custom_attribute_source"]) {
			attributeSourceInner.CustomAttributeSource = &client.CustomAttributeSource{}
			customAttributeSourceAttrs := sourceAttrs["custom_attribute_source"].(types.Object).Attributes()
			if !customAttributeSourceAttrs["attribute_contract_fulfillment"].IsNull() && !customAttributeSourceAttrs["attribute_contract_fulfillment"].IsUnknown() {
				attributeSourceInner.CustomAttributeSource.AttributeContractFulfillment = &map[string]client.AttributeFulfillmentValue{}
				for key, attributeContractFulfillmentElement := range customAttributeSourceAttrs["attribute_contract_fulfillment"].(types.Map).Elements() {
					attributeContractFulfillmentValue := client.AttributeFulfillmentValue{}
					attributeContractFulfillmentAttrs := attributeContractFulfillmentElement.(types.Object).Attributes()
					attributeContractFulfillmentSourceValue := client.SourceTypeIdKey{}
					attributeContractFulfillmentSourceAttrs := attributeContractFulfillmentAttrs["source"].(types.Object).Attributes()
					attributeContractFulfillmentSourceValue.Id = attributeContractFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
					attributeContractFulfillmentSourceValue.Type = attributeContractFulfillmentSourceAttrs["type"].(types.String).ValueString()
					attributeContractFulfillmentValue.Source = attributeContractFulfillmentSourceValue
					attributeContractFulfillmentValue.Value = attributeContractFulfillmentAttrs["value"].(types.String).ValueString()
					(*attributeSourceInner.CustomAttributeSource.AttributeContractFulfillment)[key] = attributeContractFulfillmentValue
				}
			}
			customAttributeSourceDataStoreRefValue := client.ResourceLink{}
			customAttributeSourceDataStoreRefAttrs := customAttributeSourceAttrs["data_store_ref"].(types.Object).Attributes()
			customAttributeSourceDataStoreRefValue.Id = customAttributeSourceDataStoreRefAttrs["id"].(types.String).ValueString()
			attributeSourceInner.CustomAttributeSource.DataStoreRef = customAttributeSourceDataStoreRefValue
			attributeSourceInner.CustomAttributeSource.Description = customAttributeSourceAttrs["description"].(types.String).ValueStringPointer()
			if !customAttributeSourceAttrs["filter_fields"].IsNull() && !customAttributeSourceAttrs["filter_fields"].IsUnknown() {
				attributeSourceInner.CustomAttributeSource.FilterFields = []client.FieldEntry{}
				for _, filterFieldsElement := range customAttributeSourceAttrs["filter_fields"].(types.Set).Elements() {
					filterFieldsValue := client.FieldEntry{}
					filterFieldsAttrs := filterFieldsElement.(types.Object).Attributes()
					filterFieldsValue.Name = filterFieldsAttrs["name"].(types.String).ValueString()
					filterFieldsValue.Value = filterFieldsAttrs["value"].(types.String).ValueStringPointer()
					attributeSourceInner.CustomAttributeSource.FilterFields = append(attributeSourceInner.CustomAttributeSource.FilterFields, filterFieldsValue)
				}
			}
			idAttr, ok := customAttributeSourceAttrs["id"]
			if ok {
				attributeSourceInner.CustomAttributeSource.Id = idAttr.(types.String).ValueStringPointer()
			}
			attributeSourceInner.CustomAttributeSource.Type = customAttributeSourceAttrs["type"].(types.String).ValueString()
		}
		if internaltypes.IsDefined(sourceAttrs["jdbc_attribute_source"]) {
			attributeSourceInner.JdbcAttributeSource = &client.JdbcAttributeSource{}
			jdbcAttributeSourceAttrs := sourceAttrs["jdbc_attribute_source"].(types.Object).Attributes()
			if !jdbcAttributeSourceAttrs["attribute_contract_fulfillment"].IsNull() && !jdbcAttributeSourceAttrs["attribute_contract_fulfillment"].IsUnknown() {
				attributeSourceInner.JdbcAttributeSource.AttributeContractFulfillment = &map[string]client.AttributeFulfillmentValue{}
				for key, attributeContractFulfillmentElement := range jdbcAttributeSourceAttrs["attribute_contract_fulfillment"].(types.Map).Elements() {
					attributeContractFulfillmentValue := client.AttributeFulfillmentValue{}
					attributeContractFulfillmentAttrs := attributeContractFulfillmentElement.(types.Object).Attributes()
					attributeContractFulfillmentSourceValue := client.SourceTypeIdKey{}
					attributeContractFulfillmentSourceAttrs := attributeContractFulfillmentAttrs["source"].(types.Object).Attributes()
					attributeContractFulfillmentSourceValue.Id = attributeContractFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
					attributeContractFulfillmentSourceValue.Type = attributeContractFulfillmentSourceAttrs["type"].(types.String).ValueString()
					attributeContractFulfillmentValue.Source = attributeContractFulfillmentSourceValue
					attributeContractFulfillmentValue.Value = attributeContractFulfillmentAttrs["value"].(types.String).ValueString()
					(*attributeSourceInner.JdbcAttributeSource.AttributeContractFulfillment)[key] = attributeContractFulfillmentValue
				}
			}
			if !jdbcAttributeSourceAttrs["column_names"].IsNull() && !jdbcAttributeSourceAttrs["column_names"].IsUnknown() {
				attributeSourceInner.JdbcAttributeSource.ColumnNames = []string{}
				for _, columnNamesElement := range jdbcAttributeSourceAttrs["column_names"].(types.Set).Elements() {
					attributeSourceInner.JdbcAttributeSource.ColumnNames = append(attributeSourceInner.JdbcAttributeSource.ColumnNames, columnNamesElement.(types.String).ValueString())
				}
			}
			jdbcAttributeSourceDataStoreRefValue := client.ResourceLink{}
			jdbcAttributeSourceDataStoreRefAttrs := jdbcAttributeSourceAttrs["data_store_ref"].(types.Object).Attributes()
			jdbcAttributeSourceDataStoreRefValue.Id = jdbcAttributeSourceDataStoreRefAttrs["id"].(types.String).ValueString()
			attributeSourceInner.JdbcAttributeSource.DataStoreRef = jdbcAttributeSourceDataStoreRefValue
			attributeSourceInner.JdbcAttributeSource.Description = jdbcAttributeSourceAttrs["description"].(types.String).ValueStringPointer()
			attributeSourceInner.JdbcAttributeSource.Filter = jdbcAttributeSourceAttrs["filter"].(types.String).ValueString()
			idAttr, ok := jdbcAttributeSourceAttrs["id"]
			if ok {
				attributeSourceInner.JdbcAttributeSource.Id = idAttr.(types.String).ValueStringPointer()
			}
			attributeSourceInner.JdbcAttributeSource.Schema = jdbcAttributeSourceAttrs["schema"].(types.String).ValueStringPointer()
			attributeSourceInner.JdbcAttributeSource.Table = jdbcAttributeSourceAttrs["table"].(types.String).ValueString()
			attributeSourceInner.JdbcAttributeSource.Type = jdbcAttributeSourceAttrs["type"].(types.String).ValueString()
		}
		if internaltypes.IsDefined(sourceAttrs["ldap_attribute_source"]) {
			attributeSourceInner.LdapAttributeSource = &client.LdapAttributeSource{}
			ldapAttributeSourceAttrs := sourceAttrs["ldap_attribute_source"].(types.Object).Attributes()
			if !ldapAttributeSourceAttrs["attribute_contract_fulfillment"].IsNull() && !ldapAttributeSourceAttrs["attribute_contract_fulfillment"].IsUnknown() {
				attributeSourceInner.LdapAttributeSource.AttributeContractFulfillment = &map[string]client.AttributeFulfillmentValue{}
				for key, attributeContractFulfillmentElement := range ldapAttributeSourceAttrs["attribute_contract_fulfillment"].(types.Map).Elements() {
					attributeContractFulfillmentValue := client.AttributeFulfillmentValue{}
					attributeContractFulfillmentAttrs := attributeContractFulfillmentElement.(types.Object).Attributes()
					attributeContractFulfillmentSourceValue := client.SourceTypeIdKey{}
					attributeContractFulfillmentSourceAttrs := attributeContractFulfillmentAttrs["source"].(types.Object).Attributes()
					attributeContractFulfillmentSourceValue.Id = attributeContractFulfillmentSourceAttrs["id"].(types.String).ValueStringPointer()
					attributeContractFulfillmentSourceValue.Type = attributeContractFulfillmentSourceAttrs["type"].(types.String).ValueString()
					attributeContractFulfillmentValue.Source = attributeContractFulfillmentSourceValue
					attributeContractFulfillmentValue.Value = attributeContractFulfillmentAttrs["value"].(types.String).ValueString()
					(*attributeSourceInner.LdapAttributeSource.AttributeContractFulfillment)[key] = attributeContractFulfillmentValue
				}
			}
			attributeSourceInner.LdapAttributeSource.BaseDn = ldapAttributeSourceAttrs["base_dn"].(types.String).ValueStringPointer()
			if !ldapAttributeSourceAttrs["binary_attribute_settings"].IsNull() && !ldapAttributeSourceAttrs["binary_attribute_settings"].IsUnknown() {
				attributeSourceInner.LdapAttributeSource.BinaryAttributeSettings = &map[string]client.BinaryLdapAttributeSettings{}
				for key, binaryAttributeSettingsElement := range ldapAttributeSourceAttrs["binary_attribute_settings"].(types.Map).Elements() {
					binaryAttributeSettingsValue := client.BinaryLdapAttributeSettings{}
					binaryAttributeSettingsAttrs := binaryAttributeSettingsElement.(types.Object).Attributes()
					binaryAttributeSettingsValue.BinaryEncoding = binaryAttributeSettingsAttrs["binary_encoding"].(types.String).ValueStringPointer()
					(*attributeSourceInner.LdapAttributeSource.BinaryAttributeSettings)[key] = binaryAttributeSettingsValue
				}
			}
			ldapAttributeSourceDataStoreRefValue := client.ResourceLink{}
			ldapAttributeSourceDataStoreRefAttrs := ldapAttributeSourceAttrs["data_store_ref"].(types.Object).Attributes()
			ldapAttributeSourceDataStoreRefValue.Id = ldapAttributeSourceDataStoreRefAttrs["id"].(types.String).ValueString()
			attributeSourceInner.LdapAttributeSource.DataStoreRef = ldapAttributeSourceDataStoreRefValue
			attributeSourceInner.LdapAttributeSource.Description = ldapAttributeSourceAttrs["description"].(types.String).ValueStringPointer()
			idAttr, ok := ldapAttributeSourceAttrs["id"]
			if ok {
				attributeSourceInner.LdapAttributeSource.Id = idAttr.(types.String).ValueStringPointer()
			}
			attributeSourceInner.LdapAttributeSource.MemberOfNestedGroup = ldapAttributeSourceAttrs["member_of_nested_group"].(types.Bool).ValueBoolPointer()
			if !ldapAttributeSourceAttrs["search_attributes"].IsNull() && !ldapAttributeSourceAttrs["search_attributes"].IsUnknown() {
				attributeSourceInner.LdapAttributeSource.SearchAttributes = []string{}
				for _, searchAttributesElement := range ldapAttributeSourceAttrs["search_attributes"].(types.Set).Elements() {
					attributeSourceInner.LdapAttributeSource.SearchAttributes = append(attributeSourceInner.LdapAttributeSource.SearchAttributes, searchAttributesElement.(types.String).ValueString())
				}
			}
			attributeSourceInner.LdapAttributeSource.SearchFilter = ldapAttributeSourceAttrs["search_filter"].(types.String).ValueString()
			attributeSourceInner.LdapAttributeSource.SearchScope = ldapAttributeSourceAttrs["search_scope"].(types.String).ValueString()
			attributeSourceInner.LdapAttributeSource.Type = ldapAttributeSourceAttrs["type"].(types.String).ValueString()
		}
		attributeSourceAggregation = append(attributeSourceAggregation, attributeSourceInner)
	}
	return attributeSourceAggregation
}
