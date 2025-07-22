// Copyright © 2025 Ping Identity Corporation

package oauthclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthClientResource{}
	_ resource.ResourceWithConfigure   = &oauthClientResource{}
	_ resource.ResourceWithImportState = &oauthClientResource{}
)

var (
	emptyStringSet, _            = types.SetValue(types.StringType, []attr.Value{})
	oidcPolicyDefaultObj, _      = types.ObjectValue(oidcPolicyAttrType, oidcPolicyDefaultAttrValue)
	secondarySecretsEmptyList, _ = types.ListValue(types.ObjectType{AttrTypes: secondarySecretsAttrType}, []attr.Value{})
	clientAuthDefaultObj, _      = types.ObjectValue(clientAuthAttrType, clientAuthDefaultAttrValue)

	customId = "client_id"
)

// OauthClientResource is a helper function to simplify the provider implementation.
func OauthClientResource() resource.Resource {
	return &oauthClientResource{}
}

// oauthClientResource is the resource implementation.
type oauthClientResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *oauthClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Oauth Client",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "A unique identifier the client provides to the Resource Server to identify itself. This identifier is included with every request the client makes. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Specifies whether the client is enabled. The default value is `true`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"redirect_uris": schema.SetAttribute{
				Description: "URIs to which the OAuth AS may redirect the resource owner's user agent after authorization is obtained. A redirection URI is used with the Authorization Code and Implicit grant types. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible. For example: https://.company.com/ Important: If more than one URI is added or if a single URI uses wildcards, then Authorization Code grant and token requests must contain a specific matching redirect uri parameter.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(emptyStringSet),
			},
			"grant_types": schema.SetAttribute{
				Description: "The grant types allowed for this client. The `EXTENSION` grant type applies to SAML/JWT assertion grants. Supported values are `IMPLICIT`, `AUTHORIZATION_CODE`, `RESOURCE_OWNER_CREDENTIALS`, `CLIENT_CREDENTIALS`, `REFRESH_TOKEN`, `EXTENSION`, `DEVICE_CODE`, `ACCESS_TOKEN_VALIDATION`, `CIBA`, and `TOKEN_EXCHANGE`.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("IMPLICIT",
							"AUTHORIZATION_CODE",
							"RESOURCE_OWNER_CREDENTIALS",
							"CLIENT_CREDENTIALS",
							"REFRESH_TOKEN",
							"EXTENSION",
							"DEVICE_CODE",
							"ACCESS_TOKEN_VALIDATION",
							"CIBA",
							"TOKEN_EXCHANGE",
						),
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "A descriptive name for the client instance. This name appears when the user is prompted for authorization.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description of what the client application does. This description appears when the user is prompted for authorization.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"modification_date": schema.StringAttribute{
				Description: "The time at which the client was last changed. This property is read only.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
				Description: "The time at which the client was created. This property is read only.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"logo_url": schema.StringAttribute{
				Description: "The location of the logo used on user-facing OAuth grant authorization and revocation pages.",
				Optional:    true,
				Validators: []validator.String{
					configvalidators.ValidUrl(),
				},
			},
			"default_access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The default access token manager for this client.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"restrict_to_default_access_token_manager": schema.BoolAttribute{
				Description: "Determines whether the client is restricted to using only its default access token manager. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"validate_using_all_eligible_atms": schema.BoolAttribute{
				Description: "Validates token using all eligible access token managers for the client. This setting is ignored if 'restrict_to_default_access_token_manager' is set to `true`. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"refresh_rolling": schema.StringAttribute{
				Description: "Use `ROLL` or `DONT_ROLL` to override the Roll Refresh Token Values setting on the Authorization Server Settings. `SERVER_DEFAULT` will default to the Roll Refresh Token Values setting on the Authorization Server Setting screen. Defaults to `SERVER_DEFAULT`. Supported values are `ROLL`, `DONT_ROLL`, and `SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("ROLL", "DONT_ROLL", "SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_interval_type": schema.StringAttribute{
				Description: "Use `OVERRIDE_SERVER_DEFAULT` to override the Refresh Token Rolling Interval value on the Authorization Server Settings. `SERVER_DEFAULT` will default to the Refresh Token Rolling Interval value on the Authorization Server Setting. Defaults to `SERVER_DEFAULT`. Supported values are `OVERRIDE_SERVER_DEFAULT` and `SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("OVERRIDE_SERVER_DEFAULT", "SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens. This value will override the Refresh Token Rolling Interval Value on the Authorization Server Settings.",
				Optional:    true,
			},
			"refresh_token_rolling_interval_time_unit": schema.StringAttribute{
				Description: "The refresh token rolling interval time unit. Defaults to `HOURS`. Supported values are `MINUTES`, `HOURS`, and `DAYS`. Supported in PF version `12.1` or later.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"persistent_grant_expiration_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Lifetime set globally for the OAuth AS. Defaults to `SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"persistent_grant_expiration_time": schema.Int64Attribute{
				Description: "The persistent grant expiration time. `-1` indicates an indefinite amount of time. Defaults to `0`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"persistent_grant_expiration_time_unit": schema.StringAttribute{
				Description: "The persistent grant expiration time unit. Defaults to `DAYS`. Supported values are `MINUTES`, `HOURS`, and `DAYS`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"persistent_grant_idle_timeout_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Idle Timeout set globally for the OAuth AS. Defaults to `SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. Defaults to `0`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit. Defaults to `DAYS`. Supported values are `MINUTES`, `HOURS`, and `DAYS`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"persistent_grant_reuse_type": schema.StringAttribute{
				Description: "Allows and administrator to override the Reuse Existing Persistent Access Grants for Grant Types set globally for OAuth AS. Defaults to `SERVER_DEFAULT`. Supported values are `SERVER_DEFAULT` and `OVERRIDE_SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("SERVER_DEFAULT", "OVERRIDE_SERVER_DEFAULT"),
				},
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. This value will override the Reuse Existing Persistent Access Grants for Grant Types on the Authorization Server Settings. Only `IMPLICIT` or `AUTHORIZATION_CODE` or `RESOURCE_OWNER_CREDENTIALS` are valid grant types. Supported values are `IMPLICIT`, `AUTHORIZATION_CODE`, `RESOURCE_OWNER_CREDENTIALS`, `CLIENT_CREDENTIALS`, `REFRESH_TOKEN`, `EXTENSION`, `DEVICE_CODE`, `ACCESS_TOKEN_VALIDATION`, `CIBA`, and `TOKEN_EXCHANGE`.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("IMPLICIT",
							"AUTHORIZATION_CODE",
							"RESOURCE_OWNER_CREDENTIALS",
							"CLIENT_CREDENTIALS",
							"REFRESH_TOKEN",
							"EXTENSION",
							"DEVICE_CODE",
							"ACCESS_TOKEN_VALIDATION",
							"CIBA",
							"TOKEN_EXCHANGE",
						),
					),
				},
			},
			"allow_authentication_api_init": schema.BoolAttribute{
				Description: "Set to `true` to allow this client to initiate the authentication API redirectless flow. Defaults to `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_cookieless_authentication_api": schema.BoolAttribute{
				Description: "Set to `true` to allow the authentication API redirectless flow to function without requiring any cookies. Defaults to `false`. Supported in PF version `12.1` or later.",
				Optional:    true,
				Computed:    true,
			},
			"bypass_approval_page": schema.BoolAttribute{
				Description: "Use this setting, for example, when you want to deploy a trusted application and authenticate end users via an IdP adapter or IdP connection. Defaults to `true` if `allow_authentication_api_init` is `true`, otherwise `false`.",
				Computed:    true,
				Optional:    true,
			},
			"restrict_scopes": schema.BoolAttribute{
				Description: "Restricts this client's access to specific scopes. Defaults to `true` if `allow_authentication_api_init` is `true`, otherwise `false`.",
				Computed:    true,
				Optional:    true,
			},
			"restricted_scopes": schema.SetAttribute{
				Description: "The scopes available for this client.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
			},
			"exclusive_scopes": schema.SetAttribute{
				Description: "The exclusive scopes available for this client.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
			},
			"authorization_detail_types": schema.SetAttribute{
				Description: "The authorization detail types available for this client.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
			},
			"restricted_response_types": schema.SetAttribute{
				Description: "The response types allowed for this client. If omitted all response types are available to the client.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(emptyStringSet),
			},
			"require_pushed_authorization_requests": schema.BoolAttribute{
				Description: "Determines whether pushed authorization requests are required when initiating an authorization request. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_jwt_secured_authorization_response_mode": schema.BoolAttribute{
				Description: "Determines whether JWT secured authorization response mode is required when initiating an authorization request. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether JWT Secured authorization response mode is required when initiating an authorization request. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present\n`RS256` - RSA using SHA-256\n\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256` - ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present, `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RS256",
						"RS384",
						"RS512",
						"ES256",
						"ES384",
						"ES512",
						"PS256",
						"PS384",
						"PS512",
					),
				},
			},
			"oidc_policy": schema.SingleNestedAttribute{
				Description: "Open ID Connect Policy settings. This is included in the message only when OIDC is enabled.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(oidcPolicyDefaultObj),
				Attributes: map[string]schema.Attribute{
					"id_token_signing_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm required for the ID Token.\n`NONE` - No signing algorithm\n`HS256` - HMAC using SHA-256\n`HS384` - HMAC using SHA-384\n`HS512` - HMAC using SHA-512\n`RS256` - RSA using SHA-256\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256 `- ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
						Description:         "The JSON Web Signature [JWS] algorithm required for the ID Token. `NONE` - No signing algorithm, `HS256` - HMAC using SHA-256, `HS384` - HMAC using SHA-384, `HS512` - HMAC using SHA-512, `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256. RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("NONE",
								"HS256",
								"HS384",
								"HS512",
								"RS256",
								"RS384",
								"RS512",
								"ES256",
								"ES384",
								"ES512",
								"PS256",
								"PS384",
								"PS512",
							),
						},
					},
					"id_token_encryption_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token.\n`DIR` - Direct Encryption with symmetric key\n`A128KW` - AES-128 Key Wrap\n`A192KW` - AES-192 Key Wrap\n`A256KW`- AES-256 Key Wrap\n`A128GCMKW` - AES-GCM-128 key encryption\n`A192GCMKW` - AES-GCM-192 key encryption\n`A256GCMKW` - AES-GCM-256 key encryption\n`ECDH_ES` - ECDH-ES\n`ECDH_ES_A128KW` - ECDH-ES with AES-128 Key Wrap\n`ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap\n`ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap\n`RSA_OAEP` - RSAES OAEP\n`RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
						Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token. `DIR` - Direct Encryption with symmetric key, `A128KW` - AES-128 Key Wrap, `A192KW` - AES-192 Key Wrap, `A256KW` - AES-256 Key Wrap, `A128GCMKW` - AES-GCM-128 key encryption, `A192GCMKW` - AES-GCM-192 key encryption, `A256GCMKW` - AES-GCM-256 key encryption, `ECDH_ES` - ECDH-ES, `ECDH_ES_A128KW` - ECDH-ES with AES-128 Key Wrap, `ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap, `ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap, `RSA_OAEP` - RSAES OAEP, `RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("DIR",
								"A128KW",
								"A192KW",
								"A256KW",
								"A128GCMKW",
								"A192GCMKW",
								"A256GCMKW",
								"ECDH_ES",
								"ECDH_ES_A128KW",
								"ECDH_ES_A192KW",
								"ECDH_ES_A256KW",
								"RSA_OAEP",
								"RSA_OAEP_256",
							),
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("id_token_content_encryption_algorithm")),
						},
					},
					"id_token_content_encryption_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token.\n`AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256\n`AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384\n`AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512\n`AES_128_GCM` - AES-GCM-128\n`AES_192_GCM` - AES-GCM-192\n`AES_256_GCM` - AES-GCM-256",
						Description:         "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token. `AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256, `AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384, `AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512, `AES_128_GCM` - AES-GCM-128, `AES_192_GCM` - AES-GCM-192, `AES_256_GCM` - AES-GCM-256",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("AES_128_CBC_HMAC_SHA_256",
								"AES_192_CBC_HMAC_SHA_384",
								"AES_256_CBC_HMAC_SHA_512",
								"AES_128_GCM",
								"AES_192_GCM",
								"AES_256_GCM",
							),
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("id_token_encryption_algorithm")),
						},
					},
					"policy_group": schema.SingleNestedAttribute{
						Description: "The Open ID Connect policy. A null value will represent the default policy group.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"grant_access_session_revocation_api": schema.BoolAttribute{
						Description: "Determines whether this client is allowed to access the Session Revocation API. The default is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"grant_access_session_session_management_api": schema.BoolAttribute{
						Description: "Determines whether this client is allowed to access the Session Management API. The default is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"ping_access_logout_capable": schema.BoolAttribute{
						Description: "Set this value to `true` if you wish to enable client application logout, and the client is PingAccess, or its logout endpoints follow the PingAccess path convention. The default is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"logout_uris": schema.SetAttribute{
						Description: "A list of front-channel logout URIs for this client.",
						ElementType: types.StringType,
						Optional:    true,
					},
					"pairwise_identifier_user_type": schema.BoolAttribute{
						Description: "Determines whether the subject identifier type is pairwise. The default is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"sector_identifier_uri": schema.StringAttribute{
						Description: "The URI references a file with a single JSON array of Redirect URI and JWKS URL values.",
						Optional:    true,
						Validators: []validator.String{
							configvalidators.ValidUrl(),
						},
					},
					"logout_mode": schema.StringAttribute{
						Description: "The logout mode for this client. The default is 'NONE'. Supported values are `NONE`, `PING_FRONT_CHANNEL`, `OIDC_FRONT_CHANNEL`, and `OIDC_BACK_CHANNEL`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("NONE"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"NONE",
								"PING_FRONT_CHANNEL",
								"OIDC_FRONT_CHANNEL",
								"OIDC_BACK_CHANNEL",
							),
						},
					},
					"back_channel_logout_uri": schema.StringAttribute{
						Description: "The back-channel logout URI for this client.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"post_logout_redirect_uris": schema.SetAttribute{
						Description: "URIs to which the OIDC OP may redirect the resource owner's user agent after RP-initiated logout has completed. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible. Supported in PF version `12.0` or later.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"user_info_response_content_encryption_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "The JSON Web Encryption [JWE] content-encryption algorithm for the UserInfo Response. Supported values are `AES_128_CBC_HMAC_SHA_256`, `AES_192_CBC_HMAC_SHA_384`, `AES_256_CBC_HMAC_SHA_512`, `AES_128_GCM`, `AES_192_GCM`, `AES_256_GCM`. Supported in PF version `12.2` or later.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"AES_128_CBC_HMAC_SHA_256",
								"AES_192_CBC_HMAC_SHA_384",
								"AES_256_CBC_HMAC_SHA_512",
								"AES_128_GCM",
								"AES_192_GCM",
								"AES_256_GCM",
							),
						},
					},
					"user_info_response_encryption_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the UserInfo response. Supported values are `DIR`, `A128KW`, `A192KW`, `A256KW`, `A128GCMKW`, `A192GCMKW`, `A256GCMKW`, `ECDH_ES`, `ECDH_ES_A128KW`, `ECDH_ES_A192KW`, `ECDH_ES_A256KW`, `RSA_OAEP`, `RSA_OAEP_256`. Supported in PF version `12.2` or later.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"DIR",
								"A128KW",
								"A192KW",
								"A256KW",
								"A128GCMKW",
								"A192GCMKW",
								"A256GCMKW",
								"ECDH_ES",
								"ECDH_ES_A128KW",
								"ECDH_ES_A192KW",
								"ECDH_ES_A256KW",
								"RSA_OAEP",
								"RSA_OAEP_256",
							),
						},
					},
					"user_info_response_signing_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "The JSON Web Signature [JWS] algorithm required to sign the UserInfo response. Supported values are `NONE`, `HS256`, `HS384`, `HS512`, `RS256`, `RS384`, `RS512`, `ES256`, `ES384`, `ES512`, `PS256`, `PS384`, `PS512`. Supported in PF version `12.2` or later.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"NONE",
								"HS256",
								"HS384",
								"HS512",
								"RS256",
								"RS384",
								"RS512",
								"ES256",
								"ES384",
								"ES512",
								"PS256",
								"PS384",
								"PS512",
							),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"client_auth": schema.SingleNestedAttribute{
				Description: "Client authentication settings. If this model is null, it indicates that no client authentication will be used.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(clientAuthDefaultObj),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Client authentication type. The required field for type `SECRET` is `secret`.	The required fields for type `CERTIFICATE` are `client_cert_issuer_dn` and `client_cert_subject_dn`. The required field for type `PRIVATE_KEY_JWT` is: either `jwks` or `jwks_url`.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("NONE",
								"CERTIFICATE",
								"SECRET",
								"PRIVATE_KEY_JWT",
							),
						},
					},
					"secret": schema.StringAttribute{
						Description: "Client secret for Basic Authentication. Only one of `secret` or `encrypted_secret` can be set.",
						Optional:    true,
						Sensitive:   true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"encrypted_secret": schema.StringAttribute{
						Description: "Encrypted client secret for Basic Authentication. Only one of `secret` or `encrypted_secret` can be set.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("secret")),
						},
					},
					"secondary_secrets": schema.ListNestedAttribute{
						Description: "The list of secondary client secrets that are temporarily retained.",
						Computed:    true,
						Optional:    true,
						Default:     listdefault.StaticValue(secondarySecretsEmptyList),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"secret": schema.StringAttribute{
									Description: "Secondary client secret for Basic Authentication. Either this attribute or `encrypted_secret` must be provided.",
									Optional:    true,
									Sensitive:   true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"encrypted_secret": schema.StringAttribute{
									Description: "Encrypted secondary client secret for Basic Authentication. Either this attribute or `secret` must be provided.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.String{
										stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("secret")),
									},
								},
								"expiry_time": schema.StringAttribute{
									Description: "The expiry time of the secondary secret.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
					},
					"client_cert_issuer_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Issuer DN.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"client_cert_subject_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Subject DN.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"enforce_replay_prevention": schema.BoolAttribute{
						Description: "Enforce replay prevention on JSON Web Tokens. This field is applicable only for Private Key JWT Client and Client Secret JWT Authentication.",
						Optional:    true,
					},
					"token_endpoint_auth_signing_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT and Client Secret JWT Client Authentication. All asymmetric signing algorithms are allowed for Private Key JWT if value is not present. All symmetric signing algorithms are allowed for Client Secret JWT if value is not present\n`RS256` - RSA using SHA-256\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256` - ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\n`RSASSA-PSS` is only supported with SafeNet Luna, Thales nCipher or Java 11.\n`HS256` - HMAC using SHA-256\n`HS384` - HMAC using SHA-384\n`HS512` - HMAC using SHA-512.",
						Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT and Client Secret JWT Client Authentication. All asymmetric signing algorithms are allowed for Private Key JWT if value is not present. All symmetric signing algorithms are allowed for Client Secret JWT if value is not present `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11. `HS256` - HMAC using SHA-256, `HS384` - HMAC using SHA-384, `HS512` - HMAC using SHA-512.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"RS256",
								"RS384",
								"RS512",
								"ES256",
								"ES384",
								"ES512",
								"PS256",
								"PS384",
								"PS512",
								"HS256",
								"HS384",
								"HS512",
							),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_settings": schema.SingleNestedAttribute{
				Description: "JSON Web Key Set Settings of the OAuth client. Required if private key JWT client authentication or signed requests is enabled.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"jwks_url": schema.StringAttribute{
						Description: "JSON Web Key Set (JWKS) URL of the OAuth client. Either `jwks` or `jwks_url` must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"jwks": schema.StringAttribute{
						Description: "JSON Web Key Set (JWKS) document of the OAuth client. Either `jwks` or `jwks_url` must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"extended_parameters": schema.MapNestedAttribute{
				Description: "OAuth Client Metadata can be extended to use custom Client Metadata Parameters. The names of these custom parameters should be defined in /extendedProperties.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.SetAttribute{
							Description: "A list of values",
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
			"device_flow_setting_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Device Authorization Settings set globally for the OAuth AS. Defaults to `SERVER_DEFAULT`. Supported values are `SERVER_DEFAULT` and `OVERRIDE_SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("SERVER_DEFAULT", "OVERRIDE_SERVER_DEFAULT"),
				},
			},
			"user_authorization_url_override": schema.StringAttribute{
				Description: "The URL used as `verification_url` and `verification_url_complete` values in a Device Authorization request. This property overrides the `user_authorization_url` value present in Authorization Server Settings.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"pending_authorization_timeout_override": schema.Int64Attribute{
				Description: "The `device_code` and `user_code` timeout, in seconds. This overrides the `pending_authorization_timeout` value present in Authorization Server Settings.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"device_polling_interval_override": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds. This overrides the 'device_polling_interval' value present in Authorization Server Settings.",
				Optional:    true,
			},
			"bypass_activation_code_confirmation_override": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if `verification_url_complete` is used by the end user to authorize a device. This overrides the `bypass_use_code_confirmation` value present in Authorization Server Settings.",
				Optional:    true,
			},
			"require_proof_key_for_code_exchange": schema.BoolAttribute{
				Description: "Determines whether Proof Key for Code Exchange (PKCE) is required for this client. Defaults to `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ciba_delivery_mode": schema.StringAttribute{
				Description: "The token delivery mode for the client. The default value is `POLL`. Supported values are `POLL` and `PING`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("POLL", "PING"),
				},
			},
			"ciba_notification_endpoint": schema.StringAttribute{
				Description: "The endpoint the OP will call after a successful or failed end-user authentication.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ciba_polling_interval": schema.Int64Attribute{
				Description: "The minimum amount of time in seconds that the Client must wait between polling requests to the token endpoint. The default is `0` seconds. Must be between `0` and `3600` seconds.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
					int64validator.AtMost(3600),
				},
			},
			"ciba_require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether CIBA signed requests are required for this client.",
				Optional:    true,
			},
			"ciba_request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present\n`RS256` - RSA using SHA-256\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256` - ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present, `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RS256",
						"RS384",
						"RS512",
						"ES256",
						"ES384",
						"ES512",
						"PS256",
						"PS384",
						"PS512",
					),
				},
			},
			"ciba_user_code_supported": schema.BoolAttribute{
				Description: "Determines whether the CIBA user code parameter is supported by this client.",
				Optional:    true,
			},
			"request_policy_ref": schema.SingleNestedAttribute{
				Description: "The CIBA request policy.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"token_exchange_processor_policy_ref": schema.SingleNestedAttribute{
				Description: "The Token Exchange Processor policy.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"refresh_token_rolling_grace_period_type": schema.StringAttribute{
				Description: "When specified, it overrides the global Refresh Token Grace Period defined in the Authorization Server Settings. The default value is `SERVER_DEFAULT`. Supported values are `SERVER_DEFAULT` and `OVERRIDE_SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("OVERRIDE_SERVER_DEFAULT", "SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"client_secret_retention_period_type": schema.StringAttribute{
				Description: "Use `OVERRIDE_SERVER_DEFAULT` to override the Client Secret Retention Period value on the Authorization Server Settings. `SERVER_DEFAULT` will default to the Client Secret Retention Period value on the Authorization Server Setting. Defaults to `SERVER_DEFAULT`. Supported values are `OVERRIDE_SERVER_DEFAULT` and `SERVER_DEFAULT`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("OVERRIDE_SERVER_DEFAULT", "SERVER_DEFAULT"),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is `0`, which will disable secondary client secret retention. This value will override the Client Secret Retention Period value on the Authorization Server Settings.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"client_secret_changed_time": schema.StringAttribute{
				Description: "The time at which the client secret was last changed.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_introspection_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response.\n`HS256` - HMAC using SHA-256\n`HS384` - HMAC using SHA-384\n`HS512`- HMAC using SHA-512\n`RS256` - RSA using SHA-256\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256` - ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response. `HS256` - HMAC using SHA-256, `HS384` - HMAC using SHA-384, `HS512` - HMAC using SHA-512, `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RS256",
						"RS384",
						"RS512",
						"HS256",
						"HS384",
						"HS512",
						"ES256",
						"ES384",
						"ES512",
						"PS256",
						"PS384",
						"PS512",
					),
				},
			},
			"token_introspection_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response.\n`DIR` - Direct Encryption with symmetric key\n`A128KW` - AES-128 Key Wrap\n`A192KW` - AES-192 Key Wrap\n`A256KW` - AES-256 Key Wrap\n`A128GCMKW` - AES-GCM-128 key encryption\n`A192GCMKW` - AES-GCM-192 key encryption\n`A256GCMKW` - AES-GCM-256 key encryption\n`ECDH_ES` - ECDH-ES\n`ECDH_ES_A128KW`- ECDH-ES with AES-128 Key Wrap\n`ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap\n`ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap\n`RSA_OAEP` - RSAES OAEP\n`RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response. `DIR` - Direct Encryption with symmetric key, `A128KW` - AES-128 Key Wrap, `A192KW` - AES-192 Key Wrap, `A256KW` - AES-256 Key Wrap, `A128GCMKW` - AES-GCM-128 key encryption, `A192GCMKW` - AES-GCM-192 key encryption, `A256GCMKW` - AES-GCM-256 key encryption, `ECDH_ES` - ECDH-ES, `ECDH_ES_A128KW` - ECDH-ES with AES-128 Key Wrap, `ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap, `ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap, `RSA_OAEP` - RSAES OAEP, `RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("DIR",
						"A128KW",
						"A192KW",
						"A256KW",
						"A128GCMKW",
						"A192GCMKW",
						"A256GCMKW",
						"ECDH_ES",
						"ECDH_ES_A128KW",
						"ECDH_ES_A192KW",
						"ECDH_ES_A256KW",
						"RSA_OAEP",
						"RSA_OAEP_256",
					),
				},
			},
			"token_introspection_content_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response.\n`AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256\n`AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384\n`AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512\n`AES_128_GCM` - AES-GCM-128\n`AES_192_GCM` - AES-GCM-192\n`AES_256_GCM` - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response. `AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256, `AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384, `AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512, `AES_128_GCM` - AES-GCM-128, `AES_192_GCM` - AES-GCM-192, `AES_256_GCM` - AES-GCM-256",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AES_128_CBC_HMAC_SHA_256",
						"AES_192_CBC_HMAC_SHA_384",
						"AES_256_CBC_HMAC_SHA_512",
						"AES_128_GCM",
						"AES_192_GCM",
						"AES_256_GCM",
					),
				},
			},
			"jwt_secured_authorization_response_mode_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response.\n`HS256` - HMAC using SHA-256\n`HS384` - HMAC using SHA-384\n`HS512` - HMAC using SHA-512\n`RS256` - RSA using SHA-256\n`RS384` - RSA using SHA-384\n`RS512` - RSA using SHA-512\n`ES256` - ECDSA using P256 Curve and SHA-256\n`ES384` - ECDSA using P384 Curve and SHA-384\n`ES512` - ECDSA using P521 Curve and SHA-512\n`PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\n`PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\n`PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response. `HS256` - HMAC using SHA-256, `HS384` - HMAC using SHA-384, `HS512` - HMAC using SHA-512, `RS256` - RSA using SHA-256, `RS384` - RSA using SHA-384, `RS512` - RSA using SHA-512, `ES256` - ECDSA using P256 Curve and SHA-256, `ES384` - ECDSA using P384 Curve and SHA-384, `ES512` - ECDSA using P521 Curve and SHA-512, `PS256` - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, `PS384` - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, `PS512` - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RS256",
						"RS384",
						"RS512",
						"HS256",
						"HS384",
						"HS512",
						"ES256",
						"ES384",
						"ES512",
						"PS256",
						"PS384",
						"PS512",
					),
				},
			},
			"jwt_secured_authorization_response_mode_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response.\n`DIR` - Direct Encryption with symmetric key\n`A128KW` - AES-128 Key Wrap\n`A192KW` - AES-192 Key Wrap\n`A256KW` - AES-256 Key Wrap\n`A128GCMKW` - AES-GCM-128 key encryption\n`A192GCMKW` - AES-GCM-192 key encryption\n`A256GCMKW` - AES-GCM-256 key encryption\n`ECDH_ES` - ECDH-ES\n`ECDH_ES_A128KW` - ECDH-ES with AES-128 Key Wrap\n`ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap\n`ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap\n`RSA_OAEP` - RSAES OAEP\n`RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response. `DIR` - Direct Encryption with symmetric key, `A128KW` - AES-128 Key Wrap, `A192KW` - AES-192 Key Wrap, `A256KW` - AES-256 Key Wrap, `A128GCMKW` - AES-GCM-128 key encryption, `A192GCMKW` - AES-GCM-192 key encryption, `A256GCMKW` - AES-GCM-256 key encryption, `ECDH_ES` - ECDH-ES, `ECDH_ES_A128KW` - ECDH-ES with AES-128 Key Wrap, `ECDH_ES_A192KW` - ECDH-ES with AES-192 Key Wrap, `ECDH_ES_A256KW` - ECDH-ES with AES-256 Key Wrap, `RSA_OAEP` - RSAES OAEP, `RSA_OAEP_256` - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("DIR",
						"A128KW",
						"A192KW",
						"A256KW",
						"A128GCMKW",
						"A192GCMKW",
						"A256GCMKW",
						"ECDH_ES",
						"ECDH_ES_A128KW",
						"ECDH_ES_A192KW",
						"ECDH_ES_A256KW",
						"RSA_OAEP",
						"RSA_OAEP_256",
					),
				},
			},
			"jwt_secured_authorization_response_mode_content_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response.\n`AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256\n`AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384\n`AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512\n`AES_128_GCM` - AES-GCM-128\n`AES_192_GCM` - AES-GCM-192\n`AES_256_GCM` - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response. `AES_128_CBC_HMAC_SHA_256` - Composite AES-CBC-128 HMAC-SHA-256, `AES_192_CBC_HMAC_SHA_384` - Composite AES-CBC-192 HMAC-SHA-384, `AES_256_CBC_HMAC_SHA_512` - Composite AES-CBC-256 HMAC-SHA-512, `AES_128_GCM` - AES-GCM-128, `AES_192_GCM` - AES-GCM-192, `AES_256_GCM` - AES-GCM-256",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AES_128_CBC_HMAC_SHA_256",
						"AES_192_CBC_HMAC_SHA_384",
						"AES_256_CBC_HMAC_SHA_512",
						"AES_128_GCM",
						"AES_192_GCM",
						"AES_256_GCM",
					),
				},
			},
			"require_dpop": schema.BoolAttribute{
				MarkdownDescription: "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client. Defaults to `false`.",
				Description:         "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"require_offline_access_scope_to_issue_refresh_tokens": schema.StringAttribute{
				Description: "Determines whether offline_access scope is required to issue refresh tokens by this client or not. `SERVER_DEFAULT` is the default value. Supported values are `SERVER_DEFAULT`, `NO`, and `YES`. Supported in PF version `12.1` or later.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"SERVER_DEFAULT",
						"NO",
						"YES",
					),
				},
			},
			"offline_access_require_consent_prompt": schema.StringAttribute{
				Description: "Determines whether offline_access requires the prompt parameter value to be set to 'consent' by this client or not. The value will be reset to default if the `require_offline_access_scope_to_issue_refresh_tokens` attribute is set to `SERVER_DEFAULT` or `false`. `SERVER_DEFAULT` is the default value. Supported values are `SERVER_DEFAULT`, `NO`, and `YES`. Supported in PF version `12.1` or later.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"SERVER_DEFAULT",
						"NO",
						"YES",
					),
				},
			},
			"lockout_max_malicious_actions": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of malicious actions allowed before an OAuth client is locked out. Currently, the only operation that is tracked as a malicious action is an attempt to revoke an invalid access token or refresh token. This value will override the global `MaxMaliciousActions` value on the `AccountLockingService` in the config-store. Supported in PF version `12.2` or later.",
			},
			"lockout_max_malicious_actions_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Allows an administrator to override the Max Malicious Actions configuration set globally in `AccountLockingService`. Defaults to `SERVER_DEFAULT`. Supported values are `DO_NOT_LOCKOUT`, `SERVER_DEFAULT`, `OVERRIDE_SERVER_DEFAULT`. Supported in PF version `12.2` or later.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DO_NOT_LOCKOUT",
						"SERVER_DEFAULT",
						"OVERRIDE_SERVER_DEFAULT",
					),
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *oauthClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_client"
}

