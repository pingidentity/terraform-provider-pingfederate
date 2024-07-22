package certificatesca

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
	_ datasource.DataSource              = &certificatesCAExportDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesCAExportDataSource{}
)

func CertificatesCAExportDataSource() datasource.DataSource {
	return &certificatesCAExportDataSource{}
}

type certificatesCAExportDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *certificatesCAExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates_ca_export"
}

func (r *certificatesCAExportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificatesCAExportDataSourceModel struct {
	CaId                types.String `tfsdk:"ca_id"`
	ExportedCertificate types.String `tfsdk:"exported_certificate"`
}

func (r *certificatesCAExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to retrieve the details of a trusted certificate authority.",
		Attributes: map[string]schema.Attribute{
			"ca_id": schema.StringAttribute{
				Description: "The ID of the trusted certificate authority to export.",
				Required:    true,
			},
			"exported_certificate": schema.StringAttribute{
				Description: "The exported PEM-encoded certificate.",
				Computed:    true,
			},
		},
	}
}

func (r *certificatesCAExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data certificatesCAExportDataSourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	exportRequest := r.apiClient.CertificatesCaAPI.ExportCaCertificateFile(config.AuthContext(ctx, r.providerConfig), data.CaId.ValueString())
	responseData, httpResp, err := exportRequest.Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while exporting the CA certificate", err, httpResp)
		return
	}

	// Set the exported metadata
	data.ExportedCertificate = types.StringValue(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
