package resourcelink

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
		"location": datasourceschema.StringAttribute{
			Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			Computed:    true,
			Optional:    false,
		},
	}
}
