package keypairsslserverimport

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &keyPairsSslServerImportResource{}
	_ resource.ResourceWithConfigure   = &keyPairsSslServerImportResource{}
	_ resource.ResourceWithImportState = &keyPairsSslServerImportResource{}

	rotationSettingsAttrTypes = map[string]attr.Type{
		"id":                     types.StringType,
		"creation_buffer_days":   types.Int64Type,
		"activation_buffer_days": types.Int64Type,
		"valid_days":             types.Int64Type,
		"key_algorithm":          types.StringType,
		"key_size":               types.Int64Type,
		"signature_algorithm":    types.StringType,
	}
)

// KeyPairsSslServerImportResource is a helper function to simplify the provider implementation.
func KeyPairsSslServerImportResource() resource.Resource {
	return &keyPairsSslServerImportResource{}
}

// keyPairsSslServerImportResource is the resource implementation.
type keyPairsSslServerImportResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keyPairsSslServerImportResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ImportId       types.String `tfsdk:"import_id"`
	FileData       types.String `tfsdk:"file_data"`
	Format         types.String `tfsdk:"format"`
	Password       types.String `tfsdk:"password"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
	// Computed attributes
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
	RotationSettings        types.Object `tfsdk:"rotation_settings"`
}

// GetSchema defines the schema for the resource.
func (r *keyPairsSslServerImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a file for importing an SSL key pair.",
		Attributes: map[string]schema.Attribute{
			"file_data": schema.StringAttribute{
				Description: "Base-64 encoded PKCS12 or PEM file data. In the case of PEM, the raw (non-base-64) data is also accepted. In BCFIPS mode, only PEM with PBES2 and AES or Triple DES encryption is accepted and 128-bit salt is required.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"format": schema.StringAttribute{
				Description: "Key pair file format. If specified, this field will control what file format is expected, otherwise the format will be auto-detected. In BCFIPS mode, only PEM is supported. (PKCS12, PEM)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for the file. In BCFIPS mode, the password must be at least 14 characters.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. (LOCAL, HSM)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
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

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"import_id",
		true,
		true,
		"The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.")
	resp.Schema = schema
}

func addOptionalKeyPairsSslServerImportFields(ctx context.Context, addRequest *client.KeyPairFile, plan keyPairsSslServerImportResourceModel) error {

	if internaltypes.IsDefined(plan.ImportId) {
		addRequest.Id = plan.ImportId.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CryptoProvider) {
		addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *keyPairsSslServerImportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair_ssl_server_import"
}

func (r *keyPairsSslServerImportResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readKeyPairsSslServerImportResponse(ctx context.Context, r *client.KeyPairView, state *keyPairsSslServerImportResourceModel, planFileData string, planFormat string, planPassword string) diag.Diagnostics {
	state.Id = types.StringPointerValue(r.Id)
	state.ImportId = types.StringPointerValue(r.Id)
	state.FileData = types.StringPointerValue(&planFileData)
	state.Format = types.StringPointerValue(&planFormat)
	state.Password = types.StringValue(planPassword)
	state.CryptoProvider = types.StringPointerValue(r.CryptoProvider)

	state.SerialNumber = types.StringPointerValue(r.SerialNumber)
	state.SubjectDN = types.StringPointerValue(r.SubjectDN)
	state.SubjectAlternativeNames = internaltypes.GetStringSet(r.SubjectAlternativeNames)
	state.IssuerDN = types.StringPointerValue(r.IssuerDN)
	state.ValidFrom = types.StringValue(r.GetValidFrom().Format(time.RFC3339))
	state.Expires = types.StringValue(r.GetExpires().Format(time.RFC3339))
	state.KeyAlgorithm = types.StringPointerValue(r.KeyAlgorithm)
	state.KeySize = types.Int64PointerValue(r.KeySize)
	state.SignatureAlgorithm = types.StringPointerValue(r.SignatureAlgorithm)
	state.Version = types.Int64PointerValue(r.Version)
	state.Sha1Fingerprint = types.StringPointerValue(r.Sha1Fingerprint)
	state.Sha256Fingerprint = types.StringPointerValue(r.Sha256Fingerprint)
	state.Status = types.StringPointerValue(r.Status)

	var valueFromDiags diag.Diagnostics
	state.RotationSettings, valueFromDiags = types.ObjectValueFrom(ctx, rotationSettingsAttrTypes, r.RotationSettings)
	return valueFromDiags
}

func (r *keyPairsSslServerImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan keyPairsSslServerImportResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createKeyPairsSslServerImport := client.NewKeyPairFile(plan.FileData.ValueString(), plan.Password.ValueString())
	err := addOptionalKeyPairsSslServerImportFields(ctx, createKeyPairsSslServerImport, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for KeyPair SSL Server Import", err.Error())
		return
	}

	apiCreateKeyPairsSslServerImport := r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateKeyPairsSslServerImport = apiCreateKeyPairsSslServerImport.Body(*createKeyPairsSslServerImport)
	keyPairsSslServerImportResponse, httpResp, err := r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPairExecute(apiCreateKeyPairsSslServerImport)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the KeyPair SSL Server Import", err, httpResp)
		return
	}

	// Read the response into the state
	var state keyPairsSslServerImportResourceModel
	diags = readKeyPairsSslServerImportResponse(ctx, keyPairsSslServerImportResponse, &state,
		plan.FileData.ValueString(), plan.Format.ValueString(), plan.Password.ValueString())
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *keyPairsSslServerImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state keyPairsSslServerImportResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadKeyPairsSslServerImport, httpResp, err := r.apiClient.KeyPairsSslServerAPI.GetSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ImportId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair SSL Server Import", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair SSL Server Import", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readKeyPairsSslServerImportResponse(ctx, apiReadKeyPairsSslServerImport, &state,
		state.FileData.ValueString(), state.Format.ValueString(), state.Password.ValueString())
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
// Since all non-computed attributes require replacing the resource, there is no need to implement Update for this resource.
func (r *keyPairsSslServerImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *keyPairsSslServerImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state keyPairsSslServerImportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.KeyPairsSslServerAPI.DeleteSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.ImportId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a KeyPair SSL Server Import", err, httpResp)
	}
}

func (r *keyPairsSslServerImportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("import_id"), req, resp)
}
