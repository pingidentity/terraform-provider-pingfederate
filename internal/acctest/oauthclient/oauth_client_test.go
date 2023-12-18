package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const oauthClientId = "myOauthClient"

type oauthClientResourceModel struct {
	clientId                                                      string
	grantTypes                                                    []string
	name                                                          string
	includeOptionalAttributes                                     bool
	enabled                                                       *bool
	bypassApprovalPage                                            *bool
	description                                                   *string
	logoUrl                                                       *string
	redirectUris                                                  []string
	allowAuthenticationApiInit                                    *bool
	requirePushedAuthorizationRequests                            *bool
	requireJwtSecuredAuthorizationResponseMode                    *bool
	restrictScopes                                                *bool
	restrictedScopes                                              []string
	restrictedResponseTypes                                       []string
	restrictToDefaultAccessTokenManager                           *bool
	validateUsingAllEligibleAtms                                  *bool
	oidcPolicy                                                    *client.ClientOIDCPolicy
	clientAuth                                                    *client.ClientAuth
	jwksSettings                                                  *client.JwksSettings
	requireProofKeyForCodeExchange                                *bool
	cibaDeliveryMode                                              *string
	cibaPollingInterval                                           *int64
	cibaRequireSignedRequests                                     *bool
	cibaUserCodeSupported                                         *bool
	cibaNotificationEndpoint                                      *string
	jwtSecuredAuthorizationResponseModeEncryptionAlgorithm        *string
	jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm *string
	tokenIntrospectionSigningAlgorithm                            *string
	requireSignedRequests                                         *bool
}

func TestAccOauthClient(t *testing.T) {
	resourceName := oauthClientId

	oidcPolicy := client.NewClientOIDCPolicyWithDefaults()
	oidcPolicy.IdTokenSigningAlgorithm = pointers.String("HS256")
	oidcPolicy.GrantAccessSessionRevocationApi = pointers.Bool(false)
	oidcPolicy.GrantAccessSessionSessionManagementApi = pointers.Bool(false)
	oidcPolicy.PingAccessLogoutCapable = pointers.Bool(false)
	oidcPolicy.PairwiseIdentifierUserType = pointers.Bool(true)
	oidcPolicy.SectorIdentifierUri = pointers.String("https://example.com")
	oidcPolicy.IdTokenEncryptionAlgorithm = pointers.String("A192GCMKW")
	oidcPolicy.IdTokenContentEncryptionAlgorithm = pointers.String("AES_128_CBC_HMAC_SHA_256")

	clientAuth := client.NewClientAuthWithDefaults()
	clientAuth.Type = pointers.String("SECRET")
	clientAuth.Secret = pointers.String("mySecretValue")

	jwksSettings := client.NewJwksSettingsWithDefaults()
	jwksSettings.JwksUrl = pointers.String("https://example.com")

	initialResourceModel := oauthClientResourceModel{
		clientId:                  oauthClientId,
		name:                      "initialName",
		grantTypes:                []string{"DEVICE_CODE"},
		includeOptionalAttributes: false,
	}

	updatedResourceModel := oauthClientResourceModel{
		clientId:                           oauthClientId,
		name:                               "updatedName",
		grantTypes:                         []string{"IMPLICIT", "AUTHORIZATION_CODE", "RESOURCE_OWNER_CREDENTIALS", "REFRESH_TOKEN", "EXTENSION", "DEVICE_CODE", "ACCESS_TOKEN_VALIDATION", "CIBA", "TOKEN_EXCHANGE"},
		enabled:                            pointers.Bool(false),
		includeOptionalAttributes:          true,
		bypassApprovalPage:                 pointers.Bool(true),
		description:                        pointers.String("updatedDescription"),
		logoUrl:                            pointers.String("https://example.com"),
		redirectUris:                       []string{"https://example.com"},
		allowAuthenticationApiInit:         pointers.Bool(false),
		requirePushedAuthorizationRequests: pointers.Bool(false),
		requireJwtSecuredAuthorizationResponseMode: pointers.Bool(false),
		restrictScopes:                      pointers.Bool(true),
		restrictedScopes:                    []string{"openid"},
		restrictedResponseTypes:             []string{"code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token"},
		restrictToDefaultAccessTokenManager: pointers.Bool(false),
		validateUsingAllEligibleAtms:        pointers.Bool(false),
		oidcPolicy:                          oidcPolicy,
		clientAuth:                          clientAuth,
		jwksSettings:                        jwksSettings,
		requireProofKeyForCodeExchange:      pointers.Bool(false),
		cibaDeliveryMode:                    pointers.String("PING"),
		cibaPollingInterval:                 pointers.Int64(1),
		cibaRequireSignedRequests:           pointers.Bool(true),
		cibaUserCodeSupported:               pointers.Bool(false),
		cibaNotificationEndpoint:            pointers.String("https://example.com"),
		jwtSecuredAuthorizationResponseModeEncryptionAlgorithm:        pointers.String("RSA_OAEP"),
		jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm: pointers.String("AES_128_CBC_HMAC_SHA_256"),
		tokenIntrospectionSigningAlgorithm:                            pointers.String("RS256"),
		requireSignedRequests:                                         pointers.Bool(true),
	}

	//  Client Auth and Redirect URIs are required for the resource when going back to the minimal model from the updated model
	minimalResourceModel := oauthClientResourceModel{
		clientId:                  oauthClientId,
		name:                      "updatedName",
		grantTypes:                []string{"DEVICE_CODE"},
		redirectUris:              []string{"https://example.com"},
		clientAuth:                clientAuth,
		includeOptionalAttributes: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthClient(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthClientAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthClient(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthClientAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthClient(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_client." + resourceName,
				ImportStateId:     oauthClientId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_auth.secret",
					"client_secret_changed_time",
					"modification_date",
				},
			},
			{
				// Back to minimal model
				Config: testAccOauthClient(resourceName, minimalResourceModel),
				Check:  testAccCheckExpectedOauthClientAttributes(minimalResourceModel),
			},
		},
	})
}

