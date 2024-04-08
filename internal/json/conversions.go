package json

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Convert a terraform plugin framework value into an equivalent JSON string
func FromValue(value attr.Value, skipNullOrUnknownAttrs bool) string {
	var jsonString strings.Builder

	// Simple types
	strvalue, ok := value.(basetypes.StringValue)
	if ok {
		jsonString.WriteRune('"')
		// Ensure any escaped quotes in the string are handled so that the resulting json includes a backslash
		jsonString.WriteString(strings.ReplaceAll(strvalue.ValueString(), "\"", "\\\""))
		jsonString.WriteRune('"')
	}
	boolvalue, ok := value.(basetypes.BoolValue)
	if ok {
		jsonString.WriteString(strconv.FormatBool(boolvalue.ValueBool()))
	}
	int64value, ok := value.(basetypes.Int64Value)
	if ok {
		jsonString.WriteString(strconv.FormatInt(int64value.ValueInt64(), 10))
	}
	float64value, ok := value.(basetypes.Float64Value)
	if ok {
		jsonString.WriteString(strconv.FormatFloat(float64value.ValueFloat64(), 'f', -1, 64))
	}

	// Lists and sets
	listvalue, ok := value.(basetypes.ListValue)
	if ok {
		writeArray(listvalue.Elements(), &jsonString, skipNullOrUnknownAttrs)
	}

	setvalue, ok := value.(basetypes.SetValue)
	if ok {
		writeArray(setvalue.Elements(), &jsonString, skipNullOrUnknownAttrs)
	}

	// Maps and objects
	mapvalue, ok := value.(basetypes.MapValue)
	if ok {
		writeMap(mapvalue.Elements(), &jsonString, skipNullOrUnknownAttrs, false)
	}

	objvalue, ok := value.(basetypes.ObjectValue)
	if ok {
		writeMap(objvalue.Attributes(), &jsonString, skipNullOrUnknownAttrs, true)
	}
	return jsonString.String()
}

func writeArray(values []attr.Value, builder *strings.Builder, skipNullOrUnknownAttrs bool) {
	builder.WriteRune('[')
	for i, attrValue := range values {
		builder.WriteString(FromValue(attrValue, skipNullOrUnknownAttrs))
		if i < len(values)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
}

func writeMap(values map[string]attr.Value, builder *strings.Builder, skipNullOrUnknownAttrs bool, excludeUnderscore bool) {
	builder.WriteRune('{')
	isFirst := true
	for attrName, attrValue := range values {
		if skipNullOrUnknownAttrs && !internaltypes.IsDefined(attrValue) {
			continue
		}
		if !isFirst {
			builder.WriteRune(',')
		} else {
			isFirst = false
		}
		builder.WriteRune('"')
		if excludeUnderscore {
			builder.WriteString(underscoreToCamelCase(attrName))
		} else {
			builder.WriteString(attrName)
		}
		builder.WriteString("\":")
		builder.WriteString(FromValue(attrValue, skipNullOrUnknownAttrs))
	}
	builder.WriteRune('}')
}

func underscoreToCamelCase(value string) string {
	var result strings.Builder
	upperCase := false
	for _, char := range value {
		if char == '_' {
			upperCase = true
		} else {
			if upperCase {
				result.WriteString(strings.ToUpper(string(char)))
			} else {
				result.WriteRune(char)
			}
			upperCase = false
		}
	}
	return result.String()
}
