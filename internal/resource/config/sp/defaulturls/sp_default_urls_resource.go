package spdefaulturls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

func (model *spDefaultUrlsResource) buildDefaultClientStruct() *client.SpDefaultUrls {
	result := &client.SpDefaultUrls{}
	// confirm_slo
	result.ConfirmSlo = utils.Pointer(false)
	// slo_success_url
	result.SloSuccessUrl = utils.Pointer("")
	// sso_success_url
	result.SsoSuccessUrl = utils.Pointer("")
	return result
}

func (r *spDefaultUrlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	clientData := r.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrls(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.SpDefaultUrlsAPI.UpdateSpDefaultUrlsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the spDefaultUrls", err, httpResp)
	}
}
