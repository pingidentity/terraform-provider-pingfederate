package datastore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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
			Required:    true,
			Description: "Reference to the PingOne connection this gateway uses.",
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

func toStatePingOneLdapGatewayDataStore(con context.Context, clientValue *client.DataStoreAggregation, plan basetypes.ObjectValue) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics
	pingOneLdapGatewayDataStoreAttrType := map[string]attr.Type{
		"ping_one_connection_ref":  basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"ldap_type":                basetypes.StringType{},
		"ping_one_ldap_gateway_id": basetypes.StringType{},
		"use_ssl":                  basetypes.BoolType{},
		"name":                     basetypes.StringType{},
		"binary_attributes":        basetypes.SetType{ElemType: basetypes.StringType{}},
		"type":                     basetypes.StringType{},
		"ping_one_environment_id":  basetypes.StringType{},
	}

	clientVals := clientValue.PingOneLdapGatewayDataStore
	pingOneConnectionRef, diags := resourcelink.ToState(con, &clientVals.PingOneConnectionRef)
	allDiags = append(allDiags, diags...)
	binaryAttributes := clientVals.BinaryAttributes
	binaryAttributesVal := func() types.Set {
		if len(binaryAttributes) == 0 {
			return types.SetNull(types.StringType)
		} else {
			return internaltypes.GetStringSet(clientVals.BinaryAttributes)
		}
	}
	pingOneLdapGatewayDataStoreVal := map[string]attr.Value{
		"ping_one_connection_ref":  pingOneConnectionRef,
		"ldap_type":                types.StringValue(clientVals.LdapType),
		"ping_one_ldap_gateway_id": types.StringValue(clientVals.PingOneLdapGatewayId),
		"use_ssl":                  types.BoolPointerValue(clientVals.UseSsl),
		"name":                     types.StringPointerValue(clientVals.Name),
		"binary_attributes":        binaryAttributesVal(),
		"type":                     types.StringValue("PING_ONE_LDAP_GATEWAY"),
		"ping_one_environment_id":  types.StringValue(clientVals.PingOneEnvironmentId),
	}

	pingOneLdapGatewayDataStoreObj, diags := types.ObjectValue(pingOneLdapGatewayDataStoreAttrType, pingOneLdapGatewayDataStoreVal)
	allDiags = append(allDiags, diags...)
	return pingOneLdapGatewayDataStoreObj, allDiags
}

func readPingOneLdapGatewayDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreResourceModel, plan *types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.PingOneLdapGatewayDataStore.Id)
	state.CustomId = types.StringPointerValue(r.PingOneLdapGatewayDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.PingOneLdapGatewayDataStore.MaskAttributeValues)
	state.CustomDataStore = customDataStoreEmptyStateObj
	state.JdbcDataStore = jdbcDataStoreEmptyStateObj
	state.LdapDataStore = ldapDataStoreEmptyStateObj
	state.PingOneLdapGatewayDataStore, diags = toStatePingOneLdapGatewayDataStore(ctx, r, *plan)
	return diags
}

func addOptionalPingOneLdapGatewayDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.PingOneLdapGatewayDataStore, plan dataStoreResourceModel) error {
	pingOneLdapGatewayDataStorePlan := plan.PingOneLdapGatewayDataStore.Attributes()

	if internaltypes.IsDefined(pingOneLdapGatewayDataStorePlan["use_ssl"]) {
		addRequest.PingOneLdapGatewayDataStore.UseSsl = pingOneLdapGatewayDataStorePlan["use_ssl"].(types.Bool).ValueBoolPointer()
	}

	if internaltypes.IsDefined(pingOneLdapGatewayDataStorePlan["name"]) {
		addRequest.PingOneLdapGatewayDataStore.Name = pingOneLdapGatewayDataStorePlan["name"].(types.String).ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.MaskAttributeValues) {
		addRequest.PingOneLdapGatewayDataStore.MaskAttributeValues = plan.MaskAttributeValues.ValueBoolPointer()
	}

	if internaltypes.IsDefined(pingOneLdapGatewayDataStorePlan["binary_attributes"]) {
		addRequest.PingOneLdapGatewayDataStore.BinaryAttributes = internaltypes.SetTypeToStringSet(pingOneLdapGatewayDataStorePlan["binary_attributes"].(types.Set))
	}

	if internaltypes.IsDefined(plan.CustomId) {
		addRequest.PingOneLdapGatewayDataStore.Id = plan.CustomId.ValueStringPointer()
	}

	return nil
}

