package kerberosrealmssettings

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

func (state *kerberosRealmSettingsResourceModel) readClientResponse(response *client.KerberosRealmsSettings, checkForIgnoredRequestAttrs bool) diag.Diagnostics {
	var respDiags diag.Diagnostics
	if checkForIgnoredRequestAttrs {
		respDiags = state.reportPfIgnoredAttrs(response)
	}
	// debug_log_output
	state.DebugLogOutput = types.BoolPointerValue(response.DebugLogOutput)
	// force_tcp
	state.ForceTcp = types.BoolPointerValue(response.ForceTcp)
	// kdc_retries
	if response.KdcRetries == "" {
		response.KdcRetries = "0"
	}
	kdcRetriesInt, err := strconv.ParseInt(response.KdcRetries, 10, 64)
	if err != nil {
		respDiags.AddError(providererror.InternalProviderError, "Failed to parse kdc_retries as int: "+err.Error())
	}
	state.KdcRetries = types.Int64Value(kdcRetriesInt)
	// kdc_timeout
	if response.KdcTimeout == "" {
		response.KdcTimeout = "0"
	}
	kdcTimeoutInt, err := strconv.ParseInt(response.KdcTimeout, 10, 64)
	if err != nil {
		respDiags.AddError(providererror.InternalProviderError, "Failed to parse kdc_timeout as int: "+err.Error())
	}
	state.KdcTimeout = types.Int64Value(kdcTimeoutInt)
	// key_set_retention_period_mins
	state.KeySetRetentionPeriodMins = types.Int64PointerValue(response.KeySetRetentionPeriodMins)
	return respDiags
}

func (state *kerberosRealmSettingsResourceModel) reportPfIgnoredAttrs(response *client.KerberosRealmsSettings) diag.Diagnostics {
	var respDiags diag.Diagnostics
	// If the response from PF ignored any of the values passed in, report an error
	if response.DebugLogOutput != nil && !state.DebugLogOutput.IsNull() {
		addPfIgnoredError(&respDiags, "debug_log_output", strconv.FormatBool(*response.DebugLogOutput), strconv.FormatBool(state.DebugLogOutput.ValueBool()))
	}
	if response.ForceTcp != nil && !state.ForceTcp.IsNull() {
		addPfIgnoredError(&respDiags, "force_tcp", strconv.FormatBool(*response.ForceTcp), strconv.FormatBool(state.ForceTcp.ValueBool()))
	}
	addPfIgnoredError(&respDiags, "kdc_retries", response.KdcRetries, strconv.FormatInt(state.KdcRetries.ValueInt64(), 10))
	addPfIgnoredError(&respDiags, "kdc_timeout", response.KdcTimeout, strconv.FormatInt(state.KdcTimeout.ValueInt64(), 10))
	if response.KeySetRetentionPeriodMins != nil && !state.KeySetRetentionPeriodMins.IsNull() {
		addPfIgnoredError(&respDiags, "key_set_retention_period_mins", strconv.FormatInt(*response.KeySetRetentionPeriodMins, 10), strconv.FormatInt(state.KeySetRetentionPeriodMins.ValueInt64(), 10))
	}
	return respDiags
}

func addPfIgnoredError(respDiags *diag.Diagnostics, attrName, responseValue, requestValue string) {
	if requestValue != responseValue {
		respDiags.AddError(providererror.InvalidResourceConfiguration,
			fmt.Sprintf("PingFederate failed to save the provided value for %[1]s. "+
				"Ensure you have a kerberos realm configured before modifying the pingfederate_kerberos_realm_settings.\n"+
				"PingFederate returned \"%[2]s\" after request to set %[1]s to \"%[3]s\"", attrName, responseValue, requestValue))
	}
}
