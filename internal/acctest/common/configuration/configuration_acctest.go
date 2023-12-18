package configuration

import (
	"strconv"
	"strings"

	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
)

func Hcl(configuration client.PluginConfiguration) string {
	var builder strings.Builder
	builder.WriteString("configuration = {\n")
	builder.WriteString("    tables = [\n")
	for _, table := range configuration.Tables {
		builder.WriteString("       {\n")
		builder.WriteString("           name = \"")
		builder.WriteString(table.Name)
		builder.WriteString("\"\n")
		builder.WriteString("           rows = [\n")
		for _, row := range table.Rows {
			builder.WriteString("               {\n")
			if row.DefaultRow != nil {
				builder.WriteString("                   default_row = ")
				builder.WriteString(strconv.FormatBool(*row.DefaultRow))
				builder.WriteRune('\n')
			}
			builder.WriteString("                   fields = [\n")
			for _, field := range row.Fields {
				builder.WriteString("                       {\n")
				builder.WriteString("                           name = \"")
				builder.WriteString(field.Name)
				builder.WriteString("\"\n")
				if field.Value != nil {
					builder.WriteString("                           value = \"")
					builder.WriteString(*field.Value)
					builder.WriteString("\"\n")
				}
				builder.WriteString("                       },\n")
			}
			builder.WriteString("                   ]\n")
			builder.WriteString("               }\n")
		}
		builder.WriteString("           ]\n")
		builder.WriteString("       },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("    fields = [\n")
	for _, field := range configuration.Fields {
		builder.WriteString("        {\n")
		builder.WriteString("            name = \"")
		builder.WriteString(field.Name)
		builder.WriteString("\"\n")
		if field.Value != nil {
			builder.WriteString("            value = \"")
			builder.WriteString(*field.Value)
			builder.WriteString("\"\n")
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("}\n")
	return builder.String()
}
