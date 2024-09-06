package localidentity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localIdentityProfileResource{}
	_ resource.ResourceWithConfigure   = &localIdentityProfileResource{}
	_ resource.ResourceWithImportState = &localIdentityProfileResource{}
)

var (
	authSourcesDefault, _ = types.SetValue(types.ObjectType{AttrTypes: authSourcesAttrTypes}, nil)

	authSourceUpdatePolicyDefault, _ = types.ObjectValue(authSourceUpdatePolicyAttrTypes, map[string]attr.Value{
		"store_attributes":  types.BoolValue(false),
		"retain_attributes": types.BoolValue(false),
		"update_attributes": types.BoolValue(false),
		"update_interval":   types.Int64Value(0),
	})

	emptyStringMapDefault, _ = types.MapValue(types.StringType, map[string]attr.Value{})

	emailVerificationConfigDefault, _ = types.ObjectValue(emailVerificationConfigAttrTypes, map[string]attr.Value{
		"email_verification_enabled":               types.BoolValue(false),
		"verify_email_template_name":               types.StringNull(),
		"email_verification_sent_template_name":    types.StringNull(),
		"email_verification_success_template_name": types.StringNull(),
		"email_verification_error_template_name":   types.StringNull(),
		"email_verification_type":                  types.StringNull(),
		"otp_length":                               types.Int64Null(),
		"otp_retry_attempts":                       types.Int64Null(),
		"allowed_otp_character_set":                types.StringNull(),
		"otp_time_to_live":                         types.Int64Null(),
		"email_verification_otp_template_name":     types.StringNull(),
		"otl_time_to_live":                         types.Int64Null(),
		"field_for_email_to_verify":                types.StringValue(""),
		"field_storing_verification_status":        types.StringValue(""),
		"notification_publisher_ref":               types.ObjectNull(resourcelink.AttrType()),
		"require_verified_email":                   types.BoolNull(),
		"require_verified_email_template_name":     types.StringNull(),
	})
)

// LocalIdentityProfileResource is a helper function to simplify the provider implementation.
func LocalIdentityProfileResource() resource.Resource {

	return &localIdentityProfileResource{}
}

