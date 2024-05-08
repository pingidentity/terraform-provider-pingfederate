// Code generated by ping-terraform-plugin-framework-generator

package sessionauthenticationsessionpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &authenticationSessionPolicyResource{}
	_ resource.ResourceWithConfigure   = &authenticationSessionPolicyResource{}
	_ resource.ResourceWithImportState = &authenticationSessionPolicyResource{}
)

func AuthenticationSessionPolicyResource() resource.Resource {
	return &authenticationSessionPolicyResource{}
}

type authenticationSessionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *authenticationSessionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_session_policy"
}

func (r *authenticationSessionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type authenticationSessionPolicyResourceModel struct {
	AuthenticationSource  types.Object `tfsdk:"authentication_source"`
	AuthnContextSensitive types.Bool   `tfsdk:"authn_context_sensitive"`
	EnableSessions        types.Bool   `tfsdk:"enable_sessions"`
	Id                    types.String `tfsdk:"id"`
	PolicyId              types.String `tfsdk:"policy_id"`
	IdleTimeoutMins       types.Int64  `tfsdk:"idle_timeout_mins"`
	MaxTimeoutMins        types.Int64  `tfsdk:"max_timeout_mins"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	TimeoutDisplayUnit    types.String `tfsdk:"timeout_display_unit"`
	UserDeviceType        types.String `tfsdk:"user_device_type"`
}

func (r *authenticationSessionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"authentication_source": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"source_ref": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required:    true,
								Description: "The ID of the resource.",
							},
						},
						Required:    true,
						Description: "A reference to a resource.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of this authentication source.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"IDP_ADAPTER",
								"IDP_CONNECTION",
							),
						},
					},
				},
				Required:    true,
				Description: "An authentication source (IdP adapter or IdP connection).",
			},
			"authn_context_sensitive": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Determines whether the requested authentication context is considered when deciding whether an existing session is valid for a given request. The default is false.",
			},
			"enable_sessions": schema.BoolAttribute{
				Required:    true,
				Description: "Determines whether sessions are enabled for the authentication source. This value overrides the enableSessions value from the global authentication session policy.",
			},
			"idle_timeout_mins": schema.Int64Attribute{
				Optional:    true,
				Description: "The idle timeout period, in minutes. If omitted, the value from the global authentication session policy will be used. If set to -1, the idle timeout will be set to the maximum timeout. If a value is provided for this property, a value must also be provided for maxTimeoutMins.",
			},
			"max_timeout_mins": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum timeout period, in minutes. If omitted, the value from the global authentication session policy will be used. If set to -1, sessions do not expire. If a value is provided for this property, a value must also be provided for idleTimeoutMins.",
			},
			"persistent": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Determines whether sessions for the authentication source are persistent. This value overrides the persistentSessions value from the global authentication session policy.This field is ignored if enableSessions is false.",
			},
			"timeout_display_unit": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("MINUTES"),
				Description: "The display unit for session timeout periods in the PingFederate administrative console. When the display unit is HOURS or DAYS, the timeout values in minutes must correspond to a whole number value for the specified unit.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"MINUTES",
						"HOURS",
						"DAYS",
					),
				},
			},
			"user_device_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("PRIVATE"),
				Description: "Determines the type of user device that the authentication session can be created on. If empty, the value will default to PRIVATE.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PRIVATE",
						"SHARED",
						"ANY",
					),
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"policy_id",
		false,
		false,
		"The persistent, unique ID for the session policy. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
	)

	resp.Schema = schema
}

func (model *authenticationSessionPolicyResourceModel) buildClientStruct() *client.AuthenticationSessionPolicy {
	result := &client.AuthenticationSessionPolicy{}
	// authentication_source
	authenticationSourceValue := client.AuthenticationSource{}
	authenticationSourceAttrs := model.AuthenticationSource.Attributes()
	authenticationSourceSourceRefValue := client.ResourceLink{}
	authenticationSourceSourceRefAttrs := authenticationSourceAttrs["source_ref"].(types.Object).Attributes()
	authenticationSourceSourceRefValue.Id = authenticationSourceSourceRefAttrs["id"].(types.String).ValueString()
	authenticationSourceValue.SourceRef = authenticationSourceSourceRefValue
	authenticationSourceValue.Type = authenticationSourceAttrs["type"].(types.String).ValueString()
	result.AuthenticationSource = authenticationSourceValue

	// authn_context_sensitive
	result.AuthnContextSensitive = model.AuthnContextSensitive.ValueBoolPointer()
	// enable_sessions
	result.EnableSessions = model.EnableSessions.ValueBool()
	// id
	result.Id = model.Id.ValueStringPointer()
	// idle_timeout_mins
	result.IdleTimeoutMins = model.IdleTimeoutMins.ValueInt64Pointer()
	// max_timeout_mins
	result.MaxTimeoutMins = model.MaxTimeoutMins.ValueInt64Pointer()
	// persistent
	result.Persistent = model.Persistent.ValueBoolPointer()
	// timeout_display_unit
	result.TimeoutDisplayUnit = model.TimeoutDisplayUnit.ValueStringPointer()
	// user_device_type
	result.UserDeviceType = model.UserDeviceType.ValueStringPointer()
	return result
}

func (state *authenticationSessionPolicyResourceModel) readClientResponse(response *client.AuthenticationSessionPolicy) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// authentication_source
	authenticationSourceSourceRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	authenticationSourceAttrTypes := map[string]attr.Type{
		"source_ref": types.ObjectType{AttrTypes: authenticationSourceSourceRefAttrTypes},
		"type":       types.StringType,
	}
	authenticationSourceSourceRefValue, diags := types.ObjectValue(authenticationSourceSourceRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.AuthenticationSource.SourceRef.Id),
	})
	respDiags.Append(diags...)
	authenticationSourceValue, diags := types.ObjectValue(authenticationSourceAttrTypes, map[string]attr.Value{
		"source_ref": authenticationSourceSourceRefValue,
		"type":       types.StringValue(response.AuthenticationSource.Type),
	})
	respDiags.Append(diags...)

	state.AuthenticationSource = authenticationSourceValue
	// authn_context_sensitive
	state.AuthnContextSensitive = types.BoolPointerValue(response.AuthnContextSensitive)
	// enable_sessions
	state.EnableSessions = types.BoolValue(response.EnableSessions)
	// id
	state.Id = types.StringPointerValue(response.Id)
	// idle_timeout_mins
	state.IdleTimeoutMins = types.Int64PointerValue(response.IdleTimeoutMins)
	// max_timeout_mins
	state.MaxTimeoutMins = types.Int64PointerValue(response.MaxTimeoutMins)
	// persistent
	state.Persistent = types.BoolPointerValue(response.Persistent)
	// timeout_display_unit
	state.TimeoutDisplayUnit = types.StringPointerValue(response.TimeoutDisplayUnit)
	// user_device_type
	state.UserDeviceType = types.StringPointerValue(response.UserDeviceType)
	return respDiags
}

func (r *authenticationSessionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data authenticationSessionPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData := data.buildClientStruct()
	apiCreateRequest := r.apiClient.SessionAPI.CreateSourcePolicy(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SessionAPI.CreateSourcePolicyExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Authentication Session Policy", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticationSessionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data authenticationSessionPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.SessionAPI.GetSourcePolicy(config.AuthContext(ctx, r.providerConfig), data.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the Authentication Session Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the Authentication Session Policy", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticationSessionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data authenticationSessionPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SessionAPI.UpdateSourcePolicy(config.AuthContext(ctx, r.providerConfig), data.PolicyId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SessionAPI.UpdateSourcePolicyExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authentication Session Policy", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticationSessionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data authenticationSessionPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.SessionAPI.DeleteSourcePolicy(config.AuthContext(ctx, r.providerConfig), data.PolicyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Authentication Session Policy", err, httpResp)
	}
}

func (r *authenticationSessionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to policy_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("policy_id"), req, resp)
}
