package localidentity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localIdentityIdentityProfilesResource{}
	_ resource.ResourceWithConfigure   = &localIdentityIdentityProfilesResource{}
	_ resource.ResourceWithImportState = &localIdentityIdentityProfilesResource{}
)

var (
	authSourceUpdatePolicyAttrTypes = map[string]attr.Type{
		"store_attributes":  basetypes.BoolType{},
		"retain_attributes": basetypes.BoolType{},
		"update_attributes": basetypes.BoolType{},
		"update_interval":   basetypes.Int64Type{},
	}

	authSourcesAttrTypes = map[string]attr.Type{
		"id":     basetypes.StringType{},
		"source": basetypes.StringType{},
	}

	registrationConfigAttrTypes = map[string]attr.Type{
		"captcha_enabled":                         basetypes.BoolType{},
		"captcha_provider_ref":                    basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"template_name":                           basetypes.StringType{},
		"create_authn_session_after_registration": basetypes.BoolType{},
		"username_field":                          basetypes.StringType{},
		"this_is_my_device_enabled":               basetypes.BoolType{},
		"registration_workflow":                   basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"execute_workflow":                        basetypes.StringType{},
	}

	profileConfigAttrTypes = map[string]attr.Type{
		"delete_identity_enabled": basetypes.BoolType{},
		"template_name":           basetypes.StringType{},
	}

	fieldItemAttrTypes = map[string]attr.Type{
		"type":                    basetypes.StringType{},
		"id":                      basetypes.StringType{},
		"label":                   basetypes.StringType{},
		"registration_page_field": basetypes.BoolType{},
		"profile_page_field":      basetypes.BoolType{},
		"attributes":              basetypes.MapType{ElemType: basetypes.BoolType{}},
	}

	fieldConfigAttrTypes = map[string]attr.Type{
		"fields":                        basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldItemAttrTypes}},
		"strip_space_from_unique_field": basetypes.BoolType{},
	}

	emailVerificationConfigAttrTypes = map[string]attr.Type{
		"email_verification_enabled":               basetypes.BoolType{},
		"verify_email_template_name":               basetypes.StringType{},
		"email_verification_sent_template_name":    basetypes.StringType{},
		"email_verification_success_template_name": basetypes.StringType{},
		"email_verification_error_template_name":   basetypes.StringType{},
		"email_verification_type":                  basetypes.StringType{},
		"otp_length":                               basetypes.Int64Type{},
		"otp_retry_attempts":                       basetypes.Int64Type{},
		"allowed_otp_character_set":                basetypes.StringType{},
		"otp_time_to_live":                         basetypes.Int64Type{},
		"email_verification_otp_template_name":     basetypes.StringType{},
		"otl_time_to_live":                         basetypes.Int64Type{},
		"field_for_email_to_verify":                basetypes.StringType{},
		"field_storing_verification_status":        basetypes.StringType{},
		"notification_publisher_ref":               basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"require_verified_email":                   basetypes.BoolType{},
		"require_verified_email_template_name":     basetypes.StringType{},
	}

	dsConfigAttrTypes = map[string]attr.Type{
		"base_dn":                  basetypes.StringType{},
		"type":                     basetypes.StringType{},
		"data_store_ref":           basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"data_store_mapping":       basetypes.MapType{ElemType: types.ObjectType{AttrTypes: dsMappingAttrtypes}},
		"create_pattern":           basetypes.StringType{},
		"object_class":             basetypes.StringType{},
		"auxiliary_object_classes": basetypes.SetType{ElemType: basetypes.StringType{}},
	}

	dsMappingAttrtypes = map[string]attr.Type{
		"type":     basetypes.StringType{},
		"name":     basetypes.StringType{},
		"metadata": basetypes.MapType{ElemType: basetypes.StringType{}},
	}
)

// LocalIdentityIdentityProfilesResource is a helper function to simplify the provider implementation.
func LocalIdentityIdentityProfilesResource() resource.Resource {

	return &localIdentityIdentityProfilesResource{}
}

