package datastore

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	ldapTagConfigAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"hostnames":      types.ListType{ElemType: types.StringType},
			"tags":           types.StringType,
			"default_source": types.BoolType,
		},
	}

	ldapDataStoreCommonAttrType = map[string]attr.Type{
		"hostnames":                  types.ListType{ElemType: types.StringType},
		"use_start_tls":              types.BoolType,
		"verify_host":                types.BoolType,
		"test_on_return":             types.BoolType,
		"ldap_type":                  types.StringType,
		"dns_ttl":                    types.Int64Type,
		"connection_timeout":         types.Int64Type,
		"min_connections":            types.Int64Type,
		"use_ssl":                    types.BoolType,
		"test_on_borrow":             types.BoolType,
		"ldap_dns_srv_prefix":        types.StringType,
		"name":                       types.StringType,
		"read_timeout":               types.Int64Type,
		"use_dns_srv_records":        types.BoolType,
		"max_connections":            types.Int64Type,
		"user_dn":                    types.StringType,
		"create_if_necessary":        types.BoolType,
		"binary_attributes":          types.SetType{ElemType: types.StringType},
		"max_wait":                   types.Int64Type,
		"hostnames_tags":             types.SetType{ElemType: ldapTagConfigAttrType},
		"time_between_evictions":     types.Int64Type,
		"type":                       types.StringType,
		"bind_anonymously":           types.BoolType,
		"follow_ldap_referrals":      types.BoolType,
		"client_tls_certificate_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"retry_failed_operations":    types.BoolType,
	}

	ldapDataStoreAttrType                = internaltypes.AddKeyValToMapStringAttrType(ldapDataStoreCommonAttrType, "password", types.StringType)
	ldapDataStoreEmptyStateObj           = types.ObjectNull(ldapDataStoreAttrType)
	ldapDataStoreEncryptedPassAttrType   = internaltypes.AddKeyValToMapStringAttrType(ldapDataStoreCommonAttrType, "encrypted_password", types.StringType)
	ldapDataStoreEmptyDataSourceStateObj = types.ObjectNull(ldapDataStoreEncryptedPassAttrType)
)

