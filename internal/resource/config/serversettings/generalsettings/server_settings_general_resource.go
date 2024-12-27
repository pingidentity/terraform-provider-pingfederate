package serversettingsgeneralsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsGeneralResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsGeneralResource{}
	_ resource.ResourceWithImportState = &serverSettingsGeneralResource{}
)

// ServerSettingsGeneralResource is a helper function to simplify the provider implementation.
func ServerSettingsGeneralResource() resource.Resource {
	return &serverSettingsGeneralResource{}
}

// serverSettingsGeneralResource is the resource implementation.
type serverSettingsGeneralResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsGeneralResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the general server settings.",
		Attributes: map[string]schema.Attribute{
			"datastore_validation_interval_secs": schema.Int64Attribute{
				Description: "Determines how long (in seconds) the result of testing a datastore connection is cached. The default is `300`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(300),
			},
			"disable_automatic_connection_validation": schema.BoolAttribute{
				Description: "Boolean that disables automatic connection validation when set to true. The default is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Determines the level of transaction logging for all identity provider connections. The default is `DONT_OVERRIDE`, in which case the logging level will be determined by each individual IdP connection. Options are `DONT_OVERRIDE`, `NONE`, `FULL`, `STANDARD`, `ENHANCED`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DONT_OVERRIDE"),
				Validators: []validator.String{
					stringvalidator.OneOf("DONT_OVERRIDE", "NONE", "FULL", "STANDARD", "ENHANCED"),
				},
			},
			"request_header_for_correlation_id": schema.StringAttribute{
				Description: "HTTP request header for retrieving correlation ID.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"sp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Determines the level of transaction logging for all service provider connections. The default is `DONT_OVERRIDE`, in which case the logging level will be determined by each individual SP connection. Options are `DONT_OVERRIDE`, `NONE`, `FULL`, `STANDARD`, `ENHANCED`.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("DONT_OVERRIDE"),
				Validators: []validator.String{
					stringvalidator.OneOf("DONT_OVERRIDE", "NONE", "FULL", "STANDARD", "ENHANCED"),
				},
			},
		},
	}
	resp.Schema = schema
}

func addOptionalServerSettingsGeneralFields(ctx context.Context, addRequest *client.GeneralSettings, plan serverSettingsGeneralModel) error {
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
func (r *serverSettingsGeneralResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general"
}

func (r *serverSettingsGeneralResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (m *serverSettingsGeneralModel) buildDefaultClientStruct() *client.GeneralSettings {
	return &client.GeneralSettings{
		DisableAutomaticConnectionValidation:    utils.Pointer(false),
		IdpConnectionTransactionLoggingOverride: utils.Pointer("DONT_OVERRIDE"),
		SpConnectionTransactionLoggingOverride:  utils.Pointer("DONT_OVERRIDE"),
		DatastoreValidationIntervalSecs:         utils.Pointer(int64(300)),
		RequestHeaderForCorrelationId:           utils.Pointer(""),
	}
}

func (r *serverSettingsGeneralResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsGeneralModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServerSettingsGeneral := client.NewGeneralSettings()
	err := addOptionalServerSettingsGeneralFields(ctx, createServerSettingsGeneral, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for general server settings: "+err.Error())
		return
	}

	apiCreateServerSettingsGeneral := r.apiClient.ServerSettingsAPI.UpdateGeneralSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsGeneral = apiCreateServerSettingsGeneral.Body(*createServerSettingsGeneral)
	serverSettingsGeneralResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateGeneralSettingsExecute(apiCreateServerSettingsGeneral)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the general server settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsGeneralModel
	readServerSettingsGeneralResponse(ctx, serverSettingsGeneralResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsGeneralResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsGeneralModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsGeneral, httpResp, err := r.apiClient.ServerSettingsAPI.GetGeneralSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "general server settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the general server settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readServerSettingsGeneralResponse(ctx, apiReadServerSettingsGeneral, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsGeneralResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsGeneralModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateServerSettingsGeneral := r.apiClient.ServerSettingsAPI.UpdateGeneralSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewGeneralSettings()
	err := addOptionalServerSettingsGeneralFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for general server settings: "+err.Error())
		return
	}

	updateServerSettingsGeneral = updateServerSettingsGeneral.Body(*createUpdateRequest)
	updateServerSettingsGeneralResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateGeneralSettingsExecute(updateServerSettingsGeneral)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating general server settings", err, httpResp)
		return
	}

	// Read the response
	var state serverSettingsGeneralModel
	readServerSettingsGeneralResponse(ctx, updateServerSettingsGeneralResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsGeneralResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model serverSettingsGeneralModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateGeneralSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateGeneralSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the general server settings", err, httpResp)
	}
}

func (r *serverSettingsGeneralResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState serverSettingsGeneralModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
