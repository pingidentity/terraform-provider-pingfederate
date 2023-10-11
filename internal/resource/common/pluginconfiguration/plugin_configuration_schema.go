package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Plugin instance configuration.",
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"tables": schema.ListNestedAttribute{
				Description: "List of configuration tables.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
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
									"fields": schema.ListNestedAttribute{
										Description: "The configuration fields in the row.",
										Computed:    true,
										Optional:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "The name of the configuration field.",
													Required:    true,
												},
												"value": schema.StringAttribute{
													Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
													//TODO diff
													Required: true,
												},
												"inherited": schema.BoolAttribute{
													Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
													Optional:    true,
													//TODO diff
													PlanModifiers: []planmodifier.Bool{
														boolplanmodifier.UseStateForUnknown(),
													},
												},
											},
										},
									},
									"default_row": schema.BoolAttribute{
										Description: "Whether this row is the default.",
										Optional:    true,
										//TODO diff
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
							Optional:    true,
						},
					},
				},
			},
			"fields": schema.ListNestedAttribute{
				Description: "List of configuration fields.",
				//TODO should this be computed? Usestateforunknown?
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
							Required:    true,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Optional:    true,
							//TODO diff
							/*PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},*/
						},
					},
				},
			},
			"fields_all": schema.ListNestedAttribute{
				Description: "List of configuration fields. This attribute will include any values set by default by PingFederate.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
							Required:    true,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}
