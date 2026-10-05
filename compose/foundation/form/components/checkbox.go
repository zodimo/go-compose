package components

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/checkbox"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// CheckboxComponent renders a Material 3 checkbox bound to a bool control, with
// an optional inline label and touched-gated error text beneath it.
//
// Typical usage inside a form scope:
//
//	fform.ControlFieldOf[bool](fs, "sameAsShipping", func(b *fform.FormFieldBinding[bool]) api.Composable {
//	    return components.CheckboxComponent(b, components.CheckboxWithLabel("Same as shipping"))
//	})
func CheckboxComponent(binding *fform.FormFieldBinding[bool], options ...CheckboxComponentOption) api.Composable {
	opts := DefaultCheckboxComponentOptions()
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

		rowContent := []api.Composable{
			checkbox.Checkbox(
				binding.Value(),
				func(checked bool) {
					// A control that is effectively disabled ignores writes, but
					// guarding the callback keeps disabled checkboxes inert even
					// if the underlying widget still reports a change.
					if !binding.IsEnabled() {
						return
					}
					binding.SetValue(checked)
				},
				checkbox.WithModifier(opts.Modifier),
			),
		}
		if label != "" {
			rowContent = append(rowContent, text.BodyMedium(label))
		}

		body := row.Row(
			c.Sequence(rowContent...),
			row.WithAlignment(row.Middle),
		)

		if support == "" {
			return body(c)
		}
		return column.Column(
			c.Sequence(
				body,
				text.BodyMedium(support),
				spacer.Height(opts.ErrorSpacing),
			),
			column.WithModifier(padding.Start(opts.ErrorIndent)),
		)(c)
	}
}