// localIdentityProfileResource is the resource implementation.
type localIdentityProfileResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *localIdentityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a configured local identity profile",
		Attributes: map[string]schema.Attribute{
			"profile_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the local identity profile. It can be any combination of `[a-zA-Z0-9._-]`.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					configvalidators.PingFederateId(),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The local identity profile name. Name is unique.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"apc_id": schema.SingleNestedAttribute{
				Description: "The reference to the authentication policy contract to use for this local identity profile.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"auth_sources": schema.SetNestedAttribute{
				Description: "The local identity authentication sources. Sources are unique.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(authSourcesDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The persistent, unique ID for the local identity authentication source. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.",
							Computed:    true,
							Optional:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								configvalidators.PingFederateId(),
								stringvalidator.LengthAtLeast(1),
							},
						},
						"source": schema.StringAttribute{
							Description: "The local identity authentication source. Source is unique.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"auth_source_update_policy": schema.SingleNestedAttribute{
				Description: "The attribute update policy for authentication sources.",
				Optional:    true,
				Computed:    true,
				// Default set in ModifyPlan
				Attributes: map[string]schema.Attribute{
					"store_attributes": schema.BoolAttribute{
						Description: "Whether or not to store attributes that came from authentication sources. The default value is `false`.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"retain_attributes": schema.BoolAttribute{
						Description: "Whether or not to keep attributes after user disconnects. The default value is `false`.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"update_attributes": schema.BoolAttribute{
						Description: "Whether or not to update attributes when users authenticate. The default value is `false`.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"update_interval": schema.Int64Attribute{
						Description: "The minimum number of days between updates. The default value is `0`.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
				},
			},
			"registration_enabled": schema.BoolAttribute{
				Description: "Whether the registration configuration is enabled or not. The default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"registration_config": schema.SingleNestedAttribute{
				Description: "The local identity profile registration configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"captcha_enabled": schema.BoolAttribute{
						Description: "Whether CAPTCHA is enabled or not in the registration configuration. The default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"captcha_provider_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated CAPTCHA provider.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for the registration configuration.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"create_authn_session_after_registration": schema.BoolAttribute{
						Description: "Whether to create an Authentication Session when registering a local account. The default value is `true`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
					"username_field": schema.StringAttribute{
						Description: "When creating an Authentication Session after registering a local account, PingFederate will pass the Unique ID field's value as the username. If the Unique ID value is not the username, then override which field's value will be used as the username.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"this_is_my_device_enabled": schema.BoolAttribute{
						Description: "Allows users to indicate whether their device is shared or private. In this mode, PingFederate Authentication Sessions will not be stored unless the user indicates the device is private. The default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"registration_workflow": schema.SingleNestedAttribute{
						Description: "The policy fragment to be executed as part of the registration workflow.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"execute_workflow": schema.StringAttribute{
						Description: "This setting indicates whether PingFederate should execute the workflow before or after account creation. The default is to run the registration workflow after account creation. Supported values are `BEFORE_ACCOUNT_CREATION` and `AFTER_ACCOUNT_CREATION`. Requires that `registration_workflow` is also set.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"BEFORE_ACCOUNT_CREATION", "AFTER_ACCOUNT_CREATION"}...),
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("registration_workflow")),
						},
					},
				},
			},
			"profile_config": schema.SingleNestedAttribute{
				Description: "The local identity profile management configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"delete_identity_enabled": schema.BoolAttribute{
						Description: "Whether the end user is allowed to use delete functionality. The default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for end-user profile management.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"field_config": schema.SingleNestedAttribute{
				Description: "The local identity profile field configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"fields": schema.SetNestedAttribute{
						Description: "The field configuration for the local identity profile.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of the local identity field. Supported values are `CHECKBOX`, `CHECKBOX_GROUP`, `DATE`, `DROP_DOWN`, `EMAIL`, `PHONE`, `TEXT`, and `HIDDEN`.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"CHECKBOX", "CHECKBOX_GROUP", "DATE", "DROP_DOWN", "EMAIL", "PHONE", "TEXT", "HIDDEN"}...),
									},
								},
								"id": schema.StringAttribute{
									Description: "Id of the local identity field.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"label": schema.StringAttribute{
									Description: "Label of the local identity field.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"registration_page_field": schema.BoolAttribute{
									Description: "Whether this is a registration page field or not. The default value is `false`.",
									Optional:    true,
									Computed:    true,
									// This default causes issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default:     booldefault.StaticBool(false),
								},
								"profile_page_field": schema.BoolAttribute{
									Description: "Whether this is a profile page field or not. The default value is `false`.",
									Optional:    true,
									Computed:    true,
									// This default causes issues with unexpected plans - see https://github.com/hashicorp/terraform-plugin-framework/issues/867
									// Default:     booldefault.StaticBool(false),
								},
								"attributes": schema.MapAttribute{
									Description: "Attributes of the local identity field.",
									Computed:    true,
									Optional:    true,
									// Default set in ModifyPlan
									ElementType: types.BoolType,
								},
								"options": schema.SetAttribute{
									Description: "The list of options for this selection field.",
									Computed:    true,
									Optional:    true,
									ElementType: types.StringType,
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
									},
								},
								"default_value": schema.StringAttribute{
									Description: "The default value for this field.",
									Optional:    true,
									Computed:    true,
								},
							},
						},
					},
					"strip_space_from_unique_field": schema.BoolAttribute{
						Description: "Strip leading/trailing spaces from unique ID field. The default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"email_verification_config": schema.SingleNestedAttribute{
				Description: "The local identity email verification configuration.",
				Computed:    true,
				Optional:    true,
				// Default set in ModifyPlan
				Attributes: map[string]schema.Attribute{
					"email_verification_enabled": schema.BoolAttribute{
						Description: "Whether the email ownership verification is enabled. The default value is `false`.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"verify_email_template_name": schema.StringAttribute{
						Description: "The template name for verify email. The default is `message-template-email-ownership-verification.html`.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email_verification_sent_template_name": schema.StringAttribute{
						Description: "The template name for email verification sent. The default is `local.identity.email.verification.sent.html`. Note:Only applicable if `email_verification_type` is `OTL`.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email_verification_success_template_name": schema.StringAttribute{
						Description: "The template name for email verification success. The default is `local.identity.email.verification.success.html`.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email_verification_error_template_name": schema.StringAttribute{
						Description: "The template name for email verification error. The default is `local.identity.email.verification.error.html`.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"email_verification_type": schema.StringAttribute{
						Description: "Email Verification Type. Supported values are `OTP` and `OTL`.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"OTP", "OTL"}...),
						},
					},
					"otp_length": schema.Int64Attribute{
						Description: "The OTP length generated for email verification. The default is `8`. Note: Only applicable if `email_verification_type` is `OTP`. The value must be between `5` and `100`.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(5, 100),
						},
					},
					"otp_retry_attempts": schema.Int64Attribute{
						Description: "The number of OTP retry attempts for email verification. The default is `3`. Note: Only applicable if `email_verification_type` is `OTP`.",
						Optional:    true,
					},
					"allowed_otp_character_set": schema.StringAttribute{
						Description: "The allowed character set used to generate the OTP. The default is `23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz`. Note: Only applicable if `email_verification_type` is `OTP`.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"otp_time_to_live": schema.Int64Attribute{
						Description: "Field used OTP time to live. The default is `15`. Note: Only applicable if `email_verification_type` is `OTP`.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_otp_template_name": schema.StringAttribute{
						Description: "The template name for email verification OTP verification. The default is `local.identity.email.verification.otp.html`. Note: Only applicable if `email_verification_type` is `OTP`.",
						Optional:    true,
						Computed:    true,
					},
					"otl_time_to_live": schema.Int64Attribute{
						Description: "Field used OTL time to live. The default is `1440`. Note: Only applicable if `email_verification_type` is `OTL`.",
						Computed:    true,
						Optional:    true,
					},
					"field_for_email_to_verify": schema.StringAttribute{
						Description: "Field used for email ownership verification. Note: Not required when `email_verification_enabled` is set to `false`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"field_storing_verification_status": schema.StringAttribute{
						Description: "Field used for storing email verification status. Note: Not required when `email_verification_enabled` is set to `false`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"require_verified_email": schema.BoolAttribute{
						Description: "Whether the user must verify their email address before they can complete a single sign-on transaction. The default is `false`.",
						Computed:    true,
						Optional:    true,
					},
					"require_verified_email_template_name": schema.StringAttribute{
						Description: "The template to render when the user must verify their email address before they can complete a single sign-on transaction. The default is `local.identity.email.verification.required.html`. Note: Only applicable if `email_verification_type` is OTL and `require_verified_email` is true.",
						Computed:    true,
						Optional:    true,
					},
				},
			},
			"data_store_config": schema.SingleNestedAttribute{
				Description: "The local identity profile data store configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"base_dn": schema.StringAttribute{
						Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"type": schema.StringAttribute{
						Description: "The data store config type. Supported values are `LDAP`, `PING_ONE_LDAP_GATEWAY`, `JDBC`, and `CUSTOM`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
						},
					},
					"data_store_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated data store.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
					},
					"data_store_mapping": schema.MapNestedAttribute{
						Description: "The data store mapping.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The data store attribute type. Supported values are `LDAP`, `PING_ONE_LDAP_GATEWAY`, `JDBC`, and `CUSTOM`.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
									},
								},
								"name": schema.StringAttribute{
									Description: "The data store attribute name.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"metadata": schema.MapAttribute{
									Description: "The data store attribute metadata.",
									Computed:    true,
									Optional:    true,
									Default:     mapdefault.StaticValue(emptyStringMapDefault),
									ElementType: types.StringType,
								},
							},
						},
					},
					"create_pattern": schema.StringAttribute{
						Description: "The Relative DN Pattern that will be used to create objects in the directory.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
					"object_class": schema.StringAttribute{
						Description: "The Object Class used by the new objects stored in the LDAP data store.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
					"auxiliary_object_classes": schema.SetAttribute{
						Description: "The Auxiliary Object Classes used by the new objects stored in the LDAP data store.",
						Optional:    true,
						Computed:    false,
						ElementType: types.StringType,
					},
				},
			},
			"profile_enabled": schema.BoolAttribute{
				Description: "Whether the profile configuration is enabled or not. The default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func addOptionalLocalIdentityProfileFields(ctx context.Context, addRequest *client.LocalIdentityProfile, plan localIdentityProfileModel) error {
	addRequest.Id = plan.ProfileId.ValueStringPointer()

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueString()
	}

	if internaltypes.IsDefined(plan.ApcId) {
		addRequest.ApcId = client.NewLocalIdentityProfileWithDefaults().ApcId
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ApcId, false)), &addRequest.ApcId)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AuthSources) {
		addRequest.AuthSources = client.NewLocalIdentityProfileWithDefaults().AuthSources
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AuthSources, false)), &addRequest.AuthSources)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AuthSourceUpdatePolicy) {
		addRequest.AuthSourceUpdatePolicy = client.NewLocalIdentityAuthSourceUpdatePolicy()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AuthSourceUpdatePolicy, false)), addRequest.AuthSourceUpdatePolicy)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.RegistrationEnabled) {
		addRequest.RegistrationEnabled = plan.RegistrationEnabled.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.RegistrationConfig) {
		addRequest.RegistrationConfig = client.NewRegistrationConfigWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RegistrationConfig, true)), addRequest.RegistrationConfig)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ProfileConfig) {
		addRequest.ProfileConfig = client.NewProfileConfigWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ProfileConfig, false)), addRequest.ProfileConfig)
		fmt.Println(err)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.FieldConfig) {
		addRequest.FieldConfig = client.NewFieldConfig()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.FieldConfig, true)), addRequest.FieldConfig)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.EmailVerificationConfig) {
		addRequest.EmailVerificationConfig = client.NewEmailVerificationConfigWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.EmailVerificationConfig, true)), addRequest.EmailVerificationConfig)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsNonEmptyObj(plan.DataStoreConfig) {
		addRequest.DataStoreConfig = client.NewLdapDataStoreConfigWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.DataStoreConfig, false)), addRequest.DataStoreConfig)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ProfileEnabled) {
		addRequest.ProfileEnabled = plan.ProfileEnabled.ValueBoolPointer()
	}
	return nil
}

