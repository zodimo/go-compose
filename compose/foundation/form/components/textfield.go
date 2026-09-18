// TextFieldComponent renders a bound text field driven by the form engine.
//
// Error display follows the Material 3 convention: WithError is gated on the
// control being touched, so an invalid-but-untouched field renders without error
// styling. All validation error messages (joined via errors.Join into a single
// supporting-text string) are passed as supporting text, falling back to the hint
// text when the field is valid or untouched.
package components

import (
	"sort"
	"strings"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// TextFieldComponent renders an Outlined text field bound to a form control.
// Value, validation, and touched state are all driven through the control.
func TextFieldComponent(binding *fform.FormFieldBinding[string], options ...TextFieldComponentOption) api.Composable {
	opts := DefaultTextFieldComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		supportingText := sentinel.TakeOrElseString(opts.HintText, "")

		hasError := binding.HasErrors() && binding.IsTouched()
		if hasError {
			if msg := errorMessage(binding.Errors()); msg != "" {
				supportingText = msg
			}
		}

		return c.IfLazy(
			opts.Inline,
			func() api.Composable {
				return column.Column(
					c.Sequence(
						row.Row(
							c.Sequence(
								textfield.Outlined(
									binding.Value(),
									func(s string) {
										binding.SetValue(s)
									},
									textfield.WithLabel(label),
									textfield.WithSupportingText(supportingText),
									textfield.WithError(hasError),
									textfield.WithModifier(opts.Modifier),
								),
							),
						),
					),
				)
			},
			func() api.Composable {
				return column.Column(
					c.Sequence(
						textfield.Outlined(
							binding.Value(),
							func(s string) {
								binding.SetValue(s)
							},
							textfield.WithLabel(label),
							textfield.WithSupportingText(supportingText),
							textfield.WithError(hasError),
							textfield.WithModifier(opts.Modifier),
						),
					),
				)
			},
		)(c)
	}
}

// errorMessage returns every error message in errs, sorted for determinism and
// joined with newlines. A control exposes a single Errors() entry keyed by its
// own path whose value is the errors.Join aggregation of all validator failures;
// a group exposes one entry per descendant path. Nothing is dropped.
func errorMessage(errs map[string]string) string {
	if len(errs) == 0 {
		return ""
	}
	msgs := make([]string, 0, len(errs))
	for _, msg := range errs {
		msgs = append(msgs, msg)
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n")
}
