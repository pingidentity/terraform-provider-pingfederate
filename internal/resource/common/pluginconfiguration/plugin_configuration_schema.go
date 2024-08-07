package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToSchema() schema.SingleNestedAttribute {
	fieldsSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: fieldAttrTypes}, nil)
	tablesSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: tableAttrTypes}, nil)
	fieldsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the configuration field.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
	tablesNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"rows": schema.ListNestedAttribute{
				Description: "List of table rows.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fields": schema.SetNestedAttribute{
							Description: "The configuration fields in the row.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the configuration field.",
										Required:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
										Required:    true,
										Sensitive:   true,
									},
								},
							},
						},
						"default_row": schema.BoolAttribute{
							Description: "Whether this row is the default.",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
	return schema.SingleNestedAttribute{
		Description: "Plugin instance configuration.",
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"tables": schema.SetNestedAttribute{
				Description:  "List of configuration tables.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(tablesSetDefault),
				NestedObject: tablesNestedObject,
			},
			"tables_all": schema.SetNestedAttribute{
				Description:  "List of configuration tables. This attribute will include any values set by default by PingFederate.",
				Computed:     true,
				Optional:     false,
				NestedObject: tablesNestedObject,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"fields": schema.SetNestedAttribute{
				Description:  "List of configuration fields.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(fieldsSetDefault),
				NestedObject: fieldsNestedObject,
			},
			"fields_all": schema.SetNestedAttribute{
				Description:  "List of configuration fields. This attribute will include any values set by default by PingFederate.",
				Computed:     true,
				Optional:     false,
				NestedObject: fieldsNestedObject,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
