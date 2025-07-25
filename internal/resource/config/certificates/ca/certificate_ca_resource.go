// Copyright © 2025 Ping Identity Corporation

package certificatesca

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/pemcertificates"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &certificateCAResource{}
	_ resource.ResourceWithConfigure = &certificateCAResource{}

	caResourceCustomId = "ca_id"
)

// CertificateCAResource is a helper function to simplify the provider implementation.
func CertificateCAResource() resource.Resource {
	return &certificateCAResource{}
}

// certificateCAResource is the resource implementation.
type certificateCAResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type certificatesResourceModel struct {
	CaId                    types.String `tfsdk:"ca_id"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
	Expires                 types.String `tfsdk:"expires"`
	FileData                types.String `tfsdk:"file_data"`
	FormattedFileData       types.String `tfsdk:"formatted_file_data"`
	Id                      types.String `tfsdk:"id"`
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

// GetSchema defines the schema for the resource.
func (r *certificateCAResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a trusted Certificate CA.",
		Attributes: map[string]schema.Attribute{
			"ca_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.LowercaseId(),
				},
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. This field is immutable and will trigger a replacement plan if changed.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"LOCAL", "HSM"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expires": schema.StringAttribute{
				Computed:    true,
				Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
			},
			"file_data": schema.StringAttribute{
				Description: "The certificate data in PEM format, base64-encoded. New line characters should be omitted or encoded in this value. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.ValidBase64(),
				},
			},
			"formatted_file_data": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate data in PEM format, formatted by PingFederate. This attribute is read-only.",
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

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalCaCertsFields(ctx context.Context, addRequest *client.X509File, plan certificatesResourceModel) error {
	// Empty strings are treated as equivalent to null
	addRequest.Id = plan.CaId.ValueStringPointer()
	addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
	return nil
}

// Metadata returns the resource type name.
func (r *certificateCAResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_ca"
}

func (r *certificateCAResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *certificateCAResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	// Handle drift detection for the file_data value changing outside of terraform
	var plan, state certificatesResourceModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	if internaltypes.IsDefined(plan.FileData) && internaltypes.IsDefined(state.FormattedFileData) {
		if !pemcertificates.FileDataEquivalent(plan.FileData.ValueString(), state.FormattedFileData.ValueString()) {
			// Since file_data requires replacement on change, if drift is detected we need to replace the resource
			plan.FormattedFileData = types.StringUnknown()
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
			resp.RequiresReplace = path.Paths{path.Root("formatted_file_data")}
		}
	}
}

func (state *certificatesResourceModel) readCertificateResponse(r *client.CertView, formattedPEM *string, isImportRead bool) {
	state.CaId = types.StringPointerValue(r.Id)
	state.Id = types.StringPointerValue(r.Id)
	state.CryptoProvider = types.StringPointerValue(r.CryptoProvider)
	if formattedPEM != nil {
		if isImportRead {
			// The PEM file must be in base-64 for the validator to accept it
			state.FileData = types.StringValue(base64.StdEncoding.EncodeToString([]byte(*formattedPEM)))
		}
		state.FormattedFileData = types.StringPointerValue(formattedPEM)
	} else {
		// In theory this should never happen, but if the server fails to return the PEM file,
		// just copy the value from the plan.
		state.FormattedFileData = types.StringValue(state.FileData.ValueString())
	}
	state.SerialNumber = types.StringPointerValue(r.SerialNumber)
	state.SubjectDn = types.StringPointerValue(r.SubjectDN)
	state.SubjectAlternativeNames = internaltypes.GetStringSet(r.SubjectAlternativeNames)
	state.IssuerDn = types.StringPointerValue(r.IssuerDN)
	if r.ValidFrom != nil {
		state.ValidFrom = types.StringValue(r.ValidFrom.Format(time.RFC3339))
	} else {
		state.ValidFrom = types.StringNull()
	}
	if r.Expires != nil {
		state.Expires = types.StringValue(r.Expires.Format(time.RFC3339))
	} else {
		state.Expires = types.StringNull()
	}
	state.KeyAlgorithm = types.StringPointerValue(r.KeyAlgorithm)
	state.KeySize = types.Int64PointerValue(r.KeySize)
	state.SignatureAlgorithm = types.StringPointerValue(r.SignatureAlgorithm)
	state.Version = types.Int64PointerValue(r.Version)
	state.Sha1Fingerprint = types.StringPointerValue(r.Sha1Fingerprint)
	state.Sha256Fingerprint = types.StringPointerValue(r.Sha256Fingerprint)
	state.Status = types.StringPointerValue(r.Status)
	state.CryptoProvider = types.StringPointerValue(r.CryptoProvider)
}

func (r *certificateCAResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan certificatesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createCertificate := client.NewX509File((plan.FileData.ValueString()))
	err := addOptionalCaCertsFields(ctx, createCertificate, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for a CA Certificate: "+err.Error())
		return
	}

	apiCreateCertificate := r.apiClient.CertificatesCaAPI.ImportTrustedCA(config.AuthContext(ctx, r.providerConfig))
	apiCreateCertificate = apiCreateCertificate.Body(*createCertificate)
	certificateResponse, httpResp, err := r.apiClient.CertificatesCaAPI.ImportTrustedCAExecute(apiCreateCertificate)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating a CA Certificate", err, httpResp, &caResourceCustomId)
		return
	}

	// Get the server's formatted PEM file for the cert
	var pemFile *string
	if certificateResponse.Id != nil {
		exportedPemFile, httpResp, err := r.apiClient.CertificatesCaAPI.ExportCaCertificateFile(config.AuthContext(ctx, r.providerConfig), *certificateResponse.Id).Execute()
		if err != nil {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the PEM file of a CA certificate after create", err, httpResp, &caResourceCustomId)
			return
		}
		pemFile = &exportedPemFile
	}

	// Read the response into the state
	plan.readCertificateResponse(certificateResponse, pemFile, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *certificateCAResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)
	var state certificatesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	certificateResponse, httpResp, err := r.apiClient.CertificatesCaAPI.GetTrustedCert(config.AuthContext(ctx, r.providerConfig), state.CaId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Certificate CA", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading a CA certificate", err, httpResp, &caResourceCustomId)
		}
		return
	}

	// Get the server's formatted PEM file for the cert
	exportedPemFile, httpResp, err := r.apiClient.CertificatesCaAPI.ExportCaCertificateFile(config.AuthContext(ctx, r.providerConfig), state.CaId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Certificate CA", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the PEM file of a CA certificate", err, httpResp, &caResourceCustomId)
		}
		return
	}

	// Read the response into the state
	state.readCertificateResponse(certificateResponse, &exportedPemFile, isImportRead)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *certificateCAResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes in this resource use the RequiresReplace plan modifier, so no updates can be done.
	// The PF API does not support updating a certificate CA, only creating and deleting.
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *certificateCAResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state certificatesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CertificatesCaAPI.DeleteTrustedCA(config.AuthContext(ctx, r.providerConfig), state.CaId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting a CA Certificate", err, httpResp, &caResourceCustomId)
	}
}

func (r *certificateCAResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("ca_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
