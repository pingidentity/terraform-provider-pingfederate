package idp

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdapterResource{}
	_ resource.ResourceWithConfigure   = &idpAdapterResource{}
	_ resource.ResourceWithImportState = &idpAdapterResource{}
)

// Define attribute types for object types
var (
	fieldsAttrTypes = map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
		"inherited":       types.BoolType,
	}
	tablesAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"fields": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: fieldsAttrTypes,
						},
					},
					"default_row": types.BoolType,
				},
			},
		},
		"inherited": types.BoolType,
	}
	configurationAttrTypes = map[string]attr.Type{
		"tables": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: tablesAttrTypes,
			},
		},
		"fields": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: fieldsAttrTypes,
			},
		},
	}
	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":      types.StringType,
					"pseudonym": types.BoolType,
					"masked":    types.BoolType,
				},
			},
		},
		"extended_attributes": types.SetType{
			ElemType: types.ObjectType{
				//TODO more duplication
				AttrTypes: map[string]attr.Type{
					"name":      types.StringType,
					"pseudonym": types.BoolType,
					"masked":    types.BoolType,
				},
			},
		},
		"unique_user_key_attribute": types.StringType,
		"mask_ognl_values":          types.BoolType,
		"inherited":                 types.BoolType,
	}

	attributeContractFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
			},
		},
	}

	customAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: internaltypes.ResourceLinkStateAttrType(),
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"filter_fields": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"value": types.StringType,
					"name":  types.StringType,
				},
			},
		},
	}

	jdbcAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: internaltypes.ResourceLinkStateAttrType(),
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"schema": types.StringType,
		"table":  types.StringType,
		"column_names": types.ListType{
			ElemType: types.StringType,
		},
		"filter": types.StringType,
	}

	ldapAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: internaltypes.ResourceLinkStateAttrType(),
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"search_filter":          types.StringType,
		"search_scope":           types.StringType,
		"member_of_nested_group": types.BoolType,
		"base_dn":                types.StringType,
		"binary_attribute_settings": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"binary_encoding": types.StringType,
				},
			},
		},
		"search_attributes": types.ListType{
			ElemType: types.StringType,
		},
	}

	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"custom_attribute_source": types.ObjectType{
						AttrTypes: customAttrSourceAttrTypes,
					},
					"jdbc_attribute_source": types.ObjectType{
						AttrTypes: jdbcAttrSourceAttrTypes,
					},
					"ldap_attribute_source": types.ObjectType{
						AttrTypes: ldapAttrSourceAttrTypes,
					},
				},
			},
		},
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"source": types.ObjectType{
						//TODO remove this duplication?
						AttrTypes: map[string]attr.Type{
							"type": types.StringType,
							"id":   types.StringType,
						},
					},
					"value": types.StringType,
				},
			},
		},
		"issuance_criteria": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"conditional_criteria": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"source": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type": types.StringType,
									"id":   types.StringType,
								},
							},
							"attribute_name": types.StringType,
							"condition":      types.StringType,
							"value":          types.StringType,
							"error_result":   types.StringType,
						},
					},
				},
				"expression_criteria": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"expression":   types.StringType,
							"error_result": types.StringType,
						},
					},
				},
			},
		},
		"inherited": types.BoolType,
	}
)

// IdpAdapterResource is a helper function to simplify the provider implementation.
func IdpAdapterResource() resource.Resource {
	return &idpAdapterResource{}
}

