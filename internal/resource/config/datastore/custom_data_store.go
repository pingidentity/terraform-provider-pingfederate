package datastore

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	customDataStoreAttrType = map[string]attr.Type{
		"type":                  basetypes.StringType{},
		"name":                  basetypes.StringType{},
		"plugin_descriptor_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"parent_ref":            basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"configuration":         basetypes.ObjectType{AttrTypes: pluginconfiguration.AttrType()},
	}
	customDataStoreEmptyStateObj = types.ObjectNull(customDataStoreAttrType)
)

func toSchemaCustomDataStore() schema.SingleNestedAttribute {
	customDataStoreSchema := schema.SingleNestedAttribute{}
	customDataStoreSchema.Description = "A custom data store."
	customDataStoreSchema.Default = objectdefault.StaticValue(types.ObjectNull(customDataStoreAttrType))
	customDataStoreSchema.Computed = true
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

func addOptionalCustomDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.CustomDataStore, plan dataStoreResourceModel) error {
	customDataStorePlan := plan.CustomDataStore.Attributes()

	if internaltypes.IsDefined(plan.CustomId) {
		addRequest.CustomDataStore.Id = plan.CustomId.ValueStringPointer()
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

func toStateCustomDataStore(con context.Context, clientValue *client.DataStoreAggregation, plan basetypes.ObjectValue) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics
	customDataStore := *clientValue.CustomDataStore
	pluginDescriptorRef, diags := resourcelink.ToState(con, &customDataStore.PluginDescriptorRef)
	allDiags = append(allDiags, diags...)
	parentRef, diags := resourcelink.ToState(con, customDataStore.ParentRef)
	allDiags = append(allDiags, diags...)
	configurationObject := func() (types.Object, diag.Diagnostics) {
		planConfiguration, ok := plan.Attributes()["configuration"]
		if ok {
			return pluginconfiguration.ToState(planConfiguration.(types.Object), &customDataStore.Configuration)
		} else {
			return pluginconfiguration.ToState(types.ObjectNull(pluginconfiguration.AttrType()), &customDataStore.Configuration)
		}
	}
	configurationToState, diags := configurationObject()
	allDiags = append(allDiags, diags...)
	customDataStoreVal := map[string]attr.Value{
		"type":                  types.StringValue(customDataStore.Type),
		"name":                  types.StringValue(customDataStore.Name),
		"plugin_descriptor_ref": pluginDescriptorRef,
		"parent_ref":            parentRef,
		"configuration":         configurationToState,
	}
	customDataStoreObj, diags := types.ObjectValue(customDataStoreAttrType, customDataStoreVal)
	allDiags = append(allDiags, diags...)
	return customDataStoreObj, allDiags
}

func readCustomDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreResourceModel, plan *types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.CustomDataStore.Id)
	state.CustomId = types.StringPointerValue(r.CustomDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.CustomDataStore.MaskAttributeValues)
	state.JdbcDataStore = jdbcDataStoreEmptyStateObj
	state.CustomDataStore, diags = toStateCustomDataStore(ctx, r, *plan)
	state.LdapDataStore = ldapDataStoreEmptyStateObj
	state.PingOneLdapGatewayDataStore = pingOneLdapGatewayDataStoreEmptyStateObj
	return diags
}

func createCustomDataStore(plan dataStoreResourceModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
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

	apiCreateDataStore := dsr.apiClient.DataStoresAPI.CreateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig))
	apiCreateDataStore = apiCreateDataStore.Body(createCustomDataStore)
	customDataStoreResponse, httpResp, err := dsr.apiClient.DataStoresAPI.CreateDataStoreExecute(apiCreateDataStore)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResp)
		return
	}

	// Read the response into the state
	var state dataStoreResourceModel
	diags = readCustomDataStoreResponse(con, customDataStoreResponse, &state, &plan.CustomDataStore)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updateCustomDataStore(plan dataStoreResourceModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
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

	updateCustomDataStoreRequest := dsr.apiClient.DataStoresAPI.UpdateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig), plan.Id.ValueString())
	updateCustomDataStoreRequest = updateCustomDataStoreRequest.Body(updateCustomDataStore)
	updateCustomDataStoreResponse, httpResp, err := dsr.apiClient.DataStoresAPI.UpdateDataStoreExecute(updateCustomDataStoreRequest)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating DataStore", err, httpResp)
		return
	}

	// Read the response
	var state dataStoreResourceModel
	diags = readCustomDataStoreResponse(con, updateCustomDataStoreResponse, &state, &plan.CustomDataStore)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
