package types

import client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost string
	Username  string
	Password  string
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClient      *client.APIClient
}
