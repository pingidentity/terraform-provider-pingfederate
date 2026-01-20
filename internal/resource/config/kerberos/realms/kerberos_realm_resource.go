// Copyright © 2025 Ping Identity Corporation

package kerberosrealms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kerberosRealmsResource{}
	_ resource.ResourceWithConfigure   = &kerberosRealmsResource{}
	_ resource.ResourceWithImportState = &kerberosRealmsResource{}

	emptyStringSet, _ = types.SetValue(types.StringType, nil)

	customId = "realm_id"
)

// KerberosRealmsResource is a helper function to simplify the provider implementation.
func KerberosRealmsResource() resource.Resource {
	return &kerberosRealmsResource{}
}

// kerberosRealmsResource is the resource implementation.
type kerberosRealmsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type kerberosRealmsResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	RealmId                            types.String `tfsdk:"realm_id"`
	KerberosRealmName                  types.String `tfsdk:"kerberos_realm_name"`
	ConnectionType                     types.String `tfsdk:"connection_type"`
	KeyDistributionCenters             types.Set    `tfsdk:"key_distribution_centers"`
	KerberosUsername                   types.String `tfsdk:"kerberos_username"`
	KerberosPassword                   types.String `tfsdk:"kerberos_password"`
	KerberosEncryptedPassword          types.String `tfsdk:"kerberos_encrypted_password"`
	RetainPreviousKeysOnPasswordChange types.Bool   `tfsdk:"retain_previous_keys_on_password_change"`
	SuppressDomainNameConcatenation    types.Bool   `tfsdk:"suppress_domain_name_concatenation"`
	LdapGatewayDataStoreRef            types.Object `tfsdk:"ldap_gateway_data_store_ref"`
}

// GetSchema defines the schema for the resource.
func (r *kerberosRealmsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Resource to create and manage a Kerberos Realm",
		Attributes: map[string]schema.Attribute{
			"realm_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the Kerberos Realm. It can be any combination of `[a-zA-Z0-9._-]`. This field is immutable and will trigger a replacement plan if changed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateId(),
				},
			},
			"kerberos_realm_name": schema.StringAttribute{
				Description: "The Domain/Realm name used for display in UI screens.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"connection_type": schema.StringAttribute{
				Description: "Controls how PingFederate connects to the Active Directory/Kerberos Realm. Options are `DIRECT`, `LDAP_GATEWAY`, `LOCAL_VALIDATION`. The default is `DIRECT`. `LOCAL_VALIDATION` only supported in PF version `12.2` or later.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DIRECT"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DIRECT", "LDAP_GATEWAY", "LOCAL_VALIDATION"}...),
				},
			},
			"key_distribution_centers": schema.SetAttribute{
				Description: "The Domain Controller/Key Distribution Center Host Action Names. Only applicable when `connection_type` is `DIRECT`.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(emptyStringSet),
			},
			"kerberos_username": schema.StringAttribute{
				Description: "The Domain/Realm username. Required when `connection_type` is `DIRECT` or `LOCAL_VALIDATION`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"kerberos_password": schema.StringAttribute{
				Description: "The Domain/Realm password. Required when `connection_type` is `DIRECT` or `LOCAL_VALIDATION`. Only one of this attribute and 'kerberos_encrypted_password' should be specified.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"kerberos_encrypted_password": schema.StringAttribute{
				Description: "The encrypted Domain/Realm password. Required when `connection_type` is `DIRECT` or `LOCAL_VALIDATION`. Only one of this attribute and 'kerberos_password' should be specified.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("kerberos_password")),
				},
			},
			// Computed due to dependency on connection_type, this value is not present when connection_type is LDAP_GATEWAY, default set in ModifyPlan
			"retain_previous_keys_on_password_change": schema.BoolAttribute{
				Description: "Determines whether the previous encryption keys are retained when the password is updated. Retaining the previous keys allows existing Kerberos tickets to continue to be validated. The default is `false`. Only applicable when `connection_type` is `DIRECT` or `LOCAL_VALIDATION`.",
				Computed:    true,
				Optional:    true,
			},
			// Computed due to dependency on connection_type, this value is not present when connection_type is LDAP_GATEWAY, default set in ModifyPlan
			"suppress_domain_name_concatenation": schema.BoolAttribute{
				Description: "Controls whether the KDC hostnames and the realm name are concatenated in the auto-generated `krb5.conf` file. Only applicable when `connection_type` is `DIRECT`.",
				Computed:    true,
				Optional:    true,
			},
			"ldap_gateway_data_store_ref": schema.SingleNestedAttribute{
				Description: "The LDAP gateway used by PingFederate to communicate with the Active Directory. Required when `connection_type` is `LDAP_GATEWAY`.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *kerberosRealmsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_realm"
}

