package serversettingssystemkeys

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type serverSettingsSystemKeysModel struct {
	Id       types.String `tfsdk:"id"`
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
}

func readServerSettingsSystemKeysResponse(ctx context.Context, r *client.SystemKeys, state *serverSettingsSystemKeysModel, existingId *string) diag.Diagnostics {
	var diags diag.Diagnostics
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	currentAttrs := r.GetCurrent()
	currentAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(currentAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(currentAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(currentAttrs.GetKeyData()),
	}
	currentAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, currentAttrVals)
	diags = append(diags, respDiags...)

	previousAttrs := r.GetPrevious()
	previousAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(previousAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(previousAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(previousAttrs.GetKeyData()),
	}
	previousAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, previousAttrVals)
	diags = append(diags, respDiags...)

	pendingAttrs := r.GetPending()
	pendingAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(pendingAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(pendingAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(pendingAttrs.GetKeyData()),
	}
	pendingAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, pendingAttrVals)
	diags = append(diags, respDiags...)

	state.Current = currentAttrsObjVal
	state.Pending = pendingAttrsObjVal
	state.Previous = previousAttrsObjVal
	return diags
}
