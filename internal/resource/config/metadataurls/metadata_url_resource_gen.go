// Code generated by ping-terraform-plugin-framework-generator

package metadataurls

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &metadataUrlResource{}
	_ resource.ResourceWithConfigure   = &metadataUrlResource{}
	_ resource.ResourceWithImportState = &metadataUrlResource{}
)

func MetadataUrlResource() resource.Resource {
	return &metadataUrlResource{}
}

type metadataUrlResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *metadataUrlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata_url"
}

func (r *metadataUrlResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type metadataUrlResourceModel struct {
	CertView          types.Object `tfsdk:"cert_view"`
	Name              types.String `tfsdk:"name"`
	Url               types.String `tfsdk:"url"`
	UrlId             types.String `tfsdk:"url_id"`
	ValidateSignature types.Bool   `tfsdk:"validate_signature"`
	X509File          types.Object `tfsdk:"x509_file"`
}

func (r *metadataUrlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cert_view": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"crypto_provider": schema.StringAttribute{
						Computed:    true,
						Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
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
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The persistent, unique ID for the certificate.",
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
						Validators: []validator.String{
							stringvalidator.OneOf(
								"VALID",
								"EXPIRED",
								"NOT_YET_VALID",
								"REVOKED",
							),
						},
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
				},
				Computed:    true,
				Description: "The Signature Verification Certificate details. This property is read-only and is always ignored on a POST or PUT.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name for the Metadata URL.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The Metadata URL.",
			},
			"url_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The persistent, unique ID for the Metadata Url. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"validate_signature": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Perform Metadata Signature Validation. The default value is TRUE.",
			},
			"x509_file": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"crypto_provider": schema.StringAttribute{
						Optional:    true,
						Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"LOCAL",
								"HSM",
							),
						},
					},
					"file_data": schema.StringAttribute{
						Required:    true,
						Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
					},
					"formatted_file_data": schema.StringAttribute{
						Computed:    true,
						Description: "The certificate data in PEM format, formatted by PingFederate. This attribute is read-only.",
					},
					"id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							configvalidators.LowercaseId(),
						},
					},
				},
				Optional:    true,
				Description: "Data of the Signature Verification Certificate for the Metadata URL.",
			},
		},
	}
}

func (model *metadataUrlResourceModel) buildClientStruct() (*client.MetadataUrl, diag.Diagnostics) {
	result := &client.MetadataUrl{}
	// name
	result.Name = model.Name.ValueString()
	// url
	result.Url = model.Url.ValueString()
	// url_id
	result.Id = model.UrlId.ValueStringPointer()
	// validate_signature
	result.ValidateSignature = model.ValidateSignature.ValueBoolPointer()
	// x509_file
	if !model.X509File.IsNull() {
		x509FileValue := &client.X509File{}
		x509FileAttrs := model.X509File.Attributes()
		x509FileValue.CryptoProvider = x509FileAttrs["crypto_provider"].(types.String).ValueStringPointer()
		x509FileValue.FileData = x509FileAttrs["file_data"].(types.String).ValueString()
		x509FileValue.Id = x509FileAttrs["id"].(types.String).ValueStringPointer()
		result.X509File = x509FileValue
	}

	return result, nil
}

func (state *metadataUrlResourceModel) readClientResponse(response *client.MetadataUrl) diag.Diagnostics {
	var respDiags diag.Diagnostics
	// cert_view
	certViewAttrTypes := map[string]attr.Type{
		"crypto_provider":           types.StringType,
		"expires":                   types.StringType,
		"id":                        types.StringType,
		"issuer_dn":                 types.StringType,
		"key_algorithm":             types.StringType,
		"key_size":                  types.Int64Type,
		"serial_number":             types.StringType,
		"sha1_fingerprint":          types.StringType,
		"sha256_fingerprint":        types.StringType,
		"signature_algorithm":       types.StringType,
		"status":                    types.StringType,
		"subject_alternative_names": types.ListType{ElemType: types.StringType},
		"subject_dn":                types.StringType,
		"valid_from":                types.StringType,
		"version":                   types.Int64Type,
	}
	var certViewValue types.Object
	if response.CertView == nil {
		certViewValue = types.ObjectNull(certViewAttrTypes)
	} else {
		certViewSubjectAlternativeNamesValue, diags := types.ListValueFrom(context.Background(), types.StringType, response.CertView.SubjectAlternativeNames)
		respDiags.Append(diags...)
		certViewValue, diags = types.ObjectValue(certViewAttrTypes, map[string]attr.Value{
			"crypto_provider":           types.StringPointerValue(response.CertView.CryptoProvider),
			"expires":                   types.StringValue(response.CertView.Expires.Format(time.RFC3339)),
			"id":                        types.StringPointerValue(response.CertView.Id),
			"issuer_dn":                 types.StringPointerValue(response.CertView.IssuerDN),
			"key_algorithm":             types.StringPointerValue(response.CertView.KeyAlgorithm),
			"key_size":                  types.Int64PointerValue(response.CertView.KeySize),
			"serial_number":             types.StringPointerValue(response.CertView.SerialNumber),
			"sha1_fingerprint":          types.StringPointerValue(response.CertView.Sha1Fingerprint),
			"sha256_fingerprint":        types.StringPointerValue(response.CertView.Sha256Fingerprint),
			"signature_algorithm":       types.StringPointerValue(response.CertView.SignatureAlgorithm),
			"status":                    types.StringPointerValue(response.CertView.Status),
			"subject_alternative_names": certViewSubjectAlternativeNamesValue,
			"subject_dn":                types.StringPointerValue(response.CertView.SubjectDN),
			"valid_from":                types.StringValue(response.CertView.ValidFrom.Format(time.RFC3339)),
			"version":                   types.Int64PointerValue(response.CertView.Version),
		})
		respDiags.Append(diags...)
	}

	state.CertView = certViewValue
	// name
	state.Name = types.StringValue(response.Name)
	// url
	state.Url = types.StringValue(response.Url)
	// url_id
	state.UrlId = types.StringPointerValue(response.Id)
	// validate_signature
	state.ValidateSignature = types.BoolPointerValue(response.ValidateSignature)
	// x509_file
	respDiags.Append(state.readClientResponseX509File(response)...)
	return respDiags
}

func (r *metadataUrlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data metadataUrlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.MetadataUrlsAPI.AddMetadataUrl(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.MetadataUrlsAPI.AddMetadataUrlExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the metadataUrl", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metadataUrlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data metadataUrlResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.MetadataUrlsAPI.GetMetadataUrl(config.AuthContext(ctx, r.providerConfig), data.UrlId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the metadataUrl", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the metadataUrl", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metadataUrlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data metadataUrlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.MetadataUrlsAPI.UpdateMetadataUrl(config.AuthContext(ctx, r.providerConfig), data.UrlId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.MetadataUrlsAPI.UpdateMetadataUrlExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the metadataUrl", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metadataUrlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data metadataUrlResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.MetadataUrlsAPI.DeleteMetadataUrl(config.AuthContext(ctx, r.providerConfig), data.UrlId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the metadataUrl", err, httpResp)
	}
}

func (r *metadataUrlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to url_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("url_id"), req, resp)
}
