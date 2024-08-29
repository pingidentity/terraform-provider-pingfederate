package id

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
)

func ToSchema(s *schema.Schema) {
	ToSchemaDeprecated(s, false)
}

func ToSchemaDeprecated(s *schema.Schema, deprecated bool) {
	schemaId := schema.StringAttribute{}
	schemaId.Description = "The ID of this resource."
	schemaId.Required = false
	schemaId.Optional = false
	schemaId.Computed = true
	schemaId.PlanModifiers = []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	}
	if deprecated {
		schemaId.DeprecationMessage = "This attribute is deprecated and will be removed in a future release."
	}
	s.Attributes["id"] = schemaId
}

func ToSchemaCustomId(s *schema.Schema, idName string, required, characterLimit bool, description string) {
	customId := schema.StringAttribute{}
	customId.Description = description
	customId.PlanModifiers = []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
		stringplanmodifier.RequiresReplace(),
	}
	customId.Required = required
	customId.Optional = !required
	customId.Computed = !required
	customId.Validators = []validator.String{
		stringvalidator.LengthAtLeast(1),
	}
	if characterLimit {
		customId.Validators = append(customId.Validators, configvalidators.ValidChars())
	}
	s.Attributes[idName] = customId
}
