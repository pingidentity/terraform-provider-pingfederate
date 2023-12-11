package datastore

import (
	"context"
	"encoding/json"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	jdbcTagConfigAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"connection_url": basetypes.StringType{},
			"tags":           basetypes.StringType{},
			"default_source": basetypes.BoolType{},
		},
	}

	jdbcDataStoreCommonAttrType = map[string]attr.Type{
		"max_pool_size":                basetypes.Int64Type{},
		"connection_url_tags":          basetypes.SetType{ElemType: jdbcTagConfigAttrType},
		"type":                         basetypes.StringType{},
		"name":                         basetypes.StringType{},
		"blocking_timeout":             basetypes.Int64Type{},
		"idle_timeout":                 basetypes.Int64Type{},
		"min_pool_size":                basetypes.Int64Type{},
		"driver_class":                 basetypes.StringType{},
		"connection_url":               basetypes.StringType{},
		"user_name":                    basetypes.StringType{},
		"allow_multi_value_attributes": basetypes.BoolType{},
		"validate_connection_sql":      basetypes.StringType{},
	}

	jdbcDataStoreAttrType                = internaltypes.AddKeyValToMapStringAttrType(jdbcDataStoreCommonAttrType, "password", types.StringType)
	jdbcDataStoreEmptyStateObj           = types.ObjectNull(jdbcDataStoreAttrType)
	jdbcDataStoreDataSourceAttrType      = internaltypes.AddKeyValToMapStringAttrType(jdbcDataStoreCommonAttrType, "encrypted_password", types.StringType)
	jdbcDataStoreEmptyDataSourceStateObj = types.ObjectNull(jdbcDataStoreDataSourceAttrType)
)

func toSchemaJdbcDataStore() schema.SingleNestedAttribute {
	jdbcDataStoreSchema := schema.SingleNestedAttribute{}
	jdbcDataStoreSchema.Description = "A JDBC data store."
	jdbcDataStoreSchema.Optional = true
	jdbcDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("JDBC"),
		},
		"password": schema.StringAttribute{
			Description: "The password needed to access the database. GETs will not return this attribute. To update this field, specify the new value in this attribute.",
			Optional:    true,
			Sensitive:   true,
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Omitting this attribute will set the value to a combination of the connection url and the username.",
			Computed:    true,
			Optional:    true,
		},
		"min_pool_size": schema.Int64Attribute{
			Description: "The smallest number of database connections in the connection pool for the given data store. Omitting this attribute will set the value to the connection pool default.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"max_pool_size": schema.Int64Attribute{
			Description: "The largest number of database connections in the connection pool for the given data store. Omitting this attribute will set the value to the connection pool default.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"connection_url_tags": schema.SetNestedAttribute{
			Description: "A JDBC data store's connection URLs and tags configuration. This is required if no default JDBC database location is specified.",
			Computed:    true,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"connection_url": schema.StringAttribute{
						Description: "The location of the JDBC database.",
						Required:    true,
					},
					"tags": schema.StringAttribute{
						Description: "Tags associated with this data source.",
						Optional:    true,
					},
					"default_source": schema.BoolAttribute{
						Description: "Whether this is the default connection. Defaults to false if not specified.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.AtLeastOneOf(
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url_tags"),
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url"),
				),
			},
		},
		"blocking_timeout": schema.Int64Attribute{
			Description: "The amount of time in milliseconds a request waits to get a connection from the connection pool before it fails. Omitting this attribute will set the value to the connection pool default.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"idle_timeout": schema.Int64Attribute{
			Description: "The length of time in minutes the connection can be idle in the pool before it is closed. Omitting this attribute will set the value to the connection pool default.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(0),
		},
		"driver_class": schema.StringAttribute{
			Description: "The name of the driver class used to communicate with the source database.",
			Required:    true,
		},
		"connection_url": schema.StringAttribute{
			Description: "The default location of the JDBC database. This field is required if no mapping for JDBC database location and tags are specified.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.AtLeastOneOf(
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url_tags"),
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url"),
				),
			},
		},
		"user_name": schema.StringAttribute{
			Description: "The name that identifies the user when connecting to the database.",
			Required:    true,
		},
		"allow_multi_value_attributes": schema.BoolAttribute{
			Description: "Indicates that this data store can select more than one record from a column and return the results as a multi-value attribute.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"validate_connection_sql": schema.StringAttribute{
			Description: "A simple SQL statement used by PingFederate at runtime to verify that the database connection is still active and to reconnect if needed.",
			Optional:    true,
		},
	}

	jdbcDataStoreSchema.Validators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("custom_data_store"),
			path.MatchRelative().AtParent().AtName("ldap_data_store"),
			path.MatchRelative().AtParent().AtName("ping_one_ldap_gateway_data_store"),
		),
	}

	return jdbcDataStoreSchema
}

