package oauthopenidconnectsettings

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &openIdConnectSettingsResource{}
	_ resource.ResourceWithConfigure   = &openIdConnectSettingsResource{}
	_ resource.ResourceWithImportState = &openIdConnectSettingsResource{}

	openIdConnectSettingsAttrTypes = map[string]attr.Type{
		"track_user_sessions_for_logout": types.BoolType,
		"revoke_user_session_on_logout":  types.BoolType,
		"session_revocation_lifetime":    types.Int64Type,
	}
)

// OpenIdConnectSettingsResource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectSettingsResource() resource.Resource {
	return &openIdConnectSettingsResource{}
}

// openIdConnectSettingsResource is the resource implementation.
type openIdConnectSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type openIdConnectSettingsResourceModel struct {
	Id               types.String `tfsdk:"id"`
	DefaultPolicyRef types.Object `tfsdk:"default_policy_ref"`
	SessionSettings  types.Object `tfsdk:"session_settings"`
}

// GetSchema defines the schema for the resource.
func (r *openIdConnectSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages OpenID Connect configuration settings",
		Attributes: map[string]schema.Attribute{
			"default_policy_ref": schema.SingleNestedAttribute{
				Description: "Reference to the default policy.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(types.ObjectNull(resourcelink.AttrType())),
				Attributes:  resourcelink.ToSchema(),
			},
			"session_settings": schema.SingleNestedAttribute{
				Description: "The session settings",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"track_user_sessions_for_logout": schema.BoolAttribute{
						Description: "Determines whether user sessions are tracked for logout. This property is now available under /oauth/authServerSettings and should be accessed through that resource.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"revoke_user_session_on_logout": schema.BoolAttribute{
						Description: "Determines whether the user's session is revoked on logout. This property is now available under /session/settings and should be accessed through that resource.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"session_revocation_lifetime": schema.Int64Attribute{
						Description: "How long a session revocation is tracked and stored, in minutes. This property is now available under /session/settings and should be accessed through that resource.",
						Computed:    true,
						Optional:    true,
						Default:     int64default.StaticInt64(1450),
						Validators: []validator.Int64{
							// session_revocation_lifetime must be between 1 and 43200 minutes, inclusive
							int64validator.Between(1, 43200),
						},
					},
				},
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalOpenIdConnectSettingsFields(ctx context.Context, addRequest *client.OpenIdConnectSettings, plan openIdConnectSettingsResourceModel) error {
	var err error

	if internaltypes.IsDefined(plan.DefaultPolicyRef) {
		addRequest.DefaultPolicyRef, err = resourcelink.ClientStruct(plan.DefaultPolicyRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SessionSettings) {
		addRequest.SessionSettings = &client.OIDCSessionSettings{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SessionSettings, false)), addRequest.SessionSettings)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *openIdConnectSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_open_id_connect_settings"
}

func (r *openIdConnectSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOpenIdConnectSettingsResponse(ctx context.Context, r *client.OpenIdConnectSettings, state *openIdConnectSettingsResourceModel, existingId *string) diag.Diagnostics {

	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}

	var diags, respDiags diag.Diagnostics

	state.DefaultPolicyRef, respDiags = resourcelink.ToState(ctx, r.DefaultPolicyRef)
	diags = append(diags, respDiags...)
	sessionSettings, respDiags := types.ObjectValueFrom(ctx, openIdConnectSettingsAttrTypes, r.SessionSettings)
	diags = append(diags, respDiags...)
	state.SessionSettings = sessionSettings

	// make sure all object type building appends diags
	return diags
}

func (r *openIdConnectSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan openIdConnectSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpenIdConnectSettings := client.NewOpenIdConnectSettings()
	err := addOptionalOpenIdConnectSettingsFields(ctx, createOpenIdConnectSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OpenID Connect settings", err.Error())
		return
	}
	apiCreateOpenIdConnectSettings := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOpenIdConnectSettings = apiCreateOpenIdConnectSettings.Body(*createOpenIdConnectSettings)
	openIdConnectSettingsResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettingsExecute(apiCreateOpenIdConnectSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OpenID Connect settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state openIdConnectSettingsResourceModel

	diags = readOpenIdConnectSettingsResponse(ctx, openIdConnectSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *openIdConnectSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state openIdConnectSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOpenIdConnectSettings, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.GetOIDCSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OpenID Connect settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OpenID Connect settings", err, httpResp)
		}
		return
	}

	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the response into the state
	readOpenIdConnectSettingsResponse(ctx, apiReadOpenIdConnectSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *openIdConnectSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan openIdConnectSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpenIdConnectSettings := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewOpenIdConnectSettings()
	err := addOptionalOpenIdConnectSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OpenID Connect settings", err.Error())
		return
	}

	updateOpenIdConnectSettings = updateOpenIdConnectSettings.Body(*createUpdateRequest)
	updateOpenIdConnectSettingsResponse, httpResp, err := r.apiClient.OauthOpenIdConnectAPI.UpdateOIDCSettingsExecute(updateOpenIdConnectSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OpenID Connect settings", err, httpResp)
		return
	}

	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the response
	var state openIdConnectSettingsResourceModel
	diags = readOpenIdConnectSettingsResponse(ctx, updateOpenIdConnectSettingsResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *openIdConnectSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *openIdConnectSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
