package types

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Return true if this types.String represents an empty (but non-null and non-unknown) string
func IsEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() == ""
}
func IsNonEmptyMap(m types.Map) bool {
	return !m.IsNull() && !m.IsUnknown() && m.Elements() != nil
}

func IsNonEmptyList(l types.List) bool {
	return !l.IsNull() && !l.IsUnknown() && l.Elements() != nil
}

func IsNonEmptyObj(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown() && obj.Attributes() != nil
}

func ObjContainsNoEmptyVals(obj types.Object) bool {
	for _, objVal := range obj.Attributes() {
		if !IsDefined(objVal) {
			return true
		}
	}
	return false
}

// Return true if this types.String represents a non-empty, non-null, non-unknown string
func IsNonEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() != ""
}

// Return true if this value represents a defined (non-null and non-unknown) value
func IsDefined(value attr.Value) bool {
	return value != nil && !value.IsNull() && !value.IsUnknown()
}

// Check if an attribute slice contains a value
func Contains(slice []attr.Value, value attr.Value) bool {
	for _, element := range slice {
		if element.Equal(value) {
			return true
		}
	}
	return false
}

// Check if a string slice contains a value
func StringSliceContains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

// Check if two slices representing sets are equal
func SetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Assuming there are no duplicate elements since the slices represent sets
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if bElement == aElement {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Check if two float slices representing sets are equal
func FloatSetsEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	// Assuming there are no duplicate elements since the slices represent sets
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if aElement == bElement {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Compare two given sets, match string if found in both
func MatchStringInSets(a, b []string) *string {
	for _, aElem := range a {
		for _, bElem := range b {
			if aElem == bElem {
				return &aElem
			}
		}
	}
	return nil
}

func CamelCaseToUnderscores(s string) string {
	re, _ := regexp.Compile(`([A-Z])`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToLower("_" + m[0:])
	})
	return res
}

func UnderscoresToCamelCase(s string) string {
	re, _ := regexp.Compile(`(_[A-Za-z])`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	return res
}

// Converts the basetypes.MapValue to map[string]interface{} required for PingFederate Client
func MapValuesToClientMap(mv basetypes.MapValue, con context.Context) *map[string]interface{} {
	type StringMap map[string]string
	var value StringMap
	mv.ElementsAs(con, &value, false)
	converted := map[string]interface{}{}
	for k, v := range value {
		converted[k] = v
	}
	return &converted
}

// Converts the types.Object to *map[string]interface{} required for PingFederate Client
func ObjValuesToClientMap(obj types.Object) *map[string]interface{} {
	attrs := obj.Attributes()
	converted := map[string]interface{}{}
	for key, value := range attrs {
		strvalue, ok := value.(basetypes.StringValue)
		if ok {
			if strvalue.IsNull() || strvalue.IsUnknown() {
				continue
			} else {
				converted[UnderscoresToCamelCase(key)] = strvalue.ValueString()
				continue
			}
		}
		boolvalue, ok := value.(basetypes.BoolValue)
		if ok {
			converted[UnderscoresToCamelCase(key)] = boolvalue.ValueBool()
			continue
		}
		int64value, ok := value.(basetypes.Int64Value)
		if ok {
			converted[UnderscoresToCamelCase(key)] = int64value.ValueInt64()
			continue
		}
	}

	return &converted
}

// Converts the types.Object to map[string]interface{} required for PingFederate Client
func ObjValuesToMapNoPointer(obj types.Object) map[string]interface{} {
	attrs := obj.Attributes()
	converted := map[string]interface{}{}
	for key, value := range attrs {
		strvalue, ok := value.(basetypes.StringValue)
		if ok {
			if strvalue.IsNull() || strvalue.IsUnknown() {
				continue
			} else {
				converted[UnderscoresToCamelCase(key)] = strvalue.ValueString()
				continue
			}
		}
		boolvalue, ok := value.(basetypes.BoolValue)
		if ok {
			converted[UnderscoresToCamelCase(key)] = boolvalue.ValueBool()
			continue
		}
		int64value, ok := value.(basetypes.Int64Value)
		if ok {
			converted[UnderscoresToCamelCase(key)] = int64value.ValueInt64()
			continue
		}
		float64value, ok := value.(basetypes.Float64Value)
		if ok {
			converted[UnderscoresToCamelCase(key)] = float64value.ValueFloat64()
			continue
		}
		setvalue, ok := value.(basetypes.SetValue)
		if ok {
			converted[UnderscoresToCamelCase(key)] = ConvertToPrimitive(setvalue)
		}
	}

	return converted
}

func ConvertToPrimitive(value attr.Value) interface{} {
	// Handle primitives
	strvalue, ok := value.(basetypes.StringValue)
	if ok {
		return strvalue.ValueString()
	}
	boolvalue, ok := value.(basetypes.BoolValue)
	if ok {
		return boolvalue.ValueBool()
	}
	int64value, ok := value.(basetypes.Int64Value)
	if ok {
		return int64value.ValueInt64()
	}
	float64value, ok := value.(basetypes.Float64Value)
	if ok {
		return float64value.ValueFloat64()
	}

	// Handle lists and sets
	listvalue, ok := value.(basetypes.ListValue)
	if ok {
		elements := listvalue.Elements()
		var primitiveElements []interface{}
		for _, el := range elements {
			primitiveElements = append(primitiveElements, ConvertToPrimitive(el))
		}
		return primitiveElements
	}

	setvalue, ok := value.(basetypes.SetValue)
	if ok {
		elements := setvalue.Elements()
		var primitiveElements []interface{}
		for _, el := range elements {
			primitiveElements = append(primitiveElements, ConvertToPrimitive(el))
		}
		return primitiveElements
	}

	// Handle maps
	mapvalue, ok := value.(basetypes.MapValue)
	if ok {
		mapElements := mapvalue.Elements()
		primitiveMap := map[string]interface{}{}
		for key, el := range mapElements {
			primitiveMap[UnderscoresToCamelCase(key)] = ConvertToPrimitive(el)
		}
		return primitiveMap
	}

	// Handle objects
	objvalue, ok := value.(basetypes.ObjectValue)
	if ok {
		mapElements := objvalue.Attributes()
		primitiveMap := map[string]interface{}{}
		for key, el := range mapElements {
			primitiveMap[UnderscoresToCamelCase(key)] = ConvertToPrimitive(el)
		}
		return primitiveMap
	}

	return fmt.Errorf("unable to convert given primitive type for %s", value)
}

// Converts the map[string]attr.Type to basetypes.ObjectValue required for Terraform
func MaptoObjValue(attributeTypes map[string]attr.Type, attributeValues map[string]attr.Value, diags *diag.Diagnostics) basetypes.ObjectValue {
	newObj, err := types.ObjectValue(attributeTypes, attributeValues)
	if err != nil {
		diags.AddError("ERROR: ", "An error occured while converting ")
	}
	return newObj
}

func InterfaceStringValueOrNull(value interface{}) types.String {
	if value == nil {
		return basetypes.NewStringNull()
	} else {
		return types.StringValue(value.(string))
	}
}

func InterfaceFloatSetValue(values []interface{}) []float64 {
	newFloat := make([]float64, 0, len(values))
	for _, v := range newFloat {
		newFloat = append(newFloat, float64(v))
	}

	return newFloat
}

func CreateKeysFromAttrValues(attrValues map[string]attr.Value) []string {
	attrKeys := make([]string, 0, len(attrValues))
	for k := range attrValues {
		attrKeys = append(attrKeys, k)
	}
	return attrKeys
}

func CheckListKeyMatch(k string, list []string) bool {
	for _, key_check := range list {
		if k == key_check {
			return true
		}
	}
	return false
}

// Add a keyval pair to existing map[string]attr.Type, making a deep copy and not modifying the original
func AddKeyValToMapStringAttrType(mapStringAttrType map[string]attr.Type, key string, val attr.Type) map[string]attr.Type {
	outValue := make(map[string]attr.Type)
	for k, v := range mapStringAttrType {
		outValue[k] = v
	}
	outValue[key] = val
	return outValue
}
