package oauthauthserversettings

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopeentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthServerSettingsResource{}
	_ resource.ResourceWithImportState = &oauthServerSettingsResource{}

	scopesDefault, _ = types.SetValue(types.ObjectType{AttrTypes: scopeentry.AttrTypes()}, nil)

	scopeGroupsDefault, _ = types.SetValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"scopes":      types.SetType{ElemType: types.StringType},
	}}, nil)
	persistentGrantReuseGrantTypesDefault, _ = types.SetValue(types.StringType, nil)
	allowedOriginsDefault, _                 = types.SetValue(types.StringType, nil)
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

// OauthServerSettingsResource is a helper function to simplify the provider implementation.
func OauthServerSettingsResource() resource.Resource {
	return &oauthServerSettingsResource{}
}

// oauthServerSettingsResource is the resource implementation.
type oauthServerSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *oauthServerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the OAuth authorization server settings.",
		Attributes: map[string]schema.Attribute{
			"default_scope_description": schema.StringAttribute{
				Description: "The default scope description.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
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
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope that appears when the user is prompted for authorization.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"dynamic": schema.BoolAttribute{
							Description: "True if the scope is dynamic. The default is `false`.",
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
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope that appears when the user is prompted for authorization.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"dynamic": schema.BoolAttribute{
							Description: "True if the scope is dynamic. The default is `false`.",
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
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
				Description: "Determines whether PKCE's 'plain' code challenge method will be disallowed. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_issuer_in_authorization_response": schema.BoolAttribute{
				Description: "Determines whether the authorization server's issuer value is added to the authorization response or not. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"track_user_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether user sessions are tracked for logout. The default value is `false`.",
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
			"require_offline_access_scope_to_issue_refresh_tokens": schema.BoolAttribute{
				Description: "Determines whether offline_access scope is required to issue refresh tokens or not. The default value is `false`. Supported in PF version `12.1` or later.",
				Computed:    true,
				Optional:    true,
			},
			"offline_access_require_consent_prompt": schema.BoolAttribute{
				Description: "Determines whether offline_access requires the prompt parameter value be 'consent' or not. The value will be reset to default if the `require_offline_access_scope_to_issue_refresh_tokens` attribute is set to `false`. The default value is `false`. Supported in PF version `12.1` or later.",
				Computed:    true,
				Optional:    true,
			},
			"persistent_grant_lifetime": schema.Int64Attribute{
				Description: "The persistent grant lifetime. The default value is `-1`, which indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(-1),
			},
			"persistent_grant_lifetime_unit": schema.StringAttribute{
				Description: "The persistent grant lifetime unit. Supported values are `MINUTES`, `DAYS`, and `HOURS`. The default value is `DAYS`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. The default value is `30` (days). `-1` indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(30),
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit. Supported values are `MINUTES`, `DAYS`, and `HOURS`. The default value is `DAYS`.",
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
				Description: "The roll refresh token values default policy. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds. The default value is `60`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"refresh_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens.",
				Required:    true,
			},
			"refresh_rolling_interval_time_unit": schema.StringAttribute{
				Description: "The refresh token rolling interval time unit. Supported values are `SECONDS`, `MINUTES`, and `HOURS`. The default value is `HOURS`. Supported in PF version `12.1` or later.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"SECONDS",
						"MINUTES",
						"HOURS",
					),
				},
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. Only `IMPLICIT` or `AUTHORIZATION_CODE` or `RESOURCE_OWNER_CREDENTIALS` are valid grant types.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"IMPLICIT",
							"AUTHORIZATION_CODE",
							"RESOURCE_OWNER_CREDENTIALS"),
					),
				},
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
						Description: "This is a read-only list of persistent grant attributes and includes `USER_KEY` and `USER_NAME`.",
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
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
			"bypass_authorization_for_approved_grants": schema.BoolAttribute{
				Description: "Bypass authorization for previously approved persistent grants. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_unidentified_client_ro_creds": schema.BoolAttribute{
				Description: "Allow unidentified clients to request resource owner password credentials grants. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_unidentified_client_extension_grants": schema.BoolAttribute{
				Description: "Allow unidentified clients to request extension grants. The default value is `false`.",
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
			"allowed_origins": schema.SetAttribute{
				Description: "The list of allowed origins.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(allowedOriginsDefault),
				Validators: []validator.Set{
					configvalidators.ValidUrlsSet(),
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
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					configvalidators.StartsWith("/"),
				},
			},
			"pending_authorization_timeout": schema.Int64Attribute{
				Description: "The 'device_code' and 'user_code' timeout, in seconds. The default is `600` seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(600),
			},
			"device_polling_interval": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds. The default is `5` seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
			},
			"activation_code_check_mode": schema.StringAttribute{
				Description: "Determines whether the user is prompted to enter or confirm the activation code after authenticating or before. Supported values are `AFTER_AUTHENTICATION` and `BEFORE_AUTHENTICATION`. The default value is `AFTER_AUTHENTICATION`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("AFTER_AUTHENTICATION"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"AFTER_AUTHENTICATION", "BEFORE_AUTHENTICATION"}...),
				},
			},
			"bypass_activation_code_confirmation": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device. The default is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_cookieless_user_authorization_authentication_api": schema.BoolAttribute{
				Description: "Indicates if cookies should be used for state tracking when the user authorization endpoint is operating in authentication API redirectless mode. The default is `false`. Supported in PF version `12.1` or later.",
				Optional:    true,
				Computed:    true,
			},
			"user_authorization_consent_page_setting": schema.StringAttribute{
				Description: "User Authorization Consent Page setting to use PingFederate's internal consent page or an external system. Supported values are `INTERNAL` and `ADAPTER`. The default value is `INTERNAL`.",
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"approved_scopes_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved scopes returned by the external consent page.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"approved_authorization_detail_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved authorization details returned by the external consent page.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"par_reference_timeout": schema.Int64Attribute{
				Description: "The timeout, in seconds, of the pushed authorization request reference. The default value is `60`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"par_reference_length": schema.Int64Attribute{
				Description: "The entropy of pushed authorization request references, in bytes. The default value is `24`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(24),
			},
			"par_status": schema.StringAttribute{
				Description: "The status of pushed authorization request support. Supported values are `DISABLED`, `ENABLED`, and `REQUIRED`. The default value is `ENABLED`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("ENABLED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DISABLED", "ENABLED", "REQUIRED"}...),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is `0`, which will disable secondary client secret retention.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"jwt_secured_authorization_response_mode_lifetime": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the JWT Secured authorization response. The default value is `600`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(600),
			},
			"dpop_proof_require_nonce": schema.BoolAttribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "Determines whether nonce is required in the Demonstrating Proof-of-Possession (DPoP) proof JWT. The default value is `false`. Supported in PF version `11.3` or later.",
				Computed:    true,
				Optional:    true,
			},
			"dpop_proof_lifetime_seconds": schema.Int64Attribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "The lifetime, in seconds, of the Demonstrating Proof-of-Possession (DPoP) proof JWT. The default value is `120`. Supported in PF version `11.3` or later.",
				Computed:    true,
				Optional:    true,
			},
			"dpop_proof_enforce_replay_prevention": schema.BoolAttribute{
				// Default is set in ModifyPlan below. Once only PF 11.3 and newer is supported, we can set the default in the schema here
				Description: "Determines whether Demonstrating Proof-of-Possession (DPoP) proof JWT replay prevention is enforced. The default value is `false`. Supported in PF version `11.3` or later.",
				Computed:    true,
				Optional:    true,
			},
			"bypass_authorization_for_approved_consents": schema.BoolAttribute{
				// Default is set in ModifyPlan below. Once only PF 12.0 and newer is supported, we can set the default in the schema here
				Description: "Bypass authorization for previously approved consents. The default value is `false`. Supported in PF version `12.0` or later.",
				Computed:    true,
				Optional:    true,
			},
			"consent_lifetime_days": schema.Int64Attribute{
				// Default is set in ModifyPlan below. Once only PF 12.0 and newer is supported, we can set the default in the schema here
				Description: "The consent lifetime in days. The default value is `-1`, which indicates an indefinite amount of time. Supported in PF version `12.0` or later.",
				Computed:    true,
				Optional:    true,
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func (r *oauthServerSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model oauthServerSettingsModel
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
				if strings.Count(scopeEntryName, "*") != 1 {
					resp.Diagnostics.AddAttributeError(
						path.Root("scopes"),
						providererror.InvalidAttributeConfiguration,
						fmt.Sprintf("Scope name \"%s\" must be include a single \"*\" when dynamic is set to true.", scopeEntryName))
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
					resp.Diagnostics.AddAttributeError(
						path.Root("exclusive_scopes"),
						providererror.InvalidAttributeConfiguration,
						fmt.Sprintf("Scope name \"%s\" must be prefixed with a \"*\" when dynamic is set to true.", eScopeEntryName))
				}
			}
		}
	}

	// Test if values in sets match
	matchVal := internaltypes.MatchStringInSets(scopeNames, eScopeNames)
	if matchVal != nil {
		resp.Diagnostics.AddError(
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("The scope name \"%s\" is defined in both scopes and exclusive_scopes", *matchVal))
	}

	// offline_access_require_consent_prompt can't be true if require_offline_access_scope_to_issue_refresh_tokens is false
	if model.OfflineAccessRequireConsentPrompt.ValueBool() && !model.RequireOfflineAccessScopeToIssueRefreshTokens.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("require_offline_access_scope_to_issue_refresh_tokens"),
			providererror.InvalidAttributeConfiguration,
			"require_offline_access_scope_to_issue_refresh_tokens must be set to true to set offline_access_require_consent_prompt to true")
	}
}

