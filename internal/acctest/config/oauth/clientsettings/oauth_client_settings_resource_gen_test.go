// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package oauthclientsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/accesstokenmanager"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccOauthClientSettings_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model. No computed values to check.
				Config: oauthClientSettings_MinimalHCL(),
			},
			{
				// Update to a complete model
				Config: oauthClientSettings_CompleteHCL(),
			},
			{
				// Test importing the resource
				Config:                               oauthClientSettings_CompleteHCL(),
				ResourceName:                         "pingfederate_oauth_client_settings.example",
				ImportStateVerifyIdentifierAttribute: "dynamic_client_registration.client_cert_issuer_type",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Back to minimal model
				Config: oauthClientSettings_MinimalHCL(),
			},
			{
				// Check that expected computed values are set
				Config: oauthClientSettings_ComputedCheckHCL(),
				Check:  oauthClientSettings_CheckComputedValues(),
			},
		},
	})
}

func oauthClientSettings_DependencyHcl() string {
	return `
resource "pingfederate_oauth_server_settings" "oauthSettings" {
  scopes = [
    {
      name        = "mycommonscope",
      description = "desc",
      dynamic     = false
    }
  ]
  exclusive_scopes = [
    {
      name        = "myexclusivescope",
      description = "desc",
      dynamic     = false
    }
  ]
  authorization_code_entropy = 20
  authorization_code_timeout = 50
  refresh_rolling_interval   = 1
  refresh_token_length       = 40
}

resource "pingfederate_oauth_access_token_manager" "accessTokenManager" {
  manager_id = "accessTokenManager"
  name       = "accessTokenManager"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
  }
  attribute_contract = {
    coreAttributes = []
    extended_attributes = [
      {
        name         = "extended_contract"
        multi_valued = true
      }
    ]
  }
}

resource "pingfederate_openid_connect_policy" "oidcPolicy" {
  policy_id = "oidcPolicy"
  name      = "oidcPolicy"
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.accessTokenManager.manager_id
  }
  attribute_contract = {
    extended_attributes = []
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "sub" = {
        source = {
          type = "TEXT"
        }
        value = "sub"
      }
    }
  }
}
`
}

// Minimal HCL with only required values set (no required values in this resource)
func oauthClientSettings_MinimalHCL() string {
	return fmt.Sprintf(`
	%s

resource "pingfederate_oauth_client_settings" "example" {
}
`, oauthClientSettings_DependencyHcl())
}

// HCL intended to validate expected computed values are set
func oauthClientSettings_ComputedCheckHCL() string {
	return fmt.Sprintf(`
	%s

resource "pingfederate_oauth_client_settings" "example" {
  dynamic_client_registration = {}
}
`, oauthClientSettings_DependencyHcl())
}