// localIdentityIdentityProfilesResource is the resource implementation.
type localIdentityIdentityProfilesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type localIdentityIdentityProfilesResourceModel struct {
	Id       types.String `tfsdk:"id"`
	CustomId types.String `tfsdk:"custom_id"`
	Name     types.String `tfsdk:"name"`
	ApcId    types.Object `tfsdk:"apc_id"`

	AuthSources            types.List   `tfsdk:"auth_sources"`
	AuthSourceUpdatePolicy types.Object `tfsdk:"auth_source_update_policy"`
	RegistrationEnabled    types.Bool   `tfsdk:"registration_enabled"`
	RegistrationConfig     types.Object `tfsdk:"registration_config"`

	ProfileConfig           types.Object `tfsdk:"profile_config"`
	FieldConfig             types.Object `tfsdk:"field_config"`
	EmailVerificationConfig types.Object `tfsdk:"email_verification_config"`

	DataStoreConfig types.Object `tfsdk:"data_store_config"`
	ProfileEnabled  types.Bool   `tfsdk:"profile_enabled"`
}

// GetSchema defines the schema for the resource.
func (r *localIdentityIdentityProfilesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Local Identity Identity Profiles",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The local identity profile name. Name is unique.",
				Required:    true,
			},
			"apc_id": schema.SingleNestedAttribute{
				Description: "The reference to the authentication policy contract to use for this local identity profile.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"auth_sources": schema.ListNestedAttribute{
				Description: "The local identity authentication sources. Sources are unique.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The persistent, unique ID for the local identity authentication source. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
							Computed:    true,
							Optional:    true,
						},
						"source": schema.StringAttribute{
							Description: "The local identity authentication source. Source is unique.",
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
			"auth_source_update_policy": schema.SingleNestedAttribute{
				Description: "The attribute update policy for authentication sources.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"store_attributes": schema.BoolAttribute{
						Description: "Whether or not to store attributes that came from authentication sources.",
						Optional:    true,
					},
					"retain_attributes": schema.BoolAttribute{
						Description: "Whether or not to keep attributes after user disconnects.",
						Optional:    true,
					},
					"update_attributes": schema.BoolAttribute{
						Description: "Whether or not to update attributes when users authenticate.",
						Optional:    true,
					},
					"update_interval": schema.Int64Attribute{
						Description: "The minimum number of days between updates.",
						Optional:    true,
					},
				},
			},
			"registration_enabled": schema.BoolAttribute{
				Description: "Whether the registration configuration is enabled or not.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"registration_config": schema.SingleNestedAttribute{
				Description: "The local identity profile registration configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"captcha_enabled": schema.BoolAttribute{
						Description: "Whether CAPTCHA is enabled or not in the registration configuration.",
						Computed:    true,
						Optional:    true,
					},
					"captcha_provider_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated CAPTCHA provider.",
						Computed:    true,
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for the registration configuration.",
						Required:    true,
					},
					"create_authn_session_after_registration": schema.BoolAttribute{
						Description: "Whether to create an Authentication Session when registering a local account. Default is true.",
						Computed:    true,
						Optional:    true,
					},
					"username_field": schema.StringAttribute{
						Description: "When creating an Authentication Session after registering a local account, PingFederate will pass the Unique ID field's value as the username. If the Unique ID value is not the username, then override which field's value will be used as the username.",
						Computed:    true,
						Optional:    true,
					},
					"this_is_my_device_enabled": schema.BoolAttribute{
						Description: "Allows users to indicate whether their device is shared or private. In this mode, PingFederate Authentication Sessions will not be stored unless the user indicates the device is private.",
						Computed:    true,
						Optional:    true,
					},
					"registration_workflow": schema.SingleNestedAttribute{
						Description: "The policy fragment to be executed as part of the registration workflow.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"execute_workflow": schema.StringAttribute{
						Description: "This setting indicates whether PingFederate should execute the workflow before or after account creation. The default is to run the registration workflow after account creation.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"BEFORE_ACCOUNT_CREATION", "AFTER_ACCOUNT_CREATION"}...),
						},
					},
				},
			},
			"profile_config": schema.SingleNestedAttribute{
				Description: "The local identity profile management configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"delete_identity_enabled": schema.BoolAttribute{
						Description: "Whether the end user is allowed to use delete functionality.",
						Computed:    true,
						Optional:    true,
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for end-user profile management.",
						Required:    true,
					},
				},
			},
			"field_config": schema.SingleNestedAttribute{
				Description: "The local identity profile field configuration.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"fields": schema.ListNestedAttribute{
						Description: "The field configuration for the local identity profile.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of the local identity field.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"CHECKBOX", "CHECKBOX_GROUP", "DATE", "DROP_DOWN", "EMAIL", "PHONE", "TEXT", "HIDDEN"}...),
									},
								},
								"id": schema.StringAttribute{
									Description: "Id of the local identity field.",
									Required:    true,
								},
								"label": schema.StringAttribute{
									Description: "Label of the local identity field.",
									Required:    true,
								},
								"registration_page_field": schema.BoolAttribute{
									Description: "Whether this is a registration page field or not.",
									Optional:    true,
								},
								"profile_page_field": schema.BoolAttribute{
									Description: "Whether this is a profile page field or not.",
									Optional:    true,
								},
								"attributes": schema.MapAttribute{
									Description: "Attributes of the local identity field.",
									Computed:    true,
									Optional:    true,
									ElementType: types.BoolType,
								},
							},
						},
					},
					"strip_space_from_unique_field": schema.BoolAttribute{
						Description: "Strip leading/trailing spaces from unique ID field. Default is true.",
						Computed:    true,
						Optional:    true,
					},
				},
			},
			"email_verification_config": schema.SingleNestedAttribute{
				Description: "The local identity email verification configuration.",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"email_verification_enabled": schema.BoolAttribute{
						Description: "Whether the email ownership verification is enabled.",
						Computed:    true,
						Optional:    true,
					},
					"verify_email_template_name": schema.StringAttribute{
						Description: "The template name for verify email. The default is message-template-email-ownership-verification.html.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_sent_template_name": schema.StringAttribute{
						Description: "The template name for email verification sent. The default is local.identity.email.verification.sent.html. Note:Only applicable if EmailVerificationType is OTL.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_success_template_name": schema.StringAttribute{
						Description: "The template name for email verification success. The default is local.identity.email.verification.success.html.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_error_template_name": schema.StringAttribute{
						Description: "The template name for email verification error. The default is local.identity.email.verification.error.html.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_type": schema.StringAttribute{
						Description: "Email Verification Type.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"OTP", "OTL"}...),
						},
					},
					"otp_length": schema.Int64Attribute{
						Description: "The OTP length generated for email verification. The default is 8. Note: Only applicable if EmailVerificationType is OTP.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(5, 100),
						},
					},
					"otp_retry_attempts": schema.Int64Attribute{
						Description: "The number of OTP retry attempts for email verification. The default is 3. Note: Only applicable if EmailVerificationType is OTP.",
						Optional:    true,
					},
					"allowed_otp_character_set": schema.StringAttribute{
						Description: "The allowed character set used to generate the OTP. The default is 23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz. Note: Only applicable if EmailVerificationType is OTP.",
						Optional:    true,
						Computed:    true,
					},
					"otp_time_to_live": schema.Int64Attribute{
						Description: "Field used OTP time to live. The default is 15. Note: Only applicable if EmailVerificationType is OTP.",
						Computed:    true,
						Optional:    true,
					},
					"email_verification_otp_template_name": schema.StringAttribute{
						Description: "The template name for email verification OTP verification. The default is local.identity.email.verification.otp.html. Note: Only applicable if EmailVerificationType is OTP.",
						Optional:    true,
						Computed:    true,
					},
					"otl_time_to_live": schema.Int64Attribute{
						Description: "Field used OTL time to live. The default is 1440. Note: Only applicable if EmailVerificationType is OTL.",
						Computed:    true,
						Optional:    true,
					},
					"field_for_email_to_verify": schema.StringAttribute{
						Description: "Field used for email ownership verification. Note: Not required when emailVerificationEnabled is set to false.",
						Required:    true,
					},
					"field_storing_verification_status": schema.StringAttribute{
						Description: "Field used for storing email verification status. Note: Not required when emailVerificationEnabled is set to false.",
						Required:    true,
					},
					"notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher.",
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
					},
					"require_verified_email": schema.BoolAttribute{
						Description: "Whether the user must verify their email address before they can complete a single sign-on transaction. The default is false.",
						Computed:    true,
						Optional:    true,
					},
					"require_verified_email_template_name": schema.StringAttribute{
						Description: "The template to render when the user must verify their email address before they can complete a single sign-on transaction. The default is local.identity.email.verification.required.html. Note:Only applicable if EmailVerificationType is OTL and requireVerifiedEmail is true.",
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
					},
					"type": schema.StringAttribute{
						Description: "The data store config type.",
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
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"data_store_mapping": schema.MapNestedAttribute{
						Description: "The data store mapping.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The data store attribute type.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
									},
								},
								"name": schema.StringAttribute{
									Description: "The data store attribute name.",
									Required:    true,
								},
								"metadata": schema.MapAttribute{
									Description: "The data store attribute metadata.",
									Computed:    true,
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"create_pattern": schema.StringAttribute{
						Description: "The Relative DN Pattern that will be used to create objects in the directory.",
						Required:    true,
					},
					"object_class": schema.StringAttribute{
						Description: "The Object Class used by the new objects stored in the LDAP data store.",
						Required:    true,
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
				Description: "Whether the profile configuration is enabled or not.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, false, true,
		"The persistent, unique ID for the local identity profile. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.")
	resp.Schema = schema
}

func addOptionalLocalIdentityIdentityProfilesFields(ctx context.Context, addRequest *client.LocalIdentityProfile, plan localIdentityIdentityProfilesResourceModel) error {

	if internaltypes.IsDefined(plan.CustomId) {
		addRequest.Id = plan.CustomId.ValueStringPointer()
	}

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
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.FieldConfig, false)), addRequest.FieldConfig)
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
func (r *localIdentityIdentityProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_identity_identity_profile"
}

