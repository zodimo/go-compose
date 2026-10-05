package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// SelectComponentOptions configures SelectComponent.
type SelectComponentOptions struct {
	Modifier ui.Modifier

	// Label is the floating field label. Unspecified renders no label.
	Label string
	// Placeholder is shown in the anchor when no option is selected.
	Placeholder string
	// HintText is shown beneath the field when valid. Unspecified renders no
	// supporting line.
	HintText string
}

// SelectComponentOption configures a SelectComponent.
type SelectComponentOption func(o *SelectComponentOptions)

// DefaultSelectComponentOptions returns the default select options.
func DefaultSelectComponentOptions() SelectComponentOptions {
	return SelectComponentOptions{
		Modifier:    modifier.EmptyModifier,
		Label:       sentinel.StringValueUnspecified,
		Placeholder: sentinel.StringValueUnspecified,
		HintText:    sentinel.StringValueUnspecified,
	}
}

// SelectWithModifier appends a modifier to the anchor field.
func SelectWithModifier(m ui.Modifier) SelectComponentOption {
	return func(o *SelectComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// SelectWithLabel sets the floating field label.
func SelectWithLabel(label string) SelectComponentOption {
	return func(o *SelectComponentOptions) {
		o.Label = label
	}
}

// SelectWithPlaceholder sets the text shown when no option is selected.
func SelectWithPlaceholder(placeholder string) SelectComponentOption {
	return func(o *SelectComponentOptions) {
		o.Placeholder = placeholder
	}
}

// SelectWithHintText sets the supporting hint shown when the field is valid.
func SelectWithHintText(hint string) SelectComponentOption {
	return func(o *SelectComponentOptions) {
		o.HintText = hint
	}
}
