package keypairssslclient

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &keypairsSslClientKeyResource{}
	_ resource.ResourceWithConfigure = &keypairsSslClientKeyResource{}
)

// KeypairsSslClientKeyResource is a helper function to simplify the provider implementation.
func KeypairsSslClientKeyResource() resource.Resource {
	return &keypairsSslClientKeyResource{}
}

// keypairsSslClientKeyResource is the resource implementation.
type keypairsSslClientKeyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keypairsSslClientKeyResourceModel struct {
	City                    types.String `tfsdk:"city"`
	CommonName              types.String `tfsdk:"common_name"`
	Country                 types.String `tfsdk:"country"`
	CryptoProvider          types.String `tfsdk:"crypto_provider"`
	Expires                 types.String `tfsdk:"expires"`
	FileData                types.String `tfsdk:"file_data"`
	Format                  types.String `tfsdk:"format"`
	Id                      types.String `tfsdk:"id"`
	IssuerDn                types.String `tfsdk:"issuer_dn"`
	KeyAlgorithm            types.String `tfsdk:"key_algorithm"`
	KeyId                   types.String `tfsdk:"key_id"`
	KeySize                 types.Int64  `tfsdk:"key_size"`
	Organization            types.String `tfsdk:"organization"`
	OrganizationUnit        types.String `tfsdk:"organization_unit"`
	Password                types.String `tfsdk:"password"`
	SerialNumber            types.String `tfsdk:"serial_number"`
	Sha1Fingerprint         types.String `tfsdk:"sha1_fingerprint"`
	Sha256Fingerprint       types.String `tfsdk:"sha256_fingerprint"`
	SignatureAlgorithm      types.String `tfsdk:"signature_algorithm"`
	State                   types.String `tfsdk:"state"`
	Status                  types.String `tfsdk:"status"`
	SubjectAlternativeNames types.Set    `tfsdk:"subject_alternative_names"`
	SubjectDn               types.String `tfsdk:"subject_dn"`
	ValidDays               types.Int64  `tfsdk:"valid_days"`
	ValidFrom               types.String `tfsdk:"valid_from"`
	Version                 types.Int64  `tfsdk:"version"`
}

// GetSchema defines the schema for the resource.
func (r *keypairsSslClientKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage ssl client key pairs.",
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replace plan if changed.",
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
			"file_data": schema.StringAttribute{
				Description: "Base-64 encoded PKCS12 or PEM file data. In the case of PEM, the raw (non-base-64) data is also accepted. In BCFIPS mode, only PEM with PBES2 and AES or Triple DES encryption is accepted and 128-bit salt is required. If not configured, the new key will be generated. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"format": schema.StringAttribute{
				Description: "Key pair file format. If specified, this field will control what file format is expected, otherwise the format will be auto-detected. In BCFIPS mode, only `PEM` is supported. Supported values are `PKCS12` and `PEM`. Can only be configured if `file_data` is set. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("PKCS12", "PEM"),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for the file. In BCFIPS mode, the password must be at least 14 characters. Must be configured if `file_data` is set, otherwise cannot be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Supported values are `LOCAL` and `HSM`. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("LOCAL", "HSM"),
				},
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
				Description: "The subject alternative names (SAN). Cannot be configured if `file_data` is set. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"issuer_dn": schema.StringAttribute{
				Description: "The issuer's distinguished name",
				Computed:    true,
			},
			"common_name": schema.StringAttribute{
				Description: "Common name for key pair subject. Required if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"organization": schema.StringAttribute{
				Description: "Organization for generating the key pair. Optional if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"organization_unit": schema.StringAttribute{
				Description: "Organization unit for generating the key pair. Optional if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"city": schema.StringAttribute{
				Description: "City for generating the key pair. Optional if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"state": schema.StringAttribute{
				Description: "State for generating the key pair. Optional if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"country": schema.StringAttribute{
				Description: "Country for generating the key pair. Required if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"valid_from": schema.StringAttribute{
				Description: "The start date from which the item is valid, in ISO 8601 format (UTC). This field is immutable and will trigger a replace plan if changed.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"valid_days": schema.Int64Attribute{
				Description: "Number of days the key pair will be valid for. Required if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"expires": schema.StringAttribute{
				Description: "The end date up until which the item is valid, in ISO 8601 format (UTC)",
				Computed:    true,
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The public key algorithm. Required if `file_data` is not set, otherwise can't be configured. This field is immutable and will trigger a replace plan if changed. Typically supported values are `RSA` and `EC`.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"key_size": schema.Int64Attribute{
				Description: "The public key size, in bits. Can only be configured if `file_data` is not set. If not configured and `file_data` is not set, then the default size for the key algorithm will be used. This field is immutable and will trigger a replace plan if changed. Typically supported values are `256`, `384`, and `521` for EC keys and `1024`, `2048`, and `4096` for RSA keys.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"signature_algorithm": schema.StringAttribute{
				Description: "The signature algorithm. Can only be configured if `file_data` is not set. If not configured and `file_data` is not set, then the default signature algorithm for the key algorithm will be used. This field is immutable and will trigger a replace plan if changed. Typically supported values are `SHA256withECDSA`, `SHA384withECDSA`, and `SHA512withECDSA` for EC keys, and `SHA256withRSA`, `SHA384withRSA`, and `SHA512withRSA` for RSA keys.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
	id.ToSchema(&resp.Schema)
}

