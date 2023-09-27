package acctest_test

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

const redirectValidationId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type redirectValidationResourceModel struct {
	id                                   string
	enableTargetResourceValidationForSso bool
	whiteListValidDomain                 string
	enableWreplyValidationSlo            bool
}

func TestAccRedirectValidation(t *testing.T) {
	resourceName := "myRedirectValidation"
	initialResourceModel := redirectValidationResourceModel{
		id:                                   redirectValidationId,
		enableTargetResourceValidationForSso: true,
		whiteListValidDomain:                 "example.com",
		enableWreplyValidationSlo:            false,
	}
	updatedResourceModel := redirectValidationResourceModel{
		id:                                   redirectValidationId,
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
				Config: testAccRedirectValidation(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedRedirectValidationAttributes(initialResourceModel),
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
				ImportStateId:     redirectValidationId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedirectValidation(resourceName string, resourceModel redirectValidationResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_redirect_validation" "%[1]s" {
  redirect_validation_local_settings = {
    enable_target_resource_validation_for_sso = %[2]t
    white_list = [
      {
        valid_domain = "%[3]s"
      }
    ]
  }
  redirect_validation_partner_settings = {
    enable_wreply_validation_slo = %[4]t
  }
}`, resourceName,
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
		response, _, err := testClient.RedirectValidationApi.GetRedirectValidationSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enable_target_resource_validation_for_sso", config.enableTargetResourceValidationForSso, response.RedirectValidationLocalSettings.GetEnableTargetResourceValidationForSSO())
		if err != nil {
			return err
		}

		whiteListValidDomain := response.RedirectValidationLocalSettings.WhiteList
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "valid_domain", config.whiteListValidDomain, whiteListValidDomain[0].GetValidDomain())
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enable_wreply_validation_slo", config.enableWreplyValidationSlo, response.RedirectValidationPartnerSettings.GetEnableWreplyValidationSLO())
		if err != nil {
			return err
		}

		return nil
	}
}
