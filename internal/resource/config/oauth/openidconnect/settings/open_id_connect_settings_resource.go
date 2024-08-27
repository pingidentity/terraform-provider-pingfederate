package oauthopenidconnectsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deprecatedOidcResource{}
	_ resource.ResourceWithConfigure   = &deprecatedOidcResource{}
	_ resource.ResourceWithImportState = &deprecatedOidcResource{}
)

// DeprecatedOidcResource is a helper function to simplify the provider implementation.
func DeprecatedOidcResource() resource.Resource {
	return &deprecatedOidcResource{
		impl: openidConnectSettingsResource{},
	}
}

// deprecatedOidcResource is the resource implementation.
type deprecatedOidcResource struct {
	impl openidConnectSettingsResource
}

// GetSchema defines the schema for the resource.
func (r *deprecatedOidcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_open_id_connect_settings` resource is deprecated. Use the `pingfederate_openid_connect_settings` resource instead."
}

// Metadata returns the resource type name.
func (r *deprecatedOidcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_open_id_connect_settings"
}

func (r *deprecatedOidcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *deprecatedOidcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *deprecatedOidcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deprecatedOidcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *deprecatedOidcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *deprecatedOidcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
