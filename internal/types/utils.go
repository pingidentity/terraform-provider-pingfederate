package types

import (
	"context"
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
	return !value.IsNull() && !value.IsUnknown()
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

func IsUrlFormat(s string) bool {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	//^ and $: Start and end of the string.
	// (https?:\/\/): Matches the http:// or https:// scheme. This part is optional
	// ([\da-z\.-]+): Matches the domain name, which can contain digits, alphabets (lowercase), dots, and hyphens.
	//       One or more such characters must be present.
	// (\.[a-z]{2,6}(\.[a-z]{2,6})?)?: Matches the optional TLD and sub-TLD, each containing 2 to 6 alphabets (lowercase).
	//       The entire TLD part is optional.
	// (:[0-9]{2,5})?: Matches the optional port number, which can be 2 to 5 digits long.
	// (\/[^\s\/]+[^\s]*)? matches the optional path
	//
	// The regex will match the following examples:
	// 	http://example.com
	//	https://example.com
	//	http://example.com:8080
	//	https://example.co.uk:8080
	//	https://localhost
	//	http://localhost:3000
	//	http://example.com/path/to/resource
	//  https://anotherexample.co.uk:9399/path/with/trailing/slash/
	//  localhost
	//  localhost:9999
	//
	//  Non-Matching examples
	//	http://example. (TLD ends with a dot)
	//	http://example.com:808080 (Port number exceeds 5 digits)
	//	http://example_com.com (Underscore is not allowed in domain name)
	//	ftp://example.com (Only http and https schemes are allowed if specified)

	re, _ := regexp.Compile(`^(https?:\/\/)?([\da-z\.-]+)(\.[a-z]{2,6}(\.[a-z]{2,6})?)?(:[0-9]{2,5})?(\/)?(\/[^\s]*)?$`)
	return re.MatchString(strings.ToLower(s))
}

func IsEmailFormat(s string) bool {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	// ^ and $: Start and end of the string.
	// [a-z0-9._%+-]+: Matches the local part of the email address, which can contain alphabetic characters, digits, dots, underscores, percent signs, pluses, and hyphens. One or more of these options must be present.
	// @: Matches the "@" symbol.
	// [a-z0-9.-]+: Matches the domain name, which can contain alphabetic characters, digits, dots, and hyphens. One or more of these options must be present.
	// \.: Matches the dot before the top-level domain (TLD).
	// [a-z]{2,}: Matches the TLD, which must contain at least two alphabetic characters.
	//
	// The regex will match the following example emails:
	// john.doe@example.com
	// Jane.Doe@sub.example.co
	// email+filter@gmail.com
	// email_filter@gmail.com
	// juan%sanchez@gmail.com
	// 99user@domain.com
	// user.name@domain.co.uk
	//
	// Non-Matching emails
	// @example.com (Local part is missing)
	// john.doe@.com (Domain name is missing)
	// john.doe@com (Dot before TLD is missing)
	// john.doe@domain. (TLD is missing)
	// john.doe@domain.c (TLD is too short)
	// john$doe@domain.com (invalid '$' character)
	// john..doe@example.com (Consecutive dots in the local part)
	// john.doe@.example.com (Domain starts with a dot)

	re, _ := regexp.Compile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(s))
}

func IsValidHostnameOrIp(s string) bool {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	// This implementation is not complete nor perfect, but will catch many invalid hostnames or IPs,
	// including invalid characters, whitespace, etc.
	// TODO:  expand for thoroughness and possible IPV6 validation
	// This portion matches hostnames:
	// ^([a-z0-9]+(-[a-z0-9]+)*\.)*[a-z0-9]+(-[a-z0-9]+)*$
	// ^ and $: Start and end of the string.
	// It allows for alphabetic characters, digits, and hyphens.
	// A hostname part cannot start or end with a hyphen.
	// Cannot start or end with a dot.
	// Hostname parts are separated by dots.

	re, _ := regexp.Compile(`^([a-z0-9]+(-[a-z0-9]+)*\.)*[a-z0-9]+(-[a-z0-9]+)*$`)
	return re.MatchString(strings.ToLower(s))
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

	panic("Panic reached. Unable to convert given primitive type.")
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
