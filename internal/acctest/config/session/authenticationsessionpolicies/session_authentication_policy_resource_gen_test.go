// Code generated by ping-terraform-plugin-framework-generator

package sessionauthenticationsessionpolicies_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const sessionAuthenticationPolicyPolicyId = "session_auth_policy_policy_id"

func TestAccSessionAuthenticationPolicy_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: sessionAuthenticationPolicy_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: sessionAuthenticationPolicy_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					sessionAuthenticationPolicy_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSessionAuthenticationPolicy_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: sessionAuthenticationPolicy_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: sessionAuthenticationPolicy_MinimalHCL(),
				Check:  sessionAuthenticationPolicy_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  sessionAuthenticationPolicy_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: sessionAuthenticationPolicy_CompleteHCL(),
				Check:  sessionAuthenticationPolicy_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: sessionAuthenticationPolicy_MinimalHCL(),
				Check:  sessionAuthenticationPolicy_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: sessionAuthenticationPolicy_CompleteHCL(),
				Check:  sessionAuthenticationPolicy_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               sessionAuthenticationPolicy_CompleteHCL(),
				ResourceName:                         "pingfederate_session_authentication_policy.example",
				ImportStateId:                        sessionAuthenticationPolicyPolicyId,
				ImportStateVerifyIdentifierAttribute: "policy_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func sessionAuthenticationPolicy_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_session_authentication_policy" "example" {
  policy_id = "%s"
  authentication_source = {
    source_ref = {
                id = "OTIdPJava"
              }
    type = "IDP_ADAPTER"
  }
  enable_sessions = false
}
`, sessionAuthenticationPolicyPolicyId)
}

// Maximal HCL with all values set where possible
func sessionAuthenticationPolicy_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_session_authentication_policy" "example" {
  policy_id = "%s"
  authentication_source = {
    source_ref = {
                id = "OTIdPJava"
              }
    type = "IDP_ADAPTER"
  }
  authn_context_sensitive = true
  enable_sessions = true
  idle_timeout_mins = 60
  max_timeout_mins = 480
  persistent = true
  timeout_display_unit = "HOURS"
  user_device_type = "ANY"
}
`, sessionAuthenticationPolicyPolicyId)
}

// Validate any computed values when applying minimal HCL
func sessionAuthenticationPolicy_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "authn_context_sensitive", "false"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "persistent", "false"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "timeout_display_unit", "MINUTES"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "user_device_type", "PRIVATE"),
	)
}

// Validate any computed values when applying complete HCL
func sessionAuthenticationPolicy_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "authn_context_sensitive", "true"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "persistent", "true"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "timeout_display_unit", "HOURS"),
		resource.TestCheckResourceAttr("pingfederate_session_authentication_policy.example", "user_device_type", "ANY"),
	)
}

// Delete the resource
func sessionAuthenticationPolicy_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.SessionAPI.DeleteSourcePolicy(acctest.TestBasicAuthContext(), sessionAuthenticationPolicyPolicyId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func sessionAuthenticationPolicy_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SessionAPI.DeleteSourcePolicy(acctest.TestBasicAuthContext(), sessionAuthenticationPolicyPolicyId).Execute()
	if err == nil {
		return fmt.Errorf("session_authentication_policy still exists after tests. Expected it to be destroyed")
	}
	return nil
}
