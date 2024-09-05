resource "pingfederate_password_credential_validator" "ldapUsernamePasswordCredentialValidatorExample" {
  validator_id = "ldapUnPwPCV"
  name         = "ldapUsernamePasswordCredentialValidatorExample"
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
        value = "mydatastore"
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
