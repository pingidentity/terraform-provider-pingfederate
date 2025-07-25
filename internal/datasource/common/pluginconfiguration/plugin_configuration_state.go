// Copyright © 2025 Ping Identity Corporation

package pluginconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

func ToDataSourceState(con context.Context, configuration *client.PluginConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	fieldsAttrValue, respDiags := types.ListValueFrom(con, types.ObjectType{AttrTypes: fieldAttrTypes}, configuration.Fields)
	diags.Append(respDiags...)
	tablesAttrValue, respDiags := types.ListValueFrom(con, types.ObjectType{AttrTypes: tableAttrTypes}, configuration.Tables)
	diags.Append(respDiags...)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}

	configObj, valueFromDiags := types.ObjectValue(configurationAttrTypes, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
}
