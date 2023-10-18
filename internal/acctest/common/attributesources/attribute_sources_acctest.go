package attributesources

import (
	"fmt"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
)

func JdbcHcl(attrSource *client.JdbcAttributeSource) string {
	var builder strings.Builder
	if attrSource == nil {
		return ""
	}
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

// jdbcAttributeSource := client.NewJdbcAttributeSource(
//
//	"CHANNEL_GROUP", "$${SAML_SUBJECT}", "JDBC", *client.NewResourceLink("ProvisionerDS"),
//
// )
func JdbcClientStruct(table string, filter string, attributeSourceType string, resourceLink client.ResourceLink) *client.JdbcAttributeSource {
	jdbcAttributeSource := client.NewJdbcAttributeSource(
		table, filter, attributeSourceType, resourceLink,
	)
	jdbcAttributeSource.Id = pointers.String("attributesourceid")
	jdbcAttributeSource.ColumnNames = []string{"CREATED"}
	jdbcAttributeSource.Description = pointers.String("description")
	jdbcAttributeSource.Schema = pointers.String("PUBLIC")
	return jdbcAttributeSource
}
