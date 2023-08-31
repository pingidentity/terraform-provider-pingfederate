package config

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *passwordCredentialValidatorsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Password Credential Validators",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
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
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"configuration": schema.SingleNestedAttribute{
				Description: "Plugin instance configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"tables": schema.SetNestedAttribute{
						Description: "List of configuration tables.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.SetNestedAttribute{
									Description: "List of table rows.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Set{
										setplanmodifier.UseStateForUnknown(),
										setplanmodifier.RequiresReplace(),
									},
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.SetNestedAttribute{
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
															Computed:    true,
															Optional:    true,
															PlanModifiers: []planmodifier.String{
																stringplanmodifier.UseStateForUnknown(),
																stringplanmodifier.RequiresReplace(),
															},
														},
														"encrypted_value": schema.StringAttribute{
															Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
															Computed:    true,
															Optional:    false,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Computed:    true,
															Optional:    true,
															Default:     booldefault.StaticBool(false),
															PlanModifiers: []planmodifier.Bool{
																boolplanmodifier.UseStateForUnknown(),
															},
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Computed:    true,
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
									Default:     booldefault.StaticBool(false),
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"fields": schema.SetNestedAttribute{
						Description: "List of configuration fields.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
										stringplanmodifier.RequiresReplace(),
									},
								},
								"encrypted_value": schema.StringAttribute{
									Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
									Computed:    true,
									Optional:    false,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Default:     booldefault.StaticBool(false),
									Computed:    true,
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
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of read-only attributes that are automatically populated by the password credential validator descriptor.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
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
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the password credential validator. The extended attributes are only used if the adapter supports them.",
						Computed:    true,
						Optional:    true,
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
}

