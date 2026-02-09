// Copyright © 2026 Ping Identity Corporation

package oauthclient

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	resourcelinkdatasource "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthClientDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthClientDataSource{}
)

var (
	secondarySecretsDataSourceAttrType = map[string]attr.Type{
		"encrypted_secret": types.StringType,
		"expiry_time":      types.StringType,
	}

	clientAuthDataSourceAttrType = map[string]attr.Type{
		"type":                                  types.StringType,
		"encrypted_secret":                      types.StringType,
		"secondary_secrets":                     types.ListType{ElemType: types.ObjectType{AttrTypes: secondarySecretsDataSourceAttrType}},
		"client_cert_issuer_dn":                 types.StringType,
		"client_cert_subject_dn":                types.StringType,
		"enforce_replay_prevention":             types.BoolType,
		"token_endpoint_auth_signing_algorithm": types.StringType,
	}
)

// Create a Administrative Account data source
func OauthClientDataSource() datasource.DataSource {
	return &oauthClientDataSource{}
}

// oauthClientDataSource is the datasource implementation.
type oauthClientDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *oauthClientDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes an OAuth Client.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "A unique identifier the client provides to the Resource Server to identify itself. This identifier is included with every request the client makes. For PUT requests, this field is optional and it will be overridden by the 'id' parameter of the PUT request.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Specifies whether the client is enabled. The default value is true.",
				Optional:    false,
				Computed:    true,
			},
			"redirect_uris": schema.SetAttribute{
				Description: "URIs to which the OAuth AS may redirect the resource owner's user agent after authorization is obtained. A redirection URI is used with the Authorization Code and Implicit grant types. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible.For example: https://.company.com/ Important: If more than one URI is added or if a single URI uses wildcards, then Authorization Code grant and token requests must contain a specific matching redirect uri parameter.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"grant_types": schema.SetAttribute{
				Description: "The grant types allowed for this client. The EXTENSION grant type applies to SAML/JWT assertion grants.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"name": schema.StringAttribute{
				Description: "A descriptive name for the client instance. This name appears when the user is prompted for authorization.",
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of what the client application does. This description appears when the user is prompted for authorization.",
				Optional:    false,
				Computed:    true,
			},
			"modification_date": schema.StringAttribute{
				Description: "The time at which the client was last changed. This property is read only.",
				Optional:    false,
				Computed:    true,
			},
			"creation_date": schema.StringAttribute{
				Description: "The time at which the client was created. This property is read only.",
				Optional:    false,
				Computed:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: "The location of the logo used on user-facing OAuth grant authorization and revocation pages.",
				Optional:    false,
				Computed:    true,
			},
			"default_access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "The default access token manager for this client.",
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
			},
			"restrict_to_default_access_token_manager": schema.BoolAttribute{
				Description: "Determines whether the client is restricted to using only its default access token manager. The default is false.",
				Optional:    false,
				Computed:    true,
			},
			"validate_using_all_eligible_atms": schema.BoolAttribute{
				Description: "Validates token using all eligible access token managers for the client. This setting is ignored if 'restrictToDefaultAccessTokenManager' is set to true.",
				Optional:    false,
				Computed:    true,
			},
			"refresh_rolling": schema.StringAttribute{
				Description: "Use ROLL or DONT_ROLL to override the Roll Refresh Token Values setting on the Authorization Server Settings. SERVER_DEFAULT will default to the Roll Refresh Token Values setting on the Authorization Server Setting screen. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_interval_type": schema.StringAttribute{
				Description: "Use OVERRIDE_SERVER_DEFAULT to override the Refresh Token Rolling Interval value on the Authorization Server Settings. SERVER_DEFAULT will default to the Refresh Token Rolling Interval value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens. This value will override the Refresh Token Rolling Interval Value on the Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_interval_time_unit": schema.StringAttribute{
				Description: "The refresh token rolling interval time unit. Defaults to HOURS.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_expiration_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Lifetime set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_expiration_time": schema.Int64Attribute{
				Description: "The persistent grant expiration time. -1 indicates an indefinite amount of time.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_expiration_time_unit": schema.StringAttribute{
				Description: "The persistent grant expiration time unit.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Persistent Grant Idle Timeout set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit.",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_reuse_type": schema.StringAttribute{
				Description: "Allows and administrator to override the Reuse Existing Persistent Access Grants for Grant Types set globally for OAuth AS. Defaults to SERVER_DEFAULT",
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. This value will override the Reuse Existing Persistent Access Grants for Grant Types on the Authorization Server Settings. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allow_authentication_api_init": schema.BoolAttribute{
				Description: "Set to true to allow this client to initiate the authentication API redirectless flow.",
				Optional:    false,
				Computed:    true,
			},
			"enable_cookieless_authentication_api": schema.BoolAttribute{
				Description: "Set to true to allow the authentication API redirectless flow to function without requiring any cookies.",
				Optional:    false,
				Computed:    true,
			},
			"bypass_approval_page": schema.BoolAttribute{
				Description: "Use this setting, for example, when you want to deploy a trusted application and authenticate end users via an IdP adapter or IdP connection.",
				Optional:    false,
				Computed:    true,
			},
			"restrict_scopes": schema.BoolAttribute{
				Description: "Restricts this client's access to specific scopes.",
				Optional:    false,
				Computed:    true,
			},
			"restricted_scopes": schema.SetAttribute{
				Description: "The scopes available for this client.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclusive_scopes": schema.SetAttribute{
				Description: "The exclusive scopes available for this client.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"authorization_detail_types": schema.SetAttribute{
				Description: "The authorization detail types available for this client.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"restricted_response_types": schema.SetAttribute{
				Description: "The response types allowed for this client. If omitted all response types are available to the client.",
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_pushed_authorization_requests": schema.BoolAttribute{
				Description: "Determines whether pushed authorization requests are required when initiating an authorization request. The default is false.",
				Optional:    false,
				Computed:    true,
			},
			"require_jwt_secured_authorization_response_mode": schema.BoolAttribute{
				Description: "Determines whether JWT secured authorization response mode is required when initiating an authorization request. The default is false.",
				Optional:    false,
				Computed:    true,
			},
			"require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether JWT Secured authorization response mode is required when initiating an authorization request. The default is false.",
				Optional:    false,
				Computed:    true,
			},
			"request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present\nRS256 - RSA using SHA-256\n\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Optional:            false,
				Computed:            true,
			},
			"oidc_policy": schema.SingleNestedAttribute{
				Description: "Open ID Connect Policy settings. This is included in the message only when OIDC is enabled.",
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id_token_signing_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm required for the ID Token.\nNONE - No signing algorithm\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
						Description:         "The JSON Web Signature [JWS] algorithm required for the ID Token. NONE - No signing algorithm, HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256. RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
						Optional:            false,
						Computed:            true,
					},
					"id_token_encryption_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
						Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
						Optional:            false,
						Computed:            true,
					},
					"id_token_content_encryption_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
						Description:         "The JSON Web Encryption [JWE] content encryption algorithm for the ID Token. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
						Optional:            false,
						Computed:            true,
					},
					"policy_group": schema.SingleNestedAttribute{
						Description: "The Open ID Connect policy. A null value will represent the default policy group.",
						Optional:    false,
						Computed:    true,
						Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
					},
					"grant_access_session_revocation_api": schema.BoolAttribute{
						Description: "Determines whether this client is allowed to access the Session Revocation API.",
						Optional:    false,
						Computed:    true,
					},
					"grant_access_session_session_management_api": schema.BoolAttribute{
						Description: "Determines whether this client is allowed to access the Session Management API.",
						Optional:    false,
						Computed:    true,
					},
					"ping_access_logout_capable": schema.BoolAttribute{
						Description: "Set this value to true if you wish to enable client application logout, and the client is PingAccess, or its logout endpoints follow the PingAccess path convention",
						Optional:    false,
						Computed:    true,
					},
					"logout_uris": schema.SetAttribute{
						Description: "A list of client logout URI's which will be invoked when a user logs out through one of PingFederate's SLO endpoints.",
						ElementType: types.StringType,
						Optional:    false,
						Computed:    true,
					},
					"pairwise_identifier_user_type": schema.BoolAttribute{
						Description: "Determines whether the subject identifier type is pairwise.",
						Optional:    false,
						Computed:    true,
					},
					"sector_identifier_uri": schema.StringAttribute{
						Description: "The URI references a file with a single JSON array of Redirect URI and JWKS URL values.",
						Optional:    false,
						Computed:    true,
					},
					"logout_mode": schema.StringAttribute{
						Description: "The logout mode for this client. The default is 'NONE'.",
						Optional:    false,
						Computed:    true,
					},
					"back_channel_logout_uri": schema.StringAttribute{
						Description: "The back-channel logout URI for this client.",
						Optional:    false,
						Computed:    true,
					},
					"post_logout_redirect_uris": schema.SetAttribute{
						Description: "URIs to which the OIDC OP may redirect the resource owner's user agent after RP-initiated logout has completed. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible. Supported in PF version 12.0 or later.",
						Optional:    false,
						Computed:    true,
						ElementType: types.StringType,
					},
					"user_info_response_content_encryption_algorithm": schema.StringAttribute{
						Computed:    true,
						Description: "The JSON Web Encryption [JWE] content-encryption algorithm for the UserInfo Response. Supported values are `AES_128_CBC_HMAC_SHA_256`, `AES_192_CBC_HMAC_SHA_384`, `AES_256_CBC_HMAC_SHA_512`, `AES_128_GCM`, `AES_192_GCM`, `AES_256_GCM`. Supported in PF version `12.2` or later.",
					},
					"user_info_response_encryption_algorithm": schema.StringAttribute{
						Computed:    true,
						Description: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the UserInfo response. Supported values are `DIR`, `A128KW`, `A192KW`, `A256KW`, `A128GCMKW`, `A192GCMKW`, `A256GCMKW`, `ECDH_ES`, `ECDH_ES_A128KW`, `ECDH_ES_A192KW`, `ECDH_ES_A256KW`, `RSA_OAEP`, `RSA_OAEP_256`. Supported in PF version `12.2` or later.",
					},
					"user_info_response_signing_algorithm": schema.StringAttribute{
						Computed:    true,
						Description: "The JSON Web Signature [JWS] algorithm required to sign the UserInfo response. Supported values are `NONE`, `HS256`, `HS384`, `HS512`, `RS256`, `RS384`, `RS512`, `ES256`, `ES384`, `ES512`, `PS256`, `PS384`, `PS512`. Supported in PF version `12.2` or later.",
					},
				},
			},
			"client_auth": schema.SingleNestedAttribute{
				Description: "Client authentication settings. If this model is null, it indicates that no client authentication will be used.",
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Client authentication type. The required field for type SECRET is secret.	The required fields for type CERTIFICATE are clientCertIssuerDn and clientCertSubjectDn. The required field for type PRIVATE_KEY_JWT is: either jwks or jwksUrl.",
						Optional:    false,
						Computed:    true,
					},
					"encrypted_secret": schema.StringAttribute{
						Description: "Encrypted client secret.",
						Optional:    false,
						Computed:    true,
					},
					"secondary_secrets": schema.ListNestedAttribute{
						Description: "The list of secondary client secrets that are temporarily retained.",
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"encrypted_secret": schema.StringAttribute{
									Description: "Secondary client secret for Basic Authentication. To update the secondary client secret, specify the plaintext value in this field. This field will not be populated for GET requests.",
									Optional:    false,
									Computed:    true,
								},
								"expiry_time": schema.StringAttribute{
									Description: "The expiry time of the secondary secret.",
									Optional:    false,
									Computed:    true,
								},
							},
						},
					},
					"client_cert_issuer_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Issuer DN.",
						Optional:    false,
						Computed:    true,
					},
					"client_cert_subject_dn": schema.StringAttribute{
						Description: "Client TLS Certificate Subject DN.",
						Optional:    false,
						Computed:    true,
					},
					"enforce_replay_prevention": schema.BoolAttribute{
						Description: "Enforce replay prevention on JSON Web Tokens. This field is applicable only for Private Key JWT Client Authentication.",
						Optional:    false,
						Computed:    true,
					},
					"token_endpoint_auth_signing_algorithm": schema.StringAttribute{
						MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT Client Authentication. All signing algorithms are allowed if value is not present\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
						Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT Client Authentication. All signing algorithms are allowed if value is not present, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
						Optional:            false,
						Computed:            true,
					},
				},
			},
			"jwks_settings": schema.SingleNestedAttribute{
				Description: "JSON Web Key Set Settings of the OAuth client. Required if private key JWT client authentication or signed requests is enabled.",
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"jwks_url": schema.StringAttribute{
						Description: "JSON Web Key Set (JWKS) URL of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    false,
						Computed:    true,
					},
					"jwks": schema.StringAttribute{
						Description: "JSON Web Key Set (JWKS) document of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled. If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.",
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"extended_parameters": schema.MapNestedAttribute{
				Description: "OAuth Client Metadata can be extended to use custom Client Metadata Parameters. The names of these custom parameters should be defined in /extendedProperties.",
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.SetAttribute{
							Description: "A list of values",
							Optional:    false,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"device_flow_setting_type": schema.StringAttribute{
				Description: "Allows an administrator to override the Device Authorization Settings set globally for the OAuth AS. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"user_authorization_url_override": schema.StringAttribute{
				Description: "The URL used as 'verification_url' and 'verification_url_complete' values in a Device Authorization request. This property overrides the 'userAuthorizationUrl' value present in Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"pending_authorization_timeout_override": schema.Int64Attribute{
				Description: "The 'device_code' and 'user_code' timeout, in seconds. This overrides the 'pendingAuthorizationTimeout' value present in Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"device_polling_interval_override": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds. This overrides the 'devicePollingInterval' value present in Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"bypass_activation_code_confirmation_override": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device. This overrides the 'bypassUseCodeConfirmation' value present in Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"require_proof_key_for_code_exchange": schema.BoolAttribute{
				Description: "Determines whether Proof Key for Code Exchange (PKCE) is required for this client.",
				Optional:    false,
				Computed:    true,
			},
			"ciba_delivery_mode": schema.StringAttribute{
				Description: "The token delivery mode for the client. The default value is 'POLL'.",
				Optional:    false,
				Computed:    true,
			},
			"ciba_notification_endpoint": schema.StringAttribute{
				Description: "The endpoint the OP will call after a successful or failed end-user authentication.",
				Optional:    false,
				Computed:    true,
			},
			"ciba_polling_interval": schema.Int64Attribute{
				Description: "The minimum amount of time in seconds that the Client must wait between polling requests to the token endpoint. The default is 0 seconds.",
				Optional:    false,
				Computed:    true,
			},
			"ciba_require_signed_requests": schema.BoolAttribute{
				Description: "Determines whether CIBA signed requests are required for this client.",
				Optional:    false,
				Computed:    true,
			},
			"ciba_request_object_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Optional:            false,
				Computed:            true,
			},
			"ciba_user_code_supported": schema.BoolAttribute{
				Description: "Determines whether the CIBA user code parameter is supported by this client.",
				Optional:    false,
				Computed:    true,
			},
			"request_policy_ref": schema.SingleNestedAttribute{
				Description: "The CIBA request policy.",
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
			},
			"token_exchange_processor_policy_ref": schema.SingleNestedAttribute{
				Description: "The Token Exchange Processor policy.",
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
			},
			"refresh_token_rolling_grace_period_type": schema.StringAttribute{
				Description: "When specified, it overrides the global Refresh Token Grace Period defined in the Authorization Server Settings. The default value is SERVER_DEFAULT",
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds.",
				Optional:    false,
				Computed:    true,
			},
			"client_secret_retention_period_type": schema.StringAttribute{
				Description: "Use OVERRIDE_SERVER_DEFAULT to override the Client Secret Retention Period value on the Authorization Server Settings. SERVER_DEFAULT will default to the Client Secret Retention Period value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.",
				Optional:    false,
				Computed:    true,
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention. This value will override the Client Secret Retention Period value on the Authorization Server Settings.",
				Optional:    false,
				Computed:    true,
			},
			"client_secret_changed_time": schema.StringAttribute{
				Description: "The time at which the client secret was last changed. This property is read only and is ignored on PUT and POST requests.",
				Optional:    false,
				Computed:    true,
			},
			"token_introspection_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response.\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response. HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Optional:            false,
				Computed:            true,
			},
			"token_introspection_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Optional:            false,
				Computed:            true,
			},
			"token_introspection_content_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
				Optional:            false,
				Computed:            true,
			},
			"jwt_secured_authorization_response_mode_signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response.\nHS256 - HMAC using SHA-256\nHS384 - HMAC using SHA-384\nHS512 - HMAC using SHA-512\nRS256 - RSA using SHA-256\nRS384 - RSA using SHA-384\nRS512 - RSA using SHA-512\nES256 - ECDSA using P256 Curve and SHA-256\nES384 - ECDSA using P384 Curve and SHA-384\nES512 - ECDSA using P521 Curve and SHA-512\nPS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256\nPS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384\nPS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512\nA null value will represent the default algorithm which is RS256.\nRSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Description:         "The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response. HS256 - HMAC using SHA-256, HS384 - HMAC using SHA-384, HS512 - HMAC using SHA-512, RS256 - RSA using SHA-256, RS384 - RSA using SHA-384, RS512 - RSA using SHA-512, ES256 - ECDSA using P256 Curve and SHA-256, ES384 - ECDSA using P384 Curve and SHA-384, ES512 - ECDSA using P521 Curve and SHA-512, PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256, PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384, PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512, A null value will represent the default algorithm which is RS256., RSASSA-PSS is only supported with Thales Luna, Entrust nShield Connect or Java 11.",
				Optional:            false,
				Computed:            true,
			},
			"jwt_secured_authorization_response_mode_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response.\nDIR - Direct Encryption with symmetric key\nA128KW - AES-128 Key Wrap\nA192KW - AES-192 Key Wrap\nA256KW - AES-256 Key Wrap\nA128GCMKW - AES-GCM-128 key encryption\nA192GCMKW - AES-GCM-192 key encryption\nA256GCMKW - AES-GCM-256 key encryption\nECDH_ES - ECDH-ES\nECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap\nECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap\nECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap\nRSA_OAEP - RSAES OAEP\nRSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Description:         "The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response. DIR - Direct Encryption with symmetric key, A128KW - AES-128 Key Wrap, A192KW - AES-192 Key Wrap, A256KW - AES-256 Key Wrap, A128GCMKW - AES-GCM-128 key encryption, A192GCMKW - AES-GCM-192 key encryption, A256GCMKW - AES-GCM-256 key encryption, ECDH_ES - ECDH-ES, ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap, ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap, ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap, RSA_OAEP - RSAES OAEP, RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256",
				Optional:            false,
				Computed:            true,
			},
			"jwt_secured_authorization_response_mode_content_encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response.\nAES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256\nAES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384\nAES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512\nAES_128_GCM - AES-GCM-128\nAES_192_GCM - AES-GCM-192\nAES_256_GCM - AES-GCM-256",
				Description:         "The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response. AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256, AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384, AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512, AES_128_GCM - AES-GCM-128, AES_192_GCM - AES-GCM-192, AES_256_GCM - AES-GCM-256",
				Optional:            false,
				Computed:            true,
			},
			"require_dpop": schema.BoolAttribute{
				MarkdownDescription: "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client.",
				Description:         "Determines whether Demonstrating Proof-of-Possession (DPoP) is required for this client.",
				Optional:            false,
				Computed:            true,
			},
			"require_offline_access_scope_to_issue_refresh_tokens": schema.StringAttribute{
				Description: "Determines whether offline_access scope is required to issue refresh tokens by this client or not. 'SERVER_DEFAULT' is the default value.",
				Optional:    false,
				Computed:    true,
			},
			"offline_access_require_consent_prompt": schema.StringAttribute{
				Description: "Determines whether offline_access requires the prompt parameter value to be set to 'consent' by this client or not. The value will be reset to default if the 'requireOfflineAccessScopeToIssueRefreshTokens' attribute is set to 'SERVER_DEFAULT' or 'false'. 'SERVER_DEFAULT' is the default value.",
				Optional:    false,
				Computed:    true,
			},
			"lockout_max_malicious_actions": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of malicious actions allowed before an OAuth client is locked out. Currently, the only operation that is tracked as a malicious action is an attempt to revoke an invalid access token or refresh token. This value will override the global `MaxMaliciousActions` value on the `AccountLockingService` in the config-store. Supported in PF version `12.2` or later.",
			},
			"lockout_max_malicious_actions_type": schema.StringAttribute{
				Computed:    true,
				Description: "Allows an administrator to override the Max Malicious Actions configuration set globally in `AccountLockingService`. Defaults to `SERVER_DEFAULT`. Supported values are `DO_NOT_LOCKOUT`, `SERVER_DEFAULT`, `OVERRIDE_SERVER_DEFAULT`. Supported in PF version `12.2` or later.",
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthClientDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_client"
}

