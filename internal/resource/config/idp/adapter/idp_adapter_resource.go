// Copyright © 2025 Ping Identity Corporation

package idpadapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdapterResource{}
	_ resource.ResourceWithConfigure   = &idpAdapterResource{}
	_ resource.ResourceWithImportState = &idpAdapterResource{}

	customId = "adapter_id"
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

// GetSchema defines the schema for the resource.
func (r *idpAdapterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an IdP adapter instance.",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.",
									Optional:    true,
									Computed:    true,
									// These defaults cause issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default: booldefault.StaticBool(false),
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
									Optional:    true,
									Computed:    true,
									// These defaults cause issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default: booldefault.StaticBool(false),
								},
							},
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseNonNullStateForUnknown(),
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
									Computed:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.",
									Computed:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
									Computed:    true,
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
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.",
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to `false`.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"attribute_sources": attributesources.ToSchema(0, false),
					"attribute_contract_fulfillment": schema.MapNestedAttribute{
						Description: "Defines how an attribute in an attribute contract should be populated.",
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": sourcetypeidkey.ToSchema(true),
								"value": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The value for this attribute.",
								},
							},
						},
						Validators: []validator.Map{
							configvalidators.ValidAttributeContractFulfillment(),
						},
					},
					"issuance_criteria": issuancecriteria.ToSchema(),
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"adapter_id",
		true,
		true,
		"The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed.")
	resp.Schema = schema
}

func addOptionalIdpAdapterFields(addRequest *client.IdpAdapter, plan idpAdapterModel) error {
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

		attrContractFulfillmentAttr := planAttrs["attribute_contract_fulfillment"].(types.Map)
		addRequest.AttributeMapping.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(attrContractFulfillmentAttr)

		issuanceCriteriaAttr := planAttrs["issuance_criteria"].(types.Object)
		addRequest.AttributeMapping.IssuanceCriteria = issuancecriteria.ClientStruct(issuanceCriteriaAttr)

		attributeSourcesAttr := planAttrs["attribute_sources"].(types.Set)
		addRequest.AttributeMapping.AttributeSources = []client.AttributeSourceAggregation{}
		addRequest.AttributeMapping.AttributeSources = attributesources.ClientStruct(attributeSourcesAttr)
	}

	// attribute_contract
	if !plan.AttributeContract.IsNull() && !plan.AttributeContract.IsUnknown() {
		attributeContractValue := &client.IdpAdapterAttributeContract{}
		attributeContractAttrs := plan.AttributeContract.Attributes()
		attributeContractValue.CoreAttributes = []client.IdpAdapterAttribute{}
		for _, coreAttributesElement := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
			coreAttributesValue := client.IdpAdapterAttribute{}
			coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
			coreAttributesValue.Masked = coreAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
			coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
			coreAttributesValue.Pseudonym = coreAttributesAttrs["pseudonym"].(types.Bool).ValueBoolPointer()
			attributeContractValue.CoreAttributes = append(attributeContractValue.CoreAttributes, coreAttributesValue)
		}
		if !attributeContractAttrs["extended_attributes"].IsNull() && !attributeContractAttrs["extended_attributes"].IsUnknown() {
			attributeContractValue.ExtendedAttributes = []client.IdpAdapterAttribute{}
			for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.IdpAdapterAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Masked = extendedAttributesAttrs["masked"].(types.Bool).ValueBoolPointer()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				extendedAttributesValue.Pseudonym = extendedAttributesAttrs["pseudonym"].(types.Bool).ValueBoolPointer()
				attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
		}
		attributeContractValue.MaskOgnlValues = attributeContractAttrs["mask_ognl_values"].(types.Bool).ValueBoolPointer()
		attributeContractValue.UniqueUserKeyAttribute = attributeContractAttrs["unique_user_key_attribute"].(types.String).ValueStringPointer()
		addRequest.AttributeContract = attributeContractValue
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

