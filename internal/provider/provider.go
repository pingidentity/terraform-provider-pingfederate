// Copyright © 2025 Ping Identity Corporation

package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/administrativeaccount"
	authenticationapiapplication "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi/application"
	authenticationapisettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi/settings"
	authenticationpolicies "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicies"
	authenticationpoliciesfragments "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicies/fragments"
	authenticationpoliciessettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicies/settings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicycontract"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationselector"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/captchaproviders"
	captchaproviderssettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/captchaproviders/settings"
	certificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificates/ca"
	certificatesgroups "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificates/groups"
	certificatesrevocationocspcertificates "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificates/revocation/ocspcertificates"
	certificatesrevocationsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificates/revocation/settings"
	clustersettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/cluster/settings"
	clusterstatus "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/cluster/status"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/configstore"
	configurationencryptionkeysrotate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/configurationencryptionkeys/rotate"
	connectionmetadataexport "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/connectionmetadata/export"
	datastore "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/datastore"
	extendedproperties "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/extendedproperties"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/identitystoreprovisioners"
	idpadapter "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/adapter"
	idpdefaulturls "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/defaulturls"
	idpspconnection "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/spconnection"
	idpstsrequestparameterscontracts "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/stsrequestparameterscontracts"
	idptokenprocessors "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/tokenprocessors"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idptospadaptermapping"
	incomingproxysettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/incomingproxysettings"
	kerberosrealms "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/kerberos/realms"
	kerberosrealmssettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/kerberos/realms/settings"
	keypairsoauthopenidconnect "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/oauthopenidconnect"
	keypairsoauthopenidconnectadditionalkeysets "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/oauthopenidconnect/additionalkeysets"
	keypairsigning "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/signing"
	keypairssigningcertificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/signing/certificate"
	keypairssigningrotationsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/signing/rotationsettings"
	keypairssslclient "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslclient"
	keypairssslclientcertificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslclient/certificate"
	keypairssslclientcsr "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslclient/csr"
	keypairssslserver "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslserver"
	keypairssslservercertificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslserver/certificate"
	keypairssslservercsr "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslserver/csr"
	keypairssslserversettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypairs/sslserver/settings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license"
	licenseagreement "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license/agreement"
	localidentity "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/localidentity/identityprofile"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/metadataurls"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/notificationpublishers"
	notificationpublisherssettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/notificationpublishers/settings"
	oauthaccesstokenmanager "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/accesstokenmanager"
	oauthaccesstokenmanagerssettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/accesstokenmanagers/settings"
	oauthaccesstokenmapping "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/accesstokenmapping"
	oauthauthenticationpolicycontractmappings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authenticationpolicycontractmappings"
	oauthauthserversettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings"
	oauthcibaserverpolicyrequestpolicies "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/cibaserverpolicy/requestpolicies"
	oauthcibaserverpolicysettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/cibaserverpolicy/settings"
	oauthclient "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/client"
	oauthclientregistrationpolicies "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/clientregistrationpolicies"
	oauthclientsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/clientsettings"
	oauthidpadaptermappings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/idpadaptermappings"
	oauthissuer "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/issuer"
	oauthopenidconnectpolicy "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/openidconnect/policy"
	oauthopenidconnectsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/openidconnect/settings"
	oauthoutofbandauthplugins "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/outofbandauthplugins"
	oauthresourceownercredentialsmappings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/resourceownercredentialsmappings"
	oauthtokenexchangegeneratorsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/tokenexchange/generator/settings"
	oauthtokenexchangetokengeneratormapping "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/tokenexchange/tokengeneratormapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/passwordcredentialvalidator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/pingoneconnection"
	protocolmetadatalifetimesettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/protocolmetadata/lifetimesettings"
	protocolmetadatasigningsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/protocolmetadata/signingsettings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/redirectvalidation"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/secretmanagers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings"
	serversettingsgeneralsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/generalsettings"
	serversettingslogsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/logsettings"
	serversettingssystemkeys "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/systemkeys"
	serversettingssystemkeysrotate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/systemkeys/rotate"
	serversettingswstruststssettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/wstruststssettings"
	serversettingswstruststssettingsissuercertificates "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/wstruststssettings/issuercertificates"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serviceauthentication"
	sessionapplicationsessionpolicy "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/applicationsessionpolicy"
	sessionauthenticationsessionpolicies "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/authenticationsessionpolicies"
	sessionauthenticationsessionpoliciesglobal "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/authenticationsessionpolicies/global"
	sessionsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/settings"
	spadapters "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/adapters"
	spauthenticationpolicycontractmapping "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/authenticationpolicycontractmapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/defaulturls"
	spidpconnection "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/idpconnection"
	sptargeturlmappings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/targeturlmappings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/tokenprocessortotokengeneratormapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/virtualhostnames"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfacesß
var (
	_ provider.Provider = &pingfederateProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func NewFactory(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pingfederateProvider{
			version: version,
		}
	}
}

// NewTestProvider is a helper function to simplify testing implementation.
func NewTestProvider() provider.Provider {
	return NewFactory("test")()
}

// PingFederate ProviderModel maps provider schema data to a Go type.
type pingfederateProviderModel struct {
	HttpsHost                       types.String `tfsdk:"https_host"`
	AdminApiPath                    types.String `tfsdk:"admin_api_path"`
	Username                        types.String `tfsdk:"username"`
	Password                        types.String `tfsdk:"password"`
	AccessToken                     types.String `tfsdk:"access_token"`
	ClientId                        types.String `tfsdk:"client_id"`
	ClientSecret                    types.String `tfsdk:"client_secret"`
	Scopes                          types.List   `tfsdk:"scopes"`
	TokenUrl                        types.String `tfsdk:"token_url"`
	InsecureTrustAllTls             types.Bool   `tfsdk:"insecure_trust_all_tls"`
	CACertificatePEMFiles           types.Set    `tfsdk:"ca_certificate_pem_files"`
	XBypassExternalValidationHeader types.Bool   `tfsdk:"x_bypass_external_validation_header"`
	ProductVersion                  types.String `tfsdk:"product_version"`
}

