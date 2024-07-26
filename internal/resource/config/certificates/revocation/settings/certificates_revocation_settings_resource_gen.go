// Code generated by ping-terraform-plugin-framework-generator

package certificatesrevocationsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &certificatesRevocationSettingsResource{}
	_ resource.ResourceWithConfigure   = &certificatesRevocationSettingsResource{}
	_ resource.ResourceWithImportState = &certificatesRevocationSettingsResource{}
)

func CertificatesRevocationSettingsResource() resource.Resource {
	return &certificatesRevocationSettingsResource{}
}

type certificatesRevocationSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *certificatesRevocationSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates_revocation_settings"
}

func (r *certificatesRevocationSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificatesRevocationSettingsResourceModel struct {
	CrlSettings   types.Object `tfsdk:"crl_settings"`
	OcspSettings  types.Object `tfsdk:"ocsp_settings"`
	ProxySettings types.Object `tfsdk:"proxy_settings"`
}

func (r *certificatesRevocationSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the certificate revocation settings.",
		Attributes: map[string]schema.Attribute{
			"crl_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"next_retry_mins_when_next_update_in_past": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(60),
						Description: "Next retry on next update expiration in minutes. This value defaults to `60`.",
					},
					"next_retry_mins_when_resolve_failed": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(1440),
						Description: "Next retry on resolution failure in minutes. This value defaults to `1440`.",
					},
					"treat_non_retrievable_crl_as_revoked": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Treat non retrievable CRL as revoked. This setting defaults to `false`.",
					},
					"verify_crl_signature": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
						Description: "Verify CRL signature. This setting defaults to `true`.",
					},
				},
				Optional:    true,
				Description: "Certificate revocation CRL settings. If this attribute is omitted, CRL checks are disabled.",
			},
			"ocsp_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"action_on_responder_unavailable": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("CONTINUE"),
						Description: "Action on responder unavailable. This value defaults to `CONTINUE`. Options are `CONTINUE`, `FAIL`, `FAILOVER`.",
						Validators: []validator.String{
							stringvalidator.OneOf("CONTINUE", "FAIL", "FAILOVER"),
						},
					},
					"action_on_status_unknown": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("FAIL"),
						Description: "Action on status unknown. This value defaults to `FAIL`. Options are `CONTINUE`, `FAIL`, `FAILOVER`.",
						Validators: []validator.String{
							stringvalidator.OneOf("CONTINUE", "FAIL", "FAILOVER"),
						},
					},
					"action_on_unsuccessful_response": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("FAIL"),
						Description: "Action on unsuccessful response. This value defaults to `FAIL`. Options are `CONTINUE`, `FAIL`, `FAILOVER`.",
						Validators: []validator.String{
							stringvalidator.OneOf("CONTINUE", "FAIL", "FAILOVER"),
						},
					},
					"current_update_grace_period": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: "Current update grace period in minutes. This value defaults to `5`.",
					},
					"next_update_grace_period": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: "Next update grace period in minutes. This value defaults to `5`.",
					},
					"requester_add_nonce": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Do not allow responder to use cached responses. This setting defaults to `false`.",
					},
					"responder_cert_reference": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Optional:    true,
						Description: "Resource link to OCSP responder signature verification certificate. A previously selected certificate will be deselected if this attribute is not defined.",
					},
					"responder_timeout": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: "Responder connection timeout in seconds. This value defaults to `5`.",
					},
					"responder_url": schema.StringAttribute{
						Optional:    true,
						Description: "Default responder URL. This URL is used if the certificate being checked does not specify an OCSP responder URL.",
					},
					"response_cache_period": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(48),
						Description: "Response cache period in hours. This value defaults to `48`.",
					},
				},
				Optional:    true,
				Description: "Certificate revocation OCSP settings. If this attribute is omitted, OCSP checks are disabled.",
			},
			"proxy_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Required:    true,
						Description: "Host name.",
					},
					"port": schema.Int64Attribute{
						Required:    true,
						Description: "Port number.",
					},
				},
				Optional:    true,
				Description: "If OCSP messaging is routed through a proxy server, specify the server's host (DNS name or IP address) and the port number. The same proxy information applies to CRL checking, when CRL is enabled for failover.",
			},
		},
	}
}

