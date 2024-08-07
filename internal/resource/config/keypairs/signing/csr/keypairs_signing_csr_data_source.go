package keypairssigningcsr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ datasource.DataSource              = &keypairsSigningCsrDataSource{}
	_ datasource.DataSourceWithConfigure = &keypairsSigningCsrDataSource{}
)

func KeypairsSigningCsrDataSource() datasource.DataSource {
	return &keypairsSigningCsrDataSource{}
}

type keypairsSigningCsrDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSigningCsrDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_signing_csr"
}

func (r *keypairsSigningCsrDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSigningCsrDataSourceModel struct {
	KeypairId   types.String `tfsdk:"keypair_id"`
	ExportedCsr types.String `tfsdk:"exported_csr"`
}

func (r *keypairsSigningCsrDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to generate a new certificate signing request (CSR) for a key pair.",
		Attributes: map[string]schema.Attribute{
			"keypair_id": schema.StringAttribute{
				Description: "The ID of the keypair.",
				Required:    true,
			},
			"exported_csr": schema.StringAttribute{
				Description: "The exported PEM-encoded certificate signing request.",
				Computed:    true,
			},
		},
	}
}

func (r *keypairsSigningCsrDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keypairsSigningCsrDataSourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	exportRequest := r.apiClient.KeyPairsSigningAPI.ExportCsr(config.AuthContext(ctx, r.providerConfig), data.KeypairId.ValueString())
	responseData, httpResp, err := exportRequest.Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while generating the certificate signing request.", err, httpResp)
		return
	}

	// Set the exported metadata
	data.ExportedCsr = types.StringValue(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
