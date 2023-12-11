package virtualhostnames

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

type virtualHostNamesModel struct {
	Id               types.String `tfsdk:"id"`
	VirtualHostNames types.List   `tfsdk:"virtual_host_names"`
}

// Read a VirtualHostNamesResponse object into the model struct
func readVirtualHostNamesResponse(ctx context.Context, r *client.VirtualHostNameSettings, state *virtualHostNamesModel, existingId *string) {
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.VirtualHostNames = internaltypes.GetStringList(r.VirtualHostNames)
}
