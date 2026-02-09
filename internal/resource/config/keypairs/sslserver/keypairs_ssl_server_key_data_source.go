// Copyright © 2026 Ping Identity Corporation

package keypairssslserver

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &keypairsSslServerKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &keypairsSslServerKeyDataSource{}
)

// KeypairsSslServerKeyDataSource is a helper function to simplify the provider implementation.
func KeypairsSslServerKeyDataSource() datasource.DataSource {
	return &keypairsSslServerKeyDataSource{}
}

// keypairsSslServerKeyDataSource is the data source implementation.
type keypairsSslServerKeyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keypairsSslServerKeyDataSourceModel struct {
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
	Expires                 types.String `tfsdk:"expires"`
	Id                      types.String `tfsdk:"id"`
	IssuerDn                types.String `tfsdk:"issuer_dn"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
	KeyId                   types.String `tfsdk:"key_id"`
	KeySize                 types.Int64  `tfsdk:"key_size"`
	SerialNumber            types.String `tfsdk:"serial_number"`
	Sha1Fingerprint         types.String `tfsdk:"sha1_fingerprint"`
	Sha256Fingerprint       types.String `tfsdk:"sha256_fingerprint"`
	SignatureAlgorithm      types.String `tfsdk:"signature_algorithm"`
	Status                  types.String `tfsdk:"status"`
	SubjectAlternativeNames types.Set    `tfsdk:"subject_alternative_names"`
	SubjectDn               types.String `tfsdk:"subject_dn"`
	ValidFrom               types.String `tfsdk:"valid_from"`
	Version                 types.Int64  `tfsdk:"version"`
}

// GetSchema defines the schema for the data source.
func (r *keypairsSslServerKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve a ssl server key pair.",
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate.",
				Required:    true,
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Supported values are `LOCAL` and `HSM`.",
				Computed:    true,
			},
			"serial_number": schema.StringAttribute{
				Description: "The serial number assigned by the CA",
				Computed:    true,
			},
			"subject_dn": schema.StringAttribute{
				Description: "The subject's distinguished name",
				Computed:    true,
			},
			"subject_alternative_names": schema.SetAttribute{
				Description: "The subject alternative names (SAN).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"issuer_dn": schema.StringAttribute{
				Description: "The issuer's distinguished name",
				Computed:    true,
			},
			"valid_from": schema.StringAttribute{
				Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The end date up until which the item is valid, in ISO 8601 format (UTC)",
				Computed:    true,
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The public key algorithm.",
				Computed:    true,
			},
			"key_size": schema.Int64Attribute{
				Description: "The public key size, in bits.",
				Computed:    true,
			},
			"signature_algorithm": schema.StringAttribute{
				Description: "The signature algorithm.",
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "The X.509 version to which the item conforms",
				Computed:    true,
			},
			"sha1_fingerprint": schema.StringAttribute{
				Description: "SHA-1 fingerprint in Hex encoding",
				Computed:    true,
			},
			"sha256_fingerprint": schema.StringAttribute{
				Description: "SHA-256 fingerprint in Hex encoding",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the item.",
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&resp.Schema)
}

func (state *keypairsSslServerKeyDataSourceModel) readClientResponse(response *client.KeyPairView) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringPointerValue(response.Id)
	// crypto_provider
	state.CryptoProvider = types.StringPointerValue(response.CryptoProvider)
	// expires
	if response.Expires != nil {
		state.Expires = types.StringValue(response.Expires.Format(time.RFC3339))
	} else {
		state.Expires = types.StringNull()
	}
	// issuer_dn
	state.IssuerDn = types.StringPointerValue(response.IssuerDN)
	// key_algorithm
	state.KeyAlgorithm = types.StringPointerValue(response.KeyAlgorithm)
	// key_id
	state.KeyId = types.StringPointerValue(response.Id)
	// key_size
	state.KeySize = types.Int64PointerValue(response.KeySize)
	// serial_number
	state.SerialNumber = types.StringPointerValue(response.SerialNumber)
	// sha1_fingerprint
	state.Sha1Fingerprint = types.StringPointerValue(response.Sha1Fingerprint)
	// sha256_fingerprint
	state.Sha256Fingerprint = types.StringPointerValue(response.Sha256Fingerprint)
	// signature_algorithm
	state.SignatureAlgorithm = types.StringPointerValue(response.SignatureAlgorithm)
	// status
	state.Status = types.StringPointerValue(response.Status)
	// subject_alternative_names
	state.SubjectAlternativeNames, diags = types.SetValueFrom(context.Background(), types.StringType, response.SubjectAlternativeNames)
	respDiags.Append(diags...)
	// subject_dn
	state.SubjectDn = types.StringPointerValue(response.SubjectDN)
	// valid_from
	if response.ValidFrom != nil {
		state.ValidFrom = types.StringValue(response.ValidFrom.Format(time.RFC3339))
	} else {
		state.ValidFrom = types.StringNull()
	}
	// version
	state.Version = types.Int64PointerValue(response.Version)
	return respDiags
}

// Metadata returns the data source type name.
func (r *keypairsSslServerKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_server_key"
}

func (r *keypairsSslServerKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *keypairsSslServerKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keypairsSslServerKeyDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KeyPairsSslServerAPI.GetSslServerKeyPair(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the key pair", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
