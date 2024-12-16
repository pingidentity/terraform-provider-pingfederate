// Code generated by ping-terraform-plugin-framework-generator

package captchaproviders_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const captchaProviderProviderId = "captchaProviderProviderId"

var testEnvConnId = ""

func TestAccCaptchaProvider_RemovalDrift(t *testing.T) {
	connId := os.Getenv("PF_TF_P1_CONNECTION_ID")
	envId := os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")
	testEnvConnId = connId + "|" + envId
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if connId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ID must be set for the Captcha Provider acceptance test")
			}
			if envId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for the Captcha Provider acceptance test")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: captchaProvider_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: captchaProvider_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					captchaProvider_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCaptchaProvider_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: captchaProvider_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: captchaProvider_MinimalHCL(),
				Check:  captchaProvider_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  captchaProvider_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: captchaProvider_CompleteHCL(),
				Check:  captchaProvider_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: captchaProvider_MinimalHCL(),
				Check:  captchaProvider_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: captchaProvider_CompleteHCL(),
				Check:  captchaProvider_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               captchaProvider_CompleteHCL(),
				ResourceName:                         "pingfederate_captcha_provider.example",
				ImportStateId:                        captchaProviderProviderId,
				ImportStateVerifyIdentifierAttribute: "provider_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Sensitive values aren't returned by PF, so they can't be verified
				ImportStateVerifyIgnore: []string{
					"configuration.sensitive_fields.0.value",
					"configuration.sensitive_fields.0.encrypted_value",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func captchaProvider_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_captcha_provider" "example" {
  provider_id = "%s"
  name        = "%s"
  configuration = {
    tables = [],
    fields = [
      {
        name  = "Site Key"
        value = "testSiteKey"
      },
    ]
    sensitive_fields = [
      {
        name  = "Secret Key"
        value = "testSecretKey"
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.captcha.ReCaptchaV2InvisiblePlugin"
  }
}
`, captchaProviderProviderId, captchaProviderProviderId)
}

// Maximal HCL with all values set where possible
func captchaProvider_CompleteHCL() string {
	if acctest.VersionAtLeast(version.PingFederate1200) {
		// The PingOneProtectProvider was added in PF version 12.0+
		return fmt.Sprintf(`
resource "pingfederate_captcha_provider" "example" {
  provider_id = "%s"
  name        = "%s"
  configuration = {
    tables = [],
    fields = [
      {
        name : "PingOne Environment"
        value : "%s"
      },
      {
        name : "PingOne Risk Policy"
        value : "f277d6e2-e073-018c-1b78-8be4cd16d898"
      },
      {
        "name" : "API Request Timeout",
        "value" : "2000"
      },
      {
        "name" : "Custom Proxy Host",
        "value" : ""
      },
      {
        "name" : "Custom Proxy Port",
        "value" : ""
      },
      {
        "name" : "Custom connection pool",
        "value" : "50"
      },
      {
        "name" : "Enable Risk Evaluation",
        "value" : "true"
      },
      {
        "name" : "Failure Mode",
        "value" : "Continue with fallback policy decision"
      },
      {
        "name" : "Fallback Policy Decision Value",
        "value" : "MEDIUM"
      },
      {
        "name" : "Follow Recommended Action",
        "value" : "true"
      },
      {
        "name" : "Password Encryption",
        "value" : "SHA-256"
      },
      {
        "name" : "Proxy Settings",
        "value" : "System Defaults"
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.pingone.protect.PingOneProtectProvider"
  }
}
`, captchaProviderProviderId, captchaProviderProviderId, testEnvConnId)
	} else {
		// For earlier versions use captcha v3
		return fmt.Sprintf(`
resource "pingfederate_captcha_provider" "example" {
  provider_id = "%s"
  name        = "%s"
  configuration = {
    tables = [],
    fields = [
      {
        name : "Site Key"
        value : "1234"
      },
      {
        name : "Pass Score Threshold"
        value : "0.8"
      },
      {
        name : "JavaScript File Name"
        value : "recaptcha-v3.js"
      }
    ]
    sensitive_fields = [
      {
        name : "Secret Key"
        value : "1234"
      },
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.captcha.recaptchaV3.ReCaptchaV3Plugin"
  }
}
`, captchaProviderProviderId, captchaProviderProviderId)
	}
}

// Validate any computed values when applying minimal HCL
func captchaProvider_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "id", captchaProviderProviderId),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
			map[string]string{
				"value": "recaptcha-v2-invisible.js",
			},
		),
		resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.tables_all.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.fields_all.#", "3"),
	)
}

// Validate any computed values when applying complete HCL
func captchaProvider_CheckComputedValuesComplete() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1200) {
		// The PingOneProtectProvider was added in PF version 12.0+
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "id", captchaProviderProviderId),
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.tables_all.#", "0"),
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.fields_all.#", "12"),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"value": "true",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"value": "SHA-256",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"value": "MEDIUM",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"value": "50",
				},
			),
		)
	} else {
		// For earlier versions use captcha v3
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "id", captchaProviderProviderId),
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.tables_all.#", "0"),
			resource.TestCheckResourceAttr("pingfederate_captcha_provider.example", "configuration.fields_all.#", "4"),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"name":  "Site Key",
					"value": "1234",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"name": "Secret Key",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"name":  "Pass Score Threshold",
					"value": "0.8",
				},
			),
			resource.TestCheckTypeSetElemNestedAttrs("pingfederate_captcha_provider.example", "configuration.fields_all.*",
				map[string]string{
					"name":  "JavaScript File Name",
					"value": "recaptcha-v3.js",
				},
			),
		)
	}
}

// Delete the resource
func captchaProvider_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.CaptchaProvidersAPI.DeleteCaptchaProvider(acctest.TestBasicAuthContext(), captchaProviderProviderId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func captchaProvider_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.CaptchaProvidersAPI.DeleteCaptchaProvider(acctest.TestBasicAuthContext(), captchaProviderProviderId).Execute()
	if err == nil {
		return fmt.Errorf("captcha_provider still exists after tests. Expected it to be destroyed")
	}
	return nil
}
