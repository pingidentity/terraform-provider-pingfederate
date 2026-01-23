// Copyright © 2025 Ping Identity Corporation

package keypairssslclientcsr_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/config/keypairs"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

const signingCaId = "sslclientcsrtestca"

func TestAccKeypairsSslClientCsrResponseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			keypairsSslClientCsrResponse_ImportCA(t)
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Initial CSR response import on create
				Config: keypairsSslClientCsrResponse_HCL(keypairs.CsrResponse1),
				Check:  keypairsSslClientCsrResponse_CheckComputedValuesInitial(),
			},
			{
				// Importing a second CSR response
				Config: keypairsSslClientCsrResponse_HCL(keypairs.CsrResponse2),
				Check:  keypairsSslClientCsrResponse_CheckComputedValuesUpdated(),
			},
		},
	})
}

func keypairsSslClientCsrResponse_ImportCA(t *testing.T) {
	testClient := acctest.TestClient()
	trustCaImportReq := testClient.CertificatesCaAPI.ImportTrustedCA(acctest.TestBasicAuthContext())
	trustCaImportReq = trustCaImportReq.Body(client.X509File{
		Id:       utils.Pointer(signingCaId),
		FileData: keypairs.TestCertificateAuthority,
	})
	_, httpResp, err := testClient.CertificatesCaAPI.ImportTrustedCAExecute(trustCaImportReq)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 422) {
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

func keypairsSslClientCsrResponse_HCL(fileData string) string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_ssl_client_csr_response" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  file_data  = "%s"
}
`, fileData)
}

func keypairsSslClientCsrResponse_CheckComputedValuesInitial() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "serial_number", "8770436850221969673930962783985572962"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "subject_dn", "CN=common, O=org, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "valid_from", "2025-07-30T19:13:37Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "expires", "2045-07-25T19:13:37Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "sha1_fingerprint", "3449766B67B306DE781CE2C9D4C23D527D13EB04"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "sha256_fingerprint", "6B81F1BC63C36E41A949B4A3715135939A4511EA9658FA5BA57944AE5DCA437E"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "status", "VALID"),
	)
}

func keypairsSslClientCsrResponse_CheckComputedValuesUpdated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "serial_number", "99377934021464408054402725272426356391"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "subject_dn", "CN=common, O=org, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "valid_from", "2025-07-30T19:21:47Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "expires", "2045-07-25T19:21:47Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "sha1_fingerprint", "4975FB3948D3FBD270CA25E123931BF114840EAF"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "sha256_fingerprint", "E88E5596EF241003FD2E2F702A7D832CC7F8D3FBF68F38FDF47645D83387C9C1"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_response.example", "status", "VALID"),
	)
}
