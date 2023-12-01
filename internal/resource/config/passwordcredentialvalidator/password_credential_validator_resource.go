package passwordcredentialvalidator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithConfigure   = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithImportState = &passwordCredentialValidatorResource{}
)

// PasswordCredentialValidatorResource is a helper function to simplify the provider implementation.
func PasswordCredentialValidatorResource() resource.Resource {
	return &passwordCredentialValidatorResource{}
}

// passwordCredentialValidatorResource is the resource implementation.
type passwordCredentialValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *passwordCredentialValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a password credential validator plugin instance.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: resourcelink.ToSchema(),
			},
			"configuration": pluginconfiguration.ToSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the password credential validator provides.",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of read-only attributes that are automatically populated by the password credential validator descriptor.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional attributes that can be returned by the password credential validator. The extended attributes are only used if the adapter supports them.",
						Computed:    true,
						Optional:    true,
						Default:     listdefault.StaticValue(emptyAttrList),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"validator_id",
		true,
		"The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.")
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *passwordCredentialValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_credential_validator"
}

func (r *passwordCredentialValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *passwordCredentialValidatorResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model passwordCredentialValidatorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	pluginDescriptorRefId := model.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	tableNames := []string{}
	nestedTableFieldNameValues := []string{}
	var rowCount int
	configuration := model.Configuration.Attributes()
	tables := configuration["tables"].(types.List).Elements()
	for _, table := range tables {
		tableAttrs := table.(types.Object).Attributes()
		tableName := tableAttrs["name"].(types.String).ValueString()
		tableNames = append(tableNames, tableName)
		if tableName == "Users" {
			tableRow := tableAttrs["rows"].(types.List).Elements()
			rowCount = len(tableRow)
			for _, row := range tableRow {
				rowAttrs := row.(types.Object).Attributes()
				fields := rowAttrs["fields"].(types.List).Elements()
				for _, field := range fields {
					fieldRow := field.(types.Object).Attributes()
					nestedTableFieldName := fieldRow["name"].(types.String).ValueString()
					nestedTableFieldNameValues = append(nestedTableFieldNameValues, nestedTableFieldName)
				}
			}
		}
	}

	fieldNameValues := []string{}
	fields := configuration["fields"].(types.List).Elements()
	for _, field := range fields {
		field := field.(types.Object).Attributes()
		fieldName := field["name"].(types.String).ValueString()
		fieldNameValues = append(fieldNameValues, fieldName)
	}

	switch pluginDescriptorRefId {
	case "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator":
		hasClientId, clientIdCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "Client Id", 1)
		if !hasClientId || clientIdCount == 0 {
			resp.Diagnostics.AddError("The \"Client Id\" field is required for the PingOne for Enterprise Directory Password Credential Validator", "")
		}
		hasClientSecret, clientSecretCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "Client Secret", 1)
		if !hasClientSecret || clientSecretCount == 0 {
			resp.Diagnostics.AddError("The \"Client Secret\" field is required for the PingOne for Enterprise Directory Password Credential Validator", "")
		}

	case "com.pingidentity.plugins.pcvs.p14c.PingOneForCustomersPCV":
		hasPingOneForCustomersDs, pingOneForCustomersDsCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "PingOne for Customers Datastore", 1)
		if !hasPingOneForCustomersDs || pingOneForCustomersDsCount == 0 {
			resp.Diagnostics.AddError("The \"PingOne for Customers Datastore\" field is required for the PingOne Password Credential Validator", "")
		}

	case "com.pingidentity.plugins.pcvs.pingid.PingIdPCV":
		hasAuthenticationDuringErrors, authenticationDuringErrorsCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "Authentication During Errors", 1)
		if !hasAuthenticationDuringErrors || authenticationDuringErrorsCount == 0 {
			resp.Diagnostics.AddError("The \"Authentication During Errors\" field is required for the PingID Password Credential Validator", "")
		}

	case "org.sourceid.saml20.domain.LDAPUsernamePasswordCredentialValidator":
		hasLdapDs, ldapDsCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "LDAP Datastore", 1)
		if !hasLdapDs || ldapDsCount == 0 {
			resp.Diagnostics.AddError("The \"LDAP Datastore\" field is required for the LDAP Username Password Credential Validator", "")
		}
		hasSearchBase, searchBaseCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "Search Base", 1)
		if !hasSearchBase || searchBaseCount == 0 {
			resp.Diagnostics.AddError("The \"Search Base\" field is required for the LDAP Username Password Credential Validator", "")
		}
		hasSearchFilter, searchFilterCount := internaltypes.ValidateStringCountInSlice(fieldNameValues, "Search Filter", 1)
		if !hasSearchFilter || searchFilterCount == 0 {
			resp.Diagnostics.AddError("The \"Search Filter\" field is required for the LDAP Username Password Credential Validator", "")
		}

	case "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator":
		hasRadiusDs, radiusDsCount := internaltypes.ValidateStringCountInSlice(tableNames, "RADIUS Servers", 1)
		if !hasRadiusDs || radiusDsCount == 0 {
			resp.Diagnostics.AddError("At least one \"RADIUS Servers\" table is required for the RADIUS Username Password Credential Validator", "")
		}

	case "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator":
		hasCorrectUsernameCount, usernameCount := internaltypes.ValidateStringCountInSlice(nestedTableFieldNameValues, "Username", rowCount)
		if !hasCorrectUsernameCount || usernameCount == 0 {
			resp.Diagnostics.AddError("The \"Username\" field is required for the Simple Username Password Credential Validator", "")
		}
		hasCorrectPasswordCount, passwordCount := internaltypes.ValidateStringCountInSlice(nestedTableFieldNameValues, "Password", rowCount)
		if !hasCorrectPasswordCount || passwordCount == 0 {
			resp.Diagnostics.AddError("The \"Password\" field is required for the Simple Username Password Credential Validator", "")
		}
		hasCorrectConfirmPasswordCount, confirmPasswordCount := internaltypes.ValidateStringCountInSlice(nestedTableFieldNameValues, "Confirm Password", rowCount)
		if !hasCorrectConfirmPasswordCount || confirmPasswordCount == 0 {
			resp.Diagnostics.AddError("The \"Confirm Password\" field is required for the Simple Username Password Credential Validator", "")
		}
		return
	}
}

