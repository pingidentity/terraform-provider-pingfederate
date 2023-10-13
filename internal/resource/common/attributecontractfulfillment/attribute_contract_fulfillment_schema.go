package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func ToSchema(required bool) schema.MapNestedAttribute {
	attributeContractFulfillmentSchema := schema.MapNestedAttribute{}
	attributeContractFulfillmentSchema.Description = "Defines how an attribute in an attribute contract should be populated."
	attributeContractFulfillmentSchema.NestedObject.Attributes = map[string]schema.Attribute{
		"source": sourcetypeidkey.ToSchema(),
		"value": schema.StringAttribute{
			Optional:    true,
			Description: "The value for this attribute.",
		},
	}
	if required {
		attributeContractFulfillmentSchema.Required = true
	} else {
		attributeContractFulfillmentSchema.Optional = true
	}
	return attributeContractFulfillmentSchema
}
