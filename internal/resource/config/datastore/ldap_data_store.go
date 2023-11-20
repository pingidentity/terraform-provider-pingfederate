package datastore

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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
	ldapDataStoreSchema.Description = "An LDAP Data Store"
	ldapDataStoreSchema.Default = objectdefault.StaticValue(types.ObjectNull(ldapDataStoreAttrType))
	ldapDataStoreSchema.Computed = true
	ldapDataStoreSchema.Optional = true
	ldapDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("LDAP"),
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Omitting this attribute will set the value to a combination of the connection url and the username.",
			Computed:    true,
			Optional:    true,
		},
		"read_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds a connection waits for a response to be returned before producing an error. A value of -1 causes the connection to wait indefinitely. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"hostnames": schema.SetAttribute{
			Description: "The default LDAP host names. This field is required if no mapping for host names and tags are specified.",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"verify_host": schema.BoolAttribute{
			Description: "Verifies that the presented server certificate includes the address to which the client intended to establish a connection. Omitting this attribute will set the value to true.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(true),
		},
		"test_on_return": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being returned to the pool.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
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
			Default:     int64default.StaticInt64(0),
		},
		"connection_timeout": schema.Int64Attribute{
			Description: "The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error. A value of -1 causes the pool to wait indefinitely. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"min_connections": schema.Int64Attribute{
			Description: "The smallest number of connections that can remain in each pool, without creating extra ones. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(10),
		},
		"max_connections": schema.Int64Attribute{
			Description: "The largest number of active connections that can remain in each pool without releasing extra ones. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(100),
		},
		"use_ssl": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is false.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"test_on_borrow": schema.BoolAttribute{
			Description: "Indicates whether objects are validated before being borrowed from the pool.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"ldap_dns_srv_prefix": schema.StringAttribute{
			Description: "The prefix value used to discover LDAP DNS SRV record. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("_ldap._tcp"),
		},
		"use_dns_srv_records": schema.BoolAttribute{
			Description: "Use DNS SRV Records to discover LDAP server information. The default value is false.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"create_if_necessary": schema.BoolAttribute{
			Description: "Indicates whether temporary connections can be created when the Maximum Connections threshold is reached.",
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
			Description: "The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool. Omitting this attribute or setting a value of -1 causes the pool not to wait at all and to either create a new connection or produce an error (when no connections are available).",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(-1),
		},
		"hostnames_tags": schema.SetNestedAttribute{
			Description: "A LDAP data store's host names and tags configuration. This is required if no default LDAP host names are specified.",
			Computed:    true,
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
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
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
			Description: "The frequency, in milliseconds, that the evictor cleans up the connections in the pool. A value of -1 disables the evictor. Omitting this attribute will set the value to the default value.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"user_dn": schema.StringAttribute{
			Description: "The username credential required to access the data store.",
			Required:    true,
		},
		"password": schema.StringAttribute{
			Description: "The password credential required to access the data store. GETs will not return this attribute. To update this field, specify the new value in this attribute.",
			Required:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
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

func toStateLdapDataStore(con context.Context, ldapDataStore *client.LdapDataStore, plan basetypes.ObjectValue) (types.Object, diag.Diagnostics) {
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

	//  final obj value
	ldapDataStoreAttrVal := map[string]attr.Value{
		"hostnames":              internaltypes.GetStringSet(ldapDataStore.Hostnames),
		"verify_host":            types.BoolPointerValue(ldapDataStore.VerifyHost),
		"test_on_return":         types.BoolPointerValue(ldapDataStore.TestOnReturn),
		"ldap_type":              types.StringValue(ldapDataStore.LdapType),
		"dns_ttl":                types.Int64PointerValue(ldapDataStore.DnsTtl),
		"connection_timeout":     types.Int64PointerValue(ldapDataStore.ConnectionTimeout),
		"min_connections":        types.Int64PointerValue(ldapDataStore.MinConnections),
		"use_ssl":                types.BoolPointerValue(ldapDataStore.UseSsl),
		"test_on_borrow":         types.BoolPointerValue(ldapDataStore.TestOnBorrow),
		"ldap_dns_srv_prefix":    types.StringPointerValue(ldapDataStore.LdapDnsSrvPrefix),
		"name":                   types.StringPointerValue(ldapDataStore.Name),
		"read_timeout":           types.Int64PointerValue(ldapDataStore.ReadTimeout),
		"use_dns_srv_records":    types.BoolPointerValue(ldapDataStore.UseDnsSrvRecords),
		"max_connections":        types.Int64PointerValue(ldapDataStore.MaxConnections),
		"user_dn":                userDn(),
		"create_if_necessary":    types.BoolPointerValue(ldapDataStore.CreateIfNecessary),
		"binary_attributes":      binaryAttributes,
		"max_wait":               types.Int64PointerValue(ldapDataStore.MaxWait),
		"hostnames_tags":         hostnamesTagsVal,
		"time_between_evictions": types.Int64PointerValue(ldapDataStore.TimeBetweenEvictions),
		"type":                   types.StringValue("LDAP"),
		"password":               password,
		"bind_anonymously":       types.BoolPointerValue(ldapDataStore.BindAnonymously),
		"follow_ldap_referrals":  followLdapReferrals,
	}

	ldapDataStoreObj, diags := types.ObjectValue(ldapDataStoreAttrType, ldapDataStoreAttrVal)
	allDiags = append(allDiags, diags...)
	return ldapDataStoreObj, allDiags
}

func readLdapDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreResourceModel, plan *types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.LdapDataStore.Id)
	state.DataStoreId = types.StringPointerValue(r.LdapDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.LdapDataStore.MaskAttributeValues)
	state.CustomDataStore = customDataStoreEmptyStateObj
	state.JdbcDataStore = jdbcDataStoreEmptyStateObj
	state.LdapDataStore, diags = toStateLdapDataStore(ctx, r.LdapDataStore, *plan)
	state.PingOneLdapGatewayDataStore = pingOneLdapGatewayDataStoreEmptyStateObj
	return diags
}

func addOptionalLdapDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.LdapDataStore, plan dataStoreResourceModel) error {
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
		addRequest.LdapDataStore.Hostnames = internaltypes.SetTypeToStringSet(hostnames.(types.Set))
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
		addRequest.LdapDataStore.BinaryAttributes = internaltypes.SetTypeToStringSet(binaryAttributes.(types.Set))
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

	return nil
}

func createLdapDataStore(plan dataStoreResourceModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	ldapPlan := plan.LdapDataStore.Attributes()
	ldapType := ldapPlan["ldap_type"].(types.String).ValueString()
	createLdapDataStore := client.LdapDataStoreAsDataStoreAggregation(client.NewLdapDataStore(ldapType, "LDAP"))
	err = addOptionalLdapDataStoreFields(createLdapDataStore, con, client.LdapDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createLdapDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse)
		return
	}

	// Read the response into the state
	var state dataStoreResourceModel
	diags = readLdapDataStoreResponse(con, response, &state, &plan.LdapDataStore)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updateLdapDataStore(plan dataStoreResourceModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	ldapPlan := plan.LdapDataStore.Attributes()
	ldapType := ldapPlan["ldap_type"].(types.String).ValueString()
	updateLdapDataStore := client.LdapDataStoreAsDataStoreAggregation(client.NewLdapDataStore(ldapType, "LDAP"))
	err = addOptionalLdapDataStoreFields(updateLdapDataStore, con, client.LdapDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := updateDataStore(updateLdapDataStore, dsr, con, resp, plan.Id.ValueString())
	if err != nil && (httpResponse == nil || httpResponse.StatusCode != 404) {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating DataStore", err, httpResponse)
		return
	}

	// Read the response
	var state dataStoreResourceModel
	diags = readLdapDataStoreResponse(con, response, &state, &plan.LdapDataStore)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
