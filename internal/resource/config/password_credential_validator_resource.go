package config

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithConfigure   = &passwordCredentialValidatorResource{}
	_ resource.ResourceWithImportState = &passwordCredentialValidatorResource{}
)

var (
	attrType = map[string]attr.Type{
		"name": basetypes.StringType{},
	}

	attributeContractTypes = map[string]attr.Type{
		"core_attributes":     basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"extended_attributes": basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"inherited":           basetypes.BoolType{},
	}

	emptyAttrList, _ = types.ListValue(types.ObjectType{AttrTypes: attrType}, nil)
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

type passwordCredentialValidatorResourceModel struct {
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
	Id                  types.String `tfsdk:"id"`
	CustomId            types.String `tfsdk:"custom_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *passwordCredentialValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Password Credential Validator",
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
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, true,
		"The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.")
	resp.Schema = schema
}

func addOptionalPasswordCredentialValidatorFields(ctx context.Context, addRequest *client.PasswordCredentialValidator, plan passwordCredentialValidatorResourceModel) error {
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

func readPasswordCredentialValidatorResponse(ctx context.Context, r *client.PasswordCredentialValidator, state *passwordCredentialValidatorResourceModel, configurationFromPlan basetypes.ObjectValue) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(r.Id)
	state.CustomId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, respDiags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	diags.Append(respDiags...)
	state.ParentRef, respDiags = resourcelink.ToState(ctx, r.ParentRef)
	diags.Append(respDiags...)
	state.Configuration, respDiags = pluginconfiguration.ToState(configurationFromPlan, &r.Configuration)
	diags.Append(respDiags...)

	// state.AttributeContract
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
		attributeContractCoreAttributes, respDiags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, coreAttrs)
		diags.Append(respDiags...)

		// state.AttributeContract extended_attributes
		attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
		extdAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ea := range attributeContractClientExtendedAttributes {
			extendedAttr := client.PasswordCredentialValidatorAttribute{}
			extendedAttr.Name = ea.Name
			extdAttrs = append(extdAttrs, extendedAttr)
		}
		attributeContractExtendedAttributes, respDiags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, extdAttrs)
		diags.Append(respDiags...)

		attributeContractValues := map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributes,
			"extended_attributes": attributeContractExtendedAttributes,
			"inherited":           types.BoolPointerValue(attrContract.Inherited),
		}
		state.AttributeContract, respDiags = types.ObjectValue(attributeContractTypes, attributeContractValues)
		diags.Append(respDiags...)
	}

	return diags
}

func (r *passwordCredentialValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan passwordCredentialValidatorResourceModel

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

	createPasswordCredentialValidators := client.NewPasswordCredentialValidator(plan.CustomId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err = addOptionalPasswordCredentialValidatorFields(ctx, createPasswordCredentialValidators, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}
	_, requestErr := createPasswordCredentialValidators.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Password Credential Validator: %s", requestErr.Error())
	}
	apiCreatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePasswordCredentialValidators = apiCreatePasswordCredentialValidators.Body(*createPasswordCredentialValidators)
	passwordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.CreatePasswordCredentialValidatorExecute(apiCreatePasswordCredentialValidators)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating a Password Credential Validator", err, httpResp)
		return
	}
	_, responseErr := passwordCredentialValidatorsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Password Credential Validator: %s", responseErr.Error())
	}

	// Read the response into the state
	var state passwordCredentialValidatorResourceModel

	diags = readPasswordCredentialValidatorResponse(ctx, passwordCredentialValidatorsResponse, &state, plan.Configuration)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *passwordCredentialValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state passwordCredentialValidatorResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadPasswordCredentialValidators, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
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
	diags = readPasswordCredentialValidatorResponse(ctx, apiReadPasswordCredentialValidators, &state, state.Configuration)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *passwordCredentialValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan passwordCredentialValidatorResourceModel
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

	updatePasswordCredentialValidators := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), plan.CustomId.ValueString())
	createUpdateRequest := client.NewPasswordCredentialValidator(plan.CustomId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err = addOptionalPasswordCredentialValidatorFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Password Credential Validator", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a Password Credential Validator: %s", requestErr.Error())
	}
	updatePasswordCredentialValidators = updatePasswordCredentialValidators.Body(*createUpdateRequest)
	updatePasswordCredentialValidatorsResponse, httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidatorExecute(updatePasswordCredentialValidators)
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
	diags = readPasswordCredentialValidatorResponse(ctx, updatePasswordCredentialValidatorsResponse, &plan, plan.Configuration)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *passwordCredentialValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state passwordCredentialValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.PasswordCredentialValidatorsAPI.DeletePasswordCredentialValidator(ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Password Credential Validator", err, httpResp)
	}
}

func (r *passwordCredentialValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
