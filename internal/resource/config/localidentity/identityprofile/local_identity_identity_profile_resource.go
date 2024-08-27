package localidentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localIdentityIdentityProfileResource{}
	_ resource.ResourceWithConfigure   = &localIdentityIdentityProfileResource{}
	_ resource.ResourceWithImportState = &localIdentityIdentityProfileResource{}
)

// LocalIdentityIdentityProfileResource is a helper function to simplify the provider implementation.
func LocalIdentityIdentityProfileResource() resource.Resource {
	return &localIdentityIdentityProfileResource{
		impl: localIdentityProfileResource{},
	}
}

// localIdentityIdentityProfileResource is the resource implementation.
type localIdentityIdentityProfileResource struct {
	impl localIdentityProfileResource
}

// GetSchema defines the schema for the resource.
func (r *localIdentityIdentityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_local_identity_identity_profile` resource is deprecated. Use the `pingfederate_local_identity_profile` resource instead."
}

// Metadata returns the resource type name.
func (r *localIdentityIdentityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_identity_identity_profile"
}

func (r *localIdentityIdentityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *localIdentityIdentityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *localIdentityIdentityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *localIdentityIdentityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *localIdentityIdentityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *localIdentityIdentityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
