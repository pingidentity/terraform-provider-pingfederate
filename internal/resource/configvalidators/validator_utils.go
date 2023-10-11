package configvalidators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func IsUrlFormat(aV attr.Value, validateConfigResp *resource.ValidateConfigResponse) diag.Diagnostics {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	//^ and $: Start and end of the string.
	// (https?:\/\/): Matches the http:// or https:// scheme. This part is optional
	// ([\da-z\.-]+): Matches the domain name, which can contain digits, alphabets (lowercase), dots, and hyphens.
	//       One or more such characters must be present.
	// (\.[a-z]{2,6}(\.[a-z]{2,6})?)?: Matches the optional TLD and sub-TLD, each containing 2 to 6 alphabets (lowercase).
	//       The entire TLD part is optional.
	// (:[0-9]{2,5})?: Matches the optional port number, which can be 2 to 5 digits long.
	// (\/[^\s\/]+[^\s]*)? matches the optional path
	//
	// The regex will match the following examples:
	// 	http://example.com
	//	https://example.com
	//	http://example.com:8080
	//	https://example.co.uk:8080
	//	https://localhost
	//	http://localhost:3000
	//	http://example.com/path/to/resource
	//  https://anotherexample.co.uk:9399/path/with/trailing/slash/
	//  localhost
	//  localhost:9999
	//
	//  Non-Matching examples
	//	http://example. (TLD ends with a dot)
	//	http://example.com:808080 (Port number exceeds 5 digits)
	//	http://example_com.com (Underscore is not allowed in domain name)
	//	ftp://example.com (Only http and https schemes are allowed if specified)
	basetypesStringVal := aV.(types.String)
	re, _ := regexp.Compile(`^(https?:\/\/)?([\da-z\.-]+)(\.[a-z]{2,6}(\.[a-z]{2,6})?)?(:[0-9]{2,5})?(\/)?(\/[^\s]*)?$`)
	isUrl := re.MatchString(strings.ToLower(aV.(types.String).ValueString()))
	if !internaltypes.IsNonEmptyString(basetypesStringVal) && !isUrl {
		diag := diag.Diagnostics{}
		diag.AddError("Invalid URL Format!", fmt.Sprintf("Please provide a valid origin. Origin \"%s\" needs to be in a valid URL-like format - \"http(s)//:<value>.<domain>\"", basetypesStringVal.ValueString()))
		return diag
	}
	return diag.Diagnostics{}
}

func IsEmailFormat(aV attr.Value, validateConfigResp *resource.ValidateConfigResponse) diag.Diagnostics {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	// ^ and $: Start and end of the string.
	// [a-z0-9._%+-]+: Matches the local part of the email address, which can contain alphabetic characters, digits, dots, underscores, percent signs, pluses, and hyphens. One or more of these options must be present.
	// @: Matches the "@" symbol.
	// [a-z0-9.-]+: Matches the domain name, which can contain alphabetic characters, digits, dots, and hyphens. One or more of these options must be present.
	// \.: Matches the dot before the top-level domain (TLD).
	// [a-z]{2,}: Matches the TLD, which must contain at least two alphabetic characters.
	//
	// The regex will match the following example emails:
	// john.doe@example.com
	// Jane.Doe@sub.example.co
	// email+filter@gmail.com
	// email_filter@gmail.com
	// juan%sanchez@gmail.com
	// 99user@domain.com
	// user.name@domain.co.uk
	//
	// Non-Matching emails
	// @example.com (Local part is missing)
	// john.doe@.com (Domain name is missing)
	// john.doe@com (Dot before TLD is missing)
	// john.doe@domain. (TLD is missing)
	// john.doe@domain.c (TLD is too short)
	// john$doe@domain.com (invalid '$' character)
	// john..doe@example.com (Consecutive dots in the local part)
	// john.doe@.example.com (Domain starts with a dot)
	basetypesStringVal := aV.(types.String)
	re, _ := regexp.Compile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	isEmailString := re.MatchString(strings.ToLower(basetypesStringVal.ValueString()))
	if !internaltypes.IsNonEmptyString(basetypesStringVal) && !isEmailString {
		diag := diag.Diagnostics{}
		diag.AddError("Invalid Email Format!", fmt.Sprintf("Please provide a valid email address - \"%s\" needs to be in a valid email format according to RFC 5322.  For example, \"<user>@<company>.<tld>\"", basetypesStringVal.ValueString()))
		return diag
	}
	return diag.Diagnostics{}
}

func IsValidHostnameOrIp(aV attr.Value, validateConfigResp *resource.ValidateConfigResponse) diag.Diagnostics {
	// ToLower on the MatchString call makes this case-insensitive
	// Regex Breakdown
	// This implementation is not complete nor perfect, but will catch many invalid hostnames or IPs,
	// including invalid characters, whitespace, etc.
	// TODO:  expand for thoroughness and possible IPV6 validation
	// This portion matches hostnames:
	// ^([a-z0-9]+(-[a-z0-9]+)*\.)*[a-z0-9]+(-[a-z0-9]+)*$
	// ^ and $: Start and end of the string.
	// It allows for alphabetic characters, digits, and hyphens.
	// A hostname part cannot start or end with a hyphen.
	// Cannot start or end with a dot.
	// Hostname parts are separated by dots.
	basetypesStringVal := aV.(types.String)
	re, _ := regexp.Compile(`^([a-z0-9]+(-[a-z0-9]+)*\.)*[a-z0-9]+(-[a-z0-9]+)*$`)
	isValidHostnameOrIp := re.MatchString(strings.ToLower(basetypesStringVal.ValueString()))

	if !internaltypes.IsNonEmptyString(basetypesStringVal) && !isValidHostnameOrIp {
		diag := diag.Diagnostics{}
		diag.AddError("Invalid hostname or IP!", fmt.Sprintf("Please provide a valid hostname or IP address - \"%s\" is invalid", basetypesStringVal.ValueString()))
		return diag
	}
	return diag.Diagnostics{}
}
