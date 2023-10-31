package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopeentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopegroupentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsResource{}
)

var (
	nameAttributeType = map[string]attr.Type{
		"name": basetypes.StringType{},
	}
	persistentGrantObjContractTypes = map[string]attr.Type{
		"core_attributes":     basetypes.SetType{ElemType: types.ObjectType{AttrTypes: nameAttributeType}},
		"extended_attributes": basetypes.SetType{ElemType: types.ObjectType{AttrTypes: nameAttributeType}},
	}
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

type oauthAuthServerSettingsResourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	DefaultScopeDescription                     types.String `tfsdk:"default_scope_description"`
	Scopes                                      types.Set    `tfsdk:"scopes"`
	ScopeGroups                                 types.Set    `tfsdk:"scope_groups"`
	ExclusiveScopes                             types.Set    `tfsdk:"exclusive_scopes"`
	ExclusiveScopeGroups                        types.Set    `tfsdk:"exclusive_scope_groups"`
	AuthorizationCodeTimeout                    types.Int64  `tfsdk:"authorization_code_timeout"`
	AuthorizationCodeEntropy                    types.Int64  `tfsdk:"authorization_code_entropy"`
	DisallowPlainPKCE                           types.Bool   `tfsdk:"disallow_plain_pkce"`
	IncludeIssuerInAuthorizationResponse        types.Bool   `tfsdk:"include_issuer_in_authorization_response"`
	TrackUserSessionsForLogout                  types.Bool   `tfsdk:"track_user_sessions_for_logout"`
	TokenEndpointBaseUrl                        types.String `tfsdk:"token_endpoint_base_url"`
	PersistentGrantLifetime                     types.Int64  `tfsdk:"persistent_grant_lifetime"`
	PersistentGrantLifetimeUnit                 types.String `tfsdk:"persistent_grant_lifetime_unit"`
	PersistentGrantIdleTimeout                  types.Int64  `tfsdk:"persistent_grant_idle_timeout"`
	PersistentGrantIdleTimeoutTimeUnit          types.String `tfsdk:"persistent_grant_idle_timeout_time_unit"`
	RefreshTokenLength                          types.Int64  `tfsdk:"refresh_token_length"`
	RollRefreshTokenValues                      types.Bool   `tfsdk:"roll_refresh_token_values"`
	RefreshTokenRollingGracePeriod              types.Int64  `tfsdk:"refresh_token_rolling_grace_period"`
	RefreshRollingInterval                      types.Int64  `tfsdk:"refresh_rolling_interval"`
	PersistentGrantReuseGrantTypes              types.Set    `tfsdk:"persistent_grant_reuse_grant_types"`
	PersistentGrantContract                     types.Object `tfsdk:"persistent_grant_contract"`
	BypassAuthorizationForApprovedGrants        types.Bool   `tfsdk:"bypass_authorization_for_approved_grants"`
	AllowUnidentifiedClientROCreds              types.Bool   `tfsdk:"allow_unidentified_client_ro_creds"`
	AllowUnidentifiedClientExtensionGrants      types.Bool   `tfsdk:"allow_unidentified_client_extension_grants"`
	AdminWebServicePcvRef                       types.Object `tfsdk:"admin_web_service_pcv_ref"`
	AtmIdForOAuthGrantManagement                types.String `tfsdk:"atm_id_for_oauth_grant_management"`
	ScopeForOAuthGrantManagement                types.String `tfsdk:"scope_for_oauth_grant_management"`
	AllowedOrigins                              types.List   `tfsdk:"allowed_origins"`
	UserAuthorizationUrl                        types.String `tfsdk:"user_authorization_url"`
	BypassActivationCodeConfirmation            types.Bool   `tfsdk:"bypass_activation_code_confirmation"`
	RegisteredAuthorizationPath                 types.String `tfsdk:"registered_authorization_path"`
	PendingAuthorizationTimeout                 types.Int64  `tfsdk:"pending_authorization_timeout"`
	DevicePollingInterval                       types.Int64  `tfsdk:"device_polling_interval"`
	ActivationCodeCheckMode                     types.String `tfsdk:"activation_code_check_mode"`
	UserAuthorizationConsentPageSetting         types.String `tfsdk:"user_authorization_consent_page_setting"`
	UserAuthorizationConsentAdapter             types.String `tfsdk:"user_authorization_consent_adapter"`
	ApprovedScopesAttribute                     types.String `tfsdk:"approved_scopes_attribute"`
	ApprovedAuthorizationDetailAttribute        types.String `tfsdk:"approved_authorization_detail_attribute"`
	ParReferenceTimeout                         types.Int64  `tfsdk:"par_reference_timeout"`
	ParReferenceLength                          types.Int64  `tfsdk:"par_reference_length"`
	ParStatus                                   types.String `tfsdk:"par_status"`
	ClientSecretRetentionPeriod                 types.Int64  `tfsdk:"client_secret_retention_period"`
	JwtSecuredAuthorizationResponseModeLifetime types.Int64  `tfsdk:"jwt_secured_authorization_response_mode_lifetime"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages OAuth Auth Server Settings",
		Attributes: map[string]schema.Attribute{
			"default_scope_description": schema.StringAttribute{
				Description: "The default scope description.",
				Required:    true,
			},
			"scopes": schema.SetNestedAttribute{
				Description: "The list of common scopes.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"scope_groups": schema.SetNestedAttribute{
				Description: "The list of common scope groups.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"exclusive_scope_groups": schema.SetNestedAttribute{
				Description: "The list of exclusive scope groups.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_issuer_in_authorization_response": schema.BoolAttribute{
				Description: "Determines whether the authorization server's issuer value is added to the authorization response or not. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"track_user_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether user sessions are tracked for logout. If this property is not provided on a PUT, the setting is left unchanged.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"token_endpoint_base_url": schema.StringAttribute{
				Description: "The token endpoint base URL used to validate the 'aud' claim during Private Key JWT Client Authentication.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"persistent_grant_lifetime": schema.Int64Attribute{
				Description: "The persistent grant lifetime. The default value is indefinite. -1 indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(-1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"persistent_grant_lifetime_unit": schema.StringAttribute{
				Description: "The persistent grant lifetime unit.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. The default value is 30 (days). -1 indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(30),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit. The default value is DAYS",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds. The default value is 0.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"refresh_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens, in hours.",
				Required:    true,
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"persistent_grant_contract": schema.SingleNestedAttribute{
				Description: "The persistent grant contract defines attributes that are associated with OAuth persistent grants.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "This is a read-only list of persistent grant attributes and includes USER_KEY and USER_NAME. Changes to this field will be ignored.",
						Computed:    true,
						Optional:    false,
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
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_unidentified_client_ro_creds": schema.BoolAttribute{
				Description: "Allow unidentified clients to request resource owner password credentials grants. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_unidentified_client_extension_grants": schema.BoolAttribute{
				Description: "Allow unidentified clients to request extension grants. The default value is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_for_oauth_grant_management": schema.StringAttribute{
				Description: "The OAuth scope to validate when accessing grant management service.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_origins": schema.ListAttribute{
				Description: "The list of allowed origins.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					configvalidators.ValidUrls(),
				},
			},
			"user_authorization_url": schema.StringAttribute{
				Description: "The URL used to generate 'verification_url' and 'verification_url_complete' values in a Device Authorization request",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"INTERNAL", "ADAPTER"}...),
				},
			},
			"user_authorization_consent_adapter": schema.StringAttribute{
				Description: "Adapter ID of the external consent adapter to be used for the consent page user interface.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"approved_scopes_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved scopes returned by the external consent page.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"approved_authorization_detail_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved authorization details returned by the external consent page.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"par_reference_timeout": schema.Int64Attribute{
				Description: "The timeout, in seconds, of the pushed authorization request reference. The default value is 60.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"par_reference_length": schema.Int64Attribute{
				Description: "The entropy of pushed authorization request references, in bytes. The default value is 24.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(24),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"par_status": schema.StringAttribute{
				Description: "The status of pushed authorization request support. The default value is ENABLED.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("ENABLED"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DISABLED", "ENABLED", "REQUIRED"}...),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"jwt_secured_authorization_response_mode_lifetime": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the JWT Secured authorization response. The default value is 600.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(600),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	var model oauthAuthServerSettingsResourceModel
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

func addOptionalOauthAuthServerSettingsFields(ctx context.Context, addRequest *client.AuthorizationServerSettings, plan oauthAuthServerSettingsResourceModel) error {

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

	if internaltypes.IsDefined(plan.DisallowPlainPKCE) {
		addRequest.DisallowPlainPKCE = plan.DisallowPlainPKCE.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.IncludeIssuerInAuthorizationResponse) {
		addRequest.IncludeIssuerInAuthorizationResponse = plan.IncludeIssuerInAuthorizationResponse.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.TrackUserSessionsForLogout) {
		addRequest.TrackUserSessionsForLogout = plan.TrackUserSessionsForLogout.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.TokenEndpointBaseUrl) {
		addRequest.TokenEndpointBaseUrl = plan.TokenEndpointBaseUrl.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.PersistentGrantLifetime) {
		addRequest.PersistentGrantLifetime = plan.PersistentGrantLifetime.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.PersistentGrantLifetimeUnit) {
		addRequest.PersistentGrantLifetimeUnit = plan.PersistentGrantLifetimeUnit.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.PersistentGrantIdleTimeout) {
		addRequest.PersistentGrantIdleTimeout = plan.PersistentGrantIdleTimeout.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.PersistentGrantIdleTimeoutTimeUnit) {
		addRequest.PersistentGrantIdleTimeoutTimeUnit = plan.PersistentGrantIdleTimeoutTimeUnit.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.RollRefreshTokenValues) {
		addRequest.RollRefreshTokenValues = plan.RollRefreshTokenValues.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.RefreshTokenRollingGracePeriod) {
		addRequest.RefreshTokenRollingGracePeriod = plan.RefreshTokenRollingGracePeriod.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.PersistentGrantReuseGrantTypes) {
		var slice []string
		plan.PersistentGrantReuseGrantTypes.ElementsAs(ctx, &slice, false)
		addRequest.PersistentGrantReuseGrantTypes = slice
	}

	if internaltypes.IsDefined(plan.PersistentGrantContract) {
		addRequest.PersistentGrantContract = client.NewPersistentGrantContractWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.PersistentGrantContract, false)), addRequest.PersistentGrantContract)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.BypassAuthorizationForApprovedGrants) {
		addRequest.BypassAuthorizationForApprovedGrants = plan.BypassAuthorizationForApprovedGrants.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.AllowUnidentifiedClientROCreds) {
		addRequest.AllowUnidentifiedClientROCreds = plan.AllowUnidentifiedClientROCreds.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.AllowUnidentifiedClientExtensionGrants) {
		addRequest.AllowUnidentifiedClientExtensionGrants = plan.AllowUnidentifiedClientExtensionGrants.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.AdminWebServicePcvRef) {
		addRequest.AdminWebServicePcvRef = client.NewResourceLinkWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AdminWebServicePcvRef, false)), addRequest.AdminWebServicePcvRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AtmIdForOAuthGrantManagement) {
		addRequest.AtmIdForOAuthGrantManagement = plan.AtmIdForOAuthGrantManagement.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ScopeForOAuthGrantManagement) {
		addRequest.ScopeForOAuthGrantManagement = plan.ScopeForOAuthGrantManagement.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.AllowedOrigins) {
		var slice []string
		plan.AllowedOrigins.ElementsAs(ctx, &slice, false)
		addRequest.AllowedOrigins = slice
	}

	if internaltypes.IsDefined(plan.UserAuthorizationUrl) {
		addRequest.UserAuthorizationUrl = plan.UserAuthorizationUrl.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.DevicePollingInterval) {
		addRequest.DevicePollingInterval = plan.DevicePollingInterval.ValueInt64()
	}

	if internaltypes.IsDefined(plan.ActivationCodeCheckMode) {
		addRequest.ActivationCodeCheckMode = plan.ActivationCodeCheckMode.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.UserAuthorizationConsentPageSetting) {
		addRequest.UserAuthorizationConsentPageSetting = plan.UserAuthorizationConsentPageSetting.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.UserAuthorizationConsentAdapter) {
		addRequest.UserAuthorizationConsentAdapter = plan.UserAuthorizationConsentAdapter.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ApprovedScopesAttribute) {
		addRequest.ApprovedScopesAttribute = plan.ApprovedScopesAttribute.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ApprovedAuthorizationDetailAttribute) {
		addRequest.ApprovedAuthorizationDetailAttribute = plan.ApprovedAuthorizationDetailAttribute.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ParReferenceTimeout) {
		addRequest.ParReferenceTimeout = plan.ParReferenceTimeout.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.ParReferenceLength) {
		addRequest.ParReferenceLength = plan.ParReferenceLength.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.ParStatus) {
		addRequest.ParStatus = plan.ParStatus.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ClientSecretRetentionPeriod) {
		addRequest.ClientSecretRetentionPeriod = plan.ClientSecretRetentionPeriod.ValueInt64Pointer()
	}

	if internaltypes.IsDefined(plan.JwtSecuredAuthorizationResponseModeLifetime) {
		addRequest.JwtSecuredAuthorizationResponseModeLifetime = plan.JwtSecuredAuthorizationResponseModeLifetime.ValueInt64Pointer()
	}

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

func readOauthAuthServerSettingsResponse(ctx context.Context, r *client.AuthorizationServerSettings, state *oauthAuthServerSettingsResourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
	state.DefaultScopeDescription = types.StringValue(r.DefaultScopeDescription)
	state.Scopes, respDiags = scopeentry.ToState(ctx, r.Scopes)
	diags.Append(respDiags...)
	state.ScopeGroups, respDiags = scopegroupentry.ToState(ctx, r.ScopeGroups)
	diags.Append(respDiags...)
	state.ExclusiveScopes, respDiags = scopeentry.ToState(ctx, r.ExclusiveScopes)
	diags.Append(respDiags...)
	state.ExclusiveScopeGroups, respDiags = scopegroupentry.ToState(ctx, r.ExclusiveScopeGroups)
	diags.Append(respDiags...)
	persistentGrantContract, respDiags := types.ObjectValueFrom(ctx, persistentGrantObjContractTypes, r.PersistentGrantContract)
	diags.Append(respDiags...)

	state.PersistentGrantContract = persistentGrantContract
	state.AuthorizationCodeTimeout = types.Int64Value(r.AuthorizationCodeTimeout)
	state.AuthorizationCodeEntropy = types.Int64Value(r.AuthorizationCodeEntropy)
	state.DisallowPlainPKCE = types.BoolPointerValue(r.DisallowPlainPKCE)
	state.IncludeIssuerInAuthorizationResponse = types.BoolPointerValue(r.IncludeIssuerInAuthorizationResponse)
	state.TrackUserSessionsForLogout = types.BoolPointerValue(r.TrackUserSessionsForLogout)
	state.TokenEndpointBaseUrl = types.StringPointerValue(r.TokenEndpointBaseUrl)
	state.PersistentGrantLifetime = types.Int64PointerValue(r.PersistentGrantLifetime)
	state.PersistentGrantLifetimeUnit = types.StringPointerValue(r.PersistentGrantLifetimeUnit)
	state.PersistentGrantIdleTimeout = types.Int64PointerValue(r.PersistentGrantIdleTimeout)
	state.PersistentGrantIdleTimeoutTimeUnit = types.StringPointerValue(r.PersistentGrantIdleTimeoutTimeUnit)
	state.RefreshTokenLength = types.Int64Value(r.RefreshTokenLength)
	state.RollRefreshTokenValues = types.BoolPointerValue(r.RollRefreshTokenValues)
	state.RefreshTokenRollingGracePeriod = types.Int64PointerValue(r.RefreshTokenRollingGracePeriod)
	state.RefreshRollingInterval = types.Int64Value(r.RefreshRollingInterval)
	state.PersistentGrantReuseGrantTypes = internaltypes.GetStringSet(r.PersistentGrantReuseGrantTypes)
	state.BypassAuthorizationForApprovedGrants = types.BoolPointerValue(r.BypassAuthorizationForApprovedGrants)
	state.AllowUnidentifiedClientROCreds = types.BoolPointerValue(r.AllowUnidentifiedClientROCreds)
	state.AllowUnidentifiedClientExtensionGrants = types.BoolPointerValue(r.AllowUnidentifiedClientExtensionGrants)
	state.AdminWebServicePcvRef, respDiags = resourcelink.ToState(ctx, r.AdminWebServicePcvRef)
	diags.Append(respDiags...)
	state.AtmIdForOAuthGrantManagement = types.StringPointerValue(r.AtmIdForOAuthGrantManagement)
	state.ScopeForOAuthGrantManagement = types.StringPointerValue(r.ScopeForOAuthGrantManagement)
	state.AllowedOrigins = internaltypes.GetStringList(r.AllowedOrigins)
	state.UserAuthorizationUrl = types.StringPointerValue(r.UserAuthorizationUrl)
	state.RegisteredAuthorizationPath = types.StringValue(r.RegisteredAuthorizationPath)
	state.PendingAuthorizationTimeout = types.Int64Value(r.PendingAuthorizationTimeout)
	state.DevicePollingInterval = types.Int64Value(r.DevicePollingInterval)
	state.ActivationCodeCheckMode = types.StringPointerValue(r.ActivationCodeCheckMode)
	state.BypassActivationCodeConfirmation = types.BoolValue(r.BypassActivationCodeConfirmation)
	state.UserAuthorizationConsentPageSetting = types.StringPointerValue(r.UserAuthorizationConsentPageSetting)
	state.UserAuthorizationConsentAdapter = types.StringPointerValue(r.UserAuthorizationConsentAdapter)
	state.ApprovedScopesAttribute = types.StringPointerValue(r.ApprovedScopesAttribute)
	state.ApprovedAuthorizationDetailAttribute = types.StringPointerValue(r.ApprovedAuthorizationDetailAttribute)
	state.ParReferenceTimeout = types.Int64PointerValue(r.ParReferenceTimeout)
	state.ParReferenceLength = types.Int64PointerValue(r.ParReferenceLength)
	state.ParStatus = types.StringPointerValue(r.ParStatus)
	state.ClientSecretRetentionPeriod = types.Int64PointerValue(r.ClientSecretRetentionPeriod)
	state.JwtSecuredAuthorizationResponseModeLifetime = types.Int64PointerValue(r.JwtSecuredAuthorizationResponseModeLifetime)
	return diags
}

func (r *oauthAuthServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsResourceModel

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
	_, requestErr := createOauthAuthServerSettings.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of OAuth Auth Server Settings: %s", requestErr.Error())
	}

	apiCreateOauthAuthServerSettings := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettings = apiCreateOauthAuthServerSettings.Body(*createOauthAuthServerSettings)
	oauthAuthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(apiCreateOauthAuthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Auth Server Settings", err, httpResp)
		return
	}
	_, responseErr := oauthAuthServerSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of OAuth Auth Server Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	var state oauthAuthServerSettingsResourceModel
	diags = readOauthAuthServerSettingsResponse(ctx, oauthAuthServerSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAuthServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAuthServerSettingsResourceModel

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
	// Log response JSON
	_, responseErr := apiReadOauthAuthServerSettings.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of OAuth Auth Server Settings: %s", responseErr.Error())
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
	var plan oauthAuthServerSettingsResourceModel
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
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of a OAuth Auth Server Settings: %s", requestErr.Error())
	}
	updateOauthAuthServerSettings = updateOauthAuthServerSettings.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.UpdateAuthorizationServerSettingsExecute(updateOauthAuthServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Auth Server Settings", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateOauthAuthServerSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a OAuth Auth Server Settings: %s", responseErr.Error())
	}
	// Read the response
	var state oauthAuthServerSettingsResourceModel
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
