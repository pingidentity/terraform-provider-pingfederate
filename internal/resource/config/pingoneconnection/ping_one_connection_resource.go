package pingoneconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pingOneConnectionResource{}
	_ resource.ResourceWithConfigure   = &pingOneConnectionResource{}
	_ resource.ResourceWithImportState = &pingOneConnectionResource{}
)

// PingOneConnectionResource is a helper function to simplify the provider implementation.
func PingOneConnectionResource() resource.Resource {
	return &pingOneConnectionResource{}
}

// pingOneConnectionResource is the resource implementation.
type pingOneConnectionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type pingOneConnectionResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	ConnectionId                     types.String `tfsdk:"connection_id"`
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	Active                           types.Bool   `tfsdk:"active"`
	Credential                       types.String `tfsdk:"credential"`
	EncryptedCredential              types.String `tfsdk:"encrypted_credential"`
	CredentialId                     types.String `tfsdk:"credential_id"`
	PingOneConnectionId              types.String `tfsdk:"ping_one_connection_id"`
	EnvironmentId                    types.String `tfsdk:"environment_id"`
	CreationDate                     types.String `tfsdk:"creation_date"`
	OrganizationName                 types.String `tfsdk:"organization_name"`
	Region                           types.String `tfsdk:"region"`
	PingOneManagementApiEndpoint     types.String `tfsdk:"ping_one_management_api_endpoint"`
	PingOneAuthenticationApiEndpoint types.String `tfsdk:"ping_one_authentication_api_endpoint"`
}

// GetSchema defines the schema for the resource.
func (r *pingOneConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Ping One Connection",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the PingOne Connection",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the PingOne Connection",
				Optional:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the PingOne Connection is active",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"credential": schema.StringAttribute{
				Description: "The credential for the PingOne connection.",
				Optional:    true,
			},
			"encrypted_credential": schema.StringAttribute{
				Description: "The encrypted credential for the PingOne connection.",
				Optional:    true,
			},
			"credential_id": schema.StringAttribute{
				Description: "The ID of the PingOne credential. This field is read only.",
				Optional:    false,
				Computed:    true,
			},
			"ping_one_connection_id": schema.StringAttribute{
				Description: "The ID of the PingOne connection. This field is read only.",
				Optional:    false,
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment of the PingOne credential. This field is read only.",
				Optional:    false,
				Computed:    true,
			},
			"creation_date": schema.StringAttribute{
				Description: "The creation date of the PingOne connection. This field is read only.",
				Optional:    false,
				Computed:    true,
			},
			"organization_name": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The name of the organization associated with this PingOne connection. This field is read only.",
			},
			"region": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The region of the PingOne connection. This field is read only.",
			},
			"ping_one_management_api_endpoint": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The PingOne Management API endpoint. This field is read only.",
			},
			"ping_one_authentication_api_endpoint": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The PingOne Authentication API endpoint. This field is read only.",
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "connection_id", false, "The persistent, unique ID of the connection.")

	resp.Schema = schema
}

func addOptionalPingOneConnectionFields(ctx context.Context, addRequest *client.PingOneConnection, plan pingOneConnectionResourceModel) error {
	addRequest.Id = plan.Id.ValueStringPointer()
	addRequest.Description = plan.Description.ValueStringPointer()
	addRequest.Active = plan.Active.ValueBoolPointer()
	addRequest.Credential = plan.Credential.ValueStringPointer()
	addRequest.EncryptedCredential = plan.EncryptedCredential.ValueStringPointer()
	addRequest.CredentialId = plan.CredentialId.ValueStringPointer()
	addRequest.PingOneConnectionId = plan.PingOneConnectionId.ValueStringPointer()
	addRequest.EnvironmentId = plan.EnvironmentId.ValueStringPointer()
	addRequest.OrganizationName = plan.OrganizationName.ValueStringPointer()
	addRequest.Region = plan.Region.ValueStringPointer()
	addRequest.PingOneManagementApiEndpoint = plan.PingOneManagementApiEndpoint.ValueStringPointer()
	addRequest.PingOneAuthenticationApiEndpoint = plan.PingOneAuthenticationApiEndpoint.ValueStringPointer()

	return nil

}

// Metadata returns the resource type name.
func (r *pingOneConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ping_one_connection"
}

func (r *pingOneConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readPingOneConnectionResponse(ctx context.Context, r *client.PingOneConnection, state *pingOneConnectionResourceModel) {
	state.Id = types.StringPointerValue(r.Id)
	state.ConnectionId = types.StringPointerValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringPointerValue(r.Description)
	state.Active = types.BoolPointerValue(r.Active)
	state.Credential = types.StringPointerValue(r.Credential)
	state.EncryptedCredential = types.StringPointerValue(r.EncryptedCredential)
	state.CredentialId = types.StringPointerValue(r.CredentialId)
	state.PingOneConnectionId = types.StringPointerValue(r.PingOneConnectionId)
	state.EnvironmentId = types.StringPointerValue(r.EnvironmentId)
	state.CreationDate = types.StringValue(r.CreationDate.Format(time.RFC3339))
	state.OrganizationName = types.StringPointerValue(r.OrganizationName)
	state.Region = types.StringPointerValue(r.Region)
	state.PingOneManagementApiEndpoint = types.StringPointerValue(r.PingOneManagementApiEndpoint)
	state.PingOneAuthenticationApiEndpoint = types.StringPointerValue(r.PingOneAuthenticationApiEndpoint)
}

func (r *pingOneConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pingOneConnectionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPingOneConnection := client.NewPingOneConnection(plan.Name.ValueString())
	err := addOptionalPingOneConnectionFields(ctx, createPingOneConnection, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the PingOne Connection", err.Error())
		return
	}

	apiCreatePingOneConnection := r.apiClient.PingOneConnectionsAPI.CreatePingOneConnection(config.DetermineAuthContext(ctx, r.providerConfig))
	apiCreatePingOneConnection = apiCreatePingOneConnection.Body(*createPingOneConnection)
	pingOneConnectionResponse, httpResp, err := r.apiClient.PingOneConnectionsAPI.CreatePingOneConnectionExecute(apiCreatePingOneConnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the the PingOne Connection", err, httpResp)
		return
	}

	// Read the response into the state
	var state pingOneConnectionResourceModel

	readPingOneConnectionResponse(ctx, pingOneConnectionResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *pingOneConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pingOneConnectionResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadPingOneConnection, httpResp, err := r.apiClient.PingOneConnectionsAPI.GetPingOneConnection(config.DetermineAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the the PingOne Connection", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  the PingOne Connection", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readPingOneConnectionResponse(ctx, apiReadPingOneConnection, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pingOneConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan pingOneConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePingOneConnection := r.apiClient.PingOneConnectionsAPI.UpdatePingOneConnection(config.DetermineAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewPingOneConnection(plan.Name.ValueString())
	err := addOptionalPingOneConnectionFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the PingOne Connection", err.Error())
		return
	}

	updatePingOneConnection = updatePingOneConnection.Body(*createUpdateRequest)
	updatePingOneConnectionResponse, httpResp, err := r.apiClient.PingOneConnectionsAPI.UpdatePingOneConnectionExecute(updatePingOneConnection)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the PingOne Connection", err, httpResp)
		return
	}

	// Read the response
	var state pingOneConnectionResourceModel
	readPingOneConnectionResponse(ctx, updatePingOneConnectionResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *pingOneConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pingOneConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.PingOneConnectionsAPI.DeletePingOneConnection(config.DetermineAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the PingOne Connection", err, httpResp)
	}
}

func (r *pingOneConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("connection_id"), req, resp)
}