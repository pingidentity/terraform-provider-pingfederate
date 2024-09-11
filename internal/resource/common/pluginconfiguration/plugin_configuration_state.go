package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row":      types.BoolType,
	}
	rowsAllAttrTypes = map[string]attr.Type{
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tableAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowAttrTypes}},
	}
	tablesAllAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsAllAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all":       types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables":           types.SetType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
		"tables_all":       types.SetType{ElemType: types.ObjectType{AttrTypes: tablesAllAttrTypes}},
	}
)

func AttrTypes() map[string]attr.Type {
	return configurationAttrTypes
}

// Creates state values for fields. Returns one value that only includes values specified in the plan, and a second value that includes all fields values
func toFieldsSetValue(fields []client.ConfigField, planFields, planSensitiveFields *types.Set, isImportRead bool, diags *diag.Diagnostics) (types.Set, types.Set, types.Set) {
	plannedObjValues := []attr.Value{}
	plannedSensitiveObjValues := []attr.Value{}
	allNonSensitiveObjValues := []attr.Value{}
	allSensitiveObjValues := []attr.Value{}
	allObjValues := []attr.Value{}
	planFieldsValues := map[string]*string{}
	planSensitiveFieldsValues := map[string]*string{}
	// Build up a map of all the values from the plan
	if planFields != nil {
		for _, planField := range planFields.Elements() {
			planFieldObj := planField.(types.Object)
			planFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}
	//TODO some more logic for warnings for when a non-sensitive field doesn't get returned by the PF API
	if planSensitiveFields != nil {
		for _, planField := range planSensitiveFields.Elements() {
			planFieldObj := planField.(types.Object)
			planSensitiveFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}

	for _, field := range fields {
		attrValues := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(field.Name)
		attrValues["value"] = types.StringPointerValue(field.Value)

		// If this field is in the plan, add it to the list of plan fields
		//TODO validation that you don't put the same field in both sets
		fieldAdded := false
		if planFields != nil {
			planValue, ok := planFieldsValues[field.Name]
			if ok {
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.Value == nil {
					//TODO warning
					planAttrValues["value"] = types.StringPointerValue(planValue)
				} else {
					planAttrValues["value"] = types.StringPointerValue(field.Value)
				}
				objVal, newDiags := types.ObjectValue(fieldAttrTypes, planAttrValues)
				diags.Append(newDiags...)
				plannedObjValues = append(plannedObjValues, objVal)
				fieldAdded = true
			}
		}
		if planSensitiveFields != nil && !fieldAdded {
			planValue, ok := planSensitiveFieldsValues[field.Name]
			if ok {
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.Value == nil {
					planAttrValues["value"] = types.StringPointerValue(planValue)
				} else {
					//TODO warning
					planAttrValues["value"] = types.StringPointerValue(field.Value)
				}
				objVal, newDiags := types.ObjectValue(fieldAttrTypes, planAttrValues)
				diags.Append(newDiags...)
				plannedSensitiveObjValues = append(plannedSensitiveObjValues, objVal)
			}
		}

		objVal, newDiags := types.ObjectValue(fieldAttrTypes, attrValues)
		diags.Append(newDiags...)
		allObjValues = append(allObjValues, objVal)
		if field.EncryptedValue != nil && *field.EncryptedValue != "" {
			allSensitiveObjValues = append(allSensitiveObjValues, objVal)
		} else {
			allNonSensitiveObjValues = append(allNonSensitiveObjValues, objVal)
		}
	}

	allSetVal, newDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allObjValues)
	diags.Append(newDiags...)
	var plannedSetVal, plannedSensitiveSetVal types.Set
	if isImportRead {
		// On imports, just read everything directly into the "fields" and "sensitive_fields" attributes,
		// even though there is no plan
		plannedSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, allNonSensitiveObjValues)
		diags.Append(newDiags...)
		plannedSensitiveSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, allSensitiveObjValues)
		diags.Append(newDiags...)
	} else {
		plannedSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, plannedObjValues)
		diags.Append(newDiags...)
		plannedSensitiveSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, plannedSensitiveObjValues)
		diags.Append(newDiags...)
	}
	return plannedSetVal, plannedSensitiveSetVal, allSetVal
}

