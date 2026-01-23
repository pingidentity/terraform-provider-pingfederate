// Copyright © 2025 Ping Identity Corporation

package serversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email": schema.StringAttribute{
						Description: "Contact email address.",
						Optional:    true,
						Validators: []validator.String{
							configvalidators.ValidEmail(),
							stringvalidator.LengthAtLeast(1),
						},
					},
					"first_name": schema.StringAttribute{
						Description: "Contact first name.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"last_name": schema.StringAttribute{
						Description: "Contact last name.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"phone": schema.StringAttribute{
						Description: "Contact phone number.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
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
								Description: "Email addresses where notifications are sent. Separate multiple email addresses with commas.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
									stringvalidator.LengthAtLeast(1),
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
								Description: "Email addresses where notifications are sent. Separate multiple email addresses with commas.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
									stringvalidator.LengthAtLeast(1),
								},
							},
							"initial_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when initial warning is sent (in days).",
								Optional:    true,
							},
							"final_warning_period": schema.Int64Attribute{
								Description: "Time before certificate expiration when final warning is sent (in days). Must be between `1` and `99999` days.",
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
								Description: "The mode of notification. Supported values are `NOTIFICATION_PUBLISHER` and `LOGGING_ONLY`. Set to `NOTIFICATION_PUBLISHER` to enable email notifications and server log messages. Set to `LOGGING_ONLY` to enable server log messages. Defaults to `NOTIFICATION_PUBLISHER`.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("NOTIFICATION_PUBLISHER"),
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
						Description: "Determines whether admin users are notified through email when their account is changed. Default value is `false`.",
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
								Description: "Email addresses where notifications are sent. Separate multiple email addresses with commas.",
								Required:    true,
								Validators: []validator.String{
									configvalidators.ValidEmail(),
									stringvalidator.LengthAtLeast(1),
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
						Description: "Indicates the number of days prior to certificate expiry date, the administrative console warning starts. The default value is `14` days. Supported in PF `12.0` or later.",
						Optional:    true,
						Computed:    true,
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
					},
					"expiring_certificate_administrative_console_warning_days": schema.Int64Attribute{
						Description: "Indicates the number of days past the certificate expiry date, the administrative console warning ends. The default value is `14` days. Supported in PF `12.0` or later.",
						Optional:    true,
						Computed:    true,
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
					},
					"thread_pool_exhaustion_notification_settings": schema.SingleNestedAttribute{
						Description: "Notification settings for thread pool exhaustion events. Supported in PF `12.0` or later.",
						Optional:    true,
						Computed:    true,
						Default:     objectdefault.StaticValue(types.ObjectNull(threadPoolExhaustionNotificationSettingsAttrType)),
						// Default will be set in ModifyPlan method. Once we drop support for pre-12.0 versions, we can set the default here instead.
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "Email addresses where notifications are sent. Separate multiple email addresses with commas.",
								Optional:    true,
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
								Description: "The mode of notification. Supported values are `NOTIFICATION_PUBLISHER` and `LOGGING_ONLY`. Set to `NOTIFICATION_PUBLISHER` to enable email notifications and server log messages. Set to `LOGGING_ONLY` to enable server log messages. Defaults to `LOGGING_ONLY`.",
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
					"bulkhead_alert_notification_settings": schema.SingleNestedAttribute{
						Description: "Settings for bulkhead notifications",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"email_address": schema.StringAttribute{
								Description: "Email addresses where notifications are sent. Separate multiple email addresses with commas.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"notification_publisher_ref": schema.SingleNestedAttribute{
								Description: "Reference to the associated notification publisher.",
								Optional:    true,
								Attributes:  resourcelink.ToSchema(),
							},
							"notification_mode": schema.StringAttribute{
								Description: "The mode of notification. Supported values are `NOTIFICATION_PUBLISHER` and `LOGGING_ONLY`. Set to `NOTIFICATION_PUBLISHER` to enable email notifications and server log messages. Set to `LOGGING_ONLY` to enable server log messages. Defaults to `LOGGING_ONLY`.",
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
							"thread_dump_enabled": schema.BoolAttribute{
								Description: "Generate a thread dump when a bulkhead reaches its warning threshold or is full. Default is `true`.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
				},
			},
			"roles_and_protocols": schema.SingleNestedAttribute{
				Description: "Configure roles and protocols. As of PingFederate `12.0`: This property has been deprecated and is no longer used. All Roles and protocols are always enabled.",
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
								Description: "Enable OAuth 2.0 Authorization Server (AS) Role. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable Open ID Connect. Default is `true`.",
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
								Description: "Enable Identity Provider Role. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust. Default is `true`.",
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
										Description: "Enable SAML2.0 profile. Default is `true`.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
								},
							},
							"enable_outbound_provisioning": schema.BoolAttribute{
								Description: "Enable Outbound Provisioning. Default is `true`.",
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
								Description: "Enable Service Provider Role. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_1": schema.BoolAttribute{
								Description: "Enable SAML 1.1. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_saml_1_0": schema.BoolAttribute{
								Description: "Enable SAML 1.0. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_fed": schema.BoolAttribute{
								Description: "Enable WS Federation. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_ws_trust": schema.BoolAttribute{
								Description: "Enable WS Trust. Default is `true`.",
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
										Description: "Enable SAML2.0 profile. Default is `true`.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
									"enable_xasp": schema.BoolAttribute{
										Description: "Enable Attribute Requester Mapping for X.509 Attribute Sharing Profile (XASP). Default is `true`.",
										Computed:    true,
										Optional:    false,
										Default:     booldefault.StaticBool(true),
									},
								},
							},
							"enable_open_id_connect": schema.BoolAttribute{
								Description: "Enable OpenID Connect. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
							"enable_inbound_provisioning": schema.BoolAttribute{
								Description: "Enable Inbound Provisioning. Default is `true`.",
								Computed:    true,
								Optional:    false,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
					"enable_idp_discovery": schema.BoolAttribute{
						Description: "Enable IdP Discovery. Default is `true`.",
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
							stringvalidator.LengthAtLeast(1),
							configvalidators.DoesNotEndWith("/"),
						},
					},
					"saml_2_entity_id": schema.StringAttribute{
						Description: "This ID defines your organization as the entity operating the server for SAML 2.0 transactions. It is usually defined as an organization's URL or a DNS address; for example: pingidentity.com. The SAML SourceID used for artifact resolution is derived from this ID using SHA1.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
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
						Validators: []validator.String{
							stringvalidator.LengthBetween(40, 40),
						},
						Default: stringdefault.StaticString(""),
					},
					"wsfed_realm": schema.StringAttribute{
						Description: "The URI of the realm associated with the PingFederate server. A realm represents a single unit of security administration or trust.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
				},
			},
		},
	}
	resp.Schema = schema
}

func (r *serverSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.0 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast121 := compare >= 0
	var plan *serverSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil || !internaltypes.IsDefined(plan.Notifications) {
		return
	}

	var diags diag.Diagnostics
	updatePlan := false
	planNotificationsAttrs := plan.Notifications.Attributes()

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

	if !pfVersionAtLeast121 {
		if internaltypes.IsDefined(planNotificationsAttrs["bulkhead_alert_notification_settings"]) {
			version.AddUnsupportedAttributeError("bulkhead_alert_notification_settings",
				r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
		}
	}

	// Update plan if necessary
	if updatePlan && !resp.Diagnostics.HasError() {
		planNotificationsAttrs["expired_certificate_administrative_console_warning_days"] = planExpiredCertWarningDays
		planNotificationsAttrs["expiring_certificate_administrative_console_warning_days"] = planExpiringCertWarningDays

		plan.Notifications, diags = types.ObjectValue(plan.Notifications.AttributeTypes(ctx), planNotificationsAttrs)
		resp.Diagnostics.Append(diags...)

		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

func addOptionalServerSettingsFields(addRequest *client.ServerSettings, model serverSettingsModel) {
	// contact_info
	if !model.ContactInfo.IsNull() && !model.ContactInfo.IsUnknown() {
		contactInfoValue := &client.ContactInfo{}
		contactInfoAttrs := model.ContactInfo.Attributes()
		contactInfoValue.Company = contactInfoAttrs["company"].(types.String).ValueStringPointer()
		contactInfoValue.Email = contactInfoAttrs["email"].(types.String).ValueStringPointer()
		contactInfoValue.FirstName = contactInfoAttrs["first_name"].(types.String).ValueStringPointer()
		contactInfoValue.LastName = contactInfoAttrs["last_name"].(types.String).ValueStringPointer()
		contactInfoValue.Phone = contactInfoAttrs["phone"].(types.String).ValueStringPointer()
		addRequest.ContactInfo = contactInfoValue
	}

	// federation_info
	if !model.FederationInfo.IsNull() && !model.FederationInfo.IsUnknown() {
		federationInfoValue := &client.FederationInfo{}
		federationInfoAttrs := model.FederationInfo.Attributes()
		federationInfoValue.BaseUrl = federationInfoAttrs["base_url"].(types.String).ValueStringPointer()
		federationInfoValue.Saml1xIssuerId = federationInfoAttrs["saml_1x_issuer_id"].(types.String).ValueStringPointer()
		federationInfoValue.Saml1xSourceId = federationInfoAttrs["saml_1x_source_id"].(types.String).ValueStringPointer()
		federationInfoValue.Saml2EntityId = federationInfoAttrs["saml_2_entity_id"].(types.String).ValueStringPointer()
		federationInfoValue.WsfedRealm = federationInfoAttrs["wsfed_realm"].(types.String).ValueStringPointer()
		addRequest.FederationInfo = federationInfoValue
	}

	// notifications
	if !model.Notifications.IsNull() && !model.Notifications.IsUnknown() {
		notificationsValue := &client.NotificationSettings{}
		notificationsAttrs := model.Notifications.Attributes()
		if !notificationsAttrs["account_changes_notification_publisher_ref"].IsNull() && !notificationsAttrs["account_changes_notification_publisher_ref"].IsUnknown() {
			notificationsAccountChangesNotificationPublisherRefValue := &client.ResourceLink{}
			notificationsAccountChangesNotificationPublisherRefAttrs := notificationsAttrs["account_changes_notification_publisher_ref"].(types.Object).Attributes()
			notificationsAccountChangesNotificationPublisherRefValue.Id = notificationsAccountChangesNotificationPublisherRefAttrs["id"].(types.String).ValueString()
			notificationsValue.AccountChangesNotificationPublisherRef = notificationsAccountChangesNotificationPublisherRefValue
		}
		if !notificationsAttrs["bulkhead_alert_notification_settings"].IsNull() && !notificationsAttrs["bulkhead_alert_notification_settings"].IsUnknown() {
			notificationsBulkheadAlertNotificationSettingsValue := &client.BulkheadAlertNotificationSettings{}
			notificationsBulkheadAlertNotificationSettingsAttrs := notificationsAttrs["bulkhead_alert_notification_settings"].(types.Object).Attributes()
			notificationsBulkheadAlertNotificationSettingsValue.EmailAddress = notificationsBulkheadAlertNotificationSettingsAttrs["email_address"].(types.String).ValueString()
			notificationsBulkheadAlertNotificationSettingsValue.NotificationMode = notificationsBulkheadAlertNotificationSettingsAttrs["notification_mode"].(types.String).ValueStringPointer()
			if !notificationsBulkheadAlertNotificationSettingsAttrs["notification_publisher_ref"].IsNull() && !notificationsBulkheadAlertNotificationSettingsAttrs["notification_publisher_ref"].IsUnknown() {
				notificationsBulkheadAlertNotificationSettingsNotificationPublisherRefValue := &client.ResourceLink{}
				notificationsBulkheadAlertNotificationSettingsNotificationPublisherRefAttrs := notificationsBulkheadAlertNotificationSettingsAttrs["notification_publisher_ref"].(types.Object).Attributes()
				notificationsBulkheadAlertNotificationSettingsNotificationPublisherRefValue.Id = notificationsBulkheadAlertNotificationSettingsNotificationPublisherRefAttrs["id"].(types.String).ValueString()
				notificationsBulkheadAlertNotificationSettingsValue.NotificationPublisherRef = notificationsBulkheadAlertNotificationSettingsNotificationPublisherRefValue
			}
			notificationsBulkheadAlertNotificationSettingsValue.ThreadDumpEnabled = notificationsBulkheadAlertNotificationSettingsAttrs["thread_dump_enabled"].(types.Bool).ValueBoolPointer()
			notificationsValue.BulkheadAlertNotificationSettings = notificationsBulkheadAlertNotificationSettingsValue
		}
		if !notificationsAttrs["certificate_expirations"].IsNull() && !notificationsAttrs["certificate_expirations"].IsUnknown() {
			notificationsCertificateExpirationsValue := &client.CertificateExpirationNotificationSettings{}
			notificationsCertificateExpirationsAttrs := notificationsAttrs["certificate_expirations"].(types.Object).Attributes()
			notificationsCertificateExpirationsValue.EmailAddress = notificationsCertificateExpirationsAttrs["email_address"].(types.String).ValueString()
			notificationsCertificateExpirationsValue.FinalWarningPeriod = notificationsCertificateExpirationsAttrs["final_warning_period"].(types.Int64).ValueInt64()
			notificationsCertificateExpirationsValue.InitialWarningPeriod = notificationsCertificateExpirationsAttrs["initial_warning_period"].(types.Int64).ValueInt64Pointer()
			notificationsCertificateExpirationsValue.NotificationMode = notificationsCertificateExpirationsAttrs["notification_mode"].(types.String).ValueStringPointer()
			if !notificationsCertificateExpirationsAttrs["notification_publisher_ref"].IsNull() && !notificationsCertificateExpirationsAttrs["notification_publisher_ref"].IsUnknown() {
				notificationsCertificateExpirationsNotificationPublisherRefValue := &client.ResourceLink{}
				notificationsCertificateExpirationsNotificationPublisherRefAttrs := notificationsCertificateExpirationsAttrs["notification_publisher_ref"].(types.Object).Attributes()
				notificationsCertificateExpirationsNotificationPublisherRefValue.Id = notificationsCertificateExpirationsNotificationPublisherRefAttrs["id"].(types.String).ValueString()
				notificationsCertificateExpirationsValue.NotificationPublisherRef = notificationsCertificateExpirationsNotificationPublisherRefValue
			}
			notificationsValue.CertificateExpirations = notificationsCertificateExpirationsValue
		}
		notificationsValue.ExpiredCertificateAdministrativeConsoleWarningDays = notificationsAttrs["expired_certificate_administrative_console_warning_days"].(types.Int64).ValueInt64Pointer()
		notificationsValue.ExpiringCertificateAdministrativeConsoleWarningDays = notificationsAttrs["expiring_certificate_administrative_console_warning_days"].(types.Int64).ValueInt64Pointer()
		if !notificationsAttrs["license_events"].IsNull() && !notificationsAttrs["license_events"].IsUnknown() {
			notificationsLicenseEventsValue := &client.LicenseEventNotificationSettings{}
			notificationsLicenseEventsAttrs := notificationsAttrs["license_events"].(types.Object).Attributes()
			notificationsLicenseEventsValue.EmailAddress = notificationsLicenseEventsAttrs["email_address"].(types.String).ValueString()
			if !notificationsLicenseEventsAttrs["notification_publisher_ref"].IsNull() && !notificationsLicenseEventsAttrs["notification_publisher_ref"].IsUnknown() {
				notificationsLicenseEventsNotificationPublisherRefValue := &client.ResourceLink{}
				notificationsLicenseEventsNotificationPublisherRefAttrs := notificationsLicenseEventsAttrs["notification_publisher_ref"].(types.Object).Attributes()
				notificationsLicenseEventsNotificationPublisherRefValue.Id = notificationsLicenseEventsNotificationPublisherRefAttrs["id"].(types.String).ValueString()
				notificationsLicenseEventsValue.NotificationPublisherRef = notificationsLicenseEventsNotificationPublisherRefValue
			}
			notificationsValue.LicenseEvents = notificationsLicenseEventsValue
		}
		if !notificationsAttrs["metadata_notification_settings"].IsNull() && !notificationsAttrs["metadata_notification_settings"].IsUnknown() {
			notificationsMetadataNotificationSettingsValue := &client.MetadataEventNotificationSettings{}
			notificationsMetadataNotificationSettingsAttrs := notificationsAttrs["metadata_notification_settings"].(types.Object).Attributes()
			notificationsMetadataNotificationSettingsValue.EmailAddress = notificationsMetadataNotificationSettingsAttrs["email_address"].(types.String).ValueString()
			if !notificationsMetadataNotificationSettingsAttrs["notification_publisher_ref"].IsNull() && !notificationsMetadataNotificationSettingsAttrs["notification_publisher_ref"].IsUnknown() {
				notificationsMetadataNotificationSettingsNotificationPublisherRefValue := &client.ResourceLink{}
				notificationsMetadataNotificationSettingsNotificationPublisherRefAttrs := notificationsMetadataNotificationSettingsAttrs["notification_publisher_ref"].(types.Object).Attributes()
				notificationsMetadataNotificationSettingsNotificationPublisherRefValue.Id = notificationsMetadataNotificationSettingsNotificationPublisherRefAttrs["id"].(types.String).ValueString()
				notificationsMetadataNotificationSettingsValue.NotificationPublisherRef = notificationsMetadataNotificationSettingsNotificationPublisherRefValue
			}
			notificationsValue.MetadataNotificationSettings = notificationsMetadataNotificationSettingsValue
		}
		notificationsValue.NotifyAdminUserPasswordChanges = notificationsAttrs["notify_admin_user_password_changes"].(types.Bool).ValueBoolPointer()
		if !notificationsAttrs["thread_pool_exhaustion_notification_settings"].IsNull() && !notificationsAttrs["thread_pool_exhaustion_notification_settings"].IsUnknown() {
			notificationsThreadPoolExhaustionNotificationSettingsValue := &client.ThreadPoolExhaustionNotificationSettings{}
			notificationsThreadPoolExhaustionNotificationSettingsAttrs := notificationsAttrs["thread_pool_exhaustion_notification_settings"].(types.Object).Attributes()
			notificationsThreadPoolExhaustionNotificationSettingsValue.EmailAddress = notificationsThreadPoolExhaustionNotificationSettingsAttrs["email_address"].(types.String).ValueString()
			notificationsThreadPoolExhaustionNotificationSettingsValue.NotificationMode = notificationsThreadPoolExhaustionNotificationSettingsAttrs["notification_mode"].(types.String).ValueStringPointer()
			if !notificationsThreadPoolExhaustionNotificationSettingsAttrs["notification_publisher_ref"].IsNull() && !notificationsThreadPoolExhaustionNotificationSettingsAttrs["notification_publisher_ref"].IsUnknown() {
				notificationsThreadPoolExhaustionNotificationSettingsNotificationPublisherRefValue := &client.ResourceLink{}
				notificationsThreadPoolExhaustionNotificationSettingsNotificationPublisherRefAttrs := notificationsThreadPoolExhaustionNotificationSettingsAttrs["notification_publisher_ref"].(types.Object).Attributes()
				notificationsThreadPoolExhaustionNotificationSettingsNotificationPublisherRefValue.Id = notificationsThreadPoolExhaustionNotificationSettingsNotificationPublisherRefAttrs["id"].(types.String).ValueString()
				notificationsThreadPoolExhaustionNotificationSettingsValue.NotificationPublisherRef = notificationsThreadPoolExhaustionNotificationSettingsNotificationPublisherRefValue
			}
			notificationsThreadPoolExhaustionNotificationSettingsValue.ThreadDumpEnabled = notificationsThreadPoolExhaustionNotificationSettingsAttrs["thread_dump_enabled"].(types.Bool).ValueBoolPointer()
			notificationsValue.ThreadPoolExhaustionNotificationSettings = notificationsThreadPoolExhaustionNotificationSettingsValue
		}
		addRequest.Notifications = notificationsValue
	}

	// roles_and_protocols
	if !model.RolesAndProtocols.IsNull() && !model.RolesAndProtocols.IsUnknown() {
		rolesAndProtocolsValue := &client.RolesAndProtocols{}
		rolesAndProtocolsAttrs := model.RolesAndProtocols.Attributes()
		rolesAndProtocolsValue.EnableIdpDiscovery = rolesAndProtocolsAttrs["enable_idp_discovery"].(types.Bool).ValueBoolPointer()
		if !rolesAndProtocolsAttrs["idp_role"].IsNull() && !rolesAndProtocolsAttrs["idp_role"].IsUnknown() {
			rolesAndProtocolsIdpRoleValue := &client.IdpRole{}
			rolesAndProtocolsIdpRoleAttrs := rolesAndProtocolsAttrs["idp_role"].(types.Object).Attributes()
			rolesAndProtocolsIdpRoleValue.Enable = rolesAndProtocolsIdpRoleAttrs["enable"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsIdpRoleValue.EnableOutboundProvisioning = rolesAndProtocolsIdpRoleAttrs["enable_outbound_provisioning"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsIdpRoleValue.EnableSaml10 = rolesAndProtocolsIdpRoleAttrs["enable_saml_1_0"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsIdpRoleValue.EnableSaml11 = rolesAndProtocolsIdpRoleAttrs["enable_saml_1_1"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsIdpRoleValue.EnableWsFed = rolesAndProtocolsIdpRoleAttrs["enable_ws_fed"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsIdpRoleValue.EnableWsTrust = rolesAndProtocolsIdpRoleAttrs["enable_ws_trust"].(types.Bool).ValueBoolPointer()
			if !rolesAndProtocolsIdpRoleAttrs["saml_2_0_profile"].IsNull() && !rolesAndProtocolsIdpRoleAttrs["saml_2_0_profile"].IsUnknown() {
				rolesAndProtocolsIdpRoleSaml20ProfileValue := &client.SAML20Profile{}
				rolesAndProtocolsIdpRoleSaml20ProfileAttrs := rolesAndProtocolsIdpRoleAttrs["saml_2_0_profile"].(types.Object).Attributes()
				rolesAndProtocolsIdpRoleSaml20ProfileValue.Enable = rolesAndProtocolsIdpRoleSaml20ProfileAttrs["enable"].(types.Bool).ValueBoolPointer()
				rolesAndProtocolsIdpRoleValue.Saml20Profile = rolesAndProtocolsIdpRoleSaml20ProfileValue
			}
			rolesAndProtocolsValue.IdpRole = rolesAndProtocolsIdpRoleValue
		}
		if !rolesAndProtocolsAttrs["oauth_role"].IsNull() && !rolesAndProtocolsAttrs["oauth_role"].IsUnknown() {
			rolesAndProtocolsOauthRoleValue := &client.OAuthRole{}
			rolesAndProtocolsOauthRoleAttrs := rolesAndProtocolsAttrs["oauth_role"].(types.Object).Attributes()
			rolesAndProtocolsOauthRoleValue.EnableOauth = rolesAndProtocolsOauthRoleAttrs["enable_oauth"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsOauthRoleValue.EnableOpenIdConnect = rolesAndProtocolsOauthRoleAttrs["enable_open_id_connect"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsValue.OauthRole = rolesAndProtocolsOauthRoleValue
		}
		if !rolesAndProtocolsAttrs["sp_role"].IsNull() && !rolesAndProtocolsAttrs["sp_role"].IsUnknown() {
			rolesAndProtocolsSpRoleValue := &client.SpRole{}
			rolesAndProtocolsSpRoleAttrs := rolesAndProtocolsAttrs["sp_role"].(types.Object).Attributes()
			rolesAndProtocolsSpRoleValue.Enable = rolesAndProtocolsSpRoleAttrs["enable"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableInboundProvisioning = rolesAndProtocolsSpRoleAttrs["enable_inbound_provisioning"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableOpenIDConnect = rolesAndProtocolsSpRoleAttrs["enable_open_id_connect"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableSaml10 = rolesAndProtocolsSpRoleAttrs["enable_saml_1_0"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableSaml11 = rolesAndProtocolsSpRoleAttrs["enable_saml_1_1"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableWsFed = rolesAndProtocolsSpRoleAttrs["enable_ws_fed"].(types.Bool).ValueBoolPointer()
			rolesAndProtocolsSpRoleValue.EnableWsTrust = rolesAndProtocolsSpRoleAttrs["enable_ws_trust"].(types.Bool).ValueBoolPointer()
			if !rolesAndProtocolsSpRoleAttrs["saml_2_0_profile"].IsNull() && !rolesAndProtocolsSpRoleAttrs["saml_2_0_profile"].IsUnknown() {
				rolesAndProtocolsSpRoleSaml20ProfileValue := &client.SpSAML20Profile{}
				rolesAndProtocolsSpRoleSaml20ProfileAttrs := rolesAndProtocolsSpRoleAttrs["saml_2_0_profile"].(types.Object).Attributes()
				rolesAndProtocolsSpRoleSaml20ProfileValue.Enable = rolesAndProtocolsSpRoleSaml20ProfileAttrs["enable"].(types.Bool).ValueBoolPointer()
				rolesAndProtocolsSpRoleSaml20ProfileValue.EnableXASP = rolesAndProtocolsSpRoleSaml20ProfileAttrs["enable_xasp"].(types.Bool).ValueBoolPointer()
				rolesAndProtocolsSpRoleValue.Saml20Profile = rolesAndProtocolsSpRoleSaml20ProfileValue
			}
			rolesAndProtocolsValue.SpRole = rolesAndProtocolsSpRoleValue
		}
		addRequest.RolesAndProtocols = rolesAndProtocolsValue
	}
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
	addOptionalServerSettingsFields(createServerSettings, plan)

	apiCreateServerSettings := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateServerSettings = apiCreateServerSettings.Body(*createServerSettings)
	serverSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(apiCreateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsModel
	diags = readServerSettingsResponse(ctx, serverSettingsResponse, &state, &plan)
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
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Server Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readServerSettingsResponse(ctx, apiReadServerSettings, &state, &state)
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
	addOptionalServerSettingsFields(createUpdateRequest, plan)

	updateServerSettings = updateServerSettings.Body(*createUpdateRequest)
	updateServerSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(updateServerSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsModel
	diags = readServerSettingsResponse(ctx, updateServerSettingsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	var state serverSettingsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.FederationInfo.IsNull() || state.FederationInfo.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("federation_info"),
			providererror.InternalProviderError,
			"Cannot delete the server settings resource because the federation_info configuration is missing or unknown")
	}
	resp.Diagnostics.AddWarning(providererror.ConfigurationCannotBeResetError,
		"The pingfederate_server_settings resource has been destroyed but cannot be completely returned to its original state. "+
			"The resource has been removed from Terraform state but the federation_info.base_url and federation_info.saml_2_entity_id configuration remains applied to the environment")
	resetSettings := client.NewServerSettings()
	resetSettings.FederationInfo = client.NewFederationInfo()
	resetSettings.FederationInfo.BaseUrl = state.FederationInfo.Attributes()["base_url"].(types.String).ValueStringPointer()
	resetSettings.FederationInfo.Saml2EntityId = state.FederationInfo.Attributes()["saml_2_entity_id"].(types.String).ValueStringPointer()
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateServerSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*resetSettings)
	_, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateServerSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the server settings", err, httpResp)
	}
}

func (r *serverSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState serverSettingsModel
	emptyState.ContactInfo = types.ObjectNull(contactInfoAttrType)
	emptyState.FederationInfo = types.ObjectNull(federationInfoAttrType)
	emptyState.Notifications = types.ObjectNull(notificationsAttrType)
	emptyState.RolesAndProtocols = types.ObjectNull(rolesAndProtocolsAttrType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