// Metadata returns the resource type name.
func (r *localIdentityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_identity_profile"
}

func (r *localIdentityProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *localIdentityProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *localIdentityProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var respDiags diag.Diagnostics

	// If the plan is null, this must be a destroy. Just exit early
	if plan == nil {
		return
	}

	// Set defaults that apply only when at least one of registration_enabled or profile_enabled is set to true
	if plan.RegistrationEnabled.ValueBool() || plan.ProfileEnabled.ValueBool() {
		if plan.AuthSourceUpdatePolicy.IsUnknown() {
			plan.AuthSourceUpdatePolicy = authSourceUpdatePolicyDefault
		}
		if plan.EmailVerificationConfig.IsUnknown() {
			plan.EmailVerificationConfig = emailVerificationConfigDefault
		}
	} else {
		if plan.AuthSourceUpdatePolicy.IsUnknown() {
			plan.AuthSourceUpdatePolicy = types.ObjectNull(authSourceUpdatePolicyAttrTypes)
		}
		if plan.EmailVerificationConfig.IsUnknown() {
			plan.EmailVerificationConfig = types.ObjectNull(emailVerificationConfigAttrTypes)
		}
	}

	// Email verification config inner defaults
	if internaltypes.IsDefined(plan.EmailVerificationConfig) {
		emailVerificationAttributes := plan.EmailVerificationConfig.Attributes()
		// Some defaults in the email verification config only apply if email verification is enabled
		emailVerificationEnabled := emailVerificationAttributes["email_verification_enabled"].(types.Bool).ValueBool()
		// Some attributes only apply for certain email verification types
		emailVerificationType := emailVerificationAttributes["email_verification_type"].(types.String)
		isOTP := emailVerificationType.ValueString() == "OTP"
		isOTL := emailVerificationType.ValueString() == "OTL"

		// Set unknown email verification attributes to defaults
		if emailVerificationAttributes["verify_email_template_name"].IsUnknown() {
			if emailVerificationEnabled {
				emailVerificationAttributes["verify_email_template_name"] = types.StringValue("message-template-email-ownership-verification.html")
			} else {
				emailVerificationAttributes["verify_email_template_name"] = types.StringNull()
			}
		}

		if emailVerificationAttributes["email_verification_success_template_name"].IsUnknown() {
			if emailVerificationEnabled {
				emailVerificationAttributes["email_verification_success_template_name"] = types.StringValue("local.identity.email.verification.success.html")
			} else {
				emailVerificationAttributes["email_verification_success_template_name"] = types.StringNull()
			}
		}

		if emailVerificationAttributes["email_verification_error_template_name"].IsUnknown() {
			if emailVerificationEnabled {
				emailVerificationAttributes["email_verification_error_template_name"] = types.StringValue("local.identity.email.verification.error.html")
			} else {
				emailVerificationAttributes["email_verification_error_template_name"] = types.StringNull()
			}
		}

		if emailVerificationAttributes["require_verified_email"].IsUnknown() {
			if emailVerificationEnabled {
				emailVerificationAttributes["require_verified_email"] = types.BoolValue(false)
			} else {
				emailVerificationAttributes["require_verified_email"] = types.BoolNull()
			}
		}

		// Attributes with defaults only when OTP is enabled
		if emailVerificationAttributes["allowed_otp_character_set"].IsUnknown() {
			if emailVerificationEnabled && isOTP {
				emailVerificationAttributes["allowed_otp_character_set"] = types.StringValue("23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz")
			} else {
				emailVerificationAttributes["allowed_otp_character_set"] = types.StringNull()
			}
		}

		if emailVerificationAttributes["otp_time_to_live"].IsUnknown() {
			if emailVerificationEnabled && isOTP {
				emailVerificationAttributes["otp_time_to_live"] = types.Int64Value(15)
			} else {
				emailVerificationAttributes["otp_time_to_live"] = types.Int64Null()
			}
		}

		if emailVerificationAttributes["otp_length"].IsUnknown() {
			if emailVerificationEnabled && isOTP {
				emailVerificationAttributes["otp_length"] = types.Int64Value(8)
			} else {
				emailVerificationAttributes["otp_length"] = types.Int64Null()
			}
		}

		if emailVerificationAttributes["otp_retry_attempts"].IsUnknown() {
			if emailVerificationEnabled && isOTP {
				emailVerificationAttributes["otp_retry_attempts"] = types.Int64Value(3)
			} else {
				emailVerificationAttributes["otp_retry_attempts"] = types.Int64Null()
			}
		}

		if emailVerificationAttributes["email_verification_otp_template_name"].IsUnknown() {
			if emailVerificationEnabled && isOTP {
				emailVerificationAttributes["email_verification_otp_template_name"] = types.StringValue("local.identity.email.verification.otp.html")
			} else {
				emailVerificationAttributes["email_verification_otp_template_name"] = types.StringNull()
			}
		}

		// Attributes with defaults only when OTL is enabled
		if emailVerificationAttributes["email_verification_sent_template_name"].IsUnknown() {
			if emailVerificationEnabled && isOTL {
				emailVerificationAttributes["email_verification_sent_template_name"] = types.StringValue("local.identity.email.verification.sent.html")
			} else {
				emailVerificationAttributes["email_verification_sent_template_name"] = types.StringNull()
			}
		}

		if emailVerificationAttributes["otl_time_to_live"].IsUnknown() {
			if emailVerificationEnabled && isOTL {
				emailVerificationAttributes["otl_time_to_live"] = types.Int64Value(1440)
			} else {
				emailVerificationAttributes["otl_time_to_live"] = types.Int64Null()
			}
		}

		if emailVerificationAttributes["require_verified_email_template_name"].IsUnknown() {
			if emailVerificationEnabled && isOTL {
				emailVerificationAttributes["require_verified_email_template_name"] = types.StringValue("local.identity.email.verification.required.html")
			} else {
				emailVerificationAttributes["require_verified_email_template_name"] = types.StringNull()
			}
		}

		// Update the object with the set defaults
		plan.EmailVerificationConfig, respDiags = types.ObjectValue(emailVerificationConfigAttrTypes, emailVerificationAttributes)
		resp.Diagnostics.Append(respDiags...)
	}

	// Some default for fields attributes depend on the field type
	if internaltypes.IsDefined(plan.FieldConfig) {
		fieldsList := plan.FieldConfig.Attributes()["fields"].(types.Set)
		if internaltypes.IsDefined(fieldsList) {
			fields := fieldsList.Elements()
			fieldsWithDefaults := []attr.Value{}
			for _, field := range fields {
				fieldObj := field.(types.Object)
				fieldAttrs := fieldObj.Attributes()
				if !fieldAttrs["attributes"].IsUnknown() {
					fieldsWithDefaults = append(fieldsWithDefaults, fieldObj)
					continue
				}

				defaultAttrs := map[string]attr.Value{}
				fieldType := fieldAttrs["type"].(types.String).ValueString()
				switch fieldType {
				// The other type defaults aren't implemented
				case "HIDDEN":
					defaultAttrs["Unique ID Field"] = types.BoolValue(false)
					defaultAttrs["Mask Log Values"] = types.BoolValue(false)
				case "CHECKBOX":
					defaultAttrs["Mask Log Values"] = types.BoolValue(false)
					defaultAttrs["Must Be Checked"] = types.BoolValue(false)
					defaultAttrs["Read-Only"] = types.BoolValue(false)
				case "DATE":
					defaultAttrs["Mask Log Values"] = types.BoolValue(false)
					defaultAttrs["Read-Only"] = types.BoolValue(false)
					defaultAttrs["Required"] = types.BoolValue(false)
				case "EMAIL":
					fallthrough
				case "PHONE":
					fallthrough
				case "TEXT":
					defaultAttrs["Mask Log Values"] = types.BoolValue(false)
					defaultAttrs["Read-Only"] = types.BoolValue(false)
					defaultAttrs["Required"] = types.BoolValue(false)
					defaultAttrs["Unique ID Field"] = types.BoolValue(false)
				}

				if len(defaultAttrs) > 0 {
					fieldAttrs["attributes"], respDiags = types.MapValue(types.BoolType, defaultAttrs)
					resp.Diagnostics.Append(respDiags...)
				}

				fieldWithDefaults, respDiags := types.ObjectValue(fieldObj.AttributeTypes(ctx), fieldAttrs)
				resp.Diagnostics.Append(respDiags...)
				fieldsWithDefaults = append(fieldsWithDefaults, fieldWithDefaults)
			}
			// Update the Field config with any defaults that were set
			fieldsAttrs := plan.FieldConfig.Attributes()
			fieldsAttrs["fields"], respDiags = types.SetValue(fieldsList.ElementType(ctx), fieldsWithDefaults)
			resp.Diagnostics.Append(respDiags...)
			plan.FieldConfig, respDiags = types.ObjectValue(plan.FieldConfig.AttributeTypes(ctx), fieldsAttrs)
			resp.Diagnostics.Append(respDiags...)
		}
	}

	resp.Plan.Set(ctx, plan)
}

