terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
}

resource "pingfederate_protocolmetadata_lifetimesettings" "protocolMetadataLifetimeSettingsExample" {
	cache_duration = 1440
  reload_delay = 1440
}