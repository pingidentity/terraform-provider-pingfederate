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
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tableAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields":     types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables":     types.SetType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
		"tables_all": types.SetType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}
)

func AttrTypes() map[string]attr.Type {
	return configurationAttrTypes
}

// Creates state values for fields. Returns one value that only includes values specified in the plan, and a second value that includes all fields values
func toFieldsSetValue(fields []client.ConfigField, planFields *types.Set, isImportRead bool, diags *diag.Diagnostics) (types.Set, types.Set) {
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

	for _, field := range fields {
		attrValues := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(field.Name)
		attrValues["value"] = types.StringPointerValue(field.Value)

		// If this field is in the plan, add it to the list of plan fields
		if planFields != nil {
			planValue, ok := planFieldsValues[field.Name]
			if ok {
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.Value == nil {
					planAttrValues["value"] = types.StringPointerValue(planValue)
				} else {
					planAttrValues["value"] = types.StringPointerValue(field.Value)
				}
				objVal, newDiags := types.ObjectValue(fieldAttrTypes, planAttrValues)
				diags.Append(newDiags...)
				plannedObjValues = append(plannedObjValues, objVal)
			}
		}

		objVal, newDiags := types.ObjectValue(fieldAttrTypes, attrValues)
		diags.Append(newDiags...)
		allObjValues = append(allObjValues, objVal)
	}

	allSetVal, newDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allObjValues)
	diags.Append(newDiags...)
	var plannedSetVal types.Set
	if isImportRead {
		// On imports, just read everything directly into the "fields" attribute
		plannedSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, allObjValues)
		diags.Append(newDiags...)
	} else {
		plannedSetVal, newDiags = types.SetValue(types.ObjectType{
			AttrTypes: fieldAttrTypes,
		}, plannedObjValues)
		diags.Append(newDiags...)
	}
	return plannedSetVal, allSetVal
}

func toRowsListValue(rows []client.ConfigRow, planRows *types.List, isImportRead bool, diags *diag.Diagnostics) types.List {
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
			_, attrValues["fields"] = toFieldsSetValue(row.Fields, nil, isImportRead, diags)
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
			var planRowFields *types.Set
			planRowFieldsVal, ok := planRow.Attributes()["fields"]
			if ok {
				setVal := planRowFieldsVal.(types.Set)
				planRowFields = &setVal
			}
			attrValues["fields"], _ = toFieldsSetValue(rows[i].Fields, planRowFields, isImportRead, diags)
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
		AttrTypes: tableAttrTypes,
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
	var planFields, planTables *types.Set

	planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	if ok {
		setVal := planFieldsValue.(types.Set)
		planFields = &setVal
	}
	planTablesValue, ok := configFromPlan.Attributes()["tables"]
	if ok {
		setVal := planTablesValue.(types.Set)
		planTables = &setVal
	}

	fieldsAttrValue, fieldsAllAttrValue := toFieldsSetValue(configuration.Fields, planFields, isImportRead, &diags)
	tablesAttrValue, tablesAllAttrValue := toTablesSetValue(configuration.Tables, planTables, isImportRead, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields":     fieldsAttrValue,
		"fields_all": fieldsAllAttrValue,
		"tables":     tablesAttrValue,
		"tables_all": tablesAllAttrValue,
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
