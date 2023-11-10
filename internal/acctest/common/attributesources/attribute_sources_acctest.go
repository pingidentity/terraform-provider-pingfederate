package attributesources

import (
	"fmt"
	"strconv"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
)

func JdbcHcl(attrSource *client.JdbcAttributeSource) string {
	if attrSource == nil {
		return ""
	}
	var builder strings.Builder
	if attrSource != nil {
		tf := `
		attribute_sources = [
			{
				jdbc_attribute_source = {
					data_store_ref = {
						id = "%s"
					}
					id           = "%s"
					description  = "%s"
					schema       = "%s"
					table        = "%s"
					filter       = "%s"
					column_names = %s
				}
			}
		]
	`
		builder.WriteString(fmt.Sprintf(tf,
			attrSource.DataStoreRef.Id,
			*attrSource.Id,
			*attrSource.Description,
			*attrSource.Schema,
			attrSource.Table,
			attrSource.Filter,
			acctest.StringSliceToTerraformString(attrSource.ColumnNames)))
	}
	return builder.String()
}

func JdbcClientStruct(table string, filter string, attributeSourceType string, resourceLink client.ResourceLink) *client.JdbcAttributeSource {
	jdbcAttributeSource := client.NewJdbcAttributeSource(
		table, filter, attributeSourceType, resourceLink,
	)
	jdbcAttributeSource.Id = pointers.String("jdbcattributesourceid")
	jdbcAttributeSource.ColumnNames = []string{"CREATED"}
	jdbcAttributeSource.Description = pointers.String("description")
	jdbcAttributeSource.Schema = pointers.String("PUBLIC")
	return jdbcAttributeSource
}

func LdapHcl(attrSource *client.LdapAttributeSource) string {
	if attrSource == nil {
		return ""
	}
	var builder strings.Builder
	if attrSource != nil {
		binaryAttributeSettingsHcl := ""
		for key, setting := range *attrSource.BinaryAttributeSettings {
			binaryAttributeSettingsHcl += fmt.Sprintf(`
			    "%[1]s" = {
					binary_encoding = "%[2]s"
				},
			`, key, *setting.BinaryEncoding)
		}

		tf := `
		attribute_sources = [
			{
				ldap_attribute_source = {
					member_of_nested_group = %s
					type                   = "LDAP"
					search_scope           = "%s"
					search_filter          = "%s"
					data_store_ref = {
						id = "%s"
					}
					id                = "%s"
					description       = "%s"
					base_dn           = "%s"
					search_attributes = %s
					binary_attribute_settings = {
						%s
					}
				}
			}
		]
	`
		builder.WriteString(fmt.Sprintf(tf,
			strconv.FormatBool(*attrSource.MemberOfNestedGroup),
			attrSource.SearchScope,
			attrSource.SearchFilter,
			attrSource.DataStoreRef.Id,
			*attrSource.Id,
			*attrSource.Description,
			*attrSource.BaseDn,
			acctest.StringSliceToTerraformString(attrSource.SearchAttributes),
			binaryAttributeSettingsHcl))
	}
	return builder.String()
}

func LdapClientStruct(searchFilter, searchScope string, dataStoreRef client.ResourceLink) *client.LdapAttributeSource {
	ldapAttributeSource := client.NewLdapAttributeSource(searchScope, searchFilter, "LDAP", dataStoreRef)
	ldapAttributeSource.MemberOfNestedGroup = pointers.Bool(false)
	ldapAttributeSource.Id = pointers.String("ldapattributesourceid")
	ldapAttributeSource.BaseDn = pointers.String("dc=example,dc=com")
	ldapAttributeSource.Description = pointers.String("description")
	ldapAttributeSource.SearchAttributes = []string{
		"cn",
		"sn",
		"mail",
	}
	ldapAttributeSource.BinaryAttributeSettings = &map[string]client.BinaryLdapAttributeSettings{}
	return ldapAttributeSource
}