func (r *localIdentityIdentityProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *localIdentityIdentityProfilesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model localIdentityIdentityProfilesResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	// Validates Email Verification type for Email Configuration
	if internaltypes.IsDefined(model.EmailVerificationConfig) {
		emailVerificationConfig := model.EmailVerificationConfig.Attributes()
		emailVerificationType := model.EmailVerificationConfig.Attributes()["email_verification_type"].(basetypes.StringValue).ValueString()
		switch emailVerificationType {
		case "OTP":
			if internaltypes.IsDefined(emailVerificationConfig["otl_time_to_live"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otl_time_to_live attribute is not allowed when email_verification_type is OTP. Required attributes are otp_length, otp_retry_attempts and otp_time_to_live."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["require_verified_email_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("require_verified_email_template_name is not allowed when email verification or require_verified_email is disabled or when email_verification_type is OTP."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_sent_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_sent_template_name is not allowed when email verification or require_verified_email is disabled or when email_verification_type is OTP."))
			}
		case "OTL":
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email_verification_type is OTL. Required attribute: otl_time_to_live."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_retry_attempts"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_retry_attempts attribute is not allowed when email_verification_type is OTL. Required attribute: otl_time_to_live."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["allowed_otp_character_set"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("allowed_otp_character_set attribute is not allowed when email_verification_type is OTL. Required attribute: otl_time_to_live."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["otp_length"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("otp_length attribute is not allowed when email_verification_type is OTL. Required attribute: otl_time_to_live."))
			}
			if internaltypes.IsDefined(emailVerificationConfig["email_verification_otp_template_name"]) {
				resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email_verification_otp_template_name attribute is not allowed when email_verification_type is OTL. Required attribute: otl_time_to_live."))
			}
		}
	}
	if (!model.ProfileEnabled.ValueBool()) && (!model.RegistrationEnabled.ValueBool()) {
		if internaltypes.IsDefined(model.EmailVerificationConfig) || internaltypes.IsDefined(model.DataStoreConfig) || internaltypes.IsDefined(model.FieldConfig) || internaltypes.IsDefined(model.RegistrationConfig) || internaltypes.IsDefined(model.ProfileConfig) {
			resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("email, data_store_config, field Config, registration_config and profile_config are not allowed when registration and profile are disabled."))
		}
		if internaltypes.IsDefined(model.AuthSourceUpdatePolicy) {
			resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("auth_source_update_policy is not allowed when registration and profile are disabled."))
		}
	} else {
		if (model.ProfileEnabled.ValueBool()) && (model.RegistrationEnabled.ValueBool()) {
			if !model.ProfileEnabled.ValueBool() {
				if internaltypes.IsDefined(model.FieldConfig.Attributes()["fields"]) {
					fieldObj := model.FieldConfig.Attributes()["fields"].(basetypes.ListValue)
					fieldElems := fieldObj.Elements()
					for _, fieldElem := range fieldElems {
						fieldElemAttrs := fieldElem.(basetypes.ObjectValue)
						profilePagefield := fieldElemAttrs.Attributes()["profile_page_field"].(basetypes.BoolValue)
						if (internaltypes.IsDefined(profilePagefield)) && (profilePagefield.ValueBool()) {
							resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_page_field option for the fields attribute should not be set to 'true' when profile is disabled."))
						}
						registrationPageField := fieldElemAttrs.Attributes()["registration_page_field"].(basetypes.BoolValue)
						if (internaltypes.IsDefined(registrationPageField)) && (!registrationPageField.ValueBool()) {
							resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_page_field option is required to be set to 'true' for the fields attribute when registration is the only option enabled."))
						}
					}
				}
				if !model.ProfileConfig.IsNull() {
					resp.Diagnostics.AddError("Invalid Attribute Combination!", fmt.Sprintln("profile_config is not allowed when profile is not enabled."))
				}
			}
			if (model.ProfileEnabled.ValueBool()) && (model.ProfileConfig.IsNull()) {
				resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_config is required when profile is enabled."))
			}
			if !model.RegistrationEnabled.ValueBool() {
				if internaltypes.IsDefined(model.FieldConfig.Attributes()["fields"]) {
					fieldObj := model.FieldConfig.Attributes()["fields"].(basetypes.SetValue)
					fieldElems := fieldObj.Elements()
					for _, fieldElem := range fieldElems {
						fieldElemAttrs := fieldElem.(basetypes.ObjectValue)
						registrationPageField := fieldElemAttrs.Attributes()["registration_page_field"].(basetypes.BoolValue)
						if (internaltypes.IsDefined(registrationPageField)) && (registrationPageField.ValueBool()) {
							resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_page_field option for the fields attribute should not be set to 'true' when registration is disabled."))
						}
						profilePageField := fieldElemAttrs.Attributes()["profile_page_field"].(basetypes.BoolValue)
						if (internaltypes.IsDefined(profilePageField)) && (!profilePageField.ValueBool()) {
							resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("profile_page_field option is required to be set to 'true' for the fields attribute when profile management is the only option enabled."))
						}
					}
				}
				if !model.RegistrationConfig.IsNull() {
					resp.Diagnostics.AddError("Invalid Attribute!", fmt.Sprintln("registration_config is not allowed when registration is not enabled."))
				}
			}
			if (model.RegistrationEnabled.ValueBool()) && (model.RegistrationConfig.IsNull()) {
				resp.Diagnostics.AddError("Invalid Value for Attribute!", fmt.Sprintln("registration_config is required when registration is enabled."))
			}
		}
	}
}

