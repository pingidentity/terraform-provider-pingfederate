// Copyright © 2026 Ping Identity Corporation

package sourcetypeidkey

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ToDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The attribute value source.",
		Optional:    false,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The source type of this key.",
				Optional:    false,
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
				Optional:    false,
				Computed:    true,
			},
		},
	}
}