func (r *kerberosRealmsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *kerberosRealmsResource) validatePlan(ctx context.Context, plan *kerberosRealmsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	errorMsg := "is only applicable when connection_type is set to \"DIRECT\" or \"LOCAL_VALIDATION\"."
	if plan.ConnectionType.ValueString() != "DIRECT" && plan.ConnectionType.ValueString() != "LOCAL_VALIDATION" {
		if internaltypes.IsDefined(plan.KerberosUsername) {
			diags.AddAttributeError(path.Root("kerberos_username"),
				providererror.InvalidAttributeConfiguration, "kerberos_username "+errorMsg)
		}
		if internaltypes.IsDefined(plan.KerberosPassword) {
			diags.AddAttributeError(path.Root("kerberos_password"),
				providererror.InvalidAttributeConfiguration, "kerberos_password "+errorMsg)
		}
		if internaltypes.IsDefined(plan.KerberosEncryptedPassword) {
			diags.AddAttributeError(path.Root("kerberos_encrypted_password"),
				providererror.InvalidAttributeConfiguration, "kerberos_encrypted_password "+errorMsg)
		}
		if internaltypes.IsDefined(plan.RetainPreviousKeysOnPasswordChange) {
			diags.AddAttributeError(path.Root("retain_previous_keys_on_password_change"),
				providererror.InvalidAttributeConfiguration, "retain_previous_keys_on_password_change "+errorMsg)
		}
	}
	errorMsg = "is only applicable when connection_type is set to \"DIRECT\"."
	if plan.ConnectionType.ValueString() != "DIRECT" {
		if internaltypes.IsDefined(plan.SuppressDomainNameConcatenation) {
			diags.AddAttributeError(path.Root("suppress_domain_name_concatenation"),
				providererror.InvalidAttributeConfiguration, "suppress_domain_name_concatenation "+errorMsg)
		}
		if len(plan.KeyDistributionCenters.Elements()) > 0 {
			diags.AddAttributeError(path.Root("key_distribution_centers"),
				providererror.InvalidAttributeConfiguration, "key_distribution_centers "+errorMsg)
		}
	}

	if plan.ConnectionType.ValueString() == "DIRECT" || plan.ConnectionType.ValueString() == "LOCAL_VALIDATION" {
		if plan.KerberosUsername.IsNull() {
			diags.AddAttributeError(
				path.Root("kerberos_username"),
				providererror.InvalidAttributeConfiguration,
				"kerberos_username is required when connection_type is set to \"DIRECT\" or \"LOCAL_VALIDATON\".")
		}
		if plan.KerberosPassword.IsNull() && !internaltypes.IsDefined(plan.KerberosEncryptedPassword) {
			diags.AddAttributeError(
				path.Root("kerberos_password"),
				providererror.InvalidAttributeConfiguration,
				"kerberos_password or kerberos_encrypted_password is required when connection_type is set to \"DIRECT\" or \"LOCAL_VALIDATON\".")
		}
	}

	// ldap_gateway_data_store_ref is required when connection_type is set to LDAP_GATEWAY
	if plan.ConnectionType.ValueString() == "LDAP_GATEWAY" {
		if plan.LdapGatewayDataStoreRef.IsNull() {
			diags.AddAttributeError(
				path.Root("ldap_gateway_data_store_ref"),
				providererror.InvalidAttributeConfiguration,
				"ldap_gateway_data_store_ref is required when connection_type is set to \"LDAP_GATEWAY\".")
		}
	}
	return diags
}

func (r *kerberosRealmsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *kerberosRealmsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}
	if plan.ConnectionType.ValueString() == "DIRECT" {
		if plan.RetainPreviousKeysOnPasswordChange.IsUnknown() {
			plan.RetainPreviousKeysOnPasswordChange = types.BoolValue(false)
		}
		if plan.SuppressDomainNameConcatenation.IsUnknown() {
			plan.SuppressDomainNameConcatenation = types.BoolValue(false)
		}
	}
	resp.Diagnostics.Append(r.validatePlan(ctx, plan)...)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func readKerberosRealmsResponse(ctx context.Context, r *client.KerberosRealm, state *kerberosRealmsResourceModel, plan *kerberosRealmsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Id = types.StringPointerValue(r.Id)
	state.RealmId = types.StringPointerValue(r.Id)
	state.KerberosRealmName = types.StringValue(r.KerberosRealmName)
	state.ConnectionType = types.StringPointerValue(r.ConnectionType)
	state.KeyDistributionCenters = internaltypes.GetStringSet(r.KeyDistributionCenters)
	state.KerberosUsername = types.StringPointerValue(r.KerberosUsername)
	state.KerberosPassword = types.StringPointerValue(plan.KerberosPassword.ValueStringPointer())
	if internaltypes.IsDefined(plan.KerberosEncryptedPassword) {
		state.KerberosEncryptedPassword = types.StringValue(plan.KerberosEncryptedPassword.ValueString())
	} else {
		state.KerberosEncryptedPassword = types.StringPointerValue(r.KerberosEncryptedPassword)
	}
	state.RetainPreviousKeysOnPasswordChange = types.BoolPointerValue(r.RetainPreviousKeysOnPasswordChange)
	state.SuppressDomainNameConcatenation = types.BoolPointerValue(r.SuppressDomainNameConcatenation)
	state.LdapGatewayDataStoreRef, diags = resourcelink.ToState(ctx, r.LdapGatewayDataStoreRef)

	return diags
}

