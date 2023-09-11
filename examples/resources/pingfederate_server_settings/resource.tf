resource "pingfederate_server_settings" "serverSettingsExample" {
  contact_info = {
    company    = "example company"
    email      = "adminemail@company.com"
    first_name = "Jane"
    last_name  = "Admin"
    phone      = "555-555-1222"
  }

  federation_info = {
    base_url = "https://localhost:9999"
  }

  email_server = {
    source_addr  = "emailServerAdmin@company.com"
    email_server = "myemailserver.company.com"
  }
}
