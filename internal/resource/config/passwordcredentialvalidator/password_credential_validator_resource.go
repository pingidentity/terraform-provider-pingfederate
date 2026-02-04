// Copyright © 2025 Ping Identity Corporation

package passwordcredentialvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithConfigure   = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithImportState = &passwordCredentialValidatorResource{}

	customId = "validator_id"
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
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"configuration": pluginconfiguration.ToSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the password credential validator provides.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of read-only attributes that are automatically populated by the password credential validator descriptor.",
						Computed:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseNonNullStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the password credential validator. The extended attributes are only used if the adapter supports them.",
						Computed:    true,
						Optional:    true,
						Default:     setdefault.StaticValue(emptyAttrSet),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"validator_id",
		true,
		true,
		"The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed. Must be less than 33 characters, contain no spaces, and be alphanumeric.")
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
	var model *passwordCredentialValidatorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if model == nil {
		return
	}

	if model.Configuration.IsUnknown() || model.PluginDescriptorRef.IsUnknown() {
		return
	}

	configuration := model.Configuration.Attributes()

	pluginDescriptorRefId := model.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	tablesKnown := !configuration["tables"].IsUnknown()
	var isRadiusServerTableFound bool
	if pluginDescriptorRefId == "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator" || pluginDescriptorRefId == "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator" {
		if tablesKnown {
			tables := configuration["tables"].(types.List).Elements()
			for _, table := range tables {
				tableAttrs := table.(types.Object).Attributes()
				if tableAttrs["name"].IsUnknown() {
					tablesKnown = false
				}
				tableName := tableAttrs["name"].(types.String).ValueString()
				isRadiusServerTableFound = tableName == "RADIUS Servers" || isRadiusServerTableFound
				if tableName == "Users" && pluginDescriptorRefId == "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator" {
					tableRow := tableAttrs["rows"].(types.List).Elements()
					for tableRowIndex, row := range tableRow {
						rowAttrs := row.(types.Object).Attributes()
						fields := rowAttrs["fields"].(types.Set).Elements()
						sensitiveFields := rowAttrs["sensitive_fields"].(types.Set).Elements()
						if rowAttrs["fields"].IsUnknown() || rowAttrs["sensitive_fields"].IsUnknown() {
							continue
						}
						anyUnknownNames := false
						usernameFound := false
						passwordFound := false
						confirmPasswordFound := false
						for _, field := range fields {
							fieldRow := field.(types.Object).Attributes()
							if fieldRow["name"].IsUnknown() {
								anyUnknownNames = true
								break
							}
							nestedTableFieldName := fieldRow["name"].(types.String).ValueString()
							if nestedTableFieldName == "Username" {
								usernameFound = true
							}
							if nestedTableFieldName == "Password" {
								passwordFound = true
							}
							if nestedTableFieldName == "Confirm Password" {
								confirmPasswordFound = true
							}
						}
						for _, field := range sensitiveFields {
							fieldRow := field.(types.Object).Attributes()
							if fieldRow["name"].IsUnknown() {
								anyUnknownNames = true
								break
							}
							nestedTableFieldName := fieldRow["name"].(types.String).ValueString()
							if nestedTableFieldName == "Username" {
								usernameFound = true
							}
							if nestedTableFieldName == "Password" {
								passwordFound = true
							}
							if nestedTableFieldName == "Confirm Password" {
								confirmPasswordFound = true
							}
						}
						if !usernameFound && !anyUnknownNames {
							resp.Diagnostics.AddAttributeError(
								path.Root("configuration").AtMapKey("tables"),
								providererror.InvalidAttributeConfiguration,
								"The \"Username\" field is required in the Users table for the Simple Username Password Credential Validator.\n"+
									fmt.Sprintf("Missing from row index %d in Users table", tableRowIndex))
						}
						if !passwordFound && !anyUnknownNames {
							resp.Diagnostics.AddAttributeError(
								path.Root("configuration").AtMapKey("tables"),
								providererror.InvalidAttributeConfiguration,
								"The \"Password\" field is required in the Users table for the Simple Username Password Credential Validator.\n"+
									fmt.Sprintf("Missing from row index %d in Users table", tableRowIndex))
						}
						if !confirmPasswordFound && !anyUnknownNames {
							resp.Diagnostics.AddAttributeError(
								path.Root("configuration").AtMapKey("tables"),
								providererror.InvalidAttributeConfiguration,
								"The \"Confirm Password\" field is required in the Users table for the Simple Username Password Credential Validator.\n"+
									fmt.Sprintf("Missing from row index %d in Users table", tableRowIndex))
						}
					}
				}
			}
		}
	}

	if pluginDescriptorRefId == "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator" {
		if !isRadiusServerTableFound && tablesKnown {
			resp.Diagnostics.AddAttributeError(
				path.Root("configuration").AtMapKey("tables"),
				providererror.InvalidAttributeConfiguration,
				"At least one \"RADIUS Servers\" table is required for the RADIUS Username Password Credential Validator")
		}
	}

	fieldNameMap := map[string]bool{}
	if configuration["fields"] != nil {
		fields := configuration["fields"].(types.Set).Elements()
		sensitiveFields := configuration["sensitive_fields"].(types.Set).Elements()
		anyUnknowns := configuration["fields"].IsUnknown() || configuration["sensitive_fields"].IsUnknown()
		for _, field := range fields {
			field := field.(types.Object).Attributes()
			if field["name"].IsUnknown() {
				anyUnknowns = true
				break
			}
			fieldName := field["name"].(types.String).ValueString()
			fieldNameMap[fieldName] = true
		}
		for _, field := range sensitiveFields {
			field := field.(types.Object).Attributes()
			if field["name"].IsUnknown() {
				anyUnknowns = true
				break
			}
			fieldName := field["name"].(types.String).ValueString()
			fieldNameMap[fieldName] = true
		}

		if !anyUnknowns {
			switch pluginDescriptorRefId {
			case "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator":
				_, hasClientId := fieldNameMap["Client Id"]
				if !hasClientId {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"Client Id\" field is required for the PingOne for Enterprise Directory Password Credential Validator")
				}
				_, hasClientSecret := fieldNameMap["Client Secret"]
				if !hasClientSecret {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"Client Secret\" field is required for the PingOne for Enterprise Directory Password Credential Validator")
				}

			case "com.pingidentity.plugins.pcvs.p14c.PingOneForCustomersPCV":
				_, hasPingOneForCustomersDs := fieldNameMap["PingOne For Customers Datastore"]
				if !hasPingOneForCustomersDs {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"PingOne For Customers Datastore\" field is required for the PingOne Password Credential Validator")
				}

			case "com.pingidentity.plugins.pcvs.pingid.PingIdPCV":
				_, hasAuthenticationDuringErrors := fieldNameMap["Authentication During Errors"]
				if !hasAuthenticationDuringErrors {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"Authentication During Errors\" field is required for the PingID Password Credential Validator")
				}

			case "org.sourceid.saml20.domain.LDAPUsernamePasswordCredentialValidator":
				_, hasLdapDs := fieldNameMap["LDAP Datastore"]
				if !hasLdapDs {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"LDAP Datastore\" field is required for the LDAP Username Password Credential Validator")
				}
				_, hasSearchBase := fieldNameMap["Search Base"]
				if !hasSearchBase {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"Search Base\" field is required for the LDAP Username Password Credential Validator")
				}
				_, hasSearchFilter := fieldNameMap["Search Filter"]
				if !hasSearchFilter {
					resp.Diagnostics.AddAttributeError(
						path.Root("configuration").AtMapKey("fields"),
						providererror.InvalidAttributeConfiguration,
						"The \"Search Filter\" field is required for the LDAP Username Password Credential Validator")
				}
			}
		}
	}
}

