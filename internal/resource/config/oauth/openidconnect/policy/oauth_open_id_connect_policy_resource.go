package oauthopenidconnectpolicy

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithConfigure   = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithImportState = &oauthOpenIdConnectPolicyResource{}

	attributeAttrTypes = map[string]attr.Type{
		"name":                 types.StringType,
		"include_in_id_token":  types.BoolType,
		"include_in_user_info": types.BoolType,
		"multi_valued":         types.BoolType,
	}
	attributesListAttrType = types.ListType{
		ElemType: types.ObjectType{AttrTypes: attributeAttrTypes},
	}

	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes":     attributesListAttrType,
		"extended_attributes": attributesListAttrType,
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
	}

	scopeAttributeMappingsElemAttrTypes = map[string]attr.Type{
		"values": types.ListType{ElemType: types.StringType},
	}

	scopeAttributeMappingsDefault, _ = types.MapValue(types.ObjectType{AttrTypes: scopeAttributeMappingsElemAttrTypes}, nil)
)

// OauthOpenIdConnectPolicyResource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPolicyResource() resource.Resource {
	return &oauthOpenIdConnectPolicyResource{}
}

// oauthOpenIdConnectPolicyResource is the resource implementation.
type oauthOpenIdConnectPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthOpenIdConnectPolicyResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	PolicyId                    types.String `tfsdk:"policy_id"`
	Name                        types.String `tfsdk:"name"`
	AccessTokenManagerRef       types.Object `tfsdk:"access_token_manager_ref"`
	IdTokenLifetime             types.Int64  `tfsdk:"id_token_lifetime"`
	IncludeSriInIdToken         types.Bool   `tfsdk:"include_sri_in_id_token"`
	IncludeUserInfoInIdToken    types.Bool   `tfsdk:"include_user_info_in_id_token"`
	IncludeSHashInIdToken       types.Bool   `tfsdk:"include_s_hash_in_id_token"`
	ReturnIdTokenOnRefreshGrant types.Bool   `tfsdk:"return_id_token_on_refresh_grant"`
	ReissueIdTokenInHybridFlow  types.Bool   `tfsdk:"reissue_id_token_in_hybrid_flow"`
	AttributeContract           types.Object `tfsdk:"attribute_contract"`
	AttributeMapping            types.Object `tfsdk:"attribute_mapping"`
	ScopeAttributeMappings      types.Map    `tfsdk:"scope_attribute_mappings"`
}

// GetSchema defines the schema for the resource.
func (r *oauthOpenIdConnectPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OpenID Connect Policy.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name used for display in UI screens.",
				Required:    true,
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The access token manager associated with this Open ID Connect policy.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"id_token_lifetime": schema.Int64Attribute{
				Description: "The ID Token Lifetime, in minutes. The default value is 5.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"include_sri_in_id_token": schema.BoolAttribute{
				Description: "Determines whether a Session Reference Identifier is included in the ID token.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_user_info_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the User Info is always included in the ID token",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_s_hash_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the State Hash should be included in the ID token.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"return_id_token_on_refresh_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when refresh grant is requested or not.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"reissue_id_token_in_hybrid_flow": schema.BoolAttribute{
				Description: "Determines whether a new ID Token should be returned during token request of the hybrid flow.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be returned to OAuth clients in response to requests received at the PingFederate UserInfo endpoint.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of read-only attributes (for example, sub) that are automatically populated by PingFederate.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Optional:    true,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Optional:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional attributes.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"include_in_id_token": schema.BoolAttribute{
									Description: "Attribute is included in the ID Token.",
									Optional:    true,
								},
								"include_in_user_info": schema.BoolAttribute{
									Description: "Attribute is included in the User Info.",
									Optional:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"attribute_sources":              attributesources.ToSchema(),
					"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false),
					"issuance_criteria":              issuancecriteria.ToSchema(),
				},
			},
			"scope_attribute_mappings": schema.MapNestedAttribute{
				Description: "The attribute scope mappings from scopes to attribute names.",
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(scopeAttributeMappingsDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.ListAttribute{
							Description: "A List of values.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "policy_id", false, "The policy ID used internally.")
	resp.Schema = schema
}

func addOptionalOauthOpenIdConnectPolicyFields(ctx context.Context, policy *client.OpenIdConnectPolicy, plan oauthOpenIdConnectPolicyResourceModel) error {
	if internaltypes.IsDefined(plan.IdTokenLifetime) {
		policy.IdTokenLifetime = plan.IdTokenLifetime.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeSriInIdToken) {
		policy.IncludeSriInIdToken = plan.IncludeSriInIdToken.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeUserInfoInIdToken) {
		policy.IncludeUserInfoInIdToken = plan.IncludeUserInfoInIdToken.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSHashInIdToken) {
		policy.IncludeSHashInIdToken = plan.IncludeSHashInIdToken.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ReturnIdTokenOnRefreshGrant) {
		policy.ReturnIdTokenOnRefreshGrant = plan.ReturnIdTokenOnRefreshGrant.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ReissueIdTokenInHybridFlow) {
		policy.ReissueIdTokenInHybridFlow = plan.ReissueIdTokenInHybridFlow.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ScopeAttributeMappings) {
		scopeAttributeMappings := map[string]client.ParameterValues{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ScopeAttributeMappings, false)), &scopeAttributeMappings)
		if err != nil {
			return err
		}
		policy.ScopeAttributeMappings = &scopeAttributeMappings
	}
	return nil
}

