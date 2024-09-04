package importprivatestate

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type importPrivateState struct {
	IsImport bool `json:"isImport"`
}

func MarkPrivateStateForImport(ctx context.Context, resp *resource.ImportStateResponse) diag.Diagnostics {
	value := []byte(`{"isImport": true}`)
	return resp.Private.SetKey(ctx, "import", value)
}

func IsImportRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) (bool, diag.Diagnostics) {
	var respDiags diag.Diagnostics
	importRead, diags := req.Private.GetKey(ctx, "import")
	respDiags.Append(diags...)

	isImport := false
	if importRead != nil {
		var target importPrivateState
		err := json.Unmarshal(importRead, &target)
		if err != nil {
			respDiags.AddError("Failed to unmarshal import state", err.Error())
		} else {
			isImport = target.IsImport
		}
	} else {
		respDiags.AddWarning("This is a normal read!", "")
	}
	return isImport, respDiags
}
