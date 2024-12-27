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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	jdbcTagConfigAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"connection_url": types.StringType,
			"tags":           types.StringType,
			"default_source": types.BoolType,
		},
	}

	jdbcDataStoreDataSourceAttrType = map[string]attr.Type{
		"max_pool_size":                types.Int64Type,
		"connection_url_tags":          types.SetType{ElemType: jdbcTagConfigAttrType},
		"type":                         types.StringType,
		"name":                         types.StringType,
		"blocking_timeout":             types.Int64Type,
		"idle_timeout":                 types.Int64Type,
		"min_pool_size":                types.Int64Type,
		"driver_class":                 types.StringType,
		"connection_url":               types.StringType,
		"user_name":                    types.StringType,
		"allow_multi_value_attributes": types.BoolType,
		"validate_connection_sql":      types.StringType,
		"encrypted_password":           types.StringType,
	}

	jdbcDataStoreAttrType                = internaltypes.AddKeyValToMapStringAttrType(jdbcDataStoreDataSourceAttrType, "password", types.StringType)
	jdbcDataStoreEmptyStateObj           = types.ObjectNull(jdbcDataStoreAttrType)
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
			Description: "The password needed to access the database. Either this attribute or `encrypted_password` must be specified.",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"encrypted_password": schema.StringAttribute{
			Description: "The encrypted password needed to access the database. Either this attribute or `password` must be specified.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("password")),
			},
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Defaults to a combination of the `connection_url` and `username`.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"min_pool_size": schema.Int64Attribute{
			Description: "The smallest number of database connections in the connection pool for the given data store. The default value is `10`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(10),
		},
		"max_pool_size": schema.Int64Attribute{
			Description: "The largest number of database connections in the connection pool for the given data store. The default value is `100`.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(100),
		},
		"connection_url_tags": schema.SetNestedAttribute{
			Description: "The set of connection URLs and associated tags for this JDBC data store. This is required if 'connection_url' is not provided.",
			Computed:    true,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"connection_url": schema.StringAttribute{
						Description: "The location of the JDBC database.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"tags": schema.StringAttribute{
						Description: "Tags associated with the `connection_url`. At runtime, nodes will use the first `connection_url_tags` element that has a tag that matches with `node.tags` in the run.properties file.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"default_source": schema.BoolAttribute{
						Description: "Whether this is the default connection. Default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
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
			Description: "The amount of time in milliseconds a request waits to get a connection from the connection pool before it fails. The default value is `5000` milliseconds.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(5000),
		},
		"idle_timeout": schema.Int64Attribute{
			Description: "The length of time in minutes the connection can be idle in the pool before it is closed. The default value is `5` minutes.",
			Computed:    true,
			Optional:    true,
			Default:     int64default.StaticInt64(5),
		},
		"driver_class": schema.StringAttribute{
			Description: "The name of the driver class used to communicate with the source database.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"connection_url": schema.StringAttribute{
			Description: "The default location of the JDBC database. This field is required if `connection_url_tags` is not specified.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.AtLeastOneOf(
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url_tags"),
					path.Expression.AtName(path.MatchRoot("jdbc_data_store"), "connection_url"),
				),
			},
		},
		"user_name": schema.StringAttribute{
			Description: "The name that identifies the user when connecting to the database.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"allow_multi_value_attributes": schema.BoolAttribute{
			Description: "Indicates that this data store can select more than one record from a column and return the results as a multi-value attribute. Default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"validate_connection_sql": schema.StringAttribute{
			Description: "A simple SQL statement used by PingFederate at runtime to verify that the database connection is still active and to reconnect if needed.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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
			Description: "The set of connection URLs and associated tags for this JDBC data store.",
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
						Description: "Tags associated with the connection URL. At runtime, nodes will use the first `connection_url_tags` element that has a tag that matches with node.tags in the run.properties file.",
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
		diags.AddError(providererror.InternalProviderError, "Failed to read JDBC data store from PingFederate. The response from PingFederate was nil.")
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

	password, ok := plan.JdbcDataStore.Attributes()["password"].(types.String)
	if !ok {
		password = types.StringNull()
	}

	encryptedPassword := types.StringPointerValue(jdbcDataStore.EncryptedPassword)
	if internaltypes.IsDefined(plan.JdbcDataStore.Attributes()["encrypted_password"]) {
		encryptedPassword = types.StringValue(plan.JdbcDataStore.Attributes()["encrypted_password"].(types.String).ValueString())
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
		"user_name":                    types.StringPointerValue(jdbcDataStore.UserName),
		"allow_multi_value_attributes": types.BoolPointerValue(jdbcDataStore.AllowMultiValueAttributes),
		"validate_connection_sql":      types.StringPointerValue(jdbcDataStore.ValidateConnectionSql),
		"encrypted_password":           encryptedPassword,
	}

	var toStateObjVal types.Object
	if isResource {
		jdbcAttrValue["password"] = password
		toStateObjVal, diags = types.ObjectValue(jdbcDataStoreAttrType, jdbcAttrValue)
		allDiags = append(allDiags, diags...)
	} else {
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

	userName, ok := jdbcDataStorePlan["user_name"]
	if ok {
		addRequest.JdbcDataStore.UserName = userName.(types.String).ValueStringPointer()
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

	encryptedPassword, ok := jdbcDataStorePlan["encrypted_password"]
	if ok {
		addRequest.JdbcDataStore.EncryptedPassword = encryptedPassword.(types.String).ValueStringPointer()
	}
	return nil
}

func createJdbcDataStore(plan dataStoreModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	jdbcPlan := plan.JdbcDataStore.Attributes()
	driverClass := jdbcPlan["driver_class"].(types.String).ValueString()

	createJdbcDataStore := client.JdbcDataStoreAsDataStoreAggregation(client.NewJdbcDataStore(driverClass, "JDBC"))
	err = addOptionalJdbcDataStoreFields(createJdbcDataStore, con, client.JdbcDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for DataStore: "+err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createJdbcDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpErrorCustomId(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse, &customId)
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

	updateJdbcDataStore := client.JdbcDataStoreAsDataStoreAggregation(client.NewJdbcDataStore(driverClass, "JDBC"))
	err = addOptionalJdbcDataStoreFields(updateJdbcDataStore, con, client.JdbcDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the DataStore: "+err.Error())
		return
	}

	response, httpResponse, err := updateDataStore(updateJdbcDataStore, dsr, con, resp, plan.DataStoreId.ValueString())
	if err != nil {
		config.ReportHttpErrorCustomId(con, &resp.Diagnostics, "An error occurred while updating the DataStore", err, httpResponse, &customId)
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
