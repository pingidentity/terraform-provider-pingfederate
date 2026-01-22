// Copyright © 2025 Ping Identity Corporation

package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
)

func ClientStruct(configurationObj types.Object) *client.PluginConfiguration {
	configurationValue := client.PluginConfiguration{}
	configurationAttrs := configurationObj.Attributes()
	configurationValue.Fields = fieldsFromObject(configurationAttrs["fields"].(types.Set), configurationAttrs["sensitive_fields"].(types.Set))
	configurationValue.Tables = []client.ConfigTable{}
	for _, tablesElement := range configurationAttrs["tables"].(types.List).Elements() {
		tablesValue := client.ConfigTable{}
		tablesAttrs := tablesElement.(types.Object).Attributes()
		tablesValue.Name = tablesAttrs["name"].(types.String).ValueString()
		tablesValue.Rows = []client.ConfigRow{}
		for _, rowsElement := range tablesAttrs["rows"].(types.List).Elements() {
			rowsValue := client.ConfigRow{}
			rowsAttrs := rowsElement.(types.Object).Attributes()
			rowsValue.DefaultRow = rowsAttrs["default_row"].(types.Bool).ValueBoolPointer()
			rowsValue.Fields = fieldsFromObject(rowsAttrs["fields"].(types.Set), rowsAttrs["sensitive_fields"].(types.Set))
			tablesValue.Rows = append(tablesValue.Rows, rowsValue)
		}
		configurationValue.Tables = append(configurationValue.Tables, tablesValue)
	}
	return &configurationValue
}

func fieldsFromObject(fieldsObj types.Set, sensitiveFieldsObj types.Set) []client.ConfigField {
	fields := []client.ConfigField{}
	for _, fieldsElement := range fieldsObj.Elements() {
		fieldsValue := client.ConfigField{}
		fieldsAttrs := fieldsElement.(types.Object).Attributes()
		fieldsValue.Name = fieldsAttrs["name"].(types.String).ValueString()
		fieldsValue.Value = fieldsAttrs["value"].(types.String).ValueStringPointer()
		fields = append(fields, fieldsValue)
	}
	for _, fieldsElement := range sensitiveFieldsObj.Elements() {
		fieldsValue := client.ConfigField{}
		fieldsAttrs := fieldsElement.(types.Object).Attributes()
		fieldsValue.Name = fieldsAttrs["name"].(types.String).ValueString()
		fieldsValue.Value = fieldsAttrs["value"].(types.String).ValueStringPointer()
		fieldsValue.EncryptedValue = fieldsAttrs["encrypted_value"].(types.String).ValueStringPointer()
		fields = append(fields, fieldsValue)
	}
	return fields
}
