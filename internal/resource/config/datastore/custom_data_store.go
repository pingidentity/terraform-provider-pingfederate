package datastore

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	datasourcepluginconfiguration "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	customDataStoreCommonAttrType = map[string]attr.Type{
		"type":                  types.StringType,
		"name":                  types.StringType,
		"plugin_descriptor_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"parent_ref":            types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	customDataStoreAttrType                = internaltypes.AddKeyValToMapStringAttrType(customDataStoreCommonAttrType, "configuration", types.ObjectType{AttrTypes: pluginconfiguration.AttrType()})
	customDataStoreEmptyStateObj           = types.ObjectNull(customDataStoreAttrType)
	customDataStoreDataSourceAttrType      = internaltypes.AddKeyValToMapStringAttrType(customDataStoreCommonAttrType, "configuration", types.ObjectType{AttrTypes: datasourcepluginconfiguration.AttrType()})
	customDataStoreEmptyDataSourceStateObj = types.ObjectNull(customDataStoreDataSourceAttrType)
)

func toSchemaCustomDataStore() schema.SingleNestedAttribute {
	customDataStoreSchema := schema.SingleNestedAttribute{}
	customDataStoreSchema.Description = "A custom data store."
	customDataStoreSchema.Optional = true
	customDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("CUSTOM"),
		},
		"name": schema.StringAttribute{
			Description: "The plugin instance name.",
			Required:    true,
		},
		"plugin_descriptor_ref": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
			Attributes:  resourcelink.ToSchema(),
		},
		"parent_ref": schema.SingleNestedAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
			Default:     objectdefault.StaticValue(types.ObjectNull(resourcelink.AttrType())),
			Attributes:  resourcelink.ToSchema(),
		},
		"configuration": pluginconfiguration.ToSchema(),
	}
	customDataStoreSchema.Validators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("jdbc_data_store"),
			path.MatchRelative().AtParent().AtName("ldap_data_store"),
			path.MatchRelative().AtParent().AtName("ping_one_ldap_gateway_data_store"),
		),
	}

	return customDataStoreSchema
}

func toDataSourceSchemaCustomDataStore() datasourceschema.SingleNestedAttribute {
	customDataStoreSchema := datasourceschema.SingleNestedAttribute{}
	customDataStoreSchema.Description = "A custom data store."
	customDataStoreSchema.Computed = true
	customDataStoreSchema.Optional = false
	customDataStoreSchema.Attributes = map[string]datasourceschema.Attribute{
		"type": datasourceschema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The plugin instance name.",
			Computed:    true,
			Optional:    false,
		},
		"plugin_descriptor_ref": datasourceschema.SingleNestedAttribute{
			Computed:    true,
			Optional:    false,
			Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created..)",
			Attributes:  datasourceresourcelink.ToDataSourceSchema(),
		},
		"parent_ref": datasourceschema.SingleNestedAttribute{
			Computed:    true,
			Optional:    false,
			Description: "The reference to this plugin's parent instance..)",
			Attributes:  datasourceresourcelink.ToDataSourceSchema(),
		},
		"configuration": datasourcepluginconfiguration.ToDataSourceSchema(),
	}

	return customDataStoreSchema
}

func toStateCustomDataStore(con context.Context, clientValue *client.DataStoreAggregation, plan types.Object, isResource bool) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics

	if clientValue.CustomDataStore == nil {
		diags.AddError("Failed to read custom data store from API", "The custom data store was nil")
		return types.ObjectNull(customDataStoreAttrType), diags
	}

	customDataStore := clientValue.CustomDataStore

	var customDataStoreObj types.Object
	pluginDescriptorRef, diags := resourcelink.ToState(con, &customDataStore.PluginDescriptorRef)
	allDiags = append(allDiags, diags...)
	parentRef, diags := resourcelink.ToState(con, customDataStore.ParentRef)
	allDiags = append(allDiags, diags...)
	var configurationObject types.Object
	customDataStoreVal := map[string]attr.Value{
		"type":                  types.StringValue(customDataStore.Type),
		"name":                  types.StringValue(customDataStore.Name),
		"plugin_descriptor_ref": pluginDescriptorRef,
		"parent_ref":            parentRef,
	}
	if isResource {
		planConfiguration, ok := plan.Attributes()["configuration"]
		if ok {
			configurationObject, diags = pluginconfiguration.ToState(planConfiguration.(types.Object), &customDataStore.Configuration)
			allDiags = append(allDiags, diags...)
		} else {
			configurationObject, diags = pluginconfiguration.ToState(types.ObjectNull(pluginconfiguration.AttrType()), &customDataStore.Configuration)
			allDiags = append(allDiags, diags...)
		}
		customDataStoreVal["configuration"] = configurationObject
		customDataStoreObj, diags = types.ObjectValue(customDataStoreAttrType, customDataStoreVal)
		allDiags = append(allDiags, diags...)
	} else {
		configurationObject, diags := datasourcepluginconfiguration.ToDataSourceState(con, &customDataStore.Configuration)
		allDiags = append(allDiags, diags...)
		customDataStoreVal["configuration"] = configurationObject
		customDataStoreObj, diags = types.ObjectValue(customDataStoreDataSourceAttrType, customDataStoreVal)
		allDiags = append(allDiags, diags...)
	}
	return customDataStoreObj, allDiags
}

func readCustomDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreModel, plan *types.Object, isResource bool) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.CustomDataStore.Id)
	state.DataStoreId = types.StringPointerValue(r.CustomDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.CustomDataStore.MaskAttributeValues)
	state.PingOneLdapGatewayDataStore = pingOneLdapGatewayDataStoreEmptyStateObj
	if isResource {
		state.JdbcDataStore = jdbcDataStoreEmptyStateObj
		state.CustomDataStore, diags = toStateCustomDataStore(ctx, r, *plan, true)
		state.LdapDataStore = ldapDataStoreEmptyStateObj
	} else {
		state.JdbcDataStore = jdbcDataStoreEmptyDataSourceStateObj
		state.CustomDataStore, diags = toStateCustomDataStore(ctx, r, *plan, false)
		state.LdapDataStore = ldapDataStoreEmptyDataSourceStateObj
	}
	return diags
}

func addOptionalCustomDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createCustomDataStore client.CustomDataStore, plan dataStoreModel) error {
	customDataStorePlan := plan.CustomDataStore.Attributes()

	if internaltypes.IsDefined(plan.DataStoreId) {
		addRequest.CustomDataStore.Id = plan.DataStoreId.ValueStringPointer()
	}

	parentRef := customDataStorePlan["parent_ref"]
	if internaltypes.IsNonEmptyObj(parentRef.(types.Object)) {
		parentRef, err := resourcelink.ClientStruct(parentRef.(types.Object))
		if err != nil {
			return err
		}
		addRequest.CustomDataStore.ParentRef = parentRef
	}
	return nil
}

func createCustomDataStore(plan dataStoreModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	customPlan := plan.CustomDataStore.Attributes()
	name := customPlan["name"].(types.String).ValueString()
	pluginDescriptorRef, err := resourcelink.ClientStruct(customPlan["plugin_descriptor_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create plugin descriptor reference object for DataStore", err.Error())
		return
	}

	configuration := &client.PluginConfiguration{}
	err = json.Unmarshal([]byte(internaljson.FromValue(customPlan["configuration"].(types.Object), true)), configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create configuration object for DataStore", err.Error())
		return
	}

	createCustomDataStore := client.CustomDataStoreAsDataStoreAggregation(client.NewCustomDataStore("CUSTOM", name, *pluginDescriptorRef, *configuration))
	err = addOptionalCustomDataStoreFields(createCustomDataStore, con, client.CustomDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createCustomDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse)
		return
	}
	// Read the response into the state
	var state dataStoreModel
	diags = readCustomDataStoreResponse(con, response, &state, &plan.CustomDataStore, true)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updateCustomDataStore(plan dataStoreModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	customPlan := plan.CustomDataStore.Attributes()
	pluginDescriptorRef, err := resourcelink.ClientStruct(customPlan["plugin_descriptor_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create plugin descriptor reference object for DataStore", err.Error())
		return
	}

	configuration := &client.PluginConfiguration{}
	err = json.Unmarshal([]byte(internaljson.FromValue(customPlan["configuration"].(types.Object), true)), configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create configuration object for DataStore", err.Error())
		return
	}

	name := customPlan["name"].(types.String).ValueString()
	updateCustomDataStore := client.CustomDataStoreAsDataStoreAggregation(client.NewCustomDataStore("CUSTOM", name, *pluginDescriptorRef, *configuration))
	err = addOptionalCustomDataStoreFields(updateCustomDataStore, con, client.CustomDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	response, httpResponse, err := updateDataStore(updateCustomDataStore, dsr, con, resp, plan.Id.ValueString())
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating the DataStore", err, httpResponse)
		return
	}
	// Read the response
	var state dataStoreModel
	diags = readCustomDataStoreResponse(con, response, &state, &plan.CustomDataStore, true)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