func addOptionalPasswordCredentialValidatorFields(ctx context.Context, addRequest *client.PasswordCredentialValidator, plan passwordCredentialValidatorModel) error {
	if internaltypes.IsDefined(plan.ParentRef) {
		if plan.ParentRef.Attributes()["id"].(types.String).ValueString() != "" {
			addRequest.ParentRef = client.NewResourceLinkWithDefaults()
			addRequest.ParentRef.Id = plan.ParentRef.Attributes()["id"].(types.String).ValueString()
			err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef, true)), addRequest.ParentRef)
			if err != nil {
				return err
			}
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = client.NewPasswordCredentialValidatorAttributeContractWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, true)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
		extendedAttrsLength := len(plan.AttributeContract.Attributes()["extended_attributes"].(types.List).Elements())
		if extendedAttrsLength == 0 {
			addRequest.AttributeContract.ExtendedAttributes = nil
		}
	}
	return nil
}

func (r *passwordCredentialValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan passwordCredentialValidatorModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefResLink, err := resourcelink.ClientStruct(plan.PluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", err.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfigurationWithDefaults()
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, true)), configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", err.Error())
		return
	}

	createPasswordCredentialValidators := client.NewPasswordCredentialValidator(plan.ValidatorId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err = addOptionalPasswordCredentialValidatorFields(ctx, createPasswordCredentialValidators, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}

	apiCreatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePasswordCredentialValidators = apiCreatePasswordCredentialValidators.Body(*createPasswordCredentialValidators)
	passwordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidatorExecute(apiCreatePasswordCredentialValidators)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating a Password Credential Validator", err, httpResp)
		return
	}

	// Read the response into the state
	var state passwordCredentialValidatorModel

	diags = readPasswordCredentialValidatorResponse(ctx, passwordCredentialValidatorsResponse, &state, plan.Configuration, true)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *passwordCredentialValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state passwordCredentialValidatorModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadPasswordCredentialValidators, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ValidatorId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Password Credential Validator", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Password Credential Validator", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readPasswordCredentialValidatorResponse(ctx, apiReadPasswordCredentialValidators, &state, state.Configuration, true)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *passwordCredentialValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan passwordCredentialValidatorModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefResLink, err := resourcelink.ClientStruct(plan.PluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", err.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfiguration()
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, true)), configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", err.Error())
		return
	}

	updatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ValidatorId.ValueString())
	createUpdateRequest := client.NewPasswordCredentialValidator(plan.ValidatorId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err = addOptionalPasswordCredentialValidatorFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}

	updatePasswordCredentialValidators = updatePasswordCredentialValidators.Body(*createUpdateRequest)
	updatePasswordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidatorExecute(updatePasswordCredentialValidators)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating a Password Credential Validator", err, httpResp)
		return
	}

	// Read the response
	diags = readPasswordCredentialValidatorResponse(ctx, updatePasswordCredentialValidatorsResponse, &plan, plan.Configuration, true)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *passwordCredentialValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state passwordCredentialValidatorModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.DeletePasswordCredentialValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ValidatorId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Password Credential Validator", err, httpResp)
	}
}

func (r *passwordCredentialValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("validator_id"), req, resp)
}
