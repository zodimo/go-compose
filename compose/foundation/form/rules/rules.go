// Package rules is a catalog of ready-made validators for the form engine's
// controls.
//
// Every rule is a plain fform.ValidatorFunc[T], so it composes with NewControl,
// WithCode, and the engine exactly like a hand-written validator. Rules return a
// *fform.ValidationError carrying a stable fform.ErrorCode, so callers can branch
// on the failure kind without parsing messages.
//
// Provenance: the rule set is inspired by github.com/nobl9/govy (MPL-2.0) and
// other Go validation libraries, but the implementations here are original to
// this repository. No govy source is copied, so no MPL obligations attach to
// these files.
package rules

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// Code is the ErrorCode type for rules, aliased so callers can write
// rules.CodeRequired etc. without importing the form package for the type.
type Code = fform.ErrorCode

// Error codes raised by this package. Each rule documents the code it raises.
const (
	CodeRequired      Code = "required"
	CodeOmitEmpty     Code = "omit_empty"
	CodeMinLength     Code = "min_length"
	CodeMaxLength     Code = "max_length"
	CodeLength        Code = "length"
	CodeMatchRegexp   Code = "match_regexp"
	CodeDenyRegexp    Code = "deny_regexp"
	CodeContains      Code = "contains"
	CodeExcludes      Code = "excludes"
	CodeStartsWith    Code = "starts_with"
	CodeEndsWith      Code = "ends_with"
	CodeAlpha         Code = "alpha"
	CodeAlphanumeric  Code = "alphanumeric"
	CodeNumeric       Code = "numeric"
	CodeEmail         Code = "email"
	CodeURL           Code = "url"
	CodeUUID          Code = "uuid"
	CodeIP            Code = "ip"
	CodeIPv4          Code = "ipv4"
	CodeIPv6          Code = "ipv6"
	CodeDNSLabel      Code = "dns_label"
	CodeDNSSubdomain  Code = "dns_subdomain"
	CodeBase64        Code = "base64"
	CodeHex           Code = "hex"
	CodeSemver        Code = "semver"
	CodeOneOf         Code = "one_of"
	CodeNoneOf        Code = "none_of"
	CodeEQ            Code = "eq"
	CodeNEQ           Code = "neq"
	CodeGT            Code = "gt"
	CodeGTE           Code = "gte"
	CodeLT            Code = "lt"
	CodeLTE           Code = "lte"
	CodeUnique        Code = "unique"
	CodeSliceContains Code = "slice_contains"
	CodeMapKeyPresent Code = "map_key_present"
	CodeCrossField    Code = "cross_field"
	CodeCompareFields Code = "compare_fields"
)

// fail builds a coded ValidationError.
func fail(code Code, format string, args ...any) error {
	return fform.NewValidationError(code, format, args...)
}
