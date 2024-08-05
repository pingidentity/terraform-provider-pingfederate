package idpadapter

import (
	"context"
	"encoding/json"

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
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdapterResource{}
	_ resource.ResourceWithConfigure   = &idpAdapterResource{}
	_ resource.ResourceWithImportState = &idpAdapterResource{}
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
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Required:    true,
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
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Optional:    true,
				Computed:    true,
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
		"The ID of the plugin instance. The ID cannot be modified once the instance is created.")
	resp.Schema = schema
}

func addOptionalIdpAdapterFields(ctx context.Context, addRequest *client.IdpAdapter, plan idpAdapterModel) error {
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
		addRequest.AttributeMapping.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(attrContractFulfillmentAttr)
		if err != nil {
			return err
		}

		issuanceCriteriaAttr := planAttrs["issuance_criteria"].(types.Object)
		addRequest.AttributeMapping.IssuanceCriteria, err = issuancecriteria.ClientStruct(issuanceCriteriaAttr)
		if err != nil {
			return err
		}

		attributeSourcesAttr := planAttrs["attribute_sources"].(types.Set)
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

func (r *idpAdapterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *idpAdapterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil || state == nil {
		return
	}

	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)

	resp.Plan.Set(ctx, plan)
}

func (r *idpAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdapterModel

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

	createIdpAdapter := client.NewIdpAdapter(plan.AdapterId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)
	err = addOptionalIdpAdapterFields(ctx, createIdpAdapter, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}

	apiCreateIdpAdapter := r.apiClient.IdpAdaptersAPI.CreateIdpAdapter(config.AuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapter = apiCreateIdpAdapter.Body(*createIdpAdapter)
	idpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersAPI.CreateIdpAdapterExecute(apiCreateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapter", err, httpResp)
		return
	}

	// Read the response into the state
	var state idpAdapterModel

	readResponseDiags := readIdpAdapterResponse(ctx, idpAdapterResponse, &state, &plan)
	resp.Diagnostics.Append(readResponseDiags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpAdapterModel

	diags := req.State.Get(ctx, &state)
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
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an IdpAdapter", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readResponseDiags := readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, &state)
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

	createUpdateRequest := client.NewIdpAdapter(plan.AdapterId.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)

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
	var state idpAdapterModel
	readResponseDiags := readIdpAdapterResponse(ctx, updateIdpAdapterResponse, &state, &plan)
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the IdP adapter", err, httpResp)
	}
}

func (r *idpAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to adapter_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("adapter_id"), req, resp)
}
