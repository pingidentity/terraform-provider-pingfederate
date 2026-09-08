// Copyright © 2026 Ping Identity Corporation

package pluginconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
)

// toDataSourceFieldAttrValue explicitly constructs a single field object, skipping fields
// not in the Terraform schema (e.g. inherited)
func toDataSourceFieldAttrValue(field client.ConfigField) basetypes.ObjectValue {
	return types.ObjectValueMust(fieldAttrTypes, map[string]attr.Value{
		"name":            types.StringValue(field.Name),
		"value":           types.StringPointerValue(field.Value),
		"encrypted_value": types.StringPointerValue(field.EncryptedValue),
	})
}

// toDataSourceRowAttrValue explicitly constructs a single row
func toDataSourceRowAttrValue(row client.ConfigRow) basetypes.ObjectValue {
	fields := []attr.Value{}
	for _, field := range row.Fields {
		fields = append(fields, toDataSourceFieldAttrValue(field))
	}
	return types.ObjectValueMust(rowAttrTypes, map[string]attr.Value{
		"fields":      types.SetValueMust(types.ObjectType{AttrTypes: fieldAttrTypes}, fields),
		"default_row": types.BoolPointerValue(row.DefaultRow),
	})
}

// toDataSourceTableAttrValue explicitly constructs a single table, skipping the inherited
// field on the client ConfigTable which is not in the Terraform schema
func toDataSourceTableAttrValue(table client.ConfigTable) basetypes.ObjectValue {
	rows := []attr.Value{}
	for _, row := range table.Rows {
		rows = append(rows, toDataSourceRowAttrValue(row))
	}
	return types.ObjectValueMust(tableAttrTypes, map[string]attr.Value{
		"name": types.StringValue(table.Name),
		"rows": types.ListValueMust(types.ObjectType{AttrTypes: rowAttrTypes}, rows),
	})
}

func ToDataSourceState(con context.Context, configuration *client.PluginConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	fieldsAttrValue := []attr.Value{}
	for _, field := range configuration.Fields {
		fieldsAttrValue = append(fieldsAttrValue, toDataSourceFieldAttrValue(field))
	}
	tablesAttrValue := []attr.Value{}
	for _, table := range configuration.Tables {
		tablesAttrValue = append(tablesAttrValue, toDataSourceTableAttrValue(table))
	}

	configurationAttrValue := map[string]attr.Value{
		"fields": types.SetValueMust(types.ObjectType{AttrTypes: fieldAttrTypes}, fieldsAttrValue),
		"tables": types.ListValueMust(types.ObjectType{AttrTypes: tableAttrTypes}, tablesAttrValue),
	}

	configObj, valueFromDiags := types.ObjectValue(configurationAttrTypes, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
}
