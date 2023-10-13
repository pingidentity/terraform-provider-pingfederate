package attributesources

import (
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
		builder.WriteString("  attribute_sources = [\n")
		builder.WriteString("    {\n")
		builder.WriteString("      jdbc_attribute_source = {\n")
		builder.WriteString("        data_store_ref = {\n")
		builder.WriteString("          id = \"")
		builder.WriteString(attrSource.DataStoreRef.Id)
		builder.WriteString("\"\n        }\n        id           = \"")
		builder.WriteString(*attrSource.Id)
		builder.WriteString("\"\n        description  = \"")
		builder.WriteString(*attrSource.Description)
		builder.WriteString("\"\n        schema       = \"")
		builder.WriteString(*attrSource.Schema)
		builder.WriteString("\"\n        table        = \"")
		builder.WriteString(attrSource.Table)
		builder.WriteString("\"\n        filter       = \"")
		builder.WriteString(attrSource.Filter)
		builder.WriteString("\"\n        column_names = ")
		builder.WriteString(acctest.StringSliceToTerraformString(attrSource.ColumnNames))
		builder.WriteString("\n      }\n    }\n  ]")
	}
	return builder.String()
}

func JdbcClientStruct() *client.JdbcAttributeSource {
	jdbcAttributeSource := client.NewJdbcAttributeSource(
		"CHANNEL_GROUP", "$${SAML_SUBJECT}", "JDBC", *client.NewResourceLink("ProvisionerDS"),
	)
	jdbcAttributeSource.Id = pointers.String("attributesourceid")
	jdbcAttributeSource.ColumnNames = []string{"CREATED"}
	jdbcAttributeSource.Description = pointers.String("description")
	jdbcAttributeSource.Schema = pointers.String("PUBLIC")
	return jdbcAttributeSource
}
