package resourcelink

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ToSchemaLocationUseStateForUnknown() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
		"location": schema.StringAttribute{
			DeprecationMessage: "This field is now deprecated and will be removed in a future release.",
			Description:        "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			Computed:           true,
			Optional:           false,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

func ToSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
		"location": schema.StringAttribute{
			DeprecationMessage: "This field is now deprecated and will be removed in a future release.",
			Description:        "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			Computed:           true,
			Optional:           false,
		},
	}
}

func ToSchemaNoLocation() map[string]schema.Attribute {
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
