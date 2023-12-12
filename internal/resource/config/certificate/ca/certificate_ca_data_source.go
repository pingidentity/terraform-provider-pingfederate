package certificateca

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

// Create a Administrative Account data source
func CertificateDataSource() datasource.DataSource {
	return &certificatesDataSource{}
}

// certificatesDataSource is the datasource implementation.
type certificatesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type certificateCaModel struct {
	Id                      types.String `tfsdk:"id"`
	CaId                    types.String `tfsdk:"ca_id"`
	SerialNumber            types.String `tfsdk:"serial_number"`
	SubjectDN               types.String `tfsdk:"subject_dn"`
	SubjectAlternativeNames types.Set    `tfsdk:"subject_alternative_names"`
	IssuerDN                types.String `tfsdk:"issuer_dn"`
	ValidFrom               types.String `tfsdk:"valid_from"`
	Expires                 types.String `tfsdk:"expires"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
	KeySize                 types.Int64  `tfsdk:"key_size"`
	SignatureAlgorithm      types.String `tfsdk:"signature_algorithm"`
	Version                 types.Int64  `tfsdk:"version"`
	Sha1Fingerprint         types.String `tfsdk:"sha1_fingerprint"`
	Sha256Fingerprint       types.String `tfsdk:"sha256_fingerprint"`
	Status                  types.String `tfsdk:"status"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the datasource.
func (r *certificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a trusted Certificate CA.",
		Attributes: map[string]schema.Attribute{
			"serial_number": schema.StringAttribute{
				Description: "The serial number assigned by the CA",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_dn": schema.StringAttribute{
				Description: "The subject's distinguished name",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_alternative_names": schema.SetAttribute{
				Description: "The subject alternative names (SAN)",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"issuer_dn": schema.StringAttribute{
				Description: "The issuer's distinguished name",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"valid_from": schema.StringAttribute{
				Description: "The start date from which the item is valid, in ISO 8601 format (UTC)",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The end date up until which the item is valid, in ISO 8601 format (UTC)",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The public key algorithm",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_size": schema.Int64Attribute{
				Description: "The public key size",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"signature_algorithm": schema.StringAttribute{
				Description: "The signature algorithm",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "The X.509 version to which the item conforms",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sha1_fingerprint": schema.StringAttribute{
				Description: "SHA-1 fingerprint in Hex encoding",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sha256_fingerprint": schema.StringAttribute{
				Description: "SHA-256 fingerprint in Hex encoding",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the item.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef)
	id.ToDataSourceSchemaCustomId(&schemaDef, "ca_id", true, "Unique ID for the certificate.")
	resp.Schema = schemaDef
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

// Read a CertificateResponse object into the model struct
func readCertificateResponseDataSource(ctx context.Context, r *client.CertView, state *certificateCaModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringPointerValue(r.Id)
	state.CaId = types.StringPointerValue(r.Id)
	state.SerialNumber = types.StringPointerValue(r.SerialNumber)
	state.SubjectDN = types.StringPointerValue(r.SubjectDN)
	state.SubjectAlternativeNames = internaltypes.GetStringSet(r.SubjectAlternativeNames)
	state.IssuerDN = types.StringPointerValue(r.IssuerDN)
	state.ValidFrom = types.StringValue(r.ValidFrom.Format(time.RFC3339))
	state.Expires = types.StringValue(r.Expires.Format(time.RFC3339))
	state.KeyAlgorithm = types.StringPointerValue(r.KeyAlgorithm)
	state.KeySize = types.Int64PointerValue(r.KeySize)
	state.SignatureAlgorithm = types.StringPointerValue(r.SignatureAlgorithm)
	state.Version = types.Int64PointerValue(r.Version)
	state.Sha1Fingerprint = types.StringPointerValue(r.Sha1Fingerprint)
	state.Sha256Fingerprint = types.StringPointerValue(r.Sha256Fingerprint)
	state.Status = types.StringPointerValue(r.Status)
	state.CryptoProvider = types.StringPointerValue(r.CryptoProvider)
}

// Read resource information
func (r *certificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificateCaModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadCertificate, httpResp, err := r.apiClient.CertificatesCaAPI.GetTrustedCert(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CaId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a trusted certificate CA", err, httpResp)
		return
	}

	// Read the response into the state
	readCertificateResponseDataSource(ctx, apiReadCertificate, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
