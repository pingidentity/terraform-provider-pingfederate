// Copyright © 2025 Ping Identity Corporation

package certificatesca

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ datasource.DataSource              = &certificatesCAExportDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesCAExportDataSource{}

	caExportCustomId = "ca_id"
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
	Id                  types.String `tfsdk:"id"`
}

func (r *certificatesCAExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to retrieve the details of a trusted certificate authority.",
		Attributes: map[string]schema.Attribute{
			"ca_id": schema.StringAttribute{
				Description: "The ID of the trusted certificate authority to export.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"exported_certificate": schema.StringAttribute{
				Description: "The exported PEM-encoded certificate.",
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&resp.Schema)
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
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while exporting the CA certificate", err, httpResp, &caExportCustomId)
		return
	}

	// Set the exported metadata
	data.Id = types.StringValue(data.CaId.ValueString())
	data.ExportedCertificate = types.StringValue(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
