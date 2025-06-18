// Copyright © 2025 Ping Identity Corporation

package datastore

import (
	"context"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	pingOneLdapGatewayDataStoreAttrType = map[string]attr.Type{
		"use_start_tls":            types.BoolType,
		"ping_one_connection_ref":  types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"ldap_type":                types.StringType,
		"ping_one_ldap_gateway_id": types.StringType,
		"use_ssl":                  types.BoolType,
		"name":                     types.StringType,
		"binary_attributes":        types.SetType{ElemType: types.StringType},
		"type":                     types.StringType,
		"ping_one_environment_id":  types.StringType,
	}
	pingOneLdapGatewayDataStoreEmptyStateObj = types.ObjectNull(pingOneLdapGatewayDataStoreAttrType)
)

func toSchemaPingOneLdapGatewayDataStore() schema.SingleNestedAttribute {
	pingOneLdapGatewayDataStoreSchema := schema.SingleNestedAttribute{}
	pingOneLdapGatewayDataStoreSchema.Description = "A PingOne LDAP Gateway data store."
	pingOneLdapGatewayDataStoreSchema.Optional = true
	pingOneLdapGatewayDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("PING_ONE_LDAP_GATEWAY"),
		},
		"name": schema.StringAttribute{
			Description: "The data store name with a unique value across all data sources. Defaults to `ping_one_connection_ref.id` plus `ping_one_environment_id` plus `ping_one_ldap_gateway_id`, each separated by `:`.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"use_start_tls": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using StartTLS. The default value is `false`. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`. Supported in PingFederate `12.1` and later.",
			Computed:    true,
			Optional:    true,
		},
		"ping_one_connection_ref": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Reference to the PingOne connection this gateway uses.",
			Attributes:  resourcelink.ToSchema(),
		},
		"ldap_type": schema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`. Supported values are `ACTIVE_DIRECTORY`, `ORACLE_DIRECTORY_SERVER`, `ORACLE_UNIFIED_DIRECTORY`, `UNBOUNDID_DS`, `PING_DIRECTORY`, and `GENERIC`.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACTIVE_DIRECTORY", "ORACLE_DIRECTORY_SERVER", "ORACLE_UNIFIED_DIRECTORY", "UNBOUNDID_DS", "PING_DIRECTORY", "GENERIC"),
			},
		},
		"ping_one_ldap_gateway_id": schema.StringAttribute{
			Description: "The ID of the PingOne LDAP Gateway this data store uses.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"use_ssl": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is `false`. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"binary_attributes": schema.SetAttribute{
			Description: "A list of LDAP attributes to be handled as binary data.",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
			Default:     setdefault.StaticValue(types.SetNull(types.StringType)),
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"ping_one_environment_id": schema.StringAttribute{
			Description: "The environment ID to which the gateway belongs.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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

func toDataSourceSchemaPingOneLdapGatewayDataStore() datasourceschema.SingleNestedAttribute {
	pingOneLdapGatewayDataStoreSchema := datasourceschema.SingleNestedAttribute{}
	pingOneLdapGatewayDataStoreSchema.Description = "A PingOne LDAP Gateway data store."
	pingOneLdapGatewayDataStoreSchema.Computed = true
	pingOneLdapGatewayDataStoreSchema.Optional = false
	pingOneLdapGatewayDataStoreSchema.Attributes = map[string]datasourceschema.Attribute{
		"type": datasourceschema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The data store name with a unique value across all data sources.",
			Computed:    true,
			Optional:    false,
		},
		"use_start_tls": schema.BoolAttribute{
			Description: "Connects to the LDAP data store using StartTLS.",
			Computed:    true,
			Optional:    false,
		},
		"ping_one_connection_ref": datasourceschema.SingleNestedAttribute{
			Computed:    true,
			Optional:    false,
			Description: "Reference to the PingOne connection this gateway uses.",
			Attributes:  datasourceresourcelink.ToDataSourceSchema(),
		},
		"ldap_type": datasourceschema.StringAttribute{
			Description: "A type that allows PingFederate to configure many provisioning settings automatically.",
			Computed:    true,
			Optional:    false,
		},
		"ping_one_ldap_gateway_id": datasourceschema.StringAttribute{
			Description: "The ID of the PingOne LDAP Gateway this data store uses.",
			Computed:    true,
			Optional:    false,
		},
		"use_ssl": datasourceschema.BoolAttribute{
			Description: "Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS).",
			Computed:    true,
			Optional:    false,
		},
		"binary_attributes": datasourceschema.SetAttribute{
			Description: "A list of LDAP attributes to be handled as binary data.",
			Computed:    true,
			Optional:    false,
			ElementType: types.StringType,
		},
		"ping_one_environment_id": datasourceschema.StringAttribute{
			Description: "The environment ID to which the gateway belongs.",
			Computed:    true,
			Optional:    false,
		},
	}

	return pingOneLdapGatewayDataStoreSchema
}

func toStatePingOneLdapGatewayDataStore(con context.Context, pingOneLdapGDS *client.PingOneLdapGatewayDataStore, plan types.Object) (types.Object, diag.Diagnostics) {
	var diags, allDiags diag.Diagnostics

	if pingOneLdapGDS == nil {
		diags.AddError(providererror.InternalProviderError, "Failed to read PingOne LDAP Gateway data store from PingFederate. The response from PingFederate was nil.")
		return pingOneLdapGatewayDataStoreEmptyStateObj, diags
	}

	pingOneConRefFromClient := pingOneLdapGDS.GetPingOneConnectionRef()
	pingOneConnectionRef, diags := resourcelink.ToState(con, &pingOneConRefFromClient)
	allDiags = append(allDiags, diags...)
	binaryAttributes := func() basetypes.SetValue {
		if len(pingOneLdapGDS.BinaryAttributes) > 0 {
			return internaltypes.GetStringSet(pingOneLdapGDS.BinaryAttributes)
		} else {
			return types.SetNull(types.StringType)
		}
	}
	pingOneLdapGatewayDataStoreVal := map[string]attr.Value{
		"use_start_tls":            types.BoolPointerValue(pingOneLdapGDS.UseStartTLS),
		"ping_one_connection_ref":  pingOneConnectionRef,
		"ldap_type":                types.StringValue(pingOneLdapGDS.GetLdapType()),
		"ping_one_ldap_gateway_id": types.StringValue(pingOneLdapGDS.GetPingOneLdapGatewayId()),
		"use_ssl":                  types.BoolValue(pingOneLdapGDS.GetUseSsl()),
		"name":                     types.StringValue(pingOneLdapGDS.GetName()),
		"binary_attributes":        binaryAttributes(),
		"type":                     types.StringValue("PING_ONE_LDAP_GATEWAY"),
		"ping_one_environment_id":  types.StringValue(pingOneLdapGDS.GetPingOneEnvironmentId()),
	}

	pingOneLdapGatewayDataStoreObj, diags := types.ObjectValue(pingOneLdapGatewayDataStoreAttrType, pingOneLdapGatewayDataStoreVal)
	allDiags = append(allDiags, diags...)
	return pingOneLdapGatewayDataStoreObj, allDiags
}

func readPingOneLdapGatewayDataStoreResponse(ctx context.Context, r *client.DataStoreAggregation, state *dataStoreModel, plan *types.Object, isResource bool) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringPointerValue(r.PingOneLdapGatewayDataStore.Id)
	state.DataStoreId = types.StringPointerValue(r.PingOneLdapGatewayDataStore.Id)
	state.MaskAttributeValues = types.BoolPointerValue(r.PingOneLdapGatewayDataStore.MaskAttributeValues)
	if isResource {
		state.CustomDataStore = customDataStoreEmptyStateObj
		state.JdbcDataStore = jdbcDataStoreEmptyStateObj
		state.LdapDataStore = ldapDataStoreEmptyStateObj
	} else {
		state.CustomDataStore = customDataStoreEmptyDataSourceStateObj
		state.JdbcDataStore = jdbcDataStoreEmptyDataSourceStateObj
		state.LdapDataStore = ldapDataStoreEmptyDataSourceStateObj
	}
	state.PingOneLdapGatewayDataStore, diags = toStatePingOneLdapGatewayDataStore(ctx, r.PingOneLdapGatewayDataStore, *plan)
	return diags
}

func addOptionalPingOneLdapGatewayDataStoreFields(addRequest client.DataStoreAggregation, con context.Context, createJdbcDataStore client.PingOneLdapGatewayDataStore, plan dataStoreModel) error {
	pingOneLdapGatewayDataStorePlan := plan.PingOneLdapGatewayDataStore.Attributes()

	if internaltypes.IsDefined(pingOneLdapGatewayDataStorePlan["use_start_tls"]) {
		addRequest.PingOneLdapGatewayDataStore.UseStartTLS = pingOneLdapGatewayDataStorePlan["use_start_tls"].(types.Bool).ValueBoolPointer()
	}

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
		addRequest.PingOneLdapGatewayDataStore.BinaryAttributes = internaltypes.SetTypeToStringSlice(pingOneLdapGatewayDataStorePlan["binary_attributes"].(types.Set))
	}

	if internaltypes.IsDefined(plan.DataStoreId) {
		addRequest.PingOneLdapGatewayDataStore.Id = plan.DataStoreId.ValueStringPointer()
	}

	return nil
}

func createPingOneLdapGatewayDataStore(plan dataStoreModel, con context.Context, req resource.CreateRequest, resp *resource.CreateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	pingOneLdapGatewayDsPlan := plan.PingOneLdapGatewayDataStore.Attributes()
	ldapType := pingOneLdapGatewayDsPlan["ldap_type"].(types.String).ValueString()
	pingOneConnectionRef, err := resourcelink.ClientStruct(plan.PingOneLdapGatewayDataStore.Attributes()["ping_one_connection_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to convert ping_one_connection_ref to PingOneConnectionRef: "+err.Error())
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
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for DataStore: "+err.Error())
		return
	}

	response, httpResponse, err := createDataStore(createPingOneLdapGatewayDataStore, dsr, con, resp)
	if err != nil {
		config.ReportHttpErrorCustomId(con, &resp.Diagnostics, "An error occurred while creating the DataStore", err, httpResponse, &customId)
		return
	}

	// Read the response into the state
	var state dataStoreModel
	diags = readPingOneLdapGatewayDataStoreResponse(con, response, &state, &plan.PingOneLdapGatewayDataStore, true)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}

func updatePingOneLdapGatewayDataStore(plan dataStoreModel, con context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, dsr *dataStoreResource) {
	var diags diag.Diagnostics
	var err error

	pingOneLdapGatewayDsPlan := plan.PingOneLdapGatewayDataStore.Attributes()
	ldapType := pingOneLdapGatewayDsPlan["ldap_type"].(types.String).ValueString()
	pingOneConnectionRef, err := resourcelink.ClientStruct(plan.PingOneLdapGatewayDataStore.Attributes()["ping_one_connection_ref"].(types.Object))
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to convert ping_one_connection_ref to PingOneConnectionRef: "+err.Error())
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
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the DataStore: "+err.Error())
		return
	}

	response, httpResp, err := updateDataStore(updatePingOneLdapGatewayDataStore, dsr, con, resp, plan.DataStoreId.ValueString())
	if err != nil {
		config.ReportHttpErrorCustomId(con, &resp.Diagnostics, "An error occurred while updating the DataStore", err, httpResp, &customId)
		return
	}

	// Read the response
	var state dataStoreModel
	diags = readPingOneLdapGatewayDataStoreResponse(con, response, &state, &plan.PingOneLdapGatewayDataStore, true)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(con, state)
	resp.Diagnostics.Append(diags...)
}
