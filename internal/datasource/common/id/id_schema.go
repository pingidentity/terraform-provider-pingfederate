// Copyright © 2025 Ping Identity Corporation

package id

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema(s *datasourceschema.Schema) {
	ToDataSourceSchemaCustomId(s, "id", false, "ID of this resource.")
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
