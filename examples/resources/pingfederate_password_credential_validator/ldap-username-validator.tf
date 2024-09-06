resource "pingfederate_data_store" "pingDirectoryLdapDataStore" {
  ldap_data_store = {
    name      = "PingDirectory LDAP Data Store"
    ldap_type = "PING_DIRECTORY"

    user_dn  = var.pingdirectory_bind_dn
    password = var.pingdirectory_bind_dn_password

    use_ssl = true

    hostnames = [
      "pingdirectory:636"
    ]
  }
}

resource "pingfederate_password_credential_validator" "ldapUsernamePasswordCredentialValidatorExample" {
  validator_id = "ldapUnPwPCV"
  name         = "LDAP Username Password Credential Validator"

  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.LDAPUsernamePasswordCredentialValidator"
  }

  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name  = "LDAP Datastore"
        value = pingfederate_data_store.pingDirectoryLdapDataStore.id
      },
      {
        name  = "Search Base"
        value = "cn=Users"
      },
      {
        name  = "Search Filter"
        value = "sAMAccountName=$${username}"
      },
      {
        name  = "Scope of Search"
        value = "Subtree"
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      },
      {
        name  = "Display Name Attribute"
        value = "displayName"
      },
      {
        name  = "Mail Attribute"
        value = "mail"
      },
      {
        name  = "Trim Username Spaces For Search"
        value = "true"
      },
      {
        name  = "Enable PingDirectory Detailed Password Policy Requirement Messaging"
        value = "true"
      },
      {
        name  = "Expect Password Expired Control"
        value = "false"
      }
    ]
  }
}
