package serversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsGeneralSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsGeneralSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsGeneralSettingsResource{}
)

// ServerSettingsGeneralSettingsResource is a helper function to simplify the provider implementation.
func ServerSettingsGeneralSettingsResource() resource.Resource {
	return &serverSettingsGeneralSettingsResource{}
}

// serverSettingsGeneralSettingsResource is the resource implementation.
type serverSettingsGeneralSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsGeneralSettingsResourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	DisableAutomaticConnectionValidation    types.Bool   `tfsdk:"disable_automatic_connection_validation"`
	IdpConnectionTransactionLoggingOverride types.String `tfsdk:"idp_connection_transaction_logging_override"`
	SpConnectionTransactionLoggingOverride  types.String `tfsdk:"sp_connection_transaction_logging_override"`
	DatastoreValidationIntervalSecs         types.Int64  `tfsdk:"datastore_validation_interval_secs"`
	RequestHeaderForCorrelationId           types.String `tfsdk:"request_header_for_correlation_id"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsGeneralSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a ServerSettingsGeneralSettings.",
		Attributes: map[string]schema.Attribute{
			"datastore_validation_interval_secs": schema.Int64Attribute{
				Description: "Determines how long (in seconds) the result of testing a datastore connection is cached. The default is 300.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown()},
			},
			"disable_automatic_connection_validation": schema.BoolAttribute{
				Description: "Boolean that disables automatic connection validation when set to true. The default is false.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown()},
			},
			"idp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Determines the level of transaction logging for all identity provider connections. The default is DONT_OVERRIDE, in which case the logging level will be determined by each individual IdP connection [ DONT_OVERRIDE, NONE, FULL, STANDARD, ENHANCED ]",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown()},
			},
			"request_header_for_correlation_id": schema.StringAttribute{
				Description: "HTTP request header for retrieving correlation ID.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown()},
			},
			"sp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Determines the level of transaction logging for all service provider connections. The default is DONT_OVERRIDE, in which case the logging level will be determined by each individual SP connection [ DONT_OVERRIDE, NONE, FULL, STANDARD, ENHANCED ]",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown()},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalServerSettingsGeneralSettingsFields(ctx context.Context, addRequest *client.GeneralSettings, plan serverSettingsGeneralSettingsResourceModel) error {
	if internaltypes.IsDefined(plan.DisableAutomaticConnectionValidation) {
		addRequest.DisableAutomaticConnectionValidation = plan.DisableAutomaticConnectionValidation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IdpConnectionTransactionLoggingOverride) {
		addRequest.IdpConnectionTransactionLoggingOverride = plan.IdpConnectionTransactionLoggingOverride.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SpConnectionTransactionLoggingOverride) {
		addRequest.SpConnectionTransactionLoggingOverride = plan.SpConnectionTransactionLoggingOverride.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DatastoreValidationIntervalSecs) {
		addRequest.DatastoreValidationIntervalSecs = plan.DatastoreValidationIntervalSecs.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.RequestHeaderForCorrelationId) {
		addRequest.RequestHeaderForCorrelationId = plan.RequestHeaderForCorrelationId.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *serverSettingsGeneralSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general_settings"
}

func (r *serverSettingsGeneralSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsGeneralSettingsResponse(ctx context.Context, r *client.GeneralSettings, state *serverSettingsGeneralSettingsResourceModel, existingId *string) {
	state.Id = id.GenerateUUIDToState(existingId)
	state.DisableAutomaticConnectionValidation = types.BoolPointerValue(r.DisableAutomaticConnectionValidation)
	state.IdpConnectionTransactionLoggingOverride = internaltypes.StringTypeOrNil(r.IdpConnectionTransactionLoggingOverride, true)
	state.SpConnectionTransactionLoggingOverride = internaltypes.StringTypeOrNil(r.SpConnectionTransactionLoggingOverride, true)
	state.DatastoreValidationIntervalSecs = types.Int64PointerValue(r.DatastoreValidationIntervalSecs)
	state.RequestHeaderForCorrelationId = internaltypes.StringTypeOrNil(r.RequestHeaderForCorrelationId, true)
}

func (r *serverSettingsGeneralSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsGeneralSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServerSettingsGeneralSettings := client.NewGeneralSettings()
	err := addOptionalServerSettingsGeneralSettingsFields(ctx, createServerSettingsGeneralSettings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings General Settings", err.Error())
		return
	}

	apiCreateServerSettingsGeneralSettings := r.apiClient.ServerSettingsAPI.UpdateGeneralSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsGeneralSettings = apiCreateServerSettingsGeneralSettings.Body(*createServerSettingsGeneralSettings)
	serverSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateGeneralSettingsExecute(apiCreateServerSettingsGeneralSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings General Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsGeneralSettingsResourceModel
	readServerSettingsGeneralSettingsResponse(ctx, serverSettingsGeneralSettingsResponse, &state, nil)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsGeneralSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsGeneralSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsGeneralSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetGeneralSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings General Settings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings General Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readServerSettingsGeneralSettingsResponse(ctx, apiReadServerSettingsGeneralSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsGeneralSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsGeneralSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateServerSettingsGeneralSettings := r.apiClient.ServerSettingsAPI.UpdateGeneralSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewGeneralSettings()
	err := addOptionalServerSettingsGeneralSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Server Settings General Settings", err.Error())
		return
	}

	updateServerSettingsGeneralSettings = updateServerSettingsGeneralSettings.Body(*createUpdateRequest)
	updateServerSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateGeneralSettingsExecute(updateServerSettingsGeneralSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Server Settings General Settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsGeneralSettingsResourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readServerSettingsGeneralSettingsResponse(ctx, updateServerSettingsGeneralSettingsResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsGeneralSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsGeneralSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
