package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringValuesSlice(values []string) []attr.Value {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(string(values[i]))
	}
	return setValues
}

// Get a types.Set from a slice of string
// TODO
func GetStringSet(values []string) types.Set {
	set, _ := types.SetValue(types.StringType, stringValuesSlice(values))
	return set
}

// Get a types.List from a slice of string
// TODO
func GetStringList(values []string) types.List {
	list, _ := types.ListValue(types.StringType, stringValuesSlice(values))
	return list
}

// TODO
func SetTypeToStringSet(set types.Set) []string {
	values := make([]string, 0, len(set.Elements()))
	for _, v := range set.Elements() {
		values = append(values, v.(types.String).ValueString())
	}
	return values
}

// Get a types.String from the given string pointer, handling if the pointer is nil
/*func StringTypeOrNil(str *string, useEmptyStringForNil bool) types.String {
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
}*/
