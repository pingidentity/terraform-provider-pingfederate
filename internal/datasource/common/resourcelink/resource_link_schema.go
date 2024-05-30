package resourcelink

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Description: "The ID of the resource.",
			Optional:    false,
			Computed:    true,
		},
	}
}

func ToDataSourceSchemaSingleNestedAttribute() datasourceschema.SingleNestedAttribute {
	return datasourceschema.SingleNestedAttribute{
		Description: "A reference to a resource.",
		Computed:    true,
		Optional:    false,
		Attributes:  ToDataSourceSchema(),
	}
}

func ToDataSourceSchemaSingleNestedAttributeCustomDescription(description string) datasourceschema.SingleNestedAttribute {
	singleNestedAttribute := ToDataSourceSchemaSingleNestedAttribute()
	singleNestedAttribute.Description = description
	return singleNestedAttribute
}
