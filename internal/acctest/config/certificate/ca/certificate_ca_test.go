package certificateca_test

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

const certificateId = "test"

var fileData = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1")
var fileData2 = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2")

// Attributes to test with. Add optional properties to test here if desired.
type certificatesResourceModel struct {
	id       string
	fileData string
}

func TestAccCertificate(t *testing.T) {
	resourceName := "myCertificateCa"
	initialResourceModel := certificatesResourceModel{
		id:       certificateId,
		fileData: fileData,
	}
	updatedResourceModel := certificatesResourceModel{
		id:       certificateId,
		fileData: fileData2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileData == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1 must be set for acceptance tests")
			}
			if fileData2 == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2 must be set for acceptance tests")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCertificateAttributes(initialResourceModel),
			},
			{
				Config: testAccCertificate(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCertificateAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCertificate(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_certificate_ca." + resourceName,
				ImportStateId:     certificateId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.CertificatesCaAPI.DeleteTrustedCA(ctx, updatedResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccCertificate(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCertificateAttributes(initialResourceModel),
			},
		},
	})
}

func testAccCertificate(resourceName string, resourceModel certificatesResourceModel) string {
	// Not testing with crypto_provider attribute since it requires setting up an HSM
	return fmt.Sprintf(`
resource "pingfederate_certificate_ca" "%[1]s" {
  ca_id     = "%[2]s"
  file_data = "%[3]s"
}

data "pingfederate_certificate_ca" "%[1]s" {
  ca_id = pingfederate_certificate_ca.%[1]s.ca_id
}`, resourceName,
		resourceModel.id,
		fileData,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedCertificateAttributes(config certificatesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Certificate"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CertificatesCaAPI.GetTrustedCert(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "id",
			config.id, *response.Id)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckCertificateDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.CertificatesCaAPI.DeleteTrustedCA(ctx, certificateId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Certificate", certificateId)
	}
	return nil
}
