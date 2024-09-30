package serviceauthentication

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
)

func (state *serviceAuthenticationResourceModel) readClientResponseSharedSecret(existingParentValue types.Object) types.String {
	if existingParentValue.IsNull() || existingParentValue.IsUnknown() {
		return types.StringNull()
	}

	// Get the existing sharedSecret value
	attrs := existingParentValue.Attributes()
	sharedSecret, ok := attrs["shared_secret"]
	if !ok {
		return types.StringNull()
	}

	return types.StringValue(sharedSecret.(types.String).ValueString())
}

func (r *serviceAuthenticationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this resource will be reset to the PingFederate default values.
	apiUpdateRequest := r.apiClient.ServiceAuthenticationAPI.UpdateServiceAuthentication(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(client.ServiceAuthentication{})
	_, httpResp, err := r.apiClient.ServiceAuthenticationAPI.UpdateServiceAuthenticationExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the serviceAuthentication", err, httpResp)
	}
}