// idpAdapterResource is the resource implementation.
type idpAdapterResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpAdapterResourceModel struct {
	AuthnCtxClassRef    types.String `tfsdk:"authn_ctx_class_ref"`
	Id                  types.String `tfsdk:"id"`
	CustomId            types.String `tfsdk:"custom_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
	AttributeMapping    types.Object `tfsdk:"attribute_mapping"`
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
}

// GetSchema defines the schema for the resource.
func (r *idpAdapterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Idp Adapter",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Optional:    true,
			},
			"custom_id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},

			"configuration": schema.SingleNestedAttribute{
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
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Description: "The name of the configuration field.",
															Required:    true,
														},
														"value": schema.StringAttribute{
															Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
															Optional:    true,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Optional:    true,
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Optional:    true,
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
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Optional:    true,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Optional:    true,
								},
							},
						},
					},
				},
			},

			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the IdP adapter provides.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"unique_user_key_attribute": schema.StringAttribute{
						Description: "The attribute to use for uniquely identify a user's authentication sessions.",
						Optional:    true,
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to false.",
						Optional:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
					},
				},
			},

			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					//TODO add attribute_sources
					"attribute_sources": schema.ListNestedAttribute{
						Optional:    true,
						Description: "A list of configured data stores to look up attributes from.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"custom_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from an associated data store.",
									Attributes: map[string]schema.Attribute{
										//TODO only need type on ldap dat source, others are implicit
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
												},
											},
										},
										"filter_fields": schema.ListNestedAttribute{
											Description: "The list of fields that can be used to filter a request to the custom data store.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"value": schema.StringAttribute{
														Description: "The value of this field. Whether or not the value is required will be determined by plugin validation checks.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "The name of this field.",
														Required:    true,
													},
												},
											},
										},
									},
								},
								"jdbc_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from a JDBC data store.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
												},
											},
										},
										"schema": schema.StringAttribute{
											Description: "Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.",
											Optional:    true,
										},
										"table": schema.StringAttribute{
											Description: "The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.",
											Required:    true,
										},
										"column_names": schema.ListAttribute{
											Description: "A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.",
											ElementType: types.StringType,
											Optional:    true,
										},
										"filter": schema.StringAttribute{
											Description: "The JDBC WHERE clause used to query your data store to locate a user record.",
											Required:    true,
										},
									},
								},
								"ldap_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from a LDAP data store.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
												},
											},
										},
										"search_filter": schema.StringAttribute{
											Description: "The LDAP filter that will be used to lookup the objects from the directory.",
											Required:    true,
										},
										"search_scope": schema.StringAttribute{
											Description: "Determines the node depth of the query.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf("OBJECT", "ONE_LEVEL", "SUBTREE"),
											},
										},
										"member_of_nested_group": schema.BoolAttribute{
											Description: "Set this to true to return transitive group memberships for the 'memberOf' attribute. This only applies for Active Directory data sources. All other data sources will be set to false.",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Bool{
												boolplanmodifier.UseStateForUnknown(),
											},
											Default: booldefault.StaticBool(false),
										},
										"base_dn": schema.StringAttribute{
											Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"binary_attribute_settings": schema.MapNestedAttribute{
											Description: "The advanced settings for binary LDAP attributes.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"binary_encoding": schema.StringAttribute{
														Optional:    true,
														Description: "Get the encoding type for this attribute. If not specified, the default is BASE64.",
														Validators: []validator.String{
															stringvalidator.OneOf("OBJECT", "ONE_LEVEL", "SUBTREE"),
														},
													},
												},
											},
										},
										"search_attributes": schema.ListAttribute{
											Description: "A list of LDAP attributes returned from search and available for mapping.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
							//TODO get these validators working
							/*Validators: []validator.Object{
								objectvalidator.ExactlyOneOf(
									path.MatchRelative().AtName("ldap_attribute_source"),
									path.MatchRelative().AtName("jdbc_attribute_source"),
									path.MatchRelative().AtName("custom_attribute_source")),
							},*/
						},
					},
					"attribute_contract_fulfillment": schema.MapNestedAttribute{
						Description: "A list of mappings from attribute names to their fulfillment values.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.SingleNestedAttribute{
									Description: "The attribute value source.",
									Required:    true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The source type of this key.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
											},
										},
										"id": schema.StringAttribute{
											Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
											Optional:    true,
										},
									},
								},
								"value": schema.StringAttribute{
									Description: "The value for this attribute.",
									Required:    true,
								},
							},
						},
					},
					"issuance_criteria": schema.SingleNestedAttribute{
						Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"conditional_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that checks a source attribute against a particular condition and the expected value. If the condition is true then this issuance criterion passes, otherwise the criterion fails.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										//TODO any way to share these definitions
										"source": schema.SingleNestedAttribute{
											Description: "The attribute value source.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Description: "The source type of this key.",
													Required:    true,
													Validators: []validator.String{
														stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
													},
												},
												"id": schema.StringAttribute{
													Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													Optional:    true,
												},
											},
										},
										"attribute_name": schema.StringAttribute{
											Description: "The name of the attribute to use in this issuance criterion.",
											Required:    true,
										},
										"condition": schema.StringAttribute{
											Description: "The condition that will be applied to the source attribute's value and the expected value.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"EQUALS", "EQUALS_CASE_INSENSITIVE", "EQUALS_DN", "NOT_EQUAL", "NOT_EQUAL_CASE_INSENSITIVE", "NOT_EQUAL_DN", "MULTIVALUE_CONTAINS", "MULTIVALUE_CONTAINS_CASE_INSENSITIVE", "MULTIVALUE_CONTAINS_DN", "MULTIVALUE_DOES_NOT_CONTAIN", "MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE", "MULTIVALUE_DOES_NOT_CONTAIN_DN"}...),
											},
										},
										"value": schema.StringAttribute{
											Description: "The expected value of this issuance criterion.",
											Required:    true,
										},
										"error_result": schema.StringAttribute{
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
											Optional:    true,
										},
									},
								},
							},
							"expression_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that uses a Boolean return value from an OGNL expression to determine whether or not it passes.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"expression": schema.StringAttribute{
											Required:    true,
											Description: "The OGNL expression to evaluate.",
										},
										"error_result": schema.StringAttribute{
											Optional:    true,
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
										},
									},
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
					},
				},
			},
		},
	}

	config.AddCommonSchema(&schema)
	resp.Schema = schema
}

func addOptionalIdpAdapterFields(ctx context.Context, addRequest *client.IdpAdapter, plan idpAdapterResourceModel) error {
	if internaltypes.IsDefined(plan.AuthnCtxClassRef) {
		addRequest.AuthnCtxClassRef = plan.AuthnCtxClassRef.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ParentRef) {
		addRequest.ParentRef = &client.ResourceLink{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef, false)), addRequest.ParentRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeMapping) {
		addRequest.AttributeMapping = &client.IdpAdapterContractMapping{}
		planAttrs := plan.AttributeMapping.Attributes()

		addRequest.AttributeMapping.Inherited = planAttrs["inherited"].(types.Bool).ValueBoolPointer()

		attrContractFulfillmentAttr := planAttrs["attribute_contract_fulfillment"].(types.Map)
		err := json.Unmarshal([]byte(internaljson.FromValue(attrContractFulfillmentAttr, true)), &addRequest.AttributeMapping.AttributeContractFulfillment)
		if err != nil {
			return err
		}

		issuanceCriteriaAttr := planAttrs["issuance_criteria"].(types.Object)
		addRequest.AttributeMapping.IssuanceCriteria = client.NewIssuanceCriteria()
		err = json.Unmarshal([]byte(internaljson.FromValue(issuanceCriteriaAttr, true)), addRequest.AttributeMapping.IssuanceCriteria)
		if err != nil {
			return err
		}

		attributeSourcesAttr := planAttrs["attribute_sources"].(types.List)
		addRequest.AttributeMapping.AttributeSources = []client.AttributeSourceAggregation{}
		for _, source := range attributeSourcesAttr.Elements() {
			//Determine which attribute source type this is
			sourceAttrs := source.(types.Object).Attributes()
			attributeSourceInner := client.AttributeSourceAggregation{}
			if internaltypes.IsDefined(sourceAttrs["custom_attribute_source"]) {
				attributeSourceInner.CustomAttributeSource = &client.CustomAttributeSource{}
				err = json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["custom_attribute_source"], true)), attributeSourceInner.CustomAttributeSource)
			}
			if internaltypes.IsDefined(sourceAttrs["jdbc_attribute_source"]) {
				attributeSourceInner.JdbcAttributeSource = &client.JdbcAttributeSource{}
				err = json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["jdbc_attribute_source"], true)), attributeSourceInner.JdbcAttributeSource)
			}
			if internaltypes.IsDefined(sourceAttrs["ldap_attribute_source"]) {
				attributeSourceInner.LdapAttributeSource = &client.LdapAttributeSource{}
				err = json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["ldap_attribute_source"], true)), attributeSourceInner.LdapAttributeSource)
			}
			if err != nil {
				return err
			}
			addRequest.AttributeMapping.AttributeSources = append(addRequest.AttributeMapping.AttributeSources, attributeSourceInner)
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = &client.IdpAdapterAttributeContract{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	return nil
}

// Metadata returns the resource type name.
func (r *idpAdapterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_adapters"
}

func (r *idpAdapterResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readIdpAdapterResponse(ctx context.Context, r *client.IdpAdapter, state *idpAdapterResourceModel, plan idpAdapterResourceModel, diags *diag.Diagnostics) {
	state.AuthnCtxClassRef = internaltypes.StringTypeOrNil(r.AuthnCtxClassRef, false)
	state.CustomId = types.StringValue(r.Id)
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef = internaltypes.ToStateResourceLink(ctx, &r.PluginDescriptorRef, diags)
	state.ParentRef = internaltypes.ToStateResourceLink(ctx, r.ParentRef, diags)

	var valueFromDiags diag.Diagnostics

	// Configuration
	//TODO move into common method
	//TODO unify how we handle diagnostics in these methods
	configurationAttrType := map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: config.FieldAttrTypes()}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: config.TableAttrTypes()}},
	}

	planFields := types.ListNull(types.ObjectType{AttrTypes: config.FieldAttrTypes()})
	planTables := types.ListNull(types.ObjectType{AttrTypes: config.TableAttrTypes()})

	planFieldsValue, ok := plan.Configuration.Attributes()["fields"]
	if ok {
		planFields = planFieldsValue.(types.List)
	}
	planTablesValue, ok := plan.Configuration.Attributes()["tables"]
	if ok {
		planTables = planTablesValue.(types.List)
	}

	fieldsAttrValue := config.ToFieldsListValue(r.Configuration.Fields, planFields, diags)
	tablesAttrValue := config.ToTablesListValue(r.Configuration.Tables, planTables, diags)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}
	state.Configuration, valueFromDiags = types.ObjectValue(configurationAttrType, configurationAttrValue)
	diags.Append(valueFromDiags...)

	if r.AttributeContract != nil {
		state.AttributeContract, valueFromDiags = types.ObjectValueFrom(ctx, attributeContractAttrTypes, r.AttributeContract)
		diags.Append(valueFromDiags...)
	}

	if r.AttributeMapping != nil {
		attributeMappingValues := map[string]attr.Value{
			"inherited": types.BoolPointerValue(r.AttributeMapping.Inherited),
		}

		// Build attribute_contract_fulfillment value
		attributeContractFulfillmentElementAttrTypes := attributeMappingAttrTypes["attribute_contract_fulfillment"].(types.MapType).ElemType.(types.ObjectType).AttrTypes
		attributeMappingValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(ctx,
			types.ObjectType{AttrTypes: attributeContractFulfillmentElementAttrTypes}, r.AttributeMapping.AttributeContractFulfillment)
		diags.Append(valueFromDiags...)

		// Build issuance_criteria value
		issuanceCritieraAttrTypes := attributeMappingAttrTypes["issuance_criteria"].(types.ObjectType).AttrTypes
		if r.AttributeMapping.IssuanceCriteria != nil {
			attributeMappingValues["issuance_criteria"], valueFromDiags = types.ObjectValueFrom(ctx,
				issuanceCritieraAttrTypes, r.AttributeMapping.IssuanceCriteria)
			diags.Append(valueFromDiags...)
		} else {
			attributeMappingValues["issuance_criteria"] = types.ObjectNull(issuanceCritieraAttrTypes)
		}

		// Build attribute_sources value
		attributeSourcesElementAttrTypes := attributeMappingAttrTypes["attribute_sources"].(types.ListType).ElemType.(types.ObjectType).AttrTypes
		if internaltypes.IsDefined(plan.AttributeMapping) && !internaltypes.IsDefined(plan.AttributeMapping.Attributes()["attribute_sources"]) {
			// don't return empty list if plan didn't specify any attribute sources, return null list
			attributeMappingValues["attribute_sources"] = types.ListNull(types.ObjectType{AttrTypes: attributeSourcesElementAttrTypes})
		} else {
			attrSourceElements := []attr.Value{}
			// This is assuming there won't be any default attribute sources returned by PF and that they will be returned in the same order
			planAttrSources := plan.AttributeMapping.Attributes()["attribute_sources"].(types.List).Elements()
			for i, attrSource := range r.AttributeMapping.AttributeSources {
				attrSourceValues := map[string]attr.Value{}
				if attrSource.CustomAttributeSource != nil {
					customAttrSourceValues := map[string]attr.Value{}
					customAttrSourceValues["filter_fields"], valueFromDiags = types.ListValueFrom(ctx,
						customAttrSourceAttrTypes["filter_fields"].(types.ListType).ElemType, attrSource.CustomAttributeSource.FilterFields)
					diags.Append(valueFromDiags...)

					customAttrSourceValues["type"] = types.StringValue(attrSource.CustomAttributeSource.Type)
					customAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.CustomAttributeSource.DataStoreRef)
					diags.Append(valueFromDiags...)
					customAttrSourceValues["id"] = types.StringPointerValue(attrSource.CustomAttributeSource.Id)
					customAttrSourceValues["description"] = types.StringPointerValue(attrSource.CustomAttributeSource.Description)
					customAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(ctx,
						types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.CustomAttributeSource.AttributeContractFulfillment)
					diags.Append(valueFromDiags...)
					attrSourceValues["custom_attribute_source"], valueFromDiags = types.ObjectValue(customAttrSourceAttrTypes, customAttrSourceValues)
					diags.Append(valueFromDiags...)
				} else {
					attrSourceValues["custom_attribute_source"] = types.ObjectNull(customAttrSourceAttrTypes)
				}
				if attrSource.JdbcAttributeSource != nil {
					jdbcAttrSourceValues := map[string]attr.Value{}
					jdbcAttrSourceValues["schema"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Schema)
					jdbcAttrSourceValues["table"] = types.StringValue(attrSource.JdbcAttributeSource.Table)
					jdbcAttrSourceValues["column_names"], valueFromDiags = types.ListValueFrom(ctx, types.StringType, attrSource.JdbcAttributeSource.ColumnNames)
					diags.Append(valueFromDiags...)
					jdbcAttrSourceValues["filter"] = types.StringValue(attrSource.JdbcAttributeSource.Filter)

					jdbcAttrSourceValues["type"] = types.StringValue(attrSource.JdbcAttributeSource.Type)
					jdbcAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.JdbcAttributeSource.DataStoreRef)
					diags.Append(valueFromDiags...)
					jdbcAttrSourceValues["id"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Id)
					jdbcAttrSourceValues["description"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Description)
					jdbcAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(ctx,
						types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.JdbcAttributeSource.AttributeContractFulfillment)
					diags.Append(valueFromDiags...)
					attrSourceValues["jdbc_attribute_source"], valueFromDiags = types.ObjectValue(jdbcAttrSourceAttrTypes, jdbcAttrSourceValues)
					diags.Append(valueFromDiags...)
				} else {
					attrSourceValues["jdbc_attribute_source"] = types.ObjectNull(jdbcAttrSourceAttrTypes)
				}
				if attrSource.LdapAttributeSource != nil {
					ldapAttrSourceValues := map[string]attr.Value{}
					ldapAttrSourceValues["base_dn"] = types.StringPointerValue(attrSource.LdapAttributeSource.BaseDn)
					ldapAttrSourceValues["search_scope"] = types.StringValue(attrSource.LdapAttributeSource.SearchScope)
					ldapAttrSourceValues["search_filter"] = types.StringValue(attrSource.LdapAttributeSource.SearchFilter)
					ldapAttrSourceValues["search_attributes"], valueFromDiags = types.ListValueFrom(ctx, types.StringType, attrSource.LdapAttributeSource.SearchAttributes)
					diags.Append(valueFromDiags...)
					if attrSource.LdapAttributeSource.BinaryAttributeSettings == nil ||
						!internaltypes.IsDefined(planAttrSources[i].(types.Object).Attributes()["ldap_attribute_source"].(types.Object).Attributes()["binary_attribute_settings"]) {

						ldapAttrSourceValues["binary_attribute_settings"] = types.MapNull(ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType)
					} else {
						ldapAttrSourceValues["binary_attribute_settings"], valueFromDiags = types.MapValueFrom(ctx,
							ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType, attrSource.LdapAttributeSource.BinaryAttributeSettings)
						diags.Append(valueFromDiags...)
					}
					ldapAttrSourceValues["member_of_nested_group"] = types.BoolPointerValue(attrSource.LdapAttributeSource.MemberOfNestedGroup)

					ldapAttrSourceValues["type"] = types.StringValue(attrSource.LdapAttributeSource.Type)
					ldapAttrSourceValues["data_store_ref"], valueFromDiags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.LdapAttributeSource.DataStoreRef)
					diags.Append(valueFromDiags...)
					ldapAttrSourceValues["id"] = types.StringPointerValue(attrSource.LdapAttributeSource.Id)
					ldapAttrSourceValues["description"] = types.StringPointerValue(attrSource.LdapAttributeSource.Description)
					ldapAttrSourceValues["attribute_contract_fulfillment"], valueFromDiags = types.MapValueFrom(ctx,
						types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.LdapAttributeSource.AttributeContractFulfillment)
					diags.Append(valueFromDiags...)
					attrSourceValues["ldap_attribute_source"], valueFromDiags = types.ObjectValue(ldapAttrSourceAttrTypes, ldapAttrSourceValues)
					diags.Append(valueFromDiags...)
				} else {
					attrSourceValues["ldap_attribute_source"] = types.ObjectNull(ldapAttrSourceAttrTypes)
				}
				attrSourceElement, objectValueFromDiags := types.ObjectValue(attributeSourcesElementAttrTypes, attrSourceValues)
				diags.Append(objectValueFromDiags...)
				attrSourceElements = append(attrSourceElements, attrSourceElement)
			}
			attributeMappingValues["attribute_sources"], valueFromDiags = types.ListValue(types.ObjectType{AttrTypes: attributeSourcesElementAttrTypes}, attrSourceElements)
			diags.Append(valueFromDiags...)
		}

		// Build complete attribute mapping value
		state.AttributeMapping, valueFromDiags = types.ObjectValue(attributeMappingAttrTypes, attributeMappingValues)
		diags.Append(valueFromDiags...)
	}
}

func (r *idpAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdapterResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createIdpAdapter := client.NewIdpAdapter(plan.CustomId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)
	err = addOptionalIdpAdapterFields(ctx, createIdpAdapter, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}
	requestJson, err := createIdpAdapter.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateIdpAdapter := r.apiClient.IdpAdaptersAPI.CreateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapter = apiCreateIdpAdapter.Body(*createIdpAdapter)
	idpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.CreateIdpAdapterExecute(apiCreateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapter", err, httpResp)
		return
	}
	responseJson, err := idpAdapterResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idpAdapterResourceModel

	readIdpAdapterResponse(ctx, idpAdapterResponse, &state, plan, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpAdapterResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpAdapter, httpResp, err := r.apiClient.IdpAdaptersAPI.GetIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an IdpAdapter", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadIdpAdapter.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpAdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan idpAdapterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	updateIdpAdapter := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.CustomId.ValueString())

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createUpdateRequest := client.NewIdpAdapter(plan.CustomId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)

	err = addOptionalIdpAdapterFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateIdpAdapter = updateIdpAdapter.Body(*createUpdateRequest)
	updateIdpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapterExecute(updateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdpAdapter", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateIdpAdapterResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	var state idpAdapterResourceModel
	readIdpAdapterResponse(ctx, updateIdpAdapterResponse, &state, plan, &resp.Diagnostics)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

// Delete the Idp Adapter
func (r *idpAdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idpAdapterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.IdpAdaptersAPI.DeleteIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Idp Adapter", err, httpResp)
	}
}

func (r *idpAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
