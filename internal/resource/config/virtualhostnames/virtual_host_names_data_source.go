// Copyright © 2025 Ping Identity Corporation

package virtualhostnames

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualHostNamesDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualHostNamesDataSource{}
)

// VirtualHostNamesDataSource is a helper function to simplify the provider implementation.
func VirtualHostNamesDataSource() datasource.DataSource {
	return &virtualHostNamesDataSource{}
}

// virtualHostNamesDataSource is the datasource implementation.
type virtualHostNamesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *virtualHostNamesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes settings for virtual host names.",
		Attributes: map[string]schema.Attribute{
			"virtual_host_names": schema.SetAttribute{
				Description: "List of virtual host names.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    false,
			},
		},
	}
	resp.Schema = schema
}

// Metadata returns the data source type name.
func (r *virtualHostNamesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_host_names"
}

// Configure adds the provider configured client to the data source.
func (r *virtualHostNamesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

// Read the data source state and convert it into the model
func (r *virtualHostNamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualHostNamesModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadVirtualHostNames, httpResp, err := r.apiClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Virtual Host Names", err, httpResp)
		return
	}

	// Read the response into the state
	readVirtualHostNamesResponse(ctx, apiReadVirtualHostNames, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
