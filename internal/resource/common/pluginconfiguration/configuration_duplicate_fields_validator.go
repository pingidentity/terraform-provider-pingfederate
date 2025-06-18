// Copyright © 2025 Ping Identity Corporation

package pluginconfiguration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.Object = &configurationDuplicateFieldsValidator{}

type configurationDuplicateFieldsValidator struct{}

// Configuration duplicate fields validator
func (v configurationDuplicateFieldsValidator) Description(ctx context.Context) string {
	return "Validates there are no duplicate fields in the user-defined configuration"
}

func (v configurationDuplicateFieldsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v configurationDuplicateFieldsValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	// Check fields and sensitive_fields
	fields, fieldsOk := req.ConfigValue.Attributes()["fields"]
	sensitiveFields, sensitiveFieldsOk := req.ConfigValue.Attributes()["sensitive_fields"]
	if fieldsOk && sensitiveFieldsOk {
		checkDuplicateFields(fields, sensitiveFields, "", -1, req, resp)
	}

	// Check tables.rows.fields and tables.rows.sensitive_fields
	tables, tablesOk := req.ConfigValue.Attributes()["tables"]
	if tablesOk {
		tablesObj, tablesOk := tables.(types.List)
		if tablesOk {
			for _, table := range tablesObj.Elements() {
				tableObj, tableOk := table.(types.Object)
				if tableOk {
					rows, rowsOk := tableObj.Attributes()["rows"]
					tableName, nameOk := tableObj.Attributes()["name"]
					if rowsOk && nameOk {
						rowsObj, rowsOk := rows.(types.List)
						tableNameStr, nameOk := tableName.(types.String)
						if rowsOk && nameOk {
							for rowIndex, row := range rowsObj.Elements() {
								rowObj, rowOk := row.(types.Object)
								if rowOk {
									fields, fieldsOk := rowObj.Attributes()["fields"]
									sensitiveFields, sensitiveFieldsOk := rowObj.Attributes()["sensitive_fields"]
									if fieldsOk && sensitiveFieldsOk {
										checkDuplicateFields(fields, sensitiveFields, tableNameStr.ValueString(), rowIndex, req, resp)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func checkDuplicateFields(fields attr.Value, sensitiveFields attr.Value, tableName string, rowIndex int, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	fieldsObj, fieldsOk := fields.(types.Set)
	sensitiveFieldsObj, sensitiveFieldsOk := sensitiveFields.(types.Set)
	errorMsgPrefix := "Duplicate field name in 'fields' and 'sensitive_fields': "
	if rowIndex != -1 {
		errorMsgPrefix = fmt.Sprintf("Duplicate field name in 'fields' and 'sensitive_fields' in table '%s' at row with index %d: ", tableName, rowIndex)
	}
	if fieldsOk && sensitiveFieldsOk {
		fieldNames := map[string]bool{}
		for _, field := range fieldsObj.Elements() {
			fieldObj, fieldOk := field.(types.Object)
			if fieldOk {
				fieldName, nameOk := fieldObj.Attributes()["name"]
				if nameOk {
					nameValue, nameOk := fieldName.(types.String)
					if nameOk && !nameValue.IsUnknown() {
						fieldNames[nameValue.ValueString()] = true
					}
				}
			}
		}
		for _, field := range sensitiveFieldsObj.Elements() {
			fieldObj, fieldOk := field.(types.Object)
			if fieldOk {
				fieldName, nameOk := fieldObj.Attributes()["name"]
				if nameOk {
					nameValue, nameOk := fieldName.(types.String)
					if nameOk && !nameValue.IsUnknown() {
						if _, ok := fieldNames[nameValue.ValueString()]; ok {
							resp.Diagnostics.AddAttributeError(
								req.Path,
								providererror.InvalidAttributeConfiguration,
								errorMsgPrefix+nameValue.ValueString(),
							)
						}
					}
				}
			}
		}
	}
}

func noDuplicateFields() configurationDuplicateFieldsValidator {
	return configurationDuplicateFieldsValidator{}
}
