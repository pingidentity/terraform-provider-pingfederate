package serversettingssystemkeys

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsSystemKeysResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsSystemKeysResource{}
	_ resource.ResourceWithImportState = &serverSettingsSystemKeysResource{}

	systemKeyAttrTypes = map[string]attr.Type{
		"creation_date":      types.StringType,
		"encrypted_key_data": types.StringType,
		"key_data":           types.StringType,
	}

	creationTimeDefault   = "0001-01-01T00:00:00Z"
	previousKeyDefault, _ = types.ObjectValue(systemKeyAttrTypes, map[string]attr.Value{
		"creation_date":      types.StringValue(creationTimeDefault),
		"encrypted_key_data": types.StringValue(""),
		"key_data":           types.StringValue(""),
	})
)

// ServerSettingsSystemKeysResource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysResource() resource.Resource {
	return &serverSettingsSystemKeysResource{}
}

// serverSettingsSystemKeysResource is the resource implementation.
type serverSettingsSystemKeysResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsSystemKeysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the system encryption keys.",
		Attributes: map[string]schema.Attribute{
			"current": schema.SingleNestedAttribute{
				Description: "Current SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_data")),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("encrypted_key_data")),
						},
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Previous SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(previousKeyDefault),
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
						Default:     stringdefault.StaticString(creationTimeDefault),
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_data")),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("encrypted_key_data")),
						},
					},
				},
			},
			"pending": schema.SingleNestedAttribute{
				Description: "Pending SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_data")),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("encrypted_key_data")),
						},
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addServerSettingsSystemKeysFields(ctx context.Context, addRequest *client.SystemKeys, plan serverSettingsSystemKeysModel) {

	if internaltypes.IsDefined(plan.Current) {
		currentAttrs := plan.Current.Attributes()
		encryptedKeyDataAttrcurrent := currentAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrcurrent) {
			addRequest.Current = *client.NewSystemKey()
			currentEncryptedKeyData := encryptedKeyDataAttrcurrent.ValueString()
			addRequest.Current.EncryptedKeyData = &currentEncryptedKeyData
		}
	}
	if internaltypes.IsDefined(plan.Previous) {
		previousAttrs := plan.Previous.Attributes()
		encryptedKeyDataAttrPrevious := previousAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrPrevious) {
			addRequest.Previous = client.NewSystemKey()
			previousEncryptedKeyData := encryptedKeyDataAttrPrevious.ValueString()
			addRequest.Previous.EncryptedKeyData = &previousEncryptedKeyData
		}
	}
	if internaltypes.IsDefined(plan.Pending) {
		pendingAttrs := plan.Pending.Attributes()
		encryptedKeyDataAttrPending := pendingAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrPending) {
			addRequest.Pending = *client.NewSystemKey()
			pendingEncryptedKeyData := encryptedKeyDataAttrPending.ValueString()
			addRequest.Pending.EncryptedKeyData = &pendingEncryptedKeyData
		}
	}
}

// Metadata returns the resource type name.
func (r *serverSettingsSystemKeysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys"
}

func (r *serverSettingsSystemKeysResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *serverSettingsSystemKeysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsSystemKeysModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)

	apiCreateServerSettingsSystemKeys := r.apiClient.ServerSettingsAPI.UpdateSystemKeys(config.DetermineAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsSystemKeysModel
	diags = readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsSystemKeysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsSystemKeysModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsSystemKeys, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.DetermineAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsSystemKeysResponse(ctx, apiReadServerSettingsSystemKeys, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsSystemKeysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsSystemKeysModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)

	apiCreateServerSettingsSystemKeys := r.apiClient.ServerSettingsAPI.UpdateSystemKeys(config.DetermineAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsSystemKeysModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state, id)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsSystemKeysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsSystemKeysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//  id  doesn't matter because it is a singleton resource.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
