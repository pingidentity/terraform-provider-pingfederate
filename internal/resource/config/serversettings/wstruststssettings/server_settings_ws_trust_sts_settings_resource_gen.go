// Code generated by ping-terraform-plugin-framework-generator

package serversettingswstruststssettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &serverSettingsWsTrustStsSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsWsTrustStsSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsWsTrustStsSettingsResource{}
)

func ServerSettingsWsTrustStsSettingsResource() resource.Resource {
	return &serverSettingsWsTrustStsSettingsResource{}
}

type serverSettingsWsTrustStsSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *serverSettingsWsTrustStsSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_ws_trust_sts_settings"
}

func (r *serverSettingsWsTrustStsSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type serverSettingsWsTrustStsSettingsResourceModel struct {
	BasicAuthnEnabled      types.Bool `tfsdk:"basic_authn_enabled"`
	ClientCertAuthnEnabled types.Bool `tfsdk:"client_cert_authn_enabled"`
	IssuerCerts            types.Set  `tfsdk:"issuer_certs"`
	RestrictByIssuerCert   types.Bool `tfsdk:"restrict_by_issuer_cert"`
	RestrictBySubjectDn    types.Bool `tfsdk:"restrict_by_subject_dn"`
	SubjectDns             types.Set  `tfsdk:"subject_dns"`
	Users                  types.Set  `tfsdk:"users"`
}

func (r *serverSettingsWsTrustStsSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the WS-Trust STS settings.",
		Attributes: map[string]schema.Attribute{
			"basic_authn_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Require the use of HTTP Basic Authentication to access WS-Trust STS endpoints. Requires users be populated. Default value is `false`.",
			},
			"client_cert_authn_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Require the use of Client Cert Authentication to access WS-Trust STS endpoints. Requires either `restrict_by_subject_dn` and/or `restrict_by_issuer_cert` be `true`. Default value is `false`.",
			},
			"issuer_certs": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the resource.",
						},
					},
				},
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(resourceLinkSetDefault),
				Description: "List of certificate Issuers that are used to validate certificates for access to the WS-Trust STS endpoints. Required if `restrict_by_issuer_cert` is `true`.",
			},
			"restrict_by_issuer_cert": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Restrict Access by Issuer Certificate. Ignored if `client_cert_authn_enabled` is `false`. Default value is `false`.",
			},
			"restrict_by_subject_dn": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Restrict Access by Subject DN. Ignored if `client_cert_authn_enabled` is `false`. Default value is `false`.",
			},
			"subject_dns": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(emptyStringSetDefault),
				Description: "List of Subject DNs for certificates that are allowed to authenticate to WS-Trust STS endpoints. Required if `restrict_by_subject_dn` is `true`.",
			},
			"users": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"password": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "User password. Either `password` or `encrypted_password` is required.",
						},
						"encrypted_password": schema.StringAttribute{
							Description: "Encrypted user password. Either `password` or `encrypted_password` is required.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("password")),
							},
						},
						"username": schema.StringAttribute{
							Optional:    true,
							Description: "The username.",
						},
					},
				},
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(usersSetDefault),
				Description: "List of users authorized to access WS-Trust STS endpoints when `basic_auth_enabled` is `true`. At least one users entry is required if `basic_auth_enabled` is `true`.",
			},
		},
	}
}

