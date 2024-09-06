// Code generated by ping-terraform-plugin-framework-generator

package serversettingswstruststssettingsissuercertificates

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &serverSettingsWsTrustStsSettingsIssuerCertificateResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsWsTrustStsSettingsIssuerCertificateResource{}
	_ resource.ResourceWithImportState = &serverSettingsWsTrustStsSettingsIssuerCertificateResource{}
)

func ServerSettingsWsTrustStsSettingsIssuerCertificateResource() resource.Resource {
	return &serverSettingsWsTrustStsSettingsIssuerCertificateResource{}
}

type serverSettingsWsTrustStsSettingsIssuerCertificateResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_ws_trust_sts_settings_issuer_certificate"
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel struct {
	Active                  types.Bool   `tfsdk:"active"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
	Expires                 types.String `tfsdk:"expires"`
	FileData                types.String `tfsdk:"file_data"`
	CertificateId           types.String `tfsdk:"certificate_id"`
	Id                      types.String `tfsdk:"id"`
	IssuerDn                types.String `tfsdk:"issuer_dn"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
	KeySize                 types.Int64  `tfsdk:"key_size"`
	SerialNumber            types.String `tfsdk:"serial_number"`
	Sha1Fingerprint         types.String `tfsdk:"sha1_fingerprint"`
	Sha256Fingerprint       types.String `tfsdk:"sha256_fingerprint"`
	SignatureAlgorithm      types.String `tfsdk:"signature_algorithm"`
	Status                  types.String `tfsdk:"status"`
	SubjectAlternativeNames types.List   `tfsdk:"subject_alternative_names"`
	SubjectDn               types.String `tfsdk:"subject_dn"`
	ValidFrom               types.String `tfsdk:"valid_from"`
	Version                 types.Int64  `tfsdk:"version"`
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage certificates for WS-Trust STS settings.",
		Attributes: map[string]schema.Attribute{
			"crypto_provider": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"certificate_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					configvalidators.LowercaseId(),
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
			"subject_alternative_names": schema.ListAttribute{
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
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether this an active certificate or not.",
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel) buildClientStruct() (*client.X509File, diag.Diagnostics) {
	result := &client.X509File{}
	// crypto_provider
	result.CryptoProvider = model.CryptoProvider.ValueStringPointer()
	// file_data
	result.FileData = model.FileData.ValueString()
	// certificate_id
	result.Id = model.CertificateId.ValueStringPointer()
	return result, nil
}

func (state *serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel) readClientResponse(response *client.IssuerCert) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringPointerValue(response.CertView.Id)
	// active
	state.Active = types.BoolPointerValue(response.Active)
	// crypto_provider
	state.CryptoProvider = types.StringPointerValue(response.CertView.CryptoProvider)
	// expires
	state.Expires = types.StringValue(response.CertView.Expires.Format(time.RFC3339))
	// certificate_id
	state.CertificateId = types.StringPointerValue(response.CertView.Id)
	// issuer_dn
	state.IssuerDn = types.StringPointerValue(response.CertView.IssuerDN)
	// key_algorithm
	state.KeyAlgorithm = types.StringPointerValue(response.CertView.KeyAlgorithm)
	// key_size
	state.KeySize = types.Int64PointerValue(response.CertView.KeySize)
	// serial_number
	state.SerialNumber = types.StringPointerValue(response.CertView.SerialNumber)
	// sha1_fingerprint
	state.Sha1Fingerprint = types.StringPointerValue(response.CertView.Sha1Fingerprint)
	// sha256_fingerprint
	state.Sha256Fingerprint = types.StringPointerValue(response.CertView.Sha256Fingerprint)
	// signature_algorithm
	state.SignatureAlgorithm = types.StringPointerValue(response.CertView.SignatureAlgorithm)
	// status
	state.Status = types.StringPointerValue(response.CertView.Status)
	// subject_alternative_names
	state.SubjectAlternativeNames, diags = types.ListValueFrom(context.Background(), types.StringType, response.CertView.SubjectAlternativeNames)
	respDiags.Append(diags...)
	// subject_dn
	state.SubjectDn = types.StringPointerValue(response.CertView.SubjectDN)
	// valid_from
	state.ValidFrom = types.StringValue(response.CertView.ValidFrom.Format(time.RFC3339))
	// version
	state.Version = types.Int64PointerValue(response.CertView.Version)
	return respDiags
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.ServerSettingsAPI.ImportCertificate(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.ImportCertificateExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the WS Trust issuer certificate", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.GetCert(config.AuthContext(ctx, r.providerConfig), data.CertificateId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "WS Trust Issuer Certificate", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the WS Trust issuer certificate", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This method won't be called since all non-computed attributes require replacement
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data serverSettingsWsTrustStsSettingsIssuerCertificateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.ServerSettingsAPI.DeleteCertificate(config.AuthContext(ctx, r.providerConfig), data.CertificateId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the WS Trust issuer certificate", err, httpResp)
	}
}

func (r *serverSettingsWsTrustStsSettingsIssuerCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to certificate_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("certificate_id"), req, resp)
}
