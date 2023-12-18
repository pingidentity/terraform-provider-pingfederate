package oauthauthserversettings

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopeentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsResource{}

	scopesDefault, _ = types.SetValue(types.ObjectType{AttrTypes: scopeentry.AttrTypes()}, nil)

	scopeGroupsDefault, _ = types.SetValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"scopes":      types.SetType{ElemType: types.StringType},
	}}, nil)
	persistentGrantReuseGrantTypesDefault, _ = types.SetValue(types.StringType, nil)
	allowedOriginsDefault, _                 = types.ListValue(types.StringType, nil)
	defaultCoreAttribute1, _                 = types.ObjectValue(attributeAttrTypes, map[string]attr.Value{
		"name": types.StringValue("USER_KEY"),
	})
	defaultCoreAttribute2, _ = types.ObjectValue(attributeAttrTypes, map[string]attr.Value{
		"name": types.StringValue("USER_NAME"),
	})

	coreAttributesDefault, _ = types.SetValue(attributeSetElementType, []attr.Value{
		defaultCoreAttribute1,
		defaultCoreAttribute2,
	})
	extendedAttributesDefault, _ = types.SetValue(attributeSetElementType, nil)

	persistentGrantContactDefault, _ = types.ObjectValue(map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: attributeSetElementType},
		"extended_attributes": types.SetType{ElemType: attributeSetElementType},
	}, map[string]attr.Value{
		"core_attributes":     coreAttributesDefault,
		"extended_attributes": extendedAttributesDefault,
	})
)

// OauthAuthServerSettingsResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsResource() resource.Resource {
	return &oauthAuthServerSettingsResource{}
}

