package keypairssigningcsr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource              = &keypairsSigningCsrExportResource{}
	_ resource.ResourceWithConfigure = &keypairsSigningCsrExportResource{}
)

func KeypairsSigningCsrExportResource() resource.Resource {
	return &keypairsSigningCsrExportResource{}
}

type keypairsSigningCsrExportResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSigningCsrExportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_signing_csr_export"
}

func (r *keypairsSigningCsrExportResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSigningCsrExportResourceModel struct {
	KeypairId   types.String `tfsdk:"keypair_id"`
	ExportedCsr types.String `tfsdk:"exported_csr"`
}

func (r *keypairsSigningCsrExportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *keypairsSigningCsrExportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSigningCsrExportResourceModel

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

func (r *keypairsSigningCsrExportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// PingFederate provides no read endpoint for this resource, so we'll just maintain whatever is in state
	resp.State.Raw = req.State.Raw
}

func (r *keypairsSigningCsrExportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This method won't be called since all non-computed attributes require replacement
}

func (r *keypairsSigningCsrExportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// There is no way to delete an exported CSR
	resp.Diagnostics.AddWarning("Configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}
