package id

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema(s *datasourceschema.Schema) {
	ToDataSourceSchemaCustomId(s, "id", false, "ID of this resource.")
}

func ToDataSourceSchemaDeprecated(s *datasourceschema.Schema, deprecated bool) {
	ToDataSourceSchemaCustomIdDeprecated(s, "id", false, "ID of this resource.", deprecated)
}

func ToDataSourceSchemaCustomId(s *datasourceschema.Schema, idName string, required bool, description string) {
	idAttr := datasourceschema.StringAttribute{
		Description: description,
		Computed:    !required,
		Optional:    false,
		Required:    required,
	}
	s.Attributes[idName] = idAttr
}

func ToDataSourceSchemaCustomIdDeprecated(s *datasourceschema.Schema, idName string, required bool, description string, deprecated bool) {
	idAttr := datasourceschema.StringAttribute{
		Description: description,
		Computed:    !required,
		Optional:    false,
		Required:    required,
	}
	if deprecated {
		idAttr.DeprecationMessage = "This attribute is deprecated and will be removed in a future release."
	}
	s.Attributes[idName] = idAttr
}
