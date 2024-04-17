package keypairsigninggenerate

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &keyPairSigningGenerateResource{}
	_ resource.ResourceWithConfigure   = &keyPairSigningGenerateResource{}
	_ resource.ResourceWithImportState = &keyPairSigningGenerateResource{}

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

// KeyPairSigningGenerateResource is a helper function to simplify the provider implementation.
func KeyPairSigningGenerateResource() resource.Resource {
	return &keyPairSigningGenerateResource{}
}

// keyPairSigningGenerateResource is the resource implementation.
type keyPairSigningGenerateResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keyPairSigningGenerateResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	GenerateId              types.String `tfsdk:"generate_id"`
	CommonName              types.String `tfsdk:"common_name"`
	SubjectAlternativeNames types.Set    `tfsdk:"subject_alternative_names"`
	Organization            types.String `tfsdk:"organization"`
	OrganizationUnit        types.String `tfsdk:"organization_unit"`
	City                    types.String `tfsdk:"city"`
	State                   types.String `tfsdk:"state"`
	Country                 types.String `tfsdk:"country"`
	ValidDays               types.Int64  `tfsdk:"valid_days"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
	KeySize                 types.Int64  `tfsdk:"key_size"`
	SignatureAlgorithm      types.String `tfsdk:"signature_algorithm"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
}

type keyPairSigningGenerateResourceResponseModel struct {
	Id                      types.String `tfsdk:"id"`
	GenerateId              types.String `tfsdk:"import_id"`
	Format                  types.String `tfsdk:"format"`
	Password                types.String `tfsdk:"password"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
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
func (r *keyPairSigningGenerateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Key Pairs Signing Generate",
		Attributes:  map[string]schema.Attribute{},
	}

	resp.Schema = schema
}

func addOptionalKeyPairSigningGenerateFields(ctx context.Context, addRequest *client.NewKeyPairSettings, plan keyPairSigningGenerateResourceModel) {

	if internaltypes.IsDefined(plan.GenerateId) {
		addRequest.Id = plan.GenerateId.ValueStringPointer()
	}

	var slice []string
	plan.SubjectAlternativeNames.ElementsAs(ctx, &slice, false)
	addRequest.SubjectAlternativeNames = slice

	addRequest.OrganizationUnit = plan.OrganizationUnit.ValueStringPointer()
	addRequest.City = plan.City.ValueStringPointer()
	addRequest.State = plan.State.ValueStringPointer()
	addRequest.KeySize = plan.KeySize.ValueInt64Pointer()
	addRequest.SignatureAlgorithm = plan.SignatureAlgorithm.ValueStringPointer()
	addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
}

// Metadata returns the resource type name.
func (r *keyPairSigningGenerateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair_signing_generate"
}

func (r *keyPairSigningGenerateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readKeyPairSigningGenerateResponse(ctx context.Context, r *client.KeyPairView, state *keyPairSigningGenerateResourceResponseModel) diag.Diagnostics {
	state.Id = types.StringPointerValue(r.Id)
	state.GenerateId = types.StringPointerValue(r.Id)
	state.CryptoProvider = types.StringPointerValue(r.CryptoProvider)
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

	var valueFromDiags diag.Diagnostics
	state.RotationSettings, valueFromDiags = types.ObjectValueFrom(ctx, rotationSettingsAttrTypes, r.RotationSettings)
	return valueFromDiags
}

func (r *keyPairSigningGenerateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan keyPairSigningGenerateResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createKeyPairSigningGenerate := client.NewNewKeyPairSettings(plan.CommonName.ValueString(), plan.Organization.ValueString(), plan.Country.ValueString(), plan.ValidDays.ValueInt64(), plan.KeyAlgorithm.ValueString())
	addOptionalKeyPairSigningGenerateFields(ctx, createKeyPairSigningGenerate, plan)

	apiCreateKeyPairSigningGenerate := r.apiClient.KeyPairsSslServerAPI.CreateSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateKeyPairSigningGenerate = apiCreateKeyPairSigningGenerate.Body(*createKeyPairSigningGenerate)
	keyPairSigningGenerateResponse, httpResp, err := r.apiClient.KeyPairsSslServerAPI.CreateSslServerKeyPairExecute(apiCreateKeyPairSigningGenerate)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the KeyPairSigningGenerate", err, httpResp)
		return
	}

	// Read the response into the state
	var state keyPairSigningGenerateResourceResponseModel

	diags = readKeyPairSigningGenerateResponse(ctx, keyPairSigningGenerateResponse, &state)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *keyPairSigningGenerateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state keyPairSigningGenerateResourceResponseModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadKeyPairSigningGenerate, httpResp, err := r.apiClient.KeyPairsSigningAPI.GetSigningKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.GenerateId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPairSigningGenerate", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  KeyPairSigningGenerate", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readKeyPairSigningGenerateResponse(ctx, apiReadKeyPairSigningGenerate, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// All changes forces replacement
func (r *keyPairSigningGenerateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// This config object is edit-only, so Terraform can't delete it.
func (r *keyPairSigningGenerateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *keyPairSigningGenerateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("generate_id"), req, resp)
}
