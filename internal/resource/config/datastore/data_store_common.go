// Copyright © 2026 Ping Identity Corporation

package datastore

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataStoreModel struct {
	Id                          types.String `tfsdk:"id"`
	DataStoreId                 types.String `tfsdk:"data_store_id"`
	MaskAttributeValues         types.Bool   `tfsdk:"mask_attribute_values"`
	CustomDataStore             types.Object `tfsdk:"custom_data_store"`
	JdbcDataStore               types.Object `tfsdk:"jdbc_data_store"`
	LdapDataStore               types.Object `tfsdk:"ldap_data_store"`
	PingOneLdapGatewayDataStore types.Object `tfsdk:"ping_one_ldap_gateway_data_store"`
}