func readLocalIdentityIdentityProfilesResponse(ctx context.Context, r *client.LocalIdentityProfile, state *localIdentityIdentityProfilesResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.CustomId = internaltypes.StringTypeOrNil(r.Id, false)
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
	state.AuthSources, respDiags = types.ListValue(authSourcesSliceType, authSourcesSliceAttrVal)
	diags.Append(respDiags...)

	registrationConfig := r.RegistrationConfig
	state.RegistrationConfig, respDiags = types.ObjectValueFrom(ctx, registrationConfigAttrTypes, registrationConfig)
	diags.Append(respDiags...)

	state.RegistrationEnabled = types.BoolValue(r.GetRegistrationEnabled())

	profileConfig := r.ProfileConfig
	state.ProfileConfig, respDiags = types.ObjectValueFrom(ctx, profileConfigAttrTypes, profileConfig)
	diags.Append(respDiags...)

	// field config
	fieldConfig := r.GetFieldConfig()
	fieldType := types.ObjectType{AttrTypes: fieldItemAttrTypes}
	fieldAttrsStruct := fieldConfig.GetFields()
	fieldAttrsState, respDiags := types.ListValueFrom(ctx, fieldType, fieldAttrsStruct)
	diags.Append(respDiags...)
	stripSpaceFromUniqueFieldState := types.BoolPointerValue(r.GetFieldConfig().StripSpaceFromUniqueField)
	fieldConfigAttrValues := map[string]attr.Value{
		"fields":                        fieldAttrsState,
		"strip_space_from_unique_field": stripSpaceFromUniqueFieldState,
	}
	state.FieldConfig, respDiags = types.ObjectValue(fieldConfigAttrTypes, fieldConfigAttrValues)
	diags.Append(respDiags...)

	emailVerificationConfig := r.EmailVerificationConfig
	state.EmailVerificationConfig, respDiags = types.ObjectValueFrom(ctx, emailVerificationConfigAttrTypes, emailVerificationConfig)
	diags.Append(respDiags...)

	//  data store config
	dsConfig := r.DataStoreConfig
	state.DataStoreConfig, respDiags = types.ObjectValueFrom(ctx, dsConfigAttrTypes, dsConfig)
	diags.Append(respDiags...)

	state.ProfileEnabled = types.BoolPointerValue(r.ProfileEnabled)
	return diags
}

