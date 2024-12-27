package attributesources

import (
	"fmt"
	"strconv"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
)

func Hcl(jdbcAttrSource *client.JdbcAttributeSource, ldapAttrSource *client.LdapAttributeSource) string {
	if jdbcAttrSource == nil && ldapAttrSource == nil {
		return ""
	}

	return fmt.Sprintf(`
		attribute_sources = [
			%s
			%s
		]
	`, JdbcHcl(jdbcAttrSource), LdapHcl(ldapAttrSource))
}

func JdbcHcl(attrSource *client.JdbcAttributeSource) string {
	if attrSource == nil {
		return ""
	}
	var builder strings.Builder
	if attrSource != nil {
		tf := `
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
			},
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
	jdbcAttributeSource.Description = pointers.String("jdbcdescription")
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
					binary_attribute_settings = {
			    "%[1]s" = {
					binary_encoding = "%[2]s"
				},
		}
			`, key, *setting.BinaryEncoding)
		}

		tf := `
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
					%s
				}
			},
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
	ldapAttributeSource.Description = pointers.String("ldapdescription")
	ldapAttributeSource.SearchAttributes = []string{
		"cn",
		"sn",
		"mail",
	}
	ldapAttributeSource.BinaryAttributeSettings = &map[string]client.BinaryLdapAttributeSettings{}
	return ldapAttributeSource
}

func ValidateResponseAttributes(resourceType string, resourceName *string, expectedJdbcAttrSource *client.JdbcAttributeSource,
	expectedLdapAttrSource *client.LdapAttributeSource, actualAttrSources []client.AttributeSourceAggregation) error {
	if expectedJdbcAttrSource == nil && expectedLdapAttrSource == nil {
		return nil
	}

	var err error
	for _, attributeSource := range actualAttrSources {
		if attributeSource.JdbcAttributeSource != nil && expectedJdbcAttrSource != nil {
			err = acctest.TestAttributesMatchString(resourceType, resourceName, "id",
				expectedJdbcAttrSource.DataStoreRef.Id, attributeSource.JdbcAttributeSource.DataStoreRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "description",
				*expectedJdbcAttrSource.Description, *attributeSource.JdbcAttributeSource.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "schema",
				*expectedJdbcAttrSource.Description, *attributeSource.JdbcAttributeSource.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "table",
				expectedJdbcAttrSource.Table, attributeSource.JdbcAttributeSource.Table)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringSlice(resourceType, resourceName, "column_names",
				expectedJdbcAttrSource.ColumnNames, attributeSource.JdbcAttributeSource.ColumnNames)
			if err != nil {
				return err
			}
		}
		if attributeSource.LdapAttributeSource != nil && expectedLdapAttrSource != nil {
			err = acctest.TestAttributesMatchString(resourceType, resourceName, "id",
				expectedLdapAttrSource.DataStoreRef.Id, attributeSource.LdapAttributeSource.DataStoreRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "description",
				*expectedLdapAttrSource.Description, *attributeSource.LdapAttributeSource.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "schema",
				*expectedLdapAttrSource.Description, *attributeSource.LdapAttributeSource.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, resourceName, "baseDn",
				*expectedLdapAttrSource.BaseDn, attributeSource.LdapAttributeSource.BaseDn)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, resourceName, "searchScope",
				expectedLdapAttrSource.SearchScope, attributeSource.LdapAttributeSource.SearchScope)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringSlice(resourceType, resourceName, "searchAttributes",
				expectedLdapAttrSource.SearchAttributes, attributeSource.LdapAttributeSource.SearchAttributes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
