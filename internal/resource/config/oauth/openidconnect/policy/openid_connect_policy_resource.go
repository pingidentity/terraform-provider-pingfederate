// Copyright © 2025 Ping Identity Corporation

package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/api"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributemapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &openidConnectPolicyResource{}
	_ resource.ResourceWithConfigure   = &openidConnectPolicyResource{}
	_ resource.ResourceWithImportState = &openidConnectPolicyResource{}

	customId = "policy_id"
)

// OpenidConnectPolicyResource is a helper function to simplify the provider implementation.
func OpenidConnectPolicyResource() resource.Resource {
	return &openidConnectPolicyResource{}
}

// openidConnectPolicyResource is the resource implementation.
type openidConnectPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *openidConnectPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OpenID Connect Policy.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name used for display in UI screens.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The access token manager associated with this Open ID Connect policy.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"id_token_lifetime": schema.Int64Attribute{
				Description: "The ID Token Lifetime, in minutes. The default value is `5`.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"include_sri_in_id_token": schema.BoolAttribute{
				Description: "Determines whether a Session Reference Identifier is included in the ID token. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_user_info_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the User Info is always included in the ID token. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_s_hash_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the State Hash should be included in the ID token. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"return_id_token_on_refresh_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when refresh grant is requested or not. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"reissue_id_token_in_hybrid_flow": schema.BoolAttribute{
				Description: "Determines whether a new ID Token should be returned during token request of the hybrid flow. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be returned to OAuth clients in response to requests received at the PingFederate UserInfo endpoint.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of read-only attributes (for example, sub) that are automatically populated by PingFederate.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
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
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes.",
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(emptyExtendedAttributesDefault),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
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
									Description: "Indicates whether attribute value is always returned as an array. Defaults to `false`.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
					},
				},
			},
			"attribute_mapping": attributemapping.ToSchema(true),
			"scope_attribute_mappings": schema.MapNestedAttribute{
				Description: "The attribute scope mappings from scopes to attribute names.",
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(scopeAttributeMappingsDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.SetAttribute{
							Description: "A List of values.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"include_x5t_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the X.509 thumbprint header should be included in the ID Token. Supported in PF version `11.3` or later. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				// Default is set in modify plan since it depends on PF version
			},
			"id_token_typ_header_value": schema.StringAttribute{
				Description: "ID Token Type (typ) Header Value. Supported in PF version `11.3` or later.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"return_id_token_on_token_exchange_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when token exchange is requested or not. Defaults to `false`. Supported in PF version `12.2` or later.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "policy_id", true, false, "The policy ID used internally. This field is immutable and will trigger a replacement plan if changed.")
	resp.Schema = schema
}

func (r *openidConnectPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 11.3 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast113 := compare >= 0
	// Compare to version 12.2 of PF
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1220)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast122 := compare >= 0
	var plan *oauthOpenIdConnectPolicyModel
	req.Plan.Get(ctx, &plan)
	if plan == nil {
		return
	}
	planModified := false
	// If include_x5t_in_id_token or id_token_typ_header_value is set prior to PF version 11.3, throw an error
	if !pfVersionAtLeast113 {
		if internaltypes.IsDefined(plan.IncludeX5tInIdToken) {
			version.AddUnsupportedAttributeError("include_x5t_in_id_token",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
		} else if plan.IncludeX5tInIdToken.IsUnknown() {
			plan.IncludeX5tInIdToken = types.BoolNull()
			planModified = true
		}
		if internaltypes.IsDefined(plan.IdTokenTypHeaderValue) {
			version.AddUnsupportedAttributeError("id_token_typ_header_value",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
		} else if plan.IdTokenTypHeaderValue.IsUnknown() {
			plan.IdTokenTypHeaderValue = types.StringNull()
			planModified = true
		}
	}
	// Set default if PF version is new enough
	if pfVersionAtLeast113 && plan.IncludeX5tInIdToken.IsUnknown() {
		plan.IncludeX5tInIdToken = types.BoolValue(false)
		planModified = true
	}

	if !pfVersionAtLeast122 {
		if internaltypes.IsDefined(plan.ReturnIdTokenOnTokenExchangeGrant) {
			version.AddUnsupportedAttributeError("return_id_token_on_token_exchange_grant",
				r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
		} else {
			plan.ReturnIdTokenOnTokenExchangeGrant = types.BoolNull()
			planModified = true
		}
	} else if plan.ReturnIdTokenOnTokenExchangeGrant.IsUnknown() {
		plan.ReturnIdTokenOnTokenExchangeGrant = types.BoolValue(false)
		planModified = true
	}

	if planModified {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

// Metadata returns the resource type name.
func (r *openidConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openid_connect_policy"
}

func (r *openidConnectPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (model *oauthOpenIdConnectPolicyModel) buildClientStruct() (*client.OpenIdConnectPolicy, diag.Diagnostics) {
	result := &client.OpenIdConnectPolicy{}
	var respDiags diag.Diagnostics
	var err error
	// access_token_manager_ref
	accessTokenManagerRefValue := client.ResourceLink{}
	accessTokenManagerRefAttrs := model.AccessTokenManagerRef.Attributes()
	accessTokenManagerRefValue.Id = accessTokenManagerRefAttrs["id"].(types.String).ValueString()
	result.AccessTokenManagerRef = accessTokenManagerRefValue

	// attribute_contract
	attributeContractValue := client.OpenIdConnectAttributeContract{}
	attributeContractAttrs := model.AttributeContract.Attributes()
	attributeContractValue.CoreAttributes = []client.OpenIdConnectAttribute{}
	for _, coreAttributesElement := range attributeContractAttrs["core_attributes"].(types.Set).Elements() {
		coreAttributesValue := client.OpenIdConnectAttribute{}
		coreAttributesAttrs := coreAttributesElement.(types.Object).Attributes()
		coreAttributesValue.IncludeInIdToken = coreAttributesAttrs["include_in_id_token"].(types.Bool).ValueBoolPointer()
		coreAttributesValue.IncludeInUserInfo = coreAttributesAttrs["include_in_user_info"].(types.Bool).ValueBoolPointer()
		coreAttributesValue.MultiValued = coreAttributesAttrs["multi_valued"].(types.Bool).ValueBoolPointer()
		coreAttributesValue.Name = coreAttributesAttrs["name"].(types.String).ValueString()
		attributeContractValue.CoreAttributes = append(attributeContractValue.CoreAttributes, coreAttributesValue)
	}
	attributeContractValue.ExtendedAttributes = []client.OpenIdConnectAttribute{}
	for _, extendedAttributesElement := range attributeContractAttrs["extended_attributes"].(types.Set).Elements() {
		extendedAttributesValue := client.OpenIdConnectAttribute{}
		extendedAttributesAttrs := extendedAttributesElement.(types.Object).Attributes()
		extendedAttributesValue.IncludeInIdToken = extendedAttributesAttrs["include_in_id_token"].(types.Bool).ValueBoolPointer()
		extendedAttributesValue.IncludeInUserInfo = extendedAttributesAttrs["include_in_user_info"].(types.Bool).ValueBoolPointer()
		extendedAttributesValue.MultiValued = extendedAttributesAttrs["multi_valued"].(types.Bool).ValueBoolPointer()
		extendedAttributesValue.Name = extendedAttributesAttrs["name"].(types.String).ValueString()
		attributeContractValue.ExtendedAttributes = append(attributeContractValue.ExtendedAttributes, extendedAttributesValue)
	}
	result.AttributeContract = attributeContractValue

	// attribute_mapping
	attributeMappingValue := client.AttributeMapping{}
	attributeMappingAttrs := model.AttributeMapping.Attributes()
	attributeMappingValue.AttributeContractFulfillment, err = attributecontractfulfillment.ClientStruct(attributeMappingAttrs["attribute_contract_fulfillment"].(types.Map))
	if err != nil {
		respDiags.AddError("Error building client struct for attribute_contract_fulfillment", err.Error())
	}
	attributeMappingValue.AttributeSources, err = attributesources.ClientStruct(attributeMappingAttrs["attribute_sources"].(types.List))
	if err != nil {
		respDiags.AddError("Error building client struct for attribute_sources", err.Error())
	}
	attributeMappingValue.IssuanceCriteria, err = issuancecriteria.ClientStruct(attributeMappingAttrs["issuance_criteria"].(types.Object))
	if err != nil {
		respDiags.AddError("Error building client struct for issuance_criteria", err.Error())
	}
	result.AttributeMapping = attributeMappingValue

	// id_token_lifetime
	result.IdTokenLifetime = model.IdTokenLifetime.ValueInt64Pointer()
	// id_token_typ_header_value
	result.IdTokenTypHeaderValue = model.IdTokenTypHeaderValue.ValueStringPointer()
	// include_shash_in_id_token
	result.IncludeSHashInIdToken = model.IncludeSHashInIdToken.ValueBoolPointer()
	// include_sri_in_id_token
	result.IncludeSriInIdToken = model.IncludeSriInIdToken.ValueBoolPointer()
	// include_user_info_in_id_token
	result.IncludeUserInfoInIdToken = model.IncludeUserInfoInIdToken.ValueBoolPointer()
	// include_x5t_in_id_token
	result.IncludeX5tInIdToken = model.IncludeX5tInIdToken.ValueBoolPointer()
	// name
	result.Name = model.Name.ValueString()
	// oidcthing_id
	result.Id = model.PolicyId.ValueString()
	// reissue_id_token_in_hybrid_flow
	result.ReissueIdTokenInHybridFlow = model.ReissueIdTokenInHybridFlow.ValueBoolPointer()
	// return_id_token_on_refresh_grant
	result.ReturnIdTokenOnRefreshGrant = model.ReturnIdTokenOnRefreshGrant.ValueBoolPointer()
	// return_id_token_on_token_exchange_grant
	result.ReturnIdTokenOnTokenExchangeGrant = model.ReturnIdTokenOnTokenExchangeGrant.ValueBoolPointer()
	// scope_attribute_mappings
	if !model.ScopeAttributeMappings.IsNull() {
		result.ScopeAttributeMappings = &map[string]client.ParameterValues{}
		for key, scopeAttributeMappingsElement := range model.ScopeAttributeMappings.Elements() {
			scopeAttributeMappingsValue := client.ParameterValues{}
			scopeAttributeMappingsAttrs := scopeAttributeMappingsElement.(types.Object).Attributes()
			if !scopeAttributeMappingsAttrs["values"].IsNull() {
				scopeAttributeMappingsValue.Values = []string{}
				for _, valuesElement := range scopeAttributeMappingsAttrs["values"].(types.Set).Elements() {
					scopeAttributeMappingsValue.Values = append(scopeAttributeMappingsValue.Values, valuesElement.(types.String).ValueString())
				}
			}
			(*result.ScopeAttributeMappings)[key] = scopeAttributeMappingsValue
		}
	}

	return result, respDiags
}

func (r *openidConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthOpenIdConnectPolicyModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newOIDCPolicy, diags := plan.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCreateOIDCPolicy := r.apiClient.OauthOpenIdConnectAPI.CreateOIDCPolicy(config.AuthContext(ctx, r.providerConfig))
	apiCreateOIDCPolicy = apiCreateOIDCPolicy.Body(*newOIDCPolicy)
	oidcPolicyResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.CreateOIDCPolicyExecute(apiCreateOIDCPolicy)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the OIDC Policy", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state oauthOpenIdConnectPolicyModel
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, oidcPolicyResponse, &state)
	resp.Diagnostics.Append(readResponseDiags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *openidConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthOpenIdConnectPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOIDCPolicy, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.GetOIDCPolicy(config.AuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OIDC Policy", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting an OIDC Policy", err, httpResp, &customId)
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

func (r *openidConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthOpenIdConnectPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOIDCPolicyRequest := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCPolicy(config.AuthContext(ctx, r.providerConfig), plan.PolicyId.ValueString())

	updatedPolicy, diags := plan.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOIDCPolicyRequest = updateOIDCPolicyRequest.Body(*updatedPolicy)
	updateResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCPolicyExecute(updateOIDCPolicyRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the OIDC Policy", err, httpResp, &customId)
		return
	}

	// Read the response
	var state oauthOpenIdConnectPolicyModel
	readResponseDiags := readOauthOpenIdConnectPolicyResponse(ctx, updateResponse, &state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

func (r *openidConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthOpenIdConnectPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Delete API call logic
	httpResp, err := api.ExponentialBackOffRetryDelete([]int{422},
		r.apiClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(config.AuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the OIDC Policy", err, httpResp, &customId)
	}
}

func (r *openidConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to policy_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("policy_id"), req, resp)
}