func createPingOneLdapGatewayDataStore(plan dataStoreResourceModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	pingOneLdapGatewayDsPlan := plan.PingOneLdapGatewayDataStore.Attributes()
	ldapType := pingOneLdapGatewayDsPlan["ldap_type"].(types.String).ValueString()
	pingOneConnectionRef, err := resourcelink.ClientStruct(plan.PingOneLdapGatewayDataStore.Attributes()["ping_one_connection_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert ping_one_connection_ref to PingOneConnectionRef", err.Error())
		return
	}
	pingOneEnvId := pingOneLdapGatewayDsPlan["ping_one_environment_id"].(types.String).ValueString()
	pingOneGatewayId := pingOneLdapGatewayDsPlan["ping_one_ldap_gateway_id"].(types.String).ValueString()
	createPingOneLdapGatewayDataStore := client.PingOneLdapGatewayDataStoreAsDataStoreAggregation(client.NewPingOneLdapGatewayDataStore(
		ldapType,
		*pingOneConnectionRef,
		pingOneEnvId,
		pingOneGatewayId,
		"PING_ONE_LDAP_GATEWAY",
	))
	err = addOptionalPingOneLdapGatewayDataStoreFields(createPingOneLdapGatewayDataStore, con, client.PingOneLdapGatewayDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	apiCreateDataStore := dsr.apiClient.DataStoresAPI.CreateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig))
	apiCreateDataStore = apiCreateDataStore.Body(createPingOneLdapGatewayDataStore)
	ldapDataStoreResponse, httpResp, err := dsr.apiClient.DataStoresAPI.CreateDataStoreExecute(apiCreateDataStore)
	if err != nil {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResp)
		return
	}

	// Read the response into the state
	var state dataStoreResourceModel
	diags = readPingOneLdapGatewayDataStoreResponse(con, ldapDataStoreResponse, &state, &plan.PingOneLdapGatewayDataStore)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updatePingOneLdapGatewayDataStore(plan dataStoreResourceModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	pingOneLdapGatewayDsPlan := plan.PingOneLdapGatewayDataStore.Attributes()
	ldapType := pingOneLdapGatewayDsPlan["ldap_type"].(types.String).ValueString()
	pingOneConnectionRef, err := resourcelink.ClientStruct(plan.PingOneLdapGatewayDataStore.Attributes()["ping_one_connection_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert ping_one_connection_ref to PingOneConnectionRef", err.Error())
		return
	}
	pingOneEnvId := pingOneLdapGatewayDsPlan["ping_one_environment_id"].(types.String).ValueString()
	pingOneGatewayId := pingOneLdapGatewayDsPlan["ping_one_ldap_gateway_id"].(types.String).ValueString()
	updatePingOneLdapGatewayDataStore := client.PingOneLdapGatewayDataStoreAsDataStoreAggregation(client.NewPingOneLdapGatewayDataStore(
		ldapType,
		*pingOneConnectionRef,
		pingOneEnvId,
		pingOneGatewayId,
		"PING_ONE_LDAP_GATEWAY",
	))

	err = addOptionalPingOneLdapGatewayDataStoreFields(updatePingOneLdapGatewayDataStore, con, client.PingOneLdapGatewayDataStore{}, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for DataStore", err.Error())
		return
	}

	updatePingOneLdapGatewayDataStoreRequest := dsr.apiClient.DataStoresAPI.UpdateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig), plan.Id.ValueString())
	updatePingOneLdapGatewayDataStoreRequest = updatePingOneLdapGatewayDataStoreRequest.Body(updatePingOneLdapGatewayDataStore)
	updatePingOneLdapGatewayDataStoreResponse, httpResp, err := dsr.apiClient.DataStoresAPI.UpdateDataStoreExecute(updatePingOneLdapGatewayDataStoreRequest)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(con, &resp.Diagnostics, "An error occurred while updating DataStore", err, httpResp)
		return
	}

	// Read the response
	var state dataStoreResourceModel
	diags = readPingOneLdapGatewayDataStoreResponse(con, updatePingOneLdapGatewayDataStoreResponse, &state, &plan.PingOneLdapGatewayDataStore)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