func oidcPolicyHcl(clientOidcPolicy *client.ClientOIDCPolicy) string {
	return fmt.Sprintf(`
  oidc_policy = {
    id_token_signing_algorithm                  = "%s"
    grant_access_session_revocation_api         = %t
    grant_access_session_session_management_api = %t
    ping_access_logout_capable                  = %t
    pairwise_identifier_user_type               = %t
    sector_identifier_uri                       = "%s"
    id_token_encryption_algorithm               = "%s"
    id_token_content_encryption_algorithm       = "%s"
  }
	`, *clientOidcPolicy.IdTokenSigningAlgorithm,
		*clientOidcPolicy.GrantAccessSessionRevocationApi,
		*clientOidcPolicy.GrantAccessSessionSessionManagementApi,
		*clientOidcPolicy.PingAccessLogoutCapable,
		*clientOidcPolicy.PairwiseIdentifierUserType,
		*clientOidcPolicy.SectorIdentifierUri,
		*clientOidcPolicy.IdTokenEncryptionAlgorithm,
		*clientOidcPolicy.IdTokenContentEncryptionAlgorithm)
}

func clientAuthHcl(clientAuth *client.ClientAuth) string {
	return fmt.Sprintf(`
  client_auth = {
		type   = "%s"
		secret = "%s"
  }
	`, *clientAuth.Type,
		*clientAuth.Secret)
}

func jwksSettingsHcl(jwksSettings *client.JwksSettings) string {
	return fmt.Sprintf(`
	jwks_settings = {
		jwks_url = "%s"
	}
	`, *jwksSettings.JwksUrl)
}

