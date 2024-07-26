// Code generated by ping-terraform-plugin-framework-generator

package protocolmetadatasigningsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &protocolMetadataSigningSettingsResource{}
	_ resource.ResourceWithConfigure   = &protocolMetadataSigningSettingsResource{}
	_ resource.ResourceWithImportState = &protocolMetadataSigningSettingsResource{}
)

func ProtocolMetadataSigningSettingsResource() resource.Resource {
	return &protocolMetadataSigningSettingsResource{}
}

type protocolMetadataSigningSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *protocolMetadataSigningSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protocol_metadata_signing_settings"
}

func (r *protocolMetadataSigningSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type protocolMetadataSigningSettingsResourceModel struct {
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	SigningKeyRef      types.Object `tfsdk:"signing_key_ref"`
}

func (r *protocolMetadataSigningSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage the certificate ID and algorithm used for metadata signing.",
		Attributes: map[string]schema.Attribute{
			"signature_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "Signature algorithm. If this property is unset, the default signature algorithm for the key algorithm will be used. Supported signature algorithms are available through the /keyPairs/keyAlgorithms endpoint.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"signing_key_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:    true,
				Description: "Reference to the key used for metadata signing. Refer to /keyPair/signing to get the list of available signing key pairs.",
			},
		},
	}
}

func (model *protocolMetadataSigningSettingsResourceModel) buildClientStruct() (*client.MetadataSigningSettings, diag.Diagnostics) {
	result := &client.MetadataSigningSettings{}
	// signature_algorithm
	result.SignatureAlgorithm = model.SignatureAlgorithm.ValueStringPointer()
	// signing_key_ref
	if !model.SigningKeyRef.IsNull() {
		signingKeyRefValue := &client.ResourceLink{}
		signingKeyRefAttrs := model.SigningKeyRef.Attributes()
		signingKeyRefValue.Id = signingKeyRefAttrs["id"].(types.String).ValueString()
		result.SigningKeyRef = signingKeyRefValue
	}

	return result, nil
}

func (state *protocolMetadataSigningSettingsResourceModel) readClientResponse(response *client.MetadataSigningSettings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// signature_algorithm
	state.SignatureAlgorithm = types.StringPointerValue(response.SignatureAlgorithm)
	// signing_key_ref
	signingKeyRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var signingKeyRefValue types.Object
	if response.SigningKeyRef == nil {
		signingKeyRefValue = types.ObjectNull(signingKeyRefAttrTypes)
	} else {
		signingKeyRefValue, diags = types.ObjectValue(signingKeyRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.SigningKeyRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.SigningKeyRef = signingKeyRefValue
	return respDiags
}

// Set all non-primitive attributes to null with appropriate attribute types
func (r *protocolMetadataSigningSettingsResource) emptyModel() protocolMetadataSigningSettingsResourceModel {
	var model protocolMetadataSigningSettingsResourceModel
	// signing_key_ref
	signingKeyRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	model.SigningKeyRef = types.ObjectNull(signingKeyRefAttrTypes)
	return model
}

func (r *protocolMetadataSigningSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data protocolMetadataSigningSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ProtocolMetadataAPI.UpdateSigningSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ProtocolMetadataAPI.UpdateSigningSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the protocolMetadataSigningSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *protocolMetadataSigningSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data protocolMetadataSigningSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.ProtocolMetadataAPI.GetSigningSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the protocolMetadataSigningSettings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the protocolMetadataSigningSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *protocolMetadataSigningSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data protocolMetadataSigningSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ProtocolMetadataAPI.UpdateSigningSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ProtocolMetadataAPI.UpdateSigningSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the protocolMetadataSigningSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *protocolMetadataSigningSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	emptyState := r.emptyModel()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