func toRowsListValue(rows []client.ConfigRow, planRows *types.List, isImportRead, splitSensitiveFields bool, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	objValuesWithSensitive := []attr.Value{}
	if planRows == nil || planRows.IsNull() {
		if len(rows) == 0 {
			// If the API returned no rows, treat it as null
			return types.ListNull(types.ObjectType{
				AttrTypes: rowAttrTypes,
			})
		}
		for _, row := range rows {
			attrValues := map[string]attr.Value{}
			attrValuesWithSensitive := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(row.DefaultRow)
			attrValuesWithSensitive["default_row"] = types.BoolPointerValue(row.DefaultRow)
			attrValuesWithSensitive["fields"], attrValuesWithSensitive["sensitive_fields"], attrValues["fields"] =
				toFieldsSetValue(row.Fields, nil, nil, isImportRead, diags)
			rowObjVal, newDiags := types.ObjectValue(rowAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, rowObjVal)
			rowObjValWithSensitive, newDiags := types.ObjectValue(rowAttrTypes, attrValuesWithSensitive)
			diags.Append(newDiags...)
			objValuesWithSensitive = append(objValuesWithSensitive, rowObjValWithSensitive)
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
			var planRowFields, planRowSensitiveFields *types.Set
			planRowFieldsVal, ok := planRow.Attributes()["fields"]
			if ok {
				setVal := planRowFieldsVal.(types.Set)
				planRowFields = &setVal
			}
			planRowSensitiveFieldsVal, ok := planRow.Attributes()["sensitive_fields"]
			if ok {
				setVal := planRowSensitiveFieldsVal.(types.Set)
				planRowSensitiveFields = &setVal
			}
			attrValues["fields"], attrValues["sensitive_fields"], _ = toFieldsSetValue(rows[i].Fields, planRowFields, planRowSensitiveFields, isImportRead, diags)
			rowObjVal, newDiags := types.ObjectValue(rowAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, rowObjVal)
			objValuesWithSensitive = append(objValuesWithSensitive, rowObjVal)
		}
	}
	//TODO this is all kinds of wrong I think... Need to think through the possible states coming into this method for rows
	var listVal types.List
	var newDiags diag.Diagnostics
	if splitSensitiveFields {
		listVal, newDiags = types.ListValue(types.ObjectType{
			AttrTypes: rowAttrTypes,
		}, objValuesWithSensitive)
		diags.Append(newDiags...)
	} else {
		listVal, newDiags = types.ListValue(types.ObjectType{
			AttrTypes: rowAttrTypes,
		}, objValues)
		diags.Append(newDiags...)
	}
	return listVal
}

// Creates state values for tables. Returns one value that only includes values specified in the plan, and a second value that includes all tables values
func toTablesSetValue(tables []client.ConfigTable, planTables *types.Set, isImportRead bool, diags *diag.Diagnostics) (types.Set, types.Set) {
	// List of *all* tables values to return
	finalTablesAllObjValues := []attr.Value{}
	// List of tables values to return that were expected based on the plan
	finalTablesObjValues := []attr.Value{}
	// types.Object values for tables included in the plan
	planTableObjs := map[string]types.Object{}
	if planTables == nil {
		for _, table := range tables {
			attrValues := map[string]attr.Value{}
			attrValues["name"] = types.StringValue(table.Name)
			attrValues["rows"] = toRowsListValue(table.Rows, nil, isImportRead, diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			finalTablesAllObjValues = append(finalTablesAllObjValues, tableObjValue)
		}
	} else {
		// Build up a map of all the tables included in the plan
		for _, planTable := range planTables.Elements() {
			planTableObj := planTable.(types.Object)
			planTableObjs[planTableObj.Attributes()["name"].(types.String).ValueString()] = planTableObj
		}

		for i := 0; i < len(tables); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["name"] = types.StringValue(tables[i].Name)
			// If this table was in the plan, pass in the planned rows when getting the 'rows' values in case there are some encrypted values
			// that aren't returned by the PF API
			var planTableRows *types.List
			planTable, inPlan := planTableObjs[tables[i].Name]
			if inPlan {
				planTableRowsVal, ok := planTable.Attributes()["rows"]
				if ok {
					listValue := planTableRowsVal.(types.List)
					planTableRows = &listValue
				}
			}
			attrValues["rows"] = toRowsListValue(tables[i].Rows, planTableRows, isImportRead, diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			finalTablesAllObjValues = append(finalTablesAllObjValues, tableObjValue)
			if inPlan {
				finalTablesObjValues = append(finalTablesObjValues, tableObjValue)
			}
		}
	}
	allTables, newDiags := types.SetValue(types.ObjectType{
		AttrTypes: tablesAllAttrTypes,
	}, finalTablesAllObjValues)
	diags.Append(newDiags...)
	var plannedTables types.Set
	if isImportRead {
		// On imports, just read everything directly into the "tables" attribute
		plannedTables, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: tableAttrTypes,
		}, finalTablesAllObjValues)
		diags.Append(newDiags...)
	} else {
		plannedTables, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: tableAttrTypes,
		}, finalTablesObjValues)
		diags.Append(newDiags...)
	}
	return plannedTables, allTables
}

