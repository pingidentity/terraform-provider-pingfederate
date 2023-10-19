package datastore

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	pingOneLdapGatewayDataStoreAttrType = map[string]attr.Type{
		"ping_one_connection_ref":  basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"ldap_type":                basetypes.StringType{},
		"ping_one_ldap_gateway_id": basetypes.StringType{},
		"use_ssl":                  basetypes.BoolType{},
		"name":                     basetypes.StringType{},
		"binary_attributes":        basetypes.SetType{ElemType: basetypes.StringType{}},
		"type":                     basetypes.StringType{},
		"ping_one_environment_id":  basetypes.StringType{},
	}
	pingOneLdapGatewayDataStoreEmptyStateObj = types.ObjectNull(pingOneLdapGatewayDataStoreAttrType)
)

func toSchemaPingOneLdapGatewayDataStore() schema.SingleNestedAttribute {
	pingOneLdapGatewayDataStoreSchema := schema.SingleNestedAttribute{}
	pingOneLdapGatewayDataStoreSchema.Description = "A PingOne LDAP Gateway data store."
	pingOneLdapGatewayDataStoreSchema.Default = objectdefault.StaticValue(types.ObjectNull(pingOneLdapGatewayDataStoreAttrType))
	pingOneLdapGatewayDataStoreSchema.Computed = true
	pingOneLdapGatewayDataStoreSchema.Optional = true
	pingOneLdapGatewayDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("PING_ONE_LDAP_GATEWAY"),
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Omitting this attribute will set the value to a combination of the hostname(s) and the principal.",
			Computed:    true,
			Optional:    true,
		},
		"ping_one_connection_ref": schema.SingleNestedAttribute{
			Computed:    true,
			Optional:    true,
			Description: "Reference to the PingOne connection this gateway uses.",
			Default:     objectdefault.StaticValue(types.ObjectNull(resourcelink.AttrType())),
			Attributes:  resourcelink.ToSchema(),
		},
		"ldap_type": schema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically. The value is validated against the LDAP gateway configuration in PingOne unless the header 'X-BypassExternalValidation' is set to true.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACTIVE_DIRECTORY", "ORACLE_DIRECTORY_SERVER", "ORACLE_UNIFIED_DIRECTORY", "UNBOUNDID_DS", "PING_DIRECTORY", "GENERIC"),
			},
		},
		"ping_one_ldap_gateway_id": schema.StringAttribute{
			Description: "The ID of the PingOne LDAP Gateway this data store uses.",
			Required:    true,
		},
		"use_ssl": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is false. The value is validated against the LDAP gateway configuration in PingOne unless the header 'X-BypassExternalValidation' is set to true.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"binary_attributes": schema.SetAttribute{
			Description: "A list of LDAP attributes to be handled as binary data.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"ping_one_environment_id": schema.StringAttribute{
			Description: "The environment ID that the gateway belongs to.",
			Required:    true,
		},
	}

	pingOneLdapGatewayDataStoreSchema.Validators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("custom_data_store"),
			path.MatchRelative().AtParent().AtName("jdbc_data_store"),
			path.MatchRelative().AtParent().AtName("ldap_data_store"),
		),
	}

	return pingOneLdapGatewayDataStoreSchema
}
