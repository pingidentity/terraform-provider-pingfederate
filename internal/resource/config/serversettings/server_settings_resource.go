package serversettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsResource{}
)

// ServerSettingsResource is a helper function to simplify the provider implementation.
func ServerSettingsResource() resource.Resource {
	return &serverSettingsResource{}
}

// serverSettingsResource is the resource implementation.
type serverSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
						Description: "Notification settings for certificate expiration events.",
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
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to NOTIFICATION_PUBLISHER. Supported in PF version 11.3 or later.",
								Optional:    true,
								Computed:    true,
								// Default value is set in ModifyPlan below. When PF 11.3+ is all that is supported, the default can be moved to the schema here.
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NOTIFICATION_PUBLISHER",
										"LOGGING_ONLY",
									),
								},
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
					"expired_certificate_administrative_console_warning_days": schema.Int64Attribute{
						Description: "Indicates the number of days prior to certificate expiry date, the administrative console warning starts. The default value is 14 days. Supported in PF 12.0 or later.",
						Optional:    true,
						Computed:    true,
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
					},
					"expiring_certificate_administrative_console_warning_days": schema.Int64Attribute{
						Description: "Indicates the number of days past the certificate expiry date, the administrative console warning ends. The default value is 14 days. Supported in PF 12.0 or later.",
						Optional:    true,
						Computed:    true,
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
					},
					"thread_pool_exhaustion_notification_settings": schema.SingleNestedAttribute{
						Description: "Notification settings for thread pool exhaustion events. Supported in PF 12.0 or later.",
						Optional:    true,
						Computed:    true,
						Default:     objectdefault.StaticValue(types.ObjectNull(threadPoolExhaustionNotificationSettingsAttrType)),
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "Email address where notifications are sent.",
								Required:    true,
							},
							"thread_dump_enabled": schema.BoolAttribute{
								Description: "Generate a thread dump when approaching thread pool exhaustion.",
								Optional:    true,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								Attributes:  resourcelink.ToSchema(),
							},
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to LOGGING_ONLY.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("LOGGING_ONLY"),
								Validators: []validator.String{
									stringvalidator.OneOf(
										"NOTIFICATION_PUBLISHER",
										"LOGGING_ONLY",
									),
								},
							},
						},
					},
				},
			},
			"roles_and_protocols": schema.SingleNestedAttribute{
				Description: "Configure roles and protocols. As of PingFederate 12.0: This property has been deprecated and is no longer used. All Roles and protocols are always enabled.",
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

// ValidateConfig validates the configuration of the server settings resource.
// It also checks that the email_server use_ssl and use_tls attributes are not both set to true.
func (r *serverSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	var model serverSettingsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	// Validate email_server source_addr value
	if internaltypes.IsDefined(model.EmailServer) {
		esAttrs := model.EmailServer.Attributes()
		// If email_server attribute use_ssl is set, confirm that use_tls is NOT
		esUseSSLFlag := esAttrs["use_ssl"]
		esUseTLSFlag := esAttrs["use_tls"]
		if internaltypes.IsDefined(esUseSSLFlag) {
			esUseSSLFlagValue := esUseSSLFlag.(types.Bool).ValueBool()
			if esUseSSLFlagValue && internaltypes.IsDefined(esUseTLSFlag) {
				resp.Diagnostics.AddError("Overlapping settings!", "If the email server setting \"use_ssl\" is true, \"use_tls\" cannot be set. Remove one of the two values from your resource file.")
			}
		}
	}
}

func (r *serverSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to versions 11.3 and 12.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast113 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	var plan serverSettingsModel
	req.Plan.Get(ctx, &plan)
	if !internaltypes.IsDefined(plan.Notifications) {
		return
	}

	var diags diag.Diagnostics
	updatePlan := false
	planNotificationsAttrs := plan.Notifications.Attributes()
	planCertificateExpirations := planNotificationsAttrs["certificate_expirations"].(types.Object)
	if internaltypes.IsDefined(planCertificateExpirations) {
		planCertificateExpirationsAttrs := planCertificateExpirations.Attributes()
		planNotificationMode := planCertificateExpirationsAttrs["notification_mode"].(types.String)

		// If notification_mode is set and the PF version is not new enough, throw an error
		if !pfVersionAtLeast113 {
			if internaltypes.IsDefined(planNotificationMode) {
				version.AddUnsupportedAttributeError("notifications.certificate_expirations.notification_mode",
					r.providerConfig.ProductVersion, version.PingFederate1130, &resp.Diagnostics)
			} else if planNotificationMode.IsUnknown() {
				// Set a null default when the version isn't new enough to use this attribute
				planNotificationMode = types.StringNull()
				updatePlan = true
			}
		} else if planNotificationMode.IsUnknown() { //PF version is new enough for these attributes, set defaults
			planNotificationMode = types.StringValue("NOTIFICATION_PUBLISHER")
			updatePlan = true
		}

		if updatePlan {
			planCertificateExpirationsAttrs["notification_mode"] = planNotificationMode
			planCertificateExpirations, diags = types.ObjectValue(planCertificateExpirations.AttributeTypes(ctx), planCertificateExpirationsAttrs)
			resp.Diagnostics.Append(diags...)
		}
	}

	// Check for attributes only allowed after version 12.0
	planExpiredCertWarningDays := planNotificationsAttrs["expired_certificate_administrative_console_warning_days"].(types.Int64)
	planExpiringCertWarningDays := planNotificationsAttrs["expiring_certificate_administrative_console_warning_days"].(types.Int64)
	if !pfVersionAtLeast120 {
		if internaltypes.IsDefined(planExpiredCertWarningDays) {
			version.AddUnsupportedAttributeError("expired_certificate_administrative_console_warning_days",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		} else if planExpiredCertWarningDays.IsUnknown() {
			planExpiredCertWarningDays = types.Int64Null()
			updatePlan = true
		}

		if internaltypes.IsDefined(planExpiringCertWarningDays) {
			version.AddUnsupportedAttributeError("expiring_certificate_administrative_console_warning_days",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		} else if planExpiringCertWarningDays.IsUnknown() {
			planExpiringCertWarningDays = types.Int64Null()
			updatePlan = true
		}

		if internaltypes.IsDefined(planNotificationsAttrs["thread_pool_exhaustion_notification_settings"]) {
			version.AddUnsupportedAttributeError("thread_pool_exhaustion_notification_settings",
				r.providerConfig.ProductVersion, version.PingFederate1200, &resp.Diagnostics)
		}
	} else {
		if planExpiredCertWarningDays.IsUnknown() {
			planExpiredCertWarningDays = types.Int64Value(14)
			updatePlan = true
		}
		if planExpiringCertWarningDays.IsUnknown() {
			planExpiringCertWarningDays = types.Int64Value(14)
			updatePlan = true
		}
	}

	// Update plan if necessary
	if updatePlan && !resp.Diagnostics.HasError() {
		planNotificationsAttrs["certificate_expirations"] = planCertificateExpirations
		planNotificationsAttrs["expired_certificate_administrative_console_warning_days"] = planExpiredCertWarningDays
		planNotificationsAttrs["expiring_certificate_administrative_console_warning_days"] = planExpiringCertWarningDays

		plan.Notifications, diags = types.ObjectValue(plan.Notifications.AttributeTypes(ctx), planNotificationsAttrs)
		resp.Diagnostics.Append(diags...)

		resp.Plan.Set(ctx, &plan)
	}
}

func addOptionalServerSettingsFields(ctx context.Context, addRequest *client.ServerSettings, plan serverSettingsModel) error {

	if internaltypes.IsDefined(plan.ContactInfo) {
		addRequest.ContactInfo = client.NewContactInfo()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ContactInfo, false)), addRequest.ContactInfo)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Notifications) {
		addRequest.Notifications = client.NewNotificationSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Notifications, true)), addRequest.Notifications)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.RolesAndProtocols) {
		addRequest.RolesAndProtocols = client.NewRolesAndProtocols()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RolesAndProtocols, false)), addRequest.RolesAndProtocols)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.FederationInfo) {
		addRequest.FederationInfo = client.NewFederationInfo()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.FederationInfo, false)), addRequest.FederationInfo)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.EmailServer) {
		addRequest.EmailServer = client.NewEmailServerSettingsWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.EmailServer, true)), addRequest.EmailServer)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.CaptchaSettings) {
		addRequest.CaptchaSettings = client.NewCaptchaSettingsWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.CaptchaSettings, true)), addRequest.CaptchaSettings)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *serverSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings"
}

func (r *serverSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *serverSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServerSettings := client.NewServerSettings()
	err := addOptionalServerSettingsFields(ctx, createServerSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings", err.Error())
		return
	}

	apiCreateServerSettings := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateServerSettings = apiCreateServerSettings.Body(*createServerSettings)
	serverSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(apiCreateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsModel
	diags = readServerSettingsResponse(ctx, serverSettingsResponse, &state, &plan, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read the server settings resource from the PingFederate API and update the state accordingly.
func (r *serverSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetServerSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsResponse(ctx, apiReadServerSettings, &state, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateServerSettings := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewServerSettings()
	err := addOptionalServerSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings", err.Error())
		return
	}

	updateServerSettings = updateServerSettings.Body(*createUpdateRequest)
	updateServerSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(updateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsResponse(ctx, updateServerSettingsResponse, &state, &plan, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
