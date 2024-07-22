package notificationpublishers

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
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &notificationPublishersSettingsResource{}
	_ resource.ResourceWithConfigure   = &notificationPublishersSettingsResource{}
	_ resource.ResourceWithImportState = &notificationPublishersSettingsResource{}
)

// NotificationPublishersSettingsResource is a helper function to simplify the provider implementation.
func NotificationPublishersSettingsResource() resource.Resource {
	return &notificationPublishersSettingsResource{}
}

// notificationPublishersSettingsResource is the resource implementation.
type notificationPublishersSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type notificationPublishersSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	DefaultNotificationPublisherRef types.Object `tfsdk:"default_notification_publisher_ref"`
}

// GetSchema defines the schema for the resource.
func (r *notificationPublishersSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Notification Publishers Settings",
		Attributes: map[string]schema.Attribute{
			"default_notification_publisher_ref": resourcelink.CompleteSingleNestedAttribute(
				false,
				false,
				true,
				"The default notification publisher reference",
			),
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema

}

// Metadata returns the resource type name.
func (r *notificationPublishersSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_publishers_settings"
}

func (r *notificationPublishersSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readNotificationPublishersSettingsResponse(ctx context.Context, r *client.NotificationPublishersSettings, state *notificationPublishersSettingsResourceModel, existingId *string) diag.Diagnostics {
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

func (r *notificationPublishersSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var plan notificationPublishersSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createNotificationPublishersSettings := client.NewNotificationPublishersSettings()
	createNotificationPublishersSettings.DefaultNotificationPublisherRef, err = resourcelink.ClientStruct(plan.DefaultNotificationPublisherRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add request for Notification Publishers settings", err.Error())
		return
	}

	apiCreateNotificationPublishersSettings := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateNotificationPublishersSettings = apiCreateNotificationPublishersSettings.Body(*createNotificationPublishersSettings)
	notificationPublishersSettingsResponse, httpResp, err := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettingsExecute(apiCreateNotificationPublishersSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Notification Publishers Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state notificationPublishersSettingsResourceModel

	diags = readNotificationPublishersSettingsResponse(ctx, notificationPublishersSettingsResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationPublishersSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationPublishersSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadNotificationPublishersSettings, httpResp, err := r.apiClient.NotificationPublishersAPI.GetNotificationPublishersSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Notification Publishers Settings", err, httpResp)
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
	readNotificationPublishersSettingsResponse(ctx, apiReadNotificationPublishersSettings, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationPublishersSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var plan notificationPublishersSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createUpdateRequest := client.NewNotificationPublishersSettings()
	createUpdateRequest.DefaultNotificationPublisherRef, err = resourcelink.ClientStruct(plan.DefaultNotificationPublisherRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Notification Publishers Settings", err.Error())
		return
	}
	updateNotificationPublishersSettings := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	updateNotificationPublishersSettings = updateNotificationPublishersSettings.Body(*createUpdateRequest)
	updateNotificationPublishersSettingsResponse, httpResp, err := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublishersSettingsExecute(updateNotificationPublishersSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Notification Publishers Settings", err, httpResp)
		return
	}

	// Read the response
	var state notificationPublishersSettingsResourceModel
	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readNotificationPublishersSettingsResponse(ctx, updateNotificationPublishersSettingsResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *notificationPublishersSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *notificationPublishersSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
