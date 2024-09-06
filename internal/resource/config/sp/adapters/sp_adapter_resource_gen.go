// Code generated by ping-terraform-plugin-framework-generator

package spadapters

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &spAdapterResource{}
	_ resource.ResourceWithConfigure   = &spAdapterResource{}
	_ resource.ResourceWithImportState = &spAdapterResource{}
)

func SpAdapterResource() resource.Resource {
	return &spAdapterResource{}
}

type spAdapterResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *spAdapterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_adapter"
}

func (r *spAdapterResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type spAdapterResourceModel struct {
	AdapterId             types.String `tfsdk:"adapter_id"`
	AttributeContract     types.Object `tfsdk:"attribute_contract"`
	Configuration         types.Object `tfsdk:"configuration"`
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ParentRef             types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef   types.Object `tfsdk:"plugin_descriptor_ref"`
	TargetApplicationInfo types.Object `tfsdk:"target_application_info"`
}

func (r *spAdapterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage SP adapters.",
		Attributes: map[string]schema.Attribute{
			"adapter_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Optional:    false,
						Computed:    true,
						Default:     listdefault.StaticValue(coreAttributesDefault),
						Description: "A list of read-only attributes that are automatically populated by the SP adapter descriptor.",
					},
					"extended_attributes": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Default:     listdefault.StaticValue(extendedAttributesDefault),
						Description: "A list of additional attributes that can be returned by the SP adapter. The extended attributes are only used if the adapter supports them.",
					},
				},
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(attributeContractDefault),
				Description: "A set of attributes exposed by an SP adapter.",
			},
			"configuration": pluginconfiguration.ToSchema(),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
						Description: "The ID of the resource.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
			},
			"target_application_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"application_icon_url": schema.StringAttribute{
						Optional:    true,
						Description: "The application icon URL.",
						Validators: []validator.String{
							configvalidators.ValidUrl(),
							stringvalidator.LengthAtLeast(1),
						},
					},
					"application_name": schema.StringAttribute{
						Optional:    true,
						Description: "The application name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(targetApplicationInfoDefault),
				Description: "Target Application Information exposed by an SP adapter.",
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *spAdapterResourceModel) buildClientStruct() (*client.SpAdapter, error) {
	result := &client.SpAdapter{}
	var err error
	// adapter_id
	result.Id = model.AdapterId.ValueString()
	// attribute_contract
	if !model.AttributeContract.IsNull() {
		attributeContractValue := &client.SpAdapterAttributeContract{}
		attributeContractAttrs := model.AttributeContract.Attributes()
		attributeContractValue.ExtendedAttributes = []client.SpAdapterAttribute{}
		for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.List).Elements() {
			extendedAttributesValue := client.SpAdapterAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		result.AttributeContract = attributeContractValue
	}

	// configuration
	resultConfig, err := pluginconfiguration.ClientStruct(model.Configuration)
	if err != nil {
		return nil, err
	}
	result.Configuration = *resultConfig

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

	// target_application_info
	if !model.TargetApplicationInfo.IsNull() {
		targetApplicationInfoValue := &client.SpAdapterTargetApplicationInfo{}
		targetApplicationInfoAttrs := model.TargetApplicationInfo.Attributes()
		targetApplicationInfoValue.ApplicationIconUrl = targetApplicationInfoAttrs["application_icon_url"].(types.String).ValueStringPointer()
		targetApplicationInfoValue.ApplicationName = targetApplicationInfoAttrs["application_name"].(types.String).ValueStringPointer()
		result.TargetApplicationInfo = targetApplicationInfoValue
	}

	return result, nil
}

func (state *spAdapterResourceModel) readClientResponse(response *client.SpAdapter, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
	// adapter_id
	state.AdapterId = types.StringValue(response.Id)
	// attribute_contract
	attributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: attributeContractCoreAttributesAttrTypes}
	attributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: attributeContractExtendedAttributesAttrTypes}
	attributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.ListType{ElemType: attributeContractCoreAttributesElementType},
		"extended_attributes": types.ListType{ElemType: attributeContractExtendedAttributesElementType},
	}
	var attributeContractValue types.Object
	if response.AttributeContract == nil {
		attributeContractValue = types.ObjectNull(attributeContractAttrTypes)
	} else {
		var attributeContractCoreAttributesValues []attr.Value
		for _, attributeContractCoreAttributesResponseValue := range response.AttributeContract.CoreAttributes {
			attributeContractCoreAttributesValue, diags := types.ObjectValue(attributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(attributeContractCoreAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractCoreAttributesValues = append(attributeContractCoreAttributesValues, attributeContractCoreAttributesValue)
		}
		attributeContractCoreAttributesValue, diags := types.ListValue(attributeContractCoreAttributesElementType, attributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var attributeContractExtendedAttributesValues []attr.Value
		for _, attributeContractExtendedAttributesResponseValue := range response.AttributeContract.ExtendedAttributes {
			attributeContractExtendedAttributesValue, diags := types.ObjectValue(attributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name": types.StringValue(attributeContractExtendedAttributesResponseValue.Name),
			})
			respDiags.Append(diags...)
			attributeContractExtendedAttributesValues = append(attributeContractExtendedAttributesValues, attributeContractExtendedAttributesValue)
		}
		attributeContractExtendedAttributesValue, diags := types.ListValue(attributeContractExtendedAttributesElementType, attributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		attributeContractValue, diags = types.ObjectValue(attributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributesValue,
			"extended_attributes": attributeContractExtendedAttributesValue,
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
	// target_application_info
	targetApplicationInfoAttrTypes := map[string]attr.Type{
		"application_icon_url": types.StringType,
		"application_name":     types.StringType,
	}
	var targetApplicationInfoValue types.Object
	if response.TargetApplicationInfo == nil {
		targetApplicationInfoValue = types.ObjectNull(targetApplicationInfoAttrTypes)
	} else {
		targetApplicationInfoValue, diags = types.ObjectValue(targetApplicationInfoAttrTypes, map[string]attr.Value{
			"application_icon_url": types.StringPointerValue(response.TargetApplicationInfo.ApplicationIconUrl),
			"application_name":     types.StringPointerValue(response.TargetApplicationInfo.ApplicationName),
		})
		respDiags.Append(diags...)
	}

	state.TargetApplicationInfo = targetApplicationInfoValue
	return respDiags
}

func (r *spAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data spAdapterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, err := data.buildClientStruct()
	if err != nil {
		resp.Diagnostics.AddError("Failed to build client struct for the spAdapter", err.Error())
		return
	}
	apiCreateRequest := r.apiClient.SpAdaptersAPI.CreateSpAdapter(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpAdaptersAPI.CreateSpAdapterExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the spAdapter", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var data spAdapterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.SpAdaptersAPI.GetSpAdapter(config.AuthContext(ctx, r.providerConfig), data.AdapterId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "SP Adapter", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the spAdapter", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, isImportRead)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spAdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data spAdapterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, err := data.buildClientStruct()
	if err != nil {
		resp.Diagnostics.AddError("Failed to build client struct for the spAdapter", err.Error())
		return
	}
	apiUpdateRequest := r.apiClient.SpAdaptersAPI.UpdateSpAdapter(config.AuthContext(ctx, r.providerConfig), data.AdapterId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpAdaptersAPI.UpdateSpAdapterExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the spAdapter", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spAdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data spAdapterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.SpAdaptersAPI.DeleteSpAdapter(config.AuthContext(ctx, r.providerConfig), data.AdapterId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the spAdapter", err, httpResp)
	}
}

func (r *spAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to adapter_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("adapter_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}