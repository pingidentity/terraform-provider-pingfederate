// Copyright © 2025 Ping Identity Corporation

package pluginconfiguration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
	sensitiveFieldAttrTypes = map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}

	rowsSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: sensitiveFieldAttrTypes}},
		"default_row":      types.BoolType,
	}
	rowsMergedFieldsAttrTypes = map[string]attr.Type{
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tablesSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsSensitiveFieldsSplitAttrTypes}},
	}
	tablesMergedFieldsAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsMergedFieldsAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: sensitiveFieldAttrTypes}},
		"fields_all":       types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables":           types.ListType{ElemType: types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes}},
		"tables_all":       types.ListType{ElemType: types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes}},
	}
)

func AttrTypes() map[string]attr.Type {
	return configurationAttrTypes
}

type pfConfigurationFieldsResult struct {
	plannedCleartextFields types.Set
	plannedSensitiveFields types.Set
	allCleartextFields     types.Set
	allSensitiveFields     types.Set
	allFields              types.Set
}

type pfConfigurationRowsResult struct {
	allRowsSensitiveFieldsSplit types.List
	allRowsMergedFields         types.List
}

type pfConfigurationTablesResult struct {
	plannedTables                 types.List
	allTablesSensitiveFieldsSplit types.List
	allTablesMergedFields         types.List
}

func readFieldsResponse(fields []client.ConfigField, planFields, planSensitiveFields *types.Set, diags *diag.Diagnostics) pfConfigurationFieldsResult {
	plannedCleartextFields := []attr.Value{}
	plannedSensitiveFields := []attr.Value{}
	allCleartextFields := []attr.Value{}
	allSensitiveFields := []attr.Value{}
	allFields := []attr.Value{}
	plannedFieldsValues := map[string]*string{}
	plannedSensitiveFieldsValues := map[string]*string{}
	plannedSensitiveFieldsEncryptedValues := map[string]*string{}
	// Build up a map of all the values from the plan
	if planFields != nil {
		for _, planField := range planFields.Elements() {
			planFieldObj := planField.(types.Object)
			plannedFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}
	if planSensitiveFields != nil {
		for _, planField := range planSensitiveFields.Elements() {
			planFieldObj := planField.(types.Object)
			plannedSensitiveFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
			plannedSensitiveFieldsEncryptedValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["encrypted_value"].(types.String).ValueStringPointer()
		}
	}

	for _, field := range fields {
		attrValues := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(field.Name)
		attrValues["value"] = types.StringPointerValue(field.Value)

		// If this field is in the plan, add it to the list of plan fields
		fieldAdded := false
		if planFields != nil {
			planValue, ok := plannedFieldsValues[field.Name]
			if ok {
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.EncryptedValue != nil && *field.EncryptedValue != "" {
					diags.AddAttributeWarning(
						path.Root("configuration"),
						providererror.ConfigurationWarning,
						fmt.Sprintf("Field with name %s was return encrypted by the PingFederate API. If the field is sensitive, move it to the `sensitive_fields` attribute.", field.Name),
					)
				}
				// If PF sets a default for the field when the user specifies an empty value,
				// just use the plan value and let the PF value appear in fields_all
				if field.Value == nil || (*planValue == "" && *field.Value != *planValue) {
					planAttrValues["value"] = types.StringPointerValue(planValue)
				} else {
					planAttrValues["value"] = types.StringPointerValue(field.Value)
				}
				objVal, respDiags := types.ObjectValue(fieldAttrTypes, planAttrValues)
				diags.Append(respDiags...)
				plannedCleartextFields = append(plannedCleartextFields, objVal)
				fieldAdded = true
			}
		}
		if planSensitiveFields != nil && !fieldAdded {
			planValue, ok := plannedSensitiveFieldsValues[field.Name]
			if ok {
				planEncryptedValue := plannedSensitiveFieldsEncryptedValues[field.Name]
				planAttrValues := map[string]attr.Value{}
				planAttrValues["name"] = types.StringValue(field.Name)
				if field.EncryptedValue == nil && field.Value != nil && *field.Value != "" {
					diags.AddAttributeWarning(
						path.Root("configuration"),
						providererror.ConfigurationWarning,
						fmt.Sprintf("Sensitive field with name %s was returned in cleartext by the PingFederate API. If the field is not sensitive, move it to the `fields` attribute.", field.Name),
					)
				}
				if field.Value == nil {
					planAttrValues["value"] = types.StringPointerValue(planValue)
				} else {
					planAttrValues["value"] = types.StringPointerValue(field.Value)
				}
				if planEncryptedValue != nil {
					planAttrValues["encrypted_value"] = types.StringPointerValue(planEncryptedValue)
				} else {
					planAttrValues["encrypted_value"] = types.StringPointerValue(field.EncryptedValue)
				}
				objVal, respDiags := types.ObjectValue(sensitiveFieldAttrTypes, planAttrValues)
				diags.Append(respDiags...)
				plannedSensitiveFields = append(plannedSensitiveFields, objVal)
			}
		}

		objVal, respDiags := types.ObjectValue(fieldAttrTypes, attrValues)
		diags.Append(respDiags...)
		allFields = append(allFields, objVal)
		if field.EncryptedValue != nil && *field.EncryptedValue != "" {
			sensitiveAttrValues := map[string]attr.Value{}
			sensitiveAttrValues["name"] = types.StringValue(field.Name)
			sensitiveAttrValues["value"] = types.StringPointerValue(field.Value)
			sensitiveAttrValues["encrypted_value"] = types.StringPointerValue(field.EncryptedValue)
			sensitiveObjVal, respDiags := types.ObjectValue(sensitiveFieldAttrTypes, sensitiveAttrValues)
			diags.Append(respDiags...)
			allSensitiveFields = append(allSensitiveFields, sensitiveObjVal)
		} else {
			allCleartextFields = append(allCleartextFields, objVal)
		}
	}

	plannedCleartextFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, plannedCleartextFields)
	diags.Append(respDiags...)
	plannedSensitiveFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: sensitiveFieldAttrTypes,
	}, plannedSensitiveFields)
	diags.Append(respDiags...)

	allCleartextFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allCleartextFields)
	diags.Append(respDiags...)
	allSensitiveFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: sensitiveFieldAttrTypes,
	}, allSensitiveFields)
	diags.Append(respDiags...)

	allFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allFields)
	diags.Append(respDiags...)

	return pfConfigurationFieldsResult{
		plannedCleartextFields: plannedCleartextFieldsSet,
		plannedSensitiveFields: plannedSensitiveFieldsSet,
		allCleartextFields:     allCleartextFieldsSet,
		allSensitiveFields:     allSensitiveFieldsSet,
		allFields:              allFieldsSet,
	}
}