// Configure adds the provider configured client to the data source.
func (r *oauthClientDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a OauthClientResponse object into the model struct
func readOauthClientResponseDataSource(ctx context.Context, r *client.Client, state *oauthClientModel, productVersion version.SupportedVersion) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	diags = readOauthClientResponseCommon(ctx, r, state, nil, productVersion, false)

	// state.ClientAuth
	var secondarySecretsListSlice []attr.Value
	var secondarySecretsObjToState types.List
	secondarySecretsFromClient := r.ClientAuth.GetSecondarySecrets()
	for _, secondarySecretFromClient := range secondarySecretsFromClient {
		secondarySecretAttrVal := map[string]attr.Value{}
		secondarySecretAttrVal["encrypted_secret"] = types.StringPointerValue(secondarySecretFromClient.EncryptedSecret)
		if secondarySecretFromClient.ExpiryTime != nil {
			secondarySecretAttrVal["expiry_time"] = types.StringValue(secondarySecretFromClient.ExpiryTime.Format(time.RFC3339Nano))
		} else {
			secondarySecretAttrVal["expiry_time"] = types.StringNull()
		}
		secondarySecretsAttrValObj, respDiags := types.ObjectValue(secondarySecretsDataSourceAttrType, secondarySecretAttrVal)
		diags.Append(respDiags...)
		secondarySecretsListSlice = append(secondarySecretsListSlice, secondarySecretsAttrValObj)
	}
	secondarySecretsObjToState, respDiags = types.ListValue(types.ObjectType{AttrTypes: secondarySecretsDataSourceAttrType}, secondarySecretsListSlice)
	diags.Append(respDiags...)

	clientAuthAttrValue := map[string]attr.Value{}
	clientAuthAttrValue["type"] = types.StringPointerValue(r.ClientAuth.Type)
	clientAuthAttrValue["encrypted_secret"] = types.StringPointerValue(r.ClientAuth.EncryptedSecret)
	clientAuthAttrValue["secondary_secrets"] = secondarySecretsObjToState
	clientAuthAttrValue["client_cert_issuer_dn"] = types.StringPointerValue(r.ClientAuth.ClientCertIssuerDn)
	clientAuthAttrValue["client_cert_subject_dn"] = types.StringPointerValue(r.ClientAuth.ClientCertSubjectDn)
	clientAuthAttrValue["enforce_replay_prevention"] = types.BoolPointerValue(r.ClientAuth.EnforceReplayPrevention)
	clientAuthAttrValue["token_endpoint_auth_signing_algorithm"] = types.StringPointerValue(r.ClientAuth.TokenEndpointAuthSigningAlgorithm)
	clientAuthToState, respDiags := types.ObjectValue(clientAuthDataSourceAttrType, clientAuthAttrValue)
	diags.Append(respDiags...)
	state.ClientAuth = clientAuthToState

	return diags
}

// Read resource information
func (r *oauthClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthClientModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthClient, httpResp, err := r.apiClient.OauthClientsAPI.GetOauthClientById(config.AuthContext(ctx, r.providerConfig), state.ClientId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Client", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readOauthClientResponseDataSource(ctx, apiReadOauthClient, &state, r.providerConfig.ProductVersion)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
