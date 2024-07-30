package redirectvalidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type redirectValidationModel struct {
	Id                                types.String `tfsdk:"id"`
	RedirectValidationLocalSettings   types.Object `tfsdk:"redirect_validation_local_settings"`
	RedirectValidationPartnerSettings types.Object `tfsdk:"redirect_validation_partner_settings"`
}

// Read a RedirectValidationResponse object into the model struct
func readRedirectValidationResponse(ctx context.Context, r *client.RedirectValidationSettings, state *redirectValidationModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	redirectValidationLocalSettingsObjVal, respDiags := types.ObjectValueFrom(ctx, redirectValidationLocalSettingsAttrTypes, r.RedirectValidationLocalSettings)
	diags.Append(respDiags...)
	redirectValidationPartnerSettingsObjVal, respDiags := types.ObjectValueFrom(ctx, redirectValidationPartnerSettingsAttrTypes, r.RedirectValidationPartnerSettings)
	diags.Append(respDiags...)
	state.RedirectValidationLocalSettings = redirectValidationLocalSettingsObjVal
	state.RedirectValidationPartnerSettings = redirectValidationPartnerSettingsObjVal
	return diags
}
