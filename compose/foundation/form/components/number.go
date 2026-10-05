package components

import (
	"strconv"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
	"github.com/zodimo/go-compose/state"
)

// NumberFieldComponent renders an Outlined text field bound to an int control.
//
// Because the control's value is an int, a partially-typed or unparsable entry
// cannot be represented in the control. The component therefore keeps the raw
// text in a remembered string (so the user can type "-", "1.", etc. freely) and
// writes a parsed int to the control only when the text parses. Unparsable text
// is reported as a supporting error, gated on the same touched rule as other
// fields; the control's own validators still apply to the last parsed value.
//
// A written value marks the control touched, so the error appears as the user
// leaves the field in an invalid numeric state.
//
// Limitation: an int control cannot represent "empty", so clearing this field
// falls back to displaying the control's current value (0 by default) rather
// than staying blank. Use OptionalNumberFieldComponent when the field must be
// left blank.
func NumberFieldComponent(binding *fform.FormFieldBinding[int], options ...NumberFieldComponentOption) api.Composable {
	opts := DefaultNumberFieldComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		hint := sentinel.TakeOrElseString(opts.HintText, "")

		textKey := opts.TextStateKey
		if textKey == "" {
			textKey = "numberFieldText"
		}
		rawText := state.MustRemember(c, textKey, func() string {
			return strconv.Itoa(binding.Value())
		})

		message, hasError := fieldError(binding)

		text := rawText.Get()
		if text == "" {
			text = strconv.Itoa(binding.Value())
		}
		if _, err := strconv.Atoi(text); err != nil {
			// Unparsable in-progress text takes over the supporting slot so the
			// user sees why nothing has been stored, but only after touch.
			if binding.IsTouched() {
				message = "not a valid number"
				hasError = true
			}
		}

		support := supportingText(hint, message, hasError)

		field := textfield.Outlined(
			text,
			func(s string) {
				rawText.Set(s)
				if n, err := strconv.Atoi(s); err == nil {
					binding.SetValue(n)
				} else {
					// Mark touched so invalid numeric text surfaces its error.
					binding.Control().MarkAsTouched()
				}
			},
			textfield.WithLabel(label),
			textfield.WithSupportingText(support),
			textfield.WithError(hasError),
			textfield.WithEnabled(binding.IsEnabled()),
			textfield.WithSingleLine(true),
			textfield.WithModifier(opts.Modifier),
		)

		return column.Column(c.Sequence(field))(c)
	}
}
