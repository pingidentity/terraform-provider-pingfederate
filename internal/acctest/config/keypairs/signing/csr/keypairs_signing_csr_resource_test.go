package keypairssigningcsr_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

var fileDataInitial, fileDataUpdated, fileDataCa string

const signingCaId = "singingcsrtestca"

func TestAccKeypairsSigningCsrResource(t *testing.T) {
	fileDataInitial = os.Getenv("PF_TF_ACC_TEST_CSR_RESPONSE_1")
	fileDataUpdated = os.Getenv("PF_TF_ACC_TEST_CSR_RESPONSE_2")
	fileDataCa = os.Getenv("PF_TF_ACC_TEST_CA_CERTIFICATE")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileDataInitial == "" {
				t.Fatal("PF_TF_ACC_TEST_CSR_RESPONSE_1 must be set for TestAccKeypairsSigningCsrResource")
			}
			if fileDataUpdated == "" {
				t.Fatal("PF_TF_ACC_TEST_CSR_RESPONSE_2 must be set for TestAccKeypairsSigningCsrResource")
			}
			if fileDataCa == "" {
				t.Fatal("PF_TF_ACC_TEST_CA_CERTIFICATE must be set for TestAccKeypairsSigningCsrResource")
			}
			keypairsSigningCsr_ImportCA(t)
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsSigningCsr_DeleteCA(),
		Steps: []resource.TestStep{
			{
				// Initial CSR response import on create
				Config: keypairsSigningCsr_HCL(fileDataInitial),
				Check:  keypairsSigningCsr_CheckComputedValuesInitial(),
			},
			{
				// Importing a second CSR response
				Config: keypairsSigningCsr_HCL(fileDataUpdated),
				Check:  keypairsSigningCsr_CheckComputedValuesUpdated(),
			},
		},
	})
}

func keypairsSigningCsr_ImportCA(t *testing.T) {
	testClient := acctest.TestClient()
	trustCaImportReq := testClient.CertificatesCaAPI.ImportTrustedCA(acctest.TestBasicAuthContext())
	trustCaImportReq = trustCaImportReq.Body(client.X509File{
		Id:       utils.Pointer(signingCaId),
		FileData: fileDataCa,
	})
	_, httpResp, err := testClient.CertificatesCaAPI.ImportTrustedCAExecute(trustCaImportReq)
	if err != nil {
		errorMsg := "Failed to import test CA: " + err.Error()
		if httpResp != nil {
			body, internalErr := io.ReadAll(httpResp.Body)
			if internalErr == nil {
				errorMsg += " - Detail: " + string(body)
			}
		}
		t.Error(errorMsg)
		t.FailNow()
	}
}

func keypairsSigningCsr_DeleteCA() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		_, err := testClient.CertificatesCaAPI.DeleteTrustedCA(acctest.TestBasicAuthContext(), signingCaId).Execute()
		return err
	}
}

func keypairsSigningCsr_HCL(fileData string) string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_csr" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  file_data  = "%s"
}
`, fileData)
}

func keypairsSigningCsr_CheckComputedValuesInitial() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_signing_csr.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "serial_number", "35870055780717650058227469919152395501"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "subject_dn", "CN=common, O=org, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "valid_from", "2024-07-29T15:57:40Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "expires", "2025-07-29T15:57:40Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "sha1_fingerprint", "3A34FEC4210B152AFDF1192B088E012E8475AE61"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "sha256_fingerprint", "294460C52A238B0BE701FFC0BAD142548F19C7CC6C83F2BD3982291CC0624053"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "status", "VALID"),
	)
}

func keypairsSigningCsr_CheckComputedValuesUpdated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_signing_csr.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "serial_number", "78860249853500415650095464700202533503"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "subject_dn", "CN=common, O=org, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "valid_from", "2024-07-29T16:46:30Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "expires", "2025-07-29T16:46:30Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "sha1_fingerprint", "F26E602557E3B7DFA7444904E4A28EAF94FD4F63"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "sha256_fingerprint", "F5C8404FA236325ED89C8814BE59627D0696388F6A20C1C691AE0300E46147A0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr.example", "status", "VALID"),
	)
}