func addOptionalKerberosRealmsFields(ctx context.Context, addRequest *client.KerberosRealm, plan kerberosRealmsResourceModel) error {
	var err error

	//  realm_id is a required field, so we need to set the Id to this value
	addRequest.Id = plan.RealmId.ValueStringPointer()
	addRequest.ConnectionType = plan.ConnectionType.ValueStringPointer()
	addRequest.KerberosUsername = plan.KerberosUsername.ValueStringPointer()
	addRequest.KerberosPassword = plan.KerberosPassword.ValueStringPointer()
	addRequest.KerberosEncryptedPassword = plan.KerberosEncryptedPassword.ValueStringPointer()

	if len(plan.KeyDistributionCenters.Elements()) > 0 {
		var slice []string
		plan.KeyDistributionCenters.ElementsAs(ctx, &slice, false)
		addRequest.KeyDistributionCenters = slice
	}

	addRequest.LdapGatewayDataStoreRef, err = resourcelink.ClientStruct(plan.LdapGatewayDataStoreRef)
	if err != nil {
		return err
	}

	// These are optional fields based on connection_type, so we need to check if they are defined before adding them to the request
	if internaltypes.IsDefined(plan.RetainPreviousKeysOnPasswordChange) {
		addRequest.RetainPreviousKeysOnPasswordChange = plan.RetainPreviousKeysOnPasswordChange.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressDomainNameConcatenation) {
		addRequest.SuppressDomainNameConcatenation = plan.SuppressDomainNameConcatenation.ValueBoolPointer()
	}

	return nil
}

func (r *kerberosRealmsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kerberosRealmsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createKerberosRealms := client.NewKerberosRealm(plan.KerberosRealmName.ValueString())
	err := addOptionalKerberosRealmsFields(ctx, createKerberosRealms, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for a kerberos realm: "+err.Error())
		return
	}

	apiCreateKerberosRealms := r.apiClient.KerberosRealmsAPI.CreateKerberosRealm(config.AuthContext(ctx, r.providerConfig))
	apiCreateKerberosRealms = apiCreateKerberosRealms.Body(*createKerberosRealms)
	kerberosRealmsResponse, httpResp, err := r.apiClient.KerberosRealmsAPI.CreateKerberosRealmExecute(apiCreateKerberosRealms)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating a kerberos realm", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state kerberosRealmsResourceModel

	diags = readKerberosRealmsResponse(ctx, kerberosRealmsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *kerberosRealmsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kerberosRealmsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadKerberosRealms, httpResp, err := r.apiClient.KerberosRealmsAPI.GetKerberosRealm(config.AuthContext(ctx, r.providerConfig), state.RealmId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Kerberos Realm", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting a kerberos realm", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	readKerberosRealmsResponse(ctx, apiReadKerberosRealms, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kerberosRealmsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan kerberosRealmsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateKerberosRealms := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealm(config.AuthContext(ctx, r.providerConfig), plan.RealmId.ValueString())
	createUpdateRequest := client.NewKerberosRealm(plan.KerberosRealmName.ValueString())
	err := addOptionalKerberosRealmsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for a kerberos realm: "+err.Error())
		return
	}

	updateKerberosRealms = updateKerberosRealms.Body(*createUpdateRequest)
	updateKerberosRealmsResponse, httpResp, err := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealmExecute(updateKerberosRealms)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating a kerberos realm", err, httpResp, &customId)
		return
	}

	// Read the response
	var state kerberosRealmsResourceModel
	diags = readKerberosRealmsResponse(ctx, updateKerberosRealmsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *kerberosRealmsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kerberosRealmsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.KerberosRealmsAPI.DeleteKerberosRealm(config.AuthContext(ctx, r.providerConfig), state.RealmId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting a kerberos realm", err, httpResp, &customId)
	}
}

func (r *kerberosRealmsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("realm_id"), req, resp)
}