func ToState(configFromPlan types.Object, configuration *client.PluginConfiguration, isImportRead bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planFields, planSensitiveFields, planTables *types.Set

	planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	if ok {
		setVal := planFieldsValue.(types.Set)
		planFields = &setVal
	}
	planSensitiveFieldsValue, ok := configFromPlan.Attributes()["sensitive_fields"]
	if ok {
		setVal := planSensitiveFieldsValue.(types.Set)
		planSensitiveFields = &setVal
	}
	planTablesValue, ok := configFromPlan.Attributes()["tables"]
	if ok {
		setVal := planTablesValue.(types.Set)
		planTables = &setVal
	}

	fieldsAttrValue, sensitiveFieldsAttrValue, fieldsAllAttrValue := toFieldsSetValue(configuration.Fields, planFields, planSensitiveFields, isImportRead, &diags)
	tablesAttrValue, tablesAllAttrValue := toTablesSetValue(configuration.Tables, planTables, isImportRead, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields":           fieldsAttrValue,
		"sensitive_fields": sensitiveFieldsAttrValue,
		"fields_all":       fieldsAllAttrValue,
		"tables":           tablesAttrValue,
		"tables_all":       tablesAllAttrValue,
	}
	configObj, valueFromDiags := types.ObjectValue(configurationAttrTypes, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
}

// Mark fields_all and tables_all configuration as unknown if the fields and tables have changed in the plan
func MarkComputedAttrsUnknownOnChange(planConfiguration, stateConfiguration types.Object) (types.Object, diag.Diagnostics) {
	if !internaltypes.IsDefined(planConfiguration) || !internaltypes.IsDefined(stateConfiguration) {
		return planConfiguration, nil
	}
	planConfigurationAttrs := planConfiguration.Attributes()
	planFields := planConfiguration.Attributes()["fields"]
	stateFields := stateConfiguration.Attributes()["fields"]
	if !planFields.Equal(stateFields) {
		planConfigurationAttrs["fields_all"] = types.SetUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
	}

	planTables := planConfiguration.Attributes()["tables"]
	stateTables := stateConfiguration.Attributes()["tables"]
	if !planTables.Equal(stateTables) {
		planConfigurationAttrs["tables_all"] = types.SetUnknown(types.ObjectType{AttrTypes: tableAttrTypes})
	}

	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
}

// Mark fields_all and tables_all configuration as unknown
func MarkComputedAttrsUnknown(planConfiguration types.Object) (types.Object, diag.Diagnostics) {
	if !internaltypes.IsDefined(planConfiguration) {
		return planConfiguration, nil
	}
	planConfigurationAttrs := planConfiguration.Attributes()
	planConfigurationAttrs["fields_all"] = types.SetUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
	planConfigurationAttrs["tables_all"] = types.SetUnknown(types.ObjectType{AttrTypes: tableAttrTypes})
	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
}
