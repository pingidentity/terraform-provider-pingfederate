package keypairssslclientcsr

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource              = &keypairsSslClientCsrResource{}
	_ resource.ResourceWithConfigure = &keypairsSslClientCsrResource{}
)

func KeypairsSslClientCsrResource() resource.Resource {
	return &keypairsSslClientCsrResource{}
}

type keypairsSslClientCsrResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSslClientCsrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_client_csr"
}

func (r *keypairsSslClientCsrResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSslClientCsrResourceModel struct {
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
	Expires                 types.String `tfsdk:"expires"`
	FileData                types.String `tfsdk:"file_data"`
	Id                      types.String `tfsdk:"id"`
	KeypairId               types.String `tfsdk:"keypair_id"`
	IssuerDn                types.String `tfsdk:"issuer_dn"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
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

func (r *keypairsSslClientCsrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage CSR responses for SSL client key pairs.",
		Attributes: map[string]schema.Attribute{
			"crypto_provider": schema.StringAttribute{
				Computed:    true,
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is `true`. Options are `LOCAL` or `HSM`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"LOCAL",
						"HSM",
					),
				},
			},
			"expires": schema.StringAttribute{
				Computed:    true,
				Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
			},
			"file_data": schema.StringAttribute{
				Required:    true,
				Description: "The CSR response file data in PKCS7 format or as an X.509 certificate. PEM encoding (with or without the header and footer lines) is required. New line characters should be omitted or encoded in this value.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"issuer_dn": schema.StringAttribute{
				Computed:    true,
				Description: "The issuer's distinguished name.",
			},
			"key_algorithm": schema.StringAttribute{
				Computed:    true,
				Description: "The public key algorithm.",
			},
			"key_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The public key size.",
			},
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
			"serial_number": schema.StringAttribute{
				Computed:    true,
				Description: "The serial number assigned by the CA.",
			},
			"sha1_fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-1 fingerprint in Hex encoding.",
			},
			"sha256_fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 fingerprint in Hex encoding.",
			},
			"signature_algorithm": schema.StringAttribute{
				Computed:    true,
				Description: "The signature algorithm.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the item.",
			},
			"subject_alternative_names": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The subject alternative names (SAN).",
			},
			"subject_dn": schema.StringAttribute{
				Computed:    true,
				Description: "The subject's distinguished name.",
			},
			"valid_from": schema.StringAttribute{
				Computed:    true,
				Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "The X.509 version to which the item conforms.",
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *keypairsSslClientCsrResourceModel) buildClientStruct() (*client.CSRResponse, diag.Diagnostics) {
	result := &client.CSRResponse{}
	// file_data
	result.FileData = model.FileData.ValueString()
	return result, nil
}

func (state *keypairsSslClientCsrResourceModel) readClientResponse(response *client.KeyPairView) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringPointerValue(response.Id)
	// crypto_provider
	state.CryptoProvider = types.StringPointerValue(response.CryptoProvider)
	// expires
	state.Expires = types.StringValue(response.Expires.Format(time.RFC3339))
	// keypair_id
	state.KeypairId = types.StringPointerValue(response.Id)
	// issuer_dn
	state.IssuerDn = types.StringPointerValue(response.IssuerDN)
	// key_algorithm
	state.KeyAlgorithm = types.StringPointerValue(response.KeyAlgorithm)
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
	state.ValidFrom = types.StringValue(response.ValidFrom.Format(time.RFC3339))
	// version
	state.Version = types.Int64PointerValue(response.Version)
	return respDiags
}

func (r *keypairsSslClientCsrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSslClientCsrResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.KeyPairsSslClientAPI.ImportSslClientCsrResponse(config.AuthContext(ctx, r.providerConfig), data.KeypairId.ValueString())
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KeyPairsSslClientAPI.ImportSslClientCsrResponseExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while importing the certificate signing request response", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslClientCsrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// PingFederate provides no read endpoint for this resource, so we'll just maintain whatever is in state
	resp.State.Raw = req.State.Raw
}

func (r *keypairsSslClientCsrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This method won't be called since all non-computed attributes require replacement
}

func (r *keypairsSslClientCsrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// There is no way to delete the imported CSR response
	resp.Diagnostics.AddWarning("Configuration cannot be returned to original state.  The resource has been removed from Terraform state but the configuration remains applied to the environment.", "")
}
