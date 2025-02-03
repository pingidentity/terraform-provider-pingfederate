// Copyright © 2025 Ping Identity Corporation

package attributesources

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func ClientStruct(attributeSourcesAttr basetypes.SetValue) ([]client.AttributeSourceAggregation, error) {
	attributeSourceAggregation := []client.AttributeSourceAggregation{}
	for _, source := range attributeSourcesAttr.Elements() {
		//Determine which attribute source type this is
		sourceAttrs := source.(types.Object).Attributes()
		attributeSourceInner := client.AttributeSourceAggregation{}
		if internaltypes.IsDefined(sourceAttrs["custom_attribute_source"]) {
			attributeSourceInner.CustomAttributeSource = &client.CustomAttributeSource{}
			customAttributeSourceErr := json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["custom_attribute_source"], true)), attributeSourceInner.CustomAttributeSource)
			if customAttributeSourceErr != nil {
				return nil, customAttributeSourceErr
			}
		}
		if internaltypes.IsDefined(sourceAttrs["jdbc_attribute_source"]) {
			attributeSourceInner.JdbcAttributeSource = &client.JdbcAttributeSource{}
			jdbcAttributeSourceErr := json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["jdbc_attribute_source"], true)), attributeSourceInner.JdbcAttributeSource)
			if jdbcAttributeSourceErr != nil {
				return nil, jdbcAttributeSourceErr
			}
		}
		if internaltypes.IsDefined(sourceAttrs["ldap_attribute_source"]) {
			attributeSourceInner.LdapAttributeSource = &client.LdapAttributeSource{}
			ldapAttributeSourceErr := json.Unmarshal([]byte(internaljson.FromValue(sourceAttrs["ldap_attribute_source"], true)), attributeSourceInner.LdapAttributeSource)
			if ldapAttributeSourceErr != nil {
				return nil, ldapAttributeSourceErr
			}
		}
		attributeSourceAggregation = append(attributeSourceAggregation, attributeSourceInner)
	}
	return attributeSourceAggregation, nil
}
