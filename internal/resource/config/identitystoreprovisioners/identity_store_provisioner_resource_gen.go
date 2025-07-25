// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package identitystoreprovisioners

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &identityStoreProvisionerResource{}
	_ resource.ResourceWithConfigure   = &identityStoreProvisionerResource{}
	_ resource.ResourceWithImportState = &identityStoreProvisionerResource{}

	customId = "provisioner_id"
)

func IdentityStoreProvisionerResource() resource.Resource {
	return &identityStoreProvisionerResource{}
}

type identityStoreProvisionerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *identityStoreProvisionerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_store_provisioner"
}

func (r *identityStoreProvisionerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type identityStoreProvisionerResourceModel struct {
	AttributeContract      types.Object `tfsdk:"attribute_contract"`
	Configuration          types.Object `tfsdk:"configuration"`
	GroupAttributeContract types.Object `tfsdk:"group_attribute_contract"`
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ParentRef              types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef    types.Object `tfsdk:"plugin_descriptor_ref"`
	ProvisionerId          types.String `tfsdk:"provisioner_id"`
}

func (r *identityStoreProvisionerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage identity store provisioner instances.",
		Attributes: map[string]schema.Attribute{
			"attribute_contract": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Required:    true,
						Description: "A list of identity store provisioner attributes that correspond to the attributes exposed by the identity store provisioner type.",
					},
					"core_attributes_all": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Computed:    true,
						Description: "A list of identity store provisioner attributes that correspond to the attributes exposed by the identity store provisioner type, including attributes computed by PingFederate.",
					},
					"extended_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Optional:    true,
						Description: "A list of additional attributes that can be returned by the identity store provisioner. The extended attributes are only used if the provisioner supports them.",
					},
				},
				Required:    true,
				Description: "A set of attributes exposed by an identity store provisioner.",
			},
			"configuration": pluginconfiguration.ToSchema(),
			"group_attribute_contract": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Required:    true,
						Description: "A list of identity store provisioner group attributes that correspond to the group attributes exposed by the identity store provisioner type.",
					},
					"core_attributes_all": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Computed:    true,
						Description: "A list of identity store provisioner group attributes that correspond to the group attributes exposed by the identity store provisioner type, including attributes computed by PingFederate.",
					},
					"extended_attributes": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of this attribute.",
								},
							},
						},
						Optional:    true,
						Description: "A list of additional group attributes that can be returned by the identity store provisioner. The extended group attributes are only used if the provisioner supports them.",
					},
				},
				Required:    true,
				Description: "A set of group attributes exposed by an identity store provisioner.",
			},
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
						Description: "The ID of the resource.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
			},
			"provisioner_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateIdWithCharLimit(),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (r *identityStoreProvisionerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *identityStoreProvisionerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	var state *identityStoreProvisionerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}
	var respDiags diag.Diagnostics
	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (model *identityStoreProvisionerResourceModel) buildClientStruct() (*client.IdentityStoreProvisioner, diag.Diagnostics) {
	result := &client.IdentityStoreProvisioner{}
	var respDiags diag.Diagnostics
	// attribute_contract
	if !model.AttributeContract.IsNull() {
		attributeContractValue := &client.IdentityStoreProvisionerAttributeContract{}
		attributeContractAttrs := model.AttributeContract.Attributes()
		attributeContractValue.CoreAttributes = []client.Attribute{}
		for _, coreAttributesElement := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.Attribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			attributeContractValue.CoreAttributes = append(attributeContractValue.CoreAttributes, coreAttributesValue)
		}
		attributeContractValue.ExtendedAttributes = []client.Attribute{}
		for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.Attribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		result.AttributeContract = attributeContractValue
	}

	// configuration
	configurationValue := pluginconfiguration.ClientStruct(model.Configuration)
	result.Configuration = *configurationValue

	// group_attribute_contract
	if !model.GroupAttributeContract.IsNull() {
		groupAttributeContractValue := &client.IdentityStoreProvisionerGroupAttributeContract{}
		groupAttributeContractAttrs := model.GroupAttributeContract.Attributes()
		groupAttributeContractValue.CoreAttributes = []client.GroupAttribute{}
		for _, coreAttributesElement := range groupAttributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.GroupAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			groupAttributeContractValue.CoreAttributes = append(groupAttributeContractValue.CoreAttributes, coreAttributesValue)
		}
		groupAttributeContractValue.ExtendedAttributes = []client.GroupAttribute{}
		for _, extendedAttributesElement := range groupAttributeContractAttrs["extended_attributes"].(types.Set).Elements() {
			extendedAttributesValue := client.GroupAttribute{}
			extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
			extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
			groupAttributeContractValue.ExtendedAttributes = append(groupAttributeContractValue.ExtendedAttributes, extendedAttributesValue)
		}
		result.GroupAttributeContract = groupAttributeContractValue
	}

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

	// provisioner_id
	result.Id = model.ProvisionerId.ValueString()
	return result, respDiags
}

func (state *identityStoreProvisionerResourceModel) readClientResponse(response *client.IdentityStoreProvisioner, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
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
	// provisioner_id
	state.ProvisionerId = types.StringValue(response.Id)
	// attribute_contract and group_attribute_contract
	respDiags.Append(state.readClientResponseAttributeContracts(response, isImportRead)...)
	return respDiags
}

func (r *identityStoreProvisionerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data identityStoreProvisionerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.IdentityStoreProvisionersAPI.CreateIdentityStoreProvisioner(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.IdentityStoreProvisionersAPI.CreateIdentityStoreProvisionerExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the identityStoreProvisioner", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *identityStoreProvisionerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var data identityStoreProvisionerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.IdentityStoreProvisionersAPI.GetIdentityStoreProvisioner(config.AuthContext(ctx, r.providerConfig), data.ProvisionerId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Identity Store Provisioner", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the identityStoreProvisioner", err, httpResp, &customId)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, isImportRead)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *identityStoreProvisionerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data identityStoreProvisionerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.IdentityStoreProvisionersAPI.UpdateIdentityStoreProvisioner(config.AuthContext(ctx, r.providerConfig), data.ProvisionerId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.IdentityStoreProvisionersAPI.UpdateIdentityStoreProvisionerExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the identityStoreProvisioner", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *identityStoreProvisionerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data identityStoreProvisionerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.IdentityStoreProvisionersAPI.DeleteIdentityStoreProvisioner(config.AuthContext(ctx, r.providerConfig), data.ProvisionerId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the identityStoreProvisioner", err, httpResp, &customId)
	}
}

func (r *identityStoreProvisionerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to provisioner_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("provisioner_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