func readRowsResponse(rows []client.ConfigRow, planRows *types.List, diags *diag.Diagnostics) pfConfigurationRowsResult {
	var rowsMergedFields, rowsSensitiveFieldsSplit []attr.Value
	if planRows == nil || planRows.IsNull() {
		if len(rows) == 0 {
			// If the API returned no rows, treat as null
			return pfConfigurationRowsResult{
				allRowsSensitiveFieldsSplit: types.ListNull(types.ObjectType{AttrTypes: rowsSensitiveFieldsSplitAttrTypes}),
				allRowsMergedFields:         types.ListNull(types.ObjectType{AttrTypes: rowsMergedFieldsAttrTypes}),
			}
		}
		for _, row := range rows {
			attrValues := map[string]attr.Value{}
			attrValuesSensitiveSplit := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(row.DefaultRow)
			attrValuesSensitiveSplit["default_row"] = types.BoolPointerValue(row.DefaultRow)

			rowFields := readFieldsResponse(row.Fields, nil, nil, diags)
			attrValues["fields"] = rowFields.allFields
			attrValuesSensitiveSplit["fields"] = rowFields.allCleartextFields
			attrValuesSensitiveSplit["sensitive_fields"] = rowFields.allSensitiveFields

			rowMergedFields, respDiags := types.ObjectValue(rowsMergedFieldsAttrTypes, attrValues)
			diags.Append(respDiags...)
			rowsMergedFields = append(rowsMergedFields, rowMergedFields)
			rowSensitiveFieldsSplit, respDiags := types.ObjectValue(rowsSensitiveFieldsSplitAttrTypes, attrValuesSensitiveSplit)
			diags.Append(respDiags...)
			rowsSensitiveFieldsSplit = append(rowsSensitiveFieldsSplit, rowSensitiveFieldsSplit)
		}
	} else {
		// This is assuming there are never any rows added by the PF API. If there
		// are ever rows added, this will cause a nil pointer exception trying to read
		// index i of planRowsElements.
		planRowsElements := planRows.Elements()
		for i := 0; i < len(rows); i++ {
			attrValues := map[string]attr.Value{}
			attrValuesSensitiveSplit := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(rows[i].DefaultRow)
			attrValuesSensitiveSplit["default_row"] = types.BoolPointerValue(rows[i].DefaultRow)
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

			rowFields := readFieldsResponse(rows[i].Fields, planRowFields, planRowSensitiveFields, diags)
			attrValues["fields"] = rowFields.allFields
			attrValuesSensitiveSplit["fields"] = rowFields.plannedCleartextFields
			attrValuesSensitiveSplit["sensitive_fields"] = rowFields.plannedSensitiveFields

			rowMergedFields, respDiags := types.ObjectValue(rowsMergedFieldsAttrTypes, attrValues)
			diags.Append(respDiags...)
			rowsMergedFields = append(rowsMergedFields, rowMergedFields)
			rowSensitiveFieldsSplit, respDiags := types.ObjectValue(rowsSensitiveFieldsSplitAttrTypes, attrValuesSensitiveSplit)
			diags.Append(respDiags...)
			rowsSensitiveFieldsSplit = append(rowsSensitiveFieldsSplit, rowSensitiveFieldsSplit)
		}
	}

	rowsMergedFieldsList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: rowsMergedFieldsAttrTypes,
	}, rowsMergedFields)
	diags.Append(respDiags...)
	rowsSensitiveFieldsSplitList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: rowsSensitiveFieldsSplitAttrTypes,
	}, rowsSensitiveFieldsSplit)
	diags.Append(respDiags...)
	return pfConfigurationRowsResult{
		allRowsSensitiveFieldsSplit: rowsSensitiveFieldsSplitList,
		allRowsMergedFields:         rowsMergedFieldsList,
	}
}

