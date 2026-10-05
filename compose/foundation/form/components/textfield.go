// TextFieldComponent renders a bound text field driven by the form engine.
//
// Error display follows the Material 3 convention: WithError is gated on the
// control being touched, so an invalid-but-untouched field renders without error
// styling. All validation error messages (joined via errors.Join into a single
// supporting-text string) are passed as supporting text, falling back to the hint
// text when the field is valid or untouched.
package components

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// TextFieldComponent renders an Outlined text field bound to a form control.
// Value, validation, and touched state are all driven through the control.
//
// The field is disabled when the bound control is not enabled, so a field inside
// a disabled group (own or inherited) renders read-only rather than merely
// ignoring writes.
func TextFieldComponent(binding *fform.FormFieldBinding[string], options ...TextFieldComponentOption) api.Composable {
	opts := DefaultTextFieldComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		hint := sentinel.TakeOrElseString(opts.HintText, "")

		message, hasError := fieldError(binding)
		support := supportingText(hint, message, hasError)

		fieldOpts := []textfield.TextFieldOption{
			textfield.WithLabel(label),
			textfield.WithSupportingText(support),
			textfield.WithError(hasError),
			textfield.WithEnabled(binding.IsEnabled()),
			textfield.WithModifier(opts.Modifier),
		}

		field := textfield.Outlined(
			binding.Value(),
			func(s string) {
				binding.SetValue(s)
			},
			fieldOpts...,
		)

		if opts.Inline {
			return row.Row(c.Sequence(field))(c)
		}
		return column.Column(c.Sequence(field))(c)
	}
}