// pingfederateProvider is the provider implementation.
type pingfederateProvider struct {
	version string
}

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
				MarkdownDescription: "URI for PingFederate HTTPS port. Default value can be set with the `PINGFEDERATE_PROVIDER_HTTPS_HOST` environment variable.",
				Optional:            true,
			},
			"admin_api_path": schema.StringAttribute{
				MarkdownDescription: "Path for PingFederate Admin API. Default value can be set with the `PINGFEDERATE_PROVIDER_ADMIN_API_PATH` environment variable. If no value is supplied, the value used will be `/pf-admin-api/v1`.",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token for PingFederate Admin API. Cannot be used in conjunction with username and password, or oauth. Default value can be set with the `PINGFEDERATE_PROVIDER_ACCESS_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_secret")),
					stringvalidator.ConflictsWith(path.MatchRoot("scopes")),
					stringvalidator.ConflictsWith(path.MatchRoot("token_url")),
				},
			},
			"ca_certificate_pem_files": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingFederate server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "OAuth client ID for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
					stringvalidator.AlsoRequires(path.MatchRoot("client_secret")),
					stringvalidator.AlsoRequires(path.MatchRoot("token_url")),
				},
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "OAuth client secret for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
					stringvalidator.AlsoRequires(path.MatchRoot("client_id")),
					stringvalidator.AlsoRequires(path.MatchRoot("token_url")),
				},
			},
			"insecure_trust_all_tls": schema.BoolAttribute{
				Description: "Set to true to trust any certificate when connecting to the PingFederate server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.",
				Optional:    true,
			},
			"scopes": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "OAuth scopes for access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_SCOPES` environment variable.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("access_token")),
					listvalidator.ConflictsWith(path.MatchRoot("username")),
					listvalidator.ConflictsWith(path.MatchRoot("password")),
					listvalidator.AlsoRequires(path.MatchRoot("client_id")),
					listvalidator.AlsoRequires(path.MatchRoot("client_secret")),
					listvalidator.AlsoRequires(path.MatchRoot("token_url")),
				},
			},
			"token_url": schema.StringAttribute{
				MarkdownDescription: "OAuth token URL for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
					stringvalidator.AlsoRequires(path.MatchRoot("client_id")),
					stringvalidator.AlsoRequires(path.MatchRoot("client_secret")),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for PingFederate Admin user. Must only be set with password. Cannot be used in conjunction with access_token, or oauth. Default value can be set with the `PINGFEDERATE_PROVIDER_USERNAME` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_secret")),
					stringvalidator.ConflictsWith(path.MatchRoot("scopes")),
					stringvalidator.ConflictsWith(path.MatchRoot("token_url")),
					stringvalidator.AlsoRequires(path.MatchRoot("password")),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for PingFederate Admin user. Must only be set with username. Cannot be used in conjunction with access_token, or oauth.  Default value can be set with the `PINGFEDERATE_PROVIDER_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("client_secret")),
					stringvalidator.ConflictsWith(path.MatchRoot("scopes")),
					stringvalidator.ConflictsWith(path.MatchRoot("token_url")),
					stringvalidator.AlsoRequires(path.MatchRoot("username")),
				},
			},
			"product_version": schema.StringAttribute{
				Description: "Version of the PingFederate server being configured. Default value can be set with the `PINGFEDERATE_PROVIDER_PRODUCT_VERSION` environment variable.",
				Optional:    true,
			},
			"x_bypass_external_validation_header": schema.BoolAttribute{
				Description: "Header value in request for PingFederate. When set to `true`, connectivity checks for resources such as `pingfederate_data_store` will be skipped. Default value can be set with the `PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER` environment variable.",
				Optional:    true,
			},
		},
	}
}

func addAuthAttributeDiagsError(attribute string, authenticationMethod string, envVar string, resp *provider.ConfigureResponse) {
	resp.Diagnostics.AddAttributeError(
		path.Root(attribute),
		providererror.InvalidProviderConfiguration,
		fmt.Sprintf("%s cannot be empty when using %s authentication. Either set it in the configuration or use the %s environment variable.", attribute, authenticationMethod, envVar),
	)
}

func addAttributeRequiredError(attribute, envVar string, diags *diag.Diagnostics) {
	diags.AddAttributeError(
		path.Root(attribute),
		providererror.InvalidProviderConfiguration,
		fmt.Sprintf("%s is required. Either set it in the configuration or use the %s environment variable", attribute, envVar),
	)
}

func addAttributeUnknownError(attribute, envVar string, diags *diag.Diagnostics) {
	diags.AddAttributeError(
		path.Root(attribute),
		providererror.InvalidProviderConfiguration,
		fmt.Sprintf("%s cannot be unknown. It can be set either in the configuration or with the %s environment variable", attribute, envVar),
	)
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
		addAttributeUnknownError("https_host", "PINGFEDERATE_PROVIDER_HTTPS_HOST", &resp.Diagnostics)
	} else {
		if config.HttpsHost.IsNull() {
			httpsHost = os.Getenv("PINGFEDERATE_PROVIDER_HTTPS_HOST")
		} else {
			httpsHost = config.HttpsHost.ValueString()
		}
		if httpsHost == "" {
			addAttributeRequiredError("https_host", "PINGFEDERATE_PROVIDER_HTTPS_HOST", &resp.Diagnostics)
		}
	}

	// User must provide a admin api base path to the provider
	var adminApiPath string
	if config.AdminApiPath.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("admin_api_path", "PINGFEDERATE_PROVIDER_ADMIN_API_PATH", &resp.Diagnostics)
	} else {
		if config.AdminApiPath.IsNull() {
			adminApiPath = os.Getenv("PINGFEDERATE_PROVIDER_ADMIN_API_PATH")
			if adminApiPath == "" {
				adminApiPath = "/pf-admin-api/v1"
			}
		} else {
			adminApiPath = config.AdminApiPath.ValueString()
		}
	}

	// Check if the user has provided a username to the provider
	var username string
	var hasUsername bool = false
	if config.Username.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("username", "PINGFEDERATE_PROVIDER_USERNAME", &resp.Diagnostics)
	} else {
		if config.Username.IsNull() {
			username = os.Getenv("PINGFEDERATE_PROVIDER_USERNAME")
		} else {
			username = config.Username.ValueString()
		}
		if username == "" {
			tflog.Info(ctx, "Unable to find username value")
		} else {
			hasUsername = true
		}
	}

	// Check if the user has provided a password to the provider
	var password string
	var hasPassword bool = false
	if config.Password.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("password", "PINGFEDERATE_PROVIDER_PASSWORD", &resp.Diagnostics)
	} else {
		if config.Password.IsNull() {
			password = os.Getenv("PINGFEDERATE_PROVIDER_PASSWORD")
		} else {
			password = config.Password.ValueString()
		}
		if password == "" {
			tflog.Info(ctx, "Unable to find password value")
		} else {
			hasPassword = true
		}
	}

	// Check if the user has provided an access token to the provider
	var accessToken string
	var hasAccessToken bool = false
	if config.AccessToken.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("access_token", "PINGFEDERATE_PROVIDER_ACCESS_TOKEN", &resp.Diagnostics)
	} else {
		if config.AccessToken.IsNull() {
			accessToken = os.Getenv("PINGFEDERATE_PROVIDER_ACCESS_TOKEN")
		} else {
			accessToken = config.AccessToken.ValueString()
		}
		if accessToken == "" {
			tflog.Info(ctx, "Unable to find access_token value")
		} else {
			hasAccessToken = true
		}
	}

	// Check if the user has provided an OAuth configuration to the provider
	var clientId string
	var hasClientId bool = false
	if config.ClientId.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("client_id", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID", &resp.Diagnostics)
	} else {
		if config.ClientId.IsNull() {
			clientId = os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID")
		} else {
			clientId = config.ClientId.ValueString()
		}
		if clientId == "" {
			tflog.Info(ctx, "Unable to find client_id value")
		} else {
			hasClientId = true
		}
	}

	var clientSecret string
	var hasClientSecret bool = false
	if config.ClientSecret.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("client_secret", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET", &resp.Diagnostics)
	} else {
		if config.ClientSecret.IsNull() {
			clientSecret = os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET")
		} else {
			clientSecret = config.ClientSecret.ValueString()
		}
		if clientSecret == "" {
			tflog.Info(ctx, "Unable to find client_secret value")
		} else {
			hasClientSecret = true
		}
	}

	var scopes []string
	var hasScopes bool = false
	if config.Scopes.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("scopes", "PINGFEDERATE_PROVIDER_OAUTH_SCOPES", &resp.Diagnostics)
	} else {
		if config.Scopes.IsNull() {
			envScopes := os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_SCOPES")
			if strings.Contains(envScopes, ",") {
				scopes = strings.Split(envScopes, ",")
			} else {
				if envScopes != "" {
					scopes = []string{envScopes}
				}
			}
		} else {
			config.Scopes.ElementsAs(ctx, &scopes, false)
		}
		if len(scopes) > 0 {
			hasScopes = true
		} else {
			tflog.Info(ctx, "Unable to find scopes value")
		}
	}

	var tokenUrl string
	var hasTokenUrl bool = false
	if config.TokenUrl.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		addAttributeUnknownError("token_url", "PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL", &resp.Diagnostics)
	} else {
		if config.TokenUrl.IsNull() {
			tokenUrl = os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL")
		} else {
			tokenUrl = config.TokenUrl.ValueString()
		}
		if tokenUrl == "" {
			tflog.Info(ctx, "Unable to find token_url value")
		} else {
			hasTokenUrl = true
		}
	}

	// Validate the configuration
	if !hasUsername && !hasPassword && !hasAccessToken && !hasClientId && !hasClientSecret && !hasScopes && !hasTokenUrl {
		resp.Diagnostics.AddError(
			providererror.InvalidProviderConfiguration,
			"Unable to find username and password, access_token, or OAuth required properties for configuration. "+
				"username and password, access_token, or oauth configuration required values were not supplied. Either set them in the configuration or use the PINGFEDERATE_PROVIDER_* environment variables.",
		)
	}

	// User cannot provide username and password, access token
	if (hasUsername || hasPassword) && hasAccessToken {
		resp.Diagnostics.AddError(
			providererror.InvalidProviderConfiguration,
			"Username and password cannot be used with access_token. "+
				"Only basic authentication (username and password) or access_token can be used. If you want to use access_token, remove username and password from the configuration or use the PINGFEDERATE_PROVIDER_USERNAME and PINGFEDERATE_PROVIDER_PASSWORD environment variables.",
		)
	}

	// User cannot provide username and password, OAuth configuration
	if (hasUsername || hasPassword) && (hasClientId || hasClientSecret || hasScopes || hasTokenUrl) {
		resp.Diagnostics.AddError(
			providererror.InvalidProviderConfiguration,
			"Username and password cannot be used with OAuth configuration properties. "+
				"Only basic authentication (username and password) or OAuth authentication can be used. If you want to use OAuth, remove username and password from the configuration or use the PINGFEDERATE_PROVIDER_USERNAME and PINGFEDERATE_PROVIDER_PASSWORD environment variables.",
		)
	}

	// User cannot provide access token, OAuth configuration
	if hasAccessToken && (hasClientId || hasClientSecret || hasScopes || hasTokenUrl) {
		resp.Diagnostics.AddError(
			providererror.InvalidProviderConfiguration,
			"Access token cannot be used with OAuth configuration "+
				"Only basic authentication (username and password) or access_token can be used. If you want to use basic authentication, remove access_token from the configuration or use the PINGFEDERATE_PROVIDER_ACCESS_TOKEN environment variable.",
		)
	}

	var hasBasicAuth bool = hasUsername || hasPassword
	var hasAccessTokenAuth bool = hasAccessToken
	var hasOauthConfig bool = hasClientId || hasClientSecret || hasTokenUrl
	// If user has not provided an OAuth configuration or access token, they must provide username and password
	if !(hasOauthConfig && hasAccessTokenAuth) && hasBasicAuth {
		if username == "" {
			addAuthAttributeDiagsError("username", "basic", "PINGFEDERATE_PROVIDER_USERNAME", resp)
		}

		if password == "" {
			addAuthAttributeDiagsError("password", "basic", "PINGFEDERATE_PROVIDER_PASSWORD", resp)
		}
	}

	// If user has not provided username and password or an access token, they must provide an OAuth configuration
	if !(hasBasicAuth || hasAccessTokenAuth) && hasOauthConfig {
		if clientId == "" {
			addAuthAttributeDiagsError("client_id", "OAuth", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID", resp)
		}
		if clientSecret == "" {
			addAuthAttributeDiagsError("client_secret", "OAuth", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET", resp)
		}
		if tokenUrl == "" {
			addAuthAttributeDiagsError("token_url", "OAuth", "PINGFEDEATE_PROVIDER_OAUTH_TOKEN_URL", resp)
		}
		if len(scopes) == 0 {
			tflog.Warn(ctx, "No scopes value configured.")
		}
	}

	// User must provide a product version to the provider
	var productVersion string
	var parsedProductVersion version.SupportedVersion
	var err error
	if !config.ProductVersion.IsUnknown() && !config.ProductVersion.IsNull() {
		productVersion = config.ProductVersion.ValueString()
	} else {
		productVersion = os.Getenv("PINGFEDERATE_PROVIDER_PRODUCT_VERSION")
	}

	if productVersion == "" {
		addAttributeRequiredError("product_version", "PINGFEDERATE_PROVIDER_PRODUCT_VERSION", &resp.Diagnostics)
	} else {
		// Validate the PingFederate version
		parsedProductVersion, diags = version.Parse(productVersion)
		resp.Diagnostics.Append(diags...)
	}

	// Optional attributes
	var insecureTrustAllTls bool
	if !config.InsecureTrustAllTls.IsUnknown() && !config.InsecureTrustAllTls.IsNull() {
		insecureTrustAllTls = config.InsecureTrustAllTls.ValueBool()
	} else {
		insecureTrustAllTls, err = strconv.ParseBool(os.Getenv("PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS"))
		if err != nil {
			insecureTrustAllTls = false
			tflog.Info(ctx, "Failed to parse boolean from 'PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS' environment variable, defaulting 'insecure_trust_all_tls' to false")
		}
	}

	var caCertPemFiles []string
	if !config.CACertificatePEMFiles.IsUnknown() && !config.CACertificatePEMFiles.IsNull() {
		config.CACertificatePEMFiles.ElementsAs(ctx, &caCertPemFiles, false)
	} else {
		pemFilesEnvVar := os.Getenv("PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES")
		if len(pemFilesEnvVar) == 0 {
			tflog.Info(ctx, "Did not find any certificate paths specified via the 'PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES' environment variable, using the host's root CA set")
		} else {
			caCertPemFiles = strings.Split(pemFilesEnvVar, ",")
		}
	}

	var caCertPool *x509.CertPool
	if len(caCertPemFiles) == 0 {
		tflog.Info(ctx, "No CA certs specified, using the host's root CA set")
		caCertPool = nil
	} else {
		caCertPool = x509.NewCertPool()
		for _, pemFilename := range caCertPemFiles {
			// Load CA cert
			pemFilename := filepath.Clean(pemFilename)
			caCert, err := os.ReadFile(pemFilename)
			if err != nil {
				resp.Diagnostics.AddError(providererror.InvalidProviderConfiguration,
					"Failed to read CA PEM certificate file: "+pemFilename+". "+err.Error())
			}
			tflog.Info(ctx, "Adding CA cert from file: "+pemFilename)
			if !caCertPool.AppendCertsFromPEM(caCert) {
				resp.Diagnostics.AddWarning(providererror.InvalidProviderConfiguration, "Failed to parse CA PEM certificate from file: "+pemFilename)
			}
		}
	}

	var xBypassExternalValidation bool
	var xBypassExternalValidationErr error
	if !config.XBypassExternalValidationHeader.IsUnknown() && !config.XBypassExternalValidationHeader.IsNull() {
		xBypassExternalValidation = config.XBypassExternalValidationHeader.ValueBool()
	} else {
		xBypassExternalValidation, xBypassExternalValidationErr = strconv.ParseBool(os.Getenv("PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER"))
		if xBypassExternalValidationErr != nil {
			xBypassExternalValidation = false
			tflog.Info(ctx, "Failed to parse boolean from 'PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER' environment variable, defaulting 'x_bypass_external_validation_header' to false")
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// The extra suffix for the user-agent is optional and is not considered a provider parameter.
	// We just use it directly from the environment variable, if set.
	userAgentExtraSuffix := os.Getenv("PINGFEDERATE_TF_APPEND_USER_AGENT")

	// Make the PingFederate config and API client info available during DataSource and Resource
	// type Configure methods.
	var resourceConfig internaltypes.ResourceConfiguration
	providerConfig := internaltypes.ProviderConfiguration{
		HttpsHost:      httpsHost,
		ProductVersion: parsedProductVersion,
	}

	if username != "" {
		providerConfig.Username = &username
	}

	if password != "" {
		providerConfig.Password = &password
	}

	if accessToken != "" {
		providerConfig.AccessToken = &accessToken
	}

	if hasOauthConfig {
		providerConfig.ClientId = &clientId
		providerConfig.ClientSecret = &clientSecret
		providerConfig.Scopes = scopes
		providerConfig.TokenUrl = &tokenUrl
	}

	resourceConfig.ProviderConfig = providerConfig
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.DefaultHeader["X-BypassExternalValidation"] = strconv.FormatBool(xBypassExternalValidation)
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + adminApiPath,
		},
	}
	// #nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTrustAllTls,
			RootCAs:            caCertPool,
		},
	}
	httpClient := &http.Client{Transport: tr}
	resourceConfig.ProviderConfig.Transport = tr
	clientConfig.HTTPClient = httpClient
	userAgentSuffix := fmt.Sprintf("terraform-provider-pingfederate/%s %s", p.version, productVersion)
	if userAgentExtraSuffix != "" {
		userAgentSuffix += fmt.Sprintf(" %s", userAgentExtraSuffix)
	}
	clientConfig.UserAgentSuffix = pointers.String(userAgentSuffix)
	resourceConfig.ApiClient = client.NewAPIClient(clientConfig)
	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig
	tflog.Info(ctx, "Configured PingFederate client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingfederateProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administrativeaccount.AdministrativeAccountDataSource,
		authenticationapiapplication.AuthenticationApiApplicationDataSource,
		authenticationapisettings.AuthenticationApiSettingsDataSource,
		authenticationpoliciesfragments.AuthenticationPoliciesFragmentDataSource,
		authenticationpoliciessettings.AuthenticationPoliciesSettingsDataSource,
		authenticationpolicycontract.AuthenticationPolicyContractDataSource,
		certificate.CertificatesCAExportDataSource,
		certificate.CertificateDataSource,
		clusterstatus.ClusterStatusDataSource,
		configstore.ConfigStoreDataSource,
		datastore.DataStoreDataSource,
		idpadapter.IdpAdapterDataSource,
		idpdefaulturls.IdpDefaultUrlsDataSource,
		idpspconnection.IdpSpConnectionDataSource,
		keypairsigning.KeypairsSigningKeyDataSource,
		keypairssigningcertificate.KeypairsSigningCertificateDataSource,
		keypairssslserver.KeypairsSslServerKeyDataSource,
		keypairssslclient.KeypairsSslClientKeyDataSource,
		keypairssslclientcertificate.KeypairsSslClientCertificateDataSource,
		keypairssslservercertificate.KeypairsSslServerCertificateDataSource,
		license.LicenseDataSource,
		licenseagreement.LicenseAgreementDataSource,
		localidentity.LocalIdentityProfileDataSource,
		oauthaccesstokenmanager.OauthAccessTokenManagerDataSource,
		oauthauthserversettings.OauthServerSettingsDataSource,
		oauthclient.OauthClientDataSource,
		oauthissuer.OauthIssuerDataSource,
		oauthtokenexchangetokengeneratormapping.OauthTokenExchangeTokenGeneratorMappingDataSource,
		oauthopenidconnectpolicy.OpenidConnectPolicyDataSource,
		passwordcredentialvalidator.PasswordCredentialValidatorDataSource,
		protocolmetadatalifetimesettings.ProtocolMetadataLifetimeSettingsDataSource,
		redirectvalidation.RedirectValidationDataSource,
		serversettings.ServerSettingsDataSource,
		serversettingsgeneralsettings.ServerSettingsGeneralDataSource,
		serversettingslogsettings.ServerSettingsLoggingDataSource,
		serversettingssystemkeys.ServerSettingsSystemKeysDataSource,
		sessionapplicationsessionpolicy.SessionApplicationPolicyDataSource,
		sessionauthenticationsessionpoliciesglobal.SessionAuthenticationPoliciesGlobalDataSource,
		sessionsettings.SessionSettingsDataSource,
		spauthenticationpolicycontractmapping.SpAuthenticationPolicyContractMappingDataSource,
		tokenprocessortotokengeneratormapping.TokenProcessorToTokenGeneratorMappingDataSource,
		virtualhostnames.VirtualHostNamesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *pingfederateProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		administrativeaccount.AdministrativeAccountResource,
		authenticationapiapplication.AuthenticationApiApplicationResource,
		authenticationapisettings.AuthenticationApiSettingsResource,
		authenticationpolicies.AuthenticationPoliciesResource,
		authenticationpoliciesfragments.AuthenticationPoliciesFragmentResource,
		authenticationpoliciessettings.AuthenticationPoliciesSettingsResource,
		authenticationpolicycontract.AuthenticationPolicyContractResource,
		authenticationselector.AuthenticationSelectorsResource,
		captchaproviders.CaptchaProviderResource,
		captchaproviderssettings.CaptchaProviderSettingsResource,
		certificate.CertificateCAResource,
		certificatesgroups.CertificatesGroupResource,
		certificatesrevocationocspcertificates.CertificatesRevocationOcspCertificateResource,
		certificatesrevocationsettings.CertificatesRevocationSettingsResource,
		clustersettings.ClusterSettingsResource,
		configstore.ConfigStoreResource,
		configurationencryptionkeysrotate.ConfigurationEncryptionKeysRotateResource,
		connectionmetadataexport.ConnectionMetadataExportResource,
		defaulturls.DefaultUrlsResource,
		extendedproperties.ExtendedPropertiesResource,
		identitystoreprovisioners.IdentityStoreProvisionerResource,
		idpadapter.IdpAdapterResource,
		idpspconnection.IdpSpConnectionResource,
		idpstsrequestparameterscontracts.IdpStsRequestParametersContractResource,
		idptospadaptermapping.IdpToSpAdapterMappingResource,
		idptokenprocessors.IdpTokenProcessorResource,
		incomingproxysettings.IncomingProxySettingsResource,
		kerberosrealms.KerberosRealmsResource,
		kerberosrealmssettings.KerberosRealmSettingsResource,
		keypairsoauthopenidconnect.KeypairsOauthOpenidConnectResource,
		keypairsoauthopenidconnectadditionalkeysets.KeypairsOauthOpenidConnectAdditionalKeySetResource,
		keypairsigning.KeypairsSigningKeyResource,
		keypairssigningrotationsettings.KeypairsSigningKeyRotationSettingsResource,
		keypairssslclient.KeypairsSslClientKeyResource,
		keypairssslserver.KeypairsSslServerKeyResource,
		keypairssslclientcsr.KeypairsSslClientCsrExportResource,
		keypairssslclientcsr.KeypairsSslClientCsrResource,
		keypairssslservercsr.KeypairsSslServerCsrExportResource,
		keypairssslservercsr.KeypairsSslServerCsrResource,
		keypairssslserversettings.KeypairsSslServerSettingsResource,
		datastore.DataStoreResource,
		license.LicenseResource,
		licenseagreement.LicenseAgreementResource,
		localidentity.LocalIdentityProfileResource,
		metadataurls.MetadataUrlResource,
		notificationpublisherssettings.NotificationPublisherSettingsResource,
		notificationpublishers.NotificationPublisherResource,
		oauthaccesstokenmanager.OauthAccessTokenManagerResource,
		oauthaccesstokenmanagerssettings.OauthAccessTokenManagerSettingsResource,
		oauthaccesstokenmapping.OauthAccessTokenMappingResource,
		oauthauthenticationpolicycontractmappings.OauthAuthenticationPolicyContractMappingResource,
		oauthauthserversettings.OauthServerSettingsResource,
		oauthcibaserverpolicyrequestpolicies.OauthCibaServerPolicyRequestPolicyResource,
		oauthcibaserverpolicysettings.OauthCibaServerPolicySettingsResource,
		oauthclient.OauthClientResource,
		oauthclientregistrationpolicies.OauthClientRegistrationPolicyResource,
		oauthclientsettings.OauthClientSettingsResource,
		oauthidpadaptermappings.OauthIdpAdapterMappingResource,
		oauthissuer.OauthIssuerResource,
		oauthopenidconnectpolicy.OpenidConnectPolicyResource,
		oauthopenidconnectsettings.OpenidConnectSettingsResource,
		oauthoutofbandauthplugins.OauthOutOfBandAuthPluginResource,
		oauthresourceownercredentialsmappings.OauthResourceOwnerCredentialsMappingResource,
		oauthtokenexchangegeneratorsettings.OauthTokenExchangeGeneratorSettingsResource,
		oauthtokenexchangetokengeneratormapping.OauthTokenExchangeTokenGeneratorMappingResource,
		passwordcredentialvalidator.PasswordCredentialValidatorResource,
		pingoneconnection.PingoneConnectionResource,
		protocolmetadatalifetimesettings.ProtocolMetadataLifetimeSettingsResource,
		protocolmetadatasigningsettings.ProtocolMetadataSigningSettingsResource,
		redirectvalidation.RedirectValidationResource,
		secretmanagers.SecretManagerResource,
		serversettings.ServerSettingsResource,
		serversettingsgeneralsettings.ServerSettingsGeneralResource,
		serversettingslogsettings.ServerSettingsLoggingResource,
		serversettingssystemkeysrotate.ServerSettingsSystemKeysRotateResource,
		serversettingswstruststssettingsissuercertificates.ServerSettingsWsTrustStsSettingsIssuerCertificateResource,
		serversettingswstruststssettings.ServerSettingsWsTrustStsSettingsResource,
		serviceauthentication.ServiceAuthenticationResource,
		sessionapplicationsessionpolicy.SessionApplicationPolicyResource,
		sessionauthenticationsessionpolicies.SessionAuthenticationPolicyResource,
		sessionauthenticationsessionpoliciesglobal.SessionAuthenticationPoliciesGlobalResource,
		sessionsettings.SessionSettingsResource,
		spadapters.SpAdapterResource,
		spidpconnection.SpIdpConnectionResource,
		spauthenticationpolicycontractmapping.SpAuthenticationPolicyContractMappingResource,
		sptargeturlmappings.SpTargetUrlMappingsResource,
		tokenprocessortotokengeneratormapping.TokenProcessorToTokenGeneratorMappingResource,
		virtualhostnames.VirtualHostNamesResource,
	}
}
