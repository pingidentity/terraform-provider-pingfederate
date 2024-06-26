package oauthclient_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
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
	includeExtendedParameters                                     bool
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
		includeExtendedParameters: false,
		logoUrl:                   pointers.String(""),
	}

	updatedResourceModel := oauthClientResourceModel{
		clientId:                           oauthClientId,
		name:                               "updatedName",
		grantTypes:                         []string{"IMPLICIT", "AUTHORIZATION_CODE", "RESOURCE_OWNER_CREDENTIALS", "REFRESH_TOKEN", "EXTENSION", "DEVICE_CODE", "ACCESS_TOKEN_VALIDATION", "CIBA", "TOKEN_EXCHANGE"},
		enabled:                            pointers.Bool(true),
		includeOptionalAttributes:          true,
		bypassApprovalPage:                 pointers.Bool(true),
		description:                        pointers.String("updatedDescription"),
		logoUrl:                            pointers.String("https://example.com"),
		redirectUris:                       []string{"https://example.com"},
		allowAuthenticationApiInit:         pointers.Bool(true),
		requirePushedAuthorizationRequests: pointers.Bool(true),
		requireJwtSecuredAuthorizationResponseMode: pointers.Bool(true),
		restrictScopes:                      pointers.Bool(true),
		restrictedScopes:                    []string{"openid"},
		restrictedResponseTypes:             []string{"code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token"},
		restrictToDefaultAccessTokenManager: pointers.Bool(true),
		validateUsingAllEligibleAtms:        pointers.Bool(true),
		oidcPolicy:                          oidcPolicy,
		clientAuth:                          clientAuth,
		jwksSettings:                        jwksSettings,
		requireProofKeyForCodeExchange:      pointers.Bool(true),
		cibaDeliveryMode:                    pointers.String("PING"),
		cibaPollingInterval:                 pointers.Int64(1),
		cibaRequireSignedRequests:           pointers.Bool(true),
		cibaUserCodeSupported:               pointers.Bool(true),
		cibaNotificationEndpoint:            pointers.String("https://example.com"),
		jwtSecuredAuthorizationResponseModeEncryptionAlgorithm:        pointers.String("RSA_OAEP"),
		jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm: pointers.String("AES_128_CBC_HMAC_SHA_256"),
		tokenIntrospectionSigningAlgorithm:                            pointers.String("RS256"),
		requireSignedRequests:                                         pointers.Bool(true),
		includeExtendedParameters:                                     true,
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOauthClientAttributes(initialResourceModel),
					checkPf121ComputedAttrs(resourceName),
				),
			},
			{
				// Test updating some fields
				Config: testAccOauthClient(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOauthClientAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "enabled", fmt.Sprintf("%t", *updatedResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "description", *updatedResourceModel.description),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "logo_url", *updatedResourceModel.logoUrl),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "redirect_uris.0", updatedResourceModel.redirectUris[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "allow_authentication_api_init", fmt.Sprintf("%t", *updatedResourceModel.allowAuthenticationApiInit)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_pushed_authorization_requests", fmt.Sprintf("%t", *updatedResourceModel.requirePushedAuthorizationRequests)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_jwt_secured_authorization_response_mode", fmt.Sprintf("%t", *updatedResourceModel.requireJwtSecuredAuthorizationResponseMode)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restrict_scopes", fmt.Sprintf("%t", *updatedResourceModel.restrictScopes)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_scopes.0", updatedResourceModel.restrictedScopes[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[1]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[2]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[3]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[4]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[5]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restricted_response_types.*", updatedResourceModel.restrictedResponseTypes[6]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "restrict_to_default_access_token_manager", fmt.Sprintf("%t", *updatedResourceModel.restrictToDefaultAccessTokenManager)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "validate_using_all_eligible_atms", fmt.Sprintf("%t", *updatedResourceModel.validateUsingAllEligibleAtms)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.id_token_signing_algorithm", *updatedResourceModel.oidcPolicy.IdTokenSigningAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.grant_access_session_revocation_api", fmt.Sprintf("%t", *updatedResourceModel.oidcPolicy.GrantAccessSessionRevocationApi)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.grant_access_session_session_management_api", fmt.Sprintf("%t", *updatedResourceModel.oidcPolicy.GrantAccessSessionSessionManagementApi)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.ping_access_logout_capable", fmt.Sprintf("%t", *updatedResourceModel.oidcPolicy.PingAccessLogoutCapable)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.pairwise_identifier_user_type", fmt.Sprintf("%t", *updatedResourceModel.oidcPolicy.PairwiseIdentifierUserType)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.sector_identifier_uri", *updatedResourceModel.oidcPolicy.SectorIdentifierUri),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.id_token_encryption_algorithm", *updatedResourceModel.oidcPolicy.IdTokenEncryptionAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "oidc_policy.id_token_content_encryption_algorithm", *updatedResourceModel.oidcPolicy.IdTokenContentEncryptionAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "client_auth.type", *updatedResourceModel.clientAuth.Type),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "jwks_settings.jwks_url", *updatedResourceModel.jwksSettings.JwksUrl),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_proof_key_for_code_exchange", fmt.Sprintf("%t", *updatedResourceModel.requireProofKeyForCodeExchange)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "ciba_delivery_mode", *updatedResourceModel.cibaDeliveryMode),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "ciba_polling_interval", fmt.Sprintf("%d", *updatedResourceModel.cibaPollingInterval)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "ciba_require_signed_requests", fmt.Sprintf("%t", *updatedResourceModel.cibaRequireSignedRequests)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "ciba_user_code_supported", fmt.Sprintf("%t", *updatedResourceModel.cibaUserCodeSupported)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "ciba_notification_endpoint", *updatedResourceModel.cibaNotificationEndpoint),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "jwt_secured_authorization_response_mode_encryption_algorithm", *updatedResourceModel.jwtSecuredAuthorizationResponseModeEncryptionAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "jwt_secured_authorization_response_mode_content_encryption_algorithm", *updatedResourceModel.jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "token_introspection_signing_algorithm", *updatedResourceModel.tokenIntrospectionSigningAlgorithm),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_signed_requests", fmt.Sprintf("%t", *updatedResourceModel.requireSignedRequests)),
				),
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
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthClientsAPI.DeleteOauthClient(ctx, oauthClientId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccOauthClient(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthClientAttributes(initialResourceModel),
			},
		},
	})
}

