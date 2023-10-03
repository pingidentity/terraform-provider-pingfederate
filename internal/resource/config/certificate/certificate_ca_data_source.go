package certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

// create a Certificate data source
func NewCertificateDataSource() datasource.DataSource {
	return &certificatesDataSource{}
}

// certificatesResource is the resource implementation.
type certificatesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *certificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_ca"
}

// Configure adds the provider configured client to the data source.
func (r *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificatesDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	CustomId       types.String `tfsdk:"custom_id"`
	FileData       types.String `tfsdk:"file_data"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the datasource.
func (r *certificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a CertificateCA Import.",
		Attributes: map[string]schema.Attribute{
			"custom_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate",
				Optional:    true,
				Computed:    true,
				// PlanModifiers: []planmodifier.String{
				// stringplanmodifier.UseStateForUnknown(),
				// stringplanmodifier.RequiresReplace(),
				// },
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
				Optional:    true,
				// Validators: []validator.String{
				// stringvalidator.OneOf([]string{"LOCAL", "HSM"}...),
				// },
				// PlanModifiers: []planmodifier.String{
				// stringplanmodifier.UseStateForUnknown(),
				// stringplanmodifier.RequiresReplace(),
				// },
			},
			"file_data": schema.StringAttribute{
				Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
				Required:    true,
				// PlanModifiers: []planmodifier.String{
				// stringplanmodifier.UseStateForUnknown(),
				// stringplanmodifier.RequiresReplace(),
				// },
			},
		},
	}

	config.AddCommonDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// func addOptionalCaCertsFields(ctx context.Context, addRequest *client.X509File, plan certificatesResourceModel) error {
// 	// Empty strings are treated as equivalent to null
// 	if internaltypes.IsDefined(plan.CustomId) {
// 		addRequest.Id = plan.CustomId.ValueStringPointer()
// 	}
// 	if internaltypes.IsDefined(plan.CryptoProvider) {
// 		addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
// 	}
// 	return nil
// }

func readCertificateResponseDataSource(ctx context.Context, r *client.CertView, state *certificatesDataSourceModel, diagnostics *diag.Diagnostics) {
	state.CustomId = internaltypes.StringTypeOrNil(r.Id, false)
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.CryptoProvider = internaltypes.StringTypeOrNil(r.CryptoProvider, false)
}

func (r *certificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificatesDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadCertificate, httpResp, err := r.apiClient.CertificatesCaApi.GetTrustedCert(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while looking for a Certificate", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a Certificate", err, httpResp)
		}
		return
	}

	// Log response JSON
	_, responseErr := apiReadCertificate.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of a Certificate: %s", responseErr.Error())
	}
	// Read the response into the state
	readCertificateResponseDataSource(ctx, apiReadCertificate, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