func (r *oauthClientResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *oauthClientResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model *oauthClientModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if model == nil {
		return
	}

	// Persistent Grant Expiration Validation
	if (internaltypes.IsDefined(model.PersistentGrantExpirationTime) || internaltypes.IsDefined(model.PersistentGrantExpirationTimeUnit)) &&
		!model.PersistentGrantExpirationType.IsUnknown() && model.PersistentGrantExpirationType.ValueString() != "OVERRIDE_SERVER_DEFAULT" {
		resp.Diagnostics.AddAttributeError(
			path.Root("persistent_grant_expiration_time"),
			providererror.InvalidAttributeConfiguration,
			"persistent_grant_expiration_type must be configured to \"OVERRIDE_SERVER_DEFAULT\" to modify the other persistent_grant_expiration values.")
	}

	// Refresh Token Rolling Validation
	if !model.RefreshTokenRollingIntervalType.IsUnknown() {
		if model.RefreshTokenRollingIntervalType.ValueString() == "SERVER_DEFAULT" {
			// The refresh_token_rolling_interval and refresh_token_rolling_interval_time_unit value can't be
			// configured with a non-default value when refresh_token_rolling_interval_type is set to "SERVER_DEFAULT"
			if internaltypes.IsDefined(model.RefreshTokenRollingInterval) {
				resp.Diagnostics.AddAttributeError(
					path.Root("refresh_token_rolling_interval"),
					providererror.InvalidAttributeConfiguration,
					"refresh_token_rolling_interval can only be configured if refresh_token_rolling_interval_type is set to \"OVERRIDE_SERVER_DEFAULT\".")
			}
			if internaltypes.IsDefined(model.RefreshTokenRollingIntervalTimeUnit) && model.RefreshTokenRollingIntervalTimeUnit.ValueString() != "HOURS" {
				resp.Diagnostics.AddAttributeError(
					path.Root("refresh_token_rolling_interval_time_unit"),
					providererror.InvalidAttributeConfiguration,
					"refresh_token_rolling_interval_time_unit can only be configured if refresh_token_rolling_interval_type is \"OVERRIDE_SERVER_DEFAULT\".")
			}
		} else if model.RefreshTokenRollingIntervalType.ValueString() == "OVERRIDE_SERVER_DEFAULT" && model.RefreshTokenRollingInterval.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("refresh_token_rolling_interval"),
				providererror.InvalidAttributeConfiguration,
				"refresh_token_rolling_interval must be configured when refresh_token_rolling_interval_type is \"OVERRIDE_SERVER_DEFAULT\".")
		}
	}

	//  Client Auth Defined
	var clientAuthAttributes map[string]attr.Value
	if internaltypes.IsDefined(model.ClientAuth) {
		clientAuthAttributes = model.ClientAuth.Attributes()
		if internaltypes.IsDefined(clientAuthAttributes["type"]) {
			clientAuthType := clientAuthAttributes["type"].(types.String).ValueString()
			switch clientAuthType {
			case "PRIVATE_KEY_JWT":
				errorMsg := "jwks_settings.jwks or jwks_settings.jwks_url must be defined when client_auth is configured to \"PRIVATE_KEY_JWT\"."
				if model.JwksSettings.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("jwks_settings"),
						providererror.InvalidAttributeConfiguration,
						errorMsg)
				} else if !model.JwksSettings.IsUnknown() {
					jwksSettingsAttributes := model.JwksSettings.Attributes()
					if jwksSettingsAttributes["jwks"].IsNull() && jwksSettingsAttributes["jwks_url"].IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("jwks_settings"),
							providererror.InvalidAttributeConfiguration,
							errorMsg)
					}
				}
			case "CERTIFICATE":
				if clientAuthAttributes["client_cert_subject_dn"].IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("client_auth"),
						providererror.InvalidAttributeConfiguration,
						"client_cert_subject_dn must be defined when client_auth.type is configured to \"CERTIFICATE\".")
				}
				if clientAuthAttributes["client_cert_issuer_dn"].IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("client_auth"),
						providererror.InvalidAttributeConfiguration,
						"client_cert_issuer_dn must be defined when client_auth.type is configured to \"CERTIFICATE\".")
				}
			case "SECRET":
				if clientAuthAttributes["secret"].IsNull() && clientAuthAttributes["encrypted_secret"].IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("client_auth"),
						providererror.InvalidAttributeConfiguration,
						"client_auth.secret or client_auth.encrypted_secret must be defined when client_auth.type is configured to \"SECRET\".")
				}
			}
		}
	}

	// Grant Types Validation
	// grant_types is required, don't need nil check here
	var hasCibaGrantType bool
	for _, grantType := range model.GrantTypes.Elements() {
		grantTypeVal := grantType.(types.String).ValueString()
		if grantTypeVal == "CLIENT_CREDENTIALS" {
			if model.ClientAuth.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("client_auth"),
					providererror.InvalidAttributeConfiguration,
					"client_auth must be defined when \"CLIENT_CREDENTIALS\" is included in grant_types.")
			}
		}
		if grantTypeVal == "CIBA" {
			hasCibaGrantType = true
		}
	}

	// CIBA Validation
	if !model.GrantTypes.IsUnknown() && !hasCibaGrantType && (internaltypes.IsDefined(model.CibaDeliveryMode) ||
		internaltypes.IsDefined(model.CibaNotificationEndpoint) ||
		internaltypes.IsDefined(model.CibaPollingInterval) ||
		internaltypes.IsDefined(model.CibaRequireSignedRequests) ||
		internaltypes.IsDefined(model.CibaRequestObjectSigningAlgorithm) ||
		internaltypes.IsDefined(model.CibaUserCodeSupported)) {
		resp.Diagnostics.AddError(providererror.InvalidAttributeConfiguration, "ciba attributes can only be configured when \"CIBA\" is included in grant_types.")
	}
	if hasCibaGrantType && model.CibaDeliveryMode.ValueString() == "PING" && model.CibaNotificationEndpoint.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ciba_notification_endpoint"),
			providererror.InvalidAttributeConfiguration,
			"ciba_notification_endpoint must be defined when ciba_delivery_mode is \"PING\".")
	}

	// Client Auth Validation
	// ID Token Signing Algorithm Validation when client_auth is not defined
	if model.ClientAuth.IsNull() {
		var algorithmAttributeSet []string
		if internaltypes.IsDefined(model.OidcPolicy) && !model.OidcPolicy.Attributes()["id_token_signing_algorithm"].IsUnknown() {
			algorithmAttributeSet = append(algorithmAttributeSet, model.OidcPolicy.Attributes()["id_token_signing_algorithm"].(types.String).ValueString())
		}

		if internaltypes.IsDefined(model.TokenIntrospectionSigningAlgorithm) {
			algorithmAttributeSet = append(algorithmAttributeSet, model.TokenIntrospectionSigningAlgorithm.ValueString())
		}

		if internaltypes.IsDefined(model.JwtSecuredAuthorizationResponseModeSigningAlgorithm) {
			algorithmAttributeSet = append(algorithmAttributeSet, model.JwtSecuredAuthorizationResponseModeSigningAlgorithm.ValueString())
		}

		for _, algorithmVal := range algorithmAttributeSet {
			if algorithmVal == "HS256" {
				resp.Diagnostics.AddAttributeError(
					path.Root("client_auth"),
					providererror.InvalidAttributeConfiguration,
					"client_auth must be defined when using the \"HS256\" signing algorithm")
			}
		}

		if internaltypes.IsDefined(model.TokenIntrospectionEncryptionAlgorithm) {
			resp.Diagnostics.AddAttributeError(
				path.Root("token_introspection_encryption_algorithm"),
				providererror.InvalidAttributeConfiguration,
				"client_auth must be configured when token_introspection_encryption_algorithm is configured.")
		}
	}

	// Restrict Scopes Validation
	if internaltypes.IsDefined(model.RestrictScopes) && !model.RestrictScopes.ValueBool() && model.AllowAuthenticationApiInit.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("restrict_scopes"),
			providererror.InvalidAttributeConfiguration,
			"restrict_scopes cannot be configured to false when allow_authentication_api_init is set to true.")
	}

	if len(model.RestrictedScopes.Elements()) > 0 && !model.RestrictScopes.IsUnknown() && !model.RestrictScopes.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("restricted_scopes"),
			providererror.InvalidAttributeConfiguration,
			"restrict_scopes must be set to true to configure restricted_scopes.")
	}

	// OIDC Policy Validation
	if internaltypes.IsDefined(model.OidcPolicy) {
		oidcPolicy := model.OidcPolicy.Attributes()
		pairwiseIdentifierUserType := oidcPolicy["pairwise_identifier_user_type"]
		oidcPolicySectorIdentifierUri := oidcPolicy["sector_identifier_uri"]
		if !pairwiseIdentifierUserType.IsUnknown() && !pairwiseIdentifierUserType.(types.Bool).ValueBool() &&
			internaltypes.IsDefined(oidcPolicySectorIdentifierUri) {
			resp.Diagnostics.AddAttributeError(
				path.Root("oidc_policy").AtMapKey("sector_identifier_uri"),
				providererror.InvalidAttributeConfiguration,
				"sector_identifier_uri can only be configured when pairwise_identifier_user_type is set to true.")
		}
	}

	// JWKS Settings Validation
	if model.JwksSettings.IsNull() {
		if internaltypes.IsDefined(model.TokenIntrospectionEncryptionAlgorithm) {
			resp.Diagnostics.AddAttributeError(
				path.Root("token_introspection_encryption_algorithm"),
				providererror.InvalidAttributeConfiguration,
				"token_introspection_encryption_algorithm must not be configured when jwks_settings is not configured.")
		}
		if model.RequireSignedRequests.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				path.Root("require_signed_requests"),
				providererror.InvalidAttributeConfiguration,
				"require_signed_requests must be false when jwks_settings is not configured.")
		}
	}

	// offline_access_require_consent_prompt can only be configured if require_offline_access_scope_to_issue_refresh_tokens is set to "YES"
	if !model.RequireOfflineAccessScopeToIssueRefreshTokens.IsUnknown() && model.RequireOfflineAccessScopeToIssueRefreshTokens.ValueString() != "YES" &&
		internaltypes.IsDefined(model.OfflineAccessRequireConsentPrompt) && model.OfflineAccessRequireConsentPrompt.ValueString() != "SERVER_DEFAULT" {
		resp.Diagnostics.AddAttributeError(
			path.Root("offline_access_require_consent_prompt"),
			providererror.InvalidAttributeConfiguration,
			"offline_access_require_consent_prompt can only be configured if require_offline_access_scope_to_issue_refresh_tokens is set to \"YES\".\n"+
				fmt.Sprintf("require_offline_access_scope_to_issue_refresh_tokens: %s\noffline_access_require_consent_prompt: %s", model.RequireOfflineAccessScopeToIssueRefreshTokens.ValueString(), model.OfflineAccessRequireConsentPrompt.ValueString()))
	}

	// bypass_approval_page Validation
	if internaltypes.IsDefined(model.BypassApprovalPage) && !model.BypassApprovalPage.ValueBool() && model.AllowAuthenticationApiInit.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("bypass_approval_page"),
			providererror.InvalidAttributeConfiguration,
			"bypass_approval_page cannot be configured to false when allow_authentication_api_init is set to true.")
	}

	// lockout_max_malicious_actions validation
	if internaltypes.IsDefined(model.LockoutMaxMaliciousActions) &&
		!model.LockoutMaxMaliciousActionsType.IsUnknown() && model.LockoutMaxMaliciousActionsType.ValueString() != "OVERRIDE_SERVER_DEFAULT" {
		resp.Diagnostics.AddAttributeError(
			path.Root("lockout_max_malicious_actions"),
			providererror.InvalidAttributeConfiguration,
			"lockout_max_malicious_actions_type must be configured to \"OVERRIDE_SERVER_DEFAULT\" to set lockout_max_malicious_actions.")
	}
}

