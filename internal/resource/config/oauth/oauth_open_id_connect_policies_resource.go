package oauth

import (
	"context"
	"encoding/json"
  
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthOpenIdConnectPoliciesResource{}
	_ resource.ResourceWithConfigure   = &oauthOpenIdConnectPoliciesResource{}
	_ resource.ResourceWithImportState = &oauthOpenIdConnectPoliciesResource{}
)

// OauthOpenIdConnectPoliciesResource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPoliciesResource() resource.Resource {
	return &oauthOpenIdConnectPoliciesResource{}
}

// oauthOpenIdConnectPoliciesResource is the resource implementation.
type oauthOpenIdConnectPoliciesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthOpenIdConnectPoliciesResourceModel struct {
	Id types.String `tfsdk:"id"`	
	Name types.String `tfsdk:"name"`	
	AccessTokenManagerRef types.Object `tfsdk:"access_token_manager_ref"`	
	IdTokenLifetime types.Int64 `tfsdk:"id_token_lifetime"`	
	IncludeSriInIdToken types.Bool `tfsdk:"include_sri_in_id_token"`	
	IncludeUserInfoInIdToken types.Bool `tfsdk:"include_user_info_in_id_token"`	
	IncludeSHashInIdToken types.Bool `tfsdk:"include_s_hash_in_id_token"`	
	ReturnIdTokenOnRefreshGrant types.Bool `tfsdk:"return_id_token_on_refresh_grant"`	
	ReissueIdTokenInHybridFlow types.Bool `tfsdk:"reissue_id_token_in_hybrid_flow"`	
	AttributeContract types.Object `tfsdk:"attribute_contract"`	
	AttributeMapping types.Object `tfsdk:"attribute_mapping"`	
	ScopeAttributeMappings types.Object `tfsdk:"scope_attribute_mappings"`
}

// GetSchema defines the schema for the resource.
func (r *oauthOpenIdConnectPoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages OpenID Connect policies.",
		Attributes: map[string]schema.Attribute{
	    "items": schema.SetNestedAttribute{
	    	Description: "The actual list of policies.",
	    	Computed:    true,
	    	Optional:    true,
	    	PlanModifiers: []planmodifier.Set{
	    		setplanmodifier.UseStateForUnknown(),
	    	},
			"id": schema.StringAttribute{
				Description: "The policy ID used internally.",
				Computed:    true,
				Optional: 	false,
			},
			"name": schema.StringAttribute{
				Description: "The name used for display in UI screens.",
				Computed:    true,
				Optional: 	false,
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The access token manager associated with this Open ID Connect policy.",
				Required: 	true,
				Attributes: resourceLink.ToSchema(),
			},
			"id_token_lifetime": schema.Int64Attribute{
				Description: "The ID Token Lifetime, in minutes. The default value is 5.",
	    },
	    "include_sri_in_id_token": schema.BoolAttribute{
				Description: "Determines whether a Session Reference Identifier is included in the ID token.",
			},
			"include_user_info_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the User Info is always included in the ID token",
			},
			"include_s_hash_in_id_token": schema.BoolAttribute{
				Description: "Determines whether the State Hash should be included in the ID token.",
			},
			"return_id_token_on_refresh_grant": schema.BoolAttribute{
				Description: "Determines whether an ID Token should be returned when refresh grant is requested or not.",
			},
			"reissue_id_token_in_hybrid_flow": schema.BoolAttribute{
				Description: "Determines whether a new ID Token should be returned during token request of the hybrid flow.",
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be added to an access token.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of core token attributes that are associated with the access token management plugin type. This field is read-only and is ignored on POST/PUT.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    false,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional token attributes that are associated with this access token management plugin instance.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
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
				"scope_attribute_mappings": schema.SingleNestedAttribute{
					Description: "The attribute scope mappings from scopes to attribute names.",
					Optional:    true,
					Computed:    true,
					Attributes: map[string]schema.Attribute{
						"values": schema.ListNestedAttribute{
							Description: "A List of values.",
						},
					},
				},
			},
		},
	},
}

	// Set attributes in string list
	// if setOptionalToComputed {
		// config.SetAllAttributesToOptionalAndComputed(&schema, []string{"FIX_ME"})
	// }
	// config.AddCommonSchema(&schema, false)
	// resp.Schema = schema
// }