func (model *certificatesRevocationSettingsResourceModel) buildClientStruct() (*client.CertificateRevocationSettings, diag.Diagnostics) {
	result := &client.CertificateRevocationSettings{}
	// crl_settings
	if !model.CrlSettings.IsNull() {
		crlSettingsValue := &client.CrlSettings{}
		crlSettingsAttrs := model.CrlSettings.Attributes()
		crlSettingsValue.NextRetryMinsWhenNextUpdateInPast = crlSettingsAttrs["next_retry_mins_when_next_update_in_past"].(types.Int64).ValueInt64Pointer()
		crlSettingsValue.NextRetryMinsWhenResolveFailed = crlSettingsAttrs["next_retry_mins_when_resolve_failed"].(types.Int64).ValueInt64Pointer()
		crlSettingsValue.TreatNonRetrievableCrlAsRevoked = crlSettingsAttrs["treat_non_retrievable_crl_as_revoked"].(types.Bool).ValueBoolPointer()
		crlSettingsValue.VerifyCrlSignature = crlSettingsAttrs["verify_crl_signature"].(types.Bool).ValueBoolPointer()
		result.CrlSettings = crlSettingsValue
	}

	// ocsp_settings
	if !model.OcspSettings.IsNull() {
		ocspSettingsValue := &client.OcspSettings{}
		ocspSettingsAttrs := model.OcspSettings.Attributes()
		ocspSettingsValue.ActionOnResponderUnavailable = ocspSettingsAttrs["action_on_responder_unavailable"].(types.String).ValueStringPointer()
		ocspSettingsValue.ActionOnStatusUnknown = ocspSettingsAttrs["action_on_status_unknown"].(types.String).ValueStringPointer()
		ocspSettingsValue.ActionOnUnsuccessfulResponse = ocspSettingsAttrs["action_on_unsuccessful_response"].(types.String).ValueStringPointer()
		ocspSettingsValue.CurrentUpdateGracePeriod = ocspSettingsAttrs["current_update_grace_period"].(types.Int64).ValueInt64Pointer()
		ocspSettingsValue.NextUpdateGracePeriod = ocspSettingsAttrs["next_update_grace_period"].(types.Int64).ValueInt64Pointer()
		ocspSettingsValue.RequesterAddNonce = ocspSettingsAttrs["requester_add_nonce"].(types.Bool).ValueBoolPointer()
		if !ocspSettingsAttrs["responder_cert_reference"].IsNull() {
			ocspSettingsResponderCertReferenceValue := &client.ResourceLink{}
			ocspSettingsResponderCertReferenceAttrs := ocspSettingsAttrs["responder_cert_reference"].(types.Object).Attributes()
			ocspSettingsResponderCertReferenceValue.Id = ocspSettingsResponderCertReferenceAttrs["id"].(types.String).ValueString()
			ocspSettingsValue.ResponderCertReference = ocspSettingsResponderCertReferenceValue
		}
		ocspSettingsValue.ResponderTimeout = ocspSettingsAttrs["responder_timeout"].(types.Int64).ValueInt64Pointer()
		ocspSettingsValue.ResponderUrl = ocspSettingsAttrs["responder_url"].(types.String).ValueStringPointer()
		ocspSettingsValue.ResponseCachePeriod = ocspSettingsAttrs["response_cache_period"].(types.Int64).ValueInt64Pointer()
		result.OcspSettings = ocspSettingsValue
	}

	// proxy_settings
	if !model.ProxySettings.IsNull() {
		proxySettingsValue := &client.ProxySettings{}
		proxySettingsAttrs := model.ProxySettings.Attributes()
		proxySettingsValue.Host = proxySettingsAttrs["host"].(types.String).ValueStringPointer()
		proxySettingsValue.Port = proxySettingsAttrs["port"].(types.Int64).ValueInt64Pointer()
		result.ProxySettings = proxySettingsValue
	}

	return result, nil
}