func (r *oauthServerSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to versions 11.3, 12.0, and 12.1 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast113 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast121 := compare >= 0
	var plan *oauthServerSettingsModel
	req.Plan.Get(ctx, &plan)
	if plan == nil {
		return
	}
	// If any of these fields are set by the user and the PF version is not new enough, throw an error
	if !pfVersionAtLeast113 {
		if internaltypes.IsDefined(plan.DpopProofEnforceReplayPrevention) {
			version.AddUnsupportedAttributeError("dpop_proof_enforce_replay_prevention",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
		} else if plan.DpopProofEnforceReplayPrevention.IsUnknown() {
			// Set a null default when the version isn't new enough to use this attribute
			plan.DpopProofEnforceReplayPrevention = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.DpopProofLifetimeSeconds) {
			version.AddUnsupportedAttributeError("dpop_proof_lifetime_seconds",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
		} else if plan.DpopProofLifetimeSeconds.IsUnknown() {
			plan.DpopProofLifetimeSeconds = types.Int64Null()
		}

		if internaltypes.IsDefined(plan.DpopProofRequireNonce) {
			version.AddUnsupportedAttributeError("dpop_proof_require_nonce",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
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

	// Similar logic for PF 12.0
	if !pfVersionAtLeast120 {
		if internaltypes.IsDefined(plan.BypassAuthorizationForApprovedConsents) {
			version.AddUnsupportedAttributeError("bypass_authorization_for_approved_consents",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		} else if plan.BypassAuthorizationForApprovedConsents.IsUnknown() {
			plan.BypassAuthorizationForApprovedConsents = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.ConsentLifetimeDays) {
			version.AddUnsupportedAttributeError("consent_lifetime_days",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		} else if plan.ConsentLifetimeDays.IsUnknown() {
			plan.ConsentLifetimeDays = types.Int64Null()
		}
	} else {
		if plan.BypassAuthorizationForApprovedConsents.IsUnknown() {
			plan.BypassAuthorizationForApprovedConsents = types.BoolValue(false)
		}
		if plan.ConsentLifetimeDays.IsUnknown() {
			plan.ConsentLifetimeDays = types.Int64Value(-1)
		}
	}

	// Similar logic for PF 12.1
	if !pfVersionAtLeast121 {
		if internaltypes.IsDefined(plan.RequireOfflineAccessScopeToIssueRefreshTokens) {
			version.AddUnsupportedAttributeError("require_offline_access_scope_to_issue_refresh_tokens",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else if plan.RequireOfflineAccessScopeToIssueRefreshTokens.IsUnknown() {
			plan.RequireOfflineAccessScopeToIssueRefreshTokens = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.OfflineAccessRequireConsentPrompt) {
			version.AddUnsupportedAttributeError("offline_access_require_consent_prompt",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else if plan.OfflineAccessRequireConsentPrompt.IsUnknown() {
			plan.OfflineAccessRequireConsentPrompt = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.RefreshRollingIntervalTimeUnit) {
			version.AddUnsupportedAttributeError("refresh_rolling_interval_time_unit",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else if plan.RefreshRollingIntervalTimeUnit.IsUnknown() {
			plan.RefreshRollingIntervalTimeUnit = types.StringNull()
		}

		if internaltypes.IsDefined(plan.EnableCookielessUserAuthorizationAuthenticationApi) {
			version.AddUnsupportedAttributeError("enable_cookieless_user_authorization_authentication_api",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else if plan.EnableCookielessUserAuthorizationAuthenticationApi.IsUnknown() {
			plan.EnableCookielessUserAuthorizationAuthenticationApi = types.BoolNull()
		}
	} else {
		if plan.RequireOfflineAccessScopeToIssueRefreshTokens.IsUnknown() {
			plan.RequireOfflineAccessScopeToIssueRefreshTokens = types.BoolValue(false)
		}
		if plan.OfflineAccessRequireConsentPrompt.IsUnknown() {
			plan.OfflineAccessRequireConsentPrompt = types.BoolValue(false)
		}
		if plan.RefreshRollingIntervalTimeUnit.IsUnknown() {
			plan.RefreshRollingIntervalTimeUnit = types.StringValue("HOURS")
		}
		if plan.EnableCookielessUserAuthorizationAuthenticationApi.IsUnknown() {
			plan.EnableCookielessUserAuthorizationAuthenticationApi = types.BoolValue(false)
		}
	}

	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

func addOptionalOauthServerSettingsFields(ctx context.Context, addRequest *client.AuthorizationServerSettings, plan oauthServerSettingsModel) error {

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
	addRequest.RegisteredAuthorizationPath = plan.RegisteredAuthorizationPath.ValueStringPointer()
	addRequest.BypassActivationCodeConfirmation = plan.BypassActivationCodeConfirmation.ValueBoolPointer()
	addRequest.DefaultScopeDescription = plan.DefaultScopeDescription.ValueStringPointer()
	addRequest.DevicePollingInterval = plan.DevicePollingInterval.ValueInt64Pointer()
	addRequest.PendingAuthorizationTimeout = plan.PendingAuthorizationTimeout.ValueInt64Pointer()
	addRequest.DisallowPlainPKCE = plan.DisallowPlainPKCE.ValueBoolPointer()
	addRequest.IncludeIssuerInAuthorizationResponse = plan.IncludeIssuerInAuthorizationResponse.ValueBoolPointer()
	addRequest.TrackUserSessionsForLogout = plan.TrackUserSessionsForLogout.ValueBoolPointer()
	addRequest.TokenEndpointBaseUrl = plan.TokenEndpointBaseUrl.ValueStringPointer()
	addRequest.RequireOfflineAccessScopeToIssueRefreshTokens = plan.RequireOfflineAccessScopeToIssueRefreshTokens.ValueBoolPointer()
	addRequest.OfflineAccessRequireConsentPrompt = plan.OfflineAccessRequireConsentPrompt.ValueBoolPointer()
	addRequest.PersistentGrantLifetime = plan.PersistentGrantLifetime.ValueInt64Pointer()
	addRequest.PersistentGrantLifetimeUnit = plan.PersistentGrantLifetimeUnit.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeout = plan.PersistentGrantIdleTimeout.ValueInt64Pointer()
	addRequest.PersistentGrantIdleTimeoutTimeUnit = plan.PersistentGrantIdleTimeoutTimeUnit.ValueStringPointer()
	addRequest.RollRefreshTokenValues = plan.RollRefreshTokenValues.ValueBoolPointer()
	addRequest.RefreshTokenRollingGracePeriod = plan.RefreshTokenRollingGracePeriod.ValueInt64Pointer()
	addRequest.RefreshRollingIntervalTimeUnit = plan.RefreshRollingIntervalTimeUnit.ValueStringPointer()
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
	addRequest.EnableCookielessUserAuthorizationAuthenticationApi = plan.EnableCookielessUserAuthorizationAuthenticationApi.ValueBoolPointer()
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
	addRequest.DevicePollingInterval = plan.DevicePollingInterval.ValueInt64Pointer()
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
	addRequest.BypassAuthorizationForApprovedConsents = plan.BypassAuthorizationForApprovedConsents.ValueBoolPointer()
	addRequest.ConsentLifetimeDays = plan.ConsentLifetimeDays.ValueInt64Pointer()

	return nil

}

// Metadata returns the resource type name.
func (r *oauthServerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_server_settings"
}

func (r *oauthServerSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *oauthServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthServerSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthServerSettings := client.NewAuthorizationServerSettings(plan.AuthorizationCodeTimeout.ValueInt64(), plan.AuthorizationCodeEntropy.ValueInt64(), plan.RefreshTokenLength.ValueInt64(), plan.RefreshRollingInterval.ValueInt64())
	err := addOptionalOauthServerSettingsFields(ctx, createOauthServerSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OAuth Auth Server Settings: "+err.Error())
		return
	}

	apiCreateOauthServerSettings := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthServerSettings = apiCreateOauthServerSettings.Body(*createOauthServerSettings)
	oauthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(apiCreateOauthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthServerSettingsModel
	diags = readOauthServerSettingsResponse(ctx, oauthServerSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthServerSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthServerSettings, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetAuthorizationServerSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Auth Server Settings", httpResp)
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
	diags = readOauthServerSettingsResponse(ctx, apiReadOauthServerSettings, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthServerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthServerSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	updateOauthServerSettings := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthorizationServerSettings(plan.AuthorizationCodeTimeout.ValueInt64(), plan.AuthorizationCodeEntropy.ValueInt64(), plan.RefreshTokenLength.ValueInt64(), plan.RefreshRollingInterval.ValueInt64())
	err := addOptionalOauthServerSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OAuth Auth Server Settings: "+err.Error())
		return
	}

	updateOauthServerSettings = updateOauthServerSettings.Body(*createUpdateRequest)
	updateOauthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(updateOauthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response
	var state oauthServerSettingsModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readOauthServerSettingsResponse(ctx, updateOauthServerSettingsResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthServerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthServerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
