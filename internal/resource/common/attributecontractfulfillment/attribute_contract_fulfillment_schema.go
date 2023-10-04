package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttributeContractFulfillmentSchema(required bool) schema.MapNestedAttribute {
	if required {
		return schema.MapNestedAttribute{
			Description: "Defines how an attribute in an attribute contract should be populated.",
			Required:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"source": sourcetypeidkey.SourceTypeIdKeySchema(),
					"value": schema.StringAttribute{
						Optional:    true,
						Description: "The value for this attribute.",
					},
				},
			},
		}
	} else {
		return schema.MapNestedAttribute{
			Description: "Defines how an attribute in an attribute contract should be populated.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"source": sourcetypeidkey.SourceTypeIdKeySchema(),
					"value": schema.StringAttribute{
						Optional:    true,
						Description: "The value for this attribute.",
					},
				},
			},
		}
	}
}
