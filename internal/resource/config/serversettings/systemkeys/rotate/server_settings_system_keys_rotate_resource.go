package serversettingssystemkeysrotate

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsSystemKeysRotateResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsSystemKeysRotateResource{}
	_ resource.ResourceWithImportState = &serverSettingsSystemKeysRotateResource{}
)

// ServerSettingsSystemKeysRotateResource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysRotateResource() resource.Resource {
	return &serverSettingsSystemKeysRotateResource{}
}

// serverSettingsSystemKeysRotateResource is the resource implementation.
type serverSettingsSystemKeysRotateResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsSystemKeysRotateResourceModel struct {
	Current               types.Object `tfsdk:"current"`
	Previous              types.Object `tfsdk:"previous"`
	Pending               types.Object `tfsdk:"pending"`
	RotationTriggerValues types.Map    `tfsdk:"rotation_trigger_values"`
}

// Metadata returns the resource type name.
func (r *serverSettingsSystemKeysRotateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys_rotate"
}

func (r *serverSettingsSystemKeysRotateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

// GetSchema defines the schema for the resource.
func (r *serverSettingsSystemKeysRotateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource that handles rotating the system keys.",
		Attributes: map[string]schema.Attribute{
			"current": schema.SingleNestedAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The current secret.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"encrypted_key_data": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The system key encrypted.",
					},
					"creation_date": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "Creation time of the key.",
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Previously used secret.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"encrypted_key_data": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The system key encrypted.",
					},
					"creation_date": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "Creation time of the key.",
					},
				},
			},
			"pending": schema.SingleNestedAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The next secret.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"encrypted_key_data": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The system key encrypted.",
					},
					"creation_date": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "Creation time of the key.",
					},
				},
			},
			"rotation_trigger_values": schema.MapAttribute{
				Description: "A meta-argument map of values that, if any values are changed, will force rotation of the system keys. Adding values to and removing values from the map will not trigger a key rotation. This parameter can be used to control time-based rotation using Terraform.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Rotate the keys via RequiresReplace when the trigger values change
func (r *serverSettingsSystemKeysRotateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state types.Map
	var planValues, stateValues map[string]attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rotation_trigger_values"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues = plan.Elements()

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("rotation_trigger_values"), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateValues = state.Elements()

	for k, v := range planValues {
		if stateValue, ok := stateValues[k]; ok && (v == types.StringUnknown() || !stateValue.Equal(v)) {
			resp.RequiresReplace = path.Paths{path.Root("rotation_trigger_values")}
			break
		}
	}
}

func (state *serverSettingsSystemKeysRotateResourceModel) readClientResponse(response *client.SystemKeys) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	keyAttrTypes := map[string]attr.Type{
		"encrypted_key_data": types.StringType,
		"creation_date":      types.StringType,
	}

	var currentCreationDateValue types.String
	if response.Current.CreationDate == nil {
		currentCreationDateValue = types.StringNull()
	} else {
		currentCreationDateValue = types.StringValue(response.Current.CreationDate.Format(time.RFC3339))
	}
	state.Current, diags = types.ObjectValue(keyAttrTypes, map[string]attr.Value{
		"encrypted_key_data": types.StringPointerValue(response.Current.EncryptedKeyData),
		"creation_date":      currentCreationDateValue,
	})
	respDiags.Append(diags...)

	if response.Previous != nil {
		var previousCreationDateValue types.String
		if response.Previous.CreationDate == nil {
			previousCreationDateValue = types.StringNull()
		} else {
			previousCreationDateValue = types.StringValue(response.Previous.CreationDate.Format(time.RFC3339))
		}
		state.Previous, diags = types.ObjectValue(keyAttrTypes, map[string]attr.Value{
			"encrypted_key_data": types.StringPointerValue(response.Previous.EncryptedKeyData),
			"creation_date":      previousCreationDateValue,
		})
		respDiags.Append(diags...)
	} else {
		state.Previous = types.ObjectNull(keyAttrTypes)
	}

	var pendingCreationDateValue types.String
	if response.Pending.CreationDate == nil {
		pendingCreationDateValue = types.StringNull()
	} else {
		pendingCreationDateValue = types.StringValue(response.Pending.CreationDate.Format(time.RFC3339))
	}
	state.Pending, diags = types.ObjectValue(keyAttrTypes, map[string]attr.Value{
		"encrypted_key_data": types.StringPointerValue(response.Pending.EncryptedKeyData),
		"creation_date":      pendingCreationDateValue,
	})
	respDiags.Append(diags...)

	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *serverSettingsSystemKeysRotateResource) emptyModel() serverSettingsSystemKeysRotateResourceModel {
	var model serverSettingsSystemKeysRotateResourceModel
	// keys
	keyAttrTypes := map[string]attr.Type{
		"encrypted_key_data": types.StringType,
		"creation_date":      types.StringType,
	}
	model.Current = types.ObjectNull(keyAttrTypes)
	model.Previous = types.ObjectNull(keyAttrTypes)
	model.Pending = types.ObjectNull(keyAttrTypes)
	model.RotationTriggerValues = types.MapNull(types.StringType)
	return model
}

func (r *serverSettingsSystemKeysRotateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state serverSettingsSystemKeysRotateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.RotateSystemKeys(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while rotating the system keys", err, httpResp)
		return
	}

	// Read the response into the state, maintaining the trigger values
	resp.Diagnostics.Append(state.readClientResponse(serverSettingsGeneralSettingsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *serverSettingsSystemKeysRotateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverSettingsSystemKeysRotateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the system keys", err, httpResp)
		return
	}

	// Read the response into the state, maintaining the trigger values
	resp.Diagnostics.Append(state.readClientResponse(serverSettingsGeneralSettingsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *serverSettingsSystemKeysRotateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This will only happen when adding or removing rotation trigger values. Just copy the plan into state.
	resp.State.Raw = req.Plan.Raw
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsSystemKeysRotateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsSystemKeysRotateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
