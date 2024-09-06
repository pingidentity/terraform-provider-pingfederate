package pingoneconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deprecatedPingOneConnectionResource{}
	_ resource.ResourceWithConfigure   = &deprecatedPingOneConnectionResource{}
	_ resource.ResourceWithImportState = &deprecatedPingOneConnectionResource{}
)

// DeprecatedPingOneConnectionResource is a helper function to simplify the provider implementation.
func DeprecatedPingOneConnectionResource() resource.Resource {
	return &deprecatedPingOneConnectionResource{
		impl: pingoneConnectionResource{},
	}
}

// deprecatedPingOneConnectionResource is the resource implementation.
type deprecatedPingOneConnectionResource struct {
	impl pingoneConnectionResource
}

// GetSchema defines the schema for the resource.
func (r *deprecatedPingOneConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_ping_one_connection` resource is deprecated. Use the `pingfederate_pingone_connection` resource instead."
}

// Metadata returns the resource type name.
func (r *deprecatedPingOneConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ping_one_connection"
}

func (r *deprecatedPingOneConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *deprecatedPingOneConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *deprecatedPingOneConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deprecatedPingOneConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *deprecatedPingOneConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *deprecatedPingOneConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
