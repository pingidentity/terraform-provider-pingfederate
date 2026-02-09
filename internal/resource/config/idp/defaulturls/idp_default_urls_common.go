// Copyright © 2026 Ping Identity Corporation

package idpdefaulturls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type idpDefaultUrlsModel struct {
	Id               types.String `tfsdk:"id"`
	ConfirmIdpSlo    types.Bool   `tfsdk:"confirm_idp_slo"`
	IdpSloSuccessUrl types.String `tfsdk:"idp_slo_success_url"`
	IdpErrorMsg      types.String `tfsdk:"idp_error_msg"`
}

// Read a IdpDefaultUrlsResponse object into the model struct
func readIdpDefaultUrlsResponse(ctx context.Context, r *client.IdpDefaultUrl, state *idpDefaultUrlsModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.ConfirmIdpSlo = types.BoolPointerValue(r.ConfirmIdpSlo)
	state.IdpSloSuccessUrl = types.StringPointerValue(r.IdpSloSuccessUrl)
	state.IdpErrorMsg = types.StringValue(r.IdpErrorMsg)
}
