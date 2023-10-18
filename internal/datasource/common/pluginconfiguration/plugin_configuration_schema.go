package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToDataSourceSchema() schema.SingleNestedAttribute {
	fieldsListDefault, _ := types.ListValue(types.ObjectType{AttrTypes: fieldAttrTypes}, []attr.Value{})
	return schema.SingleNestedAttribute{
		Description: "Plugin instance configuration.",
		Required:    false,
		Optional:    false,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"tables": schema.ListNestedAttribute{
				Description: "List of configuration tables.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the table.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"rows": schema.ListNestedAttribute{
							Description: "List of table rows.",
							Required:    false,
							Optional:    false,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"fields": schema.ListNestedAttribute{
										Description: "The configuration fields in the row.",
										Required:    false,
										Optional:    false,
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "The name of the configuration field.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
												"value": schema.StringAttribute{
													Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
												"inherited": schema.BoolAttribute{
													Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
											},
										},
									},
									"default_row": schema.BoolAttribute{
										Description: "Whether this row is the default.",
										Required:    false,
										Optional:    false,
										Computed:    true,
									},
								},
							},
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"fields": schema.ListNestedAttribute{
				Description: "List of configuration fields.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Default:     listdefault.StaticValue(fieldsListDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"fields_all": schema.ListNestedAttribute{
				Description: "List of configuration fields. This attribute will include any values set by default by PingFederate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