// Metadata returns the resource type name.
func (r *oauthOpenIdConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policy"
}

func (r *oauthOpenIdConnectPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readOauthOpenIdConnectPolicyResponse(ctx context.Context, response *client.OpenIdConnectPolicy, state *oauthOpenIdConnectPolicyResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(response.Id)
	state.PolicyId = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)

	state.AccessTokenManagerRef, diags = resourcelink.ToState(ctx, &response.AccessTokenManagerRef)
	respDiags.Append(diags...)

	state.IdTokenLifetime = types.Int64PointerValue(response.IdTokenLifetime)
	state.IncludeSriInIdToken = types.BoolPointerValue(response.IncludeSriInIdToken)
	state.IncludeUserInfoInIdToken = types.BoolPointerValue(response.IncludeUserInfoInIdToken)
	state.IncludeSHashInIdToken = types.BoolPointerValue(response.IncludeSHashInIdToken)
	state.ReturnIdTokenOnRefreshGrant = types.BoolPointerValue(response.ReturnIdTokenOnRefreshGrant)
	state.ReissueIdTokenInHybridFlow = types.BoolPointerValue(response.ReissueIdTokenInHybridFlow)

	state.AttributeContract, diags = types.ObjectValueFrom(ctx, attributeContractAttrTypes, response.AttributeContract)
	respDiags.Append(diags...)

	// Attribute mapping
	attributeMappingValues := map[string]attr.Value{}

	// Build attribute_contract_fulfillment value
	attributeMappingValues["attribute_contract_fulfillment"], diags = attributecontractfulfillment.ToState(ctx, response.AttributeMapping.AttributeContractFulfillment)
	respDiags.Append(diags...)

	// Build issuance_criteria value
	attributeMappingValues["issuance_criteria"], diags = issuancecriteria.ToState(ctx, response.AttributeMapping.IssuanceCriteria)
	respDiags.Append(diags...)

	// Build attribute_sources value
	attributeMappingValues["attribute_sources"], respDiags = attributesources.ToState(ctx, response.AttributeMapping.AttributeSources)
	diags.Append(respDiags...)

	// Build complete attribute mapping value
	state.AttributeMapping, diags = types.ObjectValue(attributeMappingAttrTypes, attributeMappingValues)
	respDiags.Append(diags...)

	state.ScopeAttributeMappings, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: scopeAttributeMappingsElemAttrTypes}, response.ScopeAttributeMappings)
	respDiags.Append(diags...)
	return respDiags
}

