package serversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsDataSource{}
)

// ServerSettingsDataSource is a helper function to simplify the provider implementation.
func ServerSettingsDataSource() datasource.DataSource {
	return &serverSettingsDataSource{}
}

// serverSettingsDataSource is the datasource implementation.
type serverSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the global server configuration settings",
		Attributes: map[string]schema.Attribute{
			"contact_info": schema.SingleNestedAttribute{
				Description: "Information that identifies the server.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(contactInfoDefault),
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Description: "Company name.",
						Optional:    true,
					},
					"email": schema.StringAttribute{
						Description: "Contact email address.",
						Optional:    true,
						Validators: []validator.String{
							configvalidators.ValidEmail(),
						},
					},
					"first_name": schema.StringAttribute{
						Description: "Contact first name.",
						Optional:    true,
					},
					"last_name": schema.StringAttribute{
						Description: "Contact last name.",
						Optional:    true,
					},
					"phone": schema.StringAttribute{
						Description: "Contact phone number.",
						Optional:    true,
					},
				},
			},
			"notifications": schema.SingleNestedAttribute{
				Description: "Notification settings for license and certificate expiration events.",
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(notificationsDefault),
				Attributes: map[string]schema.Attribute{
					"license_events": schema.SingleNestedAttribute{
						Description: "Settings for license event notifications.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
								},
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								Attributes:  resourcelink.ToSchema(),
							},
						},
					},
					"certificate_expirations": schema.SingleNestedAttribute{
						Description: "Settings for license event notifications.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
								},
							},
							"initial_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when initial warning is sent (in days).",
								Optional:    true,
							},
							"final_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when final warning is sent (in days).",
								Required:    true,
								Validators: []validator.Int64{
									// final_warning_period must be between 1 and 99999 days, inclusive
									int64validator.Between(1, 99999),
								},
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								Attributes:  resourcelink.ToSchema(),
							},
						},
					},
					"notify_admin_user_password_changes": schema.BoolAttribute{
						Description: "Determines whether admin users are notified through email when their account is changed.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"account_changes_notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher for admin user account changes.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"metadata_notification_settings": schema.SingleNestedAttribute{
						Description: "Settings for metadata update event notifications.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
								},
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								Attributes:  resourcelink.ToSchema(),
							},
						},
					},
				},
			},
			"roles_and_protocols": schema.SingleNestedAttribute{
				Description: "Configure roles and protocols.",
				Computed:    true,
				Optional:    false,
				Default:     objectdefault.StaticValue(rolesAndProtocolsDefault),
				Attributes: map[string]schema.Attribute{
					"oauth_role": schema.SingleNestedAttribute{
						Description: "OAuth role settings.",
						Computed:    true,
						Optional:    false,
						Default:     objectdefault.StaticValue(oauthRoleDefault),
						Attributes: map[string]schema.Attribute{
							"enable_oauth": schema.BoolAttribute{
								Description: "Enable OAuth 2.0 Authorization Server (AS) Role.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable Open ID Connect.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
					"idp_role": schema.SingleNestedAttribute{
						Description: "Identity Provider (IdP) settings.",
						Computed:    true,
						Optional:    false,
						Default:     objectdefault.StaticValue(idpRoleDefault),
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Identity Provider Role.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								Default:     objectdefault.StaticValue(idpSamlProfileDefault),
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
									"enable_auto_connect": schema.BoolAttribute{
										Description: "This property has been deprecated and no longer used.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
								},
							},
							"enable_outbound_provisioning": schema.BoolAttribute{
								Description: "Enable Outbound Provisioning.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
					"sp_role": schema.SingleNestedAttribute{
						Description: "Service Provider (SP) settings.",
						Computed:    true,
						Optional:    false,
						Default:     objectdefault.StaticValue(spRoleDefault),
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Service Provider Role.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								Default:     objectdefault.StaticValue(spSamlProfileDefault),
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
									"enable_auto_connect": schema.BoolAttribute{
										Description: "This property has been deprecated and no longer used.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
									"enable_xasp": schema.BoolAttribute{
										Description: "Enable Attribute Requester Mapping for X.509 Attribute Sharing Profile (XASP)",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
								},
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable OpenID Connect.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_inbound_provisioning": schema.BoolAttribute{
								Description: "Enable Inbound Provisioning.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
					"enable_idp_discovery": schema.BoolAttribute{
						Description: "Enable IdP Discovery.",
						Computed:    true,
						Optional:    false,
						Default:     booldefault.StaticBool(true),
					},
				},
			},
			"federation_info": schema.SingleNestedAttribute{
				Description: "Federation Info.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Description: "The fully qualified host name, port, and path (if applicable) on which the PingFederate server runs.",
						Required:    true,
						Validators: []validator.String{
							configvalidators.ValidUrl(),
						},
					},
					"saml_2_entity_id": schema.StringAttribute{
						Description: "This ID defines your organization as the entity operating the server for SAML 2.0 transactions. It is usually defined as an organization's URL or a DNS address; for example: pingidentity.com. The SAML SourceID used for artifact resolution is derived from this ID using SHA1.",
						Required:    true,
					},
					"auto_connect_entity_id": schema.StringAttribute{
						Description: "This property has been deprecated and no longer used",
						Computed:    true,
						Optional:    false,
						Default:     stringdefault.StaticString(""),
					},
					"saml_1x_issuer_id": schema.StringAttribute{
						Description: "This ID identifies your federation server for SAML 1.x transactions. As with SAML 2.0, it is usually defined as an organization's URL or a DNS address. The SourceID used for artifact resolution is derived from this ID using SHA1.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
					"saml_1x_source_id": schema.StringAttribute{
						Description: "If supplied, the Source ID value entered here is used for SAML 1.x, instead of being derived from the SAML 1.x Issuer/Audience.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
					"wsfed_realm": schema.StringAttribute{
						Description: "The URI of the realm associated with the PingFederate server. A realm represents a single unit of security administration or trust.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
				},
			},
			"email_server": schema.SingleNestedAttribute{
				Description: "Email Server Settings.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(types.ObjectNull(emailServerAttrType)),
				Attributes: map[string]schema.Attribute{
					"source_addr": schema.StringAttribute{
						Description: "The email address that appears in the 'From' header line in email messages generated by PingFederate. The address must be in valid format but need not be set up on your system.",
						Required:    true,
						Validators: []validator.String{
							configvalidators.ValidEmail(),
						},
					},
					"email_server": schema.StringAttribute{
						Description: "The IP address or hostname of your email server.",
						Required:    true,
						Validators: []validator.String{
							configvalidators.ValidHostnameOrIp(),
						},
					},
					"port": schema.Int64Attribute{
						Description: "The SMTP port on your email server. Allowable values: 1 - 65535. The default value is 25.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(25),
					},
					"ssl_port": schema.Int64Attribute{
						Description: "The secure SMTP port on your email server. This field is not active unless Use SSL is enabled. Allowable values: 1 - 65535. The default value is 465.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(465),
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"timeout": schema.Int64Attribute{
						Description: "The amount of time in seconds that PingFederate will wait before it times out connecting to the SMTP server. Allowable values: 0 - 3600. The default value is 30.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(30),
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"retry_attempts": schema.Int64Attribute{
						Description: "The number of times PingFederate tries to resend an email upon unsuccessful delivery. The default value is 2.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(2),
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"retry_delay": schema.Int64Attribute{
						Description: "The number of minutes PingFederate waits before the next retry attempt. The default value is 2.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(2),
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"use_ssl": schema.BoolAttribute{
						Description: "Requires the use of SSL/TLS on the port specified by 'sslPort'. If this option is enabled, it overrides the 'useTLS' option.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"use_tls": schema.BoolAttribute{
						Description: "Requires the use of the STARTTLS protocol on the port specified by 'port'.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"verify_hostname": schema.BoolAttribute{
						Description: "If useSSL or useTLS is enabled, this flag determines whether the email server hostname is verified against the server's SMTPS certificate.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"enable_utf8_message_headers": schema.BoolAttribute{
						Description: "Only set this flag to true if the email server supports UTF-8 characters in message headers. Otherwise, this is defaulted to false.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"use_debugging": schema.BoolAttribute{
						Description: "Turns on detailed error messages for the PingFederate server log to help troubleshoot any problems.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"username": schema.StringAttribute{
						Description: "Authorized email username. Required if the password is provided.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Description: "User password. To update the password, specify the plaintext value in this field. This field will not be populated for GET requests.",
						Computed:    true,
						Optional:    true,
						Sensitive:   true,
						Default:     stringdefault.StaticString(""),
					},
				},
			},
			"captcha_settings": schema.SingleNestedAttribute{
				Description: "Captcha Settings.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(types.ObjectNull(captchaSettingsAttrType)),
				Attributes: map[string]schema.Attribute{
					"site_key": schema.StringAttribute{
						Description: "Site key for reCAPTCHA.",
						Required:    true,
					},
					"secret_key": schema.StringAttribute{
						Description: "Secret key for reCAPTCHA. GETs will not return this attribute. To update this field, specify the new value in this attribute.",
						Required:    true,
					},
				},
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *serverSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings"
}

func (r *serverSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

// Read the server settings datasource from the PingFederate API and update the state accordingly.
func (r *serverSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
	}

	diags = readServerSettingsResponse(ctx, apiReadServerSettings, &state, &state, "server_settings_id")
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