// Maximal HCL with all values set where possible
func oauthClientSettings_CompleteHCL() string {
	versionSpecificHcl := ""
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionSpecificHcl = `
	offline_access_require_consent_prompt = "YES"
	refresh_token_rolling_interval_time_unit = "MINUTES"
	require_offline_access_scope_to_issue_refresh_tokens = "YES"
		`
	}
	if acctest.VersionAtLeast(version.PingFederate1220) {
		versionSpecificHcl += `
  lockout_max_malicious_actions_type = "OVERRIDE_SERVER_DEFAULT"
  lockout_max_malicious_actions = 5
    `
	}

	return fmt.Sprintf(`
	%s

  %s

resource "pingfederate_oauth_ciba_server_policy_request_policy" "example" {
  allow_unsigned_login_hint_token = false
  alternative_login_hint_token_issuers = [
  ]
  authenticator_ref = {
    id = "exampleCibaAuthenticator"
  }
  identity_hint_contract = {
    extended_attributes = [
    ]
  }
  identity_hint_contract_fulfillment = {
    attribute_contract_fulfillment = {
      IDENTITY_HINT_SUBJECT = {
        source = {
          id   = null
          type = "REQUEST"
        }
        value = "IDENTITY_HINT_SUBJECT"
      }
    }
  }
  identity_hint_mapping = {
    attribute_contract_fulfillment = {
      USER_KEY = {
        source = {
          id   = null
          type = "NO_MAPPING"
        }
        value = null
      }
      subject = {
        source = {
          id   = null
          type = "NO_MAPPING"
        }
        value = null
      }
    }
  }
  name                            = "oauthClientSettingsTestCibaPolicy"
  policy_id                       = "oauthClientSettingsTestCibaPolicy"
  require_token_for_identity_hint = false
  transaction_lifetime            = 120
}

resource "pingfederate_oauth_client_settings" "example" {
  depends_on = [
    pingfederate_oauth_server_settings.oauthSettings,
    pingfederate_openid_connect_policy.oidcPolicy
  ]
  # Ensures this resource will be updated before deleting the oauth access token manager or ciba server policy
  lifecycle {
    create_before_destroy = true
  }
  dynamic_client_registration = {
    allow_client_delete                          = false
    allowed_exclusive_scopes                     = ["myexclusivescope"]
    bypass_activation_code_confirmation_override = false
    ciba_polling_interval                        = 5
    ciba_require_signed_requests                 = true
    client_cert_issuer_ref = {
      id = "gdxuvcw6p95rex3go7eb3ctsb"
    }
    client_cert_issuer_type                 = "CERTIFICATE"
    client_secret_retention_period_override = 2
    client_secret_retention_period_type     = "OVERRIDE_SERVER_DEFAULT"
    default_access_token_manager_ref = {
      id = pingfederate_oauth_access_token_manager.oauthClientSettingsAtm.id
    }
    device_flow_setting_type           = "OVERRIDE_SERVER_DEFAULT"
    device_polling_interval_override   = 5
    disable_registration_access_tokens = false
    enforce_replay_prevention          = true
    initial_access_token_scope         = "mycommonscope"
    oidc_policy = {
      id_token_signing_algorithm = "ES256"
      policy_group = {
        id = "oidcPolicy"
      }
    }
    pending_authorization_timeout_override  = 5
    persistent_grant_expiration_time        = 5
    persistent_grant_expiration_time_unit   = "MINUTES"
    persistent_grant_expiration_type        = "OVERRIDE_SERVER_DEFAULT"
    persistent_grant_idle_timeout           = 3
    persistent_grant_idle_timeout_time_unit = "MINUTES"
    persistent_grant_idle_timeout_type      = "OVERRIDE_SERVER_DEFAULT"
    policy_refs = [
      {
        id = "clientRegistrationPolicy"
      }
    ]
    refresh_rolling                         = "ROLL"
    refresh_token_rolling_grace_period      = 60
    refresh_token_rolling_grace_period_type = "OVERRIDE_SERVER_DEFAULT"
    refresh_token_rolling_interval          = 10
    refresh_token_rolling_interval_type     = "OVERRIDE_SERVER_DEFAULT"
    request_policy_ref = {
      id = pingfederate_oauth_ciba_server_policy_request_policy.example.id
    }
    require_jwt_secured_authorization_response_mode = true
    require_proof_key_for_code_exchange             = true
    require_signed_requests                         = true
    restrict_common_scopes                          = true
    restrict_to_default_access_token_manager        = true
    restricted_common_scopes                        = ["mycommonscope"]
    retain_client_secret                            = true
    rotate_client_secret                            = false
    rotate_registration_access_token                = false
    token_exchange_processor_policy_ref = {
      id = "tokenexchangeprocessorpolicy"
    }
    user_authorization_url_override = "https://example.com"
	%s
  }
}
`, oauthClientSettings_DependencyHcl(),
		accesstokenmanager.TestAccessTokenManagerHCL("oauthClientSettingsAtm"),
		versionSpecificHcl)
}

// Validate any computed values when applying HCL that expects computed values
func oauthClientSettings_CheckComputedValues() resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.allow_client_delete", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.allowed_authorization_detail_types.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.allowed_exclusive_scopes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.bypass_activation_code_confirmation_override", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.ciba_polling_interval", "3"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.ciba_require_signed_requests", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.client_cert_issuer_type", "NONE"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.client_secret_retention_period_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.device_flow_setting_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.disable_registration_access_tokens", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.enforce_replay_prevention", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.persistent_grant_expiration_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.persistent_grant_idle_timeout_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.refresh_rolling", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.refresh_token_rolling_grace_period_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.refresh_token_rolling_interval_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.require_jwt_secured_authorization_response_mode", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.require_proof_key_for_code_exchange", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.require_signed_requests", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.restrict_common_scopes", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.restrict_to_default_access_token_manager", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.restricted_common_scopes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.retain_client_secret", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.rotate_client_secret", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.rotate_registration_access_token", "false"),
	}

	if acctest.VersionAtLeast(version.PingFederate1210) {
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.offline_access_require_consent_prompt", "SERVER_DEFAULT"))
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.require_offline_access_scope_to_issue_refresh_tokens", "SERVER_DEFAULT"))
	}

	if acctest.VersionAtLeast(version.PingFederate1220) {
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("pingfederate_oauth_client_settings.example", "dynamic_client_registration.lockout_max_malicious_actions_type", "SERVER_DEFAULT"))
	}

	return resource.ComposeTestCheckFunc(testCheckFuncs...)
}
