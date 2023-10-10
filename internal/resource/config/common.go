package config

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// Get schema elements common to all resources
func AddCommonSchema(s *schema.Schema) {
	s.Attributes["id"] = schema.StringAttribute{
		Description: "The ID of this resource.",
		Computed:    true,
		Required:    false,
		Optional:    false,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func AddResourceLinkSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
		"location": schema.StringAttribute{
			Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			Computed:    true,
			Optional:    false,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

func SetAllAttributesToOptionalAndComputed(s *schema.Schema, exemptAttributes []string) {
	for key, attribute := range s.Attributes {
		// If more attribute types are used by this provider, this method will need to be updated
		if !internaltypes.StringSliceContains(exemptAttributes, key) {
			stringAttr, ok := attribute.(schema.StringAttribute)
			anyOk := ok
			if ok && (!stringAttr.Computed || !stringAttr.Optional) {
				stringAttr.Required = false
				stringAttr.Optional = true
				stringAttr.Computed = true
				stringAttr.PlanModifiers = append(stringAttr.PlanModifiers, stringplanmodifier.UseStateForUnknown())
				s.Attributes[key] = stringAttr
				continue
			}
			setAttr, ok := attribute.(schema.SetAttribute)
			anyOk = ok || anyOk
			if ok && (!setAttr.Computed || !setAttr.Optional) {
				setAttr.Required = false
				setAttr.Optional = true
				setAttr.Computed = true
				setAttr.PlanModifiers = append(setAttr.PlanModifiers, setplanmodifier.UseStateForUnknown())
				s.Attributes[key] = setAttr
				continue
			}
			listAttr, ok := attribute.(schema.ListAttribute)
			anyOk = ok || anyOk
			if ok && (!listAttr.Computed || !listAttr.Optional) {
				listAttr.Required = false
				listAttr.Optional = true
				listAttr.Computed = true
				listAttr.PlanModifiers = append(listAttr.PlanModifiers, listplanmodifier.UseStateForUnknown())
				s.Attributes[key] = listAttr
				continue
			}
			boolAttr, ok := attribute.(schema.BoolAttribute)
			anyOk = ok || anyOk
			if ok && (!boolAttr.Computed || !boolAttr.Optional) {
				boolAttr.Required = false
				boolAttr.Optional = true
				boolAttr.Computed = true
				boolAttr.PlanModifiers = append(boolAttr.PlanModifiers, boolplanmodifier.UseStateForUnknown())
				s.Attributes[key] = boolAttr
				continue
			}
			intAttr, ok := attribute.(schema.Int64Attribute)
			anyOk = ok || anyOk
			if ok && (!intAttr.Computed || !intAttr.Optional) {
				intAttr.Required = false
				intAttr.Optional = true
				intAttr.Computed = true
				intAttr.PlanModifiers = append(intAttr.PlanModifiers, int64planmodifier.UseStateForUnknown())
				s.Attributes[key] = intAttr
				continue
			}
			floatAttr, ok := attribute.(schema.Float64Attribute)
			anyOk = ok || anyOk
			if ok && (!floatAttr.Computed || !floatAttr.Optional) {
				floatAttr.Required = false
				floatAttr.Optional = true
				floatAttr.Computed = true
				floatAttr.PlanModifiers = append(floatAttr.PlanModifiers, float64planmodifier.UseStateForUnknown())
				s.Attributes[key] = floatAttr
				continue
			}
			if !anyOk {
				return
			}
		}
	}
}

func ConfigurationToState(planConfiguration types.Object, configuration client.PluginConfiguration) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	configurationAttrType := map[string]attr.Type{
		"fields":     basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables":     basetypes.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}

	var planFields, planTables *types.List

	planFieldsValue, ok := planConfiguration.Attributes()["fields"]
	if ok {
		listVal := planFieldsValue.(types.List)
		planFields = &listVal
	}
	planTablesValue, ok := planConfiguration.Attributes()["tables"]
	if ok {
		listVal := planTablesValue.(types.List)
		planTables = &listVal
	}

	fieldsAttrValue, fieldsAllAttrValue := ToFieldsListValue(configuration.Fields, planFields, &diags)
	tablesAttrValue := ToTablesListValue(configuration.Tables, planTables, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields":     fieldsAttrValue,
		"fields_all": fieldsAllAttrValue,
		"tables":     tablesAttrValue,
	}
	configObj, valueFromDiags := types.ObjectValue(configurationAttrType, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
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
	// If rows is null in the plan, just return everything. Otherwise only return rows corresponding with the plan
	if planRows == nil {
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
	// If tables is null in the plan, just return everything. Otherwise only return tables corresponding with the plan
	if planTables == nil {
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
