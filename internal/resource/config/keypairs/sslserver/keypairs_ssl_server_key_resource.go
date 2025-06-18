// Copyright © 2025 Ping Identity Corporation

package keypairssslserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &keypairsSslServerKeyResource{}
	_ resource.ResourceWithConfigure = &keypairsSslServerKeyResource{}

	customId    = "key_id"
	createMutex sync.Mutex
)

// KeypairsSslServerKeyResource is a helper function to simplify the provider implementation.
func KeypairsSslServerKeyResource() resource.Resource {
	return &keypairsSslServerKeyResource{}
}

// keypairsSslServerKeyResource is the resource implementation.
type keypairsSslServerKeyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keypairsSslServerKeyResourceModel struct {
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
	RotationSettings        types.Object `tfsdk:"rotation_settings"`
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
func (r *keypairsSslServerKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage ssl server key pairs.",
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
			"rotation_settings": schema.SingleNestedAttribute{
				Description: "The local identity profile data store configuration.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
						Computed:    true,
					},
					"creation_buffer_days": schema.Int64Attribute{
						Description: "Buffer days before key pair expiration for creation of a new key pair.",
						Computed:    true,
					},
					"activation_buffer_days": schema.Int64Attribute{
						Description: "Buffer days before key pair expiration for activation of the new key pair.",
						Computed:    true,
					},
					"valid_days": schema.Int64Attribute{
						Description: "Valid days for the new key pair to be created. If this property is unset, the validity days of the original key pair will be used.",
						Computed:    true,
					},
					"key_algorithm": schema.StringAttribute{
						Description: "Key algorithm to be used while creating a new key pair. If this property is unset, the key algorithm of the original key pair will be used. Supported algorithms are available through the /keyPairs/keyAlgorithms endpoint.",
						Computed:    true,
					},
					"key_size": schema.Int64Attribute{
						Description: "Key size, in bits. If this property is unset, the key size of the original key pair will be used. Supported key sizes are available through the /keyPairs/keyAlgorithms endpoint.",
						Computed:    true,
					},
					"signature_algorithm": schema.StringAttribute{
						Description: "Required if the original key pair used SHA1 algorithm. If this property is unset, the default signature algorithm of the original key pair will be used. Supported signature algorithms are available through the /keyPairs/keyAlgorithms endpoint.",
						Computed:    true,
					},
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (r *keypairsSslServerKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *keypairsSslServerKeyResourceModel
	var config *keypairsSslServerKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	if internaltypes.IsDefined(plan.FileData) {
		// The key will be imported from file_data
		if plan.Password.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				providererror.InvalidAttributeConfiguration,
				"password must be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.CommonName) {
			resp.Diagnostics.AddAttributeError(
				path.Root("common_name"),
				providererror.InvalidAttributeConfiguration,
				"common_name cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.Organization) {
			resp.Diagnostics.AddAttributeError(
				path.Root("organization"),
				providererror.InvalidAttributeConfiguration,
				"organization cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.OrganizationUnit) {
			resp.Diagnostics.AddAttributeError(
				path.Root("organization_unit"),
				providererror.InvalidAttributeConfiguration,
				"organization_unit cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.City) {
			resp.Diagnostics.AddAttributeError(
				path.Root("city"),
				providererror.InvalidAttributeConfiguration,
				"city cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.State) {
			resp.Diagnostics.AddAttributeError(
				path.Root("state"),
				providererror.InvalidAttributeConfiguration,
				"state cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.Country) {
			resp.Diagnostics.AddAttributeError(
				path.Root("country"),
				providererror.InvalidAttributeConfiguration,
				"country cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.ValidDays) {
			resp.Diagnostics.AddAttributeError(
				path.Root("valid_days"),
				providererror.InvalidAttributeConfiguration,
				"valid_days cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.KeyAlgorithm) {
			resp.Diagnostics.AddAttributeError(
				path.Root("key_algorithm"),
				providererror.InvalidAttributeConfiguration,
				"key_algorithm cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.KeySize) {
			resp.Diagnostics.AddAttributeError(
				path.Root("key_size"),
				providererror.InvalidAttributeConfiguration,
				"key_size cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.SignatureAlgorithm) {
			resp.Diagnostics.AddAttributeError(
				path.Root("signature_algorithm"),
				providererror.InvalidAttributeConfiguration,
				"signature_algorithm cannot be configured when file_data is set")
		}
		if internaltypes.IsDefined(config.SubjectAlternativeNames) {
			resp.Diagnostics.AddAttributeError(
				path.Root("subject_alternative_names"),
				providererror.InvalidAttributeConfiguration,
				"subject_alternative_names cannot be configured when file_data is set")
		}
	} else {
		// The key will be generated
		if internaltypes.IsDefined(plan.Format) {
			resp.Diagnostics.AddAttributeError(
				path.Root("format"),
				providererror.InvalidAttributeConfiguration,
				"format cannot be configured when file_data is not set")
		}
		if internaltypes.IsDefined(plan.Password) {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				providererror.InvalidAttributeConfiguration,
				"password cannot be configured when file_data is not set")
		}
		if !internaltypes.IsDefined(plan.CommonName) {
			resp.Diagnostics.AddAttributeError(
				path.Root("common_name"),
				providererror.InvalidAttributeConfiguration,
				"common_name must be configured when file_data is not set")
		}
		if !internaltypes.IsDefined(plan.Organization) {
			resp.Diagnostics.AddAttributeError(
				path.Root("organization"),
				providererror.InvalidAttributeConfiguration,
				"organization must be configured when file_data is not set")
		}
		if !internaltypes.IsDefined(plan.Country) {
			resp.Diagnostics.AddAttributeError(
				path.Root("country"),
				providererror.InvalidAttributeConfiguration,
				"country must be configured when file_data is not set")
		}
		if !internaltypes.IsDefined(plan.ValidDays) {
			resp.Diagnostics.AddAttributeError(
				path.Root("valid_days"),
				providererror.InvalidAttributeConfiguration,
				"valid_days must be configured when file_data is not set")
		}
		if !internaltypes.IsDefined(plan.KeyAlgorithm) {
			resp.Diagnostics.AddAttributeError(
				path.Root("key_algorithm"),
				providererror.InvalidAttributeConfiguration,
				"key_algorithm must be configured when file_data is not set")
		}
	}
}

func (model *keypairsSslServerKeyResourceModel) buildGenerateClientStruct() (*client.NewKeyPairSettings, diag.Diagnostics) {
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

func (model *keypairsSslServerKeyResourceModel) buildImportClientStruct() (*client.KeyPairFile, diag.Diagnostics) {
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

func (state *keypairsSslServerKeyResourceModel) readClientResponse(response *client.KeyPairView) diag.Diagnostics {
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
	// rotation_settings
	rotationSettingsAttrTypes := map[string]attr.Type{
		"activation_buffer_days": types.Int64Type,
		"creation_buffer_days":   types.Int64Type,
		"id":                     types.StringType,
		"key_algorithm":          types.StringType,
		"key_size":               types.Int64Type,
		"signature_algorithm":    types.StringType,
		"valid_days":             types.Int64Type,
	}
	var rotationSettingsValue types.Object
	if response.RotationSettings == nil {
		rotationSettingsValue = types.ObjectNull(rotationSettingsAttrTypes)
	} else {
		rotationSettingsValue, diags = types.ObjectValue(rotationSettingsAttrTypes, map[string]attr.Value{
			"activation_buffer_days": types.Int64Value(response.RotationSettings.ActivationBufferDays),
			"creation_buffer_days":   types.Int64Value(response.RotationSettings.CreationBufferDays),
			"id":                     types.StringPointerValue(response.RotationSettings.Id),
			"key_algorithm":          types.StringPointerValue(response.RotationSettings.KeyAlgorithm),
			"key_size":               types.Int64PointerValue(response.RotationSettings.KeySize),
			"signature_algorithm":    types.StringPointerValue(response.RotationSettings.SignatureAlgorithm),
			"valid_days":             types.Int64PointerValue(response.RotationSettings.ValidDays),
		})
		respDiags.Append(diags...)
	}

	state.RotationSettings = rotationSettingsValue
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
func (r *keypairsSslServerKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_ssl_server_key"
}

func (r *keypairsSslServerKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *keypairsSslServerKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSslServerKeyResourceModel

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
		apiCreateRequest := r.apiClient.KeyPairsSslServerAPI.CreateSslServerKeyPair(config.AuthContext(ctx, r.providerConfig))
		apiCreateRequest = apiCreateRequest.Body(*clientData)
		createMutex.Lock()
		responseData, httpResp, err = r.apiClient.KeyPairsSslServerAPI.CreateSslServerKeyPairExecute(apiCreateRequest)
		createMutex.Unlock()
		if err != nil {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while generating the ssl server key", err, httpResp, &customId)
			return
		}
	} else {
		clientData, diags := data.buildImportClientStruct()
		resp.Diagnostics.Append(diags...)
		apiCreateRequest := r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPair(config.AuthContext(ctx, r.providerConfig))
		apiCreateRequest = apiCreateRequest.Body(*clientData)
		createMutex.Lock()
		responseData, httpResp, err = r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPairExecute(apiCreateRequest)
		createMutex.Unlock()
		if err != nil {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while importing the ssl server key", err, httpResp, &customId)
			return
		}
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSslServerKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keypairsSslServerKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KeyPairsSslServerAPI.GetSslServerKeyPair(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "SSL Server Key Pair", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the key pair", err, httpResp, &customId)
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
func (r *keypairsSslServerKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *keypairsSslServerKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data keypairsSslServerKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.KeyPairsSslServerAPI.DeleteSslServerKeyPair(config.AuthContext(ctx, r.providerConfig), data.KeyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the ssl server key", err, httpResp, &customId)
	}
}
