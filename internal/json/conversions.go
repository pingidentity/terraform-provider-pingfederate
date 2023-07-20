package json

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func FromValue(value attr.Value) string {
	var jsonString strings.Builder

	// Simple types
	strvalue, ok := value.(basetypes.StringValue)
	if ok {
		jsonString.WriteRune('"')
		jsonString.WriteString(strvalue.ValueString())
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
		writeArray(listvalue.Elements(), &jsonString)
	}

	setvalue, ok := value.(basetypes.SetValue)
	if ok {
		writeArray(setvalue.Elements(), &jsonString)
	}

	// Maps and objects
	mapvalue, ok := value.(basetypes.MapValue)
	if ok {
		writeMap(mapvalue.Elements(), &jsonString)
	}

	objvalue, ok := value.(basetypes.ObjectValue)
	if ok {
		writeMap(objvalue.Attributes(), &jsonString)
	}
	return jsonString.String()
}

func writeArray(values []attr.Value, builder *strings.Builder) {
	builder.WriteRune('[')
	for i, attrValue := range values {
		builder.WriteString(FromValue(attrValue))
		if i < len(values)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
}

func writeMap(values map[string]attr.Value, builder *strings.Builder) {
	builder.WriteRune('{')
	i := 0
	for attrName, attrValue := range values {
		builder.WriteRune('"')
		builder.WriteString(underscoreToCamelCase(attrName))
		builder.WriteString("\":")
		builder.WriteString(FromValue(attrValue))
		if i < len(values)-1 {
			builder.WriteRune(',')
		}
		i++

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
