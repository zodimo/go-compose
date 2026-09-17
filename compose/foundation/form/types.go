package fform

import "github.com/zodimo/go-compose/compose/ui/layout"

// Validatable is an interface for specifying if a widget is validatable.
type Validatable interface {
	Validate() error

	// SetOnValidationChanged is used to set the callback that will be triggered when the validation state changes.
	// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
	SetOnValidationChanged(func(error))
}

// StringValidator is a function signature for validating string inputs.
type StringValidator func(string) error

type Orientation layout.Axis

const (
	Horizontal Orientation = iota
	Vertical
)
