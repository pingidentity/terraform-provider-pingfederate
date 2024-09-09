package notificationpublisherssettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &notificationPublisherSettingsResource{}
	_ resource.ResourceWithConfigure   = &notificationPublisherSettingsResource{}
	_ resource.ResourceWithImportState = &notificationPublisherSettingsResource{}
)

// NotificationPublisherSettingsResource is a helper function to simplify the provider implementation.
func NotificationPublisherSettingsResource() resource.Resource {
	return &notificationPublisherSettingsResource{}
}

// notificationPublisherSettingsResource is the resource implementation.
type notificationPublisherSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type notificationPublisherSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	DefaultNotificationPublisherRef types.Object `tfsdk:"default_notification_publisher_ref"`
}

// GetSchema defines the schema for the resource.
func (r *notificationPublisherSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Notification Publisher Settings",
		Attributes: map[string]schema.Attribute{
			"default_notification_publisher_ref": resourcelink.CompleteSingleNestedAttribute(
				false,
				false,
				true,
				"The default notification publisher reference",
			),
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema

}

// Metadata returns the resource type name.
func (r *notificationPublisherSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_publisher_settings"
}

func (r *notificationPublisherSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readNotificationPublisherSettingsResponse(ctx context.Context, r *client.NotificationPublishersSettings, state *notificationPublisherSettingsResourceModel, existingId *string) diag.Diagnostics {
	var diags diag.Diagnostics

	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}
	state.DefaultNotificationPublisherRef, diags = resourcelink.ToState(ctx, r.DefaultNotificationPublisherRef)

	// make sure all object type building appends diags
	return diags
}

func (r *notificationPublisherSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var plan notificationPublisherSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createNotificationPublisherSettings := client.NewNotificationPublishersSettings()
	createNotificationPublisherSettings.DefaultNotificationPublisherRef, err = resourcelink.ClientStruct(plan.DefaultNotificationPublisherRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add request for Notification Publishers settings: "+err.Error())
		return
	}

	apiCreateNotificationPublisherSettings := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateNotificationPublisherSettings = apiCreateNotificationPublisherSettings.Body(*createNotificationPublisherSettings)
	notificationPublisherSettingsResponse, httpResp, err := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettingsExecute(apiCreateNotificationPublisherSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Notification Publishers Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state notificationPublisherSettingsResourceModel

	diags = readNotificationPublisherSettingsResponse(ctx, notificationPublisherSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationPublisherSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationPublisherSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadNotificationPublisherSettings, httpResp, err := r.apiClient.NotificationPublishersAPI.GetNotificationPublishersSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Notification Publisher Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Notification Publishers Settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readNotificationPublisherSettingsResponse(ctx, apiReadNotificationPublisherSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationPublisherSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var plan notificationPublisherSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createUpdateRequest := client.NewNotificationPublishersSettings()
	createUpdateRequest.DefaultNotificationPublisherRef, err = resourcelink.ClientStruct(plan.DefaultNotificationPublisherRef)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Notification Publishers Settings: "+err.Error())
		return
	}
	updateNotificationPublisherSettings := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettings(config.AuthContext(ctx, r.providerConfig))
	updateNotificationPublisherSettings = updateNotificationPublisherSettings.Body(*createUpdateRequest)
	updateNotificationPublisherSettingsResponse, httpResp, err := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettingsExecute(updateNotificationPublisherSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Notification Publishers Settings", err, httpResp)
		return
	}

	// Read the response
	var state notificationPublisherSettingsResourceModel
	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readNotificationPublisherSettingsResponse(ctx, updateNotificationPublisherSettingsResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *notificationPublisherSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	providererror.WarnConfigurationCannotBeReset("pingfederate_notification_publisher_settings", &resp.Diagnostics)
}

func (r *notificationPublisherSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
