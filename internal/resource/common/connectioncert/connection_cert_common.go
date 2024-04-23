package connectioncert

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ObjType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: AttrTypes(),
	}
}

func AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cert": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                        types.StringType,
				"serial_number":             types.StringType,
				"subject_dn":                types.StringType,
				"subject_alternative_names": types.SetType{ElemType: types.StringType},
				"issuer_dn":                 types.StringType,
				"valid_from":                types.StringType,
				"expires":                   types.StringType,
				"key_algorithm":             types.StringType,
				"key_size":                  types.Int64Type,
				"signature_algorithm":       types.StringType,
				"version":                   types.Int64Type,
				"sha1_fingerprint":          types.StringType,
				"sha256_fingerprint":        types.StringType,
				"status":                    types.StringType,
			},
		},
		"x509_file": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":              types.StringType,
				"file_data":       types.StringType,
				"crypto_provider": types.StringType,
			},
		},
		"active_verification_cert":    types.BoolType,
		"primary_verification_cert":   types.BoolType,
		"secondary_verification_cert": types.BoolType,
		"encryption_cert":             types.BoolType,
	}
}