// oauthAuthServerSettingsResource is the resource implementation.
type oauthAuthServerSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the OAuth authorization server settings.",
		Attributes: map[string]schema.Attribute{
			"default_scope_description": schema.StringAttribute{
				Description: "The default scope description.",
				Required:    true,
			},
			"scopes": schema.SetNestedAttribute{
				Description: "The list of common scopes.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(scopesDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope.",
							Required:    true,
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope that appears when the user is prompted for authorization.",
							Required:    true,
						},
						"dynamic": schema.BoolAttribute{
							Description: "True if the scope is dynamic. (Defaults to false)",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
			"scope_groups": schema.SetNestedAttribute{
				Description: "The list of common scope groups.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(scopeGroupsDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope group.",
							Required:    true,
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    true,
						},
						"scopes": schema.SetAttribute{
							Description: "The set of scopes for this scope group.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"exclusive_scopes": schema.SetNestedAttribute{
				Description: "The list of exclusive scopes.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(scopesDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope.",
							Required:    true,
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope that appears when the user is prompted for authorization.",
							Required:    true,
						},
						"dynamic": schema.BoolAttribute{
							Description: "True if the scope is dynamic. (Defaults to false)",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
			"exclusive_scope_groups": schema.SetNestedAttribute{
				Description: "The list of exclusive scope groups.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(scopeGroupsDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope group.",
							Required:    true,
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    true,
						},
						"scopes": schema.SetAttribute{
							Description: "The set of scopes for this scope group.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"authorization_code_timeout": schema.Int64Attribute{
				Description: "The authorization code timeout, in seconds.",
				Required:    true,
			},
			"authorization_code_entropy": schema.Int64Attribute{
				Description: "The authorization code entropy, in bytes.",
				Required:    true,
			},
			"disallow_plain_pkce": schema.BoolAttribute{
				Description: "Determines whether PKCE's 'plain' code challenge method will be disallowed. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_issuer_in_authorization_response": schema.BoolAttribute{
				Description: "Determines whether the authorization server's issuer value is added to the authorization response or not. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"track_user_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether user sessions are tracked for logout. If this property is not provided on a PUT, the setting is left unchanged.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"token_endpoint_base_url": schema.StringAttribute{
				Description: "The token endpoint base URL used to validate the 'aud' claim during Private Key JWT Client Authentication.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"persistent_grant_lifetime": schema.Int64Attribute{
				Description: "The persistent grant lifetime. The default value is indefinite. -1 indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(-1),
			},
			"persistent_grant_lifetime_unit": schema.StringAttribute{
				Description: "The persistent grant lifetime unit.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. The default value is 30 (days). -1 indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(30),
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit. The default value is DAYS",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"refresh_token_length": schema.Int64Attribute{
				Description: "The refresh token length in number of characters.",
				Required:    true,
			},
			"roll_refresh_token_values": schema.BoolAttribute{
				Description: "The roll refresh token values default policy. The default value is true.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds. The default value is 60.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"refresh_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens, in hours.",
				Required:    true,
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(persistentGrantReuseGrantTypesDefault),
				ElementType: types.StringType,
			},
			"persistent_grant_contract": schema.SingleNestedAttribute{
				Description: "The persistent grant contract defines attributes that are associated with OAuth persistent grants.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(persistentGrantContactDefault),
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "This is a read-only list of persistent grant attributes and includes USER_KEY and USER_NAME. Changes to this field will be ignored.",
						Computed:    true,
						Optional:    false,
						Default:     setdefault.StaticValue(coreAttributesDefault),
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
							},
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes for the persistent grant contract.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
							},
						},
					},
				},
			},
			"bypass_authorization_for_approved_grants": schema.BoolAttribute{
				Description: "Bypass authorization for previously approved persistent grants. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_unidentified_client_ro_creds": schema.BoolAttribute{
				Description: "Allow unidentified clients to request resource owner password credentials grants. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_unidentified_client_extension_grants": schema.BoolAttribute{
				Description: "Allow unidentified clients to request extension grants. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"admin_web_service_pcv_ref": schema.SingleNestedAttribute{
				Description: "The password credential validator reference that is used for authenticating access to the OAuth Administrative Web Service.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"atm_id_for_oauth_grant_management": schema.StringAttribute{
				Description: "The ID of the Access Token Manager used for OAuth enabled grant management.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"scope_for_oauth_grant_management": schema.StringAttribute{
				Description: "The OAuth scope to validate when accessing grant management service.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"allowed_origins": schema.ListAttribute{
				Description: "The list of allowed origins.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(allowedOriginsDefault),
				Validators: []validator.List{
					configvalidators.ValidUrls(),
				},
			},
			"user_authorization_url": schema.StringAttribute{
				Description: "The URL used to generate 'verification_url' and 'verification_url_complete' values in a Device Authorization request",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"registered_authorization_path": schema.StringAttribute{
				Description: "The Registered Authorization Path is concatenated to PingFederate base URL to generate 'verification_url' and 'verification_url_complete' values in a Device Authorization request. PingFederate listens to this path if specified",
				Required:    true,
				Validators: []validator.String{
					configvalidators.StartsWith("/"),
				},
			},
			"pending_authorization_timeout": schema.Int64Attribute{
				Description: "The 'device_code' and 'user_code' timeout, in seconds.",
				Required:    true,
			},
			"device_polling_interval": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds.",
				Required:    true,
			},
			"activation_code_check_mode": schema.StringAttribute{
				Description: "Determines whether the user is prompted to enter or confirm the activation code after authenticating or before. The default is AFTER_AUTHENTICATION.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("AFTER_AUTHENTICATION"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"AFTER_AUTHENTICATION", "BEFORE_AUTHENTICATION"}...),
				},
			},
			"bypass_activation_code_confirmation": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device.",
				Required:    true,
			},
			"user_authorization_consent_page_setting": schema.StringAttribute{
				Description: "User Authorization Consent Page setting to use PingFederate's internal consent page or an external system",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("INTERNAL"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"INTERNAL", "ADAPTER"}...),
				},
			},
			"user_authorization_consent_adapter": schema.StringAttribute{
				Description: "Adapter ID of the external consent adapter to be used for the consent page user interface.",
				Optional:    true,
			},
			"approved_scopes_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved scopes returned by the external consent page.",
				Optional:    true,
			},
			"approved_authorization_detail_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved authorization details returned by the external consent page.",
				Optional:    true,
			},
			"par_reference_timeout": schema.Int64Attribute{
				Description: "The timeout, in seconds, of the pushed authorization request reference. The default value is 60.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"par_reference_length": schema.Int64Attribute{
				Description: "The entropy of pushed authorization request references, in bytes. The default value is 24.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(24),
			},
			"par_status": schema.StringAttribute{
				Description: "The status of pushed authorization request support. The default value is ENABLED.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("ENABLED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DISABLED", "ENABLED", "REQUIRED"}...),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"jwt_secured_authorization_response_mode_lifetime": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the JWT Secured authorization response. The default value is 600.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(600),
			},
			"dpop_proof_require_nonce": schema.BoolAttribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "Determines whether nonce is required in the Demonstrating Proof-of-Possession (DPoP) proof JWT. The default value is false. Supported in PF version 11.3 or later.",
				Computed:    true,
				Optional:    true,
			},
			"dpop_proof_lifetime_seconds": schema.Int64Attribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "The lifetime, in seconds, of the Demonstrating Proof-of-Possession (DPoP) proof JWT. The default value is 120. Supported in PF version 11.3 or later.",
				Computed:    true,
				Optional:    true,
			},
			"dpop_proof_enforce_replay_prevention": schema.BoolAttribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "Determines whether Demonstrating Proof-of-Possession (DPoP) proof JWT replay prevention is enforced. The default value is false. Supported in PF version 11.3 or later.",
				Computed:    true,
				Optional:    true,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	var model oauthAuthServerSettingsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	// Scope list for comparing values in matchNameBtwnScopes variable
	scopeNames := []string{}
	// Test scope names for dynamic true, string must be prepended with *
	if internaltypes.IsDefined(model.Scopes) {
		scopeElems := model.Scopes.Elements()
		for _, scopeElem := range scopeElems {
			scopeElemObjectAttrs := scopeElem.(types.Object)
			scopeEntryName := scopeElemObjectAttrs.Attributes()["name"].(basetypes.StringValue).ValueString()
			scopeNames = append(scopeNames, scopeEntryName)
			scopeEntryIsDynamic := scopeElemObjectAttrs.Attributes()["dynamic"].(basetypes.BoolValue).ValueBool()
			if scopeEntryIsDynamic {
				if strings.Index(scopeEntryName, "*") != 0 {
					resp.Diagnostics.AddError("Scope name conflict!", fmt.Sprintf("Scope name \"%s\" must be prefixed with a \"*\" when dynamic is set to true.", scopeEntryName))
				}
			}
		}
	}

	// Test exclusive scope names for dynamic true, string must be prepended with *
	eScopeNames := []string{}
	if internaltypes.IsDefined(model.ExclusiveScopes) {
		exclusiveScopeElems := model.ExclusiveScopes.Elements()
		for _, esElem := range exclusiveScopeElems {
			esElemObjectAttrs := esElem.(types.Object)
			eScopeEntryName := esElemObjectAttrs.Attributes()["name"].(basetypes.StringValue).ValueString()
			eScopeNames = append(eScopeNames, eScopeEntryName)
			eScopeEntryIsDynamic := esElemObjectAttrs.Attributes()["dynamic"].(basetypes.BoolValue).ValueBool()
			if eScopeEntryIsDynamic {
				if strings.Index(eScopeEntryName, "*") != 0 {
					resp.Diagnostics.AddError("Exclusive scope name conflict!", fmt.Sprintf("Scope name \"%s\" must be prefixed with a \"*\" when dynamic is set to true.", eScopeEntryName))
				}
			}
		}
	}

	// Test if values in sets match
	matchVal := internaltypes.MatchStringInSets(scopeNames, eScopeNames)
	if matchVal != nil {
		resp.Diagnostics.AddError("Scope name conflict!", fmt.Sprintf("The scope name \"%s\" is already defined in another scope list", *matchVal))
	}
}