func toDataSourceSchemaJdbcDataStore() datasourceschema.SingleNestedAttribute {
	jdbcDataStoreSchema := datasourceschema.SingleNestedAttribute{}
	jdbcDataStoreSchema.Description = "A JDBC data store."
	jdbcDataStoreSchema.Computed = true
	jdbcDataStoreSchema.Optional = false
	jdbcDataStoreSchema.Attributes = map[string]datasourceschema.Attribute{
		"type": datasourceschema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
		},
		"encrypted_password": datasourceschema.StringAttribute{
			Description: "The encrypted password needed to access the database.",
			Computed:    true,
			Optional:    false,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The data store name with a unique value across all data sources.",
			Computed:    true,
			Optional:    false,
		},
		"min_pool_size": datasourceschema.Int64Attribute{
			Description: "The smallest number of database connections in the connection pool for the given data store.",
			Computed:    true,
			Optional:    false,
		},
		"max_pool_size": datasourceschema.Int64Attribute{
			Description: "The largest number of database connections in the connection pool for the given data store.",
			Computed:    true,
			Optional:    false,
		},
		"connection_url_tags": datasourceschema.SetNestedAttribute{
			Description: "A JDBC data store's connection URLs and tags configuration.",
			Computed:    true,
			Optional:    false,
			NestedObject: datasourceschema.NestedAttributeObject{
				Attributes: map[string]datasourceschema.Attribute{
					"connection_url": datasourceschema.StringAttribute{
						Description: "The location of the JDBC database.",
						Computed:    true,
						Optional:    false,
					},
					"tags": datasourceschema.StringAttribute{
						Description: "Tags associated with this data source.",
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
		"blocking_timeout": datasourceschema.Int64Attribute{
			Description: "The amount of time in milliseconds a request waits to get a connection from the connection pool before it fails.",
			Computed:    true,
			Optional:    false,
		},
		"idle_timeout": datasourceschema.Int64Attribute{
			Description: "The length of time in minutes the connection can be idle in the pool before it is closed.",
			Computed:    true,
			Optional:    false,
		},
		"driver_class": datasourceschema.StringAttribute{
			Description: "The name of the driver class used to communicate with the source database.",
			Computed:    true,
			Optional:    false,
		},
		"connection_url": datasourceschema.StringAttribute{
			Description: "The default location of the JDBC database.",
			Computed:    true,
			Optional:    false,
		},
		"user_name": datasourceschema.StringAttribute{
			Description: "The name that identifies the user when connecting to the database.",
			Computed:    true,
			Optional:    false,
		},
		"allow_multi_value_attributes": datasourceschema.BoolAttribute{
			Description: "Indicates that this data store can select more than one record from a column and return the results as a multi-value attribute.",
			Computed:    true,
			Optional:    false,
		},
		"validate_connection_sql": datasourceschema.StringAttribute{
			Description: "A simple SQL statement used by PingFederate at runtime to verify that the database connection is still active and to reconnect if needed.",
			Computed:    true,
			Optional:    false,
		},
	}

	return jdbcDataStoreSchema
}

func toStateJdbcDataStore(con context.Context, jdbcDataStore *client.JdbcDataStore, plan dataStoreModel, isResource bool) (types.Object, diag.Diagnostics) {
	var allDiags, diags diag.Diagnostics

	if jdbcDataStore == nil {
		diags.AddError("Failed to read JDBC data store from PingFederate.", "The response from PingFederate was nil.")
		return types.ObjectNull(jdbcDataStoreAttrType), diags
	}

	connectionUrlTags := func() (types.Set, diag.Diagnostics) {
		if len(jdbcDataStore.ConnectionUrlTags) > 0 {
			connectionUrlTagsSetVal, diags := types.SetValueFrom(con, jdbcTagConfigAttrType, jdbcDataStore.ConnectionUrlTags)
			return connectionUrlTagsSetVal, diags
		} else {
			connectionUrlTagsSetVal := types.SetNull(jdbcTagConfigAttrType)
			return connectionUrlTagsSetVal, diag.Diagnostics{}
		}
	}
	connectionUrlSetVal, diags := connectionUrlTags()
	allDiags = append(allDiags, diags...)

	var password basetypes.StringValue
	passwordVal, ok := plan.JdbcDataStore.Attributes()["password"].(types.String)
	if ok {
		password = passwordVal
	} else {
		password = types.StringPointerValue(pointers.String(""))
	}

	jdbcAttrValue := map[string]attr.Value{
		"type":                         types.StringValue("JDBC"),
		"blocking_timeout":             types.Int64PointerValue(jdbcDataStore.BlockingTimeout),
		"connection_url":               types.StringPointerValue(jdbcDataStore.ConnectionUrl),
		"driver_class":                 types.StringValue(jdbcDataStore.DriverClass),
		"connection_url_tags":          connectionUrlSetVal,
		"idle_timeout":                 types.Int64PointerValue(jdbcDataStore.IdleTimeout),
		"max_pool_size":                types.Int64PointerValue(jdbcDataStore.MaxPoolSize),
		"min_pool_size":                types.Int64PointerValue(jdbcDataStore.MinPoolSize),
		"name":                         types.StringPointerValue(jdbcDataStore.Name),
		"user_name":                    types.StringValue(jdbcDataStore.UserName),
		"allow_multi_value_attributes": types.BoolPointerValue(jdbcDataStore.AllowMultiValueAttributes),
		//TODO does this matter
		"validate_connection_sql": types.StringPointerValue(jdbcDataStore.ValidateConnectionSql),
	}

	var toStateObjVal types.Object
	if isResource {
		jdbcAttrValue["password"] = password
		toStateObjVal, diags = types.ObjectValue(jdbcDataStoreAttrType, jdbcAttrValue)
		allDiags = append(allDiags, diags...)
	} else {
		jdbcAttrValue["encrypted_password"] = types.StringPointerValue(jdbcDataStore.EncryptedPassword)
		toStateObjVal, diags = types.ObjectValue(jdbcDataStoreDataSourceAttrType, jdbcAttrValue)
		allDiags = append(allDiags, diags...)
	}
	return toStateObjVal, allDiags
}

func readJdbcDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreModel, plan *dataStoreModel, isResource bool) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.JdbcDataStore.Id)
	state.DataStoreId = types.StringPointerValue(r.JdbcDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.JdbcDataStore.MaskAttributeValues)
	state.PingOneLdapGatewayDataStore = pingOneLdapGatewayDataStoreEmptyStateObj
	if isResource {
		state.JdbcDataStore, diags = toStateJdbcDataStore(ctx, r.JdbcDataStore, *plan, true)
		state.CustomDataStore = customDataStoreEmptyStateObj
		state.LdapDataStore = ldapDataStoreEmptyStateObj
	} else {
		state.JdbcDataStore, diags = toStateJdbcDataStore(ctx, r.JdbcDataStore, *plan, false)
		state.CustomDataStore = customDataStoreEmptyDataSourceStateObj
		state.LdapDataStore = ldapDataStoreEmptyDataSourceStateObj
	}
	return diags
}

func addOptionalJdbcDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.JdbcDataStore, plan dataStoreModel) error {

	if internaltypes.IsDefined(plan.MaskAttributeValues) {
		addRequest.JdbcDataStore.MaskAttributeValues = plan.MaskAttributeValues.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.DataStoreId) {
		addRequest.JdbcDataStore.Id = plan.DataStoreId.ValueStringPointer()
	}

	jdbcDataStorePlan := plan.JdbcDataStore.Attributes()
	maxPoolSize, ok := jdbcDataStorePlan["max_pool_size"]
	if ok {
		addRequest.JdbcDataStore.MaxPoolSize = maxPoolSize.(types.Int64).ValueInt64Pointer()
	}

	minPoolSize, ok := jdbcDataStorePlan["min_pool_size"]
	if ok {
		addRequest.JdbcDataStore.MinPoolSize = minPoolSize.(types.Int64).ValueInt64Pointer()
	}

	connectionUrlTags, ok := jdbcDataStorePlan["connection_url_tags"]
	if ok {
		addRequest.JdbcDataStore.ConnectionUrlTags = []client.JdbcTagConfig{}
		err := json.Unmarshal([]byte(internaljson.FromValue(connectionUrlTags, true)), &addRequest.JdbcDataStore.ConnectionUrlTags)
		if err != nil {
			return err
		}
	}

	name, ok := jdbcDataStorePlan["name"]
	if ok {
		addRequest.JdbcDataStore.Name = name.(types.String).ValueStringPointer()
	}

	blockingTimeout, ok := jdbcDataStorePlan["blocking_timeout"]
	if ok {
		addRequest.JdbcDataStore.BlockingTimeout = blockingTimeout.(types.Int64).ValueInt64Pointer()
	}

	idleTimeout, ok := jdbcDataStorePlan["idle_timeout"]
	if ok {
		addRequest.JdbcDataStore.IdleTimeout = idleTimeout.(types.Int64).ValueInt64Pointer()
	}

	connectionUrl, ok := jdbcDataStorePlan["connection_url"]
	if ok {
		addRequest.JdbcDataStore.ConnectionUrl = connectionUrl.(types.String).ValueStringPointer()
	}

	allowMultiValueAttributes, ok := jdbcDataStorePlan["allow_multi_value_attributes"]
	if ok {
		addRequest.JdbcDataStore.AllowMultiValueAttributes = allowMultiValueAttributes.(types.Bool).ValueBoolPointer()
	}

	validateConnectionSql, ok := jdbcDataStorePlan["validate_connection_sql"]
	if ok {
		addRequest.JdbcDataStore.ValidateConnectionSql = validateConnectionSql.(types.String).ValueStringPointer()
	}

	password, ok := jdbcDataStorePlan["password"]
	if ok {
		addRequest.JdbcDataStore.Password = password.(types.String).ValueStringPointer()
	}
	return nil
}

