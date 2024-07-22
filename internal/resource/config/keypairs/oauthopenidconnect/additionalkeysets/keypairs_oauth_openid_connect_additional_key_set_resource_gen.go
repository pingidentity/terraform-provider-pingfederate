// Code generated by ping-terraform-plugin-framework-generator

package keypairsoauthopenidconnectadditionalkeysets

import (
	"context"

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
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

var (
	_ resource.Resource                = &keypairsOauthOpenidConnectAdditionalKeySetResource{}
	_ resource.ResourceWithConfigure   = &keypairsOauthOpenidConnectAdditionalKeySetResource{}
	_ resource.ResourceWithImportState = &keypairsOauthOpenidConnectAdditionalKeySetResource{}
)

func KeypairsOauthOpenidConnectAdditionalKeySetResource() resource.Resource {
	return &keypairsOauthOpenidConnectAdditionalKeySetResource{}
}

type keypairsOauthOpenidConnectAdditionalKeySetResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_oauth_openid_connect_additional_key_set"
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsOauthOpenidConnectAdditionalKeySetResourceModel struct {
	Description types.String `tfsdk:"description"`
	Issuers     types.List   `tfsdk:"issuers"`
	Name        types.String `tfsdk:"name"`
	SetId       types.String `tfsdk:"set_id"`
	SigningKeys types.Object `tfsdk:"signing_keys"`
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage OAuth/OpenID Connect additional signing key sets.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the key set.",
			},
			"issuers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the resource.",
						},
					},
				},
				Required:    true,
				Description: "A list of virtual issuers that will use the current key set. Once assigned to a key set, the same virtual issuer cannot be assigned to another key set instance.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The key set name.",
			},
			"set_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique ID for the key set. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"signing_keys": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"p256_active_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-256 key currently active.",
					},
					"p256_active_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for currently active P-256 key.",
					},
					"p256_previous_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-256 key previously active.",
					},
					"p256_previous_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for previously active P-256 key.",
					},
					"p256_publish_x5c_parameter": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Enable publishing of the P-256 certificate chain associated with the active key.",
					},
					"p384_active_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-384 key currently active.",
					},
					"p384_active_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for currently active P-384 key.",
					},
					"p384_previous_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-384 key previously active.",
					},
					"p384_previous_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for previously active P-384 key.",
					},
					"p384_publish_x5c_parameter": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Enable publishing of the P-384 certificate chain associated with the active key.",
					},
					"p521_active_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-521 key currently active.",
					},
					"p521_active_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for currently active P-521 key.",
					},
					"p521_previous_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the P-521 key previously active.",
					},
					"p521_previous_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for previously active P-521 key.",
					},
					"p521_publish_x5c_parameter": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable publishing of the P-521 certificate chain associated with the active key.",
					},
					"rsa_active_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Required:    true,
						Description: "Reference to the RSA key currently active.",
					},
					"rsa_active_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for currently active RSA key.",
					},
					"rsa_algorithm_active_key_ids": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key_id": schema.StringAttribute{
									Required:    true,
									Description: "Unique key identifier.",
								},
								"rsa_alg_type": schema.StringAttribute{
									Required:    true,
									Description: "The RSA signing algorithm type. The supported RSA signing algorithm types are `RS256`, `RS384`, `RS512`, `PS256`, `PS384` and `PS512`.",
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Description: "PingFederate uses the same RSA key for all RSA signing algorithms. To enable active RSA JWK entry to have unique single valued ''alg'' parameter, use this list to set a key identifier for each RSA algorithm (`RS256`, `RS384`, `RS512`, `PS256`, `PS384` and `PS512`).",
					},
					"rsa_algorithm_previous_key_ids": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key_id": schema.StringAttribute{
									Required:    true,
									Description: "Unique key identifier.",
								},
								"rsa_alg_type": schema.StringAttribute{
									Required:    true,
									Description: "The RSA signing algorithm type. The supported RSA signing algorithm types are `RS256`, `RS384`, `RS512`, `PS256`, `PS384` and `PS512`.",
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Description: "PingFederate uses the same RSA key for all RSA signing algorithms. To enable previously active RSA JWK entry to have unique single valued ''alg'' parameter, use this list to set a key identifier for each RSA algorithm (`RS256`, `RS384`, `RS512`, `PS256`, `PS384` and `PS512`).",
					},
					"rsa_previous_cert_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Reference to the RSA key previously active.",
					},
					"rsa_previous_key_id": schema.StringAttribute{
						Optional:    true,
						Description: "Key Id for previously active RSA key.",
					},
					"rsa_publish_x5c_parameter": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Enable publishing of the RSA certificate chain associated with the active key. The default value is `false`.",
					},
				},
				Required:    true,
				Description: "Setting for a OAuth/OpenID Connect signing key set while using multiple virtual issuers.",
			},
		},
	}
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.0.1 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1201)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1201 := compare >= 0
	var plan *keypairsOauthOpenidConnectAdditionalKeySetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	// If any of these fields are set by the user and the PF version is not new enough, throw an error
	if !pfVersionAtLeast1201 {
		if internaltypes.IsDefined(plan.SigningKeys) {
			p256ActiveKeyId := plan.SigningKeys.Attributes()["p256_active_key_id"]
			if internaltypes.IsDefined(p256ActiveKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p256_active_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			p256PreviousKeyId := plan.SigningKeys.Attributes()["p256_previous_key_id"]
			if internaltypes.IsDefined(p256PreviousKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p256_previous_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			p384ActiveKeyId := plan.SigningKeys.Attributes()["p384_active_key_id"]
			if internaltypes.IsDefined(p384ActiveKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p384_active_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			p384PreviousKeyId := plan.SigningKeys.Attributes()["p384_previous_key_id"]
			if internaltypes.IsDefined(p384PreviousKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p384_previous_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			p521ActiveKeyId := plan.SigningKeys.Attributes()["p521_active_key_id"]
			if internaltypes.IsDefined(p521ActiveKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p521_active_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			p521PreviousKeyId := plan.SigningKeys.Attributes()["p521_previous_key_id"]
			if internaltypes.IsDefined(p521PreviousKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.p521_previous_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			rsaActiveKeyId := plan.SigningKeys.Attributes()["rsa_active_key_id"]
			if internaltypes.IsDefined(rsaActiveKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.rsa_active_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			rsaAlgorithmActiveKeyIds := plan.SigningKeys.Attributes()["rsa_algorithm_active_key_ids"]
			if internaltypes.IsDefined(rsaAlgorithmActiveKeyIds) {
				version.AddUnsupportedAttributeError("signing_keys.rsa_algorithm_active_key_ids",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			rsaAlgorithmPreviousKeyIds := plan.SigningKeys.Attributes()["rsa_algorithm_previous_key_ids"]
			if internaltypes.IsDefined(rsaAlgorithmPreviousKeyIds) {
				version.AddUnsupportedAttributeError("signing_keys.rsa_algorithm_previous_key_ids",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
			rsaPreviousKeyId := plan.SigningKeys.Attributes()["rsa_previous_key_id"]
			if internaltypes.IsDefined(rsaPreviousKeyId) {
				version.AddUnsupportedAttributeError("signing_keys.rsa_previous_key_id",
					r.providerConfig.ProductVersion, version.PingFederate1201, &resp.Diagnostics)
			}
		}
	}
	// Set default values that can't be set in schema
	r.setConditionalDefaults(ctx, pfVersionAtLeast1201, plan, resp)
}

func (model *keypairsOauthOpenidConnectAdditionalKeySetResourceModel) buildClientStruct(versionAtLeast1201 bool) *client.AdditionalKeySet {
	result := &client.AdditionalKeySet{}
	// description
	result.Description = model.Description.ValueStringPointer()
	// issuers
	result.Issuers = []client.ResourceLink{}
	for _, issuersElement := range model.Issuers.Elements() {
		issuersValue := client.ResourceLink{}
		issuersAttrs := issuersElement.(types.Object).Attributes()
		issuersValue.Id = issuersAttrs["id"].(types.String).ValueString()
		result.Issuers = append(result.Issuers, issuersValue)
	}

	// name
	result.Name = model.Name.ValueString()
	// set_id
	result.Id = model.SetId.ValueStringPointer()
	// signing_keys
	signingKeysValue := client.SigningKeys{}
	signingKeysAttrs := model.SigningKeys.Attributes()
	if !signingKeysAttrs["p256_active_cert_ref"].IsNull() {
		signingKeysP256ActiveCertRefValue := &client.ResourceLink{}
		signingKeysP256ActiveCertRefAttrs := signingKeysAttrs["p256_active_cert_ref"].(types.Object).Attributes()
		signingKeysP256ActiveCertRefValue.Id = signingKeysP256ActiveCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P256ActiveCertRef = signingKeysP256ActiveCertRefValue
	}
	signingKeysValue.P256ActiveKeyId = signingKeysAttrs["p256_active_key_id"].(types.String).ValueStringPointer()
	if !signingKeysAttrs["p256_previous_cert_ref"].IsNull() {
		signingKeysP256PreviousCertRefValue := &client.ResourceLink{}
		signingKeysP256PreviousCertRefAttrs := signingKeysAttrs["p256_previous_cert_ref"].(types.Object).Attributes()
		signingKeysP256PreviousCertRefValue.Id = signingKeysP256PreviousCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P256PreviousCertRef = signingKeysP256PreviousCertRefValue
	}
	signingKeysValue.P256PreviousKeyId = signingKeysAttrs["p256_previous_key_id"].(types.String).ValueStringPointer()
	signingKeysValue.P256PublishX5cParameter = signingKeysAttrs["p256_publish_x5c_parameter"].(types.Bool).ValueBoolPointer()
	if !signingKeysAttrs["p384_active_cert_ref"].IsNull() {
		signingKeysP384ActiveCertRefValue := &client.ResourceLink{}
		signingKeysP384ActiveCertRefAttrs := signingKeysAttrs["p384_active_cert_ref"].(types.Object).Attributes()
		signingKeysP384ActiveCertRefValue.Id = signingKeysP384ActiveCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P384ActiveCertRef = signingKeysP384ActiveCertRefValue
	}
	signingKeysValue.P384ActiveKeyId = signingKeysAttrs["p384_active_key_id"].(types.String).ValueStringPointer()
	if !signingKeysAttrs["p384_previous_cert_ref"].IsNull() {
		signingKeysP384PreviousCertRefValue := &client.ResourceLink{}
		signingKeysP384PreviousCertRefAttrs := signingKeysAttrs["p384_previous_cert_ref"].(types.Object).Attributes()
		signingKeysP384PreviousCertRefValue.Id = signingKeysP384PreviousCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P384PreviousCertRef = signingKeysP384PreviousCertRefValue
	}
	signingKeysValue.P384PreviousKeyId = signingKeysAttrs["p384_previous_key_id"].(types.String).ValueStringPointer()
	signingKeysValue.P384PublishX5cParameter = signingKeysAttrs["p384_publish_x5c_parameter"].(types.Bool).ValueBoolPointer()
	if !signingKeysAttrs["p521_active_cert_ref"].IsNull() {
		signingKeysP521ActiveCertRefValue := &client.ResourceLink{}
		signingKeysP521ActiveCertRefAttrs := signingKeysAttrs["p521_active_cert_ref"].(types.Object).Attributes()
		signingKeysP521ActiveCertRefValue.Id = signingKeysP521ActiveCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P521ActiveCertRef = signingKeysP521ActiveCertRefValue
	}
	signingKeysValue.P521ActiveKeyId = signingKeysAttrs["p521_active_key_id"].(types.String).ValueStringPointer()
	if !signingKeysAttrs["p521_previous_cert_ref"].IsNull() {
		signingKeysP521PreviousCertRefValue := &client.ResourceLink{}
		signingKeysP521PreviousCertRefAttrs := signingKeysAttrs["p521_previous_cert_ref"].(types.Object).Attributes()
		signingKeysP521PreviousCertRefValue.Id = signingKeysP521PreviousCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.P521PreviousCertRef = signingKeysP521PreviousCertRefValue
	}
	signingKeysValue.P521PreviousKeyId = signingKeysAttrs["p521_previous_key_id"].(types.String).ValueStringPointer()
	signingKeysValue.P521PublishX5cParameter = signingKeysAttrs["p521_publish_x5c_parameter"].(types.Bool).ValueBoolPointer()
	if !signingKeysAttrs["rsa_active_cert_ref"].IsNull() {
		signingKeysRsaActiveCertRefValue := &client.ResourceLink{}
		signingKeysRsaActiveCertRefAttrs := signingKeysAttrs["rsa_active_cert_ref"].(types.Object).Attributes()
		signingKeysRsaActiveCertRefValue.Id = signingKeysRsaActiveCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.RsaActiveCertRef = signingKeysRsaActiveCertRefValue
	}
	signingKeysValue.RsaActiveKeyId = signingKeysAttrs["rsa_active_key_id"].(types.String).ValueStringPointer()

	// Key ids are added in version 12.0.1
	if versionAtLeast1201 {
		signingKeysValue.RsaAlgorithmActiveKeyIds = []client.RsaAlgKeyId{}
		for _, rsaAlgorithmActiveKeyIdsElement := range signingKeysAttrs["rsa_algorithm_active_key_ids"].(types.List).Elements() {
			rsaAlgorithmActiveKeyIdsValue := client.RsaAlgKeyId{}
			rsaAlgorithmActiveKeyIdsAttrs := rsaAlgorithmActiveKeyIdsElement.(types.Object).Attributes()
			rsaAlgorithmActiveKeyIdsValue.KeyId = rsaAlgorithmActiveKeyIdsAttrs["key_id"].(types.String).ValueString()
			rsaAlgorithmActiveKeyIdsValue.RsaAlgType = rsaAlgorithmActiveKeyIdsAttrs["rsa_alg_type"].(types.String).ValueString()
			signingKeysValue.RsaAlgorithmActiveKeyIds = append(signingKeysValue.RsaAlgorithmActiveKeyIds, rsaAlgorithmActiveKeyIdsValue)
		}
		signingKeysValue.RsaAlgorithmPreviousKeyIds = []client.RsaAlgKeyId{}
		for _, rsaAlgorithmPreviousKeyIdsElement := range signingKeysAttrs["rsa_algorithm_previous_key_ids"].(types.List).Elements() {
			rsaAlgorithmPreviousKeyIdsValue := client.RsaAlgKeyId{}
			rsaAlgorithmPreviousKeyIdsAttrs := rsaAlgorithmPreviousKeyIdsElement.(types.Object).Attributes()
			rsaAlgorithmPreviousKeyIdsValue.KeyId = rsaAlgorithmPreviousKeyIdsAttrs["key_id"].(types.String).ValueString()
			rsaAlgorithmPreviousKeyIdsValue.RsaAlgType = rsaAlgorithmPreviousKeyIdsAttrs["rsa_alg_type"].(types.String).ValueString()
			signingKeysValue.RsaAlgorithmPreviousKeyIds = append(signingKeysValue.RsaAlgorithmPreviousKeyIds, rsaAlgorithmPreviousKeyIdsValue)
		}
	}
	if !signingKeysAttrs["rsa_previous_cert_ref"].IsNull() {
		signingKeysRsaPreviousCertRefValue := &client.ResourceLink{}
		signingKeysRsaPreviousCertRefAttrs := signingKeysAttrs["rsa_previous_cert_ref"].(types.Object).Attributes()
		signingKeysRsaPreviousCertRefValue.Id = signingKeysRsaPreviousCertRefAttrs["id"].(types.String).ValueString()
		signingKeysValue.RsaPreviousCertRef = signingKeysRsaPreviousCertRefValue
	}
	signingKeysValue.RsaPreviousKeyId = signingKeysAttrs["rsa_previous_key_id"].(types.String).ValueStringPointer()
	signingKeysValue.RsaPublishX5cParameter = signingKeysAttrs["rsa_publish_x5c_parameter"].(types.Bool).ValueBoolPointer()
	result.SigningKeys = signingKeysValue

	return result
}

func (state *keypairsOauthOpenidConnectAdditionalKeySetResourceModel) readClientResponse(response *client.AdditionalKeySet, versionAtLeast1201 bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// description
	state.Description = types.StringPointerValue(response.Description)
	// issuers
	issuersAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	issuersElementType := types.ObjectType{AttrTypes: issuersAttrTypes}
	var issuersValues []attr.Value
	for _, issuersResponseValue := range response.Issuers {
		issuersValue, diags := types.ObjectValue(issuersAttrTypes, map[string]attr.Value{
			"id": types.StringValue(issuersResponseValue.Id),
		})
		respDiags.Append(diags...)
		issuersValues = append(issuersValues, issuersValue)
	}
	issuersValue, diags := types.ListValue(issuersElementType, issuersValues)
	respDiags.Append(diags...)

	state.Issuers = issuersValue
	// name
	state.Name = types.StringValue(response.Name)
	// set_id
	state.SetId = types.StringPointerValue(response.Id)
	// signing_keys
	signingKeysP256ActiveCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysP256PreviousCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysP384ActiveCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysP384PreviousCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysP521ActiveCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysP521PreviousCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysRsaActiveCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysRsaAlgorithmActiveKeyIdsAttrTypes := map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	signingKeysRsaAlgorithmActiveKeyIdsElementType := types.ObjectType{AttrTypes: signingKeysRsaAlgorithmActiveKeyIdsAttrTypes}
	signingKeysRsaAlgorithmPreviousKeyIdsAttrTypes := map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	signingKeysRsaAlgorithmPreviousKeyIdsElementType := types.ObjectType{AttrTypes: signingKeysRsaAlgorithmPreviousKeyIdsAttrTypes}
	signingKeysRsaPreviousCertRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysAttrTypes := map[string]attr.Type{
		"p256_active_cert_ref":           types.ObjectType{AttrTypes: signingKeysP256ActiveCertRefAttrTypes},
		"p256_active_key_id":             types.StringType,
		"p256_previous_cert_ref":         types.ObjectType{AttrTypes: signingKeysP256PreviousCertRefAttrTypes},
		"p256_previous_key_id":           types.StringType,
		"p256_publish_x5c_parameter":     types.BoolType,
		"p384_active_cert_ref":           types.ObjectType{AttrTypes: signingKeysP384ActiveCertRefAttrTypes},
		"p384_active_key_id":             types.StringType,
		"p384_previous_cert_ref":         types.ObjectType{AttrTypes: signingKeysP384PreviousCertRefAttrTypes},
		"p384_previous_key_id":           types.StringType,
		"p384_publish_x5c_parameter":     types.BoolType,
		"p521_active_cert_ref":           types.ObjectType{AttrTypes: signingKeysP521ActiveCertRefAttrTypes},
		"p521_active_key_id":             types.StringType,
		"p521_previous_cert_ref":         types.ObjectType{AttrTypes: signingKeysP521PreviousCertRefAttrTypes},
		"p521_previous_key_id":           types.StringType,
		"p521_publish_x5c_parameter":     types.BoolType,
		"rsa_active_cert_ref":            types.ObjectType{AttrTypes: signingKeysRsaActiveCertRefAttrTypes},
		"rsa_active_key_id":              types.StringType,
		"rsa_algorithm_active_key_ids":   types.ListType{ElemType: signingKeysRsaAlgorithmActiveKeyIdsElementType},
		"rsa_algorithm_previous_key_ids": types.ListType{ElemType: signingKeysRsaAlgorithmPreviousKeyIdsElementType},
		"rsa_previous_cert_ref":          types.ObjectType{AttrTypes: signingKeysRsaPreviousCertRefAttrTypes},
		"rsa_previous_key_id":            types.StringType,
		"rsa_publish_x5c_parameter":      types.BoolType,
	}
	var signingKeysP256ActiveCertRefValue types.Object
	if response.SigningKeys.P256ActiveCertRef == nil {
		signingKeysP256ActiveCertRefValue = types.ObjectNull(signingKeysP256ActiveCertRefAttrTypes)
	} else {
		signingKeysP256ActiveCertRefValue, diags = types.ObjectValue(signingKeysP256ActiveCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P256ActiveCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysP256PreviousCertRefValue types.Object
	if response.SigningKeys.P256PreviousCertRef == nil {
		signingKeysP256PreviousCertRefValue = types.ObjectNull(signingKeysP256PreviousCertRefAttrTypes)
	} else {
		signingKeysP256PreviousCertRefValue, diags = types.ObjectValue(signingKeysP256PreviousCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P256PreviousCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysP384ActiveCertRefValue types.Object
	if response.SigningKeys.P384ActiveCertRef == nil {
		signingKeysP384ActiveCertRefValue = types.ObjectNull(signingKeysP384ActiveCertRefAttrTypes)
	} else {
		signingKeysP384ActiveCertRefValue, diags = types.ObjectValue(signingKeysP384ActiveCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P384ActiveCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysP384PreviousCertRefValue types.Object
	if response.SigningKeys.P384PreviousCertRef == nil {
		signingKeysP384PreviousCertRefValue = types.ObjectNull(signingKeysP384PreviousCertRefAttrTypes)
	} else {
		signingKeysP384PreviousCertRefValue, diags = types.ObjectValue(signingKeysP384PreviousCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P384PreviousCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysP521ActiveCertRefValue types.Object
	if response.SigningKeys.P521ActiveCertRef == nil {
		signingKeysP521ActiveCertRefValue = types.ObjectNull(signingKeysP521ActiveCertRefAttrTypes)
	} else {
		signingKeysP521ActiveCertRefValue, diags = types.ObjectValue(signingKeysP521ActiveCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P521ActiveCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysP521PreviousCertRefValue types.Object
	if response.SigningKeys.P521PreviousCertRef == nil {
		signingKeysP521PreviousCertRefValue = types.ObjectNull(signingKeysP521PreviousCertRefAttrTypes)
	} else {
		signingKeysP521PreviousCertRefValue, diags = types.ObjectValue(signingKeysP521PreviousCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.P521PreviousCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysRsaActiveCertRefValue types.Object
	if response.SigningKeys.RsaActiveCertRef == nil {
		signingKeysRsaActiveCertRefValue = types.ObjectNull(signingKeysRsaActiveCertRefAttrTypes)
	} else {
		signingKeysRsaActiveCertRefValue, diags = types.ObjectValue(signingKeysRsaActiveCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.RsaActiveCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	var signingKeysRsaAlgorithmActiveKeyIdsValue types.List
	if versionAtLeast1201 {
		var signingKeysRsaAlgorithmActiveKeyIdsValues []attr.Value
		for _, signingKeysRsaAlgorithmActiveKeyIdsResponseValue := range response.SigningKeys.RsaAlgorithmActiveKeyIds {
			signingKeysRsaAlgorithmActiveKeyIdsValue, diags := types.ObjectValue(signingKeysRsaAlgorithmActiveKeyIdsAttrTypes, map[string]attr.Value{
				"key_id":       types.StringValue(signingKeysRsaAlgorithmActiveKeyIdsResponseValue.KeyId),
				"rsa_alg_type": types.StringValue(signingKeysRsaAlgorithmActiveKeyIdsResponseValue.RsaAlgType),
			})
			respDiags.Append(diags...)
			signingKeysRsaAlgorithmActiveKeyIdsValues = append(signingKeysRsaAlgorithmActiveKeyIdsValues, signingKeysRsaAlgorithmActiveKeyIdsValue)
		}
		signingKeysRsaAlgorithmActiveKeyIdsValue, diags = types.ListValue(signingKeysRsaAlgorithmActiveKeyIdsElementType, signingKeysRsaAlgorithmActiveKeyIdsValues)
		respDiags.Append(diags...)
	} else {
		signingKeysRsaAlgorithmActiveKeyIdsValue = types.ListNull(signingKeysRsaAlgorithmActiveKeyIdsElementType)
	}
	var signingKeysRsaAlgorithmPreviousKeyIdsValue types.List
	if versionAtLeast1201 {
		var signingKeysRsaAlgorithmPreviousKeyIdsValues []attr.Value
		for _, signingKeysRsaAlgorithmPreviousKeyIdsResponseValue := range response.SigningKeys.RsaAlgorithmPreviousKeyIds {
			signingKeysRsaAlgorithmPreviousKeyIdsValue, diags := types.ObjectValue(signingKeysRsaAlgorithmPreviousKeyIdsAttrTypes, map[string]attr.Value{
				"key_id":       types.StringValue(signingKeysRsaAlgorithmPreviousKeyIdsResponseValue.KeyId),
				"rsa_alg_type": types.StringValue(signingKeysRsaAlgorithmPreviousKeyIdsResponseValue.RsaAlgType),
			})
			respDiags.Append(diags...)
			signingKeysRsaAlgorithmPreviousKeyIdsValues = append(signingKeysRsaAlgorithmPreviousKeyIdsValues, signingKeysRsaAlgorithmPreviousKeyIdsValue)
		}
		signingKeysRsaAlgorithmPreviousKeyIdsValue, diags = types.ListValue(signingKeysRsaAlgorithmPreviousKeyIdsElementType, signingKeysRsaAlgorithmPreviousKeyIdsValues)
		respDiags.Append(diags...)
	} else {
		signingKeysRsaAlgorithmPreviousKeyIdsValue = types.ListNull(signingKeysRsaAlgorithmPreviousKeyIdsElementType)
	}
	var signingKeysRsaPreviousCertRefValue types.Object
	if response.SigningKeys.RsaPreviousCertRef == nil {
		signingKeysRsaPreviousCertRefValue = types.ObjectNull(signingKeysRsaPreviousCertRefAttrTypes)
	} else {
		signingKeysRsaPreviousCertRefValue, diags = types.ObjectValue(signingKeysRsaPreviousCertRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeys.RsaPreviousCertRef.Id),
		})
		respDiags.Append(diags...)
	}
	signingKeysValue, diags := types.ObjectValue(signingKeysAttrTypes, map[string]attr.Value{
		"p256_active_cert_ref":           signingKeysP256ActiveCertRefValue,
		"p256_active_key_id":             types.StringPointerValue(response.SigningKeys.P256ActiveKeyId),
		"p256_previous_cert_ref":         signingKeysP256PreviousCertRefValue,
		"p256_previous_key_id":           types.StringPointerValue(response.SigningKeys.P256PreviousKeyId),
		"p256_publish_x5c_parameter":     types.BoolPointerValue(response.SigningKeys.P256PublishX5cParameter),
		"p384_active_cert_ref":           signingKeysP384ActiveCertRefValue,
		"p384_active_key_id":             types.StringPointerValue(response.SigningKeys.P384ActiveKeyId),
		"p384_previous_cert_ref":         signingKeysP384PreviousCertRefValue,
		"p384_previous_key_id":           types.StringPointerValue(response.SigningKeys.P384PreviousKeyId),
		"p384_publish_x5c_parameter":     types.BoolPointerValue(response.SigningKeys.P384PublishX5cParameter),
		"p521_active_cert_ref":           signingKeysP521ActiveCertRefValue,
		"p521_active_key_id":             types.StringPointerValue(response.SigningKeys.P521ActiveKeyId),
		"p521_previous_cert_ref":         signingKeysP521PreviousCertRefValue,
		"p521_previous_key_id":           types.StringPointerValue(response.SigningKeys.P521PreviousKeyId),
		"p521_publish_x5c_parameter":     types.BoolPointerValue(response.SigningKeys.P521PublishX5cParameter),
		"rsa_active_cert_ref":            signingKeysRsaActiveCertRefValue,
		"rsa_active_key_id":              types.StringPointerValue(response.SigningKeys.RsaActiveKeyId),
		"rsa_algorithm_active_key_ids":   signingKeysRsaAlgorithmActiveKeyIdsValue,
		"rsa_algorithm_previous_key_ids": signingKeysRsaAlgorithmPreviousKeyIdsValue,
		"rsa_previous_cert_ref":          signingKeysRsaPreviousCertRefValue,
		"rsa_previous_key_id":            types.StringPointerValue(response.SigningKeys.RsaPreviousKeyId),
		"rsa_publish_x5c_parameter":      types.BoolPointerValue(response.SigningKeys.RsaPublishX5cParameter),
	})
	respDiags.Append(diags...)

	state.SigningKeys = signingKeysValue
	return respDiags
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsOauthOpenidConnectAdditionalKeySetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1201)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1201 := compare >= 0
	clientData := data.buildClientStruct(pfVersionAtLeast1201)
	apiCreateRequest := r.apiClient.KeyPairsOauthOpenIdConnectAPI.CreateKeySet(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KeyPairsOauthOpenIdConnectAPI.CreateKeySetExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the keypairsOauthOpenidConnectAdditionalKeySet", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, pfVersionAtLeast1201)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keypairsOauthOpenidConnectAdditionalKeySetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KeyPairsOauthOpenIdConnectAPI.GetKeySet(config.AuthContext(ctx, r.providerConfig), data.SetId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the keypairsOauthOpenidConnectAdditionalKeySet", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the keypairsOauthOpenidConnectAdditionalKeySet", err, httpResp)
		}
		return
	}

	// Read response into the model
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1201)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1201 := compare >= 0
	resp.Diagnostics.Append(data.readClientResponse(responseData, pfVersionAtLeast1201)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data keypairsOauthOpenidConnectAdditionalKeySetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1201)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast1201 := compare >= 0
	clientData := data.buildClientStruct(pfVersionAtLeast1201)
	apiUpdateRequest := r.apiClient.KeyPairsOauthOpenIdConnectAPI.UpdateKeySet(config.AuthContext(ctx, r.providerConfig), data.SetId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KeyPairsOauthOpenIdConnectAPI.UpdateKeySetExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the keypairsOauthOpenidConnectAdditionalKeySet", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, pfVersionAtLeast1201)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data keypairsOauthOpenidConnectAdditionalKeySetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.KeyPairsOauthOpenIdConnectAPI.DeleteKeySet(config.AuthContext(ctx, r.providerConfig), data.SetId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the keypairsOauthOpenidConnectAdditionalKeySet", err, httpResp)
	}
}

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to set_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("set_id"), req, resp)
}
