package serversettings

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
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

type serverSettingsResourceModel struct {
	Id                types.String `tfsdk:"id"`
	ContactInfo       types.Object `tfsdk:"contact_info"`
	Notifications     types.Object `tfsdk:"notifications"`
	RolesAndProtocols types.Object `tfsdk:"roles_and_protocols"`
	FederationInfo    types.Object `tfsdk:"federation_info"`
	EmailServer       types.Object `tfsdk:"email_server"`
	CaptchaSettings   types.Object `tfsdk:"captcha_settings"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Server Settings",
		Attributes: map[string]schema.Attribute{
			"contact_info": schema.SingleNestedAttribute{
				Description: "Information that identifies the server.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Description: "Company name.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"email": schema.StringAttribute{
						Description: "Contact email address.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"first_name": schema.StringAttribute{
						Description: "Contact first name.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"last_name": schema.StringAttribute{
						Description: "Contact last name.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"phone": schema.StringAttribute{
						Description: "Contact phone number.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"notifications": schema.SingleNestedAttribute{
				Description: "Notification settings for license and certificate expiration events.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"license_events": schema.SingleNestedAttribute{
						Description: "Settings for license event notifications.",
						Optional:    true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The ID of the resource.",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"location": schema.StringAttribute{
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
					"certificate_expirations": schema.SingleNestedAttribute{
						Description: "Settings for license event notifications.",
						Optional:    true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"initial_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when initial warning is sent (in days).",
								Optional:    true,
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
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
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The ID of the resource.",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"location": schema.StringAttribute{
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
					"notify_admin_user_password_changes": schema.BoolAttribute{
						Description: "Determines whether admin users are notified through email when their account is changed.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"account_changes_notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher for admin user account changes.",
						Optional:    true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"metadata_notification_settings": schema.SingleNestedAttribute{
						Description: "Settings for metadata update event notifications.",
						Optional:    true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
										"Invalid email address! Must be of the form '<address>@<company>.<domain>', where 'domain' contains only alphabetic characters and is at least 2 characters in length.",
									),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The ID of the resource.",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"location": schema.StringAttribute{
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
				},
			},
			"roles_and_protocols": schema.SingleNestedAttribute{
				Description: "Configure roles and protocols.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"oauth_role": schema.SingleNestedAttribute{
						Description: "OAuth role settings.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable_oauth": schema.BoolAttribute{
								Description: "Enable OAuth 2.0 Authorization Server (AS) Role.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable Open ID Connect.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"idp_role": schema.SingleNestedAttribute{
						Description: "Identity Provider (IdP) settings.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Identity Provider Role.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"enable_auto_connect": schema.BoolAttribute{
										Description: "This property has been deprecated and no longer used.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
							"enable_outbound_provisioning": schema.BoolAttribute{
								Description: "Enable Outbound Provisioning.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"sp_role": schema.SingleNestedAttribute{
						Description: "Service Provider (SP) settings.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Service Provider Role.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"enable_auto_connect": schema.BoolAttribute{
										Description: "This property has been deprecated and no longer used.",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"enable_xasp": schema.BoolAttribute{
										Description: "Enable Attribute Requester Mapping for X.509 Attribute Sharing Profile (XASP)",
										Computed:    true,
										Optional:    false,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable OpenID Connect.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_inbound_provisioning": schema.BoolAttribute{
								Description: "Enable Inbound Provisioning.",
								Computed:    true,
								Optional:    false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"enable_idp_discovery": schema.BoolAttribute{
						Description: "Enable IdP Discovery.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"federation_info": schema.SingleNestedAttribute{
				Description: "Federation Info.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Description: "The fully qualified host name, port, and path (if applicable) on which the PingFederate server runs.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(https?:\/\/)`),
								"Invalid entry for \"base_url\"! This value must start with 'http://' or 'https://'",
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"saml_2_entity_id": schema.StringAttribute{
						Description: "This ID defines your organization as the entity operating the server for SAML 2.0 transactions. It is usually defined as an organization's URL or a DNS address; for example: pingidentity.com. The SAML SourceID used for artifact resolution is derived from this ID using SHA1.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"auto_connect_entity_id": schema.StringAttribute{
						Description: "This property has been deprecated and no longer used",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"saml_1x_issuer_id": schema.StringAttribute{
						Description: "This ID identifies your federation server for SAML 1.x transactions. As with SAML 2.0, it is usually defined as an organization's URL or a DNS address. The SourceID used for artifact resolution is derived from this ID using SHA1.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"saml_1x_source_id": schema.StringAttribute{
						Description: "If supplied, the Source ID value entered here is used for SAML 1.x, instead of being derived from the SAML 1.x Issuer/Audience.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"wsfed_realm": schema.StringAttribute{
						Description: "The URI of the realm associated with the PingFederate server. A realm represents a single unit of security administration or trust.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"email_server": schema.SingleNestedAttribute{
				Description: "Email Server Settings.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"source_addr": schema.StringAttribute{
						Description: "The email address that appears in the 'From' header line in email messages generated by PingFederate. The address must be in valid format but need not be set up on your system.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email_server": schema.StringAttribute{
						Description: "The IP address or hostname of your email server.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
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
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"use_tls": schema.BoolAttribute{
						Description: "Requires the use of the STARTTLS protocol on the port specified by 'port'.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"verify_hostname": schema.BoolAttribute{
						Description: "If useSSL or useTLS is enabled, this flag determines whether the email server hostname is verified against the server's SMTPS certificate.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_utf8_message_headers": schema.BoolAttribute{
						Description: "Only set this flag to true if the email server supports UTF-8 characters in message headers. Otherwise, this is defaulted to false.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"use_debugging": schema.BoolAttribute{
						Description: "Turns on detailed error messages for the PingFederate server log to help troubleshoot any problems.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
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
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"captcha_settings": schema.SingleNestedAttribute{
				Description: "Captcha Settings.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"site_key": schema.StringAttribute{
						Description: "Site key for reCAPTCHA.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"secret_key": schema.StringAttribute{
						Description: "Secret key for reCAPTCHA. GETs will not return this attribute. To update this field, specify the new value in this attribute.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
	config.AddCommonSchema(&schema)
	resp.Schema = schema
}

// ValidateConfig validates the configuration of the server settings resource.
// It checks the email format of contact_info and notifications email addresses,
// the URL format of the federation_info base_url, and the email format and hostname/IP
// validity of the email_server source_addr and email_server email_server attributes.
// It also checks that the email_server use_ssl and use_tls attributes are not both set to true.
func (r *serverSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	var model serverSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	////////////////////////////////////
	// CONTACT INFO
	////////////////////////////////////
	// Validate contact_info email value
	ciEmail := model.ContactInfo.Attributes()["email"]
	if internaltypes.IsDefined(ciEmail) {
		ciEmailStringValue := ciEmail.(types.String)
		if internaltypes.IsNonEmptyString(ciEmailStringValue) && !internaltypes.IsEmailFormat(ciEmailStringValue.ValueString()) {
			resp.Diagnostics.AddError("Invalid Email Format!", fmt.Sprintf("Please provide a valid email address - \"%s\" needs to be in a valid email format according to RFC 5322.  For example, \"<user>@<company>.<tld>\"", ciEmailStringValue.ValueString()))
		}
	}

	////////////////////////////////////
	// NOTIFICATIONS
	////////////////////////////////////
	// Validate license events email_address value
	if internaltypes.IsDefined(model.Notifications) {
		notificationAttrs := model.Notifications.Attributes()

		if internaltypes.IsDefined(notificationAttrs["license_events"].(types.Object).Attributes()["email_address"]) {
			nLicEmailAddrStringValue := notificationAttrs["license_events"].(types.Object).Attributes()["email_address"].(types.String)
			if internaltypes.IsNonEmptyString(nLicEmailAddrStringValue) && !internaltypes.IsEmailFormat(nLicEmailAddrStringValue.ValueString()) {
				resp.Diagnostics.AddError("Invalid Email Format!", fmt.Sprintf("Please provide a valid email address - \"%s\" needs to be in a valid email format according to RFC 5322.  For example, \"<user>@<company>.<tld>\"", nLicEmailAddrStringValue.ValueString()))
			}
		}

		// Validate certificate_expiration events email_address value
		if internaltypes.IsDefined(notificationAttrs["certificate_expirations"]) {
			nCertEmailAddrStringValue := notificationAttrs["certificate_expirations"].(types.Object).Attributes()["email_address"].(types.String)
			if internaltypes.IsNonEmptyString(nCertEmailAddrStringValue) && !internaltypes.IsEmailFormat(nCertEmailAddrStringValue.ValueString()) {
				resp.Diagnostics.AddError("Invalid Email Format!", fmt.Sprintf("Please provide a valid email address - \"%s\" needs to be in a valid email format according to RFC 5322.  For example, \"<user>@<company>.<tld>\"", nCertEmailAddrStringValue.ValueString()))
			}
		}
	}

	////////////////////////////////////
	// FEDERATION INFO
	////////////////////////////////////
	// Validate base_url value
	if !internaltypes.IsUrlFormat(model.FederationInfo.Attributes()["base_url"].(types.String).ValueString()) {
		resp.Diagnostics.AddError("Invalid URL Format!", fmt.Sprintf("Please provide a valid origin. Origin \"%s\" needs to be in a valid URL-like format - \"http(s)//:<value>.<domain>\"", model.FederationInfo.Attributes()["base_url"].(types.String).ValueString()))
	}
	////////////////////////////////////
	// EMAIL SERVER
	////////////////////////////////////
	// Validate email_server source_addr value
	if internaltypes.IsDefined(model.EmailServer) {
		esAttrs := model.EmailServer.Attributes()
		if internaltypes.IsDefined(esAttrs["source_addr"]) && internaltypes.IsDefined(esAttrs["email_server"]) {
			if internaltypes.IsNonEmptyString(esAttrs["source_addr"].(types.String)) && internaltypes.IsNonEmptyString(esAttrs["email_server"].(types.String)) {
				// Validate source_addr host value
				if !internaltypes.IsEmailFormat(esAttrs["source_addr"].(types.String).ValueString()) {
					resp.Diagnostics.AddError("Invalid Email Format!", fmt.Sprintf("Please provide a valid email address - \"%s\" needs to be in a valid email format according to RFC 5322.  For example, \"<user>@<company>.<tld>\"", esAttrs["source_addr"].(types.String).ValueString()))
				}
				// Validate email_server host value
				if internaltypes.IsNonEmptyString(esAttrs["email_server"].(types.String)) && !internaltypes.IsValidHostnameOrIp(esAttrs["email_server"].(types.String).ValueString()) {
					resp.Diagnostics.AddError("Invalid hostname or IP!", fmt.Sprintf("Please provide a valid hostname or IP address - \"%s\" is invalid", esAttrs["email_server"].(types.String).ValueString()))
				}
			}
		}

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

func addOptionalServerSettingsFields(ctx context.Context, addRequest *client.ServerSettings, plan serverSettingsResourceModel) error {

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

	if internaltypes.ObjContainsNoEmptyVals(plan.CaptchaSettings) {
		addRequest.CaptchaSettings = client.NewCaptchaSettings()
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

func readServerSettingsResponse(ctx context.Context, r *client.ServerSettings, state *serverSettingsResourceModel, plan *serverSettingsResourceModel) {

	//////////////////////////////////////////////////
	// variables for read response
	//////////////////////////////////////////////////
	// emptyString is a variable initialized with an empty string value.
	emptyString := ""
	//TODO placeholder?
	state.Id = types.StringValue("id")

	//////////////////////////////////////////////////
	// CONTACT INFO
	//////////////////////////////////////////////////
	contactInfoAttrType := map[string]attr.Type{
		"company":    basetypes.StringType{},
		"email":      basetypes.StringType{},
		"first_name": basetypes.StringType{},
		"last_name":  basetypes.StringType{},
		"phone":      basetypes.StringType{},
	}
	// add to state
	state.ContactInfo, _ = types.ObjectValueFrom(ctx, contactInfoAttrType, r.ContactInfo)

	//////////////////////////////////////////////
	// NOTIFICATIONS
	//////////////////////////////////////////////
	// nested object
	certificateExpirationsAttrType := map[string]attr.Type{
		"email_address":              basetypes.StringType{},
		"initial_warning_period":     basetypes.Int64Type{},
		"final_warning_period":       basetypes.Int64Type{},
		"notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	// nested object
	notificationSettingsAttrType := map[string]attr.Type{
		"email_address":              basetypes.StringType{},
		"notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	// build object map for notifications from pieces above
	notificationsAttrType := map[string]attr.Type{
		"license_events":                             basetypes.ObjectType{AttrTypes: notificationSettingsAttrType},
		"certificate_expirations":                    basetypes.ObjectType{AttrTypes: certificateExpirationsAttrType},
		"notify_admin_user_password_changes":         basetypes.BoolType{},
		"account_changes_notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"metadata_notification_settings":             basetypes.ObjectType{AttrTypes: notificationSettingsAttrType},
	}

	state.Notifications, _ = types.ObjectValueFrom(ctx, notificationsAttrType, r.Notifications)

	//////////////////////////////////////////////
	// ROLES AND PROTOCOLS
	//////////////////////////////////////////////
	// nested object
	oauthRoleAttrType := map[string]attr.Type{
		"enable_oauth":           basetypes.BoolType{},
		"enable_open_id_connect": basetypes.BoolType{},
	}
	//	retrieve values for saving to state
	oauthRoleVal, _ := types.ObjectValueFrom(ctx, oauthRoleAttrType, r.RolesAndProtocols.OauthRole)

	// nested object
	idpSaml20ProfileAttrType := map[string]attr.Type{
		"enable":              basetypes.BoolType{},
		"enable_auto_connect": basetypes.BoolType{},
	}
	// retrieve values for saving to state
	idpSaml20ProfileVal, _ := types.ObjectValueFrom(ctx, idpSaml20ProfileAttrType, r.RolesAndProtocols.IdpRole.Saml20Profile)

	// nested object
	spSaml20ProfileAttrType := map[string]attr.Type{
		"enable":              basetypes.BoolType{},
		"enable_auto_connect": basetypes.BoolType{},
		"enable_xasp":         basetypes.BoolType{},
	}
	// retrieve values for saving to state
	spSaml20ProfileVal, _ := types.ObjectValueFrom(ctx, spSaml20ProfileAttrType, r.RolesAndProtocols.SpRole.Saml20Profile)

	// nested object
	idpRoleAttrType := map[string]attr.Type{
		"enable":                       basetypes.BoolType{},
		"enable_saml_1_1":              basetypes.BoolType{},
		"enable_saml_1_0":              basetypes.BoolType{},
		"enable_ws_fed":                basetypes.BoolType{},
		"enable_ws_trust":              basetypes.BoolType{},
		"saml_2_0_profile":             basetypes.ObjectType{AttrTypes: idpSaml20ProfileAttrType},
		"enable_outbound_provisioning": basetypes.BoolType{},
	}

	// retrieve values for saving to state
	idpRoleAttrValue := map[string]attr.Value{
		"enable":                       types.BoolPointerValue(r.RolesAndProtocols.IdpRole.Enable),
		"enable_saml_1_1":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableSaml11),
		"enable_saml_1_0":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableSaml10),
		"enable_ws_fed":                types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableWsFed),
		"enable_ws_trust":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableWsTrust),
		"saml_2_0_profile":             idpSaml20ProfileVal,
		"enable_outbound_provisioning": types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableOutboundProvisioning),
	}

	// save IDP role to state
	idpRoleVal, _ := types.ObjectValue(idpRoleAttrType, idpRoleAttrValue)

	// nested object
	spRoleAttrType := map[string]attr.Type{
		"enable":                      basetypes.BoolType{},
		"enable_saml_1_1":             basetypes.BoolType{},
		"enable_saml_1_0":             basetypes.BoolType{},
		"enable_ws_fed":               basetypes.BoolType{},
		"enable_ws_trust":             basetypes.BoolType{},
		"saml_2_0_profile":            basetypes.ObjectType{AttrTypes: spSaml20ProfileAttrType},
		"enable_open_id_connect":      basetypes.BoolType{},
		"enable_inbound_provisioning": basetypes.BoolType{},
	}

	// 	retrieve values for saving to state
	spRoleAttrValue := map[string]attr.Value{
		"enable":                      types.BoolPointerValue(r.RolesAndProtocols.SpRole.Enable),
		"enable_saml_1_1":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableSaml11),
		"enable_saml_1_0":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableSaml10),
		"enable_ws_fed":               types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableWsFed),
		"enable_ws_trust":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableWsTrust),
		"saml_2_0_profile":            spSaml20ProfileVal,
		"enable_open_id_connect":      types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableOpenIDConnect),
		"enable_inbound_provisioning": types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableInboundProvisioning),
	}
	// save SP role to state
	spRoleVal, _ := types.ObjectValue(spRoleAttrType, spRoleAttrValue)

	// build object map for roles and protocols from pieces above
	rolesAndProtocolsAttrType := map[string]attr.Type{
		"oauth_role":           basetypes.ObjectType{AttrTypes: oauthRoleAttrType},
		"idp_role":             basetypes.ObjectType{AttrTypes: idpRoleAttrType},
		"sp_role":              basetypes.ObjectType{AttrTypes: spRoleAttrType},
		"enable_idp_discovery": basetypes.BoolType{},
	}

	// put the values together into state
	rolesAndProtocolsAttrTypeValues := map[string]attr.Value{
		"oauth_role":           oauthRoleVal,
		"idp_role":             idpRoleVal,
		"sp_role":              spRoleVal,
		"enable_idp_discovery": types.BoolPointerValue(r.RolesAndProtocols.EnableIdpDiscovery),
	}
	state.RolesAndProtocols, _ = types.ObjectValue(rolesAndProtocolsAttrType, rolesAndProtocolsAttrTypeValues)

	//////////////////////////////////////////////
	// FEDERATION INFO
	//////////////////////////////////////////////
	federationInfoAttrType := map[string]attr.Type{
		"base_url":               basetypes.StringType{},
		"saml_2_entity_id":       basetypes.StringType{},
		"auto_connect_entity_id": basetypes.StringType{},
		"saml_1x_issuer_id":      basetypes.StringType{},
		"saml_1x_source_id":      basetypes.StringType{},
		"wsfed_realm":            basetypes.StringType{},
	}

	federationInfoAttrValue := map[string]attr.Value{
		"base_url":               types.StringPointerValue(r.FederationInfo.BaseUrl),
		"saml_2_entity_id":       types.StringPointerValue(r.FederationInfo.Saml2EntityId),
		"auto_connect_entity_id": types.StringPointerValue(&emptyString),
		"saml_1x_issuer_id":      types.StringPointerValue(r.FederationInfo.Saml1xIssuerId),
		"saml_1x_source_id":      types.StringPointerValue(r.FederationInfo.Saml1xSourceId),
		"wsfed_realm":            types.StringPointerValue(r.FederationInfo.WsfedRealm),
	}

	state.FederationInfo, _ = types.ObjectValue(federationInfoAttrType, federationInfoAttrValue)

	//////////////////////////////////////////////
	// EMAIL SERVER
	//////////////////////////////////////////////
	emailServerAttrType := map[string]attr.Type{
		"source_addr":                 basetypes.StringType{},
		"email_server":                basetypes.StringType{},
		"port":                        basetypes.Int64Type{},
		"ssl_port":                    basetypes.Int64Type{},
		"timeout":                     basetypes.Int64Type{},
		"retry_attempts":              basetypes.Int64Type{},
		"retry_delay":                 basetypes.Int64Type{},
		"use_ssl":                     basetypes.BoolType{},
		"use_tls":                     basetypes.BoolType{},
		"verify_hostname":             basetypes.BoolType{},
		"enable_utf8_message_headers": basetypes.BoolType{},
		"use_debugging":               basetypes.BoolType{},
		"username":                    basetypes.StringType{},
		"password":                    basetypes.StringType{},
	}

	// get email creds with function
	// if username and password are not set, return null values
	var getEmailCreds = func() (*string, string) {
		if plan.EmailServer.Attributes()["username"] != nil && plan.EmailServer.Attributes()["password"] != nil {
			username := plan.EmailServer.Attributes()["username"].(types.String).ValueStringPointer()
			password := plan.EmailServer.Attributes()["password"].(types.String).ValueString()
			return username, password
		} else {
			return types.StringNull().ValueStringPointer(), types.StringNull().ValueString()
		}
	}

	// retrieve values for saving to state
	username, password := getEmailCreds()
	emailServerAttrValue := map[string]attr.Value{
		"source_addr":                 types.StringValue(r.EmailServer.GetSourceAddr()),
		"email_server":                types.StringValue(r.EmailServer.GetEmailServer()),
		"port":                        types.Int64Value(r.EmailServer.GetPort()),
		"ssl_port":                    types.Int64Value(r.EmailServer.GetSslPort()),
		"timeout":                     types.Int64Value(r.EmailServer.GetTimeout()),
		"retry_attempts":              types.Int64Value(r.EmailServer.GetRetryAttempts()),
		"retry_delay":                 types.Int64Value(r.EmailServer.GetRetryDelay()),
		"use_ssl":                     types.BoolValue(r.EmailServer.GetUseSSL()),
		"use_tls":                     types.BoolValue(r.EmailServer.GetUseTLS()),
		"verify_hostname":             types.BoolValue(r.EmailServer.GetVerifyHostname()),
		"enable_utf8_message_headers": types.BoolValue(r.EmailServer.GetEnableUtf8MessageHeaders()),
		"use_debugging":               types.BoolValue(r.EmailServer.GetUseDebugging()),
		"username":                    types.StringPointerValue(username),
		"password":                    types.StringValue(password),
	}

	state.EmailServer, _ = types.ObjectValue(emailServerAttrType, emailServerAttrValue)

	//////////////////////////////////////////////
	// CAPTCHA SETTINGS
	//////////////////////////////////////////////
	captchaSettingsAttrType := map[string]attr.Type{
		"site_key":   basetypes.StringType{},
		"secret_key": basetypes.StringType{},
	}

	var getCaptchaSettingsAttrValue = func() map[string]attr.Value {
		if internaltypes.ObjContainsNoEmptyVals(plan.CaptchaSettings) {
			return map[string]attr.Value{
				"site_key":   types.StringPointerValue(r.CaptchaSettings.SiteKey),
				"secret_key": types.StringValue(plan.CaptchaSettings.Attributes()["secret_key"].(types.String).ValueString()),
			}
		} else {
			return map[string]attr.Value{
				"site_key":   types.StringPointerValue(&emptyString),
				"secret_key": types.StringValue(emptyString),
			}
		}
	}

	state.CaptchaSettings, _ = types.ObjectValue(captchaSettingsAttrType, getCaptchaSettingsAttrValue())
}

