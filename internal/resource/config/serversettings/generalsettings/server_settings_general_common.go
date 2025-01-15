package serversettingsgeneralsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
)

type serverSettingsGeneralModel struct {
	DisableAutomaticConnectionValidation    types.Bool   `tfsdk:"disable_automatic_connection_validation"`
	IdpConnectionTransactionLoggingOverride types.String `tfsdk:"idp_connection_transaction_logging_override"`
	SpConnectionTransactionLoggingOverride  types.String `tfsdk:"sp_connection_transaction_logging_override"`
	DatastoreValidationIntervalSecs         types.Int64  `tfsdk:"datastore_validation_interval_secs"`
	RequestHeaderForCorrelationId           types.String `tfsdk:"request_header_for_correlation_id"`
}

func readServerSettingsGeneralResponse(ctx context.Context, r *client.GeneralSettings, state *serverSettingsGeneralModel) {
	state.DisableAutomaticConnectionValidation = types.BoolPointerValue(r.DisableAutomaticConnectionValidation)
	state.IdpConnectionTransactionLoggingOverride = types.StringPointerValue(r.IdpConnectionTransactionLoggingOverride)
	state.SpConnectionTransactionLoggingOverride = types.StringPointerValue(r.SpConnectionTransactionLoggingOverride)
	state.DatastoreValidationIntervalSecs = types.Int64PointerValue(r.DatastoreValidationIntervalSecs)
	state.RequestHeaderForCorrelationId = types.StringPointerValue(r.RequestHeaderForCorrelationId)
}
