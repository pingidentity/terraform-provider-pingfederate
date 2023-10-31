package id

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
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

func ToSchema(s *datasourceschema.Schema) {
	schemaId := schema.StringAttribute{}
	schemaId.Description = "The ID of this resource."
	schemaId.Required = false
	schemaId.Optional = false
	schemaId.Computed = true
	schemaId.PlanModifiers = []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	}
	s.Attributes["id"] = schemaId
}

func ToDataSourceSchemaCustomId(s *datasourceschema.Schema, required bool, characterLimit bool, description string) {
	customId := schema.StringAttribute{}
	customId.Description = description
	if required {
		customId.Required = true
	} else {
		customId.Computed = true
		customId.Optional = true
	}
	if characterLimit {
		customId.Validators = []validator.String{configvalidators.ValidChars()}
	}
	s.Attributes["custom_id"] = customId
}
