package idp

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdapterResource{}
	_ resource.ResourceWithConfigure   = &idpAdapterResource{}
	_ resource.ResourceWithImportState = &idpAdapterResource{}
)

// Define attribute types for object types
var (
	// May move some of this into common package if future resources need this
	attributesAttrType = map[string]attr.Type{
		"name":      types.StringType,
		"pseudonym": types.BoolType,
		"masked":    types.BoolType,
	}

	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"core_attributes_all": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"extended_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"unique_user_key_attribute": types.StringType,
		"mask_ognl_values":          types.BoolType,
		"inherited":                 types.BoolType,
	}

	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: attributesources.ElemAttrType(),
			},
		},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria": types.ObjectType{
			AttrTypes: issuancecriteria.AttrType(),
		},
		"inherited": types.BoolType,
	}

	extendedAttributesDefault, _ = types.SetValue(types.ObjectType{
		AttrTypes: attributesAttrType,
	}, nil)
)

// IdpAdapterResource is a helper function to simplify the provider implementation.
func IdpAdapterResource() resource.Resource {
	return &idpAdapterResource{}
}

// idpAdapterResource is the resource implementation.
type idpAdapterResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpAdapterResourceModel struct {
	AuthnCtxClassRef    types.String `tfsdk:"authn_ctx_class_ref"`
	Id                  types.String `tfsdk:"id"`
	CustomId            types.String `tfsdk:"custom_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
	AttributeMapping    types.Object `tfsdk:"attribute_mapping"`
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
}

// GetSchema defines the schema for the resource.
func (r *idpAdapterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Idp Adapter",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Optional:    true,
			},
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
				Attributes:  resourcelink.ToSchema(),
			},
			"configuration": pluginconfiguration.ToSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the IdP adapter provides.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
									Computed:    true,
									// These defaults cause issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default: booldefault.StaticBool(false),
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
									Computed:    true,
									// These defaults cause issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default: booldefault.StaticBool(false),
								},
							},
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"core_attributes_all": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type. This attribute will include any values set by default by PingFederate.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Required:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Required:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(extendedAttributesDefault),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
					},
					"unique_user_key_attribute": schema.StringAttribute{
						Description: "The attribute to use for uniquely identify a user's authentication sessions.",
						Optional:    true,
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to false.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"attribute_sources":              attributesources.ToSchema(),
					"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(false, true),
					"issuance_criteria":              issuancecriteria.ToSchema(),
					"inherited": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
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

func addOptionalIdpAdapterFields(ctx context.Context, addRequest *client.IdpAdapter, plan idpAdapterResourceModel) error {
	var err error
	if internaltypes.IsDefined(plan.AuthnCtxClassRef) {
		addRequest.AuthnCtxClassRef = plan.AuthnCtxClassRef.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ParentRef) {
		addRequest.ParentRef, err = resourcelink.ClientStruct(plan.ParentRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeMapping) {
		addRequest.AttributeMapping = &client.IdpAdapterContractMapping{}
		planAttrs := plan.AttributeMapping.Attributes()

		addRequest.AttributeMapping.Inherited = planAttrs["inherited"].(types.Bool).ValueBoolPointer()

		attrContractFulfillmentAttr := planAttrs["attribute_contract_fulfillment"].(types.Map)
		addRequest.AttributeMapping.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(attrContractFulfillmentAttr)
		if err != nil {
			return err
		}

		issuanceCriteriaAttr := planAttrs["issuance_criteria"].(types.Object)
		addRequest.AttributeMapping.IssuanceCriteria, err = issuancecriteria.ClientStruct(issuanceCriteriaAttr)
		if err != nil {
			return err
		}

		attributeSourcesAttr := planAttrs["attribute_sources"].(types.List)
		addRequest.AttributeMapping.AttributeSources = []client.AttributeSourceAggregation{}
		addRequest.AttributeMapping.AttributeSources, err = attributesources.ClientStruct(attributeSourcesAttr)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = &client.IdpAdapterAttributeContract{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	return nil
}

// Metadata returns the resource type name.
func (r *idpAdapterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_adapter"
}

func (r *idpAdapterResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readIdpAdapterResponse(ctx context.Context, r *client.IdpAdapter, state *idpAdapterResourceModel, plan idpAdapterResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.AuthnCtxClassRef = internaltypes.StringTypeOrNil(r.AuthnCtxClassRef, false)
	state.CustomId = types.StringValue(r.Id)
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, diags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	respDiags.Append(diags...)
	state.ParentRef, diags = resourcelink.ToState(ctx, r.ParentRef)
	respDiags.Append(diags...)
	// Configuration
	state.Configuration, diags = pluginconfiguration.ToState(plan.Configuration, &r.Configuration)
	respDiags.Append(diags...)

	if r.AttributeContract != nil {
		attributeContractValues := map[string]attr.Value{}
		attributeContractValues["extended_attributes"], diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attributesAttrType}, r.AttributeContract.ExtendedAttributes)
		respDiags.Append(diags...)
		attributeContractValues["core_attributes_all"], diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attributesAttrType}, r.AttributeContract.CoreAttributes)
		respDiags.Append(diags...)
		attributeContractValues["unique_user_key_attribute"] = types.StringPointerValue(r.AttributeContract.UniqueUserKeyAttribute)
		attributeContractValues["mask_ognl_values"] = types.BoolPointerValue(r.AttributeContract.MaskOgnlValues)
		attributeContractValues["inherited"] = types.BoolPointerValue(r.AttributeContract.Inherited)

		// Only include core_attributes specified in the plan in the response
		if internaltypes.IsDefined(plan.AttributeContract) && internaltypes.IsDefined(plan.AttributeContract.Attributes()["core_attributes"]) {
			coreAttributes := []attr.Value{}
			planCoreAttributeNames := map[string]bool{}
			for _, planCoreAttr := range plan.AttributeContract.Attributes()["core_attributes"].(types.Set).Elements() {
				planCoreAttributeNames[planCoreAttr.(types.Object).Attributes()["name"].(types.String).ValueString()] = true
			}
			for _, coreAttr := range r.AttributeContract.CoreAttributes {
				_, attrInPlan := planCoreAttributeNames[coreAttr.Name]
				if attrInPlan {
					attrObjVal, diags := types.ObjectValueFrom(ctx, attributesAttrType, coreAttr)
					respDiags.Append(diags...)
					coreAttributes = append(coreAttributes, attrObjVal)
				}
			}
			attributeContractValues["core_attributes"], diags = types.SetValue(types.ObjectType{AttrTypes: attributesAttrType}, coreAttributes)
			respDiags.Append(diags...)
		} else {
			attributeContractValues["core_attributes"] = types.SetNull(types.ObjectType{AttrTypes: attributesAttrType})
		}

		state.AttributeContract, diags = types.ObjectValue(attributeContractAttrTypes, attributeContractValues)
		respDiags.Append(diags...)
	}

	if r.AttributeMapping != nil {
		attributeMappingValues := map[string]attr.Value{
			"inherited": types.BoolPointerValue(r.AttributeMapping.Inherited),
		}
		// The PF API won't return inherited if it is false
		if r.AttributeMapping.Inherited == nil {
			attributeMappingValues["inherited"] = types.BoolValue(false)
		}

		// Build attribute_contract_fulfillment value
		attributeMappingValues["attribute_contract_fulfillment"], diags = attributecontractfulfillment.ToState(ctx, r.AttributeMapping.AttributeContractFulfillment)
		respDiags.Append(diags...)

		// Build issuance_criteria value
		attributeMappingValues["issuance_criteria"], diags = issuancecriteria.ToState(ctx, r.AttributeMapping.IssuanceCriteria)
		respDiags.Append(diags...)

		// Build attribute_sources value
		attributeMappingValues["attribute_sources"], respDiags = attributesources.ToState(ctx, r.AttributeMapping.AttributeSources)
		diags.Append(respDiags...)

		// Build complete attribute mapping value
		state.AttributeMapping, diags = types.ObjectValue(attributeMappingAttrTypes, attributeMappingValues)
		respDiags.Append(diags...)
	}
	return respDiags
}

func (r *idpAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdapterResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createIdpAdapter := client.NewIdpAdapter(plan.CustomId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)
	err = addOptionalIdpAdapterFields(ctx, createIdpAdapter, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}

	apiCreateIdpAdapter := r.apiClient.IdpAdaptersAPI.CreateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapter = apiCreateIdpAdapter.Body(*createIdpAdapter)
	idpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.CreateIdpAdapterExecute(apiCreateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapter", err, httpResp)
		return
	}

	// Read the response into the state
	var state idpAdapterResourceModel

	readResponseDiags := readIdpAdapterResponse(ctx, idpAdapterResponse, &state, plan)
	resp.Diagnostics.Append(readResponseDiags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpAdapterResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpAdapter, httpResp, err := r.apiClient.IdpAdaptersAPI.GetIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an IdpAdapter", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an IdpAdapter", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readResponseDiags := readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpAdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan idpAdapterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	updateIdpAdapter := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.CustomId.ValueString())

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createUpdateRequest := client.NewIdpAdapter(plan.CustomId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)

	err = addOptionalIdpAdapterFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}

	updateIdpAdapter = updateIdpAdapter.Body(*createUpdateRequest)
	updateIdpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapterExecute(updateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdpAdapter", err, httpResp)
		return
	}

	// Read the response
	var state idpAdapterResourceModel
	readResponseDiags := readIdpAdapterResponse(ctx, updateIdpAdapterResponse, &state, plan)
	resp.Diagnostics.Append(readResponseDiags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

// Delete the Idp Adapter
func (r *idpAdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idpAdapterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.IdpAdaptersAPI.DeleteIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Idp Adapter", err, httpResp)
	}
}

func (r *idpAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
