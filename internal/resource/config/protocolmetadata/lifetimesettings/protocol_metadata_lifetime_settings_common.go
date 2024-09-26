package protocolmetadatalifetimesettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
)

type protocolMetadataLifetimeSettingsModel struct {
	Id            types.String `tfsdk:"id"`
	CacheDuration types.Int64  `tfsdk:"cache_duration"`
	ReloadDelay   types.Int64  `tfsdk:"reload_delay"`
}

func readProtocolMetadataLifetimeSettingsResponse(ctx context.Context, r *client.MetadataLifetimeSettings, state *protocolMetadataLifetimeSettingsModel) {
	state.CacheDuration = types.Int64Value(r.GetCacheDuration())
	state.ReloadDelay = types.Int64Value(r.GetReloadDelay())
}
