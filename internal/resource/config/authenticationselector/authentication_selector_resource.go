// Copyright © 2025 Ping Identity Corporation

package authenticationselector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
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
	_ resource.Resource                = &authenticationSelectorResource{}
	_ resource.ResourceWithConfigure   = &authenticationSelectorResource{}
	_ resource.ResourceWithImportState = &authenticationSelectorResource{}

	extendedAttributesElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}
	attributeContractAttrType = map[string]attr.Type{
		"extended_attributes": types.SetType{ElemType: extendedAttributesElemType},
	}

	customId = "selector_id"
)

// AuthenticationSelectorsResource is a helper function to simplify the provider implementation.
func AuthenticationSelectorsResource() resource.Resource {
	return &authenticationSelectorResource{}
}

// authenticationSelectorResource is the resource implementation.
type authenticationSelectorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationSelectorResourceModel struct {
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
	SelectorId          types.String `tfsdk:"selector_id"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *authenticationSelectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributeContractDefault, diags := types.ObjectValue(attributeContractAttrType, map[string]attr.Value{
		"extended_attributes": types.SetNull(extendedAttributesElemType),
	})
	resp.Diagnostics.Append(diags...)
	schema := schema.Schema{
		Description: "Manages Authentication Selectors",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. This field is immutable and will trigger a replacement plan if changed.",
				Attributes:  resourcelink.ToSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
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
			"configuration": pluginconfiguration.ToSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the Authentication Selector provides.",
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(attributeContractDefault),
				Attributes: map[string]schema.Attribute{
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A set of additional attributes that can be returned by the Authentication Selector. The extended attributes are only used if the Authentication Selector supports them.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "An attribute for the Authentication Selector attribute contract.",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "selector_id", true, true,
		"The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed.")
	resp.Schema = schema
}

func addOptionalAuthenticationSelectorsFields(addRequest *client.AuthenticationSelector, plan authenticationSelectorResourceModel) {
	// attribute_contract
	if !plan.AttributeContract.IsNull() && !plan.AttributeContract.IsUnknown() {
		attributeContractValue := &client.AuthenticationSelectorAttributeContract{}
		attributeContractAttrs := plan.AttributeContract.Attributes()
		if !attributeContractAttrs["extended_attributes"].IsNull() && !attributeContractAttrs["extended_attributes"].IsUnknown() {
			attributeContractValue.ExtendedAttributes = []client.AuthenticationSelectorAttribute{}
			for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
				extendedAttributesValue := client.AuthenticationSelectorAttribute{}
				extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
				extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
				attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
			}
		}
		addRequest.AttributeContract = attributeContractValue
	}

	// parent_ref
	if !plan.ParentRef.IsNull() {
		parentRefValue := &client.ResourceLink{}
		parentRefAttrs := plan.ParentRef.Attributes()
		parentRefValue.Id = parentRefAttrs["id"].(types.String).ValueString()
		addRequest.ParentRef = parentRefValue
	}
}

// Metadata returns the resource type name.
func (r *authenticationSelectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_selector"
}

func (r *authenticationSelectorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *authenticationSelectorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *authenticationSelectorResourceModel
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

func readAuthenticationSelectorsResponse(ctx context.Context, r *client.AuthenticationSelector, state *authenticationSelectorResourceModel, configurationFromPlan types.Object, isImportRead bool) diag.Diagnostics {
	var diags, objDiags diag.Diagnostics

	state.AttributeContract, objDiags = types.ObjectValueFrom(ctx, attributeContractAttrType, r.AttributeContract)
	diags = append(diags, objDiags...)
	state.SelectorId = types.StringValue(r.Id)
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.ParentRef, objDiags = resourcelink.ToState(ctx, r.ParentRef)
	diags = append(diags, objDiags...)
	state.PluginDescriptorRef, objDiags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	diags = append(diags, objDiags...)
	state.Configuration, objDiags = pluginconfiguration.ToState(configurationFromPlan, &r.Configuration, isImportRead)
	diags = append(diags, objDiags...)

	// make sure all object type building appends diags
	return diags
}

func (r *authenticationSelectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationSelectorResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var hasObjectErrMap = make(map[error]bool)
	pluginDescriptorRef, err := resourcelink.ClientStruct(plan.PluginDescriptorRef)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	for err, hasErr := range hasObjectErrMap {
		if hasErr {
			resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to create an Authentication Selector due to dependent object: "+err.Error())
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationSelectors := client.NewAuthenticationSelector(plan.SelectorId.ValueString(), plan.Name.ValueString(), *pluginDescriptorRef, *configuration)
	addOptionalAuthenticationSelectorsFields(createAuthenticationSelectors, plan)

	apiCreateAuthenticationSelectors := r.apiClient.AuthenticationSelectorsAPI.CreateAuthenticationSelector(config.AuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationSelectors = apiCreateAuthenticationSelectors.Body(*createAuthenticationSelectors)
	authenticationSelectorResponse, httpResp, err := r.apiClient.AuthenticationSelectorsAPI.CreateAuthenticationSelectorExecute(apiCreateAuthenticationSelectors)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the Authentication Selector", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state authenticationSelectorResourceModel

	diags = readAuthenticationSelectorsResponse(ctx, authenticationSelectorResponse, &state, plan.Configuration, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationSelectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state authenticationSelectorResourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationSelectors, httpResp, err := r.apiClient.AuthenticationSelectorsAPI.GetAuthenticationSelector(config.AuthContext(ctx, r.providerConfig), state.SelectorId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Authentication Selector", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting the Authentication Selector", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	diags = readAuthenticationSelectorsResponse(ctx, apiReadAuthenticationSelectors, &state, state.Configuration, isImportRead)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationSelectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan authenticationSelectorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var hasObjectErrMap = make(map[error]bool)
	pluginDescriptorRef, err := resourcelink.ClientStruct(plan.PluginDescriptorRef)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	configuration := pluginconfiguration.ClientStruct(plan.Configuration)

	for err, hasErr := range hasObjectErrMap {
		if hasErr {
			resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to create an Authentication Selector due to dependent object: "+err.Error())
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updateAuthenticationSelectors := r.apiClient.AuthenticationSelectorsAPI.UpdateAuthenticationSelector(config.AuthContext(ctx, r.providerConfig), plan.SelectorId.ValueString())
	createUpdateRequest := client.NewAuthenticationSelector(plan.SelectorId.ValueString(), plan.Name.ValueString(), *pluginDescriptorRef, *configuration)
	addOptionalAuthenticationSelectorsFields(createUpdateRequest, plan)

	updateAuthenticationSelectors = updateAuthenticationSelectors.Body(*createUpdateRequest)
	updateAuthenticationSelectorsResponse, httpResp, err := r.apiClient.AuthenticationSelectorsAPI.UpdateAuthenticationSelectorExecute(updateAuthenticationSelectors)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating an Authentication Selector", err, httpResp, &customId)
		return
	}

	// Read the response
	var state authenticationSelectorResourceModel
	diags = readAuthenticationSelectorsResponse(ctx, updateAuthenticationSelectorsResponse, &state, plan.Configuration, false)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationSelectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationSelectorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationSelectorsAPI.DeleteAuthenticationSelector(config.AuthContext(ctx, r.providerConfig), state.SelectorId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Selector", err, httpResp, &customId)
	}
}

func (r *authenticationSelectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("selector_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
