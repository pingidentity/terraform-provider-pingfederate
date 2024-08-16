package keypairssslservercsr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource              = &keypairsSslServerCsrExportResource{}
	_ resource.ResourceWithConfigure = &keypairsSslServerCsrExportResource{}
)

func KeypairsSslServerCsrExportResource() resource.Resource {
	return &keypairsSslServerCsrExportResource{}
}

type keypairsSslServerCsrExportResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSslServerCsrExportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_server_csr_export"
}

func (r *keypairsSslServerCsrExportResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSslServerCsrExportResourceModel struct {
	KeypairId   types.String `tfsdk:"keypair_id"`
	ExportedCsr types.String `tfsdk:"exported_csr"`
}

func (r *keypairsSslServerCsrExportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to export CSRs for SSL server key pairs.",
		Attributes: map[string]schema.Attribute{
			"keypair_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the key pair.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"exported_csr": schema.StringAttribute{
				Description: "The exported PEM-encoded certificate signing request.",
				Computed:    true,
			},
		},
	}
}

func (r *keypairsSslServerCsrExportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSslServerCsrExportResourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	exportRequest := r.apiClient.KeyPairsSslServerAPI.ExportSslServerCsr(config.AuthContext(ctx, r.providerConfig), data.KeypairId.ValueString())
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

func (r *keypairsSslServerCsrExportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// PingFederate provides no read endpoint for this resource, so we'll just maintain whatever is in state
	resp.State.Raw = req.State.Raw
}

func (r *keypairsSslServerCsrExportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This method won't be called since all non-computed attributes require replacement
}

func (r *keypairsSslServerCsrExportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// There is no way to delete an exported CSR
	resp.Diagnostics.AddWarning("Configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}
