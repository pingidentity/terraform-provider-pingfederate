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

const idpDefaultUrlsId = "id"
const idpErrMessage = "errorDetail.idpSsoFailure"

// Attributes to test with. Add optional properties to test here if desired.
type idpDefaultUrlsResourceModel struct {
	id               string
	confirmIdpSlo    bool
	idpSloSuccessUrl string
	idpErrorMsg      string
}

func TestAccIdpDefaultUrls(t *testing.T) {
	resourceName := "myIdpDefaultUrls"
	initialResourceModel := idpDefaultUrlsResourceModel{
		id:               idpDefaultUrlsId,
		confirmIdpSlo:    true,
		idpSloSuccessUrl: "https://localhost",
		idpErrorMsg:      idpErrMessage,
	}
	updatedResourceModel := idpDefaultUrlsResourceModel{
		id:               idpDefaultUrlsId,
		confirmIdpSlo:    false,
		idpSloSuccessUrl: "https://example",
		idpErrorMsg:      idpErrMessage,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIdpDefaultUrls(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpDefaultUrlsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpDefaultUrls(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpDefaultUrlsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIdpDefaultUrls(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_idp_default_urls." + resourceName,
				ImportStateId:     idpDefaultUrlsId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdpDefaultUrls(resourceName string, resourceModel idpDefaultUrlsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_default_urls" "%[1]s" {
  confirm_idp_slo     = %[2]t
  idp_error_msg       = "%[3]s"
  idp_slo_success_url = "%[4]s"
}`, resourceName,
		resourceModel.confirmIdpSlo,
		resourceModel.idpErrorMsg,
		resourceModel.idpSloSuccessUrl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpDefaultUrlsAttributes(config idpDefaultUrlsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IdpDefaultUrls"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IdpDefaultUrlsAPI.GetDefaultUrl(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "confirm_idp_slo",
			config.confirmIdpSlo, *response.ConfirmIdpSlo)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "idp_slo_success_url",
			config.idpSloSuccessUrl, *response.IdpSloSuccessUrl)
		if err != nil {
			return err
		}

		return nil
	}
}