func addOptionalOauthOpenIdConnectPoliciesFields(ctx context.Context, addRequest *client.Policies, plan oauthOpenIdConnectPoliciesResourceModel) error {
	
	if internaltypes.IsDefined(plan.Id) {
	addRequest.Id = plan.Id.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.Name) {
	addRequest.Name = plan.Name.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.AccessTokenManagerRef) {
		addRequest.AccessTokenManagerRef = client.NewAccessTokenManagerRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AccessTokenManagerRef, false)), addRequest.AccessTokenManagerRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.IdTokenLifetime) {
	addRequest.IdTokenLifetime = plan.IdTokenLifetime.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.IncludeSriInIdToken) {
	addRequest.IncludeSriInIdToken = plan.IncludeSriInIdToken.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.IncludeUserInfoInIdToken) {
	addRequest.IncludeUserInfoInIdToken = plan.IncludeUserInfoInIdToken.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.IncludeSHashInIdToken) {
	addRequest.IncludeSHashInIdToken = plan.IncludeSHashInIdToken.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.ReturnIdTokenOnRefreshGrant) {
	addRequest.ReturnIdTokenOnRefreshGrant = plan.ReturnIdTokenOnRefreshGrant.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.ReissueIdTokenInHybridFlow) {
	addRequest.ReissueIdTokenInHybridFlow = plan.ReissueIdTokenInHybridFlow.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = client.NewAttributeContract()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeMapping) {
		addRequest.AttributeMapping = client.NewAttributeMapping()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeMapping, false)), addRequest.AttributeMapping)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ScopeAttributeMappings) {
		addRequest.ScopeAttributeMappings = client.NewScopeAttributeMappings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ScopeAttributeMappings, false)), addRequest.ScopeAttributeMappings)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *oauthOpenIdConnectPoliciesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policies"
}

func (r *oauthOpenIdConnectPoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthOpenIdConnectPoliciesResponse(ctx context.Context, r *client.Policies, state *oauthOpenIdConnectPoliciesResourceModel) {
state.Id = internaltypes.StringTypeOrNil(r.Id)	
	state.Name = internaltypes.StringTypeOrNil(r.Name)	
	state.AccessTokenManagerRef = (r.AccessTokenManagerRef)	
	state.IdTokenLifetime = types.Int64Value(r.IdTokenLifetime)	
	state.IncludeSriInIdToken = types.BoolValue(r.IncludeSriInIdToken)	
	state.IncludeUserInfoInIdToken = types.BoolValue(r.IncludeUserInfoInIdToken)	
	state.IncludeSHashInIdToken = types.BoolValue(r.IncludeSHashInIdToken)	
	state.ReturnIdTokenOnRefreshGrant = types.BoolValue(r.ReturnIdTokenOnRefreshGrant)	
	state.ReissueIdTokenInHybridFlow = types.BoolValue(r.ReissueIdTokenInHybridFlow)	
	state.AttributeContract = (r.AttributeContract)	
	state.AttributeMapping = (r.AttributeMapping)	
	state.ScopeAttributeMappings = You will need to figure out what needs to go into the object(r.ScopeAttributeMappings)
}

func (r *oauthOpenIdConnectPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthOpenIdConnectPoliciesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthOpenIdConnectPolicies := client.NewOpenIdConnectPolicy()
	err := addOptionalOauthOpenIdConnectPoliciesFields(ctx, createOauthOpenIdConnectPolicies, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthOpenIdConnectPolicies", err.Error())
		return
	}
	requestJson, err := createOauthOpenIdConnectPolicies.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthOpenIdConnectPolicies := r.apiClient.OauthApi.AddOpenIdConnectPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthOpenIdConnectPolicies = apiCreateOauthOpenIdConnectPolicies.Body(*createOauthOpenIdConnectPolicies)
	oauthOpenIdConnectPoliciesResponse, httpResp, err := r.apiClient.OauthApi.AddOpenIdConnectPolicyExecute(apiCreateOauthOpenIdConnectPolicies)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthOpenIdConnectPolicies", err, httpResp)
		return
	}
	responseJson, err := oauthOpenIdConnectPoliciesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthOpenIdConnectPoliciesResourceModel

	readOauthOpenIdConnectPoliciesResponse(ctx, oauthOpenIdConnectPoliciesResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthOpenIdConnectPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthOpenIdConnectPoliciesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthOpenIdConnectPolicies, httpResp, err := r.apiClient.OauthApi.GetOpenIdConnectPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.VALUE.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OauthOpenIdConnectPolicies", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  OauthOpenIdConnectPolicies", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadOauthOpenIdConnectPolicies.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthOpenIdConnectPoliciesResponse(ctx, apiReadOauthOpenIdConnectPolicies, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthOpenIdConnectPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthOpenIdConnectPoliciesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthOpenIdConnectPoliciesResourceModel
	req.State.Get(ctx, &state)
	updateOauthOpenIdConnectPolicies := r.apiClient.OauthApi.UpdateOpenIdConnectPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.VALUE.ValueString())
	createUpdateRequest := client.NewOpenIdConnectPolicy()
	err := addOptionalOauthOpenIdConnectPoliciesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthOpenIdConnectPolicies", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthOpenIdConnectPolicies = updateOauthOpenIdConnectPolicies.Body(*createUpdateRequest)
	updateOauthOpenIdConnectPoliciesResponse, httpResp, err := r.apiClient.OauthApi.UpdateOpenIdConnectPolicyExecute(updateOauthOpenIdConnectPolicies)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthOpenIdConnectPolicies", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthOpenIdConnectPoliciesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthOpenIdConnectPoliciesResponse(ctx, updateOauthOpenIdConnectPoliciesResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthOpenIdConnectPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthOpenIdConnectPoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
