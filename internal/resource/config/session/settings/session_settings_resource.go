package sessionsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionSettingsResource{}
	_ resource.ResourceWithConfigure   = &sessionSettingsResource{}
	_ resource.ResourceWithImportState = &sessionSettingsResource{}
)

// SessionSettingsResource is a helper function to simplify the provider implementation.
func SessionSettingsResource() resource.Resource {
	return &sessionSettingsResource{}
}

// sessionSettingsResource is the resource implementation.
type sessionSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *sessionSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings for general session management.",
		Attributes: map[string]schema.Attribute{
			"track_adapter_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether adapter sessions are tracked for cleanup during single logout. The default is false.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"revoke_user_session_on_logout": schema.BoolAttribute{
				Description: "Determines whether the user's session is revoked on logout. If this property is not provided on a PUT, the setting is left unchanged.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"session_revocation_lifetime": schema.Int64Attribute{
				Description: "How long a session revocation is tracked and stored, in minutes. If this property is not provided on a PUT, the setting is left unchanged.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(40),
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalSessionSettingsFields(ctx context.Context, addRequest *client.SessionSettings, plan sessionSettingsModel) error {
	if internaltypes.IsDefined(plan.TrackAdapterSessionsForLogout) {
		addRequest.TrackAdapterSessionsForLogout = plan.TrackAdapterSessionsForLogout.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RevokeUserSessionOnLogout) {
		addRequest.RevokeUserSessionOnLogout = plan.RevokeUserSessionOnLogout.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SessionRevocationLifetime) {
		addRequest.SessionRevocationLifetime = plan.SessionRevocationLifetime.ValueInt64Pointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *sessionSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_settings"
}

func (r *sessionSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *sessionSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionSettingsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSessionSettings := client.NewSessionSettings()
	err := addOptionalSessionSettingsFields(ctx, createSessionSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Settings", err.Error())
		return
	}

	apiCreateSessionSettings := r.apiClient.SessionAPI.UpdateSessionSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSessionSettings = apiCreateSessionSettings.Body(*createSessionSettings)
	sessionSettingsResponse, httpResp, err := r.apiClient.SessionAPI.UpdateSessionSettingsExecute(apiCreateSessionSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Session Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state sessionSettingsModel
	readSessionSettingsResponse(ctx, sessionSettingsResponse, &state, nil)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sessionSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionSettingsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionSettings, httpResp, err := r.apiClient.SessionAPI.GetSessionSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Session Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Session Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionSettingsResponse(ctx, apiReadSessionSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sessionSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSessionSettings := r.apiClient.SessionAPI.UpdateSessionSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewSessionSettings()
	err := addOptionalSessionSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Settings", err.Error())
		return
	}

	updateSessionSettings = updateSessionSettings.Body(*createUpdateRequest)
	updateSessionSettingsResponse, httpResp, err := r.apiClient.SessionAPI.UpdateSessionSettingsExecute(updateSessionSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Session Settings", err, httpResp)
		return
	}

	// Get the current state to see how any attributes are changing
	var state sessionSettingsModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionSettingsResponse(ctx, updateSessionSettingsResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *sessionSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
