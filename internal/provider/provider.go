package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/administrativeaccounts"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificates"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/protocolmetadata"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfacesß
var (
	_ provider.Provider = &pingfederateProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &pingfederateProvider{}
}

// PingFederate ProviderModel maps provider schema data to a Go type.
type pingfederateProviderModel struct {
	HttpsHost types.String `tfsdk:"https_host"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
}

// pingfederateProvider is the provider implementation.
type pingfederateProvider struct{}

// Metadata returns the provider type name.
func (p *pingfederateProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingfederate"
}

// GetSchema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *pingfederateProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"https_host": schema.StringAttribute{
				MarkdownDescription: "URI for PingFederate HTTPS port",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for PingFederate Admin user",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for PingFederate Admin user",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *pingfederateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config pingfederateProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// User must provide a https host to the provider
	var httpsHost string
	if config.HttpsHost.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as https_host",
		)
	} else {
		if config.HttpsHost.IsNull() {
			httpsHost = os.Getenv("PINGFEDERATE_PROVIDER_HTTPS_HOST")
		} else {
			httpsHost = config.HttpsHost.ValueString()
		}
		if httpsHost == "" {
			resp.Diagnostics.AddError(
				"Unable to find https_host",
				"https_host cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_HTTPS_HOST environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var username string
	if config.Username.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as username",
		)
	} else {
		if config.Username.IsNull() {
			username = os.Getenv("PINGFEDERATE_PROVIDER_USERNAME")
		} else {
			username = config.Username.ValueString()
		}
		if username == "" {
			resp.Diagnostics.AddError(
				"Unable to find username",
				"username cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_USERNAME environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var password string
	if config.Password.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as password",
		)
	} else {
		if config.Password.IsNull() {
			password = os.Getenv("PINGFEDERATE_PROVIDER_PASSWORD")
		} else {
			password = config.Password.ValueString()
		}
		if password == "" {
			resp.Diagnostics.AddError(
				"Unable to find password",
				"password cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_PASSWORD environment variable.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the PingFederate config and API client info available during DataSource and Resource
	// type Configure methods.
	var resourceConfig internaltypes.ResourceConfiguration
	providerConfig := internaltypes.ProviderConfiguration{
		HttpsHost: httpsHost,
		Username:  username,
		Password:  password,
	}
	resourceConfig.ProviderConfig = providerConfig
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + "/pf-admin-api/v1",
		},
	}
	// TODO THIS IS NOT SAFE!! Eventually need to add way to trust a specific cert/signer here rather than just trusting everything
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	// #nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	clientConfig.HTTPClient = httpClient
	resourceConfig.ApiClient = client.NewAPIClient(clientConfig)
	resp.ResourceData = resourceConfig

	tflog.Info(ctx, "Configured PingFederate client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingfederateProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *pingfederateProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		administrativeaccounts.AdministrativeAccountResource,
		authenticationapi.AuthenticationApiSettingsResource,
		certificates.CertificateResource,
		config.AuthenticationPolicyContractsResource,
		config.RedirectValidationResource,
		config.VirtualHostNamesResource,
		idp.IdpDefaultUrlsResource,
		keypairs.KeyPairsSigningImportResource,
		keypairs.KeyPairsSslServerImportResource,
		license.LicenseAgreementResource,
		license.LicenseResource,
		oauth.OauthAuthServerSettingsScopesCommonScopesResource,
		oauth.OauthAuthServerSettingsScopesExclusiveScopesResource,
		oauth.OauthIssuersResource,
		protocolmetadata.ProtocolMetadataLifetimeSettingsResource,
		serversettings.ServerSettingsGeneralSettingsResource,
		session.SessionApplicationSessionPolicyResource,
		session.SessionAuthenticationSessionPoliciesGlobalResource,
		session.SessionSettingsResource,
	}
}
