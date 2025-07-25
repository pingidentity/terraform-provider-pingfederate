// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package idptokenprocessors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &idpTokenProcessorResource{}
	_ resource.ResourceWithConfigure   = &idpTokenProcessorResource{}
	_ resource.ResourceWithImportState = &idpTokenProcessorResource{}

	customId = "processor_id"
)

func IdpTokenProcessorResource() resource.Resource {
	return &idpTokenProcessorResource{}
}

type idpTokenProcessorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *idpTokenProcessorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_token_processor"
}

func (r *idpTokenProcessorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type idpTokenProcessorResourceModel struct {
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
	Configuration       types.Object `tfsdk:"configuration"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ProcessorId         types.String `tfsdk:"processor_id"`
}

func (r *idpTokenProcessorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage token processor instances.",
		Attributes: map[string]schema.Attribute{
			"attribute_contract": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"masked": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
								},
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Required: true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						Description: "A list of token processor attributes that correspond to the attributes exposed by the token processor type.",
					},
					"extended_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"masked": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
								},
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(extendedAttributesDefault),
						Description: "A list of additional attributes that can be returned by the token processor. The extended attributes are only used if the token processor supports them.",
					},
					"mask_ognl_values": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to `false`.",
					},
				},
				Optional:    true,
				Description: "A set of attributes exposed by a token processor.",
			},
			"configuration": pluginconfiguration.ToSchema(),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Optional:    true,
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource. This field is immutable and will trigger a replacement plan if changed.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. This field is immutable and will trigger a replacement plan if changed. Note: Ignored when specifying a connection's adapter override.",
			},
			"processor_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (r *idpTokenProcessorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *idpTokenProcessorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	var state *idpTokenProcessorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}
	var respDiags diag.Diagnostics
	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (model *idpTokenProcessorResourceModel) buildClientStruct() (*client.TokenProcessor, diag.Diagnostics) {
	result := &client.TokenProcessor{}
	var respDiags diag.Diagnostics
	// attribute_contract
	if !model.AttributeContract.IsNull() {
		attributeContractValue := &client.TokenProcessorAttributeContract{}
		attributeContractAttrs := model.AttributeContract.Attributes()
		attributeContractValue.CoreAttributes = []client.TokenProcessorAttribute{}
		for _, coreAttributesElement := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.TokenProcessorAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			attributeContractValue.CoreAttributes = append(attributeContractValue.CoreAttributes, coreAttributesValue)
		}
		attributeContractValue.ExtendedAttributes = []client.TokenProcessorAttribute{}
		for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.TokenProcessorAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		attributeContractValue.MaskOgnlValues = attributeContractAttrs["mask_ognl_values"].(types.Bool).ValueBoolPointer()
		result.AttributeContract = attributeContractValue
	}

	// configuration
	configurationValue := pluginconfiguration.ClientStruct(model.Configuration)
	result.Configuration = *configurationValue

	// name
	result.Name = model.Name.ValueString()
	// parent_ref
	if !model.ParentRef.IsNull() {
		parentRefValue := &client.ResourceLink{}
		parentRefAttrs := model.ParentRef.Attributes()
		parentRefValue.Id = parentRefAttrs["id"].(types.String).ValueString()
		result.ParentRef = parentRefValue
	}

	// plugin_descriptor_ref
	pluginDescriptorRefValue := client.ResourceLink{}
	pluginDescriptorRefAttrs := model.PluginDescriptorRef.Attributes()
	pluginDescriptorRefValue.Id = pluginDescriptorRefAttrs["id"].(types.String).ValueString()
	result.PluginDescriptorRef = pluginDescriptorRefValue

	// processor_id
	result.Id = model.ProcessorId.ValueString()
	return result, respDiags
}

func (state *idpTokenProcessorResourceModel) readClientResponse(response *client.TokenProcessor, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
	// attribute_contract
	attributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	attributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: attributeContractCoreAttributesAttrTypes}
	attributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"masked": types.BoolType,
		"name":   types.StringType,
	}
	attributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: attributeContractExtendedAttributesAttrTypes}
	attributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: attributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: attributeContractExtendedAttributesElementType},
		"mask_ognl_values":    types.BoolType,
	}
	var attributeContractValue types.Object
	if response.AttributeContract == nil {
		attributeContractValue = types.ObjectNull(attributeContractAttrTypes)
	} else {
		var attributeContractCoreAttributesValues []attr.Value
		for _, attributeContractCoreAttributesResponseValue := range response.AttributeContract.CoreAttributes {
			attributeContractCoreAttributesValue, diags := types.ObjectValue(attributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(attributeContractCoreAttributesResponseValue.Masked),
				"name":   types.StringValue(attributeContractCoreAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractCoreAttributesValues = append(attributeContractCoreAttributesValues, attributeContractCoreAttributesValue)
		}
		attributeContractCoreAttributesValue, diags := types.SetValue(attributeContractCoreAttributesElementType, attributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var attributeContractExtendedAttributesValues []attr.Value
		for _, attributeContractExtendedAttributesResponseValue := range response.AttributeContract.ExtendedAttributes {
			attributeContractExtendedAttributesValue, diags := types.ObjectValue(attributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"masked": types.BoolPointerValue(attributeContractExtendedAttributesResponseValue.Masked),
				"name":   types.StringValue(attributeContractExtendedAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractExtendedAttributesValues = append(attributeContractExtendedAttributesValues, attributeContractExtendedAttributesValue)
		}
		attributeContractExtendedAttributesValue, diags := types.SetValue(attributeContractExtendedAttributesElementType, attributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		attributeContractValue, diags = types.ObjectValue(attributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributesValue,
			"extended_attributes": attributeContractExtendedAttributesValue,
			"mask_ognl_values":    types.BoolPointerValue(response.AttributeContract.MaskOgnlValues),
		})
		respDiags.Append(diags...)
	}

	state.AttributeContract = attributeContractValue
	// configuration
	configurationValue, diags := pluginconfiguration.ToState(state.Configuration, &response.Configuration, isImportRead)
	respDiags.Append(diags...)

	state.Configuration = configurationValue
	// name
	state.Name = types.StringValue(response.Name)
	// parent_ref
	parentRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var parentRefValue types.Object
	if response.ParentRef == nil {
		parentRefValue = types.ObjectNull(parentRefAttrTypes)
	} else {
		parentRefValue, diags = types.ObjectValue(parentRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.ParentRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.ParentRef = parentRefValue
	// plugin_descriptor_ref
	pluginDescriptorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	pluginDescriptorRefValue, diags := types.ObjectValue(pluginDescriptorRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.PluginDescriptorRef.Id),
	})
	respDiags.Append(diags...)

	state.PluginDescriptorRef = pluginDescriptorRefValue
	// processor_id
	state.ProcessorId = types.StringValue(response.Id)
	return respDiags
}

func (r *idpTokenProcessorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data idpTokenProcessorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.IdpTokenProcessorsAPI.CreateTokenProcessor(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.IdpTokenProcessorsAPI.CreateTokenProcessorExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the idpTokenProcessor", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpTokenProcessorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var data idpTokenProcessorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.IdpTokenProcessorsAPI.GetTokenProcessor(config.AuthContext(ctx, r.providerConfig), data.ProcessorId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "IdP Token Processor", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the idpTokenProcessor", err, httpResp, &customId)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, isImportRead)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpTokenProcessorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data idpTokenProcessorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.IdpTokenProcessorsAPI.UpdateTokenProcessor(config.AuthContext(ctx, r.providerConfig), data.ProcessorId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.IdpTokenProcessorsAPI.UpdateTokenProcessorExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the idpTokenProcessor", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpTokenProcessorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data idpTokenProcessorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.IdpTokenProcessorsAPI.DeleteTokenProcessor(config.AuthContext(ctx, r.providerConfig), data.ProcessorId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the idpTokenProcessor", err, httpResp, &customId)
	}
}

func (r *idpTokenProcessorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to processor_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("processor_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
