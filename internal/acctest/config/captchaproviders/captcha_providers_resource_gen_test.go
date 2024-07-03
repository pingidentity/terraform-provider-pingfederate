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
)

const captchaProvidersProvidersId = "captchaProvidersProvidersId"

var testEnvId = ""
var testRiskPolicyId = ""

func TestAccCaptchaProviders_RemovalDrift(t *testing.T) {
	envId := os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")
	riskPolicyId := os.Getenv("PF_TF_P1_RISK_POLICY_ID")
	testEnvId = envId
	testRiskPolicyId = riskPolicyId
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if riskPolicyId == "" {
				t.Fatal("PF_TF_P1_RISK_POLICY_ID must be set")
			}
			if envId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: captchaProviders_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: captchaProviders_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					captchaProviders_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCaptchaProviders_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: captchaProviders_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: captchaProviders_MinimalHCL(),
				Check:  captchaProviders_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  captchaProviders_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: captchaProviders_CompleteHCL(),
				Check:  captchaProviders_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: captchaProviders_MinimalHCL(),
				Check:  captchaProviders_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: captchaProviders_CompleteHCL(),
				Check:  captchaProviders_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               captchaProviders_CompleteHCL(),
				ResourceName:                         "pingfederate_captcha_providers.example",
				ImportStateId:                        captchaProvidersProvidersId,
				ImportStateVerifyIdentifierAttribute: "providers_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore: []string{
					"configuration.tables",
					"configuration.fields",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func captchaProviders_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_captcha_providers" "example" {
  providers_id = "%s"
  name         = "%s"
  configuration = {
    tables = [],
    fields = [
      {
        name  = "Site Key"
        value = "testSiteKey"
      },
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
`, captchaProvidersProvidersId, captchaProvidersProvidersId)
}

// Maximal HCL with all values set where possible
func captchaProviders_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_captcha_providers" "example" {
  providers_id = "%s"
  name         = "%s"
  configuration = {
    tables = [],
    fields = [
      {
        name : "PingOne Environment"
        value : "%s"
      },
      {
        name : "PingOne Risk Policy"
        value : "%s"
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.pingone.protect.PingOneProtectProvider"
  }
}
`, captchaProvidersProvidersId, captchaProvidersProvidersId, testEnvId, testRiskPolicyId)
}

// Validate any computed values when applying minimal HCL
func captchaProviders_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.2.value", "recaptcha-v2-invisible.js"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.tables_all.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.#", "3"),
	)
}

// Validate any computed values when applying complete HCL
func captchaProviders_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.tables_all.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.#", "12"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.2.value", "true"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.3.value", "SHA-256"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.6.value", "MEDIUM"),
		resource.TestCheckResourceAttr("pingfederate_captcha_providers.example", "configuration.fields_all.11.value", "50"),
	)
}

// Delete the resource
func captchaProviders_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.CaptchaProvidersAPI.DeleteCaptchaProvider(acctest.TestBasicAuthContext(), captchaProvidersProvidersId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func captchaProviders_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.CaptchaProvidersAPI.DeleteCaptchaProvider(acctest.TestBasicAuthContext(), captchaProvidersProvidersId).Execute()
	if err == nil {
		return fmt.Errorf("captcha_providers still exists after tests. Expected it to be destroyed")
	}
	return nil
}
