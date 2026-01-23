// Copyright © 2025 Ping Identity Corporation

package configurationencryptionkeysrotate

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &configurationEncryptionKeysRotateResource{}
	_ resource.ResourceWithConfigure   = &configurationEncryptionKeysRotateResource{}
	_ resource.ResourceWithImportState = &configurationEncryptionKeysRotateResource{}
)

// ConfigurationEncryptionKeysRotateResource is a helper function to simplify the provider implementation.
func ConfigurationEncryptionKeysRotateResource() resource.Resource {
	return &configurationEncryptionKeysRotateResource{}
}

// configurationEncryptionKeysRotateResource is the resource implementation.
type configurationEncryptionKeysRotateResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type configurationEncryptionKeysRotateResourceModel struct {
	Keys                  types.List `tfsdk:"keys"`
	RotationTriggerValues types.Map  `tfsdk:"rotation_trigger_values"`
}

// Metadata returns the resource type name.
func (r *configurationEncryptionKeysRotateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_encryption_keys_rotate"
}

func (r *configurationEncryptionKeysRotateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

// GetSchema defines the schema for the resource.
func (r *configurationEncryptionKeysRotateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to handle rotating the current configuration encryption keys.",
		Attributes: map[string]schema.Attribute{
			"keys": schema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key_id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the key.",
						},
						"creation_date": schema.StringAttribute{
							Computed:    true,
							Description: "The creation date of the key.",
						},
					},
				},
				Description: "The list of Configuration Encryption Keys.",
			},
			"rotation_trigger_values": schema.MapAttribute{
				Description: "A meta-argument map of values that, if any values are changed, will force rotation of the encryption keys. Adding values to and removing values from the map will not trigger a rotation. This parameter can be used to control time-based rotation using Terraform.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Rotate the keys via RequiresReplace when the trigger values change
func (r *configurationEncryptionKeysRotateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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

func (state *configurationEncryptionKeysRotateResourceModel) readClientResponse(response *client.ConfigurationEncryptionKeys) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	keyAttrTypes := map[string]attr.Type{
		"key_id":        types.StringType,
		"creation_date": types.StringType,
	}
	var keys []attr.Value
	for _, key := range response.Items {
		var creationDateValue types.String
		if key.CreationDate == nil {
			creationDateValue = types.StringNull()
		} else {
			creationDateValue = types.StringValue(key.CreationDate.Format(time.RFC3339))
		}
		keyAttr, diags := types.ObjectValue(keyAttrTypes, map[string]attr.Value{
			"key_id":        types.StringPointerValue(key.KeyId),
			"creation_date": creationDateValue,
		})
		respDiags.Append(diags...)
		keys = append(keys, keyAttr)
	}
	state.Keys, diags = types.ListValue(types.ObjectType{AttrTypes: keyAttrTypes}, keys)
	respDiags.Append(diags...)
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *configurationEncryptionKeysRotateResource) emptyModel() configurationEncryptionKeysRotateResourceModel {
	var model configurationEncryptionKeysRotateResourceModel
	// keys
	keyAttrTypes := map[string]attr.Type{
		"key_id":        types.StringType,
		"creation_date": types.StringType,
	}
	model.Keys = types.ListNull(types.ObjectType{AttrTypes: keyAttrTypes})
	model.RotationTriggerValues = types.MapNull(types.StringType)
	return model
}

func (r *configurationEncryptionKeysRotateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state configurationEncryptionKeysRotateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ConfigurationEncryptionKeysAPI.RotateConfigurationEncryptionKey(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while rotating the encryption keys", err, httpResp)
		return
	}

	// Read the response into the state, maintaining the trigger values
	resp.Diagnostics.Append(state.readClientResponse(serverSettingsGeneralSettingsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *configurationEncryptionKeysRotateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configurationEncryptionKeysRotateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverSettingsGeneralSettingsResponse, httpResp, err := r.apiClient.ConfigurationEncryptionKeysAPI.GetConfigurationEncryptionKeys(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the encryption keys", err, httpResp)
		return
	}

	// Read the response into the state, maintaining the trigger values
	resp.Diagnostics.Append(state.readClientResponse(serverSettingsGeneralSettingsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *configurationEncryptionKeysRotateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This will only happen when adding or removing rotation trigger values. Just copy the plan into state.
	resp.State.Raw = req.Plan.Raw
}

// This config object is edit-only, so Terraform can't delete it.
func (r *configurationEncryptionKeysRotateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *configurationEncryptionKeysRotateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
