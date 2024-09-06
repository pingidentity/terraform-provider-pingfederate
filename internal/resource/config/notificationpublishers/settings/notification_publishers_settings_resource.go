package notificationpublisherssettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &notificationPublishersSettingsResource{}
	_ resource.ResourceWithConfigure   = &notificationPublishersSettingsResource{}
	_ resource.ResourceWithImportState = &notificationPublishersSettingsResource{}
)

// NotificationPublishersSettingsResource is a helper function to simplify the provider implementation.
func NotificationPublishersSettingsResource() resource.Resource {
	return &notificationPublishersSettingsResource{
		impl: notificationPublisherSettingsResource{},
	}
}

// notificationPublishersSettingsResource is the resource implementation.
type notificationPublishersSettingsResource struct {
	impl notificationPublisherSettingsResource
}

// GetSchema defines the schema for the resource.
func (r *notificationPublishersSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_notification_publishers_settings` resource is deprecated. Use the `pingfederate_notification_publisher_settings` resource instead."
}

// Metadata returns the resource type name.
func (r *notificationPublishersSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_publishers_settings"
}

func (r *notificationPublishersSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *notificationPublishersSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *notificationPublishersSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationPublishersSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *notificationPublishersSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *notificationPublishersSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
