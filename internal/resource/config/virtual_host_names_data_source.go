package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualHostNamesDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualHostNamesDataSource{}
)

// VirtualHostNamesDataSource is a helper function to simplify the provider implementation.
func NewVirtualHostNamesDataSource() datasource.DataSource {
	return &virtualHostNamesDataSource{}
}

// virtualHostNamesDataSource is the datasource implementation.
type virtualHostNamesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type virtualHostNamesDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	VirtualHostNames types.List   `tfsdk:"virtual_host_names"`
}

// GetSchema defines the schema for the datasource.
func (r *virtualHostNamesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages VirtualHostNames.",
		Attributes: map[string]schema.Attribute{
			"virtual_host_names": schema.ListAttribute{
				Description: "List of virtual host names.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
			},
		},
	}
	id.ToSchema(&schemaDef)
	resp.Schema = schemaDef
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

// Read a VirtualHostNamesResponse object into the model struct
func readVirtualHostNamesResponseDataSource(ctx context.Context, r *client.VirtualHostNameSettings, state *virtualHostNamesDataSourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
	state.VirtualHostNames = internaltypes.GetStringList(r.VirtualHostNames)
	diags.Append(respDiags...)
	return diags
}

// Read the data source state and convert it into the model
func (r *virtualHostNamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualHostNamesDataSourceModel
	var diags diag.Diagnostics

	apiReadVirtualHostNames, httpResp, err := r.apiClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting Virtual Host Names", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Virtual Host Names", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadVirtualHostNames.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Virtual Host Names: %s", responseErr.Error())
	}

	// Read the response into the state
	var id = "virtual_host_names_id"
	diags = readVirtualHostNamesResponseDataSource(ctx, apiReadVirtualHostNames, &state, &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
