// Copyright © 2025 Ping Identity Corporation

package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func IsNonEmptyObj(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown() && obj.Attributes() != nil
}

// Return true if this types.String represents a non-empty, non-null, non-unknown string
func IsNonEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() != ""
}

// Return true if this value represents a defined (non-null and non-unknown) value
func IsDefined(value attr.Value) bool {
	return value != nil && !value.IsNull() && !value.IsUnknown()
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
func StringSlicesEqual(a, b []string) bool {
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
func FloatSlicesEqual(a, b []float64) bool {
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

func UnderscoresToCamelCase(s string) string {
	re, _ := regexp.Compile(`(_[A-Za-z])`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	return res
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
	objvalue, ok := value.(types.Object)
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

// Add a keyval pair to existing map[string]attr.Type, making a deep copy and not modifying the original
func AddKeyValToMapStringAttrType(mapStringAttrType map[string]attr.Type, key string, val attr.Type) map[string]attr.Type {
	outValue := make(map[string]attr.Type)
	for k, v := range mapStringAttrType {
		outValue[k] = v
	}
	outValue[key] = val
	return outValue
}