func (r *serverSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsResourceModel

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
	_, requestErr := createServerSettings.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Server Settings: %s", requestErr.Error())
	}

	apiCreateServerSettings := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettings = apiCreateServerSettings.Body(*createServerSettings)
	serverSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(apiCreateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings", err, httpResp)
		return
	}
	_, responseErr := serverSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	var state serverSettingsResourceModel

	readServerSettingsResponse(ctx, serverSettingsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read the server settings resource from the PingFederate API and update the state accordingly.
// It retrieves the current state of the resource, sends a GET request to the PingFederate API to get the server settings,
// and updates the state with the response. If an error occurs, it logs the error and returns the error message.
// If the server settings resource is not found, it removes the resource from the state.
// It also logs the response JSON and sets the refreshed state.
// If the response is empty, it logs a warning and removes the resource from the state.
func (r *serverSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadServerSettings.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	readServerSettingsResponse(ctx, apiReadServerSettings, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state serverSettingsResourceModel
	req.State.Get(ctx, &state)
	updateServerSettings := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewServerSettings()
	err := addOptionalServerSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Server Settings: %s", requestErr.Error())
	}
	updateServerSettings = updateServerSettings.Body(*createUpdateRequest)
	updateServerSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(updateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateServerSettingsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Server Settings: %s", responseErr.Error())
	}
	// Read the response
	readServerSettingsResponse(ctx, updateServerSettingsResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Set a placeholder id value to appease terraform.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
