package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":      types.StringType,
		"value":     types.StringType,
		"inherited": types.BoolType,
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":      types.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tableAttrTypes = map[string]attr.Type{
		"name":      types.StringType,
		"rows":      types.ListType{ElemType: types.ObjectType{AttrTypes: rowAttrTypes}},
		"inherited": types.BoolType,
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields":     types.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all": types.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables":     types.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
		"tables_all": types.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}
)

func AttrTypes() map[string]attr.Type {
	return configurationAttrTypes
}

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

	for _, field := range fields {
		attrValues := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(field.Name)
		attrValues["value"] = types.StringPointerValue(field.Value)
		attrValues["inherited"] = types.BoolPointerValue(field.Inherited)

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

// Creates state values for tables. Returns one value that only includes values specified in the plan, and a second value that includes all tables values
func ToTablesListValue(tables []client.ConfigTable, planTables *types.List, diags *diag.Diagnostics) (types.List, types.List) {
	// List of *all* tables values to return
	finalTablesAllObjValues := []attr.Value{}
	// List of tables values to return that were expected based on the plan
	finalTablesObjValues := []attr.Value{}
	// types.Object values for tables included in the plan
	planTableObjs := map[string]types.Object{}
	if planTables == nil {
		for _, table := range tables {
			attrValues := map[string]attr.Value{}
			attrValues["inherited"] = types.BoolPointerValue(table.Inherited)
			attrValues["name"] = types.StringValue(table.Name)
			attrValues["rows"] = ToRowsListValue(table.Rows, nil, diags)
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
			attrValues["inherited"] = types.BoolPointerValue(tables[i].Inherited)
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
			attrValues["rows"] = ToRowsListValue(tables[i].Rows, planTableRows, diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			finalTablesAllObjValues = append(finalTablesAllObjValues, tableObjValue)
			if inPlan {
				finalTablesObjValues = append(finalTablesObjValues, tableObjValue)
			}
		}
	}
	plannedTables, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: tableAttrTypes,
	}, finalTablesObjValues)
	diags.Append(newDiags...)
	allTables, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: tableAttrTypes,
	}, finalTablesAllObjValues)
	diags.Append(newDiags...)
	return plannedTables, allTables
}

func ToState(configFromPlan types.Object, configuration *client.PluginConfiguration) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planFields, planTables *types.List

	planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	if ok {
		listVal := planFieldsValue.(types.List)
		planFields = &listVal
	}
	planTablesValue, ok := configFromPlan.Attributes()["tables"]
	if ok {
		listVal := planTablesValue.(types.List)
		planTables = &listVal
	}

	fieldsAttrValue, fieldsAllAttrValue := ToFieldsListValue(configuration.Fields, planFields, &diags)
	tablesAttrValue, tablesAllAttrValue := ToTablesListValue(configuration.Tables, planTables, &diags)

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
		planConfigurationAttrs["fields_all"] = types.ListUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
	}

	planTables := planConfiguration.Attributes()["tables"]
	stateTables := stateConfiguration.Attributes()["tables"]
	if !planTables.Equal(stateTables) {
		planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tableAttrTypes})
	}

	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
}

// Mark fields_all and tables_all configuration as unknown
func MarkComputedAttrsUnknown(planConfiguration types.Object) (types.Object, diag.Diagnostics) {
	if !internaltypes.IsDefined(planConfiguration) {
		return planConfiguration, nil
	}
	planConfigurationAttrs := planConfiguration.Attributes()
	planConfigurationAttrs["fields_all"] = types.ListUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
	planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tableAttrTypes})
	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
}