func toSchemaLdapDataStore() schema.SingleNestedAttribute {
	ldapDataStoreSchema := schema.SingleNestedAttribute{}
	ldapDataStoreSchema.Description = "An LDAP Data Store"
	ldapDataStoreSchema.Optional = true
	ldapDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("LDAP"),
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Defaults to a combination of the values of `hostnames` and `user_dn`.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"read_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds a connection waits for a response to be returned before producing an error. A value of `-1` causes the connection to wait indefinitely. Defaults to `0`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"hostnames": schema.ListAttribute{
			Description: "The default LDAP host names. This field is required if `hostnames_tags` is not specified. Failover can be configured by providing multiple host names.",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"use_start_tls": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure StartTLS encryption. The default value is `false`.",
			Computed:    true,
			Optional:    true,
		},
		"verify_host": schema.BoolAttribute{
			Description: "Verifies that the presented server certificate includes the address to which the client intended to establish a connection. Defaults to `true`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(true),
		},
		"test_on_return": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being returned to the pool. Default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"ldap_type": schema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically. The `UNBOUNDID_DS` type has been deprecated, please use the `PING_DIRECTORY` type instead. Supported values are `ACTIVE_DIRECTORY`, `ORACLE_DIRECTORY_SERVER`, `ORACLE_UNIFIED_DIRECTORY`, `PING_DIRECTORY`, `GENERIC`.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACTIVE_DIRECTORY", "ORACLE_DIRECTORY_SERVER", "ORACLE_UNIFIED_DIRECTORY", "PING_DIRECTORY", "GENERIC"),
			},
		},
		"dns_ttl": schema.Int64Attribute{
			Description: "The maximum time in milliseconds that DNS information are cached. Defaults to `0`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"connection_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error. A value of `-1` causes the pool to wait indefinitely. Defaults to `0`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"min_connections": schema.Int64Attribute{
			Description: "The smallest number of connections that can remain in each pool, without creating extra ones. Defaults to `10`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(10),
		},
		"max_connections": schema.Int64Attribute{
			Description: "The largest number of active connections that can remain in each pool without releasing extra ones. Defaults to `100`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(100),
		},
		"use_ssl": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"test_on_borrow": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being borrowed from the pool. Default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"ldap_dns_srv_prefix": schema.StringAttribute{
			Description: "The prefix value used to discover LDAP DNS SRV record. Defaults to `_ldap._tcp`.",
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("_ldap._tcp"),
		},
		"use_dns_srv_records": schema.BoolAttribute{
			Description: "Use DNS SRV Records to discover LDAP server information. The default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"create_if_necessary": schema.BoolAttribute{
			Description: "Indicates whether temporary connections can be created when the Maximum Connections threshold is reached. Default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"binary_attributes": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of LDAP attributes to be handled as binary data.",
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetNull(types.StringType)),
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"max_wait": schema.Int64Attribute{
			Description: "The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool. Setting a value of `-1` causes the pool not to wait at all and to either create a new connection or produce an error (when no connections are available). Defaults to `-1`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(-1),
		},
		"hostnames_tags": schema.SetNestedAttribute{
			Description: "The set of host names and associated tags for this LDAP data store. This is required if 'hostnames' is not provided.",
			Computed:    true,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"hostnames": schema.ListAttribute{
						Description: "The LDAP host names. Failover can be configured by providing multiple host names.",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
					"tags": schema.StringAttribute{
						Description: "Tags associated with the host names. At runtime, nodes will use the first `hostnames_tags` element that has a tag that matches with node.tags in the run.properties file.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"default_source": schema.BoolAttribute{
						Description: "Whether this is the default connection. Defaults to `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.AtLeastOneOf(
					path.Expression.AtName(path.MatchRoot("ldap_data_store"), "hostnames_tags"),
					path.Expression.AtName(path.MatchRoot("ldap_data_store"), "hostnames"),
				),
			},
		},
		"time_between_evictions": schema.Int64Attribute{
			Description: "The frequency, in milliseconds, that the evictor cleans up the connections in the pool. A value of `-1` disables the evictor. Defaults to `0`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"user_dn": schema.StringAttribute{
			Description: "The username credential required to access the data store. Mutually exclusive with `bind_anonymously` and `client_tls_certificate_ref`. `password` must also be set to use this attribute.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"password": schema.StringAttribute{
			Description: "The password credential required to access the data store. Requires `user_dn` to be set.",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"bind_anonymously": schema.BoolAttribute{
			Description: "Whether username and password are required. If `true`, then `user_dn` and `client_tls_certificate_ref` cannot be set. The default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"follow_ldap_referrals": schema.BoolAttribute{
			Description: "Follow LDAP Referrals in the domain tree. The default value is `false`. This property does not apply to PingDirectory as this functionality is configured in PingDirectory.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"client_tls_certificate_ref": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The client TLS certificate used to access the data store. If specified, authentication to the data store will be done using mutual TLS. See '/keyPairs/sslClient' to manage certificates. Supported in PF version `11.3` or later. In order to use this authentication method, you must set either `use_start_tls` or `use_ssl` to `true`. Mutually exclusive with `bind_anonymously` and `user_dn`",
			Attributes:  resourcelink.ToSchema(),
		},
		"retry_failed_operations": schema.BoolAttribute{
			Description: "Indicates whether failed operations should be retried. The default is `false`. Supported in PF version `11.3` or later.",
			Computed:    true,
			Optional:    true,
			// The default is set in ModifyPlan, since it is dependent on PF version 11.3+
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

func toDataSourceSchemaLdapDataStore() datasourceschema.SingleNestedAttribute {
	ldapDataStoreSchema := datasourceschema.SingleNestedAttribute{}
	ldapDataStoreSchema.Description = "An LDAP Data Store"
	ldapDataStoreSchema.Computed = true
	ldapDataStoreSchema.Optional = false
	ldapDataStoreSchema.Attributes = map[string]datasourceschema.Attribute{
		"type": datasourceschema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The data store name with a unique value across all data sources.",
			Computed:    true,
			Optional:    true,
		},
		"read_timeout": datasourceschema.Int64Attribute{
			Description: "The maximum number of milliseconds a connection waits for a response to be returned before producing an error.",
			Computed:    true,
			Optional:    false,
		},
		"hostnames": datasourceschema.ListAttribute{
			Description: "The default LDAP host names. Failover can be configured by providing multiple host names.",
			Computed:    true,
			Optional:    false,
			ElementType: types.StringType,
		},
		"use_start_tls": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure StartTLS encryption.",
			Computed:    true,
			Optional:    false,
		},
		"verify_host": datasourceschema.BoolAttribute{
			Description: "Verifies that the presented server certificate includes the address to which the client intended to establish a connection.",
			Computed:    true,
			Optional:    false,
		},
		"test_on_return": datasourceschema.BoolAttribute{
			Description: "Indicates whether objects are validated before being returned to the pool.",
			Computed:    true,
			Optional:    false,
		},
		"ldap_type": datasourceschema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically.",
			Computed:    true,
			Optional:    false,
		},
		"dns_ttl": datasourceschema.Int64Attribute{
			Description: "The maximum time in milliseconds that DNS information are cached.",
			Computed:    true,
			Optional:    false,
		},
		"connection_timeout": datasourceschema.Int64Attribute{
			Description: "The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error.",
			Computed:    true,
			Optional:    false,
		},
		"min_connections": datasourceschema.Int64Attribute{
			Description: "The smallest number of connections that can remain in each pool, without creating extra ones.",
			Computed:    true,
			Optional:    false,
		},
		"max_connections": datasourceschema.Int64Attribute{
			Description: "The largest number of active connections that can remain in each pool without releasing extra ones.",
			Computed:    true,
			Optional:    false,
		},
		"use_ssl": datasourceschema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS).",
			Computed:    true,
			Optional:    false,
		},
		"test_on_borrow": datasourceschema.BoolAttribute{
			Description: "Indicates whether objects are validated before being borrowed from the pool.",
			Computed:    true,
			Optional:    false,
		},
		"ldap_dns_srv_prefix": datasourceschema.StringAttribute{
			Description: "The prefix value used to discover LDAP DNS SRV record.",
			Computed:    true,
			Optional:    false,
		},
		"use_dns_srv_records": datasourceschema.BoolAttribute{
			Description: "Use DNS SRV Records to discover LDAP server information.",
			Computed:    true,
			Optional:    false,
		},
		"create_if_necessary": datasourceschema.BoolAttribute{
			Description: "Indicates whether temporary connections can be created when the Maximum Connections threshold is reached.",
			Computed:    true,
			Optional:    false,
		},
		"binary_attributes": datasourceschema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of LDAP attributes to be handled as binary data.",
			Computed:    true,
			Optional:    false,
		},
		"max_wait": datasourceschema.Int64Attribute{
			Description: "The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool.",
			Computed:    true,
			Optional:    false,
		},
		"hostnames_tags": datasourceschema.SetNestedAttribute{
			Description: "The set of host names and associated tags for this LDAP data store.",
			Computed:    true,
			Optional:    false,
			NestedObject: datasourceschema.NestedAttributeObject{
				Attributes: map[string]datasourceschema.Attribute{
					"hostnames": datasourceschema.ListAttribute{
						Description: "The LDAP host names. Failover can be configured by providing multiple host names.",
						Computed:    true,
						Optional:    false,
						ElementType: types.StringType,
					},
					"tags": datasourceschema.StringAttribute{
						Description: "Tags associated with the host names. At runtime, nodes will use the first `hostname_tags` element that has a tag that matches with node.tags in the run.properties file.",
						Computed:    true,
						Optional:    false,
					},
					"default_source": datasourceschema.BoolAttribute{
						Description: "Whether this is the default connection.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
		"time_between_evictions": datasourceschema.Int64Attribute{
			Description: "The frequency, in milliseconds, that the evictor cleans up the connections in the pool.",
			Computed:    true,
			Optional:    false,
		},
		"user_dn": datasourceschema.StringAttribute{
			Description: "The username credential required to access the data store.",
			Computed:    true,
			Optional:    false,
		},
		"encrypted_password": datasourceschema.StringAttribute{
			Description: "The encrypted password credential required to access the data store.",
			Computed:    true,
			Optional:    false,
		},
		"bind_anonymously": datasourceschema.BoolAttribute{
			Description: "Whether username and password are required.",
			Computed:    true,
			Optional:    false,
		},
		"follow_ldap_referrals": datasourceschema.BoolAttribute{
			Description: "Follow LDAP Referrals in the domain tree.",
			Computed:    true,
			Optional:    false,
		},
		"client_tls_certificate_ref": datasourceschema.SingleNestedAttribute{
			Computed:    true,
			Optional:    false,
			Description: "The client TLS certificate used to access the data store. If specified, authentication to the data store will be done using mutual TLS. See '/keyPairs/sslClient' to manage certificates. Supported in PF version `11.3` or later.",
			Attributes:  datasourceresourcelink.ToDataSourceSchema(),
		},
		"retry_failed_operations": datasourceschema.BoolAttribute{
			Description: "Indicates whether failed operations should be retried. Supported in PF version `11.3` or later.",
			Computed:    true,
			Optional:    false,
		},
	}

	return ldapDataStoreSchema
}

func toStateLdapDataStore(con context.Context, ldapDataStore *client.LdapDataStore, plan types.Object) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics

	if ldapDataStore == nil {
		diags.AddError("Failed to read Ldap data store from PingFederate.", "The response from PingFederate was nil.")
		return ldapDataStoreEmptyStateObj, diags
	}

	userDn := func() types.String {
		if ldapDataStore.BindAnonymously != nil && *ldapDataStore.BindAnonymously {
			userDnFromPlan, ok := plan.Attributes()["user_dn"]
			if ok {
				return userDnFromPlan.(types.String)
			} else {
				return types.StringNull()
			}
		} else {
			return types.StringPointerValue(ldapDataStore.UserDN)
		}
	}

	hostnamesTags := func() (basetypes.SetValue, diag.Diagnostics) {
		if len(ldapDataStore.HostnamesTags) > 0 {
			return types.SetValueFrom(con, ldapTagConfigAttrType, ldapDataStore.HostnamesTags)
		} else {
			return types.SetNull(ldapTagConfigAttrType), nil
		}
	}

	hostnamesTagsVal, diags := hostnamesTags()
	allDiags = append(allDiags, diags...)

	var followLdapReferrals types.Bool
	if ldapDataStore.LdapType == "PING_DIRECTORY" {
		followLdapReferrals = types.BoolValue(false)
	} else {
		followLdapReferrals = types.BoolPointerValue(ldapDataStore.FollowLDAPReferrals)
	}

	var password types.String
	if plan.Attributes()["password"] != nil {
		password = plan.Attributes()["password"].(types.String)
	} else {
		password = types.StringNull()
	}

	var binaryAttributes basetypes.SetValue
	if len(ldapDataStore.BinaryAttributes) > 0 {
		binaryAttributes = internaltypes.GetStringSet(ldapDataStore.BinaryAttributes)
	} else {
		binaryAttributes = types.SetNull(types.StringType)
	}

	clientTlsCertificateRef, diags := resourcelink.ToState(con, ldapDataStore.ClientTlsCertificateRef)
	allDiags = append(allDiags, diags...)

	hostnamesVal, diags := types.ListValueFrom(con, types.StringType, ldapDataStore.Hostnames)
	allDiags = append(allDiags, diags...)

	//  final obj value
	ldapDataStoreAttrVal := map[string]attr.Value{
		"hostnames":                  hostnamesVal,
		"use_start_tls":              types.BoolPointerValue(ldapDataStore.UseStartTLS),
		"verify_host":                types.BoolPointerValue(ldapDataStore.VerifyHost),
		"test_on_return":             types.BoolPointerValue(ldapDataStore.TestOnReturn),
		"ldap_type":                  types.StringValue(ldapDataStore.LdapType),
		"dns_ttl":                    types.Int64PointerValue(ldapDataStore.DnsTtl),
		"connection_timeout":         types.Int64PointerValue(ldapDataStore.ConnectionTimeout),
		"min_connections":            types.Int64PointerValue(ldapDataStore.MinConnections),
		"use_ssl":                    types.BoolPointerValue(ldapDataStore.UseSsl),
		"test_on_borrow":             types.BoolPointerValue(ldapDataStore.TestOnBorrow),
		"ldap_dns_srv_prefix":        types.StringPointerValue(ldapDataStore.LdapDnsSrvPrefix),
		"name":                       types.StringPointerValue(ldapDataStore.Name),
		"read_timeout":               types.Int64PointerValue(ldapDataStore.ReadTimeout),
		"use_dns_srv_records":        types.BoolPointerValue(ldapDataStore.UseDnsSrvRecords),
		"max_connections":            types.Int64PointerValue(ldapDataStore.MaxConnections),
		"user_dn":                    userDn(),
		"create_if_necessary":        types.BoolPointerValue(ldapDataStore.CreateIfNecessary),
		"binary_attributes":          binaryAttributes,
		"max_wait":                   types.Int64PointerValue(ldapDataStore.MaxWait),
		"hostnames_tags":             hostnamesTagsVal,
		"time_between_evictions":     types.Int64PointerValue(ldapDataStore.TimeBetweenEvictions),
		"type":                       types.StringValue("LDAP"),
		"password":                   password,
		"bind_anonymously":           types.BoolPointerValue(ldapDataStore.BindAnonymously),
		"follow_ldap_referrals":      followLdapReferrals,
		"client_tls_certificate_ref": clientTlsCertificateRef,
		"retry_failed_operations":    types.BoolPointerValue(ldapDataStore.RetryFailedOperations),
	}

	ldapDataStoreObj, diags := types.ObjectValue(ldapDataStoreAttrType, ldapDataStoreAttrVal)
	allDiags = append(allDiags, diags...)
	return ldapDataStoreObj, allDiags
}

func toDataSourceStateLdapDataStore(con context.Context, ldapDataStore *client.LdapDataStore) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics

	if ldapDataStore == nil {
		diags.AddError("Failed to read Ldap data store from PingFederate.", "The response from PingFederate was nil.")
		return ldapDataStoreEmptyStateObj, diags
	}

	hostnamesTagsVal, diags := types.SetValueFrom(con, ldapTagConfigAttrType, ldapDataStore.HostnamesTags)
	allDiags = append(allDiags, diags...)

	var followLdapReferrals types.Bool
	if ldapDataStore.LdapType == "PING_DIRECTORY" {
		followLdapReferrals = types.BoolValue(false)
	} else {
		followLdapReferrals = types.BoolPointerValue(ldapDataStore.FollowLDAPReferrals)
	}

	clientTlsCertificateRef, diags := resourcelink.ToState(con, ldapDataStore.ClientTlsCertificateRef)
	allDiags = append(allDiags, diags...)

	hostnamesVal, diags := types.ListValueFrom(con, types.StringType, ldapDataStore.Hostnames)
	allDiags = append(allDiags, diags...)

	//  final obj value
	ldapDataStoreAttrVal := map[string]attr.Value{
		"hostnames":                  hostnamesVal,
		"use_start_tls":              types.BoolPointerValue(ldapDataStore.UseStartTLS),
		"verify_host":                types.BoolPointerValue(ldapDataStore.VerifyHost),
		"test_on_return":             types.BoolPointerValue(ldapDataStore.TestOnReturn),
		"ldap_type":                  types.StringValue(ldapDataStore.LdapType),
		"dns_ttl":                    types.Int64PointerValue(ldapDataStore.DnsTtl),
		"connection_timeout":         types.Int64PointerValue(ldapDataStore.ConnectionTimeout),
		"min_connections":            types.Int64PointerValue(ldapDataStore.MinConnections),
		"use_ssl":                    types.BoolPointerValue(ldapDataStore.UseSsl),
		"test_on_borrow":             types.BoolPointerValue(ldapDataStore.TestOnBorrow),
		"ldap_dns_srv_prefix":        types.StringPointerValue(ldapDataStore.LdapDnsSrvPrefix),
		"name":                       types.StringPointerValue(ldapDataStore.Name),
		"read_timeout":               types.Int64PointerValue(ldapDataStore.ReadTimeout),
		"use_dns_srv_records":        types.BoolPointerValue(ldapDataStore.UseDnsSrvRecords),
		"max_connections":            types.Int64PointerValue(ldapDataStore.MaxConnections),
		"user_dn":                    types.StringPointerValue(ldapDataStore.UserDN),
		"create_if_necessary":        types.BoolPointerValue(ldapDataStore.CreateIfNecessary),
		"binary_attributes":          internaltypes.GetStringSet(ldapDataStore.BinaryAttributes),
		"max_wait":                   types.Int64PointerValue(ldapDataStore.MaxWait),
		"hostnames_tags":             hostnamesTagsVal,
		"time_between_evictions":     types.Int64PointerValue(ldapDataStore.TimeBetweenEvictions),
		"type":                       types.StringValue("LDAP"),
		"encrypted_password":         types.StringPointerValue(ldapDataStore.EncryptedPassword),
		"bind_anonymously":           types.BoolPointerValue(ldapDataStore.BindAnonymously),
		"follow_ldap_referrals":      followLdapReferrals,
		"client_tls_certificate_ref": clientTlsCertificateRef,
		"retry_failed_operations":    types.BoolPointerValue(ldapDataStore.RetryFailedOperations),
	}

	ldapDataStoreObj, diags := types.ObjectValue(ldapDataStoreEncryptedPassAttrType, ldapDataStoreAttrVal)
	allDiags = append(allDiags, diags...)
	return ldapDataStoreObj, allDiags
}

func readLdapDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreModel, plan *types.Object, isResource bool) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.LdapDataStore.Id)
	state.DataStoreId = types.StringPointerValue(r.LdapDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.LdapDataStore.MaskAttributeValues)
	if isResource {
		state.CustomDataStore = customDataStoreEmptyStateObj
		state.JdbcDataStore = jdbcDataStoreEmptyStateObj
		state.LdapDataStore, diags = toStateLdapDataStore(ctx, r.LdapDataStore, *plan)
	} else {
		state.CustomDataStore = customDataStoreEmptyDataSourceStateObj
		state.LdapDataStore, diags = toDataSourceStateLdapDataStore(ctx, r.LdapDataStore)
		state.JdbcDataStore = jdbcDataStoreEmptyDataSourceStateObj
	}
	state.PingOneLdapGatewayDataStore = pingOneLdapGatewayDataStoreEmptyStateObj
	return diags
}

func addOptionalLdapDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.LdapDataStore, plan dataStoreModel) error {
	ldapDataStorePlan := plan.LdapDataStore.Attributes()

	if internaltypes.IsDefined(plan.MaskAttributeValues) {
		addRequest.LdapDataStore.MaskAttributeValues = plan.MaskAttributeValues.ValueBoolPointer()
	}

	userDn, ok := ldapDataStorePlan["user_dn"]
	if ok {
		addRequest.LdapDataStore.UserDN = userDn.(types.String).ValueStringPointer()
	}

	password, ok := ldapDataStorePlan["password"]
	if ok {
		addRequest.LdapDataStore.Password = password.(types.String).ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.DataStoreId) {
		addRequest.LdapDataStore.Id = plan.DataStoreId.ValueStringPointer()
	}

	hostnames, ok := ldapDataStorePlan["hostnames"]
	if ok {
		addRequest.LdapDataStore.Hostnames = []string{}
		for _, hostname := range hostnames.(types.List).Elements() {
			addRequest.LdapDataStore.Hostnames = append(addRequest.LdapDataStore.Hostnames, hostname.(types.String).ValueString())
		}
	}

	useStartTls, ok := ldapDataStorePlan["use_start_tls"]
	if ok {
		addRequest.LdapDataStore.UseStartTLS = useStartTls.(types.Bool).ValueBoolPointer()
	}

	verifyHost, ok := ldapDataStorePlan["verify_host"]
	if ok {
		addRequest.LdapDataStore.VerifyHost = verifyHost.(types.Bool).ValueBoolPointer()
	}

	testOnReturn, ok := ldapDataStorePlan["test_on_return"]
	if ok {
		addRequest.LdapDataStore.TestOnReturn = testOnReturn.(types.Bool).ValueBoolPointer()
	}

	dnsTtl, ok := ldapDataStorePlan["dns_ttl"]
	if ok {
		addRequest.LdapDataStore.DnsTtl = dnsTtl.(types.Int64).ValueInt64Pointer()
	}

	connectionTimeout, ok := ldapDataStorePlan["connection_timeout"]
	if ok {
		addRequest.LdapDataStore.ConnectionTimeout = connectionTimeout.(types.Int64).ValueInt64Pointer()
	}

	if internaltypes.IsDefined(ldapDataStorePlan["min_connections"]) {
		addRequest.LdapDataStore.MinConnections = ldapDataStorePlan["min_connections"].(types.Int64).ValueInt64Pointer()
	}

	useSsl, ok := ldapDataStorePlan["use_ssl"]
	if ok {
		addRequest.LdapDataStore.UseSsl = useSsl.(types.Bool).ValueBoolPointer()
	}

	testOnBorrow, ok := ldapDataStorePlan["test_on_borrow"]
	if ok {
		addRequest.LdapDataStore.TestOnBorrow = testOnBorrow.(types.Bool).ValueBoolPointer()
	}

	ldapDnsSrvPrefix, ok := ldapDataStorePlan["ldap_dns_srv_prefix"]
	if ok {
		addRequest.LdapDataStore.LdapDnsSrvPrefix = ldapDnsSrvPrefix.(types.String).ValueStringPointer()
	}

	name, ok := ldapDataStorePlan["name"]
	if ok {
		addRequest.LdapDataStore.Name = name.(types.String).ValueStringPointer()
	}

	readTimeout, ok := ldapDataStorePlan["read_timeout"]
	if ok {
		addRequest.LdapDataStore.ReadTimeout = readTimeout.(types.Int64).ValueInt64Pointer()
	}

	useDnsSrvRecords, ok := ldapDataStorePlan["use_dns_srv_records"]
	if ok {
		addRequest.LdapDataStore.UseDnsSrvRecords = useDnsSrvRecords.(types.Bool).ValueBoolPointer()
	}

	if internaltypes.IsDefined(ldapDataStorePlan["max_connections"]) {
		addRequest.LdapDataStore.MaxConnections = ldapDataStorePlan["max_connections"].(types.Int64).ValueInt64Pointer()
	}

	createIfNecessary, ok := ldapDataStorePlan["create_if_necessary"]
	if ok {
		addRequest.LdapDataStore.CreateIfNecessary = createIfNecessary.(types.Bool).ValueBoolPointer()
	}

	binaryAttributes, ok := ldapDataStorePlan["binary_attributes"]
	if ok {
		addRequest.LdapDataStore.BinaryAttributes = internaltypes.SetTypeToStringSlice(binaryAttributes.(types.Set))
	}

	if internaltypes.IsDefined(ldapDataStorePlan["max_wait"]) {
		addRequest.LdapDataStore.MaxWait = ldapDataStorePlan["max_wait"].(types.Int64).ValueInt64Pointer()
	}

	hostnamesTags, ok := ldapDataStorePlan["hostnames_tags"]
	if ok {
		addRequest.LdapDataStore.HostnamesTags = []client.LdapTagConfig{}
		err := json.Unmarshal([]byte(internaljson.FromValue(hostnamesTags, true)), &addRequest.LdapDataStore.HostnamesTags)
		if err != nil {
			return err
		}
	}

	timeBetweenEvictions, ok := ldapDataStorePlan["time_between_evictions"]
	if ok {
		addRequest.LdapDataStore.TimeBetweenEvictions = timeBetweenEvictions.(types.Int64).ValueInt64Pointer()
	}

	bindAnonymously, ok := ldapDataStorePlan["bind_anonymously"]
	if ok {
		addRequest.LdapDataStore.BindAnonymously = bindAnonymously.(types.Bool).ValueBoolPointer()
	}

	followLdapReferrals, ok := ldapDataStorePlan["follow_ldap_referrals"]
	if ok {
		addRequest.LdapDataStore.FollowLDAPReferrals = followLdapReferrals.(types.Bool).ValueBoolPointer()
	}

	clientTlsCertificateRef, ok := ldapDataStorePlan["client_tls_certificate_ref"]
	if ok {
		ref, err := resourcelink.ClientStruct(clientTlsCertificateRef.(types.Object))
		if err != nil {
			return err
		}
		addRequest.LdapDataStore.ClientTlsCertificateRef = ref
	}

	retryFailedOperations, ok := ldapDataStorePlan["retry_failed_operations"]
	if ok {
		addRequest.LdapDataStore.RetryFailedOperations = retryFailedOperations.(types.Bool).ValueBoolPointer()
	}

	return nil
}