func (r *localIdentityProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model localIdentityProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if internaltypes.IsDefined(model.EmailVerificationConfig) {
		emailVerificationConfig := model.EmailVerificationConfig.Attributes()
		emailVerificationType := model.EmailVerificationConfig.Attributes()["email_verification_type"].(basetypes.StringValue).ValueString()
		// Validates Email Verification type for Email Configuration
		switch emailVerificationType {
		case "OTP":
			if internaltypes.IsDefined(emailVerificationConfig["otl_time_to_live"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otl_time_to_live attribute is not allowed when email_verification_type is OTP."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["require_verified_email_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("require_verified_email_template_name is not allowed when email verification or require_verified_email is disabled or when email_verification_type is OTP."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_sent_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_sent_template_name is not allowed when email verification or require_verified_email is disabled or when email_verification_type is OTP."))
			}
		case "OTL":
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email_verification_type is OTL."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_retry_attempts"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_retry_attempts attribute is not allowed when email_verification_type is OTL."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["allowed_otp_character_set"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("allowed_otp_character_set attribute is not allowed when email_verification_type is OTL."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email_verification_type is OTL."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_otp_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_otp_template_name attribute is not allowed when email_verification_type is OTL."))
			}
		default:
			if internaltypes.IsDefined(emailVerificationConfig["otl_time_to_live"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otl_time_to_live attribute is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["require_verified_email_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("require_verified_email_template_name is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_sent_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_sent_template_name is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_retry_attempts"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_retry_attempts attribute is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["allowed_otp_character_set"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("allowed_otp_character_set attribute is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email verification is disabled."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_otp_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_otp_template_name attribute is not allowed when  email verification is disabled."))
			}
		}
		// If email verification is enabled, some fields become required
		verificationEnabled := emailVerificationConfig["email_verification_enabled"].(types.Bool)
		if verificationEnabled.ValueBool() {
			fieldForEmailToVerify := emailVerificationConfig["field_for_email_to_verify"].(types.String)
			if fieldForEmailToVerify.IsNull() || (internaltypes.IsDefined(fieldForEmailToVerify) && fieldForEmailToVerify.ValueString() == "") {
				resp.Diagnostics.AddError("Missing Required Attribute", fmt.Sprintln("field_for_email_to_verify is required when email_verification_enabled is set to true"))
			}
			fieldStoringVerificationStatus := emailVerificationConfig["field_storing_verification_status"].(types.String)
			if fieldStoringVerificationStatus.IsNull() || (internaltypes.IsDefined(fieldStoringVerificationStatus) && fieldStoringVerificationStatus.ValueString() == "") {
				resp.Diagnostics.AddError("Missing Required Attribute", fmt.Sprintln("field_storing_verification_status is required when email_verification_enabled is set to true"))
			}
			if emailVerificationConfig["notification_publisher_ref"].IsNull() {
				resp.Diagnostics.AddError("Missing Required Attribute", fmt.Sprintln("notification_publisher_ref is required when email_verification_enabled is set to true"))
			}
		}
	}
	if !model.ProfileEnabled.ValueBool() && !model.RegistrationEnabled.ValueBool() {
		if internaltypes.IsDefined(model.EmailVerificationConfig) || internaltypes.IsDefined(model.DataStoreConfig) || internaltypes.IsDefined(model.FieldConfig) || internaltypes.IsDefined(model.RegistrationConfig) || internaltypes.IsDefined(model.ProfileConfig) {
			resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email, data_store_config, field Config, registration_config and profile_config are not allowed when registration and profile are disabled."))
		}
		if internaltypes.IsDefined(model.AuthSourceUpdatePolicy) {
			resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("auth_source_update_policy is not allowed when registration and profile are disabled."))
		}
	} else {
		if model.ProfileEnabled.ValueBool() || model.RegistrationEnabled.ValueBool() {
			if model.FieldConfig.IsNull() {
				resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("field_config is required when profile or registration is enabled."))
			}
		}
		if model.ProfileEnabled.ValueBool() && model.ProfileConfig.IsNull() {
			resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_config is required when profile is enabled."))
		}
		if model.RegistrationEnabled.ValueBool() && model.RegistrationConfig.IsNull() {
			resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_config is required when registration is enabled."))
		}
		if !model.RegistrationEnabled.ValueBool() {
			if internaltypes.IsDefined(model.FieldConfig.Attributes()["fields"]) {
				fieldObj := model.FieldConfig.Attributes()["fields"].(basetypes.SetValue)
				fieldElems := fieldObj.Elements()
				for _, fieldElem := range fieldElems {
					fieldElemAttrs := fieldElem.(types.Object)
					registrationPageField := fieldElemAttrs.Attributes()["registration_page_field"].(basetypes.BoolValue)
					if registrationPageField.ValueBool() {
						resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_page_field option for the fields attribute should not be set to 'true' when registration is disabled."))
					}
					profilePageField := fieldElemAttrs.Attributes()["profile_page_field"].(basetypes.BoolValue)
					if internaltypes.IsDefined(profilePageField) && !profilePageField.ValueBool() {
						resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_page_field option is required to be set to 'true' for the fields attribute when profile management is the only option enabled."))
					}
				}
			}
			if internaltypes.IsDefined(model.RegistrationConfig) {
				resp.Diagnostics.AddError("Invalid Attribute!", fmt.Sprintln("registration_config is not allowed when registration is not enabled."))
			}
		}
		if !model.ProfileEnabled.ValueBool() {
			if internaltypes.IsDefined(model.FieldConfig.Attributes()["fields"]) {
				fieldObj := model.FieldConfig.Attributes()["fields"].(basetypes.SetValue)
				fieldElems := fieldObj.Elements()
				for _, fieldElem := range fieldElems {
					fieldElemAttrs := fieldElem.(types.Object)
					profilePagefield := fieldElemAttrs.Attributes()["profile_page_field"].(basetypes.BoolValue)
					if profilePagefield.ValueBool() {
						resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_page_field option for the fields attribute should not be set to 'true' when profile is disabled."))
					}
					registrationPageField := fieldElemAttrs.Attributes()["registration_page_field"].(basetypes.BoolValue)
					if (internaltypes.IsDefined(registrationPageField)) && (!registrationPageField.ValueBool()) {
						resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_page_field option is required to be set to 'true' for the fields attribute when registration is the only option enabled."))
					}
				}
			}
			if internaltypes.IsDefined(model.ProfileConfig) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("profile_config is not allowed when profile is not enabled."))
			}
		}
	}
	if internaltypes.IsDefined(model.RegistrationConfig) {
		captchaEnabled := model.RegistrationConfig.Attributes()["captcha_enabled"].(types.Bool)
		captchaProviderRef := model.RegistrationConfig.Attributes()["captcha_provider_ref"].(types.Object)
		if (captchaEnabled.ValueBool() && captchaProviderRef.IsNull()) || (internaltypes.IsDefined(captchaEnabled) && !captchaEnabled.ValueBool() && internaltypes.IsDefined(captchaProviderRef)) {
			resp.Diagnostics.AddError("Invalid registration captcha settings",
				"If registration_config.captcha_enabled is set to true, then registration_config.captcha_provider_ref must be configured. If registration_config.captcha_enabled is false, then registration_config.captcha_provider_ref must not be configured.")
		}
	}

	if internaltypes.IsDefined(model.FieldConfig) && internaltypes.IsDefined(model.FieldConfig.Attributes()["fields"]) {
		fieldConfigFields := model.FieldConfig.Attributes()["fields"].(types.Set).Elements()
		if len(fieldConfigFields) > 0 {
			for _, field := range fieldConfigFields {
				fieldAttrs := field.(types.Object).Attributes()
				fieldType := fieldAttrs["type"].(types.String).ValueString()
				hasDefault := internaltypes.IsDefined(fieldAttrs["default_value"])
				hasOptions := internaltypes.IsDefined(fieldAttrs["options"])

				switch fieldType {
				case "CHECKBOX":
					// check to make sure options isn't set
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"CHECKBOX\" field type.")
					}
				case "CHECKBOX_GROUP":
					// check to make sure default_value isn't set
					if hasDefault {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"default_value\" is not applicable for the \"CHECKBOX_GROUP\" field type.")
					}
				case "DATE":
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"DATE\" field type.")
					}
				// DROP_DOWN has all options, leaving this here in case of future use
				// case "DROP_DOWN":
				case "EMAIL":
					if hasDefault {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"default_value\" is not applicable for the \"EMAIL\" field type.")
					}
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"EMAIL\" field type.")
					}
				case "HIDDEN":
					if hasDefault {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"default_value\" is not applicable for the \"HIDDEN\" field type.")
					}
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"HIDDEN\" field type.")
					}
				case "PHONE":
					if hasDefault {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"default_value\" is not applicable for the \"PHONE\" field type.")
					}
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"PHONE\" field type.")
					}
				case "TEXT":
					if hasOptions {
						resp.Diagnostics.AddError("Invalid attribute combination",
							"\"options\" is not applicable for the \"TEXT\" field type.")
					}
				}
			}
		}
	}
}

