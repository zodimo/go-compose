package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// SwitchComponentOptions configures SwitchComponent.
type SwitchComponentOptions struct {
	Modifier ui.Modifier

	// Label is rendered inline beside the switch. Unspecified renders no label.
	Label string
	// HintText is shown beneath the switch when valid. Unspecified renders no
	// supporting line.
	HintText string

	// ErrorIndent indents the error supporting text.
	ErrorIndent int
	// ErrorSpacing is the vertical gap before the error supporting text.
	ErrorSpacing int
}

// SwitchComponentOption configures a SwitchComponent.
type SwitchComponentOption func(o *SwitchComponentOptions)

// DefaultSwitchComponentOptions returns the default switch options.
func DefaultSwitchComponentOptions() SwitchComponentOptions {
	return SwitchComponentOptions{
		Modifier:     modifier.EmptyModifier,
		Label:        sentinel.StringValueUnspecified,
		HintText:     sentinel.StringValueUnspecified,
		ErrorIndent:  0,
		ErrorSpacing: 4,
	}
}

// SwitchWithModifier appends a modifier to the switch.
func SwitchWithModifier(m ui.Modifier) SwitchComponentOption {
	return func(o *SwitchComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// SwitchWithLabel sets the inline label.
func SwitchWithLabel(label string) SwitchComponentOption {
	return func(o *SwitchComponentOptions) {
		o.Label = label
	}
}

// SwitchWithHintText sets the supporting hint shown when the field is valid.
func SwitchWithHintText(hint string) SwitchComponentOption {
	return func(o *SwitchComponentOptions) {
		o.HintText = hint
	}
}

// SwitchWithErrorIndent sets the indent of the error supporting text.
func SwitchWithErrorIndent(indent int) SwitchComponentOption {
	return func(o *SwitchComponentOptions) {
		o.ErrorIndent = indent
	}
}