func (r *keypairsSslClientKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *keypairsSslClientKeyResourceModel
	var config *keypairsSslClientKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	if internaltypes.IsDefined(plan.FileData) {
		// The key will be imported from file_data
		if plan.Password.IsNull() {
			resp.Diagnostics.AddError("password must be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.CommonName) {
			resp.Diagnostics.AddError("common_name cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.Organization) {
			resp.Diagnostics.AddError("organization cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.OrganizationUnit) {
			resp.Diagnostics.AddError("organization_unit cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.City) {
			resp.Diagnostics.AddError("city cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.State) {
			resp.Diagnostics.AddError("state cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.Country) {
			resp.Diagnostics.AddError("country cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.ValidDays) {
			resp.Diagnostics.AddError("valid_days cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.KeyAlgorithm) {
			resp.Diagnostics.AddError("key_algorithm cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.KeySize) {
			resp.Diagnostics.AddError("key_size cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.SignatureAlgorithm) {
			resp.Diagnostics.AddError("signature_algorithm cannot be configured when file_data is set", "")
		}
		if internaltypes.IsDefined(config.SubjectAlternativeNames) {
			resp.Diagnostics.AddError("subject_alternative_names cannot be configured when file_data is set", "")
		}
	} else {
		// The key will be generated
		if internaltypes.IsDefined(plan.Format) {
			resp.Diagnostics.AddError("format cannot be configured when file_data is not set", "")
		}
		if internaltypes.IsDefined(plan.Password) {
			resp.Diagnostics.AddError("password cannot be configured when file_data is not set", "")
		}
		if !internaltypes.IsDefined(plan.CommonName) {
			resp.Diagnostics.AddError("common_name must be configured when file_data is not set", "")
		}
		if !internaltypes.IsDefined(plan.Organization) {
			resp.Diagnostics.AddError("organization must be configured when file_data is not set", "")
		}
		if !internaltypes.IsDefined(plan.Country) {
			resp.Diagnostics.AddError("country must be configured when file_data is not set", "")
		}
		if !internaltypes.IsDefined(plan.ValidDays) {
			resp.Diagnostics.AddError("valid_days must be configured when file_data is not set", "")
		}
		if !internaltypes.IsDefined(plan.KeyAlgorithm) {
			resp.Diagnostics.AddError("key_algorithm must be configured when file_data is not set", "")
		}
	}
}

func (model *keypairsSslClientKeyResourceModel) buildGenerateClientStruct() (*client.NewKeyPairSettings, diag.Diagnostics) {
	result := &client.NewKeyPairSettings{}
	// city
	result.City = model.City.ValueStringPointer()
	// common_name
	result.CommonName = model.CommonName.ValueString()
	// country
	result.Country = model.Country.ValueString()
	// crypto_provider
	result.CryptoProvider = model.CryptoProvider.ValueStringPointer()
	// key_algorithm
	result.KeyAlgorithm = model.KeyAlgorithm.ValueString()
	// key_id
	result.Id = model.KeyId.ValueStringPointer()
	// key_size
	if internaltypes.IsDefined(model.KeySize) {
		result.KeySize = model.KeySize.ValueInt64Pointer()
	}
	// organization
	result.Organization = model.Organization.ValueString()
	// organization_unit
	result.OrganizationUnit = model.OrganizationUnit.ValueStringPointer()
	// signature_algorithm
	result.SignatureAlgorithm = model.SignatureAlgorithm.ValueStringPointer()
	// state
	result.State = model.State.ValueStringPointer()
	// subject_alternative_names
	if !model.SubjectAlternativeNames.IsNull() {
		result.SubjectAlternativeNames = []string{}
		for _, subjectAlternativeNamesElement := range model.SubjectAlternativeNames.Elements() {
			result.SubjectAlternativeNames = append(result.SubjectAlternativeNames, subjectAlternativeNamesElement.(types.String).ValueString())
		}
	}

	// valid_days
	result.ValidDays = model.ValidDays.ValueInt64()
	return result, nil
}

func (model *keypairsSslClientKeyResourceModel) buildImportClientStruct() (*client.KeyPairFile, diag.Diagnostics) {
	result := &client.KeyPairFile{}
	// crypto_provider
	result.CryptoProvider = model.CryptoProvider.ValueStringPointer()
	// file_data
	result.FileData = model.FileData.ValueString()
	// format
	result.Format = model.Format.ValueStringPointer()
	// key_id
	result.Id = model.KeyId.ValueStringPointer()
	// password
	result.Password = model.Password.ValueString()
	return result, nil
}

func (state *keypairsSslClientKeyResourceModel) readClientResponse(response *client.KeyPairView) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringPointerValue(response.Id)
	// crypto_provider
	state.CryptoProvider = types.StringPointerValue(response.CryptoProvider)
	// expires
	state.Expires = types.StringValue(response.Expires.Format(time.RFC3339))
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

// Metadata returns the resource type name.
func (r *keypairsSslClientKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_client_key"
}

func (r *keypairsSslClientKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *keypairsSslClientKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSslClientKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	var responseData *client.KeyPairView
	var httpResp *http.Response
	var err error
	if !internaltypes.IsDefined(data.FileData) {
		clientData, diags := data.buildGenerateClientStruct()
		resp.Diagnostics.Append(diags...)
		apiCreateRequest := r.apiClient.KeyPairsSslClientAPI.CreateSslClientKeyPair(config.AuthContext(ctx, r.providerConfig))
		apiCreateRequest = apiCreateRequest.Body(*clientData)
		responseData, httpResp, err = r.apiClient.KeyPairsSslClientAPI.CreateSslClientKeyPairExecute(apiCreateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while generating the ssl client key", err, httpResp)
			return
		}
	} else {
		clientData, diags := data.buildImportClientStruct()
		resp.Diagnostics.Append(diags...)
		apiCreateRequest := r.apiClient.KeyPairsSslClientAPI.ImportSslClientKeyPair(config.AuthContext(ctx, r.providerConfig))
		apiCreateRequest = apiCreateRequest.Body(*clientData)
		responseData, httpResp, err = r.apiClient.KeyPairsSslClientAPI.ImportSslClientKeyPairExecute(apiCreateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while importing the ssl client key", err, httpResp)
			return
		}
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslClientKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keypairsSslClientKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KeyPairsSslClientAPI.GetSslClientKeyPair(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "SSL Client Key Pair", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the key pair", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
// Since all non-computed attributes require replacing the resource, there is no need to implement Update for this resource.
func (r *keypairsSslClientKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *keypairsSslClientKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data keypairsSslClientKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.KeyPairsSslClientAPI.DeleteSslClientKeyPair(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the ssl client key", err, httpResp)
	}
}
