// Copyright © 2026 Ping Identity Corporation

package keypairssslclientcertificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ datasource.DataSource              = &keypairsSslClientCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &keypairsSslClientCertificateDataSource{}

	customId = "key_id"
)

func KeypairsSslClientCertificateDataSource() datasource.DataSource {
	return &keypairsSslClientCertificateDataSource{}
}

type keypairsSslClientCertificateDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSslClientCertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_client_certificate"
}

func (r *keypairsSslClientCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSslClientCertificateDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	KeyId               types.String `tfsdk:"key_id"`
	ExportedCertificate types.String `tfsdk:"exported_certificate"`
}

// GetSchema defines the schema for the datasource.
func (r *keypairsSslClientCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve the PEM-encoded certificate from a given SSL client key pair.",
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Description: "The ID of the key pair to export.",
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

func (r *keypairsSslClientCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keypairsSslClientCertificateDataSourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	exportRequest := r.apiClient.KeyPairsSslClientAPI.ExportSslClientCertificateFile(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString())
	responseData, httpResp, err := exportRequest.Execute()
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while exporting the certificate", err, httpResp, &customId)
		return
	}

	// Set the exported metadata
	data.Id = types.StringValue(data.KeyId.ValueString())
	data.ExportedCertificate = types.StringValue(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
