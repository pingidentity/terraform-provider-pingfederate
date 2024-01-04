package oauthclient

import (
	"context"
	"encoding/json"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
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
	emptyStringSet, _           = types.SetValue(types.StringType, []attr.Value{})
	oidcPolicyDefaultObj, _     = types.ObjectValue(oidcPolicyAttrType, oidcPolicyDefaultAttrValue)
	secondarySecretsEmptySet, _ = types.SetValue(types.ObjectType{AttrTypes: secondarySecretsAttrType}, []attr.Value{})
	clientAuthDefaultObj, _     = types.ObjectValue(clientAuthAttrType, clientAuthDefaultAttrValue)
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
				Description: "A unique identifier the client provides to the Resource Server to identify itself. This identifier is included with every request the client makes. For PUT requests, this field is optional and it will be overridden by the 'id' parameter of the PUT request.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Specifies whether the client is enabled. The default value is true.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"redirect_uris": schema.SetAttribute{
				Description: "URIs to which the OAuth AS may redirect the resource owner's user agent after authorization is obtained. A redirection URI is used with the Authorization Code and Implicit grant types. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible.For example: https://.company.com/ Important: If more than one URI is added or if a single URI uses wildcards, then Authorization Code grant and token requests must contain a specific matching redirect uri parameter.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(emptyStringSet),
			},
			"grant_types": schema.SetAttribute{
				Description: "The grant types allowed for this client. The EXTENSION grant type applies to SAML/JWT assertion grants.",
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
			},
			"description": schema.StringAttribute{
				Description: "A description of what the client application does. This description appears when the user is prompted for authorization.",
				Optional:    true,
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
				Description: "Determines whether the client is restricted to using only its default access token manager. The default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"validate_using_all_eligible_atms": schema.BoolAttribute{
				Description: "Validates token using all eligible access token managers for the client. This setting is ignored if 'restrictToDefaultAccessTokenManager' is set to true.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"refresh_rolling": schema.StringAttribute{
				Description: "Use ROLL or DONT_ROLL to override the Roll Refresh Token Values setting on the Authorization Server Settings. SERVER_DEFAULT will default to the Roll Refresh Token Values setting on the Authorization Server Setting screen. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("ROLL", "DONT_ROLL", "SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_interval_type": schema.StringAttribute{
				Description: "Use OVERRIDE_SERVER_DEFAULT to override the Refresh Token Rolling Interval value on the Authorization Server Settings. SERVER_DEFAULT will default to the Refresh Token Rolling Interval value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("OVERRIDE_SERVER_DEFAULT", "SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens, in hours. This value will override the Refresh Token Rolling Interval Value on the Authorization Server Settings.",
				Optional:    true,
			},
			"persistent_grant_expiration_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Lifetime set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
			},
			"persistent_grant_expiration_time": schema.Int64Attribute{
				Description: "The persistent grant expiration time. -1 indicates an indefinite amount of time.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"persistent_grant_expiration_time_unit": schema.StringAttribute{
				Description: "The persistent grant expiration time unit.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"persistent_grant_idle_timeout_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Idle Timeout set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DAYS"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"persistent_grant_reuse_type": schema.StringAttribute{
				Description: "Allows and administrator to override the Reuse Existing Persistent Access Grants for Grant Types set globally for OAuth AS. Defaults to SERVER_DEFAULT",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("SERVER_DEFAULT", "OVERRIDE_SERVER_DEFAULT"),
				},
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. This value will override the Reuse Existing Persistent Access Grants for Grant Types on the Authorization Server Settings. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("IMPLICIT",
							"AUTHORIZATION_CODE",
							"RESOURCE_OWNER_CREDENTIALS",
							"CLIENT_CREDENTIALS",
							"REFRESH_TOKEN",
							"EXTENSION, DEVICE_CODE",
							"ACCESS_TOKEN_VALIDATION",
							"CIBA",
							"TOKEN_EXCHANGE",
						),
					),
				},
			},
			"allow_authentication_api_init": schema.BoolAttribute{
				Description: "Set to true to allow this client to initiate the authentication API redirectless flow.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"bypass_approval_page": schema.BoolAttribute{
				Description: "Use this setting, for example, when you want to deploy a trusted application and authenticate end users via an IdP adapter or IdP connection.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"restrict_scopes": schema.BoolAttribute{
				Description: "Restricts this client's access to specific scopes.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
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
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"authorization_detail_types": schema.SetAttribute{
				Description: "The authorization detail types available for this client.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"restricted_response_types": schema.SetAttribute{
				Description: "The response types allowed for this client. If omitted all response types are available to the client.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(emptyStringSet),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"require_pushed_authorization_requests": schema.BoolAttribute{
				Description: "Determines whether pushed authorization requests are required when initiating an authorization request. The default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_jwt_secured_authorization_response_mode": schema.BoolAttribute{
				Description: "Determines whether JWT secured authorization response mode is required when initiating an authorization request. The default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether JWT Secured authorization response mode is required when initiating an authorization request. The default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present\nRS256 - RSA using SHA-256\n\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
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
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm required for the ID Token.\nNONE - No signing algorithm\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
						Description:         "The JSON Web Signature [JWS] algorithm required for the ID Token. NONE - No signing algorithm, HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256. RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
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
						MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
						Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
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
						MarkdownDescription: "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
						Description:         "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
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
						Description: "Determines whether this client is allowed to access the Session Revocation API.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"grant_access_session_session_management_api": schema.BoolAttribute{
						Description: "Determines whether this client is allowed to access the Session Management API.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"ping_access_logout_capable": schema.BoolAttribute{
						Description: "Set this value to true if you wish to enable client application logout, and the client is PingAccess, or its logout endpoints follow the PingAccess path convention",
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
						Description: "Determines whether the subject identifier type is pairwise.",
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
						Description: "The logout mode for this client. The default is 'NONE'. Supported in PF version 11.3 or later.",
						Optional:    true,
						Computed:    true,
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
						Description: "The back-channel logout URI for this client. Supported in PF version 11.3 or later.",
						Optional:    true,
					},
					"post_logout_redirect_uris": schema.SetAttribute{
						Description: "URIs to which the OIDC OP may redirect the resource owner's user agent after RP-initiated logout has completed. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible. Supported in PF version 12.0 or later.",
						Optional:    true,
						ElementType: types.StringType,
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
						Description: "Client authentication type. The required field for type SECRET is secret.	The required fields for type CERTIFICATE are clientCertIssuerDn and clientCertSubjectDn. The required field for type PRIVATE_KEY_JWT is: either jwks or jwksUrl.",
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
						Description: "Client secret for Basic Authentication. To update the client secret, specify the plaintext value in this field. This field will not be populated for GET requests.",
						Optional:    true,
					},
					"secondary_secrets": schema.SetNestedAttribute{
						Description: "The list of secondary client secrets that are temporarily retained.",
						Computed:    true,
						Optional:    true,
						Default:     setdefault.StaticValue(secondarySecretsEmptySet),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"secret": schema.StringAttribute{
									Description: "Secondary client secret for Basic Authentication. To update the secondary client secret, specify the plaintext value in this field. This field will not be populated for GET requests.",
									Required:    true,
								},
								"expiry_time": schema.StringAttribute{
									Description: "The expiry time of the secondary secret.",
									Required:    true,
								},
							},
						},
					},
					"client_cert_issuer_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Issuer DN.",
						Optional:    true,
					},
					"client_cert_subject_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Subject DN.",
						Optional:    true,
					},
					"enforce_replay_prevention": schema.BoolAttribute{
						Description: "Enforce replay prevention on JSON Web Tokens. This field is applicable only for Private Key JWT Client and Client Secret JWT Authentication.",
						Optional:    true,
					},
					"token_endpoint_auth_signing_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT and Client Secret JWT Client Authentication. All asymmetric signing algorithms are allowed for Private Key JWT if value is not present. All symmetric signing algorithms are allowed for Client Secret JWT if value is not present\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512.",
						Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT and Client Secret JWT Client Authentication. All asymmetric signing algorithms are allowed for Private Key JWT if value is not present. All symmetric signing algorithms are allowed for Client Secret JWT if value is not present RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11. HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512.",
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
						Description: "JSON Web Key Set (JWKS) URL of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    true,
					},
					"jwks": schema.StringAttribute{
						Description: "JSON Web Key Set (JWKS) document of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    true,
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
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"device_flow_setting_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Device Authorization Settings set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("SERVER_DEFAULT", "OVERRIDE_SERVER_DEFAULT"),
				},
			},
			"user_authorization_url_override": schema.StringAttribute{
				Description: "The URL used as 'verification_url' and 'verification_url_complete' values in a Device Authorization request. This property overrides the 'userAuthorizationUrl' value present in Authorization Server Settings.",
				Optional:    true,
			},
			"pending_authorization_timeout_override": schema.Int64Attribute{
				Description: "The 'device_code' and 'user_code' timeout, in seconds. This overrides the 'pendingAuthorizationTimeout' value present in Authorization Server Settings.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"device_polling_interval_override": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds. This overrides the 'devicePollingInterval' value present in Authorization Server Settings.",
				Optional:    true,
			},
			"bypass_activation_code_confirmation_override": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device. This overrides the 'bypassUseCodeConfirmation' value present in Authorization Server Settings.",
				Optional:    true,
			},
			"require_proof_key_for_code_exchange": schema.BoolAttribute{
				Description: "Determines whether Proof Key for Code Exchange (PKCE) is required for this client.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ciba_delivery_mode": schema.StringAttribute{
				Description: "The token delivery mode for the client. The default value is 'POLL'.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("POLL", "PING"),
				},
			},
			"ciba_notification_endpoint": schema.StringAttribute{
				Description: "The endpoint the OP will call after a successful or failed end-user authentication.",
				Optional:    true,
			},
			"ciba_polling_interval": schema.Int64Attribute{
				Description: "The minimum amount of time in seconds that the Client must wait between polling requests to the token endpoint. The default is 0 seconds.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(3600),
				},
			},
			"ciba_require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether CIBA signed requests are required for this client.",
				Optional:    true,
			},
			"ciba_request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.",
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
				Description: "When specified, it overrides the global Refresh Token Grace Period defined in the Authorization Server Settings. The default value is SERVER_DEFAULT",
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
				Description: "Use OVERRIDE_SERVER_DEFAULT to override the Client Secret Retention Period value on the Authorization Server Settings. SERVER_DEFAULT will default to the Client Secret Retention Period value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("SERVER_DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("OVERRIDE_SERVER_DEFAULT", "SERVER_DEFAULT"),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention. This value will override the Client Secret Retention Period value on the Authorization Server Settings.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"client_secret_changed_time": schema.StringAttribute{
				Description: "The time at which the client secret was last changed. This property is read only and is ignored on PUT and POST requests.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_introspection_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response.\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response. HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
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
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
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
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
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
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response.\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response. HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11",
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
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
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
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
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
				MarkdownDescription: "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client. Supported in PF version 11.3 or later.",
				Description:         "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client. Supported in PF version 11.3 or later.",
				Optional:            true,
				Computed:            true,
				// Default set when appropriate in ModifyPlan before
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
	var model oauthClientModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	// Persistent Grant Expiration Validation
	if (internaltypes.IsDefined(model.PersistentGrantExpirationTime) || internaltypes.IsDefined(model.PersistentGrantExpirationTimeUnit)) && model.PersistentGrantExpirationType.ValueString() != "OVERRIDE_SERVER_DEFAULT" {
		resp.Diagnostics.AddError("persistent_grant_expiration_type must be configured to \"OVERRIDE_SERVER_DEFAULT\" to modify the other persistent_grant_expiration values.", "")
	}

	// Refresh Token Rolling Validation
	if (model.RefreshTokenRollingIntervalType.ValueString() == "OVERRIDE_SERVER_DEFAULT") != internaltypes.IsDefined(model.RefreshTokenRollingInterval) {
		resp.Diagnostics.AddError("refresh_token_rolling_interval must be configured when refresh_token_rolling_interval_type is \"OVERRIDE_SERVER_DEFAULT\".", "")
	}

	//  Client Auth Defined
	var clientAuthAttributes map[string]attr.Value
	clientAuthDefined := internaltypes.IsDefined(model.ClientAuth)
	if clientAuthDefined {
		clientAuthAttributes = model.ClientAuth.Attributes()
		if internaltypes.IsDefined(clientAuthAttributes["type"]) {
			clientAuthType := clientAuthAttributes["type"].(types.String).ValueString()
			switch clientAuthType {
			case "PRIVATE_KEY_JWT":
				if !internaltypes.IsNonEmptyObj(model.JwksSettings) {
					resp.Diagnostics.AddError("jwks_settings must be defined when client_auth is configured to \"PRIVATE_KEY_JWT\".", "")
				}
			case "CERTIFICATE":
				if !internaltypes.IsDefined(clientAuthAttributes["client_cert_subject_dn"]) || !internaltypes.IsDefined(clientAuthAttributes["client_cert_issuer_dn"]) {
					resp.Diagnostics.AddError("client_cert_subject_dn and client_cert_issuer_dn must be defined when client_auth is configured to \"CERTIFICATE\".", "")
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
			if clientAuthDefined {
				clientAuthType := clientAuthAttributes["type"].(types.String).ValueString()
				clientAuthSecret := clientAuthAttributes["secret"].(types.String).ValueString()
				if clientAuthType != "NONE" {
					resp.Diagnostics.AddError("client_auth.type must be set to \"SECRET\" when \"CLIENT_CREDENTIALS\" is included in grant_types.", "")
				}
				if clientAuthSecret == "" {
					resp.Diagnostics.AddError("client_auth.secret cannot be empty when \"CLIENT_CREDENTIALS\" is included in grant_types.", "")
				}
			} else if !clientAuthDefined {
				resp.Diagnostics.AddError("client_auth must be defined when \"CLIENT_CREDENTIALS\" is included in grant_types.", "")
			}
		}
		if grantTypeVal == "CIBA" {
			hasCibaGrantType = true
		}
	}

	// CIBA Validation
	if !hasCibaGrantType && (internaltypes.IsDefined(model.CibaDeliveryMode) ||
		internaltypes.IsDefined(model.CibaNotificationEndpoint) ||
		internaltypes.IsDefined(model.CibaPollingInterval) ||
		internaltypes.IsDefined(model.CibaRequireSignedRequests) ||
		internaltypes.IsDefined(model.CibaRequestObjectSigningAlgorithm) ||
		internaltypes.IsDefined(model.CibaUserCodeSupported)) {
		resp.Diagnostics.AddError("ciba attributes can only be configured when \"CIBA\" is included in grant_types.", "")
	}
	if hasCibaGrantType && (model.CibaDeliveryMode.ValueString() == "PING" && !internaltypes.IsDefined(model.CibaNotificationEndpoint)) {
		resp.Diagnostics.AddError("ciba_notification_endpoint must be defined when ciba_delivery_mode is \"PING\".", "")
	}

	// Client Auth Validation
	// ID Token Signing Algorithm Validation when client_auth is not defined
	if !internaltypes.IsDefined(model.ClientAuth) {
		var algorithmAttributeSet []string
		if internaltypes.IsDefined(model.OidcPolicy) && model.OidcPolicy.Attributes()["id_token_signing_algorithm"] != nil {
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
				resp.Diagnostics.AddError("client_auth must be defined when using the \"HS256\" signing algorithm", "")
			}
		}

		if internaltypes.IsDefined(model.TokenIntrospectionEncryptionAlgorithm) {
			resp.Diagnostics.AddError("client_auth must be configured when token_introspection_encryption_algorithm is configured.", "")
		}
	}

	// Restrict Scopes Validation
	if len(model.RestrictedScopes.Elements()) > 0 && !model.RestrictScopes.ValueBool() {
		resp.Diagnostics.AddError("restrict_scopes must be set to true to configure restricted_scopes.", "")
	}

	// OIDC Policy Validation
	if internaltypes.IsDefined(model.OidcPolicy) {
		oidcPolicy := model.OidcPolicy.Attributes()
		pairwiseIdentifierUserType := oidcPolicy["pairwise_identifier_user_type"]
		oidcPolicySectorIdentifierUri := oidcPolicy["sector_identifier_uri"]
		if (pairwiseIdentifierUserType != nil && !pairwiseIdentifierUserType.(types.Bool).ValueBool()) && internaltypes.IsDefined(oidcPolicySectorIdentifierUri) {
			resp.Diagnostics.AddError("sector_identifier_uri can only be configured when pairwise_identifier_user_type is set to true.", "")
		}
	}

	// JWKS Settings Validation
	if !internaltypes.IsDefined(model.JwksSettings) {
		if internaltypes.IsDefined(model.TokenIntrospectionEncryptionAlgorithm) {
			resp.Diagnostics.AddError("token_introspection_encryption_algorithm must not be configured when jwks_settings is not configured.", "")
		}
		if model.RequireSignedRequests.ValueBool() {
			resp.Diagnostics.AddError("require_signed_requests must be false when jwks_settings is not configured.", "")
		}
	}
}

func (r *oauthClientResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 11.3 and 12.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast113 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	var plan, state oauthClientModel
	var diags diag.Diagnostics
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	planModified := false
	// If require_dpop is set prior to PF version 11.3, throw an error
	if !pfVersionAtLeast113 {
		if internaltypes.IsDefined(plan.RequireDpop) {
			version.AddUnsupportedAttributeError("require_dpop",
				r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
		} else if plan.RequireDpop.IsUnknown() {
			// Ensure require_dpop is not unknown for older versions of PF, so that it gets passed in as nil rather than false.
			// Passing it in as false would break older versions of PF, since it is an unrecognized property.
			plan.RequireDpop = types.BoolNull()
			planModified = true
		}
	} else if plan.RequireDpop.IsUnknown() {
		// Set a default of false if the PF version is new enough
		plan.RequireDpop = types.BoolValue(false)
		planModified = true
	}
	if internaltypes.IsDefined(plan.OidcPolicy) {
		planOidcPolicyAttrs := plan.OidcPolicy.Attributes()
		// If oidc_policy.logout_mode is set prior to PF version 11.3, throw an error. Otherwise, set the PF default.
		planLogoutMode := planOidcPolicyAttrs["logout_mode"].(types.String)
		planBackChannelLogoutUri := planOidcPolicyAttrs["back_channel_logout_uri"].(types.String)
		if !pfVersionAtLeast113 {
			if internaltypes.IsDefined(planLogoutMode) {
				version.AddUnsupportedAttributeError("oidc_policy.logout_mode",
					r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
			} else if planLogoutMode.IsUnknown() {
				// Ensure logout_mode is not unknown for older versions of PF
				planOidcPolicyAttrs["logout_mode"] = types.StringNull()
				plan.OidcPolicy, diags = types.ObjectValue(plan.OidcPolicy.AttributeTypes(ctx), planOidcPolicyAttrs)
				resp.Diagnostics.Append(diags...)
				planModified = true
			}
			if internaltypes.IsDefined(planBackChannelLogoutUri) {
				version.AddUnsupportedAttributeError("oidc_policy.back_channel_logout_uri",
					r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
			}
		} else if planLogoutMode.IsUnknown() {
			// Set a default logout_mode if the PF version is new enough
			planOidcPolicyAttrs["logout_mode"] = types.StringValue("NONE")
			plan.OidcPolicy, diags = types.ObjectValue(plan.OidcPolicy.AttributeTypes(ctx), planOidcPolicyAttrs)
			resp.Diagnostics.Append(diags...)
			planModified = true
		}
		// If oidc_policy.post_logout_redirect_uris is set prior to PF version 12.0, throw an error.
		planPostLogoutRedirectUris := planOidcPolicyAttrs["post_logout_redirect_uris"].(types.Set)
		if !pfVersionAtLeast120 && internaltypes.IsDefined(planPostLogoutRedirectUris) {
			version.AddUnsupportedAttributeError("oidc_policy.post_logout_redirect_uris",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		}
	}

	// If the new plan doesn't match the state, invalidate any last-changed time values
	// See https://github.com/hashicorp/terraform-plugin-framework/issues/898 for some info on why this is needed
	req.Plan.Set(ctx, plan)
	if !req.Plan.Raw.Equal(req.State.Raw) {
		plan.ModificationDate = types.StringUnknown()
		plan.ClientSecretChangedTime = types.StringUnknown()
		planModified = true
	}

	if planModified {
		resp.Plan.Set(ctx, &plan)
	}
}

func readOauthClientResponse(ctx context.Context, r *client.Client, plan, state *oauthClientModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	diags = readOauthClientResponseCommon(ctx, r, state)

	// state.ClientAuth
	var clientAuthToState types.Object
	clientAuthFromPlan := plan.ClientAuth.Attributes()
	var secretToState basetypes.StringValue

	// state.ClientAuth.Secret
	secretVal := clientAuthFromPlan["secret"]
	if secretVal != nil && internaltypes.IsNonEmptyString(secretVal.(types.String)) {
		secretToState = types.StringValue(secretVal.(types.String).ValueString())
	} else {
		secretToState = types.StringNull()
	}

	// state.ClientAuth.Secret
	var secondarySecretsObjToState types.Set
	var secondarySecretsSetSlice []attr.Value
	secondarySecretsFromPlan := clientAuthFromPlan["secondary_secrets"]
	if secondarySecretsFromPlan != nil && len(secondarySecretsFromPlan.(types.Set).Elements()) > 0 {
		for _, secondarySecretsFromPlan := range clientAuthFromPlan["secondary_secrets"].(types.Set).Elements() {
			secondarySecretsAttrVal, respDiags := types.ObjectValueFrom(ctx, secondarySecretsAttrType, secondarySecretsFromPlan)
			diags.Append(respDiags...)
			secondarySecretsSetSlice = append(secondarySecretsSetSlice, secondarySecretsAttrVal)
		}
	}
	secondarySecretsObjToState, respDiags = types.SetValue(types.ObjectType{AttrTypes: secondarySecretsAttrType}, secondarySecretsSetSlice)
	diags.Append(respDiags...)

	// state.ClientAuth to state
	clientAuthAttrValue := map[string]attr.Value{}
	clientAuthAttrValue["type"] = types.StringPointerValue(r.ClientAuth.Type)
	clientAuthAttrValue["secret"] = secretToState
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
	addRequest.PersistentGrantExpirationType = plan.PersistentGrantExpirationType.ValueStringPointer()
	addRequest.PersistentGrantExpirationTime = plan.PersistentGrantExpirationTime.ValueInt64Pointer()
	addRequest.PersistentGrantExpirationTimeUnit = plan.PersistentGrantExpirationTimeUnit.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeoutType = plan.PersistentGrantIdleTimeoutType.ValueStringPointer()
	addRequest.PersistentGrantIdleTimeout = plan.PersistentGrantIdleTimeout.ValueInt64Pointer()
	addRequest.PersistentGrantIdleTimeoutTimeUnit = plan.PersistentGrantIdleTimeoutTimeUnit.ValueStringPointer()
	addRequest.PersistentGrantReuseType = plan.PersistentGrantReuseType.ValueStringPointer()
	addRequest.AllowAuthenticationApiInit = plan.AllowAuthenticationApiInit.ValueBoolPointer()
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

	if internaltypes.IsDefined(plan.RestrictScopes) {
		addRequest.RestrictScopes = plan.RestrictScopes.ValueBoolPointer()
		if *plan.RestrictScopes.ValueBoolPointer() && internaltypes.IsDefined(plan.RestrictedScopes) {
			var slice []string
			plan.RestrictedScopes.ElementsAs(ctx, &slice, false)
			addRequest.RestrictedScopes = slice
		}
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Client", err.Error())
		return
	}

	apiCreateOauthClient := r.apiClient.OauthClientsAPI.CreateOauthClient(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthClient = apiCreateOauthClient.Body(*createOauthClient)
	oauthClientResponse, httpResp, err := r.apiClient.OauthClientsAPI.CreateOauthClientExecute(apiCreateOauthClient)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Client", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthClientModel

	diags = readOauthClientResponse(ctx, oauthClientResponse, &plan, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthClientModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthClient, httpResp, err := r.apiClient.OauthClientsAPI.GetOauthClientById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ClientId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Client", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  OAuth Client", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readOauthClientResponse(ctx, apiReadOauthClient, &state, &state)
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

	updateOauthClient := r.apiClient.OauthClientsAPI.UpdateOauthClient(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ClientId.ValueString())
	createUpdateRequest := client.NewClient(plan.ClientId.ValueString(), grantTypes(plan.GrantTypes), plan.Name.ValueString())
	err := addOptionalOauthClientFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the OAuth Client", err.Error())
		return
	}

	updateOauthClient = updateOauthClient.Body(*createUpdateRequest)
	updateOauthClientResponse, httpResp, err := r.apiClient.OauthClientsAPI.UpdateOauthClientExecute(updateOauthClient)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the OAuth Client", err, httpResp)
		return
	}

	// Read the response
	var state oauthClientModel
	diags = readOauthClientResponse(ctx, updateOauthClientResponse, &plan, &state)
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
	httpResp, err := r.apiClient.OauthClientsAPI.DeleteOauthClient(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ClientId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Client", err, httpResp)
	}
}

func (r *oauthClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("client_id"), req, resp)
}
