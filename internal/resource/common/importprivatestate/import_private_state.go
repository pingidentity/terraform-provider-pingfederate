package importprivatestate

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

type importPrivateState struct {
	IsImport bool `json:"isImport"`
}

const importKey = "import"

func MarkPrivateStateForImport(ctx context.Context, resp *resource.ImportStateResponse) diag.Diagnostics {
	value := []byte(`{"isImport": true}`)
	return resp.Private.SetKey(ctx, importKey, value)
}

func IsImportRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) (bool, diag.Diagnostics) {
	var respDiags diag.Diagnostics
	importRead, diags := req.Private.GetKey(ctx, importKey)
	respDiags.Append(diags...)
	// Reset the value to nil
	resp.Private.SetKey(ctx, importKey, nil)

	isImport := false
	if importRead != nil {
		var target importPrivateState
		err := json.Unmarshal(importRead, &target)
		if err != nil {
			respDiags.AddError(providererror.InternalProviderError, "Failed to unmarshal import state: "+err.Error())
		} else {
			isImport = target.IsImport
		}
	}
	return isImport, respDiags
}