func (model *serverSettingsWsTrustStsSettingsResourceModel) buildClientStruct() (*client.WsTrustStsSettings, diag.Diagnostics) {
	result := &client.WsTrustStsSettings{}
	// basic_authn_enabled
	result.BasicAuthnEnabled = model.BasicAuthnEnabled.ValueBoolPointer()
	// client_cert_authn_enabled
	result.ClientCertAuthnEnabled = model.ClientCertAuthnEnabled.ValueBoolPointer()
	// issuer_certs
	result.IssuerCerts = []client.ResourceLink{}
	for _, issuerCertsElement := range model.IssuerCerts.Elements() {
		issuerCertsValue := client.ResourceLink{}
		issuerCertsAttrs := issuerCertsElement.(types.Object).Attributes()
		issuerCertsValue.Id = issuerCertsAttrs["id"].(types.String).ValueString()
		result.IssuerCerts = append(result.IssuerCerts, issuerCertsValue)
	}

	// restrict_by_issuer_cert
	result.RestrictByIssuerCert = model.RestrictByIssuerCert.ValueBoolPointer()
	// restrict_by_subject_dn
	result.RestrictBySubjectDn = model.RestrictBySubjectDn.ValueBoolPointer()
	// subject_dns
	if !model.SubjectDns.IsNull() {
		result.SubjectDns = []string{}
		for _, subjectDnsElement := range model.SubjectDns.Elements() {
			result.SubjectDns = append(result.SubjectDns, subjectDnsElement.(types.String).ValueString())
		}
	}

	// users
	result.Users = []client.UsernamePasswordCredentials{}
	for _, usersElement := range model.Users.Elements() {
		usersValue := client.UsernamePasswordCredentials{}
		usersAttrs := usersElement.(types.Object).Attributes()
		usersValue.Password = usersAttrs["password"].(types.String).ValueStringPointer()
		usersValue.EncryptedPassword = usersAttrs["encrypted_password"].(types.String).ValueStringPointer()
		usersValue.Username = usersAttrs["username"].(types.String).ValueStringPointer()
		result.Users = append(result.Users, usersValue)
	}

	return result, nil
}

func (state *serverSettingsWsTrustStsSettingsResourceModel) readClientResponse(response *client.WsTrustStsSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// basic_authn_enabled
	state.BasicAuthnEnabled = types.BoolPointerValue(response.BasicAuthnEnabled)
	// client_cert_authn_enabled
	state.ClientCertAuthnEnabled = types.BoolPointerValue(response.ClientCertAuthnEnabled)
	// issuer_certs
	issuerCertsAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	issuerCertsElementType := types.ObjectType{AttrTypes: issuerCertsAttrTypes}
	var issuerCertsValues []attr.Value
	for _, issuerCertsResponseValue := range response.IssuerCerts {
		issuerCertsValue, diags := types.ObjectValue(issuerCertsAttrTypes, map[string]attr.Value{
			"id": types.StringValue(issuerCertsResponseValue.Id),
		})
		respDiags.Append(diags...)
		issuerCertsValues = append(issuerCertsValues, issuerCertsValue)
	}
	issuerCertsValue, diags := types.SetValue(issuerCertsElementType, issuerCertsValues)
	respDiags.Append(diags...)

	state.IssuerCerts = issuerCertsValue
	// restrict_by_issuer_cert
	state.RestrictByIssuerCert = types.BoolPointerValue(response.RestrictByIssuerCert)
	// restrict_by_subject_dn
	state.RestrictBySubjectDn = types.BoolPointerValue(response.RestrictBySubjectDn)
	// subject_dns
	state.SubjectDns, diags = types.SetValueFrom(context.Background(), types.StringType, response.SubjectDns)
	respDiags.Append(diags...)
	// users
	respDiags.Append(state.readClientResponseUsers(response)...)
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *serverSettingsWsTrustStsSettingsResource) emptyModel() serverSettingsWsTrustStsSettingsResourceModel {
	var model serverSettingsWsTrustStsSettingsResourceModel
	// issuer_certs
	issuerCertsAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	issuerCertsElementType := types.ObjectType{AttrTypes: issuerCertsAttrTypes}
	model.IssuerCerts = types.SetNull(issuerCertsElementType)
	// users
	usersAttrTypes := map[string]attr.Type{
		"password":           types.StringType,
		"username":           types.StringType,
		"encrypted_password": types.StringType,
	}
	usersElementType := types.ObjectType{AttrTypes: usersAttrTypes}
	model.Users = types.SetNull(usersElementType)
	// subject_dns
	model.SubjectDns = types.SetNull(types.StringType)
	return model
}

func (r *serverSettingsWsTrustStsSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverSettingsWsTrustStsSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the serverSettingsWsTrustStsSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsWsTrustStsSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverSettingsWsTrustStsSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.GetWsTrustStsSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Server Settings WS-Trust STS Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the serverSettingsWsTrustStsSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsWsTrustStsSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data serverSettingsWsTrustStsSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the serverSettingsWsTrustStsSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverSettingsWsTrustStsSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
