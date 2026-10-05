package fform

import (
	"errors"
	"fmt"
	"sort"
)

// ErrorCode is a stable, machine-readable identifier for a validation failure.
// Codes let callers branch on *why* a field failed (localization, mapping to API
// error responses, focusing the first failing field) without parsing messages.
//
// Codes are plain strings so callers can define their own alongside the
// built-ins shipped in the rules subpackage; there is no closed enum.
type ErrorCode string

// Built-in codes for the validators the form package ships directly. The rules
// subpackage defines the codes for its own rule catalog.
const (
	// CodeRequired marks a failed required check.
	CodeRequired ErrorCode = "required"
	// CodeMinLength marks a value shorter than the minimum length.
	CodeMinLength ErrorCode = "min_length"
)

// ValidationError is the error type produced by the form's validators. It
// carries a machine-readable Code and a human-readable Message, and satisfies
// the error interface so it flows through the engine's errors.Join aggregation
// and Errors() maps unchanged.
//
// ValidationError is returned by pointer so errors.As can match it; the engine
// stores []error and joins them, so a wrapped or joined ValidationError is still
// recoverable via errors.As.
type ValidationError struct {
	// Code identifies the failure kind. May be empty for ad-hoc validators.
	Code ErrorCode
	// Message is the human-readable description shown to users.
	Message string
	// Path is the dotted form path of the failing control, filled in by the
	// binding layer when known. Validators themselves leave it empty.
	Path string
}

// NewValidationError builds a coded ValidationError with a formatted message.
func NewValidationError(code ErrorCode, format string, args ...any) *ValidationError {
	return &ValidationError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		if e.Code != "" {
			return string(e.Code)
		}
		return "validation failed"
	}
	return e.Message
}

// Unwrap returns nil: a ValidationError is a leaf. It exists so future wrapping
// (for example a cause) stays compatible with errors.Is/As.
func (e *ValidationError) Unwrap() error { return nil }

// WithPath returns a copy of the error carrying path. It is used by the binding
// layer to stamp the control's path onto errors produced without one.
func (e *ValidationError) WithPath(path string) *ValidationError {
	if e == nil {
		return nil
	}
	clone := *e
	clone.Path = path
	return &clone
}

// CodeOf reports the ErrorCode carried by err, if any. It unwraps through
// errors.Join and fmt.Errorf-wrapped chains, returning the first ValidationError
// found. ok is false when err carries no code.
func CodeOf(err error) (ErrorCode, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Code, true
	}
	return "", false
}

// CodesOf collects every distinct ErrorCode within err (following errors.Join),
// sorted and de-duplicated. It is useful for summarizing a whole tree's failures.
func CodesOf(err error) []ErrorCode {
	if err == nil {
		return nil
	}
	seen := map[ErrorCode]struct{}{}
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		// errors.Join yields an interface with Unwrap() []error.
		if multi, ok := e.(interface{ Unwrap() []error }); ok {
			for _, child := range multi.Unwrap() {
				walk(child)
			}
			return
		}
		var ve *ValidationError
		if errors.As(e, &ve) && ve.Code != "" {
			seen[ve.Code] = struct{}{}
		}
	}
	walk(err)

	codes := make([]ErrorCode, 0, len(seen))
	for code := range seen {
		codes = append(codes, code)
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i] < codes[j] })
	return codes
}

// WithCode wraps a validator so its failure carries code. The returned validator
// produces a *ValidationError when it fails, preserving the original message; on
// success it returns nil unchanged. This lets any existing ValidatorFunc gain a
// code without changing the validator interface.
func WithCode[T any](code ErrorCode, validator ValidatorFunc[T]) ValidatorFunc[T] {
	return func(value T) error {
		err := validator(value)
		if err == nil {
			return nil
		}
		// Preserve an existing code if the inner validator already set one.
		if existing, ok := CodeOf(err); ok && existing != "" {
			return err
		}
		var ve *ValidationError
		message := err.Error()
		if errors.As(err, &ve) && ve.Message != "" {
			message = ve.Message
		}
		return &ValidationError{Code: code, Message: message}
	}
}