func createLdapDataStore(plan dataStoreModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	ldapPlan := plan.LdapDataStore.Attributes()
	ldapType := ldapPlan["ldap_type"].(types.String).ValueString()
	createLdapDataStore := client.LdapDataStoreAsDataStoreAggregation(client.NewLdapDataStore(ldapType, "LDAP"))
	err = addOptionalLdapDataStoreFields(createLdapDataStore, con, client.LdapDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for DataStore: "+err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createLdapDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse)
		return
	}

	// Read the response into the state
	var state dataStoreModel
	diags = readLdapDataStoreResponse(con, response, &state, &plan.LdapDataStore, true)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updateLdapDataStore(plan dataStoreModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	ldapPlan := plan.LdapDataStore.Attributes()
	ldapType := ldapPlan["ldap_type"].(types.String).ValueString()
	updateLdapDataStore := client.LdapDataStoreAsDataStoreAggregation(client.NewLdapDataStore(ldapType, "LDAP"))
	err = addOptionalLdapDataStoreFields(updateLdapDataStore, con, client.LdapDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the DataStore: "+err.Error())
		return
	}

	response, httpResponse, err := updateDataStore(updateLdapDataStore, dsr, con, resp, plan.Id.ValueString())
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating the DataStore", err, httpResponse)
		return
	}

	// Read the response
	var state dataStoreModel
	diags = readLdapDataStoreResponse(con, response, &state, &plan.LdapDataStore, true)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
