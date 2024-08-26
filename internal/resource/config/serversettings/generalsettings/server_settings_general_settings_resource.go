package serversettingsgeneralsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsGeneralSettingsResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsGeneralSettingsResource{}
	_ resource.ResourceWithImportState = &serverSettingsGeneralSettingsResource{}
)

// ServerSettingsGeneralSettingsResource is a helper function to simplify the provider implementation.
func ServerSettingsGeneralSettingsResource() resource.Resource {
	return &serverSettingsGeneralSettingsResource{
		impl: serverSettingsGeneralResource{},
	}
}

// serverSettingsGeneralSettingsResource is the resource implementation.
type serverSettingsGeneralSettingsResource struct {
	impl serverSettingsGeneralResource
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsGeneralSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_server_settings_general_settings` resource is deprecated. Use the `pingfederate_server_settings_general` resource instead."
}

// Metadata returns the resource type name.
func (r *serverSettingsGeneralSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general_settings"
}

func (r *serverSettingsGeneralSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *serverSettingsGeneralSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *serverSettingsGeneralSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsGeneralSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsGeneralSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *serverSettingsGeneralSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
