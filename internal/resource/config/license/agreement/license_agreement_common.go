// Copyright © 2025 Ping Identity Corporation

package licenseagreement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

type licenseAgreementModel struct {
	LicenseAgreementUrl types.String `tfsdk:"license_agreement_url"`
	Accepted            types.Bool   `tfsdk:"accepted"`
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLicenseAgreementResponse(ctx context.Context, r *client.LicenseAgreementInfo, state *licenseAgreementModel) {
	state.LicenseAgreementUrl = types.StringPointerValue(r.LicenseAgreementUrl)
	state.Accepted = types.BoolPointerValue(r.Accepted)
}
