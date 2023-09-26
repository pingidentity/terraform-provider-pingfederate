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

const keyPairsSslServerImportId = "2"
const fileData = "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"
const password = "2FederateM0re"
const format = "PKCS12"

// Attributes to test with. Add optional properties to test here if desired.
type keyPairsSslServerImportResourceModel struct {
	id       string
	fileData string
	format   string
	password string
}

func TestAccKeyPairsSslServerImport(t *testing.T) {
	resourceName := "myKeyPairsSslServerImport"
	initialResourceModel := keyPairsSslServerImportResourceModel{
		id:       keyPairsSslServerImportId,
		fileData: fileData,
		format:   format,
		password: password,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckKeyPairsSslServerImportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairsSslServerImport(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedKeyPairsSslServerImportAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccKeyPairsSslServerImport(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_key_pair_ssl_server_import." + resourceName,
				ImportStateId:     keyPairsSslServerImportId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccKeyPairsSslServerImport(resourceName string, resourceModel keyPairsSslServerImportResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_key_pair_ssl_server_import" "%[1]s" {
  custom_id = "%[2]s"
  file_data = "%[3]s"
  format    = "%[4]s"
  password  = "%[5]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.fileData,
		resourceModel.format,
		resourceModel.password,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedKeyPairsSslServerImportAttributes(config keyPairsSslServerImportResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "KeyPairsSslServerImport"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.KeyPairsSslServerApi.GetSslServerKeyPair(ctx, keyPairsSslServerImportId).Execute()
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "id", config.id, *response.Id)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckKeyPairsSslServerImportDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.KeyPairsSslClientApi.DeleteSslClientKeyPair(ctx, keyPairsSslServerImportId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("KeyPairsSslServerImport", keyPairsSslServerImportId)
	}
	return nil
}