func (r *localIdentityIdentityProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan localIdentityIdentityProfilesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apcResourceLink, err := resourcelink.ClientStruct(plan.ApcId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add apc id to add request for Local Identity Identity Profile", err.Error())
		return
	}
	createLocalIdentityIdentityProfiles := client.NewLocalIdentityProfile(plan.Name.ValueString(), *apcResourceLink)
	err = addOptionalLocalIdentityIdentityProfilesFields(ctx, createLocalIdentityIdentityProfiles, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for a Local Identity Identity Profile", err.Error())
		return
	}
	apiCreateLocalIdentityIdentityProfiles := r.apiClient.LocalIdentityIdentityProfilesAPI.CreateIdentityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateLocalIdentityIdentityProfiles = apiCreateLocalIdentityIdentityProfiles.Body(*createLocalIdentityIdentityProfiles)
	localIdentityIdentityProfilesResponse, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.CreateIdentityProfileExecute(apiCreateLocalIdentityIdentityProfiles)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the LocalIdentityIdentityProfiles", err, httpResp)
		return
	}

	// Read the response into the state
	var state localIdentityIdentityProfilesResourceModel

	diags = readLocalIdentityIdentityProfilesResponse(ctx, localIdentityIdentityProfilesResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *localIdentityIdentityProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state localIdentityIdentityProfilesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadLocalIdentityIdentityProfiles, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Local Identity Profile", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Identity Profile", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readLocalIdentityIdentityProfilesResponse(ctx, apiReadLocalIdentityIdentityProfiles, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *localIdentityIdentityProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan localIdentityIdentityProfilesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateLocalIdentityIdentityProfiles := r.apiClient.LocalIdentityIdentityProfilesAPI.UpdateIdentityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.CustomId.ValueString())
	apcResourceLink, err := resourcelink.ClientStruct(plan.ApcId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add apc id to add request for Local Identity Identity Profile", err.Error())
		return
	}
	createUpdateRequest := client.NewLocalIdentityProfile(plan.Name.ValueString(), *apcResourceLink)
	err = addOptionalLocalIdentityIdentityProfilesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Identity Identity Profile", err.Error())
		return
	}
	updateLocalIdentityIdentityProfiles = updateLocalIdentityIdentityProfiles.Body(*createUpdateRequest)
	updateLocalIdentityIdentityProfilesResponse, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.UpdateIdentityProfileExecute(updateLocalIdentityIdentityProfiles)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Local Identity Identity Profile", err, httpResp)
		return
	}

	// Read the response
	diags = readLocalIdentityIdentityProfilesResponse(ctx, updateLocalIdentityIdentityProfilesResponse, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *localIdentityIdentityProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localIdentityIdentityProfilesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting Local Identity Profile", err, httpResp)
	}

}

func (r *localIdentityIdentityProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
