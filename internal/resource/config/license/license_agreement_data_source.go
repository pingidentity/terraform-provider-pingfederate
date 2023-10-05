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
	_ datasource.DataSource              = &licenseAgreementDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseAgreementDataSource{}
)

// Create a Administrative Account data source
func NewLicenseAgreementDataSource() datasource.DataSource {
	return &licenseAgreementDataSource{}
}

// licenseAgreementDataSource is the datasource implementation.
type licenseAgreementDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type licenseAgreementDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LicenseAgreementUrl types.String `tfsdk:"license_agreement_url"`
	Accepted            types.Bool   `tfsdk:"accepted"`
}

// GetSchema defines the schema for the datasource.
func (r *licenseAgreementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a LicenseAgreement.",
		Attributes: map[string]schema.Attribute{
			"license_agreement_url": schema.StringAttribute{
				Description: "URL to license agreement",
				Optional:    true,
				Computed:    true,
			},
			"accepted": schema.BoolAttribute{
				Description: "Indicates whether license agreement has been accepted. The default value is false.",
				Optional:    true,
				Computed:    true,
			},
		},
	}

	config.AddCommonDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *licenseAgreementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_agreement"
}

// Configure adds the provider configured client to the data source.
func (r *licenseAgreementDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLicenseAgreementResponseDataSource(ctx context.Context, r *client.LicenseAgreementInfo, state *licenseAgreementDataSourceModel, expectedValues *licenseAgreementDataSourceModel) {
	state.Id = types.StringValue("id")
	state.LicenseAgreementUrl = internaltypes.StringTypeOrNil(r.LicenseAgreementUrl, false)
	state.Accepted = internaltypes.BoolTypeOrNil(r.Accepted)
}

// Read resource information
func (r *licenseAgreementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseAgreementDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadLicenseAgreement, httpResp, err := r.apiClient.LicenseAPI.GetLicenseAgreement(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the License Agreement", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadLicenseAgreement.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of the License Agreement: %s", responseErr.Error())
	}

	// Read the response into the state
	readLicenseAgreementResponseDataSource(ctx, apiReadLicenseAgreement, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
