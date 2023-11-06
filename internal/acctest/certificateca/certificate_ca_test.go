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

const certificateId = "test"
const fileData = "MIIDmjCCAoICCQDncp3LMAO6YjANBgkqhkiG9w0BAQsFADCBjjELMAkGA1UEBhMCVVMxDDAKBgNVBAgMA1NKQzERMA8GA1UEBwwIc2FuIEpvc2UxFjAUBgNVBAoMDXBpbmcgSWRlbnRpdHkxDzANBgNVBAsMBkRldm9wczEWMBQGA1UEAwwNdGVycmFmb3JtdGVzdDEdMBsGCSqGSIb3DQEJARYOdGVzdEBnbWFpbC5jb20wHhcNMjMwNTMwMTU1OTE5WhcNMjQwNTI5MTU1OTE5WjCBjjELMAkGA1UEBhMCVVMxDDAKBgNVBAgMA1NKQzERMA8GA1UEBwwIc2FuIEpvc2UxFjAUBgNVBAoMDXBpbmcgSWRlbnRpdHkxDzANBgNVBAsMBkRldm9wczEWMBQGA1UEAwwNdGVycmFmb3JtdGVzdDEdMBsGCSqGSIb3DQEJARYOdGVzdEBnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDB7u+oHHQgGrZdCk74A4XJzjzhMT9MN1MJIqar+96rKogDmt3LnCh+oN5hxy0QPjrW9SiRHPZME+e6YWtBNfg21KDws2nLoH/eGmb45ObM/nApX4oFZD06ccW4zWjxuxEdKzKAMWMP60UxCZwnK99cIRMYs0x85lHhcLfTuA3VAwg95X+2FxQDk8sAdNdl1zhWaR2YS+nrmP/iheG2fT8cVLTGdklPqL9nrUDAwwUyX5I8PLsLPzJzMoXV+on4zjypNxfXt2MmuLHOGxwgxvUVRiVeCTSMo1y763OUAnds1L+uJNq1vvsD0iFwyA78I3EzaX9c5Vxhbk+3JKFD1gY1AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAGqlkRIgsAFE6/WBayYlsITtnxJooTJyZ8CHFulRMskMYdoETYUeN5FqmJ05PGUHgXX0/3fQ9RYD3Mfuupm1Vqgx8q/v5cIrBefU7zW3bjy/BMAONkPAr617NkbHAj2XC1t5YFr6Vnnx9JQoIl70slBGABPwSkahrReE5f87qkkWqVI8aiuAzu0GRkMHbv1XzGfXfVF/iK9Lq6x80tyiqL987Krw6hHPlxS4GXjwvWWO0f0GfNwENxSv6uwxvCFIp01x7LHbkPHJvMH2Z5wSZges5ZDv/rciunSZ2xYh/jGzM1gIz29DBpmayl4AwKi5/ix7p3ujCA1jdlT+nlBZ/js="
const fileData2 = "MIICMzCCAZygAwIBAgIJALiPnVsvq8dsMA0GCSqGSIb3DQEBBQUAMFMxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNVBAcTA2ZvbzEMMAoGA1UEChMDZm9vMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2ZvbzAeFw0xMzAzMTkxNTQwMTlaFw0xODAzMTgxNTQwMTlaMFMxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNVBAcTA2ZvbzEMMAoGA1UEChMDZm9vMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2ZvbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAzdGfxi9CNbMf1UUcvDQh7MYBOveIHyc0E0KIbhjK5FkCBU4CiZrbfHagaW7ZEcN0tt3EvpbOMxxc/ZQU2WN/s/wPxph0pSfsfFsTKM4RhTWD2v4fgk+xZiKd1p0+L4hTtpwnEw0uXRVd0ki6muwV5y/P+5FHUeldq+pgTcgzuK8CAwEAAaMPMA0wCwYDVR0PBAQDAgLkMA0GCSqGSIb3DQEBBQUAA4GBAJiDAAtY0mQQeuxWdzLRzXmjvdSuL9GoyT3BF/jSnpxz5/58dba8pWenv3pj4P3w5DoOso0rzkZy2jEsEitlVM2mLSbQpMM+MUVQCQoiG6W9xuCFuxSrwPISpAqEAuV4DNoxQKKWmhVv+J0ptMWD25Pnpxeq5sXzghfJnslJlQND+2kmOeUJXRmm/kEd5jhW6Y7qj/WsjTVbJmcVfewCHrPSqnI0kBBIZCe/zuf6IWUrVnZ9NA2zsmWLIodz2uFHdh1voqZiegDfqnc1zqcPGUIWVEX/r87yloqaKHee9570+sB3c4"

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
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
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
		},
	})
}

func testAccCertificate(resourceName string, resourceModel certificatesResourceModel) string {
	// Not testing with crypto_provider attribute since it requires setting up an HSM
	return fmt.Sprintf(`
resource "pingfederate_certificate_ca" "%[1]s" {
  custom_id = "%[2]s"
  file_data = "%[3]s"
}

data "pingfederate_certificate_ca" "%[1]s" {
  id = pingfederate_certificate_ca.%[1]s.custom_id
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