func testAccOauthClient(resourceName string, resourceModel oauthClientResourceModel) string {
	optionalHcl := ""
	optionalRedirectUris := ""
	optionalClientAuth := ""
	if resourceModel.redirectUris != nil {
		optionalRedirectUris = fmt.Sprintf(`
		redirect_uris = %s
		`, acctest.StringSliceToTerraformString(resourceModel.redirectUris))
	}
	if resourceModel.clientAuth != nil {
		optionalClientAuth = clientAuthHcl(resourceModel.clientAuth)
	}
	if resourceModel.includeOptionalAttributes {
		optionalHcl = fmt.Sprintf(`
		enabled = %t
		bypass_approval_page = %t
		description = "%s"
		logo_url = "%s"
		allow_authentication_api_init = %t
		require_pushed_authorization_requests = %t
		require_jwt_secured_authorization_response_mode = %t
		restrict_scopes = %t
		restricted_scopes = %s
		restricted_response_types = %s
		restrict_to_default_access_token_manager = %t
		validate_using_all_eligible_atms = %t
		%s
		%s
		require_proof_key_for_code_exchange = %t
		ciba_delivery_mode = "%s"
		ciba_polling_interval = %d
		ciba_require_signed_requests = %t
		ciba_user_code_supported = %t
		ciba_notification_endpoint = "%s"
		jwt_secured_authorization_response_mode_encryption_algorithm = "%s"
		jwt_secured_authorization_response_mode_content_encryption_algorithm = "%s"
		token_introspection_signing_algorithm = "%s"
		require_signed_requests = %t
		`,
			*resourceModel.enabled,
			*resourceModel.bypassApprovalPage,
			*resourceModel.description,
			*resourceModel.logoUrl,
			*resourceModel.allowAuthenticationApiInit,
			*resourceModel.requirePushedAuthorizationRequests,
			*resourceModel.requireJwtSecuredAuthorizationResponseMode,
			*resourceModel.restrictScopes,
			acctest.StringSliceToTerraformString(resourceModel.restrictedScopes),
			acctest.StringSliceToTerraformString(resourceModel.restrictedResponseTypes),
			*resourceModel.restrictToDefaultAccessTokenManager,
			*resourceModel.validateUsingAllEligibleAtms,
			oidcPolicyHcl(resourceModel.oidcPolicy),
			jwksSettingsHcl(resourceModel.jwksSettings),
			*resourceModel.requireProofKeyForCodeExchange,
			*resourceModel.cibaDeliveryMode,
			*resourceModel.cibaPollingInterval,
			*resourceModel.cibaRequireSignedRequests,
			*resourceModel.cibaUserCodeSupported,
			*resourceModel.cibaNotificationEndpoint,
			*resourceModel.jwtSecuredAuthorizationResponseModeEncryptionAlgorithm,
			*resourceModel.jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm,
			*resourceModel.tokenIntrospectionSigningAlgorithm,
			*resourceModel.requireSignedRequests,
		)
	}

	return fmt.Sprintf(`
resource "pingfederate_oauth_client" "%s" {
  client_id   = "%s"
  grant_types = %s
  name        = "%s"
	%s
	%s
	%s
}
data "pingfederate_oauth_client" "%s" {
  client_id = pingfederate_oauth_client.%s.client_id
}`, resourceName,
		oauthClientId,
		acctest.StringSliceToTerraformString(resourceModel.grantTypes),
		resourceModel.name,
		optionalRedirectUris,
		optionalHcl,
		optionalClientAuth,
		resourceName,
		oauthClientId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthClientAttributes(config oauthClientResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthClient"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthClientsAPI.GetOauthClientById(ctx, oauthClientId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "name", config.name, response.Name)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "grant_types", config.grantTypes, response.GrantTypes)
		if err != nil {
			return err
		}

		if config.redirectUris != nil {
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "redirect_uris", config.redirectUris, response.RedirectUris)
			if err != nil {
				return err
			}
		}

		if config.includeOptionalAttributes {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "enabled", *config.enabled, *response.Enabled)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "bypass_approval_page", *config.bypassApprovalPage, *response.BypassApprovalPage)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "description", *config.description, *response.Description)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "logo_url", *config.logoUrl, *response.LogoUrl)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "redirect_uris", config.redirectUris, response.RedirectUris)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "allow_authentication_api_init", *config.allowAuthenticationApiInit, *response.AllowAuthenticationApiInit)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_pushed_authorization_requests", *config.requirePushedAuthorizationRequests, *response.RequirePushedAuthorizationRequests)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_jwt_secured_authorization_response_mode", *config.requireJwtSecuredAuthorizationResponseMode, *response.RequireJwtSecuredAuthorizationResponseMode)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "restrict_scopes", *config.restrictScopes, *response.RestrictScopes)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "restricted_scopes", config.restrictedScopes, response.RestrictedScopes)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "restricted_response_types", config.restrictedResponseTypes, response.RestrictedResponseTypes)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "restrict_to_default_access_token_manager", *config.restrictToDefaultAccessTokenManager, *response.RestrictToDefaultAccessTokenManager)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "validate_using_all_eligible_atms", *config.validateUsingAllEligibleAtms, *response.ValidateUsingAllEligibleAtms)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_proof_key_for_code_exchange", *config.requireProofKeyForCodeExchange, *response.RequireProofKeyForCodeExchange)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "ciba_delivery_mode", *config.cibaDeliveryMode, *response.CibaDeliveryMode)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(oauthClientId), "ciba_polling_interval", *config.cibaPollingInterval, *response.CibaPollingInterval)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "ciba_require_signed_requests", *config.cibaRequireSignedRequests, *response.CibaRequireSignedRequests)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "ciba_user_code_supported", *config.cibaUserCodeSupported, *response.CibaUserCodeSupported)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "ciba_notification_endpoint", *config.cibaNotificationEndpoint, *response.CibaNotificationEndpoint)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "jwt_secured_authorization_response_mode_encryption_algorithm", *config.jwtSecuredAuthorizationResponseModeEncryptionAlgorithm, *response.JwtSecuredAuthorizationResponseModeEncryptionAlgorithm)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "jwt_secured_authorization_response_mode_content_encryption_algorithm", *config.jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm, *response.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "token_introspection_signing_algorithm", *config.tokenIntrospectionSigningAlgorithm, *response.TokenIntrospectionSigningAlgorithm)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_signed_requests", *config.requireSignedRequests, *response.RequireSignedRequests)
			if err != nil {
				return err
			}

			//  test oidc policy
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.id_token_signing_algorithm", *config.oidcPolicy.IdTokenSigningAlgorithm, *response.OidcPolicy.IdTokenSigningAlgorithm)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.grant_access_session_revocation_api", *config.oidcPolicy.GrantAccessSessionRevocationApi, *response.OidcPolicy.GrantAccessSessionRevocationApi)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.grant_access_session_session_management_api", *config.oidcPolicy.GrantAccessSessionSessionManagementApi, *response.OidcPolicy.GrantAccessSessionSessionManagementApi)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.ping_access_logout_capable", *config.oidcPolicy.PingAccessLogoutCapable, *response.OidcPolicy.PingAccessLogoutCapable)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.pairwise_identifier_user_type", *config.oidcPolicy.PairwiseIdentifierUserType, *response.OidcPolicy.PairwiseIdentifierUserType)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.sector_identifier_uri", *config.oidcPolicy.SectorIdentifierUri, *response.OidcPolicy.SectorIdentifierUri)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.id_token_encryption_algorithm", *config.oidcPolicy.IdTokenEncryptionAlgorithm, *response.OidcPolicy.IdTokenEncryptionAlgorithm)
			if err != nil {
				return err
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.id_token_content_encryption_algorithm", *config.oidcPolicy.IdTokenContentEncryptionAlgorithm, *response.OidcPolicy.IdTokenContentEncryptionAlgorithm)
			if err != nil {
				return err
			}

			// test client auth
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "client_auth.type", *config.clientAuth.Type, *response.ClientAuth.Type)
			if err != nil {
				return err
			}

			// test jwks settings
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "jwks_settings.jwks_url", *config.jwksSettings.JwksUrl, *response.JwksSettings.JwksUrl)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthClientDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthClientsAPI.DeleteOauthClient(ctx, oauthClientId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthClient", oauthClientId)
	}
	return nil
}