func (r *idpAdapterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *idpAdapterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil {
		return
	}

	// Check that any defined core and extended attributes are included in the contract fulfillment
	if internaltypes.IsDefined(plan.AttributeContract) && internaltypes.IsDefined(plan.AttributeMapping) {
		attributeContractAttrs := plan.AttributeContract.Attributes()
		attributeMappingAttrs := plan.AttributeMapping.Attributes()
		if internaltypes.IsDefined(attributeMappingAttrs["attribute_contract_fulfillment"]) {
			attributeContractFulfillmentKeys := map[string]bool{}
			for key := range attributeMappingAttrs["attribute_contract_fulfillment"].(types.Map).Elements() {
				attributeContractFulfillmentKeys[key] = true
			}

			definedAttrs := []string{}
			for _, attr := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
				attrName := attr.(types.Object).Attributes()["name"].(types.String).ValueString()
				if attrName != "" {
					definedAttrs = append(definedAttrs, attrName)
				}
			}
			for _, attr := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				attrName := attr.(types.Object).Attributes()["name"].(types.String).ValueString()
				if attrName != "" {
					definedAttrs = append(definedAttrs, attrName)
				}
			}
			for _, attrName := range definedAttrs {
				if !attributeContractFulfillmentKeys[attrName] {
					resp.Diagnostics.AddAttributeError(
						path.Root("attribute_mapping").AtMapKey("attribute_contract_fulfillment"),
						providererror.InvalidAttributeConfiguration,
						"attribute_contract_fulfillment must include all core and extended attributes. Missing attribute: "+attrName)
				}
			}
		}
	}

	if state == nil {
		return
	}

	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (r *idpAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdapterModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// plugin_descriptor_ref
	pluginDescriptorRefValue := client.ResourceLink{}
	pluginDescriptorRefAttrs := plan.PluginDescriptorRef.Attributes()
	pluginDescriptorRefValue.Id = pluginDescriptorRefAttrs["id"].(types.String).ValueString()

	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	createIdpAdapter := client.NewIdpAdapter(plan.AdapterId.ValueString(), plan.Name.ValueString(), pluginDescriptorRefValue, *configuration)
	err := addOptionalIdpAdapterFields(createIdpAdapter, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for IdpAdapter: "+err.Error())
		return
	}

	apiCreateIdpAdapter := r.apiClient.IdpAdaptersAPI.CreateIdpAdapter(config.AuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapter = apiCreateIdpAdapter.Body(*createIdpAdapter)
	idpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.CreateIdpAdapterExecute(apiCreateIdpAdapter)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapter", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state idpAdapterModel

	readResponseDiags := readIdpAdapterResponse(ctx, idpAdapterResponse, &state, &plan, false, true)
	resp.Diagnostics.Append(readResponseDiags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state idpAdapterModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadIdpAdapter, httpResp, err := r.apiClient.IdpAdaptersAPI.GetIdpAdapter(config.AuthContext(ctx, r.providerConfig), state.AdapterId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "IdP Adapter", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting an IdpAdapter", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	readResponseDiags := readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, &state, isImportRead, false)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpAdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan idpAdapterModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIdpAdapter := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapter(config.AuthContext(ctx, r.providerConfig), plan.AdapterId.ValueString())

	// plugin_descriptor_ref
	pluginDescriptorRefValue := client.ResourceLink{}
	pluginDescriptorRefAttrs := plan.PluginDescriptorRef.Attributes()
	pluginDescriptorRefValue.Id = pluginDescriptorRefAttrs["id"].(types.String).ValueString()

	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	createUpdateRequest := client.NewIdpAdapter(plan.AdapterId.ValueString(), plan.Name.ValueString(), pluginDescriptorRefValue, *configuration)

	err := addOptionalIdpAdapterFields(createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for IdpAdapter: "+err.Error())
		return
	}

	updateIdpAdapter = updateIdpAdapter.Body(*createUpdateRequest)
	updateIdpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.UpdateIdpAdapterExecute(updateIdpAdapter)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating IdpAdapter", err, httpResp, &customId)
		return
	}

	// Read the response
	var state idpAdapterModel
	readResponseDiags := readIdpAdapterResponse(ctx, updateIdpAdapterResponse, &state, &plan, false, true)
	resp.Diagnostics.Append(readResponseDiags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

// Delete the IdP Adapter
func (r *idpAdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idpAdapterModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.IdpAdaptersAPI.DeleteIdpAdapter(config.AuthContext(ctx, r.providerConfig), state.AdapterId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the IdP adapter", err, httpResp, &customId)
	}
}

func (r *idpAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to adapter_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("adapter_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