func (r *oauthAuthServerSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 11.3 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	var plan oauthAuthServerSettingsModel
	req.Plan.Get(ctx, &plan)
	// If any of these fields are set by the user and the PF version is not new enough, throw an error
	if compare < 0 {
		if internaltypes.IsDefined(plan.DpopProofEnforceReplayPrevention) {
			resp.Diagnostics.AddError("Attribute 'dpop_proof_enforce_replay_prevention' not supported by PingFederate version "+r.providerConfig.ProductVersion, "")
		} else if plan.DpopProofEnforceReplayPrevention.IsUnknown() {
			// Set a null default when the version isn't new enough to use this attribute
			plan.DpopProofEnforceReplayPrevention = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.DpopProofLifetimeSeconds) {
			resp.Diagnostics.AddError("Attribute 'dpop_proof_lifetime_seconds' not supported by PingFederate version "+r.providerConfig.ProductVersion, "")
		} else if plan.DpopProofLifetimeSeconds.IsUnknown() {
			plan.DpopProofLifetimeSeconds = types.Int64Null()
		}

		if internaltypes.IsDefined(plan.DpopProofRequireNonce) {
			resp.Diagnostics.AddError("Attribute 'dpop_proof_require_nonce' not supported by PingFederate version "+r.providerConfig.ProductVersion, "")
		} else if plan.DpopProofRequireNonce.IsUnknown() {
			plan.DpopProofRequireNonce = types.BoolNull()
		}
	} else { //PF version is new enough for these attributes, set defaults
		if plan.DpopProofEnforceReplayPrevention.IsUnknown() {
			plan.DpopProofEnforceReplayPrevention = types.BoolValue(false)
		}

		if plan.DpopProofLifetimeSeconds.IsUnknown() {
			plan.DpopProofLifetimeSeconds = types.Int64Value(120)
		}

		if plan.DpopProofRequireNonce.IsUnknown() {
			plan.DpopProofRequireNonce = types.BoolValue(false)
		}
	}

	if !resp.Diagnostics.HasError() {
		resp.Plan.Set(ctx, &plan)
	}
}