func toTablesSetValue(tables []client.ConfigTable, planTables *types.List, diags *diag.Diagnostics) pfConfigurationTablesResult {
	// List of *all* tables values to return
	allTablesMergedFields := []attr.Value{}
	// List of *all* tables values to return split into sensitive and non-sensitive fields
	allTablesSensitiveFieldsSplit := []attr.Value{}
	// List of tables values to return that were expected based on the plan
	plannedTables := []attr.Value{}
	// types.Object values for tables included in the plan
	planTableObjs := map[string]types.Object{}
	if planTables != nil {
		// Build up a map of all the tables included in the plan
		for _, planTable := range planTables.Elements() {
			planTableObj := planTable.(types.Object)
			planTableObjs[planTableObj.Attributes()["name"].(types.String).ValueString()] = planTableObj
		}
	}

	for i := 0; i < len(tables); i++ {
		attrValues := map[string]attr.Value{}
		attrValuesSensitiveSplit := map[string]attr.Value{}
		attrValues["name"] = types.StringValue(tables[i].Name)
		attrValuesSensitiveSplit["name"] = types.StringValue(tables[i].Name)
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

		tableRows := readRowsResponse(tables[i].Rows, planTableRows, diags)
		attrValues["rows"] = tableRows.allRowsMergedFields
		attrValuesSensitiveSplit["rows"] = tableRows.allRowsSensitiveFieldsSplit

		tableMergedFields, respDiags := types.ObjectValue(tablesMergedFieldsAttrTypes, attrValues)
		diags.Append(respDiags...)
		allTablesMergedFields = append(allTablesMergedFields, tableMergedFields)
		tableSensitiveFieldsSplit, respDiags := types.ObjectValue(tablesSensitiveFieldsSplitAttrTypes, attrValuesSensitiveSplit)
		diags.Append(respDiags...)
		allTablesSensitiveFieldsSplit = append(allTablesSensitiveFieldsSplit, tableSensitiveFieldsSplit)
		if inPlan {
			plannedTables = append(plannedTables, tableSensitiveFieldsSplit)
		}
	}

	allTablesMergedFieldsList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: tablesMergedFieldsAttrTypes,
	}, allTablesMergedFields)
	diags.Append(respDiags...)
	allTablesSensitiveFieldsSplitList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: tablesSensitiveFieldsSplitAttrTypes,
	}, allTablesSensitiveFieldsSplit)
	diags.Append(respDiags...)
	plannedTablesList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: tablesSensitiveFieldsSplitAttrTypes,
	}, plannedTables)
	diags.Append(respDiags...)

	return pfConfigurationTablesResult{
		plannedTables:                 plannedTablesList,
		allTablesSensitiveFieldsSplit: allTablesSensitiveFieldsSplitList,
		allTablesMergedFields:         allTablesMergedFieldsList,
	}
}

