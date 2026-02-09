// Copyright © 2026 Ping Identity Corporation

package attributemapping

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
)

func DataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
			"attribute_sources":              attributesources.ToDataSourceSchema(),
			"issuance_criteria":              issuancecriteria.ToDataSourceSchema(),
		},
		Computed:    true,
		Optional:    false,
		Description: "A list of mappings from attribute sources to attribute targets.",
	}
}
