package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/sourcetypeidkey"
)

func ToDataSourceSchema() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Description: "Defines how an attribute in an attribute contract should be populated.",
		Optional:    false,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"source": sourcetypeidkey.ToDataSourceSchema(),
				"value": schema.StringAttribute{
					Optional:    false,
					Computed:    true,
					Description: "The value for this attribute.",
				},
			},
		},
	}
}