func (r *oauthClientResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
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
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1220)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast122 := compare >= 0
	var plan *oauthClientModel
	var state *oauthClientModel
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if plan == nil {
		return
	}

	planModified := false
	if internaltypes.IsDefined(plan.OidcPolicy) {
		planOidcPolicyAttrs := plan.OidcPolicy.Attributes()
		// If oidc_policy.post_logout_redirect_uris is set prior to PF version 12.0, throw an error.
		planPostLogoutRedirectUris := planOidcPolicyAttrs["post_logout_redirect_uris"].(types.Set)
		if !pfVersionAtLeast120 && internaltypes.IsDefined(planPostLogoutRedirectUris) {
			version.AddUnsupportedAttributeError("oidc_policy.post_logout_redirect_uris",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		}
		// Check for OIDC policy attrs added in PF 12.2
		if !pfVersionAtLeast122 {
			userInfoResponseContentEncryptionAlgorithm := plan.OidcPolicy.Attributes()["user_info_response_content_encryption_algorithm"]
			if internaltypes.IsDefined(userInfoResponseContentEncryptionAlgorithm) {
				version.AddUnsupportedAttributeError("oidc_policy.user_info_response_content_encryption_algorithm",
					r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
			}
			userInfoResponseEncryptionAlgorithm := plan.OidcPolicy.Attributes()["user_info_response_encryption_algorithm"]
			if internaltypes.IsDefined(userInfoResponseEncryptionAlgorithm) {
				version.AddUnsupportedAttributeError("oidc_policy.user_info_response_encryption_algorithm",
					r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
			}
			userInfoResponseSigningAlgorithm := plan.OidcPolicy.Attributes()["user_info_response_signing_algorithm"]
			if internaltypes.IsDefined(userInfoResponseSigningAlgorithm) {
				version.AddUnsupportedAttributeError("oidc_policy.user_info_response_signing_algorithm",
					r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
			}
		}
	}

	// Version checking and default settings for attrs added in PF 12.1
	if !pfVersionAtLeast121 {
		planModified = true
		if internaltypes.IsDefined(plan.EnableCookielessAuthenticationApi) {
			version.AddUnsupportedAttributeError("enable_cookieless_authentication_api",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else {
			plan.EnableCookielessAuthenticationApi = types.BoolNull()
		}

		if internaltypes.IsDefined(plan.RefreshTokenRollingIntervalTimeUnit) {
			version.AddUnsupportedAttributeError("refresh_token_rolling_interval_time_unit",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else {
			plan.RefreshTokenRollingIntervalTimeUnit = types.StringNull()
		}

		if internaltypes.IsDefined(plan.RequireOfflineAccessScopeToIssueRefreshTokens) {
			version.AddUnsupportedAttributeError("require_offline_access_scope_to_issue_refresh_tokens",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else {
			plan.RequireOfflineAccessScopeToIssueRefreshTokens = types.StringNull()
		}

		if internaltypes.IsDefined(plan.OfflineAccessRequireConsentPrompt) {
			version.AddUnsupportedAttributeError("offline_access_require_consent_prompt",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		} else {
			plan.OfflineAccessRequireConsentPrompt = types.StringNull()
		}
	} else {
		if plan.EnableCookielessAuthenticationApi.IsUnknown() {
			plan.EnableCookielessAuthenticationApi = types.BoolValue(false)
			planModified = true
		}
		if plan.RefreshTokenRollingIntervalTimeUnit.IsUnknown() {
			plan.RefreshTokenRollingIntervalTimeUnit = types.StringValue("HOURS")
			planModified = true
		}
		if plan.RequireOfflineAccessScopeToIssueRefreshTokens.IsUnknown() {
			plan.RequireOfflineAccessScopeToIssueRefreshTokens = types.StringValue("SERVER_DEFAULT")
			planModified = true
		}
		if plan.OfflineAccessRequireConsentPrompt.IsUnknown() {
			plan.OfflineAccessRequireConsentPrompt = types.StringValue("SERVER_DEFAULT")
			planModified = true
		}
	}

	// Version checking and default settings for attrs added in PF 12.2
	if !pfVersionAtLeast122 {
		if internaltypes.IsDefined(plan.LockoutMaxMaliciousActions) {
			version.AddUnsupportedAttributeError("lockout_max_malicious_actions",
				r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
		}
		if internaltypes.IsDefined(plan.LockoutMaxMaliciousActionsType) {
			version.AddUnsupportedAttributeError("lockout_max_malicious_actions_type",
				r.providerConfig.ProductVersion, version.PingFederate1220, &resp.Diagnostics)
		} else {
			plan.LockoutMaxMaliciousActionsType = types.StringNull()
		}
	} else {
		if plan.LockoutMaxMaliciousActionsType.IsUnknown() {
			plan.LockoutMaxMaliciousActionsType = types.StringValue("SERVER_DEFAULT")
			planModified = true
		}
	}

	if plan.RestrictScopes.IsUnknown() {
		plan.RestrictScopes = types.BoolValue(plan.AllowAuthenticationApiInit.ValueBool())
		planModified = true
	}
	if plan.BypassApprovalPage.IsUnknown() {
		plan.BypassApprovalPage = types.BoolValue(plan.AllowAuthenticationApiInit.ValueBool())
		planModified = true
	}

	// Set encrypted values as necessary
	if internaltypes.IsDefined(plan.ClientAuth) && state != nil {
		clientAuthAttrs := plan.ClientAuth.Attributes()
		stateClientAuthAttrs := state.ClientAuth.Attributes()
		if clientAuthAttrs["secret"].IsNull() && clientAuthAttrs["encrypted_secret"].IsUnknown() {
			clientAuthAttrs["encrypted_secret"] = types.StringNull()
		} else if !stateClientAuthAttrs["secret"].Equal(clientAuthAttrs["secret"]) {
			clientAuthAttrs["encrypted_secret"] = types.StringUnknown()
		}
		if internaltypes.IsDefined(clientAuthAttrs["secondary_secrets"]) && internaltypes.IsDefined(stateClientAuthAttrs["secondary_secrets"]) &&
			!clientAuthAttrs["secondary_secrets"].Equal(stateClientAuthAttrs["secondary_secrets"]) {
			secondarySecrets := clientAuthAttrs["secondary_secrets"].(types.List).Elements()
			updatedSecondarySecrets := []attr.Value{}
			for _, secondarySecret := range secondarySecrets {
				secondarySecretAttrs := secondarySecret.(types.Object).Attributes()
				secondarySecretAttrs["encrypted_secret"] = types.StringUnknown()
				updatedSecondarySecret, diags := types.ObjectValue(secondarySecretsAttrType, secondarySecretAttrs)
				resp.Diagnostics.Append(diags...)
				updatedSecondarySecrets = append(updatedSecondarySecrets, updatedSecondarySecret)
			}
			finalSecondarySecrets, diags := types.ListValue(types.ObjectType{AttrTypes: secondarySecretsAttrType}, updatedSecondarySecrets)
			resp.Diagnostics.Append(diags...)
			clientAuthAttrs["secondary_secrets"] = finalSecondarySecrets
		}
		plan.ClientAuth, diags = types.ObjectValue(clientAuthAttrType, clientAuthAttrs)
		resp.Diagnostics.Append(diags...)
	}

	// If the new plan doesn't match the state, invalidate any last-changed time values
	// See https://github.com/hashicorp/terraform-plugin-framework/issues/898 for some info on why this is needed
	resp.Diagnostics.Append(req.Plan.Set(ctx, plan)...)
	if !req.Plan.Raw.Equal(req.State.Raw) {
		plan.ModificationDate = types.StringUnknown()
		plan.ClientSecretChangedTime = types.StringUnknown()
		planModified = true
	}

	if planModified {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

func readOauthClientResponse(ctx context.Context, r *client.Client, plan, state *oauthClientModel, productVersion version.SupportedVersion, isImportRead bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	diags = readOauthClientResponseCommon(ctx, r, state, plan, productVersion, isImportRead)

	// state.ClientAuth
	var clientAuthToState types.Object
	clientAuthFromPlan := plan.ClientAuth.Attributes()
	var secretToState, encryptedSecretToState basetypes.StringValue

	// state.ClientAuth.Secret
	secretVal := clientAuthFromPlan["secret"]
	if secretVal != nil && internaltypes.IsNonEmptyString(secretVal.(types.String)) {
		secretToState = types.StringValue(secretVal.(types.String).ValueString())
	} else {
		secretToState = types.StringNull()
	}

	// state.ClientAuth.EncryptedSecret
	encryptedSecretVal := clientAuthFromPlan["encrypted_secret"]
	if encryptedSecretVal != nil && internaltypes.IsDefined(encryptedSecretVal) {
		encryptedSecretToState = types.StringValue(encryptedSecretVal.(types.String).ValueString())
	} else {
		encryptedSecretToState = types.StringPointerValue(r.ClientAuth.EncryptedSecret)
	}

	// state.ClientAuth.SecondarySecrets
	var secondarySecretsObjToState types.List
	var secondarySecretsListSlice []attr.Value
	secondarySecretsFromPlan := clientAuthFromPlan["secondary_secrets"]
	if secondarySecretsFromPlan != nil && len(secondarySecretsFromPlan.(types.List).Elements()) > 0 {
		// Copy secret values from plan
		for i, secondarySecretsFromPlan := range clientAuthFromPlan["secondary_secrets"].(types.List).Elements() {
			if i < len(r.ClientAuth.SecondarySecrets) {
				planAttrs := secondarySecretsFromPlan.(types.Object).Attributes()
				expiryTime := types.StringNull()
				if r.ClientAuth.SecondarySecrets[i].ExpiryTime != nil {
					expiryTime = types.StringValue(r.ClientAuth.SecondarySecrets[i].ExpiryTime.Format(time.RFC3339Nano))
				}
				// Maintain encrypted secret from plan if included
				encryptedSecret := types.StringPointerValue(r.ClientAuth.SecondarySecrets[i].EncryptedSecret)
				if internaltypes.IsDefined(planAttrs["encrypted_secret"]) {
					encryptedSecret = types.StringValue(planAttrs["encrypted_secret"].(types.String).ValueString())
				}
				secret := types.StringNull()
				if internaltypes.IsDefined(planAttrs["secret"]) {
					secret = types.StringValue(planAttrs["secret"].(types.String).ValueString())
				}
				secondarySecretsAttrVal, respDiags := types.ObjectValue(secondarySecretsAttrType, map[string]attr.Value{
					"secret":           secret,
					"encrypted_secret": encryptedSecret,
					"expiry_time":      expiryTime,
				})
				diags.Append(respDiags...)
				secondarySecretsListSlice = append(secondarySecretsListSlice, secondarySecretsAttrVal)
			}
		}
	} else {
		// Read values directly from response
		for _, secondarySecretsFromResponse := range r.ClientAuth.SecondarySecrets {
			expiryTime := types.StringNull()
			if secondarySecretsFromResponse.ExpiryTime != nil {
				expiryTime = types.StringValue(secondarySecretsFromResponse.ExpiryTime.Format(time.RFC3339Nano))
			}
			secondarySecretsAttrVal, respDiags := types.ObjectValue(secondarySecretsAttrType, map[string]attr.Value{
				"secret":           types.StringPointerValue(secondarySecretsFromResponse.Secret),
				"encrypted_secret": types.StringPointerValue(secondarySecretsFromResponse.EncryptedSecret),
				"expiry_time":      expiryTime,
			})
			diags.Append(respDiags...)
			secondarySecretsListSlice = append(secondarySecretsListSlice, secondarySecretsAttrVal)
		}
	}
	secondarySecretsObjToState, respDiags = types.ListValue(types.ObjectType{AttrTypes: secondarySecretsAttrType}, secondarySecretsListSlice)
	diags.Append(respDiags...)

	// state.ClientAuth to state
	clientAuthAttrValue := map[string]attr.Value{}
	clientAuthAttrValue["type"] = types.StringPointerValue(r.ClientAuth.Type)
	clientAuthAttrValue["secret"] = secretToState
	clientAuthAttrValue["encrypted_secret"] = encryptedSecretToState
	clientAuthAttrValue["secondary_secrets"] = secondarySecretsObjToState
	clientAuthAttrValue["client_cert_issuer_dn"] = types.StringPointerValue(r.ClientAuth.ClientCertIssuerDn)
	clientAuthAttrValue["client_cert_subject_dn"] = types.StringPointerValue(r.ClientAuth.ClientCertSubjectDn)
	clientAuthAttrValue["enforce_replay_prevention"] = types.BoolPointerValue(r.ClientAuth.EnforceReplayPrevention)
	clientAuthAttrValue["token_endpoint_auth_signing_algorithm"] = types.StringPointerValue(r.ClientAuth.TokenEndpointAuthSigningAlgorithm)
	clientAuthToState, respDiags = types.ObjectValue(clientAuthAttrType, clientAuthAttrValue)
	diags.Append(respDiags...)
	state.ClientAuth = clientAuthToState

	return diags
}

func grantTypes(grantTypesSet types.Set) []string {
	var grantTypesSlice []string
	for _, grantType := range grantTypesSet.Elements() {
		grantTypesSlice = append(grantTypesSlice, grantType.(types.String).ValueString())
	}
	return grantTypesSlice
}

func addOptionalOauthClientFields(ctx context.Context, addRequest *client.Client, plan oauthClientModel) error {
	addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	addRequest.Description = plan.Description.ValueStringPointer()
	addRequest.LogoUrl = plan.LogoUrl.ValueStringPointer()
	addRequest.RestrictToDefaultAccessTokenManager = plan.RestrictToDefaultAccessTokenManager.ValueBoolPointer()
	addRequest.ValidateUsingAllEligibleAtms = plan.ValidateUsingAllEligibleAtms.ValueBoolPointer()
	addRequest.RefreshRolling = plan.RefreshRolling.ValueStringPointer()
	addRequest.RefreshTokenRollingIntervalType = plan.RefreshTokenRollingIntervalType.ValueStringPointer()
	addRequest.RefreshTokenRollingInterval = plan.RefreshTokenRollingInterval.ValueInt64Pointer()
	addRequest.RefreshTokenRollingIntervalTimeUnit = plan.RefreshTokenRollingIntervalTimeUnit.ValueStringPointer()
	addRequest.PersistentGrantExpirationType = plan.PersistentGrantExpirationType.ValueStringPointer()
	addRequest.PersistentGrantExpirationTime = plan.PersistentGrantExpirationTime.ValueInt64Pointer()
	addRequest.PersistentGrantExpirationTimeUnit = plan.PersistentGrantExpirationTimeUnit.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeoutType = plan.PersistentGrantIdleTimeoutType.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeout = plan.PersistentGrantIdleTimeout.ValueInt64Pointer()
	addRequest.PersistentGrantIdleTimeoutTimeUnit = plan.PersistentGrantIdleTimeoutTimeUnit.ValueStringPointer()
	addRequest.PersistentGrantReuseType = plan.PersistentGrantReuseType.ValueStringPointer()
	addRequest.AllowAuthenticationApiInit = plan.AllowAuthenticationApiInit.ValueBoolPointer()
	addRequest.EnableCookielessAuthenticationApi = plan.EnableCookielessAuthenticationApi.ValueBoolPointer()
	addRequest.BypassApprovalPage = plan.BypassApprovalPage.ValueBoolPointer()
	addRequest.RequirePushedAuthorizationRequests = plan.RequirePushedAuthorizationRequests.ValueBoolPointer()
	addRequest.RequireJwtSecuredAuthorizationResponseMode = plan.RequireJwtSecuredAuthorizationResponseMode.ValueBoolPointer()
	addRequest.RequireSignedRequests = plan.RequireSignedRequests.ValueBoolPointer()
	addRequest.RequestObjectSigningAlgorithm = plan.RequestObjectSigningAlgorithm.ValueStringPointer()
	addRequest.DeviceFlowSettingType = plan.DeviceFlowSettingType.ValueStringPointer()
	addRequest.UserAuthorizationUrlOverride = plan.UserAuthorizationUrlOverride.ValueStringPointer()
	addRequest.PendingAuthorizationTimeoutOverride = plan.PendingAuthorizationTimeoutOverride.ValueInt64Pointer()
	addRequest.DevicePollingIntervalOverride = plan.DevicePollingIntervalOverride.ValueInt64Pointer()
	addRequest.BypassActivationCodeConfirmationOverride = plan.BypassActivationCodeConfirmationOverride.ValueBoolPointer()
	addRequest.RequireProofKeyForCodeExchange = plan.RequireProofKeyForCodeExchange.ValueBoolPointer()
	addRequest.CibaDeliveryMode = plan.CibaDeliveryMode.ValueStringPointer()
	addRequest.CibaNotificationEndpoint = plan.CibaNotificationEndpoint.ValueStringPointer()
	addRequest.CibaPollingInterval = plan.CibaPollingInterval.ValueInt64Pointer()
	addRequest.CibaRequireSignedRequests = plan.CibaRequireSignedRequests.ValueBoolPointer()
	addRequest.CibaRequestObjectSigningAlgorithm = plan.CibaRequestObjectSigningAlgorithm.ValueStringPointer()
	addRequest.CibaUserCodeSupported = plan.CibaUserCodeSupported.ValueBoolPointer()
	addRequest.RefreshTokenRollingGracePeriodType = plan.RefreshTokenRollingGracePeriodType.ValueStringPointer()
	addRequest.RefreshTokenRollingGracePeriod = plan.RefreshTokenRollingGracePeriod.ValueInt64Pointer()
	addRequest.ClientSecretRetentionPeriodType = plan.ClientSecretRetentionPeriodType.ValueStringPointer()
	addRequest.ClientSecretRetentionPeriod = plan.ClientSecretRetentionPeriod.ValueInt64Pointer()
	addRequest.TokenIntrospectionSigningAlgorithm = plan.TokenIntrospectionSigningAlgorithm.ValueStringPointer()
	addRequest.TokenIntrospectionEncryptionAlgorithm = plan.TokenIntrospectionEncryptionAlgorithm.ValueStringPointer()
	addRequest.TokenIntrospectionContentEncryptionAlgorithm = plan.TokenIntrospectionContentEncryptionAlgorithm.ValueStringPointer()
	addRequest.JwtSecuredAuthorizationResponseModeSigningAlgorithm = plan.JwtSecuredAuthorizationResponseModeSigningAlgorithm.ValueStringPointer()
	addRequest.JwtSecuredAuthorizationResponseModeEncryptionAlgorithm = plan.JwtSecuredAuthorizationResponseModeEncryptionAlgorithm.ValueStringPointer()
	addRequest.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm = plan.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm.ValueStringPointer()
	addRequest.RequireDpop = plan.RequireDpop.ValueBoolPointer()
	addRequest.RestrictScopes = plan.RestrictScopes.ValueBoolPointer()
	addRequest.RequireOfflineAccessScopeToIssueRefreshTokens = plan.RequireOfflineAccessScopeToIssueRefreshTokens.ValueStringPointer()
	addRequest.OfflineAccessRequireConsentPrompt = plan.OfflineAccessRequireConsentPrompt.ValueStringPointer()
	addRequest.LockoutMaxMaliciousActions = plan.LockoutMaxMaliciousActions.ValueInt64Pointer()
	addRequest.LockoutMaxMaliciousActionsType = plan.LockoutMaxMaliciousActionsType.ValueStringPointer()

	// addRequest.RestrictedScopes
	var restrictedScopes []string
	plan.RestrictedScopes.ElementsAs(ctx, &restrictedScopes, false)
	addRequest.RestrictedScopes = restrictedScopes

	if internaltypes.IsDefined(plan.ExclusiveScopes) {
		var slice []string
		plan.ExclusiveScopes.ElementsAs(ctx, &slice, false)
		addRequest.ExclusiveScopes = slice
	}

	if internaltypes.IsDefined(plan.RedirectUris) {
		var slice []string
		plan.RedirectUris.ElementsAs(ctx, &slice, false)
		addRequest.RedirectUris = slice
	}

	if internaltypes.IsDefined(plan.DefaultAccessTokenManagerRef) {
		addRequest.DefaultAccessTokenManagerRef = client.NewResourceLinkWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.DefaultAccessTokenManagerRef, false)), addRequest.DefaultAccessTokenManagerRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.PersistentGrantReuseGrantTypes) {
		var slice []string
		plan.PersistentGrantReuseGrantTypes.ElementsAs(ctx, &slice, false)
		addRequest.PersistentGrantReuseGrantTypes = slice
	}

	if internaltypes.IsDefined(plan.AuthorizationDetailTypes) {
		var slice []string
		plan.AuthorizationDetailTypes.ElementsAs(ctx, &slice, false)
		addRequest.AuthorizationDetailTypes = slice
	}

	if internaltypes.IsDefined(plan.RestrictedResponseTypes) {
		var slice []string
		plan.RestrictedResponseTypes.ElementsAs(ctx, &slice, false)
		addRequest.RestrictedResponseTypes = slice
	}

	if internaltypes.IsDefined(plan.OidcPolicy) {

		addRequest.OidcPolicy = &client.ClientOIDCPolicy{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.OidcPolicy, true)), addRequest.OidcPolicy)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ClientAuth) {
		addRequest.ClientAuth = &client.ClientAuth{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ClientAuth, true)), addRequest.ClientAuth)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.JwksSettings) {
		addRequest.JwksSettings = &client.JwksSettings{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.JwksSettings, false)), addRequest.JwksSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ExtendedParameters) {
		mapValue := map[string]client.ParameterValues{}
		addRequest.ExtendedParameters = &mapValue
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ExtendedParameters, false)), addRequest.ExtendedParameters)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.RequestPolicyRef) {
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RequestPolicyRef, false)), addRequest.RequestPolicyRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.TokenExchangeProcessorPolicyRef) {
		addRequest.TokenExchangeProcessorPolicyRef = &client.ResourceLink{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.TokenExchangeProcessorPolicyRef, false)), addRequest.TokenExchangeProcessorPolicyRef)
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *oauthClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthClientModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthClient := client.NewClient(plan.ClientId.ValueString(), grantTypes(plan.GrantTypes), plan.Name.ValueString())
	err := addOptionalOauthClientFields(ctx, createOauthClient, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for OAuth Client: "+err.Error())
		return
	}

	apiCreateOauthClient := r.apiClient.OauthClientsAPI.CreateOauthClient(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthClient = apiCreateOauthClient.Body(*createOauthClient)
	oauthClientResponse, httpResp, err := r.apiClient.OauthClientsAPI.CreateOauthClientExecute(apiCreateOauthClient)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Client", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state oauthClientModel

	diags = readOauthClientResponse(ctx, oauthClientResponse, &plan, &state, r.providerConfig.ProductVersion, false)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state oauthClientModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthClient, httpResp, err := r.apiClient.OauthClientsAPI.GetOauthClientById(config.AuthContext(ctx, r.providerConfig), state.ClientId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Client", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting the  OAuth Client", err, httpResp, &customId)
		}
		return
	}

	// Read the response into the state
	diags = readOauthClientResponse(ctx, apiReadOauthClient, &state, &state, r.providerConfig.ProductVersion, isImportRead)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOauthClient := r.apiClient.OauthClientsAPI.UpdateOauthClient(config.AuthContext(ctx, r.providerConfig), plan.ClientId.ValueString())
	createUpdateRequest := client.NewClient(plan.ClientId.ValueString(), grantTypes(plan.GrantTypes), plan.Name.ValueString())
	err := addOptionalOauthClientFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for the OAuth Client: "+err.Error())
		return
	}

	updateOauthClient = updateOauthClient.Body(*createUpdateRequest)
	updateOauthClientResponse, httpResp, err := r.apiClient.OauthClientsAPI.UpdateOauthClientExecute(updateOauthClient)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the OAuth Client", err, httpResp, &customId)
		return
	}

	// Read the response
	var state oauthClientModel
	diags = readOauthClientResponse(ctx, updateOauthClientResponse, &plan, &state, r.providerConfig.ProductVersion, false)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthClientsAPI.DeleteOauthClient(config.AuthContext(ctx, r.providerConfig), state.ClientId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Client", err, httpResp, &customId)
	}
}

func (r *oauthClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("client_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
