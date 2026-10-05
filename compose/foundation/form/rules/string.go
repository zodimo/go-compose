package rules

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// MatchRegexp fails when the value does not match re. It raises
// CodeMatchRegexp.
func MatchRegexp(re *regexp.Regexp) fform.ValidatorFunc[string] {
	return func(value string) error {
		if re == nil {
			return nil
		}
		if !re.MatchString(value) {
			return fail(CodeMatchRegexp, "must match pattern %q", re.String())
		}
		return nil
	}
}

// DenyRegexp fails when the value matches re. It raises CodeDenyRegexp.
func DenyRegexp(re *regexp.Regexp) fform.ValidatorFunc[string] {
	return func(value string) error {
		if re == nil {
			return nil
		}
		if re.MatchString(value) {
			return fail(CodeDenyRegexp, "must not match pattern %q", re.String())
		}
		return nil
	}
}

// Contains fails when the value does not contain every substring. It raises
// CodeContains.
func Contains(substrings ...string) fform.ValidatorFunc[string] {
	return func(value string) error {
		for _, s := range substrings {
			if !strings.Contains(value, s) {
				return fail(CodeContains, "must contain %q", s)
			}
		}
		return nil
	}
}

// Excludes fails when the value contains any substring. It raises CodeExcludes.
func Excludes(substrings ...string) fform.ValidatorFunc[string] {
	return func(value string) error {
		for _, s := range substrings {
			if strings.Contains(value, s) {
				return fail(CodeExcludes, "must not contain %q", s)
			}
		}
		return nil
	}
}

// StartsWith fails when the value does not begin with any prefix. It raises
// CodeStartsWith.
func StartsWith(prefixes ...string) fform.ValidatorFunc[string] {
	return func(value string) error {
		for _, p := range prefixes {
			if strings.HasPrefix(value, p) {
				return nil
			}
		}
		return fail(CodeStartsWith, "must start with one of %v", prefixes)
	}
}

// EndsWith fails when the value does not end with any suffix. It raises
// CodeEndsWith.
func EndsWith(suffixes ...string) fform.ValidatorFunc[string] {
	return func(value string) error {
		for _, s := range suffixes {
			if strings.HasSuffix(value, s) {
				return nil
			}
		}
		return fail(CodeEndsWith, "must end with one of %v", suffixes)
	}
}

// Email fails when the value is not a parseable email address (RFC 5322,
// via net/mail). It raises CodeEmail.
func Email() fform.ValidatorFunc[string] {
	return func(value string) error {
		addr, err := mail.ParseAddress(value)
		if err != nil || addr.Address != value {
			return fail(CodeEmail, "must be a valid email address")
		}
		return nil
	}
}

// URL fails when the value is not an absolute URL with a scheme and host. It
// raises CodeURL.
func URL() fform.ValidatorFunc[string] {
	return func(value string) error {
		u, err := url.Parse(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fail(CodeURL, "must be a valid URL")
		}
		return nil
	}
}

// UUID fails when the value is not a canonical 8-4-4-4-12 hex UUID. It raises
// CodeUUID.
func UUID() fform.ValidatorFunc[string] {
	return func(value string) error {
		if !uuidPattern.MatchString(value) {
			return fail(CodeUUID, "must be a valid UUID")
		}
		return nil
	}
}

// IP fails when the value is not a valid IP address. It raises CodeIP.
func IP() fform.ValidatorFunc[string] {
	return func(value string) error {
		if net.ParseIP(value) == nil {
			return fail(CodeIP, "must be a valid IP address")
		}
		return nil
	}
}

// IPv4 fails when the value is not a valid dotted-quad IPv4 address. It raises
// CodeIPv4.
func IPv4() fform.ValidatorFunc[string] {
	return func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil || ip.To4() == nil || strings.Contains(value, ":") {
			return fail(CodeIPv4, "must be a valid IPv4 address")
		}
		return nil
	}
}

// IPv6 fails when the value is not a valid IPv6 address. It raises CodeIPv6.
func IPv6() fform.ValidatorFunc[string] {
	return func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil || ip.To4() != nil {
			return fail(CodeIPv6, "must be a valid IPv6 address")
		}
		return nil
	}
}

// DNSLabel fails when the value is not an RFC 1123 DNS label (1-63 chars, lower
// case alphanumeric or '-', starting and ending alphanumeric). It raises
// CodeDNSLabel.
func DNSLabel() fform.ValidatorFunc[string] {
	return func(value string) error {
		if !dnsLabelPattern.MatchString(value) {
			return fail(CodeDNSLabel, "must be a valid DNS label")
		}
		return nil
	}
}

// DNSSubdomain fails when the value is not an RFC 1123 DNS subdomain. It raises
// CodeDNSSubdomain.
func DNSSubdomain() fform.ValidatorFunc[string] {
	return func(value string) error {
		if len(value) == 0 || len(value) > 253 || !dnsSubdomainPattern.MatchString(value) {
			return fail(CodeDNSSubdomain, "must be a valid DNS subdomain")
		}
		return nil
	}
}

// Base64 fails when the value is not valid standard base64. It raises
// CodeBase64.
func Base64() fform.ValidatorFunc[string] {
	return func(value string) error {
		if _, err := base64.StdEncoding.DecodeString(value); err != nil {
			return fail(CodeBase64, "must be valid base64")
		}
		return nil
	}
}

// Hexadecimal fails when the value is not an even-length hex string. It raises
// CodeHex.
func Hexadecimal() fform.ValidatorFunc[string] {
	return func(value string) error {
		if len(value)%2 != 0 {
			return fail(CodeHex, "must be a valid hexadecimal string")
		}
		if _, err := hex.DecodeString(value); err != nil {
			return fail(CodeHex, "must be a valid hexadecimal string")
		}
		return nil
	}
}

// Semver fails when the value is not a semantic version (major.minor.patch with
// optional pre-release and build metadata). It raises CodeSemver.
func Semver() fform.ValidatorFunc[string] {
	return func(value string) error {
		if !semverPattern.MatchString(value) {
			return fail(CodeSemver, "must be a valid semantic version")
		}
		return nil
	}
}