func addOptionalOauthAuthServerSettingsFields(ctx context.Context, addRequest *client.AuthorizationServerSettings, plan oauthAuthServerSettingsModel) error {

	if internaltypes.IsDefined(plan.Scopes) {
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Scopes, false)), &addRequest.Scopes)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ScopeGroups) {
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ScopeGroups, false)), &addRequest.ScopeGroups)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ExclusiveScopes) {
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ExclusiveScopes, false)), &addRequest.ExclusiveScopes)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ExclusiveScopeGroups) {
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ExclusiveScopeGroups, false)), &addRequest.ExclusiveScopeGroups)
		if err != nil {
			return err
		}
	}

	addRequest.DisallowPlainPKCE = plan.DisallowPlainPKCE.ValueBoolPointer()
	addRequest.IncludeIssuerInAuthorizationResponse = plan.IncludeIssuerInAuthorizationResponse.ValueBoolPointer()
	addRequest.TrackUserSessionsForLogout = plan.TrackUserSessionsForLogout.ValueBoolPointer()
	addRequest.TokenEndpointBaseUrl = plan.TokenEndpointBaseUrl.ValueStringPointer()
	addRequest.PersistentGrantLifetime = plan.PersistentGrantLifetime.ValueInt64Pointer()
	addRequest.PersistentGrantLifetimeUnit = plan.PersistentGrantLifetimeUnit.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeout = plan.PersistentGrantIdleTimeout.ValueInt64Pointer()
	addRequest.PersistentGrantIdleTimeoutTimeUnit = plan.PersistentGrantIdleTimeoutTimeUnit.ValueStringPointer()
	addRequest.RollRefreshTokenValues = plan.RollRefreshTokenValues.ValueBoolPointer()
	addRequest.RefreshTokenRollingGracePeriod = plan.RefreshTokenRollingGracePeriod.ValueInt64Pointer()
	var persistentGrantReuseTypes []string
	plan.PersistentGrantReuseGrantTypes.ElementsAs(ctx, &persistentGrantReuseTypes, false)
	addRequest.PersistentGrantReuseGrantTypes = persistentGrantReuseTypes

	if internaltypes.IsDefined(plan.PersistentGrantContract) {
		addRequest.PersistentGrantContract = client.NewPersistentGrantContractWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.PersistentGrantContract, false)), addRequest.PersistentGrantContract)
		if err != nil {
			return err
		}
	}
	addRequest.BypassAuthorizationForApprovedGrants = plan.BypassAuthorizationForApprovedGrants.ValueBoolPointer()
	addRequest.AllowUnidentifiedClientROCreds = plan.AllowUnidentifiedClientROCreds.ValueBoolPointer()
	addRequest.AllowUnidentifiedClientExtensionGrants = plan.AllowUnidentifiedClientExtensionGrants.ValueBoolPointer()

	if internaltypes.IsDefined(plan.AdminWebServicePcvRef) {
		addRequest.AdminWebServicePcvRef = client.NewResourceLinkWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AdminWebServicePcvRef, false)), addRequest.AdminWebServicePcvRef)
		if err != nil {
			return err
		}
	}
	addRequest.AtmIdForOAuthGrantManagement = plan.AtmIdForOAuthGrantManagement.ValueStringPointer()
	addRequest.ScopeForOAuthGrantManagement = plan.ScopeForOAuthGrantManagement.ValueStringPointer()
	var allowedOrigins []string
	plan.AllowedOrigins.ElementsAs(ctx, &allowedOrigins, false)
	addRequest.AllowedOrigins = allowedOrigins
	addRequest.UserAuthorizationUrl = plan.UserAuthorizationUrl.ValueStringPointer()
	addRequest.DevicePollingInterval = plan.DevicePollingInterval.ValueInt64()
	addRequest.ActivationCodeCheckMode = plan.ActivationCodeCheckMode.ValueStringPointer()
	addRequest.UserAuthorizationConsentPageSetting = plan.UserAuthorizationConsentPageSetting.ValueStringPointer()
	addRequest.UserAuthorizationConsentAdapter = plan.UserAuthorizationConsentAdapter.ValueStringPointer()
	addRequest.ApprovedScopesAttribute = plan.ApprovedScopesAttribute.ValueStringPointer()
	addRequest.ApprovedAuthorizationDetailAttribute = plan.ApprovedAuthorizationDetailAttribute.ValueStringPointer()
	addRequest.ParReferenceTimeout = plan.ParReferenceTimeout.ValueInt64Pointer()
	addRequest.ParReferenceLength = plan.ParReferenceLength.ValueInt64Pointer()
	addRequest.ParStatus = plan.ParStatus.ValueStringPointer()
	addRequest.ClientSecretRetentionPeriod = plan.ClientSecretRetentionPeriod.ValueInt64Pointer()

	addRequest.JwtSecuredAuthorizationResponseModeLifetime = plan.JwtSecuredAuthorizationResponseModeLifetime.ValueInt64Pointer()
	addRequest.DpopProofEnforceReplayPrevention = plan.DpopProofEnforceReplayPrevention.ValueBoolPointer()
	addRequest.DpopProofLifetimeSeconds = plan.DpopProofLifetimeSeconds.ValueInt64Pointer()
	addRequest.DpopProofRequireNonce = plan.DpopProofRequireNonce.ValueBoolPointer()

	return nil

}

