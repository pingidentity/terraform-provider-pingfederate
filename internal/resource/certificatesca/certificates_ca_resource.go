package certificates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &certificatesResource{}
	_ resource.ResourceWithConfigure = &certificatesResource{}
)

// CertificateResource is a helper function to simplify the provider implementation.
func CertificateResource() resource.Resource {
	return &certificatesResource{}
}

// certificatesResource is the resource implementation.
type certificatesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type certificatesResourceModel struct {
	Id             types.String `tfsdk:"id"`
	FileData       types.String `tfsdk:"file_data"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the resource.
func (r *certificatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificatesResourceSchema(ctx, req, resp, false)
}

func certificatesResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages CetrificateCA Import.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"crypto_provider": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_data": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"file_data"})
	}
	resp.Schema = schema
}
func addOptionalCaCertsFields(ctx context.Context, addRequest *client.X509File, plan certificatesResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CryptoProvider) {
		addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
	}
	return nil
}

// Metadata returns the resource type name.
func (r *certificatesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates_ca"
}

func (r *certificatesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

// Modify plan to check if crypto_provider attribute is present in the terraform file and act accordingly
func (r *certificatesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model certificatesResourceModel
	var path path.Path
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.CryptoProvider) {
		resp.Diagnostics.AddAttributeWarning(path, "The Crypto Provider is not applicable if Hybrid HSM mode is false or if the provider is SafeNet Luna.",
			"Please remove Crypto Provider from terraform configuration if Hybrid HSM mode or Safenet Provider are not used.")
	}
}
func readCertificateResponse(ctx context.Context, r *client.CertView, state *certificatesResourceModel, expectedValues *certificatesResourceModel, diagnostics *diag.Diagnostics, createPlan types.String) {
	X509FileData := createPlan
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.CryptoProvider = internaltypes.StringTypeOrNil(r.CryptoProvider, false)
	state.FileData = types.StringValue(X509FileData.ValueString())
}

func (r *certificatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan certificatesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createCertificate := client.NewX509File((plan.FileData.ValueString()))
	err := addOptionalCaCertsFields(ctx, createCertificate, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for CA Certificates", err.Error())
		return
	}
	requestJson, err := createCertificate.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiCreateCertificate := r.apiClient.CertificatesCaApi.ImportTrustedCA(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateCertificate = apiCreateCertificate.Body(*createCertificate)
	certificateResponse, httpResp, err := r.apiClient.CertificatesCaApi.ImportTrustedCAExecute(apiCreateCertificate)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating a CA Certificate", err, httpResp)
		return
	}
	responseJson, err := certificateResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificatesResourceModel

	readCertificateResponse(ctx, certificateResponse, &state, &plan, &resp.Diagnostics, plan.FileData)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificate(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCertificate(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state certificatesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadCertificate, httpResp, err := apiClient.CertificatesCaApi.GetTrustedCert(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a Certificate", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadCertificate.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCertificateResponse(ctx, apiReadCertificate, &state, &state, &resp.Diagnostics, state.FileData)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
func (r *certificatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *certificatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteCertificate(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteCertificate(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state certificatesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := apiClient.CertificatesCaApi.DeleteTrustedCA(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a CA Certificate", err, httpResp)
		return
	}

}