func (state *certificatesRevocationSettingsResourceModel) readClientResponse(response *client.CertificateRevocationSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// crl_settings
	crlSettingsAttrTypes := map[string]attr.Type{
		"next_retry_mins_when_next_update_in_past": types.Int64Type,
		"next_retry_mins_when_resolve_failed":      types.Int64Type,
		"treat_non_retrievable_crl_as_revoked":     types.BoolType,
		"verify_crl_signature":                     types.BoolType,
	}
	var crlSettingsValue types.Object
	if response.CrlSettings == nil {
		crlSettingsValue = types.ObjectNull(crlSettingsAttrTypes)
	} else {
		crlSettingsValue, diags = types.ObjectValue(crlSettingsAttrTypes, map[string]attr.Value{
			"next_retry_mins_when_next_update_in_past": types.Int64PointerValue(response.CrlSettings.NextRetryMinsWhenNextUpdateInPast),
			"next_retry_mins_when_resolve_failed":      types.Int64PointerValue(response.CrlSettings.NextRetryMinsWhenResolveFailed),
			"treat_non_retrievable_crl_as_revoked":     types.BoolPointerValue(response.CrlSettings.TreatNonRetrievableCrlAsRevoked),
			"verify_crl_signature":                     types.BoolPointerValue(response.CrlSettings.VerifyCrlSignature),
		})
		respDiags.Append(diags...)
	}

	state.CrlSettings = crlSettingsValue
	// ocsp_settings
	ocspSettingsResponderCertReferenceAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	ocspSettingsAttrTypes := map[string]attr.Type{
		"action_on_responder_unavailable": types.StringType,
		"action_on_status_unknown":        types.StringType,
		"action_on_unsuccessful_response": types.StringType,
		"current_update_grace_period":     types.Int64Type,
		"next_update_grace_period":        types.Int64Type,
		"requester_add_nonce":             types.BoolType,
		"responder_cert_reference":        types.ObjectType{AttrTypes: ocspSettingsResponderCertReferenceAttrTypes},
		"responder_timeout":               types.Int64Type,
		"responder_url":                   types.StringType,
		"response_cache_period":           types.Int64Type,
	}
	var ocspSettingsValue types.Object
	if response.OcspSettings == nil {
		ocspSettingsValue = types.ObjectNull(ocspSettingsAttrTypes)
	} else {
		var ocspSettingsResponderCertReferenceValue types.Object
		if response.OcspSettings.ResponderCertReference == nil {
			ocspSettingsResponderCertReferenceValue = types.ObjectNull(ocspSettingsResponderCertReferenceAttrTypes)
		} else {
			ocspSettingsResponderCertReferenceValue, diags = types.ObjectValue(ocspSettingsResponderCertReferenceAttrTypes, map[string]attr.Value{
				"id": types.StringValue(response.OcspSettings.ResponderCertReference.Id),
			})
			respDiags.Append(diags...)
		}
		ocspSettingsValue, diags = types.ObjectValue(ocspSettingsAttrTypes, map[string]attr.Value{
			"action_on_responder_unavailable": types.StringPointerValue(response.OcspSettings.ActionOnResponderUnavailable),
			"action_on_status_unknown":        types.StringPointerValue(response.OcspSettings.ActionOnStatusUnknown),
			"action_on_unsuccessful_response": types.StringPointerValue(response.OcspSettings.ActionOnUnsuccessfulResponse),
			"current_update_grace_period":     types.Int64PointerValue(response.OcspSettings.CurrentUpdateGracePeriod),
			"next_update_grace_period":        types.Int64PointerValue(response.OcspSettings.NextUpdateGracePeriod),
			"requester_add_nonce":             types.BoolPointerValue(response.OcspSettings.RequesterAddNonce),
			"responder_cert_reference":        ocspSettingsResponderCertReferenceValue,
			"responder_timeout":               types.Int64PointerValue(response.OcspSettings.ResponderTimeout),
			"responder_url":                   types.StringPointerValue(response.OcspSettings.ResponderUrl),
			"response_cache_period":           types.Int64PointerValue(response.OcspSettings.ResponseCachePeriod),
		})
		respDiags.Append(diags...)
	}

	state.OcspSettings = ocspSettingsValue
	// proxy_settings
	proxySettingsAttrTypes := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}
	var proxySettingsValue types.Object
	if response.ProxySettings == nil {
		proxySettingsValue = types.ObjectNull(proxySettingsAttrTypes)
	} else {
		proxySettingsValue, diags = types.ObjectValue(proxySettingsAttrTypes, map[string]attr.Value{
			"host": types.StringPointerValue(response.ProxySettings.Host),
			"port": types.Int64PointerValue(response.ProxySettings.Port),
		})
		respDiags.Append(diags...)
	}

	state.ProxySettings = proxySettingsValue
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *certificatesRevocationSettingsResource) emptyModel() certificatesRevocationSettingsResourceModel {
	var model certificatesRevocationSettingsResourceModel
	// crl_settings
	crlSettingsAttrTypes := map[string]attr.Type{
		"next_retry_mins_when_next_update_in_past": types.Int64Type,
		"next_retry_mins_when_resolve_failed":      types.Int64Type,
		"treat_non_retrievable_crl_as_revoked":     types.BoolType,
		"verify_crl_signature":                     types.BoolType,
	}
	model.CrlSettings = types.ObjectNull(crlSettingsAttrTypes)
	// ocsp_settings
	ocspSettingsResponderCertReferenceAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	ocspSettingsAttrTypes := map[string]attr.Type{
		"action_on_responder_unavailable": types.StringType,
		"action_on_status_unknown":        types.StringType,
		"action_on_unsuccessful_response": types.StringType,
		"current_update_grace_period":     types.Int64Type,
		"next_update_grace_period":        types.Int64Type,
		"requester_add_nonce":             types.BoolType,
		"responder_cert_reference":        types.ObjectType{AttrTypes: ocspSettingsResponderCertReferenceAttrTypes},
		"responder_timeout":               types.Int64Type,
		"responder_url":                   types.StringType,
		"response_cache_period":           types.Int64Type,
	}
	model.OcspSettings = types.ObjectNull(ocspSettingsAttrTypes)
	// proxy_settings
	proxySettingsAttrTypes := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}
	model.ProxySettings = types.ObjectNull(proxySettingsAttrTypes)
	return model
}

func (r *certificatesRevocationSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data certificatesRevocationSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the certificatesRevocationSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *certificatesRevocationSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data certificatesRevocationSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.CertificatesRevocationAPI.GetRevocationSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the certificatesRevocationSettings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the certificatesRevocationSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *certificatesRevocationSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data certificatesRevocationSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.CertificatesRevocationAPI.UpdateRevocationSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the certificatesRevocationSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *certificatesRevocationSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}