// Metadata returns the resource type name.
func (r *oauthAuthServerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings"
}

func (r *oauthAuthServerSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *oauthAuthServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAuthServerSettings := client.NewAuthorizationServerSettings(plan.DefaultScopeDescription.ValueString(), plan.AuthorizationCodeTimeout.ValueInt64(), plan.AuthorizationCodeEntropy.ValueInt64(), plan.RefreshTokenLength.ValueInt64(), plan.RefreshRollingInterval.ValueInt64(), plan.RegisteredAuthorizationPath.ValueString(), plan.PendingAuthorizationTimeout.ValueInt64(), plan.DevicePollingInterval.ValueInt64(), plan.BypassActivationCodeConfirmation.ValueBool())
	err := addOptionalOauthAuthServerSettingsFields(ctx, createOauthAuthServerSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Auth Server Settings", err.Error())
		return
	}

	apiCreateOauthAuthServerSettings := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettings = apiCreateOauthAuthServerSettings.Body(*createOauthAuthServerSettings)
	oauthAuthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(apiCreateOauthAuthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAuthServerSettingsModel
	diags = readOauthAuthServerSettingsResponse(ctx, oauthAuthServerSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAuthServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAuthServerSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettings, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetAuthorizationServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Auth Server Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Auth Server Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readOauthAuthServerSettingsResponse(ctx, apiReadOauthAuthServerSettings, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthAuthServerSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	updateOauthAuthServerSettings := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthorizationServerSettings(plan.DefaultScopeDescription.ValueString(), plan.AuthorizationCodeTimeout.ValueInt64(), plan.AuthorizationCodeEntropy.ValueInt64(), plan.RefreshTokenLength.ValueInt64(), plan.RefreshRollingInterval.ValueInt64(), plan.RegisteredAuthorizationPath.ValueString(), plan.PendingAuthorizationTimeout.ValueInt64(), plan.DevicePollingInterval.ValueInt64(), plan.BypassActivationCodeConfirmation.ValueBool())
	err := addOptionalOauthAuthServerSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Auth Server Settings", err.Error())
		return
	}

	updateOauthAuthServerSettings = updateOauthAuthServerSettings.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(updateOauthAuthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response
	var state oauthAuthServerSettingsModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readOauthAuthServerSettingsResponse(ctx, updateOauthAuthServerSettingsResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthAuthServerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthAuthServerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
