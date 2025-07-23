// Copyright © 2025 Ping Identity Corporation

package serversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
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
		Description: "Describes the global server configuration settings",
		Attributes: map[string]schema.Attribute{
			"contact_info": schema.SingleNestedAttribute{
				Description: "Information that identifies the server.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Description: "Company name.",
						Computed:    true,
						Optional:    false,
					},
					"email": schema.StringAttribute{
						Description: "Contact email address.",
						Computed:    true,
						Optional:    false,
					},
					"first_name": schema.StringAttribute{
						Description: "Contact first name.",
						Computed:    true,
						Optional:    false,
					},
					"last_name": schema.StringAttribute{
						Description: "Contact last name.",
						Computed:    true,
						Optional:    false,
					},
					"phone": schema.StringAttribute{
						Description: "Contact phone number.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"notifications": schema.SingleNestedAttribute{
				Description: "Notification settings for license and certificate expiration events.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"license_events": schema.SingleNestedAttribute{
						Description: "Settings for license event notifications.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Computed:    true,
								Optional:    false,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Computed:    true,
								Optional:    false,
								Attributes:  resourcelink.ToDataSourceSchema(),
							},
						},
					},
					"certificate_expirations": schema.SingleNestedAttribute{
						Description: "Notification settings for certificate expiration events.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Computed:    true,
								Optional:    false,
							},
							"initial_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when initial warning is sent (in days).",
								Computed:    true,
								Optional:    false,
							},
							"final_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when final warning is sent (in days).",
								Computed:    true,
								Optional:    false,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Computed:    true,
								Optional:    false,
								Attributes:  resourcelink.ToDataSourceSchema(),
							},
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to NOTIFICATION_PUBLISHER.",
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"notify_admin_user_password_changes": schema.BoolAttribute{
						Description: "Determines whether admin users are notified through email when their account is changed.",
						Computed:    true,
						Optional:    false,
					},
					"account_changes_notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher for admin user account changes.",
						Computed:    true,
						Optional:    false,
						Attributes:  resourcelink.ToDataSourceSchema(),
					},
					"metadata_notification_settings": schema.SingleNestedAttribute{
						Description: "Settings for metadata update event notifications.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "The email address where notifications are sent.",
								Computed:    true,
								Optional:    false,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Computed:    true,
								Optional:    false,
								Attributes:  resourcelink.ToDataSourceSchema(),
							},
						},
					},
					"expired_certificate_administrative_console_warning_days": schema.Int64Attribute{
						Description: "Indicates the number of days prior to certificate expiry date, the administrative console warning starts. The default value is 14 days. Supported in PF 12.0 or later.",
						Optional:    true,
						Computed:    false,
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
					},
					"expiring_certificate_administrative_console_warning_days": schema.Int64Attribute{
						Description: "Indicates the number of days past the certificate expiry date, the administrative console warning ends. The default value is 14 days. Supported in PF 12.0 or later.",
						Optional:    false,
						Computed:    true,
					},
					"thread_pool_exhaustion_notification_settings": schema.SingleNestedAttribute{
						Description: "Notification settings for thread pool exhaustion events. Supported in PF 12.0 or later.",
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "Email address where notifications are sent.",
								Optional:    false,
								Computed:    true,
							},
							"thread_dump_enabled": schema.BoolAttribute{
								Description: "Generate a thread dump when approaching thread pool exhaustion.",
								Optional:    false,
								Computed:    true,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    false,
								Computed:    true,
								Attributes:  resourcelink.ToDataSourceSchema(),
							},
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to LOGGING_ONLY.",
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"bulkhead_alert_notification_settings": schema.SingleNestedAttribute{
						Description: "Settings for bulkhead notifications",
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "Email address where notifications are sent.",
								Optional:    false,
								Computed:    true,
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    false,
								Computed:    true,
								Attributes:  resourcelink.ToDataSourceSchema(),
							},
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to LOGGING_ONLY.",
								Optional:    false,
								Computed:    true,
							},
							"thread_dump_enabled": schema.BoolAttribute{
								Description: "Generate a thread dump when a bulkhead reaches its warning threshold or is full.",
								Optional:    false,
								Computed:    true,
							},
						},
					},
				},
			},
			"roles_and_protocols": schema.SingleNestedAttribute{
				Description: "Configure roles and protocols.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"oauth_role": schema.SingleNestedAttribute{
						Description: "OAuth role settings.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"enable_oauth": schema.BoolAttribute{
								Description: "Enable OAuth 2.0 Authorization Server (AS) Role.",
								Computed:    true,
								Optional:    false,
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable Open ID Connect.",
								Computed:    true,
								Optional:    false,
							},
						},
					},
					"idp_role": schema.SingleNestedAttribute{
						Description: "Identity Provider (IdP) settings.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Identity Provider Role.",
								Computed:    true,
								Optional:    false,
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
									},
								},
							},
							"enable_outbound_provisioning": schema.BoolAttribute{
								Description: "Enable Outbound Provisioning.",
								Computed:    true,
								Optional:    false,
							},
						},
					},
					"sp_role": schema.SingleNestedAttribute{
						Description: "Service Provider (SP) settings.",
						Computed:    true,
						Optional:    false,
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description: "Enable Service Provider Role.",
								Computed:    true,
								Optional:    false,
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1.",
								Computed:    true,
								Optional:    false,
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0.",
								Computed:    true,
								Optional:    false,
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation.",
								Computed:    true,
								Optional:    false,
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust.",
								Computed:    true,
								Optional:    false,
							},
							"saml_2_0_profile": schema.SingleNestedAttribute{
								Description: "SAML 2.0 Profile settings.",
								Computed:    true,
								Optional:    false,
								Attributes: map[string]schema.Attribute{
									"enable": schema.BoolAttribute{
										Description: "Enable SAML2.0 profile.",
										Computed:    true,
										Optional:    false,
									},
									"enable_xasp": schema.BoolAttribute{
										Description: "Enable Attribute Requester Mapping for X.509 Attribute Sharing Profile (XASP)",
										Computed:    true,
										Optional:    false,
									},
								},
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable OpenID Connect.",
								Computed:    true,
								Optional:    false,
							},
							"enable_inbound_provisioning": schema.BoolAttribute{
								Description: "Enable Inbound Provisioning.",
								Computed:    true,
								Optional:    false,
							},
						},
					},
					"enable_idp_discovery": schema.BoolAttribute{
						Description: "Enable IdP Discovery.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"federation_info": schema.SingleNestedAttribute{
				Description: "Federation Info.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Description: "The fully qualified host name, port, and path (if applicable) on which the PingFederate server runs.",
						Computed:    true,
						Optional:    false,
					},
					"saml_2_entity_id": schema.StringAttribute{
						Description: "This ID defines your organization as the entity operating the server for SAML 2.0 transactions. It is usually defined as an organization's URL or a DNS address; for example: pingidentity.com. The SAML SourceID used for artifact resolution is derived from this ID using SHA1.",
						Computed:    true,
						Optional:    false,
					},
					"saml_1x_issuer_id": schema.StringAttribute{
						Description: "This ID identifies your federation server for SAML 1.x transactions. As with SAML 2.0, it is usually defined as an organization's URL or a DNS address. The SourceID used for artifact resolution is derived from this ID using SHA1.",
						Computed:    true,
						Optional:    false,
					},
					"saml_1x_source_id": schema.StringAttribute{
						Description: "If supplied, the Source ID value entered here is used for SAML 1.x, instead of being derived from the SAML 1.x Issuer/Audience.",
						Computed:    true,
						Optional:    false,
					},
					"wsfed_realm": schema.StringAttribute{
						Description: "The URI of the realm associated with the PingFederate server. A realm represents a single unit of security administration or trust.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
	}
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
	apiReadServerSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetServerSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
		return
	}

	diags = readServerSettingsResponse(ctx, apiReadServerSettings, &state, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