func (r *passwordCredentialValidatorsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model passwordCredentialValidatorsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if internaltypes.IsDefined(model.AttributeContract) {
		if len(model.AttributeContract.Attributes()["extended_attributes"].(types.Set).Elements()) == 0 {
			resp.Diagnostics.AddError("Empty set!", "Please provide valid properties within extended_attributes. The set cannot be empty.\nIf no values are necessary, remove this property from your terraform file.")
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
		extendedAttrsLength := len(plan.AttributeContract.Attributes()["extended_attributes"].(types.Set).Elements())
		if extendedAttrsLength == 0 {
			addRequest.AttributeContract.ExtendedAttributes = nil
		}
	}
	return nil
}

// Metadata returns the resource type name.
func (r *passwordCredentialValidatorsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_credential_validators"
}

func (r *passwordCredentialValidatorsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readPasswordCredentialValidatorsResponse(ctx context.Context, r *client.PasswordCredentialValidator, state *passwordCredentialValidatorsResourceModel, configurationFromPlan basetypes.ObjectValue) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)

	// state.pluginDescriptorRef
	resLinkAttrTypes := map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}
	pluginDescRef := r.GetPluginDescriptorRef()
	state.PluginDescriptorRef, _ = types.ObjectValueFrom(ctx, resLinkAttrTypes, pluginDescRef)

	// state.parentRef
	parentRef := r.GetParentRef()
	state.ParentRef, _ = types.ObjectValueFrom(ctx, resLinkAttrTypes, parentRef)

	// state.Configuration
	fieldAttrType := map[string]attr.Type{
		"name":            basetypes.StringType{},
		"value":           basetypes.StringType{},
		"encrypted_value": basetypes.StringType{},
		"inherited":       basetypes.BoolType{},
	}

	rowAttrType := map[string]attr.Type{
		"fields":      basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: fieldAttrType}},
		"default_row": basetypes.BoolType{},
	}

	// configuration object
	tableAttrType := map[string]attr.Type{
		"name":      basetypes.StringType{},
		"rows":      basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: rowAttrType}},
		"inherited": basetypes.BoolType{},
	}

	configFromPlanAttrs := configurationFromPlan.Attributes()
	tables := []client.ConfigTable{}
	tablesElems := configFromPlanAttrs["tables"]
	if tablesElems != nil && len(tablesElems.(types.Set).Elements()) != 0 {
		for _, ct := range tablesElems.(types.Set).Elements() {
			tableValue := client.ConfigTable{}
			configTable := ct.(types.Object).Attributes()
			tableValue.Inherited = configTable["inherited"].(types.Bool).ValueBoolPointer()
			tableValue.Name = configTable["name"].(types.String).ValueString()
			tableValue.Inherited = configTable["inherited"].(types.Bool).ValueBoolPointer()
			tableRows := configTable["rows"].(types.Set).Elements()
			toStateTableRows := []client.ConfigRow{}
			for _, tr := range tableRows {
				tableRow := client.ConfigRow{}
				tableRowObjAttrs := tr.(types.Object).Attributes()
				tableRow.DefaultRow = tableRowObjAttrs["default_row"].(types.Bool).ValueBoolPointer()
				tableRowFields := tableRowObjAttrs["fields"].(types.Set).Elements()
				toStateTableRowFields := []client.ConfigField{}
				for _, trf := range tableRowFields {
					tableRowField := client.ConfigField{}
					tableRowFieldObjAttrs := trf.(types.Object).Attributes()
					tableRowField.Name = tableRowFieldObjAttrs["name"].(types.String).ValueString()
					tableRowField.Inherited = tableRowFieldObjAttrs["inherited"].(types.Bool).ValueBoolPointer()
					tableRowField.EncryptedValue = tableRowFieldObjAttrs["encrypted_value"].(types.String).ValueStringPointer()
					tableRowField.Value = tableRowFieldObjAttrs["value"].(types.String).ValueStringPointer()
					toStateTableRowFields = append(toStateTableRowFields, tableRowField)
				}
				tableRow.Fields = toStateTableRowFields
				toStateTableRows = append(toStateTableRows, tableRow)
			}
			tableValue.Rows = toStateTableRows
			tables = append(tables, tableValue)
		}
	}
	tableValue, _ := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: tableAttrType}, tables)

	fields := []client.ConfigField{}
	fieldsElems := configFromPlanAttrs["fields"]
	if fieldsElems != nil && len(fieldsElems.(types.Set).Elements()) != 0 {
		for _, cf := range fieldsElems.(types.Set).Elements() {
			fieldValue := client.ConfigField{}
			configFieldAttrs := cf.(types.Object).Attributes()
			fieldValue.Name = configFieldAttrs["name"].(types.String).ValueString()
			fieldValue.EncryptedValue = configFieldAttrs["encrypted_value"].(types.String).ValueStringPointer()
			fieldValue.Value = configFieldAttrs["value"].(types.String).ValueStringPointer()
			fieldValue.Inherited = configFieldAttrs["inherited"].(types.Bool).ValueBoolPointer()
			fields = append(fields, fieldValue)
		}
	}
	configFieldValue, _ := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: fieldAttrType}, fields)

	configurationAttrType := map[string]attr.Type{
		"fields": basetypes.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrType}},
		"tables": basetypes.SetType{ElemType: types.ObjectType{AttrTypes: tableAttrType}},
	}

	configurationAttrValue := map[string]attr.Value{
		"fields": configFieldValue,
		"tables": tableValue,
	}
	state.Configuration, _ = types.ObjectValue(configurationAttrType, configurationAttrValue)

	// state.AttributeContract
	attrContract := r.GetAttributeContract()

	attrType := map[string]attr.Type{
		"name": basetypes.StringType{},
	}

	// state.AttributeContract core_attributes
	attributeContractClientCoreAttributes := attrContract.CoreAttributes
	coreAttrs := []client.PasswordCredentialValidatorAttribute{}
	for _, ca := range attributeContractClientCoreAttributes {
		coreAttribute := client.PasswordCredentialValidatorAttribute{}
		coreAttribute.Name = ca.Name
		coreAttrs = append(coreAttrs, coreAttribute)
	}
	attributeContractCoreAttributes, _ := types.SetValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, coreAttrs)

	// state.AttributeContract extended_attributes
	attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
	extdAttrs := []client.PasswordCredentialValidatorAttribute{}
	for _, ea := range attributeContractClientExtendedAttributes {
		extendedAttr := client.PasswordCredentialValidatorAttribute{}
		extendedAttr.Name = ea.Name
		extdAttrs = append(extdAttrs, extendedAttr)
	}
	attributeContractExtendedAttributes, _ := types.SetValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, extdAttrs)

	attributeContractTypes := map[string]attr.Type{
		"core_attributes":     basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"extended_attributes": basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"inherited":           basetypes.BoolType{},
	}

	attributeContractValues := map[string]attr.Value{
		"core_attributes":     attributeContractCoreAttributes,
		"extended_attributes": attributeContractExtendedAttributes,
		"inherited":           types.BoolPointerValue(attrContract.Inherited),
	}
	state.AttributeContract, _ = types.ObjectValue(attributeContractTypes, attributeContractValues)
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

	createPasswordCredentialValidators := client.NewPasswordCredentialValidator(plan.Id.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalPasswordCredentialValidatorsFields(ctx, createPasswordCredentialValidators, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for PasswordCredentialValidators", err.Error())
		return
	}
	requestJson, err := createPasswordCredentialValidators.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsApi.CreatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePasswordCredentialValidators = apiCreatePasswordCredentialValidators.Body(*createPasswordCredentialValidators)
	passwordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.CreatePasswordCredentialValidatorExecute(apiCreatePasswordCredentialValidators)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the PasswordCredentialValidators", err, httpResp)
		return
	}
	responseJson, err := passwordCredentialValidatorsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordCredentialValidatorsResourceModel

	readPasswordCredentialValidatorsResponse(ctx, passwordCredentialValidatorsResponse, &state, plan.Configuration)
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
	apiReadPasswordCredentialValidators, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.GetPasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Password Credential Validator", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Credential Validator", err, httpResp)
		}
		return
	}
	// Log response JSON
	responseJson, err := apiReadPasswordCredentialValidators.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPasswordCredentialValidatorsResponse(ctx, apiReadPasswordCredentialValidators, &state, state.Configuration)

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

	// Get the current state to see how any attributes are changing
	var state passwordCredentialValidatorsResourceModel
	req.State.Get(ctx, &state)

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

	updatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsApi.UpdatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewPasswordCredentialValidator(plan.Id.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalPasswordCredentialValidatorsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for PasswordCredentialValidators", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updatePasswordCredentialValidators = updatePasswordCredentialValidators.Body(*createUpdateRequest)
	updatePasswordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.UpdatePasswordCredentialValidatorExecute(updatePasswordCredentialValidators)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating PasswordCredentialValidators", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updatePasswordCredentialValidatorsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readPasswordCredentialValidatorsResponse(ctx, updatePasswordCredentialValidatorsResponse, &state, plan.Configuration)

	// Update computed values
	diags = resp.State.Set(ctx, state)
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
	httpResp, err := r.apiClient.PasswordCredentialValidatorsApi.DeletePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a PasswordCredentialValidator", err, httpResp)
		return
	}
}

func (r *passwordCredentialValidatorsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
