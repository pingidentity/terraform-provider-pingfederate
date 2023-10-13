package keypairs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &keyPairsSigningImportDataSource{}
	_ datasource.DataSourceWithConfigure = &keyPairsSigningImportDataSource{}
)

// Create a Administrative Account data source
func NewKeyPairsSigningImportDataSource() datasource.DataSource {
	return &keyPairsSigningImportDataSource{}
}

// keyPairsSigningImportDataSource is the datasource implementation.
type keyPairsSigningImportDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keyPairsSigningImportDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
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
	RotationSettings        types.Object `tfsdk:"rotation_settings"`
}

// GetSchema defines the schema for the datasource.
func (r *keyPairsSigningImportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a KeyPairsSigningImport.",
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
			"rotation_settings": schema.SingleNestedAttribute{
				Description: "The local identity profile data store configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"creation_buffer_days": schema.Int64Attribute{
						Description: "Buffer days before key pair expiration for creation of a new key pair.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"activation_buffer_days": schema.Int64Attribute{
						Description: "Buffer days before key pair expiration for activation of the new key pair.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"valid_days": schema.Int64Attribute{
						Description: "Valid days for the new key pair to be created. If this property is unset, the validity days of the original key pair will be used.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"key_algorithm": schema.StringAttribute{
						Description: "Key algorithm to be used while creating a new key pair. If this property is unset, the key algorithm of the original key pair will be used. Supported algorithms are available through the /keyPairs/keyAlgorithms endpoint.						",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"key_size": schema.Int64Attribute{
						Description: "Key size, in bits. If this property is unset, the key size of the original key pair will be used. Supported key sizes are available through the /keyPairs/keyAlgorithms endpoint.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"signature_algorithm": schema.StringAttribute{
						Description: "Required if the original key pair used SHA1 algorithm. If this property is unset, the default signature algorithm of the original key pair will be used. Supported signature algorithms are available through the /keyPairs/keyAlgorithms endpoint.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
		},
	}
	id.AddToDataSourceSchema(&schemaDef, true, "The persistent, unique ID for the certificate.")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *keyPairsSigningImportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair_signing_import"
}

// Configure adds the provider configured client to the data source.
func (r *keyPairsSigningImportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readKeyPairsSigningImportResponseDataSource(ctx context.Context, r *client.KeyPairView, state *keyPairsSigningImportDataSourceModel, expectedValues *keyPairsSigningImportDataSourceModel) diag.Diagnostics {
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.SerialNumber = internaltypes.StringTypeOrNil(r.SerialNumber, false)
	state.SubjectDN = internaltypes.StringTypeOrNil(r.SubjectDN, false)
	state.SubjectAlternativeNames = internaltypes.GetStringSet(r.SubjectAlternativeNames)
	state.IssuerDN = internaltypes.StringTypeOrNil(r.IssuerDN, false)
	state.ValidFrom = types.StringValue(r.ValidFrom.Format(time.RFC3339))
	state.Expires = types.StringValue(r.Expires.Format(time.RFC3339))
	state.KeyAlgorithm = internaltypes.StringTypeOrNil(r.KeyAlgorithm, false)
	state.KeySize = internaltypes.Int64TypeOrNil(r.KeySize)
	state.SignatureAlgorithm = internaltypes.StringTypeOrNil(r.SignatureAlgorithm, false)
	state.Version = internaltypes.Int64TypeOrNil(r.Version)
	state.Sha1Fingerprint = internaltypes.StringTypeOrNil(r.Sha1Fingerprint, false)
	state.Sha256Fingerprint = internaltypes.StringTypeOrNil(r.Sha256Fingerprint, false)
	state.Status = internaltypes.StringTypeOrNil(r.Status, false)
	state.CryptoProvider = internaltypes.StringTypeOrNil(r.CryptoProvider, false)

	rotationSettings := r.RotationSettings
	rotationSettingsAttrTypes := map[string]attr.Type{
		"id":                     basetypes.StringType{},
		"creation_buffer_days":   basetypes.Int64Type{},
		"activation_buffer_days": basetypes.Int64Type{},
		"valid_days":             basetypes.Int64Type{},
		"key_algorithm":          basetypes.StringType{},
		"key_size":               basetypes.Int64Type{},
		"signature_algorithm":    basetypes.StringType{},
	}
	var valueFromDiags diag.Diagnostics
	state.RotationSettings, valueFromDiags = types.ObjectValueFrom(ctx, rotationSettingsAttrTypes, rotationSettings)
	return valueFromDiags
}

// Read resource information
func (r *keyPairsSigningImportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state keyPairsSigningImportDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadKeyPairsSigningImport, httpResp, err := r.apiClient.KeyPairsSigningAPI.GetSigningKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair Signing Import", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadKeyPairsSigningImport.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of the KeyPair Signing Import: %s", responseErr.Error())
	}

	// Read the response into the state
	diags = readKeyPairsSigningImportResponseDataSource(ctx, apiReadKeyPairsSigningImport, &state, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
