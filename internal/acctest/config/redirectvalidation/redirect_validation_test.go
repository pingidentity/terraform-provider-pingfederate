package redirectvalidation_test

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

// Attributes to test with. Add optional properties to test here if desired.
type redirectValidationResourceModel struct {
	enableTargetResourceValidationForSso bool
	whiteListValidDomain                 string
	enableWreplyValidationSlo            bool
}

func TestAccRedirectValidation(t *testing.T) {
	resourceName := "myRedirectValidation"
	updatedResourceModel := redirectValidationResourceModel{
		enableTargetResourceValidationForSso: false,
		whiteListValidDomain:                 "updatedexample.com",
		enableWreplyValidationSlo:            true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectValidationMinial(resourceName),
			},
			{
				// Test updating some fields
				Config: testAccRedirectValidation(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRedirectValidationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRedirectValidation(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_redirect_validation." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccRedirectValidationMinial(resourceName),
			},
		},
	})
}

func testAccRedirectValidationMinial(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_redirect_validation" "%[1]s" {
}`, resourceName)
}

func testAccRedirectValidation(resourceName string, resourceModel redirectValidationResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_redirect_validation" "%[1]s" {
  redirect_validation_local_settings = {
    enable_target_resource_validation_for_sso           = %[2]t
    enable_target_resource_validation_for_slo           = true
    enable_target_resource_validation_for_idp_discovery = true
    enable_in_error_resource_validation                 = true
    white_list = [
      {
        target_resource_sso      = true,
        target_resource_slo      = true,
        in_error_resource        = true,
        idp_discovery            = true,
        valid_domain             = "%[3]s",
        valid_path               = "/path",
        allow_query_and_fragment = true,
        require_https            = true
      },
      {
        target_resource_sso      = false,
        target_resource_slo      = true,
        in_error_resource        = true,
        idp_discovery            = true,
        valid_domain             = "anotherexample.com",
        valid_path               = "/path2",
        allow_query_and_fragment = false,
        require_https            = true
      }
    ]
  }
  redirect_validation_partner_settings = {
    enable_wreply_validation_slo = %[4]t
  }
}
data "pingfederate_redirect_validation" "%[1]s" {
  depends_on = [pingfederate_redirect_validation.%[1]s]
}
`, resourceName,
		resourceModel.enableTargetResourceValidationForSso,
		resourceModel.whiteListValidDomain,
		resourceModel.enableWreplyValidationSlo,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedRedirectValidationAttributes(config redirectValidationResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "RedirectValidation"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RedirectValidationAPI.GetRedirectValidationSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_target_resource_validation_for_sso", config.enableTargetResourceValidationForSso, response.RedirectValidationLocalSettings.GetEnableTargetResourceValidationForSSO())
		if err != nil {
			return err
		}

		whiteListValidDomain := response.RedirectValidationLocalSettings.WhiteList
		err = acctest.TestAttributesMatchString(resourceType, nil, "valid_domain", config.whiteListValidDomain, whiteListValidDomain[0].GetValidDomain())
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_wreply_validation_slo", config.enableWreplyValidationSlo, response.RedirectValidationPartnerSettings.GetEnableWreplyValidationSLO())
		if err != nil {
			return err
		}

		return nil
	}
}
