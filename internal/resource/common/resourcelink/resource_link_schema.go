package resourcelink

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ToSchemaLocationUseStateForUnknown() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
	}
}

func ToSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
	}
}

func SingleNestedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "A reference to a resource.",
		Optional:    true,
		Attributes:  ToSchema(),
	}
}

func CompleteSingleNestedAttribute(optional, computed, required bool, description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    optional,
		Computed:    computed,
		Required:    required,
		Description: description,
		Attributes:  ToSchema(),
	}
}
