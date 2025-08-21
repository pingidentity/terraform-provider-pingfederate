// Copyright © 2025 Ping Identity Corporation

package oauthaccesstokenmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAccessTokenMappingResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenMappingResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenMappingResource{}

	accessTokenMappingContext = map[string]attr.Type{
		"type":        types.StringType,
		"context_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
)

// OauthAccessTokenMappingResource is a helper function to simplify the provider implementation.
func OauthAccessTokenMappingResource() resource.Resource {
	return &oauthAccessTokenMappingResource{}
}

// oauthAccessTokenMappingResource is the resource implementation.
type oauthAccessTokenMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenMappingResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	MappingId                    types.String `tfsdk:"mapping_id"`
	Context                      types.Object `tfsdk:"context"`
	AccessTokenManagerRef        types.Object `tfsdk:"access_token_manager_ref"`
	AttributeSources             types.List   `tfsdk:"attribute_sources"`
	AttributeContractFulfillment types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria             types.Object `tfsdk:"issuance_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAccessTokenMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OAuth Access Token Mapping",
		Attributes: map[string]schema.Attribute{
			"context": schema.SingleNestedAttribute{
				Description: "The context of the OAuth Access Token Mapping. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The Access Token Mapping Context type. Options are `DEFAULT`, `PCV`, `IDP_CONNECTION`, `IDP_ADAPTER`, `AUTHENTICATION_POLICY_CONTRACT`, `CLIENT_CREDENTIALS`, `TOKEN_EXCHANGE_PROCESSOR_POLICY`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("DEFAULT", "PCV", "IDP_CONNECTION", "IDP_ADAPTER", "AUTHENTICATION_POLICY_CONTRACT", "CLIENT_CREDENTIALS", "TOKEN_EXCHANGE_PROCESSOR_POLICY"),
						},
					},
					"context_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated Access Token Mapping Context instance.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
				},
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "Reference to the access token manager this mapping is associated with. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: resourcelink.ToSchema(),
			},
			"attribute_sources": attributesources.ToSchema(0, true),
			"attribute_contract_fulfillment": schema.MapNestedAttribute{
				Description: "Defines how an attribute in an attribute contract should be populated.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": sourcetypeidkey.ToSchema(false),
						"value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The value for this attribute.",
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"issuance_criteria": issuancecriteria.ToSchema(),
			"mapping_id": schema.StringAttribute{
				Description: "The id of the Access Token Mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *oauthAccessTokenMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_mapping"
}

func (r *oauthAccessTokenMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAccessTokenMappingsResponse(ctx context.Context, r *client.AccessTokenMapping, state *oauthAccessTokenMappingResourceModel) diag.Diagnostics {
	var diags, objDiags diag.Diagnostics

	state.Id = types.StringPointerValue(r.Id)
	state.MappingId = types.StringPointerValue(r.Id)

	contextRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var contextRefValue types.Object
	if r.Context.ContextRef.Id == "" {
		contextRefValue = types.ObjectNull(contextRefAttrTypes)
	} else {
		contextRefValue, objDiags = types.ObjectValue(contextRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(r.Context.ContextRef.Id),
		})
		diags.Append(objDiags...)
	}
	contextAttrValue := map[string]attr.Value{
		"type":        types.StringValue(r.Context.Type),
		"context_ref": contextRefValue,
	}
	state.Context, objDiags = types.ObjectValue(accessTokenMappingContext, contextAttrValue)
	diags.Append(objDiags...)
	state.AccessTokenManagerRef, objDiags = resourcelink.ToState(ctx, &r.AccessTokenManagerRef)
	diags.Append(objDiags...)
	state.AttributeSources, objDiags = attributesources.ToState(ctx, r.AttributeSources)
	diags.Append(objDiags...)
	state.AttributeContractFulfillment, objDiags = attributecontractfulfillment.ToState(ctx, &r.AttributeContractFulfillment)
	diags.Append(objDiags...)
	state.IssuanceCriteria, objDiags = issuancecriteria.ToState(ctx, r.IssuanceCriteria)
	diags.Append(objDiags...)

	// make sure all object type building appends diags
	return diags
}

func (model *oauthAccessTokenMappingResourceModel) buildClientStruct() (*client.AccessTokenMapping, diag.Diagnostics) {
	result := &client.AccessTokenMapping{}
	var respDiags diag.Diagnostics
	// access_token_manager_ref
	accessTokenManagerRefValue := client.ResourceLink{}
	accessTokenManagerRefAttrs := model.AccessTokenManagerRef.Attributes()
	accessTokenManagerRefValue.Id = accessTokenManagerRefAttrs["id"].(types.String).ValueString()
	result.AccessTokenManagerRef = accessTokenManagerRefValue

	// attribute_contract_fulfillment
	result.AttributeContractFulfillment = attributecontractfulfillment.ClientStruct(model.AttributeContractFulfillment)

	// attribute_sources
	result.AttributeSources = attributesources.ClientStruct(model.AttributeSources)

	// context
	contextValue := client.AccessTokenMappingContext{}
	contextAttrs := model.Context.Attributes()
	contextContextRefValue := client.ResourceLink{}
	if internaltypes.IsDefined(contextAttrs["context_ref"]) {
		contextContextRefAttrs := contextAttrs["context_ref"].(types.Object).Attributes()
		contextContextRefValue.Id = contextContextRefAttrs["id"].(types.String).ValueString()
	}
	contextValue.ContextRef = contextContextRefValue
	contextValue.Type = contextAttrs["type"].(types.String).ValueString()
	result.Context = contextValue

	// issuance_criteria
	result.IssuanceCriteria = issuancecriteria.ClientStruct(model.IssuanceCriteria)

	// mapping_id
	result.Id = model.MappingId.ValueStringPointer()
	return result, respDiags
}

func (r *oauthAccessTokenMappingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model *oauthAccessTokenMappingResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if model == nil {
		return
	}
	if internaltypes.IsDefined(model.Context) {
		modelContextType := model.Context.Attributes()["type"].(types.String).ValueString()
		modelContextContextRef := model.Context.Attributes()["context_ref"].(types.Object)
		if (modelContextType == "DEFAULT" || modelContextType == "CLIENT_CREDENTIALS") && internaltypes.IsDefined(modelContextContextRef) {
			resp.Diagnostics.AddAttributeError(
				path.Root("context").AtMapKey("context_ref"),
				providererror.InvalidAttributeConfiguration,
				"context_ref is not required for the Access Token Mapping Context type: "+modelContextType)
		}
	}
}

func (r *oauthAccessTokenMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAccessTokenMappingResourceModel
	var err error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := plan.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateOauthAccessTokenMappings := r.apiClient.OauthAccessTokenMappingsAPI.CreateMapping(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenMappings = apiCreateOauthAccessTokenMappings.Body(*clientData)
	oauthAccessTokenMappingsResponse, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.CreateMappingExecute(apiCreateOauthAccessTokenMappings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Access Token Mapping", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAccessTokenMappingResourceModel

	diags = readOauthAccessTokenMappingsResponse(ctx, oauthAccessTokenMappingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAccessTokenMappingResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenMappings, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.GetMapping(config.AuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Access Token Mapping", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Mapping", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readOauthAccessTokenMappingsResponse(ctx, apiReadOauthAccessTokenMappings, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAccessTokenMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan oauthAccessTokenMappingResourceModel
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := plan.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateOauthAccessTokenMappings := r.apiClient.OauthAccessTokenMappingsAPI.UpdateMapping(config.AuthContext(ctx, r.providerConfig), plan.MappingId.ValueString())
	apiUpdateOauthAccessTokenMappings = apiUpdateOauthAccessTokenMappings.Body(*clientData)
	updateOauthAccessTokenMappingsResponse, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.UpdateMappingExecute(apiUpdateOauthAccessTokenMappings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the OAuth Access Token Mapping", err, httpResp)
		return
	}

	// Read the response
	var state oauthAccessTokenMappingResourceModel
	diags = readOauthAccessTokenMappingsResponse(ctx, updateOauthAccessTokenMappingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthAccessTokenMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.DeleteMapping(config.AuthContext(ctx, r.providerConfig), state.MappingId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Access Token Mapping", err, httpResp)
	}
}

func (r *oauthAccessTokenMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("mapping_id"), req, resp)
}
