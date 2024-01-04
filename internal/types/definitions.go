package types

import (
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost      string
	Username       string
	Password       string
	ProductVersion version.SupportedVersion
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClient      *client.APIClient
}
