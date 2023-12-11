package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringValuesSlice(values []string) []attr.Value {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(values[i])
	}
	return setValues
}

// Get a types.Set from a slice of string
func GetStringSet(values []string) types.Set {
	set, _ := types.SetValue(types.StringType, stringValuesSlice(values))
	return set
}

// Get a types.List from a slice of string
func GetStringList(values []string) types.List {
	list, _ := types.ListValue(types.StringType, stringValuesSlice(values))
	return list
}

func SetTypeToStringSlice(set types.Set) []string {
	values := make([]string, 0, len(set.Elements()))
	for _, v := range set.Elements() {
		values = append(values, v.(types.String).ValueString())
	}
	return values
}
