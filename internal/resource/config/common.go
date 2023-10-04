package config

import (
	"fmt"

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
