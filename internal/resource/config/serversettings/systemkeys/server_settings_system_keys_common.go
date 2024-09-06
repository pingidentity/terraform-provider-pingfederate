package serversettingssystemkeys

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serverSettingsSystemKeysModel struct {
	Id       types.String `tfsdk:"id"`
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
}
