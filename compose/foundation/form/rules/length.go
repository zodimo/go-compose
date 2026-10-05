package rules

import (
	"strings"
	"unicode"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// Required fails when value equals zero. It raises CodeRequired.
//
// For string values, prefer StringNotEmpty, which also rejects whitespace-only
// input; Required("") accepts " " as non-empty.
func Required[T comparable](zero T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value == zero {
			return fail(CodeRequired, "value is required")
		}
		return nil
	}
}

// StringNotEmpty fails when the trimmed value is empty. It raises CodeRequired,
// so callers can treat "missing" and "blank" the same way.
func StringNotEmpty() fform.ValidatorFunc[string] {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fail(CodeRequired, "must not be empty")
		}
		return nil
	}
}

// MinLength fails when the rune count is below min. It raises CodeMinLength.
// Counting runes (not bytes) keeps multibyte text correct.
func MinLength(min int) fform.ValidatorFunc[string] {
	return func(value string) error {
		if n := len([]rune(value)); n < min {
			return fail(CodeMinLength, "must be at least %d characters", min)
		}
		return nil
	}
}

// MaxLength fails when the rune count exceeds max. It raises CodeMaxLength.
func MaxLength(max int) fform.ValidatorFunc[string] {
	return func(value string) error {
		if n := len([]rune(value)); n > max {
			return fail(CodeMaxLength, "must be at most %d characters", max)
		}
		return nil
	}
}

// Length fails when the rune count is outside [min, max]. It raises CodeLength.
func Length(min, max int) fform.ValidatorFunc[string] {
	return func(value string) error {
		n := len([]rune(value))
		if n < min || n > max {
			return fail(CodeLength, "must be between %d and %d characters", min, max)
		}
		return nil
	}
}

// MinItems fails when the slice length is below min. It raises CodeMinLength.
func MinItems[T any](min int) fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		if len(value) < min {
			return fail(CodeMinLength, "must have at least %d items", min)
		}
		return nil
	}
}

// MaxItems fails when the slice length exceeds max. It raises CodeMaxLength.
func MaxItems[T any](max int) fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		if len(value) > max {
			return fail(CodeMaxLength, "must have at most %d items", max)
		}
		return nil
	}
}

// Alpha fails when the value contains anything other than ASCII letters.
// It raises CodeAlpha. Empty input passes; pair with StringNotEmpty.
func Alpha() fform.ValidatorFunc[string] {
	return runeClass(CodeAlpha, "must contain only letters", unicode.IsLetter, true)
}

// AlphaUnicode is Alpha for non-ASCII letters. It raises CodeAlpha.
func AlphaUnicode() fform.ValidatorFunc[string] {
	return runeClass(CodeAlpha, "must contain only letters", unicode.IsLetter, false)
}

// Numeric fails when the value contains anything other than ASCII digits.
// It raises CodeNumeric.
func Numeric() fform.ValidatorFunc[string] {
	return runeClass(CodeNumeric, "must contain only digits", unicode.IsDigit, true)
}

// NumericUnicode is Numeric for non-ASCII digits. It raises CodeNumeric.
func NumericUnicode() fform.ValidatorFunc[string] {
	return runeClass(CodeNumeric, "must contain only digits", unicode.IsDigit, false)
}

// Alphanumeric fails when the value contains anything other than ASCII letters
// and digits. It raises CodeAlphanumeric.
func Alphanumeric() fform.ValidatorFunc[string] {
	return runeClass(CodeAlphanumeric, "must contain only letters and digits", func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}, true)
}

// AlphanumericUnicode is Alphanumeric for non-ASCII letters and digits.
func AlphanumericUnicode() fform.ValidatorFunc[string] {
	return runeClass(CodeAlphanumeric, "must contain only letters and digits", func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}, false)
}

// runeClass builds a validator that fails when any rune is outside pred. When
// asciiOnly is true, any rune above 0x7F also fails.
func runeClass(code Code, message string, pred func(rune) bool, asciiOnly bool) fform.ValidatorFunc[string] {
	return func(value string) error {
		for _, r := range value {
			if asciiOnly && r > unicode.MaxASCII {
				return fail(code, "%s", message)
			}
			if !pred(r) {
				return fail(code, "%s", message)
			}
		}
		return nil
	}
}
