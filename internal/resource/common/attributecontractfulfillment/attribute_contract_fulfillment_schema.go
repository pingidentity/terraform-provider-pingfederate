package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func ToSchema(required, fullyComputed bool) schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Description: "Defines how an attribute in an attribute contract should be populated.",
		Required:    required,
		Optional:    !required,
		Computed:    !required,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"source": sourcetypeidkey.ToSchema(fullyComputed),
				"value": schema.StringAttribute{
					Optional:    true,
					Computed:    fullyComputed,
					Description: "The value for this attribute.",
				},
			},
		},
	}
}
