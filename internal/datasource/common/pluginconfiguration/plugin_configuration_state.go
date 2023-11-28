package pluginconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

// Creates state values for fields. Returns one value that only includes values specified in the plan, and a second value that includes all fields values
func ToFieldsListValue(fields []client.ConfigField, planFields *types.List, diags *diag.Diagnostics) (types.List, types.List) {
	plannedObjValues := []attr.Value{}
	allObjValues := []attr.Value{}
	planFieldsValues := map[string]*string{}
	// Build up a map of all the values from the plan
	if planFields != nil {
		for _, planField := range planFields.Elements() {
			planFieldObj := planField.(types.Object)
			planFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}

	// Build the list of fields
	for _, field := range fields {
		attrValues := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(field.Name)
		attrValues["value"] = types.StringPointerValue(field.Value)
		attrValues["inherited"] = types.BoolPointerValue(field.Inherited)
		attrValues["encrypted_value"] = types.StringPointerValue(field.EncryptedValue)

		// If this field is in the plan, add it to the list of plan fields
		if planFields != nil {
			planValue, ok := planFieldsValues[field.Name]
			if ok {
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.Value == nil {
					planAttrValues["value"] = types.StringPointerValue(planValue)
					planAttrValues["encrypted_value"] = types.StringPointerValue(planValue)
				} else {
					planAttrValues["value"] = types.StringPointerValue(field.Value)
					planAttrValues["encrypted_value"] = types.StringPointerValue(field.EncryptedValue)
				}
				planAttrValues["inherited"] = types.BoolPointerValue(field.Inherited)
				objVal, newDiags := types.ObjectValue(fieldAttrTypes, planAttrValues)
				diags.Append(newDiags...)
				plannedObjValues = append(plannedObjValues, objVal)
			}
		}

		objVal, newDiags := types.ObjectValue(fieldAttrTypes, attrValues)
		diags.Append(newDiags...)
		allObjValues = append(allObjValues, objVal)
	}

	allListVal, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allObjValues)
	diags.Append(newDiags...)
	plannedListVal, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, plannedObjValues)
	diags.Append(newDiags...)
	return plannedListVal, allListVal
}

func ToRowsListValue(rows []client.ConfigRow, planRows *types.List, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	if planRows == nil || planRows.IsNull() {
		if len(rows) == 0 {
			// If the API returned no rows, treat it as null
			return types.ListNull(types.ObjectType{
				AttrTypes: rowAttrTypes,
			})
		}
		for _, row := range rows {
			attrValues := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(row.DefaultRow)
			_, attrValues["fields"] = ToFieldsListValue(row.Fields, nil, diags)
			rowObjVal, newDiags := types.ObjectValue(rowAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, rowObjVal)
		}
	} else {
		// This is assuming there are never any rows added by the PF API. If there
		// are ever rows added, this will cause a nil pointer exception trying to read
		// index i of planRowsElements.
		planRowsElements := planRows.Elements()
		for i := 0; i < len(rows); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(rows[i].DefaultRow)
			planRow := planRowsElements[i].(types.Object)
			var planRowFields *types.List
			planRowFieldsVal, ok := planRow.Attributes()["fields"]
			if ok {
				listVal := planRowFieldsVal.(types.List)
				planRowFields = &listVal
			}
			attrValues["fields"], _ = ToFieldsListValue(rows[i].Fields, planRowFields, diags)
			rowObjVal, newDiags := types.ObjectValue(rowAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, rowObjVal)
		}
	}
	listVal, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: rowAttrTypes,
	}, objValues)
	diags.Append(newDiags...)
	return listVal
}

func ToTablesListValue(tables []client.ConfigTable, planTables *types.List, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	if planTables == nil || planTables.IsNull() {
		if len(tables) == 0 {
			// If the API returned no tables, treat it as null
			return types.ListNull(types.ObjectType{
				AttrTypes: tableAttrTypes,
			})
		}
		for _, table := range tables {
			attrValues := map[string]attr.Value{}
			attrValues["inherited"] = types.BoolPointerValue(table.Inherited)
			attrValues["name"] = types.StringValue(table.Name)
			attrValues["rows"] = ToRowsListValue(table.Rows, nil, diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, tableObjValue)
		}
	} else {
		// This is assuming there are never any tables added by the PF API. If there
		// are ever tables added, this will cause a nil pointer exception trying to read
		// index i of planTablesElements.
		planTablesElements := planTables.Elements()
		for i := 0; i < len(tables); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["inherited"] = types.BoolPointerValue(tables[i].Inherited)
			attrValues["name"] = types.StringValue(tables[i].Name)
			planTable := planTablesElements[i].(types.Object)
			var planTableRows *types.List
			planTableRowsVal, ok := planTable.Attributes()["rows"]
			if ok {
				listValue := planTableRowsVal.(types.List)
				planTableRows = &listValue
			}
			attrValues["rows"] = ToRowsListValue(tables[i].Rows, planTableRows, diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, tableObjValue)
		}
	}
	listVal, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: tableAttrTypes,
	}, objValues)
	diags.Append(newDiags...)
	return listVal
}

func ToDataSourceState(con context.Context, configuration *client.PluginConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags, respDiags diag.Diagnostics
	var fieldsAttrValue, tablesAttrValue types.List
	// planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	// if ok {
	// 	listVal := planFieldsValue.(types.List)
	// 	planFields = &listVal
	// }
	// planTablesValue, ok := configFromPlan.Attributes()["tables"]
	// if ok {
	// 	listVal := planTablesValue.(types.List)
	// 	planTables = &listVal
	// }

	fieldsAttrValue, respDiags = types.ListValueFrom(con, types.ObjectType{AttrTypes: fieldAttrTypes}, configuration.Fields)
	diags.Append(respDiags...)
	tablesAttrValue, respDiags = types.ListValueFrom(con, types.ObjectType{AttrTypes: tableAttrTypes}, configuration.Tables)
	diags.Append(respDiags...)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}
	configObj, valueFromDiags := types.ObjectValue(configurationAttrTypes, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
}
