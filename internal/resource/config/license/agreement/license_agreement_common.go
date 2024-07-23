package licenseagreement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type licenseAgreementModel struct {
	Id                  types.String `tfsdk:"id"`
	LicenseAgreementUrl types.String `tfsdk:"license_agreement_url"`
	Accepted            types.Bool   `tfsdk:"accepted"`
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLicenseAgreementResponse(ctx context.Context, r *client.LicenseAgreementInfo, state *licenseAgreementModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.LicenseAgreementUrl = types.StringPointerValue(r.LicenseAgreementUrl)
	state.Accepted = types.BoolPointerValue(r.Accepted)
}
