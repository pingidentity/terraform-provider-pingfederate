package certificatesrevocationsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Create a client struct with the PF-default certificate revocation settings
func (model *certificatesRevocationSettingsResourceModel) buildDefaultClientStruct() *client.CertificateRevocationSettings {
	return &client.CertificateRevocationSettings{
		OcspSettings: &client.OcspSettings{
			RequesterAddNonce:            utils.Pointer(false),
			ActionOnResponderUnavailable: utils.Pointer("CONTINUE"),
			ActionOnStatusUnknown:        utils.Pointer("FAIL"),
			ActionOnUnsuccessfulResponse: utils.Pointer("FAIL"),
			CurrentUpdateGracePeriod:     utils.Pointer(int64(5)),
			NextUpdateGracePeriod:        utils.Pointer(int64(5)),
			ResponseCachePeriod:          utils.Pointer(int64(48)),
			ResponderTimeout:             utils.Pointer(int64(5)),
			ResponderUrl:                 utils.Pointer(""),
		},
	}
}

func (r *certificatesRevocationSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model certificatesRevocationSettingsResourceModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the certificatesRevocationSettings", err, httpResp)
	}
}