// Get fields required for client as client structs
func getRequiredOauthOpenIDConnectPolicyFields(plan oauthOpenIdConnectPolicyResourceModel, diags *diag.Diagnostics) (*client.ResourceLink, *client.OpenIdConnectAttributeContract, *client.AttributeMapping) {
	var accessTokenManagerRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.AccessTokenManagerRef, false)), &accessTokenManagerRef)
	if err != nil {
		diags.AddError("Failed to read access_token_manager_ref from plan", err.Error())
		return nil, nil, nil
	}

	// attribute contract
	var attributeContract client.OpenIdConnectAttributeContract
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), &attributeContract)
	if err != nil {
		diags.AddError("Failed to read attribute_contract from plan", err.Error())
		return nil, nil, nil
	}

	// attribute mapping
	attributeMapping := client.AttributeMapping{}
	planAttrs := plan.AttributeMapping.Attributes()

	attrContractFulfillmentAttr := planAttrs["attribute_contract_fulfillment"].(types.Map)
	attributeMapping.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(attrContractFulfillmentAttr)
	if err != nil {
		diags.AddError("Failed to read attribute_mapping.attribute_contract_fulfillment from plan", err.Error())
		return nil, nil, nil
	}

	issuanceCriteriaAttr := planAttrs["issuance_criteria"].(types.Object)
	attributeMapping.IssuanceCriteria, err = issuancecriteria.ClientStruct(issuanceCriteriaAttr)
	if err != nil {
		diags.AddError("Failed to read attribute_mapping.issuance_criteria from plan", err.Error())
		return nil, nil, nil
	}

	attributeSourcesAttr := planAttrs["attribute_sources"].(types.List)
	attributeMapping.AttributeSources = []client.AttributeSourceAggregation{}
	attributeMapping.AttributeSources, err = attributesources.ClientStruct(attributeSourcesAttr)
	if err != nil {
		diags.AddError("Failed to read attribute_mapping.attribute_sources from plan", err.Error())
		return nil, nil, nil
	}

	return &accessTokenManagerRef, &attributeContract, &attributeMapping
}

func (r *oauthOpenIdConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthOpenIdConnectPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessTokenManagerRef, attributeContract, attributeMapping := getRequiredOauthOpenIDConnectPolicyFields(plan, &resp.Diagnostics)
	if accessTokenManagerRef == nil || attributeContract == nil || attributeMapping == nil {
		// Diagnostics are already added to the response in the above method, just return here
		return
	}

	newOIDCPolicy := client.NewOpenIdConnectPolicy(plan.PolicyId.ValueString(), plan.Name.ValueString(), *accessTokenManagerRef, *attributeContract, *attributeMapping)
	err := addOptionalOauthOpenIdConnectPolicyFields(ctx, newOIDCPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OIDC Policy", err.Error())
		return
	}

	apiCreateOIDCPolicy := r.apiClient.OauthOpenIdConnectAPI.CreateOIDCPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOIDCPolicy = apiCreateOIDCPolicy.Body(*newOIDCPolicy)
	oidcPolicyResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.CreateOIDCPolicyExecute(apiCreateOIDCPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OIDC Policy", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthOpenIdConnectPolicyResourceModel
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, oidcPolicyResponse, &state)
	resp.Diagnostics.Append(readResponseDiags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthOpenIdConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthOpenIdConnectPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOIDCPolicy, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.GetOIDCPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an OIDC Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OIDC Policy", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, apiReadOIDCPolicy, &state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthOpenIdConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthOpenIdConnectPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOIDCPolicyRequest := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.PolicyId.ValueString())

	accessTokenManagerRef, attributeContract, attributeMapping := getRequiredOauthOpenIDConnectPolicyFields(plan, &resp.Diagnostics)
	if accessTokenManagerRef == nil || attributeContract == nil || attributeMapping == nil {
		// Diagnostics are already added to the response in the above method, just return here
		return
	}

	updatedPolicy := client.NewOpenIdConnectPolicy(plan.PolicyId.ValueString(), plan.Name.ValueString(),
		*accessTokenManagerRef, *attributeContract, *attributeMapping)

	err := addOptionalOauthOpenIdConnectPolicyFields(ctx, updatedPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the OIDC Policy", err.Error())
		return
	}

	updateOIDCPolicyRequest = updateOIDCPolicyRequest.Body(*updatedPolicy)
	updateResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCPolicyExecute(updateOIDCPolicyRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the OIDC Policy", err, httpResp)
		return
	}

	// Read the response
	var state oauthOpenIdConnectPolicyResourceModel
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, updateResponse, &state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

func (r *oauthOpenIdConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthOpenIdConnectPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the OIDC Policy", err, httpResp)
	}
}

func (r *oauthOpenIdConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to policy_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("policy_id"), req, resp)
}
