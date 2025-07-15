// Copyright © 2025 Ping Identity Corporation

package redirectvalidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
)

type redirectValidationModel struct {
	RedirectValidationLocalSettings   types.Object `tfsdk:"redirect_validation_local_settings"`
	RedirectValidationPartnerSettings types.Object `tfsdk:"redirect_validation_partner_settings"`
}

// Read a RedirectValidationResponse object into the model struct
func readRedirectValidationResponse(ctx context.Context, r *client.RedirectValidationSettings, state *redirectValidationModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	redirectValidationLocalSettingsObjVal, respDiags := types.ObjectValueFrom(ctx, redirectValidationLocalSettingsAttrTypes, r.RedirectValidationLocalSettings)
	diags.Append(respDiags...)
	redirectValidationPartnerSettingsObjVal, respDiags := types.ObjectValueFrom(ctx, redirectValidationPartnerSettingsAttrTypes, r.RedirectValidationPartnerSettings)
	diags.Append(respDiags...)
	state.RedirectValidationLocalSettings = redirectValidationLocalSettingsObjVal
	state.RedirectValidationPartnerSettings = redirectValidationPartnerSettingsObjVal
	return diags
}