func ToState(configFromPlan types.Object, configuration *client.PluginConfiguration, isImportRead bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planFields, planSensitiveFields *types.Set
	var planTables *types.List

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
		listVal := planTablesValue.(types.List)
		planTables = &listVal
	}

	fields := readFieldsResponse(configuration.Fields, planFields, planSensitiveFields, &diags)
	tables := toTablesSetValue(configuration.Tables, planTables, &diags)

	fieldsAttrValue := fields.plannedCleartextFields
	sensitiveFieldsAttrValue := fields.plannedSensitiveFields
	tablesAttrValue := tables.plannedTables
	if isImportRead {
		fieldsAttrValue = fields.allCleartextFields
		sensitiveFieldsAttrValue = fields.allSensitiveFields
		tablesAttrValue = tables.allTablesSensitiveFieldsSplit
	}

	configurationAttrValue := map[string]attr.Value{
		"fields":           fieldsAttrValue,
		"sensitive_fields": sensitiveFieldsAttrValue,
		"fields_all":       fields.allFields,
		"tables":           tablesAttrValue,
		"tables_all":       tables.allTablesMergedFields,
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
		planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes})
	} else {
		// If there is a tables_all table with rows defined, ensure there is a table in the plan with
		// the same number of rows defined. If not, assume it's a drift, and set tables_all to unknown
		tablesAllTablesWithRows := map[string]types.Object{}
		matchesFound := map[string]bool{}
		stateTablesAll := stateConfiguration.Attributes()["tables_all"].(types.List)
		for _, table := range stateTablesAll.Elements() {
			tableObj := table.(types.Object)
			rows, ok := tableObj.Attributes()["rows"]
			if ok && !rows.IsNull() && !rows.IsUnknown() && len(rows.(types.List).Elements()) > 0 {
				name := tableObj.Attributes()["name"].(types.String).ValueString()
				tablesAllTablesWithRows[name] = tableObj
				matchesFound[name] = false
			}
		}

		// Look for tables in the plan that match the tables in tables_all
		planTables := planConfiguration.Attributes()["tables"].(types.List)
		for _, table := range planTables.Elements() {
			tableObj := table.(types.Object)
			name := tableObj.Attributes()["name"].(types.String).ValueString()
			tablesAllTable, ok := tablesAllTablesWithRows[name]
			if ok && !tableObj.IsNull() && !tableObj.IsUnknown() {
				// Compare length of rows
				tablesAllRowCount := len(tablesAllTable.Attributes()["rows"].(types.List).Elements())
				planTableRowCount := len(tableObj.Attributes()["rows"].(types.List).Elements())
				if tablesAllRowCount == planTableRowCount {
					matchesFound[name] = true
				}
			}
		}

		// If there was no match found for at least one of the tables in tables_all, set tables_all to unknown
		for _, found := range matchesFound {
			if !found {
				planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes})
				break
			}
		}
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
	planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes})
	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
}
