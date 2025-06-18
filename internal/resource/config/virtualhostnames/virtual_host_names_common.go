// Copyright © 2025 Ping Identity Corporation

package virtualhostnames

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

type virtualHostNamesModel struct {
	VirtualHostNames types.Set `tfsdk:"virtual_host_names"`
}

// Read a VirtualHostNamesResponse object into the model struct
func readVirtualHostNamesResponse(ctx context.Context, r *client.VirtualHostNameSettings, state *virtualHostNamesModel) {
	state.VirtualHostNames = internaltypes.GetStringSet(r.VirtualHostNames)
}
