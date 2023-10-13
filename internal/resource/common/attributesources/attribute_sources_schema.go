package attributesources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var attributeSourcesEmptyList, _ = types.ListValue(types.ObjectType{AttrTypes: ElemAttrType()}, []attr.Value{})

func CommonAttributeSourceSchema() map[string]schema.Attribute {
	commonAttributeSourceSchema := map[string]schema.Attribute{}
	commonAttributeSourceSchema["data_store_ref"] = schema.SingleNestedAttribute{
		Required:    true,
		Description: "Reference to the associated data store.",
		Attributes:  resourcelink.ToSchema(),
	}
	commonAttributeSourceSchema["id"] = schema.StringAttribute{
		Optional:    true,
		Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
	}
	commonAttributeSourceSchema["description"] = schema.StringAttribute{
		Optional:    true,
		Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings",
	}
	commonAttributeSourceSchema["attribute_contract_fulfillment"] = attributecontractfulfillment.ToSchema(false)
	return commonAttributeSourceSchema
}

func CustomAttributeSourceSchemaAttributes() map[string]schema.Attribute {
	customAttributeSourceSchema := CommonAttributeSourceSchema()
	customAttributeSourceSchema["type"] = schema.StringAttribute{
		Computed:    true,
		Optional:    false,
		Description: "The data store type of this attribute source.",
		Default:     stringdefault.StaticString("CUSTOM"),
	}
	customAttributeSourceSchema["filter_fields"] = schema.ListNestedAttribute{
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
	}
	return customAttributeSourceSchema
}

func JdbcAttributeSourceSchemaAttributes() map[string]schema.Attribute {
	jdbcAttributeSourceSchema := CommonAttributeSourceSchema()
	jdbcAttributeSourceSchema["type"] = schema.StringAttribute{
		Computed:    true,
		Optional:    false,
		Description: "The data store type of this attribute source.",
		Default:     stringdefault.StaticString("JDBC"),
	}
	jdbcAttributeSourceSchema["schema"] = schema.StringAttribute{
		Description: "Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.",
		Optional:    true,
	}
	jdbcAttributeSourceSchema["filter"] = schema.StringAttribute{
		Description: "The JDBC WHERE clause used to query your data store to locate a user record.",
		Required:    true,
	}
	jdbcAttributeSourceSchema["table"] = schema.StringAttribute{
		Description: "The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.",
		Required:    true,
	}
	jdbcAttributeSourceSchema["column_names"] = schema.ListAttribute{
		ElementType: basetypes.StringType{},
		Optional:    true,
		Description: "A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.",
	}
	return jdbcAttributeSourceSchema
}

func LdapAttributeSourceSchemaAttributes() map[string]schema.Attribute {
	ldapAttributeSourceSchema := CommonAttributeSourceSchema()
	ldapAttributeSourceSchema["type"] = schema.StringAttribute{
		Required:    true,
		Description: "The data store type of this attribute source.",
		Validators: []validator.String{
			stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY"}...),
		},
	}
	ldapAttributeSourceSchema["base_dn"] = schema.StringAttribute{
		Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	ldapAttributeSourceSchema["search_scope"] = schema.StringAttribute{
		Description: "Determines the node depth of the query.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("OBJECT", "ONE_LEVEL", "SUBTREE"),
		},
	}
	ldapAttributeSourceSchema["search_filter"] = schema.StringAttribute{
		Description: "The LDAP filter that will be used to lookup the objects from the directory.",
		Required:    true,
	}
	ldapAttributeSourceSchema["search_attributes"] = schema.ListAttribute{
		Description: "A list of LDAP attributes returned from search and available for mapping.",
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
	ldapAttributeSourceSchema["member_of_nested_group"] = schema.BoolAttribute{
		Description: "Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.",
		Computed:    true,
		Optional:    true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		Default: booldefault.StaticBool(false),
	}
	ldapAttributeSourceSchema["binary_attribute_settings"] = schema.SingleNestedAttribute{
		Description: "The advanced settings for binary LDAP attributes.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"binary_encoding": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"BASE64", "HEX", "SID"}...),
				},
			},
		},
	}
	return ldapAttributeSourceSchema
}

func ToSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "A list of configured data stores to look up attributes from.",
		Computed:    true,
		Optional:    true,
		Default:     listdefault.StaticValue(attributeSourcesEmptyList),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"custom_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a custom data store.",
					Optional:    true,
					Attributes:  CustomAttributeSourceSchemaAttributes(),
				},
				"jdbc_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a JDBC data store.",
					Optional:    true,
					Attributes:  JdbcAttributeSourceSchemaAttributes(),
				},
				"ldap_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a LDAP data store.",
					Optional:    true,
					Attributes:  LdapAttributeSourceSchemaAttributes(),
				},
			},
		},
	}
}