func readLocalIdentityProfileResponse(ctx context.Context, r *client.LocalIdentityProfile, state *localIdentityProfileModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringPointerValue(r.Id)
	state.ProfileId = types.StringPointerValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.ApcId, respDiags = resourcelink.ToState(ctx, &r.ApcId)
	diags.Append(respDiags...)

	// auth source update policy
	authSourceUpdatePolicy := r.AuthSourceUpdatePolicy
	state.AuthSourceUpdatePolicy, respDiags = types.ObjectValueFrom(ctx, authSourceUpdatePolicyAttrTypes, authSourceUpdatePolicy)
	diags.Append(respDiags...)

	// auth sources
	authSources := r.GetAuthSources()
	var authSourcesSliceAttrVal = []attr.Value{}
	authSourcesSliceType := types.ObjectType{AttrTypes: authSourcesAttrTypes}
	for i := 0; i < len(authSources); i++ {
		authSourcesAttrValues := map[string]attr.Value{
			"id":     types.StringPointerValue(authSources[i].Id),
			"source": types.StringPointerValue(authSources[i].Source),
		}
		authSourcesObj, respDiags := types.ObjectValue(authSourcesAttrTypes, authSourcesAttrValues)
		diags.Append(respDiags...)
		authSourcesSliceAttrVal = append(authSourcesSliceAttrVal, authSourcesObj)
	}
	state.AuthSources, respDiags = types.SetValue(authSourcesSliceType, authSourcesSliceAttrVal)
	diags.Append(respDiags...)

	state.RegistrationConfig, respDiags = types.ObjectValueFrom(ctx, registrationConfigAttrTypes, r.RegistrationConfig)
	diags.Append(respDiags...)

	state.RegistrationEnabled = types.BoolValue(r.GetRegistrationEnabled())

	state.ProfileConfig, respDiags = types.ObjectValueFrom(ctx, profileConfigAttrTypes, r.ProfileConfig)
	diags.Append(respDiags...)

	// field config
	state.FieldConfig, respDiags = types.ObjectValueFrom(ctx, fieldConfigAttrTypes, r.FieldConfig)
	diags.Append(respDiags...)

	// email verification config
	state.EmailVerificationConfig, respDiags = types.ObjectValueFrom(ctx, emailVerificationConfigAttrTypes, r.EmailVerificationConfig)
	diags.Append(respDiags...)

	//  data store config
	state.DataStoreConfig, respDiags = types.ObjectValueFrom(ctx, dsConfigAttrTypes, r.DataStoreConfig)
	diags.Append(respDiags...)

	state.ProfileEnabled = types.BoolPointerValue(r.ProfileEnabled)
	return diags
}

