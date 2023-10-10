package pluginconfiguration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":      basetypes.StringType{},
		"value":     basetypes.StringType{},
		"inherited": basetypes.BoolType{},
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": basetypes.BoolType{},
	}

	tableAttrTypes = map[string]attr.Type{
		"name":      basetypes.StringType{},
		"rows":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: rowAttrTypes}},
		"inherited": basetypes.BoolType{},
	}
)

func FieldAttrTypes() map[string]attr.Type {
	return fieldAttrTypes
}

func TableAttrTypes() map[string]attr.Type {
	return tableAttrTypes
}

func ToFieldsListValue(fields []client.ConfigField, planFields types.List, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	planFieldsElements := planFields.Elements()
	// If fields is null in the plan, just return everything. Otherwise only return fields corresponding with the plan
	//TODO We will want to change this in the future so that we store everything the server returns in state in some way
	if !internaltypes.IsDefined(planFields) {
		for _, field := range fields {
			attrValues := map[string]attr.Value{}
			attrValues["name"] = types.StringValue(field.Name)
			if field.Value == nil {
				attrValues["value"] = types.StringNull()
			} else {
				attrValues["value"] = types.StringPointerValue(field.Value)
			}
			attrValues["inherited"] = types.BoolPointerValue(field.Inherited)
			objVal, newDiags := types.ObjectValue(fieldAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, objVal)
		}
	} else {
		if len(planFieldsElements) > len(fields) {
			diags.AddError("Plan fields length is greater than response fields length",
				fmt.Sprintf("Plan fields: %d, response fields: %d", len(planFieldsElements), len(fields)))
			return types.ListNull(types.ObjectType{AttrTypes: fieldAttrTypes})
		}
		for i := 0; i < len(planFieldsElements); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["name"] = types.StringValue(fields[i].Name)
			if fields[i].Value == nil {
				// This must be an encrypted field. Use the value from the plan
				planField := planFieldsElements[i].(types.Object)
				planFieldValue := planField.Attributes()["value"].(types.String)
				attrValues["value"] = types.StringValue(planFieldValue.ValueString())
			} else {
				attrValues["value"] = types.StringPointerValue(fields[i].Value)
			}
			attrValues["inherited"] = types.BoolPointerValue(fields[i].Inherited)
			objVal, newDiags := types.ObjectValue(fieldAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, objVal)
		}
	}
	listVal, newDiags := types.ListValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, objValues)
	diags.Append(newDiags...)
	return listVal
}

func ToRowsListValue(rows []client.ConfigRow, planRows types.List, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	planRowsElements := planRows.Elements()
	// If rows is null in the plan, just return everything. Otherwise only return rows corresponding with the plan
	//TODO We will want to change this in the future so that we store everything the server returns in state in some way
	if !internaltypes.IsDefined(planRows) {
		for _, row := range rows {
			attrValues := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(row.DefaultRow)
			attrValues["fields"] = ToFieldsListValue(row.Fields, types.ListNull(types.ObjectType{AttrTypes: fieldAttrTypes}), diags)
			rowObjVal, newDiags := types.ObjectValue(rowAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, rowObjVal)
		}
	} else {
		if len(planRowsElements) > len(rows) {
			diags.AddError("Plan rows length is greater than response rows length",
				fmt.Sprintf("Plan tables: %d, response tables: %d", len(planRowsElements), len(rows)))
			return types.ListNull(types.ObjectType{AttrTypes: rowAttrTypes})
		}
		for i := 0; i < len(planRowsElements); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["default_row"] = types.BoolPointerValue(rows[i].DefaultRow)
			planRow := planRowsElements[i].(types.Object)
			planRowFields := types.ListNull(types.ObjectType{AttrTypes: fieldAttrTypes})
			planRowFieldsVal, ok := planRow.Attributes()["fields"]
			if ok {
				planRowFields = planRowFieldsVal.(types.List)
			}
			attrValues["fields"] = ToFieldsListValue(rows[i].Fields, planRowFields, diags)
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

func ToTablesListValue(tables []client.ConfigTable, planTables types.List, diags *diag.Diagnostics) types.List {
	objValues := []attr.Value{}
	planTablesElements := planTables.Elements()
	// If tables is null in the plan, just return everything. Otherwise only return tables corresponding with the plan
	//TODO We will want to change this in the future so that we store everything the server returns in state in some way
	if !internaltypes.IsDefined(planTables) {
		for _, table := range tables {
			attrValues := map[string]attr.Value{}
			attrValues["inherited"] = types.BoolPointerValue(table.Inherited)
			attrValues["name"] = types.StringValue(table.Name)
			attrValues["rows"] = ToRowsListValue(table.Rows, types.ListNull(types.ObjectType{AttrTypes: rowAttrTypes}), diags)
			tableObjValue, newDiags := types.ObjectValue(tableAttrTypes, attrValues)
			diags.Append(newDiags...)
			objValues = append(objValues, tableObjValue)
		}
	} else {
		if len(planTablesElements) > len(tables) {
			diags.AddError("Plan tables length is greater than response tables length",
				fmt.Sprintf("Plan tables: %d, response tables: %d", len(planTablesElements), len(tables)))
			return types.ListNull(types.ObjectType{AttrTypes: rowAttrTypes})
		}
		for i := 0; i < len(planTablesElements); i++ {
			attrValues := map[string]attr.Value{}
			attrValues["inherited"] = types.BoolPointerValue(tables[i].Inherited)
			attrValues["name"] = types.StringValue(tables[i].Name)
			planTable := planTablesElements[i].(types.Object)
			planTableRows := types.ListNull(types.ObjectType{AttrTypes: rowAttrTypes})
			planTableRowsVal, ok := planTable.Attributes()["rows"]
			if ok {
				planTableRows = planTableRowsVal.(types.List)
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

func ToState(configFromPlan basetypes.ObjectValue, configuration *client.PluginConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags, methodDiags diag.Diagnostics
	configurationAttrType := map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: FieldAttrTypes()}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: TableAttrTypes()}},
	}

	planFields := types.ListNull(types.ObjectType{AttrTypes: FieldAttrTypes()})
	planTables := types.ListNull(types.ObjectType{AttrTypes: TableAttrTypes()})

	planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	if ok {
		planFields = planFieldsValue.(types.List)
	}
	planTablesValue, ok := configFromPlan.Attributes()["tables"]
	if ok {
		planTables = planTablesValue.(types.List)
	}
	fieldsAttrValue := ToFieldsListValue(configuration.Fields, planFields, &diags)
	tablesAttrValue := ToTablesListValue(configuration.Tables, planTables, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}
	configurationToStateObj, methodDiags := types.ObjectValue(configurationAttrType, configurationAttrValue)
	diags.Append(methodDiags...)
	return configurationToStateObj, diags
}
