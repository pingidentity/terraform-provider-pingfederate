package localidentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &localIdentityIdentityProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &localIdentityIdentityProfileDataSource{}
)

// Create a Administrative Account data source
func NewLocalIdentityIdentityProfileDataSource() datasource.DataSource {
	return &localIdentityIdentityProfileDataSource{}
}

// localIdentityIdentityProfileDataSource is the datasource implementation.
type localIdentityIdentityProfileDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type localIdentityIdentityProfileDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ApcId                   types.Object `tfsdk:"apc_id"`
	AuthSources             types.List   `tfsdk:"auth_sources"`
	AuthSourceUpdatePolicy  types.Object `tfsdk:"auth_source_update_policy"`
	RegistrationEnabled     types.Bool   `tfsdk:"registration_enabled"`
	RegistrationConfig      types.Object `tfsdk:"registration_config"`
	ProfileConfig           types.Object `tfsdk:"profile_config"`
	FieldConfig             types.Object `tfsdk:"field_config"`
	EmailVerificationConfig types.Object `tfsdk:"email_verification_config"`
	DataStoreConfig         types.Object `tfsdk:"data_store_config"`
	ProfileEnabled          types.Bool   `tfsdk:"profile_enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *localIdentityIdentityProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages Local Identity Identity Profiles",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The local identity profile name. Name is unique.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"apc_id": schema.SingleNestedAttribute{
				Description: "The reference to the authentication policy contract to use for this local identity profile.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"auth_sources": schema.ListNestedAttribute{
				Description: "The local identity authentication sources. Sources are unique.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The persistent, unique ID for the local identity authentication source. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"source": schema.StringAttribute{
							Description: "The local identity authentication source. Source is unique.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"auth_source_update_policy": schema.SingleNestedAttribute{
				Description: "The attribute update policy for authentication sources.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"store_attributes": schema.BoolAttribute{
						Description: "Whether or not to store attributes that came from authentication sources.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"retain_attributes": schema.BoolAttribute{
						Description: "Whether or not to keep attributes after user disconnects.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"update_attributes": schema.BoolAttribute{
						Description: "Whether or not to update attributes when users authenticate.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"update_interval": schema.Int64Attribute{
						Description: "The minimum number of days between updates.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"registration_enabled": schema.BoolAttribute{
				Description: "Whether the registration configuration is enabled or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"registration_config": schema.SingleNestedAttribute{
				Description: "The local identity profile registration configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"captcha_enabled": schema.BoolAttribute{
						Description: "Whether CAPTCHA is enabled or not in the registration configuration.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"captcha_provider_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated CAPTCHA provider.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for the registration configuration.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"create_authn_session_after_registration": schema.BoolAttribute{
						Description: "Whether to create an Authentication Session when registering a local account. Default is true.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"username_field": schema.StringAttribute{
						Description: "When creating an Authentication Session after registering a local account, PingFederate will pass the Unique ID field's value as the username. If the Unique ID value is not the username, then override which field's value will be used as the username.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"this_is_my_device_enabled": schema.BoolAttribute{
						Description: "Allows users to indicate whether their device is shared or private. In this mode, PingFederate Authentication Sessions will not be stored unless the user indicates the device is private.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"registration_workflow": schema.SingleNestedAttribute{
						Description: "The policy fragment to be executed as part of the registration workflow.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"execute_workflow": schema.StringAttribute{
						Description: "This setting indicates whether PingFederate should execute the workflow before or after account creation. The default is to run the registration workflow after account creation.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"profile_config": schema.SingleNestedAttribute{
				Description: "The local identity profile management configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"delete_identity_enabled": schema.BoolAttribute{
						Description: "Whether the end user is allowed to use delete functionality.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"template_name": schema.StringAttribute{
						Description: "The template name for end-user profile management.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"field_config": schema.SingleNestedAttribute{
				Description: "The local identity profile field configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"fields": schema.ListNestedAttribute{
						Description: "The field configuration for the local identity profile.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of the local identity field.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"id": schema.StringAttribute{
									Description: "Id of the local identity field.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "Label of the local identity field.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"registration_page_field": schema.BoolAttribute{
									Description: "Whether this is a registration page field or not.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"profile_page_field": schema.BoolAttribute{
									Description: "Whether this is a profile page field or not.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"attributes": schema.MapAttribute{
									Description: "Attributes of the local identity field.",
									Required:    false,
									Optional:    false,
									Computed:    true,
									ElementType: types.BoolType,
								},
							},
						},
					},
					"strip_space_from_unique_field": schema.BoolAttribute{
						Description: "Strip leading/trailing spaces from unique ID field. Default is true.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"email_verification_config": schema.SingleNestedAttribute{
				Description: "The local identity email verification configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"email_verification_enabled": schema.BoolAttribute{
						Description: "Whether the email ownership verification is enabled.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"verify_email_template_name": schema.StringAttribute{
						Description: "The template name for verify email. The default is message-template-email-ownership-verification.html.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"email_verification_sent_template_name": schema.StringAttribute{
						Description: "The template name for email verification sent. The default is local.identity.email.verification.sent.html. Note:Only applicable if EmailVerificationType is OTL.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"email_verification_success_template_name": schema.StringAttribute{
						Description: "The template name for email verification success. The default is local.identity.email.verification.success.html.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"email_verification_error_template_name": schema.StringAttribute{
						Description: "The template name for email verification error. The default is local.identity.email.verification.error.html.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"email_verification_type": schema.StringAttribute{
						Description: "Email Verification Type.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"otp_length": schema.Int64Attribute{
						Description: "The OTP length generated for email verification. The default is 8. Note: Only applicable if EmailVerificationType is OTP.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"otp_retry_attempts": schema.Int64Attribute{
						Description: "The number of OTP retry attempts for email verification. The default is 3. Note: Only applicable if EmailVerificationType is OTP.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"allowed_otp_character_set": schema.StringAttribute{
						Description: "The allowed character set used to generate the OTP. The default is 23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz. Note: Only applicable if EmailVerificationType is OTP.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"otp_time_to_live": schema.Int64Attribute{
						Description: "Field used OTP time to live. The default is 15. Note: Only applicable if EmailVerificationType is OTP.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"email_verification_otp_template_name": schema.StringAttribute{
						Description: "The template name for email verification OTP verification. The default is local.identity.email.verification.otp.html. Note: Only applicable if EmailVerificationType is OTP.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"otl_time_to_live": schema.Int64Attribute{
						Description: "Field used OTL time to live. The default is 1440. Note: Only applicable if EmailVerificationType is OTL.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"field_for_email_to_verify": schema.StringAttribute{
						Description: "Field used for email ownership verification. Note: Not required when emailVerificationEnabled is set to false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"field_storing_verification_status": schema.StringAttribute{
						Description: "Field used for storing email verification status. Note: Not required when emailVerificationEnabled is set to false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"notification_publisher_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated notification publisher.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"require_verified_email": schema.BoolAttribute{
						Description: "Whether the user must verify their email address before they can complete a single sign-on transaction. The default is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"require_verified_email_template_name": schema.StringAttribute{
						Description: "The template to render when the user must verify their email address before they can complete a single sign-on transaction. The default is local.identity.email.verification.required.html. Note:Only applicable if EmailVerificationType is OTL and requireVerifiedEmail is true.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"data_store_config": schema.SingleNestedAttribute{
				Description: "The local identity profile data store configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"base_dn": schema.StringAttribute{
						Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The data store config type.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"data_store_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated data store.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the resource.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
							"location": schema.StringAttribute{
								Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
								Required:    false,
								Optional:    false,
								Computed:    true,
							},
						},
					},
					"data_store_mapping": schema.MapNestedAttribute{
						Description: "The data store mapping.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The data store attribute type.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "The data store attribute name.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"metadata": schema.MapAttribute{
									Description: "The data store attribute metadata.",
									Required:    false,
									Optional:    false,
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"create_pattern": schema.StringAttribute{
						Description: "The Relative DN Pattern that will be used to create objects in the directory.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"object_class": schema.StringAttribute{
						Description: "The Object Class used by the new objects stored in the LDAP data store.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"auxiliary_object_classes": schema.SetAttribute{
						Description: "The Auxiliary Object Classes used by the new objects stored in the LDAP data store.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"profile_enabled": schema.BoolAttribute{
				Description: "Whether the profile configuration is enabled or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}

	id.ToDataSourceSchema(&schemaDef, true, "The persistent, unique ID for the local identity profile. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *localIdentityIdentityProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_identity_identity_profile"
}

// Configure adds the provider configured client to the data source.
func (r *localIdentityIdentityProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLocalIdentityIdentityProfileResponseDataSource(ctx context.Context, r *client.LocalIdentityProfile, state *localIdentityIdentityProfileDataSourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
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

// Read resource information
func (r *localIdentityIdentityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state localIdentityIdentityProfileDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadLocalIdentityIdentityProfile, httpResp, err := r.apiClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Identity Profile", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadLocalIdentityIdentityProfile.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of a Local Identity Identity Profile: %s", responseErr.Error())
	}

	// Read the response into the state
	diags = readLocalIdentityIdentityProfileResponseDataSource(ctx, apiReadLocalIdentityIdentityProfile, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
