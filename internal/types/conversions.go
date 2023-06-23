package types

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func StringToTF(v string) basetypes.StringValue {
	if v == "" {
		return types.StringNull()
	} else {
		return types.StringValue(v)
	}
}

func StringToInt64Pointer(value basetypes.StringValue) *int64 {
	valueToString := value.ValueString()
	newVal, _ := strconv.ParseInt(valueToString, 10, 64)
	return &newVal
}

func BaseTypesInt64ToString(value basetypes.Int64Value) string {
	return strconv.FormatInt(value.ValueInt64(), 10)
}

func Int64PointerToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Int64ToString(value types.Int64) string {
	return strconv.FormatInt(value.ValueInt64(), 10)
}

func StringInterfaceToStringOrNull(v interface{}) basetypes.StringValue {
	stringValue, ok := v.(string)
	if ok {
		return types.StringValue(stringValue)
	} else {
		return types.StringNull()
	}
}

func BoolInterfaceToBoolOrNull(v interface{}) basetypes.BoolValue {
	boolValue, ok := v.(bool)
	if ok {
		return types.BoolValue(boolValue)
	} else {
		return types.BoolNull()
	}
}

func Int64InterfaceToInt64OrNull(v interface{}) basetypes.Int64Value {
	int64Value, ok := v.(int64)
	if ok {
		return types.Int64Value(int64Value)
	} else {
		return types.Int64Null()
	}
}

func Float64InterfaceToFloat64OrNull(v interface{}) basetypes.Float64Value {
	float64Value, ok := v.(float64)
	if ok {
		return types.Float64Value(float64(float64Value))
	} else {
		return types.Float64Null()
	}
}

// Get a types.Set from a slice of string
func GetStringSet(values []string) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(string(values[i]))
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

func InterfaceToStringSet(i interface{}) types.Set {
	values := i.([]string)
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(string(values[i]))
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

// Get a types.Set from a slice of string
func GetInterfaceStringSet(i interface{}) types.Set {
	values := i.([]interface{})
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(string(values[i].(string)))
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

// Get a types.Set from a slice of int64
func GetInt64Set(values []int64) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.Int64Value(int64(values[i]))
	}
	set, _ := types.SetValue(types.Int64Type, setValues)
	return set
}

// Get a types.Set from a slice of int64 or null set
func GetInt64SetOrNull(values []int64) types.Set {
	if len(values) >= 1 {
		setValues := make([]attr.Value, len(values))
		for i := 0; i < len(values); i++ {
			setValues[i] = types.Int64Value(int64(values[i]))
		}
		set, _ := types.SetValue(types.Int64Type, setValues)
		return set
	} else {
		return types.SetNull(types.Int64Type)
	}
}

// Get a types.Set from a slice of int64 or null set
func GetInt64InterfaceSetOrNull(i interface{}) types.Set {
	values := i.([]int64)
	if len(values) >= 1 {
		setValues := make([]attr.Value, len(values))
		for i := 0; i < len(values); i++ {
			setValues[i] = types.Int64Value(int64(values[i]))
		}
		set, _ := types.SetValue(types.Int64Type, setValues)
		return set
	} else {
		return types.SetNull(types.Int64Type)
	}
}

// Get a types.Set from a slice of float64 or null set
func GetFloat64InterfaceSetOrNull(i interface{}) types.Set {
	values := i.([]interface{})
	if i != nil && len(values) >= 1 {
		setValues := make([]attr.Value, len(values))
		for i := 0; i < len(values); i++ {
			setValues[i] = types.Float64Value(float64(values[i].(float64)))
		}
		set, _ := types.SetValue(types.Float64Type, setValues)
		return set
	} else {
		return types.SetNull(types.Float64Type)
	}
}

// Get a types.Set from a slice of strings and converts values to floats or null set
func StringInterfaceSetToFloat64SetOrNull(i interface{}) types.Set {
	values := i.([]interface{})
	if i != nil && len(values) >= 1 {
		setValues := make([]attr.Value, len(values))
		for i := 0; i < len(values); i++ {
			valueToString, _ := strconv.ParseFloat(values[i].(string), 64)
			setValues[i] = types.Float64Value(float64(valueToString))
		}
		set, _ := types.SetValue(types.Float64Type, setValues)
		return set
	} else {
		return types.SetNull(types.Float64Type)
	}
}

// Get a types.String from the given string pointer, handling if the pointer is nil
func StringTypeOrNil(str *string, useEmptyStringForNil bool) types.String {
	if str == nil {
		// If a plan was provided and is using an empty string, we should use that for a nil string in the response.
		// For PingFederate nil and empty string is equivalent, but to Terraform they are distinct. So we
		// just want to match whatever is in the plan when we get a nil string back.
		if useEmptyStringForNil {
			// Use empty string instead of null to match the plan when resetting string properties.
			// This is useful for computed values being reset to null.
			return types.StringValue("")
		} else {
			return types.StringNull()
		}
	}
	return types.StringValue(*str)
}

// Get a types.String from the given interface pointer, handling if the interface is nil
func InterfaceStringOrNil(i interface{}) string {
	if i == nil {
		return ""
	}

	return i.(string)
}

// Get a nested key value from given interface, handling if the value is nil
func GetNestedInterfaceKeyStringValue(i interface{}, nestedKey string) types.String {
	if i != nil && nestedKey != "" {
		return InterfaceStringValueOrNull(i.(map[string]interface{})[nestedKey])
	} else {
		return types.StringNull()
	}
}

// Get a nested key value from given interface, handling if the value is nil
func GetNestedInterfaceKeyBoolValuePointer(i interface{}, nestedKey string) *bool {
	boolPointer := (i.(map[string]interface{})[nestedKey]).(bool)
	return &boolPointer
}

// Get a types.Bool from the given bool pointer, handling if the pointer is nil
func BoolTypeOrNil(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// Get a types.Bool from the given interface pointer, handling if the interface is nil
func InterfaceBoolTypeOrNull(i interface{}) types.Bool {
	if i == nil {
		return types.BoolNull()
	}

	return types.BoolValue(i.(bool))
}

// Get a types.Bool from the given interface pointer, handling if the interface is nil
func InterfaceBoolPointerValue(i interface{}) *bool {
	value := i.(bool)
	return &value
}

// Get a string pointer from the given interface, handling if the interface is nil
func InterfaceStringPointerValue(i interface{}) *string {
	if i != nil {
		value := i.(string)
		return &value
	}
	return basetypes.NewStringNull().ValueStringPointer()
}

// Get a types.Int64 from the given int32 pointer, handling if the pointer is nil
func Int64TypeOrNil(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

// Get a types.Int64 from the given interface, handling if the pointer is nil
func Int64InterfaceTypeOrNil(i interface{}) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(i.(int64)))
}

// Get a types.Int64 from the given interface, handling if the pointer is nil
func InterfaceFloat64TypeOrNull(i interface{}) types.Float64 {
	if i == nil {
		return types.Float64Null()
	}

	return types.Float64Value(i.(float64))
}

// Get a types.Float64 from the given float32 pointer, handling if the pointer is nil
func Float64TypeOrNil(f *float32) types.Float64 {
	if f == nil {
		return types.Float64Null()
	}

	return types.Float64Value(float64(*f))
}

// Get types.Map from slice of Strings
func GetStringMap(m *string) types.Map {
	setValues := make(map[string]attr.Value)
	set, _ := types.MapValue(types.StringType, setValues)
	return set
}