func checkPf121ComputedAttrs(resourceName string) resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1210) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "refresh_token_rolling_interval_time_unit", "HOURS"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "enable_cookieless_authentication_api", "false"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_offline_access_scope_to_issue_refresh_tokens", "SERVER_DEFAULT"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "offline_access_require_consent_prompt", "SERVER_DEFAULT"),
		)
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "refresh_token_rolling_interval_time_unit"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "enable_cookieless_authentication_api"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "require_offline_access_scope_to_issue_refresh_tokens"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_client.%s", resourceName), "offline_access_require_consent_prompt"),
	)
}

func oidcPolicyHcl(clientOidcPolicy *client.ClientOIDCPolicy) string {
	versionedHcl := ""
	if acctest.VersionAtLeast(version.PingFederate1130) {
		versionedHcl += `
	logout_mode = "OIDC_BACK_CHANNEL"
	back_channel_logout_uri = "https://example.com"	
		`
	}
	if acctest.VersionAtLeast(version.PingFederate1200) {
		versionedHcl += `
	post_logout_redirect_uris = ["https://example.com", "https://pingidentity.com"]
		`
	}
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
	%s
  }
	`, *clientOidcPolicy.IdTokenSigningAlgorithm,
		*clientOidcPolicy.GrantAccessSessionRevocationApi,
		*clientOidcPolicy.GrantAccessSessionSessionManagementApi,
		*clientOidcPolicy.PingAccessLogoutCapable,
		*clientOidcPolicy.PairwiseIdentifierUserType,
		*clientOidcPolicy.SectorIdentifierUri,
		*clientOidcPolicy.IdTokenEncryptionAlgorithm,
		*clientOidcPolicy.IdTokenContentEncryptionAlgorithm,
		versionedHcl)
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

		if acctest.VersionAtLeast(version.PingFederate1130) {
			optionalHcl += `
		require_dpop = true	
			`
		}

		if acctest.VersionAtLeast(version.PingFederate1210) {
			optionalHcl += `
			// HCL necessary to use refresh token rolling interval time unit
		refresh_token_rolling_interval_type = "OVERRIDE_SERVER_DEFAULT"
		refresh_token_rolling_interval = 10
			// PF 12.1 attributes
		refresh_token_rolling_interval_time_unit = "MINUTES"
		enable_cookieless_authentication_api = true
		require_offline_access_scope_to_issue_refresh_tokens = "YES"
		offline_access_require_consent_prompt = "YES"
			`
		}
	}

	if resourceModel.includeExtendedParameters {
		optionalHcl += `
		extended_parameters = {
			"test" = {
				"values" = ["test"]
			}
		}`
	} else {
		optionalHcl += `
		extended_parameters = {}`
	}

	return fmt.Sprintf(`
resource "pingfederate_extended_properties" "%[1]s" {
  items = [
    {
      name         = "test"
      description  = "test"
      multi_valued = false
    }
  ]
}

resource "pingfederate_oauth_client" "%[1]s" {
  client_id   = "%[2]s"
  grant_types = %[3]s
  name        = "%[4]s"
	%[5]s
	%[6]s
	%[7]s
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

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "description", *config.description, *response.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "redirect_uris", config.redirectUris, response.RedirectUris)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_pushed_authorization_requests", *config.requirePushedAuthorizationRequests, *response.RequirePushedAuthorizationRequests)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "restrict_scopes", *config.restrictScopes, *response.RestrictScopes)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthClientId), "restricted_response_types", config.restrictedResponseTypes, response.RestrictedResponseTypes)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "validate_using_all_eligible_atms", *config.validateUsingAllEligibleAtms, *response.ValidateUsingAllEligibleAtms)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "ciba_delivery_mode", *config.cibaDeliveryMode, *response.CibaDeliveryMode)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "ciba_require_signed_requests", *config.cibaRequireSignedRequests, *response.CibaRequireSignedRequests)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "ciba_notification_endpoint", *config.cibaNotificationEndpoint, *response.CibaNotificationEndpoint)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "jwt_secured_authorization_response_mode_content_encryption_algorithm", *config.jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm, *response.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "require_signed_requests", *config.requireSignedRequests, *response.RequireSignedRequests)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.grant_access_session_revocation_api", *config.oidcPolicy.GrantAccessSessionRevocationApi, *response.OidcPolicy.GrantAccessSessionRevocationApi)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthClientId), "oidc_policy.ping_access_logout_capable", *config.oidcPolicy.PingAccessLogoutCapable, *response.OidcPolicy.PingAccessLogoutCapable)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.sector_identifier_uri", *config.oidcPolicy.SectorIdentifierUri, *response.OidcPolicy.SectorIdentifierUri)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthClientId), "oidc_policy.id_token_content_encryption_algorithm", *config.oidcPolicy.IdTokenContentEncryptionAlgorithm, *response.OidcPolicy.IdTokenContentEncryptionAlgorithm)
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
