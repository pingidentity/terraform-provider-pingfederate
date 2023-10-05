package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &licenseDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseDataSource{}
)

// Create a Administrative Account data source
func NewLicenseDataSource() datasource.DataSource {
	return &licenseDataSource{}
}

// licenseDataSource is the datasource implementation.
type licenseDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type licenseDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	FileData types.String `tfsdk:"file_data"`
}

// GetSchema defines the schema for the datasource.
func (r *licenseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a License.",
		Attributes: map[string]schema.Attribute{
			"file_data": schema.StringAttribute{
				Required: true,
			},
		},
	}

	config.AddCommonDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *licenseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

// Configure adds the provider configured client to the data source.
func (r *licenseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLicenseResponseDataSource(ctx context.Context, r *client.LicenseView, state *licenseDataSourceModel, expectedValues *licenseDataSourceModel, createPlan types.String) {
	LicenseFileData := createPlan
	state.Id = types.StringValue("id")
	state.FileData = types.StringValue(LicenseFileData.ValueString())
}

// Read resource information
func (r *licenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadLicense, httpResp, err := r.apiClient.LicenseAPI.GetLicense(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the License", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadLicense.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the request of the License: %s", responseErr.Error())
	}

	// Read the response into the state
	readLicenseResponseDataSource(ctx, apiReadLicense, &state, &state, state.FileData)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
