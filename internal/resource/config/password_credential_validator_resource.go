package config

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordCredentialValidatorsResource{}
	_ resource.ResourceWithConfigure   = &passwordCredentialValidatorsResource{}
	_ resource.ResourceWithImportState = &passwordCredentialValidatorsResource{}
)

// PasswordCredentialValidatorsResource is a helper function to simplify the provider implementation.
func PasswordCredentialValidatorsResource() resource.Resource {
	return &passwordCredentialValidatorsResource{}
}

// passwordCredentialValidatorsResource is the resource implementation.
type passwordCredentialValidatorsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type passwordCredentialValidatorsResourceModel struct {
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
	Id                  types.String `tfsdk:"id"`
	CustomId            types.String `tfsdk:"custom_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *passwordCredentialValidatorsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Password Credential Validators",
		Attributes: map[string]schema.Attribute{
			"custom_id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-zA-Z0-9_]{1,32}$"),
						"The plugin ID must be less than 33 characters, contain no spaces, and be alphanumeric.",
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes:  AddResourceLinkSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: AddResourceLinkSchema(),
			},
			"configuration": schema.SingleNestedAttribute{
				Description: "Plugin instance configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"tables": schema.ListNestedAttribute{
						Description: "List of configuration tables.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.ListNestedAttribute{
									Description: "List of table rows.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.ListNestedAttribute{
												Description: "The configuration fields in the row.",
												Computed:    true,
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Description: "The name of the configuration field.",
															Required:    true,
														},
														"value": schema.StringAttribute{
															Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
															Required:    true,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Optional:    true,
															PlanModifiers: []planmodifier.Bool{
																boolplanmodifier.UseStateForUnknown(),
															},
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Optional:    true,
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
										},
									},
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configuration fields.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Required:    true,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the password credential validator provides.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional attributes that can be returned by the password credential validator. The extended attributes are only used if the adapter supports them.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
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
					},
				},
			},
		},
	}

	AddCommonSchema(&schema)
	resp.Schema = schema
}

func (r *passwordCredentialValidatorsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model passwordCredentialValidatorsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if internaltypes.IsDefined(model.AttributeContract) {
		if len(model.AttributeContract.Attributes()["extended_attributes"].(types.List).Elements()) == 0 {
			resp.Diagnostics.AddError("Empty list!", "Please provide valid properties within extended_attributes. The list cannot be empty.\nIf no values are necessary, remove this property from your terraform file.")
		}
	}
}

func addOptionalPasswordCredentialValidatorsFields(ctx context.Context, addRequest *client.PasswordCredentialValidator, plan passwordCredentialValidatorsResourceModel) error {

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

// Metadata returns the resource type name.
func (r *passwordCredentialValidatorsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_credential_validator"
}

func (r *passwordCredentialValidatorsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readPasswordCredentialValidatorsResponse(ctx context.Context, r *client.PasswordCredentialValidator, state *passwordCredentialValidatorsResourceModel, configurationFromPlan basetypes.ObjectValue) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics

	state.Id = types.StringValue(r.Id)
	state.CustomId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)

	// state.pluginDescriptorRef
	pluginDescRef := r.GetPluginDescriptorRef()
	state.PluginDescriptorRef = internaltypes.ToStateResourceLink(ctx, &pluginDescRef, &respDiags)

	// state.parentRef
	parentRef := r.GetParentRef()
	state.ParentRef = internaltypes.ToStateResourceLink(ctx, &parentRef, &respDiags)

	// state.Configuration
	configurationAttrType := map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: FieldAttrTypes()}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: TableAttrTypes()}},
	}

	planFields := types.ListNull(types.ObjectType{AttrTypes: FieldAttrTypes()})
	planTables := types.ListNull(types.ObjectType{AttrTypes: TableAttrTypes()})

	planFieldsValue, ok := configurationFromPlan.Attributes()["fields"]
	if ok {
		planFields = planFieldsValue.(types.List)
	}
	planTablesValue, ok := configurationFromPlan.Attributes()["tables"]
	if ok {
		planTables = planTablesValue.(types.List)
	}

	fieldsAttrValue := ToFieldsListValue(r.Configuration.Fields, planFields, &diags)
	tablesAttrValue := ToTablesListValue(r.Configuration.Tables, planTables, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}
	state.Configuration, diags = types.ObjectValue(configurationAttrType, configurationAttrValue)
	respDiags.Append(diags...)

	// state.AttributeContract
	attrType := map[string]attr.Type{
		"name": basetypes.StringType{},
	}
	attributeContractTypes := map[string]attr.Type{
		"core_attributes":     basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"extended_attributes": basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"inherited":           basetypes.BoolType{},
	}

	if r.AttributeContract == nil {
		state.AttributeContract = types.ObjectNull(attributeContractTypes)
	} else {
		attrContract := r.AttributeContract

		// state.AttributeContract core_attributes
		attributeContractClientCoreAttributes := attrContract.CoreAttributes
		coreAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ca := range attributeContractClientCoreAttributes {
			coreAttribute := client.PasswordCredentialValidatorAttribute{}
			coreAttribute.Name = ca.Name
			coreAttrs = append(coreAttrs, coreAttribute)
		}
		attributeContractCoreAttributes, diags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, coreAttrs)
		respDiags.Append(diags...)

		// state.AttributeContract extended_attributes
		attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
		extdAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ea := range attributeContractClientExtendedAttributes {
			extendedAttr := client.PasswordCredentialValidatorAttribute{}
			extendedAttr.Name = ea.Name
			extdAttrs = append(extdAttrs, extendedAttr)
		}
		attributeContractExtendedAttributes, diags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, extdAttrs)
		respDiags.Append(diags...)

		attributeContractValues := map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributes,
			"extended_attributes": attributeContractExtendedAttributes,
			"inherited":           types.BoolPointerValue(attrContract.Inherited),
		}
		state.AttributeContract, diags = types.ObjectValue(attributeContractTypes, attributeContractValues)
		respDiags.Append(diags...)
	}

	return respDiags
}

func (r *passwordCredentialValidatorsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan passwordCredentialValidatorsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefId := plan.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	pluginDescRefResLink := client.NewResourceLinkWithDefaults()
	pluginDescRefResLink.Id = pluginDescRefId
	pluginDescRefErr := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), pluginDescRefResLink)
	if pluginDescRefErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", pluginDescRefErr.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfigurationWithDefaults()
	configErr := json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, true)), configuration)
	if configErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", configErr.Error())
		return
	}

	createPasswordCredentialValidators := client.NewPasswordCredentialValidator(plan.CustomId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalPasswordCredentialValidatorsFields(ctx, createPasswordCredentialValidators, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}
	_, requestErr := createPasswordCredentialValidators.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Password Credential Validator: %s", requestErr.Error())
	}
	apiCreatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsApi.CreatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePasswordCredentialValidators = apiCreatePasswordCredentialValidators.Body(*createPasswordCredentialValidators)
	passwordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.CreatePasswordCredentialValidatorExecute(apiCreatePasswordCredentialValidators)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating a Password Credential Validator", err, httpResp)
		return
	}
	_, responseErr := passwordCredentialValidatorsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Password Credential Validator: %s", responseErr.Error())
	}

	// Read the response into the state
	var state passwordCredentialValidatorsResourceModel

	diags = readPasswordCredentialValidatorsResponse(ctx, passwordCredentialValidatorsResponse, &state, plan.Configuration)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *passwordCredentialValidatorsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state passwordCredentialValidatorsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadPasswordCredentialValidators, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.GetPasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Password Credential Validator", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Password Credential Validator", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadPasswordCredentialValidators.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Password Credential Validator: %s", responseErr.Error())
	}

	// Read the response into the state
	diags = readPasswordCredentialValidatorsResponse(ctx, apiReadPasswordCredentialValidators, &state, state.Configuration)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *passwordCredentialValidatorsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan passwordCredentialValidatorsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefId := plan.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	pluginDescRefResLink := client.NewResourceLinkWithDefaults()
	pluginDescRefResLink.Id = pluginDescRefId
	pluginDescRefErr := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), pluginDescRefResLink)
	if pluginDescRefErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", pluginDescRefErr.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfiguration()
	configErr := json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, true)), configuration)
	if configErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", configErr.Error())
		return
	}

	updatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsApi.UpdatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), plan.CustomId.ValueString())
	createUpdateRequest := client.NewPasswordCredentialValidator(plan.CustomId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalPasswordCredentialValidatorsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Password Credential Validator: %s", requestErr.Error())
	}
	updatePasswordCredentialValidators = updatePasswordCredentialValidators.Body(*createUpdateRequest)
	updatePasswordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.UpdatePasswordCredentialValidatorExecute(updatePasswordCredentialValidators)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating a Password Credential Validator", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updatePasswordCredentialValidatorsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Password Credential Validator: %s", responseErr.Error())
	}
	// Read the response
	diags = readPasswordCredentialValidatorsResponse(ctx, updatePasswordCredentialValidatorsResponse, &plan, plan.Configuration)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *passwordCredentialValidatorsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state passwordCredentialValidatorsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.DeletePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Password Credential Validator", err, httpResp)
	}
}

func (r *passwordCredentialValidatorsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
