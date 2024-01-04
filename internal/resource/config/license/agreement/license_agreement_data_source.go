package licenseagreement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &licenseAgreementDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseAgreementDataSource{}
)

// Create a Administrative Account data source
func LicenseAgreementDataSource() datasource.DataSource {
	return &licenseAgreementDataSource{}
}

// licenseAgreementDataSource is the datasource implementation.
type licenseAgreementDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *licenseAgreementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a license agreement link.",
		Attributes: map[string]schema.Attribute{
			"license_agreement_url": schema.StringAttribute{
				Description: "URL to license agreement",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"accepted": schema.BoolAttribute{
				Description: "Indicates whether license agreement has been accepted. The default value is false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}

	id.ToDataSourceSchema(&schemaDef)
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

// Read resource information
func (r *licenseAgreementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseAgreementModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadLicenseAgreement, httpResp, err := r.apiClient.LicenseAPI.GetLicenseAgreement(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the license agreement", err, httpResp)
		return
	}

	// Read the response into the state
	readLicenseAgreementResponse(ctx, apiReadLicenseAgreement, &state, pointers.String("license_agreement_id"))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