func (r *passwordCredentialValidatorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *passwordCredentialValidatorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil || state == nil {
		return
	}

	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func addOptionalPasswordCredentialValidatorFields(addRequest *client.PasswordCredentialValidator, model passwordCredentialValidatorModel) {
	// attribute_contract
	if !model.AttributeContract.IsNull() && !model.AttributeContract.IsUnknown() {
		attributeContractValue := &client.PasswordCredentialValidatorAttributeContract{}
		attributeContractAttrs := model.AttributeContract.Attributes()
		if !attributeContractAttrs["core_attributes"].IsNull() && !attributeContractAttrs["core_attributes"].IsUnknown() {
			attributeContractValue.CoreAttributes = []client.PasswordCredentialValidatorAttribute{}
			for _, coreAttributesElement := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
				coreAttributesValue := client.PasswordCredentialValidatorAttribute{}
				coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
				coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
				attributeContractValue.CoreAttributes = append(attributeContractValue.CoreAttributes, coreAttributesValue)
			}
		}
		if !attributeContractAttrs["extended_attributes"].IsNull() && !attributeContractAttrs["extended_attributes"].IsUnknown() &&
			len(attributeContractAttrs["extended_attributes"].(types.Set).Elements()) > 0 {
			attributeContractValue.ExtendedAttributes = []client.PasswordCredentialValidatorAttribute{}
			for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.PasswordCredentialValidatorAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
		}
		addRequest.AttributeContract = attributeContractValue
	}

	// parent_ref
	if !model.ParentRef.IsNull() && !model.ParentRef.IsUnknown() {
		parentRefValue := &client.ResourceLink{}
		parentRefAttrs := model.ParentRef.Attributes()
		parentRefValue.Id = parentRefAttrs["id"].(types.String).ValueString()
		addRequest.ParentRef = parentRefValue
	}
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
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to build plugin descriptor ref request object: "+err.Error())
		return
	}

	// Configuration
	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	createPasswordCredentialValidators := client.NewPasswordCredentialValidator(plan.ValidatorId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	addOptionalPasswordCredentialValidatorFields(createPasswordCredentialValidators, plan)

	apiCreatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidator(config.AuthContext(ctx, r.providerConfig))
	apiCreatePasswordCredentialValidators = apiCreatePasswordCredentialValidators.Body(*createPasswordCredentialValidators)
	passwordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidatorExecute(apiCreatePasswordCredentialValidators)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating a Password Credential Validator", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state passwordCredentialValidatorModel

	diags = readPasswordCredentialValidatorResponse(ctx, passwordCredentialValidatorsResponse, &state, plan.Configuration, true, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *passwordCredentialValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state passwordCredentialValidatorModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadPasswordCredentialValidators, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(config.AuthContext(ctx, r.providerConfig), state.ValidatorId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Password Credential Validator", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting a Password Credential Validator", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	diags = readPasswordCredentialValidatorResponse(ctx, apiReadPasswordCredentialValidators, &state, state.Configuration, true, isImportRead)
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
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to build plugin descriptor ref request object: "+err.Error())
		return
	}

	// Configuration
	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	updatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidator(config.AuthContext(ctx, r.providerConfig), plan.ValidatorId.ValueString())
	createUpdateRequest := client.NewPasswordCredentialValidator(plan.ValidatorId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	addOptionalPasswordCredentialValidatorFields(createUpdateRequest, plan)

	updatePasswordCredentialValidators = updatePasswordCredentialValidators.Body(*createUpdateRequest)
	updatePasswordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidatorExecute(updatePasswordCredentialValidators)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating a Password Credential Validator", err, httpResp, &customId)
		return
	}

	// Read the response
	diags = readPasswordCredentialValidatorResponse(ctx, updatePasswordCredentialValidatorsResponse, &plan, plan.Configuration, true, false)
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
	httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.DeletePasswordCredentialValidator(config.AuthContext(ctx, r.providerConfig), state.ValidatorId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting a Password Credential Validator", err, httpResp, &customId)
	}
}

func (r *passwordCredentialValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("validator_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