func (r *localIdentityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan localIdentityProfileModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apcResourceLink, err := resourcelink.ClientStruct(plan.ApcId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add apc id to add request for a local identity profile", err.Error())
		return
	}
	createLocalIdentityProfiles := client.NewLocalIdentityProfile(plan.Name.ValueString(), *apcResourceLink)
	err = addOptionalLocalIdentityProfileFields(ctx, createLocalIdentityProfiles, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a local identity profile", err.Error())
		return
	}
	apiCreateLocalIdentityProfiles := r.apiClient.LocalIdentityIdentityProfilesAPI.CreateIdentityProfile(config.AuthContext(ctx, r.providerConfig))
	apiCreateLocalIdentityProfiles = apiCreateLocalIdentityProfiles.Body(*createLocalIdentityProfiles)
	localIdentityProfilesResponse, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.CreateIdentityProfileExecute(apiCreateLocalIdentityProfiles)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the local identity profiles", err, httpResp)
		return
	}

	// Read the response into the state
	var state localIdentityProfileModel

	diags = readLocalIdentityProfileResponse(ctx, localIdentityProfilesResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *localIdentityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state localIdentityProfileModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadLocalIdentityProfiles, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(config.AuthContext(ctx, r.providerConfig), state.ProfileId.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Local Identity Profile", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the local identity profile", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readLocalIdentityProfileResponse(ctx, apiReadLocalIdentityProfiles, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *localIdentityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan localIdentityProfileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateLocalIdentityProfiles := r.apiClient.LocalIdentityIdentityProfilesAPI.UpdateIdentityProfile(config.AuthContext(ctx, r.providerConfig), plan.ProfileId.ValueString())
	apcResourceLink, err := resourcelink.ClientStruct(plan.ApcId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add apc id to add request for the local identity profile", err.Error())
		return
	}
	createUpdateRequest := client.NewLocalIdentityProfile(plan.Name.ValueString(), *apcResourceLink)
	err = addOptionalLocalIdentityProfileFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a local identity profile", err.Error())
		return
	}
	updateLocalIdentityProfiles = updateLocalIdentityProfiles.Body(*createUpdateRequest)
	updateLocalIdentityProfilesResponse, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.UpdateIdentityProfileExecute(updateLocalIdentityProfiles)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating a local identity profile", err, httpResp)
		return
	}

	// Read the response
	diags = readLocalIdentityProfileResponse(ctx, updateLocalIdentityProfilesResponse, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *localIdentityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localIdentityProfileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(config.AuthContext(ctx, r.providerConfig), state.ProfileId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the local identity profile", err, httpResp)
	}

}

func (r *localIdentityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("profile_id"), req, resp)
}
