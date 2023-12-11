package serversettingslogsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
)

type serverSettingsLogSettingsModel struct {
	Id            types.String `tfsdk:"id"`
	LogCategories types.Set    `tfsdk:"log_categories"`
}

func readServerSettingsLogSettingsResponse(ctx context.Context, r *client.LogSettings, state *serverSettingsLogSettingsModel, existingId *string) diag.Diagnostics {
	var diags diag.Diagnostics
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}

	state.LogCategories, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: logCategoriesAttrTypes}, r.LogCategories)
	return diags
}
