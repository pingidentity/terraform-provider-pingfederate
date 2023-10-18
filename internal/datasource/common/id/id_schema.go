package id

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema(s *datasourceschema.Schema, required bool, description string) {
	idSchemaAttr := datasourceschema.StringAttribute{}
	idSchemaAttr.Description = description
	if required {
		idSchemaAttr.Required = true
	} else {
		idSchemaAttr.Computed = true
		idSchemaAttr.Required = false
		idSchemaAttr.Optional = false
	}
	s.Attributes["id"] = idSchemaAttr
}
