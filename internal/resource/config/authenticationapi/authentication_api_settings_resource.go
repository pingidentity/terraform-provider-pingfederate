package authenticationapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationApiSettingsResource{}
	_ resource.ResourceWithConfigure   = &authenticationApiSettingsResource{}
	_ resource.ResourceWithImportState = &authenticationApiSettingsResource{}
)

// AuthenticationApiSettingsResource is a helper function to simplify the provider implementation.
func AuthenticationApiSettingsResource() resource.Resource {
	return &authenticationApiSettingsResource{}
}

// authenticationApiSettingsResource is the resource implementation.
type authenticationApiSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}
type authenticationApiSettingsResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	ApiEnabled                       types.Bool   `tfsdk:"api_enabled"`
	EnableApiDescriptions            types.Bool   `tfsdk:"enable_api_descriptions"`
	RestrictAccessToRedirectlessMode types.Bool   `tfsdk:"restrict_access_to_redirectless_mode"`
	IncludeRequestContext            types.Bool   `tfsdk:"include_request_context"`
	DefaultApplicationRef            types.Object `tfsdk:"default_application_ref"`
}

// GetSchema defines the schema for the resource.
func (r *authenticationApiSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a AuthenticationApiSettings.",
		Attributes: map[string]schema.Attribute{
			"api_enabled": schema.BoolAttribute{
				Description: "Enable Authentication API",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_api_descriptions": schema.BoolAttribute{
				Description: "Enable API descriptions",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_application_ref": schema.SingleNestedAttribute{
				Description: "Enable API descriptions",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: config.AddResourceLinkSchema(),
			},
			"restrict_access_to_redirectless_mode": schema.BoolAttribute{
				Description: "Enable restrict access to redirectless mode",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_context": schema.BoolAttribute{
				Description: "Includes request context in API responses",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

func addAuthenticationApiSettingsFields(ctx context.Context, addRequest *client.AuthnApiSettings, plan authenticationApiSettingsResourceModel) error {
	if internaltypes.IsDefined(plan.ApiEnabled) {
		addRequest.ApiEnabled = plan.ApiEnabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EnableApiDescriptions) {
		addRequest.EnableApiDescriptions = plan.EnableApiDescriptions.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RestrictAccessToRedirectlessMode) {
		addRequest.RestrictAccessToRedirectlessMode = plan.RestrictAccessToRedirectlessMode.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestContext) {
		addRequest.IncludeRequestContext = plan.IncludeRequestContext.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DefaultApplicationRef) {
		addRequestNewLinkObj := internaltypes.ToRequestResourceLink(ctx, plan.DefaultApplicationRef)
		addRequest.DefaultApplicationRef = addRequestNewLinkObj
	}
	return nil

}

// Metadata returns the resource type name.
func (r *authenticationApiSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_api_settings"
}

func (r *authenticationApiSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationApiSettingsResponse(ctx context.Context, r *client.AuthnApiSettings, state *authenticationApiSettingsResourceModel, expectedValues *authenticationApiSettingsResourceModel, diags *diag.Diagnostics) {
	state.Id = types.StringValue("id")
	state.ApiEnabled = types.BoolValue(*r.ApiEnabled)
	state.EnableApiDescriptions = types.BoolValue(*r.EnableApiDescriptions)
	state.RestrictAccessToRedirectlessMode = types.BoolValue(*r.RestrictAccessToRedirectlessMode)
	state.IncludeRequestContext = types.BoolValue(*r.IncludeRequestContext)
	resourceLinkObjectValue := internaltypes.ToStateResourceLink(ctx, r.GetDefaultApplicationRef())
	state.DefaultApplicationRef = resourceLinkObjectValue
}

func (r *authenticationApiSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan authenticationApiSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get the current state to see how any attributes are changing
	updateAuthenticationApiSettings := r.apiClient.AuthenticationApiApi.UpdateAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewAuthnApiSettings()
	err := addAuthenticationApiSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationApiSettings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateAuthenticationApiSettings = updateAuthenticationApiSettings.Body(*createUpdateRequest)
	updateAuthenticationApiSettingsResponse, httpResp, err := r.apiClient.AuthenticationApiApi.UpdateAuthenticationApiSettingsExecute(updateAuthenticationApiSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AuthenticationApiSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAuthenticationApiSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	var state authenticationApiSettingsResourceModel
	readAuthenticationApiSettingsResponse(ctx, updateAuthenticationApiSettingsResponse, &state, &plan, &diags)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *authenticationApiSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAuthenticationApiSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAuthenticationApiSettings(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state authenticationApiSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationApiSettings, httpResp, err := apiClient.AuthenticationApiApi.GetAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a AuthenticationApiSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAuthenticationApiSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthenticationApiSettingsResponse(ctx, apiReadAuthenticationApiSettings, &state, &state, &diags)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationApiSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAuthenticationApiSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

// Update updates the resource and sets the updated Terraform state on success.
func updateAuthenticationApiSettings(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var plan authenticationApiSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get the current state to see how any attributes are changing
	updateAuthenticationApiSettings := apiClient.AuthenticationApiApi.UpdateAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewAuthnApiSettings()
	err := addAuthenticationApiSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationApiSettings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateAuthenticationApiSettings = updateAuthenticationApiSettings.Body(*createUpdateRequest)
	updateAuthenticationApiSettingsResponse, httpResp, err := apiClient.AuthenticationApiApi.UpdateAuthenticationApiSettingsExecute(updateAuthenticationApiSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AuthenticationApiSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAuthenticationApiSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	var state authenticationApiSettingsResourceModel
	readAuthenticationApiSettingsResponse(ctx, updateAuthenticationApiSettingsResponse, &state, &plan, &diags)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
func (r *authenticationApiSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *authenticationApiSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAuthenticationApiSettingsLocation(ctx, req, resp)
}
func importAuthenticationApiSettingsLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