func createJdbcDataStore(plan dataStoreModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	jdbcPlan := plan.JdbcDataStore.Attributes()
	driverClass := jdbcPlan["driver_class"].(types.String).ValueString()
	userName := jdbcPlan["user_name"].(types.String).ValueString()

	createJdbcDataStore := client.JdbcDataStoreAsDataStoreAggregation(client.NewJdbcDataStore(driverClass, userName, "JDBC"))
	err = addOptionalJdbcDataStoreFields(createJdbcDataStore, con, client.JdbcDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createJdbcDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse)
		return
	}

	// Read the response into the state
	var state dataStoreModel
	diags = readJdbcDataStoreResponse(con, response, &state, &plan, true)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updateJdbcDataStore(plan dataStoreModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	jdbcPlan := plan.JdbcDataStore.Attributes()
	driverClass := jdbcPlan["driver_class"].(types.String).ValueString()
	userName := jdbcPlan["user_name"].(types.String).ValueString()

	updateJdbcDataStore := client.JdbcDataStoreAsDataStoreAggregation(client.NewJdbcDataStore(driverClass, userName, "JDBC"))
	err = addOptionalJdbcDataStoreFields(updateJdbcDataStore, con, client.JdbcDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := updateDataStore(updateJdbcDataStore, dsr, con, resp, plan.Id.ValueString())
	if err != nil && (httpResponse == nil || httpResponse.StatusCode != 404) {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating DataStore", err, httpResponse)
		return
	}

	// Read the response
	var state dataStoreModel
	diags = readJdbcDataStoreResponse(con, response, &state, &plan, true)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
