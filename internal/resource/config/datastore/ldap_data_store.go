package datastore

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	ldapTagConfigAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"hostnames":      basetypes.SetType{ElemType: basetypes.StringType{}},
			"tags":           basetypes.StringType{},
			"default_source": basetypes.BoolType{},
		},
	}

	ldapDataStoreAttrType = map[string]attr.Type{
		"hostnames":              basetypes.SetType{ElemType: basetypes.StringType{}},
		"verify_host":            basetypes.BoolType{},
		"test_on_return":         basetypes.BoolType{},
		"ldap_type":              basetypes.StringType{},
		"dns_ttl":                basetypes.Int64Type{},
		"connection_timeout":     basetypes.Int64Type{},
		"min_connections":        basetypes.Int64Type{},
		"use_ssl":                basetypes.BoolType{},
		"test_on_borrow":         basetypes.BoolType{},
		"ldap_dns_srv_prefix":    basetypes.StringType{},
		"name":                   basetypes.StringType{},
		"read_timeout":           basetypes.Int64Type{},
		"use_dns_srv_records":    basetypes.BoolType{},
		"max_connections":        basetypes.Int64Type{},
		"user_dn":                basetypes.StringType{},
		"create_if_necessary":    basetypes.BoolType{},
		"binary_attributes":      basetypes.SetType{ElemType: basetypes.StringType{}},
		"max_wait":               basetypes.Int64Type{},
		"hostnames_tags":         basetypes.SetType{ElemType: ldapTagConfigAttrType},
		"time_between_evictions": basetypes.Int64Type{},
		"type":                   basetypes.StringType{},
		"password":               basetypes.StringType{},
		"bind_anonymously":       basetypes.BoolType{},
		"follow_ldap_referrals":  basetypes.BoolType{},
	}

	ldapDataStoreEmptyStateObj = types.ObjectNull(ldapDataStoreAttrType)
)

func toSchemaLdapDataStore() schema.SingleNestedAttribute {
	ldapDataStoreSchema := schema.SingleNestedAttribute{}
	ldapDataStoreSchema.Description = "A LDAP Data Store"
	ldapDataStoreSchema.Default = objectdefault.StaticValue(types.ObjectNull(ldapDataStoreAttrType))
	ldapDataStoreSchema.Computed = true
	ldapDataStoreSchema.Optional = true
	ldapDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("LDAP"),
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Omitting this attribute will set the value to a combination of the connection url and the username.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"read_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds a connection waits for a response to be returned before producing an error. A value of -1 causes the connection to wait indefinitely. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"hostnames": schema.SetAttribute{
			Description: "The default LDAP host names. This field is required if no mapping for host names and tags are specified.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"verify_host": schema.BoolAttribute{
			Description: "Verifies that the presented server certificate includes the address to which the client intended to establish a connection. Omitting this attribute will set the value to true.",
			Default:     booldefault.StaticBool(true),
			Computed:    true,
			Optional:    true,
		},
		"test_on_return": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being returned to the pool.",
			Optional:    true,
		},
		"ldap_type": schema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically. The 'UNBOUNDID_DS' type has been deprecated, please use the 'PING_DIRECTORY' type instead.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACTIVE_DIRECTORY", "ORACLE_DIRECTORY_SERVER", "ORACLE_UNIFIED_DIRECTORY", "PING_DIRECTORY", "GENERIC"),
			},
		},
		"dns_ttl": schema.Int64Attribute{
			Description: "The maximum time in milliseconds that DNS information are cached. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"connection_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error. A value of -1 causes the pool to wait indefinitely. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"min_connections": schema.Int64Attribute{
			Description: "The smallest number of connections that can remain in each pool, without creating extra ones. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"max_connections": schema.Int64Attribute{
			Description: "The largest number of active connections that can remain in each pool without releasing extra ones. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"use_ssl": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is false.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"test_on_borrow": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being borrowed from the pool.",
			Optional:    true,
		},
		"ldap_dns_srv_prefix": schema.StringAttribute{
			Description: "The prefix value used to discover LDAP DNS SRV record. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
		},
		"use_dns_srv_records": schema.BoolAttribute{
			Description: "Use DNS SRV Records to discover LDAP server information. The default value is false.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"user_dn": schema.StringAttribute{
			Description: "The username credential required to access the data store.",
			Optional:    true,
		},
		"create_if_necessary": schema.BoolAttribute{
			Description: "Indicates whether temporary connections can be created when the Maximum Connections threshold is reached.",
			Optional:    true,
		},
		"binary_attributes": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of LDAP attributes to be handled as binary data.",
			Optional:    true,
		},
		"max_wait": schema.Int64Attribute{
			Description: "The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool. Omitting this attribute or setting a value of -1 causes the pool not to wait at all and to either create a new connection or produce an error (when no connections are available).",
			Optional:    true,
		},
		"hostnames_tags": schema.SetNestedAttribute{
			Description: "A LDAP data store's host names and tags configuration. This is required if no default LDAP host names are specified.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"hostnames": schema.SetAttribute{
						Description: "The LDAP host names.",
						Required:    true,
						ElementType: types.StringType,
					},
					"tags": schema.StringAttribute{
						Description: "Tags associated with this data source.",
						Optional:    true,
					},
					"default_source": schema.BoolAttribute{
						Description: "Whether this is the default connection. Defaults to false if not specified.",
						Default:     booldefault.StaticBool(false),
						Computed:    true,
						Optional:    true,
					},
				},
			},
		},
		"time_between_evictions": schema.Int64Attribute{
			Description: "The frequency, in milliseconds, that the evictor cleans up the connections in the pool. A value of -1 disables the evictor. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"password": schema.StringAttribute{
			Description: "The password credential required to access the data store. GETs will not return this attribute. To update this field, specify the new value in this attribute.",
			Optional:    true,
			Sensitive:   true,
		},
		"bind_anonymously": schema.BoolAttribute{
			Description: "Whether username and password are required. The default value is false.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"follow_ldap_referrals": schema.BoolAttribute{
			Description: "Follow LDAP Referrals in the domain tree. The default value is false. This property does not apply to PingDirectory as this functionality is configured in PingDirectory.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
	}

	ldapDataStoreSchema.Validators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("custom_data_store"),
			path.MatchRelative().AtParent().AtName("jdbc_data_store"),
			path.MatchRelative().AtParent().AtName("ping_one_ldap_gateway_data_store"),
		),
	}

	return ldapDataStoreSchema
}
