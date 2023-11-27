package attributesources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
)

func CommonAttributeSourceDataSourceSchema() map[string]schema.Attribute {
	commonAttributeSourceDataSourceSchema := map[string]schema.Attribute{}
	commonAttributeSourceDataSourceSchema["data_store_ref"] = schema.SingleNestedAttribute{
		Optional:    false,
		Computed:    true,
		Description: "Reference to the associated data store.",
		Attributes:  resourcelink.ToDataSourceSchema(),
	}
	commonAttributeSourceDataSourceSchema["id"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
	}
	commonAttributeSourceDataSourceSchema["description"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings",
	}
	commonAttributeSourceDataSourceSchema["attribute_contract_fulfillment"] = attributecontractfulfillment.ToDataSourceSchema()
	return commonAttributeSourceDataSourceSchema
}

func CustomAttributeSourceDataSourceSchemaAttributes() map[string]schema.Attribute {
	customAttributeSourceDataSourceSchema := CommonAttributeSourceDataSourceSchema()
	customAttributeSourceDataSourceSchema["type"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The data store type of this attribute source.",
	}
	customAttributeSourceDataSourceSchema["filter_fields"] = schema.ListNestedAttribute{
		Description: "The list of fields that can be used to filter a request to the custom data store.",
		Optional:    false,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"value": schema.StringAttribute{
					Description: "The value of this field. Whether or not the value is required will be determined by plugin validation checks.",
					Optional:    false,
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "The name of this field.",
					Optional:    false,
					Computed:    true,
				},
			},
		},
	}
	return customAttributeSourceDataSourceSchema
}

func JdbcAttributeSourceDataSourceSchemaAttributes() map[string]schema.Attribute {
	jdbcAttributeSourceDataSourceSchema := CommonAttributeSourceDataSourceSchema()
	jdbcAttributeSourceDataSourceSchema["type"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The data store type of this attribute source.",
	}
	jdbcAttributeSourceDataSourceSchema["schema"] = schema.StringAttribute{
		Description: "Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.",
		Optional:    false,
		Computed:    true,
	}
	jdbcAttributeSourceDataSourceSchema["filter"] = schema.StringAttribute{
		Description: "The JDBC WHERE clause used to query your data store to locate a user record.",
		Optional:    false,
		Computed:    true,
	}
	jdbcAttributeSourceDataSourceSchema["table"] = schema.StringAttribute{
		Description: "The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.",
		Optional:    false,
		Computed:    true,
	}
	jdbcAttributeSourceDataSourceSchema["column_names"] = schema.ListAttribute{
		ElementType: basetypes.StringType{},
		Optional:    false,
		Computed:    true,
		Description: "A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.",
	}
	return jdbcAttributeSourceDataSourceSchema
}

func LdapAttributeSourceDataSourceSchemaAttributes() map[string]schema.Attribute {
	ldapAttributeSourceDataSourceSchema := CommonAttributeSourceDataSourceSchema()
	ldapAttributeSourceDataSourceSchema["type"] = schema.StringAttribute{
		Optional:    false,
		Computed:    true,
		Description: "The data store type of this attribute source.",
	}
	ldapAttributeSourceDataSourceSchema["base_dn"] = schema.StringAttribute{
		Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
		Optional:    false,
		Computed:    true,
	}
	ldapAttributeSourceDataSourceSchema["search_scope"] = schema.StringAttribute{
		Description: "Determines the node depth of the query.",
		Optional:    false,
		Computed:    true,
	}
	ldapAttributeSourceDataSourceSchema["search_filter"] = schema.StringAttribute{
		Description: "The LDAP filter that will be used to lookup the objects from the directory.",
		Optional:    false,
		Computed:    true,
	}
	ldapAttributeSourceDataSourceSchema["search_attributes"] = schema.ListAttribute{
		Description: "A list of LDAP attributes returned from search and available for mapping.",
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
	ldapAttributeSourceDataSourceSchema["member_of_nested_group"] = schema.BoolAttribute{
		Description: "Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.",
		Optional:    false,
		Computed:    true,
	}
	ldapAttributeSourceDataSourceSchema["binary_attribute_settings"] = schema.MapNestedAttribute{
		Description: "The advanced settings for binary LDAP attributes.",
		Optional:    false,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"binary_encoding": schema.StringAttribute{
					Optional:    false,
					Computed:    true,
					Description: "Get the encoding type for this attribute. If not specified, the default is BASE64.",
				},
			},
		},
	}
	return ldapAttributeSourceDataSourceSchema
}

func ToDataSourceSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: false,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"custom_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a custom data store.",
					Optional:    false,
					Computed:    true,
					Attributes:  CustomAttributeSourceDataSourceSchemaAttributes(),
				},
				"jdbc_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a JDBC data store.",
					Optional:    false,
					Computed:    true,
					Attributes:  JdbcAttributeSourceDataSourceSchemaAttributes(),
				},
				"ldap_attribute_source": schema.SingleNestedAttribute{
					Description: "The configured settings used to look up attributes from a LDAP data store.",
					Optional:    false,
					Computed:    true,
					Attributes:  LdapAttributeSourceDataSourceSchemaAttributes(),
				},
			},
		},
	}
}
